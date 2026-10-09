// Package server runs games for networked clients (ADR 006): a lobby
// hub, one room per running game, whatever the transport.
package server

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"slices"
	"strings"

	"github.com/FourSoulsRepo/four_souls/internal/protocol"
	engine "github.com/FourSoulsRepo/rules_engine"
)

// roomSeat is one place at the table. Its token lets a player come back.
type roomSeat struct {
	name   string
	token  string
	client *Client
}

type request struct {
	client *Client
	msg    []byte
	leave  bool
	attach string // the token of a client coming back
}

// Room runs one game. Only its goroutine (Run) touches the game.
type Room struct {
	in      chan request
	done    chan struct{}
	game    *engine.Game
	step    int
	seats   []roomSeat
	clients []*Client
	tokens  []string // read by the hub; fixed when the room is made
}

// NewRoom starts a game for seated players.
func NewRoom(setup engine.Setup, seats []roomSeat) (*Room, error) {
	g, _, err := engine.NewGame(setup)
	if err != nil {
		return nil, fmt.Errorf("server: %w", err)
	}
	r := &Room{in: make(chan request, 64), done: make(chan struct{}), game: g, seats: seats}
	for i, s := range seats {
		r.tokens = append(r.tokens, s.token)
		if s.client != nil {
			s.client.seat = i
			r.clients = append(r.clients, s.client)
		}
	}
	return r, nil
}

// hasToken reports whether a seat belongs to the token. Safe from any
// goroutine: tokens never change.
func (r *Room) hasToken(token string) bool { return slices.Contains(r.tokens, token) }

// send passes a request to the room unless it has stopped.
func (r *Room) send(req request) {
	select {
	case r.in <- req:
	case <-r.done:
	}
}

// Run sends everyone the first view, then handles requests until ctx
// ends; then it closes every connection.
func (r *Room) Run(ctx context.Context) {
	defer func() {
		close(r.done)
		for _, c := range r.clients {
			c.conn.Close()
		}
	}()
	r.broadcast(nil)
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
	case req.leave:
		r.drop(c)
		return
	case req.attach != "":
		r.attach(c, req.attach)
		return
	}
	env, err := protocol.Decode(req.msg)
	if err != nil {
		c.fail(0, protocol.ErrBadMessage, err.Error())
		return
	}
	switch env.Type {
	case protocol.TypeIntent:
		r.intent(c, env)
	case protocol.TypeResync:
		r.update(c, nil)
	default:
		c.fail(env.ID, protocol.ErrBadMessage, "not during a game: "+env.Type)
	}
}

// attach gives a returning player their seat back (N-08).
func (r *Room) attach(c *Client, token string) {
	i := slices.Index(r.tokens, token)
	if i < 0 {
		c.fail(0, protocol.ErrNotSeated, "no seat for this token")
		return
	}
	if old := r.seats[i].client; old != nil && old != c {
		old.seat = -1
		old.conn.Close()
	}
	r.seats[i].client = c
	c.seat = i
	r.clients = append(r.clients, c)
	r.broadcast(nil)
}

// intent applies a seated player's intent and updates everyone.
func (r *Room) intent(c *Client, env protocol.Envelope) {
	if c.seat < 0 {
		c.fail(env.ID, protocol.ErrNotSeated, "only a seated player acts")
		return
	}
	var m protocol.Intent
	if err := env.Unpack(&m); err != nil {
		c.fail(env.ID, protocol.ErrBadMessage, err.Error())
		return
	}
	in := m.Intent
	in.Player = engine.PlayerID(c.seat) // a client acts only for its own seat
	events, err := r.game.Apply(in)
	var re *engine.RuleError
	switch {
	case errors.As(err, &re):
		c.send(protocol.TypeError, env.ID, protocol.Error{Code: protocol.ErrRefused, Message: re.Reason, Rule: re.Rule})
		return
	case err != nil:
		c.fail(env.ID, protocol.ErrBadMessage, err.Error())
		return
	}
	r.step++
	r.broadcast(events)
}

// broadcast sends every client its update.
func (r *Room) broadcast(events []engine.Event) {
	for _, c := range r.clients {
		r.update(c, events)
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
	c.send(protocol.TypeUpdate, 0, protocol.Update{
		Step: r.step, Events: engine.FilterEvents(events, v), View: r.game.View(v), Allowed: allowed, Seats: seats,
	})
}

// drop forgets a connection; its seat stays for its token.
func (r *Room) drop(c *Client) {
	r.clients = slices.DeleteFunc(r.clients, func(x *Client) bool { return x == c })
	if c.seat >= 0 && r.seats[c.seat].client == c {
		r.seats[c.seat].client = nil
		r.broadcast(nil)
	}
	c.seat = -1
}

// validToken: 32 lowercase hex digits, as newToken makes them.
func validToken(s string) bool {
	if len(s) != 32 {
		return false
	}
	_, err := hex.DecodeString(s)
	return err == nil && strings.ToLower(s) == s
}

// newToken returns 128 random bits as hex (ADR 006).
func newToken() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		panic(err) // crypto/rand does not fail on supported systems
	}
	return hex.EncodeToString(b)
}
