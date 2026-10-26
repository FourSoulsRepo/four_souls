package main

import (
	"encoding/json"
	"flag"
	"os"
	"path/filepath"
	"slices"
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

var updateKeys = flag.Bool("update", false, "rewrite frontend/src/locales/server-keys.json")

// keysFile lists every key the engine, the cards, the server and the
// app send (ADR 021, item 9.4); the client's unused-key check reads it.
var keysFile = filepath.Join("frontend", "src", "locales", "server-keys.json")

// TestServerKeys writes the key registry to keysFile with -update and
// fails while the file is stale. The app links every package that
// declares keys: the engine, every card set and the server.
func TestServerKeys(t *testing.T) {
	var keys []string
	for _, k := range engine.Keys() {
		keys = append(keys, k.String())
	}
	data, err := json.MarshalIndent(keys, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	data = append(data, '\n')
	if *updateKeys {
		if err = os.WriteFile(keysFile, data, 0o600); err != nil {
			t.Fatal(err)
		}
		return
	}
	old, err := os.ReadFile(keysFile) //nolint:gosec // a fixed file in the repo
	if err != nil {
		t.Fatal(err)
	}
	if string(old) != string(data) {
		t.Errorf("%s is stale; run go test . -run TestServerKeys -update", keysFile)
	}
}

// TestServerKeysInEnglish: every key the server may send has an English
// line in en.json (ADR 021, item 10.1). The client checks the rest.
func TestServerKeysInEnglish(t *testing.T) {
	data, err := os.ReadFile(filepath.Join("frontend", "src", "locales", "en.json"))
	if err != nil {
		t.Fatal(err)
	}
	var en map[string]json.RawMessage
	if err := json.Unmarshal(data, &en); err != nil {
		t.Fatal(err)
	}
	var missing []string
	for _, k := range engine.Keys() {
		if _, ok := en[k.String()]; !ok {
			missing = append(missing, k.String())
		}
	}
	slices.Sort(missing)
	if len(missing) > 0 {
		t.Errorf("keys with no line in en.json: %q", missing)
	}
}
