package rulesengine

import "testing"

// abilitySet: fake cards built only from blocks, shaped like real ones.
var abilitySet = CardSet{Name: "abilities", Cards: []CardDef{
	{Ref: "coin_bag", Kind: TreasureCard, Outside: true, Tap: true, Abilities: []Ability{
		{Kind: Activated, Text: "↷: Gain 3¢.", Costs: []Cost{Tap()}, Effects: []Effect{GainCents(3)}},
	}},
	{Ref: "razor", Kind: TreasureCard, Outside: true, Abilities: []Ability{ // like Golden Razor Blade
		{
			Kind: Activated, Text: "Pay 5¢: Deal 1 damage to a monster or player.", Costs: []Cost{PayCents(5)},
			Targets: []TargetSpec{Choose(TargetMonsterOrPlayer)}, Effects: []Effect{DealDamage(1, 0)},
		},
	}},
	{Ref: "d6", Kind: TreasureCard, Outside: true, Tap: true, Abilities: []Ability{ // like The D6
		{
			Kind: Activated, Text: "↷: Choose a dice roll. Its controller rerolls it.", Costs: []Cost{Tap()},
			Targets: []TargetSpec{Choose(TargetDiceRoll)}, Effects: []Effect{RerollRoll(0)},
		},
	}},
	{Ref: "allowance", Kind: TreasureCard, Outside: true, Abilities: []Ability{
		{Kind: Triggered, Text: "At the start of your turn, gain 1¢.", Trigger: AtStartOfYourTurn(), Effects: []Effect{GainCents(1)}},
	}},
	{Ref: "allowance2", Kind: TreasureCard, Outside: true, Abilities: []Ability{
		{Kind: Triggered, Text: "At the start of your turn, loot 1.", Trigger: AtStartOfYourTurn(), Effects: []Effect{Loot(1)}},
	}},
	{Ref: "sale", Kind: TreasureCard, Outside: true, Statics: []Static{YouHave(StatShopPrice, -5)}}, // like Steamy Sale
	{Ref: "sword", Kind: TreasureCard, Outside: true, Statics: []Static{YouHave(StatPlayerATK, 1)}},
	{Ref: "bomb_loot", Kind: LootCard, Outside: true, Abilities: []Ability{ // like Bomb!
		{
			Kind: LootAbility, Text: "Deal 1 damage to a monster or player.",
			Targets: []TargetSpec{Choose(TargetMonsterOrPlayer)}, Effects: []Effect{DealDamage(1, 0)},
		},
	}},
	{Ref: "pills_loot", Kind: LootCard, Outside: true, Abilities: []Ability{ // like Pills!
		{Kind: LootAbility, Text: "Roll- 1-2: Loot 1. 3-4: Gain 3¢. 5-6: Lose 4¢.", Effects: []Effect{
			Roll(RollTable{}.Results(1, 2, Loot(1)).Results(3, 4, GainCents(3)).Results(5, 6, LoseCents(4))),
		}},
	}},
	{Ref: "isaac_like", Kind: CharacterCard, Outside: true, HP: 2, ATK: 1, Abilities: []Ability{
		{Kind: Activated, Text: "↷: Play an additional loot card this turn.", Costs: []Cost{Tap()}, Effects: []Effect{AddLootPlays(1)}},
	}},
	{Ref: "stoney_like", Kind: MonsterCard, Outside: true, HP: 3, DC: 3, ATK: 1, Statics: []Static{MonstersHave(StatMonsterDC, 1)}},
}}

func newAbilityGame(t *testing.T) *Game {
	t.Helper()
	g, _, err := NewGame(Setup{Seed: 11, Players: 2, Sets: []CardSet{testSet, abilitySet}})
	if err != nil {
		t.Fatal(err)
	}
	passUntilStep(t, g, StepAction)
	return g
}

func addToHand(g *Game, p PlayerID, card CardRef) ObjectID {
	id := g.newObject(card, Zone{Kind: ZoneHand}, p)
	g.Players[p].Hand = append(g.Players[p].Hand, id)
	return id
}

func activate(t *testing.T, g *Game, p PlayerID, obj ObjectID, i int) {
	t.Helper()
	if _, err := g.Apply(Intent{Player: p, Kind: IntentActivate, Objects: []ObjectID{obj}, Choice: i}); err != nil {
		t.Fatal(err)
	}
}

func resolveAll(t *testing.T, g *Game) {
	t.Helper()
	passWhile(t, g, func() bool { return len(g.Stack) > 0 })
}

func TestTapAbility(t *testing.T) {
	g := newAbilityGame(t)
	a := g.Turn.Active
	bag := giveItem(g, a, "coin_bag")
	cents := g.Players[a].Cents
	activate(t, g, a, bag, 0)
	if g.Object(bag).Charged {
		t.Error("↷ deactivates the item as the cost (R-ABIL-10)")
	}
	if len(g.Stack) != 1 || g.Stack[0].Kind != StackAbility {
		t.Fatalf("the ability goes on the stack: %+v", g.Stack)
	}
	resolveAll(t, g)
	if g.Players[a].Cents != cents+3 {
		t.Errorf("cents %d -> %d, want +3", cents, g.Players[a].Cents)
	}
	if _, err := g.Apply(Intent{Player: a, Kind: IntentActivate, Objects: []ObjectID{bag}}); err == nil {
		t.Error("used a deactivated ↷ item (R-ABIL-10)")
	}
}

func TestPaidAbilityWithTarget(t *testing.T) {
	g := newAbilityGame(t)
	a := g.Turn.Active
	razor := giveItem(g, a, "razor")
	m := putMonster(g, 0, "gaper")
	g.Players[a].Cents = 4
	if _, err := g.Apply(Intent{Player: a, Kind: IntentActivate, Objects: []ObjectID{razor}}); err == nil {
		t.Fatal("activated without enough ¢ (R-ABIL-07)")
	}
	g.Players[a].Cents = 7
	activate(t, g, a, razor, 0)
	p := g.Prompt()
	if p.Kind != PromptChoose || p.Purpose != ChooseTarget {
		t.Fatalf("prompt %+v; the target is always asked (ADR 005)", p)
	}
	choose(t, g, "gaper")
	if g.Players[a].Cents != 2 {
		t.Errorf("cents %d, want 7-5 paid on activation", g.Players[a].Cents)
	}
	resolveAll(t, g)
	if g.Object(m).Damage != 1 {
		t.Errorf("monster damage %d, want 1", g.Object(m).Damage)
	}
}

func TestCancelBeforePaying(t *testing.T) {
	g := newAbilityGame(t)
	a := g.Turn.Active
	razor := giveItem(g, a, "razor")
	putMonster(g, 0, "gaper")
	g.Players[a].Cents = 7
	activate(t, g, a, razor, 0)
	choose(t, g, "cancel")
	if g.Players[a].Cents != 7 || len(g.Stack) != 0 {
		t.Errorf("cancel: cents %d, stack %d; nothing is paid or put on the stack", g.Players[a].Cents, len(g.Stack))
	}
	if p := g.Prompt(); p.Kind != PromptPriority || p.Player != a {
		t.Errorf("after cancel, prompt %+v; priority returns to the player", p)
	}
}

func TestFizzleWhenTargetIsGone(t *testing.T) {
	g := newAbilityGame(t)
	a := g.Turn.Active
	razor := giveItem(g, a, "razor")
	m := putMonster(g, 0, "dummy") // 1 HP
	g.Players[a].Cents = 20
	activate(t, g, a, razor, 0)
	choose(t, g, "dummy")
	// In response, kill the target with a second activation's damage.
	g.Players[a].Cents = 20
	bomb := addToHand(g, a, "bomb_loot")
	if _, err := g.Apply(Intent{Player: a, Kind: IntentPlayLoot, Objects: []ObjectID{bomb}}); err != nil {
		t.Fatal(err)
	}
	choose(t, g, "dummy")
	var fizzled bool
	for range 30 {
		if len(g.Stack) == 0 {
			break
		}
		for _, e := range pass(t, g) {
			if e.Kind == EvAbilityFizzle && e.Card == "razor" {
				fizzled = true
			}
		}
	}
	if g.Object(m).Zone.Kind == ZoneInPlay {
		t.Fatal("the monster should have died from the bomb")
	}
	if !fizzled {
		t.Error("the razor's damage should fizzle: its target is gone (R-ABIL-06)")
	}
}

func TestLootAbilityTargetsWhenPlayed(t *testing.T) {
	g := newAbilityGame(t)
	a := g.Turn.Active
	b := g.next(a)
	bomb := addToHand(g, a, "bomb_loot")
	if _, err := g.Apply(Intent{Player: a, Kind: IntentPlayLoot, Objects: []ObjectID{bomb}}); err != nil {
		t.Fatal(err)
	}
	p := g.Prompt()
	if p.Purpose != ChooseTarget {
		t.Fatalf("prompt %+v; a loot's targets are chosen when played (R-ABIL-05)", p)
	}
	if !contains(g.Players[a].Hand, bomb) || g.Turn.LootPlays != 1 {
		t.Error("before the target is chosen, the card stays in hand and the loot play is unused")
	}
	if _, err := g.Apply(Intent{Player: a, Kind: IntentChoose, Choice: len(p.Options) - 2}); err != nil { // the last player
		t.Fatal(err)
	}
	resolveAll(t, g)
	if g.Players[b].Damage != 1 && g.Players[a].Damage != 1 {
		t.Error("the bomb dealt no damage to the chosen player")
	}
	if len(g.Discards[LootDeck]) != 1 {
		t.Error("the loot goes to the loot discard after resolving (R-CARD-07)")
	}
}

func TestRollAbility(t *testing.T) {
	for _, tc := range []struct {
		roll        int
		cents, hand int
	}{
		{2, 0, 1},
		{4, 3, 0},
		{6, -3, 0}, // had 3¢, loses 4: loses all it has (R-MECH-42)
	} {
		g := newAbilityGame(t)
		a := g.Turn.Active
		g.Players[a].Cents = 3
		pills := addToHand(g, a, "pills_loot")
		hand := len(g.Players[a].Hand) - 1
		g.forcedRolls = []int{tc.roll}
		if _, err := g.Apply(Intent{Player: a, Kind: IntentPlayLoot, Objects: []ObjectID{pills}}); err != nil {
			t.Fatal(err)
		}
		resolveAll(t, g)
		if got := g.Players[a].Cents - 3; got != tc.cents {
			t.Errorf("roll %d: cents %+d, want %+d (R-ABIL-24)", tc.roll, got, tc.cents)
		}
		if got := len(g.Players[a].Hand) - hand; got != tc.hand {
			t.Errorf("roll %d: hand %+d, want %+d", tc.roll, got, tc.hand)
		}
	}
}

func TestRerollADiceRoll(t *testing.T) {
	g := newAbilityGame(t)
	a := g.Turn.Active
	d6 := giveItem(g, a, "d6")
	pills := addToHand(g, a, "pills_loot")
	g.forcedRolls = []int{6, 3} // first roll loses ¢; the reroll gains
	g.Players[a].Cents = 5
	if _, err := g.Apply(Intent{Player: a, Kind: IntentPlayLoot, Objects: []ObjectID{pills}}); err != nil {
		t.Fatal(err)
	}
	// Pass until the roll itself is on the stack, then reroll it.
	passWhile(t, g, func() bool {
		for _, it := range g.Stack {
			if it.Kind == StackRoll {
				return false
			}
		}
		return true
	})
	activate(t, g, g.Prompt().Player, d6, 0)
	choose(t, g, "roll of 6")
	resolveAll(t, g)
	if g.Players[a].Cents != 8 {
		t.Errorf("cents %d; the reroll changed the same roll to 3 (R-MECH-47), want 5+3", g.Players[a].Cents)
	}
}

func TestStartOfTurnTriggers(t *testing.T) {
	g := newAbilityGame(t)
	nextP := g.next(g.Turn.Active)
	giveItem(g, nextP, "allowance")
	giveItem(g, nextP, "allowance2")
	if _, err := g.Apply(Intent{Player: g.Turn.Active, Kind: IntentEndTurn}); err != nil {
		t.Fatal(err)
	}
	cents := g.Players[nextP].Cents
	for g.Prompt().Kind != PromptChoose {
		p := g.Prompt()
		in := Intent{Player: p.Player, Kind: IntentPass}
		if p.Kind == PromptDiscard {
			in = Intent{Player: p.Player, Kind: IntentDiscard, Objects: g.Players[p.Player].Hand[:p.Count]}
		}
		if _, err := g.Apply(in); err != nil {
			t.Fatal(err)
		}
	}
	p := g.Prompt()
	if p.Purpose != ChooseTriggerOrder || p.Player != nextP || len(p.Options) != 2 {
		t.Fatalf("prompt %+v; a player with two triggers orders them (R-ABIL-15)", p)
	}
	if _, err := g.Apply(Intent{Player: nextP, Kind: IntentChoose, Choice: 0}); err != nil {
		t.Fatal(err)
	}
	resolveAll(t, g)
	if g.Players[nextP].Cents != cents+1 {
		t.Errorf("cents %d -> %d; the start-of-turn trigger gains 1¢", cents, g.Players[nextP].Cents)
	}
}

func TestExtraLootPlayOnAnotherTurn(t *testing.T) {
	g := newAbilityGame(t)
	a := g.Turn.Active
	b := g.next(a)
	isaac := g.Players[b].Character
	g.Object(isaac).Card = "isaac_like" // give b an Isaac-like character
	g.Object(isaac).Charged = true
	penny := addToHand(g, b, "penny")
	g.push(StackItem{Kind: StackRoll, Controller: a})
	pass(t, g) // a passes; b has priority on a's turn
	activate(t, g, b, isaac, 0)
	resolveAll(t, g)
	g.push(StackItem{Kind: StackRoll, Controller: a})
	pass(t, g)
	if _, err := g.Apply(Intent{Player: b, Kind: IntentPlayLoot, Objects: []ObjectID{penny}}); err != nil {
		t.Errorf("an extra loot play works on another player's turn: %v", err)
	}
}

func TestStaticAbilities(t *testing.T) {
	g := newAbilityGame(t)
	a := g.Turn.Active
	m := putMonster(g, 0, "gaper") // DC 3
	stoney := putMonster(g, 1, "stoney_like")
	if got := g.Evasion(m); got != 4 {
		t.Errorf("evasion %d; Stoney gives monsters +1 DC, want 4 (R-ABIL-28)", got)
	}
	if got := g.Evasion(stoney); got != 4 {
		t.Errorf("Stoney's own evasion %d, want 4", got)
	}
	atk := g.PlayerATK(a)
	giveItem(g, a, "sword")
	if g.PlayerATK(a) != atk+1 {
		t.Error("a +1 ATK item raises the controller's attack")
	}
	giveItem(g, a, "sale")
	g.Players[a].Cents = 5
	if _, err := g.Apply(Intent{Player: a, Kind: IntentPurchase}); err != nil {
		t.Fatal(err)
	}
	passWhile(t, g, func() bool { return true })
	if _, err := g.Apply(Intent{Player: a, Kind: IntentChoose, Choice: 0}); err != nil {
		t.Fatal(err)
	}
	if g.Players[a].Cents != 0 {
		t.Errorf("cents %d; a shop item costs 10-5 with a sale, want 0", g.Players[a].Cents)
	}
}
