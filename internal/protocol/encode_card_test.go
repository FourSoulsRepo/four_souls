package protocol

import (
	"encoding/json"
	"reflect"
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

// sampleCards covers each optional field of a card view, on and off, and
// strings that need escapes (18.17).
func sampleCards() []engine.CardView {
	return []engine.CardView{
		{ID: 1, Card: "isaac", Role: engine.RoleCharacter, Charged: true},
		{ID: 2, Card: "the_d6"},
		{ID: 3, Card: "a_dime", Damage: 2, Souls: 1},
		{ID: 4, Card: "bum_bo", Counters: []engine.Counter{{Name: "gold", Count: 2}, {Count: 0}}},
		{ID: 5, Card: "dagaz", Stats: &engine.MonsterStats{HP: 3, DC: 4, ATK: 2, PrintedHP: 3}},
		{ID: 6, Card: "tech_x", Abilities: []string{"Tap: gain 1", `Quote "x" <b>&</b>`}},
		{ID: 7, Card: "dice_bag", Dice: []int{1, 6}, Flipped: true, Copy: "a_dime", MoreCopies: []engine.CardRef{"pooter", "leaper"}},
		{ID: 8, Card: "", Role: engine.RoleMonster, Charged: false},
		{ID: 9, Card: "héllo\u2028", Abilities: []string{"tab\there"}},
	}
}

// TestAppendCardMatchesJSON: each card encodes to the bytes encoding/json
// writes for it.
func TestAppendCardMatchesJSON(t *testing.T) {
	for _, c := range sampleCards() {
		want, err := json.Marshal(c)
		if err != nil {
			t.Fatal(err)
		}
		got := appendCard(nil, c)
		if string(got) != string(want) {
			t.Errorf("card %d:\n got %s\nwant %s", c.ID, got, want)
		}
	}
}

// TestSampleCardsSetEveryField: each field of a card view is set in some
// sample, so a new field the encoder misses fails the test above
// (19.157: more_copies was never sent).
func TestSampleCardsSetEveryField(t *testing.T) {
	typ := reflect.TypeFor[engine.CardView]()
	for i := range typ.NumField() {
		set := false
		for _, c := range sampleCards() {
			if !reflect.ValueOf(c).Field(i).IsZero() {
				set = true
			}
		}
		if !set {
			t.Errorf("no sample card sets %s", typ.Field(i).Name)
		}
	}
}

// BenchmarkAppendCard and BenchmarkMarshalCard compare one card encoded by
// hand and by encoding/json.
func BenchmarkAppendCard(b *testing.B) {
	c := sampleCards()[5]
	var buf []byte
	b.ReportAllocs()
	for b.Loop() {
		buf = appendCard(buf[:0], c)
	}
}

func BenchmarkMarshalCard(b *testing.B) {
	c := sampleCards()[5]
	b.ReportAllocs()
	for b.Loop() {
		if _, err := json.Marshal(c); err != nil {
			b.Fatal(err)
		}
	}
}
