package carddb

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"unicode"
	"unicode/utf8"
)

// jsonError explains a decode error of a set file for translators (16.11):
// the file, line and column, a plain hint, and the line with a caret
// under the place. The hints match frontend/src/jsonError.ts.
func jsonError(name string, data []byte, err error) error {
	var syn *json.SyntaxError
	var typ *json.UnmarshalTypeError
	var at int
	var hint string
	switch {
	case errors.As(err, &syn):
		// Offset counts the bytes read, the wrong one too.
		off := int(syn.Offset) - 1
		if strings.Contains(syn.Error(), "unexpected end") {
			off = len(data)
		}
		at, hint = explainJSON(data, off)
	case errors.As(err, &typ):
		at = max(0, min(int(typ.Offset), len(data))-1)
		hint = typeHint(typ)
	default:
		return fmt.Errorf("%s: %w", name, err)
	}
	line, col := lineColumn(data, at)
	text := lineAt(data, line)
	pad := strings.Map(func(r rune) rune {
		if r == '\t' {
			return r
		}
		return ' '
	}, string([]rune(text)[:min(col-1, utf8.RuneCountInString(text))]))
	num := fmt.Sprint(line)
	return fmt.Errorf("%s:%d:%d: JSON error: %s\n  %s | %s\n  %s | %s^",
		name, line, col, hint, num, text, strings.Repeat(" ", len(num)), pad)
}

// explainJSON says what is wrong at off, the first wrong byte, and
// which byte to show: the comma to delete, or where a comma is missing.
func explainJSON(data []byte, off int) (int, string) {
	switch {
	case len(bytes.TrimSpace(data)) == 0:
		return 0, "the file is empty"
	case bytes.HasPrefix(data, []byte("\uFEFF")):
		return 0, "the file starts with an invisible byte order mark (BOM); save it as UTF-8 without BOM"
	case off >= len(data):
		return len(data), "the file ends too early: a closing }, ] or quote is missing"
	}
	off = max(off, 0)
	for off > 0 && isWord(data[off]) && isWord(data[off-1]) {
		off--
	}
	before := off - 1
	for before >= 0 && isSpace(data[before]) {
		before--
	}
	c := data[off]
	var prev byte
	if before >= 0 {
		prev = data[before]
	}
	switch {
	case before < 0:
		return off, "a set file starts with { and ends with }"
	case json.Valid(data[:before+1]):
		return off, "more text after the last }: delete it, or move it inside"
	case (c == '}' || c == ']') && prev == ',':
		return before, fmt.Sprintf("a comma before %c: delete this comma; the last entry has none", c)
	case c == '\'':
		return off, `JSON takes double quotes " only, not single quotes '`
	case c == '\t':
		return off, `a tab inside a text: write \t, or a space`
	case c < ' ':
		return off, `a text runs to the end of the line: its closing quote " is missing (a line break inside a text is \n)`
	case (c == '"' || c == '{' || c == '[') && (prev == '"' || prev == '}' || prev == ']' || isWord(prev)):
		return before + 1, "a comma is missing here, between two entries"
	case isWord(c) && prev == '"':
		return off, `a quote is missing, or a quote inside a text lacks its backslash: write \" for a " in a text`
	case isWord(c) && (prev == '{' || prev == ','):
		return off, `a key needs double quotes: "key": "text"`
	case isWord(c) && (prev == ':' || prev == '['):
		return off, `a text needs double quotes: "key": "text"`
	}
	r, _ := utf8.DecodeRune(data[off:])
	return off, fmt.Sprintf("%q does not fit here", string(r))
}

// typeHint names a field that holds the wrong kind of value.
func typeHint(e *json.UnmarshalTypeError) string {
	field := e.Field
	if i := strings.LastIndex(field, "."); i >= 0 {
		field = field[i+1:]
	}
	want := "another kind of value"
	if e.Type != nil {
		switch e.Type.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Float32, reflect.Float64:
			want = "a number, without quotes"
		case reflect.String:
			want = `a text in double quotes`
		case reflect.Bool:
			want = "true or false"
		case reflect.Slice, reflect.Array:
			want = "a list [ … ]"
		case reflect.Map, reflect.Struct:
			want = "an object { … }"
		default:
		}
	}
	return fmt.Sprintf("%s must be %s, not %s", field, want, e.Value)
}

// isWord reports a byte of a bare word: a letter, a digit, or UTF-8.
func isWord(b byte) bool {
	return b >= utf8.RuneSelf || b == '_' || b == '$' || unicode.IsLetter(rune(b)) || unicode.IsDigit(rune(b))
}

func isSpace(b byte) bool {
	return b == ' ' || b == '\t' || b == '\n' || b == '\r'
}

// lineColumn gives the 1-based line and column of a byte; the column
// counts characters, not bytes.
func lineColumn(data []byte, off int) (line, col int) {
	head := data[:min(off, len(data))]
	start := bytes.LastIndexByte(head, '\n') + 1
	return bytes.Count(head, []byte("\n")) + 1, utf8.RuneCount(head[start:]) + 1
}

// lineAt is a 1-based line, without its line break.
func lineAt(data []byte, line int) string {
	lines := strings.Split(string(data), "\n")
	if line > len(lines) {
		return ""
	}
	return strings.TrimSuffix(lines[line-1], "\r")
}
