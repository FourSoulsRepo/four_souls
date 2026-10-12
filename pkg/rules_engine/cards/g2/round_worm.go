package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// Round Worm (Basic Monster Card)
//
//	Each time this deals Combat Damage, the attacking player has +1 to their next attack roll this turn.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var roundWorm = engine.CardDef{
	Ref:     "round_worm",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      1,
	DC:      5,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 4}},
}
