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
	client  *Client
	msg     []byte
	leave   bool
	attach  string // the token of a client coming back
	timeout int    // a ban turn's timer ran out (its turn number)
}

// Room runs one game. Only its goroutine (Run) touches the game. It
// starts with the ban and pick phase, if the setup has one (6.6).
type Room struct {
	in      chan request
	done    chan struct{}
	setup   engine.Setup
	opts    protocol.Options
	match   *matchSetup  // the ban and pick phase; nil once playing
	game    *engine.Game // nil until the characters are chosen
	step    int
	seats   []roomSeat
	clients []*Client
	tokens  []string // read by the hub; fixed when the room is made
}

// NewRoom prepares a game for seated players; it starts in Run. The
// options must have passed checkOptions.
func NewRoom(setup engine.Setup, opts protocol.Options, seats []roomSeat) (*Room, error) {
	if setup.Players != len(seats) {
		return nil, fmt.Errorf("server: %d seats for %d players", len(seats), setup.Players)
	}
	r := &Room{in: make(chan request, 64), done: make(chan struct{}), setup: setup, opts: opts, seats: seats}
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

// Run starts the setup or the game, then handles requests until ctx
// ends; then it closes every connection.
func (r *Room) Run(ctx context.Context) {
	defer func() {
		close(r.done)
		if r.match != nil && r.match.timer != nil {
			r.match.timer.Stop()
		}
		for _, c := range r.clients {
			c.conn.Close()
		}
	}()
	r.beginSetup()
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
	case req.timeout > 0:
		r.banTimeout(req.timeout)
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
	case protocol.TypeBan:
		r.ban(c, env)
	case protocol.TypePick:
		r.pick(c, env)
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
	if r.match != nil {
		r.broadcastSetup()
	}
}

// intent applies a seated player's intent and updates everyone.
func (r *Room) intent(c *Client, env protocol.Envelope) {
	if c.seat < 0 {
		c.fail(env.ID, protocol.ErrNotSeated, "only a seated player acts")
		return
	}
	if r.game == nil {
		c.fail(env.ID, protocol.ErrBadMessage, "the game has not started: bans and picks first")
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

// broadcast sends every client its update; during the setup phase,
// the setup.
func (r *Room) broadcast(events []engine.Event) {
	for _, c := range r.clients {
		r.update(c, events)
	}
}

// seatViews lists the seats' nicknames and connections.
func (r *Room) seatViews() []protocol.Seat {
	seats := make([]protocol.Seat, len(r.seats))
	for i, s := range r.seats {
		seats[i] = protocol.Seat{Seat: i, Name: s.name, Connected: s.client != nil}
	}
	return seats
}

// update sends one client its view, events and allowed intents.
func (r *Room) update(c *Client, events []engine.Event) {
	if r.game == nil {
		if r.match != nil {
			r.sendSetup(c)
		}
		return
	}
	v := engine.Viewer{Kind: engine.ViewSpectator}
	var allowed []engine.Intent
	if c.seat >= 0 {
		v = engine.Viewer{Kind: engine.ViewPlayer, Player: engine.PlayerID(c.seat)}
		allowed = r.game.Allowed(engine.PlayerID(c.seat))
	}
	c.send(protocol.TypeUpdate, 0, protocol.Update{
		Step: r.step, Events: engine.FilterEvents(events, v), View: r.game.View(v), Allowed: allowed, Seats: r.seatViews(),
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
