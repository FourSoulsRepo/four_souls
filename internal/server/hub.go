package server

import (
	"context"
	"fmt"
	"slices"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/FourSoulsRepo/four_souls/internal/protocol"
	"github.com/FourSoulsRepo/four_souls/internal/version"
	"github.com/FourSoulsRepo/record"
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards"
)

// Conn is one client connection. Send queues a message and reports
// false when the client is gone or too slow; Close ends the connection.
// Transports implement it; the server never sees sockets (ADR 006).
type Conn interface {
	Send(msg []byte) bool
	Close()
}

// Client is a connection's handle on the server. The transport passes
// every message it reads to Receive and calls Leave when the connection
// ends. Messages go to the lobby until the client's game starts, then
// to its room.
type Client struct {
	hub  *Hub
	conn Conn
	room atomic.Pointer[Room]

	// Owned by the hub goroutine.
	hello bool
	name  string
	token string
	sets  []string
	table *table

	// Owned by the room goroutine.
	seat int
}

// Receive hands a message from the client to the lobby or its room.
func (c *Client) Receive(msg []byte) {
	if r := c.room.Load(); r != nil {
		r.send(request{client: c, msg: msg})
		return
	}
	c.hub.send(hubRequest{client: c, msg: msg})
}

// Leave tells the server the connection is gone.
func (c *Client) Leave() {
	if r := c.room.Load(); r != nil {
		r.send(request{client: c, leave: true})
	}
	c.hub.send(hubRequest{client: c, leave: true})
}

func (c *Client) send(typ string, id int, data any) {
	msg, err := protocol.Encode(typ, id, data)
	if err != nil {
		panic(err) // our own types always encode
	}
	if !c.conn.Send(msg) {
		c.conn.Close() // too slow: it reconnects and resyncs
	}
}

func (c *Client) fail(id int, code, message string) {
	c.send(protocol.TypeError, id, protocol.Error{Code: code, Message: message})
}

// table is a game in the lobby, before it starts.
type table struct {
	n     int // the game's number on this server
	id    string
	host  string
	sets  []string
	opts  protocol.Options
	seats []tableSeat
}

type tableSeat struct {
	name   string
	token  string
	client *Client
	ready  bool
}

// started is a game that left the lobby.
type started struct {
	info protocol.GameInfo
	room *Room
}

type hubRequest struct {
	client *Client
	msg    []byte
	join   bool
	leave  bool
	find   string     // a started game by ID, for a record download
	reply  chan *Room // the answer to find
	run    func(*Hub) // runs in the hub goroutine (tests)
}

// Hub is the lobby: it seats players at tables and starts their games
// (6.5). Only its goroutine (Run) touches its state.
type Hub struct {
	in      chan hubRequest
	done    chan struct{}
	ctx     context.Context //nolint:containedctx // the rooms started by Run live in it
	rooms   sync.WaitGroup
	sets    []engine.CardSet
	clients []*Client
	tables  []*table
	games   []started
	seq     int
	// second is one second of the rooms' timers; tests shorten it.
	second time.Duration
	// records is the folder for match records; "" records nothing.
	records string
	// cards is the text of every known card, for records (RP-10).
	cards []record.Card
	// seed, if not 0, makes game n's seed seed+n.
	seed uint64
}

// NewHub makes a lobby with the card sets the engine knows.
func NewHub() *Hub {
	return &Hub{in: make(chan hubRequest, 64), done: make(chan struct{}), sets: cards.Sets(), second: time.Second}
}

// Connect adds a connection; its first message must be hello.
func (h *Hub) Connect(conn Conn) *Client {
	c := &Client{hub: h, conn: conn, seat: -1}
	h.send(hubRequest{client: c, join: true})
	return c
}

func (h *Hub) send(req hubRequest) {
	select {
	case h.in <- req:
	case <-h.done:
	}
}

// Run handles the lobby until ctx ends; then it waits for the rooms and
// closes every connection.
func (h *Hub) Run(ctx context.Context) {
	h.ctx = ctx
	defer func() {
		close(h.done)
		h.rooms.Wait()
		for _, c := range h.clients {
			c.conn.Close()
		}
	}()
	for {
		select {
		case <-ctx.Done():
			return
		case req := <-h.in:
			h.handle(req)
		}
	}
}

func (h *Hub) handle(req hubRequest) {
	if req.run != nil {
		req.run(h)
		return
	}
	if req.reply != nil {
		var room *Room
		for _, g := range h.games {
			if g.info.ID == req.find {
				room = g.room
			}
		}
		req.reply <- room
		return
	}
	c := req.client
	switch {
	case req.join:
		h.clients = append(h.clients, c)
		return
	case req.leave:
		h.disconnect(c)
		return
	}
	env, err := protocol.Decode(req.msg)
	if err != nil {
		c.fail(0, protocol.ErrBadMessage, err.Error())
		return
	}
	if !c.hello && env.Type != protocol.TypeHello {
		c.fail(env.ID, protocol.ErrNoHello, "the first message must be hello")
		return
	}
	switch env.Type {
	case protocol.TypeHello:
		h.hello(c, env)
	case protocol.TypeList:
		c.send(protocol.TypeGames, env.ID, h.list())
	case protocol.TypeCreate:
		h.create(c, env)
	case protocol.TypeJoin:
		h.join(c, env)
	case protocol.TypeReady:
		h.ready(c, env)
	case protocol.TypeLeave:
		h.leaveTable(c)
		h.lobbyChanged()
	default:
		c.fail(env.ID, protocol.ErrBadMessage, "not in a game: "+env.Type)
	}
}

// hello checks the version and nickname, then takes the client back to
// its seat (token) or into the lobby.
func (h *Hub) hello(c *Client, env protocol.Envelope) {
	var m protocol.Hello
	if err := env.Unpack(&m); err != nil {
		c.fail(env.ID, protocol.ErrBadMessage, err.Error())
		return
	}
	if code := protocol.CheckVersion(m.Protocol); code != "" {
		c.fail(env.ID, code, fmt.Sprintf("protocol %d, server speaks %d", m.Protocol, protocol.Version))
		c.conn.Close()
		return
	}
	name, ok := protocol.CleanName(m.Name)
	if !ok {
		c.fail(env.ID, protocol.ErrBadName, fmt.Sprintf("a nickname has 1 to %d characters", protocol.MaxName))
		return
	}
	c.hello, c.name, c.sets = true, name, m.Sets
	v := version.Get()
	welcome := protocol.Welcome{Protocol: protocol.Version, App: v.App, Engine: v.Engine, Seat: -1}
	if m.Token != "" {
		if t, i := h.tableSeat(m.Token); t != nil {
			welcome.Token, welcome.Seat = m.Token, i
			c.token, c.table = m.Token, t
			if old := t.seats[i].client; old != nil && old != c {
				old.table = nil
				old.conn.Close()
			}
			t.seats[i].client = c
			c.send(protocol.TypeWelcome, env.ID, welcome)
			h.tableChanged(t)
			return
		}
		for _, g := range h.games {
			if g.room.hasToken(m.Token) {
				c.token = m.Token
				welcome.Token = m.Token
				c.send(protocol.TypeWelcome, env.ID, welcome)
				c.room.Store(g.room)
				g.room.send(request{client: c, attach: m.Token})
				return
			}
		}
	}
	c.token = m.Token // a returning player keeps their token, seated or not
	if !validToken(c.token) {
		c.token = newToken()
	}
	welcome.Token = c.token
	c.send(protocol.TypeWelcome, env.ID, welcome)
	c.send(protocol.TypeGames, 0, h.list())
}

func (h *Hub) tableSeat(token string) (*table, int) {
	for _, t := range h.tables {
		for i, s := range t.seats {
			if s.token == token {
				return t, i
			}
		}
	}
	return nil, -1
}

// create opens a table and seats its creator.
func (h *Hub) create(c *Client, env protocol.Envelope) {
	var m protocol.Create
	if err := env.Unpack(&m); err != nil {
		c.fail(env.ID, protocol.ErrBadMessage, err.Error())
		return
	}
	if c.table != nil {
		c.fail(env.ID, protocol.ErrAtTable, "leave your table first")
		return
	}
	if m.Seats < 2 || m.Seats > 4 {
		c.fail(env.ID, protocol.ErrBadSetup, "a game has 2 to 4 seats (GS-09)")
		return
	}
	if len(m.Sets) == 0 {
		m.Sets = []string{"b2"} // the Base Game (CD-02)
	}
	for _, code := range m.Sets {
		if _, ok := cards.Find(code); !ok {
			c.fail(env.ID, protocol.ErrBadSetup, "this server has no card set "+code)
			return
		}
	}
	if missing := lacks(c.sets, m.Sets); missing != "" {
		c.fail(env.ID, protocol.ErrMissingSets, "you do not have the card set "+missing)
		return
	}
	var opts protocol.Options // the simple setup (GS-10)
	if m.Options != nil {
		opts = *m.Options
	}
	if err := checkOptions(&opts, setsOf(m.Sets), m.Seats); err != nil {
		c.fail(env.ID, protocol.ErrBadSetup, err.Error())
		return
	}
	h.seq++
	t := &table{n: h.seq, id: "g" + strconv.Itoa(h.seq), host: c.name, sets: m.Sets, opts: opts, seats: make([]tableSeat, m.Seats)}
	h.tables = append(h.tables, t)
	h.sit(c, t, 0)
}

// join seats a client at a free seat of a table.
func (h *Hub) join(c *Client, env protocol.Envelope) {
	var m protocol.Join
	if err := env.Unpack(&m); err != nil {
		c.fail(env.ID, protocol.ErrBadMessage, err.Error())
		return
	}
	if c.table != nil {
		c.fail(env.ID, protocol.ErrAtTable, "leave your table first")
		return
	}
	i := slices.IndexFunc(h.tables, func(t *table) bool { return t.id == m.Game })
	if i < 0 {
		if slices.ContainsFunc(h.games, func(g started) bool { return g.info.ID == m.Game }) {
			c.fail(env.ID, protocol.ErrStarted, "that game has started")
		} else {
			c.fail(env.ID, protocol.ErrNoGame, "no game "+m.Game)
		}
		return
	}
	t := h.tables[i]
	if missing := lacks(c.sets, t.sets); missing != "" {
		c.fail(env.ID, protocol.ErrMissingSets, "you do not have the card set "+missing+" (CD-03)")
		return
	}
	free := slices.IndexFunc(t.seats, func(s tableSeat) bool { return s.token == "" })
	if free < 0 {
		c.fail(env.ID, protocol.ErrRoomFull, "every seat is taken")
		return
	}
	h.sit(c, t, free)
}

func (h *Hub) sit(c *Client, t *table, i int) {
	var taken []string
	for _, s := range t.seats {
		if s.name != "" {
			taken = append(taken, s.name)
		}
	}
	t.seats[i] = tableSeat{name: protocol.UniqueName(c.name, taken), token: c.token, client: c}
	c.table = t
	h.tableChanged(t)
	h.lobbyChanged()
}

// ready marks a seat ready; a full table of ready players starts.
func (h *Hub) ready(c *Client, env protocol.Envelope) {
	var m protocol.Ready
	if err := env.Unpack(&m); err != nil {
		c.fail(env.ID, protocol.ErrBadMessage, err.Error())
		return
	}
	t := c.table
	if t == nil {
		c.fail(env.ID, protocol.ErrNotSeated, "sit at a table first")
		return
	}
	for i := range t.seats {
		if t.seats[i].client == c {
			t.seats[i].ready = m.Ready
		}
	}
	for _, s := range t.seats {
		if s.client == nil || !s.ready {
			h.tableChanged(t)
			return
		}
	}
	h.start(t)
}

// start turns a table into a running game.
func (h *Hub) start(t *table) {
	seats := make([]roomSeat, len(t.seats))
	for i, s := range t.seats {
		seats[i] = roomSeat{name: s.name, token: s.token, client: s.client}
	}
	gameSeed := seed()
	if h.seed != 0 {
		gameSeed = h.seed + uint64(t.n) //nolint:gosec // n is positive
	}
	setup := engine.Setup{Seed: gameSeed, Players: len(t.seats), Sets: setsOf(t.sets), BonusSouls: !t.opts.NoBonusSouls}
	room, err := NewRoom(setup, t.opts, seats)
	if err == nil {
		room.second = h.second
		room.recDir = h.records
		room.recInfo = recordInfo{game: t.id, sets: t.sets, cards: h.cards}
		room.path = recordPath(h.records, t.id, time.Now())
	}
	if err != nil {
		for _, s := range t.seats {
			s.client.fail(0, protocol.ErrBadSetup, err.Error())
		}
		return
	}
	h.tables = slices.DeleteFunc(h.tables, func(x *table) bool { return x == t })
	h.games = append(h.games, started{info: protocol.GameInfo{ID: t.id, Host: t.host, Seats: len(t.seats), Taken: len(t.seats), Sets: t.sets, Started: true}, room: room})
	for _, s := range t.seats {
		s.client.table = nil
		s.client.room.Store(room)
	}
	h.rooms.Add(1)
	go func() {
		defer h.rooms.Done()
		room.Run(h.ctx)
	}()
	h.lobbyChanged()
}

// leaveTable frees the client's seat; an empty table closes.
func (h *Hub) leaveTable(c *Client) {
	t := c.table
	if t == nil {
		return
	}
	c.table = nil
	for i := range t.seats {
		if t.seats[i].client == c {
			t.seats[i] = tableSeat{}
		}
	}
	if slices.IndexFunc(t.seats, func(s tableSeat) bool { return s.token != "" }) < 0 {
		h.tables = slices.DeleteFunc(h.tables, func(x *table) bool { return x == t })
		return
	}
	h.tableChanged(t)
}

// disconnect forgets a connection. In the lobby its seat is freed.
func (h *Hub) disconnect(c *Client) {
	h.leaveTable(c)
	h.clients = slices.DeleteFunc(h.clients, func(x *Client) bool { return x == c })
	h.lobbyChanged()
}

// tableChanged sends the table to everyone sitting at it.
func (h *Hub) tableChanged(t *table) {
	seats := make([]protocol.Seat, len(t.seats))
	for i, s := range t.seats {
		seats[i] = protocol.Seat{Seat: i, Name: s.name, Connected: s.client != nil, Ready: s.ready}
	}
	for i, s := range t.seats {
		if s.client != nil {
			s.client.send(protocol.TypeTable, 0, protocol.Table{Game: t.id, You: i, Seats: seats, Sets: t.sets, Options: &t.opts})
		}
	}
}

// lobbyChanged sends the game list to everyone in the lobby.
func (h *Hub) lobbyChanged() {
	list := h.list()
	for _, c := range h.clients {
		if c.hello && c.table == nil && c.room.Load() == nil {
			c.send(protocol.TypeGames, 0, list)
		}
	}
}

func (h *Hub) list() protocol.Games {
	out := protocol.Games{Games: []protocol.GameInfo{}}
	for _, t := range h.tables {
		taken := 0
		for _, s := range t.seats {
			if s.token != "" {
				taken++
			}
		}
		out.Games = append(out.Games, protocol.GameInfo{ID: t.id, Host: t.host, Seats: len(t.seats), Taken: taken, Sets: t.sets})
	}
	for _, g := range h.games {
		out.Games = append(out.Games, g.info)
	}
	return out
}

// lacks returns a set of need that have does not list, or "".
func lacks(have, need []string) string {
	for _, s := range need {
		if !slices.Contains(have, s) {
			return s
		}
	}
	return ""
}

// seed is a new game's seed; the engine itself never reads the clock.
func seed() uint64 { return uint64(time.Now().UnixNano()) } //nolint:gosec // a game seed, not a secret

// setsOf looks up card sets by code; the codes were checked already.
func setsOf(codes []string) []engine.CardSet {
	var out []engine.CardSet
	for _, code := range codes {
		if set, ok := cards.Find(code); ok {
			out = append(out, set)
		}
	}
	return out
}
