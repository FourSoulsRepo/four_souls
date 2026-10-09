package carddb

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"path"
	"sort"
	"strings"
)

// Folder names inside the card data file system.
const (
	DataDir   = "data"   // one JSON file per set
	ImagesDir = "images" // card images, git-ignored (A-11)
)

// setFile is the JSON layout of data/<code>.json.
type setFile struct {
	Set   Set    `json:"set"`
	Cards []Card `json:"cards"`
}

// DB holds every loaded card.
type DB struct {
	sets  []Set
	cards map[Ref]Card
	refs  []Ref
}

// Load reads every data/*.json file of fsys and validates it.
func Load(fsys fs.FS) (*DB, error) {
	names, err := fs.Glob(fsys, path.Join(DataDir, "*.json"))
	if err != nil {
		return nil, fmt.Errorf("carddb: list sets: %w", err)
	}
	sort.Strings(names)
	db := &DB{cards: map[Ref]Card{}}
	var problems []error
	for _, name := range names {
		data, err := fs.ReadFile(fsys, name)
		if err != nil {
			return nil, fmt.Errorf("carddb: read %s: %w", name, err)
		}
		var f setFile
		if err := json.Unmarshal(data, &f); err != nil {
			return nil, fmt.Errorf("carddb: decode %s: %w", name, err)
		}
		problems = append(problems, db.add(name, f)...)
	}
	if len(problems) > 0 {
		return nil, errors.Join(problems...)
	}
	sort.Slice(db.refs, func(i, j int) bool { return db.refs[i].String() < db.refs[j].String() })
	return db, nil
}

func (db *DB) add(file string, f setFile) []error {
	var errs []error
	if f.Set.Code == "" || f.Set.Name == "" {
		errs = append(errs, fmt.Errorf("%s: set code and name are required", file))
	}
	if want := path.Join(DataDir, f.Set.Code+".json"); file != want {
		errs = append(errs, fmt.Errorf("%s: file name must be %s", file, want))
	}
	db.sets = append(db.sets, f.Set)
	for _, c := range f.Cards {
		c.Set = f.Set.Code
		if err := validate(c); err != nil {
			errs = append(errs, fmt.Errorf("%s: card %q: %w", file, c.ID, err))
			continue
		}
		ref := c.Ref()
		if old, ok := db.cards[ref]; ok {
			errs = append(errs, fmt.Errorf("%s: card %s already defined in set %s", file, ref, old.Set))
			continue
		}
		db.cards[ref] = c
		db.refs = append(db.refs, ref)
	}
	return errs
}

func validate(c Card) error {
	var errs []error
	if !idPattern.MatchString(c.ID) {
		errs = append(errs, fmt.Errorf("bad id %q", c.ID))
	}
	if c.Version == 1 || c.Version < 0 {
		errs = append(errs, errors.New("version must be omitted for the first version, or be 2 or more"))
	}
	if c.Type == "" {
		errs = append(errs, errors.New("type is required"))
	}
	if c.Copies < 1 {
		errs = append(errs, errors.New("copies must be at least 1"))
	}
	if c.Text[English].Name == "" {
		errs = append(errs, errors.New("an English name is required"))
	}
	if c.Image != "" && !safePath(c.Image) {
		errs = append(errs, fmt.Errorf("bad image path %q", c.Image))
	}
	for _, r := range c.Related {
		if _, err := ParseRef(r); err != nil {
			errs = append(errs, fmt.Errorf("related: %w", err))
		}
	}
	return errors.Join(errs...)
}

// safePath accepts a clean relative path that stays inside the images folder.
func safePath(p string) bool {
	return fs.ValidPath(p) && !strings.Contains(p, "\\")
}

// Sets returns the loaded sets in file order.
func (db *DB) Sets() []Set {
	return append([]Set(nil), db.sets...)
}

// Get returns one card version.
func (db *DB) Get(ref Ref) (Card, bool) {
	c, ok := db.cards[ref]
	return c, ok
}

// Refs returns every card reference, sorted.
func (db *DB) Refs() []Ref {
	return append([]Ref(nil), db.refs...)
}

// Has reports whether every reference is known (CD-03).
// It returns the unknown ones.
func (db *DB) Has(refs []Ref) []Ref {
	var missing []Ref
	for _, r := range refs {
		if _, ok := db.cards[r]; !ok {
			missing = append(missing, r)
		}
	}
	return missing
}

// ImagePath returns where a card's image lives inside the file system,
// or "" when the card has none.
func ImagePath(c Card) string {
	if c.Image == "" {
		return ""
	}
	return path.Join(ImagesDir, c.Image)
}
