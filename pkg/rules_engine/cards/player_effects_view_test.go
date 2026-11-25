package cards

import (
	"slices"
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

// Every effect on a player is in PlayerView.Sources (ADR 037 item 1,
// 18.30): group A (a card in the player's own area), B (a card
// elsewhere, here a room with a condition), C (a card already gone) and
// shields.

// effect is one source as a test compares it.
type effect struct {
	stat   string
	card   engine.CardRef
	amount int
	lasts  string
}

// effectsOf are player p's sources as every viewer sees them; it fails when two viewers differ (ADR 037
// item 6: effects are public).
func effectsOf(tb *enginetest.Table, p int) []effect {
	tb.T.Helper()
	var first []effect
	for i, v := range []engine.Viewer{{Kind: engine.ViewSpectator}, {Kind: engine.ViewPlayer, Player: 0}, {Kind: engine.ViewPlayer, Player: 1}} {
		var got []effect
		for _, s := range tb.G.View(v).Players[p].Sources {
			got = append(got, effect{s.Stat, s.Card, s.Amount, s.Lasts})
		}
		if i == 0 {
			first = got
		} else if !slices.Equal(got, first) {
			tb.T.Errorf("viewer %+v sees %v, the spectator %v", v, got, first)
		}
	}
	return first
}

func viewTable(t *testing.T, s engine.SituationSetup) *enginetest.Table {
	t.Helper()
	return enginetest.NewSetup(t, s, Sets()...)
}

// A: cards in the player's own area. Steamy Sale: "Shop items you
// purchase cost 5¢ less" (shop price −5, once, not once per shop item).
// Curse of the Soulless: "You can't gain souls" (a flag: amount 1).
func TestViewEffectsOwnArea(t *testing.T) {
	tb := viewTable(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "isaac", Items: []engine.CardRef{"steamy_sale", "curse_of_the_soulless"}},
		{Character: "cain"},
	}})
	want := []effect{
		{"shop_price", "steamy_sale", -5, "in_play"},
		{"no_souls", "curse_of_the_soulless", 1, "in_play"},
	}
	if got := effectsOf(tb, 0); !slices.Equal(got, want) {
		t.Errorf("Isaac's effects %v, want %v", got, want)
	}
	if got := effectsOf(tb, 1); got != nil {
		t.Errorf("Cain's effects %v, want none: the cards are Isaac's", got)
	}
}

// B: Might for the Meek, a room: "Players who control the fewest souls
// or tied for fewest have +1 ATK and may attack an additional time on
// their turn". Isaac has a soul, Cain none: only Cain gets it.
func TestViewEffectsRoomOnFewestSouls(t *testing.T) {
	tb := viewTable(t, engine.SituationSetup{
		Players: []engine.SituationPlayer{
			{Character: "isaac", Souls: []engine.CardRef{"monstro"}},
			{Character: "cain"},
		},
		Rooms: []engine.CardRef{"might_for_the_meek"},
	})
	if tb.G.SoulValue(0) == 0 {
		t.Fatal("Isaac has no soul; the test needs one")
	}
	want := []effect{
		{"atk", "might_for_the_meek", 1, "in_play"},
		{"attacks", "might_for_the_meek", 1, "in_play"},
	}
	if got := effectsOf(tb, 1); !slices.Equal(got, want) {
		t.Errorf("Cain (fewest souls) has %v, want %v", got, want)
	}
	if got := effectsOf(tb, 0); got != nil {
		t.Errorf("Isaac (more souls) has %v, want none", got)
	}
}

// C: Get Out of Jail Card: "Other players can't play loot cards or
// activate items till end of turn". The card is in the discard once it
// resolves; the effect still shows on its player, this turn only.
func TestViewEffectsGoneCard(t *testing.T) {
	tb := viewTable(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "isaac", Hand: []engine.CardRef{"get_out_of_jail_card"}},
		{Character: "cain"},
	}})
	tb.Play(0, "get_out_of_jail_card")
	tb.Settle()
	src := tb.G.View(engine.Viewer{Kind: engine.ViewSpectator}).Players[0].Sources
	if len(src) != 1 || tb.G.Object(src[0].Source).Zone.Kind == engine.ZoneInPlay || tb.G.Object(src[0].Source).Zone.Kind == engine.ZoneHand {
		t.Fatalf("sources %+v: want one, from the played card out of play", src)
	}
	want := []effect{{"lock_turn", "get_out_of_jail_card", 1, "turn"}}
	if got := effectsOf(tb, 0); !slices.Equal(got, want) {
		t.Errorf("Isaac's effects %v, want %v", got, want)
	}
	if got := effectsOf(tb, 1); got != nil {
		t.Errorf("Cain's effects %v, want none: the lock is Isaac's effect", got)
	}
	tb.EndTurn()
	if got := effectsOf(tb, 0); got != nil {
		t.Errorf("next turn: %v, want none (till end of turn)", got)
	}
}

// Shields (ADR 037 item 1.8). Soul Heart: "Prevent the next 1 damage
// they would take this turn": shield 1. Yum Heart: "Prevent the next
// instance of damage": a shield of the whole instance, no amount.
func TestViewEffectsShields(t *testing.T) {
	tb := viewTable(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "isaac", Hand: []engine.CardRef{"soul_heart"}, Items: []engine.CardRef{"yum_heart"}},
		{Character: "cain"},
	}})
	tb.Play(0, "soul_heart", "player 2 (cain)")
	tb.Settle()
	tb.Activate(0, "yum_heart", 0, "player 2 (cain)")
	tb.Settle()
	want := []effect{{"shield", "soul_heart", 1, "turn"}, {"shield", "yum_heart", 0, "turn"}}
	if got := effectsOf(tb, 1); !slices.Equal(got, want) {
		t.Errorf("Cain's effects %v, want %v", got, want)
	}
	if got := effectsOf(tb, 0); got != nil {
		t.Errorf("Isaac's effects %v, want none", got)
	}
	tb.EndTurn()
	if got := effectsOf(tb, 1); got != nil {
		t.Errorf("next turn: %v, want none (this turn only)", got)
	}
}
