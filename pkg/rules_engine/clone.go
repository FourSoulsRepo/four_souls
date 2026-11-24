package rulesengine

import (
	"maps"
	"slices"
)

// Clone returns a deep copy of the game, written out by hand: a bot
// copies the game for every move it tries (ADR 033, Decision 3).
//
// The copy saves to the same bytes and has the same checksum; it shares
// no slice, map or pointer with the game, so a change to one never shows
// in the other. Only the card definitions are shared: they never change.
// Like a save and a load, it leaves out the events of the last step and
// the test dice (forcedRolls).
//
// A new field of Game, or of a type inside it, needs a line here and in
// CloneRules (enginetest/clone.go); TestCloneCoversEveryField fails until both
// are there.
func (g *Game) Clone() (*Game, error) {
	return g.clone(true), nil
}

// clone is Clone; with keepSaved false the copy keeps no save tracking
// (saved): Imagine's worlds never save, so copying the steps since the
// last clean state, which grow all game, would be thrown away (19.213).
func (g *Game) clone(keepSaved bool) *Game {
	c := *g
	c.Players = clonePlayers(g.Players)
	c.Objects = cloneObjects(g.Objects)
	for i := range g.Decks {
		c.Decks[i] = slices.Clone(g.Decks[i])
		c.Discards[i] = slices.Clone(g.Discards[i])
	}
	c.Shop = cloneSlots(g.Shop)
	c.Monsters = cloneSlots(g.Monsters)
	c.Rooms = cloneSlots(g.Rooms)
	c.Bosses = cloneSlots(g.Bosses)
	c.Minions = cloneSlots(g.Minions)
	c.Outside = slices.Clone(g.Outside)
	c.BonusSouls = slices.Clone(g.BonusSouls)
	c.Queue = cloneActions(g.Queue)
	c.Attack = g.Attack.clone()
	c.Choice = clonePtr(g.Choice, (*Choice).clone)
	c.StartChoices = cloneEach(g.StartChoices, (*Choice).clone)
	c.Throw = clonePtr(g.Throw, func(t *Throw) Throw { return *t })
	c.FirstRoll = clonePtr(g.FirstRoll, (*FirstRoll).clone)
	c.Activating = clonePtr(g.Activating, (*Activation).clone)
	c.PendingTriggers = cloneEach(g.PendingTriggers, func(t *PendingTrigger) PendingTrigger {
		t2 := *t
		t2.Targets = slices.Clone(t.Targets)
		return t2
	})
	c.Boosts = slices.Clone(g.Boosts)
	c.Shields = slices.Clone(g.Shields)
	c.Lingering = slices.Clone(g.Lingering)
	c.Markers = slices.Clone(g.Markers)
	c.RewardChanges = slices.Clone(g.RewardChanges)
	c.ManyDice = cloneEach(g.ManyDice, func(m *ManyDice) ManyDice {
		m2 := *m
		m2.Values = slices.Clone(m.Values)
		return m2
	})
	c.Stack = cloneEach(g.Stack, (*StackItem).clone)
	c.Turn = g.Turn.clone()
	c.Waiting = g.Waiting.clone()
	c.Winners = slices.Clone(g.Winners)
	c.Challenge = clonePtr(g.Challenge, (*ChallengeState).clone)
	c.Sets = slices.Clone(g.Sets)
	c.Banned = slices.Clone(g.Banned)
	c.asks = maps.Clone(g.asks) // an askEffect holds functions only
	if keepSaved {
		c.saved.since = cloneEach(g.saved.since, (*ReplayStep).clone) // origin is never changed in place
	} else {
		c.saved = tracking{}
	}
	c.maxHP = maps.Clone(g.maxHP)
	c.events = nil
	c.forcedRolls = nil
	c.changes = nil // the event check's log (TakeSnapshot) stays with the game
	return &c
}

// cloneEach copies a slice, each element by one; nil stays nil and
// empty stays empty, so the copy saves to the same bytes.
func cloneEach[T any](s []T, one func(*T) T) []T {
	if s == nil {
		return nil
	}
	out := make([]T, len(s))
	for i := range s {
		out[i] = one(&s[i])
	}
	return out
}

// clonePtr copies the value p points to; nil stays nil.
func clonePtr[T any](p *T, one func(*T) T) *T {
	if p == nil {
		return nil
	}
	v := one(p)
	return &v
}

func clonePlayers(ps []Player) []Player {
	return cloneEach(ps, func(p *Player) Player {
		c := *p
		c.Hand = slices.Clone(p.Hand)
		c.InPlay = slices.Clone(p.InPlay)
		c.Counters = slices.Clone(p.Counters)
		c.SkipSources = slices.Clone(p.SkipSources)
		return c
	})
}

// cloneObjects copies the objects in one go; only the few with lists
// get their own.
func cloneObjects(objs []Object) []Object {
	out := slices.Clone(objs)
	for i := range out {
		o := &out[i]
		if o.Counters != nil || o.MoreGained != nil || o.Under != nil || o.ShareRewards != nil {
			o.Counters = slices.Clone(o.Counters)
			o.MoreGained = slices.Clone(o.MoreGained)
			o.Under = slices.Clone(o.Under)
			o.ShareRewards = slices.Clone(o.ShareRewards)
		}
	}
	return out
}

func cloneSlots(ss []Slot) []Slot {
	return cloneEach(ss, func(s *Slot) Slot { return Slot{Cards: slices.Clone(s.Cards)} })
}

func cloneActions(as []Action) []Action { return cloneEach(as, (*Action).clone) }

func (a *Action) clone() Action {
	c := *a
	c.Ask = clonePtr(a.Ask, (*Asking).clone)
	c.Applied = slices.Clone(a.Applied)
	return c
}

func (a *Asking) clone() Asking {
	c := *a
	c.Targets = slices.Clone(a.Targets)
	c.Answers = slices.Clone(a.Answers)
	c.Order = slices.Clone(a.Order)
	return c
}

func (a AttackState) clone() AttackState {
	a.Helpers = slices.Clone(a.Helpers)
	a.Queue = slices.Clone(a.Queue)
	a.Share = slices.Clone(a.Share)
	return a
}

func (ch *Choice) clone() Choice {
	c := *ch
	c.Replacements = slices.Clone(ch.Replacements)
	c.Objects = slices.Clone(ch.Objects)
	c.Previews = cloneEach(ch.Previews, func(p *ReplacementPreview) ReplacementPreview {
		return ReplacementPreview{Actions: cloneActions(p.Actions), Unused: slices.Clone(p.Unused), Again: p.Again}
	})
	c.Slots = slices.Clone(ch.Slots)
	c.Targets = slices.Clone(ch.Targets)
	c.Indexes = slices.Clone(ch.Indexes)
	c.Abilities = slices.Clone(ch.Abilities)
	c.Ask = clonePtr(ch.Ask, (*Asking).clone)
	c.Question = ch.Question.clone()
	c.Look = clonePtr(ch.Look, func(l *DeckLook) DeckLook { return *l })
	return c
}

// clone copies the message. A value is an int, a string, a CardValue or
// a PlayerValue (argValue): plain values, so the copy may share them.
func (m Msg) clone() Msg {
	m.Args = slices.Clone(m.Args)
	return m
}

func (f *FirstRoll) clone() FirstRoll {
	c := *f
	c.Players = slices.Clone(f.Players)
	c.Rolls = slices.Clone(f.Rolls)
	return c
}

func (a *Activation) clone() Activation {
	c := *a
	c.Chosen = slices.Clone(a.Chosen)
	return c
}

func (s *StackItem) clone() StackItem {
	c := *s
	c.Label = s.Label.clone()
	c.Targets = slices.Clone(s.Targets)
	c.Sources = cloneEach(s.Sources, func(v *StatSourceView) StatSourceView {
		c := *v
		c.When = v.When.clone()
		return c
	})
	return c
}

func (t Turn) clone() Turn {
	t.NextHelpers = slices.Clone(t.NextHelpers)
	t.MustAttackOneOf = slices.Clone(t.MustAttackOneOf)
	t.NextShare = slices.Clone(t.NextShare)
	t.LateRooms = slices.Clone(t.LateRooms)
	return t
}

func (p Prompt) clone() Prompt {
	p.Options = cloneEach(p.Options, (*Msg).cloneP)
	p.Text = p.Text.clone()
	p.Look = clonePtr(p.Look, func(l *DeckLook) DeckLook { return *l })
	p.Players = slices.Clone(p.Players)
	return p
}

func (m *Msg) cloneP() Msg { return m.clone() }

func (cs *ChallengeState) clone() ChallengeState {
	c := *cs
	c.Presents = slices.Clone(cs.Presents)
	c.Reason = cs.Reason.clone()
	c.TurnFrom = slices.Clone(cs.TurnFrom)
	c.RoundFrom = slices.Clone(cs.RoundFrom)
	return c
}

// clone copies a replay step with its own object list.
func (r *ReplayStep) clone() ReplayStep {
	c := *r
	c.Intent.Objects = slices.Clone(r.Intent.Objects)
	return c
}
