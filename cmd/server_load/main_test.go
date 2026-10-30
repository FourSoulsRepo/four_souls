package main

import (
	"bytes"
	"context"
	"io"
	"net"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/FourSoulsRepo/four_souls/internal/server"
)

func TestSpreadOf(t *testing.T) {
	got := spreadOf([]time.Duration{5, 1, 9, 3, 7})
	want := spread[time.Duration]{n: 5, min: 1, median: 5, p95: 9, p99: 9, max: 9}
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
	if p := spreadOf(hundred).p99; p != 99 {
		t.Fatalf("p99 of 1..100 = %d, want 99", p)
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

// TestCutGamesAreNotFailures: a game the stage stopped passes the stage.
func TestCutGamesAreNotFailures(t *testing.T) {
	st := stage{
		games:    2,
		results:  []gameResult{{finished: true}, {cut: true}},
		col:      newCollector(),
		finished: 1,
		cut:      1,
	}
	if problems := st.problems(config{maxResponse: time.Second}); len(problems) != 0 {
		t.Fatalf("problems %v, want none", problems)
	}
}

// TestStopAtAgainstServer runs a stage with -stop-at on a local server; the
// stage passes whether or not the other game got cut.
func TestStopAtAgainstServer(t *testing.T) {
	if testing.Short() {
		t.Skip("plays whole games")
	}
	dir := t.TempDir()
	cfg := server.DefaultConfig()
	cfg.Addr, cfg.Port, cfg.Records, cfg.Saves, cfg.Seed = "127.0.0.1", freePort(t), filepath.Join(dir, "records"), filepath.Join(dir, "saves"), 9
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	srv, err := server.Start(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = srv.Shutdown(context.Background()) }() //nolint:errcheck // test

	var out bytes.Buffer
	args := []string{"-url", "ws://" + srv.Addr(), "-players", "2", "-start", "2", "-max-games", "2", "-max-response", "10s", "-seed", "7", "-stop-at", "30"}
	if err := run(ctx, args, &out); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "result      PASS") {
		t.Fatalf("stage did not pass:\n%s", out.String())
	}
}

// TestReconnectAfterDrop: a bot whose connection is cut in the middle of a
// game takes its seat back with its token, the game finishes, and the report
// counts the reconnect (N-08).
func TestReconnectAfterDrop(t *testing.T) {
	if testing.Short() {
		t.Skip("plays whole games")
	}
	dir := t.TempDir()
	cfg := server.DefaultConfig()
	cfg.Addr, cfg.Port, cfg.Records, cfg.Saves, cfg.Seed = "127.0.0.1", freePort(t), filepath.Join(dir, "records"), filepath.Join(dir, "saves"), 11
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	srv, err := server.Start(ctx, cfg)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = srv.Shutdown(context.Background()) }() //nolint:errcheck // test

	var out bytes.Buffer
	args := []string{"-url", cutProxy(t, srv.Addr()), "-players", "2", "-start", "1", "-max-games", "1", "-max-response", "10s", "-seed", "7"}
	if err := run(ctx, args, &out); err != nil {
		t.Fatal(err)
	}
	text := out.String()
	for _, want := range []string{"1 finished, 0 failed", "reconnects  1, 0 attempts failed", "result      PASS"} {
		if !strings.Contains(text, want) {
			t.Fatalf("report lacks %q:\n%s", want, text)
		}
	}
}

// cutProxy forwards WebSocket connections to target. The second connection
// it accepts is cut once, after cutAfter bytes from the server, so the bot on
// it must reconnect. It returns the URL to dial.
func cutProxy(t *testing.T, target string) string {
	t.Helper()
	const cutAfter = 64 << 10
	var lc net.ListenConfig
	ln, err := lc.Listen(context.Background(), "tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = ln.Close() }) //nolint:errcheck // test
	go func() {
		for n := 1; ; n++ {
			client, err := ln.Accept()
			if err != nil {
				return
			}
			go forward(client, target, n == 2, cutAfter)
		}
	}()
	return "ws://" + ln.Addr().String() + "/ws"
}

// forward copies one connection both ways; with cut set, it closes both sides
// after cutAfter bytes from the server.
func forward(client net.Conn, target string, cut bool, cutAfter int64) {
	var d net.Dialer
	server, err := d.DialContext(context.Background(), "tcp", target)
	if err != nil {
		_ = client.Close() //nolint:errcheck // test
		return
	}
	go func() {
		_, _ = io.Copy(server, client) //nolint:errcheck // the connection ends either way
		_ = server.Close()             //nolint:errcheck // test
	}()
	if cut {
		_, _ = io.CopyN(client, server, cutAfter) //nolint:errcheck // the cut is the point
	} else {
		_, _ = io.Copy(client, server) //nolint:errcheck // the connection ends either way
	}
	_ = client.Close() //nolint:errcheck // test
	_ = server.Close() //nolint:errcheck // test
}

// TestRetryDelay: the app's backoff (8.24).
func TestRetryDelay(t *testing.T) {
	want := []time.Duration{0, time.Second, 2 * time.Second, 4 * time.Second, 8 * time.Second, 10 * time.Second, 10 * time.Second, 10 * time.Second}
	for n, w := range want {
		if got := retryDelay(n); got != w {
			t.Errorf("retryDelay(%d) = %v, want %v", n, got, w)
		}
	}
	if got := retryDelay(1 << 20); got != 10*time.Second {
		t.Errorf("retryDelay(large) = %v, want 10s", got)
	}
}

// TestStopAfter: stop-at percent of the games, rounded up, at least one.
func TestStopAfter(t *testing.T) {
	cases := []struct{ games, stopAt, want int }{
		{games: 2, stopAt: 30, want: 1},
		{games: 16, stopAt: 30, want: 5},
		{games: 64, stopAt: 30, want: 20},
		{games: 64, stopAt: 0, want: 0},
		{games: 4, stopAt: 100, want: 4},
	}
	for _, c := range cases {
		if got := stopAfter(c.games, c.stopAt); got != c.want {
			t.Errorf("stopAfter(%d, %d) = %d, want %d", c.games, c.stopAt, got, c.want)
		}
	}
}

// TestNextStage: the stages double and the last one is -max-games.
func TestNextStage(t *testing.T) {
	var got []int
	for games := 2; games <= 200; games = nextStage(games, 200) {
		got = append(got, games)
	}
	want := []int{2, 4, 8, 16, 32, 64, 128, 200}
	if !slices.Equal(got, want) {
		t.Errorf("stages %v, want %v", got, want)
	}
}
