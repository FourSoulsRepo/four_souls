package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Rage Creep (Basic Monster Card)
//
//	Damage this deals to the active player is also dealt to the player to their left.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var rageCreep = engine.CardDef{
	Ref:     "rage_creep",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      1,
	DC:      5,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 2}},
}
