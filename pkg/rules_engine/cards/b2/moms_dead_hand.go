package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Mom’s Dead Hand (Basic Monster Card)
//
//	When this dies, the active player may steal a non-eternal item another player controls.
var momsDeadHand = engine.CardDef{
	Ref:     "moms_dead_hand",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      2,
	DC:      5,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardTreasure, Amount: 1}, {Kind: engine.RewardCents, Amount: 4}},
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "When this dies, the active player may steal a non-eternal item another player controls.",
			Trigger: engine.WhenThisDies(),
			Effects: []engine.Effect{engine.Ask(func(c *engine.Ctx, a []int) {
				if ids := othersItems(c); a[0] >= 0 && a[0] < len(ids) {
					c.G.GainControl(c.Controller, ids[a[0]])
				}
			}, engine.Question{Text: "Steal which item?", Options: func(c *engine.Ctx, _ []int) []string {
				var out []string
				for _, id := range othersItems(c) {
					out = append(out, string(c.G.Object(id).Card))
				}
				if out == nil {
					return nil
				}
				return append(out, "none")
			}})},
		},
	},
}

// othersItems lists non-eternal items of the other players.
func othersItems(c *engine.Ctx) []engine.ObjectID {
	var out []engine.ObjectID
	for _, p := range otherPlayers(c) {
		out = append(out, nonEternalItems(c.G, p)...)
	}
	return out
}
