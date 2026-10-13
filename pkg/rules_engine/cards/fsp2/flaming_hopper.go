package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Flaming Hopper (Basic Monster Card)
//
//	This takes no combat damage on attack rolls of 6.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var flamingHopper = engine.CardDef{
	Ref:     "flaming_hopper",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      1,
	DC:      4,
	ATK:     2,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 4}},
}
