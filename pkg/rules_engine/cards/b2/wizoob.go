package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Wizoob (Basic Monster Card)
//
//	When this dies, the active player chooses a player. That player destroys a soul they control.
var wizoob = engine.CardDef{
	Ref:     "wizoob",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      3,
	DC:      5,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 3}},
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "When this dies, the active player chooses a player. That player destroys a soul they control.",
			Trigger: engine.WhenThisDies(),
			Effects: []engine.Effect{engine.Ask(wizoobDestroy, pickAPlayer("Who destroys a soul?"), wizoobSoul)},
		},
	},
}

var wizoobSoul = engine.Question{
	Text:   "Destroy which of your souls?",
	Player: func(c *engine.Ctx, a []int) engine.PlayerID { return livingPlayers(c.G, engine.NoPlayer)[a[0]] },
	Options: func(c *engine.Ctx, a []int) []string {
		if a[0] < 0 {
			return nil
		}
		var out []string
		for _, id := range souls(c.G, livingPlayers(c.G, engine.NoPlayer)[a[0]]) {
			out = append(out, string(c.G.Object(id).Card))
		}
		return out
	},
}

func wizoobDestroy(c *engine.Ctx, a []int) {
	if a[0] < 0 || a[1] < 0 {
		return
	}
	p := livingPlayers(c.G, engine.NoPlayer)[a[0]]
	c.G.DestroyObject(p, souls(c.G, p)[a[1]])
}
