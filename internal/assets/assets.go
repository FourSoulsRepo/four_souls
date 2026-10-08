// Package assets gives access to app files (B-01).
//
// Built with -tags embed, files are packed into the binary.
// Built without it (the default, used by wails dev), files are read
// from an assets folder on disk, so they can change without a rebuild.
package assets

// EnvDir overrides the assets folder in external mode.
const EnvDir = "FOUR_SOULS_ASSETS"
