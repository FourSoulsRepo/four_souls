package rulesengine

import (
	"bytes"
	"encoding/json"
	"errors"
	"slices"
	"strings"
	"testing"
)

// dutySet adds an item that lets the active player attack players at
// DC 4 (R-ATK-20), as Backstabber does, so a character can be a target.
var dutySet = CardSet{Name: "duties", Cards: []CardDef{
	{Ref: "backstab", Kind: TreasureCard, Outside: true, Statics: []Static{{Stat: StatPlayerAttackDC, Amount: 4, Applies: func(*Game, ObjectID, PlayerID, ObjectID) bool { return true }}}},
}}

// buildDuty builds hero_a (active, with backstab) and hero_b, a gaper
// and a boss in the monster slots, and extra in the setup.
func buildDuty(t *testing.T, extra string) *Game {
	t.Helper()
	g, err := buildDutySetup(extra)
	if err != nil {
		t.Fatal(err)
	}
	return g
}

func buildDutySetup(extra string) (*Game, error) {
	data := `{"players": [{"character": "hero_a", "items": ["backstab"]}, {"character": "hero_b"}], "monsters": ["gaper", "boss"]`
	if extra != "" {
		data += ", " + extra
	}
	var s SituationSetup
	if err := jsonStrict(data+"}", &s); err != nil {
		return nil, err
	}
	return BuildTable(s, testSet, dutySet)
}

// jsonStrict decodes data into v; an unknown field is an error.
func jsonStrict(data string, v any) error {
	dec := json.NewDecoder(bytes.NewReader([]byte(data)))
	dec.DisallowUnknownFields()
	return dec.Decode(v)
}

// endRule tries to end the active player's turn; "" when allowed.
func endRule(t *testing.T, g *Game) string {
	t.Helper()
	_, err := g.Apply(Intent{Player: g.Turn.Active, Kind: IntentEndTurn})
	var re *RuleError
	switch {
	case err == nil:
		return ""
	case errors.As(err, &re):
		return re.Rule
	}
	t.Fatal(err)
	return ""
}

// attackOptions declares an attack and returns the target question's
// options (R-ATK-02).
func attackOptions(t *testing.T, g *Game) []string {
	t.Helper()
	if _, err := g.Apply(Intent{Player: g.Turn.Active, Kind: IntentAttack}); err != nil {
		t.Fatal(err)
	}
	if _, err := g.settle(nil, stopAt{question: true}, false); err != nil {
		t.Fatal(err)
	}
	if g.Prompt().Kind != PromptChoose {
		t.Fatalf("prompt %+v; want the target question", g.Prompt())
	}
	return g.PlainOptions()
}

// The duty fields set this turn's duties, as the turn holds them
// (19.235). The end of the turn is refused while one is owed and the
// active player can attack (R-ATK-03); a forced target is the only one
// offered (R-ATK-03), a deck duty offers only the deck.
func TestSituationDuties(t *testing.T) {
	// No field: no duty, even with monsters in the slots (as before).
	g := buildDuty(t, "")
	if g.Turn.MustAttack != 0 || g.Turn.MustAttacks != 0 || g.Turn.MustAttackDeck != 0 {
		t.Fatalf("turn %+v; want no duty", g.Turn)
	}
	if r := endRule(t, g); r != "" {
		t.Errorf("no duty: end of turn refused by %s; want it allowed", r)
	}

	// must_attack: the boss, the first card of that ref the active
	// player may attack.
	g = buildDuty(t, `"must_attack": "boss"`)
	if g.Turn.MustAttack != g.Monsters[1].Cards[0] {
		t.Fatalf("must attack %d; want the boss", g.Turn.MustAttack)
	}
	if r := endRule(t, g); r != "R-ATK-03" {
		t.Errorf("must_attack: end of turn refused by %q; want R-ATK-03", r)
	}
	if got := attackOptions(t, g); !slices.Equal(got, []string{"boss"}) {
		t.Errorf("must_attack: target options %v; want only the boss", got)
	}
	// Once the attack is on it, the duty is met (R-ATK-22).
	if _, err := g.Apply(Intent{Player: 0, Kind: IntentChoose, Choice: g.OptionIndex("boss")}); err != nil {
		t.Fatal(err)
	}
	if g.Turn.MustAttack != 0 {
		t.Errorf("after the attack on the boss: must attack %d; want none", g.Turn.MustAttack)
	}

	// A player's character, when the active player may attack players
	// (R-ATK-20), as Lust for Blood forces.
	g = buildDuty(t, `"must_attack": "hero_b"`)
	if g.Turn.MustAttack != g.Players[1].Character || endRule(t, g) != "R-ATK-03" {
		t.Errorf("must_attack hero_b: must attack %d; want hero_b and the end refused", g.Turn.MustAttack)
	}

	// must_attacks: refused while an attack is left; with none left the
	// duty can't be met, so the turn may end ("if able").
	g = buildDuty(t, `"must_attacks": 1`)
	if g.Turn.MustAttacks != 1 || endRule(t, g) != "R-ATK-03" {
		t.Errorf("must_attacks 1: turn %+v; want 1 owed and the end refused by R-ATK-03", g.Turn)
	}
	if got := attackOptions(t, g); !slices.Contains(got, "gaper") || !slices.Contains(got, "ask.deck.monster") {
		t.Errorf("must_attacks: target options %v; want any target", got)
	}
	if g.Turn.MustAttacks != 0 {
		t.Errorf("after the attack: must attacks %d; want 0", g.Turn.MustAttacks)
	}
	g = buildDuty(t, `"must_attacks": 1, "attacks": 0`)
	if r := endRule(t, g); r != "" {
		t.Errorf("must_attacks 1, attacks 0: end refused by %s; want it allowed", r)
	}

	// must_attack_deck: the next attack goes to the monster deck only.
	g = buildDuty(t, `"must_attacks": 1, "must_attack_deck": 1`)
	if endRule(t, g) != "R-ATK-03" {
		t.Error("must_attack_deck: end of turn allowed; want it refused by R-ATK-03")
	}
	if got := attackOptions(t, g); !slices.Equal(got, []string{"ask.deck.monster"}) {
		t.Errorf("must_attack_deck: target options %v; want only the monster deck", got)
	}

	// The start and loot phases: the start-of-turn triggers set duties
	// (R-TURN-03), and they hold into the action phase.
	for _, phase := range []string{"start", "loot"} {
		g = buildDuty(t, `"phase": "`+phase+`", "must_attack": "gaper", "must_attacks": 1, "must_attack_deck": 1`)
		passAround(t, g, 2)
		if phase == "start" {
			passAround(t, g, 2) // the loot step's window
		}
		if !g.inOpenActionPhase() || g.Turn.MustAttack != g.Monsters[0].Cards[0] || g.Turn.MustAttacks != 1 || g.Turn.MustAttackDeck != 1 {
			t.Fatalf("%s: turn %+v; want the action phase with the duties", phase, g.Turn)
		}
		if r := endRule(t, g); r != "R-ATK-03" {
			t.Errorf("%s: end of turn refused by %q; want R-ATK-03", phase, r)
		}
	}

	// A declared attack: the duties are the turn's as the setup starts,
	// after the declaration counted.
	g = buildDuty(t, `"declared": "attack", "must_attacks": 1`)
	if g.Turn.MustAttacks != 1 || !g.Attack.On {
		t.Errorf("declared attack: turn %+v; want 1 more attack owed", g.Turn)
	}
}

// A bad duty field is an error that names the rule or the options.
func TestSituationDutyErrors(t *testing.T) {
	for _, c := range []struct{ extra, want string }{
		{`"must_attack": "dummy"`, `must_attack: the active player can't attack "dummy", options [gaper boss hero_b]`},
		{`"must_attack": "hero_a"`, `can't attack "hero_a"`},
		{`"must_attacks": -1`, "must_attacks -1 is below 0"},
		{`"must_attack_deck": -2`, "must_attack_deck -2 is below 0"},
		{`"phase": "end", "must_attack": "gaper"`, "phase end: the turn's end is declared, which R-ATK-03 refuses while an attack is owed"},
		{`"phase": "end", "must_attacks": 1`, "R-ATK-03 refuses"},
		{`"phase": "end", "must_attack_deck": 1`, "R-ATK-03"},
		{`"must_atack": "gaper"`, "unknown field"},
		{`"must_attack": "gaper", "stack": [{"kind": "attack_roll", "dice_roll": 3, "card": "gaper", "owner": 0}]`, `the attack roll on the stack attacks "gaper", so the duty is met (R-ATK-22)`},
	} {
		_, err := buildDutySetup(c.extra)
		if err == nil || !strings.Contains(err.Error(), c.want) {
			t.Errorf("%s: error %v; want %q", c.extra, err, c.want)
		}
	}
	// Without the item no player can be attacked.
	g, err := BuildTable(SituationSetup{Players: []SituationPlayer{{Character: "hero_a"}, {Character: "hero_b"}}, Monsters: []CardRef{"gaper"}, MustAttack: "hero_b"}, testSet)
	if err == nil || !strings.Contains(err.Error(), "options [gaper]") {
		t.Errorf("hero_b without backstab: %v, %v; want an error with options [gaper]", g, err)
	}
}
