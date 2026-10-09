package main

import (
	"bytes"
	"errors"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"unicode"

	carddb "github.com/FourSoulsRepo/card_db"
)

// StubMarker marks a card file that is still a stub. cardgen rewrites
// such files; the card author deletes the line when the card is done.
const StubMarker = "TODO(card):"

// setFile is the generated list of every card in the set; docName is
// the package comment, written once.
const (
	setFile = "set_gen.go"
	docName = "doc.go"
)

// kinds maps printed card types to engine card kinds.
var kinds = map[string]string{
	"Character Card": "CharacterCard",

	"Passive Treasure Card":      "TreasureCard",
	"Active Treasure Card":       "TreasureCard",
	"Paid Treasure Card":         "TreasureCard",
	"One-Use Treasure Card":      "TreasureCard",
	"Eternal Treasure Card":      "TreasureCard",
	"Soul Treasure Card":         "TreasureCard",
	"Wildcard Card":              "LootCard",
	"Trinket Card":               "LootCard",
	"Pill/Rune Card":             "LootCard",
	"Nickel Card":                "LootCard",
	"Bomb Card":                  "LootCard",
	"Battery Card":               "LootCard",
	"Butter Bean Card":           "LootCard",
	"Lost Soul Card":             "LootCard",
	"1 Cent Card":                "LootCard",
	"2 Cent Card":                "LootCard",
	"3 Cent Card":                "LootCard",
	"4 Cent Card":                "LootCard",
	"Dice Shard/Soul Heart Card": "LootCard",

	"Basic Monster Card":        "MonsterCard",
	"Cursed Monster Card":       "MonsterCard",
	"Holy/Charmed Monster Card": "MonsterCard",
	"Boss Card":                 "MonsterCard",
	"Epic Boss Card":            "MonsterCard",
	"Good Event Card":           "EventCard",
	"Bad Event Card":            "EventCard",
	"Curse Card":                "EventCard",

	"Bonus Soul Card": "BonusSoulCard",
	"Room Card":       "RoomCard",
}

// rewardKinds maps printed reward icons to engine reward kinds.
var rewardKinds = map[string]string{
	"Coin":     "engine.RewardCents",
	"Loot":     "engine.RewardLoot",
	"Treasure": "engine.RewardTreasure",
}

// stub is one card's generated data.
type stub struct {
	file string // file name inside the set folder
	name string // Go variable name
	src  []byte
}

// generate builds the stub of every card of one set plus the set list.
func generate(db *carddb.DB, code string) ([]stub, error) {
	var set carddb.Set
	for _, s := range db.Sets() {
		if s.Code == code {
			set = s
		}
	}
	if set.Code == "" {
		return nil, fmt.Errorf("cardgen: no set %q in card_db", code)
	}
	var cards []carddb.Card
	startingItems := map[string]bool{}
	for _, ref := range db.Refs() {
		c, _ := db.Get(ref)
		if c.Set != code {
			continue
		}
		cards = append(cards, c)
		if c.Type == "Character Card" {
			for _, r := range c.Related {
				startingItems[r] = true
			}
		}
	}

	var out []stub
	var errs []error
	for _, c := range cards {
		s, err := cardStub(code, c, startingItems[c.Ref().String()])
		if err != nil {
			errs = append(errs, err)
			continue
		}
		out = append(out, s)
	}
	if len(errs) > 0 {
		return nil, errors.Join(errs...)
	}
	list, err := setList(set, out)
	if err != nil {
		return nil, err
	}
	return append(out, list, docFile(set)), nil
}

func cardStub(pkg string, c carddb.Card, startingItem bool) (stub, error) {
	ref := c.Ref()
	kind, ok := kinds[c.Type]
	if !ok {
		return stub{}, fmt.Errorf("cardgen: %s: unknown card type %q", ref, c.Type)
	}
	en := c.Text[carddb.English]
	var b bytes.Buffer
	var todo []string

	fmt.Fprintf(&b, "package %s\n\nimport engine \"github.com/FourSoulsRepo/rules_engine\"\n\n", pkg)
	fmt.Fprintf(&b, "// %s (%s)\n", en.Name, c.Type)
	if text := printedText(en); len(text) > 0 {
		b.WriteString("//\n")
		for _, line := range text {
			fmt.Fprintf(&b, "//\t%s\n", line)
		}
	}

	fields := []string{
		fmt.Sprintf("Ref: %q", ref.String()),
		"Kind: engine." + kind,
		"Copies: " + strconv.Itoa(c.Copies),
	}
	for _, st := range c.Stats {
		n, err := strconv.Atoi(st.Value)
		if err != nil {
			todo = append(todo, fmt.Sprintf("%s is %q", st.Name, st.Value))
			continue
		}
		switch st.Name {
		case "HP", "ATK", "DC":
			fields = append(fields, st.Name+": "+strconv.Itoa(n))
		default:
			todo = append(todo, "stat "+st.Name)
		}
	}
	var rewards []string
	soul := 0
	for _, r := range c.Rewards {
		n, okN := rewardCount(r.Count)
		switch {
		case r.Kind == "Soul" && okN:
			soul += n // some cards print several soul icons
		case kind != "MonsterCard":
			// Icons on items and loot show what the card gives; the
			// effect blocks do that, not a reward box.
		case !okN || rewardKinds[r.Kind] == "":
			todo = append(todo, fmt.Sprintf("reward %q %s", r.Count, r.Kind))
		default:
			rewards = append(rewards, fmt.Sprintf("{Kind: %s, Amount: %d}", rewardKinds[r.Kind], n))
		}
	}
	if soul > 0 {
		fields = append(fields, "Soul: "+strconv.Itoa(soul))
	}
	if len(rewards) > 0 {
		fields = append(fields, "Rewards: []engine.Reward{"+strings.Join(rewards, ", ")+"}")
	}
	if kind == "CharacterCard" && len(c.Related) > 0 {
		fields = append(fields, fmt.Sprintf("StartingItem: %q", c.Related[0]))
	}
	if hasFootnote(en, "-Eternal-") {
		fields = append(fields, "Eternal: true")
	}
	if startingItem {
		fields = append(fields, "Outside: true")
	}
	if kind == "CharacterCard" || strings.Contains(strings.Join(en.Effects, "\n"), "{Tap Effect}") {
		fields = append(fields, "Tap: true")
	}

	b.WriteString("//\n")
	fmt.Fprintf(&b, "// %s implement the text above, then delete this line.\n", StubMarker)
	b.WriteString("// Until then `go run ./cmd/cardgen` rewrites this file from card_db.\n")
	for _, t := range todo {
		fmt.Fprintf(&b, "// Not generated: %s.\n", t)
	}
	name := varName(ref)
	fmt.Fprintf(&b, "var %s = engine.CardDef{\n%s,\n}\n", name, strings.Join(fields, ",\n"))

	src, err := format.Source(b.Bytes())
	if err != nil {
		return stub{}, fmt.Errorf("cardgen: %s: format: %w", ref, err)
	}
	return stub{file: fileName(ref), name: name, src: src}, nil
}

func setList(set carddb.Set, cards []stub) (stub, error) {
	var b bytes.Buffer
	b.WriteString("// Code generated by cardgen; DO NOT EDIT.\n\n")
	fmt.Fprintf(&b, "package %s\n\nimport engine \"github.com/FourSoulsRepo/rules_engine\"\n\n", set.Code)
	fmt.Fprintf(&b, "// Set is every card of %s.\n", set.Name)
	fmt.Fprintf(&b, "var Set = engine.CardSet{\nCode: %q,\nName: %q,\nCards: []engine.CardDef{\n", set.Code, set.Name)
	for _, c := range cards {
		b.WriteString(c.name + ",\n")
	}
	b.WriteString("},\n}\n")
	src, err := format.Source(b.Bytes())
	if err != nil {
		return stub{}, fmt.Errorf("cardgen: set list: %w", err)
	}
	return stub{file: setFile, src: src}, nil
}

// docFile is the package comment. It has no stub marker, so write
// creates it once and then leaves it to the set's authors.
func docFile(set carddb.Set) stub {
	src := fmt.Sprintf("// Package %s holds the cards of %s, one file per card.\n//\n"+
		"// Stubs come from card_db: go run ./cmd/cardgen -set %s\npackage %s\n", set.Code, set.Name, set.Code, set.Code)
	return stub{file: docName, src: []byte(src)}
}

// printedText is the card text as comment lines.
func printedText(t carddb.Text) []string {
	var lines []string
	for _, e := range t.Effects {
		lines = append(lines, strings.Split(e, "\n")...)
	}
	for _, f := range t.Footnotes {
		lines = append(lines, f.Title+" "+f.Text)
	}
	return lines
}

func hasFootnote(t carddb.Text, title string) bool {
	for _, f := range t.Footnotes {
		if f.Title == title {
			return true
		}
	}
	return false
}

// rewardCount reads printed counts like "3x" or "x10".
func rewardCount(s string) (int, bool) {
	n, err := strconv.Atoi(strings.Trim(s, "x "))
	return n, err == nil
}

// varName turns "the_d6" into "theD6" and "the_d6@2" into "theD6V2".
func varName(r carddb.Ref) string {
	parts := strings.Split(r.ID, "_")
	var b strings.Builder
	for i, p := range parts {
		if i > 0 {
			p = strings.ToUpper(p[:1]) + p[1:]
		}
		b.WriteString(p)
	}
	if r.Version > 1 {
		b.WriteString("V" + strconv.Itoa(r.Version))
	}
	name := b.String()
	if unicode.IsDigit(rune(name[0])) {
		name = "card" + name
	}
	return name
}

// fileName is "the_d6.go", or "the_d6_v2.go" for version 2.
func fileName(r carddb.Ref) string {
	if r.Version > 1 {
		return r.ID + "_v" + strconv.Itoa(r.Version) + ".go"
	}
	return r.ID + ".go"
}

// write puts the files into dir. Card files without the stub marker are
// finished and stay as they are; the set list is always rewritten.
func write(dir string, files []stub) (written, kept int, err error) {
	if err := os.MkdirAll(dir, 0o750); err != nil {
		return 0, 0, fmt.Errorf("cardgen: %w", err)
	}
	for _, f := range files {
		path := filepath.Join(dir, f.file)
		if f.file != setFile {
			card := f.file != docName
			old, err := os.ReadFile(path) //nolint:gosec // path is built from validated card IDs
			switch {
			case err == nil && !bytes.Contains(old, []byte(StubMarker)):
				if card {
					kept++
				}
				continue
			case err != nil && !errors.Is(err, os.ErrNotExist):
				return written, kept, fmt.Errorf("cardgen: %w", err)
			}
			if card {
				written++
			}
		}
		if err := os.WriteFile(path, f.src, 0o600); err != nil {
			return written, kept, fmt.Errorf("cardgen: %w", err)
		}
	}
	return written, kept, nil
}
