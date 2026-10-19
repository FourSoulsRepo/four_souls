package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Splash Damage (Room Card)
//
//	Each time a player deals damage to a monster, they deal 1 damage to the player to their left.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var splashDamage = engine.CardDef{
	Ref:       "splash_damage",
	Kind:      engine.RoomCard,
	Copies:    1,
	ThreePlus: true,
}
