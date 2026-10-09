package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Bob’s Brain (Passive Treasure Card)
//
//	Each time you declare an attack, roll-
//	1-2: Deal 1 damage to a monster.
//	3-4: Deal 1 damage to a player.
//	5-6: Take 1 damage.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var bobsBrain = engine.CardDef{
	Ref:    "bobs_brain",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
