package main

import (
	"os"
	"slices"
	"testing"
	"testing/fstest"

	"github.com/FourSoulsRepo/four_souls/internal/settings"
)

func TestDiceMatchSettings(t *testing.T) {
	var ids []string
	for _, d := range diceSkins {
		ids = append(ids, d.id)
	}
	if !slices.Equal(ids, settings.Dice) {
		t.Fatalf("dice %v, settings %v", ids, settings.Dice)
	}
}

func TestDiceImage(t *testing.T) {
	a := &App{cardFiles: fstest.MapFS{"images/" + diceSkins[0].file: {Data: []byte("PNG")}}}
	if got := a.Dice()[0].Image; got != "/cardimg/"+diceSkins[0].file {
		t.Fatalf("image %q", got)
	}
	// Without the private images every die is drawn plain.
	a = &App{cardFiles: os.DirFS(t.TempDir())}
	for _, d := range a.Dice() {
		if d.Image != "" || d.Name == "" {
			t.Fatalf("%+v without images", d)
		}
	}
}
