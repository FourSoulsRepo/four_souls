package main

import (
	"errors"

	"github.com/FourSoulsRepo/four_souls/internal/lan"
	"github.com/FourSoulsRepo/four_souls/internal/msgerr"
	engine "github.com/FourSoulsRepo/rules_engine"
)

// The app's own lines for the player (ADR 021); their words are in
// frontend/src/locales/en.json.
var (
	errLanNoNetwork = engine.NewKey("error.lan_no_network")
	errDeckRatios   = engine.NewKey("error.deck_ratios")
)

// playerError turns an error with words for players into a message;
// any other error stays as it is (an OS or network error, in English).
func playerError(err error) error {
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
