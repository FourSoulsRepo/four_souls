package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Pride (Boss Card)
//
//	When an attack is declared on this, the active player chooses a player. That player discards 2 loot cards.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var pride = engine.CardDef{
	Ref:     "pride",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      2,
	DC:      4,
	ATK:     1,
	Soul:    1,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 5}},
}
