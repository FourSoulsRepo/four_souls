package r

import engine "github.com/FourSoulsRepo/rules_engine"

// The Scourge (Boss Card)
//
//	When this would take damage, each living player votes either whip or whiff-
//	If whip wins, prevent the damage this would take and each non-active player takes 1 damage.
//	If whiff wins or there is a tie, the active player loots 1.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var theScourge = engine.CardDef{
	Ref:       "the_scourge",
	Kind:      engine.MonsterCard,
	Copies:    1,
	HP:        3,
	DC:        4,
	ATK:       1,
	Soul:      1,
	Rewards:   []engine.Reward{{Kind: engine.RewardLoot, Amount: 2}},
	ThreePlus: true,
}
