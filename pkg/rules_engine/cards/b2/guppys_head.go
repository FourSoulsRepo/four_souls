package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Guppy’s Head (Active Treasure Card)
//
//	{Tap Effect}Choose a player. That player gives you a loot card.
//	-Guppy- The first player to control 2 or more Guppy items gains the Soul of Guppy.
var guppysHead = engine.CardDef{
	Ref:    "guppys_head",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Activated,
			Text:    "↷: Choose a player. That player gives you a loot card.",
			Costs:   []engine.Cost{engine.Tap()},
			Targets: []engine.TargetSpec{engine.Choose(engine.TargetOtherPlayer)},
			Effects: []engine.Effect{engine.Ask(giveMeACard, giveQuestion)},
		},
	},
}

// giveQuestion is answered by the chosen player: which card they give.
var giveQuestion = engine.Question{
	Text:   "Give which loot card?",
	Player: func(c *engine.Ctx, _ []int) engine.PlayerID { return c.Targets[0].Player },
	Options: func(c *engine.Ctx, _ []int) []string {
		var out []string
		for _, id := range c.G.Players[c.Targets[0].Player].Hand {
			out = append(out, string(c.G.Object(id).Card))
		}
		return out
	},
}

func giveMeACard(c *engine.Ctx, a []int) {
	if a[0] >= 0 {
		from := c.Targets[0].Player
		c.G.GiveHandCard(from, c.Controller, c.G.Players[from].Hand[a[0]])
	}
}
