package server

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/FourSoulsRepo/four_souls/internal/protocol"
)

// getStatus asks srv's /status.
func getStatus(t *testing.T, srv *Server) protocol.ServerStatus {
	t.Helper()
	resp, err := http.Get("http://" + srv.Addr() + "/status") //nolint:noctx // test
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
	var st protocol.ServerStatus
	if err := json.NewDecoder(resp.Body).Decode(&st); err != nil {
		t.Fatal(err)
	}
	return st
}

// /status counts the tables a player may still sit at (15.100): not a
// full one, not a started game.
func TestStatusCountsOpenGames(t *testing.T) {
	srv, _ := startLAN(t, false)
	url := "ws://" + srv.Addr() + "/ws"
	if got := getStatus(t, srv); got.Open != 0 {
		t.Errorf("an empty server: %+v", got)
	}
	ann, bo, cy := dialBot(t, url, "Ann", 1), dialBot(t, url, "Bo", 2), dialBot(t, url, "Cy", 3)
	var tb protocol.Table
	ann.send(protocol.TypeCreate, protocol.Create{Seats: 2})
	ann.expect(protocol.TypeTable, &tb)
	cy.send(protocol.TypeCreate, protocol.Create{Seats: 3, Password: "pw"})
	cy.expect(protocol.TypeTable, nil)
	if got := getStatus(t, srv); got.Open != 2 {
		t.Errorf("two tables, one locked: %+v, want 2 open", got)
	}
	bo.send(protocol.TypeJoin, protocol.Join{Game: tb.Game})
	bo.expect(protocol.TypeTable, nil)
	if got := getStatus(t, srv); got.Open != 1 {
		t.Errorf("one table full: %+v, want 1 open", got)
	}
	for _, b := range []*bot{ann, bo} {
		b.send(protocol.TypeReady, protocol.Ready{Ready: true})
	}
	bo.expect(protocol.TypeUpdate, nil)
	if got := getStatus(t, srv); got.Open != 1 {
		t.Errorf("the full one started: %+v, want 1 open", got)
	}
}
