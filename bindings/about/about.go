// Package about is the binding of the app's versions and the fan-game
// notice (ADR 030).
package about

import (
	"github.com/FourSoulsRepo/four_souls/internal/legal"
	"github.com/FourSoulsRepo/four_souls/internal/version"
)

// About tells the main menu and the splash screen about the app.
type About struct{}

// New makes the binding.
func New() *About { return &About{} }

// Versions returns the app and rules engine versions for the main menu.
func (*About) Versions() version.Info {
	return version.Get()
}

// Notice returns the fan-game notice for the splash screen.
func (*About) Notice() []legal.Line {
	return legal.Notice
}
