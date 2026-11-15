package main

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"sync"

	"github.com/FourSoulsRepo/four_souls/internal/assets"
	"github.com/FourSoulsRepo/four_souls/internal/instance"
	"github.com/FourSoulsRepo/four_souls/internal/msgerr"
	"github.com/FourSoulsRepo/four_souls/internal/settings"
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards"
)

// App holds the methods the frontend calls through Wails bindings.
type App struct {
	wails     wailsContext // the app's lifetime; unbound (ADR 030, item 4)
	cardFiles fs.FS        // card data and images
	cardList  cardList
}

// NewApp creates the bindings object.
func NewApp() *App {
	return &App{cardFiles: assets.CardsFS()}
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

// appPath is where the app keeps a file or folder: the user's config
// folder, or next to the app if there is none.
func appPath(name string) string {
	dir, err := os.UserConfigDir()
	if err != nil {
		return name
	}
	return filepath.Join(dir, "FourSouls", name)
}

// wailsContext is the Wails context for runtime calls such as dialogs
// (ADR 030, items 4 and 5). It is a background context before startup
// and ends at shutdown, which stops the work bindings run in it. The
// zero value is ready. It stays apart from App while App is bound, so
// Context does not reach the frontend.
type wailsContext struct {
	mu     sync.Mutex
	ctx    context.Context //nolint:containedctx // Wails hands it over at startup
	cancel context.CancelFunc
}

// startup keeps ctx, the context Wails gives OnStartup.
func (w *wailsContext) startup(ctx context.Context) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.cancel != nil {
		w.cancel()
	}
	w.ctx, w.cancel = context.WithCancel(ctx)
}

// shutdown ends the context.
func (w *wailsContext) shutdown() {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.cancel == nil {
		w.ctx, w.cancel = context.WithCancel(context.Background())
	}
	w.cancel()
}

// Context is the Wails context; a background one before startup.
func (w *wailsContext) Context() context.Context {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.ctx == nil {
		return context.Background()
	}
	return w.ctx
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
		return msgerr.JSON([]engine.Msg{errDeckRatios.Msg(engine.Str("error", err.Error()))})
	}
	if len(short) == 0 {
		return ""
	}
	out := []engine.Msg{engine.RefuseDeckRatios.Msg(engine.N("n", players))}
	for _, sh := range short {
		out = append(out, sh.Msg())
	}
	return msgerr.JSON(out)
}
