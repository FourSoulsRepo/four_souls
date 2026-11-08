package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"slices"
	"testing"
	"testing/fstest"
)

// soundLists gives each sound's list by its JSON name.
func soundLists(t *testing.T, s Sounds) map[string][]string {
	t.Helper()
	v := reflect.ValueOf(s)
	lists := make(map[string][]string, v.NumField())
	for i := range v.NumField() {
		list, ok := v.Field(i).Interface().([]string)
		if !ok {
			t.Fatalf("Sounds.%s is not a list", v.Type().Field(i).Name)
		}
		lists[v.Type().Field(i).Tag.Get("json")] = list
	}
	return lists
}

// Every sound has its files; the new ones of 20.17 keep their variants.
func TestSoundFiles(t *testing.T) {
	for name, files := range soundLists(t, soundFiles) {
		if len(files) == 0 {
			t.Fatalf("%s has no file", name)
		}
	}
	want := []string{"assets/sounds/loot-land-1.mp3", "assets/sounds/loot-land-2.mp3", "assets/sounds/loot-land-3.mp3"}
	if !slices.Equal(soundFiles.LootLand, want) {
		t.Fatalf("loot-land: %v", soundFiles.LootLand)
	}
	if !slices.Equal(soundFiles.TurnStart, []string{"assets/sounds/turn-start.mp3"}) {
		t.Fatalf("turn-start: %v", soundFiles.TurnStart)
	}
}

func TestSoundsFromFiles(t *testing.T) {
	files := fstest.MapFS{}
	for _, list := range soundLists(t, soundFiles) {
		for _, f := range list {
			files["images/"+f] = &fstest.MapFile{Data: []byte("ID3")}
		}
	}
	got := soundLists(t, (&App{cardFiles: files}).Sounds())
	for name, list := range soundLists(t, soundFiles) {
		want := make([]string, len(list))
		for i, f := range list {
			want[i] = "/cardimg/" + f
		}
		if !slices.Equal(got[name], want) {
			t.Fatalf("%s: %v, want %v", name, got[name], want)
		}
	}
	// Without the private images the game is silent: every list is empty.
	for name, list := range soundLists(t, (&App{cardFiles: os.DirFS(t.TempDir())}).Sounds()) {
		if list == nil || len(list) != 0 {
			t.Fatalf("%s: %#v without images", name, list)
		}
	}
	// A missing variant is left out; the rest still play.
	delete(files, "images/assets/sounds/loot-land-2.mp3")
	if got := (&App{cardFiles: files}).Sounds().LootLand; !slices.Equal(got, []string{"/cardimg/assets/sounds/loot-land-1.mp3", "/cardimg/assets/sounds/loot-land-3.mp3"}) {
		t.Fatalf("loot-land without its second file: %v", got)
	}
}

// The sounds are MP3 in the private images submodule. With the submodule
// every file of every sound is there, found and served as audio; without
// it, nothing is.
func TestSoundsServed(t *testing.T) {
	files := os.DirFS("pkg/card_db")
	_, err := os.Stat("pkg/card_db/images/" + soundFiles.AttackHit[0])
	withImages := err == nil
	got := soundLists(t, (&App{cardFiles: files}).Sounds())
	for name, list := range soundLists(t, soundFiles) {
		var want []string
		for _, f := range list {
			if _, err := os.Stat("pkg/card_db/images/" + f); err == nil {
				want = append(want, "/cardimg/"+f)
			} else if withImages {
				t.Fatalf("%s: %s is missing from the images", name, f)
			}
		}
		if len(want) != len(got[name]) || (len(want) != 0 && !slices.Equal(want, got[name])) {
			t.Fatalf("%s: %v does not match the files %v", name, got[name], want)
		}
		for _, url := range got[name] {
			rec := httptest.NewRecorder()
			h := cardImages(files)(http.NotFoundHandler())
			h.ServeHTTP(rec, httptest.NewRequestWithContext(t.Context(), http.MethodGet, "http://wails.localhost"+url, nil))
			if rec.Code != http.StatusOK || rec.Header().Get("Content-Type") != "audio/mpeg" || rec.Body.Len() == 0 {
				t.Fatalf("%s: %d %q, not an MP3", url, rec.Code, rec.Header().Get("Content-Type"))
			}
		}
	}
}
