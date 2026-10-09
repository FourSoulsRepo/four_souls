package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Curse Of Pain (Curse Card)
//
//	{Curse Effect}At the start of your turn, take 1 damage.
//	-Curse- When this enters play, give this to a player. When they die, they put this into discard.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var curseOfPain = engine.CardDef{
	Ref:    "curse_of_pain",
	Kind:   engine.EventCard,
	Copies: 1,
}
