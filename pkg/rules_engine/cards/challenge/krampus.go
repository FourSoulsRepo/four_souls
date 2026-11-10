package challenge

import (
	"fmt"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/cardkit"
)

// Krampus is How the Krampus Stole Christmas (ADR 028; R-KRAMPUS). The
// boss cards hold the per-difficulty rules (R-KRAMPUS-11 to 14); the
// rule cards hold the three global modifiers and the goal (R-KRAMPUS-01,
// 07, 08, 10).
var Krampus = engine.ChallengeDef{
	ID:  "krampus",
	Set: Set.Code,
	Levels: []engine.ChallengeLevel{
		{Difficulty: engine.DifficultyNormal, Boss: "pilfering_krampus_norm", Card: "how_the_krampus_stole_christmas_norm"},
		{Difficulty: engine.DifficultyHard, Boss: "pilfering_krampus_hard", Card: "how_the_krampus_stole_christmas_hard"},
		{Difficulty: engine.DifficultyUltra, Boss: "pilfering_krampus_ultra", Card: "how_the_krampus_stole_christmas_ultra"},
	},
	Setup: krampusSetup,
}

// The challenge's own lines; their words are in en.json (ADR 021).
var (
	// ReasonPresents: all 4 presents were destroyed (R-KRAMPUS-01).
	ReasonPresents = engine.NewKey("event.challenge_lost.presents")
	// forgoAsk asks whether to forgo a reward or a soul (R-KRAMPUS-07):
	// {n} counters; the soul's {card} too.
	forgoAsk = map[string]engine.Key{
		engine.PresentCents:    engine.NewKey("card.how_the_krampus_stole_christmas.forgo_cents"),
		engine.PresentLoot:     engine.NewKey("card.how_the_krampus_stole_christmas.forgo_loot"),
		engine.PresentTreasure: engine.NewKey("card.how_the_krampus_stole_christmas.forgo_treasure"),
		engine.PresentSoul:     engine.NewKey("card.how_the_krampus_stole_christmas.forgo_soul"),
	}
	// giftAsk asks who gains a flipped present (R-KRAMPUS-10).
	giftAsk = map[string]engine.Key{
		engine.PresentCents:    engine.NewKey("card.how_the_krampus_stole_christmas.gift_cents"),
		engine.PresentLoot:     engine.NewKey("card.how_the_krampus_stole_christmas.gift_loot"),
		engine.PresentTreasure: engine.NewKey("card.how_the_krampus_stole_christmas.gift_treasure"),
		engine.PresentSoul:     engine.NewKey("card.how_the_krampus_stole_christmas.gift_soul"),
	}
)

// presentKinds are the presents in their order: ¢, loot, treasure, soul
// (R-KRAMPUS-05, Q-40).
var presentKinds = []string{engine.PresentCents, engine.PresentLoot, engine.PresentTreasure, engine.PresentSoul}

// krampusLevel is what a difficulty changes (R-KRAMPUS-02).
type krampusLevel struct {
	boss engine.CardRef
	// souls: Happy Holidays! and the boss's protection (R-KRAMPUS-10, 11).
	souls int
	// minions come when a present is destroyed (R-KRAMPUS-08).
	minions int
	// counters on the presents at the start, in presentKinds order
	// (R-KRAMPUS-05).
	counters [4]int
}

var krampusLevels = map[string]krampusLevel{
	engine.DifficultyNormal: {boss: "pilfering_krampus_norm", souls: 5, minions: 3, counters: [4]int{8, 2, 3, 3}},
	engine.DifficultyHard:   {boss: "pilfering_krampus_hard", souls: 6, minions: 4, counters: [4]int{12, 4, 3, 3}},
	engine.DifficultyUltra:  {boss: "pilfering_krampus_ultra", souls: 7, minions: 4, counters: [4]int{16, 4, 3, 3}},
}

// krampusSetup: Pilfering Krampus of the difficulty goes into a final
// boss slot; the top 4 treasure cards become the presents, face down,
// with their counters (R-KRAMPUS-04, R-KRAMPUS-05).
func krampusSetup(g *engine.Game, difficulty string) error {
	l, ok := krampusLevels[difficulty]
	if !ok {
		return fmt.Errorf("no Krampus difficulty %q", difficulty)
	}
	if g.PutInFinalBossSlot(l.boss) == 0 {
		return fmt.Errorf("%s is no monster of the game", l.boss)
	}
	for i, kind := range presentKinds {
		if !g.AddPresent(kind, l.counters[i]) {
			return fmt.Errorf("the treasure deck has no card for the %s present", kind)
		}
	}
	return nil
}

// minions lists the minions in play (R-CHAL-08).
func minions(g *engine.Game) []engine.ObjectID {
	var out []engine.ObjectID
	for _, id := range g.MonstersInPlay() {
		if g.IsMinion(id) {
			out = append(out, id)
		}
	}
	return out
}

// --- Pilfering Krampus ---

// guardedBelow: "While the number of souls controlled by all players is
// less than n or there is a minion in play, this can't be attacked or
// take damage" (R-KRAMPUS-11, Q-30).
func guardedBelow(n int) func(g *engine.Game, self engine.ObjectID) bool {
	return func(g *engine.Game, _ engine.ObjectID) bool {
		return g.SoulsInPlay() < n || len(minions(g)) > 0
	}
}

// randomMinion: "While there is a minion in play, players must attack a
// random minion once each turn if able" (Normal, R-KRAMPUS-13, Q-31).
func randomMinion(g *engine.Game, _ engine.ObjectID) []engine.ObjectID { return minions(g) }

// minionSpoils: "Players gain no rewards or souls when a minion dies"
// (R-KRAMPUS-14). The boss card prints it, so it is the boss's.
func minionSpoils(g *engine.Game, _, dying engine.ObjectID) bool { return g.IsMinion(dying) }

// takeCounters removes n counters from the present of kind.
func takeCounters(kind string, n int) engine.Effect {
	return engine.EffectFunc(func(c *engine.Ctx) { c.G.AddPresentCounters(kind, -n, engine.NoPlayer) })
}

// presentRoll is "As each turn starts, the active player rolls" with
// its table and its extra rolls (R-KRAMPUS-12). extra counts the extra
// rolls from the minions in play, once, as it resolves (Q-41).
func presentRoll(text string, extra func(minions int) int, table engine.RollTable) engine.Ability {
	return engine.Ability{
		Kind:    engine.Triggered,
		Text:    text,
		Trigger: engine.Trigger{On: engine.EvStartOfTurn, Match: func(*engine.Game, engine.ObjectID, engine.Event) bool { return true }},
		Effects: []engine.Effect{engine.RollTimes(func(c *engine.Ctx) int { return 1 + extra(len(minions(c.G))) }, table)},
	}
}

func noExtraRolls(int) int { return 0 }

// --- How the Krampus Stole Christmas ---

// forgoable is what a would-be gain could be forgone for: the present's
// kind and the counters (R-KRAMPUS-07, Q-32, Q-33). Only a monster's
// rewards and a dead card's soul; only for a present still there.
func forgoable(g *engine.Game, a engine.Action) (kind string, n int, ok bool) {
	switch {
	case a.Kind == engine.ActGainCents && a.Reward:
		kind, n = engine.PresentCents, a.Amount
	case a.Kind == engine.ActLoot && a.Reward:
		kind, n = engine.PresentLoot, a.Amount
	case a.Kind == engine.ActGainTreasure && a.Reward:
		kind, n = engine.PresentTreasure, a.Amount
	case a.Kind == engine.ActBecomeSoul:
		if g.Object(a.Object).Zone.Kind != engine.ZoneOutside || !g.CanGainSouls(a.Player) {
			return "", 0, false
		}
		kind, n = engine.PresentSoul, max(g.Def(a.Object).Soul, 1)
	default:
		return "", 0, false
	}
	return kind, n, n > 0 && a.Player != engine.NoPlayer && g.PresentLive(kind)
}

// allIWant is "All I Want for Christmas": an optional replacement of a
// monster's reward part or a soul (R-KRAMPUS-07). A forgone soul card
// goes to discard (Q-32).
var allIWant = engine.Replacement{
	Text: "If you would gain rewards from a monster or gain a soul, you may instead forego the reward or soul and put counters on the present of the respective type equal to the amount of that resource foregone.",
	When: func(g *engine.Game, _ engine.ObjectID, a engine.Action) bool {
		_, _, ok := forgoable(g, a)
		return ok
	},
	Do: func(g *engine.Game, _ engine.ObjectID, a engine.Action) []engine.Action {
		kind, n, ok := forgoable(g, a)
		if !ok {
			return []engine.Action{a}
		}
		g.Forgo(kind, n, a.Player)
		if a.Kind == engine.ActBecomeSoul {
			return []engine.Action{{Kind: engine.ActDiscardObject, Player: engine.NoPlayer, Object: a.Object}}
		}
		return nil
	},
	May: func(g *engine.Game, _ engine.ObjectID, a engine.Action) engine.Msg {
		kind, n, _ := forgoable(g, a)
		if kind == engine.PresentSoul {
			return forgoAsk[kind].Msg(engine.N("n", n), engine.CardArg("card", g.Object(a.Object).Card))
		}
		return forgoAsk[kind].Msg(engine.N("n", n))
	},
}

// whatsInTheBox is "What's in the Box?": a present reduced to 0 counters
// is destroyed; then n cards of the monster deck go into minion slots,
// one at a time (R-KRAMPUS-08). The fourth destroyed loses (R-KRAMPUS-09).
func whatsInTheBox(text string, n int) engine.Ability {
	return engine.Ability{
		Kind: engine.Triggered,
		Text: text,
		Trigger: engine.Trigger{On: engine.EvPresentCounters, Match: func(g *engine.Game, _ engine.ObjectID, e engine.Event) bool {
			return e.Amount < 0 && g.PresentLive(e.Text) && g.PresentCounters(e.Text) == 0
		}},
		Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
			kind := c.G.PresentOf(c.EventObject)
			if kind == "" || c.G.PresentCounters(kind) > 0 {
				return // counters came back before this resolved
			}
			done, all := c.G.DestroyPresent(kind)
			if !done {
				return
			}
			if all {
				c.G.LoseChallenge(ReasonPresents.Msg())
				return
			}
			for range n {
				c.G.MinionFromDeck()
			}
		})},
	}
}

// choosable lists who a present may go to: the living players, or all
// if none lives (R-ABIL-35).
func choosable(g *engine.Game) []engine.PlayerID {
	if ps := cardkit.LivingPlayers(g, engine.NoPlayer); len(ps) > 0 {
		return ps
	}
	var out []engine.PlayerID
	for _, pl := range g.Players {
		out = append(out, pl.ID)
	}
	return out
}

// happyHolidays is "Happy Holidays!": the first time the players hold n
// souls between them, every present still there flips and becomes an
// item; the active player picks who gains each, in order ¢, loot,
// treasure, soul (R-KRAMPUS-10, Q-30, Q-40).
func happyHolidays(text string, n int) engine.Ability {
	var questions []engine.Question
	for _, kind := range presentKinds {
		questions = append(questions, engine.Question{
			Text: giftAsk[kind],
			Options: func(c *engine.Ctx, _ []int) []engine.Msg {
				if c.G.PresentToGive(kind) == 0 {
					return nil
				}
				return cardkit.PlayerLabels(c.G, choosable(c.G))
			},
			Subject: func(c *engine.Ctx, _ []int) engine.ObjectID { return c.G.PresentToGive(kind) },
		})
	}
	return engine.Ability{
		Kind: engine.Triggered,
		Text: text,
		Trigger: engine.Trigger{On: engine.EvGainedSoul, Match: func(g *engine.Game, _ engine.ObjectID, _ engine.Event) bool {
			return !g.Challenge.Holidays && g.SoulsInPlay() >= n
		}},
		Effects: []engine.Effect{
			engine.EffectFunc(func(c *engine.Ctx) { c.G.FlipPresents() }),
			engine.Ask(func(c *engine.Ctx, a []int) {
				ps := choosable(c.G)
				for i, kind := range presentKinds {
					if a[i] >= 0 && a[i] < len(ps) {
						c.G.GivePresent(kind, ps[a[i]])
					}
				}
			}, questions...),
		},
	}
}

// krampusGoal: the challenge is won when Pilfering Krampus dies
// (R-KRAMPUS-01).
func krampusGoal(boss engine.CardRef) engine.Ability {
	return engine.Ability{
		Kind: engine.Triggered,
		Text: "When Pilfering Krampus dies, the players win.",
		Trigger: engine.Trigger{On: engine.EvDied, Match: func(_ *engine.Game, _ engine.ObjectID, e engine.Event) bool {
			return e.Player == engine.NoPlayer && e.Card == boss
		}},
		Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) { c.G.WinChallenge() })},
	}
}

// ruleCard is a difficulty's rule card abilities.
func ruleCard(d engine.CardDef, level string, box, holidays string) engine.CardDef {
	l := krampusLevels[level]
	d.Replacements = []engine.Replacement{allIWant}
	d.Abilities = []engine.Ability{
		whatsInTheBox(box, l.minions),
		happyHolidays(holidays, l.souls),
		krampusGoal(l.boss),
	}
	return d
}
