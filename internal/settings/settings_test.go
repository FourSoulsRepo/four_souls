package settings

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/FourSoulsRepo/four_souls/internal/protocol"
)

func TestSurvivesRestart(t *testing.T) {
	path := filepath.Join(t.TempDir(), "FourSouls", "settings.json")
	saved, err := Save(path, Settings{Nickname: "  Ann  ", Recent: []RecentServer{{Address: "10.0.0.2", Name: "Basement"}}})
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
	want := Settings{Recent: []RecentServer{}, Animation: "normal", Dice: "tear"}
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
	in := []RecentServer{{Address: " a ", Name: " A "}, {Address: "", Name: "none"}, {Address: "b", Name: "x\ny"}, {Address: "c"}, {Address: "d"}, {Address: "e"}, {Address: "f"}}
	got, err := Save(path, Settings{Recent: in})
	if err != nil {
		t.Fatal(err)
	}
	want := []RecentServer{{Address: "a", Name: "A"}, {Address: "b"}, {Address: "c"}, {Address: "d"}, {Address: "e"}}
	if !reflect.DeepEqual(got.Recent, want) {
		t.Fatalf("recent %+v, want %+v", got.Recent, want)
	}
}

// TestOldRecentLoads: a file from before N-18 keeps bare addresses.
func TestOldRecentLoads(t *testing.T) {
	path := filepath.Join(t.TempDir(), "settings.json")
	if err := os.WriteFile(path, []byte(`{"recent": ["10.0.0.2", {"address": "10.0.0.3", "name": "Basement"}]}`), 0o600); err != nil {
		t.Fatal(err)
	}
	want := []RecentServer{{Address: "10.0.0.2"}, {Address: "10.0.0.3", Name: "Basement"}}
	if got := Load(path).Recent; !reflect.DeepEqual(got, want) {
		t.Fatalf("recent %+v, want %+v", got, want)
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

func TestWindowAndMat(t *testing.T) {
	path := filepath.Join(t.TempDir(), "settings.json")
	want := Settings{Recent: []RecentServer{}, Animation: "normal", Fullscreen: true, Mat: "cathedral", Dice: "tear"}
	got, err := Save(path, Settings{Fullscreen: true, Mat: "cathedral"})
	if err != nil || !reflect.DeepEqual(got, want) {
		t.Fatalf("saved %+v, %v; want %+v", got, err, want)
	}
	if got := Load(path); !reflect.DeepEqual(got, want) {
		t.Fatalf("loaded %+v, want %+v", got, want)
	}
	for _, bad := range []Settings{{Mat: "../lava"}, {Mat: "Lava"}, {Dice: "d20"}} {
		if _, err := Save(path, bad); !errors.Is(err, ErrChoice) {
			t.Fatalf("%+v: err %v, want ErrChoice", bad, err)
		}
	}
}

// TestOldThemeIgnored: a file from before ST-07 keeps a theme; the rest
// of it still loads, and the next save drops it.
func TestOldThemeIgnored(t *testing.T) {
	path := filepath.Join(t.TempDir(), "settings.json")
	if err := os.WriteFile(path, []byte(`{"nickname": "Ann", "theme": "light", "mat": "cathedral"}`), 0o600); err != nil {
		t.Fatal(err)
	}
	got := Load(path)
	want := Settings{Nickname: "Ann", Recent: []RecentServer{}, Animation: "normal", Mat: "cathedral", Dice: "tear"}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("loaded %+v, want %+v", got, want)
	}
	if _, err := Save(path, got); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(path) //nolint:gosec // the test names the file
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(data), "theme") {
		t.Fatalf("saved file keeps the theme: %s", data)
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

// TestFeltIsDefault: Green felt is gone (ST-08); a saved one loads as
// the default mat.
func TestFeltIsDefault(t *testing.T) {
	path := filepath.Join(t.TempDir(), "settings.json")
	if err := os.WriteFile(path, []byte(`{"mat":"felt"}`), 0o600); err != nil {
		t.Fatal(err)
	}
	if got := Load(path).Mat; got != "" {
		t.Fatalf("felt reads as %q, want the default", got)
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
	in.Host = HostSetup{Seats: 3, Sets: []string{"b2", "g2", "g2", "Bad Code"}, Detailed: true, ResponseTimer: 90, NoBonusSouls: true, NoDeckRatios: true, Rollback: protocol.RollbackSoft}
	if _, err := Save(path, in); err != nil {
		t.Fatal(err)
	}
	got := Load(path).Host
	want := HostSetup{Seats: 3, Sets: []string{"b2", "g2"}, Detailed: true, ResponseTimer: 90, NoBonusSouls: true, NoDeckRatios: true, Rollback: protocol.RollbackSoft}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("host setup = %+v, want %+v", got, want)
	}
	in.Host = HostSetup{Seats: 9, ResponseTimer: 7, Rollback: "chaos"}
	if out, err := Save(path, in); err != nil || out.Host.Seats != 0 || out.Host.ResponseTimer != 0 || out.Host.Rollback != "" {
		t.Errorf("bad values: %+v, %v; want them dropped", out.Host, err)
	}
}

// End turn asks by default (V-46); the switch to stop it is kept.
func TestNoEndTurnAsk(t *testing.T) {
	path := filepath.Join(t.TempDir(), "settings.json")
	if Load(path).NoEndTurnAsk {
		t.Fatal("a fresh start does not ask before ending the turn")
	}
	if _, err := Save(path, Settings{NoEndTurnAsk: true}); err != nil {
		t.Fatal(err)
	}
	if !Load(path).NoEndTurnAsk {
		t.Fatal("the switch was lost")
	}
}

// The LAN screen keeps the last typed address, trimmed (15.97).
func TestLanAddressSurvivesRestart(t *testing.T) {
	path := filepath.Join(t.TempDir(), "settings.json")
	if _, err := Save(path, Settings{LanAddress: " 192.168.1.5:4774 "}); err != nil {
		t.Fatal(err)
	}
	if got := Load(path).LanAddress; got != "192.168.1.5:4774" {
		t.Fatalf("LAN address %q, want 192.168.1.5:4774", got)
	}
	saved, err := Save(path, Settings{LanAddress: strings.Repeat("a", maxAddress+1)})
	if err != nil {
		t.Fatal(err)
	}
	if saved.LanAddress != "" {
		t.Errorf("a too long address was kept: %q", saved.LanAddress)
	}
}

// The saved seat keeps the path it was joined by (15.99).
func TestRejoinKeepsItsPath(t *testing.T) {
	path := filepath.Join(t.TempDir(), "settings.json")
	r := Rejoin{Address: "192.168.1.5", Token: "t0k3n", Game: "g1", Via: "lan"}
	if _, err := Save(path, Settings{Rejoin: r}); err != nil {
		t.Fatal(err)
	}
	if got := Load(path).Rejoin; got != r {
		t.Fatalf("rejoin %+v, want %+v", got, r)
	}
	saved, err := Save(path, Settings{Rejoin: Rejoin{Token: "t", Via: "moon"}})
	if err != nil {
		t.Fatal(err)
	}
	if saved.Rejoin.Via != "" {
		t.Errorf("an unknown path was kept: %q", saved.Rejoin.Via)
	}
}

func TestLanguageSurvivesRestart(t *testing.T) {
	path := filepath.Join(t.TempDir(), "settings.json")
	if got := Load(path).Language; got != "" {
		t.Fatalf("default language %q, want \"\" (the OS language)", got)
	}
	for _, code := range []string{"ru", "en", "qps", ""} {
		if _, err := Save(path, Settings{Language: code}); err != nil {
			t.Fatalf("%q: %v", code, err)
		}
		if got := Load(path).Language; got != code {
			t.Fatalf("language %q after a restart, want %q", got, code)
		}
	}
}

func TestBadLanguage(t *testing.T) {
	path := filepath.Join(t.TempDir(), "settings.json")
	if _, err := Save(path, Settings{Language: "ru"}); err != nil {
		t.Fatal(err)
	}
	for _, code := range []string{"RU", "russian", "r", "ru-RU", "../x"} {
		if _, err := Save(path, Settings{Language: code}); !errors.Is(err, ErrChoice) {
			t.Fatalf("%q: err %v, want ErrChoice", code, err)
		}
	}
	if got := Load(path).Language; got != "ru" {
		t.Fatalf("language %q after bad saves, want ru", got)
	}
	// A hand-edited file with a bad code follows the OS language.
	if err := os.WriteFile(path, []byte(`{"language": "Klingon"}`), 0o600); err != nil {
		t.Fatal(err)
	}
	if got := Load(path).Language; got != "" {
		t.Fatalf("bad file language %q, want \"\"", got)
	}
}
