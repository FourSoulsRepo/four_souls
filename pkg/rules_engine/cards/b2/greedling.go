package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Greedling (Basic Monster Card)
//
//	When this dies, the active player chooses a player. They lose 7¢.
var greedling = engine.CardDef{
	Ref:     "greedling",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      2,
	DC:      5,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 7}},
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "When this dies, the active player chooses a player. They lose 7¢.",
			Trigger: engine.WhenThisDies(),
			Effects: []engine.Effect{engine.Ask(func(c *engine.Ctx, a []int) {
				if a[0] >= 0 {
					p := livingPlayers(c.G, engine.NoPlayer)[a[0]]
					c.G.LoseCentsNow(p, 7)
				}
			}, pickAPlayer("Who loses 7¢?"))},
		},
	},
}
