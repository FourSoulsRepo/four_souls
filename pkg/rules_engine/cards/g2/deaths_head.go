package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// Death’s Head (Basic Monster Card)
//
//	Monsters have +1{ATK}.
//	This can't be attacked.
//	When another monster dies, this dies.
var deathsHead = engine.CardDef{
	Ref:          "deaths_head",
	Kind:         engine.MonsterCard,
	Copies:       1,
	HP:           2,
	DC:           0,
	ATK:          0,
	Rewards:      []engine.Reward{{Kind: engine.RewardLoot, Amount: 1}},
	Unattackable: true,
	Statics:      []engine.Static{engine.MonstersHave(engine.StatMonsterATK, 1)},
	Abilities: []engine.Ability{{
		Kind: engine.Triggered,
		Text: "When another monster dies, this dies.",
		Trigger: engine.Trigger{On: engine.EvDied, Match: func(g *engine.Game, self engine.ObjectID, e engine.Event) bool {
			return e.Object != self && e.Player == engine.NoPlayer && g.Kind(e.Object) == engine.MonsterCard
		}},
		Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
			if c.G.Object(c.Source).Zone.Kind == engine.ZoneInPlay {
				c.G.KillObject(c.Source)
			}
		})},
	}},
}
