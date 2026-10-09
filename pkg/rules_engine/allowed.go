package rulesengine

// Allowed lists every intent player p may send right now (N-04). Each one
// is accepted by Apply; anything else is refused with a rule ID.
//
// For a discard prompt it returns one template intent whose Objects are
// the cards to pick from; send Prompt().Count of them.
func (g *Game) Allowed(p PlayerID) []Intent {
	if g.Over || g.Waiting.Player != p {
		return nil
	}
	var out []Intent
	try := func(in Intent) {
		if g.check(in) == nil {
			out = append(out, in)
		}
	}
	switch g.Waiting.Kind {
	case PromptPriority:
		try(Intent{Player: p, Kind: IntentPass})
		try(Intent{Player: p, Kind: IntentEndTurn})
		try(Intent{Player: p, Kind: IntentAttack})
		try(Intent{Player: p, Kind: IntentPurchase})
		for _, id := range g.Players[p].Hand {
			try(Intent{Player: p, Kind: IntentPlayLoot, Objects: []ObjectID{id}})
		}
		for _, id := range g.controlled(p) {
			for i := range g.AbilitiesOf(id) {
				try(Intent{Player: p, Kind: IntentActivate, Objects: []ObjectID{id}, Choice: i})
			}
		}
	case PromptChoose:
		for i := range g.Waiting.Options {
			try(Intent{Player: p, Kind: IntentChoose, Choice: i})
		}
	case PromptDiscard:
		out = append(out, Intent{Player: p, Kind: IntentDiscard, Objects: append([]ObjectID(nil), g.Players[p].Hand...)})
	case PromptNone, PromptGameOver:
	}
	return out
}

// CanAct reports whether p has anything to do besides passing. The server
// skips a player who cannot act (N-05, N-09).
func (g *Game) CanAct(p PlayerID) bool {
	for _, in := range g.Allowed(p) {
		if in.Kind != IntentPass {
			return true
		}
	}
	return false
}

// controlled lists the character and in-play objects of p.
func (g *Game) controlled(p PlayerID) []ObjectID {
	return append([]ObjectID{g.Players[p].Character}, g.Players[p].InPlay...)
}
