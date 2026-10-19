package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Shell Game (Room Card)
//
//	{Paid Effect}Pay 3¢: roll-
//	3: Gain 1¢.
//	4: Loot 1.
//	5: Gain 5¢.
//	6: Gain +1 treasure. Put this into discard.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var shellGame = engine.CardDef{
	Ref:    "shell_game",
	Kind:   engine.RoomCard,
	Copies: 1,
}
