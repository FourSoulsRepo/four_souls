package challenge

import engine "github.com/FourSoulsRepo/rules_engine"

// Pilfering Krampus (Challenge Boss Card)
//
//	While the number of souls controlled by all players is less than 7 or there is a minion in play, this can't be attacked or take damage.
//	As each turn starts, the active player rolls, then rolls again for each minion in play-
//	1: Remove 8 counters from the ¢ present.
//	2-3: Remove 2 counters from the loot present.
//	4-5: Remove 1 counter from the treasure present.
//	6: Remove 1 counter from the soul present.
//	Players gain no rewards or souls when a minion dies.
var pilferingKrampusUltra = engine.CardDef{
	Ref:      "pilfering_krampus_ultra",
	Kind:     engine.MonsterCard,
	Copies:   1,
	HP:       6,
	DC:       4,
	ATK:      2,
	Soul:     1,
	Rewards:  []engine.Reward{{Kind: engine.RewardTreasure, Amount: 1}},
	Outside:  true,
	Guarded:  guardedBelow(7),
	NoSpoils: minionSpoils,
	Abilities: []engine.Ability{
		presentRoll("As each turn starts, the active player rolls, then rolls again for each minion in play- 1: Remove 8 counters from the ¢ present. 2-3: Remove 2 counters from the loot present. 4-5: Remove 1 counter from the treasure present. 6: Remove 1 counter from the soul present.", func(m int) int { return m }, engine.RollTable{}.
			Results(1, 1, takeCounters(engine.PresentCents, 8)).
			Results(2, 3, takeCounters(engine.PresentLoot, 2)).
			Results(4, 5, takeCounters(engine.PresentTreasure, 1)).
			Results(6, 6, takeCounters(engine.PresentSoul, 1))),
	},
}
