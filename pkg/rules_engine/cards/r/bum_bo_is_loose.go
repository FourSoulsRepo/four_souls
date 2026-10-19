package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Bum-Bo Is Loose! (Room Card)
//
//	At the end of the turn, the active player loses 2¢.
var bumBoIsLoose = engine.CardDef{
	Ref:    "bum_bo_is_loose",
	Kind:   engine.RoomCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "At the end of the turn, the active player loses 2¢.",
		Trigger: atEndOfTheTurn(),
		Effects: []engine.Effect{engine.LoseCents(2)},
	}},
}
