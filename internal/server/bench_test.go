package server

import (
	"math/rand/v2"
	"testing"

	"github.com/FourSoulsRepo/four_souls/internal/protocol"
	engine "github.com/FourSoulsRepo/rules_engine"
)

// realUpdates plays a random game of four with the load tool's bans and
// returns the update each seat gets at each step, built as Room.update
// builds them (18.17). Shared by the encode benchmarks.
func realUpdates(tb testing.TB, steps int) []protocol.Update {
	tb.Helper()
	setup := engine.Setup{
		Seed: 11, Players: 4, Sets: setsOf([]string{"b2"}), BonusSouls: true, DeckRatios: false,
		Banned: []engine.CardRef{"blue_baby", "smelter", "portable_slot_machine"},
	}
	g, events, err := engine.NewGame(setup)
	if err != nil {
		tb.Fatal(err)
	}
	rng := rand.New(rand.NewPCG(11, 11)) //nolint:gosec // a benchmark, not a secret
	var out []protocol.Update
	for step := 1; step <= steps && !g.Over; step++ {
		for seat := range 4 {
			v := engine.Viewer{Kind: engine.ViewPlayer, Player: engine.PlayerID(seat)}
			allowed := g.Allowed(engine.PlayerID(seat))
			out = append(out, protocol.Update{
				Step: step, Events: engine.FilterEvents(events, v), View: g.View(v), Allowed: allowed,
				Aims: g.Aims(allowed), Blocked: g.BlockedTurn(engine.PlayerID(seat)),
			})
		}
		// The active player with a choice acts; the first allowed intent
		// picks a random one, so the game moves on.
		var acted bool
		for seat := range 4 {
			allowed := g.Allowed(engine.PlayerID(seat))
			if len(allowed) == 0 {
				continue
			}
			in := allowed[rng.IntN(len(allowed))]
			events, err = g.Apply(in)
			if err != nil {
				tb.Fatal(err)
			}
			acted = true
			break
		}
		if !acted {
			break
		}
	}
	return out
}

// BenchmarkEncodeRealUpdate encodes the updates of real games: the cost of
// one update on the wire (18.17). The benchmark reports ns and bytes per
// update.
func BenchmarkEncodeRealUpdate(b *testing.B) {
	updates := realUpdates(b, 400)
	if len(updates) == 0 {
		b.Fatal("no updates")
	}
	var total int
	for _, u := range updates {
		msg, err := protocol.Encode(protocol.TypeUpdate, 0, u)
		if err != nil {
			b.Fatal(err)
		}
		total += len(msg)
	}
	b.ReportMetric(float64(total)/float64(len(updates)), "bytes/update")
	b.ReportAllocs()
	b.ResetTimer()
	for i := range b.N {
		if _, err := protocol.Encode(protocol.TypeUpdate, 0, updates[i%len(updates)]); err != nil {
			b.Fatal(err)
		}
	}
}
