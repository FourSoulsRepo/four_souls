package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Bethany (Character Card)
//
//	{Tap Effect}Play an additional loot card this turn.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var bethany = engine.CardDef{
	Ref:          "bethany",
	Kind:         engine.CharacterCard,
	Copies:       1,
	HP:           2,
	ATK:          1,
	StartingItem: "book_of_virtues",
	Tap:          true,
}
