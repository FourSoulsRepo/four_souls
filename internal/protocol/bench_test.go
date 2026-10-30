package protocol

import (
	"testing"

	engine "github.com/FourSoulsRepo/rules_engine"
)

// BenchmarkEncodeUpdate encodes an update with 40 events, each with a
// message of three values: the shape of a busy step (18.17).
func BenchmarkEncodeUpdate(b *testing.B) {
	u := benchUpdate()
	b.ReportAllocs()
	for b.Loop() {
		if _, err := Encode(TypeUpdate, 0, u); err != nil {
			b.Fatal(err)
		}
	}
}

// benchUpdate is the update BenchmarkEncodeUpdate encodes.
func benchUpdate() Update {
	events := make([]engine.Event, 40)
	for i := range events {
		events[i] = engine.Event{
			Player: engine.PlayerID(i % 4),
			Object: engine.ObjectID(i + 1),
			Amount: i,
			Label: engine.Msg{Args: engine.Args{
				engine.N("n", i),
				engine.Str("text", "a soul is gained"),
				engine.PlayerArg("player", engine.PlayerID(i%4)),
			}},
		}
	}
	return Update{Step: 1, Events: events}
}
