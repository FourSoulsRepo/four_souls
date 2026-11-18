package rulesengine

import (
	"bytes"
	"embed"
	"encoding/json"
	"io/fs"
	"slices"
	"strings"
	"testing"
)

// situationFiles is the rule tree; TestSituations in
// situation_files_test.go runs it, these tests read single files.
//
//go:embed testdata/situations
var situationFiles embed.FS

func TestSituationWithoutExpectationJustAnswers(t *testing.T) {
	data := []byte(`{"title":"q","sets":["test"],"setup":{"players":[{"character":"hero_a"},{"character":"hero_b"}],"active":0},
		"action":{"player":1,"kind":"end_turn"}}`)
	ans, err := RunSituation(data, testSet)
	if err != nil {
		t.Fatal(err)
	}
	if ans.Allowed || ans.Rule == "" || ans.Reason == "" || ans.Matches != nil {
		t.Errorf("answer %+v; want a refusal with reason and no match result", ans)
	}
}

// A refusal carries the player's words as a message key with values,
// so the website shows it in its language (ADR 026); an allowed answer
// has none.
func TestSituationAnswerSay(t *testing.T) {
	data, err := fs.ReadFile(situationFiles, "testdata/situations/card/loot-needs-a-loot-play.json")
	if err != nil {
		t.Fatal(err)
	}
	ans, err := RunSituation(data, testSet, abilitySet)
	if err != nil {
		t.Fatal(err)
	}
	if ans.Allowed || ans.Say.IsZero() || ans.Reason == "" {
		t.Fatalf("answer %+v; want a refusal with Say and Reason", ans)
	}
	out, err := json.Marshal(ans)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(out, []byte(`"say":{"key":"`)) || !bytes.Contains(out, []byte(`"reason":`)) {
		t.Errorf("JSON %s; want say and reason", out)
	}

	data, err = fs.ReadFile(situationFiles, "testdata/situations/abil/razor-pays-and-damages.json")
	if err != nil {
		t.Fatal(err)
	}
	if ans, err = RunSituation(data, testSet, abilitySet); err != nil || !ans.Allowed || !ans.Say.IsZero() {
		t.Fatalf("answer %+v, %v; want allowed without Say", ans, err)
	}
	if out, err = json.Marshal(ans); err != nil || bytes.Contains(out, []byte(`"say"`)) {
		t.Errorf("JSON %s, %v; want no say", out, err)
	}
}

// A situation names a set by its card_db code or by its name.
func TestSituationSetByCode(t *testing.T) {
	coded := testSet
	coded.Code = "t1"
	for _, want := range []string{"t1", "test"} {
		data := []byte(`{"sets":["` + want + `"],"setup":{"players":[{"character":"hero_a"},{"character":"hero_b"}],"active":0},
			"action":{"player":0,"kind":"end_turn"}}`)
		if _, err := RunSituation(data, coded); err != nil {
			t.Errorf("%s: %v", want, err)
		}
	}
}

func TestSituationErrors(t *testing.T) {
	for _, bad := range []string{
		`not json`,
		`{"sets":["missing"],"setup":{"players":[{"character":"hero_a"}]},"action":{"kind":"pass"}}`,
		`{"sets":["test"],"setup":{"players":[{"character":"nope"}]},"action":{"kind":"pass"}}`,
		`{"sets":["test"],"setup":{"players":[{"character":"hero_a"}]},"action":{"kind":"dance"}}`,
		`{"sets":["test"],"setup":{"players":[{"character":"hero_a"}]},"action":{"kind":"play","card":"penny"}}`,
	} {
		if _, err := RunSituation([]byte(bad), testSet); err == nil {
			t.Errorf("no error for %s", bad)
		}
	}
}

// A stacked situation resolves last in, first out: the ability on top,
// then the loot card, then the dice roll at the bottom (12.5).
func TestStackedEventsResolveTopFirst(t *testing.T) {
	data, err := fs.ReadFile(situationFiles, "testdata/situations/stack/three-stacked-events.json")
	if err != nil {
		t.Fatal(err)
	}
	ans, err := RunSituation(data, testSet, abilitySet)
	if err != nil {
		t.Fatal(err)
	}
	var order []string
	roll := 0
	for _, e := range ans.Events {
		switch e.Kind { //nolint:exhaustive // only these two matter here
		case EvStackResolved:
			if e.Item.Kind != StackDamage { // the bomb's damage goes on the stack too
				order = append(order, string(e.Card))
			}
		case EvRollResolved:
			roll = e.Amount
		}
	}
	if want := []string{"coin_bag", "bomb_loot", ""}; !slices.Equal(order, want) || roll != 4 {
		t.Errorf("resolved %q, roll %d; want %q, roll 4", order, roll, want)
	}
}

// stackedTable builds a two-seat table with the given stack.
func stackedTable(t *testing.T, stack string, p0, p1 SituationPlayer, monsters ...CardRef) *Game {
	t.Helper()
	var items []SituationStack
	if err := json.Unmarshal([]byte(stack), &items); err != nil {
		t.Fatal(err)
	}
	g, err := BuildTable(SituationSetup{Players: []SituationPlayer{p0, p1}, Monsters: monsters, Stack: items}, testSet, abilitySet)
	if err != nil {
		t.Fatal(err)
	}
	return g
}

// A mode and a target on a stacked loot, and an ability aimed at a
// dice roll below it, use the labels of action.choices.
func TestStackedModeAndTargets(t *testing.T) {
	g := stackedTable(t, `[
		{"kind": "play", "card": "dagaz_like", "owner": 1, "choices": ["Choose a player.", "player 1 (hero_a)"]},
		{"dice_roll": 2, "owner": 0},
		{"kind": "activate", "card": "d6", "owner": 0, "choices": ["ask.roll_of 2"]}
	]`, SituationPlayer{Character: "hero_a", Items: []CardRef{"d6"}}, SituationPlayer{Character: "hero_b"})
	top := g.Stack[2]
	if top.Kind != StackAbility || len(top.Targets) != 1 || top.Targets[0].StackID != g.Stack[1].ID {
		t.Fatalf("top %+v; want The D6-like ability aimed at the roll", top)
	}
	if loot := g.Stack[0]; loot.Kind != StackLoot || loot.Mode != 1 || g.Object(loot.Source).Zone.Kind != ZoneStack {
		t.Fatalf("loot %+v; want mode 1, the card on the stack", loot)
	}
	if _, err := g.Settle(); err != nil {
		t.Fatal(err)
	}
	if g.Players[1].Cents != 1 || len(g.Stack) != 0 {
		t.Errorf("cents %d, stack %d; want 1 and empty", g.Players[1].Cents, len(g.Stack))
	}
}

func TestStackedTrigger(t *testing.T) {
	g := stackedTable(t, `[{"kind": "trigger", "card": "allowance", "owner": 0}]`,
		SituationPlayer{Character: "hero_a", Items: []CardRef{"allowance"}}, SituationPlayer{Character: "hero_b"})
	if it := g.Stack[0]; it.Kind != StackTrigger || it.Controller != 0 {
		t.Fatalf("item %+v; want player 1's trigger", it)
	}
	if _, err := g.Settle(); err != nil {
		t.Fatal(err)
	}
	if g.Players[0].Cents != 1 {
		t.Errorf("cents %d; want 1", g.Players[0].Cents)
	}
}

// A stacked attack roll means the attack is declared on that monster.
func TestStackedAttackRoll(t *testing.T) {
	g := stackedTable(t, `[{"kind": "attack_roll", "dice_roll": 6, "card": "dummy", "owner": 0}]`,
		SituationPlayer{Character: "hero_a"}, SituationPlayer{Character: "hero_b"}, "dummy")
	if !g.Attack.Started || g.Object(g.Attack.Target).Card != "dummy" || g.Turn.Attacks != 0 {
		t.Fatalf("attack %+v, attacks left %d; want an attack on dummy, none left", g.Attack, g.Turn.Attacks)
	}
	ev, err := g.Settle()
	if err != nil {
		t.Fatal(err)
	}
	died := slices.ContainsFunc(ev, func(e Event) bool { return e.Kind == EvDied && e.Card == "dummy" })
	if !died || g.Attack.On {
		t.Errorf("dummy died %v, attack on %v; want it dead and the attack over", died, g.Attack.On)
	}
}

func TestStackedErrors(t *testing.T) {
	setup := func(stack string) string {
		return `{"sets":["test","abilities"],"setup":{"players":[{"character":"hero_a","items":["coin_bag","allowance"]},
			{"character":"hero_b"}],"active":0,"monsters":["dummy"],"stack":` + stack + `},"action":{"player":0,"kind":"pass"}}`
	}
	for _, bad := range []string{
		`[{"loot": "penny", "owner": 0}]`,
		`[{"kind": "cast", "card": "penny", "owner": 0}]`,
		`[{"dice_roll": 7, "owner": 0}]`,
		`[{"dice_roll": 3, "owner": 2}]`,
		`[{"dice_roll": 3, "owner": 0, "card": "penny"}]`,
		`[{"kind": "play", "card": "nope", "owner": 0}]`,
		`[{"kind": "play", "card": "coin_bag", "owner": 0}]`,
		`[{"kind": "play", "card": "penny", "owner": 0, "choices": ["x"]}]`,
		`[{"kind": "play", "card": "bomb_loot", "owner": 0}]`,
		`[{"kind": "play", "card": "bomb_loot", "owner": 0, "choices": ["player 3 (hero_c)"]}]`,
		`[{"kind": "play", "card": "bomb_loot", "owner": 0, "choices": ["dummy", "dummy"]}]`,
		`[{"kind": "play", "card": "penny", "owner": 0, "dice_roll": 2}]`,
		`[{"kind": "activate", "card": "coin_bag", "owner": 1}]`,
		`[{"kind": "activate", "card": "coin_bag", "owner": 0, "ability": 1}]`,
		`[{"kind": "activate", "card": "allowance", "owner": 0}]`,
		`[{"kind": "trigger", "card": "coin_bag", "owner": 0}]`,
		`[{"kind": "trigger", "card": "allowance", "owner": 1}]`,
		`[{"kind": "attack_roll", "dice_roll": 3, "card": "dummy", "owner": 1}]`,
		`[{"kind": "attack_roll", "dice_roll": 3, "card": "gaper", "owner": 0}]`,
		`[{"kind": "attack_roll", "card": "dummy", "owner": 0}]`,
	} {
		if _, err := RunSituation([]byte(setup(bad)), testSet, abilitySet); err == nil {
			t.Errorf("no error for %s", bad)
		}
	}
	if _, err := RunSituation([]byte(setup(`[{"kind": "play", "card": "bomb_loot", "owner": 0, "choices": ["dummy"]}]`)), testSet, abilitySet); err != nil {
		t.Errorf("a good stack fails: %v", err)
	}
}

// buildSetup builds a table from a setup in JSON, with the test set.
func buildSetup(setup string) (*Game, error) {
	var s SituationSetup
	if err := json.Unmarshal([]byte(setup), &s); err != nil {
		return nil, err
	}
	return BuildTable(s, testSet)
}

// The 24.3 setup fields land where Go tests set them by hand.
func TestSituationSetupFields(t *testing.T) {
	g, err := buildSetup(`{
		"players": [
			{"character": "hero_a", "items": ["trinket"], "souls": ["boss"], "eternal": ["trinket", "boss"],
			 "counters": [{"card": "trinket", "count": 2}, {"card": "trinket", "name": "gold", "count": 1}, {"card": "character", "count": 4}]},
			{"character": "hero_b", "extra_loot_plays": 2}
		],
		"monsters": ["gaper", "boss"], "monster_damage": {"boss": 3},
		"bonus_souls": ["soul_x"], "counters": [{"card": "soul_x", "count": 5}, {"card": "gaper", "name": "level", "count": 1}],
		"attacks": 0, "purchases": 2,
		"decks": {"loot": ["penny"], "monster": ["dummy", "ambush_event"]},
		"rolls": [6, 1]
	}`)
	if err != nil {
		t.Fatal(err)
	}
	p0 := g.Players[0]
	trinket, soul := g.Object(p0.InPlay[0]), g.Object(p0.InPlay[1])
	if !trinket.Eternal || !soul.Eternal || trinket.CountersOf("") != 2 || trinket.CountersOf("gold") != 1 || g.Object(p0.Character).CountersOf("") != 4 {
		t.Errorf("trinket %+v, soul %+v; want eternal, 2 plain and 1 gold counter", trinket, soul)
	}
	if g.Players[1].ExtraLootPlays != 2 {
		t.Errorf("extra loot plays %d, want 2", g.Players[1].ExtraLootPlays)
	}
	gaper, boss := g.Object(g.Monsters[0].Cards[0]), g.Object(g.Monsters[1].Cards[0])
	if boss.Damage != 3 || gaper.Damage != 0 || gaper.CountersOf("level") != 1 || g.Object(g.BonusSouls[0]).CountersOf("") != 5 {
		t.Errorf("gaper %+v, boss %+v; want 3 damage on the boss, counters on the gaper and the soul", gaper, boss)
	}
	if g.Turn.Attacks != 0 || g.Turn.Purchases != 2 || g.Turn.LootPlays != 1 {
		t.Errorf("turn %+v; want 0 attacks, 2 purchases, 1 loot play", g.Turn)
	}
	top := func(d DeckKind, n int) []CardRef {
		var out []CardRef
		for _, id := range g.DeckTop(d, n) {
			out = append(out, g.Object(id).Card)
		}
		return out
	}
	if got := top(MonsterDeck, 2); !slices.Equal(got, []CardRef{"dummy", "ambush_event"}) {
		t.Errorf("monster deck top %v; want dummy, then ambush_event", got)
	}
	if got := top(LootDeck, 1); !slices.Equal(got, []CardRef{"penny"}) {
		t.Errorf("loot deck top %v; want penny", got)
	}
	if !slices.Equal(g.forcedRolls, []int{6, 1}) {
		t.Errorf("forced rolls %v; want [6 1]", g.forcedRolls)
	}
}

// A bad setup field is an error that names the item and its options.
func TestSituationSetupFieldErrors(t *testing.T) {
	const seats = `"players": [{"character": "hero_a", "items": ["trinket"]}, {"character": "hero_b"}], "monsters": ["gaper"]`
	for _, c := range []struct{ extra, want string }{
		{`"rolls": [3, 0]`, "roll 2 is 0, not 1 to 6"},
		{`"rolls": [7]`, "roll 1 is 7, not 1 to 6"},
		{`"decks": {"shop": ["trinket"]}`, `"shop" is no deck, options [treasure loot monster room]`},
		{`"decks": {"loot": ["nope"]}`, `unknown card "nope"`},
		{`"decks": {"loot": ["trinket"]}`, `"trinket" is no card of the loot deck`},
		{`"monster_damage": {"boss": 1}`, `"boss" is in no slot, options [gaper]`},
		{`"monster_damage": {"gaper": 2}`, `"gaper" has 2 HP; damage 2 is not 0 to 1`},
		{`"counters": [{"card": "trinket", "count": 1}]`, `"trinket" is in no slot and no bonus soul, options [gaper]`},
		{`"counters": [{"card": "gaper", "count": 0}]`, `"gaper" count 0 is below 1`},
		{`"counters": [{"card": "gaper", "amount": 1}]`, `counter (fields card, name, count)`},
		{`"attacks": -1`, "attacks -1 is below 0"},
		{`"purchases": -1`, "purchases -1 is below 0"},
	} {
		_, err := buildSetup(`{` + seats + `, ` + c.extra + `}`)
		if err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s: error %v; want %q", c.extra, err, c.want)
		}
	}
	for _, c := range []struct{ player, want string }{
		{`"counters": [{"card": "penny", "count": 1}]`, `player 0 counters: "penny" is not in play, options [character trinket]`},
		{`"eternal": ["character"]`, `player 0 eternal: "character" is not in play, options [trinket]`},
		{`"extra_loot_plays": -1`, "player 0 extra_loot_plays -1 is below 0"},
	} {
		_, err := buildSetup(`{"players": [{"character": "hero_a", "items": ["trinket"], ` + c.player + `}, {"character": "hero_b"}]}`)
		if err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s: error %v; want %q", c.player, err, c.want)
		}
	}
}

// razorSteps is a payload with Razor ("Pay 5¢: deal 1 damage") and
// hero_a's 7 cents; extra is put in setup, steps after "steps":.
func razorSteps(extra, steps string) string {
	return `{"sets":["test","abilities"],"setup":{"players":[{"character":"hero_a","items":["razor"],"cents":7},
		{"character":"hero_b"}],"active":0,"monsters":["gaper"]` + extra + `},"steps":` + steps + `}`
}

// A step may leave a question for a following choose step, may be
// refused, and later steps still run; each step's expect is checked
// after it, the file's at the end (24.4).
func TestSituationSteps(t *testing.T) {
	data := razorSteps("", `[
		{"player": 0, "kind": "activate", "card": "razor"},
		{"kind": "choose", "choices": ["player 2 (hero_b)"], "expect": {"allowed": true, "after": [{"player": 0, "cents": 2}, {"player": 1, "damage": 1}]}},
		{"player": 0, "kind": "activate", "card": "razor", "expect": {"allowed": false, "rule": "R-ABIL-07"}},
		{"player": 0, "kind": "end_turn"}]`)
	// R-ABIL-07: 2¢ can't pay 5¢. R-TURN-13: damage heals as the turn passes.
	end := `, "expect": {"allowed": false, "rule": "R-ABIL-07", "after": [{"player": 1, "damage": 0}]}`
	ans, err := RunSituation([]byte(strings.TrimSuffix(data, "}")+end+"}"), testSet, abilitySet)
	if err != nil {
		t.Fatal(err)
	}
	if ans.Allowed || ans.Rule != "R-ABIL-07" || ans.Matches == nil || !*ans.Matches {
		t.Fatalf("answer %+v; want the first refusal, R-ABIL-07, and a match", ans)
	}

	// Wrong values: each mismatch names its step, or the end.
	wrong := strings.Replace(data, `"damage": 1}]}}`, `"damage": 2}]}}`, 1)
	end = `, "expect": {"allowed": true, "after": [{"player": 1, "damage": 1}]}`
	ans, err = RunSituation([]byte(strings.TrimSuffix(wrong, "}")+end+"}"), testSet, abilitySet)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{
		"step 2 (choose): player 1 damage is 1, expected 2",
		"end, first refusal at step 3 (activate razor), refused by R-ABIL-07",
		"allowed is false, expected true",
		"player 1 damage is 0, expected 1",
	}
	got := strings.Join(ans.Mismatch, "\n")
	for _, w := range want {
		if !strings.Contains(got, w) {
			t.Errorf("mismatch %q; want %q in it", ans.Mismatch, w)
		}
	}
}

// Settling stops at a roll for a following throw step.
func TestSituationThrowStep(t *testing.T) {
	data := razorSteps(`,"rolls":[6,6]`, `[
		{"player": 0, "kind": "attack", "choices": ["gaper"]},
		{"kind": "throw"}]`)
	g, err := BuildTable(SituationSetup{
		Players:  []SituationPlayer{{Character: "hero_a"}, {Character: "hero_b"}},
		Monsters: []CardRef{"gaper"},
	}, testSet, abilitySet)
	if err != nil {
		t.Fatal(err)
	}
	if _, err = g.runStep(SituationAction{Player: 0, Kind: "attack", Choices: []string{"gaper"}}, stopAt{throw: true}, true); err != nil {
		t.Fatal(err)
	}
	if g.Prompt().Kind != PromptThrow {
		t.Fatalf("prompt %+v; want a roll to throw", g.Prompt())
	}
	// gaper: HP 2, DC 3, 3¢; hero_a ATK 1. The attack goes on until it
	// ends (R-ATK-11): two 6s hit twice and kill it, for 3¢.
	end := `, "expect": {"allowed": true, "after": [{"player": 0, "cents": 10, "damage": 0}]}`
	ans, err := RunSituation([]byte(strings.TrimSuffix(data, "}")+end+"}"), testSet, abilitySet)
	if err != nil || ans.Matches == nil || !*ans.Matches {
		t.Fatalf("answer %+v, %v; want a match", ans, err)
	}
}

func TestSituationStepErrors(t *testing.T) {
	for _, c := range []struct{ data, want string }{
		{`{"sets":["test"],"setup":{"players":[{"character":"hero_a"}]},"action":{"kind":"pass"},"steps":[{"kind":"pass"}]}`, "both action and steps"},
		{`{"sets":["test"],"setup":{"players":[{"character":"hero_a"}]},"steps":[]}`, "steps is empty"},
		{`{"sets":["test"],"setup":{"players":[{"character":"hero_a"}]}}`, "needs an action or steps"},
		{razorSteps(`,"priority":2`, `[{"kind":"pass"}]`), "priority 2 is not a seat, options 0 to 1"},
		{razorSteps(`,"priority":-1`, `[{"kind":"pass"}]`), "priority -1 is not a seat"},
		{razorSteps("", `[{"kind":"choose","choices":["gaper"]}]`), "step 1 (choose): choose: no question is open"},
		{razorSteps("", `[{"kind":"throw"}]`), "step 1 (throw): throw: no roll waits to be thrown"},
		{razorSteps("", `[{"kind":"pass"},{"kind":"dance"}]`), `step 2 (dance): unknown action kind "dance"`},
		{razorSteps("", `[{"kind":"activate","card":"razor"},{"kind":"choose"}]`), "step 2 (choose): choose needs choices"},
		{razorSteps("", `[{"kind":"activate","card":"razor"},{"kind":"pass"}]`), "step 1 (activate razor): a choice is needed"},
		{razorSteps("", `[{"kind":"reveal","card":"penny"}]`), "step 1 (reveal penny): reveal: \"penny\" is no card of the monster deck"},
	} {
		_, err := RunSituation([]byte(c.data), testSet, abilitySet)
		if err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("error %v; want %q", err, c.want)
		}
	}
}

// setup.priority: with an empty stack the players before the seat, from
// the active player, have passed (R-PRIO-02); when it and the rest pass,
// the window closes as in a game (R-PRIO-05). With a stack it overrides
// the top item's owner.
func TestSituationPriority(t *testing.T) {
	seat := func(n int) *int { return &n }
	players := []SituationPlayer{{Character: "hero_a"}, {Character: "hero_b"}, {Character: "hero_c"}}
	g, err := BuildTable(SituationSetup{Players: players, Active: 1, Phase: "start", Priority: seat(0)}, testSet)
	if err != nil {
		t.Fatal(err)
	}
	if w := g.Prompt(); w.Kind != PromptPriority || w.Player != 0 || g.Priority.Passes != 2 {
		t.Fatalf("prompt %+v, passes %d; want player 0 with 2 passes (players 1 and 2)", w, g.Priority.Passes)
	}
	if _, err = g.Apply(Intent{Player: 0, Kind: IntentPass}); err != nil {
		t.Fatal(err)
	}
	// All three passed with an empty stack: the game moves on (R-STACK-07)
	// to the loot step, where the active player loots 1 and gets priority
	// again (R-TURN-04, R-PRIO-02).
	w := g.Prompt()
	if w.Kind != PromptPriority || w.Player != 1 || g.Priority.Passes != 0 || g.Turn.Step != StepLootWindow || len(g.Players[1].Hand) != 1 {
		t.Errorf("prompt %+v, passes %d, step %d, hand %d; want the active player 1 in the loot step, no passes, 1 card",
			w, g.Priority.Passes, g.Turn.Step, len(g.Players[1].Hand))
	}

	// The open action phase: the active player can't pass (R-TURN-09), so
	// nobody else gets priority on an empty stack; the active seat is fine.
	_, err = BuildTable(SituationSetup{Players: players, Active: 1, Priority: seat(0)}, testSet)
	if err == nil || !strings.Contains(err.Error(), "R-TURN-09") {
		t.Errorf("open action phase, seat 0: error %v; want one naming R-TURN-09", err)
	}
	g, err = BuildTable(SituationSetup{Players: players, Active: 1, Priority: seat(1)}, testSet)
	if err != nil || g.Prompt().Player != 1 || g.Priority.Passes != 0 {
		t.Errorf("open action phase, active seat: error %v, prompt %+v; want the active player 1", err, g.Prompt())
	}

	g, err = BuildTable(SituationSetup{Players: players, Priority: seat(1), Stack: []SituationStack{{DiceRoll: 3, Owner: 2}}}, testSet)
	if err != nil {
		t.Fatal(err)
	}
	if w := g.Prompt(); w.Player != 1 || g.Priority.Passes != 0 {
		t.Errorf("prompt %+v, passes %d; want player 1 over the roll's owner, no passes", w, g.Priority.Passes)
	}
}

// passAround passes for whoever holds priority, n times.
func passAround(t *testing.T, g *Game, n int) {
	t.Helper()
	for range n {
		w := g.Prompt()
		if w.Kind != PromptPriority {
			t.Fatalf("prompt %+v; want priority to pass", w)
		}
		if _, err := g.Apply(Intent{Player: w.Player, Kind: IntentPass}); err != nil {
			t.Fatal(err)
		}
	}
}

// setup.phase and setup.declared start in a real priority window
// (24.10); when everyone passes, the game moves on as in a game.
func TestSituationWindows(t *testing.T) {
	seat := func(n int) *int { return &n }
	players := []SituationPlayer{{Character: "hero_a"}, {Character: "hero_b"}}
	build := func(s SituationSetup) *Game {
		t.Helper()
		s.Players = players
		g, err := BuildTable(s, testSet)
		if err != nil {
			t.Fatal(err)
		}
		return g
	}
	// No phase: the open action phase with one loot play, attack and
	// purchase (R-TURN-05, R-TURN-07), as before 24.10.
	g := build(SituationSetup{})
	if g.Turn.Step != StepAction || g.Turn.LootPlays != 1 || g.Turn.Attacks != 1 || g.Turn.Purchases != 1 || !g.inOpenActionPhase() {
		t.Errorf("turn %+v; want the open action phase, 1 of each", g.Turn)
	}
	// Start and loot: before the action phase there is no loot play,
	// attack or purchase yet (R-TURN-05, R-TURN-07).
	for _, phase := range []string{"start", "loot"} {
		g = build(SituationSetup{Phase: phase, Priority: seat(1)})
		if g.Turn.LootPlays != 0 || g.Turn.Attacks != 0 || g.Turn.Purchases != 0 || g.Prompt().Player != 1 || g.Priority.Passes != 1 {
			t.Errorf("%s: turn %+v, prompt %+v; want none of each, seat 1 after 1 pass", phase, g.Turn, g.Prompt())
		}
	}
	// Loot: everyone passes and the action phase begins (R-TURN-04),
	// with its loot play, attack and purchase (R-TURN-05, R-TURN-07).
	passAround(t, g, 1)
	if !g.inOpenActionPhase() || g.Turn.LootPlays != 1 || g.Turn.Attacks != 1 || g.Prompt().Player != 0 {
		t.Errorf("after the loot step: turn %+v, prompt %+v; want the open action phase", g.Turn, g.Prompt())
	}
	// End: everyone passes; the turn passes to the next player
	// (R-TURN-13), whose start phase opens (R-TURN-03).
	g = build(SituationSetup{Phase: "end"})
	if g.Turn.Step != StepEndTriggers || g.inOpenActionPhase() {
		t.Fatalf("end: turn %+v; want the end phase", g.Turn)
	}
	passAround(t, g, 2)
	if g.Turn.Active != 1 || g.Turn.Step != StepStartTriggers || g.Prompt().Player != 1 {
		t.Errorf("after the end phase: turn %+v, prompt %+v; want player 1's start phase", g.Turn, g.Prompt())
	}
	// A declared attack used one (R-TURN-07); once everyone passes, the
	// active player picks the target (R-ATK-02).
	g = build(SituationSetup{Declared: "attack", Priority: seat(1)})
	if !g.Attack.On || g.Turn.Attacks != 0 || g.Priority.Passes != 1 {
		t.Fatalf("attack: %+v, turn %+v, passes %d; want it declared and used", g.Attack, g.Turn, g.Priority.Passes)
	}
	passAround(t, g, 1)
	if w := g.Prompt(); w.Kind != PromptChoose || w.Player != 0 || g.Choice.Purpose != ChooseAttackTarget {
		t.Errorf("after the attack's window: prompt %+v; want player 0 to pick the target", w)
	}
	// A declared purchase likewise (R-SHOP-02, R-SHOP-05).
	g = build(SituationSetup{Declared: "purchase"})
	if !g.Purchase.On || g.Turn.Purchases != 0 {
		t.Fatalf("purchase: %+v, turn %+v; want it declared and used", g.Purchase, g.Turn)
	}
	passAround(t, g, 2)
	if w := g.Prompt(); w.Kind != PromptChoose || w.Player != 0 || g.Choice.Purpose != ChoosePurchase {
		t.Errorf("after the purchase's window: prompt %+v; want player 0 to pick the item", w)
	}
	// A response waits on the stack (24.10 follow-up): its owner holds
	// priority (R-PRIO-03). It resolves when everyone passes (R-STACK-03);
	// priority passes again from the active player (R-STACK-04, R-PRIO-02);
	// only then is the target or item chosen (R-ATK-02, R-SHOP-02).
	for _, c := range []struct {
		declared string
		purpose  ChoicePurpose
	}{{"attack", ChooseAttackTarget}, {"purchase", ChoosePurchase}} {
		g = build(SituationSetup{Declared: c.declared, Stack: []SituationStack{{DiceRoll: 3, Owner: 1}}})
		on := g.Attack.On || g.Purchase.On
		if !on || g.Turn.Attacks+g.Turn.Purchases != 1 || len(g.Stack) != 1 || g.Prompt().Player != 1 {
			t.Fatalf("%s with a stack: attack %+v, purchase %+v, turn %+v, prompt %+v; want it declared, used, seat 1 to respond",
				c.declared, g.Attack, g.Purchase, g.Turn, g.Prompt())
		}
		passAround(t, g, 2)
		if len(g.Stack) != 0 || g.Prompt().Kind != PromptPriority || g.Prompt().Player != 0 {
			t.Fatalf("%s: after the roll, stack %d, prompt %+v; want an empty stack, priority for player 0", c.declared, len(g.Stack), g.Prompt())
		}
		passAround(t, g, 2)
		if w := g.Prompt(); w.Kind != PromptChoose || w.Player != 0 || g.Choice.Purpose != c.purpose {
			t.Errorf("%s: after the stack: prompt %+v; want player 0 to choose", c.declared, w)
		}
	}
	// A stack with nothing declared is as before 24.10: no attack or
	// purchase is on.
	g = build(SituationSetup{Stack: []SituationStack{{DiceRoll: 3, Owner: 1}}})
	if g.Attack.On || g.Purchase.On || g.Turn.Attacks != 1 || g.Turn.Purchases != 1 {
		t.Errorf("stack only: attack %+v, purchase %+v, turn %+v; want nothing declared", g.Attack, g.Purchase, g.Turn)
	}
}

// A setup.phase or setup.declared that no real window has is an error
// that names the rule or lists the options (24.10).
func TestSituationWindowErrors(t *testing.T) {
	const seats = `"players": [{"character": "hero_a"}, {"character": "hero_b"}], "monsters": ["dummy"]`
	for _, c := range []struct{ extra, want string }{
		{`"phase": "main"`, `phase "main" is no phase, options [start loot action end]`},
		{`"declared": "end"`, `declared "end" is no declaration, options [attack purchase]`},
		{`"phase": "loot", "declared": "attack"`, "declared attack needs phase action, not loot: it is declared only in the action phase (R-ATK-01)"},
		{`"phase": "end", "declared": "purchase"`, "declared purchase needs phase action, not end: it is declared only in the action phase (R-SHOP-01)"},
		{`"declared": "attack", "stack": [{"kind": "attack_roll", "dice_roll": 3, "card": "dummy", "owner": 0}]`, "declared attack: stack item 1 is an attack roll, but the target is chosen only after the declaration's window (R-ATK-02)"},
		{`"declared": "purchase", "stack": [{"dice_roll": 3, "owner": 1}, {"kind": "attack_roll", "dice_roll": 3, "card": "dummy", "owner": 0}]`, "declared purchase: stack item 2 is an attack roll, but a purchase is declared only with an empty stack, so not during an attack (R-SHOP-01)"},
		{`"declared": "attack", "attacks": 0`, "declared attack needs one left (R-TURN-07), and attacks is 0"},
		{`"declared": "purchase", "purchases": 0`, "declared purchase needs one left (R-SHOP-05), and purchases is 0"},
		{`"phase": "start", "loot_plays": 1`, "phase start: loot_plays, attacks and purchases come with the action phase (R-TURN-05, R-TURN-07)"},
		{`"phase": "loot", "attacks": 2`, "phase loot: loot_plays, attacks and purchases come with the action phase"},
		{`"phase": "start", "purchases": 0`, "phase start: loot_plays, attacks and purchases come with the action phase"},
		{`"phase": "end", "stack": [{"kind": "attack_roll", "dice_roll": 3, "card": "dummy", "owner": 0}]`, "an attack roll needs phase action: attacks are declared only in the action phase (R-ATK-01)"},
		{`"priority": 1`, "priority 1: in the open action phase with an empty stack only the active player 0 holds priority, and they can't pass (R-TURN-09)"},
		{`"phase": "action", "priority": 1`, "(R-TURN-09)"},
	} {
		_, err := buildSetup(`{` + seats + `, ` + c.extra + `}`)
		if err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s: error %v; want %q", c.extra, err, c.want)
		}
	}
	// The same seat is fine in every other window.
	for _, extra := range []string{
		`"phase": "start"`, `"phase": "loot"`, `"phase": "end"`, `"declared": "attack"`, `"declared": "purchase"`,
		`"stack": [{"dice_roll": 3, "owner": 0}]`, `"phase": "end", "loot_plays": 0`,
		`"declared": "attack", "stack": [{"dice_roll": 3, "owner": 0}]`, `"declared": "purchase", "stack": [{"dice_roll": 3, "owner": 0}]`,
	} {
		if _, err := buildSetup(`{` + seats + `, ` + extra + `, "priority": 1}`); err != nil {
			t.Errorf("%s: %v; want no error", extra, err)
		}
	}
}

// roomSet has one room, as Beggar without the roll, with the text
// "$ Pay 1¢: Gain 2¢".
var roomSet = CardSet{Name: "rooms", Cards: []CardDef{
	{Ref: "cent_room", Kind: RoomCard, Copies: 1, Abilities: []Ability{
		{Kind: Activated, Text: "Pay 1¢: Gain 2¢.", Costs: []Cost{PayCents(1)}, Effects: []Effect{GainCents(2)}},
	}},
}}

// roomTable: hero_a is active with 1¢, hero_b holds priority with 1¢;
// cent_room is in the room slot.
func roomTable(stack, steps string) string {
	return `{"sets": ["test", "rooms"],
		"setup": {"players": [{"character": "hero_a", "cents": 1}, {"character": "hero_b", "cents": 1}], "active": 0,
			"rooms": ["cent_room"], "stack": ` + stack + `},
		"steps": ` + steps + `}`
}

// A room works as an item (24.6): an activate step may name it, and
// only the active player may use it (R-CARD-33, R-CARD-36); a stacked
// room ability is the active player's, costs paid.
func TestSituationRoom(t *testing.T) {
	run := func(data string) Answer {
		t.Helper()
		ans, err := RunSituation([]byte(data), testSet, roomSet)
		if err != nil {
			t.Fatal(err)
		}
		return ans
	}
	// The room's text: pay 1¢, gain 2¢.
	ans := run(roomTable(`[]`, `[{"player": 0, "kind": "activate", "card": "cent_room",
		"expect": {"allowed": true, "after": [{"player": 0, "cents": 2}], "resolved": ["cent_room"]}}]`))
	if ans.Matches == nil || !*ans.Matches {
		t.Errorf("active player: answer %+v; want a match", ans)
	}
	ans = run(strings.Replace(roomTable(`[]`, `[{"player": 1, "kind": "activate", "card": "cent_room"}]`), `"active": 0,`, `"active": 0, "phase": "loot", "priority": 1,`, 1))
	if ans.Allowed || ans.Rule != "R-CARD-33" { // R-CARD-33: only the active player uses a room
		t.Errorf("other player: answer %+v; want a refusal by R-CARD-33", ans)
	}
	// On the stack its 1¢ is paid: the setup shows 0¢; it resolves for 2¢.
	ans = run(strings.Replace(roomTable(`[{"kind": "activate", "card": "cent_room", "owner": 0}]`,
		`[{"player": 0, "kind": "pass", "expect": {"allowed": true, "after": [{"player": 0, "cents": 2}], "resolved": ["cent_room"]}}]`),
		`"hero_a", "cents": 1`, `"hero_a"`, 1))
	if ans.Matches == nil || !*ans.Matches {
		t.Errorf("stacked: answer %+v; want a match", ans)
	}
	_, err := RunSituation([]byte(roomTable(`[{"kind": "activate", "card": "cent_room", "owner": 1}]`, `[{"kind": "pass"}]`)), testSet, roomSet)
	if err == nil || !strings.Contains(err.Error(), "owner 1 is not the active player; only they use a room's abilities (R-CARD-33)") {
		t.Errorf("stacked for another player: error %v; want the rule", err)
	}
	_, err = RunSituation([]byte(roomTable(`[]`, `[{"player": 0, "kind": "activate", "card": "hero_b"}]`)), testSet, roomSet)
	if err == nil || !strings.Contains(err.Error(), `player 0 controls no "hero_b", and no room slot holds it`) {
		t.Errorf("unknown source: error %v; want it named", err)
	}
}
