package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Isaac’s Blessing (Room Card)
//
//	Each time a player would roll a ❶ or ❻, they may reroll it.
var isaacsBlessing = engine.CardDef{
	Ref:    "isaacs_blessing",
	Kind:   engine.RoomCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "Each time a player would roll a 1, they may reroll it.",
			Trigger: engine.WhenARollWouldBe(1),
			Effects: []engine.Effect{rollerMayReroll},
		},
		{
			Kind:    engine.Triggered,
			Text:    "Each time a player would roll a 6, they may reroll it.",
			Trigger: engine.WhenARollWouldBe(6),
			Effects: []engine.Effect{rollerMayReroll},
		},
	},
}

// rollerMayReroll: the player who rolled may reroll it.
var rollerMayReroll = engine.Ask(func(c *engine.Ctx, a []int) {
	if a[0] == 0 {
		c.G.Reroll(c.EventStack)
	}
}, engine.Question{
	Text:    "Reroll it?",
	Player:  func(c *engine.Ctx, _ []int) engine.PlayerID { return c.EventPlayer },
	Options: func(*engine.Ctx, []int) []string { return []string{"yes", "No"} },
})
