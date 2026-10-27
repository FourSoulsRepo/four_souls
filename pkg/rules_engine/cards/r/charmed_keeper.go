package r

import (
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/cardkit"
)

// Charmed Keeper (Holy/Charmed Monster Card)
//
//	When this dies, after gaining rewards, choose another player. They roll and gain ¢ equal to the result.
//	Each time this deals combat damage to a player, they lose 2¢.
var charmedKeeper = engine.CardDef{
	Ref:     "charmed_keeper",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      2,
	DC:      3,
	ATK:     2,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Roll: true}},
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "When this dies, after gaining rewards, choose another player. They roll and gain ¢ equal to the result.",
			Trigger: engine.AfterThisRewards(),
			Effects: []engine.Effect{engine.AskWithRoll(func(c *engine.Ctx, a []int) {
				if ps := othersOfActive(c.G); a[0] >= 0 && a[0] < len(ps) {
					them := *c
					them.Controller = ps[a[0]]
					them.Do(engine.Roll(centsEqualToRoll()))
				}
			}, centsEqualToRoll(), engine.Question{Text: cardkit.AskChooseAnotherPlayer, Options: func(c *engine.Ctx, _ []int) []engine.Msg {
				return playerLabels(c.G, othersOfActive(c.G))
			}})},
		},
		{
			Kind:    engine.Triggered,
			Text:    "Each time this deals combat damage to a player, they lose 2¢.",
			Takes:   engine.TakesEventPlayer,
			Trigger: cardkit.WhenThisDealsDamage(true),
			Effects: []engine.Effect{cardkit.ForRoller(engine.LoseCents(2))},
		},
	},
}
