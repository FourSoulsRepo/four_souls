package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// The card's own lines; their words are in en.json (ADR 021).
var (
	theLambNobody        = engine.NewKey("card.the_lamb.nobody")
	theLambGiveWhichSoul = engine.NewKey("card.the_lamb.give_which_soul")
	theLambWhoGivesYou   = engine.NewKey("card.the_lamb.who_gives_you_a_soul")
)

// The Lamb (Epic Boss Card)
//
//	When this dies, the active player may choose another player. They give you a soul they control.
var theLamb = engine.CardDef{
	Ref:     "the_lamb",
	Kind:    engine.MonsterCard,
	Copies:  1,
	Soul:    1,
	HP:      6,
	DC:      3,
	ATK:     6,
	Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 3}},
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "When this dies, the active player may choose another player. They give you a soul they control.",
			Takes:   engine.TakesOthers,
			Trigger: engine.WhenThisDies(),
			Effects: []engine.Effect{engine.Ask(lamb, lambPlayer, lambSoul)},
		},
	},
}

var lambPlayer = engine.Question{Text: theLambWhoGivesYou, Options: func(c *engine.Ctx, _ []int) []engine.Msg {
	var out []engine.Msg
	for _, p := range otherPlayers(c) {
		out = append(out, playerLabel(c.G, p))
	}
	return append(out, theLambNobody.Msg())
}}

var lambSoul = engine.Question{
	Text:   theLambGiveWhichSoul,
	Player: func(c *engine.Ctx, a []int) engine.PlayerID { return otherPlayers(c)[a[0]] },
	Options: func(c *engine.Ctx, a []int) []engine.Msg {
		if a[0] < 0 || a[0] >= len(otherPlayers(c)) {
			return nil
		}
		var out []engine.Msg
		for _, id := range souls(c.G, otherPlayers(c)[a[0]]) {
			out = append(out, engine.CardOption(c.G.Object(id).Card))
		}
		return out
	},
}

func lamb(c *engine.Ctx, a []int) {
	if a[1] >= 0 {
		c.G.GainControl(c.Controller, souls(c.G, otherPlayers(c)[a[0]])[a[1]])
	}
}
