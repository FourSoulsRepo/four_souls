package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Mega Satan! (Epic Boss Card)
//
//	Players can be attacked and have 4+{DC}.
//	At the end of each turn, if the active player is the only living player, they win.
//	-Indomitable- Indomitable monsters can't be covered. If they would be put in a monster slot, expand monster slots by 1 and put them in that new slot.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
var megaSatan = engine.CardDef{
	Ref:       "mega_satan",
	Kind:      engine.MonsterCard,
	Copies:    1,
	HP:        6,
	DC:        6,
	ATK:       6,
	Soul:      2,
	Rewards:   []engine.Reward{{Kind: engine.RewardTreasure, Amount: 2}},
	ThreePlus: true,
}
