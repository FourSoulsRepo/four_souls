package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// A Sack (Wildcard Card)
//
//	Loot 3.
var aSack = engine.CardDef{
	Ref:       "a_sack",
	Kind:      engine.LootCard,
	Copies:    1,
	Abilities: []engine.Ability{{Kind: engine.LootAbility, Text: "Loot 3.", Effects: []engine.Effect{engine.Loot(3)}}},
}
