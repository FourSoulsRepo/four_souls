// Command server runs Four Souls games without a UI (N-01).
//
//	server -port 4774 -records ./records
//	server -config server.json
//
// Flags override the config file. See the wiki page "Dedicated server on
// a VPS".
package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/FourSoulsRepo/four_souls/internal/legal"
	"github.com/FourSoulsRepo/four_souls/internal/server"
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
func run(ctx context.Context, args []string, out io.Writer) error {
	cfg, err := parse(args)
	if err != nil {
		return err
	}
	if err = printBanner(out); err != nil {
		return err
	}
	s, err := server.Start(ctx, cfg)
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

// parse reads the config file, if any, then lets flags override it.
func parse(args []string) (server.Config, error) {
	fs := flag.NewFlagSet("server", flag.ContinueOnError)
	def := server.DefaultConfig()
	config := fs.String("config", "", "JSON config file; flags override it")
	addr := fs.String("addr", def.Addr, "listen only on this address (default: all)")
	port := fs.Int("port", def.Port, "TCP port for players")
	records := fs.String("records", def.Records, "folder for match records")
	retention := fs.Int("retention", def.Retention, "days to keep records; 0 keeps them forever")
	saves := fs.String("saves", def.Saves, "folder for saved games")
	if err := fs.Parse(args); err != nil {
		return def, fmt.Errorf("flags: %w", err)
	}
	cfg := def
	if *config != "" {
		var err error
		if cfg, err = server.LoadConfig(*config); err != nil {
			return cfg, err
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
		}
	})
	return cfg, cfg.Check()
}

// printBanner writes the fan-game notice and the versions (L-02).
func printBanner(w io.Writer) error {
	v := version.Get()
	_, err := fmt.Fprintf(w, "%s\nFour Souls dedicated server\napp %s, rules engine %s\n",
		legal.PlainText(), v.App, v.Engine)
	if err != nil {
		return fmt.Errorf("print banner: %w", err)
	}
	return nil
}
