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
	// StepRoom comes between StepHandSize and StepCleanup; it is listed
	// last so the older steps keep their numbers (R-TURN-12).
	StepRoom
)

// Turn is the state of the current turn.
type Turn struct {
	Active      PlayerID `json:"active"`
	Number      int      `json:"number"`
	Step        Step     `json:"step"`
	LootPlays   int      `json:"loot_plays"`
	Attacks     int      `json:"attacks"`
	Purchases   int      `json:"purchases"`
	AttackRolls int      `json:"attack_rolls,omitempty"` // attack rolls resolved this turn
	// Bonus*Used count loot plays, attacks and purchases taken from
	// static bonuses ("you may attack an additional time on your turn").
	BonusLootUsed      int `json:"bonus_loot_used,omitempty"`
	BonusAttacksUsed   int `json:"bonus_attacks_used,omitempty"`
	BonusPurchasesUsed int `json:"bonus_purchases_used,omitempty"`
	// MustAttack is a monster the active player must attack this turn if
	// able (Monster Manual).
	MustAttack ObjectID `json:"must_attack,omitempty"`
	// MustAttacks is how many more attacks the active player must make;
	// MustAttackDeck of them on the monster deck.
	MustAttacks    int `json:"must_attacks,omitempty"`
	MustAttackDeck int `json:"must_attack_deck,omitempty"`
	// DeckAttacks are extra attacks only on the monster deck.
	DeckAttacks int `json:"deck_attacks,omitempty"`
	// ShopBought counts shop items purchased this turn; FreeShopItem makes
	// the first one cost 0¢ (Credit Card).
	ShopBought   int  `json:"shop_bought,omitempty"`
	FreeShopItem bool `json:"free_shop_item,omitempty"`
	EndDeclared  bool `json:"end_declared,omitempty"`
	DeathEnd     bool `json:"death_end,omitempty"` // the active player died (R-DEATH-16)
	// NextHelpers help the active player's next attack this turn;
	// NextShare also gain its rewards (Sibling Rivalry).
	NextHelpers []PlayerID `json:"next_helpers,omitempty"`
	NextShare   []PlayerID `json:"next_share,omitempty"`
	// MonsterDied: a monster died this turn, so the active player may
	// replace a room (R-TURN-12). RoomReplaced opens a priority window
	// for the new room's abilities.
	MonsterDied  bool `json:"monster_died,omitempty"`
	RoomReplaced bool `json:"room_replaced,omitempty"`
	// LateRooms entered in the end phase: their ↷ and $ abilities wait
	// for the next turn (R-TURN-18).
	LateRooms []ObjectID `json:"late_rooms,omitempty"`
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
	case StepRoom:
		g.goTo(StepCleanup)
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
		g.checkMaxHP()
		g.checkBonusSouls()
		if g.checkWin() {
			return
		}
		if g.FirstRoll != nil {
			// The roll comes before start-of-game abilities (Eden).
			g.askFirstThrow()
			continue
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
		g.goTo(StepRoom)
	case StepRoom:
		if g.Turn.RoomReplaced {
			g.Turn.RoomReplaced = false
			g.openWindow(g.Turn.Active) // the new room's "enters play" (R-TURN-18)
			return
		}
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
		g.endCopiesTillYourTurn(p)
		g.recharge(p)
	case StepStartTriggers:
		g.emit(Event{Kind: EvStartOfTurn, Player: p}) // R-TURN-03
		g.openWindow(p)
	case StepLoot:
		g.enqueue(Action{Kind: ActLoot, Player: p, Amount: 1 + g.bonus(StatLootStep, p, 0)}) // R-TURN-04
	case StepLootWindow:
		g.openWindow(p)
	case StepAction:
		g.Turn.LootPlays, g.Turn.Attacks, g.Turn.Purchases = 1, 1, 1 // R-TURN-05, R-TURN-07
		g.openWindow(p)
	case StepEndTriggers:
		g.emit(Event{Kind: EvEndOfTurn, Player: p}) // R-TURN-10
		g.openWindow(p)
	case StepHandSize:
		if extra := len(g.Players[p].Hand) - g.maxHand(p); extra > 0 {
			g.Waiting = Prompt{Kind: PromptDiscard, Player: p, Count: extra}
		}
	case StepRoom:
		g.askRoomStep(p)
	case StepCleanup:
		g.healAll()
		g.emit(Event{Kind: EvTurnEnded, Player: p})
	}
}

func (g *Game) startNextTurn() {
	next := g.next(g.Turn.Active)
	if g.ExtraTurn {
		next, g.ExtraTurn = g.Turn.Active, false
		g.emit(Event{Kind: EvExtraTurn, Player: next})
	}
	for range g.Players {
		if g.Players[next].SkipTurns == 0 {
			break
		}
		g.Players[next].SkipTurns--
		g.emit(Event{Kind: EvTurnSkipped, Player: next})
		next = g.next(next)
	}
	g.Turn = Turn{Active: next, Number: g.Turn.Number + 1, Step: StepRecharge}
	g.emit(Event{Kind: EvTurnStarted, Player: g.Turn.Active, Amount: g.Turn.Number})
}

// recharge charges what the active player controls (R-TURN-02).
func (g *Game) recharge(p PlayerID) {
	pl := &g.Players[p]
	if g.bonus(StatNoCharRecharge, p, 0) == 0 {
		g.Object(pl.Character).Charged = true
	}
	if g.bonus(StatRechargeLimit, p, 0) > 0 {
		// "You can't recharge more than 1 item during your recharge step."
		var items []ObjectID
		for _, id := range pl.InPlay {
			if o := g.Object(id); o.Role == RoleItem && !o.Charged {
				items = append(items, id)
			}
		}
		if len(items) > 1 {
			g.ask(Choice{Purpose: ChooseRecharge, Player: p, Rule: "R-TURN-02", Objects: items}, g.Labels(items))
			return
		}
		for _, id := range items {
			g.Object(id).Charged = true
		}
		g.emit(Event{Kind: EvRecharged, Player: p})
		return
	}
	for _, id := range pl.InPlay {
		if !g.def(id).OwnRechargeOnly {
			g.Object(id).Charged = true
		}
	}
	g.emit(Event{Kind: EvRecharged, Player: p})
}

// endCopiesTillYourTurn ends p's copies that last till the start of
// their next turn (Undefined).
func (g *Game) endCopiesTillYourTurn(p PlayerID) {
	for _, id := range g.Players[p].InPlay {
		o := g.Object(id)
		if o.CopyTillYourTurn {
			o.CopyOf, o.CopyTillYourTurn = "", false
		}
		if o.GainedTill == GainTillYourTurn {
			o.Gained, o.GainedTill = "", 0
		}
	}
}

// loot draws n loot cards into a hand (R-CARD-06).
func (g *Game) loot(p PlayerID, n int) {
	drawn := 0
	for range n {
		id, ok := g.drawTop(LootDeck)
		if !ok {
			break
		}
		nid := g.move(id, Zone{Kind: ZoneHand}, p)
		g.Players[p].Hand = append(g.Players[p].Hand, nid)
		g.emit(Event{Kind: EvLooted, Player: p, Object: nid, Card: g.Object(nid).Card, Private: true})
		drawn++
	}
	if drawn > 0 && g.Turn.Number > 0 {
		g.emit(Event{Kind: EvLootDone, Player: p, Amount: drawn})
	}
}

// DiscardFromHand puts a loot card from p's hand into the loot discard
// (R-MECH-25). It is for effects: a player whom abilities can't make
// discard keeps the card (Rib of Greed). Costs use discardFromHand.
func (g *Game) DiscardFromHand(p PlayerID, id ObjectID) {
	if g.protectedFromLoss(p) {
		g.emit(Event{Kind: EvPrevented, Player: p, Text: "discard"})
		return
	}
	g.discardFromHand(p, id)
}

// protectedFromLoss: abilities and the death penalty can't make p discard
// loot cards or lose ¢ (Rib of Greed).
func (g *Game) protectedFromLoss(p PlayerID) bool {
	return p != NoPlayer && g.bonus(StatNoForcedLoss, p, 0) > 0
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
		g.Players[i].TimesDamaged = 0
		g.Players[i].Activations, g.Players[i].LootPlayed = 0, 0
		g.Players[i].Rolls, g.Players[i].LastRoll, g.Players[i].LastAttack, g.Players[i].KilledBy = 0, 0, 0, 0
	}
	for _, id := range g.inPlay() {
		g.Object(id).Damage, g.Object(id).HitsThisTurn = 0, 0
	}
	g.Boosts, g.Shields, g.Lingering, g.Markers = nil, nil, nil, nil // till end of turn effects end
	for _, pl := range g.Players {
		for _, id := range append([]ObjectID{pl.Character}, pl.InPlay...) {
			o := g.Object(id)
			if o.CopyThisTurn {
				o.CopyOf, o.CopyThisTurn = "", false
			}
			if o.GainedTill == GainEndOfTurn {
				o.Gained, o.GainedTill = "", 0
			}
		}
		g.Players[pl.ID].CopyNextLoot = false
	}
	g.emit(Event{Kind: EvHealed, Player: NoPlayer})
}

// SoulValue is a player's total soul value (R-CARD-21).
func (g *Game) SoulValue(p PlayerID) int {
	total := g.characterSouls(g.Players[p].Character)
	for _, id := range g.Players[p].InPlay {
		total += g.soulsOf(id)
	}
	return total
}

// soulsOf is the soul value of a soul card: at least one (R-CARD-18).
// Other cards are no souls.
func (g *Game) soulsOf(id ObjectID) int {
	o := g.Object(id)
	if o.Role != RoleSoul {
		return 0
	}
	if d, ok := g.cards.find(o.Card); ok {
		return max(d.Soul, 1)
	}
	return 0
}

// characterSouls is a character's soul value; it counts, though it is
// not a soul (R-CARD-35).
func (g *Game) characterSouls(id ObjectID) int {
	if d, ok := g.cards.find(g.CardOf(id)); ok {
		return d.Soul
	}
	return 0
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

// PlayerAfter is the player to p's left: next in turn order.
func (g *Game) PlayerAfter(p PlayerID) PlayerID { return g.next(p) }

// PlayerBefore is the player to p's right: previous in turn order.
func (g *Game) PlayerBefore(p PlayerID) PlayerID {
	return PlayerID((int(p) + len(g.Players) - 1) % len(g.Players))
}

// maxHand is p's max hand size: the game's, or a lower cap ("Your max
// hand size is 2", R-TURN-11).
func (g *Game) maxHand(p PlayerID) int {
	limit := g.MaxHand
	for _, id := range g.inPlay() {
		for _, face := range g.facesOf(id) {
			d, ok := g.cards.find(face)
			if !ok {
				continue
			}
			for _, s := range d.Statics {
				if s.Stat == StatHandCap && s.Applies(g, id, p, 0) {
					limit = min(limit, s.amount(g, id))
				}
			}
		}
	}
	return limit
}

// WinGame ends the game with p as the winner: "That player wins." (Mom's
// Heart).
func (g *Game) WinGame(p PlayerID) {
	if g.Over {
		return
	}
	g.Over, g.Winners = true, []PlayerID{p}
	g.Waiting = Prompt{Kind: PromptGameOver, Player: NoPlayer}
	g.emit(Event{Kind: EvGameWon, Player: p})
}

// askRoomStep lets the active player put a room into discard if a
// monster died this turn (R-TURN-12). A room that put itself into
// discard doesn't count (R-TURN-17).
func (g *Game) askRoomStep(p PlayerID) {
	if !g.Turn.MonsterDied {
		return
	}
	var rooms []ObjectID
	var labels []Msg
	for _, s := range g.Rooms {
		if top, ok := s.TopOf(); ok && g.Object(top).Zone.Kind == ZoneInPlay {
			if stay := g.def(top).NoDiscard; stay != nil && stay(g, top) {
				continue
			}
			rooms = append(rooms, top)
			labels = append(labels, CardOption(g.Object(top).Card))
		}
	}
	if len(rooms) == 0 {
		return
	}
	g.ask(Choice{Purpose: ChooseRoom, Player: p, Rule: "R-TURN-12", Objects: rooms, Question: promptDiscardRoom.Msg()}, append(labels, AskKeep.Msg()))
}
