package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Golden Troll Bomb! (Good Event Card)
//
//	Choose a player. They take 2 damage. Put this in the monster deck 4 cards from the top.
var goldenTrollBomb = engine.CardDef{
	Ref:       "golden_troll_bomb",
	Kind:      engine.EventCard,
	Copies:    1,
	ThreePlus: true,
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "Choose a player. They take 2 damage. Put this in the monster deck 4 cards from the top.",
		Trigger: engine.WhenThisEntersPlay(),
		Effects: []engine.Effect{damageAChosenPlayer(2), backIntoDeck(4)},
	}},
}
