package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Ball Of Bandages (Passive Treasure Card)
//
//	You take no combat damage on attack rolls of 3.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var ballOfBandages = engine.CardDef{
	Ref:    "ball_of_bandages",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
