package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Leaper (Basic Monster Card)
//
//	Combat damage this deals is doubled on attack rolls of 1.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var leaper = engine.CardDef{
	Ref:     "leaper",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      2,
	DC:      4,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 5}},
}
