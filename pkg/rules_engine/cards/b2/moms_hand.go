package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Mom’s Hand (Basic Monster Card)
//
//	When the attacking player rolls an attack roll of 6, cancel everything that hasn't resolved and end the turn.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var momsHand = engine.CardDef{
	Ref:     "moms_hand",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      2,
	DC:      4,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 4}},
}
