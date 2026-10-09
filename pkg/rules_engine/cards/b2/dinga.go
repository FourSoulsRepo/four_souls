package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Dinga (Basic Monster Card)
//
//	When this dies on an attack roll of 6, double its rewards.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
// Not generated: reward "Roll- Gain x" Coin.
var dinga = engine.CardDef{
	Ref:    "dinga",
	Kind:   engine.MonsterCard,
	Copies: 1,
	HP:     3,
	DC:     3,
	ATK:    1,
}
