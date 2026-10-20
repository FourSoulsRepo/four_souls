package server

import (
	"net/http"
	"time"
)

// healthWait is how many hub seconds /healthz waits for the lobby.
const healthWait = 2

// healthz answers GET /healthz for Docker's healthcheck (8.19): 200
// while the lobby goroutine answers, 503 when it is stuck or stopped.
// It tells nothing about games, so anyone may ask.
func (h *Hub) healthz(w http.ResponseWriter, req *http.Request) {
	answered := make(chan struct{})
	timeout := time.NewTimer(healthWait * h.second)
	defer timeout.Stop()
	ok := false
	select {
	case h.in <- hubRequest{run: func(*Hub) { close(answered) }}:
		select {
		case <-answered:
			ok = true
		case <-timeout.C:
		case <-req.Context().Done():
		}
	case <-h.done:
	case <-timeout.C:
	case <-req.Context().Done():
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	if !ok {
		http.Error(w, "lobby not answering", http.StatusServiceUnavailable)
		return
	}
	_, _ = w.Write([]byte("ok\n")) //nolint:errcheck // the client may be gone
}
