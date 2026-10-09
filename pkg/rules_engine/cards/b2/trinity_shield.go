package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Trinity Shield (Passive Treasure Card)
//
//	Other players can't play loot cards or activate items on your turn.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var trinityShield = engine.CardDef{
	Ref:    "trinity_shield",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
