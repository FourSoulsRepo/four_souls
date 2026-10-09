package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Remote Detonator (Active Treasure Card)
//
//	{Tap Effect}Each player votes on an item in play. Destroy the item with the most votes. If there is a tie, nothing happens.
var remoteDetonator = engine.CardDef{
	Ref:    "remote_detonator",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Activated,
			Text:    "↷: Each player votes on an item in play. Destroy the item with the most votes. If there is a tie, nothing happens.",
			Costs:   []engine.Cost{engine.Tap()},
			Effects: []engine.Effect{engine.Ask(detonate, voteQuestions()...)},
		},
	},
}

// itemsInPlay lists players' items and shop items, the vote options.
func itemsInPlay(g *engine.Game) []engine.ObjectID { return append(allItems(g), g.ShopItems()...) }

// voteQuestions: each player votes, in turn order from the controller
// (R-MECH-49).
func voteQuestions() []engine.Question {
	var qs []engine.Question
	for k := range 4 {
		qs = append(qs, engine.Question{
			Text: "Vote for an item to destroy.",
			Player: func(c *engine.Ctx, _ []int) engine.PlayerID {
				return engine.PlayerID((int(c.Controller) + k) % len(c.G.Players))
			},
			Options: func(c *engine.Ctx, _ []int) []string {
				if k >= len(c.G.Players) {
					return nil
				}
				var out []string
				for _, id := range itemsInPlay(c.G) {
					out = append(out, string(c.G.Object(id).Card))
				}
				return out
			},
		})
	}
	return qs
}

func detonate(c *engine.Ctx, a []int) {
	items := itemsInPlay(c.G)
	votes := make([]int, len(items))
	for _, v := range a {
		if v >= 0 {
			votes[v]++
		}
	}
	best, tie := -1, false
	for i, n := range votes {
		switch {
		case best < 0 || n > votes[best]:
			best, tie = i, false
		case n == votes[best]:
			tie = true
		}
	}
	if best < 0 || tie || votes[best] == 0 {
		return
	}
	id := items[best]
	if o := c.G.Object(id); o.Controller == engine.NoPlayer {
		c.G.DiscardShopItem(id)
	} else {
		c.G.DestroyObject(o.Controller, id)
	}
}
