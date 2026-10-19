package r

import engine "github.com/FourSoulsRepo/rules_engine"

// The Scourge (Boss Card)
//
//	When this would take damage, each living player votes either whip or whiff-
//	If whip wins, prevent the damage this would take and each non-active player takes 1 damage.
//	If whiff wins or there is a tie, the active player loots 1.
var theScourge = engine.CardDef{
	Ref:       "the_scourge",
	Kind:      engine.MonsterCard,
	Copies:    1,
	HP:        3,
	DC:        4,
	ATK:       1,
	Soul:      1,
	Rewards:   []engine.Reward{{Kind: engine.RewardLoot, Amount: 2}},
	ThreePlus: true,
	Abilities: []engine.Ability{{
		Kind: engine.Triggered,
		Text: "When this would take damage, each living player votes either whip or whiff- If whip wins, prevent the damage this would take and each non-active player takes 1 damage. If whiff wins or there is a tie, the active player loots 1.",
		Trigger: engine.Trigger{On: engine.EvDamagePending, Match: func(_ *engine.Game, self engine.ObjectID, e engine.Event) bool {
			return e.Object == self
		}},
		Effects: []engine.Effect{vote([]string{"whip", "whiff"}, func(c *engine.Ctx, votes []int) {
			if votes[0] > votes[1] {
				c.G.CancelStackItem(c.EventStack)
				for _, p := range othersOfActive(c.G) {
					c.G.DealDamageTo(engine.Target{Player: p, IsPlayer: true}, 1, engine.NoPlayer, c.Source)
				}
				return
			}
			c.Do(engine.Loot(1))
		})},
	}},
}
