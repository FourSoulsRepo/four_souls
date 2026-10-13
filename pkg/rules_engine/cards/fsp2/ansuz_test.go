package fsp2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestAnsuz(t *testing.T) {
	tb := lootTable(t, "ansuz")
	tb.G.PutOnDeck(engine.MonsterDeck, "fly", "leech", "trite", "fatty") // fatty on top
	tb.Play(0, "ansuz", "monster deck", "fly", "leech", "trite", "fatty")
	top := tb.G.DeckTop(engine.MonsterDeck, 4)
	var got []engine.CardRef
	for _, id := range top {
		got = append(got, tb.G.Object(id).Card)
	}
	if len(got) != 4 || got[0] != "fly" || got[3] != "fatty" {
		t.Errorf("top 4 = %v, want fly, leech, trite, fatty", got)
	}
}
