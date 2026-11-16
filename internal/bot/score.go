package bot

import (
	engine "github.com/FourSoulsRepo/rules_engine"
)

// The score's weights (ADR 030, Decision 7): what a table is worth to a
// seat. Every weight is here; the arena tunes them (Decision 7.7).
const (
	// scoreWin is a won game; a lost one is -scoreWin. Far above any
	// table.
	scoreWin = 1e6

	// A soul is worth far more than everything else together.
	wSoul = 100.0
	// An item costs 10¢ (shopPrice) and stays; a cent is worth one.
	wItem  = 12.0
	wCent  = 1.0
	wLoot  = 0.75
	wCurse = -8.0
	// wHP is each HP left; wAlive is not having died this turn. Both
	// heal at the end of the turn (R-TURN-13, R-DEATH-19): they weigh
	// for the rest of it only. A death's lasting cost is its penalty,
	// which cents, loot and items count.
	wHP    = 0.5
	wAlive = 3.0

	// digShare: a monster killed lets the next card of the monster deck
	// in, which gives a soul with the deck's share of soul cards. A
	// kill is worth wSoul times that share times digShare (Decision
	// 7.9): the soul is not had yet, and others may take it first.
	digShare = 0.5

	// wAct is how much a move must beat doing nothing (a pass, or an
	// end of turn in the seat's own action phase), besides the noise of
	// its draws: a move that gains nothing, or less than a quarter coin,
	// is not made.
	wAct = 0.25
	// noiseOpen and noiseStack are how many standard errors of its gain
	// a move must also clear: on an empty stack, and on a stack. A move
	// on a stack is judged through the dice of every item below it,
	// which its own die shifts, so its gain there is mostly noise.
	noiseOpen  = 1.0
	noiseStack = 2.0
)

// Score is the table g from seat p's side (Decision 7): a win is the
// top value and a loss the bottom one; else p's own sum less the
// opponents' average (Decision 7.5).
func Score(g *engine.Game, p engine.PlayerID) float64 {
	if g.Over {
		for _, w := range g.Winners {
			if w == p {
				return scoreWin
			}
		}
		return -scoreWin
	}
	others, n := 0.0, 0
	for i := range g.Players {
		q := engine.PlayerID(i)
		if q != p {
			others += side(g, q)
			n++
		}
	}
	own := side(g, p)
	if n == 0 {
		return own
	}
	return own - others/float64(n)
}

// side is one seat's sum: souls, coins, loot cards, items, curses, HP
// and being alive (Decision 7.2 – 7.3).
func side(g *engine.Game, p engine.PlayerID) float64 {
	pl := &g.Players[p]
	s := wSoul*float64(g.SoulValue(p)) + wCent*float64(pl.Cents) + wLoot*float64(len(pl.Hand)) +
		wHP*float64(g.PlayerHP(p))
	if !pl.Dead {
		s += wAlive
	}
	for _, id := range pl.InPlay {
		switch g.Object(id).Role { //nolint:exhaustive // souls count by SoulValue; the rest is nothing
		case engine.RoleItem:
			s += wItem
		case engine.RoleCurse:
			s += wCurse
		}
	}
	return s
}

// digValue is what a monster kill is worth beyond its rewards: the
// chance that the card it lets in gives a soul (Decision 7.9). The
// share is of soul cards (a card property: Soul) in the monster deck of
// world w, all cards the seat cannot see, redealt by Imagine from the
// real unseen ones.
func digValue(w *engine.Game) float64 {
	deck := w.Decks[engine.MonsterDeck]
	if len(deck) == 0 {
		return 0
	}
	souls := 0
	for _, id := range deck {
		if w.Def(id).Soul > 0 {
			souls++
		}
	}
	return wSoul * digShare * float64(souls) / float64(len(deck))
}
