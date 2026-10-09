package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Polydactyly (Passive Treasure Card)
//
//	You may play an additional loot card on your turn.
//	You have +1{ATK} for your first attack roll each turn.
var polydactyly = engine.CardDef{
	Ref:     "polydactyly",
	Kind:    engine.TreasureCard,
	Copies:  1,
	Statics: []engine.Static{engine.YouHave(engine.StatLootPlays, 1), firstAttackRollATK},
}
