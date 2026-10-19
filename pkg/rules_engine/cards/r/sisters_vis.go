package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Sisters Vis (Boss Card)
//
//	When an attack is declared on this, the active player chooses another player. That player must make an attack roll against this after each attack roll the active player makes this attack. If this dies this turn, that player gains +1 treasure.
var sistersVis = engine.CardDef{
	Ref:     "sisters_vis",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      6,
	DC:      4,
	ATK:     2,
	Soul:    1,
	Rewards: []engine.Reward{{Kind: engine.RewardTreasure, Amount: 1}},
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "When an attack is declared on this, the active player chooses another player. That player must make an attack roll against this after each attack roll the active player makes this attack. If this dies this turn, that player gains +1 treasure.",
			Trigger: whenAttackedThis(),
			Effects: []engine.Effect{engine.Ask(func(c *engine.Ctx, a []int) {
				others := othersOfActive(c.G)
				if a[0] < 0 || a[0] >= len(others) {
					return
				}
				c.G.AddAttackHelpers(others[a[0]])
				markHelper(c, others[a[0]])
			}, engine.Question{Text: "Who must attack Sisters Vis too?", Options: func(c *engine.Ctx, _ []int) []string {
				return playerLabels(c.G, othersOfActive(c.G))
			}})},
		},
		{
			Kind:    engine.Triggered,
			Text:    "If this dies this turn, that player gains +1 treasure.",
			Trigger: engine.WhenThisDies(),
			Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
				if p, ok := markedHelper(c); ok {
					them := *c
					them.Controller = p
					them.Do(engine.GainTreasure(1))
				}
			})},
		},
	},
}
