package server

import (
	"time"

	"github.com/FourSoulsRepo/four_souls/internal/protocol"
	engine "github.com/FourSoulsRepo/rules_engine"
)

// presence relays what a player points at to the other players (15.35,
// N-15). It never reaches the engine or the record and is never
// answered: a bad or too fast one is dropped.
func (r *Room) presence(c *Client, env protocol.Envelope) {
	if c.seat < 0 || r.game == nil || !r.seats[c.seat].pace.take(time.Now(), r.second) {
		return
	}
	var m protocol.Presence
	if env.Unpack(&m) != nil {
		return
	}
	pub := r.public()
	r.relayPresence(c, protocol.Presence{Seat: c.seat, From: pub.spot(m.From, c.seat), At: pub.spot(m.At, c.seat), Point: onBoard(m.Point)})
}

// relayPresence sends a seat's presence to every other seated player.
// Spectators get none until step 9 decides (N-15).
func (r *Room) relayPresence(from *Client, p protocol.Presence) {
	for _, c := range r.clients {
		if c != from && c.seat >= 0 {
			c.send(protocol.TypePresence, 0, p)
		}
	}
}

// pace limits a player's presence messages to MaxPresence a second,
// with a burst of as many.
type pace struct {
	tokens float64
	at     time.Time
}

// take spends one message and reports whether it may go.
func (p *pace) take(now time.Time, second time.Duration) bool {
	if p.at.IsZero() {
		p.tokens = protocol.MaxPresence
	} else {
		p.tokens = min(protocol.MaxPresence, p.tokens+float64(now.Sub(p.at))/float64(second)*protocol.MaxPresence)
	}
	p.at = now
	if p.tokens < 1 {
		return false
	}
	p.tokens--
	return true
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
	for _, d := range v.Discards {
		add(d...)
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
