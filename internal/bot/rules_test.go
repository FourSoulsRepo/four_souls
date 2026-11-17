package bot

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards"
)

// enginePath is the one module the bot's code may import besides the
// standard library.
const enginePath = "github.com/FourSoulsRepo/rules_engine"

// cloneFuncs are the functions that may copy a game (Game.Clone, a call
// with no arguments). Each copies only a world from Input.World (ADR
// 030, Decision 4.4); the search of step 23.5 copies in fork.
var cloneFuncs = map[string]bool{"fork": true}

// keyFuncs are the functions that may read a card's identity (Card),
// and only to look it up in the measured table of card values, which is
// generated data (ADR 033, Decision 20.9): tableKey.
var keyFuncs = map[string]bool{"tableKey": true}

// banned are engine names the bot's code never uses: card identities
// and names (rule 5.24, ADR 033 Decision 6.3), and ways to make a game
// other than Input.World.
var banned = map[string]string{
	"CardRef":      "a card's identity",
	"Card":         "a card's identity",
	"CardOf":       "a card's identity",
	"Ref":          "a card's identity",
	"Labels":       "card names",
	"Plain":        "card names",
	"PlainOptions": "card names",
	"OptionIndex":  "card names",
	"NewGame":      "a game not from Input.World",
	"Load":         "a game not from Input.World",
	"Imagine":      "a game not from Input.World",
}

// TestBotRules checks the bot's code, test files aside (ADR 033,
// Decision 4; step 23 rules 3 and 4):
//   - it imports the standard library and the engine only: never
//     internal/server, internal/world or a card set;
//   - it names no card: no CardRef, no card field (Card only in
//     keyFuncs, to look a card up in the generated table of values), no
//     string that is a card's ref, and strings only in constants and
//     error messages;
//   - it makes no game: Clone only in cloneFuncs, no NewGame, Load or
//     Imagine; no package variable holds a game;
//   - Input holds no game, only World.
func TestBotRules(t *testing.T) {
	refs := cardRefs()
	files, err := filepath.Glob("*.go")
	if err != nil {
		t.Fatal(err)
	}
	fset := token.NewFileSet()
	checked := 0
	for _, name := range files {
		if strings.HasSuffix(name, "_test.go") {
			continue
		}
		src, err := os.ReadFile(name) //nolint:gosec // a file of this package
		if err != nil {
			t.Fatal(err)
		}
		f, err := parser.ParseFile(fset, name, src, 0)
		if err != nil {
			t.Fatal(err)
		}
		checked++
		checkImports(t.Errorf, fset, f)
		checkNames(t.Errorf, fset, f, refs)
		checkGlobals(t.Errorf, fset, f)
	}
	if checked == 0 {
		t.Fatal("no bot source files found")
	}
	checkNoGame(t, reflect.TypeFor[Input](), "Input")
}

// cardRefs lists the refs of every card the engine ships.
func cardRefs() map[string]bool {
	refs := map[string]bool{}
	sets := cards.Sets()
	for _, c := range cards.Challenges() {
		if s, ok := cards.ChallengeSet(c); ok {
			sets = append(sets, s)
		}
	}
	for _, s := range sets {
		for _, d := range s.Cards {
			refs[string(d.Ref)] = true
		}
	}
	return refs
}

func checkImports(errorf reporter, fset *token.FileSet, f *ast.File) {
	for _, im := range f.Imports {
		path, err := strconv.Unquote(im.Path.Value)
		if err != nil {
			errorf("%v", err)
			continue
		}
		std := !strings.Contains(strings.Split(path, "/")[0], ".")
		if !std && path != enginePath {
			errorf("%s: imports %s; the bot imports the standard library and the engine only", fset.Position(im.Pos()), path)
		}
	}
}

// checkNames walks every function and constant of f.
func checkNames(errorf reporter, fset *token.FileSet, f *ast.File, refs map[string]bool) {
	fn := ""
	var visit func(n ast.Node, allowStrings bool)
	visit = func(n ast.Node, allowStrings bool) {
		ast.Inspect(n, func(n ast.Node) bool {
			switch x := n.(type) {
			case *ast.FuncDecl:
				if x.Body != nil {
					fn = x.Name.Name
					visit(x.Body, false)
				}
				return false
			case *ast.GenDecl:
				if x.Tok == token.IMPORT {
					return false
				}
				if x.Tok == token.CONST {
					for _, s := range x.Specs {
						visit(s, true)
					}
					return false
				}
			case *ast.CallExpr:
				if isErrorCall(x) {
					for _, a := range x.Args {
						visit(a, true)
					}
					visit(x.Fun, allowStrings)
					return false
				}
				if sel, ok := x.Fun.(*ast.SelectorExpr); ok && sel.Sel.Name == "Clone" && len(x.Args) == 0 && !cloneFuncs[fn] {
					errorf("%s: Clone in %s; only %v copy a game, and only a world", fset.Position(x.Pos()), fn, cloneFuncs)
				}
			case *ast.SelectorExpr:
				if why, ok := banned[x.Sel.Name]; ok && (x.Sel.Name != "Card" || !keyFuncs[fn]) {
					errorf("%s: uses %s: %s", fset.Position(x.Pos()), x.Sel.Name, why)
				}
			case *ast.BasicLit:
				if x.Kind != token.STRING {
					return true
				}
				s, err := strconv.Unquote(x.Value)
				if err != nil {
					errorf("%v", err)
					return false
				}
				if refs[s] {
					errorf("%s: the string %q is a card", fset.Position(x.Pos()), s)
				}
				if !allowStrings && s != "" {
					errorf("%s: the string %q outside a constant or an error; the bot names no card", fset.Position(x.Pos()), s)
				}
			}
			return true
		})
	}
	visit(f, false)
}

// isErrorCall: fmt.Errorf or errors.New, whose strings are messages.
func isErrorCall(c *ast.CallExpr) bool {
	sel, ok := c.Fun.(*ast.SelectorExpr)
	if !ok {
		return false
	}
	pkg, ok := sel.X.(*ast.Ident)
	return ok && (pkg.Name == "fmt" && sel.Sel.Name == "Errorf" || pkg.Name == "errors" && sel.Sel.Name == "New")
}

// checkGlobals: no package variable mentions engine.Game.
func checkGlobals(errorf reporter, fset *token.FileSet, f *ast.File) {
	for _, d := range f.Decls {
		g, ok := d.(*ast.GenDecl)
		if !ok || g.Tok != token.VAR {
			continue
		}
		ast.Inspect(g, func(n ast.Node) bool {
			if sel, ok := n.(*ast.SelectorExpr); ok && sel.Sel.Name == "Game" {
				errorf("%s: a package variable holds a game", fset.Position(sel.Pos()))
			}
			return true
		})
	}
}

// checkNoGame fails when type typ holds an *engine.Game outside a func
// (World returns one; it holds none).
func checkNoGame(t *testing.T, typ reflect.Type, path string) {
	t.Helper()
	game := reflect.TypeFor[engine.Game]()
	switch typ.Kind() { //nolint:exhaustive // only containers can hold a game
	case reflect.Pointer, reflect.Slice, reflect.Array, reflect.Map:
		if typ.Elem() == game {
			t.Errorf("%s holds a game", path)
			return
		}
		checkNoGame(t, typ.Elem(), path)
	case reflect.Struct:
		if typ == game {
			t.Errorf("%s holds a game", path)
			return
		}
		for i := range typ.NumField() {
			fl := typ.Field(i)
			checkNoGame(t, fl.Type, path+"."+fl.Name)
		}
	}
}

// reporter takes one broken rule.
type reporter func(format string, args ...any)

// TestBotRulesCatch checks that the rules catch what they must, on a
// piece of bad code: each line of want is broken once at least.
func TestBotRulesCatch(t *testing.T) {
	src := `package bot
import (
	"github.com/FourSoulsRepo/four_souls/internal/server"
	engine "github.com/FourSoulsRepo/rules_engine"
)
var real *engine.Game
func peek(g *engine.Game) bool {
	c, _ := g.Clone()
	return c.Objects[0].Card == engine.CardRef("isaac") || g.Def(0).Ref == "x"
}
func tableKey(g *engine.Game) string {
	return string(g.Def(0).Ref) + string(g.Objects[0].Card)
}`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "bad.go", src, 0)
	if err != nil {
		t.Fatal(err)
	}
	var got []string
	rec := func(format string, args ...any) { got = append(got, fmt.Sprintf(format, args...)) }
	checkImports(rec, fset, f)
	checkNames(rec, fset, f, map[string]bool{"isaac": true})
	checkGlobals(rec, fset, f)
	all := strings.Join(got, "\n")
	for _, want := range []string{
		"imports github.com/FourSoulsRepo/four_souls/internal/server",
		"Clone in peek",
		"uses Card:",
		"uses CardRef:",
		"uses Ref:",
		`the string "isaac" is a card`,
		`the string "x" outside`,
		"a package variable holds a game",
	} {
		if !strings.Contains(all, want) {
			t.Errorf("no %q among:\n%s", want, all)
		}
	}
	// In tableKey (line 12) Card is allowed, Ref is not.
	ref, card := false, false
	for _, line := range got {
		if strings.HasPrefix(line, "bad.go:12:") {
			ref = ref || strings.Contains(line, "uses Ref:")
			card = card || strings.Contains(line, "uses Card:")
		}
	}
	if !ref || card {
		t.Errorf("tableKey: Ref caught %v, Card caught %v; want true, false:\n%s", ref, card, all)
	}
}
