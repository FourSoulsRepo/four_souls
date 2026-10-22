package rulesengine

import (
	"flag"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"testing"
)

// evName is the name of an event kind: EvDied, not EventCard.
var evName = regexp.MustCompile(`^Ev[A-Z]`)

var updateEventKinds = flag.Bool("update-event-kinds", false, "rewrite testdata/event_kinds.txt")

// TestEventKindsListed: testdata/event_kinds.txt lists every EventKind.
// The client's history reads the list and fails on a kind it neither
// tells nor skips with a reason (15.88, history.test.ts). After adding a
// kind: go test -run EventKinds -update-event-kinds.
func TestEventKindsListed(t *testing.T) {
	fset := token.NewFileSet()
	files, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	for _, name := range files {
		if strings.HasSuffix(name, "_test.go") {
			continue
		}
		f, perr := parser.ParseFile(fset, name, nil, 0)
		if perr != nil {
			t.Fatal(perr)
		}
		for _, d := range f.Decls {
			gd, ok := d.(*ast.GenDecl)
			if !ok || gd.Tok != token.CONST {
				continue
			}
			for _, spec := range gd.Specs {
				vs, ok := spec.(*ast.ValueSpec)
				if !ok {
					continue
				}
				typ, ok := vs.Type.(*ast.Ident)
				isKind := ok && typ.Name == "EventKind"
				for i, n := range vs.Names {
					if !isKind {
						if evName.MatchString(n.Name) {
							t.Errorf("%s: %s needs the type EventKind written out, so the list sees it", name, n.Name)
						}
						continue
					}
					lit, ok := vs.Values[i].(*ast.BasicLit)
					if !ok {
						t.Fatalf("%s: %s is not a string literal", name, n.Name)
					}
					v, uerr := strconv.Unquote(lit.Value)
					if uerr != nil {
						t.Fatal(uerr)
					}
					got = append(got, v)
				}
			}
		}
	}
	slices.Sort(got)
	path := filepath.Join("testdata", "event_kinds.txt")
	if *updateEventKinds {
		if werr := os.WriteFile(path, []byte(strings.Join(got, "\n")+"\n"), 0o600); werr != nil {
			t.Fatal(werr)
		}
		return
	}
	data, err := os.ReadFile(path) //nolint:gosec // a testdata file
	if err != nil {
		t.Fatal(err)
	}
	if want := strings.Fields(string(data)); !slices.Equal(got, want) {
		t.Errorf("%s is out of date; run go test -run EventKinds -update-event-kinds\ngot  %v\nwant %v", path, got, want)
	}
}

// TestDeathNamesTheKiller: a monster's and a player's death event name
// who killed them, for the history (15.88).
func TestDeathNamesTheKiller(t *testing.T) {
	g := toAction(t, 2)
	a := g.Turn.Active
	putMonster(g, 0, "gaper")
	g.forcedRolls = []int{6, 3} // both hit: DC 3 (R-ATK-12)
	attack(t, g)
	choose(t, g, "gaper")
	var ev []Event
	for range 200 {
		ev = append(ev, throwPending(t, g)...)
		if g.Prompt().Kind != PromptPriority || (!g.Attack.On && len(g.Stack) == 0) {
			break
		}
		ev = append(ev, pass(t, g)...)
	}
	i := slices.IndexFunc(ev, func(e Event) bool { return e.Kind == EvDied && e.Card == "gaper" })
	if i < 0 {
		t.Fatalf("no died event for the gaper in %+v", ev)
	}
	if got := ev[i].By; got != int(a)+1 {
		t.Errorf("died By %d, want the attacker %d plus 1", got, a)
	}
}
