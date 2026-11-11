package challenge

import (
	"fmt"
	"slices"

	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards/cardkit"
)

// TrickOrTreat is Trick or Treat (ADR 028; R-TRICK). The boss cards
// hold the guard, the fear counters and the death roll (R-TRICK-06 to
// 10); the two-sided rule card holds the flip, Trick's and Treat's
// abilities, and the goal (R-TRICK-01, 11 to 15). One rule card serves
// every difficulty (R-TRICK-02); Ghostly Trickery reads the difficulty.
var TrickOrTreat = engine.ChallengeDef{
	ID:  "trick_or_treat",
	Set: Set.Code,
	Levels: []engine.ChallengeLevel{
		{Difficulty: engine.DifficultyNormal, Boss: "mischievous_the_haunt_norm", Card: "trick_treat_normhardultra"},
		{Difficulty: engine.DifficultyHard, Boss: "mischievous_the_haunt_hard", Card: "trick_treat_normhardultra"},
		{Difficulty: engine.DifficultyUltra, Boss: "mischievous_the_haunt_ultra", Card: "trick_treat_normhardultra"},
	},
	Setup: trickSetup,
}

// The counters of Trick or Treat: trick on the boss, fear on players
// (R-TRICK-03, R-TRICK-07, Q-45).
const (
	TrickCounter = "trick"
	FearCounter  = "fear"
)

// The challenge's questions; their words are in en.json (ADR 021).
var (
	// removeFearAsk: which player loses a fear counter (Q-51).
	removeFearAsk = engine.NewKey("card.trick_treat.remove_fear")
	// flipAsk: flip the rule card as the turn starts (Q-43); one per
	// side, naming the side it would show.
	flipToTreatAsk = engine.NewKey("card.trick_treat.flip_to_treat")
	flipToTrickAsk = engine.NewKey("card.trick_treat.flip_to_trick")
)

// hauntLevel is what a difficulty changes (R-TRICK-03, R-TRICK-07,
// R-TRICK-13).
type hauntLevel struct {
	boss engine.CardRef
	// tricks: trick counters on the boss at the start.
	tricks int
	// rolls: Ghostly Trickery's rolls at each attack start.
	rolls int
}

var hauntLevels = map[string]hauntLevel{
	engine.DifficultyNormal: {boss: "mischievous_the_haunt_norm", tricks: 12, rolls: 1},
	engine.DifficultyHard:   {boss: "mischievous_the_haunt_hard", tricks: 14, rolls: 2},
	engine.DifficultyUltra:  {boss: "mischievous_the_haunt_ultra", tricks: 14, rolls: 3},
}

// hauntBosses are the boss cards of every difficulty, for the goal.
var hauntBosses = []engine.CardRef{"mischievous_the_haunt_norm", "mischievous_the_haunt_hard", "mischievous_the_haunt_ultra"}

// trickSetup: Mischievous The Haunt of the difficulty goes into a final
// boss slot with its trick counters; the rule card lies Trick side up
// (R-TRICK-04).
func trickSetup(g *engine.Game, difficulty string) error {
	l, ok := hauntLevels[difficulty]
	if !ok {
		return fmt.Errorf("no Trick or Treat difficulty %q", difficulty)
	}
	b := g.PutInFinalBossSlot(l.boss)
	if b == 0 {
		return fmt.Errorf("%s is no monster of the game", l.boss)
	}
	g.AddNamedCounters(b, TrickCounter, l.tricks)
	return nil
}

// finalBosses lists the final bosses in play (R-CHAL-02).
func finalBosses(g *engine.Game) []engine.ObjectID {
	var out []engine.ObjectID
	for _, s := range g.Bosses {
		if id, ok := s.TopOf(); ok {
			out = append(out, id)
		}
	}
	return out
}

// removeFear: "Remove a fear counter from a player". The active player
// (the controller) picks one with a fear counter, themselves included;
// with none anywhere, nothing happens (Q-51).
var removeFear = engine.Ask(func(c *engine.Ctx, a []int) {
	ps := c.G.PlayersWithCounters(FearCounter)
	if a[0] >= 0 && a[0] < len(ps) {
		c.G.AddPlayerCounters(ps[a[0]], FearCounter, -1)
	}
}, engine.Question{Text: removeFearAsk, Options: func(c *engine.Ctx, _ []int) []engine.Msg {
	return cardkit.PlayerLabels(c.G, c.G.PlayersWithCounters(FearCounter))
}})

// --- Mischievous The Haunt ---

// trickGuard: "While there are more trick counters on this than the
// number of souls controlled by all players, attacks can't be declared
// on this and it can't take damage" (R-TRICK-06, Q-30).
func trickGuard(g *engine.Game, self engine.ObjectID) bool {
	return g.Object(self).CountersOf(TrickCounter) > g.SoulsInPlay()
}

// fearEvasion: "For every 2 fear counters on the active player,
// monsters have +1 DC" (R-TRICK-09, Q-52).
var fearEvasion = engine.MonstersHavePerCounters(engine.StatMonsterDC, FearCounter, 2)

// fearOn: "Each time a player takes damage, put a fear counter on
// them", or on each player (Ultra Hard). One per instance (R-TRICK-07,
// Q-50).
func fearOn(text string, each bool) engine.Ability {
	return engine.Ability{
		Kind:    engine.Triggered,
		Text:    text,
		Trigger: engine.WhenAPlayerTakesDamage(),
		Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
			if !each {
				c.G.AddPlayerCounters(c.EventPlayer, FearCounter, 1)
				return
			}
			for _, pl := range c.G.Players {
				c.G.AddPlayerCounters(pl.ID, FearCounter, 1)
			}
		})},
	}
}

// deathRoll: "Each time a monster dies, the active player rolls- 4-6:
// Remove a fear counter from a player" (R-TRICK-10, Q-51).
var deathRoll = engine.Ability{
	Kind:    engine.Triggered,
	Text:    "Each time a monster dies, the active player rolls- 4-6: Remove a fear counter from a player.",
	Trigger: engine.WhenAMonsterDies(),
	Effects: []engine.Effect{engine.Roll(engine.RollTable{}.Results(4, 6, removeFear))},
}

// haunt is a difficulty's boss card abilities.
func haunt(d engine.CardDef, fearText string, each bool) engine.CardDef {
	d.Guarded = trickGuard
	d.Statics = []engine.Static{fearEvasion}
	d.Abilities = []engine.Ability{fearOn(fearText, each), deathRoll}
	return d
}

// --- Trick/Treat ---

// trickOrTreating: "As each turn starts, you may flip this" (R-TRICK-11,
// Q-43): a start-of-turn trigger; the active player says yes or no.
func trickOrTreating(ask engine.Key) engine.Ability {
	return engine.Ability{
		Kind:    engine.Triggered,
		Text:    "As each turn starts, you may flip this.",
		Trigger: engine.AtStartOfEachTurn(),
		Effects: []engine.Effect{cardkit.May(ask, engine.FlipThis())},
	}
}

// nobodysFool: "Each time you kill a monster, remove a trick counter
// from the final boss": every monster death counts (R-TRICK-12, Q-54).
var nobodysFool = engine.Ability{
	Kind:    engine.Triggered,
	Text:    "Each time you kill a monster, remove a trick counter from the final boss.",
	Trigger: engine.WhenAMonsterDies(),
	Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) {
		for _, b := range finalBosses(c.G) {
			c.G.AddNamedCounters(b, TrickCounter, -1)
		}
	})},
}

// trickRolls is Ghostly Trickery's roll count: 1, Hard 2, Ultra Hard 3
// (R-TRICK-13, Q-46).
func trickRolls(c *engine.Ctx) int {
	if c.G.Challenge == nil {
		return 1
	}
	if l, ok := hauntLevels[c.G.Challenge.Difficulty]; ok {
		return l.rolls
	}
	return 1
}

// ghostlyTrickery: "At the start of each attack against a monster, the
// active player rolls" once per difficulty's roll; each result lasts
// till end of turn (R-TRICK-13, Q-46 to Q-48, Q-53). c.EventObject is
// the attacked monster.
var ghostlyTrickery = engine.Ability{
	Kind:    engine.Triggered,
	Text:    "At the start of each attack against a monster, the active player rolls. If you are playing on hard, roll 2x instead. If you are playing on ultra hard, roll 3x instead- 1-2: Damage you would take is increased by 1 till end of turn. 3-4: Reveal cards from the top of the monster deck until you reveal a non-event monster card. The monster being attacked gains its abilities till end of turn. Put any revealed cards on the bottom of the deck. 5-6: The rewards of that monster become 1¢ till end of turn.",
	Trigger: engine.AtStartOfAttackOnMonster(),
	Effects: []engine.Effect{engine.RollTimes(trickRolls, engine.RollTable{}.
		Results(1, 2, engine.EffectFunc(func(c *engine.Ctx) { c.G.IncreaseDamageTaken(c.Controller, 1, c.Source) })).
		Results(3, 4, engine.EffectFunc(func(c *engine.Ctx) { c.G.GainRevealedMonster(c.EventObject, c.Controller) })).
		Results(5, 6, engine.EffectFunc(func(c *engine.Ctx) { c.G.SetRewardsToCents(c.EventObject, 1, c.Source) })))},
}

// youDontScareMe: "Each time you gain rewards, remove a fear counter
// from a player": once per dead card whose rewards the active player
// gains, only if it had any (R-TRICK-14, Q-54).
var youDontScareMe = engine.Ability{
	Kind: engine.Triggered,
	Text: "Each time you gain rewards, remove a fear counter from a player.",
	Trigger: engine.Trigger{On: engine.EvRewardsGained, Match: func(g *engine.Game, _ engine.ObjectID, e engine.Event) bool {
		return e.Player == g.Turn.Active && len(g.Rewards(e.Object)) > 0
	}},
	Effects: []engine.Effect{removeFear},
}

// addReward adds n of kind to the attacked monster's rewards.
func addReward(kind engine.RewardKind, n int) engine.Effect {
	return engine.EffectFunc(func(c *engine.Ctx) { c.G.AddReward(c.EventObject, kind, n, c.Source) })
}

// hallowsEveBounty: "At the start of each attack against a monster, the
// active player rolls" once at every difficulty; the result adds to the
// attacked monster's rewards till end of turn (R-TRICK-15, Q-48, Q-53).
var hallowsEveBounty = engine.Ability{
	Kind:    engine.Triggered,
	Text:    "At the start of each attack against a monster, the active player rolls- 1-2: Add 3¢ to the rewards of that monster till end of turn. 3-4: Add loot 1 to the rewards of that monster till end of turn. 5-6: Add +1 treasure to the rewards of that monster till end of turn.",
	Trigger: engine.AtStartOfAttackOnMonster(),
	Effects: []engine.Effect{engine.Roll(engine.RollTable{}.
		Results(1, 2, addReward(engine.RewardCents, 3)).
		Results(3, 4, addReward(engine.RewardLoot, 1)).
		Results(5, 6, addReward(engine.RewardTreasure, 1)))},
}

// trickGoal: the challenge is won when Mischievous The Haunt dies
// (R-TRICK-01). Both sides hold it, so a flip never drops it.
var trickGoal = engine.Ability{
	Kind: engine.Triggered,
	Text: "When Mischievous The Haunt dies, the players win.",
	Trigger: engine.Trigger{On: engine.EvDied, Match: func(_ *engine.Game, _ engine.ObjectID, e engine.Event) bool {
		return e.Player == engine.NoPlayer && slices.Contains(hauntBosses, e.Card)
	}},
	Effects: []engine.Effect{engine.EffectFunc(func(c *engine.Ctx) { c.G.WinChallenge() })},
}
