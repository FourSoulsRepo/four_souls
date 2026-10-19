package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Laser Eye! (Room Card)
//
//	Each time a player loots, they take 1 damage.
var laserEye = engine.CardDef{
	Ref:    "laser_eye",
	Kind:   engine.RoomCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "Each time a player loots, they take 1 damage.",
		Trigger: engine.Trigger{On: engine.EvLootDone, Match: anyPlayer},
		Effects: []engine.Effect{theyTake(1)},
	}},
}
