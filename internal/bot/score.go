package bot

import (
	"math"
	"slices"

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
	// An average item is worth clearly more than its price (21.16
	// tuning, 12 before): cents are spent almost only on items.
	wItem  = 15.0
	wCent  = 1.0
	wLoot  = 0.75
	wCurse = -8.0
	// wHP is each HP left; wAlive is not having died this turn. Both
	// heal at the end of the turn (R-TURN-13, R-DEATH-19): they weigh
	// for the rest of it only. A death's lasting cost is its penalty,
	// which cents, loot and items count.
	wHP    = 0.5
	wAlive = 3.0

	// wHand is the share of its measured value a loot card is worth in
	// hand (Decision 20.8.6): a card in hand is not played yet: it can
	// be discarded or stolen, few loot plays a turn. Its measured value
	// is what playing it adds; at a share of it, playing a card that
	// does its usual work beats holding it.
	wHand = 0.5

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
	// noiseAttack: standard's attack on an empty stack clears no noise
	// (Decision 8.6.2). The spread of its settles is the attack's real
	// risk, not a measuring error: an attack that wins a soul in one
	// settle of six beats ending the turn on the mean, and one standard
	// error cut every such attack (21.5 tuning). Other moves keep
	// noiseOpen: with no bar, standard used items in a loop.
	noiseAttack = 0.0
	// noiseClip: in the hard level's noise a gain counts at most four
	// souls, the most a table holds before a win. A game won or lost
	// in one world of fifty is a real part of the mean, but as noise it
	// would hide what every other world shows (Decision 9.6).
	noiseClip = 4 * wSoul

	// wBossHP is each HP a challenge's final boss has left (Decision
	// 15.3), taken from the team's sum: as a player's HP, damage heals
	// at the end of the turn, but an attack on the boss settles to its
	// end, where only a kill lasts; four times wHP, as the boss's HP is
	// the team's goal.
	wBossHP = 2.0

	// The hard level's part (Decision 7.6). wStrength is the strength
	// one opponent's weight grows by e: an opponent one soul ahead of
	// another counts e times as much. No threshold: the weights move
	// smoothly with every part of each seat's sum.
	wStrength = wSoul
	// wHold is how much more an answer still held is worth than its
	// base value (cardValue): a card or charge the seat may play on
	// another seat's turn. Spent now, it is gone against a later
	// threat. With one opponent clearly the strongest the share grows
	// up to 1 + wHoldLead times: the threat to keep it for is near. An
	// item answer (wItem, 15) is so worth 3 to 9 more while it is ready.
	wHold     = 0.2
	wHoldLead = 2.0
)

// scoring is what a level adds to the shared score (Decision 7): the
// dig value of its world, and for hard the opponents weighted by
// strength with the holding value of its answers (Decision 7.6).
// The zero scoring is Score.
type scoring struct {
	// dig is digValue of the world: what a kill adds (Decision 7.9).
	dig float64
	// lead: the opponents count by strength, not as an average.
	lead bool
	// answers are the seat's cards and charges that can answer on
	// another seat's turn (answers); each one still held counts its
	// base value times 1 + holdShare. Only with lead.
	answers []engine.ObjectID
	// flat values every loot card and item alike (wLoot, wItem), not by
	// the measured table (Decision 20): the flat level, and the
	// measurement of the table itself (Decision 20.6).
	flat bool
}

// Score is the table g from seat p's side (Decision 7): a win is the
// top value and a loss the bottom one; else p's own sum less the
// opponents' average (Decision 7.5).
func Score(g *engine.Game, p engine.PlayerID) float64 {
	return scoreAs(g, p, scoring{})
}

// scoreAs is Score as a level scores: with sc.lead the opponents count
// by strength and p's answers still held add their holding value
// (Decision 7.6). A challenge is scored for the team (teamScore).
func scoreAs(g *engine.Game, p engine.PlayerID, sc scoring) float64 {
	if g.Over {
		for _, w := range g.Winners {
			if w == p {
				return scoreWin
			}
		}
		return -scoreWin
	}
	if g.Challenge != nil {
		return teamScore(g, p, sc)
	}
	if sc.lead {
		return leadScore(g, p, sc)
	}
	others, n := 0.0, 0
	for i := range g.Players {
		q := engine.PlayerID(i)
		if q != p {
			others += side(g, q, p, sc.flat)
			n++
		}
	}
	own := side(g, p, p, sc.flat)
	if n == 0 {
		return own
	}
	return own - others/float64(n)
}

// teamScore is a challenge's table (Decision 15): both players play
// together, so every seat's sum counts for p, the partner's hand as
// seat p sees it, less each final boss's HP left. No opponents, no
// leader; won and lost are the top and the bottom (scoreAs).
func teamScore(g *engine.Game, p engine.PlayerID, sc scoring) float64 {
	s := 0.0
	for i := range g.Players {
		s += side(g, engine.PlayerID(i), p, sc.flat)
	}
	for _, id := range g.MonstersInPlay() {
		if g.IsFinalBoss(id) {
			s -= wBossHP * float64(g.HP(id))
		}
	}
	return s
}

// leadScore is the hard level's table (Decision 7.6): p's own sum and
// the holding value of its answers still held, less the opponents'
// sums weighted by strength. Each opponent's weight is exp(sum /
// wStrength), shared out over the opponents: equal opponents count as
// their average, one far ahead counts nearly alone.
func leadScore(g *engine.Game, p engine.PlayerID, sc scoring) float64 {
	var buf [4]float64
	sides := buf[:0]
	for i := range g.Players {
		sides = append(sides, side(g, engine.PlayerID(i), p, sc.flat))
	}
	own := sides[p] + holdShare(lead(sides, p))*held(g, p, sc.answers, sc.flat)
	if len(sides) == 1 {
		return own
	}
	top := slices.Max(sides)
	others, weights := 0.0, 0.0
	for q, s := range sides {
		if engine.PlayerID(q) != p {
			w := math.Exp((s - top) / wStrength)
			others += w * s
			weights += w
		}
	}
	return own - others/weights
}

// lead is how clearly one opponent of p is the strongest seat: 0 when
// none is ahead of the rest, towards 1 as one pulls away. It is the
// strongest opponent's share of the strength weights of every seat,
// p's own too, past an even share.
func lead(sides []float64, p engine.PlayerID) float64 {
	n := len(sides)
	if n < 2 {
		return 0
	}
	top := slices.Max(sides)
	sum, best := 0.0, 0.0
	for q, s := range sides {
		w := math.Exp((s - top) / wStrength)
		sum += w
		if engine.PlayerID(q) != p {
			best = max(best, w)
		}
	}
	even := 1 / float64(n)
	return min(max((best/sum-even)/(1-even), 0), 1)
}

// holdShare is the multiplier of an answer held, past its base value,
// at this lead: an answer counts cardValue × (1 + holdShare).
func holdShare(lead float64) float64 {
	return wHold * (1 + wHoldLead*lead)
}

// held is the base value (cardValue) of the answers p still holds on
// table g: a card still in its hand; a card it still controls in play
// whose ability is ready: no ↷, or charged, or recharging now as p's
// turn starts (R-TURN-02). The character has no base value, so its ↷
// adds nothing; nor does a card measured below 0 (Decision 20.8.3): a
// held answer is never a loss.
func held(g *engine.Game, p engine.PlayerID, answers []engine.ObjectID, flat bool) float64 {
	v := 0.0
	recharge := g.Turn.Active == p && g.Turn.Step == engine.StepRecharge
	for _, id := range answers {
		o := g.Object(id)
		switch {
		case slices.Contains(g.Players[p].Hand, id):
			v += max(cardValue(g, id, true, flat), 0)
		case o.Zone.Kind == engine.ZoneInPlay && o.Controller == p && (o.Charged || recharge || !g.HasTap(id)):
			v += max(cardValue(g, id, true, flat), 0)
		}
	}
	return v
}

// answers lists seat p's cards and charges that can answer: what it may
// play or use on another seat's turn (Decision 7.6), from allowed, its
// Allowed on world w. On another seat's turn that is every loot card
// and card it may play or use now. On its own turn the abilities of
// cards it controls: they work on any seat's turn, but a room's only on
// the active player's (R-CARD-36), and a loot card needs a loot play,
// which another seat's turn does not give (R-CARD-08).
func answers(w *engine.Game, p engine.PlayerID, allowed []engine.Intent) []engine.ObjectID {
	other := w.Turn.Active != p
	var out []engine.ObjectID
	for _, in := range allowed {
		if len(in.Objects) == 0 || slices.Contains(out, in.Objects[0]) {
			continue
		}
		id := in.Objects[0]
		switch in.Kind { //nolint:exhaustive // only plays and uses answer
		case engine.IntentActivate:
			if other || w.Object(id).Role != engine.RoleRoom {
				out = append(out, id)
			}
		case engine.IntentPlayLoot:
			if other {
				out = append(out, id)
			}
		}
	}
	return out
}

// side is one seat's sum: souls, coins, loot cards, items, curses, HP
// and being alive (Decision 7.2 – 7.3), as seat by sees it: the cards
// of another seat's hand count their kind's mean (Decision 20.8.2).
func side(g *engine.Game, p, by engine.PlayerID, flat bool) float64 {
	pl := &g.Players[p]
	s := wSoul*float64(g.SoulValue(p)) + wCent*float64(pl.Cents) + wHP*float64(g.PlayerHP(p))
	if !pl.Dead {
		s += wAlive
	}
	for _, id := range pl.Hand {
		s += cardValue(g, id, p == by, flat)
	}
	for _, id := range pl.InPlay {
		s += cardValue(g, id, true, flat)
	}
	return s
}

// cardValue is the base value of one card of a seat: a loot card in
// its hand, an item or a curse it has in play (Decision 7.3). It is the
// one place a card's worth is set; the hard level's holding value
// multiplies it (Decision 7.6).
//
// A loot card in hand is worth wHand times its measured value, an item
// its measured value scaled to wItem on average (values.go, Decision
// 20.8); with flat, wLoot and wItem. seen: the scoring seat
// knows the card; a card in another seat's hand is not seen. In an
// imagined world it is a guess, which the bot cannot tell from a card
// the seat remembers, so it counts its kind's mean (Decision 20.8.2).
func cardValue(g *engine.Game, id engine.ObjectID, seen, flat bool) float64 {
	o := g.Object(id)
	if o.Zone.Kind == engine.ZoneHand {
		switch {
		case flat:
			return wLoot
		case !seen:
			return wHand * cardTable.mean(kindLoot)
		}
		return wHand * cardTable.value(g, id, kindLoot)
	}
	switch o.Role { //nolint:exhaustive // souls count by SoulValue; the rest is nothing
	case engine.RoleItem:
		if flat {
			return wItem
		}
		return cardTable.value(g, id, itemKind(g, id))
	case engine.RoleCurse:
		return wCurse
	}
	return 0
}

// digValue is what a monster kill is worth beyond its rewards: the
// chance that the card it lets in gives a soul (Decision 7.9).
func digValue(w *engine.Game) float64 {
	return wSoul * digShare * soulShare(w)
}

// soulShare is the share of soul cards (a card property: Soul) in the
// monster deck of world w, all cards the seat cannot see, redealt by
// Imagine from the real unseen ones; 0 for an empty deck. digValue and
// the policy's judge of the deck's top both use it.
func soulShare(w *engine.Game) float64 {
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
	return float64(souls) / float64(len(deck))
}
