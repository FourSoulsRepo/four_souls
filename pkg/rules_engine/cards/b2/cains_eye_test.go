package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

func TestCainsEye(t *testing.T) {
	tb := enginetest.NewSetup(t, engine.SituationSetup{
		Players: []engine.SituationPlayer{seat("isaac", "cains_eye"), seat("cain")},
		Active:  1,
	}, Set)
	top := tb.G.DeckTop(engine.LootDeck, 1)[0]
	tb.EndTurn("ask.put_on_bottom " + string(tb.G.Object(top).Card))
	if tb.G.Decks[engine.LootDeck][0] != top {
		t.Error("the top card did not go to the bottom")
	}
}

func TestCainsEyeShowsTheTopCardToItsOwner(t *testing.T) {
	tb := enginetest.NewSetup(t, engine.SituationSetup{
		Players: []engine.SituationPlayer{seat("isaac", "cains_eye"), seat("cain")},
		Active:  1,
	}, Set)
	top := tb.G.DeckTop(engine.LootDeck, 1)[0]
	if _, err := tb.G.Apply(engine.Intent{Player: 1, Kind: engine.IntentEndTurn}); err != nil {
		t.Fatal(err)
	}
	for tb.G.Prompt().Kind == engine.PromptPriority {
		if _, err := tb.G.Apply(engine.Intent{Player: tb.G.Prompt().Player, Kind: engine.IntentPass}); err != nil {
			t.Fatal(err)
		}
	}
	if w := tb.G.Prompt(); w.Kind != engine.PromptChoose || w.Player != 0 {
		t.Fatalf("prompt %+v, want Isaac's choice", w)
	}
	if s := tb.G.View(engine.Viewer{Kind: engine.ViewPlayer, Player: 0}).Subject; s == nil || s.ID != top {
		t.Errorf("chooser's subject %+v, want the loot deck top (19.11)", s)
	}
	for _, v := range []engine.Viewer{{Kind: engine.ViewPlayer, Player: 1}, {Kind: engine.ViewSpectator}} {
		if s := tb.G.View(v).Subject; s != nil {
			t.Errorf("viewer %+v sees the hidden top card %+v", v, *s)
		}
	}
}
