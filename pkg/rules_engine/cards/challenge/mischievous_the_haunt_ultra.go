package challenge

import engine "github.com/FourSoulsRepo/rules_engine"

// Mischievous The Haunt (Challenge Boss Card)
//
//	This starts with 14 trick counters.
//	While there are more trick counters on this than the number of souls controlled by all players, attacks can't be declared on this and it can't take damage.
//	Each time a player takes damage, put a fear counter on each player.
//	For every 2 fear counters on the active player, monsters have +1{DC}.
//	Each time a monster dies, the active player rolls-
//	4-6: Remove a fear counter from a player.
var mischievousTheHauntUltra = haunt(engine.CardDef{
	Ref:     "mischievous_the_haunt_ultra",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      4,
	DC:      5,
	ATK:     1,
	Soul:    1,
	Rewards: []engine.Reward{{Kind: engine.RewardTreasure, Amount: 1}},
	Outside: true,
}, "Each time a player takes damage, put a fear counter on each player.", true)
