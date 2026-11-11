package cards

import (
	"slices"
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/cardkit"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

// The generic blocks of Trick or Treat (14.31), each tried with test
// cards: counters on players, evasion from them, the attack-start
// trigger, reward changes, more damage taken, gained monster
// abilities, and the flip of a two-sided rule card.

var blocksFlip = engine.NewKey("card.blocks_rules.flip")

// blocksFlipAbility: "As each turn starts, you may flip this".
var blocksFlipAbility = engine.Ability{
	Kind: engine.Triggered, Text: "As each turn starts, you may flip this.",
	Trigger: engine.AtStartOfEachTurn(),
	Effects: []engine.Effect{cardkit.May(blocksFlip, engine.FlipThis())},
}

var blocksSet = engine.CardSet{Code: "blocks_trial", Name: "blocks_trial", Cards: []engine.CardDef{
	{Ref: "blocks_hero", Kind: engine.CharacterCard, HP: 5, ATK: 1, Outside: true},
	{
		Ref: "blocks_mob", Kind: engine.MonsterCard, HP: 2, DC: 3, ATK: 1, Outside: true,
		Rewards: []engine.Reward{{Kind: engine.RewardCents, Amount: 3}},
	},
	// blocks_boss: fear on a damaged player; +1 evasion per 2 fear on
	// the active player (R-TRICK-07, R-TRICK-09).
	{
		Ref: "blocks_boss", Kind: engine.MonsterCard, HP: 4, DC: 2, ATK: 1, Soul: 1, Outside: true,
		Rewards: []engine.Reward{{Kind: engine.RewardTreasure, Amount: 1}},
		Statics: []engine.Static{engine.MonstersHavePerCounters(engine.StatMonsterDC, "fear", 2)},
		Abilities: []engine.Ability{{
			Kind: engine.Triggered, Text: "Each time a player takes damage, put a fear counter on them.",
			Trigger: engine.WhenAPlayerTakesDamage(),
			Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) { c.G.AddPlayerCounters(c.EventPlayer, "fear", 1) })},
		}},
	},
	// blocks_proud has its own attack trigger (R-ATK-25).
	{
		Ref: "blocks_proud", Kind: engine.MonsterCard, HP: 2, DC: 6, ATK: 0, Outside: true,
		Abilities: []engine.Ability{{
			Kind: engine.Triggered, Text: "When an attack is declared on this, gain 1¢.",
			Trigger: engine.Trigger{On: engine.EvAttackTarget, Match: func(_ *engine.Game, self engine.ObjectID, e engine.Event) bool { return e.Object == self }},
			Effects: []engine.Effect{engine.GainCents(1)},
		}},
	},
	{
		Ref: "blocks_spiky", Kind: engine.MonsterCard, HP: 1, DC: 1, ATK: 1, Outside: true,
		Abilities: []engine.Ability{{Kind: engine.Triggered, Text: "When this dies, gain 5¢.", Trigger: engine.WhenThisDies(), Effects: []engine.Effect{engine.GainCents(5)}}},
	},
	{
		Ref: "blocks_thorny", Kind: engine.MonsterCard, HP: 1, DC: 1, ATK: 1, Outside: true,
		Abilities: []engine.Ability{{Kind: engine.Triggered, Text: "When this dies, gain 7¢.", Trigger: engine.WhenThisDies(), Effects: []engine.Effect{engine.GainCents(7)}}},
	},
	{Ref: "blocks_event", Kind: engine.EventCard, Outside: true},
	// The rule card: Trick-like front adds loot 1 to an attacked
	// monster's rewards; its back makes them 1¢ (R-TRICK-13, R-TRICK-15).
	{
		Ref: "blocks_rules", Kind: engine.ChallengeCard, Back: "blocks_rules_back",
		Abilities: []engine.Ability{blocksFlipAbility, {
			Kind: engine.Triggered, Text: "At the start of each attack against a monster, add loot 1 to its rewards.",
			Trigger: engine.AtStartOfAttackOnMonster(),
			Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) { c.G.AddReward(c.EventObject, engine.RewardLoot, 1, c.Source) })},
		}},
	},
	{
		Ref: "blocks_rules_back", Kind: engine.ChallengeCard, Outside: true,
		Abilities: []engine.Ability{blocksFlipAbility, {
			Kind: engine.Triggered, Text: "At the start of each attack against a monster, its rewards become 1¢.",
			Trigger: engine.AtStartOfAttackOnMonster(),
			Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) { c.G.SetRewardsToCents(c.EventObject, 1, c.Source) })},
		}},
	},
}}

var blocksChallenge = engine.ChallengeDef{
	ID: "blocks_trial", Set: blocksSet.Code,
	Levels: []engine.ChallengeLevel{{Difficulty: engine.DifficultyNormal, Card: "blocks_rules"}},
}

// blocksTable builds a regular table: two blocks heroes, the monsters
// given, the bosses given in final boss slots. Plain monsters lie on
// the monster deck, so a refill draws no curse or event.
func blocksTable(t *testing.T, monsters, bosses []engine.CardRef) *enginetest.Table {
	t.Helper()
	tb := enginetest.NewSetup(t, engine.SituationSetup{
		Players:  []engine.SituationPlayer{enginetest.Seat("blocks_hero"), enginetest.Seat("blocks_hero")},
		Monsters: monsters, Bosses: bosses,
	}, append(Sets(), blocksSet)...)
	for range 6 {
		tb.G.PutOnDeck(engine.MonsterDeck, "blocks_mob")
	}
	return tb
}

// blocksChallengeTable is blocksTable as a challenge with the blocks
// rule card, Trick-like side up.
func blocksChallengeTable(t *testing.T, monsters []engine.CardRef) *enginetest.Table {
	t.Helper()
	tb := blocksTable(t, monsters, nil)
	if err := tb.G.StartChallenge(&blocksChallenge, engine.DifficultyNormal); err != nil {
		t.Fatal(err)
	}
	return tb
}

// slotTop is the top object of monster slot i.
func slotTop(tb *enginetest.Table, i int) engine.ObjectID {
	tb.T.Helper()
	id, ok := tb.G.Monsters[i].TopOf()
	if !ok {
		tb.T.Fatalf("monster slot %d is empty", i)
	}
	return id
}

// R-MECH-13, Q-45: counters on a player stay on a death and a character
// change; never below 0; the view shows them; each change is an event.
func TestPlayerCounters(t *testing.T) {
	tb := blocksTable(t, []engine.CardRef{"blocks_mob"}, nil)
	g := tb.G
	if d := g.AddPlayerCounters(1, "fear", 3); d != 3 || g.PlayerCounters(1, "fear") != 3 {
		t.Fatalf("added %d, now %d; want 3", d, g.PlayerCounters(1, "fear"))
	}
	if g.PlayerCounters(0, "fear") != 0 || g.PlayersWithCounters("fear")[0] != 1 {
		t.Error("only player 1 has fear")
	}
	g.KillPlayer(1, engine.NoPlayer)
	tb.Settle()
	g.ReplaceCharacter(1, "cain")
	if !g.Players[1].Dead || g.PlayerCounters(1, "fear") != 3 {
		t.Errorf("after a death and a new character: dead %v, fear %d; want 3", g.Players[1].Dead, g.PlayerCounters(1, "fear"))
	}
	v := g.View(engine.Viewer{Kind: engine.ViewPlayer, Player: 0})
	if c := v.Players[1].Counters; len(c) != 1 || c[0] != (engine.Counter{Name: "fear", Count: 3}) {
		t.Errorf("view counters %+v", c)
	}
	if d := g.AddPlayerCounters(1, "fear", -5); d != -3 || g.PlayerCounters(1, "fear") != 0 || len(g.Players[1].Counters) != 0 {
		t.Errorf("removed %d, left %d (%v); want -3, none", d, g.PlayerCounters(1, "fear"), g.Players[1].Counters)
	}
	if d := g.AddPlayerCounters(1, "fear", -1); d != 0 {
		t.Errorf("a removal from none changed %d", d)
	}
}

// R-TRICK-07, R-TRICK-09, Q-50, Q-52: each instance of damage to a player
// gives fear; every 2 fear on the active player give every monster +1
// evasion, all rows, capped at 6, read as each roll resolves; the V-41
// parts name the boss.
func TestEvasionFromActivePlayerCounters(t *testing.T) {
	tb := blocksTable(t, []engine.CardRef{"blocks_mob"}, []engine.CardRef{"blocks_boss"})
	g := tb.G
	mob, boss := slotTop(tb, 0), g.Bosses[0].Cards[0]
	g.AddPlayerCounters(1, "fear", 4)
	if g.Evasion(mob) != 3 {
		t.Errorf("player 1 is not active: DC %d, want 3", g.Evasion(mob))
	}
	g.AddPlayerCounters(0, "fear", 3)
	if g.Evasion(mob) != 4 || g.Evasion(boss) != 3 {
		t.Errorf("3 fear: DC %d and %d, want 4 and 3", g.Evasion(mob), g.Evasion(boss))
	}
	parts := g.Parts(engine.StatMonsterDC, engine.NoPlayer, mob)
	if len(parts) != 1 || parts[0].Source != boss || parts[0].Amount != 1 {
		t.Errorf("parts %+v, want +1 from the boss", parts)
	}
	g.AddPlayerCounters(0, "fear", 17)
	if g.Evasion(mob) != 6 {
		t.Errorf("20 fear: DC %d, want the cap 6", g.Evasion(mob))
	}
	g.AddPlayerCounters(0, "fear", -19)
	// 1 fear: DC 3. A miss gives fear 2, so the next roll needs 4: a 3
	// misses (Q-52: read as it resolves). Misses of 1 then kill the hero.
	ev := tb.Attack("blocks_mob", 2, 3, 1, 1, 1)
	if n := g.PlayerCounters(0, "fear"); n != 6 || !g.Players[0].Dead {
		t.Errorf("fear %d after five misses, want 1 + 5; dead %v", n, g.Players[0].Dead)
	}
	for _, e := range ev {
		if e.Kind == engine.EvDamaged && e.Object == mob {
			t.Error("the roll of 3 hit; the new fear should make it miss")
		}
	}
}

// Q-53, R-ATK-25: the attack-start trigger comes once the attack has a
// monster target; the target's own triggers resolve first; the first
// attack roll waits for both.
func TestAttackStartTrigger(t *testing.T) {
	tb := blocksChallengeTable(t, []engine.CardRef{"blocks_proud", "blocks_mob"})
	ev := tb.Attack("blocks_proud", 6, 6)
	proud, rules, roll := enginetest.ResolvedAt(ev, "blocks_proud"), enginetest.ResolvedAt(ev, "blocks_rules"), enginetest.FirstAttackRoll(ev)
	if proud < 0 || rules < 0 || roll < 0 || proud > rules || rules > roll {
		t.Errorf("order: proud %d, rule card %d, first roll %d; want proud < rule card < roll", proud, rules, roll)
	}
	if n := countKind(ev, engine.EvAttackStart); n != 1 {
		t.Errorf("%d attack_start events, want 1", n)
	}
	// blocks_proud had no rewards; the added loot is its only one (Q-48.5).
	if h := len(tb.G.Players[0].Hand); h != 1 {
		t.Errorf("hand %d, want the added loot 1", h)
	}
}

// Q-53: a deck attack starts once the revealed monster is in its slot;
// a revealed event starts no attack.
func TestAttackStartOnDeck(t *testing.T) {
	tb := blocksChallengeTable(t, []engine.CardRef{"blocks_mob"})
	tb.G.ForceRolls(6, 6, 6)
	ev := tb.RevealFromDeck("blocks_event")
	if countKind(ev, engine.EvAttackStart) != 0 {
		t.Error("an event revealed started an attack")
	}
	tb.G.Turn.Attacks++
	ev = tb.RevealFromDeck("blocks_spiky")
	i := slices.IndexFunc(ev, func(e engine.Event) bool { return e.Kind == engine.EvAttackStart })
	if i < 0 || tb.G.Kind(ev[i].Object) != engine.MonsterCard || ev[i].Card != "blocks_spiky" {
		t.Fatalf("no attack_start on the revealed monster: %+v", ev)
	}
	if roll := enginetest.FirstAttackRoll(ev); roll < i {
		t.Errorf("first roll at %d before attack_start at %d", roll, i)
	}
}

// R-TRICK-13, R-TRICK-15, Q-48: reward changes apply in the order made,
// last till end of turn, and the soul is no reward.
func TestRewardChanges(t *testing.T) {
	tb := blocksTable(t, []engine.CardRef{"blocks_mob", "blocks_boss"}, nil)
	g := tb.G
	mob, boss := slotTop(tb, 0), slotTop(tb, 1)
	g.AddReward(mob, engine.RewardTreasure, 1, 0)
	g.SetRewardsToCents(mob, 1, 0)
	g.AddReward(mob, engine.RewardLoot, 1, 0)
	want := []engine.Reward{{Kind: engine.RewardCents, Amount: 1}, {Kind: engine.RewardLoot, Amount: 1}}
	if got := g.Rewards(mob); !slices.Equal(got, want) || !g.RewardsChanged(mob) {
		t.Fatalf("rewards %+v, want %+v", got, want)
	}
	if got := g.Rewards(boss); len(got) != 1 || got[0].Kind != engine.RewardTreasure || g.RewardsChanged(boss) {
		t.Errorf("the boss's rewards changed: %+v", got)
	}
	tb.Attack("blocks_mob", 6, 6)
	if p := g.Players[0]; p.Cents != 1 || len(p.Hand) != 1 || len(p.InPlay) != 0 {
		t.Errorf("cents %d, hand %d, in play %d; want 1¢, 1 loot, nothing else", p.Cents, len(p.Hand), len(p.InPlay))
	}
	// 1¢ on a monster with a soul: the soul still comes (Q-48.4).
	g.SetRewardsToCents(boss, 1, 0)
	g.Turn.Attacks++
	tb.Attack("blocks_boss", 6, 6, 6, 6)
	if g.Players[0].Cents != 2 || g.SoulValue(0) != 1 {
		t.Errorf("cents %d, souls %d; want 2 and the boss's soul", g.Players[0].Cents, g.SoulValue(0))
	}
	g.AddReward(slotTop(tb, 0), engine.RewardLoot, 1, 0)
	tb.EndTurn()
	if g.RewardChanges != nil {
		t.Errorf("changes outlive the turn: %+v", g.RewardChanges)
	}
}

// R-TRICK-13, Q-46: "damage you would take is increased by 1" adds up and
// ends with the turn; Daddy Haunt uses the same rule.
func TestIncreaseDamageTaken(t *testing.T) {
	tb := blocksTable(t, []engine.CardRef{"blocks_mob"}, nil)
	g := tb.G
	g.IncreaseDamageTaken(0, 1, 0)
	g.IncreaseDamageTaken(0, 1, 0)
	g.IncreaseDamageTaken(1, 1, 0)
	tb.Attack("blocks_mob", 1, 6, 6)
	if d := g.Players[0].Damage; d != 3 {
		t.Errorf("damage %d, want 1 + 2", d)
	}
	tb.EndTurn()
	if g.Boosts != nil {
		t.Errorf("boosts outlive the turn: %+v", g.Boosts)
	}
}

// R-TRICK-13, Q-46, Q-47, Q-22: a monster gains the abilities of the
// first non-event monster revealed; the revealed cards go to the bottom;
// two gains add up; "when this dies" ones trigger; stats and rewards
// stay its own; the gains end with the turn.
func TestGainRevealedMonster(t *testing.T) {
	tb := blocksTable(t, []engine.CardRef{"blocks_mob", "blocks_mob"}, nil)
	g := tb.G
	mob := slotTop(tb, 0)
	g.PutOnDeck(engine.MonsterDeck, "blocks_spiky", "blocks_event")
	if got := g.GainRevealedMonster(mob, 0); got != "blocks_spiky" {
		t.Fatalf("gained %q, want blocks_spiky", got)
	}
	deck := g.Decks[engine.MonsterDeck]
	bottom := []engine.CardRef{g.Object(deck[0]).Card, g.Object(deck[1]).Card}
	slices.Sort(bottom)
	if !slices.Equal(bottom, []engine.CardRef{"blocks_event", "blocks_spiky"}) {
		t.Errorf("deck bottom %v, want both revealed cards", bottom)
	}
	g.PutOnDeck(engine.MonsterDeck, "blocks_thorny")
	g.GainRevealedMonster(mob, 0)
	if g.Evasion(mob) != 3 || g.HP(mob) != 2 {
		t.Errorf("DC %d, HP %d: stats must stay its own", g.Evasion(mob), g.HP(mob))
	}
	other := slotTop(tb, 1)
	g.PutOnDeck(engine.MonsterDeck, "blocks_spiky")
	g.GainRevealedMonster(other, 0)
	// The first mob in the options; both gained death triggers, ordered.
	g.ForceRolls(6, 6)
	tb.Do(engine.Intent{Player: 0, Kind: engine.IntentAttack}, "blocks_mob", "blocks_spiky: When this dies, gain 5¢.")
	if c := g.Players[0].Cents; c != 3+5+7 {
		t.Errorf("cents %d, want rewards 3 + 5 + 7", c)
	}
	tb.EndTurn()
	if o := g.Object(other); o.Gained != "" || o.MoreGained != nil {
		t.Errorf("gains outlive the turn: %q %v", o.Gained, o.MoreGained)
	}
	g.Decks[engine.MonsterDeck], g.Discards[engine.MonsterDeck] = nil, nil
	g.PutOnDeck(engine.MonsterDeck, "blocks_event")
	if got := g.GainRevealedMonster(other, 1); got != "" || len(g.Decks[engine.MonsterDeck]) != 1 {
		t.Errorf("no monster anywhere: gained %q, deck %d", got, len(g.Decks[engine.MonsterDeck]))
	}
}

// R-TRICK-11, Q-43, ADR 023: the rule card may flip as each turn starts;
// only the side up works; the view shows the face up.
func TestFlipRuleCard(t *testing.T) {
	tb := blocksChallengeTable(t, []engine.CardRef{"blocks_mob", "blocks_mob"})
	g := tb.G
	card := g.Challenge.Card
	tb.EndTurn("ask.yes")
	if !g.Object(card).Flipped || g.CardOf(card) != "blocks_rules_back" {
		t.Fatal("the rule card did not flip")
	}
	v := g.View(engine.Viewer{Kind: engine.ViewPlayer, Player: 0})
	if c := v.Challenge.Card; c.Card != "blocks_rules_back" || !c.Flipped {
		t.Errorf("view card %q flipped %v, want the back", c.Card, c.Flipped)
	}
	cents := g.Players[1].Cents
	tb.Attack("blocks_mob", 6, 6)
	if got := g.Players[1].Cents - cents; got != 1 || len(g.Players[1].Hand) != 1 {
		t.Errorf("gained %d¢, hand %d; want 1¢ from the back and the turn's loot only", got, len(g.Players[1].Hand))
	}
	tb.EndTurn("ask.no")
	if !g.Object(card).Flipped {
		t.Error("a no flipped it")
	}
}

// countKind counts the events of kind.
func countKind(events []engine.Event, kind engine.EventKind) int {
	n := 0
	for _, e := range events {
		if e.Kind == kind {
			n++
		}
	}
	return n
}
