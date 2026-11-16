// Command bot_arena plays whole seeded games between bot levels and
// prints how they went (ADR 030, Decision 18):
//
//	bot_arena -games 50
//	bot_arena -games 20 -match policy,random -match policy,policy,policy
//
// Each -match is one line-up of levels, one per seat; the seats rotate
// from game to game. Without -match it plays the default line-ups. It
// prints wins by level, unfinished games and why, steps and time per
// move, and exits 1 when a gate fails: a refused bot intent, a loop guard
// (a game turn over 1000 intents other than passes, a game over 30000),
// in any line-up; fewer than 95% of games with a winner within the turn
// cap, in a line-up of player-visible levels with no random seat; or a
// win gate: standard under 70% of 2-player games against random,
// standard under 55% of 2-player games against flat (standard with
// every card worth alike), hard under 60% of 2-player games against
// standard, one hard among three standard not over 30%. Hard thinks
// -playouts playouts a move, with a memory of the cards it saw unless
// -memory=false.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"maps"
	"os"
	"os/signal"
	"slices"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/FourSoulsRepo/four_souls/internal/arena"
	"github.com/FourSoulsRepo/four_souls/internal/bot"
	engine "github.com/FourSoulsRepo/rules_engine"
)

// minFinished is the gate of Decision 18.4: the part of games that must
// end with a winner within the turn cap.
const minFinished = 0.95

// finishGated are the levels whose line-ups must pass minFinished: the
// levels players meet. Random and policy line-ups print their share for
// reference only.
var finishGated = map[bot.Level]bool{bot.Standard: true, bot.Hard: true}

// winGate is a gate on a line-up of Seats seats, one of Level and the
// rest of Against (Decision 18.5 – 18.7): Level must win at least Min
// of the games, or more than Min when Over.
type winGate struct {
	Level, Against bot.Level
	Seats          int
	Min            float64
	Over           bool
}

// winGates are the win gates.
var winGates = []winGate{
	{Level: bot.Standard, Against: bot.Random, Seats: 2, Min: 0.70},
	{Level: bot.Hard, Against: bot.Standard, Seats: 2, Min: 0.60},
	{Level: bot.Hard, Against: bot.Standard, Seats: 4, Min: 0.30, Over: true},
	{Level: bot.Standard, Against: bot.Flat, Seats: 2, Min: 0.55},
}

// winGateOf is the win gate of a line-up, if it has one.
func winGateOf(levels []bot.Level) (winGate, bool) {
	for _, w := range winGates {
		if len(levels) != w.Seats {
			continue
		}
		mine, theirs := 0, 0
		for _, l := range levels {
			if l == w.Level {
				mine++
			}
			if l == w.Against {
				theirs++
			}
		}
		if mine == 1 && theirs == w.Seats-1 {
			return w, true
		}
	}
	return winGate{}, false
}

// fails: a share of won games misses the gate.
func (w winGate) fails(share float64) bool {
	return share < w.Min || w.Over && share <= w.Min
}

// gated: the line-up holds a finish-gated level and no random seat.
func gated(levels []bot.Level) bool {
	if slices.Contains(levels, bot.Random) {
		return false
	}
	return slices.ContainsFunc(levels, func(l bot.Level) bool { return finishGated[l] })
}

// defaultMatches are the line-ups played without -match.
var defaultMatches = []string{
	"standard,random",
	"standard,flat",
	"hard,standard",
	"hard,standard,standard,standard",
	"standard,standard",
	"standard,standard,standard",
	"standard,standard,standard,standard",
	"policy,policy",
	"policy,policy,policy",
	"policy,policy,policy,policy",
}

// errGate is a failed gate; the report says which.
var errGate = errors.New("a gate failed")

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	err := run(ctx, os.Args[1:], os.Stdout)
	if errors.Is(err, errGate) {
		stop()
		os.Exit(1)
	}
	if err != nil {
		log.Fatal(err)
	}
}

// matchFlags collects repeated -match flags.
type matchFlags []string

func (m *matchFlags) String() string { return strings.Join(*m, " ") }

func (m *matchFlags) Set(s string) error {
	*m = append(*m, s)
	return nil
}

// run parses the flags, plays every match and prints the report.
func run(ctx context.Context, args []string, out io.Writer) error {
	fs := flag.NewFlagSet("bot_arena", flag.ContinueOnError)
	games := fs.Int("games", 50, "games per line-up")
	seed := fs.Uint64("seed", 1, "seed of each line-up's first game")
	turnCap := fs.Int("cap", arena.TurnCap, "game turns after which a game counts as unfinished")
	playouts := fs.Int("playouts", arena.DefaultPlayouts, "the hard level's playouts a move")
	memory := fs.Bool("memory", true, "give each hard seat a memory of the cards it saw (Decision 19)")
	var matches matchFlags
	fs.Var(&matches, "match", "a line-up of levels, one per seat, comma-separated (repeat it)")
	if err := fs.Parse(args); err != nil {
		return err
	}
	if len(matches) == 0 {
		matches = defaultMatches
	}
	var ms []arena.Match
	for _, s := range matches {
		var levels []bot.Level
		for name := range strings.SplitSeq(s, ",") {
			l, err := bot.ParseLevel(strings.TrimSpace(name))
			if err != nil {
				return err
			}
			levels = append(levels, l)
		}
		ms = append(ms, arena.Match{Levels: levels, Games: *games, Seed: *seed, Cap: *turnCap, Playouts: *playouts, NoMemory: !*memory})
	}
	failed := false
	pr := &printer{w: out}
	for _, m := range ms {
		start := time.Now()
		r, err := arena.Play(ctx, m)
		if err != nil {
			return err
		}
		if !report(pr, r, time.Since(start)) {
			failed = true
		}
		if ctx.Err() != nil {
			pr.printf("interrupted\n")
			return ctx.Err()
		}
	}
	if pr.err != nil {
		return pr.err
	}
	if failed {
		return errGate
	}
	return nil
}

// printer writes the report and keeps the first write error.
type printer struct {
	w   io.Writer
	err error
}

func (p *printer) printf(format string, args ...any) {
	if p.err == nil {
		_, p.err = fmt.Fprintf(p.w, format, args...)
	}
}

// report prints one match; it returns false when a gate failed.
func report(out *printer, r arena.Result, took time.Duration) bool {
	m := r.Match
	out.printf("%s: %d games, %d with a winner (%.1f%%) within %d turns, in %s\n",
		m.Name(), r.Games, r.Finished, 100*r.FinishedShare(), m.Cap, took.Round(time.Millisecond))
	if slices.Contains(m.Levels, bot.Hard) {
		mem := "a memory"
		if m.NoMemory {
			mem = "no memory"
		}
		out.printf("  hard: %d playouts a move, %s\n", m.Playouts, mem)
	}
	for _, l := range slices.Sorted(maps.Keys(r.Seats)) {
		chance := float64(r.Seats[l]) / float64(len(m.Levels))
		out.printf("  %-8s wins %d of %d games (%.1f by chance)\n", l, r.Wins[l], r.Games, chance)
	}
	if r.Games > 0 && r.Moves > 0 && r.Steps > 0 {
		turns := 0
		for _, t := range r.Turns {
			turns += t
		}
		out.printf("  turns a game: median %d, max %d; %.0f intents a turn, %.0f%% passes; %s a move\n",
			r.MedianTurns(), r.MaxTurns(), float64(r.Steps)/float64(max(turns, 1)),
			100*float64(r.Passes)/float64(r.Steps), (r.Think / time.Duration(r.Moves)).Round(time.Microsecond))
		out.printf("  %d player deaths, %d monsters killed (%.2f deaths a kill)\n",
			r.PlayerDeaths, r.MonsterKills, float64(r.PlayerDeaths)/float64(max(r.MonsterKills, 1)))
	}
	for _, l := range slices.Sorted(maps.Keys(r.Levels)) {
		c, seats := r.Levels[l], max(r.Seats[l], 1)
		out.printf("  %-8s %s a move; %.0f%% of its turns with no attack; a seat's game: %.1f loot plays, %.1f purchases\n",
			l, (c.Think / time.Duration(max(c.Moves, 1))).Round(time.Microsecond),
			100*float64(c.NoAttack)/float64(max(c.Turns, 1)),
			float64(c.LootPlays)/float64(seats), float64(c.Purchases)/float64(seats))
	}
	for _, s := range r.Stalls {
		out.printf("  stall: seed %d, turn %d, souls %v, last steps %s\n", s.Seed, s.Turns, s.Souls, kinds(s.Last))
	}
	ok := true
	for _, x := range r.Refused {
		out.printf("  refused: %s\n", x)
		ok = false
	}
	if r.Games > 0 && gated(m.Levels) && r.FinishedShare() < minFinished {
		out.printf("  FAIL: under %.0f%% of games with a winner (Decision 18.4)\n", 100*minFinished)
		ok = false
	}
	if w, has := winGateOf(m.Levels); has && r.Games > 0 {
		share := float64(r.Wins[w.Level]) / float64(r.Games)
		if w.fails(share) {
			out.printf("  FAIL: %s wins %.0f%% against %s, under %.0f%% (Decision 18)\n", w.Level, 100*share, w.Against, 100*w.Min)
			ok = false
		}
	}
	if len(r.Refused) > 0 {
		out.printf("  FAIL: %d refused bot intents (Decision 18.8)\n", len(r.Refused))
	}
	if !gated(m.Levels) {
		out.printf("  (finish share for reference: no finish gate on this line-up)\n")
	}
	for _, x := range r.Guards {
		out.printf("  FAIL: loop guard: %s (Decision 18.4)\n", x)
		ok = false
	}
	return ok
}

// kindNames name the intent kinds for the stall lines.
var kindNames = map[engine.IntentKind]string{
	engine.IntentPass: "pass", engine.IntentEndTurn: "end", engine.IntentDiscard: "discard",
	engine.IntentPlayLoot: "loot", engine.IntentChoose: "choose", engine.IntentAttack: "attack",
	engine.IntentPurchase: "buy", engine.IntentActivate: "use", engine.IntentOrder: "order",
	engine.IntentThrow: "throw",
}

// kinds prints intent counts in kind order, e.g. "pass 120 end 30".
func kinds(c map[engine.IntentKind]int) string {
	var b strings.Builder
	for k := engine.IntentPass; k <= engine.IntentThrow; k++ {
		if c[k] > 0 {
			b.WriteString(kindNames[k] + " " + strconv.Itoa(c[k]) + " ")
		}
	}
	return strings.TrimSpace(b.String())
}
