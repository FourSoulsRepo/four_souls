package version

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestGet(t *testing.T) {
	got := Get()
	if got.App != App {
		t.Errorf("App = %q, want %q", got.App, App)
	}
	if got.Engine != engine.Version {
		t.Errorf("Engine = %q, want %q", got.Engine, engine.Version)
	}
}
