package rulesengine

import (
	"encoding/json"
	"reflect"
	"testing"
)

// TestArgsJSONMatchesStdlib checks the typed writer against json.Marshal
// for every value kind and for strings that need escapes (18.17).
func TestArgsJSONMatchesStdlib(t *testing.T) {
	strs := []string{
		"", "plain", `quote " and \ back`, "<b>&</b>", "tab\there\nnew\rcr",
		"\b\f\x00\x01\x1f", "  and  ", "bad \xff utf8", "héllo wörld 🂡",
	}
	for _, s := range strs {
		args := Args{Str("s", s), CardArg("c", CardRef(s)), N("n", -42), PlayerArg("p", 3)}
		got, err := args.MarshalJSON()
		if err != nil {
			t.Fatal(err)
		}
		want, err := json.Marshal(map[string]any{"s": s, "c": CardValue{Card: CardRef(s)}, "n": -42, "p": PlayerValue{Player: 3}})
		if err != nil {
			t.Fatal(err)
		}
		// The map sorts its keys; compare value by value instead.
		var a, b map[string]any
		if err := json.Unmarshal(got, &a); err != nil {
			t.Fatalf("%q: not JSON: %s", s, got)
		}
		if err := json.Unmarshal(want, &b); err != nil {
			t.Fatal(err)
		}
		for k := range b {
			if !reflect.DeepEqual(a[k], b[k]) {
				t.Errorf("%q: key %s: got %v, want %v", s, k, a[k], b[k])
			}
		}
		for _, v := range []any{s, CardValue{Card: CardRef(s)}, -42, PlayerValue{Player: 3}} {
			one, err := appendValue(nil, v)
			if err != nil {
				t.Fatal(err)
			}
			std, err := json.Marshal(v)
			if err != nil {
				t.Fatal(err)
			}
			if string(one) != string(std) {
				t.Errorf("%#v: got %s, want %s", v, one, std)
			}
		}
	}
}
