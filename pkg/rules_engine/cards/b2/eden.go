package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Eden (Character Card)
//
//	{Tap Effect}Play an additional loot card this turn.
//	When you start the game, look at the top 3 cards of the treasure deck and choose one. It becomes your starting item and gains eternal. Put the rest on the bottom of the treasure deck.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var eden = engine.CardDef{
	Ref:    "eden",
	Kind:   engine.CharacterCard,
	Copies: 1,
	HP:     2,
	ATK:    1,
	Tap:    true,
}
