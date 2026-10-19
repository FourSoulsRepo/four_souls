package r

import engine "github.com/FourSoulsRepo/rules_engine"

// The Enigma/amginE ehT (Character Card)
//
//	{Tap Effect}Play an additional loot card this turn.
//	amginE ehT: Do not start the game with this side face up.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
// Not generated: HP is "2/1".
// Not generated: ATK is "1/2".
var theEnigma = engine.CardDef{
	Ref:          "the_enigma",
	Kind:         engine.CharacterCard,
	Copies:       1,
	StartingItem: "flip",
	Tap:          true,
}
