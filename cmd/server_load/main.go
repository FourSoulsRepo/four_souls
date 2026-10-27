// Command server_load plays whole games against a running dedicated server
// and reports how much load it takes.
//
//	server_load -url ws://game.example.com:4774/ws
//	server_load -url ws://127.0.0.1:4774 -max-games 8 -max-response 500ms
//
// Each game seats -players bots at one table and plays it to the end with
// random legal moves, as the headless tests do. The games run in stages: -start
// games at once, then twice as many, up to -max-games. A stage passes when
// every game ends and no response is slower than -max-response. The first
// stage that fails ends the run. The tool only reads from the server; the
// games it opens end on the server like any other.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

// config is what the flags say.
type config struct {
	url         string
	players     int
	start       int
	maxGames    int
	sets        []string
	gameTimeout time.Duration
	maxResponse time.Duration
	seed        uint64
	stopOnFirst bool
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if err := run(ctx, os.Args[1:], os.Stdout); err != nil {
		log.Fatal(err)
	}
}

// run plays the stages until one fails, the last one ran or ctx ends, and
// prints the report.
func run(ctx context.Context, args []string, out io.Writer) error {
	cfg, err := parse(args)
	if err != nil {
		return err
	}
	p := &printer{w: out}
	p.printf("load test of %s: %d players per game, games at once from %d, doubling up to %d\n",
		cfg.url, cfg.players, cfg.start, cfg.maxGames)
	var stages []stage
	for games := cfg.start; games <= cfg.maxGames; games *= 2 {
		st := runStage(ctx, cfg, games, p)
		if ctx.Err() != nil {
			p.printf("\ninterrupted\n")
			break
		}
		st.report(p, cfg)
		stages = append(stages, st)
		if !st.passed(cfg) {
			break
		}
	}
	summary(p, cfg, stages)
	return p.firstErr()
}

// parse reads the flags and checks them.
func parse(args []string) (config, error) {
	var cfg config
	fs := flag.NewFlagSet("server_load", flag.ContinueOnError)
	fs.StringVar(&cfg.url, "url", "", "WebSocket URL of the server, e.g. ws://host:4774 (required)")
	fs.IntVar(&cfg.players, "players", 4, "players per game, 2 to 4")
	fs.IntVar(&cfg.start, "start", 2, "games at once in the first stage; each stage doubles it")
	fs.IntVar(&cfg.maxGames, "max-games", 64, "most games at once")
	sets := fs.String("sets", "b2", "card sets to play, comma-separated")
	fs.DurationVar(&cfg.gameTimeout, "game-timeout", 30*time.Minute, "a game that has not ended by then fails")
	fs.DurationVar(&cfg.maxResponse, "max-response", time.Second, "a stage fails when a response takes longer")
	fs.Uint64Var(&cfg.seed, "seed", 0, "seed for the bots' moves; 0 uses the clock")
	fs.BoolVar(&cfg.stopOnFirst, "stop-on-first", false, "when one game of a stage ends, end the others at once instead of waiting for them")
	if err := fs.Parse(args); err != nil {
		return cfg, fmt.Errorf("flags: %w", err)
	}
	cfg.sets = strings.FieldsFunc(*sets, func(r rune) bool { return r == ',' })
	return cfg, cfg.check()
}

// check reports flags that cannot work, and fills in the defaults.
func (cfg *config) check() error {
	if cfg.url == "" {
		return errors.New("-url is required, e.g. -url ws://host:4774")
	}
	u, err := url.Parse(cfg.url)
	if err != nil {
		return fmt.Errorf("-url: %w", err)
	}
	switch {
	case u.Scheme != "ws" && u.Scheme != "wss":
		return fmt.Errorf("-url %s: the scheme must be ws:// or wss://", cfg.url)
	case u.Host == "":
		return fmt.Errorf("-url %s: no host", cfg.url)
	case cfg.players < 2 || cfg.players > 4:
		return errors.New("-players: 2 to 4")
	case cfg.start < 1 || cfg.maxGames < cfg.start:
		return errors.New("-start must be at least 1, and -max-games at least -start")
	case len(cfg.sets) == 0:
		return errors.New("-sets: at least one card set")
	case cfg.gameTimeout <= 0 || cfg.maxResponse <= 0:
		return errors.New("-game-timeout and -max-response must be positive")
	}
	if u.Path == "" || u.Path == "/" {
		u.Path = "/ws"
	}
	cfg.url = u.String()
	if cfg.seed == 0 {
		cfg.seed = uint64(time.Now().UnixNano()) //nolint:gosec // a seed for the moves, not security
	}
	return nil
}
