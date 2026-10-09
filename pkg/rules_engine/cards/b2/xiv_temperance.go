package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// XIV. Temperance (Wildcard Card)
//
//	Choose one- Take 1 damage and gain 4¢. Take 2 damage and gain 8¢.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var xivTemperance = engine.CardDef{
	Ref:    "xiv_temperance",
	Kind:   engine.LootCard,
	Copies: 1,
}
