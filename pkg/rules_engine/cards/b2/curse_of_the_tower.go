package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Curse Of The Tower (Passive Treasure Card)
//
//	Each time you take damage, roll-
//	1-3: Each other player takes 1 damage.
//	4-6: Deal 1 damage to a monster.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var curseOfTheTower = engine.CardDef{
	Ref:    "curse_of_the_tower",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
