package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Mega Troll Bomb! (Bad Event Card)
//
//	Each player takes 2 damage!
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var megaTrollBomb = engine.CardDef{
	Ref:    "mega_troll_bomb",
	Kind:   engine.EventCard,
	Copies: 1,
}
