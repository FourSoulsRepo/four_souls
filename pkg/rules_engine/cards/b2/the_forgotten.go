package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// The Forgotten (Character Card)
//
//	{Tap Effect}Play an additional loot card this turn.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var theForgotten = engine.CardDef{
	Ref:          "the_forgotten",
	Kind:         engine.CharacterCard,
	Copies:       1,
	HP:           2,
	ATK:          1,
	StartingItem: "the_bone",
	Tap:          true,
}
