package fsp2

import (
	"slices"
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

func TestBossRush(t *testing.T) {
	tb := fight(t, "fatty", "fly")
	tb.G.PutOnDeck(engine.MonsterDeck, "widow", "dip", "monstro", "spider") // spider on top
	tb.RevealFromDeck("boss_rush", "ask.slot 2", "ask.slot 2")
	if top(tb, 1) != "widow" {
		t.Errorf("slot 2 = %s, want the second boss on top", top(tb, 1))
	}
	if tb.G.Turn.MustAttacks != 1 {
		t.Errorf("must attacks = %d, want 1", tb.G.Turn.MustAttacks)
	}
}

// TestBossRushOwnSlot: a boss put in Boss Rush's own slot covers what
// was there; Boss Rush still goes to discard (19.61).
func TestBossRushOwnSlot(t *testing.T) {
	tb := fight(t, "fatty", "fly")
	tb.G.PutOnDeck(engine.MonsterDeck, "widow", "dip", "monstro", "spider") // spider on top
	tb.RevealFromDeck("boss_rush", "ask.slot 1", "ask.slot 2")
	if top(tb, 0) != "monstro" || top(tb, 1) != "widow" {
		t.Fatalf("slots = %s, %s; want monstro, widow", top(tb, 0), top(tb, 1))
	}
	if !slices.Contains(discardCards(tb), "boss_rush") {
		t.Errorf("monster discard = %v, want boss_rush", discardCards(tb))
	}
	if cards := slotCards(tb, 0); !slices.Equal(cards, []engine.CardRef{"fatty", "monstro"}) {
		t.Errorf("slot 1 = %v, want [fatty monstro]", cards)
	}
	if err := tb.G.CheckInvariants(); err != nil {
		t.Error(err)
	}
	tb.G.KillObject(topID(tb, 0))
	tb.Settle()
	if top(tb, 0) != "fatty" {
		t.Errorf("slot 1 after monstro died = %s, want fatty", top(tb, 0))
	}
}

func discardCards(tb *enginetest.Table) []engine.CardRef {
	var out []engine.CardRef
	for _, id := range tb.G.Discards[engine.MonsterDeck] {
		out = append(out, tb.G.Object(id).Card)
	}
	return out
}

func slotCards(tb *enginetest.Table, i int) []engine.CardRef {
	var out []engine.CardRef
	for _, id := range tb.G.Monsters[i].Cards {
		out = append(out, tb.G.Object(id).Card)
	}
	return out
}
