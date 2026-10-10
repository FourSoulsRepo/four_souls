// Package settings keeps the player's app settings in a JSON file in
// the user's config folder, so they survive a restart.
package settings

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/FourSoulsRepo/four_souls/internal/protocol"
)

// Settings are the player's app settings.
type Settings struct {
	// Nickname is sent when joining a game (ST-06); "" until set.
	Nickname string `json:"nickname"`
	// Recent are the last addresses joined, newest first.
	Recent []string `json:"recent"`
	// Animation is the animation speed (ST-05): one of Speeds.
	Animation string `json:"animation"`
}

// Speeds are the animation speeds; the first is the default.
var Speeds = []string{"normal", "fast", "slow", "off"}

// MaxRecent is how many recent addresses are kept.
const MaxRecent = 5

// ErrAnimation means the animation speed is not one of Speeds.
var ErrAnimation = errors.New("unknown animation speed")

// ErrNickname means the nickname is empty after trimming, too long, or
// has control characters.
var ErrNickname = fmt.Errorf("nickname must be 1–%d characters", protocol.MaxName)

// Load reads the settings. A missing or broken file gives the defaults.
// Recent is never nil, so the frontend always gets a list; an unknown
// animation speed reads as the default.
func Load(path string) Settings {
	def := Settings{Recent: []string{}, Animation: Speeds[0]}
	data, err := os.ReadFile(path) //nolint:gosec // the app names the file
	if err != nil {
		return def
	}
	s := def
	if err := json.Unmarshal(data, &s); err != nil {
		return def
	}
	if s.Recent == nil {
		s.Recent = []string{}
	}
	if !slices.Contains(Speeds, s.Animation) {
		s.Animation = Speeds[0]
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
	if s.Animation == "" {
		s.Animation = Speeds[0]
	}
	if !slices.Contains(Speeds, s.Animation) {
		return Settings{}, ErrAnimation
	}
	recent := []string{}
	for _, a := range s.Recent {
		if a = strings.TrimSpace(a); a != "" && len(recent) < MaxRecent {
			recent = append(recent, a)
		}
	}
	s.Recent = recent
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
