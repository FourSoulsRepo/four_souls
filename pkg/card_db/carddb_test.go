package carddb

import (
	"errors"
	"reflect"
	"strings"
	"testing"
	"testing/fstest"
)

func TestParseRef(t *testing.T) {
	good := map[string]Ref{
		"the_d6":   {ID: "the_d6", Version: 1},
		"the_d6@2": {ID: "the_d6", Version: 2},
		"a_penny":  {ID: "a_penny", Version: 1},
	}
	for in, want := range good {
		got, err := ParseRef(in)
		if err != nil || got != want {
			t.Errorf("ParseRef(%q) = %v, %v; want %v", in, got, err, want)
		}
		if got.String() != in {
			t.Errorf("String() = %q, want %q", got.String(), in)
		}
	}
	for _, in := range []string{"", "The_D6", "the d6", "x@1", "x@0", "x@02", "x@", "x@2@3", "../x", "x__y"} {
		if _, err := ParseRef(in); !errors.Is(err, ErrBadRef) {
			t.Errorf("ParseRef(%q): want ErrBadRef, got %v", in, err)
		}
	}
}

const sampleSet = `{
  "set": {"code": "b2", "name": "Base Game V2", "order": 1, "released": "2018-12"},
  "cards": [
    {"id": "the_d6", "type": "Eternal Treasure Card", "copies": 1,
     "text": {"en": {"name": "The D6", "effects": ["{Tap Effect}Reroll a die."]}},
     "artists": [{"role": "Front character artist", "name": "Someone"}],
     "images": ["b2/the_d6.webp"], "related": ["isaac"]},
    {"id": "isaac", "type": "Character Card", "copies": 1,
     "stats": [{"name": "HP", "value": "2"}],
     "text": {"en": {"name": "Isaac"}, "ru": {"name": "Исаак"}},
     "translators": {"ru": ["Someone"]}},
    {"id": "the_d6", "version": 2, "type": "Eternal Treasure Card", "copies": 1,
     "text": {"en": {"name": "The D6"}}}
  ]
}`

func TestLoad(t *testing.T) {
	db, err := Load(fstest.MapFS{"data/b2.json": {Data: []byte(sampleSet)}})
	if err != nil {
		t.Fatal(err)
	}
	if got := len(db.Refs()); got != 3 {
		t.Fatalf("refs = %d, want 3", got)
	}
	d6, ok := db.Get(Ref{ID: "the_d6", Version: 1})
	if !ok || d6.Set != "b2" || d6.Name(English) != "The D6" || d6.Name("ru") != "The D6" {
		t.Errorf("the_d6 = %+v, %v", d6, ok)
	}
	if got := ImagePath(d6); got != "images/b2/the_d6.webp" {
		t.Errorf("image path = %q", got)
	}
	if isaac, _ := db.Get(Ref{ID: "isaac", Version: 1}); isaac.Name("ru") != "Исаак" || isaac.Name(English) != "Isaac" {
		t.Errorf("isaac names: %q, %q", isaac.Name("ru"), isaac.Name(English))
	}
	if _, ok := db.Get(Ref{ID: "the_d6", Version: 2}); !ok {
		t.Error("the_d6@2 missing")
	}
	missing := db.Has([]Ref{{ID: "isaac", Version: 1}, {ID: "maggy", Version: 1}})
	if len(missing) != 1 || missing[0].ID != "maggy" {
		t.Errorf("Has: missing = %v", missing)
	}
}

func TestLoadRejectsBadData(t *testing.T) {
	cases := map[string]string{
		"duplicate": `{"set":{"code":"b2","name":"B","order":1},"cards":[
			{"id":"x","type":"T","copies":1,"text":{"en":{"name":"X"}}},
			{"id":"x","type":"T","copies":1,"text":{"en":{"name":"X"}}}]}`,
		"version one": `{"set":{"code":"b2","name":"B","order":1},"cards":[
			{"id":"x","version":1,"type":"T","copies":1,"text":{"en":{"name":"X"}}}]}`,
		"no name": `{"set":{"code":"b2","name":"B","order":1},"cards":[
			{"id":"x","type":"T","copies":1,"text":{}}]}`,
		"no copies": `{"set":{"code":"b2","name":"B","order":1},"cards":[
			{"id":"x","type":"T","text":{"en":{"name":"X"}}}]}`,
		"escaping image": `{"set":{"code":"b2","name":"B","order":1},"cards":[
			{"id":"x","type":"T","copies":1,"text":{"en":{"name":"X"}},"images":["../secret.png"]}]}`,
		"images and copies": `{"set":{"code":"b2","name":"B","order":1},"cards":[
			{"id":"x","type":"T","copies":3,"text":{"en":{"name":"X"}},"images":["b2/x.webp","b2/x-2.webp"]}]}`,
		"no order": `{"set":{"code":"b2","name":"B"},"cards":[]}`,
		"bad related": `{"set":{"code":"b2","name":"B","order":1},"cards":[
			{"id":"x","type":"T","copies":1,"text":{"en":{"name":"X"}},"related":["Not A Ref"]}]}`,
		"wrong file name": `{"set":{"code":"g2","name":"B","order":1},"cards":[]}`,
		"translation without a name": `{"set":{"code":"b2","name":"B","order":1},"cards":[
			{"id":"x","type":"T","copies":1,"text":{"en":{"name":"X"},"ru":{"name":""}}}]}`,
		"translators without text": `{"set":{"code":"b2","name":"B","order":1},"cards":[
			{"id":"x","type":"T","copies":1,"text":{"en":{"name":"X"}},"translators":{"ru":["A"]}}]}`,
		"blank translator": `{"set":{"code":"b2","name":"B","order":1},"cards":[
			{"id":"x","type":"T","copies":1,"text":{"en":{"name":"X"},"ru":{"name":"Икс"}},"translators":{"ru":[" "]}}]}`,
	}
	for name, data := range cases {
		if _, err := Load(fstest.MapFS{"data/b2.json": {Data: []byte(data)}}); err == nil {
			t.Errorf("%s: want error", name)
		}
	}
}

func TestLoadDuplicateAcrossSets(t *testing.T) {
	card := `{"id":"x","type":"T","copies":1,"text":{"en":{"name":"X"}}}`
	fsys := fstest.MapFS{
		"data/b2.json": {Data: []byte(`{"set":{"code":"b2","name":"B","order":1},"cards":[` + card + `]}`)},
		"data/g2.json": {Data: []byte(`{"set":{"code":"g2","name":"G","order":2},"cards":[` + card + `]}`)},
	}
	_, err := Load(fsys)
	if err == nil || !strings.Contains(err.Error(), "already defined in set b2") {
		t.Errorf("want duplicate error, got %v", err)
	}
}

func TestHasImage(t *testing.T) {
	fsys := fstest.MapFS{
		"data/b2.json":           {Data: []byte(sampleSet)},
		"images/b2/the_d6.webp":  {Data: []byte("webp")},
		"images/b2/isaac.webp/x": {Data: []byte("dir, not a file")},
	}
	db, err := Load(fsys)
	if err != nil {
		t.Fatal(err)
	}
	d6, _ := db.Get(Ref{ID: "the_d6", Version: 1})
	if !HasImage(fsys, d6) {
		t.Error("the_d6 image should exist")
	}
	isaac, _ := db.Get(Ref{ID: "isaac", Version: 1})
	if HasImage(fsys, isaac) {
		t.Error("isaac has no image field, so no image")
	}
}

func TestLoadEmpty(t *testing.T) {
	db, err := Load(fstest.MapFS{})
	if err != nil || len(db.Refs()) != 0 {
		t.Errorf("empty: %v, %v", db, err)
	}
}

// reprintSets has a card in b2 and its reprints in two newer sets;
// fsp2 sorts before g2 by name but comes after it by order.
func reprintSets(g2Reprint string) fstest.MapFS {
	return fstest.MapFS{
		"data/b2.json": {Data: []byte(`{"set":{"code":"b2","name":"B","order":1},"cards":[
			{"id":"two_cents","type":"2 Cent Card","copies":2,"text":{"en":{"name":"2 Cents!"}},
			 "images":["b2/two_cents.webp","b2/two_cents-2.webp"]}]}`)},
		"data/g2.json": {Data: []byte(`{"set":{"code":"g2","name":"G","order":2},"cards":[` + g2Reprint + `]}`)},
		"data/fsp2.json": {Data: []byte(`{"set":{"code":"fsp2","name":"F","order":3},"cards":[
			{"id":"two_cents","reprint":true,"copies":1,"images":["fsp2/two_cents.webp"]}]}`)},
	}
}

func TestLoadReprints(t *testing.T) {
	db, err := Load(reprintSets(`{"id":"two_cents","reprint":true,"copies":2,"images":["g2/two_cents.webp","g2/two_cents-2.webp"]}`))
	if err != nil {
		t.Fatal(err)
	}
	var codes []string
	for _, s := range db.Sets() {
		codes = append(codes, s.Code)
	}
	if strings.Join(codes, ",") != "b2,g2,fsp2" {
		t.Errorf("sets = %v, want oldest first", codes)
	}
	if got := len(db.Refs()); got != 1 {
		t.Fatalf("refs = %d; a reprint is not a new card", got)
	}
	c, _ := db.Get(Ref{ID: "two_cents", Version: 1})
	want := []Print{
		{Set: "b2", Copies: 2, Images: []string{"b2/two_cents.webp", "b2/two_cents-2.webp"}},
		{Set: "g2", Copies: 2, Images: []string{"g2/two_cents.webp", "g2/two_cents-2.webp"}},
		{Set: "fsp2", Copies: 1, Images: []string{"fsp2/two_cents.webp"}},
	}
	if !reflect.DeepEqual(c.Prints, want) {
		t.Errorf("prints = %+v\nwant %+v", c.Prints, want)
	}
	if c.Set != "b2" || c.Copies != 2 {
		t.Errorf("card = %+v; its own set and copies stay", c)
	}
}

func TestLoadRejectsBadReprints(t *testing.T) {
	cases := map[string]string{
		"unknown card": `{"id":"three_cents","reprint":true,"copies":1}`,
		"with text":    `{"id":"two_cents","reprint":true,"copies":1,"type":"T","text":{"en":{"name":"X"}}}`,
		"no copies":    `{"id":"two_cents","reprint":true}`,
		"listed twice": `{"id":"two_cents","reprint":true,"copies":1},{"id":"two_cents","reprint":true,"copies":1}`,
	}
	for name, g2 := range cases {
		if _, err := Load(reprintSets(g2)); err == nil {
			t.Errorf("%s: want error", name)
		}
	}
	// A reprint must come after the set that defines the card.
	fsys := reprintSets(`{"id":"x","type":"T","copies":1,"text":{"en":{"name":"X"}}}`)
	fsys["data/b2.json"] = &fstest.MapFile{Data: []byte(`{"set":{"code":"b2","name":"B","order":1},"cards":[
		{"id":"x","reprint":true,"copies":1}]}`)}
	if _, err := Load(fsys); err == nil || !strings.Contains(err.Error(), "older set") {
		t.Errorf("reprint before its card: got %v", err)
	}
}

func TestLoadRejectsSharedOrder(t *testing.T) {
	fsys := fstest.MapFS{
		"data/b2.json": {Data: []byte(`{"set":{"code":"b2","name":"B","order":1},"cards":[]}`)},
		"data/g2.json": {Data: []byte(`{"set":{"code":"g2","name":"G","order":1},"cards":[]}`)},
	}
	if _, err := Load(fsys); err == nil {
		t.Error("want error for two sets with one order")
	}
}
