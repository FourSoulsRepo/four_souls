package rulesengine

import "testing"

// putMonster replaces the top of monster slot i with a new monster.
func putMonster(g *Game, i int, card CardRef) ObjectID {
	id := g.newObject(card, Zone{Kind: ZoneInPlay, Slot: MonsterSlot, Index: i}, NoPlayer)
	g.Object(id).Role = RoleMonster
	g.Monsters[i].Cards = []ObjectID{id}
	return id
}

// passWhile passes priority while the prompt is priority and cond holds.
func passWhile(t *testing.T, g *Game, cond func() bool) {
	t.Helper()
	for range 200 {
		if g.Prompt().Kind != PromptPriority || !cond() {
			return
		}
		pass(t, g)
	}
	t.Fatal("passing never ended")
}

func choose(t *testing.T, g *Game, label string) {
	t.Helper()
	p := g.Prompt()
	if p.Kind != PromptChoose {
		t.Fatalf("prompt %+v, want a choice containing %q", p, label)
	}
	for i, o := range p.Options {
		if o == label {
			if _, err := g.Apply(Intent{Player: p.Player, Kind: IntentChoose, Choice: i}); err != nil {
				t.Fatal(err)
			}
			return
		}
	}
	t.Fatalf("option %q not in %v", label, p.Options)
}

// attack declares an attack and lets everyone pass until the target choice.
func attack(t *testing.T, g *Game) {
	t.Helper()
	if _, err := g.Apply(Intent{Player: g.Turn.Active, Kind: IntentAttack}); err != nil {
		t.Fatal(err)
	}
	passWhile(t, g, func() bool { return true })
	if p := g.Prompt(); p.Kind != PromptChoose || p.Purpose != ChooseAttackTarget {
		t.Fatalf("after declaring, prompt %+v; want the target choice (R-ATK-02)", p)
	}
}

func TestKillAMonster(t *testing.T) {
	g := toAction(t, 2)
	a := g.Turn.Active
	putMonster(g, 0, "gaper")
	g.forcedRolls = []int{6, 3} // both hit: DC 3 (R-ATK-12)
	cents := g.Players[a].Cents
	attack(t, g)
	choose(t, g, "gaper")
	passWhile(t, g, func() bool { return g.Attack.On || len(g.Stack) > 0 })

	if g.Attack.On {
		t.Fatal("the attack should be over")
	}
	if got := g.Players[a].Cents - cents; got != 3 {
		t.Errorf("rewards: gained %d¢, want 3 (R-DEATH-06)", got)
	}
	if top, _ := g.Monsters[0].TopOf(); g.Object(top).Card == "gaper" && g.Object(top).Damage > 0 {
		t.Error("the dead monster is still in its slot")
	}
	if len(g.Monsters[0].Cards) == 0 {
		t.Error("the slot was not refilled (R-SHOP-06)")
	}
	if g.Turn.Attacks != 0 {
		t.Errorf("attacks left %d, want 0 (R-TURN-07)", g.Turn.Attacks)
	}
	if !g.inOpenActionPhase() || g.Prompt().Player != a {
		t.Errorf("after the attack: prompt %+v", g.Prompt())
	}
	if _, err := g.Apply(Intent{Player: a, Kind: IntentAttack}); err == nil {
		t.Error("attacked twice in one turn (R-TURN-07)")
	}
}

func TestMissDamagesTheAttacker(t *testing.T) {
	g := toAction(t, 2)
	a := g.Turn.Active
	putMonster(g, 0, "gaper")
	g.forcedRolls = []int{2} // miss: DC 3 (R-ATK-13)
	attack(t, g)
	choose(t, g, "gaper")
	// Pass until the first combat damage has resolved.
	passWhile(t, g, func() bool { return g.Players[a].Damage == 0 })
	if g.Players[a].Damage != 1 {
		t.Fatalf("player damage %d, want 1 from a miss", g.Players[a].Damage)
	}
}

func TestPlayerDeathPenalty(t *testing.T) {
	g := toAction(t, 2)
	a := g.Turn.Active
	putMonster(g, 0, "gaper")
	item := giveItem(g, a, "trinket")
	hp := g.PlayerHP(a)
	rolls := make([]int, hp)
	for i := range rolls {
		rolls[i] = 1 // every roll misses
	}
	g.forcedRolls = rolls
	cents, hand := g.Players[a].Cents, len(g.Players[a].Hand)
	attack(t, g)
	choose(t, g, "gaper")
	passWhile(t, g, func() bool { return !g.Players[a].Dead })

	// Death penalty choices (R-DEATH-14): an item, then a loot card.
	for g.Prompt().Kind == PromptPriority {
		pass(t, g)
	}
	choose(t, g, "trinket")
	p := g.Prompt()
	if p.Kind != PromptChoose || p.Purpose != ChoosePenaltyLoot {
		t.Fatalf("prompt %+v; want the loot discard of the death penalty", p)
	}
	if _, err := g.Apply(Intent{Player: a, Kind: IntentChoose, Choice: 0}); err != nil {
		t.Fatal(err)
	}
	if contains(g.Players[a].InPlay, item) {
		t.Error("the penalty item was not destroyed")
	}
	if got := len(g.Players[a].Hand); got != hand-1 {
		t.Errorf("hand %d -> %d; one loot card is discarded", hand, got)
	}
	if got := g.Players[a].Cents; got != max(cents-1, 0) {
		t.Errorf("cents %d -> %d; the penalty loses 1¢", cents, got)
	}
	if g.Object(g.Players[a].Character).Charged {
		t.Error("the character must be deactivated by the penalty")
	}
	if g.Attack.On {
		t.Error("the death cancels the attack (R-DEATH-12)")
	}
	passWhile(t, g, func() bool { return g.Turn.Active == a && g.Turn.Step != StepHandSize })
	if g.Turn.Step < StepEndTriggers && g.Turn.Active == a {
		t.Errorf("after the active player's death the turn goes to its end phase (R-DEATH-16), step %d", g.Turn.Step)
	}
}

func TestKillABossGainsASoul(t *testing.T) {
	g := toAction(t, 2)
	a := g.Turn.Active
	putMonster(g, 1, "boss")
	g.forcedRolls = []int{6, 6, 6, 6, 6, 6, 6, 6}
	attack(t, g)
	choose(t, g, "boss")
	passWhile(t, g, func() bool { return g.Attack.On || len(g.Stack) > 0 })
	if g.SoulValue(a) != 1 {
		t.Errorf("soul value %d, want 1 after killing a boss (R-CARD-14)", g.SoulValue(a))
	}
}

func TestAttackTheMonsterDeck(t *testing.T) {
	g := toAction(t, 2)
	putMonster(g, 0, "gaper")
	// Put a known monster on top of the deck.
	g.Decks[MonsterDeck] = append(g.Decks[MonsterDeck], g.newObject("dummy", DeckZone(MonsterDeck), NoPlayer))
	g.forcedRolls = []int{6}
	attack(t, g)
	choose(t, g, "monster deck")
	if p := g.Prompt(); p.Purpose != ChooseMonsterSlot {
		t.Fatalf("prompt %+v; the attacker picks the slot to cover (R-ATK-08)", p)
	}
	if _, err := g.Apply(Intent{Player: g.Turn.Active, Kind: IntentChoose, Choice: 0}); err != nil {
		t.Fatal(err)
	}
	covered := g.Monsters[0].Cards[0]
	if g.Object(covered).Zone.Kind != ZoneCovered {
		t.Fatal("the slot's monster must be covered, not in play (R-ZONE-10)")
	}
	passWhile(t, g, func() bool { return g.Attack.On || len(g.Stack) > 0 })
	top, _ := g.Monsters[0].TopOf()
	if g.Object(top).Card != "gaper" || g.Object(top).Zone.Kind != ZoneInPlay {
		t.Errorf("after killing the revealed monster the covered one is back in play (R-ZONE-11): %+v", *g.Object(top))
	}
	if top == covered {
		t.Error("an uncovered card enters play as a new object (R-ZONE-01)")
	}
}

func TestNoAttackWithAStack(t *testing.T) {
	g := toAction(t, 2)
	g.push(StackItem{Kind: StackRoll, Controller: g.Turn.Active})
	if _, err := g.Apply(Intent{Player: g.Turn.Active, Kind: IntentAttack}); err == nil {
		t.Error("declared an attack with something on the stack (R-ATK-01)")
	}
}
