package rulesengine

// Threatens reports whether stack item it may hurt p, so a player who
// waits for their turn must be asked (15.86). It does when it:
//   - targets p, or p's character, items, souls, curses, or hand cards
//     (15.84); a target already spent as a cost does not count;
//   - is damage to p or to an object p controls;
//   - may take from p with no target naming them (Ability.Takes, 15.85).
//
// An item of p's own, or of nobody (a death, rewards), never threatens.
func (g *Game) Threatens(it StackItem, p PlayerID) bool {
	if it.Controller == p || it.Controller == NoPlayer {
		return false
	}
	if it.Kind == StackDamage {
		if it.Target.IsPlayer && it.Target.Player == p {
			return true
		}
		if !it.Target.IsPlayer && g.controls(p, it.Target.Object) {
			return true
		}
	}
	for _, c := range it.Targets {
		if !c.Spent && g.chosenHits(c, p) {
			return true
		}
	}
	return g.takes(it, p)
}

// controls: id is an object p controls; a bad id is nobody's.
func (g *Game) controls(p PlayerID, id ObjectID) bool {
	return id >= 0 && int(id) < len(g.Objects) && g.Objects[id].Controller == p
}

// chosenHits: the chosen target is p, or an object p controls.
func (g *Game) chosenHits(c Chosen, p PlayerID) bool {
	switch c.Kind {
	case TargetPlayer, TargetOtherPlayer, TargetMonsterOrPlayer:
		return c.Player == p
	case TargetItem, TargetCurse, TargetYourItem, TargetYourHandCard, TargetItemOrSoul:
		return g.controls(p, c.Object)
	case TargetMonster, TargetDiceRoll, TargetStackAbility, TargetStackTrigger:
	}
	return false
}

// takes: its ability may take from p with no target naming p
// (Ability.Takes, 15.85).
func (g *Game) takes(it StackItem, p PlayerID) bool {
	d, ok := g.cards.find(it.Ability.Card)
	if !ok {
		return false
	}
	abs := d.Abilities
	if it.Ability.Lingering {
		abs = d.Lingering
	}
	if it.Ability.Index < 0 || it.Ability.Index >= len(abs) {
		return false
	}
	switch abs[it.Ability.Index].Takes {
	case TakesOthers:
		return true
	case TakesEventPlayer:
		return it.EventPlayer == p
	case TakesNobody:
	}
	return false
}
