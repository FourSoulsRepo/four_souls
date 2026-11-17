// Package msgerr carries messages for the player out of the app's
// bindings (ADR 021, ADR 030 item 6). The client renders a message in
// the player's language; each binding package keeps its own keys.
package msgerr

import (
	"encoding/json"

	engine "github.com/FourSoulsRepo/rules_engine"
)

// Error is an error the player reads in their own language: its text
// is the message as JSON, which the client's messageOf shows.
type Error struct {
	Msg   engine.Msg
	Cause error // kept for errors.Is and the log
}

// New wraps cause in the message m.
func New(m engine.Msg, cause error) error { return Error{Msg: m, Cause: cause} }

func (e Error) Error() string {
	data, err := json.Marshal(e.Msg)
	if err != nil {
		return e.Msg.String()
	}
	return string(data)
}

// Unwrap gives the cause.
func (e Error) Unwrap() error { return e.Cause }

// JSON is a list of messages as JSON, for a binding that returns text
// the client reads back (deckRatiosBlocked); "" when it can't encode.
func JSON(ms []engine.Msg) string {
	data, err := json.Marshal(ms)
	if err != nil {
		return ""
	}
	return string(data)
}
