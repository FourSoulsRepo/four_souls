package rulesengine

import (
	"encoding/json"
	"errors"
	"fmt"
	"slices"
	"strings"
)

// Situation is a shareable table position with one action to check
// (LR-04), or several steps in a row (24.4). With Expect filled in it is
// also a rules test (LR-07).
type Situation struct {
	Title  string           `json:"title"`
	Sets   []string         `json:"sets"`
	Seed   uint64           `json:"seed,omitempty"`
	Setup  SituationSetup   `json:"setup"`
	Action *SituationAction `json:"action,omitempty"`
	// Steps are actions in a row, each checked by its own Expect. A
	// situation has Action (one step, as before 24.4) or Steps.
	Steps  []SituationStep `json:"steps,omitempty"`
	Expect *Expectation    `json:"expect,omitempty"`
	// Source names the basis of Expect: a card, a rule, a question or a
	// link. A test file needs it; a payload needs none (ADR 034). Only
	// the test runner reads it (enginetest.RunSituations).
	Source []string `json:"source,omitempty"`
}

// SituationSetup is the table. Turn: the active player is in the open
// action phase with one loot play, attack and purchase left.
// LootPlays, Attacks and Purchases override those (24.3); Phase and
// Declared start it in another priority window (24.10).
type SituationSetup struct {
	Players  []SituationPlayer `json:"players"`
	Active   int               `json:"active"`
	Monsters []CardRef         `json:"monsters,omitempty"` // one per slot
	Shop     []CardRef         `json:"shop,omitempty"`     // one per slot
	Stack    []SituationStack  `json:"stack,omitempty"`    // bottom first (12.5)
	// LootPlays overrides the active player's loot plays left.
	LootPlays *int `json:"loot_plays,omitempty"`
	// BonusSouls are the active bonus souls (R-SETUP-06).
	BonusSouls []CardRef `json:"bonus_souls,omitempty"`
	// Rooms fill the room slots, one each (R-SETUP-12).
	Rooms []CardRef `json:"rooms,omitempty"`
	// Bosses and Minions fill final boss and minion slots, one each, as
	// in a challenge (R-CHAL-02, R-CHAL-07).
	Bosses  []CardRef `json:"bosses,omitempty"`
	Minions []CardRef `json:"minions,omitempty"`
	// Attacks and Purchases override the active player's attacks and
	// purchases left; bonus ones (Champion Belt) come on top.
	Attacks   *int `json:"attacks,omitempty"`
	Purchases *int `json:"purchases,omitempty"`
	// Rolls are the next dice results, in order (as Game.ForceRolls).
	Rolls []int `json:"rolls,omitempty"`
	// Decks puts cards on top of decks by name ("loot", "treasure",
	// "monster", "room"), top first, above the shuffled set cards.
	Decks map[string][]CardRef `json:"decks,omitempty"`
	// MonsterDamage is damage on a card in a monster, final boss, minion
	// or room slot, by its ref.
	MonsterDamage map[CardRef]int `json:"monster_damage,omitempty"`
	// Counters go on cards no player controls: a card in a slot or an
	// active bonus soul.
	Counters []SituationCounter `json:"counters,omitempty"`
	// Priority is the seat that holds priority. With an empty stack the
	// players before it in turn order, from the active player, have
	// passed (R-PRIO-02); with a stack it overrides the top item's owner.
	// In the open action phase only the active player holds it there
	// (R-TURN-09).
	Priority *int `json:"priority,omitempty"`
	// Phase is the step whose priority window the situation starts in
	// (24.10): "start" (after the start-of-turn triggers, R-TURN-03),
	// "loot" (after the loot step's loot, R-TURN-04), "action" (the
	// default) or "end" (after the end-of-turn triggers, R-TURN-10).
	Phase string `json:"phase,omitempty"`
	// Declared is "attack" or "purchase": the active player has declared
	// it in the action phase and not chosen the target or item yet
	// (R-ATK-02, R-SHOP-02). Stack holds the responses to it; they
	// resolve before the choice.
	Declared string `json:"declared,omitempty"`
	// MustAttack, MustAttacks and MustAttackDeck are the attack duties
	// in force this turn, as the Turn fields of those names hold them
	// (19.235). A card in a slot does not set its duty: a player may
	// have stopped it, so the setup says which hold. MustAttack names
	// what the active player must attack; MustAttacks is how many more
	// attacks they must make; MustAttackDeck how many of their next
	// attacks must go to the monster deck.
	MustAttack     CardRef `json:"must_attack,omitempty"`
	MustAttacks    int     `json:"must_attacks,omitempty"`
	MustAttackDeck int     `json:"must_attack_deck,omitempty"`
	// MustAttackOneOf names the cards the active player must make an
	// additional attack on one of (Boss Rush!, Friendly Ball, Lust for
	// Blood!; 19.241). It adds no attack: `attacks` says how many.
	MustAttackOneOf []CardRef `json:"must_attack_one_of,omitempty"`
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
	// Counters go on the player's character ("character"), items or souls.
	Counters []SituationCounter `json:"counters,omitempty"`
	// Eternal lists items or souls that are eternal, as starting items.
	Eternal []CardRef `json:"eternal,omitempty"`
	// ExtraLootPlays are loot plays from abilities, e.g. Cain's ↷ on
	// another player's turn (Player.ExtraLootPlays).
	ExtraLootPlays int `json:"extra_loot_plays,omitempty"`
}

// SituationAction is the one thing a player tries. Kind is "play",
// "activate", "attack", "purchase", "end_turn" or "pass"; "throw" throws
// the roll that waits, "choose" answers the open question, and "reveal"
// puts Card on the monster deck and attacks the deck (24.4). Choices
// answer the prompts that follow, by option label (e.g. a target's card
// name).
type SituationAction struct {
	Player  int      `json:"player"`
	Kind    string   `json:"kind"`
	Card    CardRef  `json:"card,omitempty"`
	Ability int      `json:"ability,omitempty"`
	Choices []string `json:"choices,omitempty"`
}

// SituationStep is one step of Situation.Steps: an action, and what
// should be true once it settled.
type SituationStep struct {
	SituationAction
	Expect *Expectation `json:"expect,omitempty"`
}

// Expectation is what should happen (LR-04). A field left out is not
// checked. A list of cards names every card there, in any order; an
// empty list expects none (24.5). A list or map of checks that is there
// but empty is a mismatch: it would check nothing (19.234).
type Expectation struct {
	Allowed bool             `json:"allowed"`
	Rule    string           `json:"rule,omitempty"` // the refusal's rule, when not allowed
	After   []ExpectedPlayer `json:"after,omitempty"`
	// Monsters check cards in slots.
	Monsters []ExpectedMonster `json:"monsters,omitempty"`
	// Counters on cards no player controls, as setup.counters; a count
	// of 0 expects none.
	Counters []SituationCounter `json:"counters,omitempty"`
	// Decks are the cards on top of decks by name, top first, as
	// setup.decks, [] an empty deck; Discards the whole discard piles,
	// any order.
	Decks    map[string][]CardRef `json:"decks,omitempty"`
	Discards map[string][]CardRef `json:"discards,omitempty"`
	// DeckBottoms are the cards at the bottom of decks by name, in any
	// order: several cards go to the bottom in random order (Q-22);
	// [] expects an empty deck.
	DeckBottoms map[string][]CardRef `json:"deck_bottoms,omitempty"`
	// Active checks the turn's player; LootPlays and Attacks what that
	// player can still use, bonus ones too, unlike setup (24.28).
	Active    *int `json:"active,omitempty"`
	LootPlays *int `json:"loot_plays,omitempty"`
	Attacks   *int `json:"attacks,omitempty"`
	// Over and Winners check the game's end; Winners are seats.
	Over    *bool `json:"over,omitempty"`
	Winners []int `json:"winners,omitempty"`
	// Prompt checks the open question: who is asked and its options.
	// It also makes settling stop at that question (24.6).
	Prompt *ExpectedPrompt `json:"prompt,omitempty"`
	// Resolved names stack items in the order they resolved: a card
	// ref or an item kind (dice_roll, attack_roll, …). Other items may
	// resolve between them; an empty list expects none (24.6).
	Resolved []string `json:"resolved,omitempty"`
	// Events are event kinds that must happen, or with None must not.
	Events []ExpectedEvent `json:"events,omitempty"`
	// ResolvedRolls are the values every roll resolved as, in order:
	// the final result after changes (R-DICE-06); [] expects none (24.11).
	ResolvedRolls []int `json:"resolved_rolls,omitempty"`
}

// ExpectedPlayer checks one player after everything resolved. Hand and
// Souls are counts (soul value); HandCards and SoulCards the cards.
type ExpectedPlayer struct {
	Player int   `json:"player"`
	Cents  *int  `json:"cents,omitempty"`
	Hand   *int  `json:"hand,omitempty"`
	Damage *int  `json:"damage,omitempty"`
	Souls  *int  `json:"souls,omitempty"`
	Dead   *bool `json:"dead,omitempty"`
	// HP is what is left: max HP minus damage (Game.PlayerHP).
	HP  *int `json:"hp,omitempty"`
	ATK *int `json:"atk,omitempty"`
	// Items, HandCards, SoulCards: every such card, any order.
	Items     []CardRef `json:"items,omitempty"`
	HandCards []CardRef `json:"hand_cards,omitempty"`
	SoulCards []CardRef `json:"soul_cards,omitempty"`
	// Deactivated: every deactivated item, and "character", as setup.
	Deactivated []CardRef `json:"deactivated,omitempty"`
	// Counters on the player's cards, as setup; a count of 0 expects none.
	Counters []SituationCounter `json:"counters,omitempty"`
}

// Answer is the engine's answer to a situation.
type Answer struct {
	Allowed bool   `json:"allowed"`
	Rule    string `json:"rule,omitempty"`
	Reason  string `json:"reason,omitempty"` // English, for developers
	// Say is the refusal for the player: a message key with values
	// (ADR 021), so each client shows it in its language (ADR 026).
	Say    Msg     `json:"say,omitzero"`
	Events []Event `json:"events,omitempty"`
	// Matches is set when the situation has an expectation.
	Matches  *bool    `json:"matches,omitempty"`
	Mismatch []string `json:"mismatch,omitempty"`
}

// RunSituation builds the table, runs the action or each step, resolves
// the stack after each allowed one, and compares with the expectations.
// With steps, the answer's refusal is the first refused step's.
func RunSituation(data []byte, sets ...CardSet) (Answer, error) {
	var s Situation
	if err := json.Unmarshal(data, &s); err != nil {
		return Answer{}, fmt.Errorf("situation: %w", err)
	}
	steps, err := s.stepList()
	if err != nil {
		return Answer{}, err
	}
	g, err := s.build(sets)
	if err != nil {
		return Answer{}, err
	}
	multi := s.Action == nil
	ans := Answer{Allowed: true}
	refused := ""
	var mismatch []string
	checked := s.Expect != nil
	for i, st := range steps {
		var stop stopAt
		if i+1 < len(steps) {
			stop = stopAt{question: steps[i+1].Kind == "choose", throw: steps[i+1].Kind == "throw"}
		}
		// A prompt check stops at the question it checks (24.6): the
		// step's own, or the file's on the last step.
		if st.Expect.asks() || (i == len(steps)-1 && s.Expect.asks()) {
			stop.question = true
		}
		label := fmt.Sprintf("step %d (%s)", i+1, st.label())
		res, err := g.runStep(st.SituationAction, stop, multi)
		if err != nil {
			if multi {
				return Answer{}, fmt.Errorf("situation: %s: %s", label, strings.TrimPrefix(err.Error(), "situation: "))
			}
			return Answer{}, err
		}
		ans.Events = append(ans.Events, res.events...)
		if !res.allowed {
			label += fmt.Sprintf(", refused by %s: %s", res.rule, res.reason)
		}
		if !res.allowed && ans.Allowed {
			ans.Allowed, ans.Rule, ans.Reason, ans.Say, refused = false, res.rule, res.reason, res.say, label
		}
		if st.Expect != nil {
			checked = true
			for _, m := range g.mismatches(res.allowed, res.rule, res.events, *st.Expect) {
				mismatch = append(mismatch, label+": "+m)
			}
		}
	}
	if s.Expect != nil {
		for _, m := range g.mismatches(ans.Allowed, ans.Rule, ans.Events, *s.Expect) {
			switch {
			case !multi:
				mismatch = append(mismatch, m)
			case refused != "":
				mismatch = append(mismatch, "end, first refusal at "+refused+": "+m)
			default:
				mismatch = append(mismatch, "end: "+m)
			}
		}
	}
	if checked {
		ok := len(mismatch) == 0
		ans.Matches, ans.Mismatch = &ok, mismatch
	}
	return ans, nil
}

// stepList is the action as one step, or the steps; a situation has
// one of the two (24.4).
func (s Situation) stepList() ([]SituationStep, error) {
	switch {
	case s.Action != nil && s.Steps != nil:
		return nil, errors.New("situation: has both action and steps; an action is one step, so use one of the two")
	case s.Action != nil:
		return []SituationStep{{SituationAction: *s.Action}}, nil
	case s.Steps == nil:
		return nil, errors.New("situation: needs an action or steps")
	case len(s.Steps) == 0:
		return nil, errors.New("situation: steps is empty")
	}
	return s.Steps, nil
}

// label names a step in messages: its kind and card, "play bomb".
func (a SituationAction) label() string {
	if a.Card == "" {
		return a.Kind
	}
	return a.Kind + " " + string(a.Card)
}

// stepResult is what one step did.
type stepResult struct {
	allowed      bool
	rule, reason string
	say          Msg
	events       []Event
}

// stopAt is where settling stops early: at a question once the choices
// are used up (a following "choose" step, or a prompt check), and at a
// roll waiting to be thrown (a following "throw" step).
type stopAt struct{ question, throw bool }

// runStep tries one step and, when allowed, settles with its choices,
// stopping as stop says. steps is true for Situation.Steps, which
// settle to the open action phase.
func (g *Game) runStep(a SituationAction, stop stopAt, steps bool) (stepResult, error) {
	choices := a.Choices
	var in Intent
	switch a.Kind {
	case "choose":
		w := g.Prompt()
		if w.Kind != PromptChoose && w.Kind != PromptOrder {
			return stepResult{}, fmt.Errorf("situation: choose: no question is open; the game waits for %s", waitingFor(w.Kind))
		}
		if len(choices) == 0 {
			return stepResult{}, fmt.Errorf("situation: choose needs choices, options %v", g.PlainOptions())
		}
		ev, err := g.settle(choices, stop, steps)
		return stepResult{allowed: true, events: ev}, err
	case "throw":
		w := g.Prompt()
		if w.Kind != PromptThrow {
			return stepResult{}, fmt.Errorf("situation: throw: no roll waits to be thrown; the game waits for %s", waitingFor(w.Kind))
		}
		in = Intent{Player: w.Player, Kind: IntentThrow}
	case "reveal":
		var err error
		if in, choices, err = a.reveal(g); err != nil {
			return stepResult{}, err
		}
	default:
		var err error
		if in, err = a.intent(g); err != nil {
			return stepResult{}, err
		}
	}
	ev, err := g.Apply(in)
	var re *RuleError
	switch {
	case errors.As(err, &re):
		return stepResult{rule: re.Rule, reason: re.Reason, say: re.Msg()}, nil
	case err != nil:
		return stepResult{}, err
	}
	more, err := g.settle(choices, stop, steps)
	return stepResult{allowed: true, events: append(ev, more...)}, err
}

// waitingFor names a prompt kind in step errors.
func waitingFor(k PromptKind) string {
	switch k {
	case PromptNone:
		return "nothing"
	case PromptPriority:
		return "a priority action or a pass"
	case PromptDiscard:
		return "a discard"
	case PromptGameOver:
		return "nothing: the game is over"
	case PromptChoose, PromptOrder:
		return "an answer"
	case PromptThrow:
		return "a throw"
	}
	return fmt.Sprintf("prompt %d", k)
}

// reveal puts the card on top of the monster deck and has the player
// attack the deck, placing it in monster slot 1 (as
// enginetest.RevealFromDeck); the step's choices follow.
func (a SituationAction) reveal(g *Game) (Intent, []string, error) {
	if a.Player < 0 || a.Player >= len(g.Players) {
		return Intent{}, nil, errors.New("situation: no such player")
	}
	def, ok := g.cards.find(a.Card)
	if !ok {
		return Intent{}, nil, fmt.Errorf("situation: reveal: unknown card %q", a.Card)
	}
	if d, ok := def.Kind.Deck(); !ok || d != MonsterDeck {
		return Intent{}, nil, fmt.Errorf("situation: reveal: %q is no card of the monster deck", a.Card)
	}
	if len(g.Monsters) == 0 {
		return Intent{}, nil, errors.New("situation: reveal: no monster slot to place it in")
	}
	var top CardRef
	if id, ok := g.Monsters[0].TopOf(); ok {
		top = g.Object(id).Card
	}
	g.PutOnDeck(MonsterDeck, a.Card)
	choices := append([]string{"ask.deck.monster", g.Plain(SlotOption(1, top))}, a.Choices...)
	return Intent{Player: PlayerID(a.Player), Kind: IntentAttack}, choices, nil
}

// build creates the game for a situation.
func (s Situation) build(sets []CardSet) (*Game, error) {
	names := make([]string, len(sets))
	for i, set := range sets {
		names[i] = set.Name
	}
	// A situation names a set by its code ("b2") or its name.
	for _, want := range s.Sets {
		found := false
		for _, set := range sets {
			found = found || set.Name == want || (set.Code != "" && set.Code == want)
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
		if err := g.setUpPlayer(&pl, i, sp); err != nil {
			return nil, err
		}
		g.Players = append(g.Players, pl)
	}
	// Decks: every card of the sets, shuffled; enough for refills.
	for i, d := range idx.defs {
		if deck, ok := d.Kind.Deck(); ok && !d.Outside {
			for range idx.copies(i) {
				g.AddToDeck(deck, d.Ref)
			}
		}
	}
	for d := range deckCount {
		g.ShuffleDeck(d)
	}
	if err := g.setUpDecks(s.Setup.Decks); err != nil {
		return nil, err
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
	for i, ref := range s.Setup.Rooms {
		if err := need(ref); err != nil {
			return nil, err
		}
		g.Rooms = append(g.Rooms, Slot{})
		g.putInSlot(g.newObject(ref, Zone{Kind: ZoneOutside}, NoPlayer), RoomSlot, i)
	}
	for _, ref := range s.Setup.Bosses {
		if g.PutInFinalBossSlot(ref) == 0 {
			return nil, fmt.Errorf("situation: %q is no monster for a final boss slot", ref)
		}
	}
	for i, ref := range s.Setup.Minions {
		if err := need(ref); err != nil {
			return nil, err
		}
		g.Minions = append(g.Minions, Slot{})
		g.putInSlot(g.newObject(ref, Zone{Kind: ZoneOutside}, NoPlayer), MinionSlot, i)
	}
	for _, ref := range s.Setup.BonusSouls {
		if err := need(ref); err != nil {
			return nil, err
		}
		g.BonusSouls = append(g.BonusSouls, g.newObject(ref, Zone{Kind: ZoneOutside}, NoPlayer))
	}
	if err := g.setUpSlots(s.Setup); err != nil {
		return nil, err
	}
	active := PlayerID(s.Setup.Active)
	g.Turn = Turn{Active: active, Number: 1, Step: StepAction, Entered: true, LootPlays: 1, Attacks: 1, Purchases: 1}
	if err := s.Setup.setUpTurn(&g.Turn); err != nil {
		return nil, err
	}
	if err := g.setUpDeclared(s.Setup); err != nil {
		return nil, err
	}
	if err := g.setUpRolls(s.Setup.Rolls); err != nil {
		return nil, err
	}
	holder, stackErr := g.stackUp(s.Setup.Stack, active)
	if stackErr != nil {
		return nil, stackErr
	}
	if err := g.setUpDuties(s.Setup); err != nil {
		return nil, err
	}
	g.events = nil
	return g, g.openSetupWindow(s.Setup.Priority, holder)
}

// intent turns an action into an engine intent.
func (a SituationAction) intent(g *Game) (Intent, error) {
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
		// A room in a room slot too: the engine refuses it to anyone
		// but the active player (R-CARD-33), as in a game.
		for _, id := range append(g.controlled(p), g.roomsInPlay()...) {
			if g.Object(id).Card == a.Card {
				return Intent{Player: p, Kind: IntentActivate, Objects: []ObjectID{id}, Choice: a.Ability}, nil
			}
		}
		return Intent{}, fmt.Errorf("situation: player %d controls no %q, and no room slot holds it", a.Player, a.Card)
	default:
		return Intent{}, fmt.Errorf("situation: unknown action kind %q (play, activate, attack, purchase, end_turn, pass, throw, choose, reveal)", a.Kind)
	}
}

// OrderOf answers the open order prompt by card refs, top first: each
// ref names the next card of the prompt not used yet.
func (g *Game) OrderOf(refs []string) ([]ObjectID, error) {
	if g.Waiting.Kind != PromptOrder || g.Choice == nil {
		return nil, fmt.Errorf("situation: no order prompt is open")
	}
	left := append([]ObjectID(nil), g.Choice.Objects...)
	var order []ObjectID
	for _, ref := range refs {
		i := slices.IndexFunc(left, func(id ObjectID) bool { return string(g.Object(id).Card) == ref })
		if i < 0 {
			return nil, fmt.Errorf("situation: %q is not a card left to order", ref)
		}
		order = append(order, left[i])
		left = slices.Delete(left, i, i+1)
	}
	return order, nil
}

// settle answers choices by label and lets everyone pass until the stack
// is empty and the action is over. toAction (steps) also passes empty
// windows until the open action phase (an end_turn runs to the next
// player's turn). It stops at a question with no choice left when
// stop.question, and at a roll to throw when stop.throw.
func (g *Game) settle(choices []string, stop stopAt, toAction bool) ([]Event, error) {
	var out []Event
	for range 1000 {
		w := g.Prompt()
		var in Intent
		switch {
		case g.Over:
			return out, nil
		case (w.Kind == PromptChoose || w.Kind == PromptOrder) && len(choices) == 0 && stop.question:
			return out, nil
		case w.Kind == PromptChoose:
			if len(choices) == 0 {
				return out, fmt.Errorf("situation: a choice is needed, options %v", g.PlainOptions())
			}
			i := g.OptionIndex(choices[0])
			if i < 0 {
				return out, fmt.Errorf("situation: %q is not an option of %v", choices[0], g.PlainOptions())
			}
			choices = choices[1:]
			in = Intent{Player: w.Player, Kind: IntentChoose, Choice: i}
		case w.Kind == PromptOrder:
			if len(choices) < w.Count {
				return out, fmt.Errorf("situation: an order of %d cards is needed, %d choices left", w.Count, len(choices))
			}
			order, err := g.OrderOf(choices[:w.Count])
			if err != nil {
				return out, err
			}
			choices = choices[w.Count:]
			in = Intent{Player: w.Player, Kind: IntentOrder, Objects: order}
		case w.Kind == PromptThrow && stop.throw:
			return out, nil
		case w.Kind == PromptThrow:
			in = Intent{Player: w.Player, Kind: IntentThrow}
		case w.Kind == PromptPriority && (len(g.Stack) > 0 || g.Attack.On || g.Purchase.On || (toAction && !g.inOpenActionPhase())):
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

// mismatches compares a result, its events and the table with an
// expectation.
func (g *Game) mismatches(allowed bool, rule string, events []Event, e Expectation) []string {
	var out []string
	if allowed != e.Allowed {
		out = append(out, fmt.Sprintf("allowed is %v, expected %v", allowed, e.Allowed))
	}
	if !e.Allowed && e.Rule != "" && rule != e.Rule {
		out = append(out, fmt.Sprintf("refused by %s, expected %s", rule, e.Rule))
	}
	c := checks(out)
	emptyChecks(&c, "after", len(e.After), e.After != nil)
	emptyChecks(&c, "monsters", len(e.Monsters), e.Monsters != nil)
	emptyChecks(&c, "counters", len(e.Counters), e.Counters != nil)
	emptyChecks(&c, "decks", len(e.Decks), e.Decks != nil)
	emptyChecks(&c, "discards", len(e.Discards), e.Discards != nil)
	emptyChecks(&c, "deck_bottoms", len(e.DeckBottoms), e.DeckBottoms != nil)
	emptyChecks(&c, "events", len(e.Events), e.Events != nil)
	for _, x := range e.After {
		g.checkPlayer(&c, x)
	}
	for _, x := range e.Monsters {
		g.checkMonster(&c, x)
	}
	g.checkDecks(&c, e.Decks, e.Discards)
	g.checkBottoms(&c, e.DeckBottoms)
	g.checkGame(&c, e)
	g.checkPrompt(&c, e.Prompt)
	g.checkResolved(&c, e.Resolved, events)
	checkEvents(&c, e.Events, events, len(g.Players))
	checkRolls(&c, e.ResolvedRolls, events)
	return c
}

func containsRef(list []CardRef, r CardRef) bool {
	for _, x := range list {
		if x == r {
			return true
		}
	}
	return false
}

// BuildTable builds a game from a situation setup: the active player is
// in the open action phase, or the window setup.phase and
// setup.declared name. For card tests, the sandbox and card authors.
func BuildTable(setup SituationSetup, sets ...CardSet) (*Game, error) {
	names := make([]string, len(sets))
	for i, set := range sets {
		names[i] = set.Name
	}
	return Situation{Sets: names, Setup: setup}.build(sets)
}

// Settle answers choose prompts with the given option labels, in order,
// and passes priority until the stack is empty and no attack or purchase
// is going on.
func (g *Game) Settle(choices ...string) ([]Event, error) {
	return g.settle(choices, stopAt{}, false)
}

// Run does the queued automatic work after direct calls (RefillSlots),
// then gives priority again. Only for tests and the sandbox.
func (g *Game) Run() []Event {
	if g.Waiting.Kind == PromptPriority {
		g.Priority.Deferred, g.Priority.Holder = true, g.Waiting.Player
		g.Waiting = Prompt{}
	}
	g.run()
	return g.takeEvents()
}

// ForceRolls makes the next dice rolls come out as given. Only for
// tests and the sandbox; a real game never calls it.
func (g *Game) ForceRolls(rolls ...int) {
	g.forcedRolls = append(g.forcedRolls, rolls...)
}
