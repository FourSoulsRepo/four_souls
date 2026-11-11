package challenge

import engine "github.com/FourSoulsRepo/rules_engine"

// Pilfering Krampus (Challenge Boss Card)
//
//	While the number of souls controlled by the active player is less than 3 or there is a minion in a present slot they control, this can't be attacked or take damage.
//	As each turn starts, the active player rolls-
//	1-2: Remove 3 counters from their ¢ present.
//	3-4: Remove 1 counter from their loot present.
//	5: Remove 1 counter from their treasure present.
//	6: Remove 1 counter from their soul present.
//	While there is a minion in a present slot controlled by the active player, the active player must attack a random minion once each turn if able.
//	Players gain no rewards or souls when a minion dies.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
// A question or option it asks is a key, engine.NewKey("card.pilfering_krampus_comp.<name>"),
// with its words and a note in frontend/src/locales/en.json (ADR 021).
var pilferingKrampusComp = engine.CardDef{
	Ref:     "pilfering_krampus_comp",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      5,
	DC:      4,
	ATK:     2,
	Soul:    1,
	Rewards: []engine.Reward{{Kind: engine.RewardTreasure, Amount: 1}},
	Outside: true,
}
