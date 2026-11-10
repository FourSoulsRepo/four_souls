package main

import (
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	carddb "github.com/FourSoulsRepo/card_db"
	"github.com/FourSoulsRepo/four_souls/internal/instance"
	"github.com/FourSoulsRepo/four_souls/internal/settings"
)

// TestSettingsSurviveRestart: a new App reads what the old one saved.
func TestSettingsSurviveRestart(t *testing.T) {
	path := filepath.Join(t.TempDir(), "settings.json")
	if _, err := (&App{settings: path}).SaveSettings(settings.Settings{Nickname: "Ann"}); err != nil {
		t.Fatal(err)
	}
	if got := (&App{settings: path}).Settings().Nickname; got != "Ann" {
		t.Fatalf("nickname %q after restart, want Ann", got)
	}
}

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

// TestDeckRatios: every set together fits the deck ratios; the Base Game
// alone names the groups it is short of (GS-07).
func TestDeckRatios(t *testing.T) {
	a := &App{}
	if why := a.DeckRatios(a.CardSets(), 2); why != "" {
		t.Errorf("every set: %q", why)
	}
	why := a.DeckRatios([]string{"b2"}, 2)
	if !strings.Contains(why, `{"key":"refuse.ratio.butter_beans","args":{"have":3,"need":5}}`) {
		t.Errorf("Base Game: %q", why)
	}
}

// TestCardSetsSkipChallenges: card_db has the challenge set, but the host
// screen never offers it as a card set (ADR 028).
func TestCardSetsSkipChallenges(t *testing.T) {
	db, err := carddb.Load(os.DirFS("pkg/card_db"))
	if err != nil {
		t.Fatal(err)
	}
	sets := (&App{}).CardSets()
	found := false
	for _, s := range db.Sets() {
		if s.Code == "challenge" {
			found = true
			continue
		}
		if !slices.Contains(sets, s.Code) {
			t.Errorf("set %s is missing from %v", s.Code, sets)
		}
	}
	if !found {
		t.Fatal("card_db has no challenge set")
	}
	if slices.Contains(sets, "challenge") {
		t.Errorf("the challenge set is offered: %v", sets)
	}
}

// TestChallenges: the host screen learns the challenges from the app,
// like the card sets (ADR 028, Decision 10.7).
func TestChallenges(t *testing.T) {
	got := (&App{}).Challenges()
	if len(got) != 1 || got[0].ID != "krampus" || !slices.Equal(got[0].Difficulties, []string{"normal", "hard", "ultra"}) {
		t.Errorf("challenges %+v", got)
	}
}
