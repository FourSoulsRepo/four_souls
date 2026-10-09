package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Gurdy Jr. (Boss Card)
//
//	Each time the attacking player activates an item, they take 1 damage.
var gurdyJr = engine.CardDef{
	Ref:     "gurdy_jr",
	Kind:    engine.MonsterCard,
	Copies:  1,
	Soul:    1,
	HP:      2,
	DC:      5,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardTreasure, Amount: 1}},
	Abilities: []engine.Ability{
		{
			Kind: engine.Triggered,
			Text: "Each time the attacking player activates an item, they take 1 damage.",
			Trigger: engine.Trigger{On: engine.EvActivated, Match: func(g *engine.Game, self engine.ObjectID, e engine.Event) bool {
				return g.Attack.Target == self && e.Player == g.Turn.Active && g.Object(e.Object).Role == engine.RoleItem
			}},
			Effects: []engine.Effect{engine.DealDamage(1, engine.You)},
		},
	},
}
