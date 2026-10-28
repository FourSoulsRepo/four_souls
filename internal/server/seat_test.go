package server

import (
	"encoding/json"
	"net/http"
	"net/url"
	"testing"

	"github.com/FourSoulsRepo/four_souls/internal/protocol"
)

// getSeat asks srv's /seat about token.
func getSeat(t *testing.T, srv *Server, token string) protocol.SeatInfo {
	t.Helper()
	resp, err := http.Get("http://" + srv.Addr() + "/seat?token=" + url.QueryEscape(token)) //nolint:noctx // test
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close() //nolint:errcheck // test
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status %d", resp.StatusCode)
	}
	if resp.Header.Get("Access-Control-Allow-Origin") != "*" {
		t.Errorf("no CORS header: %v", resp.Header)
	}
	var info protocol.SeatInfo
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		t.Fatal(err)
	}
	return info
}

// /seat tells a restarted app whether its token still holds a seat
// (15.99): at a table, in a running game, never in a finished one.
func TestSeatTellsALiveSeat(t *testing.T) {
	srv, _ := startLAN(t, false)
	url := "ws://" + srv.Addr() + "/ws"
	ann, bo := dialBot(t, url, "Ann", 1), dialBot(t, url, "Bo", 2)
	if got := getSeat(t, srv, ann.token); got.Live {
		t.Errorf("a token in the lobby: %+v, want not live", got)
	}
	ann.send(protocol.TypeCreate, protocol.Create{Seats: 2})
	var tb protocol.Table
	ann.expect(protocol.TypeTable, &tb)
	if got := getSeat(t, srv, ann.token); got != (protocol.SeatInfo{Live: true, Game: tb.Game, Host: "Ann"}) {
		t.Errorf("at the table: %+v", got)
	}
	bo.send(protocol.TypeJoin, protocol.Join{Game: tb.Game})
	bo.expect(protocol.TypeTable, nil)
	for _, b := range []*bot{ann, bo} {
		b.send(protocol.TypeReady, protocol.Ready{Ready: true})
	}
	bo.expect(protocol.TypeUpdate, nil)
	if got := getSeat(t, srv, bo.token); got != (protocol.SeatInfo{Live: true, Game: tb.Game, Host: "Ann", Started: true}) {
		t.Errorf("in the game: %+v", got)
	}
	for _, token := range []string{"", "not-a-token", newToken()} {
		if got := getSeat(t, srv, token); got.Live {
			t.Errorf("token %q: %+v, want not live", token, got)
		}
	}
	srv.h.send(hubRequest{run: func(h *Hub) {
		for _, g := range h.games {
			g.room.over.Store(true)
		}
	}})
	if got := getSeat(t, srv, bo.token); got.Live {
		t.Errorf("after the end: %+v, want not live", got)
	}
}

// Random seats (UC-37): a shuffle keeps every player, the same seed
// gives the same order, and over many seeds each player gets seat 0.
func TestShuffleSeats(t *testing.T) {
	names := []string{"A", "B", "C", "D"}
	first := map[string]bool{}
	for seed := uint64(1); seed <= 200; seed++ {
		seats := make([]roomSeat, len(names))
		for i, n := range names {
			seats[i] = roomSeat{name: n}
		}
		again := make([]roomSeat, len(seats))
		copy(again, seats)
		shuffleSeats(seats, seed)
		shuffleSeats(again, seed)
		seen := map[string]bool{}
		for i, s := range seats {
			if again[i].name != s.name {
				t.Fatalf("seed %d: orders differ: %v and %v", seed, seats, again)
			}
			seen[s.name] = true
		}
		if len(seen) != len(names) {
			t.Fatalf("seed %d: order %v lost a player", seed, seats)
		}
		first[seats[0].name] = true
	}
	if len(first) != len(names) {
		t.Errorf("seat 0 went to %v; want every player", first)
	}
}
