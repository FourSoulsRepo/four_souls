package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// Ragling (Basic Monster Card)
//
//	When this dies, after gaining rewards, the active player rolls-
//	1 or 6: Put this on top of the Monster Deck.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var ragling = engine.CardDef{
	Ref:     "ragling",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      2,
	DC:      3,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 1}},
}
