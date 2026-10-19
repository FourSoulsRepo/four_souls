package r

import engine "github.com/FourSoulsRepo/rules_engine"

// The Mirror (Room Card)
//
//	At the start of the turn, the active player may pay {HP} until they have 1{HP}. If they do, each time a monster dies this turn, they gain +1 treasure.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var theMirror = engine.CardDef{
	Ref:    "the_mirror",
	Kind:   engine.RoomCard,
	Copies: 1,
}
