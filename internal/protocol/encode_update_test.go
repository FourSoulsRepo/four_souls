package protocol

import (
	"encoding/json"
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

// fullUpdate sets every optional field of an update and its view, so the
// hand-written encoder covers each omitempty rule (18.17).
func fullUpdate() Update {
	top := engine.CardView{ID: 1, Card: "isaac", Charged: true}
	return Update{
		Step: 7,
		Events: []engine.Event{
			{Kind: "looted", Player: 1, Amount: 2},
			{
				Kind: "stack_resolved", Player: 1, Object: 3, Card: "a_dime", Text: "x <y>", Label: engine.Msg{}, Private: true,
				SeenBy: []engine.PlayerID{0}, Prev: 4, Source: 5, StackID: 6, Dice: []int{1, 2}, By: 2, Item: &engine.StackItem{ID: 1},
			},
		},
		View: engine.GameView{
			Players:    []engine.PlayerView{{ID: 1, Dead: true, Hand: []engine.CardView{top}}},
			Shop:       []engine.SlotView{{Top: &top}, {}},
			Monsters:   nil,
			BonusSouls: []engine.CardView{top},
			DeckSizes:  [4]int{1, 2, 3, 4}, DiscardSizes: [4]int{0, 1, 0, 2},
			Discards:   [4]*engine.CardView{&top, nil, nil, &top},
			Waiting:    engine.Prompt{Kind: 1, Player: 1, Count: 2, Options: []engine.Msg{{}}, Purpose: 1},
			WinSouls:   10,
			Over:       true,
			Winners:    []engine.PlayerID{1},
			Revealed:   &top,
			Subject:    &top,
			Ordering:   []engine.CardView{top},
			Modes:      []engine.ModeView{{Text: "a", Option: 1}},
			Passed:     []engine.PlayerID{0},
			Responding: []engine.PlayerID{1},
			Buying:     true,
			FirstRoll:  &engine.FirstRoll{Round: 1, Players: []engine.PlayerID{0, 1}, Rolls: []int{3}},
		},
		Allowed:      []engine.Intent{{Player: 1, Kind: engine.IntentPass}, {Player: 0, Kind: engine.IntentActivate, Objects: []engine.ObjectID{8, 9}, Choice: 2}},
		Aims:         []engine.Aim{{Intent: 0}},
		Blocked:      []engine.Block{{Kind: engine.IntentEndTurn}},
		Seats:        []Seat{{Seat: 0, Name: "Ann", Connected: true, Ready: true, Kicked: true, Dice: "tear"}},
		Deadline:     1760000000000,
		WaitTurn:     true,
		Pause:        &Pause{Away: []int{1}, Votes: []Vote{{Seat: 0}}, Grace: 5},
		RollbackMode: "soft",
		RollbackFrom: 3,
		Rollback:     &Rollback{Seat: 1, To: 2, Undid: 4},
		RollbackLast: &RollbackStep{Step: 6},
		RollbackVote: &RollbackVote{Step: 5},
		RollbackEnd:  &RollbackEnd{Step: 4},
		RollbackCut:  2,
	}
}

// TestEncodeFullUpdateMatchesJSON: an update with every optional field set
// encodes to the bytes encoding/json writes for it.
func TestEncodeFullUpdateMatchesJSON(t *testing.T) {
	u := fullUpdate()
	got := appendUpdate(nil, u)
	want, err := json.Marshal(u)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(want) {
		t.Errorf("full update:\n got %s\nwant %s", got, want)
	}
}

// TestEncodeEmptyUpdateMatchesJSON: the zero update, where the nil lists
// are null and the optional fields are left out.
func TestEncodeEmptyUpdateMatchesJSON(t *testing.T) {
	u := Update{}
	got := appendUpdate(nil, u)
	want, err := json.Marshal(u)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != string(want) {
		t.Errorf("empty update:\n got %s\nwant %s", got, want)
	}
}
