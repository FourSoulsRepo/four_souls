package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Challenge Room (Room Card)
//
//	Each time a monster dies, the active player rolls-
//	4-6: They must make an additional attack this turn.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var challengeRoom = engine.CardDef{
	Ref:    "challenge_room",
	Kind:   engine.RoomCard,
	Copies: 1,
}
