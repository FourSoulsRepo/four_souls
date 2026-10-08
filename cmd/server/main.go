// Command server runs Four Souls games without a UI.
package main

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/FourSoulsRepo/four_souls/internal/legal"
	"github.com/FourSoulsRepo/four_souls/internal/version"
)

func main() {
	if err := printBanner(os.Stdout); err != nil {
		log.Fatal(err)
	}
}

// printBanner writes the fan-game notice and the versions (L-02).
func printBanner(w io.Writer) error {
	v := version.Get()
	_, err := fmt.Fprintf(w, "%s\nFour Souls dedicated server\napp %s, rules engine %s\n",
		legal.PlainText(), v.App, v.Engine)
	if err != nil {
		return fmt.Errorf("print banner: %w", err)
	}
	return nil
}
