package rulesengine

// Counters on players (R-MECH-13, Q-45): they are on the player, not
// the character, so they stay when the player dies or the character
// changes. Trick or Treat's fear counters are the first (R-TRICK-07);
// any card may name its own.

// EvPlayerCounters means Player got Amount counters named Text,
// negative if removed (R-MECH-13).
const EvPlayerCounters EventKind = "player_counters"

// PlayerCounters returns how many counters named name player p has.
func (g *Game) PlayerCounters(p PlayerID, name string) int {
	return countOf(g.Players[p].Counters, name)
}

// AddPlayerCounters puts n counters named name on player p, or removes
// -n; never below 0. It returns the change made; a change emits
// player_counters.
func (g *Game) AddPlayerCounters(p PlayerID, name string, n int) int {
	pl := &g.Players[p]
	before := countOf(pl.Counters, name)
	pl.Counters = addCount(pl.Counters, name, n)
	d := countOf(pl.Counters, name) - before
	if d != 0 {
		g.emit(Event{Kind: EvPlayerCounters, Player: p, Amount: d, Text: name})
	}
	return d
}

// PlayersWithCounters lists, in seat order, the players with at least
// one counter named name: the options of "remove a fear counter from a
// player" (Q-51).
func (g *Game) PlayersWithCounters(name string) []PlayerID {
	var out []PlayerID
	for _, pl := range g.Players {
		if countOf(pl.Counters, name) > 0 {
			out = append(out, pl.ID)
		}
	}
	return out
}
