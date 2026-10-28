package main

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// progressEvery is how often the progress line is printed.
const progressEvery = 2 * time.Second

// stage is one run at a fixed number of games at once.
type stage struct {
	games    int
	results  []gameResult
	col      *collector
	elapsed  time.Duration
	finished int
	cut      int // games stopped early, when -stop-at ended the stage
}

// stopAfter is how many finished games stop a stage: stopAt percent of the
// games, rounded up, at least one. 0 means the stage waits for every game.
func stopAfter(games, stopAt int) int {
	if stopAt == 0 {
		return 0
	}
	return max(1, (games*stopAt+99)/100)
}

// runStage plays games at once and returns when every game has ended or
// failed. With -stop-at, the games still running stop once that share of the
// stage has finished. While it runs, it prints the progress of each game.
func runStage(ctx context.Context, cfg config, games int, p *printer) stage {
	col := newCollector()
	runs := make([]*gameRun, games)
	for i := range runs {
		runs[i] = newGameRun(i + 1)
	}
	stageCtx, stopStage := context.WithCancel(ctx)
	defer stopStage()
	threshold := stopAfter(games, cfg.stopAt)
	var finishedNow atomic.Int32
	var stoppedEarly atomic.Bool
	results := make([]gameResult, games)
	stopped := make(chan struct{})
	shown := make(chan struct{})
	go func() {
		showProgress(p, runs, stopped)
		close(shown)
	}()
	start := time.Now()
	var wg sync.WaitGroup
	for i, g := range runs {
		wg.Add(1)
		go func() {
			defer wg.Done()
			results[i] = playGame(stageCtx, cfg, g, col)
			if results[i].finished && threshold > 0 && int(finishedNow.Add(1)) >= threshold {
				stoppedEarly.Store(true)
				stopStage()
			}
		}()
	}
	wg.Wait()
	elapsed := time.Since(start)
	close(stopped)
	<-shown
	st := stage{games: games, results: results, col: col, elapsed: elapsed}
	for _, r := range results {
		if r.finished {
			st.finished++
		}
	}
	// Games that only stopped because the stage ended are not failures.
	for i, r := range results {
		if !r.finished && stoppedEarly.Load() && errors.Is(r.err, context.Canceled) {
			results[i].cut = true
			st.cut++
		}
	}
	return st
}

// showProgress prints where every game is, until stopped closes.
func showProgress(p *printer, runs []*gameRun, stopped <-chan struct{}) {
	tick := time.NewTicker(progressEvery)
	defer tick.Stop()
	for {
		select {
		case <-stopped:
			return
		case <-tick.C:
			p.printf("  %s\n", progressLine(runs))
		}
	}
}

// maxProgressParts is how many games the progress line lists one by one.
const maxProgressParts = 8

// progressLine has one part per game, or counts when there are many games.
func progressLine(runs []*gameRun) string {
	if len(runs) > maxProgressParts {
		var running, done, failed, stopped int
		for _, g := range runs {
			switch g.state.Load() {
			case stateFinished:
				done++
			case stateFailed:
				failed++
			case stateStopped:
				stopped++
			default:
				running++
			}
		}
		return fmt.Sprintf("%d running, %d done, %d failed, %d stopped", running, done, failed, stopped)
	}
	parts := make([]string, len(runs))
	for i, g := range runs {
		parts[i] = g.status()
	}
	return strings.Join(parts, " | ")
}
