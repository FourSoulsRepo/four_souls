package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Devil Beggar (Room Card)
//
//	{Paid Effect}Pay 1{HP}: Roll-
//	1-2: Loot 2.
//	3-4: Loot 1.
//	5-6: Take 1 damage.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var devilBeggar = engine.CardDef{
	Ref:    "devil_beggar",
	Kind:   engine.RoomCard,
	Copies: 1,
}
