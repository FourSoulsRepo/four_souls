package challenge_test

import (
	"math/rand/v2"
	"slices"
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards"
	"github.com/FourSoulsRepo/rules_engine/cards/challenge"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

// mixedSet has monsters the real sets lack, for All I Want for
// Christmas (Q-32): a mixed reward of 2¢ and 1 loot.
var mixedSet = engine.CardSet{Code: "krampus_trial", Name: "krampus_trial", Cards: []engine.CardDef{
	{
		Ref: "mixed_bag", Kind: engine.MonsterCard, HP: 1, DC: 2, ATK: 1, Soul: 1, Outside: true,
		Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 2}, {Kind: engine.RewardLoot, Amount: 1}},
	},
}}

func sets() []engine.CardSet { return append(cards.Sets(), challenge.Set, mixedSet) }

// krampus builds a Krampus table at difficulty d: Isaac (active) and
// Cain unless the seats name others, the monsters leech and fly, and
// the given minions. The presents come from the treasure deck.
func krampus(t *testing.T, d string, isaac, cain engine.SituationPlayer, minions ...engine.CardRef) *enginetest.Table {
	t.Helper()
	return krampusOn(t, d, engine.SituationSetup{
		Players:  []engine.SituationPlayer{isaac, cain},
		Monsters: []engine.CardRef{"leech", "fly"},
		Minions:  minions,
	})
}

// krampusOn builds a Krampus table from a full setup.
func krampusOn(t *testing.T, d string, s engine.SituationSetup) *enginetest.Table {
	t.Helper()
	for i, c := range []engine.CardRef{"isaac", "cain"} {
		if s.Players[i].Character == "" {
			s.Players[i].Character = c
		}
	}
	tb := enginetest.NewSetup(t, s, sets()...)
	if err := tb.G.StartChallenge(&challenge.Krampus, d); err != nil {
		t.Fatal(err)
	}
	return tb
}

// player is a player's option label.
func player(tb *enginetest.Table, p engine.PlayerID) string {
	return tb.G.Plain(engine.PlayerOption(p))
}

// attackOptions declares an attack and returns its target options.
func attackOptions(tb *enginetest.Table) []string {
	tb.T.Helper()
	if _, err := tb.G.Apply(engine.Intent{Player: tb.G.Turn.Active, Kind: engine.IntentAttack}); err != nil {
		tb.T.Fatal(err)
	}
	tb.Pass(int(tb.G.Turn.Active))
	tb.Pass(int(tb.G.PlayerAfter(tb.G.Turn.Active)))
	return tb.Options()
}

func boost(tb *enginetest.Table, p engine.PlayerID, atk int) {
	tb.G.Boosts = append(tb.G.Boosts, engine.Boost{Stat: engine.StatPlayerATK, Player: p, Amount: atk})
}

// souls is n soul cards for a seat.
func souls(n int) []engine.CardRef {
	out := make([]engine.CardRef, n)
	for i := range out {
		out[i] = "fly"
	}
	return out
}

func boss(tb *enginetest.Table) engine.ObjectID {
	tb.T.Helper()
	id, ok := tb.G.Bosses[0].TopOf()
	if !ok {
		tb.T.Fatal("no Pilfering Krampus")
	}
	return id
}

// setCounters sets the present of kind to n counters (n > 0).
func setCounters(tb *enginetest.Table, kind string, n int) {
	tb.T.Helper()
	tb.G.AddPresentCounters(kind, n-tb.G.PresentCounters(kind), engine.NoPlayer)
	if got := tb.G.PresentCounters(kind); got != n {
		tb.T.Fatalf("the %s present has %d counters, want %d", kind, got, n)
	}
}

func counters(g *engine.Game) [4]int {
	var out [4]int
	for i, k := range []string{engine.PresentCents, engine.PresentLoot, engine.PresentTreasure, engine.PresentSoul} {
		out[i] = g.PresentCounters(k)
	}
	return out
}

func minionCards(g *engine.Game) []engine.CardRef {
	var out []engine.CardRef
	for _, s := range g.Minions {
		if id, ok := s.TopOf(); ok {
			out = append(out, g.Object(id).Card)
		}
	}
	return out
}

func count(ev []engine.Event, kind engine.EventKind) int {
	n := 0
	for _, e := range ev {
		if e.Kind == kind {
			n++
		}
	}
	return n
}

func check(t *testing.T, g *engine.Game) {
	t.Helper()
	if err := g.CheckInvariants(); err != nil {
		t.Fatal(err)
	}
}

// R-KRAMPUS-04, R-KRAMPUS-05, R-KRAMPUS-06: per difficulty the boss is
// in a final boss slot; the top 4 treasure cards lie outside the game,
// face down, as the ¢, loot, treasure and soul presents with their
// counters. The view shows kinds and counters, never the cards.
func TestKrampusSetup(t *testing.T) {
	for d, want := range map[string]struct {
		boss     engine.CardRef
		counters [4]int
	}{
		engine.DifficultyNormal: {"pilfering_krampus_norm", [4]int{8, 2, 3, 3}},
		engine.DifficultyHard:   {"pilfering_krampus_hard", [4]int{12, 4, 3, 3}},
		engine.DifficultyUltra:  {"pilfering_krampus_ultra", [4]int{16, 4, 3, 3}},
	} {
		s := engine.Setup{Seed: 3, Players: 2, Sets: append(cards.Sets(), challenge.Set), Challenge: &challenge.Krampus, Difficulty: d}
		g, _, err := engine.NewGame(s)
		if err != nil {
			t.Fatal(err)
		}
		b, _ := g.Bosses[0].TopOf()
		if !g.IsFinalBoss(b) || g.Object(b).Card != want.boss {
			t.Errorf("%s: boss %s", d, g.Object(b).Card)
		}
		if got := counters(g); got != want.counters {
			t.Errorf("%s: counters %v, want %v", d, got, want.counters)
		}
		var kinds []string
		for _, p := range g.Challenge.Presents {
			kinds = append(kinds, p.Kind)
			o := g.Object(p.Object)
			if o.Zone.Kind != engine.ZoneOutside || o.Role != engine.RoleNone || g.Kind(p.Object) != engine.TreasureCard {
				t.Errorf("%s: present %+v", d, o)
			}
		}
		if !slices.Equal(kinds, []string{"cents", "loot", "treasure", "soul"}) {
			t.Errorf("%s: kinds %v", d, kinds)
		}
		v := g.View(engine.Viewer{Kind: engine.ViewPlayer, Player: 0})
		for i, p := range v.Challenge.Presents {
			if p.Card != "" || p.Counters != want.counters[i] || p.Destroyed || p.Flipped || p.Player != engine.NoPlayer {
				t.Errorf("%s: present view %+v", d, p)
			}
		}
		check(t, g)
	}
}

// R-KRAMPUS-06, Q-40: presents are not in play; no ability can choose
// them. A "destroy an item" choice lists the players' items only.
func TestPresentsNotInPlay(t *testing.T) {
	tb := krampus(t, engine.DifficultyNormal, engine.SituationPlayer{Items: []engine.CardRef{"the_d20"}}, engine.SituationPlayer{})
	for _, p := range tb.G.Challenge.Presents {
		if slices.Contains(tb.G.Players[0].InPlay, p.Object) || slices.Contains(tb.G.Players[1].InPlay, p.Object) {
			t.Fatal("a present is in a play area")
		}
	}
	if _, err := tb.G.Apply(engine.Intent{Player: 0, Kind: engine.IntentActivate, Objects: []engine.ObjectID{tb.Find(0, "the_d20")}}); err == nil {
		for _, o := range tb.Options() {
			for _, p := range tb.G.Challenge.Presents {
				if o == string(tb.G.Object(p.Object).Card) {
					t.Errorf("an ability offers the present %s", o)
				}
			}
		}
	}
}

// R-KRAMPUS-15, Q-36: the regular Krampus stays in the monster deck.
func TestRegularKrampusStays(t *testing.T) {
	s := engine.Setup{Seed: 1, Players: 2, Sets: append(cards.Sets(), challenge.Set), Challenge: &challenge.Krampus, Difficulty: engine.DifficultyNormal}
	g, _, err := engine.NewGame(s)
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, id := range g.Objects {
		found = found || (id.Card == "krampus" && id.Zone.Kind != engine.ZoneGone)
	}
	if !found {
		t.Error("the regular Krampus is out of the game")
	}
}

// R-KRAMPUS-01: seeded whole games end. Players who start with 5
// souls (Happy Holidays! at once) and always attack, Pilfering Krampus
// first, win; random players lose, by the timer or the presents.
// Replays and saves keep the challenge.
func TestKrampusWholeGames(t *testing.T) {
	for _, c := range []struct {
		seed  uint64
		d     string
		souls int
		won   bool
	}{
		{3, engine.DifficultyNormal, 5, true},
		{11, engine.DifficultyHard, 5, true},
		{1, engine.DifficultyNormal, 0, false},
		{2, engine.DifficultyHard, 0, false},
		{3, engine.DifficultyUltra, 0, false},
	} {
		g := playKrampus(t, c.seed, c.d, c.souls)
		if !g.Over || g.Challenge.Lost == c.won {
			t.Errorf("seed %d %s: over %v, lost %v (%s)", c.seed, c.d, g.Over, g.Challenge.Lost, g.Challenge.Reason.Key)
		}
		if g.Challenge.Lost && len(g.Winners) != 0 || !g.Challenge.Lost && len(g.Winners) != 2 {
			t.Errorf("seed %d %s: winners %v (R-COOP-10)", c.seed, c.d, g.Winners)
		}
	}
}

// playKrampus plays a whole Krampus game. With souls, the players start
// with that many monster souls between them and play greedy: they
// only attack, Pilfering Krampus first, and never forgo. Without, they
// pick at random.
func playKrampus(t *testing.T, seed uint64, d string, souls int) *engine.Game {
	t.Helper()
	s := engine.Setup{Seed: seed, Players: 2, Sets: append(cards.Sets(), challenge.Set), Challenge: &challenge.Krampus, Difficulty: d}
	g, _, err := engine.NewGame(s)
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < souls; {
		g.MillTop(engine.MonsterDeck)
		pile := g.Discards[engine.MonsterDeck]
		if top := pile[len(pile)-1]; g.Kind(top) == engine.MonsterCard {
			g.SoulFromDiscard(engine.MonsterDeck, g.Object(top).Card, engine.PlayerID(i%2))
			i++
		}
	}
	base, err := g.Save()
	if err != nil {
		t.Fatal(err)
	}
	rng := rand.New(rand.NewPCG(seed, 7)) //nolint:gosec // test moves
	var steps []engine.StepRecord
	for step := 0; step < 20000 && !g.Over; step++ {
		p := g.Prompt()
		in := choose(g, g.Allowed(p.Player), rng, souls > 0)
		if in.Kind == engine.IntentDiscard {
			in.Objects = in.Objects[:p.Count]
		}
		st, serr := g.Step(in)
		if serr != nil {
			t.Fatalf("step %d: %+v: %v", step, in, serr)
		}
		steps = append(steps, st)
		check(t, g)
	}
	r, err := engine.ReplayFrom(base, steps, s.Sets...)
	if err != nil {
		t.Fatalf("replay: %v", err)
	}
	a, err := g.Save()
	if err != nil {
		t.Fatal(err)
	}
	b, err := r.Save()
	if err != nil {
		t.Fatal(err)
	}
	if string(a) != string(b) {
		t.Error("the replay ends elsewhere")
	}
	if _, err := engine.Load(a, s.Sets...); err != nil {
		t.Fatalf("load: %v", err)
	}
	return g
}

func choose(g *engine.Game, allowed []engine.Intent, rng *rand.Rand, greedy bool) engine.Intent {
	if !greedy {
		var acts []engine.Intent
		for _, in := range allowed {
			if in.Kind != engine.IntentPass && in.Kind != engine.IntentEndTurn {
				acts = append(acts, in)
			}
		}
		if len(acts) > 0 && rng.IntN(3) > 0 {
			return acts[rng.IntN(len(acts))]
		}
		return allowed[rng.IntN(len(allowed))]
	}
	if g.Prompt().Kind == engine.PromptChoose {
		if i := greedyPick(g); i >= 0 {
			return allowed[i]
		}
		return allowed[rng.IntN(len(allowed))]
	}
	open := g.Turn.Step == engine.StepAction && !g.Turn.DeathEnd && !g.Attack.On && !g.Purchase.On && len(g.Stack) == 0
	kinds := []engine.IntentKind{engine.IntentPass, engine.IntentEndTurn}
	if open {
		kinds = []engine.IntentKind{engine.IntentAttack, engine.IntentEndTurn}
	}
	for _, kind := range kinds {
		for _, in := range allowed {
			if in.Kind == kind {
				return in
			}
		}
	}
	return allowed[rng.IntN(len(allowed))]
}

// greedyPick answers a choice: Pilfering Krampus as the target, else
// the monster with the most soul and the least HP; never forgo. -1:
// any option.
func greedyPick(g *engine.Game) int {
	if i := slices.Index(g.PlainOptions(), "ask.no"); i >= 0 {
		return i
	}
	c := g.Choice
	if c == nil || c.Purpose != engine.ChooseAttackTarget || len(c.Objects) == 0 {
		return -1
	}
	best, score := 0, -1000
	for i, id := range c.Objects {
		sc := 10*g.Def(id).Soul - g.HP(id)
		if g.IsFinalBoss(id) {
			sc = 1000
		}
		if sc > score {
			best, score = i, sc
		}
	}
	return best
}
