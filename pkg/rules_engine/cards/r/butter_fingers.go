package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Butter Fingers (Room Card)
//
//	At the end of each turn, the active player discards a loot card.
var butterFingers = engine.CardDef{
	Ref:    "butter_fingers",
	Kind:   engine.RoomCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "At the end of each turn, the active player discards a loot card.",
		Trigger: atEndOfTheTurn(),
		Effects: []engine.Effect{discardN(1)},
	}},
}
