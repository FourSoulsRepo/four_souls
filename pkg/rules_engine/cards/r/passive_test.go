package r

import (
	"testing"

	"github.com/FourSoulsRepo/rules_engine/enginetest"

	engine "github.com/FourSoulsRepo/rules_engine"
)

// Tests of the simple passive treasures; the others have their own files.

func TestAbel(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{seat("isaac", "abel"), seat("cain"), seat("samson")},
		Monsters: []engine.CardRef{"pooter", "fly"},
	})
	tb.G.ForceRolls(6, 5, 6, 1) // hit; Abel's 5; hit again; Abel's 1
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "pooter", "fly")
	if top(tb, 1) == "fly" {
		t.Error("Abel's 4-6 deals 1 damage to another monster: the fly dies")
	}
}

func TestAlabasterBox(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{seat("cain"), seat("isaac", "alabaster_box")}})
	box := tb.Find(1, "alabaster_box")
	tb.G.AddCountersTo(box, 1)
	tb.EndTurn() // Isaac's turn: the 2nd counter
	if got := tb.G.Players[1].Cents; got != 10 {
		t.Errorf("cents = %d, want 10 at the 2nd counter", got)
	}
	tb.G.AddCountersTo(box, 1)
	tb.EndTurn()
	tb.EndTurn() // the 4th counter
	if got := len(tb.G.Players[1].Hand); got < 4 {
		t.Errorf("hand = %d, want 4 more cards at the 4th counter", got)
	}
}

func TestBallOfBandages(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{seat("isaac", "ball_of_bandages"), seat("cain")}, Monsters: []engine.CardRef{"pooter"}})
	tb.Attack("pooter", 3, 6, 6)
	if tb.G.Players[0].Damage != 0 {
		t.Error("no combat damage on an attack roll of 3")
	}
}

func TestHaloOfFlies(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{seat("isaac", "halo_of_flies"), seat("cain")}, Monsters: []engine.CardRef{"leech"}})
	tb.Attack("leech", 1, 6)
	if tb.G.Players[0].Damage != 0 {
		t.Error("no combat damage on an attack roll of 1")
	}
}

func TestFireMind(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{seat("isaac", "fire_mind"), seat("cain")}, Monsters: []engine.CardRef{"pooter"}})
	tb.Attack("pooter", 6) // 1 ATK doubled kills the 2-HP pooter
	if top(tb, 0) == "pooter" {
		t.Error("damage doubled on a 6")
	}
	tb2 := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{seat("isaac", "fire_mind", "spoon_bender"), seat("cain")}, Monsters: []engine.CardRef{"pooter"}})
	tb2.G.ForceRolls(1)
	tb2.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "pooter", "spoon_bender")
	if !tb2.G.Players[0].Dead {
		t.Error("damage taken doubled on a 1: 2 damage kills Isaac")
	}
}

func TestBloodPuppy(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "isaac", Items: []engine.CardRef{"blood_puppy"}, Hand: []engine.CardRef{"mega_bomb"}}, seat("cain"),
	}})
	puppy := tb.Find(0, "blood_puppy")
	tb.G.AddNamedCounters(puppy, level, 3)
	tb.Play(0, "mega_bomb", foe)
	if tb.G.PlayerATK(0) != 2 {
		t.Errorf("ATK = %d, want 2 at level 4", tb.G.PlayerATK(0))
	}
}

func TestChocolateMilk(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{seat("cain"), seat("isaac", "chocolate_milk")}})
	milk := tb.Find(1, "chocolate_milk")
	tb.G.AddCountersTo(milk, 1)
	tb.EndTurn()
	if tb.G.PlayerATK(1) != 2 {
		t.Errorf("2 counters: ATK %d, want 2", tb.G.PlayerATK(1))
	}
	tb.EndTurn()
	tb.EndTurn("player 1 (cain)") // the 3rd counter: all removed, 1 damage
	if got := tb.G.Object(milk).CountersOf(""); got != 0 {
		t.Errorf("counters = %d, want 0", got)
	}
}

func TestConsolationPrize(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "cain", Souls: []engine.CardRef{"monstro", "gurdy"}}, seat("isaac", "consolation_prize"),
	}})
	tb.EndTurn()
	if got := len(tb.G.Players[1].Hand); got != 3 {
		t.Errorf("hand = %d, want 2 behind + 1 loot step", got)
	}
}

func TestDamocles(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{{Character: "isaac", Items: []engine.CardRef{"damocles"}, Cents: 10}, seat("cain")}})
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentPurchase}, "treasure deck")
	if got := len(tb.G.Players[0].InPlay); got != 3 {
		t.Errorf("items = %d, want Damocles + 2 treasures", got)
	}
}

func TestDamoclesDestroyedOnDeath(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "cain", Hand: []engine.CardRef{"mega_bomb"}}, seat("isaac", "damocles"),
	}})
	tb.Play(0, "mega_bomb", "player 2 (isaac)")
	if hasItem(tb.G, 1, "damocles") {
		t.Error("when he dies it is destroyed")
	}
}

func TestFetalHaunt(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "isaac", Items: []engine.CardRef{"fetal_haunt"}, Hand: []engine.CardRef{"bomb"}}, seat("cain"),
	}})
	tb.EndTurn("bomb", foe)
	if !handHas(tb.G, 1, "bomb") {
		t.Error("Isaac gives Cain the bomb")
	}
}

func TestLunch(t *testing.T) {
	tb := table2(t, "lunch")
	if tb.G.PlayerHP(0) != 3 {
		t.Error("+1 HP")
	}
}

func TestLodestone(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{seat("isaac", "lodestone"), seat("cain")}, Monsters: []engine.CardRef{"pooter"}})
	tb.G.ForceRolls(3)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "pooter", "pooter")
	// The 3 missed (it resolved before Lodestone); the next roll needs 2+.
	if dc := tb.G.Evasion(tb.G.Monsters[0].Cards[len(tb.G.Monsters[0].Cards)-1]); dc != 2 && top(tb, 0) == "pooter" {
		t.Errorf("the pooter has DC %d, want 3 - 1", dc)
	}
}

func TestTrickPenny(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players: []engine.SituationPlayer{{Character: "isaac", Items: []engine.CardRef{"trick_penny"}, Cents: 9}, seat("cain")},
		Shop:    []engine.CardRef{"spoon_bender"},
	})
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentPurchase}, "spoon_bender")
	if !hasItem(tb.G, 0, "spoon_bender") || tb.G.Players[0].Cents != 1 {
		t.Errorf("bought for 9¢, then +1¢: cents %d", tb.G.Players[0].Cents)
	}
	tb.EndTurn()
	if tb.G.Players[0].Cents != 2 {
		t.Error("+1¢ at the end of his turn")
	}
}

func TestTrickPennyPaid(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "isaac", Items: []engine.CardRef{"trick_penny", "the_jar"}, Cents: 3}, seat("cain"),
	}})
	tb.Activate(0, "the_jar", 0, foe)
	if tb.G.Players[0].Cents != 0 {
		t.Errorf("The Jar's 4¢ costs 3¢: cents %d", tb.G.Players[0].Cents)
	}
}

// table2 seats Isaac with the items and Cain.
func table2(t *testing.T, items ...engine.CardRef) *enginetest.Table {
	t.Helper()
	return setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{seat("isaac", items...), seat("cain")}})
}
