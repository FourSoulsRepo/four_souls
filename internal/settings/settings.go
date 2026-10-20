// Package settings keeps the player's app settings in a JSON file in
// the user's config folder, so they survive a restart.
package settings

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"slices"
	"strings"

	"github.com/FourSoulsRepo/four_souls/internal/protocol"
)

// Settings are the player's app settings.
type Settings struct {
	// Nickname is sent when joining a game (ST-06); "" until set.
	Nickname string `json:"nickname"`
	// Recent are the last servers joined, newest first.
	Recent []RecentServer `json:"recent"`
	// Animation is the animation speed (ST-05): one of Speeds.
	Animation string `json:"animation"`
	// Fullscreen starts the window fullscreen (ST-01).
	Fullscreen bool `json:"fullscreen"`
	// Theme colours menus and panels (ST-02): one of Themes.
	Theme string `json:"theme"`
	// Mat is the game mat under the table (ST-03, V-24): a mat image
	// name or one of ColorMats; see ValidMat. "" is the default: the
	// first mat image, or the first color mat without images.
	Mat string `json:"mat"`
	// Dice is the look of the player's dice (V-19): one of Dice.
	Dice string `json:"dice"`
	// NoEndTurnAsk ends the turn at once, without asking while plays are
	// left (V-46); false, the default, asks.
	NoEndTurnAsk bool `json:"no_end_turn_ask"`
	// Rejoin is the last game joined, so a restarted app takes its seat
	// back (N-08).
	Rejoin Rejoin `json:"rejoin"`
	// Host is the last game setup the player hosted, for the next one.
	// The password is never kept (N-14).
	Host HostSetup `json:"host"`
}

// RecentServer is a server joined lately (N-18).
type RecentServer struct {
	Address string `json:"address"`
	// Name is the server's name from its welcome; "" shows the address.
	Name string `json:"name"`
}

// UnmarshalJSON also reads a bare address, as older files keep them.
func (r *RecentServer) UnmarshalJSON(data []byte) error {
	var address string
	if json.Unmarshal(data, &address) == nil {
		*r = RecentServer{Address: address}
		return nil
	}
	type plain RecentServer
	var p plain
	if err := json.Unmarshal(data, &p); err != nil {
		return fmt.Errorf("recent server: %w", err)
	}
	*r = RecentServer(p)
	return nil
}

// HostSetup is a remembered host screen.
type HostSetup struct {
	Seats int `json:"seats"` // 2–4; 0: not set yet
	// Sets are the chosen card set codes; nil: not set yet (all sets).
	Sets          []string `json:"sets"`
	Detailed      bool     `json:"detailed"`
	ResponseTimer int      `json:"response_timer"` // seconds, one of ResponseTimers
	NoBonusSouls  bool     `json:"no_bonus_souls"`
	NoDeckRatios  bool     `json:"no_deck_ratios"`
}

// ResponseTimers are the host screen's response timer choices; 0 is off.
var ResponseTimers = []int{0, 60, 90, 120, 180}

// maxHostSets caps the remembered card sets.
const maxHostSets = 20

// clean drops values a host screen can't show.
func (h HostSetup) clean() HostSetup {
	if h.Seats < 2 || h.Seats > 4 {
		h.Seats = 0
	}
	if !slices.Contains(ResponseTimers, h.ResponseTimer) {
		h.ResponseTimer = 0
	}
	if h.Sets != nil {
		sets := []string{}
		for _, code := range h.Sets {
			if setCode.MatchString(code) && !slices.Contains(sets, code) && len(sets) < maxHostSets {
				sets = append(sets, code)
			}
		}
		h.Sets = sets
	}
	return h
}

// setCode is a card set code, e.g. "fsp2".
var setCode = regexp.MustCompile(`^[a-z0-9]{1,20}$`)

// Rejoin is where the last game ran and the token that holds its seat.
type Rejoin struct {
	Address string `json:"address"`
	Token   string `json:"token"`
	Game    string `json:"game"`
}

// The choices of the settings that take one of a list; the first is
// the default.
var (
	Speeds = []string{"normal", "fast", "slow", "off"}
	Themes = []string{"dark", "light"}
	Dice   = []string{
		"tear", "blue_fly", "cent", "moms_knife", "the_bible",
		"brimstone", "guppy_paw", "pentagram", "poop", "tinted_rock",
	}
)

// ColorMats are the mats drawn without images; the first is the default
// when there are no mat images.
// Image mats come from the files, so any ValidMat name is kept: a mat
// whose image is missing shows as colors.
var ColorMats = []string{"felt", "basement", "cathedral"}

// matName is a mat file name without its extension.
var matName = regexp.MustCompile(`^[a-z0-9_-]{1,40}$`)

// ValidMat reports whether name can name a mat.
func ValidMat(name string) bool {
	return matName.MatchString(name)
}

// MaxRecent is how many recent servers are kept.
const MaxRecent = 5

// ErrChoice means a setting is not one of its choices.
var ErrChoice = errors.New("unknown setting value")

// ErrNickname means the nickname is empty after trimming, too long, or
// has control characters.
var ErrNickname = fmt.Errorf("nickname must be 1–%d characters", protocol.MaxName)

// choices pairs each one-of-a-list setting with its choices.
func (s *Settings) choices() []struct {
	value   *string
	allowed []string
} {
	return []struct {
		value   *string
		allowed []string
	}{{&s.Animation, Speeds}, {&s.Theme, Themes}, {&s.Dice, Dice}}
}

// Load reads the settings. A missing or broken file gives the defaults.
// Recent is never nil, so the frontend always gets a list; an unknown
// choice reads as its default.
func Load(path string) Settings {
	def := Settings{Recent: []RecentServer{}, Animation: Speeds[0], Theme: Themes[0], Dice: Dice[0]}
	data, err := os.ReadFile(path) //nolint:gosec // the app names the file
	if err != nil {
		return def
	}
	s := def
	if err := json.Unmarshal(data, &s); err != nil {
		return def
	}
	if s.Recent == nil {
		s.Recent = []RecentServer{}
	}
	s.Host = s.Host.clean()
	for _, c := range s.choices() {
		if !slices.Contains(c.allowed, *c.value) {
			*c.value = c.allowed[0]
		}
	}
	if s.Mat != "" && !ValidMat(s.Mat) {
		s.Mat = ""
	}
	return s
}

// Save checks and writes the settings, and returns what was stored.
// The file is replaced in one step, so a crash never leaves half of it.
func Save(path string, s Settings) (Settings, error) {
	if s.Nickname != "" {
		name, ok := protocol.CleanName(s.Nickname)
		if !ok {
			return Settings{}, ErrNickname
		}
		s.Nickname = name
	}
	for _, c := range s.choices() {
		if *c.value == "" {
			*c.value = c.allowed[0]
		}
		if !slices.Contains(c.allowed, *c.value) {
			return Settings{}, fmt.Errorf("%w: %q", ErrChoice, *c.value)
		}
	}
	if s.Mat != "" && !ValidMat(s.Mat) {
		return Settings{}, fmt.Errorf("%w: %q", ErrChoice, s.Mat)
	}
	recent := []RecentServer{}
	for _, r := range s.Recent {
		r.Address = strings.TrimSpace(r.Address)
		// A name that fails the check shows the address instead.
		r.Name, _ = protocol.CleanServerName(r.Name)
		if r.Address != "" && len(recent) < MaxRecent {
			recent = append(recent, r)
		}
	}
	s.Recent = recent
	s.Host = s.Host.clean()
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return Settings{}, fmt.Errorf("settings: %w", err)
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o750); err != nil {
		return Settings{}, fmt.Errorf("settings: %w", err)
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o600); err != nil {
		return Settings{}, fmt.Errorf("settings: %w", err)
	}
	if err := os.Rename(tmp, path); err != nil {
		return Settings{}, errors.Join(fmt.Errorf("settings: %w", err), os.Remove(tmp))
	}
	return s, nil
}
