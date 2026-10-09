package b2

import (
	"strconv"

	engine "github.com/FourSoulsRepo/rules_engine"
)

// Pieces many Base Game cards share.

// extraLoot is every Base Game character's ↷ ability. It can be used on
// any player's turn, in response to anything.
var extraLoot = engine.Ability{
	Kind:    engine.Activated,
	Text:    "↷: Play an additional loot card this turn.",
	Costs:   []engine.Cost{engine.Tap()},
	Effects: []engine.Effect{engine.AddLootPlays(1)},
}

// rechargeAtEndOfTurn: "At the end of your turn, recharge this".
var rechargeAtEndOfTurn = engine.Ability{
	Kind:    engine.Triggered,
	Text:    "At the end of your turn, recharge this.",
	Trigger: engine.AtEndOfYourTurn(),
	Effects: []engine.Effect{engine.RechargeSelf()},
}

// discardOne is "Discard 1 loot card." for the controller.
var discardOne = engine.Ask(func(c *engine.Ctx, a []int) {
	if a[0] >= 0 {
		c.G.DiscardFromHand(c.Controller, c.HandCard(a[0]))
	}
}, engine.HandQuestion("Discard which loot card?"))

// damageEqualToRoll deals damage equal to the roll to target 0.
func damageEqualToRoll() engine.RollTable {
	var t engine.RollTable
	for r := 1; r <= 6; r++ {
		t = t.Results(r, r, engine.DealDamage(r, 0))
	}
	return t
}

// endYourTurn: "If it's your turn, cancel everything that hasn't
// resolved and end it".
var endYourTurn = engine.EffectFunc(func(c *engine.Ctx) {
	if c.Controller == c.G.Turn.Active {
		c.G.EndTurnNow()
	}
})

// lookMayBottom: "look at the top card of the deck. You may put it on
// the bottom".
func lookMayBottom(d engine.DeckKind) engine.Effect {
	return engine.Ask(func(c *engine.Ctx, a []int) {
		top := c.G.DeckTop(d, 1)
		if len(top) == 0 {
			return
		}
		c.G.LookAt(c.Controller, top...)
		if a[0] == 1 {
			c.G.DeckToBottom(d, top[0])
		}
	}, engine.Question{
		Text: "Put the top card on the bottom?",
		Options: func(c *engine.Ctx, _ []int) []string {
			top := c.G.DeckTop(d, 1)
			if len(top) == 0 {
				return nil
			}
			card := string(c.G.Object(top[0]).Card)
			return []string{"keep " + card + " on top", "put " + card + " on the bottom"}
		},
	})
}

// oneOnTop: "Look at the top 5 cards of the deck. Put 1 on top and the
// rest on the bottom".
func oneOnTop(d engine.DeckKind) engine.Effect {
	return engine.Ask(func(c *engine.Ctx, a []int) {
		top := c.G.DeckTop(d, 5)
		c.G.LookAt(c.Controller, top...)
		if a[0] < 0 {
			return
		}
		for i, id := range top {
			if i != a[0] {
				c.G.DeckToBottom(d, id)
			}
		}
	}, engine.Question{
		Text: "Which card stays on top? The rest go to the bottom.",
		Options: func(c *engine.Ctx, _ []int) []string {
			var out []string
			for _, id := range c.G.DeckTop(d, 5) {
				out = append(out, string(c.G.Object(id).Card))
			}
			return out
		},
	})
}

// may asks "yes or no" and runs the effects on yes: "you may …".
func may(question string, effects ...engine.Effect) engine.Effect {
	return engine.Ask(func(c *engine.Ctx, a []int) {
		if a[0] == 0 {
			c.Do(effects...)
		}
	}, engine.Question{Text: question, Options: func(*engine.Ctx, []int) []string { return []string{"yes", "no"} }})
}

// putBackInOrder: "look at the top n cards of the deck. Put them back in
// any order." Each answer picks the next card from the top.
func putBackInOrder(d engine.DeckKind, n int) engine.Effect {
	left := func(c *engine.Ctx, a []int) []engine.ObjectID {
		cards := c.G.DeckTop(d, n)
		for _, i := range a {
			if i < 0 {
				break
			}
			cards = append(cards[:i:i], cards[i+1:]...)
		}
		return cards
	}
	var qs []engine.Question
	for range n {
		qs = append(qs, engine.Question{
			Text: "Which card goes next, from the top?",
			Options: func(c *engine.Ctx, a []int) []string {
				var out []string
				for _, id := range left(c, a) {
					out = append(out, string(c.G.Object(id).Card))
				}
				return out
			},
		})
	}
	return engine.Ask(func(c *engine.Ctx, a []int) {
		cards := c.G.DeckTop(d, n)
		c.G.LookAt(c.Controller, cards...)
		var order []engine.ObjectID
		for _, i := range a {
			if i < 0 {
				break
			}
			order = append(order, cards[i])
			cards = append(cards[:i:i], cards[i+1:]...)
		}
		c.G.SetDeckTop(d, order)
	}, qs...)
}

// changeRollBy asks how much to change a roll by: "Add up to 2 to a roll".
func changeRollBy(amounts ...int) engine.Effect {
	return engine.Ask(func(c *engine.Ctx, a []int) {
		if it, ok := c.G.StackItemByID(c.Targets[0].StackID); ok {
			c.G.SetRoll(it.ID, it.Roll+amounts[a[0]])
		}
	}, engine.Question{Text: "Change the roll by?", Options: func(*engine.Ctx, []int) []string {
		var out []string
		for _, n := range amounts {
			out = append(out, strconv.Itoa(n))
		}
		return out
	}})
}

// endOfTurnIf triggers at the end of your turn when cond holds.
func endOfTurnIf(cond func(g *engine.Game, p engine.PlayerID) bool) engine.Trigger {
	return engine.Trigger{On: engine.EvEndOfTurn, Match: func(g *engine.Game, self engine.ObjectID, e engine.Event) bool {
		p := g.Object(self).Controller
		return e.Player == p && cond(g, p)
	}}
}

// firstAttackRollATK: "You have +1 ATK for your first attack roll each turn".
var firstAttackRollATK = engine.Static{
	Stat:   engine.StatPlayerATK,
	Amount: 1,
	Applies: func(g *engine.Game, self engine.ObjectID, p engine.PlayerID, _ engine.ObjectID) bool {
		return g.Object(self).Controller == p && p == g.Turn.Active && g.Attack.On && g.Turn.AttackRolls == 1
	},
}

// monstersInPlay lists the monsters on top of the monster slots.
func monstersInPlay(g *engine.Game) []engine.ObjectID {
	var out []engine.ObjectID
	for _, s := range g.Monsters {
		if top, ok := s.TopOf(); ok && g.Object(top).Role == engine.RoleMonster {
			out = append(out, top)
		}
	}
	return out
}

// livingPlayers lists living players, other than skip if skip >= 0.
func livingPlayers(g *engine.Game, skip engine.PlayerID) []engine.PlayerID {
	var out []engine.PlayerID
	for _, pl := range g.Players {
		if !pl.Dead && pl.ID != skip {
			out = append(out, pl.ID)
		}
	}
	return out
}

func playerLabel(g *engine.Game, p engine.PlayerID) string {
	return "player " + strconv.Itoa(int(p)+1) + " (" + string(g.Object(g.Players[p].Character).Card) + ")"
}

// damageAMonster: "Deal n damage to a monster", picked on resolution.
func damageAMonster(n int) engine.Effect {
	return engine.Ask(func(c *engine.Ctx, a []int) {
		if a[0] >= 0 {
			c.G.DealDamageTo(engine.Target{Object: monstersInPlay(c.G)[a[0]]}, n, c.Controller, c.Source)
		}
	}, engine.Question{Text: "Deal damage to which monster?", Options: func(c *engine.Ctx, _ []int) []string {
		var out []string
		for _, id := range monstersInPlay(c.G) {
			out = append(out, string(c.G.Object(id).Card))
		}
		return out
	}})
}

// damageAPlayer: "Deal n damage to a player" (another player if others),
// picked on resolution.
func damageAPlayer(n int, others bool) engine.Effect {
	skip := func(c *engine.Ctx) engine.PlayerID {
		if others {
			return c.Controller
		}
		return engine.NoPlayer
	}
	return engine.Ask(func(c *engine.Ctx, a []int) {
		if a[0] >= 0 {
			p := livingPlayers(c.G, skip(c))[a[0]]
			c.G.DealDamageTo(engine.Target{Player: p, IsPlayer: true}, n, c.Controller, c.Source)
		}
	}, engine.Question{Text: "Deal damage to which player?", Options: func(c *engine.Ctx, _ []int) []string {
		var out []string
		for _, p := range livingPlayers(c.G, skip(c)) {
			out = append(out, playerLabel(c.G, p))
		}
		return out
	}})
}

// rechargeAnItem: "you may recharge an item", picked on resolution; the
// last option declines.
var rechargeAnItem = engine.Ask(func(c *engine.Ctx, a []int) {
	if ids := allItems(c.G); a[0] >= 0 && a[0] < len(ids) {
		c.G.Recharge(ids[a[0]])
	}
}, engine.Question{Text: "Recharge which item?", Options: func(c *engine.Ctx, _ []int) []string {
	var out []string
	for _, id := range allItems(c.G) {
		out = append(out, string(c.G.Object(id).Card))
	}
	return append(out, "none")
}})

// allItems lists every item players control, in seat order.
func allItems(g *engine.Game) []engine.ObjectID {
	var out []engine.ObjectID
	for _, pl := range g.Players {
		for _, id := range pl.InPlay {
			if g.Object(id).Role == engine.RoleItem {
				out = append(out, id)
			}
		}
	}
	return out
}

// onYourAttackRollOf triggers when the controller's attack roll
// resolves as 6.
func onYourAttackRollOf(n int) engine.Trigger {
	return engine.Trigger{On: engine.EvRollResolved, Match: func(g *engine.Game, self engine.ObjectID, e engine.Event) bool {
		return e.Amount == n && e.Text == "attack" && e.Player == g.Object(self).Controller
	}}
}

// whenYouDealCombatDamage triggers when the controller's attack damages
// a monster.
func whenYouDealCombatDamage() engine.Trigger {
	return engine.Trigger{On: engine.EvDamaged, Match: func(g *engine.Game, self engine.ObjectID, e engine.Event) bool {
		return e.Text == "combat" && e.Object != 0 && g.Turn.Active == g.Object(self).Controller
	}}
}

// giveThisAway: "give this to another player" (the haunts).
var giveThisAway = engine.Ask(func(c *engine.Ctx, a []int) {
	if a[0] >= 0 {
		c.G.GainControl(otherPlayers(c)[a[0]], c.Source)
	}
}, engine.Question{Text: "Give this to which player?", Options: func(c *engine.Ctx, _ []int) []string {
	var out []string
	for _, p := range otherPlayers(c) {
		out = append(out, playerLabel(c.G, p))
	}
	return out
}})

// otherPlayers lists the other players, alive or not ("another player").
func otherPlayers(c *engine.Ctx) []engine.PlayerID {
	var out []engine.PlayerID
	for _, pl := range c.G.Players {
		if pl.ID != c.Controller {
			out = append(out, pl.ID)
		}
	}
	return out
}

// searchGuppy: "Search the treasure deck for a Guppy item, gain it, then
// shuffle the treasure deck".
var searchGuppy = engine.Ask(func(c *engine.Ctx, a []int) {
	if ids := guppiesInDeck(c.G); a[0] >= 0 && a[0] < len(ids) {
		c.G.TakeFromDeck(engine.TreasureDeck, ids[a[0]], c.Controller)
	}
	c.G.ShuffleDeck(engine.TreasureDeck)
}, engine.Question{Text: "Gain which Guppy item?", Options: func(c *engine.Ctx, _ []int) []string {
	var out []string
	for _, id := range guppiesInDeck(c.G) {
		out = append(out, string(c.G.Object(id).Card))
	}
	return out
}})

func guppiesInDeck(g *engine.Game) []engine.ObjectID {
	var out []engine.ObjectID
	for _, id := range g.Decks[engine.TreasureDeck] {
		if g.Def(id).Guppy {
			out = append(out, id)
		}
	}
	return out
}

// pickAPlayer asks the controller for a living player.
func pickAPlayer(text string) engine.Question {
	return engine.Question{Text: text, Options: func(c *engine.Ctx, _ []int) []string {
		var out []string
		for _, p := range livingPlayers(c.G, engine.NoPlayer) {
			out = append(out, playerLabel(c.G, p))
		}
		return out
	}}
}

// killAPlayer: "the active player kills a player".
var killAPlayer = engine.Ask(func(c *engine.Ctx, a []int) {
	if a[0] >= 0 {
		p := livingPlayers(c.G, engine.NoPlayer)[a[0]]
		c.Do(engine.EffectFunc(func(c *engine.Ctx) {
			c.Targets = []engine.Chosen{{Kind: engine.TargetPlayer, Player: p}}
			c.Do(engine.Kill(0))
		}))
	}
}, pickAPlayer("Kill which player?"))

// forcedDiscardQuestions: the controller picks a player, who then picks
// n cards from their hand, one at a time.
func forcedDiscardQuestions(n int) []engine.Question {
	qs := []engine.Question{pickAPlayer("Which player discards?")}
	for range n {
		qs = append(qs, engine.Question{
			Text:   "Discard which loot card?",
			Player: func(c *engine.Ctx, a []int) engine.PlayerID { return livingPlayers(c.G, engine.NoPlayer)[a[0]] },
			Options: func(c *engine.Ctx, a []int) []string {
				if a[0] < 0 {
					return nil
				}
				var out []string
				for _, id := range forcedLeft(c, a) {
					out = append(out, string(c.G.Object(id).Card))
				}
				return out
			},
		})
	}
	return qs
}

// forcedLeft is the chosen player's hand without the cards picked so far.
func forcedLeft(c *engine.Ctx, a []int) []engine.ObjectID {
	left := append([]engine.ObjectID(nil), c.G.Players[livingPlayers(c.G, engine.NoPlayer)[a[0]]].Hand...)
	for _, i := range a[1:] {
		if i < 0 || i >= len(left) {
			break
		}
		left = append(left[:i:i], left[i+1:]...)
	}
	return left
}

// forcedDiscard discards the cards picked by forcedDiscardQuestions.
func forcedDiscard(n int) func(c *engine.Ctx, a []int) {
	return func(c *engine.Ctx, a []int) {
		if a[0] < 0 {
			return
		}
		p := livingPlayers(c.G, engine.NoPlayer)[a[0]]
		left := append([]engine.ObjectID(nil), c.G.Players[p].Hand...)
		for _, i := range a[1 : n+1] {
			if i < 0 || i >= len(left) {
				break
			}
			c.G.DiscardFromHand(p, left[i])
			left = append(left[:i:i], left[i+1:]...)
		}
	}
}

// damageAnything: "deal n damage to a monster or player", picked on
// resolution.
func damageAnything(n int) engine.Effect {
	targets := func(g *engine.Game) []engine.Target {
		var out []engine.Target
		for _, id := range monstersInPlay(g) {
			out = append(out, engine.Target{Object: id})
		}
		for _, p := range livingPlayers(g, engine.NoPlayer) {
			out = append(out, engine.Target{Player: p, IsPlayer: true})
		}
		return out
	}
	return engine.Ask(func(c *engine.Ctx, a []int) {
		if a[0] >= 0 {
			c.G.DealDamageTo(targets(c.G)[a[0]], n, c.Controller, c.Source)
		}
	}, engine.Question{Text: "Deal damage to?", Options: func(c *engine.Ctx, _ []int) []string {
		var out []string
		for _, t := range targets(c.G) {
			if t.IsPlayer {
				out = append(out, playerLabel(c.G, t.Player))
			} else {
				out = append(out, string(c.G.Object(t.Object).Card))
			}
		}
		return out
	}})
}
