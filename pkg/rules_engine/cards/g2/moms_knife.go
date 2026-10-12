package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// Mom’s Knife (Passive Treasure Card)
//
//	{ATK}
//	Each time you deal combat damage to a monster on an attack roll of 6, deal 1 damage to each other monster.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var momsKnife = engine.CardDef{
	Ref:    "moms_knife",
	Kind:   engine.TreasureCard,
	Copies: 1,
}
