package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// XX. Judgement (Wildcard Card)
//
//	Choose the player with the most souls or tied for the most. That player destroys a soul they control.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var xxJudgement = engine.CardDef{
	Ref:    "xx_judgement",
	Kind:   engine.LootCard,
	Copies: 1,
}
