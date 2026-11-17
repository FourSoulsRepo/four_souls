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
		`line 5, top level: unknown field "zzz"; fields: title, sets, seed, setup, action, expect, source`,
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
