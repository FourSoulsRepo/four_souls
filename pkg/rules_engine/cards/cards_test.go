package cards

import (
	"bytes"
	"os"
	"strings"
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/b2"
)

// TestSetsStartGames starts games with all sets together: a fan set or an
// expansion may not be playable on its own.
func TestSetsStartGames(t *testing.T) {
	for players := 2; players <= 4; players++ {
		g, _, err := engine.NewGame(engine.Setup{Seed: 1, Players: players, Sets: Sets(), BonusSouls: true})
		if err != nil {
			t.Fatalf("%d players: %v", players, err)
		}
		if err := g.CheckInvariants(); err != nil {
			t.Fatalf("%d players: %v", players, err)
		}
	}
	if _, _, err := engine.NewGame(engine.Setup{Seed: 1, Players: 2, Sets: []engine.CardSet{b2.Set}}); err != nil {
		t.Fatalf("the Base Game alone: %v", err)
	}
}

func TestStartingItemsExist(t *testing.T) {
	for _, set := range Sets() {
		refs := map[engine.CardRef]engine.CardDef{}
		for _, c := range set.Cards {
			refs[c.Ref] = c
		}
		for _, c := range set.Cards {
			if c.StartingItem == "" {
				continue
			}
			item, ok := refs[c.StartingItem]
			if !ok || !item.Outside {
				t.Errorf("%s: starting item %s missing or not outside the game", c.Ref, c.StartingItem)
			}
		}
	}
}

func TestFind(t *testing.T) {
	if s, ok := Find("b2"); !ok || s.Name == "" || len(s.Cards) == 0 {
		t.Fatalf("Find(b2) = %v, %v", s.Name, ok)
	}
	if _, ok := Find("nope"); ok {
		t.Fatal("Find(nope) found a set")
	}
}

// FuzzSets plays random games with the real sets: every step is an
// allowed intent, the state stays valid and the game replays exactly.
//
//	go test -run '^$' -fuzz=FuzzSets -fuzztime=5m -fuzzminimizetime=5s ./cards/
func FuzzSets(f *testing.F) {
	f.Add(uint64(1), uint8(0), []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})
	f.Add(uint64(5), uint8(2), bytes.Repeat([]byte{1, 0, 3, 2}, 100))
	f.Add(uint64(9), uint8(1), bytes.Repeat([]byte{0}, 300))
	f.Fuzz(func(t *testing.T, seed uint64, players uint8, choices []byte) {
		playSets(t, engine.Setup{Seed: seed, Players: 2 + int(players%3), Sets: Sets(), BonusSouls: true, Rooms: players&4 != 0}, choices)
	})
}

// FuzzFocusedSets is FuzzSets with only the sets named in
// FOUR_SOULS_SETS (comma-separated codes; empty: all), so a new set's
// cards come up often. The room deck is on in 3 games of 4. It keeps its
// own corpus, so a run fuzzes new inputs from the start.
//
//	FOUR_SOULS_SETS=b2,r go test -run '^$' -fuzz=FuzzFocusedSets -fuzztime=10m ./cards/
func FuzzFocusedSets(f *testing.F) {
	sets := Sets()
	if env := os.Getenv("FOUR_SOULS_SETS"); env != "" {
		sets = nil
		for _, code := range strings.Split(env, ",") {
			s, ok := Find(strings.TrimSpace(code))
			if !ok {
				f.Fatalf("FOUR_SOULS_SETS: unknown set %q", code)
			}
			sets = append(sets, s)
		}
	}
	f.Add(uint64(2), uint8(4), bytes.Repeat([]byte{2, 1, 0, 3}, 100))
	f.Add(uint64(7), uint8(5), bytes.Repeat([]byte{1}, 300))
	f.Fuzz(func(t *testing.T, seed uint64, players uint8, choices []byte) {
		playSets(t, engine.Setup{Seed: seed, Players: 2 + int(players%3), Sets: sets, BonusSouls: true, Rooms: players&12 != 0}, choices)
	})
}

// playSets plays one fuzz game: every step is an allowed intent, the
// state stays valid and the game replays exactly.
func playSets(t *testing.T, s engine.Setup, choices []byte) {
	t.Helper()
	g, _, err := engine.NewGame(s)
	if err != nil {
		t.Fatal(err)
	}
	var steps []engine.StepRecord
	for i, c := range choices[:min(len(choices), 2000)] {
		if g.Over {
			break
		}
		p := g.Prompt()
		allowed := g.Allowed(p.Player)
		if len(allowed) == 0 {
			t.Fatalf("step %d: nothing allowed, prompt %+v", i, p)
		}
		in := allowed[int(c)%len(allowed)]
		if in.Kind == engine.IntentDiscard {
			in.Objects = in.Objects[:p.Count]
		}
		st, err := g.Step(in)
		if err != nil {
			t.Fatalf("step %d: allowed %+v refused: %v", i, in, err)
		}
		steps = append(steps, st)
		if err := g.CheckInvariants(); err != nil {
			t.Fatalf("step %d after %+v: %v", i, in, err)
		}
	}
	if _, err := engine.Replay(s, steps); err != nil {
		t.Fatalf("replay: %v", err)
	}
}

// TestStatMonsterBans: with the monsters that change other monsters
// banned, a fresh deal shows every monster as printed (14.10). The e2e
// tests ban the same list (e2e/tests/game.ts).
func TestStatMonsterBans(t *testing.T) {
	statMonsters := []engine.CardRef{"delirium", "stoney"}
	changed := func(bans []engine.CardRef) int {
		n := 0
		for seed := range uint64(400) {
			s := engine.Setup{Seed: seed, Players: 2, Sets: []engine.CardSet{b2.Set}, BonusSouls: true, Banned: bans}
			g, _, err := engine.NewGame(s)
			if err != nil {
				t.Fatal(err)
			}
			for _, slot := range g.View(engine.Viewer{Kind: engine.ViewSpectator}).Monsters {
				st := slot.Top.Stats
				if st != nil && (st.HP != st.PrintedHP || st.DC != st.PrintedDC || st.ATK != st.PrintedATK) {
					n++
				}
			}
		}
		return n
	}
	if changed(nil) == 0 {
		t.Fatal("no deal changed a monster; the test checks nothing")
	}
	if n := changed(statMonsters); n > 0 {
		t.Errorf("%d fresh monsters changed with %v banned; a monster is missing from the list", n, statMonsters)
	}
}
