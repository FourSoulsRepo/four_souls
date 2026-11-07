package cards

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

// TestTurnStartPriority: every character's player holds priority as
// their turn starts, and may do more than pass. The character was
// recharged, and its tap ability needs no target. So the server never
// passes for them (N-05), and the e2e turn banner test always finds
// Pass live under the banner, whatever the deal (19.125).
func TestTurnStartPriority(t *testing.T) {
	sets := Sets()
	var chars []engine.CardRef
	for _, set := range sets {
		for _, c := range set.Cards {
			// A back face is never dealt (ADR 023).
			if c.Kind == engine.CharacterCard && !c.Outside {
				chars = append(chars, c.Ref)
			}
		}
	}
	for i, c := range chars {
		other := chars[(i+1)%len(chars)]
		for seed := uint64(1); seed <= 5; seed++ {
			g, _, err := engine.NewGame(engine.Setup{Seed: seed, Players: 2, Sets: sets, Characters: []engine.CardRef{c, other}})
			if err != nil {
				t.Fatal(err)
			}
			for turn := 2; turn <= 3; turn++ {
				toTurn(t, g, turn)
				p := g.Prompt()
				if p.Kind != engine.PromptPriority || p.Player != g.Turn.Active || g.Turn.Step != engine.StepStartTriggers {
					t.Fatalf("%s seed %d turn %d: prompt %+v at step %d, want the active player's priority", c, seed, turn, p, g.Turn.Step)
				}
				if !g.CanAct(p.Player) {
					t.Fatalf("%s seed %d turn %d: player %d can only pass", c, seed, turn, p.Player)
				}
			}
		}
	}
}

// toTurn plays on until turn n starts: each prompt gets its first
// allowed answer; priority passes, and the active player ends the turn.
func toTurn(t *testing.T, g *engine.Game, n int) {
	t.Helper()
	for range 5000 {
		if g.Turn.Number >= n || g.Over {
			return
		}
		p := g.Prompt()
		allowed := g.Allowed(p.Player)
		if len(allowed) == 0 {
			t.Fatalf("nothing allowed, prompt %+v", p)
		}
		in := allowed[0]
		if p.Kind == engine.PromptPriority {
			in = engine.Intent{Player: p.Player, Kind: engine.IntentPass}
			for _, a := range allowed {
				if a.Kind == engine.IntentEndTurn {
					in = a
				}
			}
		}
		if in.Kind == engine.IntentDiscard {
			in.Objects = in.Objects[:p.Count]
		}
		if _, err := g.Apply(in); err != nil {
			t.Fatalf("%+v refused: %v", in, err)
		}
	}
	t.Fatalf("turn %d never started", n)
}
