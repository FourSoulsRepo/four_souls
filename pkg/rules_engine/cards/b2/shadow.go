package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Shadow (Passive Treasure Card)
//
//	If another player would pay the death penalty, you choose what item they would destroy and you gain any loot cards and ¢ they would lose.
var shadow = engine.CardDef{
	Ref:            "shadow",
	Kind:           engine.TreasureCard,
	Copies:         1,
	TakesPenalties: true,
}
