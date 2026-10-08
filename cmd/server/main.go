// Command server runs Four Souls games without a UI.
package main

import (
	"fmt"
	"os"

	"github.com/FourSoulsRepo/four_souls/internal/version"
)

func main() {
	v := version.Get()
	fmt.Fprintf(os.Stdout, "Four Souls dedicated server\napp %s, rules engine %s\n", v.App, v.Engine)
}
