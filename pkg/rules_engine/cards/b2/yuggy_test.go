package b2

// Rulings by Yuggy (Maestro Media's Rules Tzar, S-YUG), checked against
// the V2 card text (8.7). Each test names its source: r/FourSouls
// comment IDs (reddit.com/r/FourSouls/comments/<thread>/_/<id>) or X
// status IDs.

import (
	"slices"
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

// monsterInPlay reports whether a monster with this card is in a slot.
func monsterInPlay(g *engine.Game, card engine.CardRef) bool {
	for _, s := range g.Monsters {
		if top, ok := s.TopOf(); ok && g.Object(top).Card == card {
			return true
		}
	}
	return false
}

// chaosTable: Isaac has Chaos Card and Breakfast; Cain holds the loot
// and may play one on Isaac's turn.
func chaosTable(t *testing.T, cainHand ...engine.CardRef) *enginetest.Table {
	t.Helper()
	tb := enginetest.NewSetup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{seat("isaac", "chaos_card", "breakfast"), {Character: "cain", Hand: cainHand, Items: items("glass_cannon")}},
		Monsters: items("fly", "leech"),
	}, Set)
	tb.G.Players[1].ExtraLootPlays = 1
	return tb
}

// startChaosOnLeech puts Chaos Card's "kill the leech" on the stack and
// passes priority to Cain.
func startChaosOnLeech(tb *enginetest.Table) {
	tb.T.Helper()
	tb.Start(0, "chaos_card", 0)
	tb.Choose("Kill a player or monster.")
	tb.Choose("leech")
	tb.Pass(0)
}

// Reddit emwky5t, hqkbih5: death does not clear item abilities from the
// stack, only unresolved attack rolls; Chaos Card still kills.
func TestYuggyChaosCardAfterDeath(t *testing.T) {
	tb := chaosTable(t, "xiii_death")
	startChaosOnLeech(tb)
	tb.Play(1, "xiii_death", me, "breakfast")
	if !tb.G.Players[0].Dead {
		t.Fatal("Isaac should be dead")
	}
	if monsterInPlay(tb.G, "leech") {
		t.Error("Chaos Card did not resolve after its controller died")
	}
}

// Reddit hqm0bb7, jf48ifo: Butter Bean cancels the whole "Destroy this.
// If you do" ability: the card stays, tapped, and nothing dies.
func TestYuggyButterBeanChaosCard(t *testing.T) {
	tb := chaosTable(t, "butter_bean")
	startChaosOnLeech(tb)
	tb.Play(1, "butter_bean", "chaos_card: ↷: Destroy this. If you do, choose one- Kill a player or monster. Destroy an item or soul.")
	if !hasItem(tb.G, 0, "chaos_card") {
		t.Error("Chaos Card was destroyed though its ability was cancelled")
	}
	if tb.G.Object(tb.Find(0, "chaos_card")).Charged {
		t.Error("Chaos Card should stay tapped")
	}
	if !monsterInPlay(tb.G, "leech") {
		t.Error("the leech died though the ability was cancelled")
	}
}

// Reddit jdertem, hqkbih5: an item destroyed in response still resolves,
// but "Destroy this. If you do" fizzles: it can no longer destroy itself.
func TestYuggyChaosCardDestroyedInResponse(t *testing.T) {
	tb := chaosTable(t)
	tb.G.ForceRolls(6) // Glass Cannon's roll: recharge
	startChaosOnLeech(tb)
	tb.Activate(1, "glass_cannon", 0, "chaos_card")
	if hasItem(tb.G, 0, "chaos_card") {
		t.Fatal("Glass Cannon did not destroy Chaos Card")
	}
	if !monsterInPlay(tb.G, "leech") {
		t.Error("the leech died though Chaos Card could not destroy itself")
	}
}

// Reddit foj2hg1: Mom's Shovel may be used when the target has no soul;
// it is destroyed and steals nothing.
func TestYuggyMomsShovelNoSoul(t *testing.T) {
	tb := itemTable(t, items("moms_shovel"))
	tb.G.Recharge(tb.Find(0, "moms_shovel"))
	tb.Activate(0, "moms_shovel", 0, foe)
	if hasItem(tb.G, 0, "moms_shovel") {
		t.Error("Mom's Shovel should be destroyed even with no soul to steal")
	}
	if s := tb.G.SoulValue(0); s != 0 {
		t.Errorf("souls = %d, want 0", s)
	}
}

// Reddit h78jxje, fxaq0xe: an eternal "Destroy this. If you do" item (as
// Eden's) can never destroy itself, so its ability does nothing.
func TestYuggyEternalChaosCard(t *testing.T) {
	tb := chaosTable(t)
	tb.G.Object(tb.Find(0, "chaos_card")).Eternal = true
	tb.Activate(0, "chaos_card", 0, "Kill a player or monster.", "leech")
	if !hasItem(tb.G, 0, "chaos_card") || !monsterInPlay(tb.G, "leech") {
		t.Error("an eternal Chaos Card should do nothing")
	}
}

// Reddit jdq5fui, jdru7iy, hm7kzhc: paying HP is not taking damage
// (Fanny Pack does not loot), and HP you don't have can't be paid.
func TestYuggyPayingHPIsNotDamage(t *testing.T) {
	tb := itemTable(t, items("guppys_paw", "fanny_pack"))
	hand := len(tb.G.Players[0].Hand)
	tb.Activate(0, "guppys_paw", 0, me)
	if h := len(tb.G.Players[0].Hand); h != hand {
		t.Errorf("hand %d → %d: paying HP triggered Fanny Pack", hand, h)
	}
	if hp := tb.G.PlayerHP(0); hp != 1 {
		t.Errorf("HP = %d, want 1 after paying 1", hp)
	}
}

// Reddit jdru7iy: HP you don't have can't be paid. A dead player may
// still use Guppy's Paw, but "Pay 1 HP. If you do" then does nothing.
func TestYuggyCannotPayMissingHP(t *testing.T) {
	tb := itemTable(t, items("guppys_paw"))
	tb.G.Players[0].Damage, tb.G.Players[0].Dead = tb.G.PlayerHP(0), true
	tb.Activate(0, "guppys_paw", 0, foe)
	if n := len(tb.G.Shields); n != 0 {
		t.Errorf("%d shields: the dead player could not pay the HP", n)
	}
}

// Reddit eiv6fqd, ej99sc9: prevented damage is not taken, so "each time
// you take damage" items don't trigger (The Dead Cat, Fanny Pack).
func TestYuggyPreventedDamageIsNotTaken(t *testing.T) {
	tb := itemTable(t, items("the_dead_cat", "fanny_pack"), "bomb")
	tb.G.Object(tb.Find(0, "the_dead_cat")).Counters = []engine.Counter{{Count: 9}} // as when it enters play
	hand := len(tb.G.Players[0].Hand) - 1                                           // the bomb is played
	tb.Play(0, "bomb", me)
	if h := len(tb.G.Players[0].Hand); h != hand {
		t.Errorf("hand %d → %d: Fanny Pack looted for prevented damage", hand, h)
	}
	if d := tb.G.Players[0].Damage; d != 0 {
		t.Errorf("damage = %d, want 0", d)
	}
}

// Reddit i0q9c88: Yum Heart (V2) prevents a whole instance of damage,
// however big.
func TestYuggyYumHeartWholeInstance(t *testing.T) {
	tb := itemTable(t, items("yum_heart"), "gold_bomb")
	tb.Activate(0, "yum_heart", 0, me)
	tb.Play(0, "gold_bomb", me)
	if d := tb.G.Players[0].Damage; d != 0 {
		t.Errorf("damage = %d, want 0: the whole 3 damage is prevented", d)
	}
}

// attackRoll declares an attack on monster in a 2-player game and stops
// with the first attack roll on the stack and the attacker to act.
func attackRoll(tb *enginetest.Table, monster engine.CardRef) {
	tb.T.Helper()
	if _, err := tb.G.Apply(engine.Intent{Player: 0, Kind: engine.IntentAttack}); err != nil {
		tb.T.Fatal(err)
	}
	tb.Pass(0)
	tb.Pass(1)
	tb.Choose(string(monster))
}

// firstDamageTo is the amount of the first damage aimed at a card.
func firstDamageTo(g *engine.Game, events []engine.Event, card engine.CardRef) int {
	for _, e := range events {
		if e.Kind == engine.EvDamagePending && e.Object != 0 && g.Object(e.Object).Card == card {
			return e.Amount
		}
	}
	return 0
}

// Reddit jcb9yrj: a reroll is the same roll, so Curved Horn's +1 for the
// first attack roll still applies after The D6 rerolls it.
func TestYuggyCurvedHornAfterReroll(t *testing.T) {
	tb := enginetest.NewSetup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{seat("isaac", "curved_horn", "the_d6"), seat("cain")},
		Monsters: items("big_spider"), // 3 HP, DC 4
	}, Set)
	tb.G.ForceRolls(1, 6, 6)
	attackRoll(tb, "big_spider")
	ev := tb.Activate(0, "the_d6", 0, "ask.roll_of 1")
	if d := firstDamageTo(tb.G, ev, "big_spider"); d != 2 {
		t.Errorf("first hit dealt %d, want 2: the rerolled first roll keeps Curved Horn's +1", d)
	}
}

// ownRoll is a table with a roll of n on the stack, owned by player 0.
func ownRoll(t *testing.T, n int, seats ...engine.SituationPlayer) *enginetest.Table {
	t.Helper()
	return enginetest.NewSetup(t, engine.SituationSetup{
		Players: seats,
		Stack:   []engine.SituationStack{{DiceRoll: n, Owner: 0}},
	}, Set)
}

// Reddit jle6m6t: a roll changed down to 1 still triggers Sacred Heart
// ("when you would roll a 1").
func TestYuggySacredHeartAfterChangeToOne(t *testing.T) {
	tb := ownRoll(t, 2, seat("isaac", "sacred_heart", "book_of_belial"), seat("cain"))
	ev := tb.Activate(0, "book_of_belial", 0, "Subtract 1 from a roll.", "ask.roll_of 2", "ask.yes")
	if got := resolvedRoll(t, ev); got != 6 {
		t.Errorf("roll resolved as %d, want 6", got)
	}
}

// Reddit je719bd, je78lf9: Sacred Heart and Dad's Lost Coin trigger on
// the same 1; they go on the stack in turn order, so the other player's
// Dad's Lost Coin resolves first, and Sacred Heart then fizzles.
func TestYuggySacredHeartAndDadsLostCoin(t *testing.T) {
	tb := ownRoll(t, 1, seat("isaac", "sacred_heart"), seat("cain", "dads_lost_coin"))
	tb.G.ForceRolls(4)
	ev := tb.Settle("ask.yes", "ask.yes")
	if got := resolvedRoll(t, ev); got != 4 {
		t.Errorf("roll resolved as %d, want 4: the reroll comes first and Sacred Heart fizzles", got)
	}
}

// Reddit jc9tj7z: a roll changed to 1 by Godhead triggers Sacred Heart
// again.
func TestYuggyGodheadThenSacredHeart(t *testing.T) {
	tb := ownRoll(t, 6, seat("isaac", "sacred_heart"), seat("cain", "godhead"))
	tb.Pass(0)
	ev := tb.Activate(1, "godhead", 0, "ask.roll_of 6", "1", "ask.yes")
	if got := resolvedRoll(t, ev); got != 6 {
		t.Errorf("roll resolved as %d, want 6", got)
	}
}

// Reddit eu43m5f: Glass Cannon targets The D20, which rerolls itself in
// response; Glass Cannon then has no target and fizzles: it stays,
// tapped, and does not roll.
func TestYuggyGlassCannonLosesItsTarget(t *testing.T) {
	tb := enginetest.NewSetup(t, engine.SituationSetup{Players: []engine.SituationPlayer{seat("isaac", "glass_cannon"), seat("cain", "the_d20")}}, Set)
	tb.Start(0, "glass_cannon", 0)
	tb.Choose("the_d20")
	tb.Pass(0)
	tb.Activate(1, "the_d20", 0, "the_d20")
	if !hasItem(tb.G, 0, "glass_cannon") {
		t.Error("Glass Cannon was destroyed though it fizzled")
	}
	if tb.G.Object(tb.Find(0, "glass_cannon")).Charged {
		t.Error("Glass Cannon should stay tapped")
	}
	if h := len(tb.G.Players[0].Hand); h != 0 {
		t.Errorf("hand = %d: Glass Cannon rolled though it fizzled", h)
	}
}

// Reddit jdgwst3: damage belongs to the player whose loot or item deals
// it; Sloth makes its killer discard, here the non-active player.
func TestYuggySlothKilledByAnotherPlayer(t *testing.T) {
	tb := enginetest.NewSetup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{{Character: "isaac", Hand: items("a_penny")}, {Character: "cain", Hand: items("bomb", "a_dime")}},
		Monsters: items("sloth"),
	}, Set)
	tb.G.Players[1].ExtraLootPlays = 1
	sloth, _ := tb.G.Monsters[0].TopOf()
	tb.G.Object(sloth).Damage = 2 // one more kills it
	tb.G.ForceRolls(1)            // Isaac misses
	attackRoll(tb, "sloth")
	tb.Pass(0)
	tb.Play(1, "bomb", "sloth")
	if h := len(tb.G.Players[1].Hand); h != 0 {
		t.Errorf("Cain's hand = %d, want 0: Cain killed Sloth", h)
	}
	if h := len(tb.G.Players[0].Hand); h != 1 {
		t.Errorf("Isaac's hand = %d, want 1: Isaac did not kill Sloth", h)
	}
}

// Reddit jd9u87p: Decoy swaps The Chest away while Glass Cannon targets
// it; it is still a valid target (it did not change zone), and the soul
// it becomes belongs to its new controller.
func TestYuggyDecoyTakesTheChest(t *testing.T) {
	tb := enginetest.NewSetup(t, engine.SituationSetup{Players: []engine.SituationPlayer{seat("isaac", "glass_cannon", "the_chest"), seat("cain", "decoy")}}, Set)
	tb.G.ForceRolls(6)
	tb.Start(0, "glass_cannon", 0)
	tb.Choose("the_chest")
	tb.Pass(0)
	tb.Activate(1, "decoy", 0, "the_chest")
	if s := tb.G.SoulValue(1); s != 1 {
		t.Errorf("Cain's souls = %d, want 1: The Chest became Cain's soul", s)
	}
	if s := tb.G.SoulValue(0); s != 0 {
		t.Errorf("Isaac's souls = %d, want 0", s)
	}
}

// Reddit jdnvvux: Modeling Clay copying an active item stays tapped until
// it is recharged.
func TestYuggyModelingClayCopyStaysTapped(t *testing.T) {
	tb := enginetest.NewSetup(t, engine.SituationSetup{Players: []engine.SituationPlayer{seat("isaac", "modeling_clay"), seat("cain", "mystery_sack")}}, Set)
	tb.Activate(0, "modeling_clay", 0, "mystery_sack")
	if rule := tb.Refused(0, "modeling_clay", 0); rule != "R-ABIL-07" {
		t.Errorf("the tapped copy was refused by %q, want R-ABIL-07", rule)
	}
}

// Reddit hxgca02, eszl7bz: shop items count as in play for targeting, so
// Placebo may copy a shop item's ↷ ability.
func TestYuggyPlaceboCopiesShopItem(t *testing.T) {
	tb := enginetest.NewSetup(t, engine.SituationSetup{
		Players: []engine.SituationPlayer{seat("isaac", "placebo"), seat("cain")},
		Shop:    items("mystery_sack"),
	}, Set)
	if n := len(tb.G.AbilitiesOf(tb.Find(0, "placebo"))); n != 1 {
		t.Errorf("%d borrowed abilities, want 1: the shop's Mystery Sack", n)
	}
}

// Reddit hv6btvu: a bonus soul is gained once per game, even if it is
// destroyed later.
func TestYuggyBonusSoulOncePerGame(t *testing.T) {
	g := bonusGame(t, "soul_of_gluttony")
	for range 10 {
		g.Players[0].Hand = append(g.Players[0].Hand, g.Players[1].Hand...)
		g.Players[1].Hand = nil
	}
	if len(g.Players[0].Hand) < 10 {
		t.Skip("not enough cards")
	}
	runSome(t, g)
	var soul engine.ObjectID
	for _, id := range g.Players[0].InPlay {
		if g.Object(id).Card == "soul_of_gluttony" {
			soul = id
		}
	}
	if soul == 0 {
		t.Fatal("Soul of Gluttony was not gained")
	}
	g.DestroyObject(0, soul)
	runSome(t, g) // still 10+ loot cards
	if s := g.SoulValue(0); s != 0 {
		t.Errorf("soul value = %d, want 0: a destroyed bonus soul is not gained again", s)
	}
}

// Reddit hxr1rch: The D20 may target a shop item.
func TestYuggyD20OnShopItem(t *testing.T) {
	tb := enginetest.NewSetup(t, engine.SituationSetup{
		Players: []engine.SituationPlayer{seat("isaac", "the_d20"), seat("cain")},
		Shop:    items("breakfast"),
	}, Set)
	tb.Activate(0, "the_d20", 0, "breakfast")
	if top, ok := tb.G.Shop[0].TopOf(); !ok || tb.G.Object(top).Card == "breakfast" {
		t.Error("the shop's Breakfast was not rerolled")
	}
}

// Reddit etdouue, fmz4v53 (Q-07): a dead player can't be chosen while a
// living one can; with every other player dead, you must choose yourself.
func TestYuggyDeadPlayersCantBeChosen(t *testing.T) {
	tb := lootTable(t, "xiii_death", "bomb")
	tb.G.Players[0].ExtraLootPlays = 1
	tb.Play(0, "xiii_death", foe)
	if _, err := tb.G.Apply(engine.Intent{Player: 0, Kind: engine.IntentPlayLoot, Objects: []engine.ObjectID{tb.Find(0, "bomb")}}); err != nil {
		t.Fatal(err)
	}
	opts := tb.Options()
	if slices.Contains(opts, foe) {
		t.Errorf("dead Cain is a choice: %v", opts)
	}
	if !slices.Contains(opts, me) {
		t.Errorf("options %v: Isaac must be able to choose himself", opts)
	}
}

// S-OFF "Damage", Discord judge (2024-07-25): a dead player still takes
// damage that reaches them: nothing is marked, but Fanny Pack loots.
func TestDeadPlayerStillTakesDamage(t *testing.T) {
	tb := enginetest.NewSetup(t, engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "isaac", Hand: items("xiii_death", "xvi_the_tower")}, {Character: "cain", Items: items("fanny_pack", "breakfast")},
	}}, Set)
	tb.G.Players[0].ExtraLootPlays = 1
	tb.Play(0, "xiii_death", foe, "breakfast") // Cain's penalty: Breakfast
	hand := len(tb.G.Players[1].Hand)
	tb.G.ForceRolls(1) // each player takes 1 damage
	tb.Play(0, "xvi_the_tower")
	if h := len(tb.G.Players[1].Hand); h != hand+1 {
		t.Errorf("Cain's hand %d → %d: dead Cain took damage, so Fanny Pack loots", hand, h)
	}
	if !tb.G.Players[1].Dead {
		t.Error("Cain should still be dead, not dead twice")
	}
}
