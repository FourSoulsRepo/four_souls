package protocol

import (
	"encoding/json"
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

// samplePlayers covers the optional fields of a player view and the nil
// and empty lists, which encode differently (18.17).
func samplePlayers() []engine.PlayerView {
	top := engine.CardView{ID: 9, Card: "a_dime", Charged: true}
	return []engine.PlayerView{
		{
			ID: 0, Character: engine.CardView{ID: 1, Card: "isaac"}, HP: 2, ATK: 1, Cents: 3, Souls: 1,
			InPlay: []engine.CardView{{ID: 2, Card: "the_d6"}}, HandSize: 2, Hand: []engine.CardView{{ID: 3, Card: "x"}},
			LootPlays: 1, Attacks: 1, Purchases: 1, TreasureTop: &top, PrintedHP: 2, PrintedATK: 1,
			Sources: []engine.StatSourceView{{Stat: "hp", Amount: 1}},
		},
		{ID: 1, Dead: true, InPlay: nil, HandSize: 0, Hand: []engine.CardView{}},
		{ID: 2, InPlay: []engine.CardView{}},
	}
}

// TestAppendPlayerMatchesJSON: each player view, and each slot, encodes to
// the bytes encoding/json writes for it.
func TestAppendPlayerMatchesJSON(t *testing.T) {
	for _, p := range samplePlayers() {
		want, err := json.Marshal(p)
		if err != nil {
			t.Fatal(err)
		}
		if got := appendPlayer(nil, p); string(got) != string(want) {
			t.Errorf("player %d:\n got %s\nwant %s", p.ID, got, want)
		}
	}
	slots := []engine.SlotView{
		{},
		{Top: &engine.CardView{ID: 4, Card: "dagaz"}},
		{Covered: []engine.CardView{{ID: 5, Card: "y"}}},
		{Top: &engine.CardView{ID: 6, Card: "z"}, Covered: []engine.CardView{{ID: 7, Card: "w"}, {ID: 8, Card: "v"}}},
	}
	for _, s := range slots {
		want, err := json.Marshal(s)
		if err != nil {
			t.Fatal(err)
		}
		if got := appendSlot(nil, s); string(got) != string(want) {
			t.Errorf("slot:\n got %s\nwant %s", got, want)
		}
	}
}

// BenchmarkAppendPlayer and BenchmarkMarshalPlayer compare one player view.
func BenchmarkAppendPlayer(b *testing.B) {
	p := samplePlayers()[0]
	var buf []byte
	b.ReportAllocs()
	for b.Loop() {
		buf = appendPlayer(buf[:0], p)
	}
}

func BenchmarkMarshalPlayer(b *testing.B) {
	p := samplePlayers()[0]
	b.ReportAllocs()
	for b.Loop() {
		if _, err := json.Marshal(p); err != nil {
			b.Fatal(err)
		}
	}
}
