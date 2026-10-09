package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	carddb "github.com/FourSoulsRepo/card_db"
)

func TestVarAndFileName(t *testing.T) {
	cases := []struct {
		ref        carddb.Ref
		name, file string
	}{
		{carddb.Ref{ID: "the_d6", Version: 1}, "theD6", "the_d6.go"},
		{carddb.Ref{ID: "the_d6", Version: 2}, "theD6V2", "the_d6_v2.go"},
		{carddb.Ref{ID: "4_cents", Version: 1}, "card4Cents", "4_cents.go"},
	}
	for _, c := range cases {
		if got := varName(c.ref); got != c.name {
			t.Errorf("varName(%v) = %s, want %s", c.ref, got, c.name)
		}
		if got := fileName(c.ref); got != c.file {
			t.Errorf("fileName(%v) = %s, want %s", c.ref, got, c.file)
		}
	}
}

func TestRewardCount(t *testing.T) {
	for in, want := range map[string]int{"3x": 3, "x10": 10, "1x": 1} {
		if n, ok := rewardCount(in); !ok || n != want {
			t.Errorf("rewardCount(%q) = %d, %v", in, n, ok)
		}
	}
	for _, in := range []string{"?x", "Roll- Gain x"} {
		if _, ok := rewardCount(in); ok {
			t.Errorf("rewardCount(%q) should fail", in)
		}
	}
}

func TestGenerateBaseGame(t *testing.T) {
	db, err := carddb.Load(os.DirFS("../../pkg/card_db"))
	if err != nil {
		t.Fatal(err)
	}
	files, err := generate(db, "b2")
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != len(db.Refs())+2 { // cards, set list, doc.go
		t.Fatalf("%d files for %d cards", len(files), len(db.Refs()))
	}
	if _, err := generate(db, "zz"); err == nil {
		t.Fatal("unknown set accepted")
	}
}

func TestWriteKeepsFinishedCards(t *testing.T) {
	dir := t.TempDir()
	finished := filepath.Join(dir, "done.go")
	if err := os.WriteFile(finished, []byte("package x // done\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	stubbed := filepath.Join(dir, "stub.go")
	if err := os.WriteFile(stubbed, []byte("// "+StubMarker+" old\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	files := []stub{
		{file: "done.go", src: []byte("new")},
		{file: "stub.go", src: []byte("new")},
		{file: "fresh.go", src: []byte("new")},
		{file: setFile, src: []byte("list")},
	}
	written, kept, err := write(dir, files)
	if err != nil {
		t.Fatal(err)
	}
	if written != 2 || kept != 1 {
		t.Fatalf("written %d, kept %d; want 2, 1", written, kept)
	}
	for name, want := range map[string]string{"done.go": "done", "stub.go": "new", "fresh.go": "new", setFile: "list"} {
		got, err := os.ReadFile(filepath.Join(dir, name)) //nolint:gosec // test temp dir
		if err != nil || !strings.Contains(string(got), want) {
			t.Errorf("%s = %q, %v; want %q", name, got, err, want)
		}
	}
}
