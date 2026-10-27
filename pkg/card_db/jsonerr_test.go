package carddb

import (
	"strings"
	"testing"
	"testing/fstest"
)

// TestJSONError: a broken set file names the line and a hint (16.11).
// The cases match frontend/src/jsonError.test.ts.
func TestJSONError(t *testing.T) {
	tests := []struct {
		name, data, where, hint string
	}{
		{"an empty file", "", "1:1", "the file is empty"},
		{"spaces only", " \n", "1:1", "the file is empty"},
		{"a BOM", "\uFEFF{\"set\": {}}", "1:1", "byte order mark"},
		{"a trailing comma before }", "{\n  \"cards\": [],\n}\n", "2:14", "a comma before }: delete this comma"},
		{"a trailing comma before ]", `{"cards": [{}, ]}`, "1:14", "a comma before ]"},
		{"a missing comma between lines", "{\n  \"set\": {}\n  \"cards\": []\n}", "2:12", "a comma is missing here"},
		{"a missing comma after a number", `{"cards": [{"copies": 1 "id": "x"}]}`, "1:24", "a comma is missing here"},
		{"a missing closing quote", "{\"cards\": [{\"id\": \"x,\n\"type\": \"y\"}]}", "1:22", `closing quote " is missing`},
		{"a tab in a text", "{\"cards\": [{\"id\": \"x\ty\"}]}", "1:21", "a tab inside a text"},
		{"a quote inside a text", `{"cards": [{"id": "say "hi" now"}]}`, "1:25", `write \" for a "`},
		{"a key without quotes", "{\n  cards: []\n}", "2:3", "a key needs double quotes"},
		{"a text without quotes", `{"cards": [{"id": hello}]}`, "1:19", "a text needs double quotes"},
		{"a misspelt null", `{"cards": [{"id": nope}]}`, "1:19", "a text needs double quotes"},
		{"single quotes", `{'cards': []}`, "1:2", `double quotes " only`},
		{"the end comes too early", `{"cards": [`, "1:12", "ends too early"},
		{"text after the end", "{}\n}", "2:1", "more text after the last }"},
		{"no object", "hello", "1:1", "starts with {"},
		{"anything else", `{"cards": [],, "set": {}}`, "1:14", `"," does not fit here`},
		{"a number in quotes", "{\"cards\": [\n  {\"copies\": \"2\"}\n]}", "2:16", "copies must be a number, without quotes, not string"},
		{"an object for a list", `{"cards": {}}`, "1:11", "cards must be a list [ … ], not object"},
		{"a list for a text", `{"cards": [{"id": ["x"]}]}`, "1:", "id must be a text in double quotes"},
		{"a Cyrillic line", "{\"cards\": [{\"text\": {\"ru\": {\"name\": \"Сундук\",}}}]}", "1:45", "a comma before }"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Load(fstest.MapFS{"data/b2.json": {Data: []byte(tt.data)}})
			if err == nil {
				t.Fatal("no error")
			}
			msg := err.Error()
			if !strings.Contains(msg, "data/b2.json:"+tt.where) || !strings.Contains(msg, tt.hint) {
				t.Errorf("want %s and %q in\n%s", tt.where, tt.hint, msg)
			}
		})
	}
}

// TestJSONErrorCaret: the message shows the line with a caret under the
// place, so a translator sees it without counting columns.
func TestJSONErrorCaret(t *testing.T) {
	data := "{\n  \"set\": {\"code\": \"b2\"},\n\t\"cards\": [],\n}\n"
	_, err := Load(fstest.MapFS{"data/b2.json": {Data: []byte(data)}})
	want := "carddb: decode: data/b2.json:3:13: JSON error: a comma before }: delete this comma; the last entry has none\n" +
		"  3 | \t\"cards\": [],\n" +
		"    | \t           ^"
	if err == nil || err.Error() != want {
		t.Errorf("got\n%v\nwant\n%s", err, want)
	}
}

// TestJSONErrorEveryFile: every broken file is named at once.
func TestJSONErrorEveryFile(t *testing.T) {
	_, err := Load(fstest.MapFS{
		"data/b2.json": {Data: []byte(`{"cards": [],}`)},
		"data/g2.json": {Data: []byte(`{"cards": [] "set": {}}`)},
	})
	if err == nil || !strings.Contains(err.Error(), "data/b2.json:1:13") || !strings.Contains(err.Error(), "data/g2.json:1:13") {
		t.Errorf("want both files in %v", err)
	}
}
