package main

import (
	"bytes"
	"context"
	"net"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/FourSoulsRepo/four_souls/internal/server"
)

func TestSpreadOf(t *testing.T) {
	got := spreadOf([]time.Duration{5, 1, 9, 3, 7})
	want := spread[time.Duration]{n: 5, min: 1, median: 5, p95: 9, max: 9}
	if got != want {
		t.Fatalf("spreadOf = %+v, want %+v", got, want)
	}
	if empty := spreadOf([]int{}); empty.n != 0 {
		t.Fatalf("spreadOf(empty).n = %d", empty.n)
	}
	hundred := make([]int, 100)
	for i := range hundred {
		hundred[i] = i + 1
	}
	if p := spreadOf(hundred).p95; p != 95 {
		t.Fatalf("p95 of 1..100 = %d, want 95", p)
	}
}

func TestParseChecks(t *testing.T) {
	cases := map[string]struct {
		args []string
		ok   bool
	}{
		"no url":        {args: []string{}, ok: false},
		"http scheme":   {args: []string{"-url", "http://host"}, ok: false},
		"one player":    {args: []string{"-url", "ws://host", "-players", "1"}, ok: false},
		"max below":     {args: []string{"-url", "ws://host", "-start", "8", "-max-games", "4"}, ok: false},
		"no host":       {args: []string{"-url", "ws:///ws"}, ok: false},
		"plain address": {args: []string{"-url", "ws://host:4774"}, ok: true},
	}
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			_, err := parse(c.args)
			if (err == nil) != c.ok {
				t.Fatalf("parse(%v) error %v, want ok %v", c.args, err, c.ok)
			}
		})
	}
}

func TestParseAddsWSPath(t *testing.T) {
	cfg, err := parse([]string{"-url", "ws://host:4774"})
	if err != nil {
		t.Fatal(err)
	}
	if cfg.url != "ws://host:4774/ws" {
		t.Fatalf("url = %s", cfg.url)
	}
	if cfg.seed == 0 {
		t.Fatal("no seed picked")
	}
}

// TestRunAgainstServer plays two 2-player games at once on a local server
// and expects the stage to pass (6.10).
func TestRunAgainstServer(t *testing.T) {
	if testing.Short() {
		t.Skip("plays whole games")
	}
	dir := t.TempDir()
	cfg := server.DefaultConfig()
	cfg.Addr, cfg.Port, cfg.Records, cfg.Saves, cfg.Seed = "127.0.0.1", freePort(t), filepath.Join(dir, "records"), filepath.Join(dir, "saves"), 5
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	srv, err := server.Start(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = srv.Shutdown(context.Background()) }() //nolint:errcheck // test

	var out bytes.Buffer
	args := []string{"-url", "ws://" + srv.Addr(), "-players", "2", "-start", "2", "-max-games", "2", "-max-response", "10s", "-seed", "7"}
	if err := run(ctx, args, &out); err != nil {
		t.Fatal(err)
	}
	text := out.String()
	if !strings.Contains(text, "result      PASS") {
		t.Fatalf("stage did not pass:\n%s", text)
	}
	if !strings.Contains(text, "2 finished, 0 failed") {
		t.Fatalf("games did not finish:\n%s", text)
	}
}

func freePort(t *testing.T) int {
	t.Helper()
	var lc net.ListenConfig
	ln, err := lc.Listen(context.Background(), "tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close() //nolint:errcheck // test
	_, port, err := net.SplitHostPort(ln.Addr().String())
	if err != nil {
		t.Fatal(err)
	}
	n, err := strconv.Atoi(port)
	if err != nil {
		t.Fatal(err)
	}
	return n
}
