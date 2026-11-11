package challenge_test

import (
	"slices"
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/challenge"
)

// R-TRICK-06, Q-30: while the boss has more trick counters than the
// souls of both players, it can't be attacked and takes no damage;
// equal is no guard. Trick counters never go below 0.
func TestTrickGuard(t *testing.T) {
	tb := trick(t, engine.DifficultyNormal, engine.SituationPlayer{Souls: souls(2)}, engine.SituationPlayer{Souls: souls(1)})
	setTricks(tb, 4)
	if opts := attackOptions(tb); slices.Contains(opts, hauntNorm) || !slices.Contains(opts, "leech") {
		t.Errorf("4 tricks, 3 souls: attack options %v", opts)
	}
	tb = trick(t, engine.DifficultyNormal, engine.SituationPlayer{Souls: souls(2)}, engine.SituationPlayer{Souls: souls(1)})
	setTricks(tb, 4)
	b := boss(tb)
	tb.G.DealDamageTo(engine.Target{Object: b}, 2, 0, 0)
	tb.G.Run()
	ev := settle(tb)
	if tb.G.Object(b).Damage != 0 || count(ev, engine.EvPrevented) != 1 {
		t.Errorf("guarded: took %d damage", tb.G.Object(b).Damage)
	}
	setTricks(tb, 3)
	tb.G.DealDamageTo(engine.Target{Object: b}, 2, 0, 0)
	tb.G.Run()
	settle(tb)
	if tb.G.Object(b).Damage != 2 {
		t.Errorf("3 tricks, 3 souls: took %d damage, want 2", tb.G.Object(b).Damage)
	}
	if opts := attackOptions(tb); !slices.Contains(opts, hauntNorm) {
		t.Errorf("3 tricks, 3 souls: attack options %v", opts)
	}
	tb.G.AddNamedCounters(b, challenge.TrickCounter, -10)
	if tricks(tb) != 0 {
		t.Errorf("trick counters %d, want 0", tricks(tb))
	}
	check(t, tb.G)
}

// R-TRICK-06, R-CHAL-05: Chaos Card kills outright; a guarded boss takes
// the armor's 3 damage as none; unguarded it takes 3. A Bomb! on the
// guarded boss is prevented too.
func TestTrickGuardOutrightKill(t *testing.T) {
	for _, n := range []int{1, 0} {
		tb := trick(t, engine.DifficultyNormal, engine.SituationPlayer{Items: []engine.CardRef{"chaos_card"}, Hand: []engine.CardRef{"bomb"}}, engine.SituationPlayer{})
		setTricks(tb, n)
		if _, err := tb.G.Apply(engine.Intent{Player: 0, Kind: engine.IntentPlayLoot, Objects: []engine.ObjectID{tb.Find(0, "bomb")}}); err != nil {
			t.Fatal(err)
		}
		settle(tb, hauntNorm)
		if _, err := tb.G.Apply(engine.Intent{Player: 0, Kind: engine.IntentActivate, Objects: []engine.ObjectID{tb.Find(0, "chaos_card")}}); err != nil {
			t.Fatal(err)
		}
		settle(tb, "Kill a player or monster.", hauntNorm)
		want := 0
		if n == 0 {
			want = 4 // the Bomb's 1, then the armor's 3
		}
		b, ok := tb.G.Bosses[0].TopOf()
		if ok && tb.G.Object(b).Damage != want || !ok && want < 4 {
			t.Errorf("%d tricks: damage %d, want %d", n, tb.G.Object(b).Damage, want)
		}
		if n == 0 && (!tb.G.Over || tb.G.Challenge.Lost) {
			t.Error("4 damage killed the boss; the challenge should be won")
		}
		check(t, tb.G)
	}
}

// R-TRICK-07, Q-50: each instance of damage a player takes puts one fear
// counter on them, whatever its amount, combat or not.
func TestTrickFear(t *testing.T) {
	tb := trick(t, engine.DifficultyNormal, engine.SituationPlayer{}, engine.SituationPlayer{})
	tb.G.DealDamageTo(engine.Target{Player: 1, IsPlayer: true}, 1, 0, 0)
	tb.G.Run()
	settle(tb)
	if fear(tb, 0) != 0 || fear(tb, 1) != 1 {
		t.Errorf("Normal, non-combat: fear %d and %d, want 0 and 1", fear(tb, 0), fear(tb, 1))
	}
	attack(tb, "fly", 6, 1, 6, 1) // a miss: Isaac takes 1 combat damage
	if fear(tb, 0) != 1 || fear(tb, 1) != 1 {
		t.Errorf("Normal, combat: fear %d and %d, want 1 and 1", fear(tb, 0), fear(tb, 1))
	}
	check(t, tb.G)
}

// R-TRICK-08, Q-45: fear counters stay on a player who dies.
func TestTrickFearStays(t *testing.T) {
	tb := trick(t, engine.DifficultyNormal, engine.SituationPlayer{}, engine.SituationPlayer{})
	tb.G.AddPlayerCounters(0, challenge.FearCounter, 3)
	tb.G.KillPlayer(0, engine.NoPlayer)
	settle(tb)
	if !tb.G.Players[0].Dead || fear(tb, 0) != 3 {
		t.Errorf("dead %v, fear %d; want 3", tb.G.Players[0].Dead, fear(tb, 0))
	}
}

// R-TRICK-09, Q-52: every 2 fear counters on the active player give
// every monster, the boss too, +1 evasion, up to 6; the other player's
// fear counts not.
func TestTrickFearEvasion(t *testing.T) {
	tb := trick(t, engine.DifficultyNormal, engine.SituationPlayer{}, engine.SituationPlayer{})
	g := tb.G
	fly := tb.G.Monsters[1].Cards[len(g.Monsters[1].Cards)-1]
	g.AddPlayerCounters(1, challenge.FearCounter, 6)
	if g.Evasion(fly) != 2 {
		t.Errorf("Cain's fear: fly DC %d, want 2", g.Evasion(fly))
	}
	g.AddPlayerCounters(0, challenge.FearCounter, 5)
	if g.Evasion(fly) != 4 || g.Evasion(boss(tb)) != 6 {
		t.Errorf("5 fear: fly DC %d, boss DC %d; want 4 and the cap 6", g.Evasion(fly), g.Evasion(boss(tb)))
	}
}

// R-TRICK-10, Q-51: each monster death, the active player rolls; 4-6:
// they pick a player with a fear counter, themselves included, who
// loses one; 1-3 nothing; with no fear anywhere, no question.
func TestTrickDeathRoll(t *testing.T) {
	tb := trick(t, engine.DifficultyNormal, engine.SituationPlayer{}, engine.SituationPlayer{})
	tb.G.AddPlayerCounters(0, challenge.FearCounter, 1)
	tb.G.AddPlayerCounters(1, challenge.FearCounter, 2)
	attackAsk(tb, "fly", 6, 6, 4)
	if opts := tb.Options(); !slices.Equal(opts, []string{player(tb, 0), player(tb, 1)}) || tb.G.Prompt().Player != 0 {
		t.Fatalf("question to %d, options %v", tb.G.Prompt().Player, opts)
	}
	settle(tb, player(tb, 1))
	if fear(tb, 0) != 1 || fear(tb, 1) != 1 {
		t.Errorf("4: fear %d and %d, want 1 and 1", fear(tb, 0), fear(tb, 1))
	}
	tb.G.Turn.Attacks++
	attack(tb, "leech", 6, 6, 3)
	if fear(tb, 0) != 1 || fear(tb, 1) != 1 {
		t.Errorf("3: fear %d and %d, want unchanged", fear(tb, 0), fear(tb, 1))
	}
	tb = trick(t, engine.DifficultyNormal, engine.SituationPlayer{}, engine.SituationPlayer{})
	attack(tb, "fly", 6, 6, 6) // no fear: no question
}
