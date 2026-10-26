package main

import (
	"encoding/json"
	"errors"

	"github.com/FourSoulsRepo/four_souls/internal/lan"
	"github.com/FourSoulsRepo/four_souls/internal/protocol"
	"github.com/FourSoulsRepo/four_souls/internal/settings"
	engine "github.com/FourSoulsRepo/rules_engine"
)

// The app's own lines for the player (ADR 021); their words are in
// frontend/src/locales/en.json.
var (
	errNickname     = engine.NewKey("error.nickname_length")
	errSettingValue = engine.NewKey("error.setting_value")
	errLanNoNetwork = engine.NewKey("error.lan_no_network")
	errHosting      = engine.NewKey("error.hosting")
	errDeckRatios   = engine.NewKey("error.deck_ratios")
)

// msgError is an error the player reads in their own language: its text
// is the message as JSON, which the client's messageOf shows (ADR 021).
type msgError struct {
	msg   engine.Msg
	cause error
}

func (e msgError) Error() string {
	data, err := json.Marshal(e.msg)
	if err != nil {
		return e.msg.String()
	}
	return string(data)
}

func (e msgError) Unwrap() error { return e.cause }

// playerError turns an error with words for players into a message;
// any other error stays as it is (an OS or network error, in English).
func playerError(err error) error {
	var m engine.Msg
	switch {
	case err == nil:
		return nil
	case errors.Is(err, settings.ErrNickname):
		m = errNickname.Msg(engine.N("n", protocol.MaxName))
	case errors.Is(err, settings.ErrChoice):
		m = errSettingValue.Msg()
	case errors.Is(err, lan.ErrNoNetwork):
		m = errLanNoNetwork.Msg()
	case errors.Is(err, ErrHosting):
		m = errHosting.Msg()
	default:
		return err
	}
	return msgError{msg: m, cause: err}
}

// msgJSON is a list of messages as JSON, for a binding that returns
// text: the client reads it back (deckRatiosBlocked).
func msgJSON(ms []engine.Msg) string {
	data, err := json.Marshal(ms)
	if err != nil {
		return ""
	}
	return string(data)
}
