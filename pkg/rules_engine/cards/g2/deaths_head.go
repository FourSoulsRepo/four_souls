package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// Death’s Head (Basic Monster Card)
//
//	Monsters have +1{ATK}.
//	This can't be attacked.
//	When another monster dies, this dies.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var deathsHead = engine.CardDef{
	Ref:     "deaths_head",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      2,
	DC:      0,
	ATK:     0,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 1}},
}
