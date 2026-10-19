package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Ultra Flesh Kid! (Soul Treasure Card)
//
//	Each time a monster or player dies, put a counter on this.
//	At the end of your turn, if this has 12+ counters, it becomes a soul and loses all abilities.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var ultraFleshKid = engine.CardDef{
	Ref:    "ultra_flesh_kid",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Soul:   1,
}
