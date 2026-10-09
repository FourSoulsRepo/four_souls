package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// XV. The Devil (Wildcard Card)
//
//	Destroy an item you control. If you do, steal a non-eternal item from a player or from the shop.
var xvTheDevil = engine.CardDef{
	Ref:    "xv_the_devil",
	Kind:   engine.LootCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.LootAbility,
			Text:    "Destroy an item you control. If you do, steal a non-eternal item from a player or from the shop.",
			Targets: []engine.TargetSpec{engine.Choose(engine.TargetYourItem)},
			Effects: []engine.Effect{engine.Ask(devil, devilQuestion)},
		},
	},
}

// stealable lists non-eternal items of other players, then shop items.
func stealable(c *engine.Ctx) []engine.ObjectID {
	var out []engine.ObjectID
	for _, pl := range c.G.Players {
		if pl.ID == c.Controller {
			continue
		}
		for _, id := range pl.InPlay {
			if o := c.G.Object(id); o.Role == engine.RoleItem && !c.G.Eternal(id) {
				out = append(out, id)
			}
		}
	}
	return append(out, c.G.ShopItems()...)
}

var devilQuestion = engine.Question{
	Text: "Steal which item?",
	Options: func(c *engine.Ctx, _ []int) []string {
		var out []string
		for _, id := range stealable(c) {
			out = append(out, string(c.G.Object(id).Card))
		}
		return out
	},
}

func devil(c *engine.Ctx, a []int) {
	items := stealable(c)
	if c.G.DestroyObject(c.Controller, c.Targets[0].Object) && a[0] >= 0 {
		c.G.GainControl(c.Controller, items[a[0]])
	}
}
