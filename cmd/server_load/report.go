package main

import (
	"fmt"
	"io"
	"maps"
	"slices"
	"strconv"
	"strings"
	"sync"
	"time"
)

// printer writes the report and the progress. The first write error sticks
// and is returned at the end.
type printer struct {
	mu  sync.Mutex
	w   io.Writer
	err error
}

func (p *printer) printf(format string, args ...any) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.err == nil {
		_, p.err = fmt.Fprintf(p.w, format, args...)
	}
}

func (p *printer) firstErr() error {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.err
}

// text is the spread as "min …, median …, p95 …, max …".
func (s spread[T]) text(format func(T) string) string {
	if s.n == 0 {
		return "none"
	}
	return fmt.Sprintf("min %s, median %s, p95 %s, max %s",
		format(s.min), format(s.median), format(s.p95), format(s.max))
}

// ms is a duration rounded for the report.
func ms(d time.Duration) string {
	return d.Round(10 * time.Microsecond).String()
}

// pair is "median / max" of a latency spread.
func pair(s spread[time.Duration]) string {
	if s.n == 0 {
		return "-"
	}
	return ms(s.median) + " / " + ms(s.max)
}

// size is a byte count for the report.
func size(n int) string {
	if n >= 1<<20 {
		return fmt.Sprintf("%.1f MB", float64(n)/(1<<20))
	}
	return fmt.Sprintf("%.1f KB", float64(n)/(1<<10))
}

// problems says why a stage failed; none means it passed. A stage fails when
// a game failed or a response was slower than -max-response.
func (st stage) problems(cfg config) []string {
	var out []string
	if failed := st.games - st.finished - st.cut; failed > 0 {
		var first error
		for _, r := range st.results {
			if !r.finished && !r.cut {
				first = r.err
				break
			}
		}
		out = append(out, fmt.Sprintf("%d of %d games failed, first: %v", failed, st.games, first))
	}
	if r := st.col.responses(); r.n > 0 && r.max > cfg.maxResponse {
		out = append(out, fmt.Sprintf("the slowest response took %s, over the limit of %s", ms(r.max), cfg.maxResponse))
	}
	return out
}

func (st stage) passed(cfg config) bool {
	return len(st.problems(cfg)) == 0
}

// report prints one stage: the server, the games, the latency, the traffic
// and the other checks, then the verdict. Every bot has stopped by now.
func (st stage) report(p *printer, cfg config) {
	c := st.col
	p.printf("\n== %d games at once, %d players each ==\n", st.games, cfg.players)
	if w := c.server; w != nil {
		line(p, "server", "app %s, engine %s, protocol %d", w.App, w.Engine, w.Protocol)
		if w.Name != "" {
			line(p, "name", "%q", w.Name)
		}
	}
	line(p, "games", "%d finished, %d failed, in %s", st.finished, st.games-st.finished-st.cut, st.elapsed.Round(time.Second))
	if st.cut > 0 {
		line(p, "stopped", "%d games cut short when the stop share was reached", st.cut)
	}
	var steps []int
	var lengths []time.Duration
	for _, r := range st.results {
		if r.finished {
			steps = append(steps, r.steps)
			lengths = append(lengths, r.elapsed)
		}
	}
	line(p, "steps", "%s", spreadOf(steps).text(strconv.Itoa))
	line(p, "game time", "%s", spreadOf(lengths).text(ms))
	resp := c.responses()
	line(p, "response", "%d intents: %s", resp.n, resp.text(ms))
	ping := c.pings()
	line(p, "ping", "%d pings: %s", ping.n, ping.text(ms))
	hs := spreadOf(c.handshake)
	line(p, "handshake", "%d connections: %s", hs.n, hs.text(ms))
	line(p, "updates", "%d received, %s, %s per second; largest %s",
		c.updates, size(c.bytes), size(int(float64(c.bytes)/st.elapsed.Seconds())), size(c.largest))
	line(p, "errors", "%s", counts(c.errors))
	line(p, "notices", "%s", counts(c.notices))
	if problems := st.problems(cfg); len(problems) > 0 {
		line(p, "result", "FAIL: %s", strings.Join(problems, "; "))
	} else {
		line(p, "result", "PASS")
	}
}

// line prints one labelled row of the report.
func line(p *printer, label, format string, args ...any) {
	p.printf("%-11s %s\n", label, fmt.Sprintf(format, args...))
}

// counts lists the counts of a map by key, or "none".
func counts(m map[string]int) string {
	if len(m) == 0 {
		return "none"
	}
	var parts []string
	for _, k := range slices.Sorted(maps.Keys(m)) {
		parts = append(parts, fmt.Sprintf("%s %d", k, m[k]))
	}
	return strings.Join(parts, ", ")
}

// summary prints one row per stage and the highest load that passed.
func summary(p *printer, cfg config, stages []stage) {
	if len(stages) == 0 {
		return
	}
	p.printf("\n== summary ==\n")
	p.printf("%-7s %-10s %-22s %-22s %s\n", "games", "finished", "response median/max", "ping median/max", "result")
	best := 0
	for _, st := range stages {
		result := "FAIL"
		if st.passed(cfg) {
			result = "PASS"
			best = st.games
		}
		p.printf("%-7d %-10s %-22s %-22s %s\n", st.games, fmt.Sprintf("%d/%d", st.finished, st.games),
			pair(st.col.responses()), pair(st.col.pings()), result)
	}
	if best == 0 {
		p.printf("\nthe server did not pass even %d games at once\n", stages[0].games)
		return
	}
	p.printf("\nthe server passed %d games at once: %d players playing together\n", best, best*cfg.players)
}
