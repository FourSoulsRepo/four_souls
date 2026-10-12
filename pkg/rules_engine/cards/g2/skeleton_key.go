package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// Skeleton Key (Passive Treasure Card)
//
//	You may attack up to 99 times during your turn.
var skeletonKey = engine.CardDef{
	Ref:     "skeleton_key",
	Kind:    engine.TreasureCard,
	Copies:  1,
	Statics: []engine.Static{engine.YouHave(engine.StatAttacks, 98)}, // up to 99 attacks
}
