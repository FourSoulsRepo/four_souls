package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// Goat Head (Passive Treasure Card)
//
//	At the end of your turn, you may discard any number of loot cards, then loot equal to the number of cards discarded in this way.
var goatHead = engine.CardDef{
	Ref:    "goat_head",
	Kind:   engine.TreasureCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "At the end of your turn, you may discard any number of loot cards, then loot equal to the number of cards discarded in this way.",
			Trigger: engine.AtEndOfYourTurn(),
			Effects: []engine.Effect{engine.Ask(goatHeadDiscard, goatQuestions()...)},
		},
	},
}

// goatQuestions ask for one card at a time until "done".
func goatQuestions() []engine.Question {
	var qs []engine.Question
	for range 20 {
		qs = append(qs, engine.Question{Text: "Discard another loot card?", Options: func(c *engine.Ctx, a []int) []string {
			left := goatLeft(c, a)
			if len(a) > 0 && (a[len(a)-1] < 0 || a[len(a)-1] == len(goatLeft(c, a[:len(a)-1]))) {
				return nil // done
			}
			var out []string
			for _, id := range left {
				out = append(out, string(c.G.Object(id).Card))
			}
			if out == nil {
				return nil
			}
			return append(out, "done")
		}})
	}
	return qs
}

// goatLeft is the hand without the cards picked so far.
func goatLeft(c *engine.Ctx, a []int) []engine.ObjectID {
	left := append([]engine.ObjectID(nil), c.G.Players[c.Controller].Hand...)
	for _, i := range a {
		if i < 0 || i >= len(left) {
			break
		}
		left = append(left[:i:i], left[i+1:]...)
	}
	return left
}

func goatHeadDiscard(c *engine.Ctx, a []int) {
	left := append([]engine.ObjectID(nil), c.G.Players[c.Controller].Hand...)
	n := 0
	for _, i := range a {
		if i < 0 || i >= len(left) {
			break
		}
		c.G.DiscardFromHand(c.Controller, left[i])
		left = append(left[:i:i], left[i+1:]...)
		n++
	}
	if n > 0 {
		c.Do(engine.Loot(n))
	}
}
