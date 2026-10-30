package main

import (
	"context"
	"errors"
	"fmt"
	"math/rand/v2"
	"strconv"
	"sync"
	"time"

	"github.com/gorilla/websocket"

	"github.com/FourSoulsRepo/four_souls/internal/protocol"
	"github.com/FourSoulsRepo/four_souls/internal/version"
	engine "github.com/FourSoulsRepo/rules_engine"
)

const (
	lobbyWait   = 30 * time.Second // a lobby answer must come within this
	writeWait   = 10 * time.Second // one write to the server
	pingEvery   = time.Second      // a WebSocket ping per bot, for the ping figures
	maxRefusals = 20               // refused in a row before a bot gives up
	maxSeats    = 4                // seats per game; each bot gets its own seed
	turnLimit   = 500              // steps in one turn before a bot ends it
	maxBackoff  = 10 * time.Second // the longest wait between reconnect attempts
)

// errClosed is returned when the game has closed a bot's socket for good.
var errClosed = errors.New("the game closed the connection")

// bot is one seat of a game: a WebSocket client that plays random legal
// moves (ADR 006).
type bot struct {
	name  string
	cfg   config
	seat  int
	token string // the seat's token; a reconnect says hello with it (N-08)
	rng   *rand.Rand
	g     *gameRun
	col   *collector
	seq   int // message IDs

	mu     sync.Mutex // guards conn and closed: the game closes the socket from another goroutine
	conn   *websocket.Conn
	closed bool

	lastStep    int       // the newest step this seat has seen
	refusals    int       // refused in a row
	pending     bool      // an intent is out, waiting for its update
	pendingStep int       // the step the intent answers
	sentAt      time.Time // when the intent went out
	turnNo      int       // the turn this seat last saw
	turnStep    int       // the step that turn began at
}

// dial opens a connection for one seat and says hello; the welcome gives the
// seat its token. The bot is returned even when it fails, so the caller can
// close its socket.
func dial(ctx context.Context, cfg config, g *gameRun, seat int, col *collector) (*bot, error) {
	seed := cfg.seed + uint64(g.number*maxSeats+seat) //nolint:gosec // seat numbers are small
	b := &bot{
		name: fmt.Sprintf("bot%d", seat+1),
		cfg:  cfg,
		seat: seat,
		rng:  rand.New(rand.NewPCG(seed, 3)), //nolint:gosec // the bots' moves, not security
		g:    g,
		col:  col,
	}
	w, err := b.connect(ctx, "")
	if err != nil {
		return b, err
	}
	b.token = w.Token
	return b, nil
}

// connect opens a socket and says hello with token; an empty token asks for
// a new seat. The socket replaces the bot's old one. The welcome's handshake
// time is recorded.
func (b *bot) connect(ctx context.Context, token string) (protocol.Welcome, error) {
	var w protocol.Welcome
	start := time.Now()
	conn, resp, err := websocket.DefaultDialer.DialContext(ctx, b.cfg.url, nil)
	if resp != nil && resp.Body != nil {
		_ = resp.Body.Close() //nolint:errcheck // the handshake answer is not needed
	}
	if err != nil {
		return w, fmt.Errorf("connect to %s: %w", b.cfg.url, err)
	}
	if err := b.swap(conn); err != nil {
		return w, err
	}
	// The server's pong carries our ping's send time, so the RTT needs no state.
	conn.SetPongHandler(func(payload string) error {
		if sent, err := strconv.ParseInt(payload, 10, 64); err == nil {
			b.col.addPing(time.Since(time.Unix(0, sent)))
		}
		return nil
	})
	go b.pings(ctx, conn)
	hello := protocol.Hello{
		Protocol: protocol.Version,
		App:      version.Get().App,
		Name:     b.name,
		Role:     protocol.RolePlayer,
		Sets:     b.cfg.sets,
		Token:    token,
	}
	if err := b.send(protocol.TypeHello, hello); err != nil {
		return w, err
	}
	if err := b.expect(protocol.TypeWelcome, &w); err != nil {
		return w, err
	}
	b.col.addHandshake(time.Since(start))
	b.col.setServer(w)
	return w, nil
}

// swap makes conn the bot's socket and closes the one it replaces. A closed
// bot takes no socket: the game is over.
func (b *bot) swap(conn *websocket.Conn) error {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.closed {
		_ = conn.Close() //nolint:errcheck // the game is over; nothing to do with it
		return errClosed
	}
	if b.conn != nil {
		_ = b.conn.Close() //nolint:errcheck // the old socket is done with
	}
	b.conn = conn
	return nil
}

// close closes the bot's socket and any socket it gets later. It is safe to
// call from another goroutine.
func (b *bot) close() {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.closed = true
	if b.conn != nil {
		_ = b.conn.Close() //nolint:errcheck // closing only ends the reads; the game is decided
	}
}

// send writes one message to the server.
func (b *bot) send(typ string, data any) error {
	b.seq++
	msg, err := protocol.Encode(typ, b.seq, data)
	if err != nil {
		return err
	}
	if err := b.conn.SetWriteDeadline(time.Now().Add(writeWait)); err != nil {
		return err
	}
	return b.conn.WriteMessage(websocket.TextMessage, msg)
}

// expect reads lobby messages until one of type typ arrives and decodes it
// into v; v may be nil. An error reply fails it.
func (b *bot) expect(typ string, v any) error {
	if err := b.conn.SetReadDeadline(time.Now().Add(lobbyWait)); err != nil {
		return err
	}
	err := b.await(typ, v)
	if derr := b.conn.SetReadDeadline(time.Time{}); err == nil {
		err = derr
	}
	return err
}

func (b *bot) await(typ string, v any) error {
	for {
		_, msg, err := b.conn.ReadMessage()
		if err != nil {
			return fmt.Errorf("%s waiting for %s: %w", b.name, typ, err)
		}
		env, err := protocol.Decode(msg)
		if err != nil {
			return fmt.Errorf("%s waiting for %s: %w", b.name, typ, err)
		}
		switch env.Type {
		case typ:
			if v == nil {
				return nil
			}
			return env.Unpack(v)
		case protocol.TypeError:
			var e protocol.Error
			if err := env.Unpack(&e); err != nil {
				return err
			}
			return fmt.Errorf("%s waiting for %s: server error %s: %s", b.name, typ, e.Code, e.Message)
		}
	}
}

// play follows the game until it is over: it answers each update that
// lists something to do and records what it sees. A lost connection is
// taken back with the seat's token.
func (b *bot) play(ctx context.Context) error {
	for {
		_, msg, err := b.conn.ReadMessage()
		if err != nil {
			if ctx.Err() != nil {
				return b.g.stopReason(ctx)
			}
			if rerr := b.reconnect(ctx); rerr != nil {
				return rerr
			}
			continue
		}
		env, err := protocol.Decode(msg)
		if err != nil {
			b.col.addError("bad_frame")
			continue
		}
		over, err := b.handle(env, len(msg))
		if err != nil || over {
			return err
		}
	}
}

// reconnect takes the seat back on a new socket, as the app does after a
// lost game connection (8.24): at once, then with a longer wait each time,
// until the server seats it again or the game ends. A failed try is counted
// as a failed attempt; a seat taken back is a reconnect.
func (b *bot) reconnect(ctx context.Context) error {
	b.pending = false // an intent sent on the lost socket gets no answer
	for attempt := 0; ; attempt++ {
		select {
		case <-ctx.Done():
			return b.g.stopReason(ctx)
		case <-time.After(retryDelay(attempt)):
		}
		w, err := b.connect(ctx, b.token)
		if err == nil {
			if w.Seat != b.seat {
				return fmt.Errorf("%s came back as seat %d, not %d: the game is gone", b.name, w.Seat, b.seat)
			}
			b.col.addReconnect()
			return nil
		}
		if ctx.Err() != nil {
			return b.g.stopReason(ctx)
		}
		b.col.addFailedDial()
	}
}

// retryDelay is the wait before reconnect attempt n, from 0: at once, then 1,
// 2, 4 and 8 seconds, then every 10, as the app waits (8.24).
func retryDelay(n int) time.Duration {
	if n <= 0 {
		return 0
	}
	return min(time.Second<<min(n-1, 4), maxBackoff)
}

// pings sends a WebSocket ping every pingEvery over conn until ctx ends or
// the socket is closed. The pong handler measures each round trip.
func (b *bot) pings(ctx context.Context, conn *websocket.Conn) {
	tick := time.NewTicker(pingEvery)
	defer tick.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-tick.C:
			payload := strconv.FormatInt(time.Now().UnixNano(), 10)
			if err := conn.WriteControl(websocket.PingMessage, []byte(payload), time.Now().Add(writeWait)); err != nil {
				return
			}
		}
	}
}

// handle reads one message. It reports whether the game is over.
func (b *bot) handle(env protocol.Envelope, size int) (bool, error) {
	switch env.Type {
	case protocol.TypeUpdate:
		return b.update(env, size)
	case protocol.TypeError:
		return false, b.errorReply(env)
	case protocol.TypeNotice:
		return false, b.notice(env)
	case protocol.TypeRestarting:
		return false, fmt.Errorf("%s: the server is restarting", b.name)
	}
	return false, nil
}

// update records one step, answers it when this seat has something to do,
// and reports the end of the game.
func (b *bot) update(env protocol.Envelope, size int) (bool, error) {
	var u protocol.Update
	if err := env.Unpack(&u); err != nil {
		return false, err
	}
	b.col.addUpdate(size)
	if b.pending && u.Step > b.pendingStep {
		b.col.addResponse(time.Since(b.sentAt))
		b.pending = false
	}
	// One update can cover several steps, so a jump is normal; a step that
	// goes back is not.
	switch {
	case u.Step < b.lastStep:
		b.col.addError("out_of_order")
	case u.Step > b.lastStep:
		b.lastStep = u.Step
		b.refusals = 0
	}
	b.g.progress(u.Step, u.View.Turn.Number)
	if u.View.Over {
		return true, b.g.finish(u.Step, winnersOf(u.View.Winners))
	}
	if u.View.Turn.Number != b.turnNo {
		b.turnNo, b.turnStep = u.View.Turn.Number, u.Step
	}
	if b.pending || u.Pause != nil || len(u.Allowed) == 0 {
		return false, nil
	}
	return false, b.act(u)
}

// act sends one allowed intent. A turn that has run for turnLimit steps is
// ended, as a player would end it.
func (b *bot) act(u protocol.Update) error {
	in := pickIntent(b.rng, u.Allowed, u.Step-b.turnStep > turnLimit)
	if in.Kind == engine.IntentDiscard && u.View.Waiting.Count < len(in.Objects) {
		in.Objects = in.Objects[:u.View.Waiting.Count]
	}
	b.pending, b.pendingStep = true, u.Step
	b.sentAt = time.Now()
	return b.send(protocol.TypeIntent, protocol.Intent{Intent: in})
}

// pickIntent prefers doing something over passing two times in three, so
// games move on. With hurry set it ends the turn when it may.
func pickIntent(rng *rand.Rand, allowed []engine.Intent, hurry bool) engine.Intent {
	var acts []engine.Intent
	if hurry {
		for _, in := range allowed {
			if in.Kind == engine.IntentEndTurn {
				return in
			}
		}
	}
	for _, in := range allowed {
		if in.Kind != engine.IntentPass && in.Kind != engine.IntentEndTurn {
			acts = append(acts, in)
		}
	}
	if len(acts) > 0 && rng.IntN(3) > 0 {
		return acts[rng.IntN(len(acts))]
	}
	return allowed[rng.IntN(len(allowed))]
}

// errorReply counts an error from the server. A refused intent leaves the
// seat waiting, so it asks for the state again (resync) and tries once more.
func (b *bot) errorReply(env protocol.Envelope) error {
	var e protocol.Error
	if err := env.Unpack(&e); err != nil {
		return err
	}
	b.col.addError(e.Code)
	if b.pending {
		b.col.addResponse(time.Since(b.sentAt))
		b.pending = false
	}
	b.refusals++
	if b.refusals > maxRefusals {
		return fmt.Errorf("%s: %d errors in a row, last %s: %s", b.name, b.refusals, e.Code, e.Message)
	}
	return b.send(protocol.TypeResync, struct{}{})
}

// notice counts a notice; one that ends the game is a failure.
func (b *bot) notice(env protocol.Envelope) error {
	var n protocol.Notice
	if err := env.Unpack(&n); err != nil {
		return err
	}
	b.col.addNotice(n.Kind)
	switch n.Kind {
	case protocol.NoticeLost, protocol.NoticeIdle:
		return fmt.Errorf("the server ended the game: %s", n.Kind)
	}
	return nil
}
