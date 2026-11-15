// Package join is the binding of the Join and Online screens: games on
// the local network and what a game server says (ADR 009, ADR 030).
package join

import (
	"context"
	"errors"

	"github.com/FourSoulsRepo/four_souls/internal/lan"
	"github.com/FourSoulsRepo/four_souls/internal/msgerr"
	engine "github.com/FourSoulsRepo/rules_engine"
)

// LAN errors the player reads in their own language (ADR 021); their
// words are in frontend/src/locales/en.json.
var errLanNoNetwork = engine.NewKey("error.lan_no_network")

// api is what the binding needs from the app (ADR 030, item 3).
type api interface {
	// Context is the Wails context; it ends when the app shuts down,
	// which stops a running search or question.
	Context() context.Context
}

// Join finds games to join and asks game servers about them.
type Join struct {
	api api
}

// New makes the binding.
func New(a api) *Join {
	return &Join{api: a}
}

// lanError turns a LAN error with words for players into a message;
// any other error stays as it is (an OS or network error, in English).
func lanError(err error) error {
	var m engine.Msg
	switch {
	case err == nil:
		return nil
	case errors.Is(err, lan.ErrNoNetwork):
		m = errLanNoNetwork.Msg()
	default:
		return err
	}
	return msgerr.New(m, err)
}
