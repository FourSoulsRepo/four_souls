package main

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/FourSoulsRepo/four_souls/internal/assets"
	"github.com/FourSoulsRepo/four_souls/internal/instance"
	"github.com/FourSoulsRepo/four_souls/internal/legal"
	"github.com/FourSoulsRepo/four_souls/internal/settings"
	"github.com/FourSoulsRepo/four_souls/internal/version"
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards"
)

// App holds the methods the frontend calls through Wails bindings.
type App struct {
	ctx       context.Context //nolint:containedctx // Wails hands the app context over at startup
	settings  string          // the settings file
	cardFiles fs.FS           // card data and images
	cardList  cardList
}

// NewApp creates the bindings object; settings is its settings file.
func NewApp(settingsFile string) *App {
	return &App{settings: settingsFile, cardFiles: assets.CardsFS()}
}

// slotSettings is the settings file of this copy of the app. A new
// extra copy starts with the first copy's choices, but asks for its own
// nickname and never takes the first copy's seat.
// first is the first copy's settings file.
func slotSettings(slot *instance.Slot, first string) string {
	path := slot.File(filepath.Base(first))
	if path == first {
		return path
	}
	if _, err := os.Stat(path); errors.Is(err, fs.ErrNotExist) {
		s := settings.Load(first)
		s.Nickname, s.Rejoin = "", settings.Rejoin{}
		_, _ = settings.Save(path, s) //nolint:errcheck // a missing file gives the defaults
	}
	return path
}

// startup keeps the app context for runtime calls such as dialogs.
func (a *App) startup(ctx context.Context) { a.ctx = ctx }

// Versions returns the app and rules engine versions for the main menu.
func (a *App) Versions() version.Info {
	return version.Get()
}

// Notice returns the fan-game notice for the splash screen.
func (a *App) Notice() []legal.Line {
	return legal.Notice
}

// Settings returns the saved settings, or the defaults.
func (a *App) Settings() settings.Settings {
	return settings.Load(a.settings)
}

// SaveSettings stores the settings and returns what was stored.
func (a *App) SaveSettings(s settings.Settings) (settings.Settings, error) {
	saved, err := settings.Save(a.settings, s)
	return saved, playerError(err)
}

// CardSets returns the codes of the card sets this app has; the client
// sends them when it connects (CD-03). The challenge set is not among
// them: it is never a playable set (ADR 028).
func (a *App) CardSets() []string {
	var codes []string
	for _, s := range cards.Sets() {
		codes = append(codes, s.Code)
	}
	return codes
}

// ChallengeInfo is a challenge the host screen offers (ADR 028,
// Decision 10.7): its ID and difficulties, easiest first. The client
// names them; a server that lacks one refuses it at create.
type ChallengeInfo struct {
	ID           string   `json:"id"`
	Difficulties []string `json:"difficulties"`
}

// Challenges lists the challenges this app knows, like CardSets lists
// its sets.
func (a *App) Challenges() []ChallengeInfo {
	var out []ChallengeInfo
	for _, c := range cards.Challenges() {
		info := ChallengeInfo{ID: c.ID}
		for _, l := range c.Levels {
			info.Difficulties = append(info.Difficulties, l.Difficulty)
		}
		out = append(out, info)
	}
	return out
}

// DeckRatios says why a game with these sets and players can't use the
// deck ratios, or "" when it can (GS-07). The host screen shows it as
// the reason the box is off. The reason is a JSON list of messages
// (ADR 021): the head, {n} players, then one line per short group.
func (a *App) DeckRatios(sets []string, players int) string {
	s := engine.Setup{Players: players}
	for _, code := range sets {
		if set, ok := cards.Find(code); ok {
			s.Sets = append(s.Sets, set)
		}
	}
	short, err := engine.RatioShortfall(s)
	if err != nil {
		return msgJSON([]engine.Msg{errDeckRatios.Msg(engine.Str("error", err.Error()))})
	}
	if len(short) == 0 {
		return ""
	}
	out := []engine.Msg{engine.RefuseDeckRatios.Msg(engine.N("n", players))}
	for _, sh := range short {
		out = append(out, sh.Msg())
	}
	return msgJSON(out)
}
