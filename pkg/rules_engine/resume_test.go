package rulesengine

import (
	"bytes"
	"errors"
	"flag"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"
)

// savedGame is a game 40 steps in, saved.
func savedGame(t *testing.T) []byte {
	t.Helper()
	g, _ := recordGame(t, Setup{Seed: 9, Players: 3, Sets: []CardSet{testSet, abilitySet}}, 40)
	data, err := g.Save()
	if err != nil {
		t.Fatal(err)
	}
	return data
}

// TestResumeSameState: a game saved by an engine with the same state
// goes on (N-20).
func TestResumeSameState(t *testing.T) {
	data := savedGame(t)
	g, err := Resume(data, StateMeaning, testSet, abilitySet)
	if err != nil {
		t.Fatal(err)
	}
	if again, _ := g.Save(); !bytes.Equal(again, data) { //nolint:errcheck // compared below
		t.Error("the resumed game saves differently")
	}
}

// TestResumeRenamedField: a field the new engine calls something else
// is lost on load; the game cannot go on (N-20).
func TestResumeRenamedField(t *testing.T) {
	data := savedGame(t)
	old := bytes.Replace(data, []byte(`"max_hand":`), []byte(`"hand_limit":`), 1)
	if bytes.Equal(old, data) {
		t.Fatal("no max_hand in the saved game")
	}
	if _, err := Resume(old, StateMeaning, testSet, abilitySet); !errors.Is(err, ErrIncompatible) {
		t.Fatalf("err %v, want ErrIncompatible", err)
	}
}

// TestResumeNewField: an old game lacks the keys of fields added since.
// With omitempty a zero field writes no key either, and the game goes
// on; without it, the zero value appears and the game cannot (N-20).
func TestResumeNewField(t *testing.T) {
	data := savedGame(t)
	if bytes.Contains(data, []byte(`"compost"`)) {
		t.Fatal("compost is set; pick a field this game leaves zero")
	}
	if _, err := Resume(data, StateMeaning, testSet, abilitySet); err != nil {
		t.Errorf("no key for a zero omitempty field: %v", err)
	}
	old := bytes.Replace(data, []byte(`"win_souls":4,`), nil, 1)
	if bytes.Equal(old, data) {
		t.Fatal("no win_souls in the saved game")
	}
	if _, err := Resume(old, StateMeaning, testSet, abilitySet); !errors.Is(err, ErrIncompatible) {
		t.Errorf("no key for a field without omitempty: err %v, want ErrIncompatible", err)
	}
}

// TestResumeMeaning: a game saved before StateMeaning was raised ends.
func TestResumeMeaning(t *testing.T) {
	if _, err := Resume(savedGame(t), StateMeaning-1, testSet, abilitySet); !errors.Is(err, ErrIncompatible) {
		t.Fatalf("err %v, want ErrIncompatible", err)
	}
}

var updateStateFields = flag.Bool("update-state-fields", false, "rewrite testdata/state_fields.txt")

// TestStateFieldsOmitEmpty: every saved field without omitempty is in
// testdata/state_fields.txt. A new field needs omitempty (omitzero for a
// struct), so games saved before it still resume after a server restart
// (N-20). If its zero value is wrong for those games, raise StateMeaning
// and add the field to the list (go test -run StateFields -update-state-fields).
func TestStateFieldsOmitEmpty(t *testing.T) {
	var got []string
	seen := map[reflect.Type]bool{}
	var walk func(reflect.Type)
	walk = func(typ reflect.Type) {
		switch typ.Kind() {
		case reflect.Pointer, reflect.Slice, reflect.Array, reflect.Map:
			walk(typ.Elem())
			return
		case reflect.Struct:
		default:
			return
		}
		if seen[typ] || typ.PkgPath() != reflect.TypeFor[Game]().PkgPath() {
			return
		}
		seen[typ] = true
		for f := range typ.Fields() {
			tag := f.Tag.Get("json")
			if !f.IsExported() || tag == "-" {
				continue
			}
			opts := strings.Split(tag, ",")[1:]
			ok := slices.Contains(opts, "omitzero") || (slices.Contains(opts, "omitempty") && f.Type.Kind() != reflect.Struct)
			if !ok {
				got = append(got, typ.Name()+"."+f.Name)
			}
			walk(f.Type)
		}
	}
	walk(reflect.TypeFor[Game]())
	slices.Sort(got)
	path := filepath.Join("testdata", "state_fields.txt")
	if *updateStateFields {
		if err := os.WriteFile(path, []byte(strings.Join(got, "\n")+"\n"), 0o600); err != nil {
			t.Fatal(err)
		}
		return
	}
	data, err := os.ReadFile(path) //nolint:gosec // a testdata file
	if err != nil {
		t.Fatal(err)
	}
	want := strings.Fields(string(data))
	for _, f := range got {
		if !slices.Contains(want, f) {
			t.Errorf("%s has no omitempty (omitzero for a struct): games saved before it will not resume (N-20)", f)
		}
	}
	for _, f := range want {
		if !slices.Contains(got, f) {
			t.Errorf("%s is in %s but no longer a field without omitempty; remove it", f, path)
		}
	}
}
