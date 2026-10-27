package cards

import (
	"bytes"
	"fmt"
	"os"
	"slices"
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

// ADR 023: a back face is never dealt, even when it is picked by name.
func TestBackFaceIsNotDealt(t *testing.T) {
	_, _, err := engine.NewGame(engine.Setup{Seed: 1, Players: 2, Sets: Sets(), Characters: []engine.CardRef{"amgine_eht", "cain"}})
	if err == nil {
		t.Fatal("The Enigma's back was dealt as a character")
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
		playSets(t, engine.Setup{Seed: seed, Players: 2 + int(players%3), Sets: Sets(), BonusSouls: true, Rooms: players&4 != 0, DeckRatios: players&8 != 0}, choices)
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
	k := 0
	if len(choices) > 0 {
		k = int(choices[len(choices)-1]) * len(steps) / 255
	}
	checkUndoRedo(t, g, s, steps, k)
}

// checkUndoRedo rolls the game back to before step k from the start and
// plays the undone steps again: the game must end in the same bytes
// (ADR 020, item 7). The checksums cover the RNG, so the undo is exact.
func checkUndoRedo(t *testing.T, g *engine.Game, s engine.Setup, steps []engine.StepRecord, k int) {
	t.Helper()
	start, _, err := engine.NewGame(s)
	if err != nil {
		t.Fatal(err)
	}
	base, err := start.Save()
	if err != nil {
		t.Fatal(err)
	}
	r, err := engine.ReplayFrom(base, steps[:k], s.Sets...)
	if err != nil {
		t.Fatalf("undo to before step %d: %v", k, err)
	}
	undone, err := r.Save()
	if err != nil {
		t.Fatal(err)
	}
	if r, err = engine.ReplayFrom(undone, steps[k:], s.Sets...); err != nil {
		t.Fatalf("redo from step %d: %v", k, err)
	}
	a, errA := r.Save()
	b, errB := g.Save()
	if errA != nil || errB != nil {
		t.Fatal(errA, errB)
	}
	if !bytes.Equal(a, b) {
		t.Fatalf("undo to before step %d, then redo, ends in a different state", k)
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

// TestDeckRatios: with every set the ratios fit for any player count,
// and each group has its count in the game; the Base Game alone is short
// (R-SETUP-17, GS-07).
func TestDeckRatios(t *testing.T) {
	groups := map[engine.CardRef]engine.RatioGroup{}
	for _, s := range Sets() {
		for ref, g := range s.Groups {
			groups[ref] = g
		}
	}
	for players := 2; players <= 4; players++ {
		s := engine.Setup{Seed: uint64(players), Players: players, Sets: Sets(), BonusSouls: true, DeckRatios: true}
		short, err := engine.RatioShortfall(s)
		if err != nil || short != nil {
			t.Fatalf("%d players: shortfall %v, %v", players, short, err)
		}
		g, _, err := engine.NewGame(s)
		if err != nil {
			t.Fatalf("%d players: %v", players, err)
		}
		counts := map[engine.RatioGroup]int{}
		for _, o := range g.Objects {
			if o.Zone.Kind != engine.ZoneGone && o.Zone.Kind != engine.ZoneOutside {
				counts[groups[o.Card]]++
			}
		}
		for grp := engine.GroupPassiveItems; grp.Count() > 0; grp++ {
			if counts[grp] != grp.Count() {
				t.Errorf("%d players: %s %d cards, want %d", players, grp.Name(), counts[grp], grp.Count())
			}
		}
	}
	base := engine.Setup{Seed: 1, Players: 2, Sets: []engine.CardSet{b2.Set}, DeckRatios: true}
	short, err := engine.RatioShortfall(base)
	if err != nil {
		t.Fatal(err)
	}
	want := engine.RatioShort{Group: engine.GroupButterBeans, Have: 3, Need: 5}
	if !slices.Contains(short, want) {
		t.Errorf("Base Game shortfall %v lacks %v", short, want)
	}
	if _, _, err := engine.NewGame(base); err == nil {
		t.Error("the Base Game alone started with deck ratios")
	}
}

// TestTakes pins the abilities that take from players no target names
// (Ability.Takes, 15.85). A new card that steals, or makes a player
// discard, lose or destroy, picking them as it resolves, joins the list.
func TestTakes(t *testing.T) {
	want := []string{
		"b2/chaos/0 others", "b2/cursed_fatty/0 event", "b2/cursed_keeper_head/0 event",
		"b2/dank_globin/0 others", "b2/dead_bird/0 event", "b2/finger/0 event",
		"b2/greed/0 others", "b2/greed_event/0 others", "b2/greedling/0 others",
		"b2/keeper_head/0 event", "b2/moms_dead_hand/0 others", "b2/moms_shovel/0 others",
		"b2/monstros_tooth/0 others", "b2/pride/0 others", "b2/remote_detonator/0 others",
		"b2/ring_of_flies/0 others", "b2/scolex/0 event", "b2/tarot_cloth/0 event",
		"b2/the_lamb/0 others", "b2/wizoob/0 others", "b2/xv_the_devil/0 others",
		"g2/joker/0 others",
		"fsp2/forget_me_now/0 others", "fsp2/magnet/0 event", "fsp2/pills_fsp2/0 others",
		"r/betrayal/0 others", "r/charmed_fatty/1 event", "r/charmed_keeper/1 event",
		"r/clog/0 others", "r/dice_room_2/0 others", "r/dice_room_6/0 others",
		"r/r_key/0 others", "r/shadow_of_war/0 event",
	}
	names := map[engine.Takes]string{engine.TakesOthers: "others", engine.TakesEventPlayer: "event"}
	var got []string
	for _, set := range Sets() {
		for _, c := range set.Cards {
			for i, a := range c.Abilities {
				if a.Takes != engine.TakesNobody {
					got = append(got, fmt.Sprintf("%s/%s/%d %s", set.Code, c.Ref, i, names[a.Takes]))
				}
			}
		}
	}
	if !slices.Equal(got, want) {
		t.Errorf("marked abilities:\n%s\nwant:\n%s", strings.Join(got, "\n"), strings.Join(want, "\n"))
	}
}

// TestTapMatchesAbilities checks CardDef.Tap on treasures: it is set
// exactly when the card has a ↷ ability (R-ABIL-10), or copies one
// (Placebo). The death penalty and "deactivate an item" read it
// (R-DEATH-14, R-MECH-54; 19.65).
func TestTapMatchesAbilities(t *testing.T) {
	for _, set := range Sets() {
		for _, c := range set.Cards {
			if c.Kind != engine.TreasureCard {
				continue
			}
			tap := c.CopiesTapAbilities
			for _, a := range c.Abilities {
				if slices.Contains(a.Costs, engine.Tap()) {
					tap = true
				}
			}
			if c.Tap != tap {
				t.Errorf("%s/%s: Tap %v, ↷ ability %v", set.Code, c.Ref, c.Tap, tap)
			}
		}
	}
}
