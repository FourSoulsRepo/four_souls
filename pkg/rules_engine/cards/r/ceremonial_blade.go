package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Ceremonial Blade (Eternal Treasure Card)
//
//	{Tap Effect}Play an additional loot card this turn.
//	At the end of your turn, loot 1 and discard a loot card.
//	Each time you kill a monster, loot 1.
//	-Eternal- This can't be destroyed or put into discard.
var ceremonialBlade = engine.CardDef{
	Ref:     "ceremonial_blade",
	Kind:    engine.TreasureCard,
	Copies:  1,
	Eternal: true,
	Outside: true,
	Tap:     true,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Activated,
			Text:    "↷: Play an additional loot card this turn.",
			Costs:   []engine.Cost{engine.Tap()},
			Effects: []engine.Effect{engine.AddLootPlays(1)},
		},
		{
			Kind:    engine.Triggered,
			Text:    "At the end of your turn, loot 1 and discard a loot card.",
			Trigger: engine.AtEndOfYourTurn(),
			Effects: []engine.Effect{engine.Loot(1), discardOne},
		},
		{
			Kind: engine.Triggered,
			Text: "Each time you kill a monster, loot 1.",
			Trigger: engine.Trigger{On: engine.EvDied, Match: func(g *engine.Game, self engine.ObjectID, e engine.Event) bool {
				return youKilledMonster(g, g.Object(self).Controller, e)
			}},
			Effects: []engine.Effect{engine.Loot(1)},
		},
	},
}
