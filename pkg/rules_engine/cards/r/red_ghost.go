package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Red Ghost (Basic Monster Card)
//
//	Each time this takes damage, the active player chooses another player. That player takes 1 damage.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var redGhost = engine.CardDef{
	Ref:       "red_ghost",
	Kind:      engine.MonsterCard,
	Copies:    1,
	HP:        3,
	DC:        4,
	ATK:       1,
	Rewards:   []engine.Reward{{Kind: engine.RewardLoot, Amount: 3}},
	ThreePlus: true,
}
