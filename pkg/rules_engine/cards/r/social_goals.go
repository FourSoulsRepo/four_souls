package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Social Goals! (Room Card)
//
//	Note each goal as players complete them.
//	This room can't be put into discard till 4 goals are completed.
//	1. Play 5 loot cards.
//	2. Kill 3 monsters.
//	3. Give at least 6¢ to another player at one time.
//	4. Purchase 3 items.
//	5. Roll a 6 three times.
//	When 4 goals are completed, each player gains +2 treasure.
var socialGoals = engine.CardDef{
	Ref:       "social_goals",
	Kind:      engine.RoomCard,
	Copies:    1,
	NoDiscard: func(g *engine.Game, self engine.ObjectID) bool { return g.Object(self).CountersOf(goalsDone) < 4 },
	Abilities: []engine.Ability{
		goal("1. Play 5 loot cards.", 5, engine.Trigger{On: engine.EvLootPlayed, Match: anyPlayer}),
		goal("2. Kill 3 monsters.", 3, engine.Trigger{On: engine.EvDied, Match: func(g *engine.Game, _ engine.ObjectID, e engine.Event) bool {
			return e.Player == engine.NoPlayer && g.Object(e.Object).Killer() != engine.NoPlayer
		}}),
		goal("3. Give at least 6¢ to another player at one time.", 1, engine.Trigger{On: engine.EvGaveCents, Match: func(_ *engine.Game, _ engine.ObjectID, e engine.Event) bool {
			return e.Amount >= 6
		}}),
		goal("4. Purchase 3 items.", 3, engine.Trigger{On: engine.EvPurchased, Match: anyPlayer}),
		goal("5. Roll a 6 three times.", 3, engine.Trigger{On: engine.EvRollResolved, Match: func(g *engine.Game, self engine.ObjectID, e engine.Event) bool {
			return e.Amount == 6 && anyPlayer(g, self, e)
		}}),
	},
}

// goalsDone names the counter of completed goals; each goal counts its
// progress in its own named counter. All players add to one tally while
// the room is in play (Q-20).
const goalsDone = "goals"

// goal counts one Social Goals goal; when the 4th goal is completed,
// each player gains +2 treasure.
func goal(text string, n int, on engine.Trigger) engine.Ability {
	return engine.Ability{
		Kind:    engine.Triggered,
		Text:    text,
		Trigger: on,
		Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
			o := c.G.Object(c.Source)
			if o.CountersOf(text) >= n {
				return // done
			}
			c.G.AddNamedCounters(c.Source, text, 1)
			if o.CountersOf(text) < n {
				return
			}
			c.G.AddNamedCounters(c.Source, goalsDone, 1)
			if o.CountersOf(goalsDone) == 4 {
				c.Do(engine.EachPlayer(engine.GainTreasure(2)))
			}
		})},
	}
}
