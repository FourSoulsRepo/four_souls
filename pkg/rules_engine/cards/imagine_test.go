package cards

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math/rand/v2"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

// imagineSteps is how far each game of TestImagineKeepsTheView plays.
const imagineSteps = 300

// seatView is what seat p sees of g.
func seatView(g *engine.Game, p engine.PlayerID) engine.GameView {
	return g.View(engine.Viewer{Kind: engine.ViewPlayer, Player: p})
}

// liveCards lists the cards of every object still in the game, sorted:
// an imagined game deals the same cards, only elsewhere.
func liveCards(g *engine.Game) []engine.CardRef {
	var out []engine.CardRef
	for _, o := range g.Objects {
		if o.Zone.Kind != engine.ZoneGone {
			out = append(out, o.Card)
		}
	}
	slices.Sort(out)
	return out
}

// checkImagined imagines g for every seat and checks the copy: the seat's
// view and, when it must answer, its allowed intents are the game's; the
// copy holds the same cards and its invariants hold (ADR 033, Decision
// 2.10).
func checkImagined(g *engine.Game, seed uint64) error {
	for p := range engine.PlayerID(len(g.Players)) {
		c, err := g.Imagine(p, seed, engine.Known{})
		if err != nil {
			return err
		}
		if a, b := seatView(g, p), seatView(c, p); !reflect.DeepEqual(a, b) {
			return fmt.Errorf("seat %d: view differs: %s", p, firstDiff(a, b))
		}
		if slices.Contains(g.PromptFor(p).Who(), p) {
			if a, b := g.Allowed(p), c.Allowed(p); !reflect.DeepEqual(a, b) {
				return fmt.Errorf("seat %d: allowed differs:\n%+v\n%+v", p, a, b)
			}
		}
		if !slices.Equal(liveCards(g), liveCards(c)) {
			return fmt.Errorf("seat %d: the copy holds other cards", p)
		}
		if err := c.CheckInvariants(); err != nil {
			return fmt.Errorf("seat %d: copy: %w", p, err)
		}
	}
	return nil
}

// firstDiff shows where two values first differ, as indented JSON.
func firstDiff(a, b any) string {
	ja, errA := json.MarshalIndent(a, "", " ")
	jb, errB := json.MarshalIndent(b, "", " ")
	if errA != nil || errB != nil {
		return fmt.Sprintf("no JSON: %v, %v", errA, errB)
	}
	la, lb := strings.Split(string(ja), "\n"), strings.Split(string(jb), "\n")
	for i := range min(len(la), len(lb)) {
		if la[i] != lb[i] {
			lo := max(i-8, 0)
			return fmt.Sprintf("line %d\ngame:\n%s\ncopy:\n%s", i,
				strings.Join(la[lo:min(i+4, len(la))], "\n"), strings.Join(lb[lo:min(i+4, len(lb))], "\n"))
		}
	}
	return "same JSON; differs in Go only"
}

// TestImagineKeepsTheView: at every step of simulated games of every set,
// 2 to 4 players, and each challenge, every seat's imagined game shows
// that seat what the game shows it, and the game is not changed.
func TestImagineKeepsTheView(t *testing.T) {
	for name, s := range cloneSetups() {
		t.Run(name, func(t *testing.T) {
			g, _, err := engine.NewGame(s)
			if err != nil {
				t.Fatal(err)
			}
			check := func(step int) {
				t.Helper()
				sum, err := g.Checksum()
				if err != nil {
					t.Fatal(err)
				}
				if err := checkImagined(g, uint64(step)+1000); err != nil { //nolint:gosec // step is small
					t.Fatalf("step %d: %v", step, err)
				}
				if got, err := g.Checksum(); err != nil || got != sum {
					t.Fatalf("step %d: Imagine changed the game (err %v)", step, err)
				}
			}
			check(-1)
			playSteps(t, g, rand.New(rand.NewPCG(s.Seed, 99)), imagineSteps, check) //nolint:gosec // test moves, not security
		})
	}
}

// forgetsTheDeal checks that the copy depends only on what seat p
// knows: imagining g, and imagining an imagined g1 (same view for p,
// other hidden layout), give the same game for one seed. The whole state
// is compared: every object's card by ID (past objects too), the other
// players' option labels, the RNG.
func forgetsTheDeal(g *engine.Game, p engine.PlayerID, s1, s2 uint64, known engine.Known) error {
	g1, err := g.Imagine(p, s1, known)
	if err != nil {
		return err
	}
	a, err := g1.Imagine(p, s2, known)
	if err != nil {
		return err
	}
	b, err := g.Imagine(p, s2, known)
	if err != nil {
		return err
	}
	if va, vb := seatView(a, p), seatView(b, p); !reflect.DeepEqual(va, vb) {
		return fmt.Errorf("seat %d: view differs: %s", p, firstDiff(va, vb))
	}
	if ia, ib := a.Allowed(p), b.Allowed(p); !reflect.DeepEqual(ia, ib) {
		return fmt.Errorf("seat %d: allowed differs:\n%+v\n%+v", p, ia, ib)
	}
	// The state's JSON, not Save: a world keeps no clean state, so it
	// cannot save at a nested question.
	sa, err := json.Marshal(a)
	if err != nil {
		return err
	}
	sb, err := json.Marshal(b)
	if err != nil {
		return err
	}
	if !bytes.Equal(sa, sb) {
		var ja, jb any
		if err := json.Unmarshal(sa, &ja); err != nil {
			return err
		}
		if err := json.Unmarshal(sb, &jb); err != nil {
			return err
		}
		return fmt.Errorf("seat %d: the copy depends on the real deal: %s", p, firstDiff(ja, jb))
	}
	return nil
}

// TestImagineForgetsTheDeal: at every step of the games of
// TestImagineKeepsTheView, for every seat, the imagined game does not
// depend on where the real hidden cards lie (ADR 033, Decision 2.14).
func TestImagineForgetsTheDeal(t *testing.T) {
	for name, s := range cloneSetups() {
		t.Run(name, func(t *testing.T) {
			g, _, err := engine.NewGame(s)
			if err != nil {
				t.Fatal(err)
			}
			check := func(step int) {
				t.Helper()
				for p := range engine.PlayerID(len(g.Players)) {
					if err := forgetsTheDeal(g, p, uint64(step)+2000, uint64(step)+3000, engine.Known{}); err != nil { //nolint:gosec // step is small
						t.Fatalf("step %d: %v", step, err)
					}
				}
			}
			check(-1)
			playSteps(t, g, rand.New(rand.NewPCG(s.Seed, 99)), imagineSteps, check) //nolint:gosec // test moves, not security
		})
	}
}

// TestImagineForgetsTheDealKnown: the same with remembered cards, also
// ones that are wrong or not hidden.
func TestImagineForgetsTheDealKnown(t *testing.T) {
	g := fourAllSets(t, 40)
	shop, _ := g.Shop[0].TopOf()
	known := engine.Known{Hands: [][]engine.CardRef{nil, handCards(g, 2), handCards(g, 1)}}
	known.Tops[engine.LootDeck] = []engine.CardRef{"", deckTop(g, engine.LootDeck, 2)[1]}
	known.Bottoms[engine.MonsterDeck] = deckTop(g, engine.MonsterDeck, 1)
	known.Tops[engine.TreasureDeck] = []engine.CardRef{g.Objects[shop].Card}
	for seed := range uint64(10) {
		if err := forgetsTheDeal(g, 0, seed, seed+100, known); err != nil {
			t.Fatalf("seed %d: %v", seed, err)
		}
	}
}

// jsonTree is v as decoded JSON: maps, slices, strings, numbers.
func jsonTree(v any) (any, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, err
	}
	var out any
	return out, json.Unmarshal(b, &out)
}

// eachString calls f with every string inside a JSON tree, map keys
// too, and the path to it.
func eachString(v any, path string, f func(path, s string)) {
	switch v := v.(type) {
	case string:
		f(path, v)
	case []any:
		for i, x := range v {
			eachString(x, path+"["+strconv.Itoa(i)+"]", f)
		}
	case map[string]any:
		for k, x := range v {
			f(path+"{key}", k)
			eachString(x, path+"."+k, f)
		}
	}
}

// unknownCards lists the cards seat p cannot know now: the cards of the
// decks, of the other hands and outside the game, less every card p's
// view names and every card of a discard, below the top too (open by
// R-ZONE-07, Decision 2.15). A card with a copy p sees counts as known.
func unknownCards(g *engine.Game, p engine.PlayerID) (map[string]bool, error) {
	out := map[string]bool{}
	add := func(ids []engine.ObjectID) {
		for _, id := range ids {
			out[string(g.Objects[id].Card)] = true
		}
	}
	for _, d := range g.Decks {
		add(d)
	}
	for _, pl := range g.Players {
		if pl.ID != p {
			add(pl.Hand)
		}
	}
	add(g.Outside)
	view, err := jsonTree(seatView(g, p))
	if err != nil {
		return nil, err
	}
	eachString(view, "", func(_, s string) { delete(out, s) })
	for _, d := range g.Discards {
		for _, id := range d {
			delete(out, string(g.Objects[id].Card))
		}
	}
	return out, nil
}

// namesNoHidden checks that the copy imagined by seat p names no card
// p cannot know, outside the two places Imagine deals or renames hidden
// cards: an object's card, and the option labels of another player's
// prompt (ADR 033, Decision 2.13). Everything else is copied from the
// game as it is: question texts, previews, start choices, gains.
func namesNoHidden(g *engine.Game, p engine.PlayerID, seed uint64) error {
	unknown, err := unknownCards(g, p)
	if err != nil {
		return err
	}
	c, err := g.Imagine(p, seed, engine.Known{})
	if err != nil {
		return err
	}
	tree, err := jsonTree(c)
	if err != nil {
		return err
	}
	state, ok := tree.(map[string]any)
	if !ok {
		return fmt.Errorf("the copy is no JSON object: %T", tree)
	}
	objects, ok := state["objects"].([]any)
	if !ok {
		return fmt.Errorf("the copy has no object list: %T", state["objects"])
	}
	for _, o := range objects {
		if o, ok := o.(map[string]any); ok {
			delete(o, "card") // dealt from the pool, or renamed (Decision 2.3, 2.13.1)
		}
	}
	if a := c.Activating; a != nil && a.Ability.Card == c.Objects[a.Source].Card {
		if act, ok := state["activating"].(map[string]any); ok {
			if ab, ok := act["ability"].(map[string]any); ok {
				delete(ab, "card") // the card dealt to its source, as that object's (19.221)
			}
		}
	}
	if c.Waiting.Kind == engine.PromptChoose && c.Waiting.Player != p {
		if w, ok := state["waiting"].(map[string]any); ok {
			delete(w, "options") // renamed (Decision 2.13.2)
		}
	}
	var leaks []string
	eachString(state, "", func(path, s string) {
		if unknown[s] {
			leaks = append(leaks, fmt.Sprintf("%s = %q", path, s))
		}
	})
	if len(leaks) > 0 {
		slices.Sort(leaks)
		return fmt.Errorf("seat %d: the copy names cards it cannot know: %s", p, strings.Join(leaks, "; "))
	}
	if a := c.Activating; a != nil && a.Player != p {
		return playsOn(c, seed)
	}
	return nil
}

// playsOn plays 60 random steps on an imagined copy whose open play is
// another player's: its card was redealt, so the mode or targets picked
// so far were for another card (19.221). The copy must still play on.
func playsOn(c *engine.Game, seed uint64) (err error) {
	defer func() {
		if v := recover(); v != nil {
			err = fmt.Errorf("a copy with another player's play panics: %v", v)
		}
	}()
	rng := rand.New(rand.NewPCG(seed, 7)) //nolint:gosec // test moves, not security
	for step := 0; step < 60 && !c.Over; step++ {
		p := c.Prompt()
		allowed := c.Allowed(p.Player)
		if len(allowed) == 0 {
			return fmt.Errorf("copy step %d: nothing allowed, prompt %+v", step, p)
		}
		in := pick(rng, allowed)
		if in.Kind == engine.IntentDiscard {
			in.Objects = in.Objects[:p.Count]
		}
		if _, err := c.Apply(in); err != nil {
			return fmt.Errorf("copy step %d: allowed %+v refused: %w", step, in, err)
		}
		if err := c.CheckInvariants(); err != nil {
			return fmt.Errorf("copy step %d: %w", step, err)
		}
	}
	return nil
}

// TestImagineNamesNoHiddenCard: at every step of the games of
// TestImagineKeepsTheView, for every seat, the imagined game names a
// card the seat cannot know only where Imagine dealt or renamed it
// (ADR 033, Decision 2.13). TestImagineForgetsTheDeal cannot see a field
// Imagine copies as it is: both of its copies hold the same real name.
func TestImagineNamesNoHiddenCard(t *testing.T) {
	for name, s := range cloneSetups() {
		t.Run(name, func(t *testing.T) {
			g, _, err := engine.NewGame(s)
			if err != nil {
				t.Fatal(err)
			}
			check := func(step int) {
				t.Helper()
				for p := range engine.PlayerID(len(g.Players)) {
					if err := namesNoHidden(g, p, uint64(step)+4000); err != nil { //nolint:gosec // step is small
						t.Fatalf("step %d: %v", step, err)
					}
				}
			}
			check(-1)
			playSteps(t, g, rand.New(rand.NewPCG(s.Seed, 99)), imagineSteps, check) //nolint:gosec // test moves, not security
		})
	}
}

// TestImaginePlayout: random play of 200 steps on an imagined game never
// panics, and its invariants hold all the way.
func TestImaginePlayout(t *testing.T) {
	for name, s := range cloneSetups() {
		t.Run(name, func(t *testing.T) {
			g, _, err := engine.NewGame(s)
			if err != nil {
				t.Fatal(err)
			}
			rng := rand.New(rand.NewPCG(s.Seed, 99)) //nolint:gosec // test moves, not security
			for round := range 4 {
				playSteps(t, g, rng, 80, nil)
				if g.Over {
					break
				}
				for p := range engine.PlayerID(len(g.Players)) {
					c, err := g.Imagine(p, uint64(round), engine.Known{}) //nolint:gosec // round is small
					if err != nil {
						t.Fatal(err)
					}
					playSteps(t, c, rand.New(rand.NewPCG(uint64(round), uint64(p))), 200, func(step int) { //nolint:gosec // test moves, not security
						if err := c.CheckInvariants(); err != nil {
							t.Fatalf("round %d seat %d step %d: %v", round, p, step, err)
						}
					})
				}
			}
		})
	}
}

// fourAllSets is a 4-player game of every set after a few steps.
func fourAllSets(t *testing.T, steps int) *engine.Game {
	t.Helper()
	g, _, err := engine.NewGame(engine.Setup{Seed: 5, Players: 4, Sets: Sets(), BonusSouls: true, Rooms: true})
	if err != nil {
		t.Fatal(err)
	}
	playSteps(t, g, rand.New(rand.NewPCG(5, 99)), steps, nil) //nolint:gosec // test moves, not security
	return g
}

// handCards lists player q's hand cards.
func handCards(g *engine.Game, q engine.PlayerID) []engine.CardRef {
	var out []engine.CardRef
	for _, id := range g.Players[q].Hand {
		out = append(out, g.Objects[id].Card)
	}
	return out
}

// deckTop lists the top n cards of deck d, top first.
func deckTop(g *engine.Game, d engine.DeckKind, n int) []engine.CardRef {
	var out []engine.CardRef
	for _, id := range g.DeckTop(d, n) {
		out = append(out, g.Objects[id].Card)
	}
	return out
}

// TestImagineSeedsDiffer: two seeds deal other hands and deck tops
// differently (ADR 033, Decision 2.11); the seat's own hand stays.
func TestImagineSeedsDiffer(t *testing.T) {
	g := fourAllSets(t, 0)
	a, err := g.Imagine(0, 1, engine.Known{})
	if err != nil {
		t.Fatal(err)
	}
	b, err := g.Imagine(0, 2, engine.Known{})
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(handCards(g, 0), handCards(a, 0)) || !slices.Equal(handCards(g, 0), handCards(b, 0)) {
		t.Error("the seat's own hand changed")
	}
	for q := engine.PlayerID(1); q < 4; q++ {
		if slices.Equal(handCards(a, q), handCards(b, q)) {
			t.Errorf("player %d: both seeds deal the hand %v", q, handCards(a, q))
		}
	}
	for _, d := range []engine.DeckKind{engine.TreasureDeck, engine.LootDeck, engine.MonsterDeck, engine.RoomDeck} {
		if top := deckTop(a, d, 5); slices.Equal(top, deckTop(b, d, 5)) {
			t.Errorf("deck %v: both seeds deal the top %v", d, top)
		}
	}
	if a.RNG == b.RNG || a.RNG == g.RNG {
		t.Error("the copy's RNG is not reseeded")
	}
}

// TestImagineKnown: a known card that is hidden is dealt where the seat
// remembers it; one that is not hidden, or of another deck, is skipped
// (ADR 033, Decision 2.7 – 2.9).
func TestImagineKnown(t *testing.T) {
	g := fourAllSets(t, 0)
	hand1 := handCards(g, 1)
	lootTop := deckTop(g, engine.LootDeck, 2)
	lootBottom := g.Objects[g.Decks[engine.LootDeck][0]].Card
	// A hidden loot card elsewhere: the second card of player 2's hand.
	elsewhere := handCards(g, 2)[1]
	shop, _ := g.Shop[0].TopOf() // seat 0 sees it: not hidden
	shopCard := g.Objects[shop].Card
	treasure := deckTop(g, engine.TreasureDeck, 1)[0] // not a loot card
	known := engine.Known{Hands: [][]engine.CardRef{nil, {hand1[0], elsewhere}}}
	known.Tops[engine.TreasureDeck] = []engine.CardRef{shopCard}
	known.Tops[engine.LootDeck] = []engine.CardRef{"", lootTop[1]}
	known.Bottoms[engine.LootDeck] = []engine.CardRef{lootBottom}
	known.Tops[engine.MonsterDeck] = []engine.CardRef{treasure}
	for seed := range uint64(20) {
		c, err := g.Imagine(0, seed, known)
		if err != nil {
			t.Fatal(err)
		}
		if h := handCards(c, 1); !slices.Contains(h, hand1[0]) || !slices.Contains(h, elsewhere) {
			t.Fatalf("seed %d: player 1's hand %v lacks %s or %s", seed, h, hand1[0], elsewhere)
		}
		if got := deckTop(c, engine.LootDeck, 2)[1]; got != lootTop[1] {
			t.Fatalf("seed %d: second loot card %s, want %s", seed, got, lootTop[1])
		}
		if got := c.Objects[c.Decks[engine.LootDeck][0]].Card; got != lootBottom {
			t.Fatalf("seed %d: loot bottom %s, want %s", seed, got, lootBottom)
		}
		if top := deckTop(c, engine.TreasureDeck, 1)[0]; top == shopCard || c.Objects[shop].Card != shopCard {
			t.Fatalf("seed %d: the shop's %s was dealt to the treasure deck", seed, shopCard)
		}
		if got := deckTop(c, engine.MonsterDeck, 1)[0]; got == treasure {
			t.Fatalf("seed %d: a treasure card on the monster deck", seed)
		}
		if !slices.Equal(liveCards(g), liveCards(c)) {
			t.Fatalf("seed %d: the copy holds other cards", seed)
		}
		if a, b := seatView(g, 0), seatView(c, 0); !reflect.DeepEqual(a, b) {
			t.Fatalf("seed %d: view differs", seed)
		}
	}
}

// TestImagineNoPlayer: a seat that does not exist is an error.
func TestImagineNoPlayer(t *testing.T) {
	g := fourAllSets(t, 0)
	for _, p := range []engine.PlayerID{-1, 4} {
		if _, err := g.Imagine(p, 1, engine.Known{}); err == nil {
			t.Errorf("seat %d: no error", p)
		}
	}
}

// BenchmarkImagine imagines the game of BenchmarkClone for one seat: a
// hard bot makes one per world it tries. After 300 and 3000 steps: a
// world copies no save tracking, so a long game costs it no more
// (19.213).
func BenchmarkImagine(b *testing.B) {
	for _, steps := range []int{300, 3000} {
		b.Run(strconv.Itoa(steps), func(b *testing.B) {
			s := engine.Setup{Seed: 7, Players: 4, Sets: Sets(), BonusSouls: true, Rooms: true}
			g, _, err := engine.NewGame(s)
			if err != nil {
				b.Fatal(err)
			}
			playSteps(b, g, rand.New(rand.NewPCG(7, 99)), steps, nil) //nolint:gosec // test moves, not security
			b.ReportAllocs()
			seed := uint64(0)
			for b.Loop() {
				seed++
				if _, err := g.Imagine(1, seed, engine.Known{}); err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
