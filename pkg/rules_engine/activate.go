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
	o := g.Object(src)
	if o.Zone.Kind != ZoneInPlay || o.Controller != p || (o.Role != RoleItem && o.Role != RoleCharacter) {
		return refuse("R-ABIL-09", "only the controller uses an item's or character's abilities")
	}
	d := g.def(src)
	if i < 0 || i >= len(d.Abilities) || d.Abilities[i].Kind != Activated {
		return refuse("R-ABIL-08", "%s has no such activated ability", o.Card)
	}
	a := d.Abilities[i]
	for _, c := range a.Costs {
		if !c.canPay(g, p, src) {
			return refuse("R-ABIL-07", "cannot pay the cost (%s)", c.label())
		}
	}
	return g.targetsAvailable(a, p)
}

// targetsAvailable checks every target spec has something to choose; an
// ability with modes needs one mode whose targets are there.
func (g *Game) targetsAvailable(a Ability, p PlayerID) error {
	if len(a.Modes) > 0 {
		if len(g.availableModes(a, p)) == 0 {
			return refuse("R-ABIL-06", "no option has a valid target")
		}
		return nil
	}
	if !g.hasTargets(a.Targets, p) {
		return refuse("R-ABIL-06", "no valid target")
	}
	return nil
}

func (g *Game) hasTargets(specs []TargetSpec, p PlayerID) bool {
	for _, t := range specs {
		if opts, _ := g.targetOptions(t, p); len(opts) == 0 {
			return false
		}
	}
	return true
}

// availableModes lists the modes of a whose targets can be chosen.
func (g *Game) availableModes(a Ability, p PlayerID) []int {
	var out []int
	for i, m := range a.Modes {
		if g.hasTargets(m.Targets, p) {
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
	modes := g.availableModes(ab, a.Player)
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
		opts, labels := g.targetOptions(specs[len(a.Chosen)], a.Player)
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
		g.emit(Event{Kind: EvLootPlayed, Player: a.Player, Object: nid, Card: card})
		g.push(StackItem{Kind: StackLoot, Controller: a.Player, Source: nid, Card: card, Ability: a.Ability, Mode: a.Mode, Targets: a.Chosen})
		return
	}
	for _, c := range ab.Costs {
		c.pay(g, a.Player, a.Source)
	}
	g.emit(Event{Kind: EvActivated, Player: a.Player, Object: a.Source, Card: a.Ability.Card, Text: ab.Text})
	g.push(StackItem{
		Kind: StackAbility, Controller: a.Player, Source: a.Source, Card: a.Ability.Card,
		Ability: a.Ability, Mode: a.Mode, Targets: a.Chosen, Label: ab.Text,
	})
}

// targetOptions lists what may be chosen for a spec, with labels.
func (g *Game) targetOptions(t TargetSpec, p PlayerID) ([]Chosen, []string) {
	var opts []Chosen
	var labels []string
	others := t.Kind == TargetOtherPlayer
	addPlayers := func() {
		for _, pl := range g.Players {
			if !pl.Dead && (!others || pl.ID != p) {
				opts = append(opts, Chosen{Kind: TargetPlayer, Player: pl.ID})
				labels = append(labels, "player "+strconv.Itoa(int(pl.ID)+1)+" ("+string(g.Object(pl.Character).Card)+")")
			}
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
			if t.Where(g, p, o) {
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
		return !g.Players[c.Player].Dead
	case TargetMonster, TargetMonsterOrPlayer:
		o := g.Object(c.Object)
		return o.Zone.Kind == ZoneInPlay && o.Role == RoleMonster
	case TargetItem:
		o := g.Object(c.Object)
		return o.Zone.Kind == ZoneInPlay && o.Role == RoleItem
	case TargetDiceRoll, TargetStackAbility:
		for _, it := range g.Stack {
			if it.ID == c.StackID {
				return true
			}
		}
		return false
	case TargetCurse, TargetYourItem:
		o := g.Object(c.Object)
		return o.Zone.Kind == ZoneInPlay && o.Controller == c.Player
	}
	return false
}

// resolveAbility runs an ability's effects, or fizzles it when a target
// is gone (R-ABIL-06). A roll ability's result runs its table instead.
func (g *Game) resolveAbility(it StackItem) {
	effects := g.effectsOf(it.Ability, it.Mode, it.RollResult)
	for _, t := range it.Targets {
		if !g.stillValid(t) {
			g.emit(Event{Kind: EvAbilityFizzle, Player: it.Controller, Object: it.Source, Card: it.Card})
			return
		}
	}
	for i, e := range effects {
		e.apply(&Ctx{
			G: g, Controller: it.Controller, Source: it.Source, Targets: it.Targets,
			ref: it.Ability, mode: it.Mode, roll: it.RollResult, effect: i,
		})
	}
}

// lootPlaysFor is how many loot plays p has right now.
func (g *Game) lootPlaysFor(p PlayerID) int {
	n := g.Players[p].ExtraLootPlays
	if p == g.Turn.Active {
		n += g.Turn.LootPlays
	}
	return n
}

func (g *Game) useLootPlay(p PlayerID) {
	if p == g.Turn.Active && g.Turn.LootPlays > 0 {
		g.Turn.LootPlays--
		return
	}
	g.Players[p].ExtraLootPlays--
}
