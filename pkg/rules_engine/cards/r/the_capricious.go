package r

import engine "github.com/FourSoulsRepo/rules_engine"

// The Capricious (Character Card)
//
//	{Tap Effect}Choose one- Play an additional loot card this turn. Put a shop item into discard.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var theCapricious = engine.CardDef{
	Ref:          "the_capricious",
	Kind:         engine.CharacterCard,
	Copies:       1,
	HP:           2,
	ATK:          1,
	StartingItem: "glitch",
	Tap:          true,
}
