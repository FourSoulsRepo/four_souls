package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestPitfall(t *testing.T) {
	tb := fight(t, "fly")
	tb.G.PutOnDeck(engine.MonsterDeck, "pitfall")
	tb.Attack("fly", 6)
	if !tb.G.Turn.EndDeclared {
		t.Error("the turn ends")
	}
}
