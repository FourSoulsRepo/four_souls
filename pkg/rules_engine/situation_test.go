package rulesengine

import (
	"bytes"
	"embed"
	"encoding/json"
	"io/fs"
	"slices"
	"testing"
)

//go:embed testdata/situations/*.json
var situationFiles embed.FS

// TestSituations runs every payload in testdata/situations as a rules test
// (LR-07): each must match its expectation.
func TestSituations(t *testing.T) {
	files, err := fs.Glob(situationFiles, "testdata/situations/*.json")
	if err != nil || len(files) == 0 {
		t.Fatalf("no situations: %v", err)
	}
	for _, f := range files {
		t.Run(f, func(t *testing.T) {
			data, err := fs.ReadFile(situationFiles, f)
			if err != nil {
				t.Fatal(err)
			}
			ans, err := RunSituation(data, testSet, abilitySet)
			if err != nil {
				t.Fatal(err)
			}
			if ans.Matches == nil {
				t.Fatal("a situation test needs an expectation")
			}
			if !*ans.Matches {
				t.Errorf("mismatch: %v (answer: allowed=%v rule=%s reason=%s)", ans.Mismatch, ans.Allowed, ans.Rule, ans.Reason)
			}
		})
	}
}

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
	data, err := fs.ReadFile(situationFiles, "testdata/situations/loot-needs-a-loot-play.json")
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

	data, err = fs.ReadFile(situationFiles, "testdata/situations/razor-pays-and-damages.json")
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
	data, err := fs.ReadFile(situationFiles, "testdata/situations/three-stacked-events.json")
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
