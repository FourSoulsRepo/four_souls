package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// Fistula (Boss Card)
//
//	Each time this takes damage, it gains +1{DC} till end of turn.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var fistula = engine.CardDef{
	Ref:     "fistula",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      4,
	DC:      2,
	ATK:     1,
	Soul:    1,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 7}},
}
