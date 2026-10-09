package rulesengine

import "strconv"

// Activation is an ability (or loot card) being put on the stack while its
// targets are chosen (R-ABIL-04).
type Activation struct {
	Player  PlayerID   `json:"player"`
	Source  ObjectID   `json:"source"`
	Ability AbilityRef `json:"ability"`
	Loot    bool       `json:"loot,omitempty"` // a loot card from hand
	Chosen  []Chosen   `json:"chosen,omitempty"`
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

// targetsAvailable checks every target spec has something to choose.
func (g *Game) targetsAvailable(a Ability, p PlayerID) error {
	for _, t := range a.Targets {
		if opts, _ := g.targetOptions(t, p); len(opts) == 0 {
			return refuse("R-ABIL-06", "no valid target")
		}
	}
	return nil
}

// startActivation asks for targets one by one, then pays and pushes.
func (g *Game) startActivation(a Activation) {
	a.Window = g.Priority.Holder
	g.Activating = &a
	g.nextTarget()
}

func (g *Game) nextTarget() {
	a := g.Activating
	ab := g.ability(a.Ability)
	if len(a.Chosen) < len(ab.Targets) {
		opts, labels := g.targetOptions(ab.Targets[len(a.Chosen)], a.Player)
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
		g.Activating = nil
		g.emit(Event{Kind: EvCancelled, Player: a.Player, Object: a.Source})
		g.openWindow(a.Window)
		return
	}
	a.Chosen = append(a.Chosen, c.Targets[i])
	g.nextTarget()
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
		g.push(StackItem{Kind: StackLoot, Controller: a.Player, Source: nid, Card: card, Ability: a.Ability, Targets: a.Chosen})
		return
	}
	for _, c := range ab.Costs {
		c.pay(g, a.Player, a.Source)
	}
	g.emit(Event{Kind: EvActivated, Player: a.Player, Object: a.Source, Card: a.Ability.Card, Text: ab.Text})
	g.push(StackItem{
		Kind: StackAbility, Controller: a.Player, Source: a.Source, Card: a.Ability.Card,
		Ability: a.Ability, Targets: a.Chosen, Label: ab.Text,
	})
}

// targetOptions lists what may be chosen for a spec, with labels.
func (g *Game) targetOptions(t TargetSpec, _ PlayerID) ([]Chosen, []string) {
	var opts []Chosen
	var labels []string
	addPlayers := func() {
		for _, pl := range g.Players {
			if !pl.Dead {
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
	case TargetPlayer:
		addPlayers()
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
	}
	return opts, labels
}

// stillValid checks a target on resolution (R-ABIL-06).
func (g *Game) stillValid(c Chosen) bool {
	switch c.Kind {
	case TargetPlayer:
		return !g.Players[c.Player].Dead
	case TargetMonster, TargetMonsterOrPlayer:
		o := g.Object(c.Object)
		return o.Zone.Kind == ZoneInPlay && o.Role == RoleMonster
	case TargetItem:
		o := g.Object(c.Object)
		return o.Zone.Kind == ZoneInPlay && o.Role == RoleItem
	case TargetDiceRoll:
		for _, it := range g.Stack {
			if it.ID == c.StackID {
				return true
			}
		}
		return false
	}
	return false
}

// resolveAbility runs an ability's effects, or fizzles it when a target
// is gone (R-ABIL-06). A roll ability's result runs its table instead.
func (g *Game) resolveAbility(it StackItem) {
	ab := g.ability(it.Ability)
	effects := ab.Effects
	if it.RollResult > 0 {
		effects = nil
		for _, e := range ab.Effects {
			if r, ok := e.(rollEffect); ok {
				effects = r.table[it.RollResult-1]
			}
		}
	}
	for _, t := range it.Targets {
		if !g.stillValid(t) {
			g.emit(Event{Kind: EvAbilityFizzle, Player: it.Controller, Object: it.Source, Card: it.Card})
			return
		}
	}
	c := &Ctx{G: g, Controller: it.Controller, Source: it.Source, Targets: it.Targets}
	g.resolving = it.Ability
	for _, e := range effects {
		e.apply(c)
	}
	g.resolving = AbilityRef{}
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
