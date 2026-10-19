package r

import engine "github.com/FourSoulsRepo/rules_engine"

// The Harbingers/The Beast! (Outside Card)
//
//	The Harbingers: The active player must attack this once each turn if able.
//	If this would die, instead put a counter on this, the active player gains +1 treasure, cancels their attack and this heals to full {HP}.
//	If this has 4+ counters, flip it.
//	The Beast!: When this is flipped to this side or attacked, each player rolls-
//	Each player who rolls the lowest or tied for the lowest dies.
//	-Indomitable- Indomitable monsters can't be covered. If they would be put in a monster slot, expand monster slots by 1 and put them in that new slot.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
// Not generated: HP is "4/6".
// Not generated: DC is "4/4".
// Not generated: ATK is "1/2".
var theHarbingers = engine.CardDef{
	Ref:     "the_harbingers",
	Kind:    engine.MonsterCard,
	Copies:  1,
	Soul:    2,
	Rewards: []engine.Reward{{Kind: engine.RewardTreasure, Amount: 1}, {Kind: engine.RewardLoot, Amount: 6}},
}
