package bot

import (
	"math"
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards"
)

// TestDeckMonsterMean: the policy judges the monster deck's unknown top
// as a monster of HP deckHP, evasion deckDC and ATK deckATK. They are
// constants, not a mean of the world's deck, so a playout pays nothing
// for them (lead review, 2026-11-17). This keeps them within 1 of the
// mean monster of every set. When it fails (a new set moved the mean):
// set the constants in policy.go to the rounded means it prints, then
// run make bots.
func TestDeckMonsterMean(t *testing.T) {
	var n, hp, dc, atk float64
	for _, s := range cards.Sets() {
		for _, d := range s.Cards {
			if d.Kind != engine.MonsterCard || d.Outside || d.HP <= 0 {
				continue
			}
			n++
			hp += float64(d.HP)
			dc += float64(d.DC)
			atk += float64(d.ATK)
		}
	}
	if n == 0 {
		t.Fatal("no monster in any set")
	}
	hp, dc, atk = hp/n, dc/n, atk/n
	t.Logf("%v monsters: HP %.2f, evasion %.2f, ATK %.2f", n, hp, dc, atk)
	for _, c := range []struct {
		name       string
		mean, kept float64
	}{{"deckHP", hp, deckHP}, {"deckDC", dc, deckDC}, {"deckATK", atk, deckATK}} {
		if math.Abs(c.mean-c.kept) > 1 {
			t.Errorf("%s is %v, the mean of every set %.2f: set it to %v in policy.go, then make bots", c.name, c.kept, c.mean, math.Round(c.mean))
		}
	}
}
