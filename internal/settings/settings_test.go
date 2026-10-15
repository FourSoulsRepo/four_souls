package settings

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestSurvivesRestart(t *testing.T) {
	path := filepath.Join(t.TempDir(), "FourSouls", "settings.json")
	saved, err := Save(path, Settings{Nickname: "  Ann  ", Recent: []string{"10.0.0.2"}})
	if err != nil {
		t.Fatal(err)
	}
	if saved.Nickname != "Ann" {
		t.Fatalf("saved nickname %q, want Ann", saved.Nickname)
	}
	// A restart only has the file.
	if got := Load(path); !reflect.DeepEqual(got, saved) {
		t.Fatalf("loaded %+v, want %+v", got, saved)
	}
}

func TestDefaults(t *testing.T) {
	dir := t.TempDir()
	want := Settings{Recent: []string{}, Animation: "normal", Theme: "dark", Dice: "tear"}
	if got := Load(filepath.Join(dir, "none.json")); !reflect.DeepEqual(got, want) {
		t.Fatalf("missing file: %+v", got)
	}
	broken := filepath.Join(dir, "broken.json")
	if err := os.WriteFile(broken, []byte("{nick"), 0o600); err != nil {
		t.Fatal(err)
	}
	if got := Load(broken); !reflect.DeepEqual(got, want) {
		t.Fatalf("broken file: %+v", got)
	}
}

func TestBadNicknameKeepsFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "settings.json")
	if _, err := Save(path, Settings{Nickname: "Ann"}); err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"   ", strings.Repeat("x", 21), "a\nb"} {
		if _, err := Save(path, Settings{Nickname: name}); !errors.Is(err, ErrNickname) {
			t.Fatalf("%q: err %v, want ErrNickname", name, err)
		}
	}
	if got := Load(path).Nickname; got != "Ann" {
		t.Fatalf("nickname %q after bad saves, want Ann", got)
	}
}

func TestRecentIsCleanedAndCapped(t *testing.T) {
	path := filepath.Join(t.TempDir(), "settings.json")
	got, err := Save(path, Settings{Recent: []string{" a ", "", "b", "c", "d", "e", "f"}})
	if err != nil {
		t.Fatal(err)
	}
	if want := []string{"a", "b", "c", "d", "e"}; !reflect.DeepEqual(got.Recent, want) {
		t.Fatalf("recent %q, want %q", got.Recent, want)
	}
}

func TestAnimationSpeed(t *testing.T) {
	path := filepath.Join(t.TempDir(), "settings.json")
	got, err := Save(path, Settings{Animation: "fast"})
	if err != nil || got.Animation != "fast" {
		t.Fatalf("fast: %+v, %v", got, err)
	}
	if _, err := Save(path, Settings{Animation: "warp"}); !errors.Is(err, ErrChoice) {
		t.Fatalf("warp: err %v, want ErrChoice", err)
	}
	if got := Load(path).Animation; got != "fast" {
		t.Fatalf("animation %q after a bad save, want fast", got)
	}
	if err := os.WriteFile(path, []byte(`{"animation":"warp"}`), 0o600); err != nil {
		t.Fatal(err)
	}
	if got := Load(path).Animation; got != "normal" {
		t.Fatalf("unknown speed reads as %q, want normal", got)
	}
}

func TestWindowThemeAndMat(t *testing.T) {
	path := filepath.Join(t.TempDir(), "settings.json")
	want := Settings{Recent: []string{}, Animation: "normal", Fullscreen: true, Theme: "light", Mat: "cathedral", Dice: "tear"}
	got, err := Save(path, Settings{Fullscreen: true, Theme: "light", Mat: "cathedral"})
	if err != nil || !reflect.DeepEqual(got, want) {
		t.Fatalf("saved %+v, %v; want %+v", got, err, want)
	}
	if got := Load(path); !reflect.DeepEqual(got, want) {
		t.Fatalf("loaded %+v, want %+v", got, want)
	}
	for _, bad := range []Settings{{Theme: "pink"}, {Mat: "../lava"}, {Mat: "Lava"}, {Dice: "d20"}} {
		if _, err := Save(path, bad); !errors.Is(err, ErrChoice) {
			t.Fatalf("%+v: err %v, want ErrChoice", bad, err)
		}
	}
}

// A mat image may be added later, so any mat file name is kept; the
// table shows colors while its image is missing.
func TestImageMatKept(t *testing.T) {
	path := filepath.Join(t.TempDir(), "settings.json")
	if got, err := Save(path, Settings{Mat: "dark_room"}); err != nil || got.Mat != "dark_room" {
		t.Fatalf("saved %+v, %v", got, err)
	}
	if got := Load(path).Mat; got != "dark_room" {
		t.Fatalf("loaded mat %q", got)
	}
	if err := os.WriteFile(path, []byte(`{"mat":"../x"}`), 0o600); err != nil {
		t.Fatal(err)
	}
	if got := Load(path).Mat; got != "" {
		t.Fatalf("bad mat reads as %q, want the default", got)
	}
}

func TestRejoinSurvivesRestart(t *testing.T) {
	path := filepath.Join(t.TempDir(), "settings.json")
	r := Rejoin{Address: "192.168.1.5", Token: "t0k3n", Game: "g1"}
	if _, err := Save(path, Settings{Rejoin: r}); err != nil {
		t.Fatal(err)
	}
	if got := Load(path).Rejoin; got != r {
		t.Fatalf("rejoin %+v, want %+v", got, r)
	}
}

func TestHostSetupRemembered(t *testing.T) {
	path := filepath.Join(t.TempDir(), "settings.json")
	if got := Load(path).Host; got.Seats != 0 || got.Sets != nil {
		t.Errorf("default host setup = %+v, want unset", got)
	}
	in := Load(path)
	in.Host = HostSetup{Seats: 3, Sets: []string{"b2", "g2", "g2", "Bad Code"}, Detailed: true, ResponseTimer: 90, NoBonusSouls: true}
	if _, err := Save(path, in); err != nil {
		t.Fatal(err)
	}
	got := Load(path).Host
	want := HostSetup{Seats: 3, Sets: []string{"b2", "g2"}, Detailed: true, ResponseTimer: 90, NoBonusSouls: true}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("host setup = %+v, want %+v", got, want)
	}
	in.Host = HostSetup{Seats: 9, ResponseTimer: 7}
	if out, err := Save(path, in); err != nil || out.Host.Seats != 0 || out.Host.ResponseTimer != 0 {
		t.Errorf("bad values: %+v, %v; want them dropped", out.Host, err)
	}
}
