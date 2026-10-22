package rulesengine

import (
	"slices"
	"strconv"
	"strings"
	"testing"
)

// TestFirstPlayerRoll: every player throws in seat order; the lowest goes
// first; ties among the lowest throw again (R-SETUP-11, GS-06, 15.43).
func TestFirstPlayerRoll(t *testing.T) {
	tied := false
	for seed := range uint64(30) {
		g, ev, err := NewGame(Setup{Seed: seed, Players: 3, Sets: []CardSet{testSet}})
		if err != nil {
			t.Fatal(err)
		}
		if slices.ContainsFunc(ev, func(e Event) bool { return e.Kind == EvDiceRolled }) {
			t.Fatalf("seed %d: a die was rolled before its throw", seed)
		}
		round := []PlayerID{0, 1, 2}
		var rolls []int
		number := 1
		for g.FirstRoll != nil {
			if g.FirstRoll.Round != number {
				t.Fatalf("seed %d: round %d, want %d", seed, g.FirstRoll.Round, number)
			}
			if g.Turn.Number != 0 {
				t.Fatalf("seed %d: turn %d during the roll, want 0", seed, g.Turn.Number)
			}
			w := g.Prompt()
			if w.Kind != PromptThrow || w.Text != "first player roll" || w.Player != round[len(rolls)] {
				t.Fatalf("seed %d: prompt %+v, want player %d to throw", seed, w, round[len(rolls)])
			}
			for p := range g.Players {
				if id := PlayerID(p); id != w.Player && len(g.Allowed(id)) > 0 {
					t.Fatalf("seed %d: player %d may act during another's throw", seed, p)
				}
			}
			ev, err := g.Apply(Intent{Player: w.Player, Kind: IntentThrow})
			if err != nil {
				t.Fatal(err)
			}
			if len(ev) == 0 || ev[0].Kind != EvDiceRolled || ev[0].Player != w.Player || ev[0].Text != "first player roll" {
				t.Fatalf("seed %d: throw events %+v, want the die first", seed, ev)
			}
			rolls = append(rolls, ev[0].Amount)
			if len(rolls) < len(round) {
				continue
			}
			lowest := slices.Min(rolls)
			var low []PlayerID
			for i, v := range rolls {
				if v == lowest {
					low = append(low, round[i])
				}
			}
			if len(low) > 1 {
				tied = true
				if g.FirstRoll == nil || !slices.Equal(g.FirstRoll.Players, low) {
					t.Fatalf("seed %d: rolls %v, next round %+v, want %v", seed, rolls, g.FirstRoll, low)
				}
				tie := slices.IndexFunc(ev, func(e Event) bool { return e.Kind == EvFirstRollTied })
				if want := seatList(low); tie < 0 || ev[tie].Amount != lowest || ev[tie].Text != want {
					t.Fatalf("seed %d: tie events %+v, want %d tied by %s", seed, ev, lowest, want)
				}
				round, rolls = low, nil
				number++
				continue
			}
			first := slices.IndexFunc(ev, func(e Event) bool { return e.Kind == EvFirstPlayer })
			if first < 0 || ev[first].Player != low[0] || ev[first].Amount != lowest {
				t.Fatalf("seed %d: rolls %v, events %+v, want player %d first with %d", seed, rolls, ev, low[0], lowest)
			}
			if !slices.ContainsFunc(ev[first:], func(e Event) bool { return e.Kind == EvGameStarted && e.Player == low[0] }) {
				t.Fatalf("seed %d: no game_started after the roll", seed)
			}
		}
		if g.Turn.Number != 1 || g.Turn.Active != round[slices.Index(rolls, slices.Min(rolls))] {
			t.Fatalf("seed %d: turn %+v after rolls %v", seed, g.Turn, rolls)
		}
	}
	if !tied {
		t.Error("no seed tied; the reroll went untested")
	}
}

func seatList(ps []PlayerID) string {
	seats := make([]string, len(ps))
	for i, p := range ps {
		seats[i] = strconv.Itoa(int(p))
	}
	return strings.Join(seats, ",")
}

// TestFirstPlayerRollRefusesOthers: only the asked player throws.
func TestFirstPlayerRollRefusesOthers(t *testing.T) {
	g, _, err := NewGame(Setup{Seed: 1, Players: 2, Sets: []CardSet{testSet}})
	if err != nil {
		t.Fatal(err)
	}
	for _, in := range []Intent{{Player: 1, Kind: IntentThrow}, {Player: 0, Kind: IntentPass}, {Player: 0, Kind: IntentEndTurn}} {
		if _, err := g.Apply(in); err == nil {
			t.Errorf("%+v accepted during player 0's first-player throw", in)
		}
	}
}

// TestFirstPlayerRollReplays: the throws are steps of the record, and a
// replay draws the same dice (ADR 007, ADR 014).
func TestFirstPlayerRollReplays(t *testing.T) {
	s := Setup{Seed: 4, Players: 4, Sets: []CardSet{testSet}}
	g, _, err := NewGame(s)
	if err != nil {
		t.Fatal(err)
	}
	var steps []StepRecord
	for g.FirstRoll != nil {
		st, serr := g.Step(Intent{Player: g.Prompt().Player, Kind: IntentThrow})
		if serr != nil {
			t.Fatal(serr)
		}
		steps = append(steps, st)
	}
	h, err := Replay(s, steps)
	if err != nil {
		t.Fatal(err)
	}
	if h.Turn.Active != g.Turn.Active || h.Turn.Number != 1 {
		t.Errorf("replayed turn %+v, want %+v", h.Turn, g.Turn)
	}
}

// TestFirstPlayerRollInTheView: everyone sees who still rolls and the
// dice so far, so a client that joins late can show them.
func TestFirstPlayerRollInTheView(t *testing.T) {
	g, _, err := NewGame(Setup{Seed: 1, Players: 2, Sets: []CardSet{testSet}})
	if err != nil {
		t.Fatal(err)
	}
	ev, err := g.Apply(Intent{Player: 0, Kind: IntentThrow})
	if err != nil {
		t.Fatal(err)
	}
	fr := g.View(Viewer{Kind: ViewPlayer, Player: 1}).FirstRoll
	if fr == nil || fr.Round != 1 || !slices.Equal(fr.Players, []PlayerID{0, 1}) || !slices.Equal(fr.Rolls, []int{ev[0].Amount}) {
		t.Fatalf("view first roll %+v, want both players and the die %d", fr, ev[0].Amount)
	}
	throwFirstRoll(t, g)
	if g.View(Viewer{Kind: ViewPlayer, Player: 1}).FirstRoll != nil {
		t.Error("the view keeps the first roll after the game started")
	}
}
