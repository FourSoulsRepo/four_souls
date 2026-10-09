package rulesengine

import (
	"reflect"
	"testing"
)

func TestRNGIsDeterministic(t *testing.T) {
	a, b := NewRNG(42), NewRNG(42)
	for range 1000 {
		if a.Intn(97) != b.Intn(97) {
			t.Fatal("same seed gave different numbers")
		}
	}
	c := NewRNG(43)
	same := 0
	for range 100 {
		if a.D6() == c.D6() {
			same++
		}
	}
	if same > 50 {
		t.Errorf("different seeds look too similar: %d/100 equal rolls", same)
	}
}

func TestRNGIsFair(t *testing.T) {
	r := NewRNG(7)
	var counts [6]int
	const n = 60000
	for range n {
		counts[r.D6()-1]++
	}
	for face, c := range counts {
		if c < n/6*9/10 || c > n/6*11/10 {
			t.Errorf("face %d came up %d times out of %d", face+1, c, n)
		}
	}
}

// A-08: the state must never contain Go maps.
func TestGameHasNoMaps(t *testing.T) {
	var check func(t reflect.Type, path string)
	seen := map[reflect.Type]bool{}
	check = func(t reflect.Type, path string) {
		if seen[t] {
			return
		}
		seen[t] = true
		switch t.Kind() { //nolint:exhaustive // only containers matter here
		case reflect.Map:
			panic("map in game state at " + path)
		case reflect.Pointer, reflect.Slice, reflect.Array:
			check(t.Elem(), path+"[]")
		case reflect.Struct:
			for i := range t.NumField() {
				f := t.Field(i)
				check(f.Type, path+"."+f.Name)
			}
		}
	}
	defer func() {
		if r := recover(); r != nil {
			t.Fatal(r)
		}
	}()
	check(reflect.TypeFor[Game](), "Game")
}

func TestMoveCreatesNewObject(t *testing.T) {
	g := &Game{RNG: NewRNG(1)}
	g.AddToDeck(LootDeck, "a_penny")
	id, ok := g.drawTop(LootDeck)
	if !ok {
		t.Fatal("deck is empty")
	}
	nid := g.discard(id, LootDeck)
	if nid == id {
		t.Fatal("a card that changes zone must become a new object (R-ZONE-01)")
	}
	if g.Object(id).Zone.Kind != ZoneGone {
		t.Error("the old object should be gone")
	}
	if g.Object(nid).Card != "a_penny" || g.Object(nid).Zone != DiscardZone(LootDeck) {
		t.Errorf("new object = %+v", *g.Object(nid))
	}
}

func TestEmptyDeckRefillsFromDiscard(t *testing.T) {
	g := &Game{RNG: NewRNG(1)}
	g.AddToDeck(LootDeck, "a", "b", "c")
	for range 3 {
		id, _ := g.drawTop(LootDeck)
		g.discard(id, LootDeck)
	}
	if len(g.Decks[LootDeck]) != 0 || len(g.Discards[LootDeck]) != 3 {
		t.Fatalf("deck %v discard %v", g.Decks[LootDeck], g.Discards[LootDeck])
	}
	if _, ok := g.drawTop(LootDeck); !ok {
		t.Fatal("empty deck must refill from its discard (R-ZONE-05)")
	}
	if len(g.Decks[LootDeck]) != 2 || len(g.Discards[LootDeck]) != 0 {
		t.Errorf("after refill: deck %d, discard %d", len(g.Decks[LootDeck]), len(g.Discards[LootDeck]))
	}
	g2 := &Game{RNG: NewRNG(1)}
	if _, ok := g2.drawTop(LootDeck); ok {
		t.Error("drawing from an empty deck and discard must fail")
	}
}

func TestShuffleIsDeterministic(t *testing.T) {
	order := func(seed uint64) []CardRef {
		g := &Game{RNG: NewRNG(seed)}
		g.AddToDeck(TreasureDeck, "a", "b", "c", "d", "e", "f", "g", "h")
		g.ShuffleDeck(TreasureDeck)
		var out []CardRef
		for _, id := range g.Decks[TreasureDeck] {
			out = append(out, g.Object(id).Card)
		}
		return out
	}
	if !reflect.DeepEqual(order(5), order(5)) {
		t.Error("same seed must shuffle the same way")
	}
	if reflect.DeepEqual(order(5), order(6)) {
		t.Error("different seeds should usually shuffle differently")
	}
}

func TestCloneAndChecksum(t *testing.T) {
	g := &Game{RNG: NewRNG(3), Players: []Player{{ID: 0, Cents: 3}}}
	g.AddToDeck(LootDeck, "a", "b")
	sum1, err := g.Checksum()
	if err != nil {
		t.Fatal(err)
	}
	c, err := g.Clone()
	if err != nil {
		t.Fatal(err)
	}
	sum2, err := c.Checksum()
	if err != nil {
		t.Fatal(err)
	}
	if sum1 != sum2 {
		t.Error("a clone must have the same checksum")
	}
	c.Players[0].Cents = 10
	if g.Players[0].Cents != 3 {
		t.Error("changing the clone changed the original")
	}
	sum3, err := c.Checksum()
	if err != nil {
		t.Fatal(err)
	}
	if sum3 == sum1 {
		t.Error("a changed state must change the checksum")
	}
	c.RNG.D6()
	sum4, err := c.Checksum()
	if err != nil {
		t.Fatal(err)
	}
	if sum4 == sum3 {
		t.Error("RNG position is part of the state")
	}
}
