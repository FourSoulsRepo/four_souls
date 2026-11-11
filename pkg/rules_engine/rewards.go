package rulesengine

// Changes to a monster's rewards till end of turn: "Add loot 1 to the
// rewards of that monster" (Hallows' Eve Bounty, R-TRICK-15) and "The
// rewards of that monster become 1¢" (Ghostly Trickery, R-TRICK-13).
// They apply in the order made (Q-48): becoming ¢ replaces the whole
// reward with every earlier addition; later additions add to it. The
// soul is no reward: a soul icon stays (Q-48.4). A doubling (Dinga)
// doubles the changed rewards.

// RewardChange is one change to the rewards of Object till end of turn.
// Set: the rewards become Amount¢; otherwise Amount of Kind is added.
type RewardChange struct {
	Object ObjectID   `json:"object,omitempty"`
	Set    bool       `json:"set,omitempty"`
	Kind   RewardKind `json:"kind,omitempty"`
	Amount int        `json:"amount,omitempty"`
}

// EvRewardsChanged means Object's rewards changed till end of turn by
// Source. Text is "add cents", "add loot" or "add treasure" (Amount of
// it added), or "set cents" (they became Amount¢).
const EvRewardsChanged EventKind = "rewards_changed"

// rewardCodes name the reward kinds in events.
var rewardCodes = [...]string{RewardCents: "cents", RewardLoot: "loot", RewardTreasure: "treasure"}

// AddReward adds n of kind to the rewards of monster id till end of
// turn (R-TRICK-15); a monster without rewards gets it as its only one
// (Q-48.5). source did it.
func (g *Game) AddReward(id ObjectID, kind RewardKind, n int, source ObjectID) {
	g.changeRewards(RewardChange{Object: id, Kind: kind, Amount: n}, "add "+rewardCodes[kind], source)
}

// SetRewardsToCents makes the rewards of monster id n¢ till end of turn,
// whatever was added before (R-TRICK-13, Q-48). source did it.
func (g *Game) SetRewardsToCents(id ObjectID, n int, source ObjectID) {
	g.changeRewards(RewardChange{Object: id, Set: true, Kind: RewardCents, Amount: n}, "set cents", source)
}

func (g *Game) changeRewards(c RewardChange, text string, source ObjectID) {
	if g.Object(c.Object).Zone.Kind != ZoneInPlay {
		return
	}
	g.RewardChanges = append(g.RewardChanges, c)
	g.emit(Event{Kind: EvRewardsChanged, Player: NoPlayer, Object: c.Object, Card: g.Object(c.Object).Card, Amount: c.Amount, Text: text, Source: source})
}

// Rewards are the rewards object id gives now: its card's (R-CARD-31),
// with this turn's changes in the order made (Q-48).
func (g *Game) Rewards(id ObjectID) []Reward {
	out := g.def(id).Rewards
	changed := false
	for _, c := range g.RewardChanges {
		if c.Object != id {
			continue
		}
		if !changed {
			out, changed = append([]Reward(nil), out...), true
		}
		if c.Set {
			out = []Reward{{Kind: RewardCents, Amount: c.Amount}}
			continue
		}
		out = append(out, Reward{Kind: c.Kind, Amount: c.Amount})
	}
	return out
}

// RewardsChanged reports whether id's rewards differ from its card's
// this turn.
func (g *Game) RewardsChanged(id ObjectID) bool {
	for _, c := range g.RewardChanges {
		if c.Object == id {
			return true
		}
	}
	return false
}

// moveRewardChanges carries a dying monster's changes to its dead card,
// which gives the rewards (R-DEATH-04, R-DEATH-06).
func (g *Game) moveRewardChanges(from, to ObjectID) {
	for i := range g.RewardChanges {
		if g.RewardChanges[i].Object == from {
			g.RewardChanges[i].Object = to
		}
	}
}

// rewardViews are the rewards id gives now, for the view (14.33).
func (g *Game) rewardViews(id ObjectID) []RewardView {
	rs := g.Rewards(id)
	out := make([]RewardView, 0, len(rs))
	for _, r := range rs {
		out = append(out, RewardView{Kind: rewardCodes[r.Kind], Amount: r.Amount, Roll: r.Roll})
	}
	return out
}
