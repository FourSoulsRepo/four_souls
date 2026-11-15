package join

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"

	"github.com/FourSoulsRepo/four_souls/internal/lan"
)

// No network is a message the player reads; other errors stay as they
// are (ADR 021).
func TestLanError(t *testing.T) {
	if lanError(nil) != nil {
		t.Error("nil error gave a message")
	}
	err := lanError(fmt.Errorf("find: %w", lan.ErrNoNetwork))
	var m struct {
		Key string `json:"key"`
	}
	if jerr := json.Unmarshal([]byte(err.Error()), &m); jerr != nil || m.Key != "error.lan_no_network" {
		t.Errorf("no network: %q, %v", err.Error(), jerr)
	}
	other := errors.New("socket")
	if got := lanError(other); !errors.Is(got, other) || got.Error() != "socket" {
		t.Errorf("other error: %v", got)
	}
}
