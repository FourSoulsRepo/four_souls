package rulesengine

import (
	"fmt"
	"reflect"
	"slices"
	"testing"
)

// playLoot plays a loot from p's hand, picks label if asked, and resolves.
func playLoot(t *testing.T, g *Game, p PlayerID, card CardRef, label string) {
	t.Helper()
	id := addToHand(g, p, card)
	if _, err := g.Apply(Intent{Player: p, Kind: IntentPlayLoot, Objects: []ObjectID{id}}); err != nil {
		t.Fatal(err)
	}
	if label != "" {
		choose(t, g, label)
	}
	resolveAll(t, g)
}

// playerLabel is how a choice names player p.
func playerLabel(g *Game, p PlayerID) string {
	return fmt.Sprintf("player %d (%s)", p+1, g.Object(g.Players[p].Character).Card)
}

func TestPartsListEachSourceOfAStat(t *testing.T) {
	g := newAbilityGame(t)
	a := g.Turn.Active
	sword := giveItem(g, a, "sword")
	playLoot(t, g, a, "strength_like", playerLabel(g, a))
	parts := g.Parts(StatPlayerATK, a, 0)
	if len(parts) == 0 || g.Object(parts[0].Source).Card != "strength_like" {
		t.Fatalf("parts = %+v, want the played loot first", parts)
	}
	want := []Part{
		{Stat: StatPlayerATK, Source: parts[0].Source, Amount: 1, Lasts: LastsTurn},
		{Stat: StatPlayerATK, Source: sword, Amount: 1, Lasts: LastsInPlay},
	}
	if !reflect.DeepEqual(parts, want) {
		t.Fatalf("parts = %+v, want %+v", parts, want)
	}
	if got := g.bonus(StatPlayerATK, a, 0); got != sumParts(parts) {
		t.Errorf("bonus = %d, parts sum to %d", got, sumParts(parts))
	}
	if got, printed := g.PlayerATK(a), g.def(g.Players[a].Character).ATK; got != printed+2 {
		t.Errorf("ATK = %d, want printed %d + 2", got, printed)
	}
}

func TestNextAttackRollShowsOnTheAttackerTillUsed(t *testing.T) {
	g := newAbilityGame(t)
	a := g.Turn.Active
	b := g.next(a)
	playLoot(t, g, a, "next_roll_loot", "")
	gv := g.View(Viewer{Kind: ViewSpectator})
	got := gv.Players[a].Sources
	want := []StatSourceView{{Stat: "attack_roll", Source: got[0].Source, Card: "next_roll_loot", Amount: 1, Lasts: "next_roll"}}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("active player's sources = %+v, want %+v", got, want)
	}
	if s := gv.Players[b].Sources; s != nil {
		t.Errorf("anyone's next attack roll shows on the active player only, got %+v on the other", s)
	}
	if n := g.useNextAttackRoll(a); n != 1 {
		t.Fatalf("the roll used %d, want 1", n)
	}
	if s := g.View(Viewer{Kind: ViewSpectator}).Players[a].Sources; s != nil {
		t.Errorf("a used next-roll boost still shows: %+v", s)
	}
}

func TestViewListsPlayerSourcesForEveryone(t *testing.T) {
	g := newAbilityGame(t)
	a := g.Turn.Active
	sword := giveItem(g, a, "sword")
	d := g.def(g.Players[a].Character)
	want := []StatSourceView{{Stat: "atk", Source: sword, Card: "sword", Amount: 1, Lasts: "in_play"}}
	for _, v := range []Viewer{{Kind: ViewPlayer, Player: a}, {Kind: ViewPlayer, Player: g.next(a)}, {Kind: ViewSpectator}} {
		pv := g.View(v).Players[a]
		if !reflect.DeepEqual(pv.Sources, want) {
			t.Errorf("viewer %+v: sources = %+v, want %+v", v, pv.Sources, want)
		}
		if pv.PrintedHP != d.HP || pv.PrintedATK != d.ATK {
			t.Errorf("viewer %+v: printed %d/%d, want %d/%d", v, pv.PrintedHP, pv.PrintedATK, d.HP, d.ATK)
		}
		if pv.ATK != pv.PrintedATK+1 {
			t.Errorf("viewer %+v: ATK %d is not printed + sources", v, pv.ATK)
		}
	}
}

func TestViewHidesASourceInAHand(t *testing.T) {
	g := newAbilityGame(t)
	a := g.Turn.Active
	b := g.next(a)
	hidden := addToHand(g, b, "bomb_loot")
	g.Boosts = append(g.Boosts, Boost{Stat: StatPlayerHP, Player: a, Amount: 1, Source: hidden})
	got := g.View(Viewer{Kind: ViewPlayer, Player: a}).Players[a].Sources
	want := []StatSourceView{{Stat: "hp", Amount: 1, Lasts: "turn"}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("sources = %+v, want no card or id for a card in a hand: %+v", got, want)
	}
}

func TestViewListsMonsterSources(t *testing.T) {
	g := newAbilityGame(t)
	a := g.Turn.Active
	stoney := putMonster(g, 0, "stoney_like")
	m := putMonster(g, 1, "gaper")
	g.Boosts = append(g.Boosts, Boost{Stat: StatAttackedDC, Player: a, Amount: -1, Source: stoney})
	g.Attack.On, g.Attack.Target = true, m
	s := g.View(Viewer{Kind: ViewSpectator}).Monsters[1].Top.Stats
	want := []StatSourceView{
		{Stat: "dc", Source: stoney, Card: "stoney_like", Amount: 1, Lasts: "in_play"},
		{Stat: "dc", Source: stoney, Card: "stoney_like", Amount: -1, Lasts: "turn"},
	}
	if !reflect.DeepEqual(s.Sources, want) {
		t.Fatalf("sources = %+v, want %+v", s.Sources, want)
	}
	if s.DC != s.PrintedDC+1-1 {
		t.Errorf("DC %d is not printed %d + sources", s.DC, s.PrintedDC)
	}
}

// The e2e save helper writes this stat by number (e2e/tests/saves.ts).
func TestNextAttackRollStatNumber(t *testing.T) {
	if StatNextAttackRoll != 18 {
		t.Errorf("StatNextAttackRoll = %d; update NEXT_ATTACK_ROLL in e2e/tests/saves.ts", StatNextAttackRoll)
	}
}

// The odds a client shows match the hit check, for every die face
// (15.94): roll changes, a next roll's, DC changes and their clamp, and
// Eternal D6's 2 or 3 that becomes a 1.
func TestAttackOddsMatchTheHitCheck(t *testing.T) {
	cases := []struct {
		name   string
		dc     int                      // added to Gaper's DC 3
		boosts func(a PlayerID) []Boost // the attacker's changes
	}{
		{"no change", 0, nil},
		{"+1 to attack rolls", 0, func(a PlayerID) []Boost { return []Boost{{Stat: StatAttackRoll, Player: a, Amount: 1}} }},
		{"-2 to anyone's next attack roll", 0, func(PlayerID) []Boost { return []Boost{{Stat: StatNextAttackRoll, Player: NoPlayer, Amount: -2}} }},
		{"DC 6", 3, nil},
		{"DC 8 clamps to 6", 5, nil},
		{"DC 1", -4, nil},
		{"a 2 or 3 becomes a 1", 0, func(PlayerID) []Boost { return []Boost{{Stat: StatRollToOne, Player: NoPlayer, Amount: 1}} }},
		{"+2 to rolls against DC 1", -2, func(a PlayerID) []Boost { return []Boost{{Stat: StatRoll, Player: a, Amount: 2}} }},
	}
	for _, c := range cases {
		for d := 1; d <= 6; d++ {
			g := toAction(t, 2)
			a := g.Turn.Active
			m := putMonster(g, 0, "gaper")
			if c.boosts != nil {
				g.Boosts = append(g.Boosts, c.boosts(a)...)
			}
			if c.dc != 0 {
				g.Boosts = append(g.Boosts, Boost{Stat: StatMonsterDC, Player: NoPlayer, Object: m, Amount: c.dc})
			}
			attack(t, g)
			choose(t, g, "gaper")
			odds := g.View(Viewer{Kind: ViewSpectator}).Attack.Odds
			if odds == nil {
				t.Fatalf("%s: no odds once the attack is aimed", c.name)
			}
			want := slices.Contains(odds.Hits, d)
			g.forcedRolls = []int{d}
			resolved := false
			for range 20 {
				for _, e := range pass(t, g) {
					resolved = resolved || e.Kind == EvRollResolved
				}
				if resolved {
					break
				}
			}
			if !resolved {
				t.Fatalf("%s, die %d: the roll never resolved", c.name, d)
			}
			// A hit puts its combat damage to the monster on the stack.
			top := g.Stack[len(g.Stack)-1]
			hit := top.Kind == StackDamage && top.Target.Object == m
			if hit != want {
				t.Errorf("%s, die %d: hit = %v, odds hits %v", c.name, d, hit, odds.Hits)
			}
		}
	}
}

// The odds name the roll changes and count Champion Belt-like ATK for
// the first roll before it is made (15.94, 15.54 note).
func TestAttackOddsNameTheirSources(t *testing.T) {
	g := toAction(t, 2)
	a := g.Turn.Active
	putMonster(g, 0, "gaper")
	g.Boosts = append(g.Boosts, Boost{Stat: StatAttackRoll, Player: a, Amount: 1, Source: g.Players[a].Character})
	attack(t, g)
	choose(t, g, "gaper")
	odds := g.View(Viewer{Kind: ViewSpectator}).Attack.Odds
	if odds.DC != 3 || !reflect.DeepEqual(odds.Hits, []int{2, 3, 4, 5, 6}) {
		t.Errorf("odds %+v, want DC 3, hits on 2+", odds)
	}
	if odds.ATK != g.PlayerATK(a) {
		t.Errorf("ATK %d, want %d", odds.ATK, g.PlayerATK(a))
	}
	if len(odds.Sources) != 1 || odds.Sources[0].Stat != "attack_roll" || odds.Sources[0].Lasts != "turn" {
		t.Errorf("sources %+v, want the +1 to attack rolls", odds.Sources)
	}
}

// A thrown attack roll shows what will change it and the result; the
// resolved roll's event names the same changes (15.94).
func TestThrownRollNamesItsChanges(t *testing.T) {
	g := toAction(t, 2)
	a := g.Turn.Active
	putMonster(g, 0, "gaper")
	sword := giveItem(g, a, "sword")
	g.Boosts = append(g.Boosts, Boost{Stat: StatAttackRoll, Player: a, Amount: 1, Source: sword})
	attack(t, g)
	choose(t, g, "gaper")
	g.forcedRolls = []int{2}
	throwPending(t, g)
	var roll *StackItem
	for _, it := range g.View(Viewer{Kind: ViewSpectator}).Stack {
		if it.Kind == StackRoll {
			roll = &it
		}
	}
	want := []StatSourceView{{Stat: "attack_roll", Source: sword, Card: "sword", Amount: 1, Lasts: "turn"}}
	if roll == nil || roll.Roll != 2 || roll.Final != 3 || !reflect.DeepEqual(roll.Sources, want) {
		t.Fatalf("roll item %+v, want 2 that becomes 3 by the sword's +1", roll)
	}
	var got []StatSourceView
	for range 20 {
		for _, e := range pass(t, g) {
			if e.Kind == EvRollResolved {
				got = e.Sources
			}
		}
		if got != nil {
			break
		}
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("roll_resolved sources %+v, want %+v", got, want)
	}
}
