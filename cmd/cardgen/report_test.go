package main

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"

	carddb "github.com/FourSoulsRepo/card_db"
	"github.com/FourSoulsRepo/rules_engine/cards"
)

func TestCardStatus(t *testing.T) {
	dir := t.TempDir()
	files := map[string]string{
		"stub.go":        "// " + StubMarker + " implement\n",
		"done.go":        "package x\n",
		"tested.go":      "package x\n",
		"tested_test.go": "package x\n",
	}
	for name, src := range files {
		if err := os.WriteFile(filepath.Join(dir, name), []byte(src), 0o600); err != nil {
			t.Fatal(err)
		}
	}
	for id, want := range map[string]Status{"stub": NotStarted, "done": Implemented, "tested": Tested, "missing": NotStarted} {
		got, err := cardStatus(dir, carddb.Ref{ID: id, Version: 1})
		if err != nil || got != want {
			t.Errorf("%s: %v, %v; want %v", id, got, err, want)
		}
	}
}

// TestReportIsCurrent fails when a card changed but the report was not
// regenerated: go run ./cmd/cardgen -report
func TestReportIsCurrent(t *testing.T) {
	db, err := carddb.Load(os.DirFS("../../pkg/card_db"))
	if err != nil {
		t.Fatal(err)
	}
	want, err := report(db, "../../pkg/rules_engine/cards")
	if err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(filepath.Join("../..", reportFile))
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got, want) {
		t.Fatalf("%s is out of date; run: go run ./cmd/cardgen -report", reportFile)
	}
}

// TestThreePlusMatchesCardDB: the engine leaves out the cards card_db
// marks with the 3+ player symbol (R-CARD-32), and no others.
func TestThreePlusMatchesCardDB(t *testing.T) {
	db, err := carddb.Load(os.DirFS("../../pkg/card_db"))
	if err != nil {
		t.Fatal(err)
	}
	for _, set := range cards.Sets() {
		for _, d := range set.Cards {
			ref, err := carddb.ParseRef(string(d.Ref))
			if err != nil {
				t.Fatal(err)
			}
			c, ok := db.Get(ref)
			if !ok {
				continue
			}
			if c.ThreePlus != d.ThreePlus {
				t.Errorf("%s: card_db three_plus %v, engine ThreePlus %v", d.Ref, c.ThreePlus, d.ThreePlus)
			}
		}
	}
}
