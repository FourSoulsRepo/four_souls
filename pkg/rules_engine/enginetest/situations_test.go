package enginetest_test

import (
	"strings"
	"testing"
	"testing/fstest"

	"github.com/FourSoulsRepo/rules_engine/cards/b2"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

// razorFile is a b2 test file: Razor Blade ("↷: Deal 1 damage to a
// player.") on Cain. source and expect are filled in per case.
func razorFile(source, expect string) string {
	return `{
		"title": "Razor Blade on Cain",` + source + `
		"sets": ["b2"],
		"setup": {"players": [{"character": "isaac", "items": ["razor_blade"]}, {"character": "cain"}], "active": 0},
		"action": {"player": 0, "kind": "activate", "card": "razor_blade", "choices": ["player 2 (cain)"]}` + expect + `
	}`
}

// check runs one file of a tree at path, with b2, and returns its problems.
func check(tree enginetest.Tree, path, data string) []string {
	files := fstest.MapFS{enginetest.SituationDir + "/" + path: {Data: []byte(data)}}
	return enginetest.CheckSituation(tree, files, path, b2.Set)
}

const (
	cardSource = `"source": ["card razor_blade"],`
	// The card text: Cain takes 1 damage.
	rightExpect = `, "expect": {"allowed": true, "after": [{"player": 1, "damage": 1}]}`
)

func TestCheckSituationPasses(t *testing.T) {
	if got := check(enginetest.CardTree, "razor_blade/hits-cain.json", razorFile(cardSource, rightExpect)); len(got) != 0 {
		t.Errorf("problems %q; want none", got)
	}
}

// A wrong expect fails, naming the file, the step and every check.
func TestCheckSituationWrongExpect(t *testing.T) {
	wrong := `, "expect": {"allowed": true, "after": [{"player": 1, "damage": 2}, {"player": 0, "damage": 1}]}`
	got := check(enginetest.CardTree, "razor_blade/hits-cain.json", razorFile(cardSource, wrong))
	if len(got) != 2 {
		t.Fatalf("problems %q; want two failed checks", got)
	}
	for i, want := range []string{"player 1 damage is 1, expected 2", "player 0 damage is 0, expected 1"} {
		for _, part := range []string{"razor_blade/hits-cain.json", "step 1 (activate razor_blade)", want} {
			if !strings.Contains(got[i], part) {
				t.Errorf("problem %q; want %q in it", got[i], part)
			}
		}
	}
}

// A test file needs source and expect; a payload needs neither, but a
// file in a tree is a test file.
func TestCheckSituationNeedsSourceAndExpect(t *testing.T) {
	for _, c := range []struct{ source, expect, want string }{
		{"", rightExpect, "needs source"},
		{`"source": [],`, rightExpect, "needs source"},
		{cardSource, "", "needs expect"},
	} {
		got := check(enginetest.CardTree, "razor_blade/hits-cain.json", razorFile(c.source, c.expect))
		if len(got) != 1 || !strings.Contains(got[0], c.want) {
			t.Errorf("source %q, expect %q: problems %q; want one with %q", c.source, c.expect, got, c.want)
		}
	}
}

func TestCheckSituationSourceForms(t *testing.T) {
	for _, good := range []string{"card razor_blade", "R-ATK-01", "Q-07", "R-07 #4", "https://example.com/ruling"} {
		src := `"source": ["` + good + `"],`
		if got := check(enginetest.CardTree, "razor_blade/hits-cain.json", razorFile(src, rightExpect)); len(got) != 0 {
			t.Errorf("source %q: problems %q; want none", good, got)
		}
	}
	for _, bad := range []string{"card nope", "R-ATK-99", "Q-999", "R-07 #99", "the FAQ", "http://example.com", "R-atk-01"} {
		src := `"source": ["card razor_blade", "` + bad + `"],`
		got := check(enginetest.CardTree, "razor_blade/hits-cain.json", razorFile(src, rightExpect))
		if len(got) != 1 || !strings.Contains(got[0], bad) {
			t.Errorf("source %q: problems %q; want one naming it", bad, got)
		}
	}
}

// A file outside the layout fails before it runs (ADR 033 items 5, 6).
func TestCheckSituationLayout(t *testing.T) {
	data := razorFile(cardSource, rightExpect)
	for _, c := range []struct {
		tree enginetest.Tree
		path string
		ok   bool
	}{
		{enginetest.CardTree, "razor_blade/hits-cain.json", true},
		{enginetest.CardTree, "interactions/razor-blade-cain.json", true},
		{enginetest.CardTree, "hits-cain.json", false},
		{enginetest.CardTree, "razor_blade/more/hits-cain.json", false},
		{enginetest.CardTree, "razor_blade/hits_cain.json", false},
		{enginetest.CardTree, "razor_blade/Hits-Cain.json", false},
		{enginetest.CardTree, "razor_blade/hits-cain.txt", false},
		{enginetest.CardTree, "nope/hits-cain.json", false},
		{enginetest.CardTree, "atk/hits-cain.json", false},
		{enginetest.RuleTree, "abil/hits-cain.json", true},
		{enginetest.RuleTree, "attack/hits-cain.json", false},
		{enginetest.RuleTree, "razor_blade/hits-cain.json", false},
		{enginetest.RuleTree, "interactions/hits-cain.json", false},
	} {
		got := check(c.tree, c.path, data)
		if c.ok != (len(got) == 0) {
			t.Errorf("%s: problems %q; want ok %v", c.path, got, c.ok)
		}
	}
}

func TestSituationFiles(t *testing.T) {
	files := fstest.MapFS{
		enginetest.SituationDir + "/razor_blade/hits-cain.json": {},
		enginetest.SituationDir + "/razor_blade/more/deep.json": {},
		"other/file.json": {},
	}
	got, err := enginetest.SituationFiles(files)
	if err != nil || len(got) != 2 || got[0] != "razor_blade/hits-cain.json" || got[1] != "razor_blade/more/deep.json" {
		t.Errorf("files %q, %v; want the two under %s", got, err, enginetest.SituationDir)
	}
	if _, err := enginetest.SituationFiles(fstest.MapFS{}); err == nil {
		t.Error("no error for a missing tree")
	}
}

// A test file is read strictly: a misspelt field fails, naming it; a
// payload outside the tree is read as before.
func TestCheckSituationUnknownField(t *testing.T) {
	data := strings.Replace(razorFile(cardSource, rightExpect), `"active": 0`, `"active": 0, "rols": [6]`, 1)
	got := check(enginetest.CardTree, "razor_blade/hits-cain.json", data)
	if len(got) != 2 || !strings.Contains(got[0], `setup: unknown field "rols"; did you mean "rolls"?`) || !strings.Contains(got[1], "Situation-payload.md") {
		t.Errorf("problems %q; want the field, its place, the nearest name and the field list", got)
	}
}

// An unknown field names its line and path; with no near name, the
// fields of that place are listed.
func TestCheckSituationUnknownFieldPlace(t *testing.T) {
	data := "{\n\"title\": \"t\",\n\"setup\": {\"players\": [{}, {\n\"centz\": 1}]},\n\"zzz\": 1}"
	got := check(enginetest.CardTree, "razor_blade/hits-cain.json", data)
	want := []string{
		`line 4, setup.players[1]: unknown field "centz"; did you mean "cents"?`,
		`line 5, top level: unknown field "zzz"; fields: title, sets, seed, setup, action, steps, expect, source`,
	}
	if len(got) != 3 {
		t.Fatalf("problems %q; want two fields and the field list", got)
	}
	for i, w := range want {
		if !strings.HasSuffix(got[i], w) {
			t.Errorf("problem %q; want it to end with %q", got[i], w)
		}
	}
}

// razorSteps is razorFile with steps: Razor Blade on Cain, then the
// end of the turn; step is the first step's extra fields.
func razorSteps(step, expect string) string {
	return `{
		"title": "Razor Blade on Cain, then end the turn",` + cardSource + `
		"sets": ["b2"],
		"setup": {"players": [{"character": "isaac", "items": ["razor_blade"]}, {"character": "cain"}], "active": 0},
		"steps": [
			{"player": 0, "kind": "activate", "card": "razor_blade", "choices": ["player 2 (cain)"]` + step + `},
			{"player": 0, "kind": "end_turn"}
		]` + expect + `
	}`
}

// A step's mismatch names the file and the step; the end's names the end.
func TestCheckSituationStepMismatch(t *testing.T) {
	// The card text: Cain takes 1 damage. R-TURN-13: it heals as the turn passes.
	right := razorSteps(`, "expect": {"allowed": true, "after": [{"player": 1, "damage": 1}]}`, `, "expect": {"allowed": true, "after": [{"player": 1, "damage": 0}]}`)
	if got := check(enginetest.CardTree, "razor_blade/hits-cain.json", right); len(got) != 0 {
		t.Fatalf("problems %q; want none", got)
	}
	wrong := razorSteps(`, "expect": {"allowed": true, "after": [{"player": 1, "damage": 2}]}`, `, "expect": {"allowed": true, "after": [{"player": 1, "damage": 1}]}`)
	got := check(enginetest.CardTree, "razor_blade/hits-cain.json", wrong)
	want := []string{
		"razor_blade/hits-cain.json, step 1 (activate razor_blade): player 1 damage is 1, expected 2",
		"razor_blade/hits-cain.json, end: player 1 damage is 0, expected 1",
	}
	if len(got) != len(want) {
		t.Fatalf("problems %q; want %q", got, want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("problem %q; want %q", got[i], want[i])
		}
	}
}

// The strict reader walks steps too: a typo names the step's path.
func TestCheckSituationStepUnknownField(t *testing.T) {
	data := razorSteps(`, "choises": []`, `, "expect": {"allowed": true}`)
	data = strings.Replace(data, `{"player": 0, "kind": "end_turn"}`, `{"player": 0, "kind": "end_turn", "zzz": 1}`, 1)
	got := check(enginetest.CardTree, "razor_blade/hits-cain.json", data)
	want := []string{
		`steps[0]: unknown field "choises"; did you mean "choices"?`,
		`steps[1]: unknown field "zzz"; fields: player, kind, card, ability, choices, expect`,
	}
	if len(got) != 3 {
		t.Fatalf("problems %q; want two fields and the field list", got)
	}
	for i, w := range want {
		if !strings.HasSuffix(got[i], w) {
			t.Errorf("problem %q; want it to end with %q", got[i], w)
		}
	}
}

// The strict reader covers the 23.5 checks: a typo in a player, a
// monster or a counter check names its place and the nearest field.
func TestCheckSituationExpectUnknownField(t *testing.T) {
	expect := `, "expect": {"allowed": true, "after": [{"player": 1, "hand_card": []}],
		"monsters": [{"card": "fly", "evasoin": 2}], "counters": [{"card": "fly", "cont": 1}], "winer": [0]}`
	got := check(enginetest.CardTree, "razor_blade/hits-cain.json", razorFile(cardSource, expect))
	want := []string{
		`expect.after[0]: unknown field "hand_card"; did you mean "hand_cards"?`,
		`expect.monsters[0]: unknown field "evasoin"; did you mean "evasion"?`,
		`expect.counters[0]: unknown field "cont"; did you mean "count"?`,
		`expect: unknown field "winer"; did you mean "winners"?`,
	}
	if len(got) != len(want)+1 {
		t.Fatalf("problems %q; want %d fields and the field list", got, len(want))
	}
	for i, w := range want {
		if !strings.HasSuffix(got[i], w) {
			t.Errorf("problem %q; want it to end with %q", got[i], w)
		}
	}
}

// The strict reader covers the 23.6 checks: prompt, resolved, events.
func TestCheckSituationOrderUnknownField(t *testing.T) {
	expect := `, "expect": {"allowed": true, "prompt": {"player": 0, "optons": []},
		"resolved": ["razor_blade"], "events": [{"kind": "damaged", "plyer": 1}], "resolve": []}`
	got := check(enginetest.CardTree, "razor_blade/hits-cain.json", razorFile(cardSource, expect))
	want := []string{
		`expect.prompt: unknown field "optons"; did you mean "options"?`,
		`expect.events[0]: unknown field "plyer"; did you mean "player"?`,
		`expect: unknown field "resolve"; did you mean "resolved"?`,
	}
	if len(got) != len(want)+1 {
		t.Fatalf("problems %q; want %d fields and the field list", got, len(want))
	}
	for i, w := range want {
		if !strings.HasSuffix(got[i], w) {
			t.Errorf("problem %q; want it to end with %q", got[i], w)
		}
	}
}
