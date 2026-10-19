package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Laser Eye! (Room Card)
//
//	Each time a player loots, they take 1 damage.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var laserEye = engine.CardDef{
	Ref:    "laser_eye",
	Kind:   engine.RoomCard,
	Copies: 1,
}
