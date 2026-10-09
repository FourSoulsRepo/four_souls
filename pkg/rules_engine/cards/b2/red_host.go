package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Red Host (Basic Monster Card)
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var redHost = engine.CardDef{
	Ref:     "red_host",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      2,
	DC:      3,
	ATK:     2,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 5}},
}
