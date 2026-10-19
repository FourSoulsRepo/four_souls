package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Turdlings (Boss Card)
//
//	This only takes combat damage on attack rolls of 2.
//	When the attacking player makes their second attack roll this turn, after combat damage, cancel the attack.
var turdlings = engine.CardDef{
	Ref:     "turdlings",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      2,
	DC:      2,
	ATK:     2,
	Soul:    1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 2}},
	CombatMod: func(_ *engine.Game, _ engine.ObjectID, roll int, hit bool, n int) int {
		if hit && roll != 2 {
			return 0
		}
		return n
	},
	Abilities: []engine.Ability{{
		Kind: engine.Triggered,
		Text: "When the attacking player makes their second attack roll this turn, after combat damage, cancel the attack.",
		Trigger: engine.Trigger{On: engine.EvRollResolved, Match: func(g *engine.Game, self engine.ObjectID, e engine.Event) bool {
			// The roll is counted after this event: 1 means this is the 2nd.
			return e.Text == "attack" && e.Player == g.Turn.Active && g.Attack.Target == self && g.Turn.AttackRolls == 1
		}},
		Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) { c.G.CancelAttackAfterDamage() })},
	}},
}
