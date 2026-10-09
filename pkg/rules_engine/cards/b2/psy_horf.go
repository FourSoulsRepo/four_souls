package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Psy Horf (Basic Monster Card)
//
//	When this dies, the active player recharges each item they control.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var psyHorf = engine.CardDef{
	Ref:     "psy_horf",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      1,
	DC:      5,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 1}},
}
