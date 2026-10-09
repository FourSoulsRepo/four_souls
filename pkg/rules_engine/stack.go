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
	Target     Target    `json:"target"`
	// Attack marks attack rolls and combat damage; they leave the stack
	// when the attack ends (R-ATK-15).
	Attack bool `json:"attack,omitempty"`
	// Ability is the ability of loot, activated and triggered items.
	Ability AbilityRef `json:"ability"`
	Targets []Chosen   `json:"targets,omitempty"`
	// RollFor is the roll ability waiting for this roll (R-ABIL-24).
	RollFor AbilityRef `json:"roll_for"`
	// RollResult is set on the trigger that reads a roll's result.
	RollResult int `json:"roll_result,omitempty"`
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
	start := it.Controller
	if start == NoPlayer {
		start = g.Turn.Active // the game's items: active player first (R-PRIO-02)
	}
	g.givePriority(start)
	return it.ID
}

// resolveTop resolves the top item (R-STACK-02, R-STACK-03). Priority then
// passes again, starting with the active player (R-STACK-04, R-PRIO-02).
func (g *Game) resolveTop() {
	n := len(g.Stack)
	it := g.Stack[n-1]
	g.Stack = g.Stack[:n-1]
	g.emit(Event{Kind: EvStackResolved, Player: it.Controller, Object: it.Source, Card: it.Card, Amount: it.Roll, Text: it.Label})
	// Priority after a resolution starts with the active player
	// (R-STACK-04); items pushed by the resolution set it again.
	g.givePriority(g.Turn.Active)
	switch it.Kind {
	case StackLoot:
		// The loot ability happens, then the loot goes to the loot
		// discard (R-CARD-07).
		if it.Ability.Card != "" {
			g.resolveAbility(it)
		}
		g.discard(it.Source, LootDeck)
	case StackRoll:
		// The result is final once it resolves (R-DICE-06).
		g.emit(Event{Kind: EvRollResolved, Player: it.Controller, Amount: it.Roll})
		if it.Attack {
			g.resolveAttackRoll(it)
		}
		if it.RollFor.Card != "" {
			// The roll ability's result trigger goes on the stack (R-ABIL-24).
			g.push(StackItem{
				Kind: StackTrigger, Controller: it.Controller, Source: it.Source, Card: it.RollFor.Card,
				Ability: it.RollFor, RollResult: it.Roll, Label: "roll result",
			})
		}
	case StackDamage:
		g.resolveDamage(it)
	case StackDeath:
		g.resolveDeath(it)
	case StackAbility, StackTrigger:
		g.resolveAbility(it)
	}
	if len(g.Queue) > 0 && g.Waiting.Kind == PromptPriority {
		// Queued steps (e.g. death rewards) happen before anyone acts.
		g.givePriority(g.Priority.Holder)
	}
}

// roll rolls a D6 for p and puts the roll on the stack (R-DICE-02).
func (g *Game) roll(p PlayerID, label string) int {
	r := g.d6()
	g.emit(Event{Kind: EvDiceRolled, Player: p, Amount: r, Text: label})
	return g.push(StackItem{Kind: StackRoll, Controller: p, Roll: r, Label: label})
}

// playLoot moves a loot card from hand to the stack (R-MECH-45).
func (g *Game) playLoot(p PlayerID, id ObjectID) {
	g.Players[p].Hand = remove(g.Players[p].Hand, id)
	g.useLootPlay(p)
	nid := g.move(id, Zone{Kind: ZoneStack}, p)
	card := g.Object(nid).Card
	g.emit(Event{Kind: EvLootPlayed, Player: p, Object: nid, Card: card})
	g.push(StackItem{Kind: StackLoot, Controller: p, Source: nid, Card: card})
}

// givePriority opens a window for p, or waits until the action queue is
// empty: queued steps (such as death steps) happen before anyone acts
// (R-DEATH-10).
func (g *Game) givePriority(p PlayerID) {
	if len(g.Queue) > 0 {
		g.Priority = Priority{Deferred: true, Holder: p}
		g.Waiting = Prompt{}
		return
	}
	g.openWindow(p)
}
