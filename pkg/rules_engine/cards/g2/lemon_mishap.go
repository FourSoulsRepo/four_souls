package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// Lemon Mishap (Active Treasure Card)
//
//	{Tap Effect}Till end of turn, each time a player rolls an attack roll of 1 or 2, deal 1 damage to each monster.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var lemonMishap = engine.CardDef{
	Ref:    "lemon_mishap",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
}
