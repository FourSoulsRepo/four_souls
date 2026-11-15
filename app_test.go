package main

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/FourSoulsRepo/four_souls/internal/instance"
	"github.com/FourSoulsRepo/four_souls/internal/settings"
)

// TestSecondCopySettings: a second copy keeps the first copy's choices
// but has its own nickname and seat.
func TestSecondCopySettings(t *testing.T) {
	dir := t.TempDir()
	first := filepath.Join(dir, "settings.json")
	if _, err := settings.Save(first, settings.Settings{
		Nickname: "Ann", Mat: "cathedral",
		Rejoin: settings.Rejoin{Address: "localhost", Token: "t1", Game: "g1"},
	}); err != nil {
		t.Fatal(err)
	}
	a, err := instance.Acquire(dir)
	if err != nil {
		t.Fatal(err)
	}
	b, err := instance.Acquire(dir)
	if err != nil {
		t.Fatal(err)
	}
	if got := slotSettings(a, first); got != first {
		t.Errorf("first copy uses %q, want %q", got, first)
	}
	path := slotSettings(b, first)
	if path == first {
		t.Fatal("second copy shares the first copy's settings")
	}
	s := settings.Load(path)
	if s.Nickname != "" || s.Rejoin.Token != "" || s.Mat != "cathedral" {
		t.Errorf("second copy settings %+v; want the Cathedral mat, no nickname, no seat", s)
	}
	if settings.Load(first).Rejoin.Token != "t1" {
		t.Error("the first copy lost its seat")
	}
}

// The Wails context is a background one before startup; it ends at
// shutdown, with or without a startup (ADR 030, item 4).
func TestWailsContext(t *testing.T) {
	var a App
	if a.Context().Done() != nil {
		t.Error("ends before startup")
	}
	type key struct{}
	a.startup(context.WithValue(context.Background(), key{}, "wails"))
	if a.Context().Value(key{}) != "wails" {
		t.Error("lost the startup context")
	}
	a.shutdown()
	if a.Context().Err() == nil {
		t.Error("alive after shutdown")
	}
	var unstarted App
	unstarted.shutdown()
	if unstarted.Context().Err() == nil {
		t.Error("alive after shutdown without startup")
	}
}
