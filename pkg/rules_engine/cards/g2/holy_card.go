package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// Holy Card (Wildcard Card)
//
//	Choose a player. The next time that player would die this turn, prevent death. If it's their turn, cancel everything that hasn't resolved and end it.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var holyCard = engine.CardDef{
	Ref:    "holy_card",
	Kind:   engine.LootCard,
	Copies: 1,
}
