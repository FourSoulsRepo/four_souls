package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// Succubus (Paid Treasure Card)
//
//	{Paid Effect}Pay 3¢:
//	You may play up to 2 additional loot cards this turn.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var succubus = engine.CardDef{
	Ref:    "succubus",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
