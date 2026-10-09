package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// The Chest (Soul Treasure Card)
//
//	if this would be destroyed, it becomes a soul instead.
var theChest = engine.CardDef{
	Ref:               "the_chest",
	Kind:              engine.TreasureCard,
	Copies:            1,
	Soul:              1,
	SoulWhenDestroyed: true,
}
