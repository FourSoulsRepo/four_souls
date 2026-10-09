package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Two Of Clubs (Active Treasure Card)
//
//	{Tap Effect}Choose a player. Till end of turn, if they would loot any number of loot cards, they loot double that number instead.
var twoOfClubs = engine.CardDef{
	Ref:    "two_of_clubs",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Activated,
			Text:    "↷: Choose a player. Till end of turn, if they would loot any number of loot cards, they loot double that number instead.",
			Costs:   []engine.Cost{engine.Tap()},
			Targets: []engine.TargetSpec{engine.Choose(engine.TargetPlayer)},
			Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
				c.G.Boosts = append(c.G.Boosts, engine.Boost{Stat: engine.StatLootDouble, Player: c.Targets[0].Player, Amount: 1})
			})},
		},
	},
}
