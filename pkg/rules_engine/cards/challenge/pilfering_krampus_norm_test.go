package challenge_test

import (
	"errors"
	"slices"
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

const normBoss = "pilfering_krampus_norm"

// R-KRAMPUS-11, Q-30: below 5 souls between both players, or with a
// minion in play, Pilfering Krampus can't be attacked and takes no
// damage. Character soul values count (The Lost).
func TestKrampusGuarded(t *testing.T) {
	tb := krampus(t, engine.DifficultyNormal, engine.SituationPlayer{Souls: souls(3)}, engine.SituationPlayer{Souls: souls(1)})
	if opts := attackOptions(tb); slices.Contains(opts, normBoss) || !slices.Contains(opts, "leech") {
		t.Errorf("4 souls: attack options %v", opts)
	}
	tb = krampus(t, engine.DifficultyNormal, engine.SituationPlayer{Souls: souls(4)}, engine.SituationPlayer{})
	b := boss(tb)
	tb.G.DealDamageTo(engine.Target{Object: b}, 2, 0, 0)
	tb.G.Run()
	ev := tb.Settle()
	if tb.G.Object(b).Damage != 0 || count(ev, engine.EvPrevented) != 1 {
		t.Errorf("4 souls: Krampus took %d damage", tb.G.Object(b).Damage)
	}

	// 4 souls and The Lost's soul value: 5 (Q-30).
	tb = krampus(t, engine.DifficultyNormal, engine.SituationPlayer{Souls: souls(4)}, engine.SituationPlayer{Character: "the_lost"})
	if got := tb.G.SoulsInPlay(); got != 5 {
		t.Fatalf("souls in play %d, want 5", got)
	}
	if opts := attackOptions(tb); !slices.Contains(opts, normBoss) {
		t.Errorf("5 souls: attack options %v", opts)
	}

	tb = krampus(t, engine.DifficultyNormal, engine.SituationPlayer{Souls: souls(5)}, engine.SituationPlayer{}, "pooter")
	tb.G.Turn.RandomAttackMade = true // the Normal duty aside
	if opts := attackOptions(tb); slices.Contains(opts, normBoss) {
		t.Errorf("5 souls and a minion: attack options %v", opts)
	}
	check(t, tb.G)
}

// R-KRAMPUS-11, R-CHAL-05: killed outright by an ability (Chaos Card),
// the boss takes 3 damage instead; guarded, it takes none.
func TestKrampusChaosCard(t *testing.T) {
	for _, n := range []int{4, 5} {
		tb := krampus(t, engine.DifficultyNormal, engine.SituationPlayer{Souls: souls(n), Items: []engine.CardRef{"chaos_card"}}, engine.SituationPlayer{})
		tb.Activate(0, "chaos_card", 0, "Kill a player or monster.", normBoss)
		want := 0
		if n == 5 {
			want = 3
		}
		if got := tb.G.Object(boss(tb)).Damage; got != want {
			t.Errorf("%d souls: Krampus has %d damage, want %d", n, got, want)
		}
		check(t, tb.G)
	}
}

// R-KRAMPUS-12 (Normal): as Cain's turn starts, he rolls once; the
// table removes counters from a present.
func TestKrampusRollNormal(t *testing.T) {
	for roll, want := range map[int][4]int{
		1: {4, 2, 3, 3}, 2: {4, 2, 3, 3}, 3: {8, 1, 3, 3}, 4: {8, 1, 3, 3}, 5: {8, 2, 2, 3}, 6: {8, 2, 3, 2},
	} {
		tb := krampus(t, engine.DifficultyNormal, engine.SituationPlayer{}, engine.SituationPlayer{})
		tb.G.ForceRolls(roll, 1, 1)
		tb.EndTurn()
		if got := counters(tb.G); got != want {
			t.Errorf("roll %d: counters %v, want %v", roll, got, want)
		}
		check(t, tb.G)
	}
}

// R-KRAMPUS-13, Q-31: with a minion in play, the active player can't
// end the turn before attacking; the attack goes to a random minion,
// with no target and no deck to pick. Then the turn may end.
func TestKrampusForcedAttack(t *testing.T) {
	tb := krampus(t, engine.DifficultyNormal, engine.SituationPlayer{}, engine.SituationPlayer{}, "pooter", "dip")
	_, err := tb.G.Apply(engine.Intent{Player: 0, Kind: engine.IntentEndTurn})
	var re *engine.RuleError
	if !errors.As(err, &re) || re.Msg().Key.String() != "refuse.must_attack" {
		t.Fatalf("end turn with a minion: %v", err)
	}
	boost(tb, 0, 4)
	tb.G.ForceRolls(6)
	ev := tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack})
	if count(ev, engine.EvRandomPick) != 1 || len(tb.G.Minions) != 1 {
		t.Errorf("the attack went to no random minion; minions %v", minionCards(tb.G))
	}
	if _, err := tb.G.Apply(engine.Intent{Player: 0, Kind: engine.IntentEndTurn}); err != nil {
		t.Errorf("the duty is met: %v", err)
	}
	check(t, tb.G)
}

// Q-31: a minion that comes later brings the duty; without an attack
// left there is none. Q-31, Q-36: with another forced attack the
// attacker picks which comes first; the random one is an option.
func TestKrampusForcedAttackLater(t *testing.T) {
	tb := krampus(t, engine.DifficultyNormal, engine.SituationPlayer{}, engine.SituationPlayer{})
	if blocked(tb) {
		t.Fatal("no minion: the turn may end")
	}
	tb.G.PutOnDeck(engine.MonsterDeck, "pooter")
	tb.G.MinionFromDeck()
	if !blocked(tb) {
		t.Error("a minion came: the turn can't end")
	}
	tb.G.Turn.Attacks = 0
	if blocked(tb) {
		t.Error("no attack left: the turn may end")
	}

	tb = krampus(t, engine.DifficultyNormal, engine.SituationPlayer{}, engine.SituationPlayer{}, "pooter")
	leech := tb.G.MonstersInPlay()[0]
	tb.G.Turn.MustAttack = leech // Monster Manual, the regular Krampus
	opts := attackOptions(tb)
	if !slices.Equal(opts, []string{"leech", "ask.attack_random"}) {
		t.Fatalf("options %v", opts)
	}
	tb.Choose("ask.attack_random")
	if tgt := tb.G.Attack.Target; !tb.G.IsMinion(tgt) || tb.G.Turn.MustAttack != leech {
		t.Errorf("the random attack goes to a minion, leech still owed: target %d", tgt)
	}
}

func blocked(tb *enginetest.Table) bool {
	c, err := tb.G.Clone()
	if err != nil {
		tb.T.Fatal(err)
	}
	_, err = c.Apply(engine.Intent{Player: c.Turn.Active, Kind: engine.IntentEndTurn})
	return err != nil
}

// R-KRAMPUS-14: a dying minion gives no rewards and no soul; its card
// goes to discard. A monster in a monster slot gives them as usual.
func TestKrampusMinionSpoils(t *testing.T) {
	tb := krampus(t, engine.DifficultyNormal, engine.SituationPlayer{Items: []engine.CardRef{"chaos_card"}}, engine.SituationPlayer{}, "gurdy")
	hand := len(tb.G.Players[0].Hand)
	tb.Activate(0, "chaos_card", 0, "Kill a player or monster.", "gurdy")
	if len(tb.G.Minions) != 0 || len(tb.G.Players[0].Hand) != hand || tb.G.SoulValue(0) != 0 {
		t.Errorf("minion gurdy: minions %v, hand %d, souls %d", minionCards(tb.G), len(tb.G.Players[0].Hand), tb.G.SoulValue(0))
	}
	d := tb.G.Discards[engine.MonsterDeck]
	if len(d) == 0 || tb.G.Object(d[len(d)-1]).Card != "gurdy" {
		t.Error("the minion's soul card goes to discard")
	}
	boost(tb, 0, 4)
	tb.G.PutOnDeck(engine.MonsterDeck, "fly") // the refill
	tb.G.ForceRolls(6)                        // a monster slot: its loot reward, not forgone
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "leech", "ask.no")
	if len(tb.G.Players[0].Hand) != hand+1 {
		t.Errorf("leech's reward: hand %d", len(tb.G.Players[0].Hand))
	}
	check(t, tb.G)
}
