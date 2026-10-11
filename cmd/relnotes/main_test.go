package main

import (
	"errors"
	"os"
	"strings"
	"testing"
)

const changelog = `# Changelog

## Unreleased

- Next things.

## [0.2.0] - 2026-11-01

- Saved games.

## 0.1.0 - 2026-10-20

- First release.
`

func TestNotes(t *testing.T) {
	got, err := notes([]byte(changelog), "0.2.0", "0.1.3")
	if err != nil {
		t.Fatal(err)
	}
	want := "**App 0.2.0 · Rules engine 0.1.3**\n\n- Saved games.\n"
	if got != want {
		t.Fatalf("notes:\n%s\nwant:\n%s", got, want)
	}
	if got, err := notes([]byte(changelog), "0.1.0", "0.1.0"); err != nil || !strings.Contains(got, "First release.") {
		t.Errorf("0.1.0 notes: %q, %v", got, err)
	}
}

func TestTrialUsesUnreleased(t *testing.T) {
	got, err := notes([]byte(changelog), "0.3.0-test", "0.1.3")
	if err != nil || !strings.Contains(got, "Next things.") {
		t.Fatalf("trial notes %q, %v", got, err)
	}
}

func TestRefuses(t *testing.T) {
	if _, err := notes([]byte(changelog), "v1.0", "0.1.0"); !errors.Is(err, ErrVersion) {
		t.Errorf("bad version: err %v", err)
	}
	if _, err := notes([]byte(changelog), "9.9.9", "0.1.0"); err == nil {
		t.Error("a release without its changelog section must fail")
	}
}

// TestChangelogHasUnreleased: the repo's changelog always has a section
// for trial releases.
func TestChangelogHasUnreleased(t *testing.T) {
	if _, err := notes(mustRead(t, "../../CHANGELOG.md"), "0.0.0-test", "0.0.0"); err != nil {
		t.Fatal(err)
	}
}

func mustRead(t *testing.T, path string) []byte {
	t.Helper()
	data, err := os.ReadFile(path) //nolint:gosec // a fixed repo file
	if err != nil {
		t.Fatal(err)
	}
	return data
}
