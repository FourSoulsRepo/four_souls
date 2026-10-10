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
// restarts the response timer (6.7). applied: an intent changed the game.
func (r *Room) after(events []engine.Event, applied bool) {
	if applied {
		r.step++
	}
	for range maxAutoSteps {
		in, ok := r.autoIntent()
		if !ok {
			break
		}
		more, err := r.apply(in)
		if err != nil {
			break // the engine said no; wait for the player
		}
		events = append(events, more...)
		r.step++
	}
	r.restartTimer()
	r.broadcast(events)
	if r.game.Over {
		r.closeRecord(true)
	}
}

// autoIntent is a pass the server makes for the waiting player: one who
// can do nothing but pass (N-05, N-09), or one skipping the stack they
// saw (N-06).
func (r *Room) autoIntent() (engine.Intent, bool) {
	g := r.game
	if g.Over || r.paused() {
		return engine.Intent{}, false
	}
	w := g.Prompt()
	if w.Player >= 0 && int(w.Player) < len(r.seats) && r.seats[w.Player].kicked {
		return r.forcedIntent() // a kicked seat stays empty (N-08)
	}
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
	r.after(nil, false)
}

// restartTimer starts the response timer for whoever must answer now.
func (r *Room) restartTimer() {
	if r.timer != nil {
		r.timer.Stop()
		r.timer = nil
	}
	r.timerGen++
	r.deadline = time.Time{}
	if r.opts.ResponseTimer == 0 || r.game.Over || r.paused() {
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
	if gen != r.timerGen || r.game == nil || r.game.Over || r.paused() {
		return
	}
	in, ok := r.forcedIntent()
	if !ok {
		return
	}
	events, err := r.apply(in)
	if err != nil {
		return
	}
	r.after(events, true)
}

// paused: a running game waits for a seated player who left (N-08).
func (r *Room) paused() bool {
	return r.game != nil && !r.game.Over && len(r.away()) > 0
}

// away lists the seats whose player left and was not kicked.
func (r *Room) away() []int {
	var out []int
	for i, s := range r.seats {
		if s.client == nil && !s.kicked {
			out = append(out, i)
		}
	}
	return out
}

// resume continues a game once nobody is away any more.
func (r *Room) resume() {
	if r.game == nil {
		return
	}
	if !r.paused() {
		r.votes = nil
		r.after(nil, false)
		return
	}
	r.broadcast(nil)
}

// vote records a connected player's vote while the game is paused.
// When more than half of them vote to kick, every missing player is
// kicked; their seats stay empty (N-08).
func (r *Room) vote(c *Client, env protocol.Envelope) {
	var m protocol.Vote
	if err := env.Unpack(&m); err != nil {
		c.fail(env.ID, protocol.ErrBadMessage, err.Error())
		return
	}
	if c.seat < 0 || !r.paused() {
		c.fail(env.ID, protocol.ErrBadMessage, "nothing to vote on")
		return
	}
	r.votes = slices.DeleteFunc(r.votes, func(v protocol.Vote) bool { return v.Seat == c.seat })
	r.votes = append(r.votes, protocol.Vote{Seat: c.seat, Kick: m.Kick})
	connected, kicks := 0, 0
	for _, s := range r.seats {
		if s.client != nil {
			connected++
		}
	}
	for _, v := range r.votes {
		if v.Kick {
			kicks++
		}
	}
	if kicks*2 <= connected {
		r.broadcast(nil)
		return
	}
	for _, i := range r.away() {
		r.seats[i].kicked = true
	}
	r.votes = nil
	r.after(nil, false)
}

// forcedIntent is the minimal answer for the waiting player, when the
// player is out of time or kicked: pass or end the turn; a discard
// takes the first cards; a choice takes the first option.
func (r *Room) forcedIntent() (engine.Intent, bool) {
	w := r.game.Prompt()
	allowed := r.game.Allowed(w.Player)
	if len(allowed) == 0 {
		return engine.Intent{}, false
	}
	switch w.Kind {
	case engine.PromptPriority:
		for _, kind := range []engine.IntentKind{engine.IntentPass, engine.IntentEndTurn} {
			if i := slices.IndexFunc(allowed, func(x engine.Intent) bool { return x.Kind == kind }); i >= 0 {
				return allowed[i], true
			}
		}
		return engine.Intent{}, false
	case engine.PromptDiscard:
		in := allowed[0]
		in.Objects = in.Objects[:w.Count]
		return in, true
	case engine.PromptChoose:
		return allowed[0], true
	case engine.PromptNone, engine.PromptGameOver:
		return engine.Intent{}, false
	}
	return engine.Intent{}, false
}
