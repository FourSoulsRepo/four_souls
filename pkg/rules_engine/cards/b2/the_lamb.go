package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// The Lamb (Epic Boss Card)
//
//	When this dies, the active player may choose another player. They give you a soul they control.
var theLamb = engine.CardDef{
	Ref:     "the_lamb",
	Kind:    engine.MonsterCard,
	Copies:  1,
	Soul:    2,
	HP:      6,
	DC:      3,
	ATK:     6,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 3}},
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "When this dies, the active player may choose another player. They give you a soul they control.",
			Trigger: engine.WhenThisDies(),
			Effects: []engine.Effect{engine.Ask(lamb, lambPlayer, lambSoul)},
		},
	},
}

var lambPlayer = engine.Question{Text: "Who gives you a soul?", Options: func(c *engine.Ctx, _ []int) []string {
	var out []string
	for _, p := range otherPlayers(c) {
		out = append(out, playerLabel(c.G, p))
	}
	return append(out, "nobody")
}}

var lambSoul = engine.Question{
	Text:   "Give which soul?",
	Player: func(c *engine.Ctx, a []int) engine.PlayerID { return otherPlayers(c)[a[0]] },
	Options: func(c *engine.Ctx, a []int) []string {
		if a[0] < 0 || a[0] >= len(otherPlayers(c)) {
			return nil
		}
		var out []string
		for _, id := range souls(c.G, otherPlayers(c)[a[0]]) {
			out = append(out, string(c.G.Object(id).Card))
		}
		return out
	},
}

func lamb(c *engine.Ctx, a []int) {
	if a[1] >= 0 {
		c.G.GainControl(c.Controller, souls(c.G, otherPlayers(c)[a[0]])[a[1]])
	}
}
