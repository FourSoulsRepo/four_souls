package website

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/FourSoulsRepo/four_souls/internal/cardinfo"
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards"
)

// The payloads in testdata are what cmd/website/smoke.mjs sends to the
// real .wasm; here they run through the same Go code.
func TestCheckPayloads(t *testing.T) {
	files, err := filepath.Glob("testdata/*.json")
	if err != nil || len(files) == 0 {
		t.Fatalf("no payloads: %v", err)
	}
	for _, f := range files {
		t.Run(f, func(t *testing.T) {
			data, err := os.ReadFile(f) //nolint:gosec // test files from Glob
			if err != nil {
				t.Fatal(err)
			}
			out := Check(string(data))
			var ans struct {
				engine.Answer
				Error string `json:"error"`
			}
			if err := json.Unmarshal([]byte(out), &ans); err != nil || ans.Error != "" {
				t.Fatalf("%s: %v", out, err)
			}
			if ans.Matches == nil || !*ans.Matches {
				t.Errorf("answer %+v; want it to match the expectation", ans)
			}
			if !ans.Allowed && ans.Say.IsZero() {
				t.Errorf("refusal without a message key: %+v", ans)
			}
		})
	}
}

func TestCheckError(t *testing.T) {
	for _, bad := range []string{"not json", `{"sets":["test"],"setup":{"players":[{"character":"isaac"}]},"action":{"kind":"pass"}}`} {
		var out struct {
			Error string `json:"error"`
		}
		if err := json.Unmarshal([]byte(Check(bad)), &out); err != nil || out.Error == "" {
			t.Errorf("%s: %+v, %v; want an error answer", bad, out, err)
		}
	}
}

func TestRecover(t *testing.T) {
	out := func() (out string) {
		defer Recover(&out)
		panic("boom")
	}()
	if !strings.Contains(out, `"error":"internal error: boom"`) {
		t.Errorf("got %s", out)
	}
}

// The site shows cards of its sets only, as text: no image URL.
func TestCards(t *testing.T) {
	var got []cardinfo.CardInfo
	if err := json.Unmarshal([]byte(Cards(os.DirFS("../../pkg/card_db"))), &got); err != nil {
		t.Fatal(err)
	}
	seen := map[string]bool{}
	for _, c := range got {
		seen[c.Set] = true
		if c.Image != "" || len(c.Images) != 0 {
			t.Fatalf("%s has images: %q %v", c.Ref, c.Image, c.Images)
		}
	}
	for _, code := range Sets {
		if !seen[code] {
			t.Errorf("no cards of set %s", code)
		}
	}
	if len(seen) != len(Sets) {
		t.Errorf("sets %v; want only %v", seen, Sets)
	}
}

func TestSetsAreEngineSets(t *testing.T) {
	for _, code := range Sets {
		if _, ok := cards.Find(code); !ok {
			t.Errorf("set %s is not an engine set", code)
		}
	}
}

func TestVersionAndNotice(t *testing.T) {
	if v := Version(); !strings.Contains(v, `"engine":"`+engine.Version+`"`) {
		t.Errorf("version %s", v)
	}
	if n := Notice(); !strings.Contains(n, "unofficial") {
		t.Errorf("notice %s", n)
	}
	if s := SetsJSON(); !strings.HasPrefix(s, `["b2"`) {
		t.Errorf("sets %s", s)
	}
}
