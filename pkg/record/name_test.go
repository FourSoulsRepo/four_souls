package record

import (
	"strings"
	"testing"
	"time"
	"unicode/utf8"
)

func TestSafeName(t *testing.T) {
	for in, want := range map[string]string{
		"Ann":                         "Ann",
		"  Big  Bob  ":                "Big-Bob",
		"../../etc/passwd":            "etc-passwd",
		`C:\Windows\x`:                "C-Windows-x",
		"a\x00b\nc\td":                "a-b-c-d",
		"evil\u202egpj.exe":           "evil-gpj-exe",
		"Иван_Петров":                 "Иван-Петров",
		"<>:\"|?*":                    "",
		"..":                          "",
		"CON":                         "CON",
		"x" + strings.Repeat("y", 40): "x" + strings.Repeat("y", 19),
		"aaaaaaaaaaaaaaaaaaa b":       "aaaaaaaaaaaaaaaaaaa",
	} {
		if got := SafeName(in); got != want {
			t.Errorf("SafeName(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestFileName(t *testing.T) {
	at := time.Date(2026, 10, 10, 12, 0, 0, 0, time.UTC)
	got := FileName(at, "g3", []string{"Ann", "../Bob", "***"})
	if want := "2026-10-10-120000-g3_Ann_Bob_p3.fsrec"; got != want {
		t.Fatalf("FileName = %q, want %q", got, want)
	}
	long := strings.Repeat("ж", 20)
	got = FileName(at, "g12", []string{long, long, long, long})
	if utf8.RuneCountInString(got) > maxFile+len(Ext) || len(got) > 255 {
		t.Fatalf("FileName is too long: %d runes, %d bytes", utf8.RuneCountInString(got), len(got))
	}
	if SafeFileName(got) != got {
		t.Fatalf("SafeFileName changed a server name: %q", SafeFileName(got))
	}
}

func TestSafeFileName(t *testing.T) {
	for in, want := range map[string]string{
		"2026-10-10-120000-g3_Ann_Bob.fsrec": "2026-10-10-120000-g3_Ann_Bob.fsrec",
		"../../.bashrc":                      "bashrc.fsrec",
		`..\..\evil.fsrec`:                   "evil.fsrec",
		"run.sh":                             "run-sh.fsrec",
		"/":                                  "",
		"___":                                "",
		"":                                   "",
	} {
		if got := SafeFileName(in); got != want {
			t.Errorf("SafeFileName(%q) = %q, want %q", in, got, want)
		}
	}
}
