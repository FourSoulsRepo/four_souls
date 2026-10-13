package fsp2

import engine "github.com/FourSoulsRepo/rules_engine"

// Brain (Basic Monster Card)
//
//	Each time the attacking player rolls an attack roll of 6, they take 1 damage.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var brain = engine.CardDef{
	Ref:     "brain",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      2,
	DC:      3,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 1}},
}
