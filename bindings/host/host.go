// Package host is the binding that hosts a game inside the app (N-01,
// ADR 030).
package host

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/FourSoulsRepo/four_souls/internal/msgerr"
	"github.com/FourSoulsRepo/four_souls/internal/server"
	engine "github.com/FourSoulsRepo/rules_engine"
)

// errHosting is ErrHosting in the player's words (ADR 021); its text
// is in frontend/src/locales/en.json.
var errHosting = engine.NewKey("error.hosting")

// Host runs the game server inside the app, so a player can host a
// game for friends (N-01). The host joins it as a normal client.
type Host struct {
	records, saves string // the server's folders

	mu  sync.Mutex
	srv *server.Server
	cfg server.Config
}

// New makes the binding; the hosted server keeps its match records
// and saved games in these folders.
func New(records, saves string) *Host {
	return &Host{records: records, saves: saves}
}

// Info tells the host how friends reach the game.
type Info struct {
	Running   bool     `json:"running"`
	Port      int      `json:"port"`
	Addresses []string `json:"addresses"` // give one of these to friends
}

// ErrHosting means a game is hosted already.
var ErrHosting = errors.New("a game is hosted already")

// Start runs the server on the default port; the host then creates a
// game in the lobby like any player.
func (h *Host) Start() (Info, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.srv != nil {
		return h.info(), msgerr.New(errHosting.Msg(), ErrHosting)
	}
	cfg := server.DefaultConfig()
	cfg.Records, cfg.Retention = h.records, 0 // kept until the user deletes them (RP-07)
	cfg.Saves = h.saves
	cfg.Discovery = true // friends on the LAN find the game (N-13)
	cfg.HostOnly = true  // friends join; only this app opens games (N-17)
	// TestSeeds stays off: no player picks the deal (ADR 032).
	srv, err := server.Start(context.Background(), cfg)
	if err != nil {
		return Info{}, err //nolint:wrapcheck // shown to the player as is
	}
	h.srv, h.cfg = srv, cfg
	return h.info(), nil
}

// Stop ends the hosted game within 10 seconds. The app calls it when
// it closes too.
func (h *Host) Stop() error {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.srv == nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err := h.srv.Shutdown(ctx)
	h.srv = nil
	return err //nolint:wrapcheck // shown to the player as is
}

// Info reports whether a game is hosted and where.
func (h *Host) Info() Info {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.info()
}

func (h *Host) info() Info {
	if h.srv == nil {
		return Info{}
	}
	return Info{Running: true, Port: h.cfg.Port, Addresses: server.LocalAddresses()}
}
