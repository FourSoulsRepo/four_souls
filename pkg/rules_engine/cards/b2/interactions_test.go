package b2

// Interaction tests: card combinations and rulings from step 3 (TS-04).
// Each test names its source: the re-checked old rulings (R-07) and the
// S-RU FAQ rulings in docs/rules/open-questions.md, and Yuggy's rulings.

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

// R-07 #1: Monster Manual forces the attack target but gives no attack.
func TestRulingMonsterManualGivesNoAttack(t *testing.T) {
	tb := itemTable(t, items("monster_manual"))
	tb.Activate(0, "monster_manual", 0, "leech")
	if n := tb.G.Turn.Attacks; n != 1 {
		t.Errorf("attacks = %d, want 1", n)
	}
}

// R-07 #2: a player dies at most once per turn (R-DEATH-17).
func TestRulingDeadPlayerDiesOnce(t *testing.T) {
	tb := enginetest.NewSetup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "isaac", Hand: items("xiii_death", "xiii_death")}, {Character: "cain", Cents: 3},
	}}, Set)
	tb.G.Players[0].ExtraLootPlays = 1
	tb.Play(0, "xiii_death", foe)
	tb.Play(0, "xiii_death", foe)
	if c := tb.G.Players[1].Cents; c != 2 {
		t.Errorf("Cain's cents = %d, want 2: one penalty only", c)
	}
}

// R-07 #3: a copy takes abilities, not counters (Modeling Clay, Tech X).
func TestRulingCopyTakesNoCounters(t *testing.T) {
	tb := itemTable(t, items("modeling_clay", "tech_x"))
	tb.G.Object(tb.Find(0, "tech_x")).Counters = []engine.Counter{{Count: 3}}
	tb.Activate(0, "modeling_clay", 0, "tech_x")
	if rule := tb.Refused(0, "modeling_clay", 1); rule != "R-ABIL-07" {
		t.Errorf("the copy used Tech X's counters: refused by %q", rule)
	}
}

// R-07 #4: a dead player can be damaged, but nothing is marked (R-MECH-16).
func TestRulingDamageToDeadPlayer(t *testing.T) {
	tb := lootTable(t, "xiii_death", "bomb")
	tb.G.Players[0].ExtraLootPlays = 1
	tb.Play(0, "xiii_death", foe)
	before := tb.G.Players[1].Damage
	tb.Play(0, "bomb", foe)
	if d := tb.G.Players[1].Damage; d != before || !tb.G.Players[1].Dead {
		t.Errorf("damage %d, want %d", d, before)
	}
}

// R-07 #5: Eden's eternal Glass Cannon is not destroyed by its own roll.
func TestRulingEdenGlassCannon(t *testing.T) {
	tb := enginetest.NewSetup(t, engine.SituationSetup{Players: []engine.SituationPlayer{seat("eden", "glass_cannon"), seat("cain", "breakfast")}}, Set)
	tb.G.Object(tb.Find(0, "glass_cannon")).Eternal = true
	tb.G.ForceRolls(1)
	tb.Activate(0, "glass_cannon", 0, "breakfast")
	if !hasItem(tb.G, 0, "glass_cannon") || len(tb.G.Players[0].Hand) != 2 {
		t.Error("the eternal cannon should stay, and Eden loots 2")
	}
}

// R-07 #6: Steamy Sale lowers shop items, not the top of the deck.
func TestRulingSteamySaleOnlyShop(t *testing.T) {
	tb := enginetest.NewSetup(t, engine.SituationSetup{Players: []engine.SituationPlayer{{Character: "isaac", Items: items("steamy_sale"), Cents: 5}, seat("cain")}}, Set)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentPurchase}, "treasure deck")
	if c := tb.G.Players[0].Cents; c != 5 {
		t.Errorf("cents = %d: the deck top cost less than 10¢", c)
	}
}

// R-07 #10: death cancels combat (R-DEATH-12, R-ATK-15).
func TestRulingDeathCancelsCombat(t *testing.T) {
	tb := monsterTable(t, items("conjoined_fatty"), seat("isaac"), seat("cain"))
	tb.Attack("conjoined_fatty", 1)
	if tb.G.Attack.On || !tb.G.Players[0].Dead {
		t.Error("the attack goes on after Isaac died")
	}
}

// R-07 #11: a bomb's damage is not combat damage (Lust ignores it).
func TestRulingBombsAreNotCombat(t *testing.T) {
	tb := enginetest.NewSetup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{{Character: "isaac", Hand: items("bomb")}, seat("cain")},
		Monsters: items("lust"),
	}, Set)
	tb.Play(0, "bomb", "lust")
	if d := tb.G.Players[0].Damage; d != 0 {
		t.Errorf("damage = %d, want 0: Lust reacts only to combat damage", d)
	}
}

// S-RU FAQ #1: Two of Clubs twice still only doubles.
func TestRulingTwoOfClubsTwice(t *testing.T) {
	tb := itemTable(t, items("two_of_clubs", "bum_friend"), "a_dime")
	tb.Activate(0, "two_of_clubs", 0, me)
	tb.G.Recharge(tb.Find(0, "two_of_clubs"))
	tb.Activate(0, "two_of_clubs", 0, me)
	tb.Activate(0, "bum_friend", 0, "a_dime") // loot 1, doubled once; one goes back
	if h := len(tb.G.Players[0].Hand); h != 2 {
		t.Errorf("hand = %d, want 2", h)
	}
}

// Yuggy (X, 2024-08-09): a roll pushed above 6 counts as 6; minus 1 is 5.
func TestRulingRollAboveSix(t *testing.T) {
	tb := rollTable(t, 6, seat("judas", "spoon_bender", "book_of_belial"), seat("cain"))
	tb.Pass(1)
	tb.Start(0, "spoon_bender", 0)
	tb.Choose("roll of 6")
	tb.Pass(0)
	tb.Pass(1) // the +1 resolves: still 6
	tb.Activate(0, "book_of_belial", 0, "Subtract 1 from a roll.", "roll of 6")
	if r, ok := firstRoll(tb.G); ok && r != 5 {
		t.Errorf("roll = %d, want 5", r)
	}
}

func firstRoll(g *engine.Game) (int, bool) {
	for _, it := range g.Stack {
		if it.Kind == engine.StackRoll {
			return it.Roll, true
		}
	}
	return 0, false
}
