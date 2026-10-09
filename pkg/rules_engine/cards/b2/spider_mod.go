package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Spider Mod (Passive Treasure Card)
//
//	Each time a player rolls a ❺, you may put a monster not being attacked into discard and replace it with the top card of the monster deck.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var spiderMod = engine.CardDef{
	Ref:    "spider_mod",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
