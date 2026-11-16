package bot

import (
	"math"
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards"
)

// TestValuesFile: the embedded table parses, and every loot card and
// item of every set has a value: its measured one, or its kind's mean
// (Decision 20.8.1, 20.9). A card a deck holds is in the table; a card
// that starts outside the game (a starting item) and a challenge's card
// fall back to the mean.
func TestValuesFile(t *testing.T) {
	f, err := ParseValues(valuesJSON)
	if err != nil {
		t.Fatal(err)
	}
	if f.Note == "" || f.Samples == 0 || len(f.Cards) == 0 {
		t.Fatalf("the table is empty: run make bot-values (%+v)", f)
	}
	sets := cards.Sets()
	inDecks := len(sets)
	for _, c := range cards.Challenges() {
		if s, ok := cards.ChallengeSet(c); ok {
			sets = append(sets, s)
		}
	}
	for i, s := range sets {
		for _, d := range s.Cards {
			var kinds []valueKind
			item := kindPassive
			if d.Tap {
				item = kindActive
			}
			switch {
			case d.Kind == engine.LootCard:
				kinds = append(kinds, kindLoot)
				if d.Trinket {
					kinds = append(kinds, item)
				}
			case d.Kind == engine.TreasureCard:
				kinds = append(kinds, item)
			}
			for _, k := range kinds {
				v, ok := cardTable.cards[k][string(d.Ref)]
				switch {
				case ok && (math.IsNaN(v) || math.IsInf(v, 0)):
					t.Errorf("%s (%s): value %v", d.Ref, kindNames[k], v)
				case !ok && i < inDecks && !d.Outside:
					t.Errorf("%s (%s) of set %s is not in the table: run make bot-values BOT_VALUES=-missing", d.Ref, kindNames[k], s.Code)
				}
			}
		}
	}
	for k, m := range cardTable.means {
		if math.IsNaN(m) || m < 0 {
			t.Errorf("%s mean %v", kindNames[k], m)
		}
	}
}

// TestValueFallback: a card the table lacks is worth its kind's mean,
// in hand and in play; with flat the old values hold.
func TestValueFallback(t *testing.T) {
	defer func(v *values) { cardTable = v }(cardTable)
	cardTable = newValues(ValueFile{Cards: []ValueEntry{
		{Key: "a_penny", Kind: KindLoot, Value: 1, N: 1e7},
		{Key: "bomb", Kind: KindLoot, Value: 3, N: 1e7},
		{Key: "breakfast", Kind: KindPassive, Value: 20, N: 1e7},
		{Key: "the_battery", Kind: KindActive, Value: 30, N: 1e7},
	}})
	g := table(t, engine.SituationSetup{
		Players: []engine.SituationPlayer{
			{Character: "isaac", Items: []engine.CardRef{"the_d6", "boomerang"}, Hand: []engine.CardRef{"a_nickel", "a_penny"}},
			{Character: "cain", Items: []engine.CardRef{"breakfast"}, Hand: []engine.CardRef{"a_penny"}},
		},
	})
	pl0, pl1 := g.Players[0], g.Players[1]
	d6 := pl0.InPlay[0] // a starting item: outside every deck
	for _, c := range []struct {
		name     string
		id       engine.ObjectID
		seen     bool
		want     float64
		wantFlat float64
	}{
		{"own unknown loot", pl0.Hand[0], true, 2, wLoot},
		{"own known loot", pl0.Hand[1], true, 1, wLoot},
		{"other's loot", pl1.Hand[0], false, 2, wLoot},
		{"starting item", d6, true, 30, wItem},
		{"unknown active item", pl0.InPlay[1], true, 30, wItem},
		{"known passive item", pl1.InPlay[0], true, 20, wItem},
	} {
		if got := cardValue(g, c.id, c.seen, false); math.Abs(got-c.want) > 1e-4 {
			t.Errorf("%s: %v, want %v", c.name, got, c.want)
		}
		if got := cardValue(g, c.id, c.seen, true); got != c.wantFlat {
			t.Errorf("%s flat: %v, want %v", c.name, got, c.wantFlat)
		}
	}
	if !g.HasTap(d6) || !g.HasTap(pl0.InPlay[1]) {
		t.Fatal("the test's active items have no ↷")
	}
	// Seat 0 sees its own hand by card; seat 1 sees it by the mean.
	if a, b := side(g, 0, 0, false), side(g, 0, 1, false); math.Abs(a-b-(2+1-2-2)) > 1e-4 {
		t.Errorf("own hand %v, seen by another %v", a, b)
	}
}

// TestValueShrink: a card's value is pulled to its kind's mean by
// valueShrink samples; a kind with no sample keeps the flat value.
func TestValueShrink(t *testing.T) {
	v := newValues(ValueFile{Cards: []ValueEntry{
		{Key: "a", Kind: KindLoot, Value: 4, N: 8},
		{Key: "b", Kind: KindLoot, Value: 0, N: 1000},
		{Key: "c", Kind: KindLoot, Value: 100, N: 0},
	}})
	if m := v.mean(kindLoot); m != 2 {
		t.Errorf("loot mean %v, want 2: a card with no sample is not in it", m)
	}
	if a := v.cards[kindLoot]["a"]; math.Abs(a-(8*4+valueShrink*2)/(8+valueShrink)) > 1e-9 {
		t.Errorf("a: %v", a)
	}
	if c := v.cards[kindLoot]["c"]; c != 2 {
		t.Errorf("c with no sample: %v, want the mean", c)
	}
	if b := v.cards[kindLoot]["b"]; b > 0.02 {
		t.Errorf("b with many samples: %v, want near 0", b)
	}
	if v.mean(kindActive) != wItem || v.mean(kindPassive) != wItem {
		t.Errorf("item means %v %v, want the flat %v", v.mean(kindActive), v.mean(kindPassive), wItem)
	}
}
