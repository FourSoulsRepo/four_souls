package rulesengine

import (
	"slices"
	"strconv"
)

// Target is what a stack item aims at.
type Target struct {
	Object ObjectID `json:"object,omitempty"`
	Player PlayerID `json:"player"`
	// IsPlayer is true when the target is Player, otherwise Object.
	IsPlayer bool `json:"is_player,omitempty"`
}

// AttackState tracks the attack in progress (R-ATK).
type AttackState struct {
	// LastRoll is the latest attack roll's result.
	LastRoll int      `json:"last_roll,omitempty"`
	On       bool     `json:"on,omitempty"`      // an attack was declared
	Started  bool     `json:"started,omitempty"` // a target was chosen
	Target   ObjectID `json:"target,omitempty"`
	// Revealed is the monster-deck card waiting for a slot (R-ATK-08).
	Revealed ObjectID `json:"revealed,omitempty"`
	// Deck: the attack was declared on the monster deck; it is not an
	// attack on a monster even when a monster is revealed (R-ATK-05).
	Deck bool `json:"deck,omitempty"`
	// Helpers make an attack roll after each attack roll of the active
	// player, in turn order (R-ATK-21); Queue holds those still to roll
	// this round. Share also gain the rewards if the target dies this
	// attack (Sibling Rivalry).
	// Rolls counts the attack rolls made against the target this attack
	// ("every other attack roll made against it": Null).
	Rolls   int        `json:"rolls,omitempty"`
	Helpers []PlayerID `json:"helpers,omitempty"`
	Queue   []PlayerID `json:"queue,omitempty"`
	Share   []PlayerID `json:"share,omitempty"`
}

// Events about combat and death.
const (
	EvAttackDeclared EventKind = "attack_declared"
	EvAttackTarget   EventKind = "attack_target"
	EvAttackEnded    EventKind = "attack_ended"
	EvDamaged        EventKind = "damaged"
	EvDied           EventKind = "died"
	// EvObjectDied: an attackable item or room died (R-MECH-05); Text is
	// "item" or "room". Monster-death triggers don't see it.
	EvObjectDied     EventKind = "object_died"
	EvGainedSoul     EventKind = "gained_soul"
	EvGainedTreasure EventKind = "gained_treasure"
	// EvGainedItem: Player gained the item Object: from the treasure deck,
	// bought, searched for, or a trinket; never stolen or given
	// (R-MECH-35). "When you gain or purchase this" looks at it.
	EvGainedItem  EventKind = "gained_item"
	EvDestroyed   EventKind = "destroyed"
	EvDeactivated EventKind = "deactivated"
)

func (g *Game) def(id ObjectID) CardDef {
	d, _ := g.cards.find(g.CardOf(id)) //nolint:errcheck // every object's card has a definition (checked at setup)
	return d
}

// CardOf is the card an object acts as: the card it copies, if any, or
// its own (Diplopia, Modeling Clay). Where it goes when it leaves play
// still follows its own card.
func (g *Game) CardOf(id ObjectID) CardRef {
	if o := g.Object(id); o.CopyOf != "" {
		return o.CopyOf
	}
	return g.Object(id).Card
}

// Eternal reports whether an object is eternal: printed or gained
// (R-ABIL-18).
func (g *Game) Eternal(id ObjectID) bool {
	o := g.Object(id)
	return (g.def(id).Eternal && !o.LostEternal) || o.Eternal || g.boosted(StatEternal, id)
}

// boosted reports whether the object has a till-end-of-turn boost of stat.
func (g *Game) boosted(stat Stat, id ObjectID) bool {
	for _, b := range g.Boosts {
		if b.Stat == stat && b.Object == id {
			return true
		}
	}
	return false
}

// HP is an object's remaining health; an attacked player's character
// stands for the player (R-ATK-20).
func (g *Game) HP(id ObjectID) int {
	if p, ok := g.characterOf(id); ok {
		return g.PlayerHP(p)
	}
	return max(g.def(id).HP+g.monsterBonus(StatMonsterHP, id)-g.Object(id).Damage, 0)
}

// characterOf: id is a player's character; p is that player.
func (g *Game) characterOf(id ObjectID) (PlayerID, bool) {
	o := g.Object(id)
	if o.Role != RoleCharacter || o.Controller == NoPlayer || g.Players[o.Controller].Character != id {
		return NoPlayer, false
	}
	return o.Controller, true
}

// monsterBonus is bonus for a monster stat; items and rooms that can be
// attacked are not monsters (R-MECH-04).
func (g *Game) monsterBonus(stat Stat, id ObjectID) int {
	if g.Object(id).Role != RoleMonster {
		return 0
	}
	return g.bonus(stat, NoPlayer, id)
}

// AttackedPlayer is the player being attacked, if the attack is on a
// player (R-ATK-20).
func (g *Game) AttackedPlayer() (PlayerID, bool) {
	if !g.Attack.Started {
		return NoPlayer, false
	}
	return g.characterOf(g.Attack.Target)
}

// PlayerAttackDC is the evasion of player p for the active player's
// attack, and whether p may be attacked at all: the newest effect that
// allows it sets it (R-MECH-06).
func (g *Game) PlayerAttackDC(p PlayerID) (int, bool) {
	if p == g.Turn.Active || g.Players[p].Dead {
		return 0, false
	}
	for i := len(g.Boosts) - 1; i >= 0; i-- {
		if b := g.Boosts[i]; b.Stat == StatPlayerAttackDC && b.Object == 0 && b.Player == p {
			return b.Amount, true // an effect of this turn is the newest
		}
	}
	dc, newest := 0, ObjectID(-1)
	for _, id := range g.inPlay() {
		d, ok := g.cards.find(g.CardOf(id))
		if !ok {
			continue
		}
		for _, s := range d.Statics {
			if s.Stat == StatPlayerAttackDC && s.Applies(g, id, p, 0) && id > newest {
				dc, newest = s.amount(g, id), id
			}
		}
	}
	return dc, newest >= 0
}

// PlayerHP is a player's remaining health (R-CARD-25, R-MECH-19).
func (g *Game) PlayerHP(p PlayerID) int {
	return max(g.def(g.Players[p].Character).HP+g.bonus(StatPlayerHP, p, 0)-g.Players[p].Damage, 0)
}

// Evasion is the dice check to hit a monster, between 1 and 6 (R-ATK-18).
func (g *Game) Evasion(id ObjectID) int {
	if p, ok := g.characterOf(id); ok {
		dc, _ := g.PlayerAttackDC(p)
		return min(max(dc, 1), 6)
	}
	dc := g.def(id).DC + g.monsterBonus(StatMonsterDC, id)
	if g.Attack.On && g.Attack.Target == id {
		dc += g.bonus(StatAttackedDC, g.Turn.Active, 0) // "monsters you attack have -1 DC"
	}
	return min(max(dc, 1), 6)
}

// PlayerATK is a player's attack (R-CARD-25).
func (g *Game) PlayerATK(p PlayerID) int {
	return max(g.def(g.Players[p].Character).ATK+g.bonus(StatPlayerATK, p, 0), 0)
}

// MonsterATK is a monster's attack.
func (g *Game) MonsterATK(id ObjectID) int {
	if p, ok := g.characterOf(id); ok {
		return g.PlayerATK(p)
	}
	return max(g.def(id).ATK+g.monsterBonus(StatMonsterATK, id), 0)
}

// declareAttack: priority passes before a target is chosen (R-ATK-02).
func (g *Game) declareAttack(p PlayerID) {
	if g.Turn.MustAttacks > 0 {
		g.Turn.MustAttacks--
	}
	if g.Turn.DeckAttacks > 0 && g.attacksLeft() <= g.Turn.DeckAttacks {
		// Only deck attacks are left: this one is on the monster deck.
		g.Turn.DeckAttacks--
		g.Turn.MustAttackDeck++
	}
	if g.Turn.Attacks > 0 {
		g.Turn.Attacks--
	} else {
		g.Turn.BonusAttacksUsed++
	}
	g.Attack = AttackState{On: true}
	g.emit(Event{Kind: EvAttackDeclared, Player: p})
	g.openWindow(p)
}

// askAttackTarget offers the monsters in play and the monster deck.
func (g *Game) askAttackTarget() {
	if g.Turn.MustAttackDeck > 0 {
		g.Turn.MustAttackDeck--
		g.ask(Choice{Purpose: ChooseAttackTarget, Player: g.Turn.Active, Rule: "R-ATK-02", Deck: true}, []string{"monster deck"})
		return
	}
	if m := g.Turn.MustAttack; m != 0 && g.Object(m).Zone.Kind == ZoneInPlay {
		g.ask(Choice{Purpose: ChooseAttackTarget, Player: g.Turn.Active, Rule: "R-ATK-02", Objects: []ObjectID{m}}, g.targetLabels([]ObjectID{m}))
		return
	}
	targets := g.attackTargets()
	labels := append(g.targetLabels(targets), "monster deck")
	g.ask(Choice{Purpose: ChooseAttackTarget, Player: g.Turn.Active, Rule: "R-ATK-02", Objects: targets, Deck: true}, labels)
}

// attackTargets lists what the active player may attack besides the
// monster deck: monsters, attackable items and rooms (R-MECH-04), and
// players an effect lets them attack (R-ATK-20).
func (g *Game) attackTargets() []ObjectID {
	var out []ObjectID
	for _, s := range g.Monsters {
		if top, ok := s.TopOf(); ok && g.Object(top).Role == RoleMonster && !g.def(top).Unattackable {
			out = append(out, top)
		}
	}
	for _, pl := range g.Players {
		for _, id := range pl.InPlay {
			if g.Object(id).Role == RoleItem && g.def(id).Attackable {
				out = append(out, id)
			}
		}
	}
	for _, s := range g.Rooms {
		if top, ok := s.TopOf(); ok && g.def(top).Attackable {
			out = append(out, top)
		}
	}
	for _, p := range g.turnOrderFrom(g.Turn.Active)[1:] {
		if _, ok := g.PlayerAttackDC(p); ok {
			out = append(out, g.Players[p].Character)
		}
	}
	return out
}

// targetLabels names attack targets: a card, or "player N (character)".
func (g *Game) targetLabels(ids []ObjectID) []string {
	out := g.labels(ids)
	for i, id := range ids {
		if p, ok := g.characterOf(id); ok {
			out[i] = "player " + strconv.Itoa(int(p)+1) + " (" + string(g.Object(id).Card) + ")"
		}
	}
	return out
}

// attackDeck reveals the top monster card; the attacker picks its slot
// (R-ATK-08).
func (g *Game) attackDeck() {
	id, ok := g.drawTop(MonsterDeck)
	if !ok {
		g.endAttack()
		return
	}
	g.Attack.Revealed, g.Attack.Deck = g.move(id, Zone{Kind: ZoneOutside}, NoPlayer), true
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
	id := g.coverMonsterSlot(revealed, slot)
	if g.Object(id).Role == RoleEvent {
		// The event's abilities trigger; the attack is over (R-ATK-09).
		g.endAttack()
		g.enterMonsterSlot(id)
		return
	}
	g.startAttack(id)
}

// startAttack begins attack rolls against a target (R-ATK-11).
func (g *Game) startAttack(target ObjectID) {
	g.Attack.Started, g.Attack.Target = true, target
	// Sibling Rivalry may answer the declaration: its helpers join once
	// the target is chosen (S-YUG).
	g.AddAttackHelpers(g.Turn.NextHelpers...)
	g.Attack.Share = append(g.Attack.Share, g.Turn.NextShare...)
	g.Turn.NextHelpers, g.Turn.NextShare = nil, nil
	if g.Turn.MustAttack == target {
		g.Turn.MustAttack = 0 // the forced attack is made (Monster Manual, Krampus)
	}
	g.emit(Event{Kind: EvAttackTarget, Player: g.Turn.Active, Object: target, Card: g.Object(target).Card})
	g.attackRoll()
}

func (g *Game) attackRoll() { g.attackRollBy(g.Turn.Active) }

// attackRollBy pushes an attack roll of p against the attack's target;
// helpers roll too (R-ATK-21).
func (g *Game) attackRollBy(p PlayerID) {
	g.push(StackItem{
		Kind: StackRoll, Controller: p, Label: "attack roll", Attack: true,
		Target: Target{Object: g.Attack.Target},
	})
}

// resolveAttackRoll: hit or miss (R-ATK-12, R-ATK-13).
func (g *Game) resolveAttackRoll(it StackItem) {
	if !g.Attack.Started || it.Target.Object != g.Attack.Target {
		return
	}
	p := it.Controller
	if p == g.Turn.Active {
		g.Turn.AttackRolls++ // e.g. "your first attack roll each turn"
		g.Attack.Queue = g.helpersToRoll()
	}
	g.Attack.LastRoll = it.Roll
	g.Attack.Rolls++
	target := it.Target.Object
	tp, onPlayer := g.characterOf(target)
	mod := g.def(target).CombatMod
	if onPlayer {
		mod = nil // the attacked player's own cards change it below (R-ATK-24)
	}
	if it.Roll >= g.Evasion(target) {
		atk := g.PlayerATK(p)
		if mod != nil {
			atk = mod(g, target, it.Roll, true, atk)
		}
		atk = g.attackerCombatMod(p, it.Roll, true, atk)
		dt := Target{Object: target}
		if onPlayer {
			atk = g.attackerCombatMod(tp, it.Roll, false, atk) // damage the attacked player takes
			dt = Target{Player: tp, IsPlayer: true}
		}
		if atk > 0 {
			g.push(StackItem{
				Kind: StackDamage, Controller: p, Source: g.Players[p].Character, Amount: atk, Label: "combat damage", Attack: true,
				Target: dt,
			})
		}
		return
	}
	// Nobody deals 0 damage (R-MECH-20).
	atk := g.MonsterATK(target)
	if mod != nil {
		atk = mod(g, target, it.Roll, false, atk)
	}
	ctrl := NoPlayer
	if onPlayer {
		atk = g.attackerCombatMod(tp, it.Roll, true, atk) // damage the attacked player deals
		ctrl = tp
	}
	atk = g.attackerCombatMod(p, it.Roll, false, atk)
	if atk > 0 {
		g.push(StackItem{
			Kind: StackDamage, Controller: ctrl, Source: target, Amount: atk,
			Label: "combat damage", Attack: true, Target: Target{Player: p, IsPlayer: true},
		})
	}
}

// helpersToRoll lists the living helpers, in turn order from the active
// player, who roll after the active player's attack roll (R-ATK-21).
func (g *Game) helpersToRoll() []PlayerID {
	var out []PlayerID
	tp, _ := g.characterOf(g.Attack.Target)
	for _, p := range g.turnOrderFrom(g.Turn.Active)[1:] {
		if slices.Contains(g.Attack.Helpers, p) && !g.Players[p].Dead && p != tp {
			out = append(out, p)
		}
	}
	return out
}

// AddAttackHelpers makes players roll after each attack roll of the
// active player this attack (Team Up, R-ATK-21).
func (g *Game) AddAttackHelpers(ps ...PlayerID) {
	for _, p := range ps {
		if p != g.Turn.Active && !slices.Contains(g.Attack.Helpers, p) {
			g.Attack.Helpers = append(g.Attack.Helpers, p)
		}
	}
}

// HelpNextAttack makes p a helper of the active player's next attack
// this turn; with share p also gains the rewards if the target dies
// (Sibling Rivalry).
func (g *Game) HelpNextAttack(p PlayerID, share bool) {
	if p == g.Turn.Active {
		return
	}
	g.Turn.NextHelpers = append(g.Turn.NextHelpers, p)
	if share {
		g.Turn.NextShare = append(g.Turn.NextShare, p)
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
	for len(g.Attack.Queue) > 0 {
		p := g.Attack.Queue[0]
		g.Attack.Queue = g.Attack.Queue[1:]
		if !g.Players[p].Dead && g.PlayerHP(p) > 0 { // R-ATK-16
			g.attackRollBy(p)
			return
		}
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
	if it.Amount = g.shield(it.Target, it.Amount); it.Amount <= 0 {
		return
	}
	for _, id := range g.inPlay() { // e.g. Dry Baby, The Dead Cat
		if mod := g.def(id).DamageMod; mod != nil {
			if it.Amount = mod(g, id, it.Target, it.Amount, it.Attack); it.Amount <= 0 {
				g.emit(Event{Kind: EvPrevented, Player: it.Target.player(), Object: it.Target.Object, Source: id})
				return
			}
		}
	}
	if it.Target.IsPlayer {
		p := it.Target.Player
		hp := g.PlayerHP(p)
		if hp == 0 || g.Players[p].Dead {
			// The player still takes the damage: nothing is marked and they
			// can't die again, but "when you take damage" triggers (R-MECH-51).
			g.Players[p].TimesDamaged++
			g.emit(Event{Kind: EvDamaged, Player: p, Amount: it.Amount, Source: it.Source})
			return
		}
		n := min(it.Amount, hp)
		g.Players[p].Damage += n
		g.Players[p].TimesDamaged++
		e := Event{Kind: EvDamaged, Player: p, Amount: n, Source: it.Source, By: int(it.Controller) + 1}
		if it.Attack {
			e.Text = "combat"
		}
		g.emit(e)
		if g.PlayerHP(p) == 0 {
			g.Players[p].KilledBy = int(it.Controller) + 1
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
	o.HitsThisTurn++
	if g.HP(id) == 0 {
		o.KilledBy = int(it.Controller) + 1
		if it.Attack {
			o.KilledOn = g.Attack.LastRoll
		}
	}
	e := Event{Kind: EvDamaged, Player: NoPlayer, Object: id, Card: o.Card, Amount: n, Source: it.Source, By: int(it.Controller) + 1}
	if it.Attack {
		e.Text = "combat" // dealt by the active player's attack
	}
	g.emit(e)
	if g.HP(id) == 0 && !g.Eternal(id) { // R-DEATH-03
		g.push(StackItem{Kind: StackDeath, Controller: NoPlayer, Label: "death", Target: Target{Object: id}})
	}
}

// shield uses the first shield that protects t against n damage and
// returns the damage left (R-MECH-46).
func (g *Game) shield(t Target, n int) int {
	for i, s := range g.Shields {
		st := s.Target
		if st.IsPlayer != t.IsPlayer || (t.IsPlayer && st.Player != t.Player) || (!t.IsPlayer && st.Object != t.Object) {
			continue
		}
		g.Shields = append(g.Shields[:i:i], g.Shields[i+1:]...)
		prevented := n
		switch {
		case s.Cap > 0:
			prevented = max(n-s.Cap, 0)
		case s.Amount > 0:
			prevented = min(n, s.Amount)
		}
		g.emit(Event{Kind: EvPrevented, Player: t.Player, Object: t.Object, Amount: prevented, Source: s.Source})
		return n - prevented
	}
	return n
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

// deathText marks the death event of an item or room that could be
// attacked; a monster's death has none (R-MECH-05).
func (g *Game) deathText(id ObjectID) string {
	switch g.kindOf(id) { //nolint:exhaustive // only these two can be attacked besides monsters
	case TreasureCard:
		return "item"
	case RoomCard:
		return "room"
	}
	return ""
}

// monsterDeath follows the monster death steps (R-DEATH-04 to R-DEATH-10).
// Steps after the first are queued actions: they do not use the stack.
func (g *Game) monsterDeath(id ObjectID) {
	o := g.Object(id)
	card, killedBy := o.Card, o.KilledBy
	var share []PlayerID
	if g.Attack.Target == id {
		share = g.Attack.Share // Sibling Rivalry: died this attack
		g.endAttack()
	}
	text := g.deathText(id)
	inSlot := o.Zone.Slot == MonsterSlot || o.Zone.Slot == RoomSlot
	var z Zone
	if inSlot && o.Controller == NoPlayer {
		z = g.takeFromSlot(id)
	} else if o.Controller != NoPlayer {
		// An attackable item: it is destroyed (R-MECH-05).
		g.Players[o.Controller].InPlay = remove(g.Players[o.Controller].InPlay, id)
		inSlot = false
	}
	holding := g.move(id, Zone{Kind: ZoneOutside}, NoPlayer) // R-DEATH-04
	h := g.Object(holding)
	h.KilledBy, h.ShareRewards = killedBy, share
	before := len(g.PendingTriggers)
	kind := EvDied
	if text != "" {
		kind = EvObjectDied
	}
	g.emit(Event{Kind: kind, Player: NoPlayer, Object: holding, Card: card, Prev: id, Text: text})
	if text == "item" {
		g.emit(Event{Kind: EvDestroyed, Player: o.Controller, Object: holding, Card: card, Prev: id})
	}
	triggered := len(g.PendingTriggers) > before
	if inSlot {
		// The card below was covered when this died: it misses the death (R-ZONE-11).
		g.uncover(z)
	}
	if triggered {
		// "When this dies" triggers resolve before the rewards (R-DEATH-05).
		g.push(StackItem{Kind: StackDeathStep, Controller: NoPlayer, Label: "rewards", Target: Target{Object: holding}, Amount: stepRewards})
		return
	}
	g.rewardsStep(holding)
}

// The monster death steps that can wait under triggers on the stack.
const (
	stepRewards = 1 // R-DEATH-06
	stepSoul    = 2 // R-DEATH-08, R-DEATH-09
)

// rewardsStep gives the active player the rewards (R-DEATH-06); then
// "after rewards" triggers trigger (R-DEATH-07).
func (g *Game) rewardsStep(holding ObjectID) {
	active := g.receiver(holding)
	times := 1
	if g.Object(holding).DoubleRewards || g.bonus(StatDoubleRewards, NoPlayer, 0) > 0 {
		times = 2 // e.g. Dinga; doubled once however many apply (R-ABIL-36)
	}
	for _, r := range g.def(holding).Rewards {
		if r.Roll {
			g.push(StackItem{Kind: StackRoll, Controller: active, Source: holding, Label: "reward roll", Reward: int(r.Kind) + 1, Amount: times})
			continue
		}
		g.enqueue(Action{Kind: r.action(), Player: active, Amount: r.Amount * times})
	}
	for _, p := range g.Object(holding).ShareRewards {
		// Sibling Rivalry: the helper gains the rewards too.
		for _, r := range g.def(holding).Rewards {
			if !r.Roll {
				g.enqueue(Action{Kind: r.action(), Player: p, Amount: r.Amount * times})
			}
		}
	}
	g.enqueue(Action{Kind: ActRewardsDone, Player: active, Object: holding})
}

// receiver gets a dead card's rewards and soul: the active player
// (R-CARD-13), or its killer for "The player who kills this gains its
// rewards" (Ultra Greed).
func (g *Game) receiver(holding ObjectID) PlayerID {
	if k := g.Object(holding).Killer(); g.def(holding).RewardsToKiller && k != NoPlayer {
		return k
	}
	return g.Turn.Active
}

// soulStep makes the dead monster a soul of the active player or puts it
// into discard, unless an effect moved it already; slots refill
// (R-DEATH-08, R-DEATH-09).
func (g *Game) soulStep(holding ObjectID) {
	if g.Object(holding).Zone.Kind == ZoneOutside {
		if g.def(holding).Soul > 0 {
			g.enqueue(Action{Kind: ActBecomeSoul, Player: g.receiver(holding), Object: holding})
		} else {
			g.enqueue(Action{Kind: ActDiscardObject, Player: NoPlayer, Object: holding})
		}
	}
	g.enqueue(Action{Kind: ActRefillSlots, Player: NoPlayer})
}

// playerDeath follows the player death steps (R-DEATH-12 to R-DEATH-16).
func (g *Game) playerDeath(p PlayerID) {
	g.Players[p].Dead = true // R-DEATH-17
	defer g.dropCurses(p)
	g.emit(Event{Kind: EvDied, Player: p, Object: g.Players[p].Character, Card: g.Object(g.Players[p].Character).Card})
	if p == g.Turn.Active {
		g.endAttack() // R-DEATH-12
		g.Turn.DeathEnd = true
	}
	if tp, ok := g.AttackedPlayer(); ok && tp == p {
		g.endAttack() // killing the attacked player ends the attack (R-ATK-20)
	}
	if n := len(g.PendingTriggers); n > 0 && g.triggeredBy(EvDied, p) {
		// "When a player dies" triggers resolve first: the penalty waits
		// under them on the stack (R-DEATH-13).
		g.push(StackItem{Kind: StackPenalty, Controller: NoPlayer, Label: "death penalty", Target: Target{Player: p, IsPlayer: true}})
		return
	}
	g.payPenalty(p)
}

// triggeredBy reports whether a waiting trigger came from this player's
// death event.
func (g *Game) triggeredBy(kind EventKind, p PlayerID) bool {
	for _, t := range g.PendingTriggers {
		if t.On == kind && t.EventPlayer == p {
			return true
		}
	}
	return false
}

// payPenalty queues the death penalty (R-DEATH-14); "after paying
// penalties" triggers look at the event that follows it.
func (g *Game) payPenalty(p PlayerID) {
	g.enqueue(
		Action{Kind: ActPenaltyItem, Player: p},
		Action{Kind: ActPenaltyLoot, Player: p},
		g.penaltyCents(p),
		Action{Kind: ActDeactivateTaps, Player: p},
		Action{Kind: ActPenaltyDone, Player: p},
	)
}

// removeFromSlot takes an object off the top of its slot; the card below,
// if any, enters play again as a new object (R-ZONE-11).
func (g *Game) removeFromSlot(id ObjectID) {
	g.uncover(g.takeFromSlot(id))
}

// takeFromSlot takes an object off the top of its slot and returns the
// slot's zone; the card below stays covered until uncover.
func (g *Game) takeFromSlot(id ObjectID) Zone {
	z := g.Object(id).Zone
	slot := g.slotAt(z)
	slot.Cards = remove(slot.Cards, id)
	return z
}

// uncover puts the covered card on top of the slot at z into play.
func (g *Game) uncover(z Zone) {
	slot := g.slotAt(z)
	if below, ok := slot.TopOf(); ok && g.Object(below).Zone.Kind == ZoneCovered {
		slot.Cards = slot.Cards[:len(slot.Cards)-1]
		g.putInSlot(below, z.Slot, z.Index)
	}
}

func (g *Game) slotAt(z Zone) *Slot {
	switch z.Slot {
	case ShopSlot:
		return &g.Shop[z.Index]
	case MonsterSlot:
		return &g.Monsters[z.Index]
	default:
		return &g.Rooms[z.Index]
	}
}

// destroyItem destroys an item a player controls (R-MECH-23). It goes to
// the discard of its own deck: a trinket is a loot card.
func (g *Game) destroyItem(p PlayerID, id ObjectID) { g.DestroyObject(p, id) }

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
			if g.def(id).Indomitable && g.Turn.Number > 0 {
				g.putInSlot(id, MonsterSlot, g.newMonsterSlot()) // R-ABIL-23; slot i still needs a card
				continue
			}
			nid := g.putInSlot(id, MonsterSlot, i)
			if g.Object(nid).Role == RoleEvent {
				g.enterMonsterSlot(nid) // it stays until its abilities are done
				break
			}
		}
	}
}

// d6 draws a die's value when it is thrown; tests may force results.
func (g *Game) d6() int {
	if len(g.forcedRolls) > 0 {
		r := g.forcedRolls[0]
		g.forcedRolls = g.forcedRolls[1:]
		return r
	}
	return g.RNG.D6()
}

// attacksLeft is how many attacks the active player may still declare.
func (g *Game) attacksLeft() int {
	return g.Turn.Attacks + max(g.bonus(StatAttacks, g.Turn.Active, 0)-g.Turn.BonusAttacksUsed, 0)
}

// shadowOf is the first other player, in turn order from p, whose card
// takes p's death penalty (Shadow), or NoPlayer.
func (g *Game) shadowOf(p PlayerID) PlayerID {
	for _, q := range g.turnOrderFrom(p)[1:] {
		for _, id := range g.Players[q].InPlay {
			if g.Object(id).Role == RoleItem && g.def(id).TakesPenalties {
				return q
			}
		}
	}
	return NoPlayer
}

// penaltyCents is the penalty's 1¢: lost, or gained by a Shadow.
func (g *Game) penaltyCents(p PlayerID) Action {
	if s := g.shadowOf(p); s != NoPlayer && !g.protectedFromLoss(p) {
		return Action{Kind: ActStealCents, Player: s, From: p, Amount: 1}
	}
	return Action{Kind: ActLoseCents, Player: p, Amount: 1}
}

// Def returns the card definition an object acts as.
func (g *Game) Def(id ObjectID) CardDef { return g.def(id) }

// attackerCombatMod applies the CombatMod of the attacker's character and
// items: "You take no combat damage on attack rolls of 1" (Infamy).
func (g *Game) attackerCombatMod(p PlayerID, roll int, hit bool, n int) int {
	ids := append([]ObjectID{g.Players[p].Character}, g.Players[p].InPlay...)
	for _, id := range ids {
		if o := g.Object(id); o.Role == RoleSoul {
			continue
		}
		if mod := g.def(id).CombatMod; mod != nil {
			n = mod(g, id, roll, hit, n)
		}
	}
	return n
}

// redirector is another player who may choose which monster the active
// player attacks (Mom's Eye Shadow), first in turn order.
func (g *Game) redirector() (PlayerID, bool) {
	for _, p := range g.turnOrderFrom(g.Turn.Active) {
		if p != g.Turn.Active && g.bonus(StatChooseAttack, p, 0) > 0 {
			return p, true
		}
	}
	return NoPlayer, false
}

// askRedirect lets q choose the monster the attack goes to; the declared
// one is the first option. Only attacks on a monster, never the deck
// (S-YUG).
func (g *Game) askRedirect(q PlayerID, declared ObjectID) {
	options := []ObjectID{declared}
	for _, s := range g.Monsters {
		if top, ok := s.TopOf(); ok && top != declared && g.Object(top).Role == RoleMonster && !g.def(top).Unattackable {
			options = append(options, top)
		}
	}
	g.ask(Choice{Purpose: ChooseRedirect, Player: q, Rule: "R-ATK-03", Objects: options, Question: "Which monster do they attack?"}, g.labels(options))
}

// GivePenaltyItem makes p, if p dies this turn, give player to the item
// they would destroy for the death penalty instead ("they give the active player
// the item they would destroy": Lust for Blood).
func (g *Game) GivePenaltyItem(p, to PlayerID) {
	g.Boosts = append(g.Boosts, Boost{Stat: StatPenaltyItemTo, Player: p, Amount: int(to) + 1})
}

// penaltyItemTo is who gets p's death-penalty item this turn, or NoPlayer.
func (g *Game) penaltyItemTo(p PlayerID) PlayerID {
	for _, b := range g.Boosts {
		if b.Stat == StatPenaltyItemTo && b.Object == 0 && b.Player == p {
			return PlayerID(b.Amount - 1)
		}
	}
	return NoPlayer
}
