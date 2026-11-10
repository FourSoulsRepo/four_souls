package cards

import (
	"bytes"
	"slices"
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

// slotsSet is a test challenge's set (14.21): a final boss and a rule
// card. Each turn's end puts the top monster card into a minion slot,
// so fuzzing meets minions and events in minion slots; the boss's
// death wins.
var slotsSet = engine.CardSet{Code: "slots_trial", Name: "slots_trial", Cards: []engine.CardDef{
	{
		Ref: "slots_boss", Kind: engine.MonsterCard, HP: 6, DC: 4, ATK: 1, Soul: 1, Outside: true,
		Rewards: []engine.Reward{{Kind: engine.RewardTreasure, Amount: 1}},
	},
	{Ref: "slots_rules", Kind: engine.ChallengeCard, Abilities: []engine.Ability{
		{
			Kind: engine.Triggered, Text: "At the end of each turn, put the top card of the monster deck into a minion slot.",
			Trigger: engine.Trigger{On: engine.EvEndOfTurn, Match: func(*engine.Game, engine.ObjectID, engine.Event) bool { return true }},
			Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) { c.G.MinionFromDeck() })},
		},
		{
			Kind: engine.Triggered, Text: "When the final boss dies, the players win.",
			Trigger: engine.Trigger{On: engine.EvDied, Match: func(_ *engine.Game, _ engine.ObjectID, e engine.Event) bool { return e.Card == "slots_boss" }},
			Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) { c.G.WinChallenge() })},
		},
	}},
}}

// slotsChallenge puts its boss into a final boss slot and the first two
// monsters of the monster deck into minion slots (R-CHAL-02, R-CHAL-09).
var slotsChallenge = engine.ChallengeDef{
	ID: "slots_trial", Set: slotsSet.Code,
	Levels: []engine.ChallengeLevel{{Difficulty: engine.DifficultyNormal, Boss: "slots_boss", Card: "slots_rules"}},
	Setup: func(g *engine.Game, _ string) error {
		g.PutInFinalBossSlot("slots_boss")
		var monsters []engine.ObjectID
		for i := len(g.Decks[engine.MonsterDeck]) - 1; i >= 0 && len(monsters) < 2; i-- {
			if id := g.Decks[engine.MonsterDeck][i]; g.Kind(id) == engine.MonsterCard {
				monsters = append(monsters, id)
			}
		}
		for _, id := range monsters {
			g.PutInMinionSlot(id)
		}
		return nil
	},
}

func slotsSetup(seed uint64, rooms bool) engine.Setup {
	return engine.Setup{
		Seed: seed, Players: 2, Sets: append(Sets(), slotsSet), Rooms: rooms,
		Challenge: &slotsChallenge, Difficulty: engine.DifficultyNormal,
	}
}

// FuzzChallengeSets is FuzzSets in a challenge with a final boss and
// minions (14.21). Flags 1 to 3 (bits 0 and 1) play Krampus at Normal,
// Hard or Ultra Hard (14.22); 0 the test challenge. Bit 2 adds rooms.
//
//	go test -run '^$' -fuzz=FuzzChallengeSets -fuzztime=5m ./cards/
func FuzzChallengeSets(f *testing.F) {
	f.Add(uint64(1), uint8(0), []byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9})
	f.Add(uint64(5), uint8(4), bytes.Repeat([]byte{1, 0, 3, 2}, 100))
	f.Add(uint64(9), uint8(0), bytes.Repeat([]byte{0}, 300))
	f.Add(uint64(3), uint8(4), bytes.Repeat([]byte{2, 1}, 300))
	f.Add(uint64(2), uint8(1), bytes.Repeat([]byte{1, 0, 3, 2}, 300))
	f.Add(uint64(4), uint8(6), bytes.Repeat([]byte{2, 1}, 400))
	f.Add(uint64(7), uint8(3), bytes.Repeat([]byte{0, 3, 1}, 400))
	f.Fuzz(func(t *testing.T, seed uint64, flags uint8, choices []byte) {
		rooms := flags&4 != 0
		if level := flags & 3; level > 0 {
			s := krampusSetup(seed, krampusLevels[level-1])
			s.Rooms = rooms
			playSets(t, s, choices)
			return
		}
		playSets(t, slotsSetup(seed, rooms), choices)
	})
}

// R-CHAL-02, R-CHAL-09: the challenge's setup puts its boss into a final
// boss slot and each minion into a slot of its own, apart from the
// monster slots; the view shows both rows, a regular game neither.
func TestChallengeSlotsSetup(t *testing.T) {
	g, _, err := engine.NewGame(slotsSetup(1, false))
	if err != nil {
		t.Fatal(err)
	}
	if len(g.Bosses) != 1 || len(g.Minions) != 2 || len(g.Monsters) != 2 {
		t.Fatalf("bosses %d, minions %d, monsters %d", len(g.Bosses), len(g.Minions), len(g.Monsters))
	}
	boss, _ := g.Bosses[0].TopOf()
	if g.Object(boss).Card != "slots_boss" || !g.IsFinalBoss(boss) || g.IsMinion(boss) {
		t.Errorf("boss %+v", g.Object(boss))
	}
	for i, s := range g.Minions {
		m, _ := s.TopOf()
		if !g.IsMinion(m) || g.Object(m).Zone.Index != i {
			t.Errorf("minion %d: %+v", i, g.Object(m))
		}
	}
	checkInvariants(t, g)
	v := g.View(engine.Viewer{Kind: engine.ViewSpectator})
	if len(v.Bosses) != 1 || v.Bosses[0].Top.Card != "slots_boss" || len(v.Minions) != 2 || len(v.Monsters) != 2 {
		t.Errorf("view: bosses %+v, minions %+v", v.Bosses, v.Minions)
	}
	r, _, err := engine.NewGame(engine.Setup{Seed: 1, Players: 2, Sets: Sets()})
	if err != nil {
		t.Fatal(err)
	}
	if v := r.View(engine.Viewer{Kind: engine.ViewSpectator}); v.Bosses != nil || v.Minions != nil {
		t.Error("a regular game shows challenge slots")
	}
	data, err := g.Save()
	if err != nil {
		t.Fatal(err)
	}
	l, err := engine.Load(data, append(Sets(), slotsSet)...)
	if err != nil || len(l.Bosses) != 1 || len(l.Minions) != 2 {
		t.Fatalf("load: %v", err)
	}
}

// slotsTable: Isaac (active) and Cain, the monsters leech and fly, the
// final boss Gurdy (HP 5) and the given minions.
func slotsTable(t *testing.T, isaac engine.SituationPlayer, minions ...engine.CardRef) *enginetest.Table {
	t.Helper()
	return slotsTable2(t, isaac, engine.SituationPlayer{}, minions...)
}

// slotsTable2 is slotsTable with Cain's cards too.
func slotsTable2(t *testing.T, isaac, cain engine.SituationPlayer, minions ...engine.CardRef) *enginetest.Table {
	t.Helper()
	isaac.Character, cain.Character = "isaac", "cain"
	return enginetest.NewSetup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{isaac, cain},
		Monsters: []engine.CardRef{"leech", "fly"},
		Bosses:   []engine.CardRef{"gurdy"},
		Minions:  minions,
	}, Sets()...)
}

func bossOf(tb *enginetest.Table) engine.ObjectID {
	tb.T.Helper()
	id, ok := tb.G.Bosses[0].TopOf()
	if !ok {
		tb.T.Fatal("the final boss slot is empty")
	}
	return id
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

func checkInvariants(t *testing.T, g *engine.Game) {
	t.Helper()
	if err := g.CheckInvariants(); err != nil {
		t.Fatal(err)
	}
}

// R-CHAL-03, R-CHAL-07: the final boss and the minions may be attacked
// and chosen by abilities that choose monsters.
func TestBossAndMinionsAreMonsters(t *testing.T) {
	tb := slotsTable(t, engine.SituationPlayer{}, "pooter")
	if got := len(tb.G.MonstersInPlay()); got != 4 {
		t.Fatalf("%d monsters in play, want leech, fly, gurdy, pooter", got)
	}
	if _, err := tb.G.Apply(engine.Intent{Player: 0, Kind: engine.IntentAttack}); err != nil {
		t.Fatal(err)
	}
	tb.Pass(0)
	tb.Pass(1)
	opts := tb.Options()
	for _, want := range []string{"leech", "fly", "gurdy", "pooter"} {
		if !slices.Contains(opts, want) {
			t.Errorf("attack options %v lack %s", opts, want)
		}
	}
	tb.G.ForceRolls(6, 6, 6, 6, 6) // Isaac's ATK 1 against HP 5
	tb.Choose("gurdy")
	ev := tb.Settle()
	hits := 0
	for _, e := range ev {
		if e.Kind == engine.EvDamaged && e.Card == "gurdy" && e.Amount == 1 {
			hits++
		}
	}
	if _, ok := tb.G.Bosses[0].TopOf(); ok || hits != 5 {
		t.Errorf("5 hits of 1 kill gurdy: hits %d", hits)
	}
	tb = slotsTable(t, engine.SituationPlayer{Items: []engine.CardRef{"chaos_card"}}, "pooter")
	tb.Activate(0, "chaos_card", 0, "Kill a player or monster.", "pooter")
	if len(tb.G.Minions) != 0 {
		t.Error("Chaos Card kills the minion")
	}
	checkInvariants(t, tb.G)
}

// R-CHAL-05: killed outright by an ability, the final boss takes 3
// damage instead; a minion dies as usual (R-CHAL-11).
func TestBossArmorKill(t *testing.T) {
	tb := slotsTable(t, engine.SituationPlayer{Items: []engine.CardRef{"chaos_card"}}, "pooter")
	tb.Activate(0, "chaos_card", 0, "Kill a player or monster.", "gurdy")
	boss := tb.G.Object(bossOf(tb))
	if boss.Card != "gurdy" || boss.Damage != 3 || tb.G.HP(boss.ID) != 2 {
		t.Fatalf("gurdy %+v, HP %d: want 3 damage and alive", boss, tb.G.HP(boss.ID))
	}
	tb.G.KillObject(boss.ID) // e.g. Euthanasia, Pandora's Box
	tb.G.Run()
	tb.Settle()
	if _, ok := tb.G.Bosses[0].TopOf(); ok {
		t.Fatal("3 more damage kills gurdy (HP 2)")
	}
	checkInvariants(t, tb.G)
}

// R-CHAL-05, R-CHAL-12: one instance of more than 3 non-combat damage
// becomes 3; an event put in a minion slot resolves, goes to discard,
// and its slot is removed.
func TestBossArmorDamageAndMinionEvent(t *testing.T) {
	tb := slotsTable(t, engine.SituationPlayer{}, "pooter")
	tb.G.PutOnDeck(engine.MonsterDeck, "tnt_barrel")
	if nid := tb.G.MinionFromDeck(); tb.G.Object(nid).Role != engine.RoleEvent || len(tb.G.Minions) != 2 {
		t.Fatalf("tnt_barrel did not come as an event in a new minion slot: %+v", tb.G.Object(nid))
	}
	tb.G.ForceRolls(5) // 5-6: that monster takes 4 damage
	tb.G.Run()
	tb.Settle("gurdy")
	if boss := tb.G.Object(bossOf(tb)); boss.Damage != 3 {
		t.Errorf("gurdy has %d damage, want 4 capped to 3", boss.Damage)
	}
	if got := minionCards(tb.G); !slices.Equal(got, []engine.CardRef{"pooter"}) {
		t.Errorf("minions %v, want the event's slot removed", got)
	}
	d := tb.G.Discards[engine.MonsterDeck]
	if len(d) == 0 || tb.G.Object(d[len(d)-1]).Card != "tnt_barrel" {
		t.Error("the event goes to the monster discard")
	}
	checkInvariants(t, tb.G)
}

// R-CHAL-05, R-CHAL-06: combat damage is never capped; a killed final
// boss dies as usual: rewards, its soul; its slot stays, empty, and is
// never refilled.
func TestBossCombatDamageAndDeath(t *testing.T) {
	tb := slotsTable(t, engine.SituationPlayer{})
	tb.G.Boosts = append(tb.G.Boosts, engine.Boost{Stat: engine.StatPlayerATK, Player: 0, Amount: 4})
	tb.Attack("gurdy", 6)
	if _, ok := tb.G.Bosses[0].TopOf(); ok || len(tb.G.Bosses) != 1 {
		t.Fatalf("5 combat damage kills gurdy (HP 5); bosses %+v", tb.G.Bosses)
	}
	souls := 0
	for _, id := range tb.G.Players[0].InPlay {
		if o := tb.G.Object(id); o.Role == engine.RoleSoul && o.Card == "gurdy" {
			souls++
		}
	}
	if souls != 1 {
		t.Error("Isaac gains gurdy as a soul")
	}
	tb.G.RefillSlots()
	tb.G.Run()
	tb.Settle()
	if _, ok := tb.G.Bosses[0].TopOf(); ok {
		t.Error("the final boss slot was refilled")
	}
	checkInvariants(t, tb.G)
}

// R-CHAL-04, R-CHAL-10, Q-37: effects that would move a final boss or a
// minion do nothing to it: put into discard (Ehwaz, Void), to the deck
// bottom (Emergency Contact), under a card (Friendly Ball), into a deck.
// Monster slot choices (Key) never offer their slots.
func TestBossAndMinionsCantLeave(t *testing.T) {
	tb := slotsTable(t, engine.SituationPlayer{Hand: []engine.CardRef{"ehwaz"}}, "pooter")
	leech := tb.G.MonstersInPlay()[0]
	tb.G.PutOnDeck(engine.MonsterDeck, "fly", "fly") // the refills
	tb.Play(0, "ehwaz")
	if tb.G.Object(leech).Zone.Kind == engine.ZoneInPlay {
		t.Error("Ehwaz: the leech stays")
	}
	if tb.G.Object(bossOf(tb)).Card != "gurdy" || !slices.Equal(minionCards(tb.G), []engine.CardRef{"pooter"}) {
		t.Errorf("Ehwaz moved the boss or the minion: minions %v", minionCards(tb.G))
	}

	tb = slotsTable(t, engine.SituationPlayer{Items: []engine.CardRef{"void", "friendly_ball"}}, "pooter")
	tb.Activate(0, "void", 0, "Put a monster not being attacked into discard.", "gurdy")
	tb.Activate(0, "friendly_ball", 0, "pooter")
	if tb.G.Object(bossOf(tb)).Card != "gurdy" || !slices.Equal(minionCards(tb.G), []engine.CardRef{"pooter"}) {
		t.Errorf("Void or Friendly Ball moved them: minions %v", minionCards(tb.G))
	}
	boss := bossOf(tb)
	tb.G.PutIntoDeck(engine.MonsterDeck, boss, 0)
	tb.G.SlotToDeckBottom(boss)
	if tb.G.Object(boss).Zone.Kind != engine.ZoneInPlay {
		t.Error("the boss went into the deck")
	}

	tb = slotsTable2(t, engine.SituationPlayer{}, engine.SituationPlayer{Hand: []engine.CardRef{"emergency_contact"}}, "pooter")
	tb.G.ForceRolls(1)
	if _, err := tb.G.Apply(engine.Intent{Player: 0, Kind: engine.IntentAttack}); err != nil {
		t.Fatal(err)
	}
	tb.Pass(0)
	tb.Pass(1)
	tb.Choose("gurdy")
	tb.ThrowPending()
	tb.Pass(0)
	tb.G.Players[1].ExtraLootPlays = 1 // as from his ↷
	tb.Play(1, "emergency_contact")
	if tb.G.Attack.On || tb.G.Object(bossOf(tb)).Card != "gurdy" {
		t.Error("Emergency Contact: the attack is cancelled; the boss stays")
	}
	checkInvariants(t, tb.G)

	tb = slotsTable(t, engine.SituationPlayer{Hand: []engine.CardRef{"key"}}, "pooter")
	if _, err := tb.G.Apply(engine.Intent{Player: 0, Kind: engine.IntentPlayLoot, Objects: []engine.ObjectID{tb.Find(0, "key")}}); err != nil {
		t.Fatal(err)
	}
	if _, err := tb.G.Settle("Put the top card of the monster deck in a monster slot not being attacked."); err == nil {
		t.Fatal("Key asked for no slot")
	}
	if got := tb.Options(); !slices.Equal(got, []string{"ask.slot_card 1 leech", "ask.slot_card 2 fly"}) {
		t.Errorf("Key offers %v, want the monster slots only", got)
	}
}

// R-CHAL-08, R-CHAL-11: minions are the monsters in minion slots only;
// a dead minion's slot is removed and the later ones shift down; minion
// slots are never refilled.
func TestMinionSlotsRemoved(t *testing.T) {
	tb := slotsTable(t, engine.SituationPlayer{}, "pooter", "dip", "squirt")
	for _, id := range tb.G.MonstersInPlay() {
		want := slices.Contains([]engine.CardRef{"pooter", "dip", "squirt"}, tb.G.Object(id).Card)
		if tb.G.IsMinion(id) != want {
			t.Errorf("%s: minion %v", tb.G.Object(id).Card, !want)
		}
	}
	squirt, _ := tb.G.Minions[2].TopOf()
	tb.G.Boosts = append(tb.G.Boosts, engine.Boost{Stat: engine.StatPlayerATK, Player: 0, Amount: 4})
	tb.Attack("dip", 6) // dip: HP 1
	if got := minionCards(tb.G); !slices.Equal(got, []engine.CardRef{"pooter", "squirt"}) {
		t.Fatalf("minions %v after dip died", got)
	}
	if z := tb.G.Object(squirt).Zone; z.Kind != engine.ZoneInPlay || z.Index != 1 {
		t.Errorf("squirt's zone %+v, want slot 2 shifted to index 1", z)
	}
	tb.G.RefillSlots()
	tb.G.Run()
	tb.Settle()
	if len(tb.G.Minions) != 2 {
		t.Errorf("%d minion slots after a refill, want 2", len(tb.G.Minions))
	}
	tb.G.PutOnDeck(engine.MonsterDeck, "fly")
	if nid := tb.G.MinionFromDeck(); !tb.G.IsMinion(nid) || tb.G.Object(nid).Zone.Index != 2 {
		t.Errorf("a new minion gets a new slot: %+v", tb.G.Object(nid))
	}
	checkInvariants(t, tb.G)
}
