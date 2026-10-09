package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Two Of Clubs (Active Treasure Card)
//
//	{Tap Effect}Choose a player. Till end of turn, if they would loot any number of loot cards, they loot double that number instead.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var twoOfClubs = engine.CardDef{
	Ref:    "two_of_clubs",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
}
