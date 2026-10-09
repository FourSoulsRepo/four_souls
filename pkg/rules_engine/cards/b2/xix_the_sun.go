package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// XIX. The Sun (Wildcard Card)
//
//	Put this on the bottom of the loot deck. If you do, take an extra turn after this one if it's your turn.
var xixTheSun = engine.CardDef{
	Ref:    "xix_the_sun",
	Kind:   engine.LootCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.LootAbility,
			Text:    "Put this on the bottom of the loot deck. If you do, take an extra turn after this one if it's your turn.",
			Effects: []engine.Effect{engine.ThisToLootBottomExtraTurn()},
		},
	},
}
