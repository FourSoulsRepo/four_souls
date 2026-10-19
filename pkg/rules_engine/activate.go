package rulesengine

import "strconv"

// Activation is an ability (or loot card) being put on the stack while its
// targets are chosen (R-ABIL-04).
type Activation struct {
	Player  PlayerID   `json:"player"`
	Source  ObjectID   `json:"source"`
	Ability AbilityRef `json:"ability"`
	Loot    bool       `json:"loot,omitempty"` // a loot card from hand
	// Mode is the chosen "choose one-" option; -1 until it is chosen.
	Mode   int      `json:"mode"`
	Chosen []Chosen `json:"chosen,omitempty"`
	// Window is the priority holder to return to if the player cancels.
	Window PlayerID `json:"window"`
}

// Events about abilities.
const (
	EvActivated EventKind = "activated"
	EvCancelled EventKind = "cancelled"
)

// activatable returns why p cannot activate ability i of src, or nil.
func (g *Game) activatable(p PlayerID, src ObjectID, i int) error {
	// A client may send any ID; an unknown one must not panic (A-16).
	if src < 0 || int(src) >= len(g.Objects) {
		return refuse("R-ABIL-09", "only the controller uses an item's or character's abilities")
	}
	o := g.Object(src)
	if o.Zone.Kind == ZoneInPlay && o.Role == RoleRoom {
		if p != g.Turn.Active {
			return deny("R-CARD-36", "Only on your turn", "only the active player uses a room's abilities")
		}
		if contains(g.Turn.LateRooms, src) {
			return deny("R-TURN-18", "Usable from next turn", "a room that entered in the end phase can't be used this turn")
		}
	} else if o.Zone.Kind != ZoneInPlay || o.Controller != p || (o.Role != RoleItem && o.Role != RoleCharacter && !g.soulKeepsAbilities(src)) {
		return refuse("R-ABIL-09", "only the controller uses an item's or character's abilities")
	}
	if g.bonus(StatOneAbility, p, 0) > 0 && g.Players[p].Activations > 0 {
		return deny("R-ABIL-12", "Only one ability each turn", "only one ability each turn")
	}
	refs := g.AbilitiesOf(src)
	if i < 0 || i >= len(refs) || g.ability(refs[i]).Kind != Activated {
		return refuse("R-ABIL-08", "%s has no such activated ability", o.Card)
	}
	a := g.ability(refs[i])
	for _, c := range a.Costs {
		if !c.canPay(g, p, src) {
			return deny("R-ABIL-07", costSay(c), "cannot pay the cost (%s)", c.label())
		}
	}
	return g.targetsAvailable(a, p, src)
}

// targetsAvailable checks every target spec has something to choose; an
// ability with modes needs one mode whose targets are there.
func (g *Game) targetsAvailable(a Ability, p PlayerID, src ObjectID) error {
	if len(a.Modes) > 0 {
		if len(g.availableModes(a, p, src)) == 0 {
			return deny("R-ABIL-06", "No valid target now", "no option has a valid target")
		}
		return nil
	}
	if !g.hasTargets(a.Targets, p, src) {
		return deny("R-ABIL-06", "No valid target now", "no valid target")
	}
	return nil
}

func (g *Game) hasTargets(specs []TargetSpec, p PlayerID, src ObjectID) bool {
	for _, t := range specs {
		if opts, _ := g.targetOptions(t, p, src); len(opts) == 0 {
			return false
		}
	}
	return true
}

// availableModes lists the modes of a whose targets can be chosen.
func (g *Game) availableModes(a Ability, p PlayerID, src ObjectID) []int {
	var out []int
	for i, m := range a.Modes {
		if g.hasTargets(m.Targets, p, src) {
			out = append(out, i)
		}
	}
	return out
}

// startActivation asks for the mode and targets one by one, then pays
// and pushes.
func (g *Game) startActivation(a Activation) {
	a.Window = g.Priority.Holder
	a.Mode = -1
	g.Activating = &a
	ab := g.ability(a.Ability)
	if len(ab.Modes) == 0 {
		a.Mode = 0
		g.nextTarget()
		return
	}
	modes := g.availableModes(ab, a.Player, a.Source)
	labels := make([]string, 0, len(modes)+1)
	for _, m := range modes {
		labels = append(labels, ab.Modes[m].Text)
	}
	labels = append(labels, "cancel")
	g.ask(Choice{Purpose: ChooseMode, Player: a.Player, Rule: "R-ABIL-04", Indexes: modes}, labels)
}

// chooseMode records the mode, or cancels before anything is paid.
func (g *Game) chooseMode(c Choice, i int) {
	if i == len(c.Indexes) {
		g.cancelActivation()
		return
	}
	g.Activating.Mode = c.Indexes[i]
	g.nextTarget()
}

// targetSpecs are the targets to choose for an ability and mode.
func (g *Game) targetSpecs(ref AbilityRef, mode int) []TargetSpec {
	ab := g.ability(ref)
	if len(ab.Modes) > 0 {
		return ab.Modes[mode].Targets
	}
	return ab.Targets
}

func (g *Game) nextTarget() {
	a := g.Activating
	specs := g.targetSpecs(a.Ability, a.Mode)
	if len(a.Chosen) < len(specs) {
		opts, labels := g.targetOptions(specs[len(a.Chosen)], a.Player, a.Source)
		labels = append(labels, "cancel")
		g.ask(Choice{Purpose: ChooseTarget, Player: a.Player, Rule: "R-ABIL-04", Targets: opts}, labels)
		return
	}
	g.finishActivation()
}

// chooseTarget records a target, or cancels before anything is paid.
func (g *Game) chooseTarget(c Choice, i int) {
	a := g.Activating
	if i == len(c.Targets) {
		g.cancelActivation()
		return
	}
	a.Chosen = append(a.Chosen, c.Targets[i])
	g.nextTarget()
}

func (g *Game) cancelActivation() {
	a := g.Activating
	g.Activating = nil
	g.emit(Event{Kind: EvCancelled, Player: a.Player, Object: a.Source})
	g.openWindow(a.Window)
}

func (g *Game) finishActivation() {
	a := *g.Activating
	g.Activating = nil
	ab := g.ability(a.Ability)
	if a.Loot {
		g.Players[a.Player].Hand = remove(g.Players[a.Player].Hand, a.Source)
		g.useLootPlay(a.Player)
		nid := g.move(a.Source, Zone{Kind: ZoneStack}, a.Player)
		card := g.Object(nid).Card
		g.Players[a.Player].LootPlayed++
		g.emit(Event{Kind: EvLootPlayed, Player: a.Player, Object: nid, Card: card})
		g.push(StackItem{Kind: StackLoot, Controller: a.Player, Source: nid, Card: card, Ability: a.Ability, Mode: a.Mode, Targets: a.Chosen})
		if pl := &g.Players[a.Player]; pl.CopyNextLoot && !g.def(nid).Trinket {
			// A copy of the loot goes on the stack above it (Blank Card).
			pl.CopyNextLoot = false
			g.push(StackItem{
				Kind: StackAbility, Controller: a.Player, Source: nid, Card: card, Ability: a.Ability,
				Mode: a.Mode, Targets: a.Chosen, Label: "copy of " + string(card),
			})
		}
		return
	}
	for _, c := range ab.Costs {
		if cc, ok := c.(chosenCost); ok {
			cc.payChosen(g, a.Player, a.Chosen)
			continue
		}
		c.pay(g, a.Player, a.Source)
	}
	g.Players[a.Player].Activations++
	g.emit(Event{Kind: EvActivated, Player: a.Player, Object: a.Source, Card: a.Ability.Card, Text: ab.Text})
	g.push(StackItem{
		Kind: StackAbility, Controller: a.Player, Source: a.Source, Card: a.Ability.Card,
		Ability: a.Ability, Mode: a.Mode, Targets: a.Chosen, Label: ab.Text,
	})
}

// targetOptions lists what may be chosen for a spec, with labels.
func (g *Game) targetOptions(t TargetSpec, p PlayerID, src ObjectID) ([]Chosen, []string) {
	var opts []Chosen
	var labels []string
	others := t.Kind == TargetOtherPlayer
	addPlayers := func() {
		var eligible []PlayerID
		for _, pl := range g.Players {
			if !others || pl.ID != p {
				eligible = append(eligible, pl.ID)
			}
		}
		for _, id := range g.choosablePlayers(eligible) {
			opts = append(opts, Chosen{Kind: TargetPlayer, Player: id})
			labels = append(labels, "player "+strconv.Itoa(int(id)+1)+" ("+string(g.Object(g.Players[id].Character).Card)+")")
		}
	}
	addMonsters := func() {
		for _, s := range g.Monsters {
			if top, ok := s.TopOf(); ok && g.Object(top).Role == RoleMonster {
				opts = append(opts, Chosen{Kind: TargetMonster, Player: NoPlayer, Object: top})
				labels = append(labels, string(g.Object(top).Card))
			}
		}
	}
	switch t.Kind {
	case TargetPlayer, TargetOtherPlayer:
		addPlayers() // the chosen kind is TargetPlayer either way
	case TargetMonster:
		addMonsters()
	case TargetMonsterOrPlayer:
		addMonsters()
		addPlayers()
	case TargetItem:
		for _, pl := range g.Players {
			for _, id := range pl.InPlay {
				if g.Object(id).Role == RoleItem {
					opts = append(opts, Chosen{Kind: TargetItem, Player: pl.ID, Object: id})
					labels = append(labels, string(g.Object(id).Card))
				}
			}
		}
		// Shop items are items in play too (R-CARD-01, R-ZONE-02).
		for _, id := range g.ShopItems() {
			opts = append(opts, Chosen{Kind: TargetItem, Player: NoPlayer, Object: id})
			labels = append(labels, string(g.Object(id).Card))
		}
	case TargetDiceRoll:
		for _, it := range g.Stack {
			if it.Kind == StackRoll {
				opts = append(opts, Chosen{Kind: TargetDiceRoll, Player: it.Controller, StackID: it.ID})
				labels = append(labels, "roll of "+strconv.Itoa(it.Roll))
			}
		}
	case TargetStackAbility:
		for _, it := range g.Stack {
			if g.cancellable(it) {
				opts = append(opts, Chosen{Kind: TargetStackAbility, Player: it.Controller, StackID: it.ID})
				label := string(it.Card)
				if it.Label != "" {
					label += ": " + it.Label
				}
				labels = append(labels, label)
			}
		}
	case TargetStackTrigger:
		for _, it := range g.Stack {
			if it.Kind == StackTrigger {
				opts = append(opts, Chosen{Kind: TargetStackTrigger, Player: it.Controller, StackID: it.ID})
				label := string(it.Card)
				if it.Label != "" {
					label += ": " + it.Label
				}
				labels = append(labels, label)
			}
		}
	case TargetYourHandCard:
		for _, id := range g.Players[p].Hand {
			opts = append(opts, Chosen{Kind: TargetYourHandCard, Player: p, Object: id})
			labels = append(labels, string(g.Object(id).Card))
		}
	case TargetItemOrSoul:
		for _, pl := range g.Players {
			for _, id := range pl.InPlay {
				if r := g.Object(id).Role; r == RoleItem || r == RoleSoul {
					opts = append(opts, Chosen{Kind: TargetItemOrSoul, Player: pl.ID, Object: id})
					labels = append(labels, string(g.Object(id).Card))
				}
			}
		}
	case TargetCurse, TargetYourItem:
		role := RoleCurse
		if t.Kind == TargetYourItem {
			role = RoleItem
		}
		for _, pl := range g.Players {
			if t.Kind == TargetYourItem && pl.ID != p {
				continue
			}
			for _, id := range pl.InPlay {
				if g.Object(id).Role == role {
					opts = append(opts, Chosen{Kind: t.Kind, Player: pl.ID, Object: id})
					labels = append(labels, string(g.Object(id).Card))
				}
			}
		}
	}
	if t.Where != nil {
		keptOpts, keptLabels := opts[:0], labels[:0]
		for i, o := range opts {
			if t.Where(g, src, o) {
				keptOpts, keptLabels = append(keptOpts, o), append(keptLabels, labels[i])
			}
		}
		opts, labels = keptOpts, keptLabels
	}
	return opts, labels
}

// cancellable: an item's ↷ or $ ability, or a loot card being played.
func (g *Game) cancellable(it StackItem) bool {
	switch it.Kind { //nolint:exhaustive // only these two can be cancelled this way
	case StackLoot:
		return true
	case StackAbility:
		return g.Object(it.Source).Role == RoleItem
	}
	return false
}

// stillValid checks a target on resolution (R-ABIL-06).
func (g *Game) stillValid(c Chosen) bool {
	switch c.Kind {
	case TargetPlayer, TargetOtherPlayer:
		return true // dead players stay valid targets (R-07 #4)
	case TargetMonster, TargetMonsterOrPlayer:
		o := g.Object(c.Object)
		return o.Zone.Kind == ZoneInPlay && o.Role == RoleMonster
	case TargetItem:
		o := g.Object(c.Object)
		return o.Zone.Kind == ZoneInPlay && o.Role == RoleItem
	case TargetDiceRoll, TargetStackAbility, TargetStackTrigger:
		for _, it := range g.Stack {
			if it.ID == c.StackID {
				return true
			}
		}
		return false
	case TargetCurse, TargetYourItem, TargetItemOrSoul:
		o := g.Object(c.Object)
		return o.Zone.Kind == ZoneInPlay && o.Controller == c.Player
	case TargetYourHandCard:
		return true // spent by the cost already
	}
	return false
}

// resolveAbility runs an ability's effects, or fizzles it when a target
// is gone (R-ABIL-06). A roll ability's result runs its table instead.
func (g *Game) resolveAbility(it StackItem) {
	effects := g.effectsOf(it.Source, it.Ability, it.Mode, it.RollResult)
	defer g.afterEventAbility(it.Source)
	// A "would roll N" trigger needs the roll to still be N (R-DICE-07).
	if it.Kind == StackTrigger && it.Checked > 0 && !g.rollStillIs(it.EventStack, it.Checked-1) {
		g.emit(Event{Kind: EvAbilityFizzle, Player: it.Controller, Object: it.Source, Card: it.Card})
		return
	}
	for _, t := range it.Targets {
		// A roll's result was checked when the roll ability resolved.
		if t.Spent || it.RollResult > 0 {
			continue
		}
		if !g.stillValid(t) {
			g.emit(Event{Kind: EvAbilityFizzle, Player: it.Controller, Object: it.Source, Card: it.Card})
			return
		}
	}
	for i, e := range effects {
		e.apply(&Ctx{
			G: g, Controller: it.Controller, Source: it.Source, Targets: it.Targets,
			EventPlayer: it.EventPlayer, EventAmount: it.EventAmount, EventStack: it.EventStack, EventObject: it.EventObject, EventSource: it.EventSource, EventBy: it.EventBy,
			ref: it.Ability, mode: it.Mode, roll: it.RollResult, effect: i,
		})
	}
}

// lootPlaysFor is how many loot plays p has right now.
func (g *Game) lootPlaysFor(p PlayerID) int {
	n := g.Players[p].ExtraLootPlays
	if p == g.Turn.Active {
		n += g.Turn.LootPlays + max(g.bonus(StatLootPlays, p, 0)-g.Turn.BonusLootUsed, 0)
	}
	return n
}

func (g *Game) useLootPlay(p PlayerID) {
	if p == g.Turn.Active && g.Turn.LootPlays > 0 {
		g.Turn.LootPlays--
		return
	}
	if p == g.Turn.Active && g.bonus(StatLootPlays, p, 0) > g.Turn.BonusLootUsed {
		g.Turn.BonusLootUsed++
		return
	}
	g.Players[p].ExtraLootPlays--
}

// AbilitiesOf lists the abilities an object can use, in the order an
// Activate intent's Choice counts them: its card's abilities, or, for a
// card that borrows them (Placebo), every ↷ ability of the other
// non-eternal items in play: players' items, then shop items, which count
// as in play for targeting (S-YUG).
func (g *Game) AbilitiesOf(id ObjectID) []AbilityRef {
	if !g.def(id).CopiesTapAbilities {
		var refs []AbilityRef
		for _, face := range g.facesOf(id) {
			if d, ok := g.cards.find(face); ok {
				for i := range d.Abilities {
					refs = append(refs, AbilityRef{Card: face, Index: i})
				}
			}
		}
		return refs
	}
	var items []ObjectID
	for _, pl := range g.Players {
		items = append(items, pl.InPlay...)
	}
	for _, slot := range g.Shop {
		if top, ok := slot.TopOf(); ok {
			items = append(items, top)
		}
	}
	var refs []AbilityRef
	for _, other := range items {
		if other == id || g.Object(other).Role != RoleItem || g.Eternal(other) {
			continue
		}
		for i, a := range g.def(other).Abilities {
			if a.Kind == Activated && hasTap(a) {
				refs = append(refs, AbilityRef{Card: g.CardOf(other), Index: i})
			}
		}
	}
	return refs
}

func hasTap(a Ability) bool {
	for _, c := range a.Costs {
		if _, ok := c.(tapCost); ok {
			return true
		}
	}
	return false
}

// afterEventAbility sends an event to discard once its ability is done,
// after what the ability queued (R-CARD-15).
func (g *Game) afterEventAbility(src ObjectID) {
	if o := g.Object(src); o.Role == RoleEvent && o.Zone.Kind == ZoneInPlay {
		if g.Object(src).Zone.Slot == MonsterSlot && !g.def(src).Curse && !g.waitingRoll(src) {
			g.enqueue(Action{Kind: ActFinishEvent, Player: NoPlayer, Object: src})
		}
	}
}

// waitingRoll: a roll or roll result of this source is still on the stack.
func (g *Game) waitingRoll(src ObjectID) bool {
	for _, it := range g.Stack {
		if it.Source == src && (it.Kind == StackRoll || it.RollFor.Card != "") {
			return true
		}
	}
	return false
}

// choosablePlayers keeps the living players: a dead player can't be
// chosen while a living one can (R-ABIL-35). Untargeted effects ("each
// player", "the player to your left") still reach the dead (R-MECH-16).
func (g *Game) choosablePlayers(ids []PlayerID) []PlayerID {
	var living []PlayerID
	for _, id := range ids {
		if !g.Players[id].Dead {
			living = append(living, id)
		}
	}
	if len(living) == 0 {
		return ids
	}
	return living
}

// tapAbilitiesInPlay lists every ↷ ability of the items in play, shop
// items included, for a card that copies one ("? Card"; R-CARD-34).
func (g *Game) tapAbilitiesInPlay(self ObjectID) []AbilityRef {
	var items []ObjectID
	for _, pl := range g.Players {
		items = append(items, pl.InPlay...)
	}
	for _, slot := range g.Shop {
		if top, ok := slot.TopOf(); ok {
			items = append(items, top)
		}
	}
	var refs []AbilityRef
	for _, id := range items {
		if id == self || g.Object(id).Role != RoleItem {
			continue
		}
		for i, a := range g.def(id).Abilities {
			if a.Kind == Activated && hasTap(a) {
				refs = append(refs, AbilityRef{Card: g.CardOf(id), Index: i})
			}
		}
	}
	return refs
}
