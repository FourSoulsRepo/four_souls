package r

import engine "github.com/FourSoulsRepo/rules_engine"

// Helpers of the room cards. "You" on a room is the active player, who
// controls its abilities (R-CARD-36).

// playersHave changes a stat of every player: "Players have +1 HP".
func playersHave(s engine.Stat, n int) engine.Static {
	return playersWith(s, n, func(*engine.Game, engine.PlayerID) bool { return true })
}

// activePlayerHas changes a stat of the active player.
func activePlayerHas(s engine.Stat, n int) engine.Static {
	return playersWith(s, n, func(g *engine.Game, p engine.PlayerID) bool { return p == g.Turn.Active })
}

// playersWith changes a stat of the players who match.
func playersWith(s engine.Stat, n int, match func(g *engine.Game, p engine.PlayerID) bool) engine.Static {
	return engine.Static{Stat: s, Amount: n, Applies: func(g *engine.Game, _ engine.ObjectID, p engine.PlayerID, _ engine.ObjectID) bool {
		return p != engine.NoPlayer && match(g, p)
	}}
}

// itemCount and soulCount count what a player controls.
func itemCount(g *engine.Game, p engine.PlayerID) int { return len(itemsOf(g, p)) }
func soulCount(g *engine.Game, p engine.PlayerID) int { return g.SoulValue(p) }

// mostOf: p has the most, or tied for the most.
func mostOf(count func(*engine.Game, engine.PlayerID) int) func(*engine.Game, engine.PlayerID) bool {
	return func(g *engine.Game, p engine.PlayerID) bool {
		for _, pl := range g.Players {
			if count(g, pl.ID) > count(g, p) {
				return false
			}
		}
		return true
	}
}

// fewestOf: p has the fewest, or tied for the fewest.
func fewestOf(count func(*engine.Game, engine.PlayerID) int) func(*engine.Game, engine.PlayerID) bool {
	return func(g *engine.Game, p engine.PlayerID) bool {
		for _, pl := range g.Players {
			if count(g, pl.ID) < count(g, p) {
				return false
			}
		}
		return true
	}
}

// atStartOfTheTurn and atEndOfTheTurn trigger on every turn.
func atStartOfTheTurn() engine.Trigger {
	return engine.Trigger{On: engine.EvStartOfTurn, Match: anyPlayer}
}

func atEndOfTheTurn() engine.Trigger {
	return engine.Trigger{On: engine.EvEndOfTurn, Match: anyPlayer}
}

// anyPlayer matches an event about any player.
func anyPlayer(_ *engine.Game, _ engine.ObjectID, e engine.Event) bool {
	return e.Player != engine.NoPlayer
}

// eachOtherThanDead runs the effects for each player but the one who
// died, in turn order from the active player.
func eachOtherThanDead(effects ...engine.Effect) engine.Effect {
	return engine.EffectFunc(func(c *engine.Ctx) {
		for _, p := range c.G.TurnOrderFrom(c.G.Turn.Active) {
			if p != c.EventPlayer {
				them := *c
				them.Controller = p
				them.Do(effects...)
			}
		}
	})
}

// theyTake: the event's player takes n damage from this.
func theyTake(n int) engine.Effect {
	return engine.EffectFunc(func(c *engine.Ctx) {
		c.G.DealDamageTo(engine.Target{Player: c.EventPlayer, IsPlayer: true}, n, engine.NoPlayer, c.Source)
	})
}

// youTake: the controller takes n damage from this.
func youTake(n int) engine.Effect {
	return engine.EffectFunc(func(c *engine.Ctx) {
		c.G.DealDamageTo(engine.Target{Player: c.Controller, IsPlayer: true}, n, engine.NoPlayer, c.Source)
	})
}

// discardThisRoom is "Put this into discard".
var discardThisRoom = engine.EffectFunc(func(c *engine.Ctx) { c.G.DiscardRoom(c.Source) })

// diceRoom is a Dice Room: the "enters play" effects, then "At the end
// of the turn, put this into discard." That discard is not the room
// step (R-TURN-17).
func diceRoom(text string, effects ...engine.Effect) []engine.Ability {
	return []engine.Ability{
		{Kind: engine.Triggered, Text: text, Trigger: engine.WhenThisEntersPlay(), Effects: effects},
		{
			Kind:    engine.Triggered,
			Text:    "At the end of the turn, put this into discard.",
			Trigger: atEndOfTheTurn(),
			Effects: []engine.Effect{discardThisRoom},
		},
	}
}

// diceRoomTakes is diceRoom for a room that takes the players' hands
// (Ability.Takes, 15.85).
func diceRoomTakes(text string, effects ...engine.Effect) []engine.Ability {
	abs := diceRoom(text, effects...)
	abs[0].Takes = engine.TakesOthers
	return abs
}

// rerollAll rerolls each of the items; each reroll is its own gain
// (S-DIS: Damocles, Dice Room).
func rerollAll(g *engine.Game, ids []engine.ObjectID) {
	for _, id := range ids {
		g.RerollItem(id)
	}
}

// shopItems lists the shop items.
func shopItems(g *engine.Game) []engine.ObjectID {
	var out []engine.ObjectID
	for _, s := range g.Shop {
		if top, ok := s.TopOf(); ok {
			out = append(out, top)
		}
	}
	return out
}

// discardShop puts each shop item into discard.
var discardShop = engine.EffectFunc(func(c *engine.Ctx) {
	for _, id := range shopItems(c.G) {
		c.G.DiscardShopItem(id)
	}
})

// eachMonsterToDiscard puts each monster into discard; with notAttacked, not
// the one being attacked. The slots are refilled.
func eachMonsterToDiscard(notAttacked bool) engine.Effect {
	return engine.EffectFunc(func(c *engine.Ctx) {
		for _, id := range monstersInPlay(c.G) {
			if !notAttacked || !c.G.Attack.On || c.G.Attack.Target != id {
				c.G.DiscardMonster(id)
			}
		}
		c.G.RefillSlots()
	})
}

// monsterOrRoomDecks are the options of "the monster or room deck".
func monsterOrRoomDecks(g *engine.Game) []engine.DeckKind {
	out := []engine.DeckKind{engine.MonsterDeck}
	if len(g.Rooms) > 0 {
		out = append(out, engine.RoomDeck)
	}
	return out
}

// deckLabels names decks as the prompts do.
func deckLabels(ds []engine.DeckKind) []string {
	out := make([]string, len(ds))
	for i, d := range ds {
		out[i] = d.String() + " deck"
	}
	return out
}

// roomsAndMonstersToDiscard lists the rooms and the monsters not being
// attacked that may be put into discard.
func roomsAndMonstersToDiscard(g *engine.Game) []engine.ObjectID {
	var out []engine.ObjectID
	for _, id := range monstersInPlay(g) {
		if !g.Attack.On || g.Attack.Target != id {
			out = append(out, id)
		}
	}
	for _, s := range g.Rooms {
		if top, ok := s.TopOf(); ok && g.Object(top).Role == engine.RoleRoom {
			if stay := g.Def(top).NoDiscard; stay == nil || !stay(g, top) {
				out = append(out, top)
			}
		}
	}
	return out
}

// roomOrMonsterToDiscard: "put a room or monster not being attacked into
// discard"; with may, "none" is an option.
func roomOrMonsterToDiscard(may bool) engine.Effect {
	return engine.Ask(func(c *engine.Ctx, a []int) {
		ids := roomsAndMonstersToDiscard(c.G)
		if a[0] < 0 || a[0] >= len(ids) {
			return
		}
		if c.G.Object(ids[a[0]]).Role == engine.RoleRoom {
			c.G.DiscardRoom(ids[a[0]])
			return
		}
		c.G.DiscardMonster(ids[a[0]])
		c.G.RefillSlots()
	}, engine.Question{Text: "Put which into discard?", Options: func(c *engine.Ctx, _ []int) []string {
		out := labelsOf(c.G, roomsAndMonstersToDiscard(c.G))
		if may {
			out = append(out, "none")
		}
		return out
	}})
}
