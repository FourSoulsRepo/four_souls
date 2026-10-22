package protocol

import (
	"bytes"
	"encoding/json"
	"flag"
	"os"
	"strings"
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/b2"
)

var update = flag.Bool("update", false, "rewrite testdata/golden.jsonl")

func ptr[T any](v T) *T { return &v }

// TestGolden pins the JSON of every message (ADR 006). If it fails, the
// wire format changed: if a field was renamed or removed, bump Version;
// then run with -update.
func TestGolden(t *testing.T) {
	g, _, err := engine.NewGame(engine.Setup{Seed: 1, Players: 2, Sets: []engine.CardSet{b2.Set}, Characters: []engine.CardRef{"isaac", "cain"}})
	if err != nil {
		t.Fatal(err)
	}
	p := g.Prompt().Player
	view := g.View(engine.Viewer{Kind: engine.ViewPlayer, Player: p})
	messages := []struct {
		typ  string
		id   int
		data any
	}{
		{TypeHello, 1, Hello{Protocol: Version, App: "0.1.0", Name: "Ann", Role: RolePlayer, Sets: []string{"b2"}, Dice: "tear"}},
		// A fixed sample: an engine version bump is not a wire change.
		{TypeWelcome, 1, Welcome{Protocol: Version, App: "0.1.0", Engine: "0.1.0", Token: "t0k3n", Seat: 0, Presence: true, Name: "Basement"}},
		{TypeIntent, 2, Intent{Intent: engine.Intent{Kind: engine.IntentPass}}},
		{TypeResync, 3, struct{}{}},
		{TypeUpdate, 0, Update{Step: 1, Events: []engine.Event{{Kind: engine.EvLooted, Player: 0, Card: "a_dime"}}, View: view, Allowed: g.Allowed(p), Aims: []engine.Aim{{Intent: 0, Targets: []engine.OptionView{{Object: 320, Player: engine.NoPlayer}, {Player: 1, IsPlayer: true}}, Reason: "Pick a monster or a player."}}, Blocked: []engine.Block{{Kind: engine.IntentAttack, Reason: "Not your turn"}, {Kind: engine.IntentPlayLoot, Object: 314, Reason: "Loot is played on your turn"}}, Seats: []Seat{{Seat: 0, Name: "Ann", Connected: true, Dice: "tear"}, {Seat: 1, Name: "Bo"}}}},
		{TypeError, 2, Error{Code: ErrRefused, Message: "no loot play available", Rule: "R-CARD-08"}},
		{TypeList, 4, struct{}{}},
		{TypeGames, 4, Games{Games: []GameInfo{{ID: "g1", Host: "Ann", Seats: 4, Taken: 2, Sets: []string{"b2"}}}}},
		{TypeCreate, 5, Create{Seats: 3, Sets: []string{"b2"}, Options: &Options{Picking: PickDraft, DraftSize: 3, BanRounds: 1, BanTimer: 30, HostBans: []engine.CardRef{"eden"}, ResponseTimer: 90}}},
		{TypeJoin, 6, Join{Game: "g1"}},
		{TypeReady, 7, Ready{Ready: true}},
		{TypeLeave, 8, struct{}{}},
		{TypeTable, 0, Table{Game: "g1", You: 1, Seats: []Seat{{Seat: 0, Name: "Ann", Connected: true, Ready: true, Dice: "tear"}, {Seat: 1, Name: "Bo", Connected: true}}, Sets: []string{"b2"}, Options: &Options{Picking: PickRandom}}},
		{TypeSetup, 0, Setup{Phase: PhaseBan, Round: 1, Turn: 0, Pool: []engine.CardRef{"isaac", "cain"}, Banned: []Ban{{Seat: -1, Card: "eden"}}, Deadline: 1760000000000, Seats: []Seat{{Seat: 0, Name: "Ann"}}}},
		{TypeBan, 9, BanCard{Card: "cain"}},
		{TypePick, 10, PickCard{Card: "isaac"}},
		{TypeVote, 12, Vote{Kick: true}},
		{TypeSave, 13, struct{}{}},
		{TypeSaved, 0, Saved{Save: "g1"}},
		{TypeLoad, 14, Load{Save: "g1"}},
		{TypeWaitTurn, 15, WaitTurn{On: true}},
		{TypeIntent, 16, Intent{Intent: engine.Intent{Kind: engine.IntentThrow}}},
		{TypePresence, 17, Presence{From: &Spot{Hand: 2}, Point: &Point{X: 0.25, Y: 0.5}}},
		{TypePresence, 0, Presence{Seat: 1, From: &Spot{Object: 314}, At: &Spot{Player: ptr(engine.PlayerID(0))}}},
		{TypePresence, 0, Presence{Seat: 1, At: &Spot{Deck: ptr(engine.MonsterDeck)}}},
		{TypePresence, 0, Presence{Seat: 1, At: &Spot{StackID: 3}}},
	}
	var got bytes.Buffer
	for _, m := range messages {
		b, encErr := Encode(m.typ, m.id, m.data)
		if encErr != nil {
			t.Fatal(encErr)
		}
		got.Write(b)
		got.WriteByte('\n')
	}
	const file = "testdata/golden.jsonl"
	if *update {
		if werr := os.WriteFile(file, got.Bytes(), 0o600); werr != nil {
			t.Fatal(werr)
		}
	}
	want, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(got.Bytes(), want) {
		t.Errorf("the wire format changed (ADR 006): bump protocol.Version unless fields were only added, then run go test ./internal/protocol -update")
	}
}

func TestRoundTrip(t *testing.T) {
	b, err := Encode(TypeHello, 5, Hello{Protocol: Version, Name: "Bo", Role: RoleSpectator})
	if err != nil {
		t.Fatal(err)
	}
	e, err := Decode(b)
	if err != nil || e.Type != TypeHello || e.ID != 5 {
		t.Fatalf("decode: %+v, %v", e, err)
	}
	var h Hello
	if err := e.Unpack(&h); err != nil || h.Name != "Bo" || h.Role != RoleSpectator {
		t.Fatalf("unpack: %+v, %v", h, err)
	}
	if _, err := Decode([]byte(`{"id":1}`)); err == nil {
		t.Error("a message without a type was accepted")
	}
	if _, err := Decode([]byte(`not json`)); err == nil {
		t.Error("bad JSON was accepted")
	}
	var raw json.RawMessage = []byte(`{"name": 7}`)
	if err := (Envelope{Type: TypeHello, Data: raw}).Unpack(&h); err == nil {
		t.Error("a wrong field type was accepted")
	}
}

func TestCheckVersion(t *testing.T) {
	if CheckVersion(Version) != "" || CheckVersion(Version-1) != ErrClientOutdated || CheckVersion(Version+1) != ErrServerOutdated {
		t.Error("CheckVersion")
	}
}

func TestNames(t *testing.T) {
	for in, want := range map[string]string{"  Ann ": "Ann", "Олег": "Олег"} {
		if got, ok := CleanName(in); !ok || got != want {
			t.Errorf("CleanName(%q) = %q, %v", in, got, ok)
		}
	}
	for _, bad := range []string{"", "   ", "a\x07b", "ThisNameIsWayTooLongForUs"} {
		if _, ok := CleanName(bad); ok {
			t.Errorf("CleanName(%q) accepted", bad)
		}
	}
	// A server name may be longer than a nickname (N-18).
	if got, ok := CleanServerName(" ThisNameIsWayTooLongForUs "); !ok || got != "ThisNameIsWayTooLongForUs" {
		t.Errorf("CleanServerName = %q, %v", got, ok)
	}
	if _, ok := CleanServerName(strings.Repeat("x", MaxServerName+1)); ok {
		t.Error("CleanServerName accepted a too long name")
	}
	if got := UniqueName("Ann", []string{"Ann", "Ann (2)"}); got != "Ann (3)" {
		t.Errorf("UniqueName = %q", got)
	}
}

func TestCleanDice(t *testing.T) {
	for id, want := range map[string]string{
		"tear":                  "tear",
		"holy_set-2":            "holy_set-2",
		"":                      "",
		"Tear":                  "",
		"../x":                  "",
		strings.Repeat("a", 33): "",
	} {
		if got := CleanDice(id); got != want {
			t.Errorf("CleanDice(%q) = %q, want %q", id, got, want)
		}
	}
}
