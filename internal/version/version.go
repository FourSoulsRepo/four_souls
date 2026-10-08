// Package version reports the app and rules engine versions.
package version

import engine "github.com/FourSoulsRepo/rules_engine"

// App is the app version, stamped at link time:
//
//	-ldflags "-X github.com/FourSoulsRepo/four_souls/internal/version.App=1.2.3"
var App = "dev"

// Info holds both versions shown to users.
type Info struct {
	App    string `json:"app"`
	Engine string `json:"engine"`
}

// Get returns the current versions.
func Get() Info {
	return Info{App: App, Engine: engine.Version}
}
