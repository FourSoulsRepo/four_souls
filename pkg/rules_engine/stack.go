package rulesengine

import "slices"

// StackKind is what a stack item is (R-STACK-01).
type StackKind int

// The stack item kinds.
const (
	StackLoot      StackKind = iota // a played loot card (R-CARD-07)
	StackAbility                    // an activated ability (R-ABIL-08)
	StackTrigger                    // a triggered ability (R-ABIL-14)
	StackRoll                       // a dice roll (R-DICE-02)
	StackDamage                     // damage aimed at a target (R-MECH-15)
	StackDeath                      // a pending death (R-DEATH-01)
	StackPenalty                    // a death penalty waiting for "when a player dies" triggers (R-DEATH-13)
	StackDeathStep                  // a monster death step waiting for triggers (R-DEATH-05, R-DEATH-07)
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
	// EventPlayer and EventAmount: the triggering event, for triggers.
	EventPlayer PlayerID `json:"event_player,omitempty"`
	EventAmount int      `json:"event_amount,omitempty"`
	EventStack  int      `json:"event_stack,omitempty"`
	EventObject ObjectID `json:"event_object,omitempty"`
	// EventSource is what caused the event, e.g. the monster that dealt damage.
	EventSource ObjectID `json:"event_source,omitempty"`
	// EventBy is the player who dealt the damage, plus 1 (Event.By).
	EventBy int `json:"event_by,omitempty"`
	// Reward is a roll reward's kind + 1: the result is how much is gained.
	Reward int `json:"reward,omitempty"`
	// Checked is the roll value (+1) that "would roll" triggers last saw
	// (R-DICE-05); the roll only resolves once nothing changes it. On a
	// "would roll" trigger it is the value (+1) the trigger saw: it
	// fizzles if the roll no longer has it (R-DICE-07).
	Checked int `json:"checked,omitempty"`
	// Text is the words of the ability it uses, or of its chosen mode,
	// for display. Only views set it (15.21, S-07).
	Text string `json:"text,omitempty"`
	// SourceCard and TargetCard are the cards of Source and Target, for
	// display, while anyone may see them. Only views set them (15.22).
	SourceCard CardRef `json:"source_card,omitempty"`
	TargetCard CardRef `json:"target_card,omitempty"`
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
	if it.Kind == StackDeath && it.Target.IsPlayer && g.useDeathWard(it.Target.Player) {
		return 0
	}
	g.StackSeq++
	it.ID = g.StackSeq
	if it.Kind == StackRoll && it.Controller != NoPlayer && it.Roll == 0 && len(g.ManyDice) == 0 && g.bonus(StatTwoDiceLeftChooses, NoPlayer, 0) > 0 {
		// "If a player would roll a dice, they instead roll 2 dice. The
		// player to their left chooses one of the rolls as the result."
		g.ManyDice = append(g.ManyDice, ManyDice{Chooser: g.next(it.Controller), Dice: 2})
	}
	if it.Kind == StackRoll && it.Controller != NoPlayer && len(g.ManyDice) > 0 && g.ManyDice[0].StackID == 0 {
		// "They instead roll 4 dice. You choose one of the rolls." The
		// dice are drawn with the throw.
		g.ManyDice[0].StackID = it.ID
	}
	if it.Kind == StackRoll && it.Roll == 0 {
		// The roll waits for its throw before anyone gets priority; its
		// value is drawn only then (15.59).
		g.enqueue(Action{Kind: ActThrow, Player: it.Controller, Amount: it.ID})
	}
	g.Stack = append(g.Stack, it)
	g.emit(Event{Kind: EvStackAdded, Player: it.Controller, Object: it.Source, Card: it.Card, Amount: it.Roll, Text: it.Label})
	// "Would take damage" and "would die" triggers look at these.
	switch it.Kind { //nolint:exhaustive // only damage and death have "would" triggers
	case StackDamage:
		e := Event{Kind: EvDamagePending, Player: it.Target.player(), Object: it.Target.Object, Amount: it.Amount, StackID: it.ID, Source: it.Source}
		if it.Attack {
			e.Text = "combat" // as on the damaged event
		}
		g.emit(e)
	case StackDeath:
		g.emit(Event{Kind: EvDeathPending, Player: it.Target.player(), Object: it.Target.Object, StackID: it.ID})
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
	if top := &g.Stack[n-1]; top.Kind == StackRoll && top.Checked != top.Roll+1 {
		// The roll tries to resolve: "would roll N" triggers trigger now; if
		// any did, they go on the stack first (R-DICE-05).
		top.Checked = top.Roll + 1
		before := len(g.PendingTriggers)
		g.emit(Event{Kind: EvRollWouldResolve, Player: top.Controller, Amount: top.Roll, StackID: top.ID})
		if len(g.PendingTriggers) > before {
			g.openWindow(g.Turn.Active)
			return
		}
	}
	it := g.Stack[n-1]
	g.Stack = g.Stack[:n-1]
	shown := g.shown(it)
	g.emit(Event{
		Kind: EvStackResolved, Player: it.Controller, Object: it.Source, Card: it.Card, Amount: it.Roll, Text: it.Label,
		StackID: it.ID, Item: &shown,
	})
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
			g.discardFromStack(it)
		}
	case StackRoll:
		// Continuous roll changes apply, then the result is final (R-DICE-06).
		if it.Controller != NoPlayer {
			bonus := g.bonus(StatRoll, it.Controller, 0)
			if it.Attack {
				bonus += g.bonus(StatAttackRoll, it.Controller, 0)
			}
			if it.Attack {
				bonus += g.useNextAttackRoll(it.Controller)
			}
			it.Roll = min(max(it.Roll+bonus, 1), 6)
			if (it.Roll == 2 || it.Roll == 3) && g.bonus(StatRollToOne, NoPlayer, 0) > 0 {
				it.Roll = 1 // a replacement as the roll resolves (R-DICE-12)
			}
		}
		re := Event{Kind: EvRollResolved, Player: it.Controller, Amount: it.Roll}
		if it.Attack {
			re.Text = "attack"
		}
		g.emit(re)
		if it.Controller != NoPlayer {
			pl := &g.Players[it.Controller]
			if pl.Rolls > 0 && pl.LastRoll == it.Roll {
				g.emit(Event{Kind: EvRollRepeated, Player: it.Controller, Amount: it.Roll}) // R-DICE-08
			}
			pl.Rolls, pl.LastRoll = pl.Rolls+1, it.Roll
			if it.Attack {
				if pl.LastAttack == it.Roll {
					g.emit(Event{Kind: EvAttackRollRepeated, Player: it.Controller, Amount: it.Roll}) // R-DICE-09
				}
				pl.LastAttack = it.Roll
			}
		}
		if it.Attack {
			g.resolveAttackRoll(it)
		}
		if it.Reward > 0 {
			r := Reward{Kind: RewardKind(it.Reward - 1)}
			g.enqueue(Action{Kind: r.action(), Player: it.Controller, Amount: it.Roll * max(it.Amount, 1)})
		}
		if it.RollFor.Card != "" {
			// The roll ability's result trigger goes on the stack (R-ABIL-24).
			g.push(StackItem{
				Kind: StackTrigger, Controller: it.Controller, Source: it.Source, Card: it.RollFor.Card,
				Ability: it.RollFor, Mode: it.Mode, Targets: it.Targets, RollResult: it.Roll, Label: "roll result",
				EventPlayer: it.EventPlayer, EventAmount: it.EventAmount, EventStack: it.EventStack, EventObject: it.EventObject, EventSource: it.EventSource, EventBy: it.EventBy,
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
	case StackDeathStep:
		if it.Amount == stepRewards {
			g.rewardsStep(it.Target.Object)
		} else {
			g.soulStep(it.Target.Object)
		}
	}
	if (len(g.Queue) > 0 || len(g.PendingTriggers) > 0) && g.Waiting.Kind == PromptPriority {
		// Queued steps (e.g. death rewards) happen and new triggers go on
		// the stack before anyone acts (R-ABIL-14).
		g.openWindow(g.Priority.Holder)
	}
}

// roll puts a roll of p on the stack; p throws its die first (R-DICE-02).
func (g *Game) roll(p PlayerID, label string) int {
	return g.push(StackItem{Kind: StackRoll, Controller: p, Label: label})
}

// Throw is a roll on the stack waiting for its die to be thrown (15.59).
// Reroll marks a throw that rolls the roll again (R-MECH-47).
type Throw struct {
	StackID int  `json:"stack_id"`
	Reroll  bool `json:"reroll,omitempty"`
}

// askThrow asks the roller of a queued throw to throw. A roll that left
// the stack before its throw is never thrown; the game's own rolls are
// thrown at once.
func (g *Game) askThrow(a Action) {
	it, ok := g.StackItemByID(a.Amount)
	if !ok || it.Kind != StackRoll {
		return
	}
	t := Throw{StackID: it.ID, Reroll: a.Reroll}
	if it.Controller == NoPlayer {
		g.throwDie(t)
		return
	}
	g.Throw = &t
	g.Waiting = Prompt{Kind: PromptThrow, Player: it.Controller, Text: it.Label}
}

// throw answers the throw prompt.
func (g *Game) throw() {
	t := *g.Throw
	g.Throw, g.Waiting = nil, Prompt{}
	g.throwDie(t)
}

// throwDie draws the value of a thrown roll (R-DICE-01). Many dice
// (R-DICE-11) are drawn here too; the chooser then picks the result.
func (g *Game) throwDie(t Throw) {
	i := slices.IndexFunc(g.Stack, func(it StackItem) bool { return it.ID == t.StackID })
	if i < 0 {
		return
	}
	it := &g.Stack[i]
	it.Roll = g.d6()
	text := it.Label
	if t.Reroll {
		text = "reroll"
	}
	r, p := it.Roll, it.Controller
	if len(g.ManyDice) > 0 && g.ManyDice[0].StackID == t.StackID && g.ManyDice[0].Values == nil {
		md := &g.ManyDice[0]
		md.Values = []int{r}
		for range md.Dice - 1 {
			md.Values = append(md.Values, g.d6())
		}
		g.enqueue(Action{Kind: ActChooseDie, Player: md.Chooser, Amount: t.StackID})
	}
	g.emit(Event{Kind: EvDiceRolled, Player: p, Amount: r, Text: text, StackID: t.StackID})
}

// playLoot moves a loot card from hand to the stack (R-MECH-45).
func (g *Game) playLoot(p PlayerID, id ObjectID) {
	g.Players[p].Hand = remove(g.Players[p].Hand, id)
	g.useLootPlay(p)
	nid := g.move(id, Zone{Kind: ZoneStack}, p)
	card := g.Object(nid).Card
	g.Players[p].LootPlayed++
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
	g.enterAsItem(p, nid)
	g.emit(Event{Kind: EvEnteredPlay, Player: p, Object: nid, Card: g.Object(nid).Card})
	g.emit(Event{Kind: EvGainedItem, Player: p, Object: nid, Card: g.Object(nid).Card})
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
		if o := g.Object(it.Source); o.Role == RoleEvent && !g.waitingRoll(it.Source) {
			g.enqueue(Action{Kind: ActFinishEvent, Player: NoPlayer, Object: it.Source}) // R-CARD-15
		}
		if it.Kind == StackLoot && g.Object(it.Source).Zone.Kind == ZoneStack {
			g.discardFromStack(it)
		}
		return
	}
}

// discardFromStack puts a loot card that left the stack into the loot
// discard. Prev is its object on the stack, so the client flies it from
// there (15.74).
func (g *Game) discardFromStack(it StackItem) {
	nid := g.discard(it.Source, LootDeck)
	g.emit(Event{Kind: EvDiscarded, Player: it.Controller, Object: nid, Card: it.Card, Prev: it.Source, Text: "stack"})
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

// StackItemByID returns the stack item with the given ID.
func (g *Game) StackItemByID(id int) (StackItem, bool) {
	for _, it := range g.Stack {
		if it.ID == id {
			return it, true
		}
	}
	return StackItem{}, false
}

// Reroll rolls a dice roll on the stack again (R-MECH-47): the roll's
// own player throws again before anyone gets priority (15.59).
func (g *Game) Reroll(stackID int) {
	for i := range g.Stack {
		if it := &g.Stack[i]; it.ID == stackID && it.Kind == StackRoll {
			it.Roll = 0
			g.enqueue(Action{Kind: ActThrow, Player: it.Controller, Amount: it.ID, Reroll: true})
		}
	}
}

// useDeathWard prevents p's death with a ward from this turn, if p has
// one: "The next time that player would die this turn, prevent it. If
// it's their turn, cancel everything that hasn't resolved and end it."
// (Holy Mantle, Holy Card; R-MECH-07, R-MECH-52).
func (g *Game) useDeathWard(p PlayerID) bool {
	for i, b := range g.Boosts {
		if b.Stat != StatDeathWard || b.Object != 0 || b.Player != p {
			continue
		}
		g.Boosts = append(g.Boosts[:i:i], g.Boosts[i+1:]...)
		pl := &g.Players[p]
		maxHP := g.PlayerHP(p) + pl.Damage
		pl.Damage = min(pl.Damage, max(maxHP-1, 0))
		g.emit(Event{Kind: EvPrevented, Player: p, Text: "death"})
		if p == g.Turn.Active {
			g.EndTurnNow()
		}
		return true
	}
	return false
}

// useNextAttackRoll returns and uses up the "next attack roll this turn"
// boosts that apply to p's attack roll: p's own and anyone's (R-ATK-19).
func (g *Game) useNextAttackRoll(p PlayerID) int {
	total := 0
	kept := g.Boosts[:0]
	for _, b := range g.Boosts {
		if b.Stat == StatNextAttackRoll && b.Object == 0 && (b.Player == p || b.Player == NoPlayer) {
			total += b.Amount
			continue
		}
		kept = append(kept, b)
	}
	g.Boosts = kept
	return total
}

// RerollWithDice rerolls a roll on the stack with several dice; chooser
// picks the result ("Its controller rerolls it, but rolls 2 dice
// instead. … That player chooses one of the rolls": Classic Roller).
func (g *Game) RerollWithDice(stackID, dice int, chooser PlayerID) {
	if _, ok := g.StackItemByID(stackID); !ok {
		return
	}
	g.ManyDice = append([]ManyDice{{Chooser: chooser, Dice: dice, StackID: stackID}}, g.ManyDice...)
	g.Reroll(stackID)
}
