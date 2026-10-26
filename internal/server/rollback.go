package server

import (
	"fmt"
	"log"
	"slices"

	"github.com/FourSoulsRepo/four_souls/internal/protocol"
	engine "github.com/FourSoulsRepo/rules_engine"
)

// A game with a rollback mode keeps one base state and the live steps
// since it; a rollback rebuilds an old state by replay (ADR 020, items
// 7 and 8). 17.3 does Godmode; Soft comes in 17.4, and keeping the base
// over a server restart in 17.5.

// liveStep is one applied step since the base. Its events are not kept:
// a replay needs only the intent and the checksum.
type liveStep struct {
	n   int // the step number; numbers never repeat
	rec engine.StepRecord
}

// setBase makes the game as it is now the base: the start, a loaded
// save, or a game kept over a restart. Nothing before it can be rolled
// back. A game without a rollback mode keeps no base.
func (r *Room) setBase() {
	r.base, r.baseStep, r.live = nil, r.step, nil
	if r.opts.Rollback == "" || r.game == nil {
		return
	}
	state, err := r.game.Save()
	if err != nil {
		log.Printf("server: rollback base of %s: %v; no rollback", r.recInfo.game, err)
		return
	}
	r.base = state
}

// keepLive adds an applied step to the live steps. Its number is r.steps:
// a rollback takes a number too, so the server's step numbers and the
// record's stay the same.
func (r *Room) keepLive(st engine.StepRecord) {
	if r.base != nil {
		r.live = append(r.live, liveStep{n: r.steps, rec: engine.StepRecord{Intent: st.Intent, Checksum: st.Checksum}})
	}
}

// rollback undoes a live step and every step after it, at a player's
// request (ADR 020, item 6). Godmode: any seated player, any live step,
// no vote; only the game end stops it.
func (r *Room) rollback(c *Client, env protocol.Envelope) {
	var m protocol.RollbackTo
	if err := env.Unpack(&m); err != nil {
		c.fail(env.ID, protocol.ErrBadMessage, err.Error())
		return
	}
	switch {
	case c.seat < 0:
		c.fail(env.ID, protocol.ErrNotSeated, "only a seated player rolls back")
		return
	case r.game == nil:
		c.fail(env.ID, protocol.ErrBadMessage, "the game has not started: bans and picks first")
		return
	case r.base == nil || r.opts.Rollback != protocol.RollbackGod:
		c.fail(env.ID, protocol.ErrRollback, "this game has no rollback (ADR 020)")
		return
	case r.game.Over:
		c.fail(env.ID, protocol.ErrRollback, "the game is over")
		return
	case r.paused():
		c.fail(env.ID, protocol.ErrPaused, "the game waits for a player who left")
		return
	}
	i := slices.IndexFunc(r.live, func(s liveStep) bool { return s.n == m.Step })
	if i < 0 {
		c.fail(env.ID, protocol.ErrRollback, fmt.Sprintf("step %d cannot be rolled back: it is undone, or from before the base", m.Step))
		return
	}
	steps := make([]engine.StepRecord, i)
	for j, s := range r.live[:i] {
		steps[j] = s.rec
	}
	g, err := engine.ReplayFrom(r.base, steps, r.setup.Sets...)
	if err != nil {
		// A replay that differs leaves the game as it is (item 7.4).
		log.Printf("server: rollback of %s to before step %d: %v", r.recInfo.game, m.Step, err)
		c.fail(env.ID, protocol.ErrRollback, "the rollback failed; the game goes on")
		return
	}
	undid := len(r.live) - i
	r.game, r.live = g, r.live[:i]
	// The rollback takes the next number: no step number repeats, and
	// every client sees a new step (item 8.1).
	r.steps++
	r.step++
	if r.rec != nil {
		if err := r.rec.Rollback(m.Step, c.seat, undid); err != nil {
			log.Printf("server: record %s: %v; recording stops", r.path, err)
			r.rec = nil
		}
	}
	// "Wait for my turn" would pass the restored prompt away at once.
	for i := range r.seats {
		r.seats[i].wait = false
	}
	// The restored table goes out alone, with no events: the clients
	// jump to it (item 8.5). Passes the server makes for a player who
	// can only pass come after, as new steps.
	r.rolled = &protocol.Rollback{Seat: c.seat, To: m.Step, Undid: undid}
	r.restartTimer()
	r.broadcast(nil)
	r.rolled = nil
	if _, ok := r.autoIntent(); ok {
		r.after(nil, false)
	}
}
