package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// ? Card (Wildcard Card)
//
//	As you play this, choose an item. This copies one of that item's ↷ abilities.
var questionmarkCard = engine.CardDef{
	Ref:             "questionmark_card",
	Kind:            engine.LootCard,
	Copies:          1,
	CopiesTapOnPlay: true, // "As you play this, choose an item. This copies one of that item's ↷ abilities."
}
