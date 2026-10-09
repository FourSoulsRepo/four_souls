package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Cursed Psy Horf (Cursed Monster Card)
//
//	{Curse Effect}Each time a player activates an item, they take 1 damage.
var cursedPsyHorf = engine.CardDef{
	Ref:     "cursed_psy_horf",
	Kind:    engine.MonsterCard,
	Copies:  1,
	HP:      1,
	DC:      5,
	ATK:     1,
	Rewards: []engine.Reward{{Kind: engine.RewardLoot, Amount: 2}},
	Abilities: []engine.Ability{
		{
			Kind: engine.Triggered,
			Text: "Each time a player activates an item, they take 1 damage.",
			Trigger: engine.Trigger{On: engine.EvActivated, Match: func(g *engine.Game, _ engine.ObjectID, e engine.Event) bool {
				return g.Object(e.Object).Role == engine.RoleItem
			}},
			Effects: []engine.Effect{forRoller(engine.DealDamage(1, engine.You))},
		},
	},
}
