package r

import (
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/cardkit"
)

// The card's own lines; their words are in en.json (ADR 021).
var (
	planetariumGain1Treasure = engine.NewKey("card.planetarium.gain_1_treasure")
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
		Effects: []engine.Effect{cardkit.May(planetariumGain1Treasure, engine.GainTreasure(1))},
	}},
}
