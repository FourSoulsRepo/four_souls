package cards

import (
	"bytes"
	"fmt"
	"math/rand/v2"
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

// playSteps plays up to n random allowed steps on g, as simulate does,
// and calls each after every step; it stops early when the game ends.
func playSteps(t testing.TB, g *engine.Game, rng *rand.Rand, n int, each func(step int)) {
	t.Helper()
	for step := 0; step < n && !g.Over; step++ {
		p := g.Prompt()
		allowed := g.Allowed(p.Player)
		if len(allowed) == 0 {
			t.Fatalf("step %d: nothing allowed, prompt %+v", step, p)
		}
		in := pick(rng, allowed)
		if in.Kind == engine.IntentDiscard {
			in.Objects = in.Objects[:p.Count]
		}
		if _, err := g.Apply(in); err != nil {
			t.Fatalf("step %d: allowed %+v refused: %v", step, in, err)
		}
		if each != nil {
			each(step)
		}
	}
}

// cloneSetups are games of every set (b2 alone, b2 with each other set,
// all sets at 2 to 4 players) and of each challenge at each difficulty.
func cloneSetups() map[string]engine.Setup {
	out := map[string]engine.Setup{}
	all := Sets()
	for i, s := range all {
		sets := []engine.CardSet{all[0]}
		if i > 0 {
			sets = append(sets, s)
		}
		out["set "+s.Code] = engine.Setup{Seed: uint64(10 + i), Players: 2 + i%3, Sets: sets, BonusSouls: true, Rooms: true}
	}
	for n := 2; n <= 4; n++ {
		out[fmt.Sprintf("all sets, %d players", n)] = engine.Setup{
			Seed: uint64(20 + n), Players: n, Sets: all, BonusSouls: true, Rooms: n != 3, DeckRatios: n == 3,
		}
	}
	for i, c := range Challenges() {
		for j, d := range krampusLevels {
			s := krampusSetup(uint64(30+3*i+j), d)
			s.Challenge = c
			out[c.ID+" "+d] = s
		}
	}
	return out
}

// TestCloneSavesTheSame: a copy has the game's checksum and save bytes,
// and shares no list, map or pointer with it, every 10 steps of
// simulated games (ADR 030, Decision 3.2).
func TestCloneSavesTheSame(t *testing.T) {
	for name, s := range cloneSetups() {
		t.Run(name, func(t *testing.T) {
			g, _, err := engine.NewGame(s)
			if err != nil {
				t.Fatal(err)
			}
			check := func(step int) {
				t.Helper()
				c, err := g.Clone()
				if err != nil {
					t.Fatal(err)
				}
				if err := sameGame(g, c); err != nil {
					t.Fatalf("step %d: %v", step, err)
				}
				if err := enginetest.NoSharing(g, c); err != nil {
					t.Fatalf("step %d: %v", step, err)
				}
			}
			check(-1)
			playSteps(t, g, rand.New(rand.NewPCG(s.Seed, 99)), 800, func(step int) { //nolint:gosec // test moves, not security
				if step%10 == 0 {
					check(step)
				}
			})
			check(800)
		})
	}
}

// TestCloneLeavesTheGame: 50 steps on a copy never change the game.
func TestCloneLeavesTheGame(t *testing.T) {
	for name, s := range cloneSetups() {
		t.Run(name, func(t *testing.T) {
			g, _, err := engine.NewGame(s)
			if err != nil {
				t.Fatal(err)
			}
			rng := rand.New(rand.NewPCG(s.Seed, 99)) //nolint:gosec // test moves, not security
			for round := range 6 {
				playSteps(t, g, rng, 60, nil)
				if g.Over {
					break
				}
				before, err := g.Save()
				if err != nil {
					t.Fatal(err)
				}
				sum, err := g.Checksum()
				if err != nil {
					t.Fatal(err)
				}
				c, err := g.Clone()
				if err != nil {
					t.Fatal(err)
				}
				playSteps(t, c, rand.New(rand.NewPCG(uint64(round), 7)), 50, nil) //nolint:gosec // test moves, not security
				after, err := g.Save()
				if err != nil {
					t.Fatal(err)
				}
				if got, err := g.Checksum(); err != nil || got != sum || !bytes.Equal(before, after) {
					t.Fatalf("round %d: 50 steps on the copy changed the game (err %v)", round, err)
				}
				if err := g.CheckInvariants(); err != nil {
					t.Fatalf("round %d: %v", round, err)
				}
			}
		})
	}
}

// sameGame says whether c saves to the bytes of g and has its checksum.
func sameGame(g, c *engine.Game) error {
	a, err := g.Save()
	if err != nil {
		return err
	}
	b, err := c.Save()
	if err != nil {
		return err
	}
	if !bytes.Equal(a, b) {
		return fmt.Errorf("save differs:\n%s\n%s", a, b)
	}
	sa, err := g.Checksum()
	if err != nil {
		return err
	}
	sb, err := c.Checksum()
	if err != nil {
		return err
	}
	if sa != sb {
		return fmt.Errorf("checksum %x, copy %x", sa, sb)
	}
	return nil
}

// BenchmarkClone copies a 4-player game with every set after 300
// random steps: the copy a bot makes for each try (ADR 030, Decision 3).
func BenchmarkClone(b *testing.B) {
	s := engine.Setup{Seed: 7, Players: 4, Sets: Sets(), BonusSouls: true, Rooms: true}
	g, _, err := engine.NewGame(s)
	if err != nil {
		b.Fatal(err)
	}
	playSteps(b, g, rand.New(rand.NewPCG(7, 99)), 300, nil) //nolint:gosec // test moves, not security
	b.ReportAllocs()
	for b.Loop() {
		if _, err := g.Clone(); err != nil {
			b.Fatal(err)
		}
	}
}
