package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Lust For Blood! (Good Event Card)
//
//	Choose a player. You must make an additional attack against them. They have 4+{DC} for the attack.
//	When that player dies this turn, they give the active player the item they would destroy for the death penalty.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var lustForBlood = engine.CardDef{
	Ref:    "lust_for_blood",
	Kind:   engine.EventCard,
	Copies: 1,
}
