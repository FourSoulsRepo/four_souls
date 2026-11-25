package challenge_test

import (
	"maps"
	"slices"
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

// The flip, Nobody's Fool, Ghostly Trickery, You Don't Scare Me,
// Hallows' Eve Bounty and the goal are in
// testdata/situations/trick_treat_normhardultra (24.27). What stays
// here: the flip question's text, the view, and the abilities_gained
// event's fields.

// R-TRICK-11, Q-43: as Cain's turn starts, the flip question is the
// rule card's own (its text key).
func TestTrickFlipQuestionText(t *testing.T) {
	tb := trick(t, engine.DifficultyNormal, engine.SituationPlayer{}, engine.SituationPlayer{})
	g := tb.G
	if _, err := g.Apply(engine.Intent{Player: 0, Kind: engine.IntentEndTurn}); err != nil {
		t.Fatal(err)
	}
	for range 100 {
		tb.ThrowPending()
		if w := g.Prompt(); w.Kind == engine.PromptPriority {
			tb.Pass(int(w.Player))
			continue
		}
		break
	}
	if w := g.Prompt(); w.Player != 1 || g.Turn.Active != 1 || w.Text.Key.String() != "card.trick_treat.flip_to_treat" {
		t.Fatalf("flip asked of %d on %d's turn: %+v", w.Player, g.Turn.Active, w.Text)
	}
}

// 18.29, ADR 037: the view sends what makes a player take more damage.
// Expected from the card texts: Ghostly Trickery on Hard rolls 2x,
// each "1-2: Damage you would take is increased by 1 till end of turn";
// Daddy Haunt: "take that much damage +1 instead". So +3: two +1s of
// this turn from the Trick card, +1 in play from Daddy Haunt.
func TestGhostlyTrickeryView(t *testing.T) {
	tb := trick(t, engine.DifficultyHard, engine.SituationPlayer{Items: []engine.CardRef{"daddy_haunt"}}, engine.SituationPlayer{})
	attack(tb, "fly", 1, 2, 6, 6, 6)
	haunt := tb.Find(0, "daddy_haunt")
	type part struct {
		source engine.ObjectID
		card   engine.CardRef
		amount int
		lasts  string
	}
	want := map[part]int{
		{tb.G.Challenge.Card, trickCard, 1, "turn"}: 2,
		{haunt, "daddy_haunt", 1, "in_play"}:        1,
	}
	for _, v := range []engine.Viewer{{Kind: engine.ViewSpectator}, {Kind: engine.ViewPlayer, Player: 1}} {
		got := map[part]int{}
		sum := 0
		for _, s := range tb.G.View(v).Players[0].Sources {
			if s.Stat == "damage_taken" {
				got[part{s.Source, s.Card, s.Amount, s.Lasts}]++
				sum += s.Amount
			}
		}
		if sum != 3 || !maps.Equal(got, want) {
			t.Errorf("viewer %+v: damage taken %d from %v, want 3 from %v", v, sum, got, want)
		}
	}
}

// R-TRICK-13, Q-47: the abilities_gained event names the revealed
// card (Text) and the attacked monster (Card). What the gained ability
// does is in
// trick_treat_normhardultra/ghostly-trickery-3-4-gains-a-revealed-monsters-abilities.json.
func TestGhostlyTrickeryGainEvent(t *testing.T) {
	tb := trick(t, engine.DifficultyNormal, engine.SituationPlayer{}, engine.SituationPlayer{})
	tb.G.PutOnDeck(engine.MonsterDeck, "boom_fly")
	ev := attack(tb, "fly", 3, 6, 1)
	gained := slices.IndexFunc(ev, func(e engine.Event) bool { return e.Kind == engine.EvAbilitiesGained })
	if gained < 0 || ev[gained].Text != "boom_fly" || ev[gained].Card != "fly" {
		t.Fatalf("no abilities_gained of boom_fly: %+v", ev)
	}
}
