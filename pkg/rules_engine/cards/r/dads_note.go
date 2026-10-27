package r

import (
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/cardkit"
)

// The card's own lines; their words are in en.json (ADR 021).
var (
	dadsNotePutWhichMonster = engine.NewKey("card.dads_note.put_which_monster_card_from")
)

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
		Effects: []engine.Effect{discardMonsterToSlot(3, true)},
	}},
}

// discardMonsterToSlot asks n times for a non-event monster card in the
// monster discard and a slot not being attacked for it, then places them
// in order. With upTo the player may stop early ("up to 3").
func discardMonsterToSlot(n int, upTo bool) engine.Effect {
	var qs []engine.Question
	for i := range n {
		qs = append(qs, engine.Question{Text: dadsNotePutWhichMonster, Options: func(c *engine.Ctx, a []int) []engine.Msg {
			if i > 0 && (a[2*i-2] < 0 || a[2*i-1] < 0) {
				return nil // done
			}
			left := discardMonstersLeft(c, a[:2*i])
			if len(left) == 0 {
				return nil
			}
			if !upTo {
				return labelsOf(c.G, left)
			}
			return append(labelsOf(c.G, left), engine.AskDone.Msg())
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
	q := slotQuestion(cardkit.AskPutItInWhichMonster)
	opts := q.Options
	q.Options = func(c *engine.Ctx, a []int) []engine.Msg {
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
