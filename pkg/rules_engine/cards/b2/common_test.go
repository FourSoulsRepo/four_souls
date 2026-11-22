package b2

import (
	"slices"
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/enginetest"
)

// seat is a player with a character and items.
var seat = enginetest.Seat

// table builds a b2 table; player 0 is active.
func table(t *testing.T, seats ...engine.SituationPlayer) *enginetest.Table {
	t.Helper()
	return enginetest.New(t, Set, seats...)
}

// testExtraLootOnAnotherTurn checks a character's ↷ on another player's
// turn: that player gets one more loot play (Player.ExtraLootPlays). The
// rest of the ↷ is in <character>/extra-loot-play-in-response.json.
func testExtraLootOnAnotherTurn(t *testing.T, character engine.CardRef) {
	t.Helper()
	tb := table(t, seat(character), seat(character))
	tb.Start(0, character, 0)
	tb.Pass(0)
	tb.Activate(1, character, 0) // in response, on another player's turn
	if got := tb.G.Players[1].ExtraLootPlays; got != 1 {
		t.Errorf("other player's extra loot plays = %d, want 1", got)
	}
}

// newGame starts a real game with fixed characters.
func newGame(t *testing.T, characters ...engine.CardRef) *engine.Game {
	t.Helper()
	g, _, err := engine.NewGame(engine.Setup{Seed: 7, Players: len(characters), Sets: []engine.CardSet{Set}, Characters: characters})
	if err != nil {
		t.Fatal(err)
	}
	enginetest.ThrowFirstRoll(t, g)
	return g
}

// hasItem reports whether player p controls an item with this card.
func hasItem(g *engine.Game, p int, card engine.CardRef) bool {
	for _, id := range g.Players[p].InPlay {
		if o := g.Object(id); o.Card == card && o.Role == engine.RoleItem {
			return true
		}
	}
	return false
}

// monsterTable builds a b2 table with the given monsters in slots.
func monsterTable(t *testing.T, monsters []engine.CardRef, seats ...engine.SituationPlayer) *enginetest.Table {
	t.Helper()
	return enginetest.NewSetup(t, engine.SituationSetup{Players: seats, Monsters: monsters}, Set)
}

// lootTable: player 0 (Isaac) holds the cards, player 1 is Cain; a fly
// and a leech are in the monster slots.
func lootTable(t *testing.T, hand ...engine.CardRef) *enginetest.Table {
	t.Helper()
	return enginetest.NewSetup(t, engine.SituationSetup{
		Players: []engine.SituationPlayer{
			{Character: "isaac", Hand: hand},
			seat("cain"),
		},
		Monsters: []engine.CardRef{"fly", "leech"},
	}, Set)
}

const (
	me  = "player 1 (isaac)"
	foe = "player 2 (cain)"
)

// inDiscard reports whether the loot discard holds the card.
func inDiscard(g *engine.Game, card engine.CardRef) bool {
	for _, id := range g.Discards[engine.LootDeck] {
		if g.Object(id).Card == card {
			return true
		}
	}
	return false
}

// testGain plays a coin card and checks the cents.
func testGain(t *testing.T, card engine.CardRef, n int) {
	t.Helper()
	tb := lootTable(t, card)
	tb.Play(0, card)
	if got := tb.G.Players[0].Cents; got != n {
		t.Errorf("cents = %d, want %d", got, n)
	}
	if !inDiscard(tb.G, card) {
		t.Error("the loot card is not in the loot discard (R-CARD-07)")
	}
}

// itemTable: player 0 (Isaac) has the items and the hand; player 1 is
// Cain; a fly and a leech are in the monster slots.
func itemTable(t *testing.T, items []engine.CardRef, hand ...engine.CardRef) *enginetest.Table {
	t.Helper()
	return enginetest.NewSetup(t, engine.SituationSetup{
		Players: []engine.SituationPlayer{
			{Character: "isaac", Items: items, Hand: hand},
			seat("cain"),
		},
		Monsters: []engine.CardRef{"fly", "leech"},
	}, Set)
}

// items is a shorthand for a list of cards.
func items(cards ...engine.CardRef) []engine.CardRef { return cards }

// slayTable: the monster in slot 1, Isaac and Cain with the given setups.
func slayTable(t *testing.T, monster engine.CardRef, isaac, cain engine.SituationPlayer) *enginetest.Table {
	t.Helper()
	return enginetest.NewSetup(t, engine.SituationSetup{Players: []engine.SituationPlayer{isaac, cain}, Monsters: items(monster)}, Set)
}

// bonusGame starts a real game with the bonus souls set aside; b2 has
// exactly 3, so all of them are in play.
func bonusGame(t *testing.T, soul engine.CardRef) *engine.Game {
	t.Helper()
	g, _, err := engine.NewGame(engine.Setup{Seed: 3, Players: 2, Sets: []engine.CardSet{Set}, BonusSouls: true, Characters: items("isaac", "cain")})
	if err != nil {
		t.Fatal(err)
	}
	if !slices.ContainsFunc(g.BonusSouls, func(id engine.ObjectID) bool { return g.Object(id).Card == soul }) {
		t.Fatalf("%s is not an active bonus soul", soul)
	}
	return g
}

// runSome applies one allowed intent so the engine runs its checks.
func runSome(t *testing.T, g *engine.Game) {
	t.Helper()
	p := g.Prompt().Player
	allowed := g.Allowed(p)
	if len(allowed) == 0 {
		t.Fatal("nothing allowed")
	}
	in := allowed[0]
	if in.Kind == engine.IntentDiscard {
		in.Objects = in.Objects[:g.Prompt().Count]
	}
	if _, err := g.Apply(in); err != nil {
		t.Fatal(err)
	}
}

// viewDice are the die numbers the spectator view shows on player p's
// card ref (15.67).
func viewDice(tb *enginetest.Table, p int, ref engine.CardRef) []int {
	tb.T.Helper()
	for _, c := range tb.G.View(engine.Viewer{Kind: engine.ViewSpectator}).Players[p].InPlay {
		if c.Card == ref {
			return c.Dice
		}
	}
	tb.T.Fatalf("%s is not in play", ref)
	return nil
}

// rollTable has a roll of n on the stack, owned by player 1.
func rollTable(t *testing.T, n int, seats ...engine.SituationPlayer) *enginetest.Table {
	t.Helper()
	return enginetest.NewSetup(t, engine.SituationSetup{
		Players: seats,
		Stack:   []engine.SituationStack{{DiceRoll: n, Owner: 1}},
	}, Set)
}

// resolvedRoll is the result of the first roll that resolved.
func resolvedRoll(t *testing.T, events []engine.Event) int {
	t.Helper()
	for _, e := range events {
		if e.Kind == engine.EvRollResolved {
			return e.Amount
		}
	}
	t.Fatal("no roll resolved")
	return 0
}

// sackFirst puts the roll result on the stack first, so the other card's
// "each time a player rolls" trigger resolves first (R-ABIL-15, 19.233).
const sackFirst = "mystery_sack: ↷: Roll- 1-2: Loot 1. 3-4: Gain 4¢."
