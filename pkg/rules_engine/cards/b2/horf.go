package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Horf (Basic Monster Card)
//
//	Combat damage this deals is increased by 1 on attack rolls of 2.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var horf = engine.CardDef{
	Ref:     "horf",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      1,
	DC:      4,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 3}},
}
