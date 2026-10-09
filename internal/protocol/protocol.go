package protocol

import (
	"encoding/json"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"

	engine "github.com/FourSoulsRepo/rules_engine"
)

// Version is the protocol version (ADR 006). Adding a field keeps it;
// renaming or removing a field of any message, or of the engine's view,
// event or intent JSON, or changing its meaning, bumps it.
const Version = 1

// DefaultPort is the server's TCP port unless configured.
const DefaultPort = 4774

// MaxClientMessage is the largest message a client may send.
const MaxClientMessage = 64 << 10

// Message types.
const (
	TypeHello   = "hello"   // client → server, first message
	TypeWelcome = "welcome" // server → client, answer to hello
	TypeIntent  = "intent"  // client → server: something to do
	TypeResync  = "resync"  // client → server: send me a fresh update
	TypeUpdate  = "update"  // server → client, after every applied step
	TypeError   = "error"   // server → client

	// Lobby (6.5).
	TypeList   = "list"   // client → server: which games are there?
	TypeGames  = "games"  // server → client: the games on this server
	TypeCreate = "create" // client → server: open a new game and sit down
	TypeJoin   = "join"   // client → server: sit down at a game
	TypeReady  = "ready"  // client → server: ready, or not
	TypeLeave  = "leave"  // client → server: leave the table before it starts
	TypeTable  = "table"  // server → client: the table you sit at

	// Match setup (6.6).
	TypeSetup = "setup" // server → client: bans and picks before the game
	TypeBan   = "ban"   // client → server: ban a character on your turn
	TypePick  = "pick"  // client → server: pick one of your offered characters

	// Turn flow (6.7).
	TypeSkipAll = "skip_all" // client → server: pass on the visible stack (N-06)
)

// Envelope wraps every message. ID is set by the client on requests;
// the server copies it into replies.
type Envelope struct {
	Type string          `json:"type"`
	ID   int             `json:"id,omitempty"`
	Data json.RawMessage `json:"data,omitempty"`
}

// Role is how a client takes part.
type Role string

// The roles; spectators come with step 9.
const (
	RolePlayer    Role = "player"
	RoleSpectator Role = "spectator"
	RoleJudge     Role = "judge"
)

// Hello opens a connection. Token is set when reconnecting (N-08).
type Hello struct {
	Protocol int      `json:"protocol"`
	App      string   `json:"app"`
	Name     string   `json:"name"`
	Role     Role     `json:"role"`
	Sets     []string `json:"sets"`            // card set codes the client has (CD-03)
	Token    string   `json:"token,omitempty"` // to reclaim a seat
}

// Welcome accepts a client.
type Welcome struct {
	Protocol int    `json:"protocol"`
	App      string `json:"app"`
	Engine   string `json:"engine"`
	Token    string `json:"token"`
	Seat     int    `json:"seat"` // -1 when not seated
}

// Intent asks the server to do something; the server fills in the
// player from the connection's seat.
type Intent struct {
	Intent engine.Intent `json:"intent"`
}

// Update is the state for one client after a step (N-04): the new
// events, the full view and what this client may do now.
type Update struct {
	Step    int             `json:"step"`
	Events  []engine.Event  `json:"events"`
	View    engine.GameView `json:"view"`
	Allowed []engine.Intent `json:"allowed"`
	Seats   []Seat          `json:"seats"`
	// Deadline is when the response timer answers for the waiting player,
	// Unix ms; 0: no timer (N-07).
	Deadline int64 `json:"deadline,omitempty"`
	// SkipAll is set while this player skips the visible stack (N-06).
	SkipAll bool `json:"skip_all,omitempty"`
}

// SkipAll turns "skip all" on or off (N-06).
type SkipAll struct {
	On bool `json:"on"`
}

// Seat is one seat at the table: who sits there and whether they are
// connected now.
type Seat struct {
	Seat      int    `json:"seat"`
	Name      string `json:"name"`
	Connected bool   `json:"connected"`
	Ready     bool   `json:"ready,omitempty"` // in the lobby
}

// GameInfo is one game in the lobby list.
type GameInfo struct {
	ID      string   `json:"id"`
	Host    string   `json:"host"` // the nickname of who created it
	Seats   int      `json:"seats"`
	Taken   int      `json:"taken"`
	Sets    []string `json:"sets"`
	Started bool     `json:"started"`
}

// Games lists the games on the server.
type Games struct {
	Games []GameInfo `json:"games"`
}

// Create opens a game with 2 to 4 seats and the given card sets.
// Options are the detailed setup; leaving them out is the simple one
// (GS-10).
type Create struct {
	Seats   int      `json:"seats"`
	Sets    []string `json:"sets"`
	Options *Options `json:"options,omitempty"`
}

// Picking modes (GS-01).
const (
	PickRandom = "random" // each player gets a random character
	PickDraft  = "draft"  // each player is offered DraftSize and picks one
)

// Options is the detailed match setup (GS-01 to GS-04, GS-10).
type Options struct {
	Picking   string           `json:"picking,omitempty"`    // PickRandom (default) or PickDraft
	DraftSize int              `json:"draft_size,omitempty"` // characters offered in a draft; default 3
	BanRounds int              `json:"ban_rounds,omitempty"` // each player bans one character per round
	BanTimer  int              `json:"ban_timer,omitempty"`  // seconds per ban; 0: no timer
	HostBans  []engine.CardRef `json:"host_bans,omitempty"`  // never dealt (GS-02)
	// NoBonusSouls plays without the 3 bonus souls (R-SETUP-06).
	NoBonusSouls bool `json:"no_bonus_souls,omitempty"`
	// ResponseTimer is seconds to answer a prompt: 0 is off, else at
	// least 60; an allowance for animations is added (N-07).
	ResponseTimer int `json:"response_timer,omitempty"`
}

// Ban is one banned character; Seat is -1 for the host's list.
type Ban struct {
	Seat int            `json:"seat"`
	Card engine.CardRef `json:"card"`
}

// Setup phases.
const (
	PhaseBan  = "ban"
	PhasePick = "pick"
)

// Setup is the ban and pick phase before the game, sent to each seat
// after every change. Bans are public; offers only go to their player.
type Setup struct {
	Phase    string           `json:"phase"`
	Round    int              `json:"round,omitempty"`    // ban round, from 1
	Turn     int              `json:"turn"`               // the seat that bans now; -1 when everyone picks
	Pool     []engine.CardRef `json:"pool"`               // characters still in the game
	Banned   []Ban            `json:"banned"`             // in order
	Offers   []engine.CardRef `json:"offers,omitempty"`   // your draft choices
	Picked   []bool           `json:"picked,omitempty"`   // who has picked, per seat
	Deadline int64            `json:"deadline,omitempty"` // Unix ms when the ban turn ends; 0: no timer
	Seats    []Seat           `json:"seats"`
}

// BanCard bans a character on the player's ban turn.
type BanCard struct {
	Card engine.CardRef `json:"card"`
}

// PickCard picks one of the player's offered characters.
type PickCard struct {
	Card engine.CardRef `json:"card"`
}

// Join sits down at a game.
type Join struct {
	Game string `json:"game"`
}

// Ready marks the player ready to start, or not.
type Ready struct {
	Ready bool `json:"ready"`
}

// Table is the table a player sits at, before the game starts. The
// game starts by itself when every seat is taken and everyone is ready.
type Table struct {
	Game    string   `json:"game"`
	You     int      `json:"you"` // your seat
	Seats   []Seat   `json:"seats"`
	Sets    []string `json:"sets"`
	Options *Options `json:"options,omitempty"`
}

// Error codes.
const (
	ErrClientOutdated = "client_outdated"
	ErrServerOutdated = "server_outdated"
	ErrBadMessage     = "bad_message"
	ErrRefused        = "refused" // the engine refused an intent; Rule says why
	ErrNotSeated      = "not_seated"
	ErrRoomFull       = "room_full"
	ErrBadName        = "bad_name"
	ErrNoHello        = "no_hello" // the first message must be hello
	ErrNoGame         = "no_game"
	ErrStarted        = "started"      // the game has started already
	ErrMissingSets    = "missing_sets" // the client lacks a card set (CD-03)
	ErrBadSetup       = "bad_setup"
	ErrAtTable        = "at_table" // already sitting at a table
	ErrNotYourTurn    = "not_your_turn"
	ErrBadCard        = "bad_card"
)

// Error reports a problem with a message. Rule is the rules ID when the
// engine refused an intent.
type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Rule    string `json:"rule,omitempty"`
}

// Encode wraps data into an envelope of the given type.
func Encode(typ string, id int, data any) ([]byte, error) {
	raw, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("protocol: encode %s: %w", typ, err)
	}
	out, err := json.Marshal(Envelope{Type: typ, ID: id, Data: raw})
	if err != nil {
		return nil, fmt.Errorf("protocol: encode %s: %w", typ, err)
	}
	return out, nil
}

// Decode reads an envelope.
func Decode(msg []byte) (Envelope, error) {
	var e Envelope
	if err := json.Unmarshal(msg, &e); err != nil {
		return Envelope{}, fmt.Errorf("protocol: %w", err)
	}
	if e.Type == "" {
		return Envelope{}, fmt.Errorf("protocol: message without a type")
	}
	return e, nil
}

// Unpack reads an envelope's data into v.
func (e Envelope) Unpack(v any) error {
	if err := json.Unmarshal(e.Data, v); err != nil {
		return fmt.Errorf("protocol: %s: %w", e.Type, err)
	}
	return nil
}

// CheckVersion compares a client's protocol version with ours and
// returns the error code to send, or "".
func CheckVersion(client int) string {
	switch {
	case client < Version:
		return ErrClientOutdated
	case client > Version:
		return ErrServerOutdated
	default:
		return ""
	}
}

// MaxName is the longest nickname, in characters (ST-06).
const MaxName = 20

// CleanName checks a nickname: trimmed, 1 to MaxName characters, no
// control characters. It returns the cleaned name and whether it is ok.
func CleanName(name string) (string, bool) {
	name = strings.TrimSpace(name)
	n := utf8.RuneCountInString(name)
	if n == 0 || n > MaxName || !utf8.ValidString(name) {
		return "", false
	}
	for _, r := range name {
		if unicode.IsControl(r) {
			return "", false
		}
	}
	return name, true
}

// UniqueName makes name differ from the taken ones: "Ann", then "Ann (2)".
func UniqueName(name string, taken []string) string {
	free := func(s string) bool { return !slices.Contains(taken, s) }
	if free(name) {
		return name
	}
	for i := 2; ; i++ {
		if s := name + " (" + strconv.Itoa(i) + ")"; free(s) {
			return s
		}
	}
}
