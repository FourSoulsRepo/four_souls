package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// Parabite (Basic Monster Card)
//
//	Each time this takes combat damage, the attacking player has -1 to their next attack roll this turn.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var parabite = engine.CardDef{
	Ref:     "parabite",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      2,
	DC:      3,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 1}},
}
