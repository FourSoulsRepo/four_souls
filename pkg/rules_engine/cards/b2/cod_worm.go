package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Cod Worm (Basic Monster Card)
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var codWorm = engine.CardDef{
	Ref:     "cod_worm",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      2,
	DC:      5,
	ATK:     0,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 3}},
}
