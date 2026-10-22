package server

import (
	"log"
	"slices"
	"time"

	"github.com/FourSoulsRepo/four_souls/internal/protocol"
)

// Idle limits (N-20): a game nobody is connected to ends after a while,
// and every game ends a day after it started, so stuck and endless
// games do not live on the server forever.

// checkLimits ends the game when it ran into an idle limit.
func (r *Room) checkLimits() {
	if r.over.Load() {
		return // a finished game stays for its record download
	}
	now := r.now()
	if slices.ContainsFunc(r.seats, func(s roomSeat) bool { return s.client != nil }) {
		r.emptySince = time.Time{}
	} else if r.emptySince.IsZero() {
		r.emptySince = now
	}
	switch {
	case r.gameLimit > 0 && !r.born.IsZero() && now.Sub(r.born) >= r.gameLimit:
		r.endIdle(protocol.IdleLong)
	case r.emptyLimit > 0 && !r.emptySince.IsZero() && now.Sub(r.emptySince) >= r.emptyLimit:
		r.endIdle(protocol.IdleEmpty)
	}
}

// endIdle ends the game unfinished. Connected players hear why at once,
// the others when they come back; everyone goes back to the lobby.
func (r *Room) endIdle(reason string) {
	n := protocol.Notice{Kind: protocol.NoticeIdle, Game: r.recInfo.game, Reason: reason}
	log.Printf("server: game %s ended by the %s limit", n.Game, reason)
	var away []string
	for _, s := range r.seats {
		if s.client != nil {
			s.client.send(protocol.TypeNotice, 0, n)
		} else if !s.kicked {
			away = append(away, s.token)
		}
	}
	if r.hub != nil {
		r.hub.send(hubRequest{run: func(h *Hub) { h.remember(n, away) }})
	}
	r.stopped = true // Run ends the record and sends everyone to the lobby
}

// checkTables closes a table nobody sat at for the empty limit: only a
// table kept over a restart can be empty. It also forgets old notices.
func (h *Hub) checkTables() {
	h.forgetLost()
	now := h.now()
	closed := false
	h.tables = slices.DeleteFunc(h.tables, func(t *table) bool {
		if slices.ContainsFunc(t.seats, func(s tableSeat) bool { return s.client != nil }) {
			t.emptySince = time.Time{}
			return false
		}
		if t.emptySince.IsZero() {
			t.emptySince = now
		}
		if h.emptyLimit > 0 && now.Sub(t.emptySince) >= h.emptyLimit {
			log.Printf("server: table %s closed by the empty limit", t.id)
			closed = true
			return true
		}
		return false
	})
	if closed {
		h.lobbyChanged()
	}
}
