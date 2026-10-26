package rulesengine

import (
	"errors"
	"slices"
)

// Block says why player p can't do something now, in a few words
// (15.29, F-17, F-18, C-06). Kind is the intent it would be: a loot
// play or an ability use names its card in Object; an attack, a
// purchase or the end of the turn has no object.
type Block struct {
	Kind   IntentKind `json:"kind"`
	Object ObjectID   `json:"object"`
	Reason Msg        `json:"reason"`
}

// Blocked lists what p can't do now and why: each loot card in their
// hand, each card of theirs with an activated ability, each room with
// one, and the turn actions. Anything Allowed is left out. A card that
// would be usable with priority says to wait for it.
func (g *Game) Blocked(p PlayerID) []Block {
	if g.Over || p < 0 || int(p) >= len(g.Players) {
		return nil
	}
	allowed := g.Allowed(p)
	ok := func(kind IntentKind, id ObjectID) bool {
		return slices.ContainsFunc(allowed, func(in Intent) bool {
			return in.Kind == kind && (len(in.Objects) == 0 || in.Objects[0] == id)
		})
	}
	wait := g.waitSay(p)
	// Each intent is checked as if p had priority: a refusal then is the
	// real reason; none means p only waits for it.
	was := g.Waiting
	g.Waiting = Prompt{Kind: PromptPriority, Player: p}
	defer func() { g.Waiting = was }()
	why := func(in Intent) Msg {
		var re *RuleError
		switch err := g.check(in); {
		case err == nil:
			return wait
		case errors.As(err, &re) && !re.Say.IsZero():
			return re.Say
		default:
			return refuseNotUsable.Msg()
		}
	}
	var out []Block
	for _, kind := range []IntentKind{IntentAttack, IntentPurchase, IntentEndTurn} {
		if ok(kind, 0) {
			continue
		}
		reason := g.turnSay(p)
		if reason.IsZero() {
			reason = why(Intent{Player: p, Kind: kind})
		}
		out = append(out, Block{Kind: kind, Reason: reason})
	}
	for _, id := range g.Players[p].Hand {
		if !ok(IntentPlayLoot, id) {
			out = append(out, Block{Kind: IntentPlayLoot, Object: id, Reason: why(Intent{Player: p, Kind: IntentPlayLoot, Objects: []ObjectID{id}})})
		}
	}
	for _, id := range append(g.controlled(p), g.roomsInPlay()...) {
		first := slices.IndexFunc(g.AbilitiesOf(id), func(r AbilityRef) bool { return g.ability(r).Kind == Activated })
		if first < 0 || ok(IntentActivate, id) {
			continue
		}
		out = append(out, Block{Kind: IntentActivate, Object: id, Reason: why(Intent{Player: p, Kind: IntentActivate, Objects: []ObjectID{id}, Choice: first})})
	}
	return out
}

// waitSay is why p can't act now although the rules would let them.
func (g *Game) waitSay(p PlayerID) Msg {
	if g.Waiting.Player == p && g.Waiting.Kind != PromptPriority {
		return refuseAnswerPrompt.Msg()
	}
	return refuseWaitPriority.Msg()
}

// turnSay is why p can't attack, buy or end the turn now, when the
// turn itself is the reason; no line when it is not.
func (g *Game) turnSay(p PlayerID) Msg {
	switch {
	case p != g.Turn.Active:
		return refuseNotYourTurn.Msg()
	case g.Turn.Step != StepAction:
		return refuseActionPhase.Msg()
	case g.Turn.EndDeclared || g.Turn.DeathEnd:
		return refuseTurnEnding.Msg()
	case g.Attack.On:
		return refuseFinishAttack.Msg()
	case g.Purchase.On:
		return refuseFinishPurchase.Msg()
	case len(g.Stack) > 0:
		return refuseWaitStack.Msg()
	default:
		return Msg{}
	}
}

// noLootSay is why p has no loot play: none is left on their turn;
// on another's turn they have none.
func (g *Game) noLootSay(p PlayerID) Msg {
	if p != g.Turn.Active {
		return refuseLootYourTurn.Msg()
	}
	return refuseNoLootPlays.Msg()
}

// costSay is why a cost can't be paid, for the player.
func costSay(c Cost) Msg {
	switch c := c.(type) {
	case tapCost:
		return refuseCostUsed.Msg()
	case centsCost:
		return refuseCostCents.Msg()
	case hpCost:
		return refuseCostHP.Msg()
	case destroySelfCost:
		return refuseCostEternal.Msg()
	case removeCountersCost:
		return refuseCostCounters.Msg()
	case discardChosenCost:
		return refuseCostDiscard.Msg()
	case condCost:
		return c.say.Msg() // "Not during an attack"
	default:
		return refuseNotUsable.Msg()
	}
}
