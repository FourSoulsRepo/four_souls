package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Experimental Treatment (Soul Treasure Card)
//
//	{Paid Effect}Destroy this: roll-
//	1: You die.
//	2: Loot 3.
//	3: Gain 9¢.
//	4: Gain +1 treasure.
//	5: Gain +2 treasure.
//	6: This becomes a soul. Gain it.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var experimentalTreatment = engine.CardDef{
	Ref:    "experimental_treatment",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Soul:   1,
}
