package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/FourSoulsRepo/four_souls/internal/protocol"
)

// seatWait is how many hub seconds /seat waits for the lobby.
const seatWait = 2

// seatInfo tells whether token holds a seat here: at a table, or in a
// game that has not finished. Run on the hub goroutine.
func (h *Hub) seatInfo(token string) protocol.SeatInfo {
	if t, _ := h.tableSeat(token); t != nil {
		return protocol.SeatInfo{Live: true, Game: t.id, Host: t.host}
	}
	for _, g := range h.games {
		if g.room.hasToken(token) && !g.room.over.Load() {
			return protocol.SeatInfo{Live: true, Game: g.info.ID, Host: g.info.Host, Started: true}
		}
	}
	return protocol.SeatInfo{}
}

// seat answers GET /seat?token=… with a SeatInfo (15.99). A restarted
// app asks before it offers Rejoin; any page may ask, as the app's
// frontend is a page of its own origin, so the answer allows any.
func (h *Hub) seat(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	token := req.URL.Query().Get("token")
	info := protocol.SeatInfo{}
	if validToken(token) {
		reply := make(chan protocol.SeatInfo, 1)
		timeout := time.NewTimer(seatWait * h.second)
		defer timeout.Stop()
		select {
		case h.in <- hubRequest{run: func(h *Hub) { reply <- h.seatInfo(token) }}:
			select {
			case info = <-reply:
			case <-timeout.C:
				http.Error(w, "lobby not answering", http.StatusServiceUnavailable)
				return
			}
		case <-h.done:
			http.Error(w, "server stopping", http.StatusServiceUnavailable)
			return
		case <-timeout.C:
			http.Error(w, "lobby not answering", http.StatusServiceUnavailable)
			return
		case <-req.Context().Done():
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(info) //nolint:errcheck,errchkjson // the client may be gone
}
