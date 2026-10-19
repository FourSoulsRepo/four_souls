package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Charmed Dip (Holy/Charmed Monster Card)
//
//	When this dies, after gaining rewards, choose another player. They recharge an item they control.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var charmedDip = engine.CardDef{
	Ref:     "charmed_dip",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      1,
	DC:      4,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 2}},
}
