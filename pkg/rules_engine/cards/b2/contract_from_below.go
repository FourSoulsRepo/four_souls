package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Contract From Below (Paid Treasure Card)
//
//	{Paid Effect}Destroy 2 items you control:
//	steal a non-eternal item from a player.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var contractFromBelow = engine.CardDef{
	Ref:    "contract_from_below",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
