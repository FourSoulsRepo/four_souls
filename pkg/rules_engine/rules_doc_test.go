package rulesengine

import (
	"embed"
	"fmt"
	"io/fs"
	"regexp"
	"sort"
	"strconv"
	"testing"
)

//go:embed docs/rules/*.md
var rulesDocs embed.FS

var (
	// ruleDef is a rule definition: "**R-STACK-03**".
	ruleDef = regexp.MustCompile(`\*\*(R-[A-Z]+-\d+)\*\*`)
	// ruleRef is any mention of a rule ID.
	ruleRef = regexp.MustCompile(`\bR-[A-Z]+-\d+\b`)
	ruleID  = regexp.MustCompile(`^R-([A-Z]+)-(\d{2})$`)
)

func TestRuleIDs(t *testing.T) {
	files, err := fs.Glob(rulesDocs, "docs/rules/*.md")
	if err != nil || len(files) == 0 {
		t.Fatalf("no rules files: %v", err)
	}
	defined := map[string]string{} // id -> file
	refs := map[string][]string{}  // id -> files that mention it
	numbers := map[string][]int{}  // topic -> numbers
	for _, f := range files {
		data, err := fs.ReadFile(rulesDocs, f)
		if err != nil {
			t.Fatal(err)
		}
		for _, m := range ruleDef.FindAllStringSubmatch(string(data), -1) {
			id := m[1]
			parts := ruleID.FindStringSubmatch(id)
			if parts == nil {
				t.Errorf("%s: bad rule id %s (want R-TOPIC-NN)", f, id)
				continue
			}
			if prev, dup := defined[id]; dup {
				t.Errorf("%s: %s already defined in %s", f, id, prev)
			}
			defined[id] = f
			n, err := strconv.Atoi(parts[2])
			if err != nil {
				t.Errorf("%s: bad number in %s: %v", f, id, err)
				continue
			}
			numbers[parts[1]] = append(numbers[parts[1]], n)
		}
		for _, id := range ruleRef.FindAllString(string(data), -1) {
			refs[id] = append(refs[id], f)
		}
	}
	for id, from := range refs {
		if _, ok := defined[id]; !ok {
			t.Errorf("%v: reference to undefined rule %s", from, id)
		}
	}
	for topic, ns := range numbers {
		sort.Ints(ns)
		for i, n := range ns {
			if n != i+1 {
				t.Errorf("%s: rules are not numbered 01..%02d without gaps (found %s)", topic, len(ns), fmt.Sprint(ns))
				break
			}
		}
	}
	t.Logf("%d rules in %d topics", len(defined), len(numbers))
}
