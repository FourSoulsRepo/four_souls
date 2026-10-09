package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Sloth (Boss Card)
//
//	When this dies, the player that killed it discards their hand.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var sloth = engine.CardDef{
	Ref:     "sloth",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      3,
	DC:      4,
	ATK:     1,
	Soul:    1,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 1}},
}
