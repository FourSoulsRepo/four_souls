package server

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"testing"
	"time"

	"github.com/FourSoulsRepo/four_souls/internal/protocol"
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/b2"
)

// writeSave puts a saved game of a fresh 2-player game into dir.
func writeSave(t *testing.T, dir, id, engineVersion string) {
	t.Helper()
	g, _, err := engine.NewGame(engine.Setup{Seed: 5, Players: 2, Sets: []engine.CardSet{b2.Set}, BonusSouls: true})
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

// TestBotsSaveAndLoad: a saved game keeps its bots' levels and numbers.
// Loaded, its table seats them at once, ready; nobody takes a bot's
// seat, not even a player called "Bot 1"; the game starts when its
// player is back, and the bots play on. A server that does not allow a
// saved bot's level refuses the load, as it refuses seating that bot
// (ADR 033, Decisions 16.1, 16.2, 17.1; 23.9).
func TestBotsSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	h, _ := botHub(t, dir, time.Millisecond)
	ann, _ := arrive(t, h, "Ann", "")
	ann.say(protocol.TypeCreate, protocol.Create{Seats: 3})
	ann.next(protocol.TypeTable, nil)
	ann.seatBot(1, protocol.BotStandard)
	ann.seatBot(2, protocol.BotHard)
	ann.say(protocol.TypeReady, protocol.Ready{Ready: true})
	saved, asked := "", false
	for saved == "" {
		e := ann.nextOf(protocol.TypeUpdate, protocol.TypeSaved)
		if e.Type == protocol.TypeSaved {
			var s protocol.Saved
			if err := e.Unpack(&s); err != nil {
				t.Fatal(err)
			}
			saved = s.Save
			continue
		}
		var u protocol.Update
		if err := e.Unpack(&u); err != nil {
			t.Fatal(err)
		}
		if u.View.Over {
			t.Fatal("the game ended before the save")
		}
		if !asked && u.Step > 10 && u.View.Turn.Active != 0 {
			ann.say(protocol.TypeSave, struct{}{}) // in a bot's turn
			asked = true
			continue
		}
		if in, ok := passOrForced(u); ok {
			ann.say(protocol.TypeIntent, protocol.Intent{Intent: in})
		}
	}
	f, err := readSave(filepath.Join(dir, "saves"), saved)
	if err != nil {
		t.Fatal(err)
	}
	want := []savedSeat{{Name: "Ann", Token: f.Seats[0].Token}, {Name: "Bot 1", Bot: protocol.BotStandard, BotN: 1}, {Name: "Bot 2", Bot: protocol.BotHard, BotN: 2}}
	if !slices.Equal(f.Seats, want) || f.Seats[0].Token == "" {
		t.Fatalf("saved seats %+v, want %+v", f.Seats, want)
	}

	// A server of standard bots only refuses it.
	strict := NewHub()
	strict.saves, strict.bots = h.saves, []string{protocol.BotStandard}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go strict.Run(ctx)
	sa, _ := arrive(t, strict, "Ann", "")
	sa.say(protocol.TypeLoad, protocol.Load{Save: saved})
	if code, key := sa.refusal(); code != protocol.ErrBadSetup || key != "error.bot_level" {
		t.Errorf("a hard bot on a server of standard bots: %s %s", code, key)
	}

	// "Bot 1" is no player of the save.
	impostor, _ := arrive(t, h, "Bot 1", "")
	impostor.say(protocol.TypeLoad, protocol.Load{Save: saved})
	if code := impostor.failure(); code != protocol.ErrNotInSave {
		t.Errorf("a player called Bot 1 loads: %s", code)
	}

	ann.say(protocol.TypeLoad, protocol.Load{Save: saved})
	var tb protocol.Table
	ann.next(protocol.TypeTable, &tb)
	for i, s := range tb.Seats[1:] {
		if s.Bot != want[i+1].Bot || s.BotN != i+1 || s.Name != want[i+1].Name || !s.Ready {
			t.Errorf("seat %d of the loaded table %+v, want a ready bot", i+1, s)
		}
	}
	impostor.say(protocol.TypeJoin, protocol.Join{Game: tb.Game})
	if code := impostor.failure(); code != protocol.ErrNotInSave {
		t.Errorf("a player called Bot 1 joins: %s", code)
	}
	ann.say(protocol.TypeReady, protocol.Ready{Ready: true})
	var u protocol.Update
	ann.next(protocol.TypeUpdate, &u)
	if u.Step != f.Step || u.Pause != nil {
		t.Fatalf("the loaded game at step %d, pause %+v; want step %d, going on", u.Step, u.Pause, f.Step)
	}
	turn := u.View.Turn.Number
	for u.View.Turn.Number < turn+3 { // a whole turn of a bot
		if u.View.Over {
			return
		}
		if in, ok := passOrForced(u); ok {
			ann.say(protocol.TypeIntent, protocol.Intent{Intent: in})
		}
		u = protocol.Update{}
		ann.next(protocol.TypeUpdate, &u)
		if u.Pause != nil {
			t.Fatalf("the loaded game paused: %+v", u.Pause)
		}
	}
}

// nextOf waits for the next message of one of the given types.
func (p *player) nextOf(types ...string) protocol.Envelope {
	p.t.Helper()
	for {
		select {
		case msg := <-p.conn.out:
			e, err := protocol.Decode(msg)
			if err != nil {
				p.t.Fatal(err)
			}
			if slices.Contains(types, e.Type) {
				return e
			}
		case <-time.After(5 * time.Second):
			p.t.Fatalf("no %v message", types)
			return protocol.Envelope{}
		}
	}
}
