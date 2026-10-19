package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Sisters Vis (Boss Card)
//
//	When an attack is declared on this, the active player chooses another player. That player must make an attack roll against this after each attack roll the active player makes this attack. If this dies this turn, that player gains +1 treasure.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var sistersVis = engine.CardDef{
	Ref:     "sisters_vis",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      6,
	DC:      4,
	ATK:     2,
	Soul:    1,
	Rewards: []engine.Reward{{Kind: engine.RewardTreasure, Amount: 1}},
}
