package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// 1-Up (Passive Treasure Card)
//
//	When you would die on your turn, destroy this. If you do, prevent death, heal to full {HP}, and cancel your attack. You may attack an additional time this turn.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var card1Up = engine.CardDef{
	Ref:    "1_up",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
