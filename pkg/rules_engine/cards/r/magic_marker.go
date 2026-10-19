package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Magic Marker (Wildcard Card)
//
//	Change a number in the effect text of a card in play or loot being played by 1 till end of turn. The number can't go below 1 or above 6.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var magicMarker = engine.CardDef{
	Ref:    "magic_marker",
	Kind:   engine.LootCard,
	Copies: 1,
}
