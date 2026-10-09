package rulesengine

// Step is where the turn is (R-TURN).
type Step int

// The turn steps, in order.
const (
	StepRecharge      Step = iota // R-TURN-02, nobody has priority
	StepStartTriggers             // R-TURN-03
	StepLoot                      // R-TURN-04
	StepLootWindow                // priority after the loot (R-TURN-04)
	StepAction                    // R-TURN-05 to R-TURN-09
	StepEndTriggers               // R-TURN-10
	StepHandSize                  // R-TURN-11, nobody has priority
	StepCleanup                   // R-TURN-13
)

// Turn is the state of the current turn.
type Turn struct {
	Active      PlayerID `json:"active"`
	Number      int      `json:"number"`
	Step        Step     `json:"step"`
	LootPlays   int      `json:"loot_plays"`
	Attacks     int      `json:"attacks"`
	Purchases   int      `json:"purchases"`
	EndDeclared bool     `json:"end_declared,omitempty"`
	DeathEnd    bool     `json:"death_end,omitempty"` // the active player died (R-DEATH-16)
	// entered is true once the current step's automatic work is done.
	Entered bool `json:"entered,omitempty"`
}

// Priority is an open priority window (R-PRIO).
type Priority struct {
	Open   bool     `json:"open"`
	Holder PlayerID `json:"holder"`
	Passes int      `json:"passes"` // passes in a row
	// Deferred: priority goes to Holder once the action queue is empty.
	Deferred bool `json:"deferred,omitempty"`
}

// next returns the player after p in turn order (R-TURN-01).
func (g *Game) next(p PlayerID) PlayerID {
	return PlayerID((int(p) + 1) % len(g.Players))
}

// openWindow gives priority to start (R-PRIO-02, R-PRIO-03).
func (g *Game) openWindow(start PlayerID) {
	if len(g.PendingTriggers) > 0 || len(g.Queue) > 0 {
		// Triggers go on the stack and queued steps happen before anyone
		// gets priority (R-ABIL-14, R-DEATH-10).
		g.Priority = Priority{Deferred: true, Holder: start}
		g.Waiting = Prompt{}
		return
	}
	g.Priority = Priority{Open: true, Holder: start}
	g.Waiting = Prompt{Kind: PromptPriority, Player: start}
}

// openWindowAfter opens a window where p has already passed, e.g. after
// declaring the end of the turn.
func (g *Game) openWindowAfter(p PlayerID) {
	g.openWindow(g.next(p))
	g.Priority.Passes = 1
	g.closeIfAllPassed()
}

// pass moves priority on; when everyone passed in a row the window closes
// (R-PRIO-05).
func (g *Game) pass() {
	g.Priority.Passes++
	if g.closeIfAllPassed() {
		return
	}
	g.Priority.Holder = g.next(g.Priority.Holder)
	g.Waiting = Prompt{Kind: PromptPriority, Player: g.Priority.Holder}
}

func (g *Game) closeIfAllPassed() bool {
	if g.Priority.Passes < len(g.Players) {
		return false
	}
	g.Priority = Priority{}
	g.Waiting = Prompt{}
	if len(g.Stack) > 0 {
		g.resolveTop() // R-STACK-03
		return true
	}
	g.windowClosed()
	return true
}

// windowClosed moves the game on once all players passed with an empty
// stack (R-STACK-07). The stack itself comes in step 4.4.
func (g *Game) windowClosed() {
	switch g.Turn.Step { //nolint:exhaustive // other steps never open a window
	case StepStartTriggers:
		g.goTo(StepLoot)
	case StepLootWindow:
		g.goTo(StepAction)
	case StepAction:
		switch {
		case g.Turn.DeathEnd, g.Turn.EndDeclared:
			// After an active player's death the turn goes to its end phase
			// once the stack has resolved (R-DEATH-16).
			g.goTo(StepEndTriggers)
		case g.Purchase.On:
			g.askPurchase() // R-SHOP-02
		case g.Attack.On && !g.Attack.Started:
			g.askAttackTarget() // R-ATK-02
		case g.Attack.Started:
			g.continueAttack() // R-ATK-11
		default:
			// The action phase only ends when the active player ends it (R-TURN-09).
			g.openWindow(g.Turn.Active)
		}
	case StepEndTriggers:
		g.goTo(StepHandSize)
	}
}

func (g *Game) goTo(s Step) {
	g.Turn.Step = s
	g.Turn.Entered = false
}

// inOpenActionPhase is true when the active player may attack, purchase or
// end the turn: action phase, empty stack (R-TURN-06).
func (g *Game) inOpenActionPhase() bool {
	return g.Turn.Step == StepAction && !g.Turn.EndDeclared && !g.Turn.DeathEnd &&
		!g.Attack.On && !g.Purchase.On && len(g.Stack) == 0
}

// maxRunSteps stops a run loop that never settles; that is always a bug.
const maxRunSteps = 100000

// run does automatic work until the engine needs input.
func (g *Game) run() {
	for steps := 0; !g.Over && g.Waiting.Kind == PromptNone; steps++ {
		if steps > maxRunSteps {
			panic("rulesengine: the game does not settle (bug)")
		}
		if g.checkWin() {
			return
		}
		if len(g.Queue) > 0 {
			g.processQueue()
			continue
		}
		if len(g.PendingTriggers) > 0 {
			g.placeNextTrigger()
			continue
		}
		if g.Priority.Deferred {
			g.openWindow(g.Priority.Holder)
			continue
		}
		if g.Turn.Entered {
			// The step finished its work without asking anything.
			g.advanceFrom(g.Turn.Step)
			continue
		}
		g.Turn.Entered = true
		g.enterStep()
	}
}

func (g *Game) advanceFrom(s Step) {
	switch s { //nolint:exhaustive // steps with windows advance in windowClosed
	case StepRecharge:
		g.goTo(StepStartTriggers)
	case StepLoot:
		g.goTo(StepLootWindow)
	case StepHandSize:
		g.goTo(StepCleanup)
	case StepCleanup:
		g.startNextTurn()
	case StepStartTriggers, StepLootWindow, StepAction, StepEndTriggers:
		// These steps move on when their priority window closes; without
		// an open window, give priority again instead of spinning.
		g.openWindow(g.Turn.Active)
	}
}

func (g *Game) enterStep() {
	p := g.Turn.Active
	switch g.Turn.Step {
	case StepRecharge:
		g.recharge(p)
	case StepStartTriggers:
		g.emit(Event{Kind: EvStartOfTurn, Player: p}) // R-TURN-03
		g.openWindow(p)
	case StepLoot:
		g.enqueue(Action{Kind: ActLoot, Player: p, Amount: 1}) // R-TURN-04
	case StepLootWindow:
		g.openWindow(p)
	case StepAction:
		g.Turn.LootPlays, g.Turn.Attacks, g.Turn.Purchases = 1, 1, 1 // R-TURN-05, R-TURN-07
		g.openWindow(p)
	case StepEndTriggers:
		g.emit(Event{Kind: EvEndOfTurn, Player: p}) // R-TURN-10
		g.openWindow(p)
	case StepHandSize:
		if extra := len(g.Players[p].Hand) - g.MaxHand; extra > 0 {
			g.Waiting = Prompt{Kind: PromptDiscard, Player: p, Count: extra}
		}
	case StepCleanup:
		g.healAll()
		g.emit(Event{Kind: EvTurnEnded, Player: p})
	}
}

func (g *Game) startNextTurn() {
	g.Turn = Turn{Active: g.next(g.Turn.Active), Number: g.Turn.Number + 1, Step: StepRecharge}
	g.emit(Event{Kind: EvTurnStarted, Player: g.Turn.Active, Amount: g.Turn.Number})
}

// recharge charges what the active player controls (R-TURN-02).
func (g *Game) recharge(p PlayerID) {
	pl := &g.Players[p]
	g.Object(pl.Character).Charged = true
	for _, id := range pl.InPlay {
		g.Object(id).Charged = true
	}
	g.emit(Event{Kind: EvRecharged, Player: p})
}

// loot draws n loot cards into a hand (R-CARD-06).
func (g *Game) loot(p PlayerID, n int) {
	for range n {
		id, ok := g.drawTop(LootDeck)
		if !ok {
			return
		}
		nid := g.move(id, Zone{Kind: ZoneHand}, p)
		g.Players[p].Hand = append(g.Players[p].Hand, nid)
		g.emit(Event{Kind: EvLooted, Player: p, Object: nid, Card: g.Object(nid).Card, Private: true})
	}
}

func (g *Game) discardFromHand(p PlayerID, id ObjectID) {
	g.Players[p].Hand = remove(g.Players[p].Hand, id)
	nid := g.discard(id, LootDeck)
	g.emit(Event{Kind: EvDiscarded, Player: p, Object: nid, Card: g.Object(nid).Card})
}

// healAll heals every object with HP, dead players too (R-TURN-13).
func (g *Game) healAll() {
	for i := range g.Players {
		g.Players[i].Damage = 0
		g.Players[i].Dead = false // alive again (R-DEATH-19)
		g.Players[i].ExtraLootPlays = 0
	}
	for i := range g.Objects {
		if g.Objects[i].Zone.Kind == ZoneInPlay {
			g.Objects[i].Damage = 0
		}
	}
	g.emit(Event{Kind: EvHealed, Player: NoPlayer})
}

// SoulValue is a player's total soul value (R-CARD-21).
func (g *Game) SoulValue(p PlayerID) int {
	total := 0
	for _, id := range g.Players[p].InPlay {
		o := g.Object(id)
		if o.Role != RoleSoul {
			continue
		}
		if d, ok := g.cards.find(o.Card); ok {
			total += max(d.Soul, 1)
		}
	}
	return total
}

// checkWin ends the game when players reach the soul target (R-WIN-01 to
// R-WIN-03). It runs between resolutions, without the stack (R-WIN-02).
func (g *Game) checkWin() bool {
	var winners []PlayerID
	for _, pl := range g.Players {
		if g.SoulValue(pl.ID) >= g.WinSouls {
			winners = append(winners, pl.ID)
		}
	}
	if len(winners) == 0 {
		return false
	}
	g.Over, g.Winners = true, winners
	g.Waiting = Prompt{Kind: PromptGameOver, Player: NoPlayer}
	for _, w := range winners {
		g.emit(Event{Kind: EvGameWon, Player: w})
	}
	return true
}
