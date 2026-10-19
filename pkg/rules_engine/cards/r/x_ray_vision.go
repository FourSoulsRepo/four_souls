package r

import engine "github.com/FourSoulsRepo/rules_engine"

// X-Ray Vision (Active Treasure Card)
//
//	{Tap Effect}Look at a player's hand and the top card of a deck. You may put the top card of a deck into discard.
//	At the end of your turn, recharge this.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var xRayVision = engine.CardDef{
	Ref:    "x_ray_vision",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
}
