package server

import (
	"time"

	"github.com/FourSoulsRepo/four_souls/internal/protocol"
	engine "github.com/FourSoulsRepo/rules_engine"
)

// presence relays what a player points at to everyone else at the
// game (15.35, N-15). It never reaches the engine or the record and is
// never answered: a bad one is dropped, and so is every one while the
// server has presence off. One that comes too soon waits (sample).
func (r *Room) presence(c *Client, env protocol.Envelope) {
	if !r.presenceOn || c.seat < 0 || r.game == nil {
		return
	}
	var m protocol.Presence
	if env.Unpack(&m) != nil {
		return
	}
	pub := r.public()
	r.sample(c, protocol.Presence{
		Seat:  c.seat,
		From:  pub.spot(m.From, c.seat),
		At:    pub.spot(m.At, c.seat),
		Point: onBoard(m.Point),
		Held:  pub.spot(m.Held, c.seat),
	})
}

// relayPresence sends a seat's presence to every other client of the
// game: the other players and, from step 9, spectators (N-15).
func (r *Room) relayPresence(from *Client, p protocol.Presence) {
	if !r.presenceOn {
		return
	}
	for _, c := range r.clients {
		if c != from {
			c.send(protocol.TypePresence, 0, p)
		}
	}
}

// sampler paces one seat's presence (20.13, ADR 027 "Other players
// live"): at most one per gap. One that comes sooner waits for the gap
// to end; a newer one replaces it, so the last one always goes.
type sampler struct {
	last  time.Time          // when the last one went
	next  *protocol.Presence // the one waiting, if any
	from  *Client            // its sender
	timer *time.Timer        // ends the wait
}

// presenceGap is the least time after one presence before p may go:
// shorter while aiming, so the arrow follows closely. second is the
// room's second; tests shorten it.
func presenceGap(p protocol.Presence, second time.Duration) time.Duration {
	gap := protocol.PresenceGap
	if p.From != nil {
		gap = protocol.PresenceAimGap
	}
	return time.Duration(float64(second) * gap.Seconds())
}

// sample relays p now, if its gap has passed; else it waits for it.
func (r *Room) sample(c *Client, p protocol.Presence) {
	s := &r.seats[c.seat].sample
	s.next, s.from = &p, c
	if s.timer == nil {
		r.flushPresence(c.seat)
	}
}

// flushPresence relays a seat's waiting presence once its gap has
// passed, else sets a timer for when it does.
func (r *Room) flushPresence(seat int) {
	s := &r.seats[seat].sample
	s.timer = nil
	if s.next == nil || s.from == nil || s.from.seat != seat {
		s.next, s.from = nil, nil
		return
	}
	now := time.Now()
	if wait := s.last.Add(presenceGap(*s.next, r.second)).Sub(now); wait > 0 {
		s.timer = time.AfterFunc(wait, func() { r.send(request{presence: seat + 1}) })
		return
	}
	r.relayPresence(s.from, *s.next)
	s.last, s.next, s.from = now, nil, nil
}

// dropPresence forgets a seat's waiting presence.
func (r *Room) dropPresence(seat int) {
	s := &r.seats[seat].sample
	if s.timer != nil {
		s.timer.Stop()
	}
	s.timer, s.next, s.from = nil, nil, nil
}

// publicTable is what every player can see now: the spots a presence
// message may name.
type publicTable struct {
	objects map[engine.ObjectID]bool
	stack   map[int]bool
	hands   []int // hand sizes, per seat
}

// public is the table as a spectator sees it, kept for the step.
func (r *Room) public() *publicTable {
	if r.pub != nil && r.pubGame == r.game && r.pubStep == r.step {
		return r.pub
	}
	v := r.game.View(engine.Viewer{Kind: engine.ViewSpectator})
	pub := &publicTable{objects: map[engine.ObjectID]bool{}, stack: map[int]bool{}}
	add := func(cards ...engine.CardView) {
		for _, cv := range cards {
			pub.objects[cv.ID] = true
		}
	}
	slots := func(slots []engine.SlotView) {
		for _, s := range slots {
			if s.Top != nil {
				add(*s.Top)
			}
			add(s.Covered...)
		}
	}
	for _, p := range v.Players {
		add(p.Character)
		add(p.InPlay...)
		add(p.Hand...) // empty for a spectator; a public hand would show here
		pub.hands = append(pub.hands, p.HandSize)
	}
	slots(v.Shop)
	slots(v.Monsters)
	slots(v.Rooms)
	for _, d := range v.Discards {
		if d != nil {
			add(*d) // only the top card of a pile is in the view
		}
	}
	add(v.BonusSouls...)
	if v.Revealed != nil {
		add(*v.Revealed)
	}
	for _, it := range v.Stack {
		pub.stack[it.ID] = true
		if it.Source != 0 {
			pub.objects[it.Source] = true
		}
	}
	r.pub, r.pubGame, r.pubStep = pub, r.game, r.step
	return pub
}

// spot keeps s if it names exactly one thing every player can see, or
// a slot of the sender's own hand; else nil (A-16).
func (pub *publicTable) spot(s *protocol.Spot, seat int) *protocol.Spot {
	if s == nil {
		return nil
	}
	set, ok := 0, true
	if s.Object != 0 {
		set++
		ok = ok && pub.objects[s.Object]
	}
	if s.Hand != 0 {
		set++
		ok = ok && s.Hand > 0 && s.Hand <= pub.hands[seat]
	}
	if s.Player != nil {
		set++
		ok = ok && *s.Player >= 0 && int(*s.Player) < len(pub.hands)
	}
	if s.Deck != nil {
		set++
		ok = ok && *s.Deck >= engine.TreasureDeck && *s.Deck <= engine.RoomDeck
	}
	if s.StackID != 0 {
		set++
		ok = ok && pub.stack[s.StackID]
	}
	if set != 1 || !ok {
		return nil
	}
	out := *s
	return &out
}

// onBoard keeps a point on the board frame, moved inside it.
func onBoard(p *protocol.Point) *protocol.Point {
	if p == nil {
		return nil
	}
	return &protocol.Point{X: min(max(p.X, 0), 1), Y: min(max(p.Y, 0), 1)}
}
