package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Holy Psy Horf (Holy/Charmed Monster Card)
//
//	At the end of each turn, the active player loots 1.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var holyPsyHorf = engine.CardDef{
	Ref:     "holy_psy_horf",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      3,
	DC:      3,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 3}},
}
