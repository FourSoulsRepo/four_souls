package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/FourSoulsRepo/four_souls/internal/protocol"
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards"
)

// Config is the server's settings, from flags or a JSON file (6.3).
type Config struct {
	Addr      string `json:"addr"`      // listen address; empty: all
	Port      int    `json:"port"`      // TCP port
	Records   string `json:"records"`   // folder for match records (6.9)
	Retention int    `json:"retention"` // days to keep records; 0 keeps them
	// Players is the size of the one game until the lobby exists (6.5).
	Players int `json:"players"`
}

// DefaultConfig is used for settings the flags and file leave out.
func DefaultConfig() Config {
	return Config{Port: protocol.DefaultPort, Records: "records", Retention: 30, Players: 2}
}

// LoadConfig reads a JSON config file over the defaults.
func LoadConfig(path string) (Config, error) {
	c := DefaultConfig()
	data, err := os.ReadFile(path) //nolint:gosec // the operator names the file
	if err != nil {
		return c, fmt.Errorf("server: config: %w", err)
	}
	if err := json.Unmarshal(data, &c); err != nil {
		return c, fmt.Errorf("server: config %s: %w", path, err)
	}
	return c, c.Check()
}

// Check reports settings that cannot work.
func (c Config) Check() error {
	switch {
	case c.Port < 1 || c.Port > 65535:
		return fmt.Errorf("server: port %d is not 1 to 65535", c.Port)
	case c.Retention < 0:
		return errors.New("server: retention is days, 0 or more")
	case c.Players < 2 || c.Players > 4:
		return fmt.Errorf("server: %d players; a game is for 2 to 4", c.Players)
	}
	return nil
}

// Server listens for clients and runs their rooms.
type Server struct {
	cfg    Config
	http   *http.Server
	ln     net.Listener
	cancel context.CancelFunc
	rooms  sync.WaitGroup
}

// Start listens and runs the server until Shutdown or until ctx ends.
// Until the lobby (6.5) it runs one game of cfg.Players seats.
func Start(ctx context.Context, cfg Config) (*Server, error) {
	if err := cfg.Check(); err != nil {
		return nil, err
	}
	room, err := NewRoom(engine.Setup{Seed: seed(), Players: cfg.Players, Sets: cards.Sets(), BonusSouls: true})
	if err != nil {
		return nil, err
	}
	var lc net.ListenConfig
	ln, err := lc.Listen(ctx, "tcp", net.JoinHostPort(cfg.Addr, strconv.Itoa(cfg.Port)))
	if err != nil {
		return nil, fmt.Errorf("server: %w", err)
	}
	ctx, cancel := context.WithCancel(ctx)
	s := &Server{cfg: cfg, ln: ln, cancel: cancel}
	s.rooms.Add(1)
	go func() {
		defer s.rooms.Done()
		room.Run(ctx)
	}()
	mux := http.NewServeMux()
	mux.Handle("/ws", Handler(room))
	s.http = &http.Server{Handler: mux, ReadHeaderTimeout: 10 * time.Second}
	go func() { _ = s.http.Serve(ln) }() //nolint:errcheck // ends with ErrServerClosed on Shutdown
	return s, nil
}

// Addr is the address the server listens on.
func (s *Server) Addr() string { return s.ln.Addr().String() }

// Shutdown stops accepting clients and stops every room. Saving running
// records joins here with 6.9.
func (s *Server) Shutdown(ctx context.Context) error {
	err := s.http.Shutdown(ctx)
	s.cancel()
	s.rooms.Wait()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("server: shutdown: %w", err)
	}
	return nil
}

// seed is a new game's seed; the engine itself never reads the clock.
func seed() uint64 { return uint64(time.Now().UnixNano()) } //nolint:gosec // a game seed, not a secret
