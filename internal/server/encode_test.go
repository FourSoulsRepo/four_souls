package server

import (
	"encoding/json"
	"testing"

	"github.com/FourSoulsRepo/four_souls/internal/protocol"
)

// referenceEncode is the encoding before the hand-written encoders: the
// data by encoding/json, the envelope spliced around it (18.17).
func referenceEncode(t *testing.T, typ string, id int, data any) []byte {
	t.Helper()
	raw, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}
	head, err := json.Marshal(protocol.Envelope{Type: typ, ID: id})
	if err != nil {
		t.Fatal(err)
	}
	out := append(head[:len(head)-1], `,"data":`...)
	out = append(out, raw...)
	return append(out, '}')
}

// TestEncodeUpdateMatchesJSON: every update of a real game encodes to the
// bytes encoding/json writes for it, so the wire format is unchanged.
func TestEncodeUpdateMatchesJSON(t *testing.T) {
	for i, u := range realUpdates(t, 400) {
		got, err := protocol.Encode(protocol.TypeUpdate, 0, u)
		if err != nil {
			t.Fatal(err)
		}
		want := referenceEncode(t, protocol.TypeUpdate, 0, u)
		if string(got) != string(want) {
			t.Fatalf("update %d differs:\n got %s\nwant %s", i, got, want)
		}
	}
}
