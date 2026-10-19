package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Conjoined Twin (Room Card)
//
//	Each time a monster dies, the player who killed it kills another player.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var conjoinedTwin = engine.CardDef{
	Ref:       "conjoined_twin",
	Kind:      engine.RoomCard,
	Copies:    1,
	ThreePlus: true,
}
