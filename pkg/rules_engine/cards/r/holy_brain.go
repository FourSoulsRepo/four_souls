package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Holy Brain (Holy/Charmed Monster Card)
//
//	Each time a player rolls a ❸, each monster gains -1{DC} till end of turn.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var holyBrain = engine.CardDef{
	Ref:     "holy_brain",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      2,
	DC:      3,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 1}},
}
