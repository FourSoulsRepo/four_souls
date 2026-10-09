package carddb

import (
	"errors"
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
  "set": {"code": "b2", "name": "Base Game V2"},
  "cards": [
    {"id": "the_d6", "type": "Eternal Treasure Card", "copies": 1,
     "text": {"en": {"name": "The D6", "effects": ["{Tap Effect}Reroll a die."]}},
     "artists": [{"role": "Front character artist", "name": "Someone"}],
     "image": "b2/the_d6.webp", "related": ["isaac"]},
    {"id": "isaac", "type": "Character Card", "copies": 1,
     "stats": [{"name": "HP", "value": "2"}],
     "text": {"en": {"name": "Isaac"}}},
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
		"duplicate": `{"set":{"code":"b2","name":"B"},"cards":[
			{"id":"x","type":"T","copies":1,"text":{"en":{"name":"X"}}},
			{"id":"x","type":"T","copies":1,"text":{"en":{"name":"X"}}}]}`,
		"version one": `{"set":{"code":"b2","name":"B"},"cards":[
			{"id":"x","version":1,"type":"T","copies":1,"text":{"en":{"name":"X"}}}]}`,
		"no name": `{"set":{"code":"b2","name":"B"},"cards":[
			{"id":"x","type":"T","copies":1,"text":{}}]}`,
		"no copies": `{"set":{"code":"b2","name":"B"},"cards":[
			{"id":"x","type":"T","text":{"en":{"name":"X"}}}]}`,
		"escaping image": `{"set":{"code":"b2","name":"B"},"cards":[
			{"id":"x","type":"T","copies":1,"text":{"en":{"name":"X"}},"image":"../secret.png"}]}`,
		"bad related": `{"set":{"code":"b2","name":"B"},"cards":[
			{"id":"x","type":"T","copies":1,"text":{"en":{"name":"X"}},"related":["Not A Ref"]}]}`,
		"wrong file name": `{"set":{"code":"g2","name":"B"},"cards":[]}`,
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
		"data/b2.json": {Data: []byte(`{"set":{"code":"b2","name":"B"},"cards":[` + card + `]}`)},
		"data/g2.json": {Data: []byte(`{"set":{"code":"g2","name":"G"},"cards":[` + card + `]}`)},
	}
	_, err := Load(fsys)
	if err == nil || !strings.Contains(err.Error(), "already defined in set b2") {
		t.Errorf("want duplicate error, got %v", err)
	}
}

func TestLoadEmpty(t *testing.T) {
	db, err := Load(fstest.MapFS{})
	if err != nil || len(db.Refs()) != 0 {
		t.Errorf("empty: %v, %v", db, err)
	}
}
