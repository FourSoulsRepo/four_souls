package main

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/FourSoulsRepo/four_souls/internal/protocol"
)

// CheckSeat tries each base in order and asks /seat with the token
// (15.99); anything but an http(s) server address is refused.
func TestCheckSeat(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/seat" {
			http.NotFound(w, r)
			return
		}
		_ = json.NewEncoder(w).Encode(protocol.SeatInfo{Live: r.URL.Query().Get("token") == "t0k3n", Host: "Ann"}) //nolint:errcheck // test
	}))
	defer srv.Close()
	a := &App{}
	// A base that does not answer is skipped, as wss before ws.
	got, err := a.CheckSeat([]string{"http://127.0.0.1:1", srv.URL}, "t0k3n")
	if err != nil || !got.Live || got.Host != "Ann" {
		t.Fatalf("got %+v, %v; want Ann's live seat", got, err)
	}
	if got, err := a.CheckSeat([]string{srv.URL + "/"}, "other"); err != nil || got.Live {
		t.Errorf("another token: %+v, %v", got, err)
	}
	if _, err := a.CheckSeat([]string{"http://127.0.0.1:1"}, "t0k3n"); err == nil {
		t.Error("no server answered, want an error")
	}
	for _, bad := range [][]string{nil, {"file:///etc/passwd"}, {"http://"}, {srv.URL + "/record"}, {srv.URL + "?x=1"}} {
		if _, err := a.CheckSeat(bad, "t0k3n"); !errors.Is(err, ErrServerURL) {
			t.Errorf("%q: err %v, want ErrServerURL", bad, err)
		}
	}
}
