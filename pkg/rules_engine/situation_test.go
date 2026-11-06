package rulesengine

import (
	"bytes"
	"embed"
	"encoding/json"
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
