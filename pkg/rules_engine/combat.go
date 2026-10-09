package rulesengine

import "strconv"

// Target is what a stack item aims at.
type Target struct {
	Object ObjectID `json:"object,omitempty"`
	Player PlayerID `json:"player"`
	// IsPlayer is true when the target is Player, otherwise Object.
	IsPlayer bool `json:"is_player,omitempty"`
}

// AttackState tracks the attack in progress (R-ATK).
type AttackState struct {
	On      bool     `json:"on,omitempty"`      // an attack was declared
	Started bool     `json:"started,omitempty"` // a target was chosen
	Target  ObjectID `json:"target,omitempty"`
	// Revealed is the monster-deck card waiting for a slot (R-ATK-08).
	Revealed ObjectID `json:"revealed,omitempty"`
}

// Events about combat and death.
const (
	EvAttackDeclared EventKind = "attack_declared"
	EvAttackTarget   EventKind = "attack_target"
	EvAttackEnded    EventKind = "attack_ended"
	EvDamaged        EventKind = "damaged"
	EvDied           EventKind = "died"
	EvGainedSoul     EventKind = "gained_soul"
	EvGainedTreasure EventKind = "gained_treasure"
	EvDestroyed      EventKind = "destroyed"
	EvDeactivated    EventKind = "deactivated"
)

func (g *Game) def(id ObjectID) CardDef {
	d, _ := g.cards.find(g.Object(id).Card) //nolint:errcheck // every object's card has a definition (checked at setup)
	return d
}

// eternal reports whether an object is eternal: printed or gained
// (R-ABIL-18).
func (g *Game) eternal(id ObjectID) bool {
	return g.def(id).Eternal || g.Object(id).Eternal
}

// HP is an object's remaining health.
func (g *Game) HP(id ObjectID) int {
	return max(g.def(id).HP+g.bonus(StatMonsterHP, NoPlayer, id)-g.Object(id).Damage, 0)
}

// PlayerHP is a player's remaining health (R-CARD-25, R-MECH-19).
func (g *Game) PlayerHP(p PlayerID) int {
	return max(g.def(g.Players[p].Character).HP+g.bonus(StatPlayerHP, p, 0)-g.Players[p].Damage, 0)
}

// Evasion is the dice check to hit a monster, between 1 and 6 (R-ATK-18).
func (g *Game) Evasion(id ObjectID) int {
	return min(max(g.def(id).DC+g.bonus(StatMonsterDC, NoPlayer, id), 1), 6)
}

// PlayerATK is a player's attack (R-CARD-25).
func (g *Game) PlayerATK(p PlayerID) int {
	return max(g.def(g.Players[p].Character).ATK+g.bonus(StatPlayerATK, p, 0), 0)
}

// MonsterATK is a monster's attack.
func (g *Game) MonsterATK(id ObjectID) int {
	return max(g.def(id).ATK+g.bonus(StatMonsterATK, NoPlayer, id), 0)
}

// declareAttack: priority passes before a target is chosen (R-ATK-02).
func (g *Game) declareAttack(p PlayerID) {
	g.Turn.Attacks--
	g.Attack = AttackState{On: true}
	g.emit(Event{Kind: EvAttackDeclared, Player: p})
	g.openWindow(p)
}

// askAttackTarget offers the monsters in play and the monster deck.
func (g *Game) askAttackTarget() {
	var monsters []ObjectID
	for _, s := range g.Monsters {
		if top, ok := s.TopOf(); ok && g.Object(top).Role == RoleMonster {
			monsters = append(monsters, top)
		}
	}
	labels := append(g.labels(monsters), "monster deck")
	g.ask(Choice{Purpose: ChooseAttackTarget, Player: g.Turn.Active, Rule: "R-ATK-02", Objects: monsters, Deck: true}, labels)
}

// attackDeck reveals the top monster card; the attacker picks its slot
// (R-ATK-08).
func (g *Game) attackDeck() {
	id, ok := g.drawTop(MonsterDeck)
	if !ok {
		g.endAttack()
		return
	}
	g.Attack.Revealed = g.move(id, Zone{Kind: ZoneOutside}, NoPlayer)
	g.emit(Event{Kind: EvCardRevealed, Player: g.Turn.Active, Object: g.Attack.Revealed, Card: g.Object(g.Attack.Revealed).Card})
	slots := make([]int, len(g.Monsters))
	labels := make([]string, len(g.Monsters))
	for i, s := range g.Monsters {
		slots[i] = i
		labels[i] = "slot " + strconv.Itoa(i+1)
		if top, ok := s.TopOf(); ok {
			labels[i] += ": " + string(g.Object(top).Card)
		}
	}
	g.ask(Choice{Purpose: ChooseMonsterSlot, Player: g.Turn.Active, Rule: "R-ATK-08", Slots: slots}, labels)
}

// placeRevealed covers a monster slot with the revealed card (R-ATK-08).
func (g *Game) placeRevealed(slot int) {
	revealed := g.Attack.Revealed
	g.Attack.Revealed = 0
	if top, ok := g.Monsters[slot].TopOf(); ok {
		o := g.Object(top)
		o.Zone.Kind, o.Role = ZoneCovered, RoleNone // covered: not in play (R-ZONE-10)
	}
	id := g.putInSlot(revealed, MonsterSlot, slot)
	if g.Object(id).Role == RoleEvent {
		// Events trigger on entering play (step 4.7), then go to discard;
		// the attack is over (R-ATK-09).
		g.removeFromSlot(id)
		g.discard(id, MonsterDeck)
		g.endAttack()
		return
	}
	g.startAttack(id)
}

// startAttack begins attack rolls against a target (R-ATK-11).
func (g *Game) startAttack(target ObjectID) {
	g.Attack.Started, g.Attack.Target = true, target
	g.emit(Event{Kind: EvAttackTarget, Player: g.Turn.Active, Object: target, Card: g.Object(target).Card})
	g.attackRoll()
}

func (g *Game) attackRoll() {
	p := g.Turn.Active
	r := g.d6()
	g.emit(Event{Kind: EvDiceRolled, Player: p, Amount: r, Text: "attack roll"})
	g.push(StackItem{
		Kind: StackRoll, Controller: p, Roll: r, Label: "attack roll", Attack: true,
		Target: Target{Object: g.Attack.Target},
	})
}

// resolveAttackRoll: hit or miss (R-ATK-12, R-ATK-13).
func (g *Game) resolveAttackRoll(it StackItem) {
	if !g.Attack.Started || it.Target.Object != g.Attack.Target {
		return
	}
	p := g.Turn.Active
	if it.Roll >= g.Evasion(it.Target.Object) {
		if atk := g.PlayerATK(p); atk > 0 {
			g.push(StackItem{
				Kind: StackDamage, Controller: p, Amount: atk, Label: "combat damage", Attack: true,
				Target: Target{Object: it.Target.Object},
			})
		}
		return
	}
	// Nobody deals 0 damage (R-MECH-20).
	if atk := g.MonsterATK(it.Target.Object); atk > 0 {
		g.push(StackItem{
			Kind: StackDamage, Controller: NoPlayer, Source: it.Target.Object, Amount: atk,
			Label: "combat damage", Attack: true, Target: Target{Player: p, IsPlayer: true},
		})
	}
}

// continueAttack runs between attack rolls: the attack goes on until one
// side dies or has 0 HP (R-ATK-14).
func (g *Game) continueAttack() {
	t := g.Attack.Target
	if g.PlayerHP(g.Turn.Active) == 0 || g.Object(t).Zone.Kind != ZoneInPlay || g.HP(t) == 0 {
		g.endAttack()
		g.openWindow(g.Turn.Active)
		return
	}
	g.attackRoll()
}

// endAttack ends the attack and removes its unresolved rolls and combat
// damage from the stack (R-ATK-15).
func (g *Game) endAttack() {
	if !g.Attack.On {
		return
	}
	g.Attack = AttackState{}
	kept := g.Stack[:0]
	for _, it := range g.Stack {
		if !it.Attack {
			kept = append(kept, it)
		}
	}
	g.Stack = kept
	g.emit(Event{Kind: EvAttackEnded, Player: g.Turn.Active})
}

// resolveDamage marks damage; an object at 0 HP gets its death on the
// stack (R-MECH-15, R-MECH-16, R-DEATH-01).
func (g *Game) resolveDamage(it StackItem) {
	if g.useShield(it.Target) {
		return
	}
	if it.Target.IsPlayer {
		p := it.Target.Player
		hp := g.PlayerHP(p)
		if hp == 0 || g.Players[p].Dead {
			return
		}
		n := min(it.Amount, hp)
		g.Players[p].Damage += n
		g.emit(Event{Kind: EvDamaged, Player: p, Amount: n})
		if g.PlayerHP(p) == 0 {
			g.push(StackItem{Kind: StackDeath, Controller: NoPlayer, Label: "death", Target: Target{Player: p, IsPlayer: true}})
		}
		return
	}
	id := it.Target.Object
	o := g.Object(id)
	if o.Zone.Kind != ZoneInPlay || g.HP(id) == 0 {
		return
	}
	n := min(it.Amount, g.HP(id))
	o.Damage += n
	g.emit(Event{Kind: EvDamaged, Player: NoPlayer, Object: id, Card: o.Card, Amount: n})
	if g.HP(id) == 0 && !g.eternal(id) { // R-DEATH-03
		g.push(StackItem{Kind: StackDeath, Controller: NoPlayer, Label: "death", Target: Target{Object: id}})
	}
}

// useShield prevents damage to t if a shield protects it (R-MECH-46).
func (g *Game) useShield(t Target) bool {
	for i, s := range g.Shields {
		if s.IsPlayer == t.IsPlayer && ((t.IsPlayer && s.Player == t.Player) || (!t.IsPlayer && s.Object == t.Object)) {
			g.Shields = append(g.Shields[:i:i], g.Shields[i+1:]...)
			g.emit(Event{Kind: EvPrevented, Player: t.Player, Object: t.Object})
			return true
		}
	}
	return false
}

// resolveDeath: the object or player dies (R-DEATH-02).
func (g *Game) resolveDeath(it StackItem) {
	if it.Target.IsPlayer {
		if !g.Players[it.Target.Player].Dead {
			g.playerDeath(it.Target.Player)
		}
		return
	}
	if g.Object(it.Target.Object).Zone.Kind == ZoneInPlay {
		g.monsterDeath(it.Target.Object)
	}
}

// monsterDeath follows the monster death steps (R-DEATH-04 to R-DEATH-10).
// Steps after the first are queued actions: they do not use the stack.
func (g *Game) monsterDeath(id ObjectID) {
	card := g.Object(id).Card
	if g.Attack.Target == id {
		g.endAttack()
	}
	g.removeFromSlot(id)
	holding := g.move(id, Zone{Kind: ZoneOutside}, NoPlayer) // R-DEATH-04
	g.emit(Event{Kind: EvDied, Player: NoPlayer, Object: holding, Card: card})
	active := g.Turn.Active
	d := g.def(holding)
	for _, r := range d.Rewards { // R-DEATH-06
		g.enqueue(Action{Kind: r.action(), Player: active, Amount: r.Amount})
	}
	if d.Soul > 0 { // R-DEATH-08
		g.enqueue(Action{Kind: ActBecomeSoul, Player: active, Object: holding})
	} else {
		g.enqueue(Action{Kind: ActDiscardObject, Player: NoPlayer, Object: holding})
	}
	g.enqueue(Action{Kind: ActRefillSlots, Player: NoPlayer}) // R-DEATH-09
}

// playerDeath follows the player death steps (R-DEATH-12 to R-DEATH-16).
func (g *Game) playerDeath(p PlayerID) {
	g.Players[p].Dead = true // R-DEATH-17
	g.emit(Event{Kind: EvDied, Player: p, Object: g.Players[p].Character, Card: g.Object(g.Players[p].Character).Card})
	if p == g.Turn.Active {
		g.endAttack() // R-DEATH-12
		g.Turn.DeathEnd = true
	}
	g.enqueue( // death penalty (R-DEATH-14)
		Action{Kind: ActPenaltyItem, Player: p},
		Action{Kind: ActPenaltyLoot, Player: p},
		Action{Kind: ActLoseCents, Player: p, Amount: 1},
		Action{Kind: ActDeactivateTaps, Player: p},
	)
}

// removeFromSlot takes an object off the top of its slot; the card below,
// if any, enters play again as a new object (R-ZONE-11).
func (g *Game) removeFromSlot(id ObjectID) {
	z := g.Object(id).Zone
	var slot *Slot
	switch z.Slot {
	case ShopSlot:
		slot = &g.Shop[z.Index]
	case MonsterSlot:
		slot = &g.Monsters[z.Index]
	case RoomSlot:
		slot = &g.Rooms[z.Index]
	}
	slot.Cards = remove(slot.Cards, id)
	if below, ok := slot.TopOf(); ok && g.Object(below).Zone.Kind == ZoneCovered {
		slot.Cards = slot.Cards[:len(slot.Cards)-1]
		g.putInSlot(below, z.Slot, z.Index)
	}
}

// destroyItem destroys an item a player controls (R-MECH-23).
func (g *Game) destroyItem(p PlayerID, id ObjectID) {
	if g.eternal(id) {
		return // R-ABIL-18
	}
	g.Players[p].InPlay = remove(g.Players[p].InPlay, id)
	nid := g.discard(id, TreasureDeck)
	g.emit(Event{Kind: EvDestroyed, Player: p, Object: nid, Card: g.Object(nid).Card, Prev: id})
}

// refillSlots fills empty slots from their decks (R-SHOP-06); events met
// while refilling monster slots are resolved until a monster sits there
// (R-SHOP-08).
func (g *Game) refillSlots() {
	for i := range g.Shop {
		if len(g.Shop[i].Cards) == 0 {
			if id, ok := g.drawTop(TreasureDeck); ok {
				g.putInSlot(id, ShopSlot, i)
			}
		}
	}
	for i := range g.Monsters {
		for tries := 0; len(g.Monsters[i].Cards) == 0 && tries < 100; tries++ {
			id, ok := g.drawTop(MonsterDeck)
			if !ok {
				break
			}
			nid := g.putInSlot(id, MonsterSlot, i)
			if g.Object(nid).Role == RoleEvent {
				// Event abilities come with step 4.7; then it is discarded.
				g.removeFromSlot(nid)
				g.discard(nid, MonsterDeck)
			}
		}
	}
}

// d6 rolls a die; tests may force results.
func (g *Game) d6() int {
	if len(g.forcedRolls) > 0 {
		r := g.forcedRolls[0]
		g.forcedRolls = g.forcedRolls[1:]
		return r
	}
	return g.RNG.D6()
}
