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
}

// ask opens a choose prompt.
func (g *Game) ask(c Choice, labels []string) {
	g.Choice = &c
	g.Waiting = Prompt{Kind: PromptChoose, Player: c.Player, Options: labels, Purpose: c.Purpose}
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
		g.destroyItem(c.Player, c.Objects[i])
	case ChoosePenaltyLoot:
		g.discardFromHand(c.Player, c.Objects[i])
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
