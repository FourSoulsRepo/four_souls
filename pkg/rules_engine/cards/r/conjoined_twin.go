package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Conjoined Twin (Room Card)
//
//	Each time a monster dies, the player who killed it kills another player.
var conjoinedTwin = engine.CardDef{
	Ref:       "conjoined_twin",
	Kind:      engine.RoomCard,
	Copies:    1,
	ThreePlus: true,
	Abilities: []engine.Ability{{
		Kind: engine.Triggered,
		Text: "Each time a monster dies, the player who killed it kills another player.",
		Trigger: engine.Trigger{On: engine.EvDied, Match: func(g *engine.Game, _ engine.ObjectID, e engine.Event) bool {
			return e.Player == engine.NoPlayer && g.Object(e.Object).Killer() != engine.NoPlayer
		}},
		Effects: []engine.Effect{engine.Ask(func(c *engine.Ctx, a []int) {
			killer := c.G.Object(c.EventObject).Killer()
			if ps := othersOf(c.G, killer); a[0] >= 0 && a[0] < len(ps) {
				c.G.KillPlayer(ps[a[0]], killer)
			}
		}, engine.Question{
			Text:   "Kill which player?",
			Player: func(c *engine.Ctx, _ []int) engine.PlayerID { return c.G.Object(c.EventObject).Killer() },
			Options: func(c *engine.Ctx, _ []int) []string {
				return playerLabels(c.G, othersOf(c.G, c.G.Object(c.EventObject).Killer()))
			},
		})},
	}},
}
