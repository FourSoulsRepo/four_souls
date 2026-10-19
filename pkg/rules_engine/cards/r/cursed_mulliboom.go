package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Cursed Mulliboom (Cursed Monster Card)
//
//	Each time a monster dies, this deals 3 damage to a player chosen at random.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var cursedMulliboom = engine.CardDef{
	Ref:     "cursed_mulliboom",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      1,
	DC:      4,
	ATK:     3,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 2}},
}
