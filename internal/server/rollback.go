package server

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"slices"
	"time"

	"github.com/FourSoulsRepo/four_souls/internal/protocol"
	engine "github.com/FourSoulsRepo/rules_engine"
)

// A game with a rollback mode keeps one base state and the live steps
// since it; a rollback rebuilds an old state by replay (ADR 020, items
// 7 and 8). Godmode rolls back at once; Soft stops at barriers and asks
// the other players (items 3 to 5). A server restart keeps the base
// and the live steps if they replay to the kept state (item 10).

// liveStep is one applied step since the base. Its events are not kept:
// a replay needs only the intent and the checksum.
type liveStep struct {
	n   int // the step number; numbers never repeat
	rec engine.StepRecord
	// barrier: the step had a barrier event; Soft stops here (item 4).
	barrier bool
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
		r.live = append(r.live, liveStep{
			n: r.steps, rec: engine.StepRecord{Intent: st.Intent, Checksum: st.Checksum},
			barrier: slices.ContainsFunc(st.Events, engine.Barrier),
		})
	}
}

// keepHistory puts the base and the live steps into a restart file
// (item 10.1). A game without a base keeps none: Disabled has nothing.
func (r *Room) keepHistory(f *restartFile) error {
	if r.base == nil {
		return nil
	}
	f.Base, f.BaseStep, f.Live = r.base, r.baseStep, make([]keptStep, len(r.live))
	for i, s := range r.live {
		in, err := json.Marshal(s.rec.Intent)
		if err != nil {
			return fmt.Errorf("live step %d: %w", s.n, err)
		}
		f.Live[i] = keptStep{N: s.n, Intent: in, Checksum: s.rec.Checksum, Barrier: s.barrier}
	}
	return nil
}

// restoreHistory takes back the rollback history of a game kept over a
// restart (items 10.2 to 10.6). If it does not check out, the kept
// state is the new base: the restart is a barrier, and every update
// says where it cut the history (rollback_cut).
func (r *Room) restoreHistory(f *restartFile) {
	r.setBase()
	if r.base == nil {
		return // no rollback mode
	}
	if len(f.Base) == 0 {
		r.cut = r.step // a file without the history (before 17.5)
		return
	}
	live, err := r.replayKept(f)
	if err != nil {
		log.Printf("server: rollback history of %s ends at the restart: %v", f.Game, err)
		r.cut = r.step
		return
	}
	r.base, r.baseStep, r.live = f.Base, f.BaseStep, live
}

// replayKept resumes the kept base on this engine and applies the kept
// live steps; the result must save to the kept state's bytes. The old
// checksums are not compared: a new engine may sum otherwise, and the
// final bytes decide. The live steps take this engine's checksums, so
// a later rollback checks against them. A step stays a barrier if it
// was one, or is one now.
func (r *Room) replayKept(f *restartFile) ([]liveStep, error) {
	g, err := engine.Resume(f.Base, f.Meaning, r.setup.Sets...)
	if err != nil {
		return nil, fmt.Errorf("the base: %w", err)
	}
	live := make([]liveStep, len(f.Live))
	for i, k := range f.Live {
		var in engine.Intent
		if jerr := json.Unmarshal(k.Intent, &in); jerr != nil {
			return nil, fmt.Errorf("live step %d: %w", k.N, jerr)
		}
		st, serr := g.Step(in)
		if serr != nil {
			return nil, fmt.Errorf("live step %d: %w", k.N, serr)
		}
		live[i] = liveStep{
			n: k.N, rec: engine.StepRecord{Intent: in, Checksum: st.Checksum},
			barrier: k.Barrier || slices.ContainsFunc(st.Events, engine.Barrier),
		}
	}
	state, err := g.Save()
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(state, f.State) {
		return nil, errors.New("the live steps replay to another state")
	}
	return live, nil
}

// softVote is an open Soft vote (ADR 020, item 5). The game waits
// while it is open.
type softVote struct {
	step     int // the target
	seat     int // who opened it
	author   int // the target's author; never a voter
	undo     int // live steps it would undo
	votes    []protocol.RollbackBallot
	timer    *time.Timer // with the response timer only (item 5.10)
	gen      int
	deadline time.Time
}

// rollbackFrom is the first step a rollback may target: the first
// since the base; in Soft, the first after the last barrier (item 3.1).
func (r *Room) rollbackFrom() int {
	if r.opts.Rollback == protocol.RollbackSoft {
		for i := len(r.live) - 1; i >= 0; i-- {
			if r.live[i].barrier {
				return r.live[i].n + 1
			}
		}
	}
	return r.baseStep + 1
}

// softTarget: a Soft rollback may target the live step, with its
// barrier rule and never a lone pass or throw (item 3.2).
func (r *Room) softTarget(s liveStep) bool {
	k := s.rec.Intent.Kind
	return s.n >= r.rollbackFrom() && k != engine.IntentPass && k != engine.IntentThrow
}

// softSteps lists the Soft targets now with their authors, for the
// clients' buttons.
func (r *Room) softSteps() []protocol.RollbackStep {
	var out []protocol.RollbackStep
	for _, s := range r.live {
		if r.softTarget(s) {
			out = append(out, protocol.RollbackStep{Step: s.n, Seat: int(s.rec.Intent.Player)})
		}
	}
	return out
}

// voters are the connected seats other than the author (item 5.4).
func (r *Room) voters(author int) []int {
	var out []int
	for i, s := range r.seats {
		if i != author && s.client != nil {
			out = append(out, i)
		}
	}
	return out
}

// rollback undoes a live step and every step after it, at a player's
// request (ADR 020, item 6). Godmode: any seated player, any live step,
// no vote; only the game end stops it. Soft: a target after the last
// barrier, and a vote of the others (items 3 and 5).
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
	case r.base == nil || r.opts.Rollback == "":
		c.fail(env.ID, protocol.ErrRollback, "this game has no rollback (ADR 020)")
		return
	case r.game.Over:
		c.fail(env.ID, protocol.ErrRollback, "the game is over")
		return
	case r.paused():
		c.fail(env.ID, protocol.ErrPaused, "the game waits for a player who left")
		return
	case r.rbVote != nil:
		c.fail(env.ID, protocol.ErrRollback, "a rollback vote is open")
		return
	}
	i := slices.IndexFunc(r.live, func(s liveStep) bool { return s.n == m.Step })
	if i < 0 {
		c.fail(env.ID, protocol.ErrRollback, fmt.Sprintf("step %d cannot be rolled back: it is undone, or from before the base", m.Step))
		return
	}
	if r.opts.Rollback == protocol.RollbackGod {
		if err := r.rollBackTo(c.seat, i); err != nil {
			c.fail(env.ID, protocol.ErrRollback, err.Error())
		}
		return
	}
	author := int(r.live[i].rec.Intent.Player)
	switch {
	case !r.softTarget(r.live[i]):
		c.fail(env.ID, protocol.ErrRollback, fmt.Sprintf("step %d cannot be rolled back in Soft: a barrier came after it, or it is a pass or a throw", m.Step))
		return
	case c.seat == author && !m.Ask:
		c.fail(env.ID, protocol.ErrRollback, "in Soft the author never rolls back alone; ask the others (ADR 020, item 5.1)")
		return
	case len(r.voters(author)) == 0:
		c.fail(env.ID, protocol.ErrRollback, "no other player is connected to vote (item 5.7)")
		return
	}
	v := &softVote{step: m.Step, seat: c.seat, author: author, undo: len(r.live) - i}
	if c.seat != author {
		v.votes = []protocol.RollbackBallot{{Seat: c.seat, Yes: true}} // item 5.3
	}
	r.rbVote = v
	r.restartTimer() // the game waits: no response timer
	if r.opts.ResponseTimer > 0 {
		d := time.Duration(r.opts.ResponseTimer) * r.second
		r.voteGen++
		v.gen, v.deadline = r.voteGen, time.Now().Add(d)
		gen := v.gen
		v.timer = time.AfterFunc(d, func() { r.send(request{voteUp: gen}) })
	}
	r.tally()
}

// rollbackVote records a voter's yes or no; a voter may change it.
func (r *Room) rollbackVote(c *Client, env protocol.Envelope) {
	var m protocol.RollbackVoteAnswer
	if err := env.Unpack(&m); err != nil {
		c.fail(env.ID, protocol.ErrBadMessage, err.Error())
		return
	}
	switch {
	case r.rbVote == nil:
		c.fail(env.ID, protocol.ErrRollback, "no rollback vote is open")
		return
	case c.seat < 0 || !slices.Contains(r.voters(r.rbVote.author), c.seat):
		c.fail(env.ID, protocol.ErrRollback, "only the other connected players vote; never the author (item 5.4)")
		return
	}
	v := r.rbVote
	v.votes = slices.DeleteFunc(v.votes, func(b protocol.RollbackBallot) bool { return b.Seat == c.seat })
	v.votes = append(v.votes, protocol.RollbackBallot{Seat: c.seat, Yes: m.Yes})
	r.tally()
}

// rollbackCancel ends the vote at its opener's request (item 5.13).
func (r *Room) rollbackCancel(c *Client, env protocol.Envelope) {
	if r.rbVote == nil || c.seat < 0 || c.seat != r.rbVote.seat {
		c.fail(env.ID, protocol.ErrRollback, "only whoever opened the rollback vote cancels it")
		return
	}
	r.endVote(protocol.VoteCancelled)
}

// voteTimeUp declines a vote whose time ran out (item 5.11).
func (r *Room) voteTimeUp(gen int) {
	if r.rbVote != nil && r.rbVote.gen == gen {
		r.endVote(protocol.VoteTimeout)
	}
}

// need is how many voters must say yes: more than half (item 5.5).
func need(voters int) int { return voters/2 + 1 }

// tally ends the vote once the majority is reached or out of reach
// (item 5.12); else everyone sees the votes so far.
func (r *Room) tally() {
	v := r.rbVote
	voters := r.voters(v.author)
	yes, no := 0, 0
	for _, b := range v.votes {
		switch {
		case !slices.Contains(voters, b.Seat):
		case b.Yes:
			yes++
		default:
			no++
		}
	}
	n := need(len(voters))
	switch {
	case yes >= n:
		r.closeVote()
		i := slices.IndexFunc(r.live, func(s liveStep) bool { return s.n == v.step })
		if err := r.rollBackTo(v.seat, i); err != nil {
			r.tellEnd(v, protocol.VoteFailed)
		}
	case len(voters)-no < n:
		r.endVote(protocol.VoteDeclined)
	default:
		r.broadcast(nil)
	}
}

// closeVote forgets the open vote and stops its timer.
func (r *Room) closeVote() {
	if r.rbVote == nil {
		return
	}
	if r.rbVote.timer != nil {
		r.rbVote.timer.Stop()
	}
	r.rbVote = nil
}

// endVote ends the vote with no rollback; a declined target may be
// asked again (item 5.14). The game goes on: the passes it held back
// come now, and the response timer starts again.
func (r *Room) endVote(why string) {
	v := r.rbVote
	if v == nil {
		return
	}
	r.closeVote()
	r.tellEnd(v, why)
}

// tellEnd sends everyone the table after a vote that ended with no
// rollback, and why.
func (r *Room) tellEnd(v *softVote, why string) {
	r.voteEnd = &protocol.RollbackEnd{Step: v.step, Seat: v.seat, Why: why}
	r.after(nil, false)
	r.voteEnd = nil
}

// voteView is the open vote as the clients see it.
func (r *Room) voteView() *protocol.RollbackVote {
	v := r.rbVote
	if v == nil {
		return nil
	}
	voters := r.voters(v.author)
	out := &protocol.RollbackVote{
		Step: v.step, Seat: v.seat, Author: v.author, Undo: v.undo,
		Voters: voters, Votes: slices.Clone(v.votes), Need: need(len(voters)),
	}
	if !v.deadline.IsZero() {
		out.Deadline = v.deadline.UnixMilli()
	}
	return out
}

// rollBackTo rebuilds the game from before live step i: seat rolled it
// back, alone in Godmode or by a vote in Soft.
func (r *Room) rollBackTo(seat, i int) error {
	target := r.live[i].n
	steps := make([]engine.StepRecord, i)
	for j, s := range r.live[:i] {
		steps[j] = s.rec
	}
	g, err := engine.ReplayFrom(r.base, steps, r.setup.Sets...)
	if err != nil {
		// A replay that differs leaves the game as it is (item 7.4).
		log.Printf("server: rollback of %s to before step %d: %v", r.recInfo.game, target, err)
		return errors.New("the rollback failed; the game goes on")
	}
	undid := len(r.live) - i
	r.game, r.live = g, r.live[:i]
	// The rollback takes the next number: no step number repeats, and
	// every client sees a new step (item 8.1).
	r.steps++
	r.step++
	if r.rec != nil {
		if err := r.rec.Rollback(target, seat, undid); err != nil {
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
	r.rolled = &protocol.Rollback{Seat: seat, To: target, Undid: undid}
	r.restartTimer()
	r.broadcast(nil)
	r.rolled = nil
	if _, ok := r.autoIntent(); ok {
		r.after(nil, false)
	}
	return nil
}
