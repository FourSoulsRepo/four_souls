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
	name  string // path inside the file system
}

// DB holds every loaded card.
type DB struct {
	sets  []Set
	cards map[Ref]Card // back faces too (ADR 023)
	refs  []Ref        // fronts only: a back is never dealt
}

// Load reads every data/*.json file of fsys and validates it.
// Sets load oldest first; a reprint must follow its card's set (ADR 010).
func Load(fsys fs.FS) (*DB, error) {
	names, err := fs.Glob(fsys, path.Join(DataDir, "*.json"))
	if err != nil {
		return nil, fmt.Errorf("carddb: list sets: %w", err)
	}
	sort.Strings(names)
	var files []setFile
	var problems, broken []error
	for _, name := range names {
		data, err := fs.ReadFile(fsys, name)
		if err != nil {
			return nil, fmt.Errorf("carddb: read %s: %w", name, err)
		}
		var f setFile
		if err := json.Unmarshal(data, &f); err != nil {
			broken = append(broken, jsonError(name, data, err))
			continue
		}
		problems = append(problems, checkSet(name, f)...)
		f.name = name
		files = append(files, f)
	}
	// A broken file hides its cards, so other checks would mislead.
	if len(broken) > 0 {
		return nil, fmt.Errorf("carddb: decode: %w", errors.Join(broken...))
	}
	sort.SliceStable(files, func(i, j int) bool { return files[i].Set.Order < files[j].Set.Order })
	db := &DB{cards: map[Ref]Card{}}
	for i, f := range files {
		if i > 0 && f.Set.Order == files[i-1].Set.Order {
			problems = append(problems, fmt.Errorf("%s: order %d is also used by set %s", f.name, f.Set.Order, files[i-1].Set.Code))
		}
		db.sets = append(db.sets, f.Set)
		problems = append(problems, db.add(f)...)
	}
	for _, f := range files {
		problems = append(problems, db.addReprints(f)...)
	}
	if len(problems) > 0 {
		return nil, errors.Join(problems...)
	}
	sort.Slice(db.refs, func(i, j int) bool { return db.refs[i].String() < db.refs[j].String() })
	return db, nil
}

func checkSet(file string, f setFile) []error {
	var errs []error
	if f.Set.Code == "" || f.Set.Name == "" {
		errs = append(errs, fmt.Errorf("%s: set code and name are required", file))
	}
	if want := path.Join(DataDir, f.Set.Code+".json"); file != want {
		errs = append(errs, fmt.Errorf("%s: file name must be %s", file, want))
	}
	if f.Set.Order < 1 {
		errs = append(errs, fmt.Errorf("%s: set order must be 1 or more", file))
	}
	return errs
}

// add loads the cards a set defines; reprints come later.
func (db *DB) add(f setFile) []error {
	var errs []error
	for _, c := range f.Cards {
		if c.Reprint {
			continue
		}
		c.Set = f.Set.Code
		if err := validate(c); err != nil {
			errs = append(errs, fmt.Errorf("%s: card %q: %w", f.name, c.ID, err))
			continue
		}
		ref := c.Ref()
		if old, ok := db.cards[ref]; ok {
			errs = append(errs, fmt.Errorf("%s: card %s already defined in set %s", f.name, ref, old.Set))
			continue
		}
		c.Prints = []Print{{Set: c.Set, Copies: c.Copies, Images: c.Images}}
		db.cards[ref] = c
		db.refs = append(db.refs, ref)
		if b, ok := c.BackCard(); ok {
			if old, ok := db.cards[b.Ref()]; ok {
				errs = append(errs, fmt.Errorf("%s: back %s of %s already defined in set %s", f.name, b.ID, ref, old.Set))
				continue
			}
			db.cards[b.Ref()] = b
		}
	}
	return errs
}

// addReprints adds a set's copies to cards of older sets.
func (db *DB) addReprints(f setFile) []error {
	var errs []error
	for _, r := range f.Cards {
		if !r.Reprint {
			continue
		}
		ref := r.Ref()
		if err := validateReprint(r); err != nil {
			errs = append(errs, fmt.Errorf("%s: reprint %s: %w", f.name, ref, err))
			continue
		}
		c, ok := db.cards[ref]
		if !ok || db.order(c.Set) >= f.Set.Order {
			errs = append(errs, fmt.Errorf("%s: reprint %s: no such card in an older set", f.name, ref))
			continue
		}
		if c.Prints[len(c.Prints)-1].Set == f.Set.Code {
			errs = append(errs, fmt.Errorf("%s: reprint %s listed twice", f.name, ref))
			continue
		}
		c.Prints = append(c.Prints, Print{Set: f.Set.Code, Copies: r.Copies, Images: r.Images})
		db.cards[ref] = c
	}
	return errs
}

// order returns the order of a loaded set.
func (db *DB) order(code string) int {
	for _, s := range db.sets {
		if s.Code == code {
			return s.Order
		}
	}
	return 0
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
	if c.Text[English].Name == "" {
		errs = append(errs, errors.New("an English name is required"))
	}
	errs = append(errs, validateLangs(c)...)
	if err := validateCopies(c.Copies, c.Images); err != nil {
		errs = append(errs, err)
	}
	for _, r := range c.Related {
		if _, err := ParseRef(r); err != nil {
			errs = append(errs, fmt.Errorf("related: %w", err))
		}
	}
	if c.Back != nil {
		if err := validateBack(c); err != nil {
			errs = append(errs, fmt.Errorf("back: %w", err))
		}
	}
	return errors.Join(errs...)
}

// validateBack checks a two-sided card's back face (ADR 023): its own
// id and English name, one image at most, and a name in every language.
func validateBack(c Card) error {
	b := c.Back
	var errs []error
	if !idPattern.MatchString(b.ID) || b.ID == c.ID {
		errs = append(errs, fmt.Errorf("bad id %q", b.ID))
	}
	if b.Text[English].Name == "" {
		errs = append(errs, errors.New("an English name is required"))
	}
	for lang, t := range b.Text {
		if t.Name == "" {
			errs = append(errs, fmt.Errorf("text %q has no name", lang))
		}
	}
	if len(b.Images) > 1 {
		errs = append(errs, errors.New("one image at most"))
	}
	for _, img := range b.Images {
		if !safePath(img) {
			errs = append(errs, fmt.Errorf("bad image path %q", img))
		}
	}
	return errors.Join(errs...)
}

// validateLangs checks the translations (LO-02, CD-06): every language
// has a name, and translators belong to a language with text.
func validateLangs(c Card) []error {
	var errs []error
	for lang, t := range c.Text {
		if t.Name == "" {
			errs = append(errs, fmt.Errorf("text %q has no name", lang))
		}
	}
	for lang, names := range c.Translators {
		if _, ok := c.Text[lang]; !ok {
			errs = append(errs, fmt.Errorf("translators for %q, which has no text", lang))
		}
		for _, n := range names {
			if strings.TrimSpace(n) == "" {
				errs = append(errs, fmt.Errorf("a blank translator for %q", lang))
			}
		}
	}
	return errs
}

// validateReprint accepts only an id, copies and images: the text lives
// with the card in its own set.
func validateReprint(r Card) error {
	var errs []error
	if !idPattern.MatchString(r.ID) {
		errs = append(errs, fmt.Errorf("bad id %q", r.ID))
	}
	if r.Type != "" || len(r.Text) > 0 || len(r.Stats) > 0 || len(r.Rewards) > 0 ||
		len(r.Artists) > 0 || len(r.Translators) > 0 || len(r.Related) > 0 || r.Rebalanced || r.ThreePlus || r.Back != nil {
		errs = append(errs, errors.New("a reprint has only id, version, copies and images"))
	}
	if err := validateCopies(r.Copies, r.Images); err != nil {
		errs = append(errs, err)
	}
	return errors.Join(errs...)
}

// validateCopies checks the copy count against the images: none, one
// for every copy, or one per copy.
func validateCopies(copies int, images []string) error {
	var errs []error
	if copies < 1 {
		errs = append(errs, errors.New("copies must be at least 1"))
	}
	if len(images) > 1 && len(images) != copies {
		errs = append(errs, fmt.Errorf("%d images for %d copies; give one, or one per copy", len(images), copies))
	}
	for _, img := range images {
		if !safePath(img) {
			errs = append(errs, fmt.Errorf("bad image path %q", img))
		}
	}
	return errors.Join(errs...)
}

// safePath accepts a clean relative path that stays inside the images folder.
func safePath(p string) bool {
	return fs.ValidPath(p) && !strings.Contains(p, "\\")
}

// Sets returns the loaded sets, oldest first.
func (db *DB) Sets() []Set {
	return append([]Set(nil), db.sets...)
}

// Get returns one card version.
func (db *DB) Get(ref Ref) (Card, bool) {
	c, ok := db.cards[ref]
	return c, ok
}

// Refs returns every card reference, sorted. Back faces are left out;
// Get finds them, and the front's Back names them (ADR 023).
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

// HasImage reports whether the card's English image file exists in fsys.
// Clients show a blank card with text when it does not (A-11, SP-05).
func HasImage(fsys fs.FS, c Card) bool {
	return isFile(fsys, ImagePath(c, English))
}

// FindImage returns the card's first image in lang as a path inside
// images/, e.g. "ru/b2/the_d6.webp". It falls back to English (LO-06)
// and returns "" when neither file exists (ADR 021).
func FindImage(fsys fs.FS, c Card, lang Lang) string {
	for _, l := range []Lang{lang, English} {
		if p := ImagePath(c, l); isFile(fsys, p) {
			return strings.TrimPrefix(p, ImagesDir+"/")
		}
	}
	return ""
}

// ImagePath returns where a card's first image in lang lives inside the
// file system, images/<lang>/<set>/<file>, or "" when the card has none.
// Card data names images without the language (ADR 021).
func ImagePath(c Card, lang Lang) string {
	if len(c.Images) == 0 || !isLang(string(lang)) {
		return ""
	}
	return path.Join(ImagesDir, string(lang), c.Images[0])
}

// ImageLangs lists the languages that have an image folder, sorted.
// Other folders, such as images/assets, are not languages.
func ImageLangs(fsys fs.FS) []Lang {
	entries, err := fs.ReadDir(fsys, ImagesDir)
	if err != nil {
		return nil
	}
	var langs []Lang
	for _, e := range entries {
		if e.IsDir() && isLang(e.Name()) {
			langs = append(langs, Lang(e.Name()))
		}
	}
	return langs
}

// isLang reports whether s looks like a language code: 2–3 small letters.
func isLang(s string) bool {
	if len(s) < 2 || len(s) > 3 {
		return false
	}
	for _, r := range s {
		if r < 'a' || r > 'z' {
			return false
		}
	}
	return true
}

func isFile(fsys fs.FS, p string) bool {
	if p == "" {
		return false
	}
	info, err := fs.Stat(fsys, p)
	return err == nil && !info.IsDir()
}
