// Package server runs games for networked clients (ADR 006): one room
// per game, whatever the transport.
package server

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"

	"github.com/FourSoulsRepo/four_souls/internal/protocol"
	"github.com/FourSoulsRepo/four_souls/internal/version"
	engine "github.com/FourSoulsRepo/rules_engine"
)

// Conn is one client connection. Send queues a message and reports
// false when the client is gone or too slow; Close ends the connection.
// Transports implement it; the room never sees sockets (ADR 006).
type Conn interface {
	Send(msg []byte) bool
	Close()
}

// Client is a connection's handle on the room. The transport passes
// every message it reads to Receive and calls Leave when the connection
// ends.
type Client struct {
	room *Room
	conn Conn
	// Set by the room goroutine only.
	hello bool
	seat  int // -1: not seated
	role  protocol.Role
}

// Receive hands a message from the client to the room.
func (c *Client) Receive(msg []byte) { c.room.send(request{client: c, msg: msg}) }

// Leave tells the room the connection is gone.
func (c *Client) Leave() { c.room.send(request{client: c, leave: true}) }

// seat is one place at the table. Its token lets a player come back.
type seat struct {
	name   string
	token  string
	client *Client
}

type request struct {
	client *Client
	msg    []byte
	join   bool
	leave  bool
}

// Room runs one game. Only its goroutine (Run) touches the game.
type Room struct {
	in      chan request
	done    chan struct{}
	game    *engine.Game
	step    int
	seats   []seat
	clients []*Client
}

// NewRoom sets up a game; its seats are free until players join.
func NewRoom(setup engine.Setup) (*Room, error) {
	g, _, err := engine.NewGame(setup)
	if err != nil {
		return nil, fmt.Errorf("server: %w", err)
	}
	return &Room{
		in:    make(chan request, 64),
		done:  make(chan struct{}),
		game:  g,
		seats: make([]seat, setup.Players),
	}, nil
}

// Connect adds a connection; its first message must be hello.
func (r *Room) Connect(conn Conn) *Client {
	c := &Client{room: r, conn: conn, seat: -1}
	r.send(request{client: c, join: true})
	return c
}

// send passes a request to the room unless it has stopped.
func (r *Room) send(req request) {
	select {
	case r.in <- req:
	case <-r.done:
	}
}

// Run handles requests until ctx ends, then closes every connection.
func (r *Room) Run(ctx context.Context) {
	defer func() {
		close(r.done)
		for _, c := range r.clients {
			c.conn.Close()
		}
	}()
	for {
		select {
		case <-ctx.Done():
			return
		case req := <-r.in:
			r.handle(req)
		}
	}
}

func (r *Room) handle(req request) {
	c := req.client
	switch {
	case req.join:
		r.clients = append(r.clients, c)
		return
	case req.leave:
		r.drop(c)
		return
	}
	env, err := protocol.Decode(req.msg)
	if err != nil {
		r.fail(c, 0, protocol.ErrBadMessage, err.Error(), "")
		return
	}
	if !c.hello && env.Type != protocol.TypeHello {
		r.fail(c, env.ID, protocol.ErrNoHello, "the first message must be hello", "")
		return
	}
	switch env.Type {
	case protocol.TypeHello:
		r.hello(c, env)
	case protocol.TypeIntent:
		r.intent(c, env)
	case protocol.TypeResync:
		r.update(c, nil)
	default:
		r.fail(c, env.ID, protocol.ErrBadMessage, "unknown message type "+env.Type, "")
	}
}

// hello seats a player: back in their seat with a token, or in the
// first free one.
func (r *Room) hello(c *Client, env protocol.Envelope) {
	var h protocol.Hello
	if err := env.Unpack(&h); err != nil {
		r.fail(c, env.ID, protocol.ErrBadMessage, err.Error(), "")
		return
	}
	if code := protocol.CheckVersion(h.Protocol); code != "" {
		r.fail(c, env.ID, code, fmt.Sprintf("protocol %d, server speaks %d", h.Protocol, protocol.Version), "")
		r.drop(c)
		c.conn.Close()
		return
	}
	name, ok := protocol.CleanName(h.Name)
	if !ok {
		r.fail(c, env.ID, protocol.ErrBadName, fmt.Sprintf("a nickname has 1 to %d characters", protocol.MaxName), "")
		return
	}
	i := r.seatFor(h.Token)
	if i < 0 {
		r.fail(c, env.ID, protocol.ErrRoomFull, "every seat is taken", "")
		return
	}
	s := &r.seats[i]
	if s.client != nil && s.client != c {
		s.client.seat = -1 // an older connection of the same player
		s.client.conn.Close()
	}
	if s.token == "" {
		s.token = newToken()
		s.name = protocol.UniqueName(name, r.names())
	}
	s.client = c
	c.hello, c.seat, c.role = true, i, protocol.RolePlayer
	v := version.Get()
	r.reply(c, protocol.TypeWelcome, env.ID, protocol.Welcome{Protocol: protocol.Version, App: v.App, Engine: v.Engine, Token: s.token, Seat: i})
	r.broadcast(nil) // everyone sees the seat taken
}

// seatFor finds the seat a token belongs to, or the first free seat.
func (r *Room) seatFor(token string) int {
	if token != "" {
		for i, s := range r.seats {
			if s.token == token {
				return i
			}
		}
	}
	for i, s := range r.seats {
		if s.token == "" {
			return i
		}
	}
	return -1
}

func (r *Room) names() []string {
	var out []string
	for _, s := range r.seats {
		if s.name != "" {
			out = append(out, s.name)
		}
	}
	return out
}

// intent applies a seated player's intent and updates everyone.
func (r *Room) intent(c *Client, env protocol.Envelope) {
	if c.seat < 0 {
		r.fail(c, env.ID, protocol.ErrNotSeated, "only a seated player acts", "")
		return
	}
	var m protocol.Intent
	if err := env.Unpack(&m); err != nil {
		r.fail(c, env.ID, protocol.ErrBadMessage, err.Error(), "")
		return
	}
	in := m.Intent
	in.Player = engine.PlayerID(c.seat) // a client acts only for its own seat
	events, err := r.game.Apply(in)
	var re *engine.RuleError
	switch {
	case errors.As(err, &re):
		r.fail(c, env.ID, protocol.ErrRefused, re.Reason, re.Rule)
		return
	case err != nil:
		r.fail(c, env.ID, protocol.ErrBadMessage, err.Error(), "")
		return
	}
	r.step++
	r.broadcast(events)
}

// broadcast sends every client its update.
func (r *Room) broadcast(events []engine.Event) {
	for _, c := range r.clients {
		if c.hello {
			r.update(c, events)
		}
	}
}

// update sends one client its view, events and allowed intents.
func (r *Room) update(c *Client, events []engine.Event) {
	v := engine.Viewer{Kind: engine.ViewSpectator}
	var allowed []engine.Intent
	if c.seat >= 0 {
		v = engine.Viewer{Kind: engine.ViewPlayer, Player: engine.PlayerID(c.seat)}
		allowed = r.game.Allowed(engine.PlayerID(c.seat))
	}
	seats := make([]protocol.Seat, len(r.seats))
	for i, s := range r.seats {
		seats[i] = protocol.Seat{Seat: i, Name: s.name, Connected: s.client != nil}
	}
	r.reply(c, protocol.TypeUpdate, 0, protocol.Update{
		Step: r.step, Events: engine.FilterEvents(events, v), View: r.game.View(v), Allowed: allowed, Seats: seats,
	})
}

// drop forgets a connection; its seat stays reserved for its token.
func (r *Room) drop(c *Client) {
	for i, x := range r.clients {
		if x == c {
			r.clients = append(r.clients[:i:i], r.clients[i+1:]...)
			break
		}
	}
	if c.seat >= 0 && r.seats[c.seat].client == c {
		r.seats[c.seat].client = nil
		r.broadcast(nil)
	}
	c.seat = -1
}

func (r *Room) fail(c *Client, id int, code, message, rule string) {
	r.reply(c, protocol.TypeError, id, protocol.Error{Code: code, Message: message, Rule: rule})
}

// reply sends a message; a client that cannot keep up is disconnected.
func (r *Room) reply(c *Client, typ string, id int, data any) {
	msg, err := protocol.Encode(typ, id, data)
	if err != nil {
		panic(err) // our own types always encode
	}
	if !c.conn.Send(msg) {
		c.conn.Close()
	}
}

// newToken returns 128 random bits as hex (ADR 006).
func newToken() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		panic(err) // crypto/rand does not fail on supported systems
	}
	return hex.EncodeToString(b)
}
