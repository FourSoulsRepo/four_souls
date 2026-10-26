package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// The card's own lines; their words are in en.json (ADR 021).
var (
	spiderModReplaceWhichMonster = engine.NewKey("card.spider_mod.replace_which_monster")
)

// Spider Mod (Passive Treasure Card)
//
//	Each time a player rolls a ❺, you may put a monster not being attacked into discard and replace it with the top card of the monster deck.
var spiderMod = engine.CardDef{
	Ref:    "spider_mod",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "Each time a player rolls a 5, you may put a monster not being attacked into discard and replace it with the top card of the monster deck.",
			Trigger: engine.OnRollOf(5),
			Effects: []engine.Effect{engine.Ask(func(c *engine.Ctx, a []int) {
				if ids := notAttacked(c.G); a[0] >= 0 && a[0] < len(ids) {
					c.G.DiscardMonster(ids[a[0]])
					c.G.RefillSlots()
				}
			}, engine.Question{Text: spiderModReplaceWhichMonster, Options: func(c *engine.Ctx, _ []int) []engine.Msg {
				var out []engine.Msg
				for _, id := range notAttacked(c.G) {
					out = append(out, engine.CardOption(c.G.Object(id).Card))
				}
				if out == nil {
					return nil
				}
				return append(out, engine.AskNone.Msg())
			}})},
		},
	},
}

func notAttacked(g *engine.Game) []engine.ObjectID {
	var out []engine.ObjectID
	for _, id := range monstersInPlay(g) {
		if id != g.Attack.Target {
			out = append(out, id)
		}
	}
	return out
}
