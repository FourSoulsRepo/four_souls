package rulesengine

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"slices"
)

// SituationCounter is a number of counters on one card (24.3). Card is
// the card's ref, or "character" for a player's character; Name names
// the counters ("gold", a level), empty for plain counters.
type SituationCounter struct {
	Card  CardRef `json:"card"`
	Name  string  `json:"name,omitempty"`
	Count int     `json:"count"`
}

// UnmarshalJSON rejects unknown fields, so a misspelt field is an error,
// not a counter of 0.
func (c *SituationCounter) UnmarshalJSON(data []byte) error {
	type plain SituationCounter
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	var p plain
	if err := dec.Decode(&p); err != nil {
		return fmt.Errorf("counter (fields card, name, count): %w", err)
	}
	*c = SituationCounter(p)
	return nil
}

// setUpPlayer applies a seat's counters, eternal cards and extra loot
// plays once its cards are built; seat is its index, as in errors.
func (g *Game) setUpPlayer(pl *Player, seat int, sp SituationPlayer) error {
	if sp.ExtraLootPlays < 0 {
		return fmt.Errorf("situation: player %d extra_loot_plays %d is below 0", seat, sp.ExtraLootPlays)
	}
	pl.ExtraLootPlays = sp.ExtraLootPlays
	// refs are the seat's cards in play by ref; "character" first.
	refs := []CardRef{"character"}
	ids := []ObjectID{pl.Character}
	for _, id := range pl.InPlay {
		refs, ids = append(refs, g.Object(id).Card), append(ids, id)
	}
	find := func(field string, ref CardRef, from int) (ObjectID, error) {
		if i := slices.Index(refs[from:], ref); i >= 0 {
			return ids[from+i], nil
		}
		return 0, fmt.Errorf("situation: player %d %s: %q is not in play, options %v", seat, field, ref, refs[from:])
	}
	for _, ref := range sp.Eternal {
		id, err := find("eternal", ref, 1) // an item or soul, not the character
		if err != nil {
			return err
		}
		g.Object(id).Eternal = true
	}
	for _, c := range sp.Counters {
		id, err := find("counters", c.Card, 0)
		if err != nil {
			return err
		}
		if err := putCounters(g.Object(id), c); err != nil {
			return fmt.Errorf("situation: player %d counters: %w", seat, err)
		}
	}
	return nil
}

func putCounters(o *Object, c SituationCounter) error {
	if c.Count < 1 {
		return fmt.Errorf("%q count %d is below 1", c.Card, c.Count)
	}
	o.addCounters(c.Name, c.Count)
	return nil
}

// setUpDecks puts the setup's cards on top of their decks, top first.
func (g *Game) setUpDecks(decks map[string][]CardRef) error {
	for name := range decks {
		if !slices.Contains(deckNames[:], name) {
			return fmt.Errorf("situation: decks: %q is no deck, options %v", name, deckNames)
		}
	}
	for d := range deckCount {
		refs := decks[d.String()]
		for _, ref := range refs {
			def, ok := g.cards.find(ref)
			if !ok {
				return fmt.Errorf("situation: decks %s: unknown card %q", d, ref)
			}
			if k, ok := def.Kind.Deck(); !ok || k != d {
				return fmt.Errorf("situation: decks %s: %q is no card of the %s deck", d, ref, d)
			}
		}
		for i := len(refs) - 1; i >= 0; i-- {
			g.PutOnDeck(d, refs[i])
		}
	}
	return nil
}

// setUpSlots puts damage and counters on cards in slots and on bonus
// souls. A ref names the first such card, in slot order.
func (g *Game) setUpSlots(s SituationSetup) error {
	refs, ids, slotted := g.slotCards() // bonus souls take counters, not damage
	damaged := make([]CardRef, 0, len(s.MonsterDamage))
	for ref := range s.MonsterDamage {
		damaged = append(damaged, ref)
	}
	slices.Sort(damaged)
	for _, ref := range damaged {
		i := slices.Index(refs[:slotted], ref)
		if i < 0 {
			return fmt.Errorf("situation: monster_damage: %q is in no slot, options %v", ref, refs[:slotted])
		}
		n, hp := s.MonsterDamage[ref], g.objectMaxHP(ids[i])
		if n < 0 || n >= hp {
			return fmt.Errorf("situation: monster_damage: %q has %d HP; damage %d is not 0 to %d", ref, hp, n, hp-1)
		}
		g.Object(ids[i]).Damage = n
	}
	for _, c := range s.Counters {
		i := slices.Index(refs, c.Card)
		if i < 0 {
			return fmt.Errorf("situation: counters: %q is in no slot and no bonus soul, options %v", c.Card, refs)
		}
		if err := putCounters(g.Object(ids[i]), c); err != nil {
			return fmt.Errorf("situation: counters: %w", err)
		}
	}
	return nil
}

// setupPhases are the values of setup.phase and the step each starts
// in, its priority window open (24.10). The step's own work is done:
// start is after the recharge and the start-of-turn triggers
// (R-TURN-03), loot after the active player looted (R-TURN-04), end
// after the end-of-turn triggers (R-TURN-10). When everyone passes, the
// game moves on as in a game (R-STACK-07).
var setupPhases = []struct {
	name string
	step Step
}{{"start", StepStartTriggers}, {"loot", StepLootWindow}, {"action", StepAction}, {"end", StepEndTriggers}}

// phaseStep is the step setup.phase names; none is the action phase.
func (s SituationSetup) phaseStep() (Step, error) {
	if s.Phase == "" {
		return StepAction, nil
	}
	names := make([]string, len(setupPhases))
	for i, p := range setupPhases {
		if p.name == s.Phase {
			return p.step, nil
		}
		names[i] = p.name
	}
	return 0, fmt.Errorf("situation: phase %q is no phase, options %v", s.Phase, names)
}

// setUpTurn sets the turn's step and applies its overrides. loot_plays
// is unchecked, as before 24.3, so old payloads keep their answers.
func (s SituationSetup) setUpTurn(t *Turn) error {
	step, err := s.phaseStep()
	if err != nil {
		return err
	}
	t.Step = step
	switch step { //nolint:exhaustive // the other steps are no setup phase
	case StepStartTriggers, StepLootWindow:
		// The action phase gives the loot play, attack and purchase
		// (R-TURN-05, R-TURN-07); before it there are none.
		if s.LootPlays != nil || s.Attacks != nil || s.Purchases != nil {
			return fmt.Errorf("situation: phase %s: loot_plays, attacks and purchases come with the action phase (R-TURN-05, R-TURN-07); leave them out", s.Phase)
		}
		t.LootPlays, t.Attacks, t.Purchases = 0, 0, 0
		return nil
	case StepEndTriggers:
		t.EndDeclared = true // the action phase ended (R-TURN-09)
	}
	if s.LootPlays != nil {
		t.LootPlays = *s.LootPlays
	}
	for _, o := range []struct {
		name string
		n    *int
		to   *int
	}{{"attacks", s.Attacks, &t.Attacks}, {"purchases", s.Purchases, &t.Purchases}} {
		if o.n == nil {
			continue
		}
		if *o.n < 0 {
			return fmt.Errorf("situation: %s %d is below 0", o.name, *o.n)
		}
		*o.to = *o.n
	}
	return nil
}

// setUpRolls forces the next dice results.
func (g *Game) setUpRolls(rolls []int) error {
	for i, r := range rolls {
		if r < 1 || r > 6 {
			return fmt.Errorf("situation: rolls: roll %d is %d, not 1 to 6", i+1, r)
		}
	}
	g.ForceRolls(rolls...)
	return nil
}

// setUpDeclared declares the setup's attack or purchase, as the active
// player does in the action phase with an empty stack (R-ATK-01,
// R-SHOP-01): it uses one up, and priority passes before the target or
// item is chosen (R-ATK-02, R-SHOP-02). Nothing triggers and nothing is
// paid; the setup shows what the declaration left, as for the stack.
// setup.stack holds the responses to it: they resolve first, and once
// everyone passes on the empty stack the active player chooses, as in a
// game (R-STACK-07). An attack roll is no such response.
func (g *Game) setUpDeclared(s SituationSetup) error {
	var rule, limit, counts string
	var left func() int
	switch s.Declared {
	case "":
		return nil
	case "attack":
		rule, limit, counts, left = "R-ATK-01", "R-TURN-07", "attacks", g.attacksLeft
	case "purchase":
		rule, limit, counts, left = "R-SHOP-01", "R-SHOP-05", "purchases", g.purchasesLeft
	default:
		return fmt.Errorf("situation: declared %q is no declaration, options [attack purchase]", s.Declared)
	}
	switch {
	case g.Turn.Step != StepAction:
		return fmt.Errorf("situation: declared %s needs phase action, not %s: it is declared only in the action phase (%s)", s.Declared, s.Phase, rule)
	case left() < 1:
		return fmt.Errorf("situation: declared %s needs one left (%s), and %s is 0", s.Declared, limit, counts)
	}
	if i := slices.IndexFunc(s.Stack, func(st SituationStack) bool { return st.Kind == "attack_roll" }); i >= 0 {
		if s.Declared == "attack" {
			return fmt.Errorf("situation: declared attack: stack item %d is an attack roll, but the target is chosen only after the declaration's window (R-ATK-02), and attack rolls go against the target (R-ATK-11)", i+1)
		}
		return fmt.Errorf("situation: declared purchase: stack item %d is an attack roll, but a purchase is declared only with an empty stack, so not during an attack (R-SHOP-01)", i+1)
	}
	if s.Declared == "attack" {
		g.spendAttack()
		g.Attack = AttackState{On: true}
	} else {
		g.spendPurchase()
		g.Purchase = PurchaseState{On: true}
	}
	return nil
}

// setUpDuties sets the attack duties in force this turn (19.235), as
// the turn holds them once the setup's window is open, so after a
// declared attack or a stacked attack roll has counted (R-TURN-07). A card in a slot sets no
// duty (owner, 2026-11-23): a player may have stopped it, so the setup
// names the duties that hold. They hold from the start phase, where the
// start-of-turn triggers have set them (R-TURN-03); in the end phase
// the turn's end is declared, which R-ATK-03 refuses while one is owed.
func (g *Game) setUpDuties(s SituationSetup) error {
	for _, n := range []struct {
		name string
		n    int
	}{{"must_attacks", s.MustAttacks}, {"must_attack_deck", s.MustAttackDeck}} {
		if n.n < 0 {
			return fmt.Errorf("situation: %s %d is below 0", n.name, n.n)
		}
	}
	if s.MustAttack == "" && s.MustAttacks == 0 && s.MustAttackDeck == 0 && len(s.MustAttackOneOf) == 0 {
		return nil
	}
	if g.Turn.Step == StepEndTriggers {
		return errors.New("situation: phase end: the turn's end is declared, which R-ATK-03 refuses while an attack is owed; " +
			"leave out must_attack, must_attacks, must_attack_deck and must_attack_one_of")
	}
	targets := g.attackTargets()
	refs := make([]CardRef, len(targets))
	for i, id := range targets {
		refs[i] = g.Object(id).Card
	}
	// target finds the first card named ref the active player may
	// attack now.
	target := func(field string, ref CardRef) (ObjectID, error) {
		for i, r := range refs {
			if r == ref {
				if g.Attack.Started && g.Attack.Target == targets[i] {
					return 0, fmt.Errorf("situation: %s: the attack roll on the stack attacks %q, so the duty is met (R-ATK-22); leave it out", field, ref)
				}
				return targets[i], nil
			}
		}
		return 0, fmt.Errorf("situation: %s: the active player can't attack %q, options %v", field, ref, refs)
	}
	if s.MustAttack != "" {
		id, err := target("must_attack", s.MustAttack)
		if err != nil {
			return err
		}
		g.Turn.MustAttack = id
	}
	// must_attack_one_of names cards, not objects (19.241); a setup
	// builds no covered card, so each is one the active player can
	// attack now.
	for _, ref := range s.MustAttackOneOf {
		if slices.Contains(g.Turn.MustAttackOneOfCards, ref) {
			return fmt.Errorf("situation: must_attack_one_of: %q is named twice", ref)
		}
		if _, err := target("must_attack_one_of", ref); err != nil {
			return err
		}
		g.Turn.MustAttackOneOfCards = append(g.Turn.MustAttackOneOfCards, ref)
	}
	g.Turn.MustAttacks, g.Turn.MustAttackDeck = s.MustAttacks, s.MustAttackDeck
	return nil
}

// openSetupWindow opens the first window: for holder, the top item's owner
// or the active player, unless the setup names the seat (24.4). With an
// empty stack, the players from the active player up to that seat have
// passed (R-PRIO-02), so when it and the rest pass the window closes as
// in a game (R-PRIO-05). In the open action phase nobody passes on an
// empty stack: the active player can't (R-TURN-09), so the seat must be
// theirs (24.10).
func (g *Game) openSetupWindow(seat *int, holder PlayerID) error {
	if seat == nil {
		g.openWindow(holder)
		return nil
	}
	n := len(g.Players)
	if *seat < 0 || *seat >= n {
		return fmt.Errorf("situation: priority %d is not a seat, options 0 to %d", *seat, n-1)
	}
	if PlayerID(*seat) != g.Turn.Active && g.inOpenActionPhase() {
		return fmt.Errorf("situation: priority %d: in the open action phase with an empty stack only the active player %d holds priority, and they can't pass (R-TURN-09); "+
			"for another seat use phase start, loot or end, declared attack or purchase, or a stack", *seat, g.Turn.Active)
	}
	g.openWindow(PlayerID(*seat))
	if len(g.Stack) == 0 {
		g.Priority.Passes = (*seat - int(g.Turn.Active) + n) % n
	}
	return nil
}
