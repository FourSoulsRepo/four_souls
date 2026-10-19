package fsp2

import (
	"slices"
	"strconv"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/cardkit"
)

// Pieces many cards share; they live in cardkit for every set.
var (
	extraLoot      = cardkit.ExtraLoot
	discardOne     = cardkit.DiscardOne
	rechargeAnItem = cardkit.RechargeAnItem
	damageAPlayer  = cardkit.DamageAPlayer
	damageAMonster = cardkit.DamageAMonster
	damageAnything = cardkit.DamageAnything
	changeRollBy   = cardkit.ChangeRollBy
	forcedDiscard  = cardkit.ForcedDiscard
	giveThisAway   = cardkit.GiveThisAway
	may            = cardkit.May
	monstersInPlay = cardkit.MonstersInPlay
	livingPlayers  = cardkit.LivingPlayers
	playerLabel    = cardkit.PlayerLabel

	attackRollOnThis    = cardkit.AttackRollOnThis
	whenThisTakesDamage = cardkit.WhenThisTakesDamage
	whenThisDealsDamage = cardkit.WhenThisDealsDamage
	forRoller           = cardkit.ForRoller
	pickAPlayer         = cardkit.PickAPlayer
	putBackInOrder      = cardkit.PutBackInOrder

	forcedDiscardQuestions = cardkit.ForcedDiscardQuestions
	onYourAttackRollOf     = cardkit.OnYourAttackRollOf
	itemsOf                = cardkit.ItemsOf
	freeSlots              = cardkit.FreeSlots
)

// whenAnotherPlayerDies triggers each time another player dies.
func whenAnotherPlayerDies() engine.Trigger {
	return engine.Trigger{On: engine.EvDied, Match: func(g *engine.Game, self engine.ObjectID, e engine.Event) bool {
		return e.Player != engine.NoPlayer && e.Player != g.Object(self).Controller
	}}
}

// reorderADeck is "Look at the top n cards of a deck and put them back in
// any order": the deck first, then one order question (19.16).
func reorderADeck(n int) engine.Effect {
	return engine.Ask(func(c *engine.Ctx, a []int) {
		if a[0] < 0 {
			return
		}
		c.G.SetDeckTop(c.Deck(a[0]), c.Order)
	}, engine.DeckQuestion("Look at the top cards of which deck?"), engine.Question{
		Text: cardkit.PutBackText,
		Order: func(c *engine.Ctx, a []int) []engine.ObjectID {
			if a[0] < 0 {
				return nil
			}
			return c.G.DeckTop(c.Deck(a[0]), n)
		},
	})
}

// rerollAnItem is "Reroll an item" (yours only, if set), picked on
// resolution; eternal items can't be rerolled (R-MECH-48).
func rerollAnItem(yours bool) engine.Effect {
	items := func(c *engine.Ctx) []engine.ObjectID {
		p := engine.NoPlayer
		if yours {
			p = c.Controller
		}
		var out []engine.ObjectID
		for _, id := range itemsOf(c.G, p) {
			if !c.G.Eternal(id) {
				out = append(out, id)
			}
		}
		return out
	}
	return engine.Ask(func(c *engine.Ctx, a []int) {
		if ids := items(c); a[0] >= 0 && a[0] < len(ids) {
			c.G.RerollItem(ids[a[0]])
		}
	}, engine.Question{Text: "Reroll which item?", Options: func(c *engine.Ctx, _ []int) []string {
		var out []string
		for _, id := range items(c) {
			out = append(out, string(c.G.Object(id).Card))
		}
		return out
	}})
}

// rerollYourItems is "Reroll each item you control".
var rerollYourItems = engine.EffectFunc(func(c *engine.Ctx) {
	for _, id := range itemsOf(c.G, c.Controller) {
		if !c.G.Eternal(id) {
			c.G.RerollItem(id)
		}
	}
})

// discardHandThenLoot is "Discard your hand, then loot n".
func discardHandThenLoot(n int) engine.Effect {
	return engine.EffectFunc(func(c *engine.Ctx) {
		for _, id := range append([]engine.ObjectID(nil), c.G.Players[c.Controller].Hand...) {
			c.G.DiscardFromHand(c.Controller, id)
		}
		c.Do(engine.Loot(n))
	})
}

// cancelALoot is "Cancel the effect of a loot being played", picked on
// resolution among the loot cards on the stack.
var cancelALoot = engine.Ask(func(c *engine.Ctx, a []int) {
	if ids := lootsOnStack(c); a[0] >= 0 && a[0] < len(ids) {
		c.G.CancelStackItem(ids[a[0]])
	}
}, engine.Question{Text: "Cancel which loot card?", Options: func(c *engine.Ctx, _ []int) []string {
	var out []string
	for _, id := range lootsOnStack(c) {
		it, _ := c.G.StackItemByID(id)
		out = append(out, string(it.Card))
	}
	return out
}})

func lootsOnStack(c *engine.Ctx) []int {
	var out []int
	for _, it := range c.G.Stack {
		if it.Kind == engine.StackLoot && it.Source != c.Source {
			out = append(out, it.ID)
		}
	}
	return out
}

// whenYouMiss triggers when the controller's attack roll misses.
func whenYouMiss() engine.Trigger {
	return engine.Trigger{On: engine.EvRollResolved, Match: func(g *engine.Game, self engine.ObjectID, e engine.Event) bool {
		return e.Text == "attack" && e.Player == g.Object(self).Controller && g.Attack.Started && e.Amount < g.Evasion(g.Attack.Target)
	}}
}

// isAttackRoll: a dice roll target that is an attack roll.
func isAttackRoll(g *engine.Game, _ engine.ObjectID, c engine.Chosen) bool {
	it, ok := g.StackItemByID(c.StackID)
	return ok && it.Attack
}

// killYou is "die!" for the controller.
var killYou = engine.EffectFunc(func(c *engine.Ctx) {
	if !c.G.Players[c.Controller].Dead {
		them := *c
		them.Targets = []engine.Chosen{{Kind: engine.TargetPlayer, Player: c.Controller}}
		them.Do(engine.Kill(0))
	}
})

// monstersDCOnYourTurn is "Monsters have +1 DC on your turn".
var monstersDCOnYourTurn = engine.Static{Stat: engine.StatMonsterDC, Amount: 1, Applies: func(g *engine.Game, self engine.ObjectID, _ engine.PlayerID, _ engine.ObjectID) bool {
	return g.Turn.Active == g.Object(self).Controller
}}

// whenYouDealDamage triggers each time the controller deals damage to a
// monster or another player.
func whenYouDealDamage() engine.Trigger {
	return engine.Trigger{On: engine.EvDamaged, Match: func(g *engine.Game, self engine.ObjectID, e engine.Event) bool {
		me := g.Object(self).Controller
		return e.By == int(me)+1 && e.Player != me // not "take damage" from your own effects
	}}
}

// whenYouKill triggers when a monster or player the controller killed dies.
func whenYouKill() engine.Trigger {
	return engine.Trigger{On: engine.EvDied, Match: func(g *engine.Game, self engine.ObjectID, e engine.Event) bool {
		me := g.Object(self).Controller
		if e.Player != engine.NoPlayer {
			return e.Player != me && g.Players[e.Player].KilledBy == int(me)+1
		}
		return e.Prev != 0 && g.Object(e.Prev).KilledBy == int(me)+1
	}}
}

// combatDamageToAMonster triggers when the controller's attack deals
// combat damage to a monster; c.EventObject is the monster.
func combatDamageToAMonster() engine.Trigger {
	return engine.Trigger{On: engine.EvDamaged, Match: func(g *engine.Game, self engine.ObjectID, e engine.Event) bool {
		return e.Text == "combat" && e.Player == engine.NoPlayer && e.Object != 0 && g.Turn.Active == g.Object(self).Controller
	}}
}

// damageAnotherThan deals 1 damage to a monster or player other than the
// monster that was just damaged.
var damageAnotherThan = engine.Ask(func(c *engine.Ctx, a []int) {
	if ts := otherTargets(c); a[0] >= 0 && a[0] < len(ts) {
		c.G.DealDamageTo(ts[a[0]], 1, c.Controller, c.Source)
	}
}, engine.Question{Text: "Deal 1 damage to?", Options: func(c *engine.Ctx, _ []int) []string {
	var out []string
	for _, t := range otherTargets(c) {
		if t.IsPlayer {
			out = append(out, playerLabel(c.G, t.Player))
		} else {
			out = append(out, string(c.G.Object(t.Object).Card))
		}
	}
	return out
}})

func otherTargets(c *engine.Ctx) []engine.Target {
	var out []engine.Target
	for _, id := range monstersInPlay(c.G) {
		if id != c.EventObject {
			out = append(out, engine.Target{Object: id})
		}
	}
	for _, p := range livingPlayers(c.G, engine.NoPlayer) {
		out = append(out, engine.Target{Player: p, IsPlayer: true})
	}
	return out
}

// eachPlayerDestroysASoul: each player, in turn order, picks a soul they
// control to destroy.
var eachPlayerDestroysASoul = engine.EffectFunc(func(c *engine.Ctx) {
	p := c.G.Turn.Active
	for range c.G.Players {
		pl := p
		them := *c
		them.Controller = pl
		them.Do(engine.Ask(func(c *engine.Ctx, a []int) {
			if ids := soulsOf(c.G, c.Controller); a[0] >= 0 && a[0] < len(ids) {
				c.G.DestroyObject(c.Controller, ids[a[0]])
			}
		}, engine.Question{Text: "Destroy which soul?", Options: func(c *engine.Ctx, _ []int) []string {
			var out []string
			for _, id := range soulsOf(c.G, c.Controller) {
				out = append(out, string(c.G.Object(id).Card))
			}
			return out
		}}))
		p = c.G.PlayerAfter(p)
	}
})

func soulsOf(g *engine.Game, p engine.PlayerID) []engine.ObjectID {
	var out []engine.ObjectID
	for _, id := range g.Players[p].InPlay {
		if g.Object(id).Role == engine.RoleSoul {
			out = append(out, id)
		}
	}
	return out
}

// giveLootThenPlay: target 0 gives the controller a loot card of their
// choice; it is revealed, then played for free if able.
var giveLootThenPlay = engine.Ask(func(c *engine.Ctx, a []int) {
	from := c.Targets[0].Player
	hand := c.G.Players[from].Hand
	if a[0] < 0 || a[0] >= len(hand) {
		return
	}
	id := c.G.GiveHandCard(from, c.Controller, hand[a[0]])
	for _, pl := range c.G.Players {
		c.G.LookAt(pl.ID, id) // "Reveal it"
	}
	c.G.PlayLootFree(c.Controller, id)
}, engine.Question{
	Text:   "Give which loot card?",
	Player: func(c *engine.Ctx, _ []int) engine.PlayerID { return c.Targets[0].Player },
	Options: func(c *engine.Ctx, _ []int) []string {
		var out []string
		for _, id := range c.G.Players[c.Targets[0].Player].Hand {
			out = append(out, string(c.G.Object(id).Card))
		}
		return out
	},
})

// divorce: target 0 gives half their cents and loot cards, rounded down,
// then an item, both of their choice.
var divorce = engine.EffectFunc(func(c *engine.Ctx) {
	from := c.Targets[0].Player
	n := c.G.Players[from].Cents / 2
	c.G.Players[from].Cents -= n
	c.G.Players[c.Controller].Cents += n
	cards := len(c.G.Players[from].Hand) / 2
	qs := make([]engine.Question, 0, cards+1)
	for range cards {
		qs = append(qs, engine.Question{
			Text:   "Give which loot card?",
			Player: func(c *engine.Ctx, _ []int) engine.PlayerID { return c.Targets[0].Player },
			Options: func(c *engine.Ctx, _ []int) []string {
				var out []string
				for _, id := range c.G.Players[c.Targets[0].Player].Hand {
					out = append(out, string(c.G.Object(id).Card))
				}
				return out
			},
		})
	}
	c.Do(engine.Ask(func(c *engine.Ctx, a []int) {
		// Each answer indexes the hand as it was; give from the highest
		// index down so earlier indexes stay valid.
		hand := append([]engine.ObjectID(nil), c.G.Players[from].Hand...)
		picked := map[int]bool{}
		for _, i := range a {
			if i >= 0 && i < len(hand) {
				picked[i] = true
			}
		}
		for i := len(hand) - 1; i >= 0; i-- {
			if picked[i] {
				c.G.GiveHandCard(from, c.Controller, hand[i])
			}
		}
	}, qs...), engine.Ask(func(c *engine.Ctx, a []int) {
		if ids := itemsOf(c.G, from); a[0] >= 0 && a[0] < len(ids) {
			c.G.GainControl(c.Controller, ids[a[0]])
		}
	}, engine.Question{
		Text:   "Give which item?",
		Player: func(c *engine.Ctx, _ []int) engine.PlayerID { return c.Targets[0].Player },
		Options: func(c *engine.Ctx, _ []int) []string {
			var out []string
			for _, id := range itemsOf(c.G, c.Targets[0].Player) {
				out = append(out, string(c.G.Object(id).Card))
			}
			return out
		},
	}))
})

// rerollRollersItem: "you may reroll an item they control" for the
// player who rolled.
var rerollRollersItem = engine.Ask(func(c *engine.Ctx, a []int) {
	if ids := rerollable(c.G, c.EventPlayer); a[0] >= 0 && a[0] < len(ids) {
		c.G.RerollItem(ids[a[0]])
	}
}, engine.Question{Text: "Reroll which of their items?", Options: func(c *engine.Ctx, _ []int) []string {
	var out []string
	for _, id := range rerollable(c.G, c.EventPlayer) {
		out = append(out, string(c.G.Object(id).Card))
	}
	if len(out) == 0 {
		return nil
	}
	return append(out, "none")
}})

func rerollable(g *engine.Game, p engine.PlayerID) []engine.ObjectID {
	if p == engine.NoPlayer {
		return nil
	}
	var out []engine.ObjectID
	for _, id := range itemsOf(g, p) {
		if !g.Eternal(id) {
			out = append(out, id)
		}
	}
	return out
}

// deactivateAnItem: "you may deactivate an item".
var deactivateAnItem = engine.Ask(func(c *engine.Ctx, a []int) {
	if ids := chargedItems(c.G); a[0] >= 0 && a[0] < len(ids) {
		c.G.Object(ids[a[0]]).Charged = false
	}
}, engine.Question{Text: "Deactivate which item?", Options: func(c *engine.Ctx, _ []int) []string {
	var out []string
	for _, id := range chargedItems(c.G) {
		out = append(out, string(c.G.Object(id).Card))
	}
	if len(out) == 0 {
		return nil
	}
	return append(out, "none")
}})

func chargedItems(g *engine.Game) []engine.ObjectID {
	var out []engine.ObjectID
	for _, id := range itemsOf(g, engine.NoPlayer) {
		if g.Object(id).Charged {
			out = append(out, id)
		}
	}
	return out
}

// lookMayDiscard is "Look at the top card of a deck. You may put it into
// discard or put it back on top".
var lookMayDiscard = engine.Ask(func(c *engine.Ctx, a []int) {
	if a[0] < 0 {
		return
	}
	d := c.Deck(a[0])
	c.G.LookAt(c.Controller, c.G.DeckTop(d, 1)...)
	if a[1] == 1 {
		c.G.MillTop(d)
	}
}, engine.DeckQuestion("Look at the top card of which deck?"), engine.Question{
	Text: "Put it into discard?",
	Subject: func(c *engine.Ctx, a []int) engine.ObjectID {
		return c.G.DeckTop(c.Deck(a[0]), 1)[0]
	},
	Options: func(c *engine.Ctx, a []int) []string {
		if a[0] < 0 || len(c.G.DeckTop(c.Deck(a[0]), 1)) == 0 {
			return nil
		}
		card := string(c.G.Object(c.G.DeckTop(c.Deck(a[0]), 1)[0]).Card)
		return []string{"keep " + card + " on top", "put " + card + " into discard"}
	},
})

// storeRoll keeps the roll's result for Telepathy for Dummies.
func storeRoll() engine.RollTable {
	var t engine.RollTable
	for r := 1; r <= 6; r++ {
		t = t.Results(r, r, engine.EffectFunc(func(c *engine.Ctx) {
			clearStoredRoll(c.G, c.Source)
			c.G.Boosts = append(c.G.Boosts, engine.Boost{Stat: engine.StatUsed, Player: engine.NoPlayer, Object: c.Source, Amount: -r})
		}))
	}
	return t
}

// storedRoll is the stored result, 0 if none; it is kept as a negative
// "used" mark so it ends with the turn.
func storedRoll(g *engine.Game, self engine.ObjectID) int {
	for _, b := range g.Boosts {
		if b.Stat == engine.StatUsed && b.Object == self && b.Amount < 0 {
			return -b.Amount
		}
	}
	return 0
}

func clearStoredRoll(g *engine.Game, self engine.ObjectID) {
	kept := g.Boosts[:0]
	for _, b := range g.Boosts {
		if b.Stat != engine.StatUsed || b.Object != self || b.Amount >= 0 {
			kept = append(kept, b)
		}
	}
	g.Boosts = kept
}

// revealGiveCurses: reveal the top 6 monster cards; each curse goes to a
// player of the controller's choice; the rest go to the bottom.
var revealGiveCurses = engine.EffectFunc(func(c *engine.Ctx) {
	top := c.G.DeckTop(engine.MonsterDeck, 6)
	for _, pl := range c.G.Players {
		c.G.LookAt(pl.ID, top...)
	}
	var curses, rest []engine.ObjectID
	for _, id := range top {
		if c.G.Def(id).Curse {
			curses = append(curses, id)
		} else {
			rest = append(rest, id)
		}
	}
	qs := make([]engine.Question, 0, len(curses))
	for _, id := range curses {
		card := string(c.G.Object(id).Card)
		qs = append(qs, engine.Question{Text: "Give " + card + " to which player?", Subject: func(*engine.Ctx, []int) engine.ObjectID { return id }, Options: func(c *engine.Ctx, _ []int) []string {
			var out []string
			for _, p := range livingPlayers(c.G, engine.NoPlayer) {
				out = append(out, playerLabel(c.G, p))
			}
			return out
		}})
	}
	c.Do(engine.Ask(func(c *engine.Ctx, a []int) {
		ps := livingPlayers(c.G, engine.NoPlayer)
		for i, id := range curses {
			if i < len(a) && a[i] >= 0 && a[i] < len(ps) {
				c.G.GiveCurseFromDeck(id, ps[a[i]])
			}
		}
		for _, id := range rest {
			c.G.DeckToBottom(engine.MonsterDeck, id)
		}
	}, qs...))
})

// thisGains gives this monster +n of a stat till end of turn.
func thisGains(stat engine.Stat, n int) engine.Effect {
	return engine.EffectFunc(func(c *engine.Ctx) {
		c.G.Boosts = append(c.G.Boosts, engine.Boost{Stat: stat, Player: engine.NoPlayer, Object: c.Source, Amount: n})
	})
}

// healThis heals this monster by n.
func healThis(n int) engine.Effect {
	return engine.EffectFunc(func(c *engine.Ctx) { c.G.HealObject(c.Source, n) })
}

// damageActive deals n damage from this monster to the active player.
func damageActive(n int) engine.Effect {
	return engine.EffectFunc(func(c *engine.Ctx) {
		c.G.DealDamageTo(engine.Target{Player: c.G.Turn.Active, IsPlayer: true}, n, engine.NoPlayer, c.Source)
	})
}

// eachPlayerTakes deals n damage from this card to each living player.
func eachPlayerTakes(n int) engine.Effect {
	return engine.EffectFunc(func(c *engine.Ctx) {
		for _, p := range livingPlayers(c.G, engine.NoPlayer) {
			c.G.DealDamageTo(engine.Target{Player: p, IsPlayer: true}, n, engine.NoPlayer, c.Source)
		}
	})
}

// damageNeighbour deals 1 damage to the player to the active player's
// left (the next in turn order, R-TURN-01) or right.
func damageNeighbour(left bool) engine.Effect {
	return engine.EffectFunc(func(c *engine.Ctx) {
		p := c.G.PlayerBefore(c.G.Turn.Active)
		if left {
			p = c.G.PlayerAfter(c.G.Turn.Active)
		}
		c.G.DealDamageTo(engine.Target{Player: p, IsPlayer: true}, 1, engine.NoPlayer, c.Source)
	})
}

// mustAttackThis is "The active player must attack this once each turn
// if able".
var mustAttackThis = engine.EffectFunc(func(c *engine.Ctx) {
	if c.G.Object(c.Source).Zone.Kind == engine.ZoneInPlay && (!c.G.Attack.Started || c.G.Attack.Target != c.Source) {
		c.G.Turn.MustAttack = c.Source
	}
})

// mostSouls lists the players who control the most souls, or tied.
func mostSouls(g *engine.Game) []engine.PlayerID {
	most := 0
	for _, pl := range g.Players {
		most = max(most, g.SoulValue(pl.ID))
	}
	var out []engine.PlayerID
	for _, pl := range g.Players {
		if g.SoulValue(pl.ID) == most {
			out = append(out, pl.ID)
		}
	}
	return out
}

// revealTwoBosses reveals monster cards till 2 bosses (monsters with a soul,
// S-YUG); they go to monster slots not being attacked, the rest into
// discard; the active player must attack once more.
var revealTwoBosses = engine.EffectFunc(func(c *engine.Ctx) {
	var bosses []engine.ObjectID
	for len(bosses) < 2 {
		top := c.G.DeckTop(engine.MonsterDeck, 1)
		if len(top) == 0 || slices.Contains(bosses, top[0]) {
			break // the deck ran out of other cards
		}
		for _, pl := range c.G.Players {
			c.G.LookAt(pl.ID, top...)
		}
		if d := c.G.Def(top[0]); d.Kind == engine.MonsterCard && d.Soul > 0 {
			bosses = append(bosses, top[0])
			c.G.DeckToBottom(engine.MonsterDeck, top[0]) // set aside until placed
			continue
		}
		c.G.MillTop(engine.MonsterDeck)
	}
	qs := make([]engine.Question, 0, len(bosses))
	for _, id := range bosses {
		card := string(c.G.Object(id).Card)
		qs = append(qs, engine.Question{Text: "Put " + card + " in which monster slot?", Subject: func(*engine.Ctx, []int) engine.ObjectID { return id }, Options: func(c *engine.Ctx, _ []int) []string {
			var out []string
			for _, i := range freeSlots(c.G) {
				out = append(out, "slot "+strconv.Itoa(i+1))
			}
			return out
		}})
	}
	c.Do(engine.Ask(func(c *engine.Ctx, a []int) {
		slots := freeSlots(c.G)
		for i, id := range bosses {
			if i < len(a) && a[i] >= 0 && a[i] < len(slots) {
				c.G.PlaceFromDeck(id, slots[a[i]])
			}
		}
		c.G.ForceAttacks(1, false)
	}, qs...))
})
