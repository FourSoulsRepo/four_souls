package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Rainbow Baby (Active Treasure Card)
//
//	{Tap Effect}Choose one- Each player loots 1. Each player gains 2¢. Each player takes 1 damage.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var rainbowBaby = engine.CardDef{
	Ref:    "rainbow_baby",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
}
