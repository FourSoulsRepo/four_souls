package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Cain (Character Card)
//
//	{Tap Effect}Play an additional loot card this turn.
//	If you control this as the game starts, you go first.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var cain = engine.CardDef{
	Ref:          "cain",
	Kind:         engine.CharacterCard,
	Copies:       1,
	HP:           2,
	ATK:          1,
	StartingItem: "sleight_of_hand",
	Tap:          true,
}
