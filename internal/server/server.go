package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	carddb "github.com/FourSoulsRepo/card_db"
	"github.com/FourSoulsRepo/four_souls/internal/assets"
	"github.com/FourSoulsRepo/four_souls/internal/lan"
	"github.com/FourSoulsRepo/four_souls/internal/protocol"
	"github.com/FourSoulsRepo/four_souls/internal/version"
	"github.com/FourSoulsRepo/record"
)

// Config is the server's settings, from flags or a JSON file (6.3).
type Config struct {
	Addr      string `json:"addr"`      // listen address; empty: all
	Port      int    `json:"port"`      // TCP port
	Records   string `json:"records"`   // folder for match records (6.9)
	Retention int    `json:"retention"` // days to keep records; 0 keeps them
	Saves     string `json:"saves"`     // folder for saved games (N-11)
	// Seed makes game seeds reproducible: game n gets Seed+n. 0 seeds
	// from the clock. For tests and bug reports.
	Seed uint64 `json:"seed,omitempty"`
	// Discovery answers LAN probes while a game waits (N-13, ADR 009).
	// Only the app's host sets it; a dedicated server never announces.
	Discovery bool `json:"-"`
	// HostOnly lets only a client on this computer create or load a
	// game: the app's host. Friends join (N-17). A dedicated server
	// lets everyone create.
	HostOnly bool `json:"-"`
	// Presence relays live aiming and hover (15.35, ADR 016). An
	// operator short of CPU or bandwidth may turn it off.
	Presence bool `json:"presence"`
	// Name is shown in the app instead of the address (N-18); empty
	// shows the address. Up to protocol.MaxServerName characters.
	Name string `json:"name"`
}

// DefaultConfig is used for settings the flags and file leave out.
func DefaultConfig() Config {
	return Config{Port: protocol.DefaultPort, Records: "records", Retention: 30, Saves: "saves", Presence: true}
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
	case c.Name != "" && !validName(c.Name):
		return fmt.Errorf("server: name must be 1–%d characters, without control characters", protocol.MaxServerName)
	}
	return nil
}

// validName reports whether a server name passes CleanServerName.
func validName(name string) bool {
	_, ok := protocol.CleanServerName(name)
	return ok
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
	hub.records, hub.cards, hub.seed, hub.saves = cfg.Records, cardTexts(), cfg.Seed, cfg.Saves
	hub.hostOnly, hub.presence = cfg.HostOnly, cfg.Presence
	hub.name, _ = protocol.CleanServerName(cfg.Name)
	s.hub.Add(2)
	go func() {
		defer s.hub.Done()
		pruneLoop(ctx.Done(), cfg.Records, cfg.Retention)
	}()
	go func() {
		defer s.hub.Done()
		hub.Run(ctx)
	}()
	if cfg.Discovery {
		port := cfg.Port
		if a, ok := ln.Addr().(*net.TCPAddr); ok {
			port = a.Port // the real one when cfg.Port is 0
		}
		if err := lan.Listen(ctx, port, waitingGames(hub, port)); err != nil {
			log.Printf("server: LAN discovery is off: %v", err) // joining by address still works
		}
	}
	mux := http.NewServeMux()
	mux.Handle("/ws", Handler(hub))
	mux.HandleFunc("/record", hub.downloadRecord)
	mux.HandleFunc("/healthz", hub.healthz)
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

// cardTexts reads the card data for records; without it records carry
// no card text (RP-10).
func cardTexts() []record.Card {
	db, err := carddb.Load(assets.CardsFS())
	if err != nil {
		return nil
	}
	var out []record.Card
	for _, ref := range db.Refs() {
		c, _ := db.Get(ref)
		en := c.Text[carddb.English]
		out = append(out, record.Card{Ref: ref.String(), Name: en.Name, Type: c.Type, Effects: en.Effects})
	}
	return out
}

// waitingGames answers a LAN probe with the games that wait for
// players; none waiting means no answer (ADR 009).
func waitingGames(hub *Hub, port int) func() (lan.Answer, bool) {
	id := newToken()
	return func() (lan.Answer, bool) {
		reply := make(chan []lan.Game, 1)
		hub.send(hubRequest{run: func(h *Hub) {
			var out []lan.Game
			for _, g := range h.list().Games {
				if !g.Started && g.Taken < g.Seats {
					out = append(out, lan.Game{Host: g.Host, Seats: g.Seats, Taken: g.Taken, Locked: g.Locked})
				}
			}
			reply <- out
		}})
		var games []lan.Game
		select {
		case games = <-reply:
		case <-hub.done:
		}
		v := version.Get()
		return lan.Answer{ID: id, Protocol: protocol.Version, App: v.App, Port: port, Games: games}, len(games) > 0
	}
}
