package enginetest

import (
	"encoding/json"
	"errors"
	"fmt"
	"os" //nolint:depguard // only tests import enginetest, so the wasm engine never does; the dump writes files (24.7)
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"sync"
	"unicode"

	engine "github.com/FourSoulsRepo/rules_engine"
)

// DumpEnv names the folder a Table writes its recording to (24.7):
// SITUATION_DUMP=.cache/situations go test ./cards/b2 -run TestDinga.
// A relative folder is from the repository root. Each table writes
// <folder>/<package>/<TestName>.json when its test ends: a draft
// situation file with title, sets, setup and steps, an empty source,
// no expect (ADR 034, item 11), and a todo list the author clears.
const DumpEnv = "SITUATION_DUMP"

// Dump is the file a recording writes. It is a draft, not a situation
// file: Todo is no situation field, so the runner rejects the file until
// the author removes it, fills Source and writes expect from the rules.
type Dump struct {
	Title  string                `json:"title"`
	Source []string              `json:"source"`
	Sets   []string              `json:"sets"`
	Setup  engine.SituationSetup `json:"setup"`
	Steps  []DumpStep            `json:"steps"`
	Todo   []string              `json:"todo"`
}

// DumpStep is one recorded step, in the situation steps form. Player is
// left out where the runner does not read it (choose, throw).
type DumpStep struct {
	Player  *int           `json:"player,omitempty"`
	Kind    string         `json:"kind"`
	Card    engine.CardRef `json:"card,omitempty"`
	Ability int            `json:"ability,omitempty"`
	Choices []string       `json:"choices,omitempty"`
}

// recorder records one Table: its setup, the hand-set state the setup
// can say, every die thrown, and a step per Table helper. Helpers nest
// (Activate calls Do), so only the outermost one records (depth 1).
type recorder struct {
	tb    *Table
	name  string // the test's name
	dir   string
	sets  []engine.CardSet
	setup engine.SituationSetup
	// base is the table as observed right after the build, and decks its
	// decks by ref, bottom first; the first boundary compares with them.
	base  engine.SituationSetup
	decks [][]engine.CardRef
	steps []DumpStep
	rolls []int
	todo  []string
	depth int
	// dropped counts passes left out since the last step, thrown dice
	// thrown outside a helper.
	dropped, thrown int
	// started: the first boundary is done, so setup is final.
	started bool
	// logLen is the game's applied-intent log length at the last helper's
	// end; last is the flat state then.
	logLen int
	last   map[string]string
}

// startRecorder records tb when DumpEnv is set, else returns nil.
func startRecorder(tb *Table, setup engine.SituationSetup, sets []engine.CardSet) *recorder {
	dir := os.Getenv(DumpEnv)
	if dir == "" {
		return nil
	}
	r := &recorder{tb: tb, name: tb.T.Name(), dir: dir, sets: sets, setup: cloneSetup(tb.T, setup), base: observe(tb.G), decks: deckRefs(tb.G)}
	r.logLen, _ = logLen(tb.G)
	tb.T.Cleanup(r.write)
	return r
}

// enter starts a Table helper. At the outermost helper it records what
// happened since the last one, then record records the helper's own
// step. The returned func ends the helper.
func (tb *Table) enter(record func(r *recorder)) func() {
	r := tb.rec
	if r == nil {
		return func() {}
	}
	r.depth++
	if r.depth == 1 {
		r.boundary(false)
		if record != nil {
			record(r)
		}
	}
	return func() {
		r.depth--
		if r.depth == 0 {
			r.logLen, _ = logLen(tb.G)
			r.last = flatten(tb.G)
		}
	}
}

// saw notes the dice in events: every die thrown becomes setup.rolls.
func (tb *Table) saw(events []engine.Event) {
	if tb.rec == nil {
		return
	}
	for _, e := range events {
		switch {
		case e.Kind != engine.EvDiceRolled || e.Text == "first":
		case len(e.Dice) > 0: // many dice: all of them (R-DICE-11)
			tb.rec.rolls = append(tb.rec.rolls, e.Dice...)
		default:
			tb.rec.rolls = append(tb.rec.rolls, e.Amount)
		}
	}
}

// add appends a step. A play or activation while the stack holds items
// is a response, which a step can't be (owner answer 7).
func (r *recorder) add(s DumpStep) {
	r.flush()
	r.steps = append(r.steps, s)
	if len(r.steps) > 1 && len(r.tb.G.Stack) > 0 && (s.Kind == "activate" || s.Kind == "play") {
		r.note("step %d (%s %s) responds to the stack; steps can't (owner answer 7): start from setup.stack instead", len(r.steps), s.Kind, s.Card)
	}
}

// pass records a pass by p. After a step, a pass in the middle of the
// stack, an attack or a purchase is that step's settling, which the
// runner does itself, so it is left out.
func (r *recorder) pass(p int) {
	g := r.tb.G
	open := g.Turn.Step == engine.StepAction && !g.Turn.EndDeclared && !g.Attack.On && !g.Purchase.On && len(g.Stack) == 0
	if len(r.steps) > 0 && !open {
		r.dropped++
		return
	}
	r.add(step(p, "pass", "", nil))
	r.settled("Pass")
}

// flush notes the passes left out since the last step.
func (r *recorder) flush() {
	if r.dropped > 0 {
		r.note("after step %d: %d passes left out, as the runner settles each step; a check in the middle of the stack stays Go", len(r.steps), r.dropped)
	}
	if r.thrown > 0 {
		r.note("after step %d: %d dice thrown outside a helper are not in rolls", len(r.steps), r.thrown)
	}
	r.dropped, r.thrown = 0, 0
}

// note adds a todo line.
func (r *recorder) note(format string, args ...any) {
	r.todo = append(r.todo, fmt.Sprintf(format, args...))
}

func seatOf(p int) *int { return &p }

// step is a step of kind for player p.
func step(p int, kind string, card engine.CardRef, choices []string) DumpStep {
	return DumpStep{Player: seatOf(p), Kind: kind, Card: card, Choices: slices.Clone(choices)}
}

// settled notes a helper that does not settle, unlike a runner step.
func (r *recorder) settled(helper string) {
	r.note("step %d (%s): %s does not settle, but the runner settles every step; check the steps after it "+
		"(owner answer 7: a response in the middle of a stack is setup.stack)", len(r.steps), helper, helper)
}

// settle records Table.Settle: its choices answer the last step's
// questions, as the runner settles a step with its choices.
func (r *recorder) settle(choices []string) {
	if n := len(r.steps); n > 0 {
		r.steps[n-1].Choices = append(r.steps[n-1].Choices, choices...)
		return
	}
	if w := r.tb.G.Prompt(); w.Kind == engine.PromptPriority {
		r.add(step(int(w.Player), "pass", "", choices))
		r.note("step 1 (Settle): recorded as a pass by the priority holder")
		return
	}
	r.note("Settle%q before any step: no step says it", choices)
}

// intent records an intent as a step. live is true before it is
// applied, so a choose prompt's label can be read; false for one found
// in the log afterwards.
func (r *recorder) intent(in engine.Intent, choices []string, live bool) {
	g := r.tb.G
	p := int(in.Player)
	var card engine.CardRef
	if len(in.Objects) > 0 && int(in.Objects[0]) < len(g.Objects) {
		card = g.Object(in.Objects[0]).Card
	}
	switch in.Kind {
	case engine.IntentActivate:
		s := step(p, "activate", card, choices)
		s.Ability = in.Choice
		r.add(s)
	case engine.IntentPlayLoot:
		r.add(step(p, "play", card, choices))
	case engine.IntentAttack:
		r.add(step(p, "attack", "", choices))
	case engine.IntentPurchase:
		r.add(step(p, "purchase", "", choices))
	case engine.IntentEndTurn:
		r.add(step(p, "end_turn", "", choices))
	case engine.IntentPass:
		r.add(step(p, "pass", "", choices))
	case engine.IntentThrow:
		r.add(DumpStep{Kind: "throw", Choices: slices.Clone(choices)})
	case engine.IntentChoose:
		label := fmt.Sprintf("#%d", in.Choice)
		if opts := g.PlainOptions(); live && in.Choice >= 0 && in.Choice < len(opts) {
			label = opts[in.Choice]
		} else {
			r.note("step %d (choose): option index %d was chosen outside a helper; write its label", len(r.steps)+1, in.Choice)
		}
		r.add(DumpStep{Kind: "choose", Choices: append([]string{label}, choices...)})
	default: // discard, order: no step kind
		r.note("after step %d: player %d's intent %d (objects %v) has no step kind", len(r.steps), p, in.Kind, in.Objects)
	}
}

// boundary records what happened between two helpers: intents applied
// directly (tb.G.Apply, tb.G.Settle) become steps; before the first
// step, hand-set state becomes setup or todo; between steps, any state
// change is a todo. final is the test's end, where writes don't count.
func (r *recorder) boundary(final bool) {
	g := r.tb.G
	log, ok := appliedIntents(g)
	if !ok {
		r.note("enginetest can't read the applied intents (Game.saved.since); direct G.Apply calls are missing")
	}
	if len(log) < r.logLen {
		r.note("after step %d: the applied-intent log was reset (a Save?); direct calls before it may be missing", len(r.steps))
		r.logLen = 0
	}
	direct := log[min(r.logLen, len(log)):]
	r.logLen = len(log)
	switch {
	case !r.started && len(direct) > 0:
		// The game kept its state from just before its first applied
		// intent (Game.saved.origin), hand-set state included (24.11).
		before, err := originGame(g, r.sets)
		if err == nil && len(direct) != len(log) {
			err = errors.New("helpers applied intents before them")
		}
		if err != nil {
			r.note("intents applied directly before the first helper: hand-set state before them is not in setup (%v)", err)
			break
		}
		r.fillSetup(before)
	case !r.started:
		r.fillSetup(g)
	case len(direct) == 0 && !final:
		next := len(r.steps) + 1
		r.todo = append(r.todo, diffFlat(r.last, flatten(g), func(group, parts string) string {
			return fmt.Sprintf("before step %d: %s changed outside a Table helper: %s", next, group, parts)
		})...)
	}
	r.started = true
	r.direct(direct)
}

// direct records intents applied outside the helpers. Passes and throws
// after a step are its settling, which the runner does itself.
func (r *recorder) direct(log []engine.Intent) {
	if len(log) == 0 {
		return
	}
	first := len(r.steps) + 1
	for _, in := range log {
		switch {
		case len(r.steps) > 0 && in.Kind == engine.IntentPass:
			r.dropped++
		case len(r.steps) > 0 && in.Kind == engine.IntentThrow:
			r.thrown++
		default:
			r.intent(in, nil, false)
		}
	}
	switch last := len(r.steps); {
	case last == first:
		r.note("step %d comes from a direct G.Apply or G.Settle call; check it", first)
	case last > first:
		r.note("steps %d to %d come from direct G.Apply or G.Settle calls; check them", first, last)
	}
}

// fillSetup copies the hand-set state the setup can say into it, then
// rebuilds the table from it: what still differs is a todo. g is the
// table before its first step.
func (r *recorder) fillSetup(g *engine.Game) {
	cur := observe(g)
	changed(reflect.ValueOf(&r.setup).Elem(), reflect.ValueOf(r.base), reflect.ValueOf(cur),
		"Players", "Decks", "Stack", "Priority", "Phase", "Declared", "Rolls")
	for i := range min(len(r.setup.Players), len(cur.Players)) {
		changed(reflect.ValueOf(&r.setup.Players[i]).Elem(), reflect.ValueOf(r.base.Players[i]), reflect.ValueOf(cur.Players[i]))
	}
	for d, now := range deckRefs(g) {
		was := r.decks[d]
		if len(now) <= len(was) || !slices.Equal(now[:len(was)], was) {
			continue // a deck that only grew on top is a setup deck; else a todo below
		}
		top := slices.Clone(now[len(was):])
		slices.Reverse(top)
		name := engine.DeckKind(d).String()
		if r.setup.Decks == nil {
			r.setup.Decks = map[string][]engine.CardRef{}
		}
		r.setup.Decks[name] = append(top, r.setup.Decks[name]...)
	}
	rebuilt, err := engine.BuildTable(r.setup, r.sets...)
	if err != nil {
		r.note("setup: the dumped setup does not build: %v", err)
		return
	}
	r.todo = append(r.todo, diffFlat(flatten(rebuilt), flatten(g), func(group, parts string) string {
		return fmt.Sprintf("hand-set before step 1, and no setup field says it: %s (the dumped setup's → the test's): %s", group, parts)
	})...)
}

// changed sets each field of dst whose value differs between base and
// cur to cur's, except the fields named in skip.
func changed(dst, base, cur reflect.Value, skip ...string) {
	for i := range dst.NumField() {
		name := dst.Type().Field(i).Name
		if slices.Contains(skip, name) || reflect.DeepEqual(base.Field(i).Interface(), cur.Field(i).Interface()) {
			continue
		}
		dst.Field(i).Set(cur.Field(i))
	}
}

// observe reads the table as a setup: the fields a setup can say.
func observe(g *engine.Game) engine.SituationSetup {
	s := engine.SituationSetup{Active: int(g.Turn.Active)}
	for _, pl := range g.Players {
		ch := g.Object(pl.Character)
		sp := engine.SituationPlayer{
			Character: ch.Card, Cents: pl.Cents, Damage: pl.Damage, Dead: pl.Dead,
			ExtraLootPlays: pl.ExtraLootPlays, Counters: counters("character", ch.Counters),
		}
		if !ch.Charged {
			sp.Deactivated = append(sp.Deactivated, "character")
		}
		for _, id := range pl.InPlay {
			o := g.Object(id)
			switch o.Role { // curses and the rest: no setup field
			case engine.RoleItem:
				sp.Items = append(sp.Items, o.Card)
				if !o.Charged {
					sp.Deactivated = append(sp.Deactivated, o.Card)
				}
			case engine.RoleSoul:
				sp.Souls = append(sp.Souls, o.Card)
			default:
				continue
			}
			sp.Counters = append(sp.Counters, counters(o.Card, o.Counters)...)
			if o.Eternal {
				sp.Eternal = append(sp.Eternal, o.Card)
			}
		}
		for _, id := range pl.Hand {
			sp.Hand = append(sp.Hand, g.Object(id).Card)
		}
		s.Players = append(s.Players, sp)
	}
	s.Monsters, s.Shop, s.Rooms = tops(g, g.Monsters), tops(g, g.Shop), tops(g, g.Rooms)
	s.Bosses, s.Minions = tops(g, g.Bosses), tops(g, g.Minions)
	var slotted []engine.ObjectID
	for _, row := range [][]engine.Slot{g.Monsters, g.Bosses, g.Minions, g.Rooms} {
		for _, sl := range row {
			if id, ok := sl.TopOf(); ok {
				slotted = append(slotted, id)
			}
		}
	}
	for _, id := range slotted {
		o := g.Object(id)
		if _, seen := s.MonsterDamage[o.Card]; o.Damage > 0 && !seen {
			if s.MonsterDamage == nil {
				s.MonsterDamage = map[engine.CardRef]int{}
			}
			s.MonsterDamage[o.Card] = o.Damage
		}
		s.Counters = append(s.Counters, counters(o.Card, o.Counters)...)
	}
	for _, id := range g.BonusSouls {
		o := g.Object(id)
		s.BonusSouls = append(s.BonusSouls, o.Card)
		s.Counters = append(s.Counters, counters(o.Card, o.Counters)...)
	}
	loot, attacks, purchases := g.Turn.LootPlays, g.Turn.Attacks, g.Turn.Purchases
	s.LootPlays, s.Attacks, s.Purchases = &loot, &attacks, &purchases
	s.MustAttack = mustAttack(g)
	s.MustAttacks, s.MustAttackDeck = g.Turn.MustAttacks, g.Turn.MustAttackDeck
	for _, id := range g.Turn.MustAttackOneOf {
		if ref := attackableRef(g, id); ref != "" {
			s.MustAttackOneOf = append(s.MustAttackOneOf, ref)
		}
	}
	return s
}

// mustAttack is the ref of the card the active player must attack this
// turn (19.235), when setup.must_attack can name it: a card in play the
// active player may attack. Otherwise the rebuilt table lacks the duty,
// and the diff makes it a todo.
func mustAttack(g *engine.Game) engine.CardRef {
	return attackableRef(g, g.Turn.MustAttack)
}

// attackableRef is the ref of id when a setup duty can name it: a card
// in play the active player may attack (19.235, 19.241). A covered one
// (R-ZONE-10) has none.
func attackableRef(g *engine.Game, id engine.ObjectID) engine.CardRef {
	if id == 0 || g.Object(id).Zone.Kind != engine.ZoneInPlay || !g.CanBeAttacked(id) {
		return ""
	}
	for _, pl := range g.Players {
		if pl.Character == id {
			if _, ok := g.PlayerAttackDC(pl.ID); !ok {
				return ""
			}
		}
	}
	return g.Object(id).Card
}

func counters(card engine.CardRef, cs []engine.Counter) []engine.SituationCounter {
	var out []engine.SituationCounter
	for _, c := range cs {
		out = append(out, engine.SituationCounter{Card: card, Name: c.Name, Count: c.Count})
	}
	return out
}

// tops is the top card of each slot; an empty slot has none.
func tops(g *engine.Game, row []engine.Slot) []engine.CardRef {
	var out []engine.CardRef
	for _, sl := range row {
		if id, ok := sl.TopOf(); ok {
			out = append(out, g.Object(id).Card)
		}
	}
	return out
}

// deckRefs is each deck by ref, bottom first.
func deckRefs(g *engine.Game) [][]engine.CardRef {
	out := make([][]engine.CardRef, len(g.Decks))
	for d, ids := range g.Decks {
		for _, id := range ids {
			out[d] = append(out[d], g.Object(id).Card)
		}
	}
	return out
}

// cloneSetup copies a setup, so the test's later changes to its slices
// don't reach the recording.
func cloneSetup(t interface{ Fatal(...any) }, s engine.SituationSetup) engine.SituationSetup {
	data, err := json.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	var out engine.SituationSetup
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatal(err)
	}
	return out
}

// appliedIntents reads the intents the game logged since its last
// clean save (Game.saved.since, unexported; 19.177). It is the
// recorder's only view of tb.G.Apply calls; reflect reads it so the
// engine needs no test hook. ok is false if the field is gone.
func appliedIntents(g *engine.Game) ([]engine.Intent, bool) {
	v := sinceLog(g)
	if !v.IsValid() {
		return nil, false
	}
	out := make([]engine.Intent, v.Len())
	for i := range out {
		in := v.Index(i).FieldByName("Intent")
		out[i] = engine.Intent{
			Player: engine.PlayerID(in.FieldByName("Player").Int()),
			Kind:   engine.IntentKind(in.FieldByName("Kind").Int()),
			Choice: int(in.FieldByName("Choice").Int()),
		}
		objs := in.FieldByName("Objects")
		for j := range objs.Len() {
			out[i].Objects = append(out[i].Objects, engine.ObjectID(objs.Index(j).Int()))
		}
	}
	return out, true
}

// originGame is the game as it was just before its first applied
// intent: the engine keeps that state as its save origin (beforeApply,
// Game.saved.origin, unexported), so a direct G.Apply as the test's
// first step still shows the state the test set by hand before it.
func originGame(g *engine.Game, sets []engine.CardSet) (*engine.Game, error) {
	v := reflect.ValueOf(g).Elem().FieldByName("saved")
	if v.IsValid() {
		v = v.FieldByName("origin")
	}
	if !v.IsValid() || v.Kind() != reflect.Slice || v.Len() == 0 {
		return nil, errors.New("the game kept no state from before them (Game.saved.origin)")
	}
	before, err := engine.Load(v.Bytes(), sets...)
	if err != nil {
		return nil, fmt.Errorf("the state from before them does not load: %w", err)
	}
	return before, nil
}

func logLen(g *engine.Game) (int, bool) {
	v := sinceLog(g)
	if !v.IsValid() {
		return 0, false
	}
	return v.Len(), true
}

func sinceLog(g *engine.Game) reflect.Value {
	v := reflect.ValueOf(g).Elem().FieldByName("saved")
	if !v.IsValid() {
		return v
	}
	v = v.FieldByName("since")
	if v.Kind() != reflect.Slice {
		return reflect.Value{}
	}
	return v
}

var objectIDType = reflect.TypeFor[engine.ObjectID]()

// flatten is the game's exported state as path: value, zero values left
// out. An object ID is replaced by the object's own fields, so two
// tables built in another order compare equal. RNG and the object list
// itself are left out.
func flatten(g *engine.Game) map[string]string {
	f := flat{g: g, out: map[string]string{}}
	v := reflect.ValueOf(g).Elem()
	for i := range v.NumField() {
		sf := v.Type().Field(i)
		if sf.IsExported() && sf.Name != "RNG" && sf.Name != "Objects" {
			f.walk(v.Field(i), sf.Name, 0)
		}
	}
	return f.out
}

type flat struct {
	g   *engine.Game
	out map[string]string
}

const maxDepth = 24

func (f flat) walk(v reflect.Value, path string, depth int) {
	if depth > maxDepth {
		return
	}
	if v.Type() == objectIDType {
		id := engine.ObjectID(v.Int())
		if id < 0 || int(id) >= len(f.g.Objects) {
			f.out[path] = fmt.Sprint(id)
			return
		}
		f.fields(reflect.ValueOf(f.g.Object(id)).Elem(), path, depth+1)
		return
	}
	switch v.Kind() {
	case reflect.Pointer, reflect.Interface:
		if !v.IsNil() {
			f.walk(v.Elem(), path, depth+1)
		}
	case reflect.Struct:
		f.fields(v, path, depth)
	case reflect.Slice, reflect.Array:
		for i := range v.Len() {
			f.walk(v.Index(i), fmt.Sprintf("%s[%d]", path, i), depth+1)
		}
	case reflect.Map:
		for _, k := range v.MapKeys() {
			f.walk(v.MapIndex(k), fmt.Sprintf("%s[%v]", path, k), depth+1)
		}
	case reflect.Func, reflect.Chan, reflect.UnsafePointer:
	default:
		if !v.IsZero() {
			f.out[path] = fmt.Sprint(v)
		}
	}
}

// fields walks a struct's exported fields; an object's ID is left out.
func (f flat) fields(v reflect.Value, path string, depth int) {
	for i := range v.NumField() {
		sf := v.Type().Field(i)
		if !sf.IsExported() || (sf.Name == "ID" && sf.Type == objectIDType) {
			continue
		}
		f.walk(v.Field(i), path+"."+sf.Name, depth+1)
	}
}

// maxDiff caps the parts shown per group and the groups per diff.
const maxDiff = 4

// diffFlat lists where got differs from want, a line per group of
// paths: a top field and its index, as Players[1] or Boosts[0]. line
// gets the group and its parts, "Cents none → 3".
func diffFlat(want, got map[string]string, line func(group, parts string) string) []string {
	var paths []string
	for p, v := range got {
		if want[p] != v {
			paths = append(paths, p)
		}
	}
	for p := range want {
		if _, ok := got[p]; !ok {
			paths = append(paths, p)
		}
	}
	slices.Sort(paths)
	var groups []string
	parts := map[string][]string{}
	for _, p := range paths {
		g, rest := splitGroup(p)
		if _, ok := parts[g]; !ok {
			groups = append(groups, g)
		}
		parts[g] = append(parts[g], fmt.Sprintf("%s %s → %s", rest, shown(want[p]), shown(got[p])))
	}
	var out []string
	for i, g := range groups {
		if i == 2*maxDiff {
			out = append(out, fmt.Sprintf("… and %d more changed fields", len(groups)-i))
			break
		}
		ps := parts[g]
		if len(ps) > maxDiff {
			ps = append(ps[:maxDiff:maxDiff], fmt.Sprintf("%d more", len(parts[g])-maxDiff))
		}
		out = append(out, line(g, strings.Join(ps, "; ")))
	}
	return out
}

// splitGroup splits a path after its top field and that field's first
// index: "Players[1].SkipTurns" is "Players[1]" and "SkipTurns".
func splitGroup(path string) (group, rest string) {
	i := strings.IndexByte(path, '.')
	if i < 0 {
		return path, "value"
	}
	return path[:i], path[i+1:]
}

func shown(v string) string {
	if v == "" {
		return "none"
	}
	return v
}

// write saves the recording when the test ends.
func (r *recorder) write() {
	r.boundary(true)
	r.flush()
	d := Dump{Title: title(r.name), Source: []string{}, Sets: setCodes(r.sets), Setup: r.setup, Steps: r.steps}
	if d.Steps == nil {
		d.Steps = []DumpStep{}
	}
	if len(r.rolls) > 0 {
		d.Setup.Rolls = r.rolls
	}
	d.Todo = append(d.Todo,
		"a draft of "+r.name+" (24.7): write title as a question, fill source, and write expect from the Go test's checks "+
			"and their sources, never from a run (CLAUDE.md 5.25); then remove todo and save it under testdata/situations")
	if len(r.rolls) > 0 {
		d.Todo = append(d.Todo, "setup.rolls are every die the test threw, forced or not")
	}
	if len(r.steps) == 0 {
		d.Todo = append(d.Todo, "no step was recorded: the runner needs at least one")
	}
	if r.tb.T.Failed() {
		d.Todo = append(d.Todo, "the Go test failed while recording")
	}
	d.Todo = append(d.Todo, r.todo...)
	data, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		r.tb.T.Errorf("%s: %v", DumpEnv, err)
		return
	}
	path, err := dumpPath(r.dir, r.name)
	if err == nil {
		err = os.MkdirAll(filepath.Dir(path), 0o750)
	}
	if err == nil {
		err = os.WriteFile(path, append(data, '\n'), 0o600)
	}
	if err != nil {
		r.tb.T.Errorf("%s: %v", DumpEnv, err)
	}
}

func setCodes(sets []engine.CardSet) []string {
	out := make([]string, len(sets))
	for i, s := range sets {
		out[i] = s.Name
		if s.Code != "" {
			out[i] = s.Code
		}
	}
	return out
}

// title turns a test name into words: TestBookOfSin/two_rolls is
// "Book Of Sin: two rolls".
func title(name string) string {
	parts := strings.Split(strings.TrimPrefix(name, "Test"), "/")
	for i, p := range parts {
		var b strings.Builder
		rs := []rune(p)
		for j, c := range rs {
			if j > 0 && unicode.IsUpper(c) && (unicode.IsLower(rs[j-1]) || unicode.IsDigit(rs[j-1])) {
				b.WriteRune(' ')
			}
			b.WriteRune(c)
		}
		parts[i] = strings.ReplaceAll(b.String(), "_", " ")
	}
	return strings.Join(parts, ": ")
}

var (
	dumpMu    sync.Mutex
	dumpPaths = map[string]bool{}
)

// dumpPath is <dir>/<package folder>/<test name>.json; a second table
// in one test gets -2, and so on. A relative dir is from the repository
// root, since go test runs in the package folder.
func dumpPath(dir, name string) (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	if !filepath.IsAbs(dir) {
		dir = filepath.Join(repoRoot(wd), dir)
	}
	base := filepath.Join(dir, filepath.Base(wd), filepath.FromSlash(name))
	dumpMu.Lock()
	defer dumpMu.Unlock()
	path := base + ".json"
	for n := 2; dumpPaths[path]; n++ {
		path = fmt.Sprintf("%s-%d.json", base, n)
	}
	dumpPaths[path] = true
	return path, nil
}

// repoRoot is the nearest folder up from dir with a .git entry, or dir.
func repoRoot(dir string) string {
	for d := dir; ; d = filepath.Dir(d) {
		if _, err := os.Stat(filepath.Join(d, ".git")); err == nil {
			return d
		}
		if filepath.Dir(d) == d {
			return dir
		}
	}
}
