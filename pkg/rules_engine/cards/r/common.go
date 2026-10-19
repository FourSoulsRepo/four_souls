package r

import (
	"slices"
	"strings"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/cardkit"
)

// Pieces many cards share; they live in cardkit for every set.
var (
	itemsOf        = cardkit.ItemsOf
	freeSlots      = cardkit.FreeSlots
	notAttacked    = cardkit.NotAttacked
	inShop         = cardkit.InShop
	putBackInOrder = cardkit.PutBackInOrder
	otherPlayers   = cardkit.OtherPlayers
	damageAMonster = cardkit.DamageAMonster
	playerLabel    = cardkit.PlayerLabel
	damageAnything = cardkit.DamageAnything
	monstersInPlay = cardkit.MonstersInPlay
	discardOne     = cardkit.DiscardOne
	extraLoot      = cardkit.ExtraLoot
)

// eachOtherThan deals n damage from the controller to each player but
// skip and the controller, in turn order from the controller (R-MECH-28).
func eachOtherThan(c *engine.Ctx, skip engine.PlayerID, n int) {
	order := []engine.PlayerID{}
	p := c.Controller
	for range c.G.Players {
		if p != skip && p != c.Controller {
			order = append(order, p)
		}
		p = c.G.PlayerAfter(p)
	}
	// The first in turn order resolves first: it goes on the stack last.
	for i := len(order) - 1; i >= 0; i-- {
		c.G.DealDamageTo(engine.Target{Player: order[i], IsPlayer: true}, n, c.Controller, c.Source)
	}
}

// stealFromShop: player p steals a shop item of their choice; with may
// they may decline.
func stealFromShop(p engine.PlayerID, may bool) engine.Effect {
	return engine.Ask(func(c *engine.Ctx, a []int) {
		if items := c.G.ShopItems(); a[0] >= 0 && a[0] < len(items) {
			c.G.GainControl(p, items[a[0]])
		}
	}, engine.Question{
		Text:   "Steal which shop item?",
		Player: func(*engine.Ctx, []int) engine.PlayerID { return p },
		Options: func(c *engine.Ctx, _ []int) []string {
			out := labelsOf(c.G, c.G.ShopItems())
			if may && len(out) > 0 {
				out = append(out, "none")
			}
			return out
		},
	})
}

// labelsOf lists the cards of objects.
func labelsOf(g *engine.Game, ids []engine.ObjectID) []string {
	out := make([]string, len(ids))
	for i, id := range ids {
		out[i] = string(g.Object(id).Card)
	}
	return out
}

// playerLabels names players as the prompts do.
func playerLabels(g *engine.Game, ps []engine.PlayerID) []string {
	out := make([]string, len(ps))
	for i, p := range ps {
		out[i] = playerLabel(g, p)
	}
	return out
}

// othersInOrder lists the other players in turn order from the
// controller.
func othersInOrder(c *engine.Ctx) []engine.PlayerID {
	var out []engine.PlayerID
	p := c.G.PlayerAfter(c.Controller)
	for range len(c.G.Players) - 1 {
		out = append(out, p)
		p = c.G.PlayerAfter(p)
	}
	return out
}

// eachOtherMay asks each other player, in turn order, a yes-or-no
// question; do gets the players who said yes ("each other player may
// choose to loot 1").
func eachOtherMay(question string, do func(c *engine.Ctx, yes []engine.PlayerID)) engine.Effect {
	qs := make([]engine.Question, 3) // at most 3 other players
	for i := range qs {
		qs[i] = engine.Question{
			Text:   question,
			Player: func(c *engine.Ctx, _ []int) engine.PlayerID { return othersInOrder(c)[i] },
			Options: func(c *engine.Ctx, _ []int) []string {
				if i >= len(othersInOrder(c)) {
					return nil
				}
				return []string{"yes", "no"}
			},
		}
	}
	return engine.Ask(func(c *engine.Ctx, a []int) {
		var yes []engine.PlayerID
		for i, p := range othersInOrder(c) {
			if i < len(a) && a[i] == 0 {
				yes = append(yes, p)
			}
		}
		do(c, yes)
	}, qs...)
}

// eachGivesYouALoot: each of ps gives the controller a loot card of
// their choice, in order.
func eachGivesYouALoot(ps []engine.PlayerID) engine.Effect {
	qs := make([]engine.Question, len(ps))
	for i, p := range ps {
		qs[i] = engine.Question{
			Text:    "Give which loot card?",
			Player:  func(*engine.Ctx, []int) engine.PlayerID { return p },
			Options: func(c *engine.Ctx, _ []int) []string { return labelsOf(c.G, c.G.Players[p].Hand) },
		}
	}
	return engine.Ask(func(c *engine.Ctx, a []int) {
		for i, p := range ps {
			if hand := c.G.Players[p].Hand; a[i] >= 0 && a[i] < len(hand) {
				c.G.GiveHandCard(p, c.Controller, hand[a[i]])
			}
		}
	}, qs...)
}

// discardAShopItem puts a shop item of the controller's choice into
// discard; its slot is refilled.
var discardAShopItem = shopItemToDiscard(false)

// mayDiscardAShopItem is discardAShopItem with "none".
var mayDiscardAShopItem = shopItemToDiscard(true)

func shopItemToDiscard(may bool) engine.Effect {
	return engine.Ask(func(c *engine.Ctx, a []int) {
		if items := c.G.ShopItems(); a[0] >= 0 && a[0] < len(items) {
			c.G.DiscardShopItem(items[a[0]])
		}
	}, engine.Question{Text: "Put which shop item into discard?", Options: func(c *engine.Ctx, _ []int) []string {
		out := labelsOf(c.G, c.G.ShopItems())
		if may && len(out) > 0 {
			out = append(out, "none")
		}
		return out
	}})
}

// cardsInPlay lists the non-eternal cards in play: characters, items,
// shop items, monsters and rooms.
func cardsInPlay(g *engine.Game) []engine.ObjectID {
	var out []engine.ObjectID
	for _, pl := range g.Players {
		out = append(out, pl.Character)
		for _, id := range pl.InPlay {
			if g.Object(id).Role == engine.RoleItem {
				out = append(out, id)
			}
		}
	}
	out = append(out, g.ShopItems()...)
	for _, rows := range [][]engine.Slot{g.Monsters, g.Rooms} {
		for _, s := range rows {
			if top, ok := s.TopOf(); ok {
				out = append(out, top)
			}
		}
	}
	kept := out[:0]
	for _, id := range out {
		if !g.Eternal(id) {
			kept = append(kept, id)
		}
	}
	return kept
}

// changeACounter: "Choose a non-eternal card in play. Put a counter on it
// of a type already on it or remove a counter from it" (Marbles). With
// items, only items, and a generic counter may always be added
// (Playdough Cookie); with none, the controller may choose nothing.
func changeACounter(items bool) engine.Effect { return changeCounters(items, nil) }

// changeACounterOr is changeACounter for items: "Put a counter on that
// item or remove one from it. Otherwise, …" runs otherwise on "none"
// (Playdough Cookie).
func changeACounterOr(otherwise engine.Effect) engine.Effect { return changeCounters(true, otherwise) }

func changeCounters(items bool, otherwise engine.Effect) engine.Effect {
	cards := func(g *engine.Game) []engine.ObjectID {
		var out []engine.ObjectID
		for _, id := range cardsInPlay(g) {
			o := g.Object(id)
			if items && o.Role != engine.RoleItem {
				continue
			}
			if items || len(o.Counters) > 0 {
				out = append(out, id)
			}
		}
		return out
	}
	type change struct {
		name string
		n    int
	}
	changes := func(g *engine.Game, id engine.ObjectID) []change {
		o := g.Object(id)
		var out []change
		generic := false
		for _, ct := range o.Counters {
			out = append(out, change{ct.Name, 1}, change{ct.Name, -1})
			generic = generic || ct.Name == ""
		}
		if items && !generic {
			out = append(out, change{"", 1})
		}
		return out
	}
	name := func(n string) string {
		if n == "" {
			return "counter"
		}
		return n + " counter"
	}
	return engine.Ask(func(c *engine.Ctx, a []int) {
		ids := cards(c.G)
		if a[0] < 0 || a[0] >= len(ids) || a[1] < 0 {
			if otherwise != nil {
				c.Do(otherwise)
			}
			return
		}
		ch := changes(c.G, ids[a[0]])[a[1]]
		c.G.AddNamedCounters(ids[a[0]], ch.name, ch.n)
	}, engine.Question{Text: "Change the counters of which card?", Options: func(c *engine.Ctx, _ []int) []string {
		out := labelsOf(c.G, cards(c.G))
		if items && len(out) > 0 {
			out = append(out, "none")
		}
		return out
	}}, engine.Question{Text: "Which change?", Options: func(c *engine.Ctx, a []int) []string {
		ids := cards(c.G)
		if a[0] < 0 || a[0] >= len(ids) {
			return nil
		}
		var out []string
		for _, ch := range changes(c.G, ids[a[0]]) {
			if ch.n > 0 {
				out = append(out, "add a "+name(ch.name))
			} else {
				out = append(out, "remove a "+name(ch.name))
			}
		}
		return out
	}})
}

// killYou kills the controller ("You die").
var killYou = engine.EffectFunc(func(c *engine.Ctx) {
	if !c.G.Players[c.Controller].Dead {
		them := *c
		them.Targets = []engine.Chosen{{Kind: engine.TargetPlayer, Player: c.Controller}}
		them.Do(engine.Kill(0))
	}
})

// discardHandLoot is "discard your hand and loot n".
func discardHandLoot(n int) engine.Effect {
	return engine.EffectFunc(func(c *engine.Ctx) {
		for _, id := range append([]engine.ObjectID(nil), c.G.Players[c.Controller].Hand...) {
			c.G.DiscardFromHand(c.Controller, id)
		}
		c.Do(engine.Loot(n))
	})
}

// level is the name of level counters (R-MECH-41).
const level = "level"

// levelUp puts a level counter on this.
var levelUp = engine.EffectFunc(func(c *engine.Ctx) { c.G.AddNamedCounters(c.Source, level, 1) })

// atStartOfYourTurnAtLevel triggers at the start of the controller's turn
// while this has at least n levels.
func atStartOfYourTurnAtLevel(n int) engine.Trigger {
	return engine.Trigger{On: engine.EvStartOfTurn, Match: func(g *engine.Game, self engine.ObjectID, e engine.Event) bool {
		o := g.Object(self)
		return e.Player == o.Controller && o.CountersOf(level) >= n
	}}
}

// whenCounterReaches triggers when the nth counter is put on this.
func whenCounterReaches(n int) engine.Trigger {
	return engine.Trigger{On: engine.EvCounters, Match: func(g *engine.Game, self engine.ObjectID, e engine.Event) bool {
		if e.Object != self || e.Amount <= 0 || e.Text != "" {
			return false
		}
		now := g.Object(self).CountersOf("")
		return now-e.Amount < n && n <= now
	}}
}

// noCombatDamageTakenOn: "You take no combat damage on attack rolls of …".
func noCombatDamageTakenOn(rolls ...int) func(*engine.Game, engine.ObjectID, int, bool, int) int {
	return func(_ *engine.Game, _ engine.ObjectID, roll int, hit bool, n int) int {
		if !hit && slices.Contains(rolls, roll) {
			return 0
		}
		return n
	}
}

// combatDamageToAMonster triggers when the controller deals combat damage
// to a monster; c.EventObject is the monster.
func combatDamageToAMonster() engine.Trigger {
	return engine.Trigger{On: engine.EvDamaged, Match: func(g *engine.Game, self engine.ObjectID, e engine.Event) bool {
		me := g.Object(self).Controller
		return e.Text == "combat" && e.Player == engine.NoPlayer && e.Object != 0 && e.By == int(me)+1
	}}
}

// damageOtherThanEvent deals 1 damage to a monster or player other than
// the monster of the event ("another monster or player").
var damageOtherThanEvent = engine.Ask(func(c *engine.Ctx, a []int) {
	if ts := targetsBut(c); a[0] >= 0 && a[0] < len(ts) {
		c.G.DealDamageTo(ts[a[0]], 1, c.Controller, c.Source)
	}
}, engine.Question{Text: "Deal 1 damage to?", Options: func(c *engine.Ctx, _ []int) []string {
	var out []string
	for _, t := range targetsBut(c) {
		if t.IsPlayer {
			out = append(out, playerLabel(c.G, t.Player))
		} else {
			out = append(out, string(c.G.Object(t.Object).Card))
		}
	}
	return out
}})

func targetsBut(c *engine.Ctx) []engine.Target {
	var out []engine.Target
	for _, id := range monstersInPlay(c.G) {
		if id != c.EventObject {
			out = append(out, engine.Target{Object: id})
		}
	}
	for _, p := range cardkit.LivingPlayers(c.G, engine.NoPlayer) {
		out = append(out, engine.Target{Player: p, IsPlayer: true})
	}
	return out
}

// soulsBehind is how many souls the player with the most souls has more
// than p.
func soulsBehind(g *engine.Game, p engine.PlayerID) int {
	most := 0
	for _, pl := range g.Players {
		most = max(most, g.SoulValue(pl.ID))
	}
	return most - g.SoulValue(p)
}

// revealTenTakeNamed reveals the top 10 loot cards; the controller takes
// those match picks, the rest go to the bottom of the loot deck.
func revealTenTakeNamed(match func(engine.CardRef) bool) engine.Effect {
	return engine.EffectFunc(func(c *engine.Ctx) {
		top := c.G.DeckTop(engine.LootDeck, 10)
		for _, pl := range c.G.Players {
			c.G.LookAt(pl.ID, top...)
		}
		for _, id := range top {
			if match(c.G.Object(id).Card) {
				c.G.DeckToHand(engine.LootDeck, id, c.Controller)
			} else {
				c.G.DeckToBottom(engine.LootDeck, id)
			}
		}
	})
}

// lootDiscardToHand puts a card that match picks from the loot discard
// into the controller's hand; the controller chooses which.
func lootDiscardToHand(match func(engine.CardRef) bool) engine.Effect {
	cards := func(g *engine.Game) []engine.ObjectID {
		var out []engine.ObjectID
		for _, id := range g.Discards[engine.LootDeck] {
			if match(g.Object(id).Card) {
				out = append(out, id)
			}
		}
		return out
	}
	return engine.Ask(func(c *engine.Ctx, a []int) {
		if ids := cards(c.G); a[0] >= 0 && a[0] < len(ids) {
			c.G.DiscardToHand(engine.LootDeck, ids[a[0]], c.Controller)
		}
	}, engine.Question{Text: "Take which card from the loot discard?", Options: func(c *engine.Ctx, _ []int) []string {
		return labelsOf(c.G, cards(c.G))
	}})
}

// mayDeckQuestion asks for a deck, or "none"; check the answer with
// deckChosen before reading it with c.Deck.
func mayDeckQuestion(text string) engine.Question {
	q := engine.DeckQuestion(text)
	opts := q.Options
	q.Options = func(c *engine.Ctx, a []int) []string { return append(opts(c, a), "none") }
	return q
}

// deckChosen: answer a picked a deck, not "none".
func deckChosen(c *engine.Ctx, a int) bool { return a >= 0 && a < len(c.G.DecksInPlay()) }

// othersOfActive lists the players other than the active one, in turn
// order from the active player.
func othersOfActive(g *engine.Game) []engine.PlayerID {
	return g.TurnOrderFrom(g.Turn.Active)[1:]
}

// playersCanBeAttacked: "The active player may attack other players.
// Attacked players have n+ DC" (R-ATK-20).
func playersCanBeAttacked(n int) engine.Static {
	return engine.Static{Stat: engine.StatPlayerAttackDC, Amount: n, Applies: func(*engine.Game, engine.ObjectID, engine.PlayerID, engine.ObjectID) bool {
		return true
	}}
}

// whenAttackedThis triggers when an attack is declared on this, also
// when it was revealed by an attack on the deck (R-ATK-06).
func whenAttackedThis() engine.Trigger {
	return engine.Trigger{On: engine.EvAttackTarget, Match: func(_ *engine.Game, self engine.ObjectID, e engine.Event) bool {
		return e.Object == self
	}}
}

// whenYouKillAPlayer triggers when another player the controller killed
// dies.
func whenYouKillAPlayer() engine.Trigger {
	return engine.Trigger{On: engine.EvDied, Match: func(g *engine.Game, self engine.ObjectID, e engine.Event) bool {
		me := g.Object(self).Controller
		return e.Player != engine.NoPlayer && e.Player != me && g.Players[e.Player].KilledBy == int(me)+1
	}}
}

// markHelper remembers, till end of turn, the player this card
// chose (Sisters Vis); markedHelper reads it.
func markHelper(c *engine.Ctx, p engine.PlayerID) {
	cardkit.MarkUsedOn(c.G, c.Source, int(p)+1)
}

func markedHelper(c *engine.Ctx) (engine.PlayerID, bool) {
	for _, b := range c.G.Boosts {
		if b.Stat == engine.StatUsed && b.Object == c.Source && b.Amount > 0 {
			return engine.PlayerID(b.Amount - 1), true
		}
	}
	return engine.NoPlayer, false
}

// randomPlayerTakes: this deals n damage to a living player chosen at
// random; non-active players only if nonActive.
func randomPlayerTakes(n int, nonActive bool) engine.Effect {
	players := func(g *engine.Game) []engine.PlayerID {
		skip := engine.NoPlayer
		if nonActive {
			skip = g.Turn.Active
		}
		return cardkit.LivingPlayers(g, skip)
	}
	return engine.Ask(func(c *engine.Ctx, a []int) {
		if ps := players(c.G); a[0] >= 0 && a[0] < len(ps) {
			c.G.DealDamageTo(engine.Target{Player: ps[a[0]], IsPlayer: true}, n, engine.NoPlayer, c.Source)
		}
	}, engine.Question{Text: "A player at random", Random: true, Options: func(c *engine.Ctx, _ []int) []string {
		return playerLabels(c.G, players(c.G))
	}})
}

// damageActive: this deals n damage to the active player.
func damageActive(n int) engine.Effect {
	return engine.EffectFunc(func(c *engine.Ctx) {
		c.G.DealDamageTo(engine.Target{Player: c.G.Turn.Active, IsPlayer: true}, n, engine.NoPlayer, c.Source)
	})
}

// afterRewardsAnotherPlayer is the charmed monsters' "When this dies,
// after gaining rewards, choose another player. They …": the active
// player chooses; the effects run as that player's.
func afterRewardsAnotherPlayer(they string, effects ...engine.Effect) engine.Ability {
	return engine.Ability{
		Kind:    engine.Triggered,
		Text:    "When this dies, after gaining rewards, choose another player. " + they,
		Trigger: engine.AfterThisRewards(),
		Effects: []engine.Effect{engine.Ask(func(c *engine.Ctx, a []int) {
			ps := othersOfActive(c.G)
			if a[0] < 0 || a[0] >= len(ps) {
				return
			}
			them := *c
			them.Controller = ps[a[0]]
			them.Do(effects...)
		}, engine.Question{Text: "Choose another player.", Options: func(c *engine.Ctx, _ []int) []string {
			return playerLabels(c.G, othersOfActive(c.G))
		}})},
	}
}

// rechargeYourItem: the controller recharges a deactivated item they
// control.
var rechargeYourItem = engine.Ask(func(c *engine.Ctx, a []int) {
	if ids := deactivatedItems(c.G, c.Controller); a[0] >= 0 && a[0] < len(ids) {
		c.G.Recharge(ids[a[0]])
	}
}, engine.Question{Text: "Recharge which item?", Options: func(c *engine.Ctx, _ []int) []string {
	return labelsOf(c.G, deactivatedItems(c.G, c.Controller))
}})

func deactivatedItems(g *engine.Game, p engine.PlayerID) []engine.ObjectID {
	var out []engine.ObjectID
	for _, id := range itemsOf(g, p) {
		if !g.Object(id).Charged {
			out = append(out, id)
		}
	}
	return out
}

// centsEqualToRoll: the roller gains ¢ equal to the result.
func centsEqualToRoll() engine.RollTable {
	var t engine.RollTable
	for r := 1; r <= 6; r++ {
		t = t.Results(r, r, engine.GainCents(r))
	}
	return t
}

// damageNeighbour: this deals 1 damage to the player to the active
// player's left (next in turn order, R-TURN-01) or right.
func damageNeighbour(left bool) engine.Effect {
	return engine.EffectFunc(func(c *engine.Ctx) {
		p := c.G.PlayerBefore(c.G.Turn.Active)
		if left {
			p = c.G.PlayerAfter(c.G.Turn.Active)
		}
		c.G.DealDamageTo(engine.Target{Player: p, IsPlayer: true}, 1, engine.NoPlayer, c.Source)
	})
}

// thisGains gives this monster n of a stat till end of turn.
func thisGains(stat engine.Stat, n int) engine.Effect {
	return engine.EffectFunc(func(c *engine.Ctx) {
		c.G.Boosts = append(c.G.Boosts, engine.Boost{Stat: stat, Player: engine.NoPlayer, Object: c.Source, Amount: n})
	})
}

// killUpTo: the active player kills up to n other players, one at a
// time; "none" stops (S-YUG: up to includes none).
func killUpTo(n int) engine.Effect {
	left := func(c *engine.Ctx, a []int) []engine.PlayerID {
		ps := othersOfActive(c.G)
		for _, x := range a {
			if x < 0 || x >= len(ps) {
				return nil
			}
			ps = append(ps[:x:x], ps[x+1:]...)
		}
		return ps
	}
	qs := make([]engine.Question, n)
	for i := range qs {
		qs[i] = engine.Question{Text: "Kill which player?", Options: func(c *engine.Ctx, a []int) []string {
			ps := left(c, a[:i])
			if len(ps) == 0 {
				return nil
			}
			return append(playerLabels(c.G, ps), "none")
		}}
	}
	return engine.Ask(func(c *engine.Ctx, a []int) {
		ps := othersOfActive(c.G)
		var picked []engine.PlayerID
		for _, x := range a {
			if x < 0 || x >= len(ps) {
				break
			}
			picked = append(picked, ps[x])
			ps = append(ps[:x:x], ps[x+1:]...)
		}
		for _, p := range picked {
			them := *c
			them.Targets = []engine.Chosen{{Kind: engine.TargetPlayer, Player: p}}
			them.Do(engine.Kill(0))
		}
	}, qs...)
}

// cursesInDeck lists the curse cards in the monster deck.
func cursesInDeck(g *engine.Game) []engine.ObjectID {
	var out []engine.ObjectID
	for _, id := range g.Decks[engine.MonsterDeck] {
		if g.Def(id).Curse {
			out = append(out, id)
		}
	}
	return out
}

// slotQuestionAfterPick asks for a free monster slot once the question
// before it picked something.
func slotQuestionAfterPick(text string) engine.Question {
	q := slotQuestion(text)
	opts := q.Options
	q.Options = func(c *engine.Ctx, a []int) []string {
		if a[len(a)-1] < 0 {
			return nil
		}
		return opts(c, a)
	}
	return q
}

// vote: each living player, in turn order from the active player, votes
// for one option (R-MECH-49); decide gets the votes per option.
func vote(options []string, decide func(c *engine.Ctx, votes []int)) engine.Effect {
	voters := func(c *engine.Ctx) []engine.PlayerID { return cardkit.LivingPlayers(c.G, engine.NoPlayer) }
	inOrder := func(c *engine.Ctx) []engine.PlayerID {
		var out []engine.PlayerID
		for _, p := range c.G.TurnOrderFrom(c.G.Turn.Active) {
			if slices.Contains(voters(c), p) {
				out = append(out, p)
			}
		}
		return out
	}
	qs := make([]engine.Question, 4)
	for i := range qs {
		qs[i] = engine.Question{
			Text:   "Vote: " + strings.Join(options, " or ") + "?",
			Player: func(c *engine.Ctx, _ []int) engine.PlayerID { return inOrder(c)[i] },
			Options: func(c *engine.Ctx, _ []int) []string {
				if i >= len(inOrder(c)) {
					return nil
				}
				return options
			},
		}
	}
	return engine.Ask(func(c *engine.Ctx, a []int) {
		votes := make([]int, len(options))
		for _, x := range a {
			if x >= 0 && x < len(votes) {
				votes[x]++
			}
		}
		decide(c, votes)
	}, qs...)
}

// rerollAnyNumber: the controller rerolls any number of their
// non-eternal items, picked one at a time (up to max).
func rerollAnyNumber(most int) engine.Effect {
	items := func(c *engine.Ctx) []engine.ObjectID {
		var out []engine.ObjectID
		for _, id := range itemsOf(c.G, c.Controller) {
			if !c.G.Eternal(id) {
				out = append(out, id)
			}
		}
		return out
	}
	left := func(c *engine.Ctx, a []int) ([]engine.ObjectID, bool) {
		ids := items(c)
		for _, x := range a {
			if x < 0 || x >= len(ids) {
				return nil, true
			}
			ids = append(ids[:x:x], ids[x+1:]...)
		}
		return ids, false
	}
	qs := make([]engine.Question, most)
	for i := range qs {
		qs[i] = engine.Question{Text: "Reroll which item?", Options: func(c *engine.Ctx, a []int) []string {
			ids, done := left(c, a[:i])
			if done || len(ids) == 0 {
				return nil
			}
			return append(labelsOf(c.G, ids), "done")
		}}
	}
	return engine.Ask(func(c *engine.Ctx, a []int) {
		ids := items(c)
		var picked []engine.ObjectID
		for _, x := range a {
			if x < 0 || x >= len(ids) {
				break
			}
			picked = append(picked, ids[x])
			ids = append(ids[:x:x], ids[x+1:]...)
		}
		for _, id := range picked {
			c.G.RerollItem(id)
		}
	}, qs...)
}

// damageAChosenPlayer is "Choose a player. They take n damage.".
func damageAChosenPlayer(n int) engine.Effect {
	return engine.Ask(func(c *engine.Ctx, a []int) {
		if ps := cardkit.LivingPlayers(c.G, engine.NoPlayer); a[0] >= 0 && a[0] < len(ps) {
			c.G.DealDamageTo(engine.Target{Player: ps[a[0]], IsPlayer: true}, n, c.Controller, c.Source)
		}
	}, cardkit.PickAPlayer("Choose a player."))
}

// backIntoDeck puts this event into the monster deck n cards from the
// top; its slot is refilled.
func backIntoDeck(n int) engine.Effect {
	return engine.EffectFunc(func(c *engine.Ctx) {
		if c.G.Object(c.Source).Zone.Kind == engine.ZoneInPlay {
			c.G.PutIntoDeck(engine.MonsterDeck, c.Source, n)
			c.G.RefillSlots()
		}
	})
}

// monstersAndPlayers lists the monsters in play and the living players,
// as chosen targets.
func monstersAndPlayers(g *engine.Game) []engine.Chosen {
	var out []engine.Chosen
	for _, id := range monstersInPlay(g) {
		out = append(out, engine.Chosen{Kind: engine.TargetMonster, Player: engine.NoPlayer, Object: id})
	}
	for _, p := range cardkit.LivingPlayers(g, engine.NoPlayer) {
		out = append(out, engine.Chosen{Kind: engine.TargetPlayer, Player: p})
	}
	return out
}

// chosenLabels names chosen targets as the prompts do.
func chosenLabels(g *engine.Game, cs []engine.Chosen) []string {
	out := make([]string, len(cs))
	for i, c := range cs {
		if c.Kind == engine.TargetPlayer {
			out[i] = playerLabel(g, c.Player)
		} else {
			out[i] = string(g.Object(c.Object).Card)
		}
	}
	return out
}

// discardN: the controller discards n loot cards of their choice.
func discardN(n int) engine.Effect {
	qs := make([]engine.Question, n)
	for i := range qs {
		qs[i] = engine.Question{Text: "Discard which loot card?", Options: func(c *engine.Ctx, a []int) []string {
			return labelsOf(c.G, handLeft(c, a[:i]))
		}}
	}
	return engine.Ask(func(c *engine.Ctx, a []int) {
		left := append([]engine.ObjectID(nil), c.G.Players[c.Controller].Hand...)
		for _, x := range a {
			if x < 0 || x >= len(left) {
				return
			}
			c.G.DiscardFromHand(c.Controller, left[x])
			left = append(left[:x:x], left[x+1:]...)
		}
	}, qs...)
}

// mustAttackThis is "The active player must attack this once each turn
// if able".
var mustAttackThis = engine.EffectFunc(func(c *engine.Ctx) {
	if c.G.Object(c.Source).Zone.Kind == engine.ZoneInPlay && (!c.G.Attack.Started || c.G.Attack.Target != c.Source) {
		c.G.Turn.MustAttack = c.Source
	}
})

// lowestRollDies: each player rolls; each player who rolls the lowest or
// tied for the lowest dies (The Beast!). The last result to resolve
// compares them; results are kept as marks on this card.
var lowestRollDies = engine.EachPlayer(engine.Roll(lowestTable()))

func lowestTable() engine.RollTable {
	var t engine.RollTable
	for r := 1; r <= 6; r++ {
		t = t.Results(r, r, engine.EffectFunc(func(c *engine.Ctx) {
			cardkit.MarkUsedOn(c.G, c.Source, 100+10*r+int(c.Controller))
			results := map[engine.PlayerID]int{}
			for _, b := range c.G.Boosts {
				if b.Stat == engine.StatUsed && b.Object == c.Source && b.Amount >= 100 {
					results[engine.PlayerID((b.Amount-100)%10)] = (b.Amount - 100) / 10
				}
			}
			if len(results) < len(c.G.Players) {
				return
			}
			low := 7
			for _, v := range results {
				low = min(low, v)
			}
			// Deaths go on the stack in reverse turn order (R-MECH-28).
			order := c.G.TurnOrderFrom(c.G.Turn.Active)
			for i := len(order) - 1; i >= 0; i-- {
				if p := order[i]; results[p] == low {
					them := *c
					them.Targets = []engine.Chosen{{Kind: engine.TargetPlayer, Player: p}}
					them.Do(engine.Kill(0))
				}
			}
			clearMarks(c.G, c.Source)
		}))
	}
	return t
}

// clearMarks removes the roll marks lowestTable keeps on id.
func clearMarks(g *engine.Game, id engine.ObjectID) {
	kept := g.Boosts[:0]
	for _, b := range g.Boosts {
		if b.Stat != engine.StatUsed || b.Object != id || b.Amount < 100 {
			kept = append(kept, b)
		}
	}
	g.Boosts = kept
}

// othersOf lists the players other than p, in turn order from p.
func othersOf(g *engine.Game, p engine.PlayerID) []engine.PlayerID {
	return g.TurnOrderFrom(p)[1:]
}

// otherMonsters lists the monsters in play but m.
func otherMonsters(g *engine.Game, m engine.ObjectID) []engine.ObjectID {
	var out []engine.ObjectID
	for _, id := range monstersInPlay(g) {
		if id != m {
			out = append(out, id)
		}
	}
	return out
}

// nonEternalItems lists the non-eternal items p controls.
func nonEternalItems(g *engine.Game, p engine.PlayerID) []engine.ObjectID {
	var out []engine.ObjectID
	for _, id := range itemsOf(g, p) {
		if !g.Eternal(id) {
			out = append(out, id)
		}
	}
	return out
}

// damageUpTo deals n damage to each of up to most different monsters or
// players, picked one at a time.
func damageUpTo(most, n int) engine.Effect {
	left := func(c *engine.Ctx, a []int) ([]engine.Chosen, bool) {
		ts := monstersAndPlayers(c.G)
		for _, x := range a {
			if x < 0 || x >= len(ts) {
				return nil, true
			}
			ts = append(ts[:x:x], ts[x+1:]...)
		}
		return ts, false
	}
	qs := make([]engine.Question, most)
	for i := range qs {
		qs[i] = engine.Question{Text: "Deal damage to?", Options: func(c *engine.Ctx, a []int) []string {
			ts, done := left(c, a[:i])
			if done || len(ts) == 0 {
				return nil
			}
			return append(chosenLabels(c.G, ts), "done")
		}}
	}
	return engine.Ask(func(c *engine.Ctx, a []int) {
		ts := monstersAndPlayers(c.G)
		for _, x := range a {
			if x < 0 || x >= len(ts) {
				return
			}
			t := engine.Target{Object: ts[x].Object}
			if ts[x].Kind == engine.TargetPlayer {
				t = engine.Target{Player: ts[x].Player, IsPlayer: true}
			}
			c.G.DealDamageTo(t, n, c.Controller, c.Source)
			ts = append(ts[:x:x], ts[x+1:]...)
		}
	}, qs...)
}
