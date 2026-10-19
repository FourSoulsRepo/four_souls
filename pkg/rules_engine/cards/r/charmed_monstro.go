package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Charmed Monstro (Boss Card)
//
//	If the active player would gain this soul, they instead choose another player. That player gains this soul.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var charmedMonstro = engine.CardDef{
	Ref:     "charmed_monstro",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      4,
	DC:      3,
	ATK:     1,
	Soul:    1,
	Rewards: []engine.Reward{{Kind: engine.RewardTreasure, Amount: 2}},
}
