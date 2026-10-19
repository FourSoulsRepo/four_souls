package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Social Goals! (Room Card)
//
//	Note each goal as players complete them.
//	This room can't be put into discard till 4 goals are completed.
//	1. Play 5 loot cards.
//	2. Kill 3 monsters.
//	3. Give at least 6¢ to another player at one time.
//	4. Purchase 3 items.
//	5. Roll a 6 three times.
//	When 4 goals are completed, each player gains +2 treasure.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var socialGoals = engine.CardDef{
	Ref:    "social_goals",
	Kind:   engine.RoomCard,
	Copies: 1,
}
