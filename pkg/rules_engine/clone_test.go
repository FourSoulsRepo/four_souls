package rulesengine_test

import (
	"reflect"
	"slices"
	"strings"
	"testing"
	"unsafe"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

// TestCloneCoversEveryField walks the types of Game: each struct type
// is in enginetest.CloneRules with exactly its fields, a plain field holds no list,
// map or pointer, and a marked one does.
func TestCloneCoversEveryField(t *testing.T) {
	seen := map[string]bool{}
	var walk func(rt reflect.Type)
	walk = func(rt reflect.Type) {
		for {
			switch rt.Kind() {
			case reflect.Pointer, reflect.Slice, reflect.Array:
				rt = rt.Elem()
				continue
			case reflect.Map:
				walk(rt.Key())
				rt = rt.Elem()
				continue
			default:
			}
			break
		}
		if rt.Kind() != reflect.Struct || seen[rt.Name()] {
			return
		}
		seen[rt.Name()] = true
		rule, ok := enginetest.CloneRules[rt.Name()]
		if !ok {
			t.Errorf("type %s is not in enginetest.CloneRules: copy it in Game.Clone and list its fields", rt.Name())
			return
		}
		listed := strings.Fields(rule)
		var fields []string
		for i := range rt.NumField() {
			fields = append(fields, rt.Field(i).Name)
		}
		names := make([]string, len(listed))
		for i, f := range listed {
			names[i] = strings.TrimRight(f, "*!")
		}
		if !slices.Equal(names, fields) {
			t.Errorf("%s: fields %v, enginetest.CloneRules lists %v: a new field needs a copy rule in Game.Clone and here", rt.Name(), fields, names)
			return
		}
		for i, f := range listed {
			ft := rt.Field(i).Type
			switch {
			case strings.HasSuffix(f, "!"):
			case strings.HasSuffix(f, "*"):
				if !holdsRefs(ft) {
					t.Errorf("%s.%s is a plain value: drop its *", rt.Name(), names[i])
				}
				walk(ft)
			default:
				if holdsRefs(ft) {
					t.Errorf("%s.%s holds a list, map or pointer: copy it in Game.Clone and mark it *", rt.Name(), names[i])
				}
				walk(ft)
			}
		}
	}
	walk(reflect.TypeFor[engine.Game]())
	if t.Failed() {
		return // a type left unwalked is no stale entry
	}
	for name := range enginetest.CloneRules {
		if !seen[name] {
			t.Errorf("enginetest.CloneRules lists %s, which Game no longer holds", name)
		}
	}
}

// holdsRefs says whether a value of type rt shares memory when copied.
// A string never changes, so it is a plain value.
func holdsRefs(rt reflect.Type) bool {
	switch rt.Kind() {
	case reflect.Slice, reflect.Map, reflect.Pointer, reflect.Interface, reflect.Func, reflect.Chan, reflect.UnsafePointer:
		return true
	case reflect.Array:
		return holdsRefs(rt.Elem())
	case reflect.Struct:
		for i := range rt.NumField() {
			if holdsRefs(rt.Field(i).Type) {
				return true
			}
		}
	default:
	}
	return false
}

// TestCloneCopiesEveryField fills a game so that every slice, map and
// pointer in it, at every depth, holds something, then checks the copy
// shares none of them: a field marked * that Clone forgets fails here,
// even if no simulated game ever fills it.
func TestCloneCopiesEveryField(t *testing.T) {
	var g engine.Game
	fill(reflect.ValueOf(&g).Elem())
	c, err := g.Clone()
	if err != nil {
		t.Fatal(err)
	}
	if err := enginetest.NoSharing(&g, c); err != nil {
		t.Fatal(err)
	}
}

// fill gives every slice one element, every map one entry and every
// pointer a value, through all nested types; fields marked ! stay
// empty.
func fill(v reflect.Value) {
	switch v.Kind() {
	case reflect.Slice:
		v.Set(reflect.MakeSlice(v.Type(), 1, 1))
		fill(v.Index(0))
	case reflect.Pointer:
		v.Set(reflect.New(v.Type().Elem()))
		fill(v.Elem())
	case reflect.Map:
		k := reflect.New(v.Type().Key()).Elem()
		e := reflect.New(v.Type().Elem()).Elem()
		fill(k)
		fill(e)
		v.Set(reflect.MakeMap(v.Type()))
		v.SetMapIndex(k, e)
	case reflect.Array:
		for i := range v.Len() {
			fill(v.Index(i))
		}
	case reflect.Struct:
		for i := range v.NumField() {
			if mark, _ := enginetest.CloneMark(v.Type(), v.Type().Field(i).Name); mark == "!" {
				continue
			}
			f := v.Field(i)
			if !f.CanSet() { // an unexported field, e.g. maxHP
				f = reflect.NewAt(f.Type(), unsafe.Pointer(f.UnsafeAddr())).Elem() //nolint:gosec // test only: fills unexported state fields
			}
			fill(f)
		}
	default: // a plain value stays zero
	}
}
