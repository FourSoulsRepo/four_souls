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
// 2: the order prompt (19.16); an older client cannot answer it.
// 3: the throw prompt and intent (15.60).
// 4: no "skip all" message or field (15.84).
// 5: rollback (ADR 020); an older client would show a rolled back
// table as if the undone steps still stood.
const Version = 5

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
	TypeWaitTurn = "wait_turn" // client → server: pass until the turn before mine ends (N-16)

	// Disconnects (6.8).
	TypeVote = "vote" // client → server: wait for the missing player, or kick them

	// Save and continue later (6.11).
	TypeSave  = "save"  // client → server: the host saves the game; it stops
	TypeSaved = "saved" // server → client: the game was saved under this ID
	TypeLoad  = "load"  // client → server: open a saved game's table again

	// Live aiming and hover (15.35, N-15): never an intent or an event.
	TypePresence = "presence" // client → server → the other players

	// Server restarts and game ends (N-20). An older client ignores them.
	TypeRestarting = "restarting" // server → client: back in a moment; reconnect
	TypeNotice     = "notice"     // server → client: something to tell the player

	// Rollback (ADR 020).
	TypeRollback = "rollback" // client → server: undo a step and every step after it
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
	Dice     string   `json:"dice,omitempty"`  // the player's dice look (V-19)
}

// Welcome accepts a client.
type Welcome struct {
	Protocol int    `json:"protocol"`
	App      string `json:"app"`
	Engine   string `json:"engine"`
	Token    string `json:"token"`
	Seat     int    `json:"seat"` // -1 when not seated
	// Create is set when this client may open a game here: always on a
	// dedicated server, only for the host's own app on a hosted one (N-17).
	Create bool `json:"create,omitempty"`
	// Presence is set when the server relays live aiming and hover; a
	// client sends none otherwise (15.35, ADR 016).
	Presence bool `json:"presence,omitempty"`
	// Name is the server's name, set by its operator; the app shows it
	// instead of the address (N-18). Empty on an app's own host.
	Name string `json:"name,omitempty"`
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
	// Aims are the targets an allowed play or use asks first, by its
	// index in Allowed, to light while a card is dragged (15.25).
	Aims []engine.Aim `json:"aims,omitempty"`
	// Blocked says why this player's other cards and turn actions can't
	// be used now (15.29).
	Blocked []engine.Block `json:"blocked,omitempty"`
	Seats   []Seat         `json:"seats"`
	// Deadline is when the response timer answers for the waiting player,
	// Unix ms; 0: no timer (N-07).
	Deadline int64 `json:"deadline,omitempty"`
	// WaitTurn is set while this player waits for their turn (N-16).
	WaitTurn bool `json:"wait_turn,omitempty"`
	// Pause is set while the game waits for a disconnected player (N-08).
	Pause *Pause `json:"pause,omitempty"`
	// RollbackMode is the game's rollback mode (ADR 020); "": none.
	RollbackMode string `json:"rollback_mode,omitempty"`
	// RollbackFrom is the first step a rollback may target: the first
	// since the base state (ADR 020, item 7). 0 while none may.
	RollbackFrom int `json:"rollback_from,omitempty"`
	// Rollback is set on the update a rollback made: the table is the
	// one from before step To. Step is the rollback's own number.
	Rollback *Rollback `json:"rollback,omitempty"`
}

// Rollback tells who rolled the game back to before step To, and how
// many live steps that undid (ADR 020, item 8).
type Rollback struct {
	Seat  int `json:"seat"`
	To    int `json:"to"`
	Undid int `json:"undid"`
}

// RollbackTo asks to undo step Step and every live step after it.
type RollbackTo struct {
	Step int `json:"step"`
}

// Pause is a game waiting for disconnected players, and the vote.
type Pause struct {
	Away  []int  `json:"away"`  // seats whose player is gone
	Votes []Vote `json:"votes"` // the connected players' votes so far
	// Grace is when a kick vote opens, Unix ms, after a server restart:
	// the players are still reconnecting (N-20). 0: open now.
	Grace int64 `json:"grace,omitempty"`
}

// Vote is one player's vote while the game is paused: Kick true kicks
// the missing players, false waits longer.
type Vote struct {
	Seat int  `json:"seat"`
	Kick bool `json:"kick"`
}

// WaitTurn turns "wait for my turn" on or off (N-16).
type WaitTurn struct {
	On bool `json:"on"`
}

// Seat is one seat at the table: who sits there and whether they are
// connected now.
type Seat struct {
	Seat      int    `json:"seat"`
	Name      string `json:"name"`
	Connected bool   `json:"connected"`
	Ready     bool   `json:"ready,omitempty"`  // in the lobby
	Kicked    bool   `json:"kicked,omitempty"` // the seat is empty for good (N-08)
	// Dice is the seat's dice look (V-19); "" draws the default one.
	Dice string `json:"dice,omitempty"`
}

// Presence is what a player points at now (N-15): the card they aim
// from, and the card or spot the arrow or pointer is over. The engine
// never sees it and records never keep it. Each one replaces the
// sender's last; an empty one means nothing. The server fills in Seat,
// keeps only spots every player can see, and drops the rest.
type Presence struct {
	Seat  int    `json:"seat"`            // the sender; set by the server
	From  *Spot  `json:"from,omitempty"`  // what aims; nil: not aiming
	At    *Spot  `json:"at,omitempty"`    // what the arrow or pointer is over
	Point *Point `json:"point,omitempty"` // the pointer, while over no spot
}

// Spot is one thing on the table a player points at; exactly one field
// is set. A card in the sender's hand is its slot, never its object, so
// the others learn nothing about it (A-16).
type Spot struct {
	Object  engine.ObjectID  `json:"object,omitempty"`   // a card every player can see
	Hand    int              `json:"hand,omitempty"`     // the sender's hand slot, from 1
	Player  *engine.PlayerID `json:"player,omitempty"`   // a player
	Deck    *engine.DeckKind `json:"deck,omitempty"`     // a deck
	StackID int              `json:"stack_id,omitempty"` // a stack item
}

// Point is a spot on the board frame, 0 to 1 from its top left corner.
type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

// MaxPresence is how many presence messages a player may send each
// second; the server drops the rest.
const MaxPresence = 20

// GameInfo is one game in the lobby list.
type GameInfo struct {
	ID      string   `json:"id"`
	Host    string   `json:"host"` // the nickname of who created it
	Seats   int      `json:"seats"`
	Taken   int      `json:"taken"`
	Sets    []string `json:"sets"`
	Started bool     `json:"started"`
	Locked  bool     `json:"locked,omitempty"` // joining needs the password (N-14)
	// Rollback is the game's rollback mode (ADR 020, item 1.4).
	Rollback string `json:"rollback,omitempty"`
}

// SeatInfo answers GET /seat?token=… (15.99): whether the token still
// holds a seat at a table or in a game that has not finished, so a
// restarted app offers Rejoin. Only the token's holder knows it.
type SeatInfo struct {
	Live    bool   `json:"live"`
	Game    string `json:"game,omitempty"`
	Host    string `json:"host,omitempty"` // the nickname of who created it
	Started bool   `json:"started,omitempty"`
}

// ServerStatus answers GET /status (15.100): how many games wait for
// players, so the Online screen's Recent list shows the server is up.
// Any lobby client sees as much.
type ServerStatus struct {
	Open int `json:"open"`
}

// Games lists the games on the server, and the saved ones.
type Games struct {
	Games []GameInfo `json:"games"`
	Saves []SaveInfo `json:"saves,omitempty"`
}

// SaveInfo is one saved game (N-11): who played, when it was saved.
type SaveInfo struct {
	ID     string   `json:"id"`
	Host   string   `json:"host"`
	Seats  []string `json:"seats"` // nicknames
	Saved  int64    `json:"saved"` // Unix ms
	Engine string   `json:"engine"`
}

// Saved tells the players their game was saved and stopped.
type Saved struct {
	Save string `json:"save"`
}

// Notice kinds (N-20).
const (
	// NoticeEngine: the game went on after a restart on a newer engine.
	NoticeEngine = "engine_updated"
	// NoticeLost: the new server could not read the game; it ended.
	NoticeLost = "game_lost"
	// NoticeIdle: the game ran into an idle limit and ended.
	NoticeIdle = "game_idle"
)

// Why a game ran into an idle limit (NoticeIdle).
const (
	IdleEmpty = "empty" // nobody was connected for too long
	IdleLong  = "long"  // the game ran for the longest time allowed
)

// Notice tells the player something the game screen does not show.
type Notice struct {
	Kind   string `json:"kind"`
	Game   string `json:"game,omitempty"`
	Reason string `json:"reason,omitempty"` // of NoticeIdle: IdleEmpty or IdleLong
}

// Load opens a saved game's table; its players take their seats back.
type Load struct {
	Save string `json:"save"`
}

// Create opens a game with 2 to 4 seats and the given card sets.
// Options are the detailed setup; leaving them out is the simple one
// (GS-10).
type Create struct {
	Seats   int      `json:"seats"`
	Sets    []string `json:"sets"`
	Options *Options `json:"options,omitempty"`
	// Password, if set, is needed to join (N-14). It is not an option,
	// so tables, records and saves never carry it.
	Password string `json:"password,omitempty"`
}

// MaxPassword is the longest game password, in bytes.
const MaxPassword = 64

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
	HostBans  []engine.CardRef `json:"host_bans,omitempty"`  // any card; never in the game (GS-02, GS-13)
	// NoBonusSouls plays without the 3 bonus souls (R-SETUP-06).
	NoBonusSouls bool `json:"no_bonus_souls,omitempty"`
	// NoDeckRatios plays with whole decks instead of the deck ratios
	// (R-SETUP-17, GS-07). The server sets it when the chosen sets have
	// too few cards for the ratios.
	NoDeckRatios bool `json:"no_deck_ratios,omitempty"`
	// ResponseTimer is seconds to answer a prompt: 0 is off, else at
	// least 60; an allowance for animations is added (N-07).
	ResponseTimer int `json:"response_timer,omitempty"`
	// Rollback is the rollback mode (ADR 020): "" is Disabled, else
	// RollbackSoft or RollbackGod.
	Rollback string `json:"rollback,omitempty"`
}

// Rollback modes (ADR 020, item 1).
const (
	RollbackSoft = "soft" // the other players vote; barriers stop it (17.4)
	RollbackGod  = "god"  // Godmode: anyone, any step since the base, no vote
)

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
	Game     string `json:"game"`
	Password string `json:"password,omitempty"` // for a locked game (N-14)
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
	// Loaded is the save the table continues; its seats wait for their
	// players (N-11).
	Loaded string `json:"loaded,omitempty"`
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
	ErrPaused         = "paused" // the game waits for a disconnected player
	ErrNotHost        = "not_host"
	ErrNoSave         = "no_save"
	ErrSaveVersion    = "save_version" // the save is from another engine version
	ErrNotInSave      = "not_in_save"  // the player did not play the saved game
	ErrBadCard        = "bad_card"
	ErrPassword       = "password" // missing or wrong game password (N-14)
	ErrRollback       = "rollback" // the rollback is not allowed now (ADR 020)
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

// MaxServerName is the longest server name, in characters (N-18).
const MaxServerName = 40

// CleanName checks a nickname: trimmed, 1 to MaxName characters, no
// control characters. It returns the cleaned name and whether it is ok.
func CleanName(name string) (string, bool) {
	return cleanText(name, MaxName)
}

// CleanServerName checks a server name as CleanName does a nickname,
// with up to MaxServerName characters (N-18).
func CleanServerName(name string) (string, bool) {
	return cleanText(name, MaxServerName)
}

// cleanText trims name and checks it has 1 to most characters and no
// control characters.
func cleanText(name string, most int) (string, bool) {
	name = strings.TrimSpace(name)
	n := utf8.RuneCountInString(name)
	if n == 0 || n > most || !utf8.ValidString(name) {
		return "", false
	}
	for _, r := range name {
		if unicode.IsControl(r) {
			return "", false
		}
	}
	return name, true
}

// MaxDice is the longest dice look ID, in bytes.
const MaxDice = 32

// CleanDice checks a dice look ID from hello: up to MaxDice lowercase
// letters, digits, '_' or '-'. Anything else is "", the default look.
// The server keeps unknown IDs; each client draws them plain (V-19).
func CleanDice(id string) string {
	if len(id) > MaxDice {
		return ""
	}
	for _, r := range id {
		if (r < 'a' || r > 'z') && (r < '0' || r > '9') && r != '_' && r != '-' {
			return ""
		}
	}
	return id
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
