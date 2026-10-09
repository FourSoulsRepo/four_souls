// Command cardgen writes engine card stubs from card_db data (step 5.1)
// and the card status report (step 5.2).
//
// From the repo root:
//
//	go run ./cmd/cardgen -set b2   # stubs for one set
//	go run ./cmd/cardgen -report   # rewrite the status report
//
// Stubs go into pkg/rules_engine/cards/<set>/, one file per card, plus
// the set list set_gen.go. A card file that no longer has the stub marker
// is finished and is never touched again.
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	carddb "github.com/FourSoulsRepo/card_db"
)

// reportFile is where -report writes, relative to the repo root.
const reportFile = "pkg/rules_engine/docs/card-status.md"

func main() {
	set := flag.String("set", "", "write stubs for this set code, e.g. b2")
	doReport := flag.Bool("report", false, "rewrite "+reportFile)
	data := flag.String("data", "pkg/card_db", "card_db folder with data/*.json")
	out := flag.String("out", "pkg/rules_engine/cards", "folder of the set packages")
	flag.Parse()
	if *set == "" && !*doReport {
		flag.Usage()
		os.Exit(2)
	}

	db, err := carddb.Load(os.DirFS(*data))
	if err != nil {
		log.Fatal(err)
	}
	if *set != "" {
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
	if *doReport {
		md, err := report(db, *out)
		if err != nil {
			log.Fatal(err)
		}
		if err := os.WriteFile(reportFile, md, 0o600); err != nil {
			log.Fatal(err)
		}
		fmt.Println("wrote", reportFile)
	}
}
