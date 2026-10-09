package rulesengine

import "slices"

// AbilityKind is how an ability is used (R-ABIL).
type AbilityKind int

// The ability kinds.
const (
	Activated   AbilityKind = iota // ↷ or $ ability (R-ABIL-08)
	LootAbility                    // on a loot card, done when it resolves (R-ABIL-17)
	Triggered                      // goes on the stack when its trigger happens (R-ABIL-13)
)

// Ability is one ability of a card, built from blocks: costs, targets,
// effects and, for triggered abilities, a trigger. See the wiki page
// "Effect blocks" for every block with an example.
type Ability struct {
	Kind    AbilityKind
	Text    string // shown to players, e.g. "↷: Gain 3¢."
	Costs   []Cost
	Targets []TargetSpec
	Effects []Effect
	Trigger Trigger
	// Modes are "choose one-" options, picked on activation before the
	// targets (R-ABIL-04). An ability with modes uses the targets and
	// effects of the chosen mode instead of its own.
	Modes []Mode
}

// Mode is one "choose one-" option.
type Mode struct {
	Text    string
	Targets []TargetSpec
	Effects []Effect
}

// Ctx is what an effect sees while it resolves.
type Ctx struct {
	G          *Game
	Controller PlayerID
	Source     ObjectID
	Targets    []Chosen
	// EventPlayer and EventAmount describe the event that triggered a
	// triggered ability: who rolled, how much damage was taken.
	EventPlayer PlayerID
	EventAmount int
	// EventStack is the stack item the event was about: the roll a
	// "would roll" trigger may change.
	EventStack int

	// Where the effect is, so an Ask can find it again.
	ref    AbilityRef
	mode   int
	roll   int
	effect int
}

// Chosen is one chosen target.
type Chosen struct {
	Kind    TargetKind `json:"kind"`
	Player  PlayerID   `json:"player"`
	Object  ObjectID   `json:"object,omitempty"`
	StackID int        `json:"stack_id,omitempty"`
	// Spent: a cost used this choice up (a discarded card, a given
	// item); it is not checked again on resolution.
	Spent bool `json:"spent,omitempty"`
}

// --- Costs (R-ABIL-07) ---

// Cost is paid when an ability is activated; it must be payable.
type Cost interface {
	canPay(g *Game, p PlayerID, self ObjectID) bool
	pay(g *Game, p PlayerID, self ObjectID)
	label() string
}

type tapCost struct{}

// Tap deactivates the object: the ↷ cost (R-ABIL-10).
func Tap() Cost { return tapCost{} }

func (tapCost) canPay(g *Game, _ PlayerID, self ObjectID) bool { return g.Object(self).Charged }
func (tapCost) pay(g *Game, _ PlayerID, self ObjectID)         { g.Object(self).Charged = false }
func (tapCost) label() string                                  { return "↷" }

type centsCost struct{ n int }

// PayCents pays n¢ to the game's pool (R-MECH-44).
func PayCents(n int) Cost { return centsCost{n} }

func (c centsCost) canPay(g *Game, p PlayerID, _ ObjectID) bool { return g.Players[p].Cents >= c.n }
func (c centsCost) pay(g *Game, p PlayerID, _ ObjectID) {
	g.Players[p].Cents -= c.n
	g.emit(Event{Kind: EvPaid, Player: p, Amount: c.n})
}
func (c centsCost) label() string { return "pay ¢" }

type destroySelfCost struct{}

// DestroySelf destroys the object as the cost.
func DestroySelf() Cost { return destroySelfCost{} }

func (destroySelfCost) canPay(g *Game, _ PlayerID, self ObjectID) bool { return !g.Eternal(self) }
func (destroySelfCost) pay(g *Game, p PlayerID, self ObjectID)         { g.destroyItem(p, self) }
func (destroySelfCost) label() string                                  { return "destroy this" }

type removeCountersCost struct{ n int }

// RemoveCounters removes n counters from the object: "Remove 2 counters
// from this:".
func RemoveCounters(n int) Cost { return removeCountersCost{n} }

func (c removeCountersCost) canPay(g *Game, _ PlayerID, self ObjectID) bool {
	return g.Object(self).CountersOf("") >= c.n
}

func (c removeCountersCost) pay(g *Game, p PlayerID, self ObjectID) {
	g.Object(self).addCounters("", -c.n)
	g.emit(Event{Kind: EvCounters, Player: p, Object: self, Card: g.Object(self).Card, Amount: -c.n})
}
func (c removeCountersCost) label() string { return "remove counters" }

// --- Targets (R-ABIL-04 to R-ABIL-06) ---

// TargetKind is what a target may be.
type TargetKind int

// The target kinds.
const (
	TargetPlayer          TargetKind = iota // a living player
	TargetMonster                           // a monster in play
	TargetMonsterOrPlayer                   // either
	TargetItem                              // an item a player controls
	TargetDiceRoll                          // a dice roll on the stack
	TargetOtherPlayer                       // a living player other than you
	TargetStackAbility                      // a ↷ or $ ability of an item, or a loot being played, on the stack
	TargetCurse                             // a curse a player has
	TargetYourItem                          // an item you control
	TargetYourHandCard                      // a loot card in your hand (for costs)
	TargetItemOrSoul                        // an item or a soul a player controls
)

// TargetSpec says what to choose; the engine always asks (ADR 005).
// Where, if set, narrows the options further; self is the ability's own
// object (the loot card, for a loot ability).
type TargetSpec struct {
	Kind  TargetKind
	Where func(g *Game, self ObjectID, c Chosen) bool
}

// Choose builds a target spec: Choose(TargetMonsterOrPlayer).
func Choose(k TargetKind) TargetSpec { return TargetSpec{Kind: k} }

// ChooseWhere builds a target spec with a filter: "the player with the
// most souls".
func ChooseWhere(k TargetKind, where func(g *Game, self ObjectID, c Chosen) bool) TargetSpec {
	return TargetSpec{Kind: k, Where: where}
}

// --- Effects ---

// Effect does one thing when an ability resolves, by queueing actions or
// putting things on the stack, so replacements and choices apply.
type Effect interface {
	apply(c *Ctx)
}

// EffectFunc is a custom effect for a card the blocks do not cover.
type EffectFunc func(c *Ctx)

func (f EffectFunc) apply(c *Ctx) { f(c) }

type lootEffect struct{ n int }

// Loot makes the controller loot n (R-CARD-06).
func Loot(n int) Effect { return lootEffect{n} }

func (e lootEffect) apply(c *Ctx) {
	c.G.enqueue(Action{Kind: ActLoot, Player: c.Controller, Amount: e.n})
}

type gainCentsEffect struct{ n int }

// GainCents gives the controller n¢ (R-MECH-36).
func GainCents(n int) Effect { return gainCentsEffect{n} }

func (e gainCentsEffect) apply(c *Ctx) {
	c.G.enqueue(Action{Kind: ActGainCents, Player: c.Controller, Amount: e.n})
}

type loseCentsEffect struct{ n int }

// LoseCents takes n¢ from the controller (R-MECH-42).
func LoseCents(n int) Effect { return loseCentsEffect{n} }

func (e loseCentsEffect) apply(c *Ctx) {
	c.G.enqueue(Action{Kind: ActLoseCents, Player: c.Controller, Amount: e.n})
}

type gainTreasureEffect struct{ n int }

// GainTreasure gives the controller the top n treasure cards (R-CARD-02).
func GainTreasure(n int) Effect { return gainTreasureEffect{n} }

func (e gainTreasureEffect) apply(c *Ctx) {
	c.G.enqueue(Action{Kind: ActGainTreasure, Player: c.Controller, Amount: e.n})
}

type damageEffect struct{ n, target int }

// DealDamage puts n damage on the stack aimed at target number t
// (R-MECH-15). DealDamage(n, You) is "take n damage" (R-MECH-18).
func DealDamage(n, t int) Effect { return damageEffect{n, t} }

func (e damageEffect) apply(c *Ctx) {
	if e.n <= 0 {
		return // nobody deals or takes 0 damage (R-MECH-20)
	}
	t := c.target(e.target)
	tgt := Target{Object: t.Object}
	if t.Kind == TargetPlayer {
		tgt = Target{Player: t.Player, IsPlayer: true}
	}
	c.G.push(StackItem{Kind: StackDamage, Controller: c.Controller, Source: c.Source, Amount: e.n, Label: "damage", Target: tgt})
}

type extraLootEffect struct{ n int }

// AddLootPlays gives the controller n more loot plays this turn.
func AddLootPlays(n int) Effect { return extraLootEffect{n} }

func (e extraLootEffect) apply(c *Ctx) {
	if c.Controller == c.G.Turn.Active {
		c.G.Turn.LootPlays += e.n
	} else {
		c.G.Players[c.Controller].ExtraLootPlays += e.n
	}
}

type rerollEffect struct{ target int }

// RerollRoll rerolls a dice roll on the stack: the same roll gets a new
// result (R-MECH-47).
func RerollRoll(t int) Effect { return rerollEffect{t} }

func (e rerollEffect) apply(c *Ctx) {
	id := c.Targets[e.target].StackID
	for i := range c.G.Stack {
		if c.G.Stack[i].ID == id {
			c.G.Stack[i].Roll = c.G.d6()
			c.G.emit(Event{Kind: EvDiceRolled, Player: c.G.Stack[i].Controller, Amount: c.G.Stack[i].Roll, Text: "reroll"})
		}
	}
}

// RollTable maps roll results to effects: index 0 is a roll of 1.
type RollTable [6][]Effect

type rollEffect struct{ table RollTable }

// Roll makes the controller roll; when the roll resolves, the effects for
// its result happen (R-ABIL-24).
func Roll(t RollTable) Effect { return rollEffect{t} }

func (e rollEffect) apply(c *Ctx) {
	r := c.G.d6()
	c.G.emit(Event{Kind: EvDiceRolled, Player: c.Controller, Amount: r, Text: "roll"})
	c.G.push(StackItem{
		Kind: StackRoll, Controller: c.Controller, Source: c.Source, Roll: r, Label: "roll",
		RollFor: c.ref, Mode: c.mode, Targets: c.Targets,
		EventPlayer: c.EventPlayer, EventAmount: c.EventAmount, EventStack: c.EventStack,
	})
}

// Results fills a roll table: Results(1, 2, Loot(1)) for "1-2: Loot 1".
func (t RollTable) Results(from, to int, effects ...Effect) RollTable {
	for r := from; r <= to; r++ {
		t[r-1] = effects
	}
	return t
}

// --- More effects (step 5.3) ---

type rechargeSelfEffect struct{}

// RechargeSelf recharges the ability's own object (R-MECH-22).
func RechargeSelf() Effect { return rechargeSelfEffect{} }

func (rechargeSelfEffect) apply(c *Ctx) {
	o := c.G.Object(c.Source)
	if o.Zone.Kind != ZoneInPlay || o.Charged {
		return
	}
	o.Charged = true
	c.G.emit(Event{Kind: EvRecharged, Player: c.Controller, Object: c.Source, Card: o.Card})
}

type addCounterEffect struct{ n int }

// AddCounters puts n counters on the ability's own object.
func AddCounters(n int) Effect { return addCounterEffect{n} }

func (e addCounterEffect) apply(c *Ctx) {
	o := c.G.Object(c.Source)
	if o.Zone.Kind != ZoneInPlay {
		return
	}
	o.addCounters("", e.n)
	c.G.emit(Event{Kind: EvCounters, Player: c.Controller, Object: c.Source, Card: o.Card, Amount: e.n})
}

type modifyRollEffect struct{ n, target int }

// ModifyRoll adds n (or subtracts, if negative) to a dice roll on the
// stack; a result stays between 1 and 6 (R-DICE-01, R-DICE-04).
func ModifyRoll(n, t int) Effect { return modifyRollEffect{n, t} }

func (e modifyRollEffect) apply(c *Ctx) {
	id := c.Targets[e.target].StackID
	for i := range c.G.Stack {
		if it := &c.G.Stack[i]; it.ID == id {
			it.Roll = min(max(it.Roll+e.n, 1), 6)
			c.G.emit(Event{Kind: EvRollChanged, Player: it.Controller, Amount: it.Roll})
		}
	}
}

type becomeSoulEffect struct{}

// BecomeSoul turns the ability's own object into a soul of its
// controller; as a soul it has no abilities (R-CARD-18).
func BecomeSoul() Effect { return becomeSoulEffect{} }

func (becomeSoulEffect) apply(c *Ctx) {
	o := c.G.Object(c.Source)
	if o.Zone.Kind != ZoneInPlay || o.Controller != c.Controller {
		return
	}
	o.Role, o.Charged, o.Counters = RoleSoul, false, nil
	c.G.emit(Event{Kind: EvGainedSoul, Player: c.Controller, Object: c.Source, Card: o.Card})
}

type stealCentsEffect struct{ n, target int }

// StealCents takes up to n¢ from the target player and gives them to
// the controller (R-MECH-38).
func StealCents(n, t int) Effect { return stealCentsEffect{n, t} }

func (e stealCentsEffect) apply(c *Ctx) {
	c.G.enqueue(Action{Kind: ActStealCents, Player: c.Controller, From: c.target(e.target).Player, Amount: e.n})
}

type boostEffect struct {
	player, monster Stat // the stat for a player or a monster target
	n, target       int
}

// GainATKThisTurn gives a player or monster +n ATK till end of turn
// (R-TURN-13 ends it).
func GainATKThisTurn(n, t int) Effect { return boostEffect{StatPlayerATK, StatMonsterATK, n, t} }

// GainHPThisTurn gives a player or monster +n HP till end of turn.
func GainHPThisTurn(n, t int) Effect { return boostEffect{StatPlayerHP, StatMonsterHP, n, t} }

// GainRollBonusThisTurn gives a player +n to their dice rolls till end
// of turn; it applies as each roll resolves (R-DICE-06).
func GainRollBonusThisTurn(n, t int) Effect { return boostEffect{StatRoll, StatRoll, n, t} }

func (e boostEffect) apply(c *Ctx) {
	t := c.target(e.target)
	b := Boost{Stat: e.monster, Player: NoPlayer, Object: t.Object, Amount: e.n}
	if t.Kind == TargetPlayer {
		b = Boost{Stat: e.player, Player: t.Player, Amount: e.n}
	}
	c.G.Boosts = append(c.G.Boosts, b)
	c.G.emit(Event{Kind: EvBoosted, Player: b.Player, Object: b.Object, Amount: e.n, Text: b.Stat.String()})
}

type shieldEffect struct{ n, target int }

// PreventNextDamage prevents the next instance of damage the target
// would take this turn (R-MECH-46).
func PreventNextDamage(t int) Effect { return shieldEffect{0, t} }

// PreventDamage prevents up to n of the next instance of damage the
// target would take this turn: "Prevent the next 1 damage".
func PreventDamage(n, t int) Effect { return shieldEffect{n, t} }

func (e shieldEffect) apply(c *Ctx) {
	t := c.target(e.target)
	s := Shield{Target: Target{Object: t.Object}, Amount: e.n, Source: c.Source}
	if t.Kind == TargetPlayer {
		s.Target = Target{Player: t.Player, IsPlayer: true}
	}
	c.G.Shields = append(c.G.Shields, s)
	c.G.emit(Event{Kind: EvShielded, Player: s.Target.Player, Object: s.Target.Object, Amount: e.n})
}

type rechargeEffect struct{ target int }

// RechargeTarget recharges the target item ("Recharge an item").
func RechargeTarget(t int) Effect { return rechargeEffect{t} }

func (e rechargeEffect) apply(c *Ctx) { c.G.Recharge(c.target(e.target).Object) }

type rechargeAllEffect struct{ target int }

// RechargeItemsOf recharges each item the target player controls.
func RechargeItemsOf(t int) Effect { return rechargeAllEffect{t} }

func (e rechargeAllEffect) apply(c *Ctx) {
	for _, id := range c.G.Players[c.target(e.target).Player].InPlay {
		if c.G.Object(id).Role == RoleItem {
			c.G.Recharge(id)
		}
	}
}

// Recharge turns an object in play upright (R-MECH-22).
func (g *Game) Recharge(id ObjectID) {
	o := g.Object(id)
	if o.Zone.Kind != ZoneInPlay || o.Charged {
		return
	}
	o.Charged = true
	g.emit(Event{Kind: EvRecharged, Player: o.Controller, Object: id, Card: o.Card})
}

type killEffect struct{ target int }

// Kill kills the target player or monster: its death goes on the stack
// (R-MECH-23, R-DEATH-01).
func Kill(t int) Effect { return killEffect{t} }

func (e killEffect) apply(c *Ctx) {
	t := c.target(e.target)
	if t.Kind == TargetPlayer {
		c.G.Players[t.Player].Damage += c.G.PlayerHP(t.Player) // HP to 0 (R-MECH-23)
		c.G.push(StackItem{Kind: StackDeath, Controller: NoPlayer, Label: "death", Target: Target{Player: t.Player, IsPlayer: true}})
		return
	}
	if c.G.Eternal(t.Object) {
		return // R-DEATH-03
	}
	c.G.Object(t.Object).Damage = c.G.def(t.Object).HP + c.G.bonus(StatMonsterHP, NoPlayer, t.Object)
	c.G.push(StackItem{Kind: StackDeath, Controller: NoPlayer, Label: "death", Target: Target{Object: t.Object}})
}

// You is the target index that means the ability's controller, for text
// without a target: DealDamage(1, You) is "Take 1 damage".
const You = -1

// target returns chosen target t, or the controller for You.
func (c *Ctx) target(t int) Chosen {
	if t == You {
		return Chosen{Kind: TargetPlayer, Player: c.Controller}
	}
	return c.Targets[t]
}

// effectsOf returns the effects that run for an ability: those of the
// chosen mode, or of a roll result.
func (g *Game) effectsOf(ref AbilityRef, mode, roll int) []Effect {
	ab := g.ability(ref)
	effects := ab.Effects
	if len(ab.Modes) > 0 {
		effects = ab.Modes[mode].Effects
	}
	if roll > 0 {
		for _, e := range effects {
			if r, ok := e.(rollEffect); ok {
				return r.table[roll-1]
			}
		}
		return nil
	}
	return effects
}

type cancelEffect struct{ target int }

// CancelTarget cancels the target stack item: "Cancel the ↷ or $
// ability of an item or a loot being played".
func CancelTarget(t int) Effect { return cancelEffect{t} }

func (e cancelEffect) apply(c *Ctx) { c.G.CancelStackItem(c.target(e.target).StackID) }

type destroyEffect struct{ target int }

// DestroyTarget destroys the target item or curse (R-MECH-23).
func DestroyTarget(t int) Effect { return destroyEffect{t} }

func (e destroyEffect) apply(c *Ctx) {
	t := c.target(e.target)
	c.G.DestroyObject(t.Player, t.Object)
}

type eachPlayerEffect struct {
	effects []Effect
	others  bool
}

// EachPlayer runs the effects once for each player, in turn order from
// the controller, as if that player controlled them (R-MECH-27): "Each
// player gains 1¢" is EachPlayer(GainCents(1)). Damage to You goes on
// the stack in reverse order, so it resolves in turn order (R-MECH-28).
func EachPlayer(effects ...Effect) Effect { return eachPlayerEffect{effects: effects} }

func (e eachPlayerEffect) apply(c *Ctx) {
	order := c.G.turnOrderFrom(c.Controller)
	damage := false
	for _, x := range e.effects {
		if _, ok := x.(damageEffect); ok {
			damage = true
		}
	}
	if damage {
		slices.Reverse(order)
	}
	for _, p := range order {
		if e.others && p == c.Controller {
			continue
		}
		each := *c
		each.Controller = p
		each.Do(e.effects...)
	}
}

type monstersDamageEffect struct{ n int }

// EachMonsterTakesDamage puts n damage on the stack for each monster in
// play, in slot order.
func EachMonsterTakesDamage(n int) Effect { return monstersDamageEffect{n} }

func (e monstersDamageEffect) apply(c *Ctx) {
	for _, s := range c.G.Monsters {
		if top, ok := s.TopOf(); ok && c.G.Object(top).Role == RoleMonster {
			c.G.push(StackItem{Kind: StackDamage, Controller: c.Controller, Source: c.Source, Amount: e.n, Label: "damage", Target: Target{Object: top}})
		}
	}
}

type addAttacksEffect struct{ n, target int }

// AddAttacks lets the target player attack n more times this turn, if
// it is their turn: "may attack an additional time this turn".
func AddAttacks(n, t int) Effect { return addAttacksEffect{n, t} }

func (e addAttacksEffect) apply(c *Ctx) {
	if c.target(e.target).Player == c.G.Turn.Active {
		c.G.Turn.Attacks += e.n
	}
}

type sunEffect struct{}

// ThisToLootBottomExtraTurn puts the loot card on the bottom of the loot
// deck; if it does and it is the controller's turn, they take an extra
// turn after this one (The Sun).
func ThisToLootBottomExtraTurn() Effect { return sunEffect{} }

func (sunEffect) apply(c *Ctx) {
	g := c.G
	if g.Object(c.Source).Zone.Kind != ZoneStack {
		return
	}
	nid := g.move(c.Source, DeckZone(LootDeck), NoPlayer)
	g.Decks[LootDeck] = append([]ObjectID{nid}, g.Decks[LootDeck]...)
	g.emit(Event{Kind: EvMovedToDeck, Player: c.Controller, Object: nid, Card: g.Object(nid).Card, Text: "loot deck bottom"})
	if c.Controller == g.Turn.Active {
		g.ExtraTurn = true
	}
}

type preventDeathEffect struct{}

// PreventYourDeath removes the controller's death from the stack and
// heals them to 1 HP (R-MECH-46): "Prevent death".
func PreventYourDeath() Effect { return preventDeathEffect{} }

func (preventDeathEffect) apply(c *Ctx) {
	g := c.G
	for _, it := range g.Stack {
		if it.Kind == StackDeath && it.Target.IsPlayer && it.Target.Player == c.Controller {
			g.CancelStackItem(it.ID)
			pl := &g.Players[c.Controller]
			maxHP := g.PlayerHP(c.Controller) + pl.Damage
			pl.Damage = min(pl.Damage, max(maxHP-1, 0))
			g.emit(Event{Kind: EvPrevented, Player: c.Controller, Text: "death"})
			return
		}
	}
}

// NotThis is a target filter: anything but the ability's own object
// ("another item").
func NotThis(_ *Game, self ObjectID, c Chosen) bool { return c.Object != self }

// CapNextDamage reduces the next instance of damage the target takes
// this turn to at most n: "is reduced to 1".
func CapNextDamage(n, t int) Effect { return capEffect{n, t} }

type capEffect struct{ n, target int }

func (e capEffect) apply(c *Ctx) {
	t := c.target(e.target)
	s := Shield{Target: Target{Object: t.Object}, Cap: e.n, Source: c.Source}
	if t.Kind == TargetPlayer {
		s.Target = Target{Player: t.Player, IsPlayer: true}
	}
	c.G.Shields = append(c.G.Shields, s)
	c.G.emit(Event{Kind: EvShielded, Player: s.Target.Player, Object: s.Target.Object, Amount: e.n, Text: "cap"})
}

// EachOtherPlayer is EachPlayer without the controller.
func EachOtherPlayer(effects ...Effect) Effect {
	return eachPlayerEffect{effects: effects, others: true}
}

// DestroyThis destroys the ability's own object as an effect and runs
// then only if it was destroyed: "↷: Destroy this. If you do, …".
func DestroyThis(then ...Effect) Effect { return destroySelfThen{then} }

type destroySelfThen struct{ then []Effect }

func (e destroySelfThen) apply(c *Ctx) {
	if c.G.DestroyObject(c.Controller, c.Source) {
		c.Do(e.then...)
	}
}

type rerollItemEffect struct{ target int }

// RerollTarget rerolls the target item (R-MECH-48).
func RerollTarget(t int) Effect { return rerollItemEffect{t} }

func (e rerollItemEffect) apply(c *Ctx) { c.G.RerollItem(c.target(e.target).Object) }

// chosenCost is a cost paid with something chosen on activation, like a
// target (R-ABIL-07): "Discard a loot card:".
type chosenCost interface {
	payChosen(g *Game, p PlayerID, chosen []Chosen)
}

type discardChosenCost struct{ target int }

// DiscardChosen discards the hand card chosen as target t, as a cost:
// Targets: Choose(TargetYourHandCard), Costs: DiscardChosen(0).
func DiscardChosen(t int) Cost { return discardChosenCost{t} }

func (discardChosenCost) canPay(g *Game, p PlayerID, _ ObjectID) bool {
	return len(g.Players[p].Hand) > 0
}
func (discardChosenCost) pay(*Game, PlayerID, ObjectID) {}
func (discardChosenCost) label() string                 { return "discard a loot card" }
func (c discardChosenCost) payChosen(g *Game, p PlayerID, chosen []Chosen) {
	g.DiscardFromHand(p, chosen[c.target].Object)
	chosen[c.target].Spent = true
}

type giveChosenCost struct{ item, player int }

// GiveChosen gives the item chosen as target item to the player chosen
// as target player, as a cost: "Give an item you control to another
// player:".
func GiveChosen(item, player int) Cost { return giveChosenCost{item, player} }

func (giveChosenCost) canPay(*Game, PlayerID, ObjectID) bool { return true }
func (giveChosenCost) pay(*Game, PlayerID, ObjectID)         {}
func (giveChosenCost) label() string                         { return "give an item" }
func (c giveChosenCost) payChosen(g *Game, _ PlayerID, chosen []Chosen) {
	g.GainControl(chosen[c.player].Player, chosen[c.item].Object)
	chosen[c.item].Spent = true
}

type destroyChosenCost struct{ targets []int }

// DestroyChosen destroys the items chosen as the given targets, as a
// cost: "Destroy 2 items you control:".
func DestroyChosen(targets ...int) Cost { return destroyChosenCost{targets} }

func (destroyChosenCost) canPay(*Game, PlayerID, ObjectID) bool { return true }
func (destroyChosenCost) pay(*Game, PlayerID, ObjectID)         {}
func (destroyChosenCost) label() string                         { return "destroy items" }
func (c destroyChosenCost) payChosen(g *Game, p PlayerID, chosen []Chosen) {
	for _, t := range c.targets {
		g.DestroyObject(p, chosen[t].Object)
		chosen[t].Spent = true
	}
}

// NotChosen is a target filter: not something already chosen for this
// activation ("2 items" means two different ones).
func NotChosen(g *Game, _ ObjectID, c Chosen) bool {
	if g.Activating == nil {
		return true
	}
	for _, x := range g.Activating.Chosen {
		if x.Object == c.Object && x.Kind == c.Kind {
			return false
		}
	}
	return true
}

// Ability returns the definition of an ability.
func (g *Game) Ability(ref AbilityRef) Ability { return g.ability(ref) }
