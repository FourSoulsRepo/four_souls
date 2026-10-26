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
	"sync/atomic"
	"time"

	"github.com/FourSoulsRepo/four_souls/internal/protocol"
	"github.com/FourSoulsRepo/record"
	engine "github.com/FourSoulsRepo/rules_engine"
)

// roomSeat is one place at the table. Its token lets a player come back.
type roomSeat struct {
	name   string
	token  string
	client *Client
	wait   bool // "wait for my turn": pass until the turn before ends (N-16)
	// waitSeq is the last stack ID when the wait began; a newer item
	// that targets the player ends it (15.84).
	waitSeq int
	kicked  bool   // empty for good after a vote (N-08)
	dice    string // the dice look (V-19); kept while the player is away
	pace    pace   // presence messages (15.35)
}

type request struct {
	client  *Client
	msg     []byte
	leave   bool
	attach  string // the token of a client coming back
	dice    string // the returning client's dice look (V-19)
	timeout int    // a ban turn's timer ran out (its turn number)
	respond int    // the response timer ran out (its number)
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

	// second is one second of the ban and response timers; tests shorten it.
	second time.Duration
	votes  []protocol.Vote // while paused (N-08)

	// The match record (6.9). path is fixed when the room is made; over
	// is read by the hub for downloads and hellos (19.32).
	recDir   string
	recInfo  recordInfo
	steps    int // applied steps, for the record of a saved game
	saves    string
	loaded   *saveFile // the save this room continues, if any
	stopped  bool      // saved: the room ends, its players go back to the lobby
	hub      *Hub
	rec      *record.Writer
	path     string
	over     atomic.Bool
	timer    *time.Timer // the response timer (N-07)
	timerGen int
	deadline time.Time

	// pub is the public table presence checks against, for pubStep of
	// pubGame (15.35).
	pub     *publicTable
	pubGame *engine.Game
	pubStep int
	// presenceOn relays presence; the server's operator may turn it off.
	presenceOn bool

	// Kept over a server restart (N-20).
	n         int          // the game's number on the server
	host      string       // the nickname that opened the table
	born      time.Time    // when the game left the lobby
	kept      *restartFile // the restart file this room continues, if any
	setupWait bool         // kept bans and picks start over once all are back
	notice    string       // told to each returning player once (N-20)
	noticed   []bool       // by seat
	grace     time.Time    // no kick vote before it: players come back

	// Idle limits (N-20): the hub's, with its clock.
	emptyLimit time.Duration
	gameLimit  time.Duration
	emptySince time.Time // when the last player left
	now        func() time.Time
	limitTick  time.Duration

	// Rollback (ADR 020, item 7): the base state, its step number, and
	// the live steps since; only when the game's mode allows rollback.
	base     []byte
	baseStep int
	live     []liveStep
	rolled   *protocol.Rollback // the rollback the next broadcast tells
}

// NewRoom prepares a game for seated players; it starts in Run. The
// options must have passed checkOptions.
func NewRoom(setup engine.Setup, opts protocol.Options, seats []roomSeat) (*Room, error) {
	if setup.Players != len(seats) {
		return nil, fmt.Errorf("server: %d seats for %d players", len(seats), setup.Players)
	}
	r := &Room{in: make(chan request, 64), done: make(chan struct{}), setup: setup, opts: opts, seats: seats, second: time.Second, now: time.Now, limitTick: time.Minute}
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

// seatOf is the seat a token belongs to, or -1.
func (r *Room) seatOf(token string) int { return slices.Index(r.tokens, token) }

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
		r.closeRecord(false)
		close(r.done)
		if r.stopped {
			r.backToLobby()
			return
		}
		if r.match != nil && r.match.timer != nil {
			r.match.timer.Stop()
		}
		if r.timer != nil {
			r.timer.Stop()
		}
		for _, c := range r.clients {
			c.conn.Close()
		}
	}()
	switch {
	case r.kept != nil:
		r.resumeKept()
	case r.loaded != nil:
		r.resumeLoaded()
	default:
		r.beginSetup()
	}
	tick := time.NewTicker(r.limitTick)
	defer tick.Stop()
	for !r.stopped {
		select {
		case <-ctx.Done():
			r.keep()
			return
		case req := <-r.in:
			r.handle(req)
		case <-tick.C:
			r.checkLimits()
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
		r.attach(c, req.attach, req.dice)
		return
	case req.timeout > 0:
		r.banTimeout(req.timeout)
		return
	case req.respond > 0:
		r.timeUp(req.respond)
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
	case protocol.TypeWaitTurn:
		r.waitTurn(c, env)
	case protocol.TypeVote:
		r.vote(c, env)
	case protocol.TypeSave:
		r.save(c, env)
	case protocol.TypePresence:
		r.presence(c, env)
	case protocol.TypeRollback:
		r.rollback(c, env)
	default:
		c.fail(env.ID, protocol.ErrBadMessage, "not during a game: "+env.Type)
	}
}

// attach gives a returning player their seat back (N-08).
func (r *Room) attach(c *Client, token, dice string) {
	i := slices.Index(r.tokens, token)
	if i < 0 || r.seats[i].kicked {
		c.fail(0, protocol.ErrNotSeated, "no seat for this token")
		c.conn.Close()
		return
	}
	if old := r.seats[i].client; old != nil && old != c {
		old.seat = -1
		old.conn.Close()
	}
	r.seats[i].client, r.seats[i].dice = c, dice
	c.seat = i
	r.clients = append(r.clients, c)
	if r.kept != nil {
		r.tellEngine(c)
	}
	if r.setupWait && r.everyoneBack() {
		r.setupWait = false
		r.beginSetup()
		return
	}
	if r.match != nil {
		r.broadcastSetup()
		return
	}
	r.resume()
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
	if r.paused() {
		c.fail(env.ID, protocol.ErrPaused, "the game waits for a player who left")
		return
	}
	var m protocol.Intent
	if err := env.Unpack(&m); err != nil {
		c.fail(env.ID, protocol.ErrBadMessage, err.Error())
		return
	}
	in := m.Intent
	in.Player = engine.PlayerID(c.seat) // a client acts only for its own seat
	events, err := r.apply(in)
	var re *engine.RuleError
	switch {
	case errors.As(err, &re):
		c.send(protocol.TypeError, env.ID, protocol.Error{Code: protocol.ErrRefused, Message: re.Reason, Rule: re.Rule})
		return
	case err != nil:
		c.fail(env.ID, protocol.ErrBadMessage, err.Error())
		return
	}
	r.after(events, true)
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
		seats[i] = protocol.Seat{Seat: i, Name: s.name, Connected: s.client != nil, Kicked: s.kicked, Dice: s.dice}
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
	var blocked []engine.Block
	if c.seat >= 0 {
		v = engine.Viewer{Kind: engine.ViewPlayer, Player: engine.PlayerID(c.seat)}
		allowed = r.game.Allowed(engine.PlayerID(c.seat))
		blocked = r.game.Blocked(engine.PlayerID(c.seat))
	}
	u := protocol.Update{
		Step: r.step, Events: engine.FilterEvents(events, v), View: r.game.View(v), Allowed: allowed, Seats: r.seatViews(),
		Aims: r.game.Aims(allowed), Blocked: blocked,
	}
	if !r.deadline.IsZero() {
		u.Deadline = r.deadline.UnixMilli()
	}
	if c.seat >= 0 {
		u.WaitTurn = r.seats[c.seat].wait
	}
	if r.paused() {
		u.Pause = &protocol.Pause{Away: r.away(), Votes: slices.Clone(r.votes)}
		if time.Now().Before(r.grace) {
			u.Pause.Grace = r.grace.UnixMilli()
		}
	}
	if r.base != nil {
		u.RollbackMode, u.RollbackFrom, u.Rollback = r.opts.Rollback, r.baseStep+1, r.rolled
	}
	c.send(protocol.TypeUpdate, 0, u)
}

// drop forgets a connection; its seat stays for its token. During a
// game this pauses it for everyone (N-08).
func (r *Room) drop(c *Client) {
	r.clients = slices.DeleteFunc(r.clients, func(x *Client) bool { return x == c })
	if c.seat >= 0 && r.seats[c.seat].client == c {
		r.seats[c.seat].client = nil
		r.seats[c.seat].wait = false
		r.relayPresence(c, protocol.Presence{Seat: c.seat}) // its arrow goes
		if r.paused() {
			r.restartTimer() // stops it while paused
		}
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

// backToLobby hands the room's players back to the lobby after a save.
func (r *Room) backToLobby() {
	for _, c := range r.clients {
		c.room.Store(nil)
	}
	if r.hub == nil {
		return
	}
	r.hub.send(hubRequest{run: func(h *Hub) {
		h.games = slices.DeleteFunc(h.games, func(g started) bool { return g.room == r })
		h.lobbyChanged()
	}})
}
