package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// XX. Judgement (Wildcard Card)
//
//	Choose the player with the most souls or tied for the most. That player destroys a soul they control.
var xxJudgement = engine.CardDef{
	Ref:    "xx_judgement",
	Kind:   engine.LootCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.LootAbility,
			Text:    "Choose the player with the most souls or tied for the most. That player destroys a soul they control.",
			Targets: []engine.TargetSpec{engine.ChooseWhere(engine.TargetPlayer, mostSouls)},
			Effects: []engine.Effect{engine.Ask(judgement, judgementQuestion)},
		},
	},
}

func mostSouls(g *engine.Game, _ engine.ObjectID, c engine.Chosen) bool {
	for _, pl := range g.Players {
		if g.SoulValue(pl.ID) > g.SoulValue(c.Player) {
			return false
		}
	}
	return true
}

func souls(g *engine.Game, p engine.PlayerID) []engine.ObjectID {
	var out []engine.ObjectID
	for _, id := range g.Players[p].InPlay {
		if g.Object(id).Role == engine.RoleSoul {
			out = append(out, id)
		}
	}
	return out
}

// judgementQuestion is answered by the chosen player.
var judgementQuestion = engine.Question{
	Text:   "Destroy which of your souls?",
	Player: func(c *engine.Ctx, _ []int) engine.PlayerID { return c.Targets[0].Player },
	Options: func(c *engine.Ctx, _ []int) []string {
		var out []string
		for _, id := range souls(c.G, c.Targets[0].Player) {
			out = append(out, string(c.G.Object(id).Card))
		}
		return out
	},
}

func judgement(c *engine.Ctx, a []int) {
	if a[0] >= 0 {
		p := c.Targets[0].Player
		c.G.DestroyObject(p, souls(c.G, p)[a[0]])
	}
}
