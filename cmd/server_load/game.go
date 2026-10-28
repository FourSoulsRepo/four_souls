package main

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"sync"
	"sync/atomic"
	"time"

	"github.com/FourSoulsRepo/four_souls/internal/protocol"
	engine "github.com/FourSoulsRepo/rules_engine"
)

// Game states, shown in the progress line.
const (
	stateRunning int32 = iota
	stateFinished
	stateFailed
	stateStopped // cut short by the stage, not a failure
)

// gameRun is one game the tool plays: its number in the stage, its progress
// and what the seats agreed on at the end.
type gameRun struct {
	number int
	start  time.Time
	limit  time.Duration
	cancel context.CancelFunc
	step   atomic.Int64 // the newest step a seat has seen
	turn   atomic.Int64 // the turn number at that step
	state  atomic.Int32

	mu      sync.Mutex
	err     error // the first failure
	ended   bool
	steps   int
	winners []int
}

// gameResult is what a finished or failed game left behind.
type gameResult struct {
	finished bool
	cut      bool // stopped early: -stop-at ended its stage
	err      error
	steps    int
	winners  []int
	elapsed  time.Duration
}

func newGameRun(number int) *gameRun {
	return &gameRun{number: number, start: time.Now()}
}

// fail records the first failure and stops the game's bots.
func (g *gameRun) fail(err error) {
	g.mu.Lock()
	if g.err == nil {
		g.err = err
		g.state.Store(stateFailed)
		if errors.Is(err, context.Canceled) {
			g.state.Store(stateStopped)
		}
	}
	g.mu.Unlock()
	g.cancel()
}

// finish records the end as one seat sees it. Every seat must see the same
// last step and winners.
func (g *gameRun) finish(step int, winners []int) error {
	g.mu.Lock()
	defer g.mu.Unlock()
	if !g.ended {
		g.ended, g.steps, g.winners = true, step, winners
		return nil
	}
	if step != g.steps || !slices.Equal(winners, g.winners) {
		return fmt.Errorf("the seats disagree on the end: step %d with winners %v, then step %d with winners %v",
			g.steps, g.winners, step, winners)
	}
	return nil
}

// progress records where the game is, for the progress line.
func (g *gameRun) progress(step, turn int) {
	g.step.Store(int64(step))
	g.turn.Store(int64(turn))
}

// status is the game's part of the progress line.
func (g *gameRun) status() string {
	switch g.state.Load() {
	case stateFinished:
		return fmt.Sprintf("#%d done", g.number)
	case stateFailed:
		return fmt.Sprintf("#%d FAILED", g.number)
	case stateStopped:
		return fmt.Sprintf("#%d stopped", g.number)
	default:
		return fmt.Sprintf("#%d turn %d step %d", g.number, g.turn.Load(), g.step.Load())
	}
}

// result reads what the game left behind.
func (g *gameRun) result() gameResult {
	g.mu.Lock()
	defer g.mu.Unlock()
	return gameResult{
		finished: g.err == nil && g.ended,
		err:      g.err,
		steps:    g.steps,
		winners:  g.winners,
		elapsed:  time.Since(g.start),
	}
}

// stopReason says why a game's bot stopped reading before the game ended.
// A timeout names the turn and step the game was stuck at.
func (g *gameRun) stopReason(ctx context.Context) error {
	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return fmt.Errorf("not over after %s, stuck at turn %d, step %d", g.limit, g.turn.Load(), g.step.Load())
	}
	return ctx.Err()
}

// playGame seats the bots at a new table, plays the game to its end and
// returns what happened. It stops the game after cfg.gameTimeout.
func playGame(parent context.Context, cfg config, g *gameRun, col *collector) gameResult {
	ctx, cancel := context.WithTimeout(parent, cfg.gameTimeout)
	defer cancel()
	g.cancel, g.limit = cancel, cfg.gameTimeout
	bots, err := openTable(ctx, cfg, g, col)
	defer closeAll(bots)
	if err != nil {
		if ctx.Err() != nil {
			err = g.stopReason(ctx)
		}
		g.fail(err)
		return g.result()
	}
	// A failure or the timeout cancels ctx; closing the sockets wakes every bot.
	go func() {
		<-ctx.Done()
		closeAll(bots)
	}()
	var wg sync.WaitGroup
	for _, b := range bots {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := b.play(ctx); err != nil {
				g.fail(err)
			}
		}()
	}
	wg.Wait()
	res := g.result()
	if res.finished {
		g.state.Store(stateFinished)
	}
	return res
}

// openTable dials the host, opens a table with it, then dials each other
// seat and has it join at once. A bot that waits without reading gets the
// lobby's list updates piled up, and the server drops it; so a bot joins
// right after its welcome. Every seat is ready at the end, so the game starts
// by itself (6.5). The bots it returns are open even when it fails.
func openTable(ctx context.Context, cfg config, g *gameRun, col *collector) ([]*bot, error) {
	var bots []*bot
	open := func(seat int) (*bot, error) {
		b, err := dial(ctx, cfg, g, seat, col)
		if b != nil {
			bots = append(bots, b)
		}
		return b, err
	}
	host, err := open(0)
	if err != nil {
		return bots, err
	}
	if err := host.send(protocol.TypeCreate, protocol.Create{Seats: cfg.players, Sets: cfg.sets}); err != nil {
		return bots, err
	}
	var tb protocol.Table
	if err := host.expect(protocol.TypeTable, &tb); err != nil {
		return bots, err
	}
	for seat := 1; seat < cfg.players; seat++ {
		b, err := open(seat)
		if err != nil {
			return bots, err
		}
		if err := b.send(protocol.TypeJoin, protocol.Join{Game: tb.Game}); err != nil {
			return bots, err
		}
		if err := b.expect(protocol.TypeTable, nil); err != nil {
			return bots, err
		}
	}
	for _, b := range bots {
		if err := b.send(protocol.TypeReady, protocol.Ready{Ready: true}); err != nil {
			return bots, err
		}
	}
	return bots, nil
}

// closeAll closes every bot's socket.
func closeAll(bots []*bot) {
	for _, b := range bots {
		_ = b.conn.Close() //nolint:errcheck // closing only ends the reads; the game is decided
	}
}

// winnersOf lists the winning seats as numbers.
func winnersOf(ids []engine.PlayerID) []int {
	out := make([]int, len(ids))
	for i, id := range ids {
		out[i] = int(id)
	}
	return out
}
