package challenge

import engine "github.com/FourSoulsRepo/rules_engine"

// Mischievous The Haunt (Challenge Boss Card)
//
//	While the active player controls less than 3 souls, attacks can't be declared on this and it can't take damage.
//	Each time a player dies, put a fear counter on them.
//	For every 2 fear counters on the active player, monsters have +1{DC}.
//	Each time a monster dies, the active player chooses a player. Remove a fear counter from them.
//
// TODO(card): implement the text above, then delete this line.
// Until then `go run ./cmd/cardgen` rewrites this file from card_db.
// A question or option it asks is a key, engine.NewKey("card.mischievous_the_haunt_comp.<name>"),
// with its words and a note in frontend/src/locales/en.json (ADR 021).
var mischievousTheHauntComp = engine.CardDef{
	Ref:     "mischievous_the_haunt_comp",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      4,
	DC:      5,
	ATK:     1,
	Soul:    1,
	Rewards: []engine.Reward{{Kind: engine.RewardTreasure, Amount: 1}},
	Outside: true,
}
