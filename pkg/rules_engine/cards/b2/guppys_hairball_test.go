package b2

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestGuppysHairball(t *testing.T) {
	tb := monsterTable(t, []engine.CardRef{"leech"}, seat("isaac", "guppys_hairball"), seat("cain"))
	tb.Attack("leech", 1, 6, 6) // miss: the hairball rolls 6 and prevents 1 of 2; then kill
	if d := tb.G.Players[0].Damage; d != 1 {
		t.Errorf("damage = %d, want 1", d)
	}
}
