package rulesengine

import "slices"

// Card hooks that guard monsters, force an attack at random, or take a
// dead monster's spoils away (ADR 028, Decision 6.4). Pilfering Krampus
// uses them (R-KRAMPUS-11, R-KRAMPUS-13, R-KRAMPUS-14); the engine
// knows no challenge by name.

// askAttackRandom is the option of an owed random attack (Q-31).
var askAttackRandom = NewKey("ask.attack_random")

// CanBeAttacked reports whether id may be an attack's target: it is
// neither "can't be attacked" (Unattackable, gained too: Q-47) nor
// guarded now (Guarded).
func (g *Game) CanBeAttacked(id ObjectID) bool {
	for _, face := range g.facesOf(id) {
		if d, ok := g.cards.find(face); ok && d.Unattackable {
			return false
		}
	}
	return !g.guarded(id)
}

// guarded: id is in play and its Guarded hook holds; it can't be
// attacked or take damage.
func (g *Game) guarded(id ObjectID) bool {
	f := g.def(id).Guarded
	return f != nil && g.Object(id).Zone.Kind == ZoneInPlay && f(g, id)
}

// noSpoils reports whether a card in play takes the rewards and the
// soul of the dying monster away (NoSpoils, R-KRAMPUS-14).
func (g *Game) noSpoils(dying ObjectID) bool {
	for _, id := range g.inPlay() {
		if f := g.def(id).NoSpoils; f != nil && f(g, id, dying) {
			return true
		}
	}
	return false
}

// owedTargets lists what the active player must attack now (R-ATK-03):
// the MustAttack card while it is in play (a covered card is not,
// R-ZONE-10) and can be attacked, and every attack target whose card
// is one of MustAttackOneOfCards (19.241). Those are named by card, so
// a boss that comes up uncovered as a new object (R-ZONE-01) still
// counts; one in a discard, a deck or souls is no attack target. With
// several owed, the attacker picks (R-ATK-03).
func (g *Game) owedTargets() []ObjectID {
	var out []ObjectID
	if id := g.Turn.MustAttack; id != 0 && g.Object(id).Zone.Kind == ZoneInPlay && g.CanBeAttacked(id) {
		out = append(out, id)
	}
	if len(g.Turn.MustAttackOneOfCards) == 0 {
		return out
	}
	for _, id := range g.attackTargets() {
		if !contains(out, id) && slices.Contains(g.Turn.MustAttackOneOfCards, g.Object(id).Card) {
			out = append(out, id)
		}
	}
	return out
}

// randomAttackTargets lists the monsters the active player must attack
// one of at random (MustAttackRandom) that can be attacked now; nil once
// the duty is met this turn (Q-31).
func (g *Game) randomAttackTargets() []ObjectID {
	if g.Turn.RandomAttackMade {
		return nil
	}
	var out []ObjectID
	for _, id := range g.inPlay() {
		f := g.def(id).MustAttackRandom
		if f == nil {
			continue
		}
		for _, t := range f(g, id) {
			if !contains(out, t) && g.Object(t).Zone.Kind == ZoneInPlay && g.CanBeAttacked(t) {
				out = append(out, t)
			}
		}
	}
	return out
}

// randomAttackDue: p is the active player, owes the random attack, and
// is able to attack something other than the monster deck (Q-31): the
// turn can't end.
func (g *Game) randomAttackDue(p PlayerID) bool {
	return p == g.Turn.Active && g.canAttack(p) && g.attacksLeft() > g.Turn.DeckAttacks && len(g.randomAttackTargets()) > 0
}

// attackRandom: the game picks the attack's target among ids at random
// (Q-31); the attacker picks nothing.
func (g *Game) attackRandom(ids []ObjectID) {
	if len(ids) == 0 {
		g.askAttackTarget() // none left to attack: the usual choice
		return
	}
	opts := make([]Msg, len(ids))
	for j, id := range ids {
		opts[j] = CardOption(g.Object(id).Card)
	}
	i := g.randomPick(opts)
	g.emit(Event{Kind: EvRandomPick, Player: g.Turn.Active, Amount: i, Label: opts[i], Object: ids[i]})
	g.attackChosen(ids[i])
}

// attackChosen starts the attack on id, unless another player may
// redirect it first (Mom's Eye Shadow).
func (g *Game) attackChosen(id ObjectID) {
	if q, ok := g.redirector(); ok && g.Object(id).Role == RoleMonster {
		g.askRedirect(q, id)
		return
	}
	g.startAttack(id)
}

func btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}
