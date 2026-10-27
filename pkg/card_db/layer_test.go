package carddb

import (
	"io/fs"
	"reflect"
	"testing"
	"testing/fstest"
)

// A translator's own images lie over the built-in ones (16.12).
func TestWithImages(t *testing.T) {
	d6 := Card{Images: []string{"b2/the_d6.webp"}}
	isaac := Card{Images: []string{"b2/isaac.webp"}}
	cain := Card{Images: []string{"b2/cain.webp"}}
	base := fstest.MapFS{
		"data/b2.json":             {Data: []byte(sampleSet)},
		"images/en/b2/the_d6.webp": {Data: []byte("en")},
		"images/en/b2/isaac.webp":  {Data: []byte("en")},
		"images/de/b2/isaac.webp":  {Data: []byte("built-in de")},
		"images/assets/die.webp":   {Data: []byte("die")},
	}
	extra := fstest.MapFS{
		"de/b2/the_d6.webp": {Data: []byte("own de")},
		"de/b2/isaac.webp":  {Data: []byte("own de")},
		"pl/b2/cain.webp":   {Data: []byte("own pl")},
		"assets/die.webp":   {Data: []byte("not a language")},
		"data/b2.json":      {Data: []byte("not a language")},
		"README.md":         {Data: []byte("a file")},
	}
	fsys := WithImages(base, extra)
	cases := []struct {
		name string
		card Card
		lang Lang
		want string
	}{
		{"own image", d6, "de", "de/b2/the_d6.webp"},
		{"own image wins over a built-in one", isaac, "de", "de/b2/isaac.webp"},
		{"own language with no built-in folder", cain, "pl", "pl/b2/cain.webp"},
		{"no own image: English", d6, "pl", "en/b2/the_d6.webp"},
		{"no image anywhere", cain, "de", ""},
		{"English stays", d6, English, "en/b2/the_d6.webp"},
	}
	for _, tc := range cases {
		if got := FindImage(fsys, tc.card, tc.lang); got != tc.want {
			t.Errorf("%s: %q, want %q", tc.name, got, tc.want)
		}
	}
	read := func(name string) string {
		data, err := fs.ReadFile(fsys, name)
		if err != nil {
			return err.Error()
		}
		return string(data)
	}
	for name, want := range map[string]string{
		"images/de/b2/isaac.webp": "own de",
		"images/en/b2/isaac.webp": "en",
		"images/assets/die.webp":  "die",
		"data/b2.json":            sampleSet,
		"images/pl/b2/cain.webp":  "own pl",
	} {
		if got := read(name); got != want {
			t.Errorf("read %s: %q, want %q", name, got, want)
		}
	}
	// Paths that could leave the images folder never reach the own images.
	for _, name := range []string{"images/../data/b2.json", "images/de/../pl/b2/cain.webp", "images/de\\b2\\the_d6.webp", "images//de/b2/the_d6.webp", "/images/de/b2/the_d6.webp"} {
		if sub, ok := extraPath(name); ok {
			t.Errorf("extraPath(%q) = %q, want refused", name, sub)
		}
		if data, err := fs.ReadFile(fsys, name); err == nil {
			t.Errorf("read %s: %q, want an error", name, data)
		}
	}
	if got := ImageLangs(fsys); !reflect.DeepEqual(got, []Lang{"de", English, "pl"}) {
		t.Errorf("ImageLangs = %v, want [de en pl]", got)
	}
	if info, err := fs.Stat(fsys, "images/pl/b2"); err != nil || !info.IsDir() {
		t.Errorf("a folder only the own images have: %v, %v", info, err)
	}
	if _, err := Load(fsys); err != nil {
		t.Errorf("card data through the layer: %v", err)
	}
}

// Without built-in images, as in a public clone, own images still show
// and other cards have none (A-11).
func TestWithImagesNoBuiltIn(t *testing.T) {
	d6 := Card{Images: []string{"b2/the_d6.webp"}}
	fsys := WithImages(fstest.MapFS{"data/b2.json": {Data: []byte(sampleSet)}},
		fstest.MapFS{"de/b2/the_d6.webp": {Data: []byte("own de")}})
	if got := FindImage(fsys, d6, "de"); got != "de/b2/the_d6.webp" {
		t.Errorf("own image: %q", got)
	}
	if got := FindImage(fsys, d6, English); got != "" {
		t.Errorf("no English image: %q", got)
	}
	if got := ImageLangs(fsys); !reflect.DeepEqual(got, []Lang{"de"}) {
		t.Errorf("ImageLangs = %v, want [de]", got)
	}
	empty := WithImages(fstest.MapFS{}, fstest.MapFS{})
	if got := ImageLangs(empty); len(got) != 0 {
		t.Errorf("ImageLangs of nothing = %v", got)
	}
}
