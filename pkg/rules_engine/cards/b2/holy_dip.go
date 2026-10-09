package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Holy Dip (Holy/Charmed Monster Card)
//
//	Each time a player rolls a ❶, they gain 1¢.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var holyDip = engine.CardDef{
	Ref:     "holy_dip",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      1,
	DC:      4,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 1}},
}
