package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Forget Me Now (One-Use Treasure Card)
//
//	{Tap Effect}Destroy this. If you do, each player destroys a soul they control.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var forgetMeNow = engine.CardDef{
	Ref:    "forget_me_now",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
}
