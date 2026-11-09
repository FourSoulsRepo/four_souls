package server

import (
	"bytes"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

// TestAppendFrame checks the frame header of each length class (RFC 6455,
// section 5.2): a short text frame, a 16-bit and a 64-bit length.
func TestAppendFrame(t *testing.T) {
	cases := []struct {
		n      int
		header []byte
	}{
		{0, []byte{0x81, 0}},
		{125, []byte{0x81, 125}},
		{126, []byte{0x81, 126, 0, 126}},
		{65535, []byte{0x81, 126, 0xff, 0xff}},
		{65536, []byte{0x81, 127, 0, 0, 0, 0, 0, 1, 0, 0}},
	}
	for _, c := range cases {
		msg := bytes.Repeat([]byte{'x'}, c.n)
		got := appendFrame(nil, msg)
		want := append(append([]byte{}, c.header...), msg...)
		if !bytes.Equal(got, want) {
			t.Errorf("length %d: got % x, want the header % x and the payload", c.n, got[:min(len(got), len(c.header))], c.header)
		}
	}
}

// TestIdleConnectionGetsPings checks that the server pings a connection
// that sends and gets nothing, so a proxy keeps it open (19.148).
func TestIdleConnectionGetsPings(t *testing.T) {
	old := pingEvery
	pingEvery = 20 * time.Millisecond
	t.Cleanup(func() { pingEvery = old })
	h := startHub(t)
	srv := httptest.NewServer(Handler(h))
	defer srv.Close()
	url := "ws" + strings.TrimPrefix(srv.URL, "http")
	ws, resp, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		t.Fatal(err)
	}
	_ = resp.Body.Close() //nolint:errcheck // test
	defer ws.Close()      //nolint:errcheck // test
	pings := make(chan struct{}, 16)
	ws.SetPingHandler(func(string) error {
		select {
		case pings <- struct{}{}:
		default:
		}
		return nil
	})
	// Reading runs the ping handler; the client sends nothing.
	go func() {
		for {
			if _, _, rerr := ws.ReadMessage(); rerr != nil {
				return
			}
		}
	}()
	deadline := time.After(5 * time.Second)
	for range 3 {
		select {
		case <-pings:
		case <-deadline:
			t.Fatal("an idle connection got fewer than 3 pings")
		}
	}
}
