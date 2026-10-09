package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Eden’s Blessing (Passive Treasure Card)
//
//	At the end of your turn, if you have 0¢, gain 6¢.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var edensBlessing = engine.CardDef{
	Ref:    "edens_blessing",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
