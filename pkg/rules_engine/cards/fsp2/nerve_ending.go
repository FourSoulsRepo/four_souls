package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Nerve Ending (Basic Monster Card)
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var nerveEnding = engine.CardDef{
	Ref:     "nerve_ending",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      4,
	DC:      2,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 3}},
}
