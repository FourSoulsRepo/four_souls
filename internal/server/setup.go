package server

import (
	"fmt"
	"math/rand/v2"
	"slices"
	"time"

	"github.com/FourSoulsRepo/four_souls/internal/protocol"
	engine "github.com/FourSoulsRepo/rules_engine"
)

// banSecond is one second of the ban timer; tests make it shorter.
var banSecond = time.Second

// Limits of the detailed setup.
const (
	maxBanRounds = 3
	minDraft     = 2
	maxDraft     = 5
	minBanTimer  = 5
	maxBanTimer  = 300
)

// matchSetup is the ban and pick phase before a game (6.6, GS-01 to
// GS-04). Only the room goroutine touches it.
type matchSetup struct {
	opts   protocol.Options
	pool   []engine.CardRef // characters still in the game
	banned []protocol.Ban
	phase  string
	round  int   // ban round, from 1
	order  []int // ban turns of this round, by seat
	turn   int   // index into order
	gen    int   // counts ban turns, so a late timer is ignored
	timer  *time.Timer
	end    time.Time
	offers [][]engine.CardRef
	picks  []engine.CardRef
	rng    *rand.Rand
}

// characters lists every character of the sets, in a fixed order.
func characters(sets []engine.CardSet) []engine.CardRef {
	var out []engine.CardRef
	for _, s := range sets {
		for _, d := range s.Cards {
			if d.Kind == engine.CharacterCard {
				out = append(out, d.Ref)
			}
		}
	}
	return out
}

// checkOptions fills in defaults and checks the detailed setup against
// the card sets and the number of players.
func checkOptions(o *protocol.Options, sets []engine.CardSet, players int) error {
	if o.Picking == "" {
		o.Picking = protocol.PickRandom
	}
	if o.Picking != protocol.PickRandom && o.Picking != protocol.PickDraft {
		return fmt.Errorf("picking is %q or %q", protocol.PickRandom, protocol.PickDraft)
	}
	if o.DraftSize == 0 {
		o.DraftSize = 3
	}
	switch {
	case o.Picking == protocol.PickDraft && (o.DraftSize < minDraft || o.DraftSize > maxDraft):
		return fmt.Errorf("a draft offers %d to %d characters", minDraft, maxDraft)
	case o.BanRounds < 0 || o.BanRounds > maxBanRounds:
		return fmt.Errorf("0 to %d ban rounds", maxBanRounds)
	case o.BanTimer != 0 && (o.BanTimer < minBanTimer || o.BanTimer > maxBanTimer):
		return fmt.Errorf("a ban timer is off or %d to %d seconds", minBanTimer, maxBanTimer)
	}
	all := characters(sets)
	for _, b := range o.HostBans {
		if !slices.Contains(all, b) {
			return fmt.Errorf("%s is not a character of these sets", b)
		}
	}
	need := players
	if o.Picking == protocol.PickDraft {
		need = players * o.DraftSize
	}
	if left := len(all) - len(o.HostBans) - o.BanRounds*players; left < need {
		return fmt.Errorf("%d characters left after bans; the setup needs %d", left, need)
	}
	return nil
}

// beginSetup starts the setup phase: host bans out, then ban rounds,
// then picks. Without bans or a draft the characters are dealt at once.
func (r *Room) beginSetup() {
	m := &matchSetup{opts: r.opts, rng: rand.New(rand.NewPCG(r.setup.Seed, r.setup.Seed^0x5eed))} //nolint:gosec // dealing, not security
	for _, c := range characters(r.setup.Sets) {
		if slices.Contains(r.opts.HostBans, c) {
			m.banned = append(m.banned, protocol.Ban{Seat: -1, Card: c}) // GS-02
			continue
		}
		m.pool = append(m.pool, c)
	}
	r.match = m
	if r.opts.BanRounds > 0 {
		m.phase = protocol.PhaseBan
		r.nextBanRound()
		return
	}
	r.deal()
}

// nextBanRound starts a ban round in snake order: 1-2-3-4, then 4-3-2-1
// (GS-04).
func (r *Room) nextBanRound() {
	m := r.match
	m.round++
	m.order = make([]int, len(r.seats))
	for i := range m.order {
		m.order[i] = i
	}
	if m.round%2 == 0 {
		slices.Reverse(m.order)
	}
	m.turn = 0
	r.startBanTurn()
}

// startBanTurn starts the ban timer of the current turn, if any.
func (r *Room) startBanTurn() {
	m := r.match
	m.gen++
	m.end = time.Time{}
	if m.opts.BanTimer > 0 {
		d := time.Duration(m.opts.BanTimer) * banSecond
		m.end = time.Now().Add(d)
		gen := m.gen
		m.timer = time.AfterFunc(d, func() { r.send(request{timeout: gen}) })
	}
	r.broadcastSetup()
}

// ban handles a player's ban on their turn.
func (r *Room) ban(c *Client, env protocol.Envelope) {
	m := r.match
	if m == nil || m.phase != protocol.PhaseBan {
		c.fail(env.ID, protocol.ErrBadMessage, "no ban phase now")
		return
	}
	if c.seat != m.order[m.turn] {
		c.fail(env.ID, protocol.ErrNotYourTurn, "it is another player's ban")
		return
	}
	var b protocol.BanCard
	if err := env.Unpack(&b); err != nil {
		c.fail(env.ID, protocol.ErrBadMessage, err.Error())
		return
	}
	i := slices.Index(m.pool, b.Card)
	if i < 0 {
		c.fail(env.ID, protocol.ErrBadCard, string(b.Card)+" is not in the pool")
		return
	}
	m.pool = slices.Delete(m.pool, i, i+1)
	m.banned = append(m.banned, protocol.Ban{Seat: c.seat, Card: b.Card})
	r.endBanTurn()
}

// banTimeout: the turn's timer ran out; that player makes no ban (GS-04).
func (r *Room) banTimeout(gen int) {
	if m := r.match; m != nil && m.phase == protocol.PhaseBan && m.gen == gen {
		r.endBanTurn()
	}
}

func (r *Room) endBanTurn() {
	m := r.match
	if m.timer != nil {
		m.timer.Stop()
		m.timer = nil
	}
	m.turn++
	switch {
	case m.turn < len(m.order):
		r.startBanTurn()
	case m.round < m.opts.BanRounds:
		r.nextBanRound()
	default:
		r.deal()
	}
}

// deal gives out characters: one each at random, or a draft offer each
// to pick from (GS-01, GS-03).
func (r *Room) deal() {
	m := r.match
	m.rng.Shuffle(len(m.pool), func(i, j int) { m.pool[i], m.pool[j] = m.pool[j], m.pool[i] })
	n := len(r.seats)
	if m.opts.Picking != protocol.PickDraft {
		r.play(m.pool[:n])
		return
	}
	m.phase = protocol.PhasePick
	m.offers = make([][]engine.CardRef, n)
	m.picks = make([]engine.CardRef, n)
	for i := range n {
		m.offers[i] = slices.Clone(m.pool[i*m.opts.DraftSize : (i+1)*m.opts.DraftSize])
	}
	r.broadcastSetup()
}

// pick records a player's draft choice; when all picked, play starts.
func (r *Room) pick(c *Client, env protocol.Envelope) {
	m := r.match
	if m == nil || m.phase != protocol.PhasePick || c.seat < 0 {
		c.fail(env.ID, protocol.ErrBadMessage, "no pick phase now")
		return
	}
	var p protocol.PickCard
	if err := env.Unpack(&p); err != nil {
		c.fail(env.ID, protocol.ErrBadMessage, err.Error())
		return
	}
	if !slices.Contains(m.offers[c.seat], p.Card) {
		c.fail(env.ID, protocol.ErrBadCard, string(p.Card)+" is not one of your offers")
		return
	}
	m.picks[c.seat] = p.Card
	if slices.Contains(m.picks, "") {
		r.broadcastSetup()
		return
	}
	r.play(m.picks)
}

// play creates the game with the chosen characters.
func (r *Room) play(characters []engine.CardRef) {
	s := r.setup
	s.Characters = slices.Clone(characters)
	g, _, err := engine.NewGame(s)
	if err != nil {
		panic(fmt.Sprintf("server: the checked setup failed: %v", err))
	}
	r.game, r.match = g, nil
	r.broadcast(nil)
}

// broadcastSetup sends every client the setup state.
func (r *Room) broadcastSetup() {
	for _, c := range r.clients {
		r.sendSetup(c)
	}
}

func (r *Room) sendSetup(c *Client) {
	m := r.match
	msg := protocol.Setup{
		Phase: m.phase, Round: m.round, Turn: -1, Pool: slices.Clone(m.pool),
		Banned: slices.Clone(m.banned), Seats: r.seatViews(),
	}
	if m.phase == protocol.PhaseBan {
		msg.Turn = m.order[m.turn]
		if !m.end.IsZero() {
			msg.Deadline = m.end.UnixMilli()
		}
	}
	if m.phase == protocol.PhasePick {
		msg.Picked = make([]bool, len(m.picks))
		for i, p := range m.picks {
			msg.Picked[i] = p != ""
		}
		if c.seat >= 0 {
			msg.Offers = m.offers[c.seat] // only your own offers
		}
	}
	c.send(protocol.TypeSetup, 0, msg)
}
