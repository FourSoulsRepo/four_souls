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
	StackPenalty                  // a death penalty waiting for "when a player dies" triggers (R-DEATH-13)
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
	Mode    int        `json:"mode,omitempty"` // the chosen "choose one-" option
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
	// "Would take damage" and "would die" triggers look at these.
	switch it.Kind { //nolint:exhaustive // only damage and death have "would" triggers
	case StackDamage:
		g.emit(Event{Kind: EvDamagePending, Player: it.Target.player(), Object: it.Target.Object, Amount: it.Amount})
	case StackDeath:
		g.emit(Event{Kind: EvDeathPending, Player: it.Target.player(), Object: it.Target.Object})
	}
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
		switch {
		case g.Object(it.Source).Zone.Kind != ZoneStack:
			// The effect moved the card already (The Sun).
		case g.def(it.Source).Trinket:
			g.trinketToPlay(it.Controller, it.Source) // R-ABIL-19
		default:
			g.discard(it.Source, LootDeck)
		}
	case StackRoll:
		// Continuous roll changes apply, then the result is final (R-DICE-06).
		if it.Controller != NoPlayer {
			it.Roll = min(max(it.Roll+g.bonus(StatRoll, it.Controller, 0), 1), 6)
		}
		g.emit(Event{Kind: EvRollResolved, Player: it.Controller, Amount: it.Roll})
		if it.Attack {
			g.resolveAttackRoll(it)
		}
		if it.RollFor.Card != "" {
			// The roll ability's result trigger goes on the stack (R-ABIL-24).
			g.push(StackItem{
				Kind: StackTrigger, Controller: it.Controller, Source: it.Source, Card: it.RollFor.Card,
				Ability: it.RollFor, Mode: it.Mode, Targets: it.Targets, RollResult: it.Roll, Label: "roll result",
			})
		}
	case StackDamage:
		g.resolveDamage(it)
	case StackDeath:
		g.resolveDeath(it)
	case StackAbility, StackTrigger:
		g.resolveAbility(it)
	case StackPenalty:
		g.payPenalty(it.Target.Player)
	}
	if (len(g.Queue) > 0 || len(g.PendingTriggers) > 0) && g.Waiting.Kind == PromptPriority {
		// Queued steps (e.g. death rewards) happen and new triggers go on
		// the stack before anyone acts (R-ABIL-14).
		g.openWindow(g.Priority.Holder)
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

// SetRoll changes a dice roll on the stack to n, between 1 and 6
// (R-DICE-04): "change the result of a dice roll".
func (g *Game) SetRoll(stackID, n int) {
	for i := range g.Stack {
		if it := &g.Stack[i]; it.ID == stackID && it.Kind == StackRoll {
			it.Roll = min(max(n, 1), 6)
			g.emit(Event{Kind: EvRollChanged, Player: it.Controller, Amount: it.Roll})
		}
	}
}

// player is the target player, or NoPlayer for an object.
func (t Target) player() PlayerID {
	if t.IsPlayer {
		return t.Player
	}
	return NoPlayer
}

// trinketToPlay puts a resolved trinket into play as an item of p.
func (g *Game) trinketToPlay(p PlayerID, id ObjectID) {
	nid := g.move(id, Zone{Kind: ZoneInPlay}, p)
	o := g.Object(nid)
	o.Role, o.Charged = RoleItem, true
	g.Players[p].InPlay = append(g.Players[p].InPlay, nid)
	g.emit(Event{Kind: EvEnteredPlay, Player: p, Object: nid, Card: o.Card})
}

// CancelStackItem removes an item from the stack without resolving it;
// a cancelled loot card goes to the loot discard (R-MECH-32).
func (g *Game) CancelStackItem(id int) {
	for i, it := range g.Stack {
		if it.ID != id {
			continue
		}
		g.Stack = append(g.Stack[:i:i], g.Stack[i+1:]...)
		g.emit(Event{Kind: EvCancelled, Player: it.Controller, Object: it.Source, Card: it.Card, Text: it.Label})
		if it.Kind == StackLoot && g.Object(it.Source).Zone.Kind == ZoneStack {
			g.discard(it.Source, LootDeck)
		}
		return
	}
}

// EndTurnNow cancels everything that has not resolved and ends the
// active player's turn: "End the turn. Cancel everything that hasn't
// resolved." The turn goes to its end phase.
func (g *Game) EndTurnNow() {
	g.endAttack()
	g.Purchase = PurchaseState{}
	for len(g.Stack) > 0 {
		g.CancelStackItem(g.Stack[len(g.Stack)-1].ID)
	}
	if g.Turn.Step < StepAction {
		// From the start phase too: the action phase closes at once.
		g.Turn.Step, g.Turn.Entered = StepAction, true
	}
	g.Turn.EndDeclared = true
	g.emit(Event{Kind: EvTurnEndedEarly, Player: g.Turn.Active})
}
