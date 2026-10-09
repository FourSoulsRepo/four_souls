package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// We Need To Go Deeper! (Good Event Card)
//
//	Put any number of non-event monster cards in discard on top of the monster deck.
//	The active player may attack an additional time this turn.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var weNeedToGoDeeper = engine.CardDef{
	Ref:    "we_need_to_go_deeper",
	Kind:   engine.EventCard,
	Copies: 1,
}
