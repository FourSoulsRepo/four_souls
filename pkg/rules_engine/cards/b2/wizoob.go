package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Wizoob (Basic Monster Card)
//
//	When this dies, the active player chooses a player. That player destroys a soul they control.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var wizoob = engine.CardDef{
	Ref:     "wizoob",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      3,
	DC:      5,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 3}},
}
