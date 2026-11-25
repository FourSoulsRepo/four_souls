package r

import (
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/cardkit"
)

// The card's own lines; their words are in en.json (ADR 021).
var friendlyBallNeedsNone = engine.NewKey("card.friendly_ball.needs_none")

// Friendly Ball (Active Treasure Card)
//
//	{Tap Effect}Put a monster not being attacked under this if there are no cards under this.
//	{Tap Effect}Put a monster from under this in a monster slot not being attacked. The active player must make an additional attack on it this turn.
var friendlyBall = engine.CardDef{
	Ref:    "friendly_ball",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Tap:    true,
	Abilities: []engine.Ability{
		{
			Kind:  engine.Activated,
			Text:  "↷: Put a monster not being attacked under this if there are no cards under this.",
			Costs: []engine.Cost{engine.Tap()},
			// "if there are no cards under this": usable only then, a
			// limit the card states (R-ABIL-08; owner, 19.242).
			Needs: []engine.Need{engine.OnlyIf(friendlyBallNeedsNone, func(g *engine.Game, _ engine.PlayerID, self engine.ObjectID) bool {
				return len(g.Object(self).Under) == 0
			})},
			Targets: []engine.TargetSpec{engine.ChooseWhere(engine.TargetMonster, notAttacked)},
			Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
				if len(c.G.Object(c.Source).Under) == 0 {
					c.G.PutUnder(c.Source, c.Targets[0].Object)
				}
			})},
		},
		{
			Kind:  engine.Activated,
			Text:  "↷: Put a monster from under this in a monster slot not being attacked. The active player must make an additional attack on it this turn.",
			Costs: []engine.Cost{engine.Tap()},
			// "a monster from under this" is its target (R-ABIL-05); with
			// none it has no valid target (R-ABIL-40; owner, 19.242).
			Needs: []engine.Need{engine.ActsOn(func(g *engine.Game, _ engine.PlayerID, self engine.ObjectID) bool {
				return len(g.Object(self).Under) > 0
			})},
			Effects: []engine.Effect{engine.Ask(func(c *engine.Ctx, a []int) {
				under := c.G.Object(c.Source).Under
				if len(under) == 0 || a[0] < 0 {
					return
				}
				if nid := c.G.FromUnder(c.Source, under[0], freeSlots(c.G)[a[0]]); nid != 0 && c.G.Object(nid).Role == engine.RoleMonster {
					c.G.ForceAttackOn(nid) // "an additional attack on it" (19.241)
				}
			}, engine.Question{Text: cardkit.AskPutItInWhichMonster, Options: func(c *engine.Ctx, _ []int) []engine.Msg {
				if len(c.G.Object(c.Source).Under) == 0 {
					return nil
				}
				var out []engine.Msg
				for _, i := range freeSlots(c.G) {
					out = append(out, engine.SlotOption(i+1, ""))
				}
				return out
			}})},
		},
	},
}
