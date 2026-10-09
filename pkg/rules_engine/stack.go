package rulesengine

// StackKind is what a stack item is (R-STACK-01).
type StackKind int

// The stack item kinds.
const (
	StackLoot    StackKind = iota // a played loot card (R-CARD-07)
	StackAbility                  // an activated ability (R-ABIL-08)
	StackTrigger                  // a triggered ability (R-ABIL-14)
	StackRoll                     // a dice roll (R-DICE-02)
	StackDamage                   // damage aimed at a target (R-MECH-15)
	StackDeath                    // a pending death (R-DEATH-01)
)

// StackItem is one thing waiting on the stack.
type StackItem struct {
	ID         int       `json:"id"`
	Kind       StackKind `json:"kind"`
	Controller PlayerID  `json:"controller"`
	Source     ObjectID  `json:"source"`         // the card or object it comes from
	Card       CardRef   `json:"card,omitempty"` // for display
	Roll       int       `json:"roll,omitempty"` // current result of a roll
	Amount     int       `json:"amount,omitempty"`
	Label      string    `json:"label,omitempty"`
}

// Events about the stack.
const (
	EvStackAdded    EventKind = "stack_added"
	EvStackResolved EventKind = "stack_resolved"
	EvLootPlayed    EventKind = "loot_played"
)

// push puts an item on top of the stack; priority passes starting with
// its controller (R-PRIO-03).
func (g *Game) push(it StackItem) int {
	g.StackSeq++
	it.ID = g.StackSeq
	g.Stack = append(g.Stack, it)
	g.emit(Event{Kind: EvStackAdded, Player: it.Controller, Object: it.Source, Card: it.Card, Amount: it.Roll, Text: it.Label})
	g.openWindow(it.Controller)
	return it.ID
}

// resolveTop resolves the top item (R-STACK-02, R-STACK-03). Priority then
// passes again, starting with the active player (R-STACK-04, R-PRIO-02).
func (g *Game) resolveTop() {
	n := len(g.Stack)
	it := g.Stack[n-1]
	g.Stack = g.Stack[:n-1]
	switch it.Kind {
	case StackLoot:
		// Loot abilities come with effect blocks (step 4.7); then the
		// loot goes to the loot discard (R-CARD-07).
		g.discard(it.Source, LootDeck)
	case StackRoll:
		// The result is final once it resolves (R-DICE-06).
	case StackAbility, StackTrigger, StackDamage, StackDeath:
		// Filled in by steps 4.5 to 4.7.
	}
	g.emit(Event{Kind: EvStackResolved, Player: it.Controller, Object: it.Source, Card: it.Card, Amount: it.Roll, Text: it.Label})
	if !g.Over {
		g.openWindow(g.Turn.Active)
	}
}

// roll rolls a D6 for p and puts the roll on the stack (R-DICE-02).
func (g *Game) roll(p PlayerID, label string) int {
	r := g.RNG.D6()
	g.emit(Event{Kind: EvDiceRolled, Player: p, Amount: r, Text: label})
	return g.push(StackItem{Kind: StackRoll, Controller: p, Roll: r, Label: label})
}

// playLoot moves a loot card from hand to the stack (R-MECH-45).
func (g *Game) playLoot(p PlayerID, id ObjectID) {
	g.Players[p].Hand = remove(g.Players[p].Hand, id)
	g.Turn.LootPlays--
	nid := g.move(id, Zone{Kind: ZoneStack}, p)
	card := g.Object(nid).Card
	g.emit(Event{Kind: EvLootPlayed, Player: p, Object: nid, Card: card})
	g.push(StackItem{Kind: StackLoot, Controller: p, Source: nid, Card: card})
}
