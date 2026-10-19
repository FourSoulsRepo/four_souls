package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Curse Of The Hunted (Curse Card)
//
//	Monsters have +1{ATK} on your turn.
//	-Curse- When this enters play, give this to a player. When they die, they put this into discard.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var curseOfTheHunted = engine.CardDef{
	Ref:    "curse_of_the_hunted",
	Kind:   engine.EventCard,
	Copies: 1,
	Curse:  true,
}
