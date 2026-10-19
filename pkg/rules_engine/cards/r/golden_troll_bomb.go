package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Golden Troll Bomb! (Good Event Card)
//
//	Choose a player. They take 2 damage. Put this in the monster deck 4 cards from the top.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var goldenTrollBomb = engine.CardDef{
	Ref:       "golden_troll_bomb",
	Kind:      engine.EventCard,
	Copies:    1,
	ThreePlus: true,
}
