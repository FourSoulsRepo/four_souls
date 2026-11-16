package protocol

import (
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strings"
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

// TestViewFieldsInProtocolTS checks that the client knows every field
// of the game update (TS-07, 22.4): each JSON field the server sends in
// an Update, the view, events and prompts included, is declared in
// frontend/src/protocol.ts. viewFields.test.ts then checks that the
// client reads each one. Lobby messages are out of scope.
func TestViewFieldsInProtocolTS(t *testing.T) {
	src, err := os.ReadFile("../../frontend/src/protocol.ts")
	if err != nil {
		t.Fatal(err)
	}
	ts := tsInterfaces(string(src))
	root, ok := ts["Update"]
	if !ok {
		t.Fatal("protocol.ts has no interface Update")
	}
	w := fieldWalk{ts: ts, seen: map[reflect.Type]bool{}}
	w.walk(reflect.TypeFor[Update](), tsMembers(root), "Update")
	if len(w.missing) > 0 {
		t.Errorf("protocol.ts lacks fields the server sends; add them (an unread one goes on the list in viewFields.test.ts):\n%s", strings.Join(w.missing, "\n"))
	}
	if w.fields < 100 {
		t.Errorf("only %d fields walked; the walk is broken", w.fields)
	}
}

// fieldWalk walks a Go type and its TS twin side by side.
type fieldWalk struct {
	ts      map[string]string // interface name → body
	seen    map[reflect.Type]bool
	missing []string
	fields  int
}

// msgType is a text the client turns into words whole (ADR 021); its
// own fields are not view data.
var msgType = reflect.TypeFor[engine.Msg]()

func (w *fieldWalk) walk(typ reflect.Type, members map[string]string, at string) {
	if w.seen[typ] {
		return
	}
	w.seen[typ] = true
	for i := range typ.NumField() {
		f := typ.Field(i)
		if !f.IsExported() {
			continue
		}
		name, _, _ := strings.Cut(f.Tag.Get("json"), ",")
		if name == "-" {
			continue
		}
		if name == "" {
			name = f.Name
		}
		w.fields++
		path := at + "." + name
		tsType, ok := members[name]
		if !ok {
			w.missing = append(w.missing, fmt.Sprintf("  %s (Go %s.%s)", path, typ.Name(), f.Name))
			continue
		}
		elem := f.Type
		for elem.Kind() == reflect.Pointer || elem.Kind() == reflect.Slice || elem.Kind() == reflect.Array {
			elem = elem.Elem()
		}
		if elem.Kind() != reflect.Struct || elem == msgType {
			continue
		}
		inner, err := w.resolve(tsType)
		if err != nil {
			w.missing = append(w.missing, fmt.Sprintf("  %s: %v", path, err))
			continue
		}
		w.walk(elem, inner, path)
	}
}

// resolve finds the members of a TS field type: an inline object, or
// the one interface of this file the type names (CardView[] | null).
func (w *fieldWalk) resolve(typ string) (map[string]string, error) {
	if i := strings.Index(typ, "{"); i >= 0 {
		end := matchBrace(typ, i)
		if end < 0 {
			return nil, fmt.Errorf("unclosed object type %q", typ)
		}
		return tsMembers(typ[i+1 : end]), nil
	}
	var found []string
	for _, id := range tsIdent.FindAllString(typ, -1) {
		if _, ok := w.ts[id]; ok {
			found = append(found, id)
		}
	}
	if len(found) != 1 {
		return nil, fmt.Errorf("TS type %q names no single interface of protocol.ts", typ)
	}
	return tsMembers(w.ts[found[0]]), nil
}

var (
	tsInterface = regexp.MustCompile(`^export interface (\w+)\b.*\{`)
	tsIdent     = regexp.MustCompile(`\b[A-Z]\w*\b`)
	tsMember    = regexp.MustCompile(`^(\w+)\??\s*:\s*(.+)$`)
	tsComment   = regexp.MustCompile(`//.*$`)
)

// tsInterfaces returns each exported interface's body, comments removed.
// It reads protocol.ts's own style: one declaration per line, braces
// that match within the interface.
func tsInterfaces(src string) map[string]string {
	out := map[string]string{}
	lines := strings.Split(src, "\n")
	for i := 0; i < len(lines); i++ {
		m := tsInterface.FindStringSubmatch(lines[i])
		if m == nil {
			continue
		}
		var body strings.Builder
		depth := 1
		for i++; i < len(lines) && depth > 0; i++ {
			line := tsComment.ReplaceAllString(lines[i], "")
			depth += strings.Count(line, "{") - strings.Count(line, "}")
			if depth > 0 {
				body.WriteString(line + "\n")
			}
		}
		i--
		out[m[1]] = body.String()
	}
	return out
}

// tsMembers splits a body into field name → type, at nesting depth 0.
func tsMembers(body string) map[string]string {
	out := map[string]string{}
	depth, start := 0, 0
	add := func(s string) {
		if m := tsMember.FindStringSubmatch(strings.TrimSpace(s)); m != nil {
			out[m[1]] = strings.TrimSpace(m[2])
		}
	}
	for i, r := range body {
		switch r {
		case '{', '(':
			depth++
		case '}', ')':
			depth--
		case ';', '\n':
			if depth == 0 {
				add(body[start:i])
				start = i + 1
			}
		}
	}
	add(body[start:])
	return out
}

// matchBrace returns the index of the brace that closes s[open].
func matchBrace(s string, open int) int {
	depth := 0
	for i := open; i < len(s); i++ {
		switch s[i] {
		case '{':
			depth++
		case '}':
			depth--
			if depth == 0 {
				return i
			}
		}
	}
	return -1
}

// TestTSMembers pins the small TS reader the check above uses.
func TestTSMembers(t *testing.T) {
	ts := tsInterfaces("export interface A {\n  x?: number; // {\n  y: { a: number; b?: B }[];\n  z: (B | null)[];\n}\n\nexport interface B {\n  c: string;\n}\n")
	a := tsMembers(ts["A"])
	if len(a) != 3 || a["x"] != "number" || a["z"] != "(B | null)[]" {
		t.Fatalf("members of A: %v", a)
	}
	w := fieldWalk{ts: ts}
	y, err := w.resolve(a["y"])
	if err != nil || len(y) != 2 || y["b"] != "B" {
		t.Fatalf("inline type: %v, %v", y, err)
	}
	z, err := w.resolve(a["z"])
	if err != nil || z["c"] != "string" {
		t.Fatalf("named type: %v, %v", z, err)
	}
}
