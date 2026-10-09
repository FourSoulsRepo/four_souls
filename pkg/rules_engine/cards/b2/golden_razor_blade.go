package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Golden Razor Blade (Paid Treasure Card)
//
//	{Paid Effect}Pay 5¢:
//	Deal 1 damage to a monster or player.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var goldenRazorBlade = engine.CardDef{
	Ref:    "golden_razor_blade",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
