package instance

import (
	"errors"
	"path/filepath"
	"testing"
)

func TestSlots(t *testing.T) {
	dir := t.TempDir()
	a, err := Acquire(dir)
	if err != nil {
		t.Fatal(err)
	}
	b, err := Acquire(dir)
	if err != nil {
		t.Fatal(err)
	}
	if a.N != 1 || b.N != 2 {
		t.Fatalf("slots %d, %d; want 1, 2", a.N, b.N)
	}
	if got, want := a.File("settings.json"), filepath.Join(dir, "settings.json"); got != want {
		t.Errorf("slot 1 file %q, want %q", got, want)
	}
	if got, want := b.File("settings.json"), filepath.Join(dir, "settings-2.json"); got != want {
		t.Errorf("slot 2 file %q, want %q", got, want)
	}
	// A closed copy frees its slot for the next one.
	if err = a.Release(); err != nil {
		t.Fatal(err)
	}
	c, err := Acquire(dir)
	if err != nil {
		t.Fatal(err)
	}
	if c.N != 1 {
		t.Errorf("after release got slot %d, want 1", c.N)
	}
}

func TestBusy(t *testing.T) {
	dir := t.TempDir()
	for range MaxSlots {
		if _, err := Acquire(dir); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := Acquire(dir); !errors.Is(err, ErrBusy) {
		t.Errorf("err %v, want ErrBusy", err)
	}
}
