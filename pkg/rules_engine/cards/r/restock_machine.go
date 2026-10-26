package r

import (
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/cardkit"
)

// The card's own lines; their words are in en.json (ADR 021).
var (
	restockMachineMonsters = engine.NewKey("card.restock_machine.monsters")
	restockMachineShop     = engine.NewKey("card.restock_machine.shop")
)

// Restock Machine (Room Card)
//
//	{Paid Effect}Pay 2¢:
//	Put each shop item or each monster not being attacked into discard.
var restockMachine = engine.CardDef{
	Ref:    "restock_machine",
	Kind:   engine.RoomCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind:  engine.Activated,
		Text:  "$ Pay 2¢: Put each shop item or each monster not being attacked into discard.",
		Costs: []engine.Cost{engine.PayCents(2)},
		Effects: []engine.Effect{engine.Ask(func(c *engine.Ctx, a []int) {
			if a[0] == 0 {
				c.Do(discardShop)
			} else {
				c.Do(eachMonsterToDiscard(true))
			}
		}, engine.Question{Text: cardkit.AskPutWhichIntoDiscard, Options: func(*engine.Ctx, []int) []engine.Msg {
			return engine.Options(restockMachineShop, restockMachineMonsters)
		}})},
	}},
}
