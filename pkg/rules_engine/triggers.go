package rulesengine

// AbilityRef names one ability of one card.
type AbilityRef struct {
	Card  CardRef `json:"card"`
	Index int     `json:"index"`
}

// Trigger says when a triggered ability triggers (R-ABIL-13).
type Trigger struct {
	On    EventKind
	Match func(g *Game, self ObjectID, e Event) bool
}

// Events used by triggers.
const (
	EvStartOfTurn   EventKind = "start_of_turn" // start-of-turn triggers (R-TURN-03)
	EvEndOfTurn     EventKind = "end_of_turn"   // end-of-turn triggers (R-TURN-10)
	EvRollResolved  EventKind = "roll_resolved" // a roll became final (R-DICE-06)
	EvAbilityFizzle EventKind = "fizzled"       // R-ABIL-06
)

// AtStartOfYourTurn triggers at the start of the controller's turn.
func AtStartOfYourTurn() Trigger {
	return Trigger{On: EvStartOfTurn, Match: func(g *Game, self ObjectID, e Event) bool {
		return e.Player == g.Object(self).Controller
	}}
}

// AtEndOfYourTurn triggers at the end of the controller's turn.
func AtEndOfYourTurn() Trigger {
	return Trigger{On: EvEndOfTurn, Match: func(g *Game, self ObjectID, e Event) bool {
		return e.Player == g.Object(self).Controller
	}}
}

// WhenAMonsterDies triggers when any monster dies (R-DEATH-05).
func WhenAMonsterDies() Trigger {
	return Trigger{On: EvDied, Match: func(_ *Game, _ ObjectID, e Event) bool {
		return e.Player == NoPlayer
	}}
}

// OnRollOf triggers when any roll resolves as n (circled number, R-ABIL-16).
func OnRollOf(n int) Trigger {
	return Trigger{On: EvRollResolved, Match: func(_ *Game, _ ObjectID, e Event) bool {
		return e.Amount == n
	}}
}

// WhenThisIsDestroyed triggers when this object is destroyed.
func WhenThisIsDestroyed() Trigger {
	return Trigger{On: EvDestroyed, Match: func(_ *Game, self ObjectID, e Event) bool {
		return e.Prev == self
	}}
}

// PendingTrigger waits to go on the stack (R-ABIL-14).
type PendingTrigger struct {
	Ability    AbilityRef `json:"ability"`
	Source     ObjectID   `json:"source"`
	Controller PlayerID   `json:"controller"`
}

// collectTriggers finds triggered abilities that match an event. Objects
// in play are checked, plus the object the event moved away (e.Prev), so
// "when this is destroyed" still sees itself.
func (g *Game) collectTriggers(e Event) {
	check := func(id ObjectID) {
		o := g.Object(id)
		def, ok := g.cards.find(o.Card)
		if !ok {
			return
		}
		for i, a := range def.Abilities {
			if a.Kind == Triggered && a.Trigger.On == e.Kind && a.Trigger.Match != nil && a.Trigger.Match(g, id, e) {
				g.PendingTriggers = append(g.PendingTriggers, PendingTrigger{
					Ability: AbilityRef{Card: o.Card, Index: i}, Source: id, Controller: o.Controller,
				})
			}
		}
	}
	for i := range g.Objects {
		if g.Objects[i].Zone.Kind == ZoneInPlay {
			check(g.Objects[i].ID)
		}
	}
	if e.Prev != 0 && g.Object(e.Prev).Zone.Kind != ZoneInPlay {
		check(e.Prev)
	}
}

// placeNextTrigger puts waiting triggers on the stack in rule order: the
// game's first (ordered by the active player), then each player's in turn
// order from the active player; a player with several orders them
// (R-ABIL-15). The engine asks whenever there is an order to choose.
func (g *Game) placeNextTrigger() {
	group := g.nextTriggerGroup()
	if len(group) == 1 {
		g.pushTrigger(group[0])
		return
	}
	chooser := g.PendingTriggers[group[0]].Controller
	if chooser == NoPlayer {
		chooser = g.Turn.Active
	}
	labels := make([]string, len(group))
	for i, gi := range group {
		labels[i] = g.abilityText(g.PendingTriggers[gi].Ability)
	}
	g.ask(Choice{Purpose: ChooseTriggerOrder, Player: chooser, Rule: "R-ABIL-15", Indexes: group}, labels)
}

// nextTriggerGroup returns the indexes of the triggers that go next.
func (g *Game) nextTriggerGroup() []int {
	order := []PlayerID{NoPlayer}
	p := g.Turn.Active
	for range g.Players {
		order = append(order, p)
		p = g.next(p)
	}
	for _, c := range order {
		var group []int
		for i, t := range g.PendingTriggers {
			if t.Controller == c {
				group = append(group, i)
			}
		}
		if len(group) > 0 {
			return group
		}
	}
	return nil
}

func (g *Game) pushTrigger(i int) {
	t := g.PendingTriggers[i]
	g.PendingTriggers = append(g.PendingTriggers[:i:i], g.PendingTriggers[i+1:]...)
	g.push(StackItem{
		Kind: StackTrigger, Controller: t.Controller, Source: t.Source, Card: t.Ability.Card,
		Ability: t.Ability, Label: g.abilityText(t.Ability),
	})
}

func (g *Game) ability(ref AbilityRef) Ability {
	d, ok := g.cards.find(ref.Card)
	if !ok || ref.Index >= len(d.Abilities) {
		panic("rulesengine: unknown ability " + string(ref.Card))
	}
	return d.Abilities[ref.Index]
}

func (g *Game) abilityText(ref AbilityRef) string {
	if t := g.ability(ref).Text; t != "" {
		return string(ref.Card) + ": " + t
	}
	return string(ref.Card)
}
