package r

import (
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/cardkit"
)

// The card's own lines; their words are in en.json (ADR 021).
var (
	classicRollerWhichPlayerChooses = engine.NewKey("card.classic_roller.which_player_chooses_the_result")
)

// Classic Roller (Eternal Treasure Card)
//
//	{Tap Effect}Choose a dice roll. Its controller rerolls it, but rolls 2 dice instead. They choose another player. That player chooses one of the rolls as the result.
//	At the end of your turn, recharge this.
//	-Eternal- This can't be destroyed or put into discard.
var classicRoller = engine.CardDef{
	Ref:       "classic_roller",
	Kind:      engine.TreasureCard,
	Copies:    1,
	Eternal:   true,
	ThreePlus: true,
	Outside:   true,
	Tap:       true,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Activated,
			Text:    "↷: Choose a dice roll. Its controller rerolls it, but rolls 2 dice instead. They choose another player. That player chooses one of the rolls as the result.",
			Costs:   []engine.Cost{engine.Tap()},
			Targets: []engine.TargetSpec{engine.Choose(engine.TargetDiceRoll)},
			Effects: []engine.Effect{engine.Ask(func(c *engine.Ctx, a []int) {
				roller := c.Targets[0].Player
				others := othersOf(c.G, roller)
				if a[0] >= 0 && a[0] < len(others) {
					c.G.RerollWithDice(c.Targets[0].StackID, 2, others[a[0]])
				}
			}, engine.Question{
				Text:   classicRollerWhichPlayerChooses,
				Player: func(c *engine.Ctx, _ []int) engine.PlayerID { return c.Targets[0].Player },
				Options: func(c *engine.Ctx, _ []int) []engine.Msg {
					return playerLabels(c.G, othersOf(c.G, c.Targets[0].Player))
				},
			})},
		},
		cardkit.RechargeAtEndOfTurn,
	},
}
