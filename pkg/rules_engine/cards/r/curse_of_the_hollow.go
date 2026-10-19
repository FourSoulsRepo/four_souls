package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Curse Of The Hollow (Curse Card)
//
//	At the end of your turn, discard x loot cards and lose x cents, where x is the number of souls you control.
//	-Curse- When this enters play, give this to a player. When they die, they put this into discard.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var curseOfTheHollow = engine.CardDef{
	Ref:    "curse_of_the_hollow",
	Kind:   engine.EventCard,
	Copies: 1,
	Curse:  true,
}
