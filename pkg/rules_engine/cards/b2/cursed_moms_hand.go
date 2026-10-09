package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Cursed Mom’s Hand (Cursed Monster Card)
//
//	{Curse Effect}When the active player rolls a 6, cancel everything that hasn't resolved and end the turn.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var cursedMomsHand = engine.CardDef{
	Ref:     "cursed_moms_hand",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      2,
	DC:      4,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 4}},
}
