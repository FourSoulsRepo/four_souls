package r

import (
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/cardkit"
)

// Planetarium (Room Card)
//
//	At the start of the turn, the active player may gain +1 treasure.
var planetarium = engine.CardDef{
	Ref:    "planetarium",
	Kind:   engine.RoomCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "At the start of the turn, the active player may gain +1 treasure.",
		Trigger: atStartOfTheTurn(),
		Effects: []engine.Effect{cardkit.May("Gain +1 treasure?", engine.GainTreasure(1))},
	}},
}
