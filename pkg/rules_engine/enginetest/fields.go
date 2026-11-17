package enginetest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

// unknownFields lists every field in data that type t has no place for:
// its line, its path in the file, and the nearest known field or the
// list of fields there. A syntax error ends the walk; the strict decode
// reports it.
func unknownFields(data []byte, t reflect.Type) []string {
	w := fieldWalker{dec: json.NewDecoder(bytes.NewReader(data)), data: data}
	w.value(t, "")
	return w.out
}

var anyType = reflect.TypeFor[any]()

type fieldWalker struct {
	dec  *json.Decoder
	data []byte
	out  []string
	err  error
}

// value walks one JSON value that should fit type t.
func (w *fieldWalker) value(t reflect.Type, path string) {
	for t.Kind() == reflect.Pointer {
		t = t.Elem()
	}
	tok, err := w.dec.Token()
	if err != nil {
		w.err = err
		return
	}
	switch tok {
	case json.Delim('['):
		elem := anyType
		if t.Kind() == reflect.Slice || t.Kind() == reflect.Array {
			elem = t.Elem()
		}
		for i := 0; w.err == nil && w.dec.More(); i++ {
			w.value(elem, fmt.Sprintf("%s[%d]", path, i))
		}
		w.close()
	case json.Delim('{'):
		for w.err == nil && w.dec.More() {
			w.member(t, path)
		}
		w.close()
	}
}

// member walks one key and its value inside an object of type t.
func (w *fieldWalker) member(t reflect.Type, path string) {
	tok, err := w.dec.Token()
	if err != nil {
		w.err = err
		return
	}
	key, ok := tok.(string)
	if !ok { // the decoder only gives string keys
		w.err = fmt.Errorf("object key %v is not a string", tok)
		return
	}
	at := w.dec.InputOffset()
	inner := key
	if path != "" {
		inner = path + "." + key
	}
	switch t.Kind() {
	case reflect.Struct:
		ft, ok := fieldType(t, key)
		if !ok {
			w.unknown(t, key, path, at)
			ft = anyType
		}
		w.value(ft, inner)
	case reflect.Map:
		w.value(t.Elem(), inner)
	default:
		w.value(anyType, inner)
	}
}

func (w *fieldWalker) close() {
	if _, err := w.dec.Token(); err != nil {
		w.err = err
	}
}

// unknown records a field t lacks, with the nearest name it has.
func (w *fieldWalker) unknown(t reflect.Type, key, path string, at int64) {
	line := 1 + bytes.Count(w.data[:at], []byte("\n"))
	where := path
	if where == "" {
		where = "top level"
	}
	names := fieldNames(t)
	hint := "fields: " + strings.Join(names, ", ")
	if near := nearest(key, names); near != "" {
		hint = fmt.Sprintf("did you mean %q?", near)
	}
	w.out = append(w.out, fmt.Sprintf("line %d, %s: unknown field %q; %s", line, where, key, hint))
}

// fieldType finds the struct field JSON maps key to; like encoding/json,
// it ignores case.
func fieldType(t reflect.Type, key string) (reflect.Type, bool) {
	for f := range t.Fields() {
		if name := jsonName(f); name != "" && strings.EqualFold(name, key) {
			return f.Type, true
		}
	}
	return nil, false
}

func fieldNames(t reflect.Type) []string {
	var names []string
	for f := range t.Fields() {
		if name := jsonName(f); name != "" {
			names = append(names, name)
		}
	}
	return names
}

// jsonName is a field's JSON name, or "" for one JSON skips.
func jsonName(f reflect.StructField) string {
	if !f.IsExported() {
		return ""
	}
	name, _, _ := strings.Cut(f.Tag.Get("json"), ",")
	switch name {
	case "-":
		return ""
	case "":
		return f.Name
	}
	return name
}

// nearest is the name within two edits of key, or "".
func nearest(key string, names []string) string {
	best, bestDist := "", 3
	for _, n := range names {
		if d := editDistance(strings.ToLower(key), n); d < bestDist {
			best, bestDist = n, d
		}
	}
	return best
}

// editDistance is the Levenshtein distance between a and b.
func editDistance(a, b string) int {
	prev := make([]int, len(b)+1)
	for j := range prev {
		prev[j] = j
	}
	for i := 1; i <= len(a); i++ {
		cur := make([]int, len(b)+1)
		cur[0] = i
		for j := 1; j <= len(b); j++ {
			cost := 1
			if a[i-1] == b[j-1] {
				cost = 0
			}
			cur[j] = min(prev[j]+1, cur[j-1]+1, prev[j-1]+cost)
		}
		prev = cur
	}
	return prev[len(b)]
}
