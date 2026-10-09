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
)

// Config is the server's settings, from flags or a JSON file (6.3).
type Config struct {
	Addr      string `json:"addr"`      // listen address; empty: all
	Port      int    `json:"port"`      // TCP port
	Records   string `json:"records"`   // folder for match records (6.9)
	Retention int    `json:"retention"` // days to keep records; 0 keeps them
}

// DefaultConfig is used for settings the flags and file leave out.
func DefaultConfig() Config {
	return Config{Port: protocol.DefaultPort, Records: "records", Retention: 30}
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
	}
	return nil
}

// Server listens for clients and runs their rooms.
type Server struct {
	cfg    Config
	http   *http.Server
	ln     net.Listener
	cancel context.CancelFunc
	hub    sync.WaitGroup
}

// Start listens and runs the lobby and its games until Shutdown or
// until ctx ends.
func Start(ctx context.Context, cfg Config) (*Server, error) {
	if err := cfg.Check(); err != nil {
		return nil, err
	}
	var lc net.ListenConfig
	ln, err := lc.Listen(ctx, "tcp", net.JoinHostPort(cfg.Addr, strconv.Itoa(cfg.Port)))
	if err != nil {
		return nil, fmt.Errorf("server: %w", err)
	}
	ctx, cancel := context.WithCancel(ctx)
	s := &Server{cfg: cfg, ln: ln, cancel: cancel}
	hub := NewHub()
	s.hub.Add(1)
	go func() {
		defer s.hub.Done()
		hub.Run(ctx)
	}()
	mux := http.NewServeMux()
	mux.Handle("/ws", Handler(hub))
	s.http = &http.Server{Handler: mux, ReadHeaderTimeout: 10 * time.Second}
	go func() { _ = s.http.Serve(ln) }() //nolint:errcheck // ends with ErrServerClosed on Shutdown
	return s, nil
}

// Addr is the address the server listens on.
func (s *Server) Addr() string { return s.ln.Addr().String() }

// Shutdown stops accepting clients and stops the lobby and every game.
// Saving running records joins here with 6.9.
func (s *Server) Shutdown(ctx context.Context) error {
	err := s.http.Shutdown(ctx)
	s.cancel()
	s.hub.Wait()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("server: shutdown: %w", err)
	}
	return nil
}

// LocalAddresses lists this computer's IPv4 addresses that friends may
// reach: LAN and virtual LAN (Tailscale, ZeroTier, …), not loopback.
func LocalAddresses() []string {
	var out []string
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}
		for _, a := range addrs {
			if ipn, ok := a.(*net.IPNet); ok {
				if ip4 := ipn.IP.To4(); ip4 != nil && !ip4.IsLinkLocalUnicast() {
					out = append(out, ip4.String())
				}
			}
		}
	}
	return out
}
