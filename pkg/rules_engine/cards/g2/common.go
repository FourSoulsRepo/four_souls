package g2

import (
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/cardkit"
)

// The card's own lines; their words are in en.json (ADR 021).
var (
	askPutWhichEvent = engine.NewKey("ask.put_which_event_in_a")
)

// Pieces many cards share; they live in cardkit for every set.
var (
	extraLoot           = cardkit.ExtraLoot
	rechargeAtEndOfTurn = cardkit.RechargeAtEndOfTurn
	livingPlayers       = cardkit.LivingPlayers
	rechargeAnItem      = cardkit.RechargeAnItem
	discardOne          = cardkit.DiscardOne
	killAMonster        = cardkit.KillAMonster
	monstersInPlay      = cardkit.MonstersInPlay
	may                 = cardkit.May
	whenThisDealsDamage = cardkit.WhenThisDealsDamage
	whenThisTakesDamage = cardkit.WhenThisTakesDamage
	attackRollOnThis    = cardkit.AttackRollOnThis
	putBackInOrder      = cardkit.PutBackInOrder
	notAttacked         = cardkit.NotAttacked
	inShop              = cardkit.InShop
	discardTarget       = cardkit.DiscardTarget
)

// centsEqualToRoll: "that player gains ¢ equal to the result" for
// target 0.
func centsEqualToRoll() engine.RollTable {
	var t engine.RollTable
	for r := 1; r <= 6; r++ {
		t = t.Results(r, r, engine.EffectFunc(func(c *engine.Ctx) {
			them := *c
			them.Controller = c.Targets[0].Player
			them.Do(engine.GainCents(r))
		}))
	}
	return t
}

// pickAnotherPlayer asks the controller for another living player.
func pickAnotherPlayer(text engine.Key) engine.Question {
	return engine.Question{Text: text, Options: func(c *engine.Ctx, _ []int) []engine.Msg {
		var out []engine.Msg
		for _, p := range livingPlayers(c.G, c.Controller) {
			out = append(out, cardkit.PlayerLabel(c.G, p))
		}
		return out
	}}
}

// forTarget runs effects as if target 0, a player, controlled them:
// "they gain …".
func forTarget(effects ...engine.Effect) engine.Effect {
	return engine.EffectFunc(func(c *engine.Ctx) {
		them := *c
		them.Controller = c.Targets[0].Player
		them.Do(effects...)
	})
}

// rechargeYourItems: "Recharge each item you control".
var rechargeYourItems = engine.EffectFunc(func(c *engine.Ctx) {
	for _, id := range c.G.Players[c.Controller].InPlay {
		if c.G.Object(id).Role == engine.RoleItem {
			c.G.Recharge(id)
		}
	}
})

// attackedDC: "Monsters you attack have +n DC till end of turn".
func attackedDC(n int) engine.Effect {
	return engine.EffectFunc(func(c *engine.Ctx) {
		c.G.Boosts = append(c.G.Boosts, engine.Boost{Source: c.Source, Stat: engine.StatAttackedDC, Player: c.Controller, Amount: n})
	})
}

// combatDamageOnRollOf triggers when the controller takes combat damage
// from a monster on an attack roll of n; c.EventSource is the monster.
func combatDamageOnRollOf(n int) engine.Trigger {
	return engine.Trigger{On: engine.EvDamaged, Match: func(g *engine.Game, self engine.ObjectID, e engine.Event) bool {
		return e.Text == "combat" && e.Player == g.Object(self).Controller && g.Attack.LastRoll == n
	}}
}

// rechargeYourItem: "you may recharge an item you control", one at a time.
var rechargeYourItem = engine.Ask(func(c *engine.Ctx, a []int) {
	if ids := yourDeactivated(c); a[0] >= 0 && a[0] < len(ids) {
		c.G.Recharge(ids[a[0]])
	}
}, engine.Question{Text: cardkit.AskRechargeWhichItem, Options: func(c *engine.Ctx, _ []int) []engine.Msg {
	var out []engine.Msg
	for _, id := range yourDeactivated(c) {
		out = append(out, engine.CardOption(c.G.Object(id).Card))
	}
	if len(out) == 0 {
		return nil
	}
	return append(out, engine.AskNone.Msg())
}})

func yourDeactivated(c *engine.Ctx) []engine.ObjectID {
	var out []engine.ObjectID
	for _, id := range c.G.Players[c.Controller].InPlay {
		if o := c.G.Object(id); o.Role == engine.RoleItem && !o.Charged {
			out = append(out, id)
		}
	}
	return out
}

// freeSlots lists monster slots whose monster is not being attacked.
func freeSlots(g *engine.Game) []int {
	var out []int
	for i, s := range g.Monsters {
		if top, ok := s.TopOf(); !ok || !g.Attack.On || top != g.Attack.Target {
			out = append(out, i)
		}
	}
	return out
}

// coverASlot: "put the top card of the monster deck in a monster slot not
// being attacked"; optional adds a "none" option.
func coverASlot(optional bool) engine.Effect {
	return engine.Ask(func(c *engine.Ctx, a []int) {
		if slots := freeSlots(c.G); a[0] >= 0 && a[0] < len(slots) {
			c.G.CoverMonsterSlot(slots[a[0]])
		}
	}, engine.Question{Text: cardkit.AskWhichMonsterSlot, Options: func(c *engine.Ctx, _ []int) []engine.Msg {
		var out []engine.Msg
		for _, i := range freeSlots(c.G) {
			out = append(out, engine.SlotOption(i+1, ""))
		}
		if optional {
			out = append(out, engine.AskNone.Msg())
		}
		return out
	}})
}

// lookAtADeck is "Look at the top card of a deck".
var lookAtADeck = engine.Ask(func(c *engine.Ctx, a []int) {
	if a[0] >= 0 {
		c.G.LookAt(c.Controller, c.G.DeckTop(c.Deck(a[0]), 1)...)
	}
}, engine.DeckQuestion(cardkit.AskLookTop))

// targetsYours: a stack ability or loot that targets an item or a dice
// roll the chooser controls.
func targetsYours(g *engine.Game, self engine.ObjectID, c engine.Chosen) bool {
	it, ok := g.StackItemByID(c.StackID)
	if !ok {
		return false
	}
	me := g.Object(self).Controller
	for _, t := range it.Targets {
		if roll, ok := g.StackItemByID(t.StackID); ok && t.StackID != 0 && roll.Kind == engine.StackRoll && roll.Controller == me {
			return true
		}
		if t.StackID == 0 && t.Object != 0 && g.Object(t.Object).Role == engine.RoleItem && g.Object(t.Object).Controller == me {
			return true
		}
	}
	return false
}

// noCombatDamageTakenOn is "You take no combat damage on attack rolls of n".
func noCombatDamageTakenOn(n int) func(*engine.Game, engine.ObjectID, int, bool, int) int {
	return func(_ *engine.Game, _ engine.ObjectID, roll int, hit bool, dmg int) int {
		if !hit && roll == n {
			return 0
		}
		return dmg
	}
}

// combatDamageToAMonster triggers when the controller deals combat damage
// to a monster, on an attack roll of roll (0: any roll).
func combatDamageToAMonster(roll int) engine.Trigger {
	return engine.Trigger{On: engine.EvDamaged, Match: func(g *engine.Game, self engine.ObjectID, e engine.Event) bool {
		return e.Text == "combat" && e.Player == engine.NoPlayer && e.Object != 0 &&
			g.Turn.Active == g.Object(self).Controller && (roll == 0 || g.Attack.LastRoll == roll)
	}}
}

// damageOtherMonsters: "deal 1 damage to each other monster" than the one
// that was damaged.
var damageOtherMonsters = engine.EffectFunc(func(c *engine.Ctx) {
	for _, id := range monstersInPlay(c.G) {
		if id != c.EventObject {
			c.G.DealDamageTo(engine.Target{Object: id}, 1, c.Controller, c.Source)
		}
	}
})

// namedMom: Mom!, Mom's Heart! or It Lives!
func namedMom(g *engine.Game, _ engine.ObjectID, c engine.Chosen) bool {
	switch g.Object(c.Object).Card {
	case "mom", "moms_heart", "it_lives":
		return true
	}
	return false
}

// ouija: "Reveal the top 4 cards of the monster deck. You may put an event
// card revealed this way in a monster slot not being attacked. Put the
// rest on the bottom of the deck in a random order".
var ouija = engine.Ask(func(c *engine.Ctx, a []int) {
	top := c.G.DeckTop(engine.MonsterDeck, 4)
	for _, pl := range c.G.Players {
		c.G.LookAt(pl.ID, top...) // revealed to everyone
	}
	events := ouijaEvents(c.G, top)
	placed := engine.ObjectID(-1)
	if slots := freeSlots(c.G); a[0] >= 0 && a[0] < len(events) && a[1] >= 0 && a[1] < len(slots) {
		placed = events[a[0]]
		c.G.PlaceFromDeck(placed, slots[a[1]])
	}
	var rest []engine.ObjectID
	for _, id := range top {
		if id != placed {
			rest = append(rest, id)
		}
	}
	c.G.ToBottomRandom(engine.MonsterDeck, rest)
}, engine.Question{Text: askPutWhichEvent, Options: func(c *engine.Ctx, _ []int) []engine.Msg {
	var out []engine.Msg
	for _, id := range ouijaEvents(c.G, c.G.DeckTop(engine.MonsterDeck, 4)) {
		out = append(out, engine.CardOption(c.G.Object(id).Card))
	}
	if len(out) == 0 {
		return nil
	}
	return append(out, engine.AskNone.Msg())
}}, engine.Question{Text: cardkit.AskWhichMonsterSlot, Options: func(c *engine.Ctx, a []int) []engine.Msg {
	if a[0] < 0 || a[0] >= len(ouijaEvents(c.G, c.G.DeckTop(engine.MonsterDeck, 4))) {
		return nil
	}
	var out []engine.Msg
	for _, i := range freeSlots(c.G) {
		out = append(out, engine.SlotOption(i+1, ""))
	}
	return out
}})

func ouijaEvents(g *engine.Game, ids []engine.ObjectID) []engine.ObjectID {
	var out []engine.ObjectID
	for _, id := range ids {
		if g.Kind(id) == engine.EventCard {
			out = append(out, id)
		}
	}
	return out
}

// healThis heals this monster by n.
func healThis(n int) engine.Effect {
	return engine.EffectFunc(func(c *engine.Ctx) { c.G.HealObject(c.Source, n) })
}

// thisGains gives this monster +n of a stat till end of turn.
func thisGains(stat engine.Stat, n int) engine.Effect {
	return engine.EffectFunc(func(c *engine.Ctx) {
		c.G.Boosts = append(c.G.Boosts, engine.Boost{Source: c.Source, Stat: stat, Player: engine.NoPlayer, Object: c.Source, Amount: n})
	})
}

// backOnTop puts the dead monster on top of the monster deck.
var backOnTop = engine.EffectFunc(func(c *engine.Ctx) {
	if c.G.Object(c.EventObject).Zone.Kind == engine.ZoneOutside {
		c.G.PutIntoDeck(engine.MonsterDeck, c.EventObject, 0)
	}
})
