// Package arena plays whole games between bot levels and counts the
// results (ADR 030, Decision 18). cmd/bot_arena prints them; tests use
// it to check that no bot intent is refused.
//
// The arena holds the real game. Each seat gets only its worlds, made by
// internal/world from the game as the server's bot driver makes them; a
// seat never sees the game itself.
package arena

import (
	"context"
	"fmt"
	"math/rand/v2"
	"runtime"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/FourSoulsRepo/four_souls/internal/bot"
	"github.com/FourSoulsRepo/four_souls/internal/world"
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards"
)

// The limits of a game (Decision 18.4).
const (
	// TurnCap ends a game with no winner after this many game turns
	// (Turn.Number); it then counts as unfinished.
	TurnCap = 150
	// MaxTurnIntents and MaxGameIntents catch a loop: a game turn with
	// more applied intents other than passes, or a game with more
	// applied intents, ends the game as a guard hit, a failed gate. A
	// turn counts no passes: in a 4-player game every stack item waits
	// for four, so a long turn of real plays (a paid ability used 40
	// times) passed 1000 intents with no loop.
	MaxTurnIntents = 1000
	MaxGameIntents = 30000
	// DefaultPlayouts is the hard level's playouts a move: what a weak
	// server gives in its think time (Decision 18). The arena counts
	// playouts, never time, so its results repeat (Decision 9.10).
	DefaultPlayouts = 150
)

// Match is a series of games between one line-up of levels.
type Match struct {
	// Levels are the seats' levels, 2 to 4; game k rotates them by k
	// seats, so each level plays from each seat (Decision 18.1).
	Levels []bot.Level
	Games  int
	// Seed is the first game's seed; game k plays Seed+k.
	Seed uint64
	// Cap is the turn cap; 0 is TurnCap.
	Cap int
	// TurnIntents and GameIntents are the loop guards; 0 is
	// MaxTurnIntents and MaxGameIntents.
	TurnIntents, GameIntents int
	// Playouts is the hard level's playouts a move; 0 is
	// DefaultPlayouts.
	Playouts int
	// NoMemory plays the hard seats without a memory (Decision 19), to
	// measure what it is worth; else each hard seat has one.
	NoMemory bool
	// Watch, when set, sees the real game of game k before each bot
	// move, at step. It runs on that game's goroutine, games side by
	// side, and must copy what it keeps: make bot-values keeps positions
	// of arena games (Decision 20.3).
	Watch func(k, step int, g *engine.Game)
}

// Name is the line-up, e.g. "policy-policy".
func (m Match) Name() string {
	names := make([]string, len(m.Levels))
	for i, l := range m.Levels {
		names[i] = string(l)
	}
	return strings.Join(names, "-")
}

// Result counts the games of a match.
type Result struct {
	Match    Match
	Games    int
	Finished int // ended with a winner within the cap
	// Wins counts won games by level; a shared win counts for each
	// winner's level. Seats counts the seats each level held, so
	// Wins/Seats·players is its share against chance.
	Wins   map[bot.Level]int
	Seats  map[bot.Level]int
	Steps  int // applied intents
	Passes int // of them passes
	// Turns is each game's length in game turns, in game order.
	Turns []int
	Moves int // intents chosen by a bot (all of them)
	Think time.Duration
	// PlayerDeaths and MonsterKills count deaths of players and of
	// monsters, to judge the bots' fights.
	PlayerDeaths int
	MonsterKills int
	// Levels counts each level's own play, to compare levels in the
	// same games.
	Levels map[bot.Level]LevelCount
	// Refused lists every bot intent the engine refused (Decision 18.8).
	Refused []string
	// Stalls describes each game that ended at the cap.
	Stalls []Stall
	// Guards lists each game a loop guard ended, with its seed.
	Guards []string
}

// LevelCount is one level's play over a match: its seats' moves and
// turns. Divided by Result.Seats it is a seat's share a game.
type LevelCount struct {
	Moves int
	Think time.Duration
	// Turns are the turns one of its seats was the active player;
	// NoAttack those of them with no attack.
	Turns, NoAttack int
	// LootPlays and Purchases are its seats' loot plays (on any turn)
	// and purchases.
	LootPlays, Purchases int
}

// count adds one bot move of level l, intent in, thought for d.
func (r *Result) count(l bot.Level, in engine.Intent, d time.Duration) {
	c := r.Levels[l]
	c.Moves++
	c.Think += d
	if in.Kind == engine.IntentPlayLoot {
		c.LootPlays++
	}
	if in.Kind == engine.IntentPurchase {
		c.Purchases++
	}
	r.Levels[l] = c
}

// countTurn adds one finished turn of a seat of level l; attacked: the
// seat declared an attack in it.
func (r *Result) countTurn(l bot.Level, attacked bool) {
	c := r.Levels[l]
	c.Turns++
	if !attacked {
		c.NoAttack++
	}
	r.Levels[l] = c
}

// Stall is a game with no winner within the cap: where it stood.
type Stall struct {
	Seed  uint64
	Turns int
	// Souls are each seat's soul value at the cap.
	Souls []int
	// Last counts the intent kinds of the last steps: what the bots
	// kept doing.
	Last map[engine.IntentKind]int
}

// FinishedShare is the part of games that ended with a winner.
func (r Result) FinishedShare() float64 {
	if r.Games == 0 {
		return 0
	}
	return float64(r.Finished) / float64(r.Games)
}

// Play plays the match's games, side by side on every core; the result
// is the same as one after another. ctx stops games not yet started.
func Play(ctx context.Context, m Match) (Result, error) {
	if m.Cap == 0 {
		m.Cap = TurnCap
	}
	if m.TurnIntents == 0 {
		m.TurnIntents = MaxTurnIntents
	}
	if m.GameIntents == 0 {
		m.GameIntents = MaxGameIntents
	}
	if m.Playouts == 0 {
		m.Playouts = DefaultPlayouts
	}
	if n := len(m.Levels); n < 2 || n > 4 {
		return Result{}, fmt.Errorf("arena: %d seats; a game is for 2 to 4", n)
	}
	games := make([]Result, m.Games)
	errs := make([]error, m.Games)
	next := make(chan int)
	var wg sync.WaitGroup
	for range min(runtime.GOMAXPROCS(0), max(m.Games, 1)) {
		wg.Go(func() {
			for k := range next {
				games[k] = Result{Wins: map[bot.Level]int{}, Seats: map[bot.Level]int{}, Levels: map[bot.Level]LevelCount{}}
				errs[k] = playGame(ctx, m, k, &games[k])
			}
		})
	}
	for k := range m.Games {
		if ctx.Err() != nil {
			break
		}
		next <- k
	}
	close(next)
	wg.Wait()
	r := Result{Match: m, Wins: map[bot.Level]int{}, Seats: map[bot.Level]int{}, Levels: map[bot.Level]LevelCount{}}
	for k, g := range games {
		if errs[k] != nil {
			return r, errs[k]
		}
		r.add(g)
	}
	return r, nil
}

// add counts one game's result into r.
func (r *Result) add(g Result) {
	r.Games += g.Games
	r.Finished += g.Finished
	for l, n := range g.Wins {
		r.Wins[l] += n
	}
	for l, n := range g.Seats {
		r.Seats[l] += n
	}
	r.Steps += g.Steps
	r.Passes += g.Passes
	r.Turns = append(r.Turns, g.Turns...)
	r.PlayerDeaths += g.PlayerDeaths
	r.MonsterKills += g.MonsterKills
	for l, c := range g.Levels {
		a := r.Levels[l]
		a.Moves += c.Moves
		a.Think += c.Think
		a.Turns += c.Turns
		a.NoAttack += c.NoAttack
		a.LootPlays += c.LootPlays
		a.Purchases += c.Purchases
		r.Levels[l] = a
	}
	r.Moves += g.Moves
	r.Think += g.Think
	r.Refused = append(r.Refused, g.Refused...)
	r.Stalls = append(r.Stalls, g.Stalls...)
	r.Guards = append(r.Guards, g.Guards...)
}

// lastSteps is how many steps a stall's Last looks back.
const lastSteps = 300

// playGame plays game k of the match and adds it to r.
func playGame(ctx context.Context, m Match, k int, r *Result) error {
	seed := m.Seed + uint64(k) //nolint:gosec // k counts games from 0
	n := len(m.Levels)
	s := engine.Setup{
		Seed: seed, Players: n, Sets: cards.Sets(), BonusSouls: true,
		Rooms: seed%2 == 1, DeckRatios: seed%3 == 0,
	}
	g, setup, err := engine.NewGame(s)
	if err != nil {
		return fmt.Errorf("arena: game %d: %w", seed, err)
	}
	levels := make([]bot.Level, n)
	rngs := make([]*rand.Rand, n)
	for i := range n {
		levels[i] = m.Levels[(i+k)%n]
		r.Seats[levels[i]]++
		rngs[i] = rand.New(rand.NewPCG(seed, uint64(i)+1)) //nolint:gosec // bot moves, not security
	}
	mems := memories(levels, m.NoMemory)
	observe(g, mems, setup)
	worldSeeds := rand.New(rand.NewPCG(seed, 0)) //nolint:gosec // imagined deals, not security
	var kinds []engine.IntentKind
	step, turn, inTurn := 0, g.Turn.Number, 0
	attacked := false // the active player attacked this turn
	guard := ""
	for ; !g.Over && g.Turn.Number <= m.Cap; step++ {
		if g.Turn.Number != turn {
			turn, inTurn = g.Turn.Number, 0
		}
		if inTurn > m.TurnIntents {
			guard = fmt.Sprintf("game %d: turn %d over %d intents other than passes", seed, turn, m.TurnIntents)
			break
		}
		if step >= m.GameIntents {
			guard = fmt.Sprintf("game %d: over %d intents, turn %d", seed, m.GameIntents, turn)
			break
		}
		if m.Watch != nil {
			m.Watch(k, step, g)
		}
		p := g.Prompt().Who()[0]
		var known engine.Known
		if mems[p] != nil {
			known = mems[p].Known()
		}
		worlds, err := world.For(g, p, known, worldSeeds.Uint64())
		if err != nil {
			return fmt.Errorf("arena: game %d step %d: %w", seed, step, err)
		}
		start := time.Now()
		in, err := bot.Choose(ctx, bot.Input{
			Seat: p, Level: levels[p], Rand: rngs[p], World: bot.World(worlds),
			Budget: bot.Budget{Playouts: m.Playouts},
		})
		took := time.Since(start)
		r.Think += took
		r.Moves++
		active, number := g.Turn.Active, g.Turn.Number
		var events []engine.Event
		if err == nil {
			events, err = g.Apply(in)
		}
		if err != nil {
			r.Refused = append(r.Refused, fmt.Sprintf("game %d step %d seat %d (%s): %+v: %v", seed, step, p, levels[p], in, err))
			in = forced(g, p)
			if events, err = g.Apply(in); err != nil {
				return fmt.Errorf("arena: game %d step %d: the forced answer %+v: %w", seed, step, in, err)
			}
		}
		observe(g, mems, events)
		r.countDeaths(events)
		r.count(levels[p], in, took)
		attacked = attacked || in.Kind == engine.IntentAttack && p == active
		if g.Turn.Number != number {
			if number > 0 { // turn 0: the start-of-game choices
				r.countTurn(levels[active], attacked)
			}
			attacked = false
		}
		if in.Kind == engine.IntentPass {
			r.Passes++
		} else {
			inTurn++
		}
		kinds = append(kinds, in.Kind)
	}
	if g.Over {
		r.countTurn(levels[g.Turn.Active], attacked)
	}
	r.Games++
	r.Steps += step
	r.Turns = append(r.Turns, min(g.Turn.Number, m.Cap))
	if guard != "" {
		r.Guards = append(r.Guards, guard)
	}
	if g.Over {
		r.Finished++
		for _, w := range g.Winners {
			r.Wins[levels[w]]++
		}
		return nil
	}
	st := Stall{Seed: seed, Turns: min(g.Turn.Number, m.Cap), Last: map[engine.IntentKind]int{}}
	for _, pl := range g.Players {
		st.Souls = append(st.Souls, g.SoulValue(pl.ID))
	}
	for _, kd := range kinds[max(len(kinds)-lastSteps, 0):] {
		st.Last[kd]++
	}
	r.Stalls = append(r.Stalls, st)
	return nil
}

// memories are the seats' memories: one for each hard seat, unless
// none is wanted (Decision 19.13).
func memories(levels []bot.Level, none bool) []*world.Memory {
	out := make([]*world.Memory, len(levels))
	for i, l := range levels {
		if l == bot.Hard && !none {
			out[i] = world.NewMemory(engine.PlayerID(i))
		}
	}
	return out
}

// observe feeds each memory the step as its seat's client gets it: the
// seat's view and the events filtered for it (Decision 19.3).
func observe(g *engine.Game, mems []*world.Memory, events []engine.Event) {
	var base engine.ViewBase
	built := false
	for i, mem := range mems {
		if mem == nil {
			continue
		}
		if !built {
			base, built = g.ViewBase(), true
		}
		v := engine.Viewer{Kind: engine.ViewPlayer, Player: engine.PlayerID(i)}
		mem.Observe(g.ViewFor(base, v), engine.FilterEvents(events, v))
	}
}

// countDeaths counts the players and monsters that died in events.
func (r *Result) countDeaths(events []engine.Event) {
	for _, e := range events {
		if e.Kind != engine.EvDied {
			continue
		}
		if e.Player == engine.NoPlayer {
			r.MonsterKills++
		} else {
			r.PlayerDeaths++
		}
	}
}

// MedianTurns is the middle game length, in game turns.
func (r Result) MedianTurns() int {
	if len(r.Turns) == 0 {
		return 0
	}
	s := slices.Sorted(slices.Values(r.Turns))
	return s[len(s)/2]
}

// MaxTurns is the longest game, in game turns.
func (r Result) MaxTurns() int {
	if len(r.Turns) == 0 {
		return 0
	}
	return slices.Max(r.Turns)
}

// forced is the answer the arena sends for a seat whose bot failed: the
// first allowed intent, a discard cut to its count (as the server's
// forced answer, N-07).
func forced(g *engine.Game, p engine.PlayerID) engine.Intent {
	allowed := g.Allowed(p)
	if len(allowed) == 0 {
		return engine.Intent{Player: p}
	}
	in := allowed[0]
	if in.Kind == engine.IntentDiscard {
		in.Objects = slices.Clone(in.Objects[:g.PromptFor(p).Count])
	}
	return in
}
