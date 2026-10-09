package protocol

import (
	"bytes"
	"encoding/json"
	"flag"
	"os"
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards"
)

var update = flag.Bool("update", false, "rewrite testdata/golden.jsonl")

// TestGolden pins the JSON of every message (ADR 006). If it fails, the
// wire format changed: bump Version, then run with -update.
func TestGolden(t *testing.T) {
	g, _, err := engine.NewGame(engine.Setup{Seed: 1, Players: 2, Sets: cards.Sets(), Characters: []engine.CardRef{"isaac", "cain"}})
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
		{TypeHello, 1, Hello{Protocol: Version, App: "0.1.0", Name: "Ann", Role: RolePlayer, Sets: []string{"b2"}}},
		{TypeWelcome, 1, Welcome{Protocol: Version, App: "0.1.0", Engine: engine.Version, Token: "t0k3n", Seat: 0}},
		{TypeIntent, 2, Intent{Intent: engine.Intent{Kind: engine.IntentPass}}},
		{TypeResync, 3, struct{}{}},
		{TypeUpdate, 0, Update{Step: 1, Events: []engine.Event{{Kind: engine.EvLooted, Player: 0, Card: "a_dime"}}, View: view, Allowed: g.Allowed(p)}},
		{TypeError, 2, Error{Code: ErrRefused, Message: "no loot play available", Rule: "R-CARD-08"}},
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
		t.Errorf("the wire format changed (ADR 006): bump protocol.Version, then run go test ./internal/protocol -update")
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
