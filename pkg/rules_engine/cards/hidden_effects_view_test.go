package cards

import (
	"slices"
	"strings"
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

// Hidden effects that are no stats are in PlayerView.Sources too (ADR 037
// items 3.17, 3.18, 18.36): skipped and extra turns, a copied next loot,
// lingering abilities, and game-wide effects on the player who made them.

// hidden is one source as these tests compare it, with its game-wide mark.
type hidden struct {
	stat     string
	card     engine.CardRef
	amount   int
	lasts    string
	everyone bool
}

// hiddenOf are player p's sources as every viewer sees them; it fails
// when two viewers differ (ADR 037 item 6: effects are public).
func hiddenOf(tb *enginetest.Table, p int) []hidden {
	tb.T.Helper()
	var first []hidden
	for i, v := range []engine.Viewer{{Kind: engine.ViewSpectator}, {Kind: engine.ViewPlayer, Player: 0}, {Kind: engine.ViewPlayer, Player: 1}} {
		var got []hidden
		for _, s := range tb.G.View(v).Players[p].Sources {
			got = append(got, hidden{s.Stat, s.Card, s.Amount, s.Lasts, s.Everyone})
		}
		if i == 0 {
			first = got
		} else if !slices.Equal(got, first) {
			tb.T.Errorf("viewer %+v sees %v, the spectator %v", v, got, first)
		}
	}
	return first
}

// sorted orders sources by card, for effects whose order the rules leave
// open (two triggers of one event).
func sorted(h []hidden) []hidden {
	return slices.SortedFunc(slices.Values(h), func(a, b hidden) int { return strings.Compare(string(a.card), string(b.card)) })
}

func twoSeats(isaac engine.SituationPlayer) engine.SituationSetup {
	isaac.Character = "isaac"
	return engine.SituationSetup{Players: []engine.SituationPlayer{isaac, {Character: "cain"}}}
}

// Famine: "When this dies, the active player skips their next turn." The
// skip lasts till that turn; then it is gone.
func TestViewHiddenFamineSkip(t *testing.T) {
	s := twoSeats(engine.SituationPlayer{})
	s.Monsters = []engine.CardRef{"famine", "fly"}
	tb := plainTable(t, s)
	tb.Attack("famine", 6, 6) // Famine: HP 2, DC 3; Isaac: ATK 1
	want := []hidden{{"skip_turn", "famine", 1, "next_turn", false}}
	if got := hiddenOf(tb, 0); !slices.Equal(got, want) {
		t.Fatalf("Isaac's effects %v, want %v", got, want)
	}
	if got := hiddenOf(tb, 1); got != nil {
		t.Errorf("Cain's effects %v, want none: Isaac was the active player", got)
	}
	// Two skips from one card add up: amount 2, then one per skipped turn.
	famine := tb.G.View(engine.Viewer{Kind: engine.ViewSpectator}).Players[0].Sources[0].Source
	tb.G.SkipNextTurn(0, famine)
	if got := hiddenOf(tb, 0); !slices.Equal(got, []hidden{{"skip_turn", "famine", 2, "next_turn", false}}) {
		t.Errorf("two skips by Famine: %v, want one source with amount 2", got)
	}
	tb.G.Players[0].SkipTurns, tb.G.Players[0].SkipSources = 1, tb.G.Players[0].SkipSources[:1]
	tb.EndTurn()
	if got := hiddenOf(tb, 0); !slices.Equal(got, want) {
		t.Errorf("on Cain's turn Isaac's effects %v, want %v: his turn is not skipped yet", got, want)
	}
	tb.EndTurn() // Isaac's turn is skipped: Cain again
	if tb.G.Turn.Active != 1 {
		t.Fatalf("active player %d, want Cain: Isaac skips", tb.G.Turn.Active)
	}
	if got := hiddenOf(tb, 0); got != nil {
		t.Errorf("after the skip Isaac's effects %v, want none", got)
	}
}

// Haunted!: "Each time a player gains a soul, they skip their next
// turn." With Famine's own skip Isaac skips two turns, one per card.
func TestViewHiddenHauntedSkips(t *testing.T) {
	s := twoSeats(engine.SituationPlayer{})
	s.Monsters = []engine.CardRef{"famine", "fly"}
	s.Rooms = []engine.CardRef{"haunted"}
	tb := plainTable(t, s)
	tb.Attack("famine", 6, 6)
	want := []hidden{{"skip_turn", "famine", 1, "next_turn", false}, {"skip_turn", "haunted", 1, "next_turn", false}}
	if got := hiddenOf(tb, 0); !slices.Equal(sorted(got), want) {
		t.Errorf("Isaac's effects %v, want %v", got, want)
	}
}

// Overflow!: "Each player rolls- 1-2: They skip their next turn." Cain
// rolls a 1, Isaac a 3. A skip in a state from before 18.36 has no
// source: it shows first, with no card.
func TestViewHiddenOverflowSkip(t *testing.T) {
	s := twoSeats(engine.SituationPlayer{})
	s.Monsters = []engine.CardRef{"fly", "leech"}
	tb := plainTable(t, s)
	tb.G.PutOnDeck(engine.MonsterDeck, "overflow")
	tb.Attack("fly", 6, 3, 1)
	want := []hidden{{"skip_turn", "overflow", 1, "next_turn", false}}
	if got := hiddenOf(tb, 1); !slices.Equal(got, want) {
		t.Fatalf("Cain's effects %v, want %v", got, want)
	}
	if got := hiddenOf(tb, 0); got != nil {
		t.Errorf("Isaac's effects %v, want none: he rolled a 3", got)
	}
	tb.G.Players[1].SkipTurns++ // an older skip, saved without a source
	want = []hidden{{"skip_turn", "", 1, "next_turn", false}, {"skip_turn", "overflow", 1, "next_turn", false}}
	if got := hiddenOf(tb, 1); !slices.Equal(got, want) {
		t.Errorf("Cain's effects %v, want %v", got, want)
	}
}

// XIX. The Sun: "Put this on the bottom of the loot deck. If you do, take
// an extra turn after this one if it's your turn." It shows this turn.
func TestViewHiddenSunExtraTurn(t *testing.T) {
	tb := plainTable(t, twoSeats(engine.SituationPlayer{Hand: []engine.CardRef{"xix_the_sun"}}))
	tb.Play(0, "xix_the_sun")
	tb.Settle()
	want := []hidden{{"extra_turn", "xix_the_sun", 1, "turn", false}}
	if got := hiddenOf(tb, 0); !slices.Equal(got, want) {
		t.Fatalf("Isaac's effects %v, want %v", got, want)
	}
	if got := hiddenOf(tb, 1); got != nil {
		t.Errorf("Cain's effects %v, want none", got)
	}
	tb.EndTurn()
	if tb.G.Turn.Active != 0 {
		t.Fatalf("active player %d, want Isaac's extra turn", tb.G.Turn.Active)
	}
	if got := hiddenOf(tb, 0); got != nil {
		t.Errorf("in the extra turn Isaac's effects %v, want none", got)
	}
}

// Blank Card: "The next time you play a non-trinket, non-ambush loot card
// this turn, copy it." It shows till that loot is played.
func TestViewHiddenBlankCard(t *testing.T) {
	tb := plainTable(t, twoSeats(engine.SituationPlayer{Items: []engine.CardRef{"blank_card"}, Hand: []engine.CardRef{"a_penny"}}))
	tb.Activate(0, "blank_card", 0)
	tb.Settle()
	want := []hidden{{"copy_next_loot", "blank_card", 1, "turn", false}}
	if got := hiddenOf(tb, 0); !slices.Equal(got, want) {
		t.Fatalf("Isaac's effects %v, want %v", got, want)
	}
	tb.Play(0, "a_penny")
	tb.Settle()
	if got := hiddenOf(tb, 0); got != nil {
		t.Errorf("after the loot Isaac's effects %v, want none", got)
	}
}

// Black Heart: "Choose a player. Prevent the next 1 damage they would take
// this turn. Till end of turn, when that player dies, …": a shield and a
// lingering effect on the chosen player, not on the card's player.
func TestViewHiddenBlackHeart(t *testing.T) {
	tb := plainTable(t, twoSeats(engine.SituationPlayer{Hand: []engine.CardRef{"black_heart"}}))
	tb.Play(0, "black_heart", "player 2 (cain)")
	tb.Settle()
	want := []hidden{{"shield", "black_heart", 1, "turn", false}, {"lingering", "black_heart", 1, "turn", false}}
	if got := hiddenOf(tb, 1); !slices.Equal(got, want) {
		t.Fatalf("Cain's effects %v, want %v", got, want)
	}
	if got := hiddenOf(tb, 0); got != nil {
		t.Errorf("Isaac's effects %v, want none", got)
	}
	tb.EndTurn()
	if got := hiddenOf(tb, 1); got != nil {
		t.Errorf("next turn Cain's effects %v, want none (till end of turn)", got)
	}
}

// The Mirror: "At the start of the turn, the active player may pay HP …
// If they do, each time a monster dies this turn, they gain +1
// treasure." Cain pays: the effect is his, this turn.
func TestViewHiddenMirror(t *testing.T) {
	s := twoSeats(engine.SituationPlayer{})
	s.Rooms = []engine.CardRef{"the_mirror"}
	tb := plainTable(t, s)
	tb.EndTurn("ask.yes")
	want := []hidden{{"lingering", "the_mirror", 1, "turn", false}}
	if got := hiddenOf(tb, 1); !slices.Equal(got, want) {
		t.Fatalf("Cain's effects %v, want %v", got, want)
	}
	if got := hiddenOf(tb, 0); got != nil {
		t.Errorf("Isaac's effects %v, want none", got)
	}
}

// Game-wide effects show on the player who made them, for every player
// (ADR 037 item 3.18); the other player has none.
func TestViewHiddenGameWide(t *testing.T) {
	cases := []struct {
		name string
		// Lemon Mishap: "Till end of turn, each time a player rolls an
		// attack roll of 1 or 2, …". Plan C: "At the end of the turn,
		// kill all players." Eternal D6: "Till end of turn, if a player
		// would roll a 2 or 3, it becomes a 1 instead." A Lucky Penny!:
		// "Rewards are doubled till end of turn." Compost: "The next time
		// a player would loot, …": till then, maybe past this turn.
		// Mutant Spider: "The next time a player would roll a dice, they
		// instead roll 4 dice. You choose …": till the next roll.
		isaac engine.SituationPlayer
		use   func(tb *enginetest.Table)
		want  hidden
	}{
		{
			"lemon_mishap",
			engine.SituationPlayer{Items: []engine.CardRef{"lemon_mishap"}},
			func(tb *enginetest.Table) { tb.Activate(0, "lemon_mishap", 0) },
			hidden{"lingering", "lemon_mishap", 1, "turn", true},
		},
		{
			"plan_c",
			engine.SituationPlayer{Items: []engine.CardRef{"plan_c"}},
			func(tb *enginetest.Table) { tb.Activate(0, "plan_c", 0, "fatty") },
			hidden{"lingering", "plan_c", 1, "turn", true},
		},
		{
			"eternal_d6",
			engine.SituationPlayer{Items: []engine.CardRef{"eternal_d6"}},
			func(tb *enginetest.Table) {
				tb.Activate(0, "eternal_d6", 0)
				tb.G.Object(tb.Find(0, "eternal_d6")).Charged = true
				tb.Activate(0, "eternal_d6", 0)
				tb.Activate(0, "eternal_d6", 1)
			},
			hidden{"roll_to_one", "eternal_d6", 1, "turn", true},
		},
		{
			"a_lucky_penny",
			engine.SituationPlayer{Hand: []engine.CardRef{"a_lucky_penny"}},
			func(tb *enginetest.Table) { tb.Play(0, "a_lucky_penny") },
			hidden{"double_rewards", "a_lucky_penny", 1, "turn", true},
		},
		{
			"compost",
			engine.SituationPlayer{Items: []engine.CardRef{"compost"}},
			func(tb *enginetest.Table) { tb.Activate(0, "compost", 0) },
			hidden{"compost", "compost", 1, "next_loot", true},
		},
		{
			"mutant_spider",
			engine.SituationPlayer{Items: []engine.CardRef{"mutant_spider"}},
			func(tb *enginetest.Table) { tb.Activate(0, "mutant_spider", 0) },
			hidden{"many_dice", "mutant_spider", 4, "next_roll", true},
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			s := twoSeats(c.isaac)
			s.Monsters = []engine.CardRef{"fatty", "fly"}
			tb := plainTable(t, s)
			c.use(tb)
			tb.Settle()
			if got := hiddenOf(tb, 0); !slices.Equal(got, []hidden{c.want}) {
				t.Errorf("Isaac's effects %v, want %v", got, []hidden{c.want})
			}
			if got := hiddenOf(tb, 1); got != nil {
				t.Errorf("Cain's effects %v, want none: Isaac made it", got)
			}
		})
	}
}

// Compost lasts till a player loots: Cain's loot at the start of his turn
// uses it.
func TestViewHiddenCompostUsed(t *testing.T) {
	tb := plainTable(t, twoSeats(engine.SituationPlayer{Items: []engine.CardRef{"compost"}}))
	tb.Activate(0, "compost", 0)
	tb.Settle()
	tb.EndTurn()
	if got := hiddenOf(tb, 0); got != nil {
		t.Errorf("after Cain's loot Isaac's effects %v, want none", got)
	}
}

// Red Candle's chosen number shows as a die on its card (ADR 037 item 8,
// 15.67), not on any player.
func TestViewHiddenRedCandleOnNoPlayer(t *testing.T) {
	tb := plainTable(t, twoSeats(engine.SituationPlayer{Items: []engine.CardRef{"red_candle"}}))
	tb.Activate(0, "red_candle", 0, "3")
	tb.Settle()
	for p := range 2 {
		if got := hiddenOf(tb, p); got != nil {
			t.Errorf("player %d's effects %v, want none: the die is on the card", p, got)
		}
	}
	for _, c := range tb.G.View(engine.Viewer{Kind: engine.ViewSpectator}).Players[0].InPlay {
		if c.Card == "red_candle" && !slices.Equal(c.Dice, []int{3}) {
			t.Errorf("Red Candle's dice %v, want [3]", c.Dice)
		}
	}
}

// Sibling Rivalry: "Choose a non-active player. The next time the active
// player declares an attack this turn, the chosen player must make an
// attack roll after each attack roll the active player makes for the
// attack. If that monster dies this attack, the chosen player also gains
// the rewards." Once the card is tapped nothing on the table names the
// chosen player, so the effect shows on them only, till that attack ends
// (18.38).
func TestViewHiddenSiblingRivalry(t *testing.T) {
	s := engine.SituationSetup{Players: []engine.SituationPlayer{
		{Character: "isaac", Items: []engine.CardRef{"sibling_rivalry"}},
		{Character: "cain"},
		{Character: "maggy"},
	}}
	s.Monsters = []engine.CardRef{"pooter", "fly"}
	tb := plainTable(t, s)
	tb.Activate(0, "sibling_rivalry", 0, "player 2 (cain)")
	want := []hidden{{"next_attack_helper", "sibling_rivalry", 1, "turn", false}}
	check := func(when string) {
		t.Helper()
		if got := hiddenOf(tb, 1); !slices.Equal(got, want) {
			t.Errorf("%s Cain's effects %v, want %v", when, got, want)
		}
		for _, p := range []int{0, 2} {
			if got := hiddenOf(tb, p); got != nil {
				t.Errorf("%s player %d's effects %v, want none: only Cain was chosen", when, p+1, got)
			}
		}
	}
	check("before the attack")
	// The attack is open: Cain rolls after Isaac's rolls, still shown.
	// Isaac and Cain hit with ATK 1 each: the Pooter (HP 2) dies.
	tb.G.ForceRolls(6, 6)
	if _, err := tb.G.Apply(engine.Intent{Player: 0, Kind: engine.IntentAttack}); err != nil {
		t.Fatal(err)
	}
	for p := range 3 {
		tb.Pass(p)
	}
	tb.Choose("pooter")
	if !tb.G.Attack.Started {
		t.Fatal("the attack on the Pooter has not started")
	}
	check("during the attack")
	pooter := tb.G.Attack.Target
	tb.Settle()
	if tb.G.Attack.On || tb.G.Object(pooter).Zone.Kind == engine.ZoneInPlay {
		t.Fatal("the attack is still on, or the Pooter lives")
	}
	if got := hiddenOf(tb, 1); got != nil {
		t.Errorf("after the attack Cain's effects %v, want none: it was that attack only", got)
	}
}

// plainTable is viewTable whose refills and treasure gains draw plain
// cards first, so no test depends on the deal: Clotty has no abilities,
// Book of Sin only a ↷ one.
func plainTable(t *testing.T, s engine.SituationSetup) *enginetest.Table {
	t.Helper()
	tb := viewTable(t, s)
	for range 6 {
		tb.G.PutOnDeck(engine.MonsterDeck, "clotty")
		tb.G.PutOnDeck(engine.TreasureDeck, "book_of_sin")
	}
	return tb
}

// TestEveryLingeringIsShown: a lingering ability shows on a player
// (ShownOn), or is on this list with why not (5.24, 18.36).
func TestEveryLingeringIsShown(t *testing.T) {
	notShown := map[engine.CardRef]string{
		"red_candle": "its chosen number is a die on its card (ADR 037 item 8)",
		"poker_chip": "a replacement puts it on the stack at once; it never lasts",
	}
	for _, set := range Sets() {
		for _, c := range set.Cards {
			for i, a := range c.Lingering {
				if (a.ShownOn == nil) != (notShown[c.Ref] != "") {
					t.Errorf("%s lingering %d: ShownOn set %v; on the not-shown list %q", c.Ref, i, a.ShownOn != nil, notShown[c.Ref])
				}
			}
		}
	}
}
