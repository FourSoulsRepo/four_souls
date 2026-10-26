package b2

import (
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/cardkit"
)

// The D10 (Passive Treasure Card)
//
//	Each time a player rolls a ❸, you may put the top card of the Monster Deck in a monster slot not being attacked.
var theD10 = engine.CardDef{
	Ref:    "the_d10",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "Each time a player rolls a 3, you may put the top card of the monster deck in a monster slot not being attacked.",
			Trigger: engine.OnRollOf(3),
			Effects: []engine.Effect{engine.Ask(func(c *engine.Ctx, a []int) {
				if slots := freeSlots(c.G); a[0] >= 0 && a[0] < len(slots) {
					c.G.CoverMonsterSlot(slots[a[0]])
				}
			}, engine.Question{Text: cardkit.AskWhichMonsterSlot, Options: func(c *engine.Ctx, _ []int) []engine.Msg {
				var out []engine.Msg
				for _, i := range freeSlots(c.G) {
					out = append(out, engine.SlotOption(i+1, ""))
				}
				return append(out, engine.AskNone.Msg())
			}})},
		},
	},
}

// freeSlots lists monster slots whose monster is not being attacked.
func freeSlots(g *engine.Game) []int {
	var out []int
	for i, s := range g.Monsters {
		if top, ok := s.TopOf(); !ok || top != g.Attack.Target {
			out = append(out, i)
		}
	}
	return out
}
