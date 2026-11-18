//go:build !linux || !cgo

package diag

// webkitVersion is unknown here: no WebKitGTK, or no cgo to ask it.
func webkitVersion() string { return "" }
