package msgerr

import (
	"errors"
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

// The text is the message as JSON; the cause stays reachable.
func TestError(t *testing.T) {
	cause := errors.New("disk full")
	err := New(engine.RefuseDeckRatios.Msg(engine.N("n", 3)), cause)
	if got, want := err.Error(), `{"key":"refuse.deck_ratios","args":{"n":3}}`; got != want {
		t.Errorf("text %s, want %s", got, want)
	}
	if !errors.Is(err, cause) {
		t.Error("cause lost")
	}
	var m Error
	if !errors.As(err, &m) || m.Msg.Key != engine.RefuseDeckRatios {
		t.Errorf("message lost: %v", err)
	}
}

func TestJSON(t *testing.T) {
	ms := []engine.Msg{engine.RefuseDeckRatios.Msg(engine.N("n", 2)), engine.RefuseDeckRatios.Msg()}
	if got, want := JSON(ms), `[{"key":"refuse.deck_ratios","args":{"n":2}},{"key":"refuse.deck_ratios"}]`; got != want {
		t.Errorf("JSON %s, want %s", got, want)
	}
	if got := JSON(nil); got != "null" {
		t.Errorf("no messages: %s", got)
	}
}
