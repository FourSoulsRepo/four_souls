package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Curse Of Impulse (Curse Card)
//
//	{Curse Effect}At the end of your turn, deactivate each item you control and your character.
//	-Curse- When this enters play, give this to a player. When they die, they put this into discard.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var curseOfImpulse = engine.CardDef{
	Ref:    "curse_of_impulse",
	Kind:   engine.EventCard,
	Copies: 1,
	Curse:  true,
}
