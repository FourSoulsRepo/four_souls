package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

// TestHealthz: 200 while the lobby answers, 503 when it is stuck or
// stopped (8.19).
func TestHealthz(t *testing.T) {
	h := NewHub()
	h.second = time.Millisecond
	ctx, cancel := context.WithCancel(context.Background())
	go h.Run(ctx)
	defer cancel()
	srv := httptest.NewServer(http.HandlerFunc(h.healthz))
	defer srv.Close()
	get := func() int {
		resp, err := http.Get(srv.URL) //nolint:noctx // test
		if err != nil {
			t.Fatal(err)
		}
		_ = resp.Body.Close() //nolint:errcheck // test
		return resp.StatusCode
	}
	if code := get(); code != http.StatusOK {
		t.Fatalf("a running lobby: %d", code)
	}
	release := make(chan struct{})
	h.send(hubRequest{run: func(*Hub) { <-release }})
	if code := get(); code != http.StatusServiceUnavailable {
		t.Errorf("a stuck lobby: %d", code)
	}
	close(release)
	cancel()
	<-h.done
	if code := get(); code != http.StatusServiceUnavailable {
		t.Errorf("a stopped lobby: %d", code)
	}
}
