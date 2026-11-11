package challenge

import engine "github.com/FourSoulsRepo/rules_engine"

// Mischievous The Haunt (Challenge Boss Card)
//
//	This starts with 12 trick counters.
//	While there are more trick counters on this than the number of souls controlled by all players, attacks can't be declared on this and it can't take damage.
//	Each time a player takes damage, put a fear counter on them.
//	For every 2 fear counters on the active player, monsters have +1{DC}.
//	Each time a monster dies, the active player rolls-
//	4-6: Remove a fear counter from a player.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
// A question or option it asks is a key, engine.NewKey("card.mischievous_the_haunt_norm.<name>"),
// with its words and a note in frontend/src/locales/en.json (ADR 021).
var mischievousTheHauntNorm = engine.CardDef{
	Ref:     "mischievous_the_haunt_norm",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      4,
	DC:      5,
	ATK:     1,
	Soul:    1,
	Rewards: []engine.Reward{{Kind: engine.RewardTreasure, Amount: 1}},
	Outside: true,
}
