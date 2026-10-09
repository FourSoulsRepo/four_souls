package rulesengine

import (
	"embed"
	"io/fs"
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
