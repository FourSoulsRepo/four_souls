package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Null (Basic Monster Card)
//
//	This takes no combat damage on every other attack roll made against it.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var null = engine.CardDef{
	Ref:     "null",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      3,
	DC:      4,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 6}},
}
