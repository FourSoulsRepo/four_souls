package challenge_test

import (
	"reflect"
	"slices"
	"strings"
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards"
	"github.com/FourSoulsRepo/rules_engine/cards/challenge"
)

// setup.game (24.34): the game type of a situation.

// gameTable is a 2-player table of Isaac and Cain with game gm, the
// leech and fly in slots, and the treasure deck top named.
func gameTable(gm *engine.SituationGame) engine.SituationSetup {
	return engine.SituationSetup{
		Players:  []engine.SituationPlayer{{Character: "isaac"}, {Character: "cain"}},
		Monsters: []engine.CardRef{"leech", "fly"},
		Decks:    map[string][]engine.CardRef{"treasure": {"boomerang", "breakfast", "box", "blank_card", "book_of_sin"}},
		Game:     gm,
	}
}

func build(t *testing.T, s engine.SituationSetup) *engine.Game {
	t.Helper()
	g, err := engine.BuildTable(s, trickSets()...)
	if err != nil {
		t.Fatal(err)
	}
	check(t, g)
	return g
}

// A regular game is the default: no game, and type regular, build the
// same table, with no challenge and no final boss.
func TestGameRegularIsTheDefault(t *testing.T) {
	none := build(t, gameTable(nil))
	regular := build(t, gameTable(&engine.SituationGame{Type: engine.GameRegular}))
	for _, g := range []*engine.Game{none, regular} {
		if g.Challenge != nil || len(g.Bosses) != 0 {
			t.Errorf("a regular game: challenge %v, %d final boss slots", g.Challenge, len(g.Bosses))
		}
	}
	a, err := none.Save()
	if err != nil {
		t.Fatal(err)
	}
	b, err := regular.Save()
	if err != nil {
		t.Fatal(err)
	}
	if string(a) != string(b) {
		t.Error("type regular builds another table than no game")
	}
}

// R-COOP-02, R-KRAMPUS-04, R-KRAMPUS-05: a Krampus challenge starts with
// the timer at 8, the rule card of the difficulty, Pilfering Krampus of
// the difficulty in a final boss slot, and the top 4 treasure cards as
// the ¢, loot, treasure and soul presents with the table's counters.
func TestGameKrampusStarts(t *testing.T) {
	for _, c := range []struct {
		d        string
		counters [4]int // R-KRAMPUS-05
	}{
		{engine.DifficultyNormal, [4]int{8, 2, 3, 3}},
		{engine.DifficultyHard, [4]int{12, 4, 3, 3}},
		{engine.DifficultyUltra, [4]int{16, 4, 3, 3}},
	} {
		g := build(t, gameTable(&engine.SituationGame{Type: engine.GameChallenge, ID: "krampus", Difficulty: c.d}))
		ch := g.Challenge
		suffix := map[string]string{engine.DifficultyNormal: "norm", engine.DifficultyHard: "hard", engine.DifficultyUltra: "ultra"}[c.d]
		if ch == nil || ch.ID != "krampus" || ch.Difficulty != c.d || ch.Timer != 8 || ch.Starter != 0 {
			t.Fatalf("%s: challenge %+v", c.d, ch)
		}
		if got := g.CardOf(ch.Card); got != engine.CardRef("how_the_krampus_stole_christmas_"+suffix) {
			t.Errorf("%s: rule card %s", c.d, got)
		}
		if got := minionCards(g); len(got) != 0 {
			t.Errorf("%s: minions %v", c.d, got)
		}
		if id, ok := g.Bosses[0].TopOf(); len(g.Bosses) != 1 || !ok || g.Object(id).Card != engine.CardRef("pilfering_krampus_"+suffix) {
			t.Errorf("%s: final boss slots %v", c.d, g.Bosses)
		}
		if got := counters(g); got != c.counters {
			t.Errorf("%s: present counters %v, want %v", c.d, got, c.counters)
		}
		for i, want := range []engine.CardRef{"boomerang", "breakfast", "box", "blank_card"} {
			if p := ch.Presents[i]; g.Object(p.Object).Card != want || g.Object(p.Object).Zone.Kind != engine.ZoneOutside {
				t.Errorf("%s: present %s is %s, want %s outside the game", c.d, p.Kind, g.Object(p.Object).Card, want)
			}
		}
		if top := g.Decks[engine.TreasureDeck]; g.Object(top[len(top)-1]).Card != "book_of_sin" {
			t.Errorf("%s: the treasure deck top is not the 5th named card", c.d)
		}
		// As the Go tests start it: the same challenge state.
		tb := krampusOn(t, c.d, gameTable(nil))
		if want := tb.G.Challenge; !reflect.DeepEqual(ch, want) {
			t.Errorf("%s: setup.game %+v, StartChallenge %+v", c.d, ch, want)
		}
	}
}

// R-TRICK-03, R-TRICK-04: Mischievous The Haunt of the difficulty in a
// final boss slot with 12, 14 or 14 trick counters; the Trick/Treat
// card Trick side up; the timer at 8 (R-COOP-02); no presents.
func TestGameTrickStarts(t *testing.T) {
	for _, c := range []struct {
		d, boss string
		tricks  int
	}{
		{engine.DifficultyNormal, hauntNorm, 12},
		{engine.DifficultyHard, "mischievous_the_haunt_hard", 14},
		{engine.DifficultyUltra, hauntUltra, 14},
	} {
		g := build(t, gameTable(&engine.SituationGame{Type: engine.GameChallenge, ID: "trick_or_treat", Difficulty: c.d}))
		ch := g.Challenge
		if ch == nil || ch.ID != "trick_or_treat" || ch.Difficulty != c.d || ch.Timer != 8 || len(ch.Presents) != 0 {
			t.Fatalf("%s: challenge %+v", c.d, ch)
		}
		if g.CardOf(ch.Card) != trickCard || g.Object(ch.Card).Flipped {
			t.Errorf("%s: rule card %s, flipped %v", c.d, g.CardOf(ch.Card), g.Object(ch.Card).Flipped)
		}
		id, _ := g.Bosses[0].TopOf()
		if g.Object(id).Card != engine.CardRef(c.boss) || g.Object(id).CountersOf(challenge.TrickCounter) != c.tricks {
			t.Errorf("%s: boss %s with %d tricks", c.d, g.Object(id).Card, g.Object(id).CountersOf(challenge.TrickCounter))
		}
	}
}

// Each error names what is wrong and the options.
func TestGameErrors(t *testing.T) {
	kr := func(f func(*engine.SituationGame)) *engine.SituationGame {
		gm := &engine.SituationGame{Type: engine.GameChallenge, ID: "krampus", Difficulty: engine.DifficultyNormal}
		f(gm)
		return gm
	}
	two := 2
	for _, c := range []struct {
		name string
		s    func(s *engine.SituationSetup)
		sets []engine.CardSet
		want string
	}{
		{"no type", func(s *engine.SituationSetup) { s.Game = &engine.SituationGame{} }, nil, "type is missing; options [regular challenge]"},
		{"unknown type", func(s *engine.SituationSetup) { s.Game = &engine.SituationGame{Type: "coop"} }, nil, `type "coop" is no game type; options [regular challenge]`},
		{"team", func(s *engine.SituationSetup) { s.Game = &engine.SituationGame{Type: "team"} }, nil, "type team is not in the game yet (Q-03, parked)"},
		{"regular id", func(s *engine.SituationSetup) { s.Game = &engine.SituationGame{Type: "regular", ID: "krampus"} }, nil, "id fits only type challenge, not regular"},
		{"regular holidays", func(s *engine.SituationSetup) { s.Game = &engine.SituationGame{Type: "regular", Holidays: true} }, nil, "holidays fits only type challenge"},
		{"regular rule card", func(s *engine.SituationSetup) { s.Game = &engine.SituationGame{Type: "regular", RuleCard: trickCard} }, nil, "rule_card fits only type challenge"},
		{"no id", func(s *engine.SituationSetup) { s.Game = kr(func(g *engine.SituationGame) { g.ID = "" }) }, nil, "a challenge needs id; options [krampus trick_or_treat]"},
		{"unknown id", func(s *engine.SituationSetup) { s.Game = kr(func(g *engine.SituationGame) { g.ID = "greed" }) }, nil, `challenge "greed" is in none of the sets; options [krampus trick_or_treat]`},
		{"no challenge set", func(s *engine.SituationSetup) { s.Game = kr(func(*engine.SituationGame) {}) }, cards.Sets(), `challenge "krampus" is in none of the sets; options []`},
		{"no difficulty", func(s *engine.SituationSetup) { s.Game = kr(func(g *engine.SituationGame) { g.Difficulty = "" }) }, nil, "challenge krampus needs difficulty; options [normal hard ultra]"},
		{"unknown difficulty", func(s *engine.SituationSetup) { s.Game = kr(func(g *engine.SituationGame) { g.Difficulty = "easy" }) }, nil, `has no difficulty "easy"; options [normal hard ultra]`},
		{"3 players", func(s *engine.SituationSetup) {
			s.Game = kr(func(*engine.SituationGame) {})
			s.Players = append(s.Players, engine.SituationPlayer{Character: "samson"})
		}, nil, "is for 2 players, not 3 (R-COOP-01)"},
		{"bosses", func(s *engine.SituationSetup) {
			s.Game = kr(func(*engine.SituationGame) {})
			s.Bosses = []engine.CardRef{"pilfering_krampus_norm"}
		}, nil, "puts its own final boss out (R-CHAL-02); leave bosses out"},
		{"bonus souls", func(s *engine.SituationSetup) {
			s.Game = kr(func(*engine.SituationGame) {})
			s.BonusSouls = []engine.CardRef{"soul_of_greed"}
		}, nil, "without bonus souls (R-COOP-02)"},
		{"presents in trick", func(s *engine.SituationSetup) {
			s.Game = &engine.SituationGame{Type: "challenge", ID: "trick_or_treat", Difficulty: "normal", Holidays: true}
		}, nil, "challenge trick_or_treat has no presents"},
		{"unknown present", func(s *engine.SituationSetup) {
			s.Game = kr(func(g *engine.SituationGame) { g.Presents = []engine.SituationPresent{{Kind: "coal", Counters: 1}} })
		}, nil, `"coal" is no present; options [cents loot treasure soul]`},
		{"present twice", func(s *engine.SituationSetup) {
			s.Game = kr(func(g *engine.SituationGame) {
				g.Presents = []engine.SituationPresent{{Kind: "loot", Counters: 1}, {Kind: "loot", Destroyed: true}}
			})
		}, nil, "presents: loot is named twice"},
		{"present below 0", func(s *engine.SituationSetup) {
			s.Game = kr(func(g *engine.SituationGame) { g.Presents = []engine.SituationPresent{{Kind: "loot", Counters: -1}} })
		}, nil, "counters -1 is below 1; a present at 0 is destroyed (R-KRAMPUS-08)"},
		{"present names nothing", func(s *engine.SituationSetup) {
			s.Game = kr(func(g *engine.SituationGame) { g.Presents = []engine.SituationPresent{{Kind: "loot"}} })
		}, nil, "loot: names nothing to change"},
		{"destroyed with counters", func(s *engine.SituationSetup) {
			s.Game = kr(func(g *engine.SituationGame) {
				g.Presents = []engine.SituationPresent{{Kind: "loot", Counters: 2, Destroyed: true}}
			})
		}, nil, "a destroyed present has no counters and goes to nobody"},
		{"to without holidays", func(s *engine.SituationSetup) {
			s.Game = kr(func(g *engine.SituationGame) { g.Presents = []engine.SituationPresent{{Kind: "loot", To: &two}} })
		}, nil, "loot: to: a present is gained only at Happy Holidays! (R-KRAMPUS-10); add holidays"},
		{
			"holidays without to", func(s *engine.SituationSetup) { s.Game = kr(func(g *engine.SituationGame) { g.Holidays = true }) }, nil,
			"holidays: the cents present names no player who gained it (to), or destroyed (R-KRAMPUS-10)",
		},
		{"to no seat", func(s *engine.SituationSetup) {
			s.Game = kr(func(g *engine.SituationGame) {
				g.Holidays = true
				g.Presents = []engine.SituationPresent{{Kind: "cents", To: &two}}
			})
		}, nil, "cents: to 2 is not a seat, options 0 to 1"},
		{
			"rule card", func(s *engine.SituationSetup) { s.Game = kr(func(g *engine.SituationGame) { g.RuleCard = trickCard }) }, nil,
			`rule_card "trick_treat_normhardultra" is no side of the rule card; options [how_the_krampus_stole_christmas_norm]`,
		},
		{"fear 0", func(s *engine.SituationSetup) {
			s.Players[0].Counters = []engine.SituationCounter{{Card: engine.PlayerCounters, Name: challenge.FearCounter}}
		}, nil, `player 0 counters: "player" count 0 is below 1`},
	} {
		s := gameTable(nil)
		c.s(&s)
		sets := c.sets
		if sets == nil {
			sets = trickSets()
		}
		_, err := engine.BuildTable(s, sets...)
		if err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s: error %v, want %q", c.name, err, c.want)
		}
	}
}

// R-KRAMPUS-08, R-KRAMPUS-09: presents set up with counters left, and a
// destroyed one in the treasure discard.
func TestGamePresentsSetUp(t *testing.T) {
	g := build(t, gameTable(&engine.SituationGame{
		Type: engine.GameChallenge, ID: "krampus", Difficulty: engine.DifficultyNormal,
		Presents: []engine.SituationPresent{{Kind: engine.PresentCents, Counters: 1}, {Kind: engine.PresentLoot, Destroyed: true}},
	}))
	if got, want := counters(g), [4]int{1, 0, 3, 3}; got != want {
		t.Errorf("counters %v, want %v", got, want)
	}
	if g.PresentLive(engine.PresentLoot) || !g.Challenge.Presents[1].Destroyed {
		t.Error("the loot present is not destroyed")
	}
	if d := g.Discards[engine.TreasureDeck]; len(d) != 1 || g.Object(d[0]).Card != "breakfast" {
		t.Errorf("the treasure discard is %v, want the loot present's breakfast", d)
	}
	// One more destroyed present is no loss: only the fourth ends it (R-KRAMPUS-09).
	if g.Over {
		t.Error("the game is over")
	}
}

// R-KRAMPUS-10, Q-40: after Happy Holidays! each present not destroyed
// is a charged item of the player who gained it, with no counters.
func TestGameHolidaysSetUp(t *testing.T) {
	zero, one := 0, 1
	g := build(t, gameTable(&engine.SituationGame{
		Type: engine.GameChallenge, ID: "krampus", Difficulty: engine.DifficultyNormal, Holidays: true,
		Presents: []engine.SituationPresent{
			{Kind: engine.PresentCents, To: &zero},
			{Kind: engine.PresentLoot, Destroyed: true},
			{Kind: engine.PresentTreasure, To: &one},
			{Kind: engine.PresentSoul, To: &one},
		},
	}))
	if !g.Challenge.Holidays {
		t.Error("no Happy Holidays!")
	}
	for p, want := range [][]engine.CardRef{{"boomerang"}, {"box", "blank_card"}} {
		var got []engine.CardRef
		for _, id := range g.Players[p].InPlay {
			if o := g.Object(id); o.Role == engine.RoleItem && o.Charged && len(o.Counters) == 0 {
				got = append(got, o.Card)
			}
		}
		if !slices.Equal(got, want) {
			t.Errorf("player %d items %v, want %v", p, got, want)
		}
	}
}

// R-TRICK-04, R-TRICK-11: rule_card names the side up; Treat is the
// Trick/Treat card flipped.
func TestGameRuleCardSide(t *testing.T) {
	for _, side := range []engine.CardRef{trickCard, treatCard} {
		g := build(t, gameTable(&engine.SituationGame{Type: engine.GameChallenge, ID: "trick_or_treat", Difficulty: engine.DifficultyNormal, RuleCard: side}))
		if got := g.CardOf(g.Challenge.Card); got != side || g.Object(g.Challenge.Card).Flipped != (side == treatCard) {
			t.Errorf("rule_card %s: side up %s", side, got)
		}
	}
}

// R-TRICK-06: the boss's trick counters in setup.counters are those
// left, not more on top of its start (R-TRICK-03); 0 leaves none.
// Fear counters go on the player (R-TRICK-08).
func TestGameBossCountersAndFear(t *testing.T) {
	for _, n := range []int{3, 0} {
		s := gameTable(&engine.SituationGame{Type: engine.GameChallenge, ID: "trick_or_treat", Difficulty: engine.DifficultyNormal})
		s.Counters = []engine.SituationCounter{{Card: hauntNorm, Name: challenge.TrickCounter, Count: n}}
		s.Players[1].Counters = []engine.SituationCounter{{Card: engine.PlayerCounters, Name: challenge.FearCounter, Count: 2}}
		g := build(t, s)
		id, _ := g.Bosses[0].TopOf()
		if got := g.Object(id).CountersOf(challenge.TrickCounter); got != n {
			t.Errorf("tricks %d, want %d", got, n)
		}
		if g.PlayerCounters(0, challenge.FearCounter) != 0 || g.PlayerCounters(1, challenge.FearCounter) != 2 {
			t.Errorf("fear %d and %d, want 0 and 2", g.PlayerCounters(0, challenge.FearCounter), g.PlayerCounters(1, challenge.FearCounter))
		}
	}
}

// expect.presents and expect.after[].counters on "player": a wrong
// value is a mismatch that names it. The pass is refused in the open
// action phase (R-TURN-09); only the checks matter here.
func TestGameChecks(t *testing.T) {
	payload := `{
  "title": "checks", "sets": ["b2", "challenge"],
  "setup": {
    "game": {"type": "challenge", "id": "krampus", "difficulty": "normal",
      "presents": [{"kind": "loot", "destroyed": true}]},
    "players": [{"character": "isaac", "counters": [{"card": "player", "name": "fear", "count": 1}]}, {"character": "cain"}],
    "active": 0, "monsters": ["leech", "fly"]
  },
  "action": {"player": 0, "kind": "pass"},
  "expect": {"allowed": false,
    "presents": [{"kind": "cents", "counters": 7}, {"kind": "loot", "destroyed": false}, {"kind": "soul", "to": 1}, {"kind": "coal", "counters": 1}],
    "after": [{"player": 0, "counters": [{"card": "player", "name": "fear", "count": 2}]}, {"player": 1, "counters": [{"card": "player", "name": "fear", "count": 0}]}]}
}`
	ans, err := engine.RunSituation([]byte(payload), trickSets()...)
	if err != nil {
		t.Fatal(err)
	}
	want := []string{
		"player 0 fear counters on the player are 1, expected 2",
		"present cents counters is 8, expected 7",
		"present loot destroyed is true, expected false",
		"present soul is not an item of player 1",
		`presents: "coal" is no present, options [cents loot treasure soul]`,
	}
	if !slices.Equal(ans.Mismatch, want) {
		t.Errorf("mismatches\n%q\nwant\n%q", ans.Mismatch, want)
	}
	regular := strings.Replace(payload, `"type": "challenge", "id": "krampus", "difficulty": "normal",
      "presents": [{"kind": "loot", "destroyed": true}]`, `"type": "regular"`, 1)
	ans, err = engine.RunSituation([]byte(regular), trickSets()...)
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Contains(ans.Mismatch, "presents: the game has no presents") {
		t.Errorf("a regular game: mismatches %q", ans.Mismatch)
	}
}

// The test sets reach the files of this tree: its runner passes them,
// so a file names them in sets (24.34).
func TestTrialSetsReachFiles(t *testing.T) {
	payload := `{"title": "trial", "sets": ["b2", "challenge", "krampus_trial", "trick_trial"],
  "setup": {"game": {"type": "challenge", "id": "krampus", "difficulty": "normal"},
    "players": [{"character": "isaac"}, {"character": "cain"}], "monsters": ["mixed_bag", "bare_mob"]},
  "action": {"player": 0, "kind": "attack"}, "expect": {"allowed": true, "prompt": {"options": ["mixed_bag", "bare_mob", "ask.deck.monster"]}}}`
	ans, err := engine.RunSituation([]byte(payload), trickSets()...)
	if err != nil {
		t.Fatal(err)
	}
	if ans.Matches == nil || !*ans.Matches {
		t.Errorf("mismatches %q", ans.Mismatch)
	}
}
