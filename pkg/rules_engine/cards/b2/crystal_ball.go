package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Crystal Ball (Active Treasure Card)
//
//	{Tap Effect}Before a dice is rolled, choose a number. If the next roll is that number, loot 3.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var crystalBall = engine.CardDef{
	Ref:    "crystal_ball",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
}
