package r

import (
	"strconv"

	engine "github.com/FourSoulsRepo/rules_engine"
)

// Auction Gavel (Active Treasure Card)
//
//	{Tap Effect}Choose a shop item. Other players may bid ¢ for it, starting at 1 or more ¢. In turn order, other players may top the high bid. Bidding ends when the high bid stands. The high bidder pays you the number of ¢ they bid and steals the chosen item.
var auctionGavel = engine.CardDef{
	Ref:       "auction_gavel",
	Kind:      engine.TreasureCard,
	Copies:    1,
	Tap:       true,
	ThreePlus: true,
	Abilities: []engine.Ability{{
		Kind:    engine.Activated,
		Text:    "↷: Choose a shop item. Other players may bid ¢ for it, starting at 1 or more ¢. In turn order, other players may top the high bid. Bidding ends when the high bid stands. The high bidder pays you the number of ¢ they bid and steals the chosen item.",
		Costs:   []engine.Cost{engine.Tap()},
		Targets: []engine.TargetSpec{engine.ChooseWhere(engine.TargetItem, inShop)},
		Effects: []engine.Effect{auction()},
	}},
}

// maxBids bounds the bidding questions; each bid raises by at least 1¢.
const maxBids = 60

// auctionState replays the answers: the bidder whose turn it is, or
// NoPlayer when bidding is over, and the high bid.
func auctionState(c *engine.Ctx, a []int) (next, high engine.PlayerID, bid int) {
	others := othersInOrder(c)
	high = engine.NoPlayer
	passes, turn := 0, 0
	for _, x := range a {
		p := others[turn%len(others)]
		if x > 0 { // option 0 passes; option k bids bid+k
			bid, high, passes = bid+x, p, 0
		} else {
			passes++
		}
		turn++
	}
	// Bidding ends when the high bid stands: every other bidder passed in a
	// row, or everyone passed without a bid.
	if len(others) == 0 || (high == engine.NoPlayer && passes >= len(others)) || (high != engine.NoPlayer && passes >= len(others)-1) {
		return engine.NoPlayer, high, bid
	}
	p := others[turn%len(others)]
	if p == high {
		turn++
		p = others[turn%len(others)]
	}
	return p, high, bid
}

func auction() engine.Effect {
	qs := make([]engine.Question, maxBids)
	for i := range qs {
		qs[i] = engine.Question{
			Text: "Bid for the shop item?",
			Player: func(c *engine.Ctx, a []int) engine.PlayerID {
				next, _, _ := auctionState(c, a)
				return next
			},
			Options: func(c *engine.Ctx, a []int) []string {
				next, _, bid := auctionState(c, a)
				if next == engine.NoPlayer {
					return nil
				}
				out := []string{"pass"}
				for n := bid + 1; n <= c.G.Players[next].Cents; n++ {
					out = append(out, "bid "+strconv.Itoa(n)+"¢")
				}
				return out
			},
		}
	}
	return engine.Ask(func(c *engine.Ctx, a []int) {
		var answers []int
		for _, x := range a {
			if x >= 0 {
				answers = append(answers, x)
			}
		}
		_, high, bid := auctionState(c, answers)
		item := c.Targets[0].Object
		if high == engine.NoPlayer || !inShop(c.G, c.Source, c.Targets[0]) || c.G.Players[high].Cents < bid {
			return
		}
		c.G.Players[high].Cents -= bid
		c.G.Players[c.Controller].Cents += bid
		c.G.GainControl(high, item)
	}, qs...)
}
