package fsp2

import (
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/cardkit"
)

// The card's own lines; their words are in en.json (ADR 021).
var (
	telepathyForDummiesChangeYourRoll = engine.NewKey("card.telepathy_for_dummies.change_your_roll_to_the")
)

// Telepathy For Dummies (Passive Treasure Card)
//
//	At the start of your turn, roll-
//	You may change the result of your next roll this turn to this result.
var telepathyForDummies = engine.CardDef{
	Ref:    "telepathy_for_dummies",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Dice: func(g *engine.Game, self engine.ObjectID) []int {
		if n := storedRoll(g, self); n > 0 {
			return []int{n}
		}
		return nil
	},
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "At the start of your turn, roll- You may change the result of your next roll this turn to this result.",
			Trigger: engine.AtStartOfYourTurn(),
			Effects: []engine.Effect{engine.Roll(storeRoll())},
		},
		{
			Kind: engine.Triggered,
			Text: "You may change the result of your next roll this turn to the stored result.",
			Trigger: engine.Trigger{On: engine.EvRollWouldResolve, Match: func(g *engine.Game, self engine.ObjectID, e engine.Event) bool {
				return e.Player == g.Object(self).Controller && storedRoll(g, self) > 0 && !cardkit.UsedOn(g, self, 100+e.StackID)
			}},
			Effects: []engine.Effect{engine.Ask(func(c *engine.Ctx, a []int) {
				n := storedRoll(c.G, c.Source)
				clearStoredRoll(c.G, c.Source) // only the next roll
				it, ok := c.G.StackItemByID(c.EventStack)
				if ok && a[0] == 0 && n > 0 {
					c.G.SetRoll(it.ID, n, c.Source)
				}
			}, engine.Question{Text: telepathyForDummiesChangeYourRoll, Options: func(*engine.Ctx, []int) []engine.Msg { return engine.Options(engine.AskYes, engine.AskNo) }})},
		},
	},
}
