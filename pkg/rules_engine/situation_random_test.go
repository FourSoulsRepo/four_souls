package rulesengine

import (
	"strconv"
	"strings"
	"testing"
)

// randomSet has a card with two picks at random and one that makes
// the attack's target random (24.30).
var randomSet = CardSet{Name: "random", Cards: []CardDef{
	{Ref: "darts", Kind: TreasureCard, Outside: true, Tap: true, Abilities: []Ability{{
		Kind: Activated, Text: "↷: Deal 1 damage to a player chosen at random, then 1 damage to a player chosen at random.",
		Costs: []Cost{Tap()},
		Effects: []Effect{Ask(func(c *Ctx, a []int) {
			for _, i := range a {
				c.G.DealDamageTo(Target{Player: PlayerID(i), IsPlayer: true}, 1, c.Controller, c.Source)
			}
		}, dartPick, dartPick)},
	}}},
	{Ref: "lure", Kind: TreasureCard, Outside: true, MustAttackRandom: func(g *Game, _ ObjectID) []ObjectID {
		var out []ObjectID
		for _, s := range g.Monsters {
			if id, ok := s.TopOf(); ok {
				out = append(out, id)
			}
		}
		return out
	}},
}}

var dartPick = Question{Random: true, Options: func(c *Ctx, _ []int) []Msg {
	var out []Msg
	for _, pl := range c.G.Players {
		out = append(out, PlayerOption(pl.ID))
	}
	return out
}}

// darts is a situation where hero_a throws the darts at three players.
func darts(seed int, random string) []byte {
	return []byte(`{"title": "darts", "sets": ["test", "random"], "seed": ` + strconv.Itoa(seed) + `,
		"setup": {"players": [{"character": "hero_a", "items": ["darts"]}, {"character": "hero_b"}, {"character": "hero_c"}], "active": 0` + random + `},
		"action": {"player": 0, "kind": "activate", "card": "darts"}}`)
}

// damages are the players' damage after the situation's events.
func damages(t *testing.T, data []byte) []int {
	t.Helper()
	ans, err := RunSituation(data, testSet, randomSet)
	if err != nil {
		t.Fatal(err)
	}
	out := make([]int, 3)
	for _, e := range ans.Events {
		if e.Kind == EvDamaged {
			out[e.Player] += e.Amount
		}
	}
	return out
}

// A named pick is used: both darts hit the players the setup names.
func TestSituationRandomNamedPicks(t *testing.T) {
	for seed := range 10 {
		got := damages(t, darts(seed, `, "random": ["player 3 (hero_c)", "player 3 (hero_c)"]`))
		if got[0] != 0 || got[1] != 0 || got[2] != 2 {
			t.Errorf("seed %d: damage %v; want both darts on player 3: [0 0 2]", seed, got)
		}
	}
}

// Picks after the list are random from the seed, as dice after rolls:
// over a few seeds the second dart hits more than one player.
func TestSituationRandomAfterTheList(t *testing.T) {
	seen := map[int]bool{}
	for seed := range 20 {
		got := damages(t, darts(seed, `, "random": ["player 2 (hero_b)"]`))
		if got[1] < 1 {
			t.Fatalf("seed %d: damage %v; want the first dart on player 2", seed, got)
		}
		for p, d := range got {
			if d > 0 && (p != 1 || d == 2) {
				seen[p] = true
			}
		}
	}
	if len(seen) < 2 {
		t.Errorf("the second dart hit only %v over 20 seeds; want a random player", seen)
	}
}

// A pick that is no option is an error naming the options; so is an
// empty one.
func TestSituationRandomBadPick(t *testing.T) {
	_, err := RunSituation(darts(1, `, "random": ["player 2 (hero_c)"]`), testSet, randomSet)
	want := `random: pick 1 "player 2 (hero_c)" is not an option of ["player 1 (hero_a)" "player 2 (hero_b)" "player 3 (hero_c)"]`
	if err == nil || !strings.Contains(err.Error(), want) {
		t.Errorf("error %v; want %q", err, want)
	}
	_, err = RunSituation(darts(1, `, "random": ["player 2 (hero_b)", ""]`), testSet, randomSet)
	if err == nil || !strings.Contains(err.Error(), "random: pick 2 is empty") {
		t.Errorf("error %v; want pick 2 empty", err)
	}
}

// A random attack's target (MustAttackRandom, Q-31) reads the same
// list: the picked card is attacked.
func TestSituationRandomAttackTarget(t *testing.T) {
	for _, c := range []struct {
		pick string
		want CardRef
	}{{"gaper", "gaper"}, {"boss", "boss"}} {
		g, err := BuildTable(SituationSetup{
			Players:  []SituationPlayer{{Character: "hero_a", Items: []CardRef{"lure"}}, {Character: "hero_b"}},
			Monsters: []CardRef{"gaper", "boss"},
			Rolls:    []int{6, 6, 6, 6}, // every attack roll hits, so nobody dies
			Random:   []string{c.pick},
		}, testSet, randomSet)
		if err != nil {
			t.Fatal(err)
		}
		ev, err := g.Apply(Intent{Player: 0, Kind: IntentAttack})
		if err != nil {
			t.Fatal(err)
		}
		more, err := g.Settle()
		if err != nil {
			t.Fatal(err)
		}
		ev = append(ev, more...)
		var got CardRef
		for _, e := range ev {
			if e.Kind == EvRandomPick {
				got = g.Object(e.Object).Card
			}
		}
		if got != c.want || g.takePickErr() != nil {
			t.Errorf("random %q: picked %q; want %q", c.pick, got, c.want)
		}
	}
	g, err := BuildTable(SituationSetup{
		Players:  []SituationPlayer{{Character: "hero_a", Items: []CardRef{"lure"}}, {Character: "hero_b"}},
		Monsters: []CardRef{"gaper", "boss"},
		Rolls:    []int{6, 6, 6, 6},
		Random:   []string{"dummy"},
	}, testSet, randomSet)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := g.Apply(Intent{Player: 0, Kind: IntentAttack}); err != nil {
		t.Fatal(err)
	}
	if _, err := g.Settle(); err != nil {
		t.Fatal(err)
	}
	if err := g.takePickErr(); err == nil || !strings.Contains(err.Error(), `"dummy" is not an option of ["gaper" "boss"]`) {
		t.Errorf("error %v; want dummy refused with the options", err)
	}
}
