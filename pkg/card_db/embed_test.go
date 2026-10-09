//go:build embed

package carddb

import (
	"io/fs"
	"testing"
)

func TestEmbedded(t *testing.T) {
	fsys, ok := Embedded()
	if !ok {
		t.Fatal("Embedded reports false in an embed build")
	}
	db, err := Load(fsys)
	if err != nil {
		t.Fatal(err)
	}
	if len(db.Refs()) == 0 {
		t.Fatal("no cards embedded")
	}
	// The submodule's git data must never be packed into the binary.
	if _, err := fs.Stat(fsys, "images/.git"); err == nil {
		t.Error("images/.git is embedded")
	}
	with := 0
	for _, ref := range db.Refs() {
		c, _ := db.Get(ref)
		if HasImage(fsys, c) {
			with++
		}
	}
	t.Logf("%d cards, %d with embedded images", len(db.Refs()), with)
}
