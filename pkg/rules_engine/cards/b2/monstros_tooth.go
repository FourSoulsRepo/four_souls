package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// The card's own lines; their words are in en.json (ADR 021).
var (
	monstrosToothDestroyWhichOf = engine.NewKey("card.monstros_tooth.destroy_which_of_your_items")
	monstrosToothWhichPlayer    = engine.NewKey("card.monstros_tooth.which_player")
)

// Monstro’s Tooth (Passive Treasure Card)
//
//	At the start of your turn, choose a player at random. That player destroys an item they control.
var monstrosTooth = engine.CardDef{
	Ref:    "monstros_tooth",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "At the start of your turn, choose a player at random. That player destroys an item they control.",
			Takes:   engine.TakesOthers,
			Trigger: engine.AtStartOfYourTurn(),
			Effects: []engine.Effect{engine.Ask(toothDestroy, toothPlayer, toothItem)},
		},
	},
}

var toothPlayer = engine.Question{Random: true, Text: monstrosToothWhichPlayer, Options: func(c *engine.Ctx, _ []int) []engine.Msg {
	var out []engine.Msg
	for _, pl := range c.G.Players {
		out = append(out, playerLabel(c.G, pl.ID))
	}
	return out
}}

var toothItem = engine.Question{
	Text:   monstrosToothDestroyWhichOf,
	Player: func(_ *engine.Ctx, a []int) engine.PlayerID { return engine.PlayerID(a[0]) },
	Options: func(c *engine.Ctx, a []int) []engine.Msg {
		var out []engine.Msg
		for _, id := range itemsOf(c.G, engine.PlayerID(a[0])) {
			out = append(out, engine.CardOption(c.G.Object(id).Card))
		}
		return out
	},
}

func itemsOf(g *engine.Game, p engine.PlayerID) []engine.ObjectID {
	var out []engine.ObjectID
	for _, id := range g.Players[p].InPlay {
		if g.Object(id).Role == engine.RoleItem {
			out = append(out, id)
		}
	}
	return out
}

func toothDestroy(c *engine.Ctx, a []int) {
	if a[1] >= 0 {
		p := engine.PlayerID(a[0])
		c.G.DestroyObject(p, itemsOf(c.G, p)[a[1]])
	}
}
