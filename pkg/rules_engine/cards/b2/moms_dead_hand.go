package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Mom’s Dead Hand (Basic Monster Card)
//
//	When this dies, the active player may steal a non-eternal item another player controls.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var momsDeadHand = engine.CardDef{
	Ref:     "moms_dead_hand",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      2,
	DC:      5,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardTreasure, Amount: 1}, {Kind: engine.RewardCents, Amount: 4}},
}
