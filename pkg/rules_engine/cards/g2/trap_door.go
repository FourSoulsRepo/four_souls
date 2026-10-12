package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// Trap Door! (Good Event Card)
//
//	Roll-
//	1-3: Look at the top 3 cards of the monster deck and put them back in any order. You may attack an additional time this turn.
//	4-6: Take 1 damage.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var trapDoor = engine.CardDef{
	Ref:    "trap_door",
	Kind:   engine.EventCard,
	Copies: 1,
}
