package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Placebo (Active Treasure Card)
//
//	{Tap Effect}This copies a ↷ ability of a non-eternal item.
var placebo = engine.CardDef{
	Ref:                "placebo",
	Kind:               engine.TreasureCard,
	Copies:             1,
	Tap:                true,
	CopiesTapAbilities: true,
}
