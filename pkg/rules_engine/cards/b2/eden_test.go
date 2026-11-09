package b2

import (
	"slices"
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestEden(t *testing.T) {
	testExtraLoot(t, "eden")
	// Eden's choice comes before the first-player roll (R-SETUP-09).
	g, _, err := engine.NewGame(engine.Setup{Seed: 7, Players: 2, Sets: []engine.CardSet{Set}, Characters: []engine.CardRef{"isaac", "eden"}})
	if err != nil {
		t.Fatal(err)
	}
	p := g.Prompt()
	if p.Kind != engine.PromptChoose || p.Player != 1 || len(p.Options) != 3 {
		t.Fatalf("first prompt %+v, want Eden's player choosing among 3", p)
	}
	// The roll starts only after her choice: no roll in the view yet (19.142).
	if g.FirstRoll != nil {
		t.Fatalf("first-player roll %+v during Eden's choice, want none", g.FirstRoll)
	}
	if v := g.View(engine.Viewer{Player: 0}); v.FirstRoll != nil {
		t.Fatalf("the view sends the roll %+v during Eden's choice", v.FirstRoll)
	}
	top := g.DeckTop(engine.TreasureDeck, 3)
	bottom := len(g.Decks[engine.TreasureDeck])
	if _, err := g.Apply(engine.Intent{Player: 1, Kind: engine.IntentChoose, Choice: 1}); err != nil {
		t.Fatal(err)
	}
	chosen := g.Object(top[1]).Card
	var item *engine.Object
	for _, id := range g.Players[1].InPlay {
		if g.Object(id).Card == chosen {
			item = g.Object(id)
		}
	}
	if item == nil || !item.Eternal || !item.Charged {
		t.Fatalf("starting item %s: %+v, want an eternal charged item", chosen, item)
	}
	// No order question: the other two go to the bottom at random (Q-22).
	deck := g.Decks[engine.TreasureDeck]
	if len(deck) != bottom-1 {
		t.Errorf("treasure deck has %d cards, want %d", len(deck), bottom-1)
	}
	for i, id := range deck[:2] {
		if c := g.Object(id).Card; c != g.Object(top[0]).Card && c != g.Object(top[2]).Card {
			t.Errorf("bottom card %d is %s, not one of the other two", i, c)
		}
	}
	if g.Prompt().Kind != engine.PromptThrow || g.Prompt().Player != 0 || g.FirstRoll == nil || g.FirstRoll.Round != 1 {
		t.Errorf("after the choice the game waits on %+v, roll %+v, want the first roll", g.Prompt(), g.FirstRoll)
	}
}

// TestEdenBeforeCain: with Cain at the table, Eden still chooses first;
// Cain's start comes after her choice (R-SETUP-09, 19.142).
func TestEdenBeforeCain(t *testing.T) {
	g, ev, err := engine.NewGame(engine.Setup{Seed: 3, Players: 2, Sets: []engine.CardSet{Set}, Characters: []engine.CardRef{"cain", "eden"}})
	if err != nil {
		t.Fatal(err)
	}
	for _, e := range ev {
		if e.Kind == engine.EvFirstPlayer || e.Kind == engine.EvGameStarted {
			t.Fatalf("event %s before Eden's choice", e.Kind)
		}
	}
	if p := g.Prompt(); p.Kind != engine.PromptChoose || p.Player != 1 || g.Turn.Number != 0 {
		t.Fatalf("prompt %+v in turn %d, want Eden's choice before turn 1", p, g.Turn.Number)
	}
	ev, err = g.Apply(engine.Intent{Player: 1, Kind: engine.IntentChoose, Choice: 0})
	if err != nil {
		t.Fatal(err)
	}
	started := false
	for _, e := range ev {
		started = started || (e.Kind == engine.EvGameStarted && e.Player == 0)
	}
	if !started || g.Turn.Number != 1 || g.Turn.Active != 0 {
		t.Fatalf("after the choice: turn %+v, started %v, want Cain's turn 1", g.Turn, started)
	}
}

// TestEdensAtOnce: several Edens choose at the same time, each among
// their own cards; the first-player roll waits for all (R-SETUP-09,
// 20.22). Either may answer first; the game replays the same.
func TestEdensAtOnce(t *testing.T) {
	for _, first := range []engine.PlayerID{0, 2} {
		s := engine.Setup{Seed: 11, Players: 3, Sets: []engine.CardSet{Set}, Characters: []engine.CardRef{"eden", "isaac", "eden"}}
		g, _, err := engine.NewGame(s)
		if err != nil {
			t.Fatal(err)
		}
		top := g.DeckTop(engine.TreasureDeck, 6)
		w := g.Prompt()
		if w.Kind != engine.PromptChoose || w.Player != 0 || !slices.Equal(w.Players, []engine.PlayerID{0, 2}) || g.FirstRoll != nil {
			t.Fatalf("prompt %+v, roll %+v, want both Edens choosing, no roll", w, g.FirstRoll)
		}
		// Each sees their own three cards: the top three, then the next.
		for i, p := range []engine.PlayerID{0, 2} {
			v := g.View(engine.Viewer{Player: p})
			if v.Waiting.Player != p || len(v.Options) != 3 || v.Waiting.Look == nil || v.Waiting.Look.Count != 3 {
				t.Fatalf("player %d sees %+v, options %+v", p, v.Waiting, v.Options)
			}
			for j, o := range v.Options {
				if o.Object != top[3*i+j] {
					t.Errorf("player %d option %d is %d, want %d", p, j, o.Object, top[3*i+j])
				}
			}
			if len(g.Allowed(p)) != 3 {
				t.Errorf("player %d may send %v, want 3 picks", p, g.Allowed(p))
			}
		}
		if v := g.View(engine.Viewer{Player: 1}); v.Waiting.Player != 0 || len(v.Waiting.Options) != 0 || len(v.Options) != 0 || v.FirstRoll != nil {
			t.Fatalf("Isaac sees %+v, options %+v, roll %+v", v.Waiting, v.Options, v.FirstRoll)
		}
		if _, err := g.Apply(engine.Intent{Player: 1, Kind: engine.IntentChoose}); err == nil {
			t.Fatal("Isaac answered an Eden's choice")
		}
		var steps []engine.StepRecord
		pick := func(p engine.PlayerID, i int) {
			t.Helper()
			st, err := g.Step(engine.Intent{Player: p, Kind: engine.IntentChoose, Choice: i})
			if err != nil {
				t.Fatalf("player %d picks %d: %v", p, i, err)
			}
			steps = append(steps, st)
			if err := g.CheckInvariants(); err != nil {
				t.Fatal(err)
			}
		}
		second := 2 - first
		pick(first, 1)
		w = g.Prompt()
		if w.Player != second || !slices.Equal(w.Players, []engine.PlayerID{second}) || g.FirstRoll != nil {
			t.Fatalf("after one pick: prompt %+v, roll %+v, want only player %d choosing", w, g.FirstRoll, second)
		}
		if len(g.Allowed(first)) != 0 {
			t.Errorf("player %d may answer twice", first)
		}
		pick(second, 2)
		if w := g.Prompt(); w.Kind != engine.PromptThrow || g.FirstRoll == nil || len(g.StartChoices) != 0 {
			t.Fatalf("after both picks: prompt %+v, roll %+v, want the roll", w, g.FirstRoll)
		}
		picked := map[engine.PlayerID]int{first: 1, second: 2}
		for i, p := range []engine.PlayerID{0, 2} {
			want := g.Object(top[3*i+picked[p]]).Card
			found := false
			for _, id := range g.Players[p].InPlay {
				found = found || (g.Object(id).Card == want && g.Object(id).Eternal)
			}
			if !found {
				t.Errorf("player %d lacks the eternal %s", p, want)
			}
		}
		if _, err := engine.Replay(s, steps); err != nil {
			t.Fatalf("replay with player %d first: %v", first, err)
		}
	}
}
