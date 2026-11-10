package challenge

import engine "github.com/FourSoulsRepo/rules_engine"

// Pilfering Krampus (Challenge Boss Card)
//
//	While the number of souls controlled by all players is less than 5 or there is a minion in play, this can't be attacked or take damage.
//	As each turn starts, the active player rolls-
//	1-2: Remove 4 counters from the ¢ present.
//	3-4: Remove 1 counter from the loot present.
//	5: Remove 1 counter from the treasure present.
//	6: Remove 1 counter from the soul present.
//	While there is a minion in play, players must attack a random minion once each turn if able.
//	Players gain no rewards or souls when a minion dies.
var pilferingKrampusNorm = engine.CardDef{
	Ref:              "pilfering_krampus_norm",
	Kind:             engine.MonsterCard,
	Copies:           1,
	HP:               5,
	DC:               4,
	ATK:              2,
	Soul:             1,
	Rewards:          []engine.Reward{{Kind: engine.RewardTreasure, Amount: 1}},
	Outside:          true,
	Guarded:          guardedBelow(5),
	NoSpoils:         minionSpoils,
	MustAttackRandom: randomMinion,
	Abilities: []engine.Ability{
		presentRoll("As each turn starts, the active player rolls- 1-2: Remove 4 counters from the ¢ present. 3-4: Remove 1 counter from the loot present. 5: Remove 1 counter from the treasure present. 6: Remove 1 counter from the soul present.", noExtraRolls, engine.RollTable{}.
			Results(1, 2, takeCounters(engine.PresentCents, 4)).
			Results(3, 4, takeCounters(engine.PresentLoot, 1)).
			Results(5, 5, takeCounters(engine.PresentTreasure, 1)).
			Results(6, 6, takeCounters(engine.PresentSoul, 1))),
	},
}
