package rulesengine

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Situation is a shareable table position with one action to check
// (LR-04). With Expect filled in it is also a rules test (LR-07).
type Situation struct {
	Title  string          `json:"title"`
	Sets   []string        `json:"sets"`
	Seed   uint64          `json:"seed,omitempty"`
	Setup  SituationSetup  `json:"setup"`
	Action SituationAction `json:"action"`
	Expect *Expectation    `json:"expect,omitempty"`
}

// SituationSetup is the table. Turn: the active player is in the open
// action phase with one loot play, attack and purchase left.
type SituationSetup struct {
	Players  []SituationPlayer `json:"players"`
	Active   int               `json:"active"`
	Monsters []CardRef         `json:"monsters,omitempty"` // one per slot
	Shop     []CardRef         `json:"shop,omitempty"`     // one per slot
	Stack    []SituationStack  `json:"stack,omitempty"`    // bottom first
	// LootPlays overrides the active player's loot plays left.
	LootPlays *int `json:"loot_plays,omitempty"`
}

// SituationPlayer is one seat.
type SituationPlayer struct {
	Character CardRef   `json:"character"`
	Items     []CardRef `json:"items,omitempty"`
	Hand      []CardRef `json:"hand,omitempty"`
	Cents     int       `json:"cents,omitempty"`
	Souls     []CardRef `json:"souls,omitempty"`
	Damage    int       `json:"damage,omitempty"`
	Dead      bool      `json:"dead,omitempty"`
	// Deactivated lists items (or "character") that start deactivated.
	Deactivated []CardRef `json:"deactivated,omitempty"`
}

// SituationStack is one item already on the stack.
type SituationStack struct {
	DiceRoll int `json:"dice_roll"`
	Owner    int `json:"owner"`
}

// SituationAction is the one thing a player tries. Kind is "play",
// "activate", "attack", "purchase", "end_turn" or "pass". Choices answer
// the prompts that follow, by option label (e.g. a target's card name).
type SituationAction struct {
	Player  int      `json:"player"`
	Kind    string   `json:"kind"`
	Card    CardRef  `json:"card,omitempty"`
	Ability int      `json:"ability,omitempty"`
	Choices []string `json:"choices,omitempty"`
}

// Expectation is what should happen (LR-04).
type Expectation struct {
	Allowed bool             `json:"allowed"`
	Rule    string           `json:"rule,omitempty"` // the refusal's rule, when not allowed
	After   []ExpectedPlayer `json:"after,omitempty"`
}

// ExpectedPlayer checks one player after everything resolved.
type ExpectedPlayer struct {
	Player int   `json:"player"`
	Cents  *int  `json:"cents,omitempty"`
	Hand   *int  `json:"hand,omitempty"`
	Damage *int  `json:"damage,omitempty"`
	Souls  *int  `json:"souls,omitempty"`
	Dead   *bool `json:"dead,omitempty"`
}

// Answer is the engine's answer to a situation.
type Answer struct {
	Allowed bool    `json:"allowed"`
	Rule    string  `json:"rule,omitempty"`
	Reason  string  `json:"reason,omitempty"`
	Events  []Event `json:"events,omitempty"`
	// Matches is set when the situation has an expectation.
	Matches  *bool    `json:"matches,omitempty"`
	Mismatch []string `json:"mismatch,omitempty"`
}

// RunSituation builds the table, tries the action, resolves the stack
// when allowed, and compares with the expectation if there is one.
func RunSituation(data []byte, sets ...CardSet) (Answer, error) {
	var s Situation
	if err := json.Unmarshal(data, &s); err != nil {
		return Answer{}, fmt.Errorf("situation: %w", err)
	}
	g, err := s.build(sets)
	if err != nil {
		return Answer{}, err
	}
	in, err := s.intent(g)
	if err != nil {
		return Answer{}, err
	}
	ans := Answer{}
	ev, err := g.Apply(in)
	var re *RuleError
	switch {
	case errors.As(err, &re):
		ans.Rule, ans.Reason = re.Rule, re.Reason
	case err != nil:
		return Answer{}, err
	default:
		ans.Allowed = true
		ans.Events = ev
		more, err := g.settle(s.Action.Choices)
		if err != nil {
			return Answer{}, err
		}
		ans.Events = append(ans.Events, more...)
	}
	if s.Expect != nil {
		ans.compare(g, *s.Expect)
	}
	return ans, nil
}

// build creates the game for a situation.
func (s Situation) build(sets []CardSet) (*Game, error) {
	names := make([]string, len(sets))
	for i, set := range sets {
		names[i] = set.Name
	}
	for _, want := range s.Sets {
		found := false
		for _, n := range names {
			found = found || n == want
		}
		if !found {
			return nil, fmt.Errorf("situation needs card set %q", want)
		}
	}
	idx, err := newCardIndex(sets...)
	if err != nil {
		return nil, err
	}
	g := &Game{RNG: NewRNG(s.Seed), cards: idx, MaxHand: defaultHand, WinSouls: defaultSouls, Sets: names}
	need := func(ref CardRef) error {
		if _, ok := idx.find(ref); !ok {
			return fmt.Errorf("situation: unknown card %q", ref)
		}
		return nil
	}
	if len(s.Setup.Players) < 1 || s.Setup.Active < 0 || s.Setup.Active >= len(s.Setup.Players) {
		return nil, errors.New("situation: needs players and a valid active player")
	}
	for i, sp := range s.Setup.Players {
		p := PlayerID(i)
		if err := need(sp.Character); err != nil {
			return nil, err
		}
		ch := g.newObject(sp.Character, Zone{Kind: ZoneInPlay}, p)
		g.Object(ch).Role, g.Object(ch).Charged = RoleCharacter, !containsRef(sp.Deactivated, "character")
		pl := Player{ID: p, Character: ch, Cents: sp.Cents, Damage: sp.Damage, Dead: sp.Dead}
		for _, ref := range sp.Items {
			if err := need(ref); err != nil {
				return nil, err
			}
			id := g.newObject(ref, Zone{Kind: ZoneInPlay}, p)
			g.Object(id).Role, g.Object(id).Charged = RoleItem, !containsRef(sp.Deactivated, ref)
			pl.InPlay = append(pl.InPlay, id)
		}
		for _, ref := range sp.Souls {
			if err := need(ref); err != nil {
				return nil, err
			}
			id := g.newObject(ref, Zone{Kind: ZoneInPlay}, p)
			g.Object(id).Role = RoleSoul
			pl.InPlay = append(pl.InPlay, id)
		}
		for _, ref := range sp.Hand {
			if err := need(ref); err != nil {
				return nil, err
			}
			pl.Hand = append(pl.Hand, g.newObject(ref, Zone{Kind: ZoneHand}, p))
		}
		g.Players = append(g.Players, pl)
	}
	// Decks: every card of the sets, shuffled; enough for refills.
	for _, d := range idx.defs {
		if deck, ok := d.Kind.Deck(); ok && !d.Outside {
			for range max(d.Copies, 1) {
				g.AddToDeck(deck, d.Ref)
			}
		}
	}
	for d := range deckCount {
		g.ShuffleDeck(d)
	}
	for i, ref := range s.Setup.Monsters {
		if err := need(ref); err != nil {
			return nil, err
		}
		g.Monsters = append(g.Monsters, Slot{})
		g.putInSlot(g.newObject(ref, Zone{Kind: ZoneOutside}, NoPlayer), MonsterSlot, i)
	}
	for i, ref := range s.Setup.Shop {
		if err := need(ref); err != nil {
			return nil, err
		}
		g.Shop = append(g.Shop, Slot{})
		g.putInSlot(g.newObject(ref, Zone{Kind: ZoneOutside}, NoPlayer), ShopSlot, i)
	}
	active := PlayerID(s.Setup.Active)
	g.Turn = Turn{Active: active, Number: 1, Step: StepAction, Entered: true, LootPlays: 1, Attacks: 1, Purchases: 1}
	if s.Setup.LootPlays != nil {
		g.Turn.LootPlays = *s.Setup.LootPlays
	}
	holder := active
	for _, st := range s.Setup.Stack {
		owner := PlayerID(st.Owner)
		g.StackSeq++
		g.Stack = append(g.Stack, StackItem{ID: g.StackSeq, Kind: StackRoll, Controller: owner, Roll: st.DiceRoll, Label: "roll"})
		holder = owner
	}
	g.events = nil
	g.openWindow(holder)
	return g, nil
}

// intent turns the situation's action into an engine intent.
func (s Situation) intent(g *Game) (Intent, error) {
	a := s.Action
	p := PlayerID(a.Player)
	if a.Player < 0 || a.Player >= len(g.Players) {
		return Intent{}, errors.New("situation: no such player")
	}
	switch a.Kind {
	case "pass":
		return Intent{Player: p, Kind: IntentPass}, nil
	case "end_turn":
		return Intent{Player: p, Kind: IntentEndTurn}, nil
	case "attack":
		return Intent{Player: p, Kind: IntentAttack}, nil
	case "purchase":
		return Intent{Player: p, Kind: IntentPurchase}, nil
	case "play":
		for _, id := range g.Players[p].Hand {
			if g.Object(id).Card == a.Card {
				return Intent{Player: p, Kind: IntentPlayLoot, Objects: []ObjectID{id}}, nil
			}
		}
		return Intent{}, fmt.Errorf("situation: %q is not in player %d's hand", a.Card, a.Player)
	case "activate":
		for _, id := range g.controlled(p) {
			if g.Object(id).Card == a.Card {
				return Intent{Player: p, Kind: IntentActivate, Objects: []ObjectID{id}, Choice: a.Ability}, nil
			}
		}
		return Intent{}, fmt.Errorf("situation: player %d controls no %q", a.Player, a.Card)
	default:
		return Intent{}, fmt.Errorf("situation: unknown action kind %q", a.Kind)
	}
}

// settle answers choices by label and lets everyone pass until the stack
// is empty and the action is over.
func (g *Game) settle(choices []string) ([]Event, error) {
	var out []Event
	for range 1000 {
		w := g.Prompt()
		var in Intent
		switch {
		case g.Over:
			return out, nil
		case w.Kind == PromptChoose:
			if len(choices) == 0 {
				return out, fmt.Errorf("situation: a choice is needed, options %v", w.Options)
			}
			i := -1
			for j, o := range w.Options {
				if o == choices[0] {
					i = j
				}
			}
			if i < 0 {
				return out, fmt.Errorf("situation: %q is not an option of %v", choices[0], w.Options)
			}
			choices = choices[1:]
			in = Intent{Player: w.Player, Kind: IntentChoose, Choice: i}
		case w.Kind == PromptPriority && (len(g.Stack) > 0 || g.Attack.On || g.Purchase.On):
			in = Intent{Player: w.Player, Kind: IntentPass}
		default:
			return out, nil
		}
		ev, err := g.Apply(in)
		if err != nil {
			return out, err
		}
		out = append(out, ev...)
	}
	return out, errors.New("situation: the action never settled")
}

func (ans *Answer) compare(g *Game, e Expectation) {
	if ans.Allowed != e.Allowed {
		ans.Mismatch = append(ans.Mismatch, fmt.Sprintf("allowed is %v, expected %v", ans.Allowed, e.Allowed))
	}
	if !e.Allowed && e.Rule != "" && ans.Rule != e.Rule {
		ans.Mismatch = append(ans.Mismatch, fmt.Sprintf("refused by %s, expected %s", ans.Rule, e.Rule))
	}
	for _, x := range e.After {
		if x.Player < 0 || x.Player >= len(g.Players) {
			ans.Mismatch = append(ans.Mismatch, fmt.Sprintf("no player %d", x.Player))
			continue
		}
		pl := g.Players[x.Player]
		checkInt := func(name string, want *int, got int) {
			if want != nil && *want != got {
				ans.Mismatch = append(ans.Mismatch, fmt.Sprintf("player %d %s is %d, expected %d", x.Player, name, got, *want))
			}
		}
		checkInt("cents", x.Cents, pl.Cents)
		checkInt("hand", x.Hand, len(pl.Hand))
		checkInt("damage", x.Damage, pl.Damage)
		checkInt("souls", x.Souls, g.SoulValue(pl.ID))
		if x.Dead != nil && *x.Dead != pl.Dead {
			ans.Mismatch = append(ans.Mismatch, fmt.Sprintf("player %d dead is %v, expected %v", x.Player, pl.Dead, *x.Dead))
		}
	}
	ok := len(ans.Mismatch) == 0
	ans.Matches = &ok
}

func containsRef(list []CardRef, r CardRef) bool {
	for _, x := range list {
		if x == r {
			return true
		}
	}
	return false
}
