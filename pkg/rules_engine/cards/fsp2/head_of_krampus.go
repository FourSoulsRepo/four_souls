package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Head Of Krampus (Active Treasure Card)
//
//	{Tap Effect}Roll-
//	1-3: Deal 1 damage to each player.
//	4-6: Deal 1 damage to each monster.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var headOfKrampus = engine.CardDef{
	Ref:    "head_of_krampus",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
}
