package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Mom’s Eye (Basic Monster Card)
//
//	When this dies, the active player may look at a player's hand.
var momsEye = engine.CardDef{
	Ref:     "moms_eye",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      1,
	DC:      4,
	ATK:     2,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 1}},
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "When this dies, the active player may look at a player's hand.",
			Trigger: engine.WhenThisDies(),
			Effects: []engine.Effect{engine.Ask(func(c *engine.Ctx, a []int) {
				if ps := livingPlayers(c.G, engine.NoPlayer); a[0] >= 0 && a[0] < len(ps) {
					c.G.LookAt(c.Controller, c.G.Players[ps[a[0]]].Hand...)
				}
			}, pickAPlayer("Look at whose hand?"))},
		},
	},
}
