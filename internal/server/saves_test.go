package server

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/FourSoulsRepo/four_souls/internal/protocol"
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards"
)

// writeSave puts a saved game of a fresh 2-player game into dir.
func writeSave(t *testing.T, dir, id, engineVersion string) {
	t.Helper()
	g, _, err := engine.NewGame(engine.Setup{Seed: 5, Players: 2, Sets: cards.Sets(), BonusSouls: true})
	if err != nil {
		t.Fatal(err)
	}
	state, err := g.Save()
	if err != nil {
		t.Fatal(err)
	}
	f := saveFile{
		Format: saveFormat, Engine: engineVersion, Game: id, Sets: []string{"b2"},
		Seats: []savedSeat{{Name: "Ann", Token: newToken()}, {Name: "Bob", Token: newToken()}}, State: state,
	}
	data, err := json.Marshal(f)
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(dir, 0o750); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, id+saveExt), data, 0o600); err != nil {
		t.Fatal(err)
	}
}

func startHubWithSaves(t *testing.T, saves string) *Hub {
	t.Helper()
	h := NewHub()
	h.saves = saves
	ctx, cancel := context.WithCancel(context.Background())
	go h.Run(ctx)
	t.Cleanup(cancel)
	return h
}

func TestLoadChecks(t *testing.T) {
	dir := t.TempDir()
	writeSave(t, dir, "g1", engine.Version)
	writeSave(t, dir, "g2", "0.0.0-older")
	h := startHubWithSaves(t, dir)

	ann, _ := arrive(t, h, "Ann", "")
	var games protocol.Games
	ann.next(protocol.TypeGames, &games)
	if len(games.Saves) != 2 {
		t.Fatalf("saves %+v", games.Saves)
	}
	ann.say(protocol.TypeLoad, protocol.Load{Save: "g2"})
	if code := ann.failure(); code != protocol.ErrSaveVersion {
		t.Errorf("an older save: %s", code)
	}
	ann.say(protocol.TypeLoad, protocol.Load{Save: "../g1"})
	if code := ann.failure(); code != protocol.ErrNoSave {
		t.Errorf("a path: %s", code)
	}
	cy, _ := arrive(t, h, "Cy", "")
	cy.say(protocol.TypeLoad, protocol.Load{Save: "g1"})
	if code := cy.failure(); code != protocol.ErrNotInSave {
		t.Errorf("a stranger: %s", code)
	}

	// Ann's app lost its token: she comes back by nickname.
	ann.say(protocol.TypeLoad, protocol.Load{Save: "g1"})
	var tb protocol.Table
	ann.next(protocol.TypeTable, &tb)
	if tb.You != 0 || tb.Loaded != "g1" || tb.Seats[1].Name != "Bob" || tb.Seats[1].Connected {
		t.Fatalf("table %+v", tb)
	}
	cy.say(protocol.TypeJoin, protocol.Join{Game: "g1"})
	if code := cy.failure(); code != protocol.ErrNotInSave {
		t.Errorf("a stranger joins: %s", code)
	}
	bob, _ := arrive(t, h, "Bob", "")
	bob.say(protocol.TypeJoin, protocol.Join{Game: "g1"})
	bob.next(protocol.TypeTable, &tb)
	ann.say(protocol.TypeReady, protocol.Ready{Ready: true})
	bob.say(protocol.TypeReady, protocol.Ready{Ready: true})
	var u protocol.Update
	bob.next(protocol.TypeUpdate, &u)
	if u.Seats[0].Name != "Ann" || len(u.View.Players) != 2 {
		t.Fatalf("continued game %+v", u.Seats)
	}
	if _, err := os.Stat(filepath.Join(dir, "g1"+saveExt)); !os.IsNotExist(err) {
		t.Error("the save was not removed once the game went on")
	}
}
