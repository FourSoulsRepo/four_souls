package r

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

func TestBobbyBomb(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{{Character: "isaac", Cents: 10}, seat("cain")}})
	tb.G.PutOnDeck(engine.LootDeck, "mega_bomb", "a_nickel", "bomb")
	tb.G.PutOnDeck(engine.TreasureDeck, "bobby_bomb")
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentPurchase}, "treasure deck")
	if !handHas(tb.G, 0, "bomb") || !handHas(tb.G, 0, "mega_bomb") || handHas(tb.G, 0, "a_nickel") {
		t.Error("the bombs of the top 10 go to his hand; the nickel to the bottom")
	}
}

func TestBobbyBombStartOfTurn(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{seat("cain"), seat("isaac", "bobby_bomb")}})
	tb.G.PutOnDeck(engine.LootDeck, "troll_bombs")
	tb.G.MillTop(engine.LootDeck)
	tb.EndTurn("troll_bombs")
	if !handHas(tb.G, 1, "troll_bombs") {
		t.Error("a bomb card comes back from the loot discard")
	}
}

func TestMomsBottleOfPills(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{{Character: "isaac", Cents: 10}, seat("cain")}})
	tb.G.PutOnDeck(engine.LootDeck, "pills", "bomb", "pills_r")
	tb.G.PutOnDeck(engine.TreasureDeck, "moms_bottle_of_pills")
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentPurchase}, "treasure deck")
	if !handHas(tb.G, 0, "pills") || !handHas(tb.G, 0, "pills_r") || handHas(tb.G, 0, "bomb") {
		t.Error("each Pills! goes to his hand")
	}
}

func TestGnawedLeafNotYourTurn(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "cain", Hand: []engine.CardRef{"mega_bomb"}}, seat("isaac", "gnawed_leaf"), seat("samson"),
	}})
	tb.Play(0, "mega_bomb", "player 2 (isaac)", "no")
	if tb.G.Players[1].Damage != 0 {
		t.Error("damage is prevented while it's not his turn")
	}
}

func TestGnawedLeafCancels(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "cain", Hand: []engine.CardRef{"two_of_spades", "two_of_spades"}},
		{Character: "isaac", Items: []engine.CardRef{"gnawed_leaf"}, Hand: []engine.CardRef{"bomb"}},
		seat("samson"),
	}})
	tb.G.Turn.LootPlays = 2
	tb.Play(0, "two_of_spades", "player 2 (isaac)", "yes")
	if len(tb.G.Players[0].Hand) != 1 {
		t.Fatal("the first loot aimed at him is cancelled")
	}
	tb.Play(0, "two_of_spades", "player 2 (isaac)")
	if len(tb.G.Players[0].Hand) != 1 {
		t.Error("the second resolves: Cain loots 1 for Isaac's 1 card")
	}
}

func TestHandicappedPlacard(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{seat("isaac", "handicapped_placard"), {Character: "cain", Souls: []engine.CardRef{"monstro", "gurdy"}}},
		Monsters: []engine.CardRef{"leech"},
	})
	if dc := tb.G.Evasion(tb.G.Monsters[0].Cards[0]); dc != 2 {
		t.Errorf("DC = %d, want 4 - 2 on his turn", dc)
	}
}

func TestJudasShadow(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "cain", Hand: []engine.CardRef{"mega_bomb"}}, seat("isaac", "judas_shadow"),
	}, Shop: []engine.CardRef{"spoon_bender"}})
	tb.Play(0, "mega_bomb", "player 2 (isaac)", "judas_shadow")
	if got := tb.G.Players[1].Cents; got != 2 {
		t.Errorf("2 counters, then destroyed as the penalty: cents %d, want 2", got)
	}
	if got := len(tb.G.Players[1].Hand); got != 2 {
		t.Errorf("hand = %d, want 2", got)
	}
}

func TestKeepersBox(t *testing.T) {
	tb := setup(t, engine.SituationSetup{
		Players: []engine.SituationPlayer{{Character: "cain", Cents: 10}, seat("isaac", "keepers_box")},
		Shop:    []engine.CardRef{"spoon_bender"},
	})
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentPurchase}, "spoon_bender")
	if tb.G.Players[1].Cents != 2 || len(tb.G.Players[1].Hand) != 1 {
		t.Error("another player bought a shop item: 2¢ and loot 1")
	}
	tb2 := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{{Character: "isaac", Cents: 10}, seat("cain")}})
	tb2.G.PutOnDeck(engine.TreasureDeck, "keepers_box")
	shop := len(tb2.G.Shop)
	tb2.Do(engine.Intent{Player: 0, Kind: engine.IntentPurchase}, "treasure deck")
	if len(tb2.G.Shop) != shop+2 {
		t.Error("it entered play: 2 more shop slots")
	}
}

func TestLilChest(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{seat("isaac", "lil_chest"), seat("cain")}, Monsters: []engine.CardRef{"leech"}})
	tb.Attack("leech", 6, 2) // the kill; Lil Chest's 2
	if got := tb.G.Players[0].Cents; got != 2 {
		t.Errorf("cents = %d, want 2", got)
	}
}

func TestPlaydoughCookie(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{seat("cain"), seat("isaac", "playdough_cookie", "spoon_bender")}})
	tb.EndTurn("none")
	if tb.G.Players[1].Cents != 1 {
		t.Error("no item chosen: gain 1¢")
	}
	tb.EndTurn()
	tb.EndTurn("spoon_bender", "add a counter")
	if tb.G.Object(tb.Find(1, "spoon_bender")).CountersOf("") != 1 {
		t.Error("a counter on the Spoon Bender")
	}
}

func TestRockBottom(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "isaac", Cents: 10, Hand: []engine.CardRef{"bomb", "bomb"}}, seat("cain"),
	}})
	tb.G.PutOnDeck(engine.TreasureDeck, "rock_bottom")
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentPurchase}, "treasure deck")
	if got := tb.G.Object(tb.Find(0, "rock_bottom")).CountersOf(""); got != 2 {
		t.Fatalf("counters = %d, want his 2 cards", got)
	}
	tb.G.Players[0].Hand = tb.G.Players[0].Hand[:1]
	tb.EndTurn()
	if got := len(tb.G.Players[0].Hand); got != 3 {
		t.Errorf("hand = %d, want 1 + 2: fewer cards than counters", got)
	}
}

func TestSpelunkerHat(t *testing.T) {
	tb := setup(t, engine.SituationSetup{Players: []engine.SituationPlayer{seat("cain"), seat("isaac", "spelunker_hat")}})
	tb.G.PutOnDeck(engine.TreasureDeck, "spoon_bender", "the_battery", "boomerang")
	tb.EndTurn("treasure deck", "spoon_bender", "the_battery", "boomerang")
	if top := tb.G.DeckTop(engine.TreasureDeck, 1); tb.G.Object(top[0]).Card != "spoon_bender" {
		t.Error("the Spoon Bender is on top now")
	}
}
