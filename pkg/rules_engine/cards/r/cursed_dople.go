package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Cursed Dople (Cursed Monster Card)
//
//	Each time this takes damage, it deals 1 damage to each non-active player.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var cursedDople = engine.CardDef{
	Ref:       "cursed_dople",
	Kind:      engine.MonsterCard,
	Copies:    1,
	HP:        3,
	DC:        4,
	ATK:       2,
	Rewards:   []engine.Reward{{Kind: engine.RewardCents, Roll: true}},
	ThreePlus: true,
}
