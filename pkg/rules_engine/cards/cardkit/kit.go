// Package cardkit holds pieces many cards share, in every set.
package cardkit

import (
	"slices"
	"strconv"

	engine "github.com/FourSoulsRepo/rules_engine"
)

// ExtraLoot is every Base Game character's ↷ ability. It can be used on
// any player's turn, in response to anything.
var ExtraLoot = engine.Ability{
	Kind:    engine.Activated,
	Text:    "↷: Play an additional loot card this turn.",
	Costs:   []engine.Cost{engine.Tap()},
	Effects: []engine.Effect{engine.AddLootPlays(1)},
}

// RechargeAtEndOfTurn is "At the end of your turn, recharge this".
var RechargeAtEndOfTurn = engine.Ability{
	Kind:    engine.Triggered,
	Text:    "At the end of your turn, recharge this.",
	Trigger: engine.AtEndOfYourTurn(),
	Effects: []engine.Effect{engine.RechargeSelf()},
}

// DiscardOne is "Discard 1 loot card." for the controller.
var DiscardOne = engine.Ask(func(c *engine.Ctx, a []int) {
	if a[0] >= 0 {
		c.G.DiscardFromHand(c.Controller, c.HandCard(a[0]))
	}
}, engine.HandQuestion("Discard which loot card?"))

// DamageEqualToRoll deals damage equal to the roll to target 0.
func DamageEqualToRoll() engine.RollTable {
	var t engine.RollTable
	for r := 1; r <= 6; r++ {
		t = t.Results(r, r, engine.DealDamage(r, 0))
	}
	return t
}

// EndYourTurn is "If it's your turn, cancel everything that hasn't
// resolved and end it".
var EndYourTurn = engine.EffectFunc(func(c *engine.Ctx) {
	if c.Controller == c.G.Turn.Active {
		c.G.EndTurnNow()
	}
})

// LookMayBottom is "look at the top card of the deck. You may put it on
// the bottom".
func LookMayBottom(d engine.DeckKind) engine.Effect {
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
		Text:    "Put the top card on the bottom?",
		Subject: func(c *engine.Ctx, _ []int) engine.ObjectID { return c.G.DeckTop(d, 1)[0] },
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

// OneOnTop is "Look at the top 5 cards of the deck. Put 1 on top and the
// rest on the bottom". The rest go in random order (Q-22).
func OneOnTop(d engine.DeckKind) engine.Effect {
	return engine.Ask(func(c *engine.Ctx, a []int) {
		top := c.G.DeckTop(d, 5)
		c.G.LookAt(c.Controller, top...)
		if a[0] < 0 || a[0] >= len(top) {
			return
		}
		c.G.ToBottomRandom(d, append(top[:a[0]:a[0]], top[a[0]+1:]...))
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

// May asks "yes or no" and runs the effects on yes: "you may …".
func May(question string, effects ...engine.Effect) engine.Effect {
	return engine.Ask(func(c *engine.Ctx, a []int) {
		if a[0] == 0 {
			c.Do(effects...)
		}
	}, engine.Question{Text: question, Options: func(*engine.Ctx, []int) []string { return []string{"yes", "No"} }})
}

// PutBackInOrder is "look at the top n cards of the deck. Put them back in
// any order." One order question puts them back, top first (19.16).
func PutBackInOrder(d engine.DeckKind, n int) engine.Effect {
	return engine.Ask(func(c *engine.Ctx, _ []int) {
		c.G.SetDeckTop(d, c.Order)
	}, engine.Question{
		Text:  PutBackText,
		Order: func(c *engine.Ctx, _ []int) []engine.ObjectID { return c.G.DeckTop(d, n) },
	})
}

// PutBackText is the question of every "put them back in any order".
const PutBackText = "Put the cards back in any order"

// ChangeRollBy asks how much to change a roll by: "Add up to 2 to a roll".
func ChangeRollBy(amounts ...int) engine.Effect {
	return engine.Ask(func(c *engine.Ctx, a []int) {
		if it, ok := c.G.StackItemByID(c.Targets[0].StackID); ok {
			c.G.SetRoll(it.ID, it.Roll+amounts[a[0]], c.Source)
		}
	}, engine.Question{Text: "Change the roll by?", Options: func(*engine.Ctx, []int) []string {
		var out []string
		for _, n := range amounts {
			out = append(out, strconv.Itoa(n))
		}
		return out
	}})
}

// EndOfTurnIf triggers at the end of your turn when cond holds.
func EndOfTurnIf(cond func(g *engine.Game, p engine.PlayerID) bool) engine.Trigger {
	return engine.Trigger{On: engine.EvEndOfTurn, Match: func(g *engine.Game, self engine.ObjectID, e engine.Event) bool {
		p := g.Object(self).Controller
		return e.Player == p && cond(g, p)
	}}
}

// FirstAttackRollATK is "You have +1 ATK for your first attack roll each turn".
var FirstAttackRollATK = engine.Static{
	Stat:   engine.StatPlayerATK,
	Amount: 1,
	When:   "first attack roll",
	Applies: func(g *engine.Game, self engine.ObjectID, p engine.PlayerID, _ engine.ObjectID) bool {
		return g.Object(self).Controller == p && p == g.Turn.Active && g.Attack.On && g.Turn.AttackRolls == 1
	},
}

// MonstersInPlay lists the monsters on top of the monster slots.
func MonstersInPlay(g *engine.Game) []engine.ObjectID {
	var out []engine.ObjectID
	for _, s := range g.Monsters {
		if top, ok := s.TopOf(); ok && g.Object(top).Role == engine.RoleMonster {
			out = append(out, top)
		}
	}
	return out
}

// LivingPlayers lists living players, other than skip if skip >= 0.
func LivingPlayers(g *engine.Game, skip engine.PlayerID) []engine.PlayerID {
	var out []engine.PlayerID
	for _, pl := range g.Players {
		if !pl.Dead && pl.ID != skip {
			out = append(out, pl.ID)
		}
	}
	return out
}

// PlayerLabel names a player as the prompts do: "player 1 (isaac)".
func PlayerLabel(g *engine.Game, p engine.PlayerID) string {
	return "player " + strconv.Itoa(int(p)+1) + " (" + string(g.Object(g.Players[p].Character).Card) + ")"
}

// DamageAMonster is "Deal n damage to a monster", picked on resolution.
func DamageAMonster(n int) engine.Effect {
	return engine.Ask(func(c *engine.Ctx, a []int) {
		if a[0] >= 0 {
			c.G.DealDamageTo(engine.Target{Object: MonstersInPlay(c.G)[a[0]]}, n, c.Controller, c.Source)
		}
	}, engine.Question{Text: "Deal damage to which monster?", Options: func(c *engine.Ctx, _ []int) []string {
		var out []string
		for _, id := range MonstersInPlay(c.G) {
			out = append(out, string(c.G.Object(id).Card))
		}
		return out
	}})
}

// DamageAPlayer is "Deal n damage to a player" (another player if others),
// picked on resolution.
func DamageAPlayer(n int, others bool) engine.Effect {
	skip := func(c *engine.Ctx) engine.PlayerID {
		if others {
			return c.Controller
		}
		return engine.NoPlayer
	}
	return engine.Ask(func(c *engine.Ctx, a []int) {
		if a[0] >= 0 {
			p := LivingPlayers(c.G, skip(c))[a[0]]
			c.G.DealDamageTo(engine.Target{Player: p, IsPlayer: true}, n, c.Controller, c.Source)
		}
	}, engine.Question{Text: "Deal damage to which player?", Options: func(c *engine.Ctx, _ []int) []string {
		var out []string
		for _, p := range LivingPlayers(c.G, skip(c)) {
			out = append(out, PlayerLabel(c.G, p))
		}
		return out
	}})
}

// RechargeAnItem is "you may recharge an item", picked on resolution; the
// last option declines.
var RechargeAnItem = engine.Ask(func(c *engine.Ctx, a []int) {
	if ids := AllItems(c.G); a[0] >= 0 && a[0] < len(ids) {
		c.G.Recharge(ids[a[0]])
	}
}, engine.Question{Text: "Recharge which item?", Options: func(c *engine.Ctx, _ []int) []string {
	var out []string
	for _, id := range AllItems(c.G) {
		out = append(out, string(c.G.Object(id).Card))
	}
	return append(out, "none")
}})

// AllItems lists every item players control, in seat order.
func AllItems(g *engine.Game) []engine.ObjectID {
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

// OnYourAttackRollOf triggers when the controller's attack roll
// resolves as 6.
func OnYourAttackRollOf(n int) engine.Trigger {
	return engine.Trigger{On: engine.EvRollResolved, Match: func(g *engine.Game, self engine.ObjectID, e engine.Event) bool {
		return e.Amount == n && e.Text == "attack" && e.Player == g.Object(self).Controller
	}}
}

// WhenYouDealCombatDamage triggers when the controller's attack damages
// a monster.
func WhenYouDealCombatDamage() engine.Trigger {
	return engine.Trigger{On: engine.EvDamaged, Match: func(g *engine.Game, self engine.ObjectID, e engine.Event) bool {
		return e.Text == "combat" && e.Object != 0 && g.Turn.Active == g.Object(self).Controller
	}}
}

// GiveThisAway is "give this to another player" (the haunts).
var GiveThisAway = engine.Ask(func(c *engine.Ctx, a []int) {
	if a[0] >= 0 {
		c.G.GainControl(OtherPlayers(c)[a[0]], c.Source)
	}
}, engine.Question{Text: "Give this to which player?", Options: func(c *engine.Ctx, _ []int) []string {
	var out []string
	for _, p := range OtherPlayers(c) {
		out = append(out, PlayerLabel(c.G, p))
	}
	return out
}})

// OtherPlayers lists the other players, alive or not ("another player").
func OtherPlayers(c *engine.Ctx) []engine.PlayerID {
	var out []engine.PlayerID
	for _, pl := range c.G.Players {
		if pl.ID != c.Controller {
			out = append(out, pl.ID)
		}
	}
	return out
}

// SearchGuppy is "Search the treasure deck for a Guppy item, gain it, then
// shuffle the treasure deck".
var SearchGuppy = engine.Ask(func(c *engine.Ctx, a []int) {
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

// PickAPlayer asks the controller for a living player.
func PickAPlayer(text string) engine.Question {
	return engine.Question{Text: text, Options: func(c *engine.Ctx, _ []int) []string {
		var out []string
		for _, p := range LivingPlayers(c.G, engine.NoPlayer) {
			out = append(out, PlayerLabel(c.G, p))
		}
		return out
	}}
}

// KillAPlayer is "the active player kills a player".
var KillAPlayer = engine.Ask(func(c *engine.Ctx, a []int) {
	if a[0] >= 0 {
		p := LivingPlayers(c.G, engine.NoPlayer)[a[0]]
		c.Do(engine.EffectFunc(func(c *engine.Ctx) {
			c.Targets = []engine.Chosen{{Kind: engine.TargetPlayer, Player: p}}
			c.Do(engine.Kill(0))
		}))
	}
}, PickAPlayer("Kill which player?"))

// ForcedDiscardQuestions is the controller picks a player, who then picks
// n cards from their hand, one at a time.
func ForcedDiscardQuestions(n int) []engine.Question {
	qs := []engine.Question{PickAPlayer("Which player discards?")}
	for range n {
		qs = append(qs, engine.Question{
			Text:   "Discard which loot card?",
			Player: func(c *engine.Ctx, a []int) engine.PlayerID { return LivingPlayers(c.G, engine.NoPlayer)[a[0]] },
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
	left := append([]engine.ObjectID(nil), c.G.Players[LivingPlayers(c.G, engine.NoPlayer)[a[0]]].Hand...)
	for _, i := range a[1:] {
		if i < 0 || i >= len(left) {
			break
		}
		left = append(left[:i:i], left[i+1:]...)
	}
	return left
}

// ForcedDiscard discards the cards picked by ForcedDiscardQuestions.
func ForcedDiscard(n int) func(c *engine.Ctx, a []int) {
	return func(c *engine.Ctx, a []int) {
		if a[0] < 0 {
			return
		}
		p := LivingPlayers(c.G, engine.NoPlayer)[a[0]]
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

// DamageAnything is "deal n damage to a monster or player", picked on
// resolution.
func DamageAnything(n int) engine.Effect {
	targets := func(g *engine.Game) []engine.Target {
		var out []engine.Target
		for _, id := range MonstersInPlay(g) {
			out = append(out, engine.Target{Object: id})
		}
		for _, p := range LivingPlayers(g, engine.NoPlayer) {
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
				out = append(out, PlayerLabel(c.G, t.Player))
			} else {
				out = append(out, string(c.G.Object(t.Object).Card))
			}
		}
		return out
	}})
}

// AttackRollOnThis triggers when the attacking player's attack roll
// against this monster resolves as n.
func AttackRollOnThis(n int) engine.Trigger {
	return engine.Trigger{On: engine.EvRollResolved, Match: func(g *engine.Game, self engine.ObjectID, e engine.Event) bool {
		return e.Amount == n && e.Text == "attack" && g.Attack.Target == self
	}}
}

// PlayerRollOf triggers when a player's roll resolves as n.
func PlayerRollOf(n int) engine.Trigger {
	return engine.Trigger{On: engine.EvRollResolved, Match: func(_ *engine.Game, _ engine.ObjectID, e engine.Event) bool {
		return e.Amount == n && e.Player != engine.NoPlayer
	}}
}

// WhenThisTakesDamage triggers when this monster takes damage; combat
// only, if set.
func WhenThisTakesDamage(combat bool) engine.Trigger {
	return engine.Trigger{On: engine.EvDamaged, Match: func(_ *engine.Game, self engine.ObjectID, e engine.Event) bool {
		return e.Object == self && (!combat || e.Text == "combat")
	}}
}

// WhenThisDealsDamage triggers when this monster damages a player;
// combat only, if set. c.EventPlayer is the player.
func WhenThisDealsDamage(combat bool) engine.Trigger {
	return engine.Trigger{On: engine.EvDamaged, Match: func(_ *engine.Game, self engine.ObjectID, e engine.Event) bool {
		return e.Source == self && e.Player != engine.NoPlayer && (!combat || e.Text == "combat")
	}}
}

// ForRoller runs effects as if the event's player controlled them: "they
// gain 1¢" on a monster card.
func ForRoller(effects ...engine.Effect) engine.Effect {
	return engine.EffectFunc(func(c *engine.Ctx) {
		if c.EventPlayer == engine.NoPlayer {
			return
		}
		them := *c
		them.Controller = c.EventPlayer
		them.Do(effects...)
	})
}

// RollerDiscards is "they discard a loot card".
var RollerDiscards = ForRoller(DiscardOne)

// EachMonsterGains gives each monster +1 of a stat till end of turn.
func EachMonsterGains(stat engine.Stat) engine.Effect {
	return engine.EffectFunc(func(c *engine.Ctx) {
		for _, id := range MonstersInPlay(c.G) {
			c.G.Boosts = append(c.G.Boosts, engine.Boost{Source: c.Source, Stat: stat, Player: engine.NoPlayer, Object: id, Amount: 1})
		}
	})
}

// GuppyItems counts the Guppy items p controls.
func GuppyItems(g *engine.Game, p engine.PlayerID) int {
	n := 0
	for _, id := range g.Players[p].InPlay {
		if g.Object(id).Role == engine.RoleItem && g.Def(id).Guppy {
			n++
		}
	}
	return n
}

// KillAMonster is "Kill a monster", picked on resolution.
var KillAMonster = engine.Ask(func(c *engine.Ctx, a []int) {
	if ids := MonstersInPlay(c.G); a[0] >= 0 && a[0] < len(ids) {
		c.G.KillObject(ids[a[0]])
	}
}, engine.Question{Text: "Kill which monster?", Options: func(c *engine.Ctx, _ []int) []string {
	var out []string
	for _, id := range MonstersInPlay(c.G) {
		out = append(out, string(c.G.Object(id).Card))
	}
	return out
}})

// Used reports whether the object's "once each turn" ability was used
// this turn.
func Used(g *engine.Game, id engine.ObjectID) bool {
	for _, b := range g.Boosts {
		if b.Stat == engine.StatUsed && b.Object == id {
			return true
		}
	}
	return false
}

// MarkUsed marks the object's "once each turn" ability as used; the mark
// ends with the turn.
func MarkUsed(g *engine.Game, id engine.ObjectID) {
	g.Boosts = append(g.Boosts, engine.Boost{Stat: engine.StatUsed, Player: engine.NoPlayer, Object: id, Amount: 1})
}

// UsedOn reports whether the object's ability was used this turn for n,
// e.g. a roll's stack ID ("once per roll") or a chosen number.
func UsedOn(g *engine.Game, id engine.ObjectID, n int) bool {
	for _, b := range g.Boosts {
		if b.Stat == engine.StatUsed && b.Object == id && b.Amount == n {
			return true
		}
	}
	return false
}

// MarkUsedOn marks the object's ability as used this turn for n.
func MarkUsedOn(g *engine.Game, id engine.ObjectID, n int) {
	g.Boosts = append(g.Boosts, engine.Boost{Stat: engine.StatUsed, Player: engine.NoPlayer, Object: id, Amount: n})
}

// ItemsOf lists the items p controls, or every player's for NoPlayer.
func ItemsOf(g *engine.Game, p engine.PlayerID) []engine.ObjectID {
	var out []engine.ObjectID
	for _, pl := range g.Players {
		if p != engine.NoPlayer && pl.ID != p {
			continue
		}
		for _, id := range pl.InPlay {
			if g.Object(id).Role == engine.RoleItem {
				out = append(out, id)
			}
		}
	}
	return out
}

// FreeSlots lists monster slots whose monster is not being attacked.
func FreeSlots(g *engine.Game) []int {
	var out []int
	for i, s := range g.Monsters {
		if top, ok := s.TopOf(); !ok || !g.Attack.On || top != g.Attack.Target {
			out = append(out, i)
		}
	}
	return out
}

// WhenYouGainThis triggers when the controller gains this item: from the
// treasure deck, by a purchase, or by a search; not when it is given or
// stolen (R-MECH-35: "When you gain or purchase this").
func WhenYouGainThis() engine.Trigger {
	return engine.Trigger{On: engine.EvGainedItem, Match: func(_ *engine.Game, self engine.ObjectID, e engine.Event) bool {
		return e.Object == self
	}}
}

// WhenAPlayerGainsATreasure triggers when a player gains a treasure: a card
// from the treasure deck, bought from it too; not a shop item (R-CARD-37).
func WhenAPlayerGainsATreasure() engine.Trigger {
	return engine.Trigger{On: engine.EvGainedTreasure, Match: func(_ *engine.Game, _ engine.ObjectID, e engine.Event) bool {
		return e.Text == "" || e.Text == "searched"
	}}
}

// LookMayDiscard is "Look at the top card of a deck. You may put it into
// discard or put it back on top".
var LookMayDiscard = engine.Ask(func(c *engine.Ctx, a []int) {
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

// NotAttacked is a target filter: a monster that is not the one being attacked.
func NotAttacked(g *engine.Game, _ engine.ObjectID, c engine.Chosen) bool {
	return !g.Attack.On || g.Attack.Target != c.Object
}

// InShop is a target filter: an item in a shop slot.
func InShop(g *engine.Game, _ engine.ObjectID, c engine.Chosen) bool {
	return slices.Contains(g.ShopItems(), c.Object)
}

// DiscardTarget puts target 0, a monster or a shop item, into discard;
// its slot is refilled.
var DiscardTarget = engine.EffectFunc(func(c *engine.Ctx) {
	id := c.Targets[0].Object
	if slices.Contains(c.G.ShopItems(), id) {
		c.G.DiscardShopItem(id) // refills its slot
		return
	}
	c.G.DiscardMonster(id)
	c.G.RefillSlots()
})
