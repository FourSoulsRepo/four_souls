package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Angelic Intervention (Room Card)
//
//	Each time the active player attacks the top of the monster deck, after putting it in a monster slot, they may cancel their attack.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var angelicIntervention = engine.CardDef{
	Ref:    "angelic_intervention",
	Kind:   engine.RoomCard,
	Copies: 1,
}
