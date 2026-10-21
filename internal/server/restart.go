package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"time"

	"github.com/FourSoulsRepo/four_souls/internal/protocol"
	"github.com/FourSoulsRepo/four_souls/internal/version"
	"github.com/FourSoulsRepo/record"
	engine "github.com/FourSoulsRepo/rules_engine"
)

// A server that stops keeps its tables and games in restart files; the
// next start brings them back and their players reconnect by token
// (N-20, ADR 019).

// restartFolder is the folder in the saves folder for restart files.
const restartFolder = "restart"

// restartExt is the extension of restart files.
const restartExt = ".fsrestart"

// lostExt is added to the restart file of a game that could not go on.
const lostExt = ".lost"

// restartWarning is how long the clients get to read the restarting
// message before their connections close.
const restartWarning = 300 * time.Millisecond

// What a restart file keeps.
const (
	keptTable = "table" // a lobby table, before its game started
	keptSetup = "setup" // bans and picks; they start over
	keptGame  = "game"  // a running game
)

// restartFile is one table or game kept over a restart. It holds every
// hand, so it stays on the server (RP-05).
type restartFile struct {
	saveFile

	Kind    string    `json:"kind"`
	N       int       `json:"n"` // the game's number on the server
	Host    string    `json:"host"`
	Started time.Time `json:"started"` // when the game left the lobby
	// Seed is the setup's: bans and picks deal the same offers again.
	Seed uint64 `json:"seed,omitempty"`
	// Password of a lobby table (N-14); the file is deleted on load.
	Password string `json:"password,omitempty"`
	// Meaning is the saving engine's StateMeaning (engine.Resume).
	Meaning int `json:"meaning,omitempty"`
}

// voteGrace is how many seconds a game kept over a restart takes no kick
// vote: its players are still reconnecting.
const voteGrace = 60

// lostFor is how long the players of a game that could not go on after
// a restart are told so when they reconnect.
const lostFor = 30 * time.Minute

// lostGame is a game the new server could not go on with: its players
// get an apology when they come back (N-20).
type lostGame struct {
	game   string
	tokens []string // the players not told yet
	until  time.Time
}

// restartPath names the restart file of a game.
func restartPath(dir, game string) string { return filepath.Join(dir, game+restartExt) }

// writeRestart writes f to dir.
func writeRestart(dir string, f *restartFile) error {
	v := version.Get()
	f.Format, f.App, f.Engine, f.Protocol, f.Saved = saveFormat, v.App, v.Engine, protocol.Version, time.Now().UTC()
	data, err := json.Marshal(f) //nolint:gosec // G117: the table's password goes on with it; the file is 0600 on the server
	if err != nil {
		return fmt.Errorf("restart file: %w", err)
	}
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return fmt.Errorf("restart file: %w", err)
	}
	if err := os.WriteFile(restartPath(dir, f.Game), data, 0o600); err != nil {
		return fmt.Errorf("restart file: %w", err)
	}
	return nil
}

// readRestarts reads every restart file of dir, oldest game first.
func readRestarts(dir string) []*restartFile {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	var out []*restartFile
	for _, e := range entries {
		if !strings.HasSuffix(e.Name(), restartExt) {
			continue
		}
		path := filepath.Join(dir, e.Name())
		data, err := os.ReadFile(path) //nolint:gosec // a file of the server's own folder
		var f restartFile
		if err == nil {
			err = json.Unmarshal(data, &f)
		}
		if err == nil && (f.Format != saveFormat || f.Game == "" || len(f.Seats) == 0) {
			err = errors.New("not a restart file this server can read")
		}
		if err != nil {
			log.Printf("server: restart file %s: %v", e.Name(), err)
			continue
		}
		out = append(out, &f)
	}
	slices.SortFunc(out, func(a, b *restartFile) int { return a.N - b.N })
	return out
}

// warnRestart tells every client the server restarts; the app then
// reconnects by itself (N-20).
func (h *Hub) warnRestart() {
	for _, c := range h.clients {
		if c.hello {
			c.send(protocol.TypeRestarting, 0, struct{}{})
		}
	}
}

// keepTables writes the lobby tables to restart files. A saved game's
// table is not kept: its save is still there to load (N-11).
func (h *Hub) keepTables() {
	for _, t := range h.tables {
		if t.loaded != nil {
			continue
		}
		f := &restartFile{
			saveFile: saveFile{Game: t.id, Sets: t.sets, Options: t.opts},
			Kind:     keptTable, N: t.n, Host: t.host, Password: t.password,
		}
		for _, s := range t.seats {
			f.Seats = append(f.Seats, savedSeat{Name: s.name, Token: s.token})
		}
		if err := writeRestart(h.restart, f); err != nil {
			log.Printf("server: keep table %s: %v", t.id, err)
		}
	}
}

// restore brings back the tables and games of the restart files. Their
// seats wait for their players' tokens.
func (h *Hub) restore() {
	for _, f := range readRestarts(h.restart) {
		h.seq = max(h.seq, f.N)
		switch f.Kind {
		case keptTable:
			h.restoreTable(f)
		case keptSetup, keptGame:
			h.restoreRoom(f)
		default:
			log.Printf("server: restart file %s: kind %q", f.Game, f.Kind)
		}
	}
}

func (h *Hub) restoreTable(f *restartFile) {
	t := &table{n: f.N, id: f.Game, host: f.Host, sets: f.Sets, opts: f.Options, password: f.Password, seats: make([]tableSeat, len(f.Seats))}
	for i, s := range f.Seats {
		t.seats[i] = tableSeat{name: s.Name, token: s.Token}
	}
	h.tables = append(h.tables, t)
	if err := os.Remove(restartPath(h.restart, f.Game)); err != nil {
		log.Printf("server: restart file %s: %v", f.Game, err)
	}
}

func (h *Hub) restoreRoom(f *restartFile) {
	var g *engine.Game
	if f.Kind == keptGame {
		var err error
		if g, err = engine.Resume(f.State, f.Meaning, setsOf(f.Sets)...); err != nil {
			h.lose(f, err)
			return
		}
	}
	seats := make([]roomSeat, len(f.Seats))
	for i, s := range f.Seats {
		seats[i] = roomSeat{name: s.Name, token: s.Token, kicked: s.Kicked, dice: s.Dice}
	}
	setup := engine.Setup{
		Seed: f.Seed, Players: len(seats), Sets: setsOf(f.Sets), BonusSouls: !f.Options.NoBonusSouls,
		Banned: f.Options.HostBans, DeckRatios: !f.Options.NoDeckRatios,
	}
	room, err := NewRoom(setup, f.Options, seats)
	if err != nil {
		log.Printf("server: restart file %s: %v", f.Game, err)
		return
	}
	h.setUpRoom(room, f.N, f.Game, f.Host, f.Sets)
	room.kept, room.path, room.born, room.game = f, f.Record, f.Started, g
	if g != nil && f.Engine != engine.Version {
		room.notice = protocol.NoticeEngine // the rules may have changed
	}
	h.run(room, protocol.GameInfo{ID: f.Game, Host: f.Host, Seats: len(seats), Taken: len(seats), Sets: f.Sets, Started: true})
}

// keep writes the room to a restart file instead of ending its record;
// a finished game is not kept.
func (r *Room) keep() {
	if r.hub == nil || r.hub.restart == "" || r.over.Load() {
		return
	}
	f := &restartFile{
		saveFile: saveFile{Game: r.recInfo.game, Sets: r.recInfo.sets, Options: r.opts, Record: r.path},
		N:        r.n, Host: r.host, Started: r.born, Seed: r.setup.Seed,
	}
	for _, s := range r.seats {
		f.Seats = append(f.Seats, savedSeat{Name: s.name, Token: s.token, Kicked: s.kicked, Dice: s.dice})
	}
	if r.game == nil {
		f.Kind = keptSetup
	} else {
		state, err := r.game.Save()
		if err != nil {
			log.Printf("server: keep %s: %v", f.Game, err)
			return
		}
		f.Kind, f.State, f.Meaning, f.Step, f.Steps = keptGame, state, engine.StateMeaning, r.step, r.steps
	}
	if err := writeRestart(r.hub.restart, f); err != nil {
		log.Printf("server: keep %s: %v", f.Game, err)
		return
	}
	if r.rec != nil {
		if err := r.rec.Stop(); err != nil {
			log.Printf("server: record %s: %v", r.path, err)
		}
		r.rec = nil // no end line: the game goes on
	}
}

// resumeKept starts a room from its restart file; the hub has loaded its
// game. Bans and picks start over once every player is back; a game
// waits paused for them.
func (r *Room) resumeKept() {
	f := r.kept
	r.noticed = make([]bool, len(r.seats))
	if f.Kind == keptSetup {
		r.setupWait = true
		r.forgetKept()
		return
	}
	r.step, r.steps = f.Step, f.Steps
	r.grace = time.Now().Add(voteGrace * r.second)
	if r.recDir != "" && f.Record != "" {
		if w, err := record.Append(f.Record, f.Steps, time.Now().UTC()); err == nil {
			r.rec = w
		} else {
			log.Printf("server: record %s: %v", f.Record, err)
		}
	}
	r.forgetKept()
	r.after(nil, false)
}

// lose ends a kept game the new server cannot go on with. Its record
// gets its end; the restart file stays for the operator as .lost. Its
// players get an apology when they come back.
func (h *Hub) lose(f *restartFile, why error) {
	log.Printf("server: game %s cannot go on after the restart: %v", f.Game, why)
	if h.records != "" && f.Record != "" {
		now := time.Now().UTC()
		w, err := record.Append(f.Record, f.Steps, now)
		if err == nil {
			err = w.Close(record.End{Ended: now, Finished: false})
		}
		if err != nil {
			log.Printf("server: record %s: %v", f.Record, err)
		}
	}
	path := restartPath(h.restart, f.Game)
	if err := os.Rename(path, path+lostExt); err != nil {
		log.Printf("server: restart file %s: %v", f.Game, err)
	}
	lost := lostGame{game: f.Game, until: time.Now().Add(lostFor)}
	for _, s := range f.Seats {
		lost.tokens = append(lost.tokens, s.Token)
	}
	h.lost = append(h.lost, lost)
}

// tellLost sends the apology to a player of a lost game, once.
func (h *Hub) tellLost(c *Client) {
	now := time.Now()
	h.lost = slices.DeleteFunc(h.lost, func(l lostGame) bool { return now.After(l.until) || len(l.tokens) == 0 })
	for i := range h.lost {
		l := &h.lost[i]
		if j := slices.Index(l.tokens, c.token); j >= 0 {
			l.tokens = slices.Delete(l.tokens, j, j+1)
			c.send(protocol.TypeNotice, 0, protocol.Notice{Kind: protocol.NoticeLost, Game: l.game})
		}
	}
}

// tellEngine tells a returning player, once, that the game went on
// with a newer engine.
func (r *Room) tellEngine(c *Client) {
	if r.notice == "" || c.seat < 0 || r.noticed[c.seat] {
		return
	}
	r.noticed[c.seat] = true
	c.send(protocol.TypeNotice, 0, protocol.Notice{Kind: r.notice, Game: r.recInfo.game})
}

// forgetKept deletes the room's restart file once it runs again.
func (r *Room) forgetKept() {
	if err := os.Remove(restartPath(r.hub.restart, r.kept.Game)); err != nil {
		log.Printf("server: restart file %s: %v", r.kept.Game, err)
	}
}

// everyoneBack: every seat has its player again.
func (r *Room) everyoneBack() bool {
	return !slices.ContainsFunc(r.seats, func(s roomSeat) bool { return s.client == nil })
}
