package rulesengine

// ChoicePurpose says what a choose prompt is for.
type ChoicePurpose int

// The choices the engine asks for. It never picks for a player, even
// when only one option exists (ADR 005).
const (
	ChooseReplacement  ChoicePurpose = iota // which replacement first (R-ABIL-33)
	ChooseAttackTarget                      // a monster or the monster deck (R-ATK-02)
	ChooseMonsterSlot                       // where a revealed monster goes (R-ATK-08)
	ChoosePenaltyItem                       // the death penalty item (R-DEATH-14)
	ChoosePenaltyLoot                       // the death penalty discard (R-DEATH-14)
	ChoosePurchase                          // a shop item or the treasure deck (R-SHOP-02)
	ChooseTarget                            // a target for an ability (R-ABIL-04)
	ChooseTriggerOrder                      // which trigger goes on the stack next (R-ABIL-15)
	ChooseMode                              // a "choose one-" option (R-ABIL-04)
	ChooseAnswer                            // a question while an ability resolves (R-ABIL-05)
	ChooseStartingItem                      // a start-of-game choice, e.g. Eden (R-SETUP-09)
	ChooseCursed                            // who gains a curse (R-ABIL-20)
	ChooseLowest                            // which card goes to the very bottom of a deck
)

// Choice is an open question to one player. Options are listed in the
// prompt; the answer is an index into them.
type Choice struct {
	Purpose      ChoicePurpose    `json:"purpose"`
	Player       PlayerID         `json:"player"`
	Rule         string           `json:"rule"`
	Replacements []ReplacementRef `json:"replacements,omitempty"`
	Objects      []ObjectID       `json:"objects,omitempty"`
	// Deck adds the monster deck as the last option (attack targets).
	Deck bool `json:"deck,omitempty"`
	// Slots are monster slot indexes (where a revealed card goes).
	Slots []int `json:"slots,omitempty"`
	// Targets are the target options of ChooseTarget; one more option
	// after them means "cancel".
	Targets []Chosen `json:"targets,omitempty"`
	// Indexes are pending-trigger indexes for ChooseTriggerOrder, or
	// mode indexes for ChooseMode (one more option means "cancel").
	Indexes []int `json:"indexes,omitempty"`
	// Ask is the effect waiting for this answer (ChooseAnswer).
	Ask *Asking `json:"ask,omitempty"`
	// Question is the text of the question, for display.
	Question string `json:"question,omitempty"`
	// Owner is whose death penalty a penalty choice is for; To, if not
	// NoPlayer, gains the discarded loot card instead (Shadow).
	Owner PlayerID `json:"owner"`
	To    PlayerID `json:"to"`
}

// ask opens a choose prompt.
func (g *Game) ask(c Choice, labels []string) {
	g.Choice = &c
	g.Waiting = Prompt{Kind: PromptChoose, Player: c.Player, Options: labels, Purpose: c.Purpose, Text: c.Question}
}

// answer carries out the chosen option of the open choice.
func (g *Game) answer(i int) {
	c := *g.Choice
	g.Choice, g.Waiting = nil, Prompt{}
	switch c.Purpose {
	case ChooseReplacement:
		g.applyReplacement(c.Replacements[i])
	case ChooseAttackTarget:
		if c.Deck && i == len(c.Objects) {
			g.attackDeck()
			return
		}
		g.startAttack(c.Objects[i])
	case ChooseMonsterSlot:
		g.placeRevealed(c.Slots[i])
	case ChoosePenaltyItem:
		g.destroyItem(c.Owner, c.Objects[i])
	case ChoosePenaltyLoot:
		if c.To != NoPlayer {
			g.GiveHandCard(c.Owner, c.To, c.Objects[i])
		} else {
			g.discardFromHand(c.Owner, c.Objects[i])
		}
	case ChooseTarget:
		g.chooseTarget(c, i)
	case ChooseTriggerOrder:
		g.pushTrigger(c.Indexes[i])
	case ChooseMode:
		g.chooseMode(c, i)
	case ChooseAnswer:
		g.answerAsk(c, i)
	case ChooseStartingItem:
		g.chooseStartingItem(c, i)
	case ChooseLowest:
		g.chooseLowest(c, i)
	case ChooseCursed:
		g.giveCurse(c.Objects[0], PlayerID(c.Slots[i]))
	case ChoosePurchase:
		if c.Deck && i == len(c.Objects) {
			g.purchase(0, true)
		} else {
			g.purchase(c.Objects[i], false)
		}
		g.givePriority(g.Turn.Active)
	}
}

// labels lists the cards of objects, for prompt options.
func (g *Game) labels(ids []ObjectID) []string {
	out := make([]string, len(ids))
	for i, id := range ids {
		out[i] = string(g.Object(id).Card)
	}
	return out
}
