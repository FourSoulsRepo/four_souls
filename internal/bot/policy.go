package bot

import (
	"math/rand/v2"
	"slices"

	engine "github.com/FourSoulsRepo/rules_engine"
)

// The playout policy (Decision 6): weights on the allowed intents, then
// a weighted random pick. No copies, so it is cheap enough for every
// seat of every playout. Weights read kinds and public numbers only:
// coins, prices, HP, ATK, evasion, souls on the table. Every weight is
// here.
const (
	// shopPrice is what a shop item or the treasure deck's top costs
	// (R-SHOP-03); the policy buys only with that much.
	shopPrice = 10

	// In the active player's open action phase.
	wEndTurn  = 1.0  // ends the turn once nothing else draws
	wAttack   = 0.5  // any attack: souls come from monsters
	wAttackOd = 10.0 // plus this times the best odds of a kill
	wPurchase = 30.0 // with the coins for it: buy before a fight
	wLootOwn  = 3.0
	wUseOwn   = 0.7 // divided by 1 + abilities used this turn

	// Anywhere else: mostly let the stack resolve.
	wPass     = 12.0
	wLootElse = 0.4
	wUseElse  = 0.3

	// Attack targets: the odds counted twice, times 1 + wSoulFactor
	// for each soul; wTargetMin keeps every target possible.
	wTargetMin  = 0.05
	wSoulFactor = 4.0
	// The monster deck's top is judged as an average monster of every
	// set (TestDeckMonsterMean keeps these near the cards' means); its
	// souls are the world's share of soul cards there (soulShare).
	deckHP, deckDC, deckATK = 3, 4, 1

	// Other choices.
	wDeckBuy     = 0.6 // the treasure deck's top, against 1 for a shop item
	wCurseSelf   = 0.05
	wCurseLeader = 1.0 // plus one per soul of the seat that gets it
)

// policyPick picks one allowed intent of seat p by the policy weights
// and completes it (Decision 6).
func policyPick(g *engine.Game, p engine.PlayerID, allowed []engine.Intent, rng *rand.Rand) engine.Intent {
	ws := policyWeights(g, p, allowed)
	return complete(g, p, allowed[weighted(ws, rng)], rng)
}

// policyBest is the allowed intent of seat p with the top policy
// weight, ties drawn at random, completed.
func policyBest(g *engine.Game, p engine.PlayerID, allowed []engine.Intent, rng *rand.Rand) engine.Intent {
	ws := policyWeights(g, p, allowed)
	top := slices.Max(ws)
	for i, w := range ws {
		if w < top {
			ws[i] = 0
		}
	}
	return complete(g, p, allowed[weighted(ws, rng)], rng)
}

// weighted is an index drawn with chance in proportion to its weight;
// with no positive weight, any index.
func weighted(ws []float64, rng *rand.Rand) int {
	sum := 0.0
	for _, w := range ws {
		sum += max(w, 0)
	}
	if sum <= 0 {
		return rng.IntN(len(ws))
	}
	x := rng.Float64() * sum
	for i, w := range ws {
		x -= max(w, 0)
		if x < 0 {
			return i
		}
	}
	return len(ws) - 1
}

// policyWeights is the weight of each allowed intent of seat p.
func policyWeights(g *engine.Game, p engine.PlayerID, allowed []engine.Intent) []float64 {
	ws := make([]float64, len(allowed))
	w := g.PromptFor(p)
	if w.Kind == engine.PromptChoose {
		return choiceWeights(g, p, w, allowed)
	}
	open := false // the active player's own action phase, stack empty
	for _, in := range allowed {
		if in.Kind == engine.IntentEndTurn {
			open = true
		}
	}
	pl := g.Players[p]
	for i, in := range allowed {
		switch in.Kind { //nolint:exhaustive // other prompts have one intent
		case engine.IntentPass:
			ws[i] = wPass
		case engine.IntentEndTurn:
			ws[i] = wEndTurn
		case engine.IntentAttack:
			ws[i] = wAttack + wAttackOd*max(bestOdds(g, p), deckOdds(g, p))
		case engine.IntentPurchase:
			ws[i] = pick(pl.Cents >= shopPrice, wPurchase, 0)
		case engine.IntentPlayLoot:
			ws[i] = pick(open, wLootOwn, wLootElse)
		case engine.IntentActivate:
			ws[i] = pick(open, wUseOwn, wUseElse) / float64(1+pl.Activations)
		default:
			ws[i] = 1
		}
	}
	return ws
}

// choiceWeights weighs the options of a choose prompt by its purpose.
// An option that cancels the move gets no weight, so a play is never
// taken back and tried again for ever.
func choiceWeights(g *engine.Game, p engine.PlayerID, w engine.Prompt, allowed []engine.Intent) []float64 {
	ws := make([]float64, len(allowed))
	for i := range ws {
		ws[i] = 1
	}
	c := g.Choice
	if c == nil || c.Player != p || c.Purpose != w.Purpose {
		return ws // a start-of-game choice of several seats: any
	}
	n := len(w.Options)
	for i, in := range allowed {
		o := in.Choice
		switch c.Purpose { //nolint:exhaustive // the rest: any option
		case engine.ChooseAttackTarget:
			ws[i] = attackTargetWeight(g, p, c, o)
		case engine.ChoosePurchase:
			ws[i] = pick(o < len(c.Objects), 1, wDeckBuy)
		case engine.ChooseTarget:
			ws[i] = pick(o >= len(c.Targets) && n > len(c.Targets), 0, 1)
		case engine.ChooseMode:
			ws[i] = pick(o >= len(c.Indexes) && n > len(c.Indexes), 0, 1)
		case engine.ChooseCopyAbility:
			ws[i] = pick(o >= len(c.Abilities), 0, 1)
		case engine.ChooseCursed, engine.ChoosePassCurse:
			if o < len(c.Slots) {
				q := engine.PlayerID(c.Slots[o])
				ws[i] = pick(q == p, wCurseSelf, wCurseLeader+float64(g.SoulValue(q)))
			}
		}
	}
	return ws
}

// attackTargetWeight weighs option o of an attack target choice: the
// monsters, then the monster deck, then the owed random attack.
func attackTargetWeight(g *engine.Game, p engine.PlayerID, c *engine.Choice, o int) float64 {
	if o < len(c.Objects) {
		m := c.Objects[o]
		return targetWeight(attackOdds(g, p, m), float64(g.Def(m).Soul))
	}
	if c.Deck && o == len(c.Objects) {
		return targetWeight(deckOdds(g, p), soulShare(g))
	}
	return wTargetMin // the random attack
}

// targetWeight is the weight of an attack target with these odds and
// souls: its expected souls, and its rewards, the odds counted twice.
func targetWeight(odds, souls float64) float64 {
	return wTargetMin + odds*odds*(1+wSoulFactor*souls)
}

// deckOdds is p's chance against the monster deck's unknown top, judged
// as an average monster.
func deckOdds(g *engine.Game, p engine.PlayerID) float64 {
	return odds(g, p, deckHP, deckDC, deckATK)
}

// bestOdds is the best chance of p to kill a monster in play now.
func bestOdds(g *engine.Game, p engine.PlayerID) float64 {
	best := 0.0
	for _, m := range g.MonstersInPlay() {
		if g.CanBeAttacked(m) {
			best = max(best, attackOdds(g, p, m))
		}
	}
	return best // no monster: 0, and deckOdds decides
}

// attackOdds is the chance that p kills monster m before m kills p,
// from the public numbers: p's ATK and HP, m's HP, evasion and ATK. A
// roll of the evasion or more hits (R-ATK-18); each miss costs p the
// monster's ATK in HP.
func attackOdds(g *engine.Game, p engine.PlayerID, m engine.ObjectID) float64 {
	return odds(g, p, g.HP(m), g.Evasion(m), g.MonsterATK(m))
}

// odds is the chance that p kills a monster of HP hp, evasion dc and
// ATK matk before it kills p.
func odds(g *engine.Game, p engine.PlayerID, hp, dc, matk int) float64 {
	atk := g.PlayerATK(p)
	if atk <= 0 {
		return 0
	}
	hits := (max(hp, 1) + atk - 1) / atk
	// A sure hit is never sure: an ability may raise the evasion in
	// the fight, so a roll of 1 counts as a miss.
	hit := float64(7-min(max(dc, 2), 6)) / 6
	misses := 1 << 10 // a monster with no ATK never kills
	if matk > 0 {
		misses = (max(g.PlayerHP(p), 1) + matk - 1) / matk
	}
	return winBefore(hits, misses, hit)
}

// winBefore is the chance of hits successes before misses failures,
// each try a success with chance h.
func winBefore(hits, misses int, h float64) float64 {
	if h >= 1 {
		return 1
	}
	misses = min(misses, 64) // past this the sum no longer moves
	// sum over j failures (j < misses) of C(hits-1+j, j) h^hits (1-h)^j
	term := 1.0
	for range hits {
		term *= h
	}
	sum := term
	for j := 1; j < misses; j++ {
		term *= float64(hits-1+j) / float64(j) * (1 - h)
		sum += term
	}
	return min(sum, 1)
}

// pick is a when cond holds, else b.
func pick(cond bool, a, b float64) float64 {
	if cond {
		return a
	}
	return b
}
