package fsp2

import (
	"slices"
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestMutantSpider(t *testing.T) {
	tb := armed(t, items("mutant_spider"), "trite", "fatty")
	tb.Activate(0, "mutant_spider", 0)
	tb.G.ForceRolls(1, 2, 6, 3)
	if _, err := tb.G.Apply(attackIntent(tb)); err != nil {
		t.Fatal(err)
	}
	var dice []int
	settleAny2 := func() {
		for range 20 {
			w := tb.G.Prompt()
			if w.Kind == engine.PromptThrow {
				for _, e := range tb.ThrowPending() {
					if e.Kind == engine.EvDiceRolled {
						dice = e.Dice
					}
				}
				continue
			}
			if len(w.Options) == 4 {
				tb.Choose("6")
				return
			}
			if w.Text == "" && len(w.Options) > 0 {
				tb.Choose("trite")
				continue
			}
			tb.Pass(int(w.Player))
		}
	}
	settleAny2()
	tb.Settle()
	if !slices.Equal(dice, []int{1, 2, 6, 3}) {
		t.Errorf("the throw shows all 4 dice (19.39): %v", dice)
	}
	if top(tb, 0) == "trite" {
		t.Error("the chosen 6 hits the trite (DC 5)")
	}
}
