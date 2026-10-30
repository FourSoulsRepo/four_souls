// Command server runs Four Souls games without a UI (N-01).
//
//	server -port 4774 -records ./records
//	server -config server.json
//	FOUR_SOULS_PORT=4774 FOUR_SOULS_RECORDS=/data/records server
//
// Every flag has an ENV variable: FOUR_SOULS_ and the flag's name in
// capitals (N-19). ENV overrides the config file; flags override both.
// See the wiki page "Dedicated server on a VPS".
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/FourSoulsRepo/four_souls/internal/legal"
	"github.com/FourSoulsRepo/four_souls/internal/server"
	"github.com/FourSoulsRepo/four_souls/internal/tracing"
	"github.com/FourSoulsRepo/four_souls/internal/version"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if err := run(ctx, os.Args[1:], os.Stdout); err != nil {
		log.Fatal(err)
	}
}

// run starts the server with the given arguments and stops it when ctx
// ends.
func run(ctx context.Context, args []string, out io.Writer) (err error) {
	cfg, trace, err := parseAll(args)
	if err != nil {
		return err
	}
	if err = printBanner(out); err != nil {
		return err
	}
	stopTrace, err := tracing.Start(ctx, trace)
	if err != nil {
		return err
	}
	// The trace and profiles are written when the server has stopped.
	defer func() { err = errors.Join(err, stopTrace()) }()
	cfg.KeepGames = true // an update or restart does not end the games (N-20)
	// Shutdown stops the server, not the signal: the players hear of the
	// restart before their games stop.
	s, err := server.Start(context.WithoutCancel(ctx), cfg)
	if err != nil {
		return err
	}
	if _, err = fmt.Fprintf(out, "listening on %s (ws://…/ws); records in %s\n", s.Addr(), cfg.Records); err != nil {
		return fmt.Errorf("print: %w", err)
	}
	<-ctx.Done()
	// Shutting down gets its own time, after the signal ended ctx.
	shut, cancel := context.WithTimeout(context.WithoutCancel(ctx), 10*time.Second)
	defer cancel()
	if err = s.Shutdown(shut); err != nil {
		return err
	}
	_, err = fmt.Fprintln(out, "stopped")
	return err
}

// envPrefix starts the ENV variable of every flag (N-19).
const envPrefix = "FOUR_SOULS_"

// parse reads the arguments into the server's config; see parseAll.
func parse(args []string) (server.Config, error) {
	cfg, _, err := parseAll(args)
	return cfg, err
}

// parseAll reads the config file, if any, then lets ENV and flags override
// it. It also reads the tracing options (ADR 024).
func parseAll(args []string) (server.Config, tracing.Options, error) {
	fs := flag.NewFlagSet("server", flag.ContinueOnError)
	def := server.DefaultConfig()
	config := fs.String("config", "", "JSON config file; ENV and flags override it")
	addr := fs.String("addr", def.Addr, "listen only on this address (default: all)")
	port := fs.Int("port", def.Port, "TCP port for players")
	records := fs.String("records", def.Records, "folder for match records")
	retention := fs.Int("retention", def.Retention, "days to keep records; 0 keeps them forever")
	saves := fs.String("saves", def.Saves, "folder for saved games")
	presence := fs.Bool("presence", def.Presence, "relay live aiming and hover; false saves CPU and bandwidth")
	name := fs.String("name", def.Name, "server name players see instead of the address")
	emptyMinutes := fs.Int("empty_minutes", def.EmptyMinutes, "end a game nobody is connected to for this many minutes; 0 never")
	gameHours := fs.Int("game_hours", def.GameHours, "end a game this many hours after it started; 0 never")
	trace := tracing.Flags(fs)
	// ENV is set like a flag before the flags, so a flag wins. An empty
	// value counts as unset, as compose passes one for a blank setting.
	var envErr error
	fs.VisitAll(func(f *flag.Flag) {
		key := envPrefix + strings.ToUpper(f.Name)
		if v := os.Getenv(key); v != "" && envErr == nil {
			if err := fs.Set(f.Name, v); err != nil {
				envErr = fmt.Errorf("env %s=%q: %w", key, v, err)
			}
		}
	})
	if envErr != nil {
		return def, tracing.Options{}, envErr
	}
	if err := fs.Parse(args); err != nil {
		return def, tracing.Options{}, fmt.Errorf("flags: %w", err)
	}
	cfg := def
	if *config != "" {
		var err error
		if cfg, err = server.LoadConfig(*config); err != nil {
			return cfg, tracing.Options{}, err
		}
	}
	fs.Visit(func(f *flag.Flag) {
		switch f.Name {
		case "addr":
			cfg.Addr = *addr
		case "port":
			cfg.Port = *port
		case "records":
			cfg.Records = *records
		case "retention":
			cfg.Retention = *retention
		case "saves":
			cfg.Saves = *saves
		case "presence":
			cfg.Presence = *presence
		case "name":
			cfg.Name = *name
		case "empty_minutes":
			cfg.EmptyMinutes = *emptyMinutes
		case "game_hours":
			cfg.GameHours = *gameHours
		}
	})
	return cfg, *trace, cfg.Check()
}

// printBanner writes the fan-game notice and the versions (L-02).
func printBanner(w io.Writer) error {
	v := version.Get()
	_, err := fmt.Fprintf(w, "%s\nFour Souls dedicated server\napp %s, rules engine %s\ncpus %d, GOMAXPROCS %d\n",
		legal.PlainText(), v.App, v.Engine, runtime.NumCPU(), runtime.GOMAXPROCS(0))
	if err != nil {
		return fmt.Errorf("print banner: %w", err)
	}
	return nil
}
