//go:build js && wasm

// Command website exposes the rules engine to the rules website.
// It runs in the browser as WebAssembly; there is no web server.
package main

import (
	"syscall/js"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func main() {
	js.Global().Set("fourSoulsEngineVersion", engine.Version)

	// Keep the Go runtime alive so JavaScript can call exported functions.
	select {}
}
