package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Curse Of Empathy (Curse Card)
//
//	Each time a monster dies, discard a loot card and lose 1¢.
//	-Curse- When this enters play, give this to a player. When they die, they put this into discard.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var curseOfEmpathy = engine.CardDef{
	Ref:    "curse_of_empathy",
	Kind:   engine.EventCard,
	Copies: 1,
	Curse:  true,
}
