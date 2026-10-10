package server

import (
	"log"
	"net/http"
	"path/filepath"
	"slices"
	"time"

	"github.com/FourSoulsRepo/four_souls/internal/protocol"
	"github.com/FourSoulsRepo/four_souls/internal/version"
	"github.com/FourSoulsRepo/record"
	engine "github.com/FourSoulsRepo/rules_engine"
)

// recordInfo is what a room needs to start its record.
type recordInfo struct {
	game  string
	sets  []string
	cards []record.Card
}

// apply applies an intent and appends it to the record (RP-01, RP-03).
func (r *Room) apply(in engine.Intent) ([]engine.Event, error) {
	st, err := r.game.Step(in)
	if err != nil {
		return nil, err //nolint:wrapcheck // the engine's RuleError is checked by the caller
	}
	r.steps++
	if r.rec != nil {
		if err := r.rec.Step(st.Intent, st.Events, st.Checksum); err != nil {
			log.Printf("server: record %s: %v; recording stops", r.path, err)
			r.rec = nil
		}
	}
	return st.Events, nil
}

// openRecord starts the match record when the game starts (A-12).
func (r *Room) openRecord(s engine.Setup) {
	if r.recDir == "" {
		return
	}
	chars := make([]string, len(s.Characters))
	for i, c := range s.Characters {
		chars[i] = string(c)
	}
	seats := make([]record.Seat, len(r.seats))
	for i, st := range r.seats {
		seats[i] = record.Seat{Seat: i, Name: st.name}
	}
	v := version.Get()
	now := time.Now()
	h := record.Header{
		App: v.App, Engine: v.Engine, Protocol: protocol.Version, Game: r.recInfo.game, Started: now.UTC(),
		Setup: record.Setup{Seed: s.Seed, Players: s.Players, Sets: r.recInfo.sets, BonusSouls: s.BonusSouls, Characters: chars},
		Seats: seats, Cards: r.recInfo.cards,
	}
	w, err := record.Create(r.path, h)
	if err != nil {
		log.Printf("server: no record for %s: %v", r.recInfo.game, err)
		return
	}
	r.rec = w
}

// closeRecord ends the record: finished when the game has a winner,
// unfinished when the server stops first (RP-08).
func (r *Room) closeRecord(finished bool) {
	if r.rec == nil {
		return
	}
	end := record.End{Ended: time.Now().UTC(), Finished: finished}
	if finished {
		for _, w := range r.game.Winners {
			end.Winners = append(end.Winners, int(w))
		}
	}
	if err := r.rec.Close(end); err != nil {
		log.Printf("server: record %s: %v", r.path, err)
	}
	r.rec = nil
	if finished {
		r.over.Store(true)
	}
}

// recordPath names a game's record: started time and game ID.
func recordPath(dir, game string, now time.Time) string {
	return filepath.Join(dir, now.UTC().Format("2006-01-02-150405")+"-"+game+record.Ext)
}

// downloadRecord serves a finished game's record to one of its players
// (RP-05): GET /record?game=g1&token=….
func (h *Hub) downloadRecord(w http.ResponseWriter, req *http.Request) {
	game, token := req.URL.Query().Get("game"), req.URL.Query().Get("token")
	reply := make(chan *Room, 1)
	h.send(hubRequest{find: game, reply: reply})
	var room *Room
	select {
	case room = <-reply:
	case <-h.done:
	}
	switch {
	case room == nil:
		http.Error(w, "no such game", http.StatusNotFound)
	case !slices.Contains(room.tokens, token) || token == "":
		http.Error(w, "only the game's players get its record", http.StatusForbidden)
	case !room.over.Load():
		http.Error(w, "the record is ready when the game is over (RP-05)", http.StatusConflict)
	default:
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Disposition", `attachment; filename="`+filepath.Base(room.path)+`"`)
		http.ServeFile(w, req, room.path)
	}
}

// pruneLoop deletes old records now and every hour (RP-07).
func pruneLoop(done <-chan struct{}, dir string, days int) {
	prune := func() {
		if n, err := record.Prune(dir, days, time.Now()); err != nil {
			log.Printf("server: pruning records: %v", err)
		} else if n > 0 {
			log.Printf("server: deleted %d records older than %d days", n, days)
		}
	}
	prune()
	t := time.NewTicker(time.Hour)
	defer t.Stop()
	for {
		select {
		case <-done:
			return
		case <-t.C:
			prune()
		}
	}
}
