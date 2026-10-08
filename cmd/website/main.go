//go:build js && wasm

// Command website exposes the rules engine to the rules website.
// It runs in the browser as WebAssembly; there is no web server.
package main

func main() {
	// Keep the Go runtime alive so JavaScript can call exported functions.
	select {}
}
