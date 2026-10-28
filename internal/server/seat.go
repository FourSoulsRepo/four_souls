package server

import (
	"encoding/json"
	"math/rand/v2"
	"net/http"
	"time"

	"github.com/FourSoulsRepo/four_souls/internal/protocol"
)

// seatWait is how many hub seconds /seat waits for the lobby.
const seatWait = 2

// shuffleSeats puts a game's seats in a random order (Random seats),
// drawn from the game's seed.
func shuffleSeats(seats []roomSeat, seed uint64) {
	rng := rand.New(rand.NewPCG(seed, seed^0x5ea7)) //nolint:gosec // seating, not security
	rng.Shuffle(len(seats), func(i, j int) { seats[i], seats[j] = seats[j], seats[i] })
}

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
		var ok bool
		if info, ok = askHub(h, w, req, seatWait, func(h *Hub) protocol.SeatInfo { return h.seatInfo(token) }); !ok {
			return
		}
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(info) //nolint:errcheck,errchkjson // the client may be gone
}

// askHub runs run on the hub goroutine and returns its answer. It waits
// at most wait hub seconds; else it answers 503 itself and returns
// false, as it does when the server stops or the client is gone.
func askHub[T any](h *Hub, w http.ResponseWriter, req *http.Request, wait int, run func(*Hub) T) (T, bool) {
	var none T
	reply := make(chan T, 1)
	timeout := time.NewTimer(time.Duration(wait) * h.second)
	defer timeout.Stop()
	select {
	case h.in <- hubRequest{run: func(h *Hub) { reply <- run(h) }}:
		select {
		case v := <-reply:
			return v, true
		case <-timeout.C:
		}
	case <-h.done:
		http.Error(w, "server stopping", http.StatusServiceUnavailable)
		return none, false
	case <-timeout.C:
	case <-req.Context().Done():
		return none, false
	}
	http.Error(w, "lobby not answering", http.StatusServiceUnavailable)
	return none, false
}
