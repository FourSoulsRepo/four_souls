package r

import (
	"testing"

	"github.com/FourSoulsRepo/rules_engine/enginetest"

	engine "github.com/FourSoulsRepo/rules_engine"
)

// Attacks on players, attackable items and helpers (R-ATK-20, R-ATK-21,
// R-MECH-04), through the cards that allow them.

func sirenTable(t *testing.T, isaacItems, cainItems []engine.CardRef) *enginetest.Table {
	t.Helper()
	return setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{seat("isaac", isaacItems...), seat("cain", cainItems...)},
		Monsters: []engine.CardRef{"the_siren", "fly"},
	})
}

func TestTheSirenAttackAPlayer(t *testing.T) {
	tb := sirenTable(t, nil, nil)
	tb.G.ForceRolls(4, 4)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, foe)
	if !tb.G.Players[1].Dead {
		t.Fatal("two hits on DC 4 kill Cain (2 HP)")
	}
	if tb.G.Attack.On {
		t.Error("killing the attacked player ends the attack")
	}
}

func TestTheSirenMissHurtsTheAttacker(t *testing.T) {
	tb := sirenTable(t, nil, nil)
	tb.G.ForceRolls(3, 6, 6)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, foe)
	if tb.G.Players[0].Damage != 1 {
		t.Errorf("a miss: Cain deals his 1 ATK to Isaac; damage %d", tb.G.Players[0].Damage)
	}
}

func TestBackstabber(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		seat("isaac", "backstabber"), {Character: "cain", Souls: []engine.CardRef{"monstro"}}, seat("samson"),
	}})
	tb.G.ForceRolls(3, 3)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, foe)
	if !tb.G.Players[1].Dead || len(tb.G.Players[0].Hand) != 1 {
		t.Error("Cain (more souls) can be attacked at 3+; killing him loots 1")
	}
	for _, o := range tb.Options() {
		if o == "player 3 (samson)" {
			t.Error("Samson has no more souls: not a target")
		}
	}
}

func TestHaloOfFliesWhenAttacked(t *testing.T) {
	tb := sirenTable(t, nil, []engine.CardRef{"halo_of_flies"})
	tb.G.ForceRolls(6, 4, 4)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, foe, "halo_of_flies")
	if tb.G.Players[1].Damage != 0 && !tb.G.Players[1].Dead {
		t.Error("a 6 deals no damage to Cain (R-ATK-24)")
	}
	if !tb.G.Players[1].Dead {
		t.Error("the 4s still hit")
	}
}

func TestPunchingBag(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{seat("isaac"), seat("cain", "punching_bag")}})
	if tb.G.PlayerHP(1) != 4 {
		t.Errorf("Cain's HP = %d, want 2 + 2", tb.G.PlayerHP(1))
	}
	tb.G.ForceRolls(5, 5)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "punching_bag")
	if hasItem(tb.G, 1, "punching_bag") {
		t.Fatal("killed, it is destroyed (R-MECH-05)")
	}
	if tb.G.Players[0].Cents != 5 {
		t.Errorf("the active player gains its 5¢; cents %d", tb.G.Players[0].Cents)
	}
}

func TestUltraGreedTeamUp(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{seat("isaac"), seat("cain"), seat("samson")},
		Monsters: []engine.CardRef{"ultra_greed"},
	})
	tb.G.Object(tb.G.Monsters[0].Cards[0]).Damage = 7 // 1 HP left
	// Isaac misses (1 damage to him), Cain misses, Samson hits and kills.
	tb.G.ForceRolls(1, 1, 5)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "ultra_greed")
	if tb.G.Players[2].Cents != 40 || tb.G.SoulValue(2) != 1 {
		t.Errorf("Samson killed it: he gains the 40¢ and the soul; cents %d, souls %d", tb.G.Players[2].Cents, tb.G.SoulValue(2))
	}
	if tb.G.Players[0].Damage != 1 || tb.G.Players[1].Damage != 1 {
		t.Error("each miss hurts the player who rolled it")
	}
}

func TestSistersVis(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{seat("isaac"), seat("cain"), seat("samson")},
		Monsters: []engine.CardRef{"sisters_vis"},
	})
	tb.G.Object(tb.G.Monsters[0].Cards[0]).Damage = 4 // 2 HP left
	tb.G.ForceRolls(6, 6)                             // Isaac hits, then Samson hits and kills
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "sisters_vis", "player 3 (samson)")
	if tb.G.SoulValue(0) != 1 {
		t.Error("the active player gains the soul")
	}
	if len(tb.G.Players[2].InPlay) != 1 {
		t.Error("Samson killed it with Isaac: he gains +1 treasure")
	}
}

func TestSiblingRivalry(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{seat("isaac"), seat("cain", "sibling_rivalry")},
		Monsters: []engine.CardRef{"pooter"},
	})
	tb.G.ForceRolls(6, 6) // Isaac hits, Cain hits: the pooter (2 HP) dies
	if _, err := tb.G.Apply(engine.Intent{Player: 0, Kind: engine.IntentAttack}); err != nil {
		t.Fatal(err)
	}
	tb.Pass(0)
	tb.Activate(1, "sibling_rivalry", 0, foe, "pooter") // in answer to the declaration (S-YUG)
	if top(tb, 0) == "pooter" {
		t.Fatal("the pooter should be dead")
	}
	if len(tb.G.Players[0].Hand) != 1 || len(tb.G.Players[1].Hand) != 1 {
		t.Errorf("both gain its reward (loot 1): hands %d and %d", len(tb.G.Players[0].Hand), len(tb.G.Players[1].Hand))
	}
}

func TestMegaSatanIndomitable(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{seat("isaac"), seat("cain"), seat("samson")},
		Monsters: []engine.CardRef{"fly", "leech"},
	})
	tb.G.PutOnDeck(engine.MonsterDeck, "mega_satan")
	tb.G.ForceRolls(6)
	tb.Attack("fly", 6) // the fly dies; Mega Satan would refill its slot
	if len(tb.G.Monsters) != 3 {
		t.Fatalf("monster slots = %d, want 3: Mega Satan takes a new one", len(tb.G.Monsters))
	}
	if top(tb, 2) != "mega_satan" || top(tb, 0) == "fly" {
		t.Error("Mega Satan in the new slot; slot 1 refilled with another card")
	}
}

func TestMegaSatanWins(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{seat("isaac"), {Character: "cain", Dead: true}, {Character: "samson", Dead: true}},
		Monsters: []engine.CardRef{"mega_satan", "fly"},
	})
	if _, err := tb.G.Apply(engine.Intent{Player: 0, Kind: engine.IntentEndTurn}); err != nil {
		t.Fatal(err)
	}
	for range 20 {
		if w := tb.G.Prompt(); w.Kind == engine.PromptPriority && !tb.G.Over {
			tb.Pass(int(w.Player))
		}
	}
	if !tb.G.Over || len(tb.G.Winners) != 1 || tb.G.Winners[0] != 0 {
		t.Error("the only living player wins at the end of the turn")
	}
}

func TestLustForBlood(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{seat("isaac"), seat("cain", "spoon_bender")},
		Monsters: []engine.CardRef{"fly", "leech"},
	})
	tb.G.PutOnDeck(engine.MonsterDeck, "lust_for_blood")
	tb.G.ForceRolls(6)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "fly", foe) // Lust for Blood! fills the slot
	if _, err := tb.G.Apply(engine.Intent{Player: 0, Kind: engine.IntentEndTurn}); err == nil {
		t.Fatal("Isaac must attack Cain first")
	}
	tb.G.ForceRolls(4, 4)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, foe, "spoon_bender")
	if !tb.G.Players[1].Dead || !hasItem(tb.G, 0, "spoon_bender") {
		t.Error("Cain dies at 4+ and gives Isaac the item he would destroy")
	}
}
