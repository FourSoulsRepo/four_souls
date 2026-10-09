package server

import (
	"slices"
	"time"

	"github.com/FourSoulsRepo/four_souls/internal/protocol"
	engine "github.com/FourSoulsRepo/rules_engine"
)

// Response timer limits (N-07).
const (
	minResponseTimer = 60
	maxResponseTimer = 600
	// animationAllowance is added to the response timer so clients can
	// finish showing what happened.
	animationAllowance = 5
)

// maxAutoSteps bounds the automatic passes after one intent.
const maxAutoSteps = 1000

// after runs automatic passes, then sends everyone the result and
// restarts the response timer (6.7).
func (r *Room) after(events []engine.Event) {
	r.step++
	for range maxAutoSteps {
		in, ok := r.autoIntent()
		if !ok {
			break
		}
		more, err := r.game.Apply(in)
		if err != nil {
			break // the engine said no; wait for the player
		}
		events = append(events, more...)
		r.step++
	}
	r.restartTimer()
	r.broadcast(events)
}

// autoIntent is a pass the server makes for the waiting player: one who
// can do nothing but pass (N-05, N-09), or one skipping the stack they
// saw (N-06).
func (r *Room) autoIntent() (engine.Intent, bool) {
	g := r.game
	if g.Over {
		return engine.Intent{}, false
	}
	w := g.Prompt()
	if w.Kind != engine.PromptPriority {
		return engine.Intent{}, false
	}
	pass := engine.Intent{Player: w.Player, Kind: engine.IntentPass}
	if !g.CanAct(w.Player) {
		return pass, true
	}
	s := &r.seats[w.Player]
	if s.skip == nil {
		return engine.Intent{}, false
	}
	if len(g.Stack) == 0 {
		s.skip = nil // "skip all" lasts until the stack is empty
		return engine.Intent{}, false
	}
	for _, it := range g.Stack {
		if !slices.Contains(s.skip, it.ID) {
			s.skip = nil // something new: control returns to the player
			return engine.Intent{}, false
		}
	}
	return pass, true
}

// skipAll turns "skip all" on for the stack the player sees now.
func (r *Room) skipAll(c *Client, env protocol.Envelope) {
	var m protocol.SkipAll
	if err := env.Unpack(&m); err != nil {
		c.fail(env.ID, protocol.ErrBadMessage, err.Error())
		return
	}
	if c.seat < 0 || r.game == nil {
		c.fail(env.ID, protocol.ErrNotSeated, "only a seated player in a game")
		return
	}
	s := &r.seats[c.seat]
	s.skip = nil
	if m.On && len(r.game.Stack) > 0 {
		s.skip = []int{}
		for _, it := range r.game.Stack {
			s.skip = append(s.skip, it.ID)
		}
	}
	r.after(nil)
}

// restartTimer starts the response timer for whoever must answer now.
func (r *Room) restartTimer() {
	if r.timer != nil {
		r.timer.Stop()
		r.timer = nil
	}
	r.timerGen++
	r.deadline = time.Time{}
	if r.opts.ResponseTimer == 0 || r.game.Over {
		return
	}
	d := time.Duration(r.opts.ResponseTimer+animationAllowance) * r.second
	r.deadline = time.Now().Add(d)
	gen := r.timerGen
	r.timer = time.AfterFunc(d, func() { r.send(request{respond: gen}) })
}

// timeUp answers for a player whose response timer ran out (N-07): pass
// or end the turn; a discard takes the first cards; a choice takes the
// first option.
func (r *Room) timeUp(gen int) {
	if gen != r.timerGen || r.game == nil || r.game.Over {
		return
	}
	w := r.game.Prompt()
	allowed := r.game.Allowed(w.Player)
	var in engine.Intent
	switch w.Kind {
	case engine.PromptPriority:
		i := slices.IndexFunc(allowed, func(x engine.Intent) bool { return x.Kind == engine.IntentPass })
		if i < 0 {
			i = slices.IndexFunc(allowed, func(x engine.Intent) bool { return x.Kind == engine.IntentEndTurn })
		}
		if i < 0 {
			return
		}
		in = allowed[i]
	case engine.PromptDiscard:
		in = allowed[0]
		in.Objects = in.Objects[:w.Count]
	case engine.PromptChoose:
		in = allowed[0]
	case engine.PromptNone, engine.PromptGameOver:
		return
	}
	events, err := r.game.Apply(in)
	if err != nil {
		return
	}
	r.after(events)
}
