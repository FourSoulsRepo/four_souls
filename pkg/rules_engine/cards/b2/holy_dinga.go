package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Holy Dinga (Holy/Charmed Monster Card)
//
//	Each time a player rolls a ❻, they heal 1{HP}.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
// Not generated: reward "Roll- Gain x" Coin.
var holyDinga = engine.CardDef{
	Ref:    "holy_dinga",
	Kind:   engine.MonsterCard,
	Copies: 1,
	HP:     3,
	DC:     3,
	ATK:    1,
}
