package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Ceremonial Blade (Eternal Treasure Card)
//
//	{Tap Effect}Play an additional loot card this turn.
//	At the end of your turn, loot 1 and discard a loot card.
//	Each time you kill a monster, loot 1.
//	-Eternal- This can't be destroyed or put into discard.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var ceremonialBlade = engine.CardDef{
	Ref:     "ceremonial_blade",
	Kind:    engine.TreasureCard,
	Copies:  1,
	Eternal: true,
	Outside: true,
	Tap:     true,
}
