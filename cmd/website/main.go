//go:build js && wasm

// Command website exposes the rules engine to the rules website.
// It runs in the browser as WebAssembly; there is no web server.
// Build it with -tags embed, so the card data is inside (make site).
package main

import (
	"syscall/js"

	"github.com/FourSoulsRepo/four_souls/internal/website"
)

func main() {
	api := map[string]any{
		"version": export(func(js.Value) string { return website.Version() }),
		"sets":    export(func(js.Value) string { return website.SetsJSON() }),
		"cards":   export(func(js.Value) string { return website.EmbeddedCards() }),
		"check": export(func(args js.Value) string {
			if args.Length() < 1 || args.Index(0).Type() != js.TypeString {
				return `{"error":"check needs the payload as a string"}`
			}
			return website.Check(args.Index(0).String())
		}),
		"notice": export(func(js.Value) string { return website.Notice() }),
	}
	js.Global().Set("fourSouls", js.ValueOf(api))

	// Keep the Go runtime alive so JavaScript can call exported functions.
	select {}
}

// export wraps f for JavaScript: it gets the call's arguments as an
// array and returns a JSON string; a panic becomes an error answer.
func export(f func(args js.Value) string) js.Func {
	return js.FuncOf(func(_ js.Value, args []js.Value) any {
		var out string
		func() {
			defer website.Recover(&out)
			list := make([]any, len(args))
			for i, a := range args {
				list[i] = a
			}
			out = f(js.ValueOf(list))
		}()
		return out
	})
}
