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
// in hand and in play; a hand card counts wHand of its value; items are
// scaled so their mean is wItem; with flat the old values hold.
func TestValueFallback(t *testing.T) {
	defer func(v *values) { cardTable = v }(cardTable)
	cardTable = newValues(ValueFile{Cards: []ValueEntry{
		{Key: "a_penny", Kind: KindLoot, Value: 1, N: 10},
		{Key: "bomb", Kind: KindLoot, Value: 3, N: 10},
		{Key: "breakfast", Kind: KindPassive, Value: 20, N: 10},
		{Key: "the_battery", Kind: KindActive, Value: 30, N: 10},
	}})
	scale := wItem / 25 // the mean of the two items
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
		{"own unknown loot", pl0.Hand[0], true, wHand * 2, wLoot},
		{"own known loot", pl0.Hand[1], true, wHand * 1, wLoot},
		{"other's loot", pl1.Hand[0], false, wHand * 2, wLoot},
		{"starting item", d6, true, 30 * scale, wItem},
		{"unknown active item", pl0.InPlay[1], true, 30 * scale, wItem},
		{"known passive item", pl1.InPlay[0], true, 20 * scale, wItem},
	} {
		if got := cardValue(g, c.id, c.seen, false); math.Abs(got-c.want) > 1e-9 {
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
	if a, b := side(g, 0, 0, false), side(g, 0, 1, false); math.Abs(a-b-wHand*(2+1-2-2)) > 1e-9 {
		t.Errorf("own hand %v, seen by another %v", a, b)
	}
}

// TestValueShrink: a card's value is weighed against its kind's mean by
// precision, with τ² the kind's spread less its noise (at least
// tauFloor2); a card with no sample gets the mean, one with no error
// keeps its value; items are scaled to a mean of wItem; a kind with no
// sample keeps the flat value.
func TestValueShrink(t *testing.T) {
	v := newValues(ValueFile{Cards: []ValueEntry{
		{Key: "a", Kind: KindLoot, Value: 4, N: 8, SE: 1},
		{Key: "b", Kind: KindLoot, Value: 0, N: 8, SE: 1},
		{Key: "c", Kind: KindLoot, Value: 100, N: 0},
		{Key: "d", Kind: KindLoot, Value: 2, N: 8},
	}})
	// Mean 2; spread (4 + 4 + 0) / 2 = 4, less mean se² 2/3: τ² = 10/3.
	m, tau2 := 2.0, 4-2.0/3
	if got := v.mean(kindLoot); got != m {
		t.Errorf("loot mean %v, want %v: a card with no sample is not in it", got, m)
	}
	if got := v.tau2[kindLoot]; math.Abs(got-tau2) > 1e-9 {
		t.Errorf("τ² %v, want %v", got, tau2)
	}
	if a, want := v.cards[kindLoot]["a"], (4+m/tau2)/(1+1/tau2); math.Abs(a-want) > 1e-9 {
		t.Errorf("a: %v, want %v", a, want)
	}
	if c := v.cards[kindLoot]["c"]; c != m {
		t.Errorf("c with no sample: %v, want the mean", c)
	}
	if d := v.cards[kindLoot]["d"]; d != 2 {
		t.Errorf("d with no error: %v, want its own 2", d)
	}
	if v.mean(kindActive) != wItem || v.mean(kindPassive) != wItem {
		t.Errorf("item means %v %v, want the flat %v", v.mean(kindActive), v.mean(kindPassive), wItem)
	}
	// Two items alike: τ² at its floor, both pulled to the mean; the
	// items scaled so their mean is wItem, a bad one stays bad.
	w := newValues(ValueFile{Cards: []ValueEntry{
		{Key: "x", Kind: KindPassive, Value: 6, N: 8, SE: 3},
		{Key: "y", Kind: KindPassive, Value: 6, N: 8, SE: 3},
		{Key: "z", Kind: KindActive, Value: 9, N: 8},
		{Key: "u", Kind: KindActive, Value: -1, N: 8},
	}})
	if w.tau2[kindPassive] != tauFloor2 {
		t.Errorf("τ² of equal cards %v, want the floor", w.tau2[kindPassive])
	}
	scale := wItem / 5 // (6 + 6 + 9 - 1) / 4
	if x := w.cards[kindPassive]["x"]; math.Abs(x-6*scale) > 1e-9 {
		t.Errorf("x: %v, want %v", x, 6*scale)
	}
	if z, u := w.cards[kindActive]["z"], w.cards[kindActive]["u"]; math.Abs(z-9*scale) > 1e-9 || math.Abs(u+scale) > 1e-9 {
		t.Errorf("z %v, u %v: want %v and %v", z, u, 9*scale, -scale)
	}
}
