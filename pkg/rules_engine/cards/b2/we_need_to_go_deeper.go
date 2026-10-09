package b2

import engine "github.com/FourSoulsRepo/rules_engine"

// We Need To Go Deeper! (Good Event Card)
//
//	Put any number of non-event monster cards in discard on top of the monster deck.
//	The active player may attack an additional time this turn.
var weNeedToGoDeeper = engine.CardDef{
	Ref:    "we_need_to_go_deeper",
	Kind:   engine.EventCard,
	Copies: 1,
	Abilities: []engine.Ability{
		{
			Kind:    engine.Triggered,
			Text:    "Put any number of non-event monster cards in discard on top of the monster deck. The active player may attack an additional time this turn.",
			Trigger: engine.WhenThisEntersPlay(),
			Effects: []engine.Effect{engine.Ask(deeper, deeperQuestions()...), engine.AddAttacks(1, engine.You)},
		},
	},
}

// deeperQuestions pick monster cards from the discard one at a time,
// until "done".
func deeperQuestions() []engine.Question {
	var qs []engine.Question
	for range 12 {
		qs = append(qs, engine.Question{Text: "Put another monster card on top of the deck?", Options: func(c *engine.Ctx, a []int) []string {
			if len(a) > 0 && (a[len(a)-1] < 0 || a[len(a)-1] == len(deeperLeft(c, a[:len(a)-1]))) {
				return nil // done
			}
			var out []string
			for _, id := range deeperLeft(c, a) {
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

func deeperLeft(c *engine.Ctx, a []int) []engine.ObjectID {
	left := nonEventDiscards(c.G)
	for _, i := range a {
		if i < 0 || i >= len(left) {
			break
		}
		left = append(left[:i:i], left[i+1:]...)
	}
	return left
}

func deeper(c *engine.Ctx, a []int) {
	left := nonEventDiscards(c.G)
	for _, i := range a {
		if i < 0 || i >= len(left) {
			break
		}
		c.G.DiscardToDeckTop(engine.MonsterDeck, left[i])
		left = append(left[:i:i], left[i+1:]...)
	}
}
