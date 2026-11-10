package challenge

import engine "github.com/FourSoulsRepo/rules_engine"

// Pilfering Krampus (Challenge Boss Card)
//
//	While the number of souls controlled by all players is less than 6 or there is a minion in play, this can't be attacked or take damage.
//	As each turn starts, the active player rolls, then rolls again for every 2 minions in play-
//	1-2: Remove 5 counters from the ¢ present.
//	3-4: Remove 2 counters from the loot present.
//	5: Remove 1 counter from the treasure present.
//	6: Remove 1 counter from the soul present.
//	Players gain no rewards or souls when a minion dies.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
// A question or option it asks is a key, engine.NewKey("card.pilfering_krampus_hard.<name>"),
// with its words and a note in frontend/src/locales/en.json (ADR 021).
var pilferingKrampusHard = engine.CardDef{
	Ref:     "pilfering_krampus_hard",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      6,
	DC:      4,
	ATK:     2,
	Soul:    1,
	Rewards: []engine.Reward{{Kind: engine.RewardTreasure, Amount: 1}},
	Outside: true,
}
