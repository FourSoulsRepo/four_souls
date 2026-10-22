package rulesengine

import (
	"fmt"
	"reflect"
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

// sum adds up the parts' amounts.
func sum(parts []Part) int {
	n := 0
	for _, pt := range parts {
		n += pt.Amount
	}
	return n
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
	if got := g.bonus(StatPlayerATK, a, 0); got != sum(parts) {
		t.Errorf("bonus = %d, parts sum to %d", got, sum(parts))
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
