package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// Imp (Basic Monster Card)
//
//	Each time the attacking player rolls an attack roll of 3, this gains +1{DC} till end of turn.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var imp = engine.CardDef{
	Ref:     "imp",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      3,
	DC:      5,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Roll: true}},
}
