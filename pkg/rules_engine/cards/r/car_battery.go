package r

import (
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/cardkit"
)

// Car Battery (Paid Treasure Card)
//
//	This can't be recharged except by its own abilities.
//	{Tap Effect}Recharge up to 3 other items in play.
//	{Paid Effect}Pay 7¢: Recharge this.
var carBattery = engine.CardDef{
	Ref:             "car_battery",
	Kind:            engine.TreasureCard,
	Copies:          1,
	Tap:             true,
	OwnRechargeOnly: true,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Activated,
			Text:    "↷: Recharge up to 3 other items in play.",
			Costs:   []engine.Cost{engine.Tap()},
			Effects: []engine.Effect{rechargeUpTo(3)},
		},
		{
			Kind:    engine.Activated,
			Text:    "$ Pay 7¢: Recharge this.",
			Costs:   []engine.Cost{engine.PayCents(7)},
			Effects: []engine.Effect{engine.RechargeSelf()},
		},
	},
}

// rechargeUpTo recharges up to n other deactivated items, one at a time.
func rechargeUpTo(n int) engine.Effect {
	items := func(c *engine.Ctx) []engine.ObjectID {
		var out []engine.ObjectID
		for _, id := range itemsOf(c.G, engine.NoPlayer) {
			if id != c.Source && !c.G.Object(id).Charged {
				out = append(out, id)
			}
		}
		return out
	}
	qs := make([]engine.Question, n)
	for i := range qs {
		qs[i] = engine.Question{Text: cardkit.AskRechargeWhichItem, Options: func(c *engine.Ctx, a []int) []engine.Msg {
			if i > 0 && a[i-1] < 0 {
				return nil
			}
			left := items(c)
			for _, x := range a[:i] {
				if x >= len(left) {
					return nil // done
				}
				left = append(left[:x:x], left[x+1:]...)
			}
			if len(left) == 0 {
				return nil
			}
			return append(labelsOf(c.G, left), engine.AskDone.Msg())
		}}
	}
	return engine.Ask(func(c *engine.Ctx, a []int) {
		left := items(c)
		var picked []engine.ObjectID
		for _, x := range a {
			if x < 0 || x >= len(left) {
				break
			}
			picked = append(picked, left[x])
			left = append(left[:x:x], left[x+1:]...)
		}
		for _, id := range picked {
			c.G.Recharge(id)
		}
	}, qs...)
}
