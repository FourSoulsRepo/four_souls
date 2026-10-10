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

// saveExt is the extension of saved games (N-11).
const saveExt = ".fssave"

// saveFormat is the saved game format version.
const saveFormat = 1

// saveFile is an unfinished game kept to continue later (N-11). It
// holds every hand, so it stays on the server (RP-05).
type saveFile struct {
	Format   int              `json:"format"`
	App      string           `json:"app"`
	Engine   string           `json:"engine"`
	Protocol int              `json:"protocol"`
	Game     string           `json:"game"`
	Saved    time.Time        `json:"saved"`
	Sets     []string         `json:"sets"`
	Options  protocol.Options `json:"options"`
	Seats    []savedSeat      `json:"seats"`
	Step     int              `json:"step"`   // the step number to go on from
	Steps    int              `json:"steps"`  // steps in the record so far
	Record   string           `json:"record"` // the record file to append to
	State    json.RawMessage  `json:"state"`  // the engine's Save
}

type savedSeat struct {
	Name  string `json:"name"`
	Token string `json:"token"`
}

// save writes the game to the saves folder and stops the room: only
// the host (seat 1) saves, at any moment of a running game (N-11).
func (r *Room) save(c *Client, env protocol.Envelope) {
	switch {
	case c.seat != 0:
		c.fail(env.ID, protocol.ErrNotHost, "only the host saves the game")
		return
	case r.game == nil || r.game.Over:
		c.fail(env.ID, protocol.ErrBadMessage, "no running game to save")
		return
	case r.saves == "":
		c.fail(env.ID, protocol.ErrBadMessage, "this server does not keep saves")
		return
	}
	state, err := r.game.Save()
	if err != nil {
		c.fail(env.ID, protocol.ErrBadMessage, err.Error())
		return
	}
	v := version.Get()
	f := saveFile{
		Format: saveFormat, App: v.App, Engine: v.Engine, Protocol: protocol.Version,
		Game: r.recInfo.game, Saved: time.Now().UTC(), Sets: r.recInfo.sets, Options: r.opts,
		Step: r.step, Steps: r.steps, Record: r.path, State: state,
	}
	for _, s := range r.seats {
		f.Seats = append(f.Seats, savedSeat{Name: s.name, Token: s.token})
	}
	data, err := json.Marshal(f)
	if err != nil {
		c.fail(env.ID, protocol.ErrBadMessage, err.Error())
		return
	}
	if err := os.MkdirAll(r.saves, 0o750); err != nil {
		c.fail(env.ID, protocol.ErrBadMessage, err.Error())
		return
	}
	if err := os.WriteFile(savePath(r.saves, f.Game), data, 0o600); err != nil {
		c.fail(env.ID, protocol.ErrBadMessage, err.Error())
		return
	}
	r.closeRecord(false)
	for _, x := range r.clients {
		x.send(protocol.TypeSaved, 0, protocol.Saved{Save: f.Game})
	}
	r.stopped = true // Run ends; everyone goes back to the lobby
}

func savePath(dir, game string) string { return filepath.Join(dir, game+saveExt) }

// readSave reads a saved game by ID.
func readSave(dir, id string) (*saveFile, error) {
	if id == "" || strings.ContainsAny(id, `/\.`) {
		return nil, errors.New("no such save")
	}
	data, err := os.ReadFile(savePath(dir, id)) //nolint:gosec // the ID has no path separators
	if err != nil {
		return nil, fmt.Errorf("no such save: %w", err)
	}
	var f saveFile
	if err := json.Unmarshal(data, &f); err != nil || f.Format != saveFormat {
		return nil, errors.New("not a save this server can read")
	}
	return &f, nil
}

// listSaves describes the saved games for the lobby.
func listSaves(dir string) []protocol.SaveInfo {
	if dir == "" {
		return nil
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil
	}
	var out []protocol.SaveInfo
	for _, e := range entries {
		id, ok := strings.CutSuffix(e.Name(), saveExt)
		if !ok {
			continue
		}
		f, err := readSave(dir, id)
		if err != nil {
			continue
		}
		info := protocol.SaveInfo{ID: id, Saved: f.Saved.UnixMilli(), Engine: f.Engine}
		for _, s := range f.Seats {
			info.Seats = append(info.Seats, s.Name)
		}
		if len(info.Seats) > 0 {
			info.Host = info.Seats[0]
		}
		out = append(out, info)
	}
	return out
}

// load opens a saved game's table: its seats wait for their players,
// who come back by token, or by nickname (N-11).
func (h *Hub) load(c *Client, env protocol.Envelope) {
	var m protocol.Load
	if err := env.Unpack(&m); err != nil {
		c.fail(env.ID, protocol.ErrBadMessage, err.Error())
		return
	}
	if c.table != nil {
		c.fail(env.ID, protocol.ErrAtTable, "leave your table first")
		return
	}
	if slices.ContainsFunc(h.tables, func(t *table) bool { return t.loaded != nil && t.loaded.Game == m.Save }) {
		h.joinLoaded(c, env, m.Save)
		return
	}
	f, err := readSave(h.saves, m.Save)
	if err != nil {
		c.fail(env.ID, protocol.ErrNoSave, err.Error())
		return
	}
	if f.Engine != engine.Version {
		c.fail(env.ID, protocol.ErrSaveVersion, fmt.Sprintf("saved with rules engine %s; this server runs %s (N-11)", f.Engine, engine.Version))
		return
	}
	if missing := lacks(c.sets, f.Sets); missing != "" {
		c.fail(env.ID, protocol.ErrMissingSets, "you do not have the card set "+missing)
		return
	}
	t := &table{id: f.Game, host: f.Seats[0].Name, sets: f.Sets, opts: f.Options, loaded: f, seats: make([]tableSeat, len(f.Seats))}
	for i, s := range f.Seats {
		t.seats[i] = tableSeat{name: s.Name, token: s.Token}
	}
	i := t.claim(c)
	if i < 0 {
		c.fail(env.ID, protocol.ErrNotInSave, "you did not play this game")
		return
	}
	h.tables = append(h.tables, t)
	h.seatLoaded(c, t, i)
}

// joinLoaded seats a player of a saved game at its open table.
func (h *Hub) joinLoaded(c *Client, env protocol.Envelope, id string) {
	i := slices.IndexFunc(h.tables, func(t *table) bool { return t.loaded != nil && t.loaded.Game == id })
	t := h.tables[i]
	if missing := lacks(c.sets, t.sets); missing != "" {
		c.fail(env.ID, protocol.ErrMissingSets, "you do not have the card set "+missing+" (CD-03)")
		return
	}
	seat := t.claim(c)
	if seat < 0 {
		c.fail(env.ID, protocol.ErrNotInSave, "no seat of this saved game is yours")
		return
	}
	h.seatLoaded(c, t, seat)
}

func (h *Hub) seatLoaded(c *Client, t *table, i int) {
	t.seats[i].client, t.seats[i].ready = c, false
	t.seats[i].token = c.token // the seat now answers to this connection's token
	c.table = t
	h.tableChanged(t)
	h.lobbyChanged()
}

// claim finds the seat of a saved game that is c's: its token, or else
// a free seat with c's nickname. -1: none.
func (t *table) claim(c *Client) int {
	for i, s := range t.seats {
		if s.client == nil && s.token == c.token {
			return i
		}
	}
	for i, s := range t.seats {
		if s.client == nil && s.name == c.name {
			return i
		}
	}
	return -1
}

// resumeLoaded starts a room from a save instead of a new game.
func (r *Room) resumeLoaded() {
	f := r.loaded
	g, err := engine.Load(f.State, r.setup.Sets...)
	if err != nil {
		log.Printf("server: load %s: %v", f.Game, err)
		r.stopped = true
		return
	}
	r.game, r.step, r.steps, r.path = g, f.Step, f.Steps, f.Record
	if r.recDir != "" && f.Record != "" {
		if w, err := record.Append(f.Record, f.Steps, time.Now().UTC()); err == nil {
			r.rec = w
		} else {
			log.Printf("server: record %s: %v", f.Record, err)
		}
	}
	if err := os.Remove(savePath(r.saves, f.Game)); err != nil {
		log.Printf("server: save %s: %v", f.Game, err)
	}
	r.after(nil, false)
}
