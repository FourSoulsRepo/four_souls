package g2

import engine "github.com/FourSoulsRepo/rules_engine"

// Lemon Mishap (Active Treasure Card)
//
//	{Tap Effect}Till end of turn, each time a player rolls an attack roll of 1 or 2, deal 1 damage to each monster.
var lemonMishap = engine.CardDef{
	Ref:    "lemon_mishap",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
	Abilities: []engine.Ability{{
		Kind:    engine.Activated,
		Text:    "↷: Till end of turn, each time a player rolls an attack roll of 1 or 2, deal 1 damage to each monster.",
		Costs:   []engine.Cost{engine.Tap()},
		Effects: []engine.Effect{engine.Linger(0)},
	}},
	Lingering: []engine.Ability{{
		Kind: engine.Triggered,
		Text: "Each time a player rolls an attack roll of 1 or 2, deal 1 damage to each monster.",
		Trigger: engine.Trigger{On: engine.EvRollResolved, Match: func(_ *engine.Game, _ engine.ObjectID, e engine.Event) bool {
			return e.Text == "attack" && (e.Amount == 1 || e.Amount == 2)
		}},
		Effects: []engine.Effect{engine.EachMonsterTakesDamage(1)},
	}},
}
