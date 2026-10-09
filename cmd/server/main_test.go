package main

import (
	"bytes"
	"context"
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
	if err != nil || cfg.Port != protocol.DefaultPort || cfg.Retention != 30 {
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
	if _, err := parse([]string{"-port", "0"}); err == nil {
		t.Error("port 0 accepted")
	}
	if _, err := parse([]string{"-players", "5"}); err == nil {
		t.Error("5 players accepted")
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
func TestRun(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	var out syncBuffer
	done := make(chan error, 1)
	go func() { done <- run(ctx, []string{"-addr", "127.0.0.1", "-port", "47741"}, &out) }()

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
	msg, err := protocol.Encode(protocol.TypeHello, 1, protocol.Hello{Protocol: protocol.Version, Name: "Ann"})
	if err != nil {
		t.Fatal(err)
	}
	if err := ws.WriteMessage(websocket.TextMessage, msg); err != nil {
		t.Fatal(err)
	}
	if _, got, err := ws.ReadMessage(); err != nil || !strings.Contains(string(got), `"welcome"`) {
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
	for {
		if _, _, err := ws.ReadMessage(); err != nil {
			if strings.Contains(err.Error(), "timeout") {
				t.Error("the connection is still open after shutdown")
			}
			break
		}
	}
	for _, want := range []string{"Four Souls dedicated server", "listening on 127.0.0.1:47741", "stopped"} {
		if !strings.Contains(out.String(), want) {
			t.Errorf("output lacks %q:\n%s", want, out.String())
		}
	}
}
