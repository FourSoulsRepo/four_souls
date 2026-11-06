package rulesengine

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
)

// SituationStack is one item already on the stack (12.5). Kind is empty
// for a dice roll (the first form, `{"dice_roll": 4, "owner": 1}`),
// "play" for a loot card, "activate" for an activated ability,
// "trigger" for a triggered ability and "attack_roll" for an attack roll.
// Card and Ability pick the card as in an action; Choices are the mode
// and targets chosen when it went on the stack, by option label.
//
// An item is as it is once played: costs are paid and the setup shows
// what they left (a stacked loot is not in hand, a tapped item is in
// "deactivated"). Nothing is paid again and nothing triggers on its
// way onto the stack.
type SituationStack struct {
	Kind     string   `json:"kind,omitempty"`
	DiceRoll int      `json:"dice_roll"`
	Owner    int      `json:"owner"`
	Card     CardRef  `json:"card,omitempty"`
	Ability  int      `json:"ability,omitempty"`
	Choices  []string `json:"choices,omitempty"`
}

// UnmarshalJSON rejects unknown fields, so a misspelt item is an error,
// not a dice roll of 0.
func (s *SituationStack) UnmarshalJSON(data []byte) error {
	type plain SituationStack
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	var p plain
	if err := dec.Decode(&p); err != nil {
		return fmt.Errorf("stack item: %w", err)
	}
	*s = SituationStack(p)
	return nil
}

// stackUp puts the setup's items on the stack, bottom first, and returns
// who gets priority: the top item's owner, or holder for an empty stack.
func (g *Game) stackUp(items []SituationStack, holder PlayerID) (PlayerID, error) {
	for i, st := range items {
		it, err := g.stackedItem(st)
		if err != nil {
			return 0, fmt.Errorf("situation: stack item %d: %w", i+1, err)
		}
		g.StackSeq++
		it.ID = g.StackSeq
		g.Stack = append(g.Stack, it)
		holder = it.Controller
	}
	return holder, nil
}

// stackedItem builds one item with the constructors normal play uses.
func (g *Game) stackedItem(st SituationStack) (StackItem, error) {
	if st.Owner < 0 || st.Owner >= len(g.Players) {
		return StackItem{}, fmt.Errorf("owner %d is not a seat", st.Owner)
	}
	p := PlayerID(st.Owner)
	rolls := st.Kind == "" || st.Kind == "attack_roll"
	switch {
	case rolls && (st.DiceRoll < 1 || st.DiceRoll > 6):
		return StackItem{}, fmt.Errorf("dice_roll %d is not 1 to 6", st.DiceRoll)
	case !rolls && st.DiceRoll != 0:
		return StackItem{}, fmt.Errorf("kind %q has no dice_roll", st.Kind)
	case st.Kind != "activate" && st.Kind != "trigger" && st.Ability != 0:
		return StackItem{}, fmt.Errorf("kind %q has no ability", st.Kind)
	case st.Kind == "" && (st.Card != "" || len(st.Choices) > 0):
		return StackItem{}, errors.New(`a dice roll has no card or choices; set "kind"`)
	case (st.Kind == "attack_roll" || st.Kind == "trigger") && len(st.Choices) > 0:
		return StackItem{}, fmt.Errorf("kind %q has no choices", st.Kind)
	}
	if st.Kind != "" {
		if _, ok := g.cards.find(st.Card); !ok {
			return StackItem{}, fmt.Errorf("unknown card %q", st.Card)
		}
	}
	switch st.Kind {
	case "":
		return StackItem{Kind: StackRoll, Controller: p, Roll: st.DiceRoll, Label: eventRoll.Msg()}, nil
	case "play":
		return g.stackedLoot(p, st)
	case "activate":
		return g.stackedAbility(p, st)
	case "trigger":
		return g.stackedTrigger(p, st)
	case "attack_roll":
		return g.stackedAttackRoll(p, st)
	default:
		return StackItem{}, fmt.Errorf("unknown kind %q (play, activate, trigger, attack_roll, or none for a dice roll)", st.Kind)
	}
}

// stackedLoot: the loot card is on the stack, not in any hand.
func (g *Game) stackedLoot(p PlayerID, st SituationStack) (StackItem, error) {
	d, _ := g.cards.find(st.Card)
	if d.Kind != LootCard {
		return StackItem{}, fmt.Errorf("%q is not a loot card", st.Card)
	}
	if d.CopiesTapOnPlay {
		return StackItem{}, fmt.Errorf("%q copies an ability when played; it can't start on the stack", st.Card)
	}
	src := g.newObject(st.Card, Zone{Kind: ZoneStack}, p)
	a := Activation{Player: p, Source: src, Loot: true}
	if i := lootIndex(d); i >= 0 {
		a.Ability = AbilityRef{Card: d.Ref, Index: i}
		var err error
		if a.Mode, a.Chosen, err = g.stackedChoices(a.Ability, p, src, st.Choices); err != nil {
			return StackItem{}, err
		}
	} else if len(st.Choices) > 0 {
		return StackItem{}, fmt.Errorf("%q has nothing to choose", st.Card)
	}
	return g.activationItem(a, src), nil
}

// stackedAbility: an activated ability of a card the owner controls.
func (g *Game) stackedAbility(p PlayerID, st SituationStack) (StackItem, error) {
	i := slices.IndexFunc(g.controlled(p), func(id ObjectID) bool { return g.Object(id).Card == st.Card })
	if i < 0 {
		return StackItem{}, fmt.Errorf("owner %d controls no %q", st.Owner, st.Card)
	}
	src := g.controlled(p)[i]
	refs := g.AbilitiesOf(src)
	if st.Ability < 0 || st.Ability >= len(refs) || g.ability(refs[st.Ability]).Kind != Activated {
		return StackItem{}, fmt.Errorf("%q has no activated ability %d", st.Card, st.Ability)
	}
	a := Activation{Player: p, Source: src, Ability: refs[st.Ability]}
	var err error
	if a.Mode, a.Chosen, err = g.stackedChoices(a.Ability, p, src, st.Choices); err != nil {
		return StackItem{}, err
	}
	return g.activationItem(a, src), nil
}

// stackedTrigger: a triggered ability of a card the owner controls, or
// of a monster or room in a slot, whose "you" is the active player. Its
// event's player is the owner.
func (g *Game) stackedTrigger(p PlayerID, st SituationStack) (StackItem, error) {
	cards := g.controlled(p)
	for _, rows := range [][]Slot{g.Monsters, g.Rooms} {
		for _, s := range rows {
			if top, ok := s.TopOf(); ok {
				cards = append(cards, top)
			}
		}
	}
	i := slices.IndexFunc(cards, func(id ObjectID) bool { return g.Object(id).Card == st.Card })
	if i < 0 || g.abilityController(cards[i]) != p {
		return StackItem{}, fmt.Errorf("owner %d controls no %q", st.Owner, st.Card)
	}
	src := cards[i]
	refs := g.AbilitiesOf(src)
	if st.Ability < 0 || st.Ability >= len(refs) || g.ability(refs[st.Ability]).Kind != Triggered {
		return StackItem{}, fmt.Errorf("%q has no triggered ability %d", st.Card, st.Ability)
	}
	return triggerItem(PendingTrigger{Ability: refs[st.Ability], Source: src, Controller: p, EventPlayer: p}), nil
}

// stackedAttackRoll: the active player attacks the monster card; the
// attack is declared and its target chosen, which used one attack.
func (g *Game) stackedAttackRoll(p PlayerID, st SituationStack) (StackItem, error) {
	if p != g.Turn.Active {
		return StackItem{}, fmt.Errorf("owner %d is not the active player; only they attack here", st.Owner)
	}
	target, found := ObjectID(0), false
	for _, s := range g.Monsters {
		if top, ok := s.TopOf(); ok && !found && g.Object(top).Card == st.Card && g.Object(top).Role == RoleMonster {
			target, found = top, true
		}
	}
	switch {
	case !found:
		return StackItem{}, fmt.Errorf("%q is not a monster in a monster slot", st.Card)
	case g.Attack.Started && g.Attack.Target != target:
		return StackItem{}, errors.New("every attack roll must attack the same monster")
	case !g.Attack.Started:
		g.Attack = AttackState{On: true, Started: true, Target: target}
		if g.Turn.Attacks > 0 {
			g.Turn.Attacks--
		}
	}
	it := g.attackRollItem(p)
	it.Roll = st.DiceRoll
	return it, nil
}

// stackedChoices reads the mode and targets of an ability by option
// label, as the prompts of a normal play would offer them.
func (g *Game) stackedChoices(ref AbilityRef, p PlayerID, src ObjectID, labels []string) (int, []Chosen, error) {
	pick := func(what string, msgs []Msg) (int, error) {
		plains := make([]string, len(msgs))
		for i, m := range msgs {
			plains[i] = g.Plain(m)
		}
		if len(labels) == 0 {
			return 0, fmt.Errorf("needs a %s, options %v", what, plains)
		}
		i := slices.Index(plains, labels[0])
		if i < 0 {
			return 0, fmt.Errorf("%q is not a %s of %v", labels[0], what, plains)
		}
		labels = labels[1:]
		return i, nil
	}
	ab := g.ability(ref)
	mode := 0
	if len(ab.Modes) > 0 {
		modes := g.availableModes(ab, p, src)
		msgs := make([]Msg, len(modes))
		for i, m := range modes {
			msgs[i] = AskWords.Msg(Str("text", ab.Modes[m].Text))
		}
		i, err := pick("mode", msgs)
		if err != nil {
			return 0, nil, err
		}
		mode = modes[i]
	}
	var chosen []Chosen
	for _, spec := range g.targetSpecs(ref, mode) {
		opts, msgs := g.targetOptions(spec, p, src)
		i, err := pick("target", msgs)
		if err != nil {
			return 0, nil, err
		}
		chosen = append(chosen, opts[i])
	}
	if len(labels) > 0 {
		return 0, nil, fmt.Errorf("choices %v left over", labels)
	}
	return mode, chosen, nil
}
