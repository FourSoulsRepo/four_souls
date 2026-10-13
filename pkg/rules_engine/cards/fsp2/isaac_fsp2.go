package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Isaac! (Epic Boss Card)
//
//	Each time this takes damage, the active player chooses a living player. This deals 1 damage to that player.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var isaacFsp2 = engine.CardDef{
	Ref:       "isaac_fsp2",
	Kind:      engine.MonsterCard,
	Copies:    1,
	HP:        7,
	DC:        3,
	ATK:       1,
	Soul:      2,
	Rewards:   []engine.Reward{{Kind: engine.RewardTreasure, Amount: 1}},
	ThreePlus: true,
}
