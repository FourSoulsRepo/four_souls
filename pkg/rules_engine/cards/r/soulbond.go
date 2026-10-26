package r

import (
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/cardkit"
)

// The card's own lines; their words are in en.json (ADR 021).
var (
	soulbondWhichPlayerDies = engine.NewKey("card.soulbond.which_player_dies")
)

// Soulbond (Eternal Treasure Card)
//
//	Each time you die, choose another player. That player dies.
//	You don't lose ¢ or discard loot cards when paying the death penalty.
//	-Eternal- This can't be destroyed or put into discard.
var soulbond = engine.CardDef{
	Ref:       "soulbond",
	Kind:      engine.TreasureCard,
	Copies:    1,
	Eternal:   true,
	ThreePlus: true,
	Outside:   true,
	Replacements: []engine.Replacement{{
		Text: "You don't lose ¢ or discard loot cards when paying the death penalty.",
		When: func(g *engine.Game, self engine.ObjectID, a engine.Action) bool {
			me := g.Object(self).Controller
			return a.Player == me && (a.Kind == engine.ActPenaltyLoot || (a.Kind == engine.ActLoseCents && a.Penalty))
		},
		Do: func(*engine.Game, engine.ObjectID, engine.Action) []engine.Action { return nil },
	}},
	Abilities: []engine.Ability{{
		Kind:    engine.Triggered,
		Text:    "Each time you die, choose another player. That player dies.",
		Trigger: engine.WhenYouDie(),
		Effects: []engine.Effect{engine.Ask(func(c *engine.Ctx, a []int) {
			if ps := cardkit.LivingPlayers(c.G, c.Controller); a[0] >= 0 && a[0] < len(ps) {
				them := *c
				them.Targets = []engine.Chosen{{Kind: engine.TargetPlayer, Player: ps[a[0]]}}
				them.Do(engine.Kill(0))
			}
		}, engine.Question{Text: soulbondWhichPlayerDies, Options: func(c *engine.Ctx, _ []int) []engine.Msg {
			return playerLabels(c.G, cardkit.LivingPlayers(c.G, c.Controller))
		}})},
	}},
}
