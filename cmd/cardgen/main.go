// Command cardgen writes engine card stubs from card_db data (step 5.1).
//
// From the repo root:
//
//	go run ./cmd/cardgen -set b2
//
// It writes one file per card into pkg/rules_engine/cards/<set>/ and the
// set list set_gen.go. A card file that no longer has the stub marker is
// finished and is never touched again.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	carddb "github.com/FourSoulsRepo/card_db"
)

func main() {
	set := flag.String("set", "b2", "set code, e.g. b2")
	data := flag.String("data", "pkg/card_db", "card_db folder with data/*.json")
	out := flag.String("out", "pkg/rules_engine/cards", "folder of the set packages")
	flag.Parse()

	db, err := carddb.Load(os.DirFS(*data))
	if err != nil {
		log.Fatal(err)
	}
	files, err := generate(db, *set)
	if err != nil {
		log.Fatal(err)
	}
	dir := filepath.Join(*out, *set)
	written, kept, err := write(dir, files)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("%s: %d stubs written, %d finished cards kept\n", dir, written, kept)
}
