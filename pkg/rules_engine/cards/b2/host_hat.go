package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Host Hat (Active Treasure Card)
//
//	{Tap Effect}Prevent the next 1 damage you would take this turn. When you prevent damage this way, deal 1 damage to another player.
var hostHat = engine.CardDef{
	Ref:    "host_hat",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Activated,
			Text:    "↷: Prevent the next 1 damage you would take this turn. When you prevent damage this way, deal 1 damage to another player.",
			Costs:   []engine.Cost{engine.Tap()},
			Effects: []engine.Effect{engine.PreventDamage(1, engine.You)},
		},
		{
			Kind:    engine.Triggered,
			Text:    "When you prevent damage this way, deal 1 damage to another player.",
			Trigger: whenThisPrevents(),
			Effects: []engine.Effect{damageAPlayer(1, true)},
		},
	},
}

// whenThisPrevents triggers when a shield made by this object prevents
// damage.
func whenThisPrevents() engine.Trigger {
	return engine.Trigger{On: engine.EvPrevented, Match: func(_ *engine.Game, self engine.ObjectID, e engine.Event) bool {
		return e.Source == self && e.Amount > 0
	}}
}
