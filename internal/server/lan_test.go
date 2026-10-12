package server

import (
	"context"
	"testing"
	"time"

	"github.com/FourSoulsRepo/four_souls/internal/lan"
	"github.com/FourSoulsRepo/four_souls/internal/protocol"
)

// startLAN starts a server on a free port; discovery as asked.
func startLAN(t *testing.T, discovery bool) (*Server, int) {
	t.Helper()
	cfg := DefaultConfig()
	cfg.Addr, cfg.Port, cfg.Records, cfg.Saves, cfg.Discovery = "127.0.0.1", freePort(t), t.TempDir(), t.TempDir(), discovery
	srv, err := Start(t.Context(), cfg)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = srv.Shutdown(context.Background()) }) //nolint:errcheck // test
	return srv, cfg.Port
}

func find(t *testing.T, port int) []lan.Found {
	t.Helper()
	found, err := lan.Find(t.Context(), port, 300*time.Millisecond)
	if err != nil {
		t.Fatal(err)
	}
	return found
}

// A hosting app answers while a game waits, and stops when it starts
// (N-13, ADR 009).
func TestLANFindsOnlyWaitingGames(t *testing.T) {
	srv, port := startLAN(t, true)
	if found := find(t, port); len(found) != 0 {
		t.Fatalf("found %+v with no game, want nothing", found)
	}
	url := "ws://" + srv.Addr() + "/ws"
	ann, bo := dialBot(t, url, "Ann", 1), dialBot(t, url, "Bo", 2)
	ann.send(protocol.TypeCreate, protocol.Create{Seats: 2})
	var tb protocol.Table
	ann.expect(protocol.TypeTable, &tb)

	found := find(t, port)
	if len(found) != 1 || len(found[0].Games) != 1 {
		t.Fatalf("found %+v, want Ann's game", found)
	}
	if g := found[0].Games[0]; g.Host != "Ann" || g.Seats != 2 || g.Taken != 1 {
		t.Errorf("game %+v, want Ann's with 1 of 2 seats", g)
	}
	if found[0].Protocol != protocol.Version {
		t.Errorf("protocol %d, want %d", found[0].Protocol, protocol.Version)
	}

	bo.send(protocol.TypeJoin, protocol.Join{Game: tb.Game})
	bo.expect(protocol.TypeTable, nil)
	if found := find(t, port); len(found) != 0 {
		t.Errorf("found %+v with every seat taken, want nothing", found)
	}
}

func TestDedicatedServerNeverAnswers(t *testing.T) {
	srv, port := startLAN(t, false)
	ann := dialBot(t, "ws://"+srv.Addr()+"/ws", "Ann", 1)
	ann.send(protocol.TypeCreate, protocol.Create{Seats: 2})
	ann.expect(protocol.TypeTable, nil)
	if found := find(t, port); len(found) != 0 {
		t.Errorf("a dedicated server answered: %+v", found)
	}
}
