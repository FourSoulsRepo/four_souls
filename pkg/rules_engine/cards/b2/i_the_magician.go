package b2

import (
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/cardkit"
)

// The card's own lines; their words are in en.json (ADR 021).
var (
	iTheMagicianChangeTheRoll = engine.NewKey("card.i_the_magician.change_the_roll_to")
)

// I. The Magician (Wildcard Card)
//
//	Change the result of a dice roll to a number of your choosing.
var iTheMagician = engine.CardDef{
	Ref:    "i_the_magician",
	Kind:   engine.LootCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.LootAbility,
			Text:    "Change the result of a dice roll to a number of your choosing.",
			Targets: []engine.TargetSpec{engine.Choose(engine.TargetDiceRoll)},
			Effects: []engine.Effect{engine.Ask(func(c *engine.Ctx, a []int) {
				c.G.SetRoll(c.Targets[0].StackID, a[0]+1, c.Source)
			}, engine.Question{
				Text:    iTheMagicianChangeTheRoll,
				Options: cardkit.DieFaces,
			})},
		},
	},
}
