package rulesengine

import (
	"slices"
	"strings"
	"testing"
)

// expectTable is a situation for the 23.5 checks: hero_a plays a
// Bomb!-like card on the dummy (HP 1), which dies; boss refills its
// slot from the deck. expect is the step's expect.
func expectTable(expect string) string {
	return `{"sets": ["test", "abilities"],
		"setup": {
			"players": [
				{"character": "hero_a", "items": ["razor", "sword", "coin_bag"], "deactivated": ["coin_bag"],
				 "hand": ["bomb_loot", "penny"], "souls": ["boss"], "counters": [{"card": "razor", "count": 2}]},
				{"character": "hero_b"}
			],
			"active": 0,
			"monsters": ["gaper", "dummy"], "monster_damage": {"gaper": 1},
			"counters": [{"card": "gaper", "name": "level", "count": 1}],
			"decks": {"monster": ["boss", "gaper"]}
		},
		"steps": [{"player": 0, "kind": "play", "card": "bomb_loot", "choices": ["dummy"], "expect": ` + expect + `}]}`
}

// rightExpect follows from the cards: Bomb! deals 1 damage to the dummy
// (HP 1), so it dies (R-DEATH-02) into the monster discard; the slot
// refills from the deck top (R-SHOP-06); the loot goes to the loot
// discard (R-CARD-07) and uses the turn's one loot play (R-TURN-05,
// R-CARD-08). sword: +1 ATK.
const rightExpect = `{"allowed": true,
	"after": [{"player": 0, "hp": 2, "atk": 2, "items": ["coin_bag", "razor", "sword"], "hand_cards": ["penny"],
		"soul_cards": ["boss"], "deactivated": ["coin_bag"], "counters": [{"card": "razor", "count": 2}, {"card": "character", "count": 0}]},
		{"player": 1, "items": [], "hand_cards": [], "deactivated": []}],
	"monsters": [
		{"slot": 1, "card": "gaper", "damage": 1, "hp": 1, "atk": 1, "evasion": 3, "dead": false},
		{"slot": 2, "card": "boss"},
		{"card": "dummy", "dead": true}
	],
	"counters": [{"card": "gaper", "name": "level", "count": 1}, {"card": "gaper", "count": 0}],
	"decks": {"monster": ["gaper"]},
	"discards": {"loot": ["bomb_loot"], "monster": ["dummy"], "treasure": []},
	"active": 0, "loot_plays": 0, "attacks": 1, "over": false, "winners": []}`

func runExpect(t *testing.T, data string) []string {
	t.Helper()
	ans, err := RunSituation([]byte(data), testSet, abilitySet)
	if err != nil {
		t.Fatal(err)
	}
	if ans.Matches == nil {
		t.Fatal("no match result")
	}
	return ans.Mismatch
}

// Every 23.5 check passes when the values are right.
func TestSituationExpectChecksPass(t *testing.T) {
	if got := runExpect(t, expectTable(rightExpect)); len(got) != 0 {
		t.Errorf("mismatch %q; want none", got)
	}
}

// Each check fails alone on a wrong value, and says what it found.
func TestSituationExpectChecksFail(t *testing.T) {
	for _, c := range []struct{ expect, want string }{
		{`"after": [{"player": 0, "hp": 3}]`, "player 0 hp is 2, expected 3"},
		{`"after": [{"player": 0, "atk": 1}]`, "player 0 atk is 2, expected 1"},
		{`"after": [{"player": 0, "items": ["razor", "sword"]}]`, "player 0 items are [razor sword coin_bag], expected [razor sword]"},
		{`"after": [{"player": 1, "items": ["razor"]}]`, "player 1 items are [], expected [razor]"},
		{`"after": [{"player": 0, "hand_cards": ["bomb_loot"]}]`, "player 0 hand_cards are [penny], expected [bomb_loot]"},
		{`"after": [{"player": 0, "soul_cards": []}]`, "player 0 soul_cards are [boss], expected []"},
		{`"after": [{"player": 0, "deactivated": ["character"]}]`, "player 0 deactivated are [coin_bag], expected [character]"},
		{`"after": [{"player": 0, "counters": [{"card": "razor", "count": 3}]}]`, "player 0 counters on razor are 2, expected 3"},
		{`"after": [{"player": 0, "counters": [{"card": "razor", "name": "gold", "count": 1}]}]`, "player 0 gold counters on razor are 0, expected 1"},
		{`"after": [{"player": 0, "counters": [{"card": "penny", "count": 1}]}]`, `player 0 counters: "penny" is not there, options [character razor sword coin_bag boss]`},
		{`"after": [{"player": 2, "hp": 1}]`, "no player 2"},
		{`"monsters": [{"slot": 1, "card": "boss"}]`, `monster slot 1 is "gaper", expected "boss"`},
		{`"monsters": [{"slot": 3, "card": "boss"}]`, "no monster slot 3, slots 1 to 2"},
		{`"monsters": [{"card": "gaper", "damage": 0}]`, "monster gaper damage is 1, expected 0"},
		{`"monsters": [{"slot": 1, "card": "gaper", "hp": 2}]`, "monster slot 1 hp is 1, expected 2"},
		{`"monsters": [{"card": "gaper", "atk": 2}]`, "monster gaper atk is 1, expected 2"},
		{`"monsters": [{"card": "gaper", "evasion": 4}]`, "monster gaper evasion is 3, expected 4"},
		{`"monsters": [{"card": "gaper", "dead": true}]`, "monster gaper dead is false, expected true"},
		{`"monsters": [{"card": "dummy", "dead": false}]`, "monster dummy dead is true, expected false"},
		{`"monsters": [{"card": "dummy", "damage": 0}]`, "monster dummy is in no slot, options [gaper boss]"},
		{`"monsters": [{"damage": 0}]`, "monsters: a check needs card"},
		{`"counters": [{"card": "gaper", "name": "level", "count": 2}]`, "level counters on gaper are 1, expected 2"},
		{`"decks": {"monster": ["boss"]}`, "monster deck top is [gaper], expected [boss]"},
		{`"decks": {"shop": []}`, `"shop" is no deck, options [treasure loot monster room]`},
		{`"discards": {"loot": []}`, "loot discards are [bomb_loot], expected []"},
		{`"discards": {"monster": ["gaper"]}`, "monster discards are [dummy], expected [gaper]"},
		{`"active": 1`, "active is 0, expected 1"},
		{`"loot_plays": 1`, "loot_plays is 0, expected 1"},
		{`"attacks": 0`, "attacks is 1, expected 0"},
		{`"over": true`, "over is false, expected true"},
		{`"winners": [0]`, "winners are [], expected [0]"},
	} {
		got := runExpect(t, expectTable(`{"allowed": true, `+c.expect+`}`))
		if want := "step 1 (play bomb_loot): " + c.want; len(got) != 1 || got[0] != want {
			t.Errorf("%s: mismatch %q; want %q", c.expect, got, want)
		}
	}
}

// overTable: hero_a has 3 souls and plays a Bomb!-like card on a boss
// left with 1 HP; its soul is the fourth.
func overTable(expect string) string {
	return `{"sets": ["test", "abilities"],
		"setup": {"players": [{"character": "hero_a", "hand": ["bomb_loot"], "souls": ["boss", "boss", "boss"]}, {"character": "hero_b"}],
			"monsters": ["boss"], "monster_damage": {"boss": 3}},
		"steps": [{"player": 0, "kind": "play", "card": "bomb_loot", "choices": ["boss"]}],
		"expect": ` + expect + `}`
}

// R-WIN-01: a player with 4 souls wins at once.
func TestSituationExpectGameOver(t *testing.T) {
	if got := runExpect(t, overTable(`{"allowed": true, "over": true, "winners": [0], "after": [{"player": 0, "souls": 4}]}`)); len(got) != 0 {
		t.Errorf("mismatch %q; want none", got)
	}
	got := runExpect(t, overTable(`{"allowed": true, "over": false, "winners": [1]}`))
	want := []string{"end: over is true, expected false", "end: winners are [0], expected [1]"}
	if !slices.Equal(got, want) {
		t.Errorf("mismatch %q; want %q", got, want)
	}
}

// Old payloads keep their meaning: hand and souls stay counts.
func TestSituationExpectOldCounts(t *testing.T) {
	got := runExpect(t, expectTable(`{"allowed": true, "after": [{"player": 0, "hand": 1, "souls": 1}]}`))
	if len(got) != 0 {
		t.Errorf("mismatch %q; want none", got)
	}
	got = runExpect(t, expectTable(`{"allowed": true, "after": [{"player": 0, "hand": 2, "souls": 0}]}`))
	if s := strings.Join(got, "\n"); !strings.Contains(s, "player 0 hand is 1, expected 2") || !strings.Contains(s, "player 0 souls is 1, expected 0") {
		t.Errorf("mismatch %q; want hand and souls as counts", got)
	}
}
