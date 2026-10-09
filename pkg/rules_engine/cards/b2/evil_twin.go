package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Evil Twin (Basic Monster Card)
//
//	Damage dealt to this is also dealt to the player to the active player's left.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var evilTwin = engine.CardDef{
	Ref:     "evil_twin",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      3,
	DC:      5,
	ATK:     2,
	Rewards: []engine.Reward{{Kind: engine.RewardTreasure, Amount: 1}},
}
