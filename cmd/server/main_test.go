package main

import (
	"bytes"
	"context"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"

	"github.com/FourSoulsRepo/four_souls/internal/protocol"
)

func TestParse(t *testing.T) {
	cfg, err := parse(nil)
	if err != nil || cfg.Port != protocol.DefaultPort || cfg.Retention != 30 || !cfg.Presence {
		t.Fatalf("defaults %+v, %v", cfg, err)
	}
	file := filepath.Join(t.TempDir(), "server.json")
	if werr := os.WriteFile(file, []byte(`{"port": 5000, "retention": 7}`), 0o600); werr != nil {
		t.Fatal(werr)
	}
	cfg, err = parse([]string{"-config", file, "-retention", "0"})
	if err != nil || cfg.Port != 5000 || cfg.Retention != 0 {
		t.Fatalf("file and flag %+v, %v", cfg, err)
	}
	if werr := os.WriteFile(file, []byte(`{"presence": false}`), 0o600); werr != nil {
		t.Fatal(werr)
	}
	if cfg, err = parse([]string{"-config", file}); err != nil || cfg.Presence {
		t.Fatalf("presence off in the file: %+v, %v", cfg, err)
	}
	if cfg, err = parse([]string{"-config", file, "-presence"}); err != nil || !cfg.Presence {
		t.Fatalf("presence on by flag over the file: %+v, %v", cfg, err)
	}
	if werr := os.WriteFile(file, []byte(`{"name": "Basement"}`), 0o600); werr != nil {
		t.Fatal(werr)
	}
	if cfg, err = parse([]string{"-config", file}); err != nil || cfg.Name != "Basement" {
		t.Fatalf("name in the file: %+v, %v", cfg, err)
	}
	if cfg, err = parse([]string{"-config", file, "-name", "Cathedral"}); err != nil || cfg.Name != "Cathedral" {
		t.Fatalf("name by flag over the file: %+v, %v", cfg, err)
	}
	for _, bad := range []string{" ", "a\tb", strings.Repeat("x", protocol.MaxServerName+1)} {
		if _, err := parse([]string{"-name", bad}); err == nil {
			t.Errorf("name %q accepted", bad)
		}
	}
	if _, err := parse([]string{"-port", "0"}); err == nil {
		t.Error("port 0 accepted")
	}
	// Idle limits (N-20): on by default, 0 turns one off, never negative.
	if lim, lerr := parse(nil); lerr != nil || lim.EmptyMinutes != 30 || lim.GameHours != 24 {
		t.Fatalf("limit defaults %+v, %v", lim, lerr)
	}
	if lim, lerr := parse([]string{"-empty_minutes", "0", "-game_hours", "48"}); lerr != nil || lim.EmptyMinutes != 0 || lim.GameHours != 48 {
		t.Fatalf("limit flags %+v, %v", lim, lerr)
	}
	if _, lerr := parse([]string{"-game_hours", "-1"}); lerr == nil {
		t.Error("a negative limit accepted")
	}
	// Bots (ADR 033, Decision 17): every level by default; off, standard
	// or all; nothing else.
	if b, berr := parse(nil); berr != nil || b.Bots != "all" {
		t.Fatalf("bots default %+v, %v", b, berr)
	}
	for _, mode := range []string{"off", "standard", "all"} {
		if b, berr := parse([]string{"-bots", mode}); berr != nil || b.Bots != mode {
			t.Errorf("-bots %s: %+v, %v", mode, b, berr)
		}
	}
	if _, berr := parse([]string{"-bots", "hard"}); berr == nil {
		t.Error("-bots hard accepted")
	}
	if werr := os.WriteFile(file, []byte(`{"bots": "off"}`), 0o600); werr != nil {
		t.Fatal(werr)
	}
	if b, berr := parse([]string{"-config", file}); berr != nil || b.Bots != "off" {
		t.Fatalf("bots off in the file: %+v, %v", b, berr)
	}
	// The hard bot's think time (Decision 17.2.1): 1.5 s by default; a
	// duration from 100 ms to 10 s.
	if b, berr := parse(nil); berr != nil || b.BotThink != "1.5s" {
		t.Fatalf("bot_think default %+v, %v", b, berr)
	}
	if b, berr := parse([]string{"-bot_think", "800ms"}); berr != nil || b.BotThink != "800ms" {
		t.Errorf("-bot_think 800ms: %+v, %v", b, berr)
	}
	for _, bad := range []string{"50ms", "11s", "fast", "0"} {
		if _, berr := parse([]string{"-bot_think", bad}); berr == nil {
			t.Errorf("-bot_think %s accepted", bad)
		}
	}
	if werr := os.WriteFile(file, []byte(`{"bot_think": "3s"}`), 0o600); werr != nil {
		t.Fatal(werr)
	}
	if b, berr := parse([]string{"-config", file}); berr != nil || b.BotThink != "3s" {
		t.Fatalf("bot_think in the file: %+v, %v", b, berr)
	}
}

// TestParseTestSeeds: test seeds are off unless the flag or ENV turns
// them on; a config file cannot (ADR 032).
func TestParseTestSeeds(t *testing.T) {
	if cfg, err := parse(nil); err != nil || cfg.TestSeeds {
		t.Fatalf("default %+v, %v", cfg, err)
	}
	if cfg, err := parse([]string{"-test_seeds"}); err != nil || !cfg.TestSeeds {
		t.Fatalf("flag %+v, %v", cfg, err)
	}
	file := filepath.Join(t.TempDir(), "server.json")
	if err := os.WriteFile(file, []byte(`{"test_seeds": true, "TestSeeds": true}`), 0o600); err != nil {
		t.Fatal(err)
	}
	if cfg, err := parse([]string{"-config", file}); err != nil || cfg.TestSeeds {
		t.Fatalf("file %+v, %v", cfg, err)
	}
	t.Setenv("FOUR_SOULS_TEST_SEEDS", "true")
	if cfg, err := parse(nil); err != nil || !cfg.TestSeeds {
		t.Fatalf("env %+v, %v", cfg, err)
	}
}

// TestParseEnv: every flag has a FOUR_SOULS_ variable; ENV overrides the
// file and flags override ENV (N-19).
// The tracing flags read like the others, and ENV sets them too (ADR 024).
func TestParseTracing(t *testing.T) {
	if _, tr, err := parseAll(nil); err != nil || tr.Dir != "" || tr.Addr != "" {
		t.Fatalf("defaults %+v, %v", tr, err)
	}
	if _, tr, err := parseAll([]string{"-trace_dir", "traces", "-pprof_addr", "127.0.0.1:6060"}); err != nil || tr.Dir != "traces" || tr.Addr != "127.0.0.1:6060" {
		t.Fatalf("flags %+v, %v", tr, err)
	}
	t.Setenv("FOUR_SOULS_TRACE_DIR", "env-traces")
	if _, tr, err := parseAll(nil); err != nil || tr.Dir != "env-traces" {
		t.Fatalf("env %+v, %v", tr, err)
	}
}

func TestParseEnv(t *testing.T) {
	dir := t.TempDir()
	file := filepath.Join(dir, "server.json")
	if err := os.WriteFile(file, []byte(`{"port": 5000, "retention": 7, "name": "File"}`), 0o600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("FOUR_SOULS_CONFIG", file)
	t.Setenv("FOUR_SOULS_ADDR", "127.0.0.1")
	t.Setenv("FOUR_SOULS_PORT", "5001")
	t.Setenv("FOUR_SOULS_RECORDS", "/data/records")
	t.Setenv("FOUR_SOULS_SAVES", "/data/saves")
	t.Setenv("FOUR_SOULS_PRESENCE", "false")
	t.Setenv("FOUR_SOULS_NAME", "Basement")
	t.Setenv("FOUR_SOULS_RETENTION", "")
	t.Setenv("FOUR_SOULS_EMPTY_MINUTES", "45")
	t.Setenv("FOUR_SOULS_GAME_HOURS", "0")
	t.Setenv("FOUR_SOULS_BOTS", "standard")
	t.Setenv("FOUR_SOULS_BOT_THINK", "2s")
	cfg, err := parse(nil)
	if err != nil {
		t.Fatal(err)
	}
	want := cfg
	want.Addr, want.Port, want.Records, want.Saves, want.Presence, want.Name = "127.0.0.1", 5001, "/data/records", "/data/saves", false, "Basement"
	want.Retention = 7 // an empty variable is unset: the file's value stays
	want.EmptyMinutes, want.GameHours, want.Bots, want.BotThink = 45, 0, "standard", "2s"
	if cfg != want {
		t.Fatalf("env over the file: got %+v, want %+v", cfg, want)
	}
	if cfg, err = parse([]string{"-port", "5002", "-presence", "-name", "Cathedral"}); err != nil || cfg.Port != 5002 || !cfg.Presence || cfg.Name != "Cathedral" {
		t.Fatalf("flags over env: %+v, %v", cfg, err)
	}
	for key, bad := range map[string]string{"FOUR_SOULS_PORT": "abc", "FOUR_SOULS_PRESENCE": "maybe", "FOUR_SOULS_RETENTION": "-x", "FOUR_SOULS_NAME": "a\tb"} {
		t.Run(key, func(t *testing.T) {
			t.Setenv(key, bad)
			if _, err := parse(nil); err == nil {
				t.Errorf("%s=%q accepted", key, bad)
			} else if key != "FOUR_SOULS_NAME" && !strings.Contains(err.Error(), key) {
				t.Errorf("error %q does not name %s", err, key)
			}
		})
	}
	t.Setenv("FOUR_SOULS_CONFIG", filepath.Join(dir, "missing.json"))
	if _, err := parse(nil); err == nil {
		t.Error("a missing config file from env accepted")
	}
}

// syncBuffer is a bytes.Buffer safe for the server goroutine and the test.
type syncBuffer struct {
	mu sync.Mutex
	b  bytes.Buffer
}

func (s *syncBuffer) Write(p []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.b.Write(p)
}

func (s *syncBuffer) String() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.b.String()
}

// TestRun starts the server, connects a client, and stops it cleanly.
// The welcome carries the server's name (N-18). The signal that stops
// it lets the client hear of the restart first (N-20).
func TestRun(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	var out syncBuffer
	done := make(chan error, 1)
	dir := t.TempDir()
	args := []string{"-addr", "127.0.0.1", "-port", "47741", "-name", "Basement", "-records", filepath.Join(dir, "records"), "-saves", filepath.Join(dir, "saves")}
	go func() { done <- run(ctx, args, &out) }()

	var ws *websocket.Conn
	for range 50 {
		c, resp, err := websocket.DefaultDialer.Dial("ws://127.0.0.1:47741/ws", nil)
		if err == nil {
			_ = resp.Body.Close() //nolint:errcheck // test
			ws = c
			break
		}
		time.Sleep(50 * time.Millisecond)
	}
	if ws == nil {
		t.Fatal("could not connect")
	}
	// Docker's healthcheck asks /healthz (8.19).
	if resp, err := http.Get("http://127.0.0.1:47741/healthz"); err != nil || resp.StatusCode != http.StatusOK { //nolint:noctx // test
		t.Fatalf("healthz: %v, %v", resp, err)
	} else {
		_ = resp.Body.Close() //nolint:errcheck // test
	}
	msg, err := protocol.Encode(protocol.TypeHello, 1, protocol.Hello{Protocol: protocol.Version, Name: "Ann"})
	if err != nil {
		t.Fatal(err)
	}
	if err := ws.WriteMessage(websocket.TextMessage, msg); err != nil {
		t.Fatal(err)
	}
	if _, got, err := ws.ReadMessage(); err != nil || !strings.Contains(string(got), `"welcome"`) || !strings.Contains(string(got), `"name":"Basement"`) {
		t.Fatalf("got %s, %v", got, err)
	}
	cancel()
	select {
	case err := <-done:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(15 * time.Second):
		t.Fatal("the server did not stop")
	}
	// Messages sent before the shutdown may still be buffered; then the
	// connection must end.
	if err := ws.SetReadDeadline(time.Now().Add(5 * time.Second)); err != nil {
		t.Fatal(err)
	}
	restarting := false
	for {
		_, got, err := ws.ReadMessage()
		if err != nil {
			if strings.Contains(err.Error(), "timeout") {
				t.Error("the connection is still open after shutdown")
			}
			break
		}
		restarting = restarting || strings.Contains(string(got), `"type":"restarting"`)
	}
	if !restarting {
		t.Error("no restarting message before the connection closed")
	}
	for _, want := range []string{"Four Souls dedicated server", "listening on 127.0.0.1:47741", "stopped"} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("output lacks %q:\n%s", want, out.String())
		}
	}
}
