package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Dad’s Note (Wildcard Card)
//
//	Choose up to 3 non-event monster cards in discard. Put them in one or more monster slots not being attacked.
//
// Cards in discard can't be targets here; they are chosen as it resolves.
var dadsNote = engine.CardDef{
	Ref:    "dads_note",
	Kind:   engine.LootCard,
	Copies: 1,
	Abilities: []engine.Ability{{
		Kind:    engine.LootAbility,
		Text:    "Choose up to 3 non-event monster cards in discard. Put them in one or more monster slots not being attacked.",
		Effects: []engine.Effect{dadsNotePick(3)},
	}},
}

// dadsNotePick asks n times for a monster card in the discard (or done)
// and a slot for it, then places them in order.
func dadsNotePick(n int) engine.Effect {
	var qs []engine.Question
	for i := range n {
		qs = append(qs, engine.Question{Text: "Put which monster card from the discard in a slot?", Options: func(c *engine.Ctx, a []int) []string {
			if i > 0 && (a[2*i-2] < 0 || a[2*i-1] < 0) {
				return nil // done
			}
			left := discardMonstersLeft(c, a[:2*i])
			if len(left) == 0 {
				return nil
			}
			return append(labelsOf(c.G, left), "done")
		}}, slotQuestionAfter(2*i))
	}
	return engine.Ask(func(c *engine.Ctx, a []int) {
		left := discardMonsters(c.G)
		slots := freeSlots(c.G)
		for i := 0; i+1 < len(a); i += 2 {
			if a[i] < 0 || a[i] >= len(left) || a[i+1] < 0 || a[i+1] >= len(slots) {
				return
			}
			c.G.PlaceFromDiscard(left[a[i]], slots[a[i+1]])
			left = append(left[:a[i]:a[i]], left[a[i]+1:]...)
		}
	}, qs...)
}

// slotQuestionAfter asks for a slot once the card question before it
// (answer k) picked a card.
func slotQuestionAfter(k int) engine.Question {
	q := slotQuestion("Put it in which monster slot?")
	opts := q.Options
	q.Options = func(c *engine.Ctx, a []int) []string {
		if a[k] < 0 || a[k] >= len(discardMonstersLeft(c, a[:k])) {
			return nil
		}
		return opts(c, a)
	}
	return q
}

// discardMonsters lists the non-event monster cards in the monster discard.
func discardMonsters(g *engine.Game) []engine.ObjectID {
	var out []engine.ObjectID
	for _, id := range g.Discards[engine.MonsterDeck] {
		if g.Def(id).Kind == engine.MonsterCard {
			out = append(out, id)
		}
	}
	return out
}

// discardMonstersLeft is discardMonsters without the cards picked in the
// pairs of answers so far.
func discardMonstersLeft(c *engine.Ctx, a []int) []engine.ObjectID {
	left := discardMonsters(c.G)
	for i := 0; i+1 < len(a); i += 2 {
		if a[i] < 0 || a[i] >= len(left) {
			return left
		}
		left = append(left[:a[i]:a[i]], left[a[i]+1:]...)
	}
	return left
}
