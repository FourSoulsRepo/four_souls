package rulesengine

import "fmt"

// PromptKind is what the engine waits for.
type PromptKind int

// The prompt kinds; more come with later steps (ADR 005).
const (
	PromptNone     PromptKind = iota
	PromptPriority            // the player may act or pass (R-PRIO-01)
	PromptDiscard             // the player discards Count loot cards
	PromptGameOver            // nothing more to do
	// PromptChoose asks the player to pick one of Options; Purpose says
	// what for (see Choice).
	PromptChoose
)

// Prompt names who must answer and what kind of answer is expected.
type Prompt struct {
	Kind    PromptKind    `json:"kind"`
	Player  PlayerID      `json:"player"`
	Count   int           `json:"count,omitempty"`
	Options []string      `json:"options,omitempty"`
	Purpose ChoicePurpose `json:"purpose,omitempty"`
	Text    string        `json:"text,omitempty"` // the question, if any
}

// IntentKind is what a player wants to do.
type IntentKind int

// The intents so far.
const (
	IntentPass     IntentKind = iota // pass priority (R-PRIO-05)
	IntentEndTurn                    // declare the end of the turn (R-TURN-06)
	IntentDiscard                    // discard the chosen loot cards
	IntentPlayLoot                   // play the loot card in Objects[0] (R-CARD-08)
	IntentChoose                     // pick option Choice of the prompt
	IntentAttack                     // declare an attack (R-ATK-01)
	IntentPurchase                   // declare a purchase (R-SHOP-01)
	IntentActivate                   // use ability Choice of Objects[0] (R-ABIL-08)
)

// Intent is one player's request. The engine checks it against the
// current prompt; nothing changes if it is not allowed.
type Intent struct {
	Player  PlayerID   `json:"player"`
	Kind    IntentKind `json:"kind"`
	Objects []ObjectID `json:"objects,omitempty"`
	Choice  int        `json:"choice,omitempty"`
}

// RuleError explains why an intent is not allowed, citing a rule ID.
type RuleError struct {
	Rule   string
	Reason string
}

func (e *RuleError) Error() string { return e.Reason + " (" + e.Rule + ")" }

func refuse(rule, format string, args ...any) error {
	return &RuleError{Rule: rule, Reason: fmt.Sprintf(format, args...)}
}

// Prompt returns what the engine waits for.
func (g *Game) Prompt() Prompt { return g.Waiting }

// Apply checks an intent and, if allowed, runs the game until it needs
// input again. It returns what happened, in order. On error the state is
// unchanged.
func (g *Game) Apply(in Intent) ([]Event, error) {
	if err := g.check(in); err != nil {
		return nil, err
	}
	g.events = nil
	g.apply(in)
	g.run()
	return g.takeEvents(), nil
}

func (g *Game) check(in Intent) error {
	w := g.Waiting
	if g.Over {
		return refuse("R-WIN-03", "the game is over")
	}
	if in.Player != w.Player {
		return refuse("R-PRIO-01", "player %d does not have priority; player %d does", in.Player, w.Player)
	}
	switch in.Kind {
	case IntentPass:
		if w.Kind != PromptPriority {
			return refuse("R-PRIO-01", "nothing to pass")
		}
		if g.inOpenActionPhase() && in.Player == g.Turn.Active {
			return refuse("R-TURN-09", "in the action phase the active player ends the turn instead of passing")
		}
	case IntentEndTurn:
		if w.Kind != PromptPriority || !g.inOpenActionPhase() || in.Player != g.Turn.Active {
			return refuse("R-TURN-06", "only the active player can end the turn, in the action phase, with an empty stack")
		}
	case IntentDiscard:
		if w.Kind != PromptDiscard {
			return refuse("R-TURN-11", "no discard is asked for")
		}
		if len(in.Objects) != w.Count {
			return refuse("R-TURN-11", "discard exactly %d loot cards", w.Count)
		}
		if !g.allInHand(in.Player, in.Objects) {
			return refuse("R-MECH-25", "you can only discard different loot cards from your own hand")
		}
	case IntentPlayLoot:
		if w.Kind != PromptPriority {
			return refuse("R-CARD-08", "loot can only be played with priority")
		}
		if len(in.Objects) != 1 || !contains(g.Players[in.Player].Hand, in.Objects[0]) {
			return refuse("R-CARD-08", "play one loot card from your hand")
		}
		if g.lootPlaysFor(in.Player) < 1 {
			return refuse("R-CARD-08", "no loot play available")
		}
		if ab, ok := g.lootAbility(in.Objects[0]); ok {
			if err := g.targetsAvailable(ab, in.Player, in.Objects[0]); err != nil {
				return err
			}
		}
	case IntentActivate:
		if w.Kind != PromptPriority {
			return refuse("R-ABIL-08", "abilities are used with priority")
		}
		if len(in.Objects) != 1 {
			return refuse("R-ABIL-08", "activate one object's ability")
		}
		if err := g.activatable(in.Player, in.Objects[0], in.Choice); err != nil {
			return err
		}
	case IntentAttack:
		if w.Kind != PromptPriority || !g.inOpenActionPhase() || in.Player != g.Turn.Active {
			return refuse("R-ATK-01", "only the active player attacks, in the action phase, with an empty stack")
		}
		if g.attacksLeft() < 1 {
			return refuse("R-TURN-07", "no attack left this turn")
		}
	case IntentPurchase:
		if w.Kind != PromptPriority || !g.inOpenActionPhase() || in.Player != g.Turn.Active {
			return refuse("R-SHOP-01", "only the active player purchases, in the action phase, with an empty stack")
		}
		if g.purchasesLeft() < 1 {
			return refuse("R-SHOP-05", "no purchase left this turn")
		}
	case IntentChoose:
		if w.Kind != PromptChoose || g.Choice == nil {
			return refuse("R-PRIO-01", "no choice is asked for")
		}
		if in.Choice < 0 || in.Choice >= len(w.Options) {
			return refuse(g.Choice.Rule, "choose one of the %d options", len(w.Options))
		}
	default:
		return refuse("R-PRIO-01", "unknown intent")
	}
	return nil
}

func (g *Game) allInHand(p PlayerID, ids []ObjectID) bool {
	seen := []ObjectID{}
	for _, id := range ids {
		if contains(seen, id) || !contains(g.Players[p].Hand, id) {
			return false
		}
		seen = append(seen, id)
	}
	return true
}

func (g *Game) apply(in Intent) {
	switch in.Kind {
	case IntentPass:
		g.emit(Event{Kind: EvPassed, Player: in.Player})
		g.pass()
	case IntentEndTurn:
		// Declaring the end passes priority (R-TURN-08): the others may respond.
		g.Turn.EndDeclared = true
		g.openWindowAfter(in.Player)
	case IntentDiscard:
		for _, id := range in.Objects {
			g.discardFromHand(in.Player, id)
		}
		g.Waiting = Prompt{}
	case IntentPlayLoot:
		if _, ok := g.lootAbility(in.Objects[0]); ok {
			d := g.def(in.Objects[0])
			g.startActivation(Activation{
				Player: in.Player, Source: in.Objects[0], Loot: true,
				Ability: AbilityRef{Card: d.Ref, Index: lootIndex(d)},
			})
		} else {
			g.playLoot(in.Player, in.Objects[0])
		}
	case IntentActivate:
		g.startActivation(Activation{
			Player: in.Player, Source: in.Objects[0],
			Ability: AbilityRef{Card: g.Object(in.Objects[0]).Card, Index: in.Choice},
		})
	case IntentChoose:
		g.answer(in.Choice)
	case IntentAttack:
		g.declareAttack(in.Player)
	case IntentPurchase:
		g.declarePurchase(in.Player)
	}
}

func (g *Game) emit(e Event) {
	g.events = append(g.events, e)
	g.collectTriggers(e)
}

// lootAbility returns the loot ability of a loot card, if it has one.
func (g *Game) lootAbility(id ObjectID) (Ability, bool) {
	d := g.def(id)
	if i := lootIndex(d); i >= 0 {
		return d.Abilities[i], true
	}
	return Ability{}, false
}

func lootIndex(d CardDef) int {
	for i, a := range d.Abilities {
		if a.Kind == LootAbility {
			return i
		}
	}
	return -1
}

func (g *Game) takeEvents() []Event {
	out := g.events
	g.events = nil
	return out
}

func contains(list []ObjectID, id ObjectID) bool {
	for _, x := range list {
		if x == id {
			return true
		}
	}
	return false
}

func remove(list []ObjectID, id ObjectID) []ObjectID {
	for i, x := range list {
		if x == id {
			return append(list[:i:i], list[i+1:]...)
		}
	}
	return list
}
