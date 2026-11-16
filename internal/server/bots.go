package server

import (
	"context"
	"log"
	"math/rand/v2"
	"slices"
	"time"

	botai "github.com/FourSoulsRepo/four_souls/internal/bot"
	"github.com/FourSoulsRepo/four_souls/internal/protocol"
	"github.com/FourSoulsRepo/four_souls/internal/world"
	engine "github.com/FourSoulsRepo/rules_engine"
)

// The bot driver (ADR 030, Decisions 11 – 14): after every broadcast it
// notes each bot seat the game waits for, and Run starts one thinker
// for it under the room's context once the request is handled. A thinker
// runs bot.Choose on its own goroutine, on worlds the room copied for
// it; its answer comes back as a room request with the step it was
// for. Any newer step, a pause, a vote or the game end drops it.

// The pace of a bot (Decision 13), in the room's seconds.
const (
	// botAnswerPace: a pass or an answer to a prompt, after the prompt came.
	botAnswerPace = 0.7
	// botActionPace: an action in the bot's own turn, plus up to one more.
	botActionPace = 1.0
)

// hardThink is the hard bot's time budget per move (Decision 9.8). The
// flag -bot_think comes with the thinker limit (21.10).
const hardThink = 1500 * time.Millisecond

// thinker is one open bot move.
type thinker struct {
	id     int // unique in the room: an answer names it
	step   int // the room's step it thinks about
	done   <-chan struct{}
	cancel context.CancelFunc
}

// botAnswer is a thinker's move, sent back to the room.
type botAnswer struct {
	seat, id, step int
	intent         engine.Intent
	err            error // no move: the room sends the forced answer
}

// botLevel is the bot package's level for a seat's protocol level.
//
// stand-in until line 1 merges (21.8): this branch's bot package knows
// only random and policy, so standard and hard play the policy until
// the levels of line 1 arrive; ParseLevel then knows them and the
// stand-in is never reached. Remove it at the merge.
func botLevel(level string) botai.Level {
	if l, err := botai.ParseLevel(level); err == nil {
		return l
	}
	return botai.Policy // stand-in until line 1 merges (21.8)
}

// botBudget is the search budget of a level: hard thinks for the room's
// botThink; standard needs none (Decision 8.7).
func (r *Room) botBudget(level string) botai.Budget {
	if level == protocol.BotHard {
		return botai.Budget{Time: r.botThink}
	}
	return botai.Budget{}
}

// hasBots reports whether a bot sits at the table.
func (r *Room) hasBots() bool {
	return slices.ContainsFunc(r.seats, func(s roomSeat) bool { return s.bot != "" })
}

// botRng is the room's source of the bots' seeds; room goroutine only.
func (r *Room) botRng() *rand.Rand {
	if r.botRand == nil {
		r.botRand = rand.New(rand.NewPCG(r.setup.Seed, r.setup.Seed^0xb075)) //nolint:gosec // bot play, not security
	}
	return r.botRand
}

// botSecond is a part of the room's second.
func (r *Room) botSecond(f float64) time.Duration {
	return time.Duration(f * float64(r.second))
}

// remember feeds each hard seat's memory with what its client would get
// now: its view and its events (ADR 030, Decision 19.3).
//
// memory from 21.15: world.Memory does not exist in this branch yet.
// The merge fills this in: for each hard seat, feed r.memory[i] with
// ViewFor(base, its viewer) and FilterEvents(events, its viewer), then
// set r.known[i] from it. Until then r.known stays empty: a hard bot
// imagines every card it cannot see now.
func (r *Room) remember(_ []engine.Event, _ engine.ViewBase) {
	if len(r.known) != len(r.seats) {
		r.known = make([]engine.Known, len(r.seats))
	}
}

// drive drops the thinkers a broadcast made stale and notes the bot
// seats that must think now (Decision 12). Run starts their thinkers
// (startThinkers) with its context, once the request is handled.
func (r *Room) drive() {
	r.due = r.due[:0]
	if r.game == nil || !r.hasBots() {
		return
	}
	start := time.Now()
	defer r.noteBotWork(start)
	if len(r.thinkers) != len(r.seats) {
		r.thinkers = make([]*thinker, len(r.seats))
	}
	if len(r.known) != len(r.seats) {
		r.known = make([]engine.Known, len(r.seats))
	}
	if r.game.Over || r.paused() || r.rbVote != nil {
		r.stopThinkers()
		return
	}
	for i := range r.seats {
		wanted := r.wantsMove(i)
		if t := r.thinkers[i]; t != nil && (!wanted || t.step != r.step) {
			t.cancel()
			r.thinkers[i] = nil
		}
		if wanted && r.thinkers[i] == nil {
			r.due = append(r.due, i)
		}
	}
}

// wantsMove reports whether the game waits for bot seat i now.
func (r *Room) wantsMove(i int) bool {
	g := r.game
	if g == nil || g.Over || r.paused() || r.rbVote != nil || r.seats[i].bot == "" {
		return false
	}
	p := engine.PlayerID(i)
	return slices.Contains(g.Prompt().Who(), p) && len(g.Allowed(p)) > 0
}

// startThinkers starts a thinker for each seat drive noted, under the
// room's context: the room's end ends them all.
func (r *Room) startThinkers(ctx context.Context) {
	if len(r.due) == 0 {
		return
	}
	start := time.Now()
	defer r.noteBotWork(start)
	for _, i := range r.due {
		if r.thinkers[i] == nil && r.wantsMove(i) {
			r.think(ctx, i)
		}
	}
	r.due = r.due[:0]
}

// stopThinkers cancels every open thinker.
func (r *Room) stopThinkers() {
	for i, t := range r.thinkers {
		if t != nil {
			t.cancel()
			r.thinkers[i] = nil
		}
	}
}

// think starts a thinker for bot seat i. The room copies the game here,
// on its own goroutine (world.For); the thinker gets only the worlds.
// Its context comes from the room's: the room's steps cancel it alone,
// and the room's end stops every thinker.
func (r *Room) think(ctx context.Context, i int) {
	s := r.seats[i]
	p := engine.PlayerID(i)
	rng := r.botRng()
	seed, own := rng.Uint64(), r.game.PromptFor(p).Kind == engine.PromptPriority && r.game.Turn.Active == p
	answerPace, actionPace := r.botSecond(botAnswerPace), r.botSecond(botActionPace+rng.Float64())
	worlds, err := world.For(r.game, p, r.known[i], seed)
	ctx, cancel := context.WithCancel(ctx)
	r.thinkN++
	t := &thinker{id: r.thinkN, step: r.step, done: ctx.Done(), cancel: cancel}
	r.thinkers[i] = t
	in := botai.Input{
		Seat: p, Level: botLevel(s.bot), Rand: rand.New(rand.NewPCG(seed, rng.Uint64())), //nolint:gosec // bot play, not security
		Budget: r.botBudget(s.bot), World: botai.World(worlds),
	}
	go func() {
		begin := time.Now()
		a := &botAnswer{seat: i, id: t.id, step: t.step, err: err}
		if err == nil {
			a.intent, a.err = botai.Choose(ctx, in)
		}
		pace := answerPace
		if own && a.err == nil && a.intent.Kind != engine.IntentPass {
			pace = actionPace
		}
		wait := time.NewTimer(pace - time.Since(begin)) // thinking counts inside the pace
		defer wait.Stop()
		select {
		case <-ctx.Done():
			return
		case <-wait.C:
		}
		if ctx.Err() == nil {
			r.send(request{bot: a})
		}
	}()
}

// botMove applies a thinker's answer like a client's intent: record,
// rollback history, automatic passes (Decision 12.7). An answer for an
// older step, or one that comes while the game is over, paused or
// voting, is dropped. An error or a refused intent gives the bot's
// fallback move, and is logged (Decision 12.10): the response timer may
// be off, so a bot never waits for it.
func (r *Room) botMove(a *botAnswer) {
	if a.seat < 0 || a.seat >= len(r.thinkers) {
		return
	}
	t := r.thinkers[a.seat]
	if t == nil || t.id != a.id {
		return // dropped already
	}
	r.thinkers[a.seat] = nil
	if r.game == nil || r.game.Over || r.paused() || r.rbVote != nil || a.step != r.step {
		r.drive()
		return
	}
	p := engine.PlayerID(a.seat)
	name := r.seats[a.seat].name
	if a.err == nil {
		in := a.intent
		in.Player = p
		events, err := r.apply(in)
		if err == nil {
			r.after(events, true)
			return
		}
		log.Printf("server: game %s: %s: the engine refused %+v: %v; the fallback move goes", r.recInfo.game, name, in, err)
	} else {
		log.Printf("server: game %s: %s has no move: %v; the fallback move goes", r.recInfo.game, name, a.err)
	}
	for _, in := range r.botFallbacks(p) {
		in.Player = p
		events, err := r.apply(in)
		if err == nil {
			r.after(events, true)
			return
		}
		// Allowed gave what the engine refuses: an engine bug.
		log.Printf("server: game %s: engine bug: seed %d, step %d: %s's allowed %+v was refused: %v", r.recInfo.game, r.setup.Seed, r.step, name, in, err)
	}
	r.drive() // nothing went: the bot thinks again
}

// botFallbacks are a bot's moves when its own is missing or refused:
// first the forced answer of the response timer, then each allowed
// intent. forcedIntent completes a discard; the other allowed intents
// are complete as they are. Priority's forced answer is a pass or the
// end of the turn; when both are forbidden (a must-attack, R-ABIL-12),
// an allowed action goes.
func (r *Room) botFallbacks(p engine.PlayerID) []engine.Intent {
	var out []engine.Intent
	if in, ok := r.forcedIntent(p); ok {
		out = append(out, in)
	}
	if r.game.PromptFor(p).Kind != engine.PromptDiscard {
		out = append(out, r.game.Allowed(p)...)
	}
	return out
}

// noteBotWork keeps the longest time the driver held the room: the
// work only bots cost it, the copies for the thinkers among it (21.8).
// A bot's move then costs the room what a player's intent does; Run
// logs a slow one as "a bot's move".
func (r *Room) noteBotWork(start time.Time) {
	d := time.Since(start)
	if int64(d) > r.botWork.Load() {
		r.botWork.Store(int64(d))
	}
}

// botBanTurn schedules a bot's ban on its ban turn (Decision 14.1).
func (r *Room) botBanTurn() bool {
	m := r.match
	if r.seats[m.order[m.turn]].bot == "" {
		return false
	}
	gen := m.gen
	m.timer = time.AfterFunc(r.botSecond(botAnswerPace), func() { r.send(request{botBan: gen}) })
	return true
}

// botBan bans one random character of the pool for the bot whose ban
// turn it is (Decision 14.1).
func (r *Room) botBan(gen int) {
	m := r.match
	if m == nil || m.phase != protocol.PhaseBan || m.gen != gen || len(m.pool) == 0 {
		return
	}
	seat := m.order[m.turn]
	if r.seats[seat].bot == "" {
		return
	}
	r.banCard(seat, r.botRng().IntN(len(m.pool)))
}

// botPicks picks one random character of each bot's draft offer
// (Decision 14.2).
func (r *Room) botPicks() {
	m := r.match
	for i, s := range r.seats {
		if s.bot != "" && len(m.offers[i]) > 0 {
			m.picks[i] = m.offers[i][r.botRng().IntN(len(m.offers[i]))]
		}
	}
}
