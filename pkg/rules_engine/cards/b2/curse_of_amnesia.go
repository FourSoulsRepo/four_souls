package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Curse Of Amnesia (Curse Card)
//
//	{Curse Effect}At the end of your turn, discard 2 loot cards.
//	-Curse- When this enters play, give this to a player. When they die, they put this into discard.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var curseOfAmnesia = engine.CardDef{
	Ref:    "curse_of_amnesia",
	Kind:   engine.EventCard,
	Copies: 1,
}
