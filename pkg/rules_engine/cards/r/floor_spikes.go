package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Floor Spikes (Room Card)
//
//	Each time a player rolls a ❶, they take 1 damage.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var floorSpikes = engine.CardDef{
	Ref:    "floor_spikes",
	Kind:   engine.RoomCard,
	Copies: 1,
}
