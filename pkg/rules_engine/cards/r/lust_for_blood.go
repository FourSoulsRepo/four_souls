package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Lust For Blood! (Good Event Card)
//
//	Choose a player. You must make an additional attack against them. They have 4+{DC} for the attack.
//	When that player dies this turn, they give the active player the item they would destroy for the death penalty.
var lustForBlood = engine.CardDef{
	Ref:     "lust_for_blood",
	Kind:    engine.EventCard,
	Copies:  1,
	Rewards: []engine.Reward{{Kind: engine.RewardTreasure, Amount: 1}},
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "Choose a player. You must make an additional attack against them. They have 4+ DC for the attack. When that player dies this turn, they give the active player the item they would destroy for the death penalty.",
			Trigger: engine.WhenThisEntersPlay(),
			Effects: []engine.Effect{engine.Ask(func(c *engine.Ctx, a []int) {
				others := othersInOrder(c)
				if a[0] < 0 || a[0] >= len(others) {
					return
				}
				p := others[a[0]]
				c.G.Boosts = append(c.G.Boosts, engine.Boost{Stat: engine.StatPlayerAttackDC, Player: p, Amount: 4})
				c.G.ForceAttacks(1, false)
				c.G.Turn.MustAttack = c.G.Players[p].Character
				c.G.GivePenaltyItem(p, c.G.Turn.Active)
			}, engine.Question{Text: "Which player must be attacked?", Options: func(c *engine.Ctx, _ []int) []string {
				return playerLabels(c.G, othersInOrder(c))
			}})},
		},
	},
}
