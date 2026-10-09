package main

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/FourSoulsRepo/four_souls/internal/server"
)

// Host runs the game server inside the app, so a player can host a
// game for friends (N-01). The host joins it as a normal client.
type Host struct {
	mu  sync.Mutex
	srv *server.Server
	cfg server.Config
}

// HostInfo tells the host how friends reach the game.
type HostInfo struct {
	Running   bool     `json:"running"`
	Port      int      `json:"port"`
	Addresses []string `json:"addresses"` // give one of these to friends
}

// ErrHosting means a game is hosted already.
var ErrHosting = errors.New("a game is hosted already")

// Start runs the server on the default port; the host then creates a
// game in the lobby like any player.
func (h *Host) Start() (HostInfo, error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.srv != nil {
		return h.info(), ErrHosting
	}
	cfg := server.DefaultConfig()
	srv, err := server.Start(context.Background(), cfg)
	if err != nil {
		return HostInfo{}, err //nolint:wrapcheck // shown to the player as is
	}
	h.srv, h.cfg = srv, cfg
	return h.info(), nil
}

// Stop ends the hosted game.
func (h *Host) Stop() error { return h.stop(context.Background()) }

// stop ends the hosted game within ctx, at most 10 seconds.
func (h *Host) stop(ctx context.Context) error {
	h.mu.Lock()
	defer h.mu.Unlock()
	if h.srv == nil {
		return nil
	}
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	err := h.srv.Shutdown(ctx)
	h.srv = nil
	return err //nolint:wrapcheck // shown to the player as is
}

// Info reports whether a game is hosted and where.
func (h *Host) Info() HostInfo {
	h.mu.Lock()
	defer h.mu.Unlock()
	return h.info()
}

func (h *Host) info() HostInfo {
	if h.srv == nil {
		return HostInfo{}
	}
	return HostInfo{Running: true, Port: h.cfg.Port, Addresses: server.LocalAddresses()}
}
