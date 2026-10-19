package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Cursed Soul (Soul Treasure Card)
//
//	At the start of your turn, take 1 damage and put a counter on this. Then, if this has 6+ counters, it becomes a soul and loses all abilities.
//	At the end of your turn, give this to another player.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var cursedSoul = engine.CardDef{
	Ref:    "cursed_soul",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Soul:   1,
}
