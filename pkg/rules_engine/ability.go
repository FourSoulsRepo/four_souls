package rulesengine

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
}

// Ctx is what an effect sees while it resolves.
type Ctx struct {
	G          *Game
	Controller PlayerID
	Source     ObjectID
	Targets    []Chosen
}

// Chosen is one chosen target.
type Chosen struct {
	Kind    TargetKind `json:"kind"`
	Player  PlayerID   `json:"player"`
	Object  ObjectID   `json:"object,omitempty"`
	StackID int        `json:"stack_id,omitempty"`
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

func (destroySelfCost) canPay(g *Game, _ PlayerID, self ObjectID) bool { return !g.def(self).Eternal }
func (destroySelfCost) pay(g *Game, p PlayerID, self ObjectID)         { g.destroyItem(p, self) }
func (destroySelfCost) label() string                                  { return "destroy this" }

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
)

// TargetSpec says what to choose; the engine always asks (ADR 005).
type TargetSpec struct {
	Kind TargetKind
}

// Choose builds a target spec: Choose(TargetMonsterOrPlayer).
func Choose(k TargetKind) TargetSpec { return TargetSpec{Kind: k} }

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

// DealDamage puts n damage on the stack aimed at target number t (R-MECH-15).
func DealDamage(n, t int) Effect { return damageEffect{n, t} }

func (e damageEffect) apply(c *Ctx) {
	t := c.Targets[e.target]
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
		RollFor: c.abilityRef(),
	})
}

// abilityRef finds the ability of the current context; set by resolve.
func (c *Ctx) abilityRef() AbilityRef { return c.G.resolving }

// Results fills a roll table: Results(1, 2, Loot(1)) for "1-2: Loot 1".
func (t RollTable) Results(from, to int, effects ...Effect) RollTable {
	for r := from; r <= to; r++ {
		t[r-1] = effects
	}
	return t
}
