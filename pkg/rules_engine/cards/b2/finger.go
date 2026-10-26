package b2

import (
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/cardkit"
)

// The card's own lines; their words are in en.json (ADR 021).
var (
	fingerForWhichOf  = engine.NewKey("card.finger.for_which_of_their_items")
	fingerSwapWhichOf = engine.NewKey("card.finger.swap_which_of_your_items")
)

// Finger! (Passive Treasure Card)
//
//	Each time a player rolls a ❷, you may swap a non-eternal item you control with a non-eternal item they control.
var finger = engine.CardDef{
	Ref:    "finger",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "Each time a player rolls a 2, you may swap a non-eternal item you control with a non-eternal item they control.",
			Takes:   engine.TakesEventPlayer,
			Trigger: engine.OnRollOf(2),
			Effects: []engine.Effect{engine.Ask(fingerSwap, fingerMine, fingerTheirs)},
		},
	},
}

func nonEternalItems(g *engine.Game, p engine.PlayerID) []engine.ObjectID {
	var out []engine.ObjectID
	for _, id := range g.Players[p].InPlay {
		if g.Object(id).Role == engine.RoleItem && !g.Eternal(id) {
			out = append(out, id)
		}
	}
	return out
}

var fingerMine = engine.Question{Text: fingerSwapWhichOf, Options: func(c *engine.Ctx, _ []int) []engine.Msg {
	if c.EventPlayer == engine.NoPlayer || c.EventPlayer == c.Controller || len(nonEternalItems(c.G, c.EventPlayer)) == 0 {
		return nil
	}
	var out []engine.Msg
	for _, id := range nonEternalItems(c.G, c.Controller) {
		out = append(out, engine.CardOption(c.G.Object(id).Card))
	}
	if out == nil {
		return nil
	}
	return append(out, cardkit.AskDontSwap.Msg())
}}

var fingerTheirs = engine.Question{Text: fingerForWhichOf, Options: func(c *engine.Ctx, a []int) []engine.Msg {
	if a[0] < 0 || a[0] >= len(nonEternalItems(c.G, c.Controller)) {
		return nil
	}
	var out []engine.Msg
	for _, id := range nonEternalItems(c.G, c.EventPlayer) {
		out = append(out, engine.CardOption(c.G.Object(id).Card))
	}
	return out
}}

func fingerSwap(c *engine.Ctx, a []int) {
	if a[1] < 0 {
		return
	}
	mine, theirs := nonEternalItems(c.G, c.Controller)[a[0]], nonEternalItems(c.G, c.EventPlayer)[a[1]]
	c.G.GainControl(c.EventPlayer, mine)
	c.G.GainControl(c.Controller, theirs)
}
