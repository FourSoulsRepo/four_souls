package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Envy (Boss Card)
//
//	When this dies, the active player must make an additional attack.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var envy = engine.CardDef{
	Ref:     "envy",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      2,
	DC:      5,
	ATK:     1,
	Soul:    1,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 1}},
}
