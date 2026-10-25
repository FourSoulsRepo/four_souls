package server

import (
	"encoding/json"
	"net/http"

	"github.com/FourSoulsRepo/four_souls/internal/protocol"
)

// statusWait is how many hub seconds /status waits for the lobby.
const statusWait = 2

// serverStatus counts the tables with a free seat; a started game is
// no table. Run on the hub goroutine.
func (h *Hub) serverStatus() protocol.ServerStatus {
	st := protocol.ServerStatus{}
	for _, t := range h.tables {
		for _, s := range t.seats {
			if s.token == "" {
				st.Open++
				break
			}
		}
	}
	return st
}

// status answers GET /status with a ServerStatus (15.100). The Online
// screen asks every Recent server when it opens. It tells only what the
// lobby lists to anyone who connects, so any page may ask.
func (h *Hub) status(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	st, ok := askHub(h, w, req, statusWait, (*Hub).serverStatus)
	if !ok {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(st) //nolint:errcheck,errchkjson // the client may be gone
}
