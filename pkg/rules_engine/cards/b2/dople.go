package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Dople (Basic Monster Card)
//
//	Damage dealt to this is also dealt to the player to the active player's right.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var dople = engine.CardDef{
	Ref:     "dople",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      2,
	DC:      4,
	ATK:     2,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 7}},
}
