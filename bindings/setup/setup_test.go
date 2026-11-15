package setup

import (
	"os"
	"slices"
	"strings"
	"testing"

	carddb "github.com/FourSoulsRepo/card_db"
)

// TestDeckRatios: every set together fits the deck ratios; the Base Game
// alone names the groups it is short of (GS-07).
func TestDeckRatios(t *testing.T) {
	s := New()
	if why := s.DeckRatios(s.CardSets(), 2); why != "" {
		t.Errorf("every set: %q", why)
	}
	why := s.DeckRatios([]string{"b2"}, 2)
	if !strings.Contains(why, `{"key":"refuse.ratio.butter_beans","args":{"have":3,"need":5}}`) {
		t.Errorf("Base Game: %q", why)
	}
}

// TestCardSetsSkipChallenges: card_db has the challenge set, but the host
// screen never offers it as a card set (ADR 028).
func TestCardSetsSkipChallenges(t *testing.T) {
	db, err := carddb.Load(os.DirFS("../../pkg/card_db"))
	if err != nil {
		t.Fatal(err)
	}
	sets := New().CardSets()
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
	got := New().Challenges()
	if len(got) != 2 || got[0].ID != "krampus" || got[1].ID != "trick_or_treat" {
		t.Fatalf("challenges %+v", got)
	}
	for _, c := range got {
		if !slices.Equal(c.Difficulties, []string{"normal", "hard", "ultra"}) {
			t.Errorf("%s: difficulties %v", c.ID, c.Difficulties)
		}
	}
}
