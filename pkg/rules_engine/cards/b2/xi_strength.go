package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// XI. Strength (Wildcard Card)
//
//	Choose a player.
//	They gain +1{ATK} till end of turn and may attack an additional time this turn.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var xiStrength = engine.CardDef{
	Ref:    "xi_strength",
	Kind:   engine.LootCard,
	Copies: 1,
}
