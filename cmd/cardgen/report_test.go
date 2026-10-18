package main

import (
	"bytes"
	"maps"
	"os"
	"path/filepath"
	"testing"

	carddb "github.com/FourSoulsRepo/card_db"
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards"
)

func TestCardStatus(t *testing.T) {
	dir := t.TempDir()
	files := map[string]string{
		"stub.go":        "// " + StubMarker + " implement\n",
		"done.go":        "package x\n",
		"tested.go":      "package x\n",
		"tested_test.go": "package x\n",
	}
	for name, src := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(src), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	for id, want := range map[string]Status{"stub": NotStarted, "done": Implemented, "tested": Tested, "missing": NotStarted} {
		got, err := cardStatus(dir, carddb.Ref{ID: id, Version: 1})
		if err != nil || got != want {
			t.Errorf("%s: %v, %v; want %v", id, got, err, want)
		}
	}
}

// TestReportIsCurrent fails when a card changed but the report was not
// regenerated: go run ./cmd/cardgen -report
func TestReportIsCurrent(t *testing.T) {
	db, err := carddb.Load(os.DirFS("../../pkg/card_db"))
	if err != nil {
		t.Fatal(err)
	}
	want, err := report(db, "../../pkg/rules_engine/cards")
	if err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(filepath.Join("../..", reportFile))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("%s is out of date; run: go run ./cmd/cardgen -report", reportFile)
	}
}

// TestThreePlusMatchesCardDB: the engine leaves out the cards card_db
// marks with the 3+ player symbol (R-CARD-32), and no others.
func TestThreePlusMatchesCardDB(t *testing.T) {
	db, err := carddb.Load(os.DirFS("../../pkg/card_db"))
	if err != nil {
		t.Fatal(err)
	}
	for _, set := range cards.Sets() {
		for _, d := range set.Cards {
			ref, err := carddb.ParseRef(string(d.Ref))
			if err != nil {
				t.Fatal(err)
			}
			c, ok := db.Get(ref)
			if !ok {
				continue
			}
			if c.ThreePlus != d.ThreePlus {
				t.Errorf("%s: card_db three_plus %v, engine ThreePlus %v", d.Ref, c.ThreePlus, d.ThreePlus)
			}
		}
	}
}

// soulCount sums the soul icons in a card's reward box.
func soulCount(c carddb.Card) int {
	n := 0
	for _, r := range c.Rewards {
		if k, ok := rewardCount(r.Count); ok && r.Kind == "Soul" {
			n += k
		}
	}
	return n
}

// TestBossesHaveSouls: every boss gives a soul. foursouls.com left the
// soul off Famine; overrides/b2.json puts it back.
func TestBossesHaveSouls(t *testing.T) {
	db, err := carddb.Load(os.DirFS("../../pkg/card_db"))
	if err != nil {
		t.Fatal(err)
	}
	for _, ref := range db.Refs() {
		c, _ := db.Get(ref)
		if c.Type == "Boss Card" && soulCount(c) == 0 {
			t.Errorf("%s: boss without a soul", ref)
		}
	}
}

// TestSoulsMatchCardDB: the engine's Soul matches the card_db reward box.
func TestSoulsMatchCardDB(t *testing.T) {
	db, err := carddb.Load(os.DirFS("../../pkg/card_db"))
	if err != nil {
		t.Fatal(err)
	}
	for _, set := range cards.Sets() {
		for _, d := range set.Cards {
			ref, err := carddb.ParseRef(string(d.Ref))
			if err != nil {
				t.Fatal(err)
			}
			c, ok := db.Get(ref)
			if !ok {
				continue
			}
			if want := soulCount(c); d.Soul != want {
				t.Errorf("%s: card_db souls %d, engine Soul %d", d.Ref, want, d.Soul)
			}
		}
	}
}

// rewardSums totals the printed cents, loot and treasure of a reward box.
// Rolled counts and other icons are left out.
func rewardSums(c carddb.Card) map[string]int {
	out := map[string]int{}
	for _, r := range c.Rewards {
		if k, ok := rewardCount(r.Count); ok && rewardKinds[r.Kind] != "" {
			out[r.Kind] += k
		}
	}
	return out
}

// TestRewardsMatchCardDB: a monster's fixed rewards match the card_db
// reward box. foursouls.com added extra icons to Ring of Flies and
// Mom's Dead Hand; overrides/b2.json fixes them.
func TestRewardsMatchCardDB(t *testing.T) {
	db, err := carddb.Load(os.DirFS("../../pkg/card_db"))
	if err != nil {
		t.Fatal(err)
	}
	names := map[engine.RewardKind]string{engine.RewardCents: "Coin", engine.RewardLoot: "Loot", engine.RewardTreasure: "Treasure"}
	for _, set := range cards.Sets() {
		for _, d := range set.Cards {
			ref, err := carddb.ParseRef(string(d.Ref))
			if err != nil {
				t.Fatal(err)
			}
			c, ok := db.Get(ref)
			if !ok || d.Kind != engine.MonsterCard {
				continue
			}
			got := map[string]int{}
			for _, r := range d.Rewards {
				if !r.Roll {
					got[names[r.Kind]] += r.Amount
				}
			}
			if want := rewardSums(c); !maps.Equal(got, want) {
				t.Errorf("%s: card_db rewards %v, engine %v", d.Ref, want, got)
			}
		}
	}
}
