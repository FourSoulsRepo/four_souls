package cards

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

// loadCheck saves and loads a game while a nested ask waits (19.177,
// 19.143): the loaded twin must be the same game, and must step the
// same, until no nested ask waits. Its save keeps the last clean state
// and replays the steps since it, so the nested question is rebuilt by
// the card's own code.
type loadCheck struct {
	twin  *engine.Game
	sets  []engine.CardSet
	saves int
	cards map[engine.CardRef]bool // the cards whose nested asks were saved
}

func newLoadCheck(sets []engine.CardSet) *loadCheck {
	return &loadCheck{sets: sets, cards: map[engine.CardRef]bool{}}
}

// after runs once g has applied in: with Step if step, else Apply.
func (c *loadCheck) after(t *testing.T, g *engine.Game, in engine.Intent, step bool) {
	t.Helper()
	if c.twin != nil {
		var err error
		if step {
			_, err = c.twin.Step(in)
		} else {
			_, err = c.twin.Apply(in)
		}
		if err != nil {
			t.Fatalf("loaded at a nested ask: %+v refused: %v", in, err)
		}
		c.same(t, g, "after a step")
		if !g.NestedAsk() {
			c.twin = nil
		}
	}
	if c.twin != nil || !g.NestedAsk() {
		return
	}
	data, err := g.Save()
	if err != nil {
		t.Fatalf("save at a nested ask: %v", err)
	}
	if c.twin, err = engine.Load(data, c.sets...); err != nil {
		t.Fatalf("load at a nested ask: %v", err)
	}
	c.same(t, g, "after the load")
	c.saves++
	if ch := g.Choice; ch != nil && ch.Ask != nil {
		c.cards[g.CardOf(ch.Ask.Source)] = true
	}
}

// same fails unless the twin's whole state is g's.
func (c *loadCheck) same(t *testing.T, g *engine.Game, when string) {
	t.Helper()
	a, errA := g.Checksum()
	b, errB := c.twin.Checksum()
	if errA != nil || errB != nil {
		t.Fatal(errA, errB)
	}
	if a != b {
		t.Fatalf("saved at a nested ask, loaded: the state differs %s", when)
	}
}
