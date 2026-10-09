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
}

// Seat is one seat at the table: who sits there and whether they are
// connected now.
type Seat struct {
	Seat      int    `json:"seat"`
	Name      string `json:"name"`
	Connected bool   `json:"connected"`
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
