package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Cursed Psy Horf (Cursed Monster Card)
//
//	{Curse Effect}Each time a player activates an item, they take 1 damage.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var cursedPsyHorf = engine.CardDef{
	Ref:     "cursed_psy_horf",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      1,
	DC:      5,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 2}},
}
