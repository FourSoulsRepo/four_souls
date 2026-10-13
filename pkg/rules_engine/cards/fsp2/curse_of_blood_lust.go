package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Curse Of Blood Lust (Curse Card)
//
//	{Curse Effect}You must attack on your turn if able.
//	-Curse- When this enters play, give this to a player. When they die, they put this into discard.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var curseOfBloodLust = engine.CardDef{
	Ref:    "curse_of_blood_lust",
	Kind:   engine.EventCard,
	Copies: 1,
	Curse:  true,
}
