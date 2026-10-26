package r

import engine "github.com/FourSoulsRepo/rules_engine"

// The card's own lines; their words are in en.json (ADR 021).
var (
	gelloDealTheDamage = engine.NewKey("card.gello.deal_the_damage_to_which")
)

// Gello (Eternal Treasure Card)
//
//	{Tap Effect}Choose a monster being attacked. Heal that monster to full {HP}, then deal damage equal to the number of {HP} healed in this way to another monster. If it's not your turn, cancel the attack and the active player may make an additional attack this turn.
//	-Eternal- This can't be destroyed or put into discard.
var gello = engine.CardDef{
	Ref:     "gello",
	Kind:    engine.TreasureCard,
	Copies:  1,
	Eternal: true,
	Outside: true,
	Tap:     true,
	Abilities: []engine.Ability{{
		Kind:  engine.Activated,
		Text:  "↷: Choose a monster being attacked. Heal that monster to full HP, then deal damage equal to the number of HP healed in this way to another monster. If it's not your turn, cancel the attack and the active player may make an additional attack this turn.",
		Costs: []engine.Cost{engine.Tap()},
		Targets: []engine.TargetSpec{engine.ChooseWhere(engine.TargetMonster, func(g *engine.Game, _ engine.ObjectID, c engine.Chosen) bool {
			return g.Attack.Started && g.Attack.Target == c.Object && g.HP(c.Object) > 0 // S-DIS: not at 0 HP
		})},
		Effects: []engine.Effect{engine.Ask(func(c *engine.Ctx, a []int) {
			m := c.Targets[0].Object
			healed := c.G.Object(m).Damage
			c.G.HealObject(m, healed)
			if others := otherMonsters(c.G, m); healed > 0 && a[0] >= 0 && a[0] < len(others) {
				c.G.DealDamageTo(engine.Target{Object: others[a[0]]}, healed, c.Controller, c.Source)
			}
			if c.Controller != c.G.Turn.Active {
				c.Do(engine.CancelAttack())
				c.G.Turn.Attacks++
			}
		}, engine.Question{Text: gelloDealTheDamage, Options: func(c *engine.Ctx, _ []int) []engine.Msg {
			if c.G.Object(c.Targets[0].Object).Damage == 0 {
				return nil
			}
			return labelsOf(c.G, otherMonsters(c.G, c.Targets[0].Object))
		}})},
	}},
}
