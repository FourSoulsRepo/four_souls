package r

import (
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/cardkit"
)

// The card's own lines; their words are in en.json (ADR 021).
var (
	friendlyBallNeedsCard = engine.NewKey("card.friendly_ball.needs_card")
	friendlyBallNeedsNone = engine.NewKey("card.friendly_ball.needs_none")
)

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
			Kind: engine.Activated,
			Text: "↷: Put a monster not being attacked under this if there are no cards under this.",
			Costs: []engine.Cost{engine.If(friendlyBallNeedsNone, func(g *engine.Game, _ engine.PlayerID, self engine.ObjectID) bool {
				return len(g.Object(self).Under) == 0
			}), engine.Tap()},
			Targets: []engine.TargetSpec{engine.ChooseWhere(engine.TargetMonster, notAttacked)},
			Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
				if len(c.G.Object(c.Source).Under) == 0 {
					c.G.PutUnder(c.Source, c.Targets[0].Object)
				}
			})},
		},
		{
			Kind: engine.Activated,
			Text: "↷: Put a monster from under this in a monster slot not being attacked. The active player must make an additional attack on it this turn.",
			Costs: []engine.Cost{engine.If(friendlyBallNeedsCard, func(g *engine.Game, _ engine.PlayerID, self engine.ObjectID) bool {
				return len(g.Object(self).Under) > 0
			}), engine.Tap()},
			Effects: []engine.Effect{engine.Ask(func(c *engine.Ctx, a []int) {
				under := c.G.Object(c.Source).Under
				if len(under) == 0 || a[0] < 0 {
					return
				}
				if nid := c.G.FromUnder(c.Source, under[0], freeSlots(c.G)[a[0]]); nid != 0 && c.G.Object(nid).Role == engine.RoleMonster {
					c.G.ForceAttacks(1, false)
					c.G.Turn.MustAttack = nid
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
