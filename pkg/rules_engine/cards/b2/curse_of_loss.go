package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Curse Of Loss (Curse Card)
//
//	{Curse Effect}When you die, destroy a soul you control.
//	-Curse- When this enters play, give this to a player. When they die, they put this into discard.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var curseOfLoss = engine.CardDef{
	Ref:    "curse_of_loss",
	Kind:   engine.EventCard,
	Copies: 1,
}
