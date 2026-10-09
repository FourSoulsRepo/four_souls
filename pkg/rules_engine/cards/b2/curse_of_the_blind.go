package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Curse Of The Blind (Curse Card)
//
//	{Curse Effect}Monsters have +1{DC} on your turn.
//	-Curse- When this enters play, give this to a player. When they die, they put this into discard.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var curseOfTheBlind = engine.CardDef{
	Ref:    "curse_of_the_blind",
	Kind:   engine.EventCard,
	Copies: 1,
}
