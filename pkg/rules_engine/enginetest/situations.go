package enginetest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/fs"
	"reflect"
	"regexp"
	"slices"
	"strings"
	"sync"
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/docs"
)

// SituationDir is where a package keeps its situation files (ADR 034).
const SituationDir = "testdata/situations"

// Tree is the kind of a situation tree: which folders it may hold.
type Tree int

const (
	// CardTree is a set package's tree: a folder per card ref of its sets,
	// and interactions/ for several cards together.
	CardTree Tree = iota
	// RuleTree is the engine's tree: a folder per rule topic in lower
	// case, as R-ATK is atk/.
	RuleTree
	// InteractionTree is the cross-set tree in cards/: interactions/
	// only, so a card's own files stay in its set's tree (19.238).
	InteractionTree
)

// RunSituations runs a card tree: every file under testdata/situations
// in files, one subtest each, named by its path in the tree. A set
// package embeds its tree and calls it from a test (ADR 034):
//
//	//go:embed testdata/situations
//	var situations embed.FS
//
//	func TestSituations(t *testing.T) { enginetest.RunSituations(t, situations, Set) }
func RunSituations(t *testing.T, files fs.FS, sets ...engine.CardSet) {
	t.Helper()
	run(t, CardTree, files, sets)
}

// RunInteractionSituations runs the cross-set tree, as RunSituations;
// it holds only interactions/.
func RunInteractionSituations(t *testing.T, files fs.FS, sets ...engine.CardSet) {
	t.Helper()
	run(t, InteractionTree, files, sets)
}

// RunRuleSituations runs the engine's rule tree, as RunSituations.
func RunRuleSituations(t *testing.T, files fs.FS, sets ...engine.CardSet) {
	t.Helper()
	run(t, RuleTree, files, sets)
}

func run(t *testing.T, tree Tree, files fs.FS, sets []engine.CardSet) {
	t.Helper()
	paths, err := SituationFiles(files)
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range paths {
		t.Run(strings.TrimSuffix(p, ".json"), func(t *testing.T) {
			for _, problem := range CheckSituation(tree, files, p, sets...) {
				t.Error(problem)
			}
		})
	}
}

// SituationFiles lists every file under testdata/situations in files,
// by its path in the tree: "razor_blade/hits-a-player.json".
func SituationFiles(files fs.FS) ([]string, error) {
	var out []string
	err := fs.WalkDir(files, SituationDir, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			out = append(out, strings.TrimPrefix(p, SituationDir+"/"))
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("situation files: %w", err)
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("no situation files in %s", SituationDir)
	}
	return out, nil
}

// CheckSituation checks one file of a tree, by its path in the tree:
// its place in the layout, its source and its expectation. It returns
// every problem; none means it passes.
func CheckSituation(tree Tree, files fs.FS, path string, sets ...engine.CardSet) []string {
	if problem := checkLayout(tree, path, sets); problem != "" {
		return []string{path + ": " + problem}
	}
	data, err := fs.ReadFile(files, SituationDir+"/"+path)
	if err != nil {
		return []string{err.Error()}
	}
	// A test file is read strictly: a misspelt field is an error, not
	// a field left out. A payload is read as before (rule 2 of step 24).
	if bad := unknownFields(data, reflect.TypeFor[engine.Situation]()); len(bad) > 0 {
		for i, b := range bad {
			bad[i] = path + ": " + b
		}
		return append(bad, path+": the fields are in docs/wiki/Situation-payload.md")
	}
	var s engine.Situation
	dec := json.NewDecoder(bytes.NewReader(data))
	dec.DisallowUnknownFields()
	if err = dec.Decode(&s); err != nil {
		return []string{fmt.Sprintf("%s: %v; the fields are in docs/wiki/Situation-payload.md", path, err)}
	}
	var out []string
	if s.Expect == nil {
		out = append(out, path+": a test file needs expect")
	}
	if len(s.Source) == 0 {
		out = append(out, path+": a test file needs source: the card, rule or ruling its expect comes from")
	}
	for _, src := range s.Source {
		if problem := checkSource(src, sets); problem != "" {
			out = append(out, fmt.Sprintf("%s: source %q: %s", path, src, problem))
		}
	}
	if s.Expect == nil {
		return out
	}
	ans, err := engine.RunSituation(data, sets...)
	if err != nil {
		return append(out, fmt.Sprintf("%s: %v", path, err))
	}
	if s.Action == nil {
		// Steps: each mismatch names its step, or the end (24.4).
		for _, m := range ans.Mismatch {
			out = append(out, path+", "+m)
		}
	} else {
		step := strings.TrimSpace(fmt.Sprintf("step 1 (%s %s", s.Action.Kind, s.Action.Card)) + ")"
		if ans.Rule != "" {
			step += fmt.Sprintf(", refused by %s: %s", ans.Rule, ans.Reason)
		}
		for _, m := range ans.Mismatch {
			out = append(out, fmt.Sprintf("%s, %s: %s", path, step, m))
		}
	}
	// A test file uses every roll and pick it lists (24.32); a payload
	// may leave some, so only this runner reads Unused.
	for _, u := range ans.Unused {
		out = append(out, path+", end: "+u)
	}
	return out
}

// caseName is a file name in kebab case: "pays-5-cents.json".
var caseName = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*\.json$`)

// checkLayout checks a path against ADR 034 items 5 and 6: a folder,
// then a .json file in kebab case, no deeper levels.
func checkLayout(tree Tree, path string, sets []engine.CardSet) string {
	parts := strings.Split(path, "/")
	if len(parts) != 2 {
		return "outside the layout: want <folder>/<case>.json, no deeper levels (ADR 034)"
	}
	if !caseName.MatchString(parts[1]) {
		return fmt.Sprintf("%q: want a .json file named in kebab case, e.g. pays-5-cents.json", parts[1])
	}
	folder := parts[0]
	switch tree {
	case CardTree:
		if folder != "interactions" && !hasCard(sets, folder) {
			return fmt.Sprintf("folder %q: want a card ref of sets %v, or interactions", folder, setNames(sets))
		}
	case InteractionTree:
		if folder != "interactions" {
			return fmt.Sprintf("folder %q: the cross-set tree holds only interactions; a card's files go in its set's tree", folder)
		}
	case RuleTree:
		topics, err := ruleTopics()
		if err != nil {
			return err.Error()
		}
		if !slices.Contains(topics, folder) {
			return fmt.Sprintf("folder %q: want a rule topic, one of %v", folder, topics)
		}
	}
	return ""
}

var (
	ruleForm     = regexp.MustCompile(`^R-[A-Z]+-\d+$`)
	questionForm = regexp.MustCompile(`^(Q-\d+|R-07 #\d+)$`)
	linkForm     = regexp.MustCompile(`^https://\S+$`)
)

// checkSource checks one source item against ADR 034 item 9.
func checkSource(src string, sets []engine.CardSet) string {
	switch {
	case strings.HasPrefix(src, "card "):
		if ref := strings.TrimPrefix(src, "card "); !hasCard(sets, ref) {
			return fmt.Sprintf("no card %q in sets %v", ref, setNames(sets))
		}
	case ruleForm.MatchString(src):
		ids, err := ruleIDs()
		if err != nil {
			return err.Error()
		}
		if !ids[src] {
			return "no such rule in pkg/rules_engine/docs/rules"
		}
	case questionForm.MatchString(src):
		ids, err := questions()
		if err != nil {
			return err.Error()
		}
		if !ids[src] {
			return "no such question or old ruling in open-questions.md"
		}
	case linkForm.MatchString(src):
	default:
		return "want card <ref>, R-<TOPIC>-<NN>, Q-<NN>, R-07 #<N> or an https:// link"
	}
	return ""
}

var (
	ruleIDs   = sync.OnceValues(docs.RuleIDs)
	questions = sync.OnceValues(docs.Questions)
	// ruleTopics are the topics of the rule IDs in lower case, sorted.
	ruleTopics = sync.OnceValues(func() ([]string, error) {
		ids, err := ruleIDs()
		if err != nil {
			return nil, err
		}
		var out []string
		for id := range ids {
			topic := strings.ToLower(strings.Split(id, "-")[1])
			if !slices.Contains(out, topic) {
				out = append(out, topic)
			}
		}
		slices.Sort(out)
		return out, nil
	})
)

// hasCard reports whether one of the sets defines or reprints ref.
func hasCard(sets []engine.CardSet, ref string) bool {
	for _, set := range sets {
		for _, c := range set.Cards {
			if string(c.Ref) == ref {
				return true
			}
		}
		for _, r := range set.Reprints {
			if string(r.Ref) == ref {
				return true
			}
		}
	}
	return false
}

func setNames(sets []engine.CardSet) []string {
	out := make([]string, len(sets))
	for i, set := range sets {
		out[i] = set.Name
	}
	return out
}
