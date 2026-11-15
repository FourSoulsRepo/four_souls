package main

import (
	"context"
	"errors"
	"io/fs"
	"os"
	"path/filepath"
	"sync"

	"github.com/FourSoulsRepo/four_souls/internal/instance"
	"github.com/FourSoulsRepo/four_souls/internal/settings"
)

// App is the app's lifetime; it is not bound (ADR 030, item 4). It
// keeps the Wails context for runtime calls such as dialogs. Context is
// a background context before startup and ends at shutdown, which stops
// the work bindings run in it (ADR 030, item 5). The zero value is ready.
type App struct {
	mu     sync.Mutex
	ctx    context.Context //nolint:containedctx // Wails hands it over at startup
	cancel context.CancelFunc
}

// NewApp creates the app.
func NewApp() *App {
	return &App{}
}

// startup keeps ctx, the context Wails gives OnStartup.
func (a *App) startup(ctx context.Context) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.cancel != nil {
		a.cancel()
	}
	a.ctx, a.cancel = context.WithCancel(ctx)
}

// shutdown ends the context.
func (a *App) shutdown() {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.cancel == nil {
		a.ctx, a.cancel = context.WithCancel(context.Background())
	}
	a.cancel()
}

// Context is the Wails context; a background one before startup.
func (a *App) Context() context.Context {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.ctx == nil {
		return context.Background()
	}
	return a.ctx
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
