package rulesengine

import (
	"errors"
	"fmt"
	"slices"
)

// The game types of setup.game (24.34).
const (
	GameRegular   = "regular"
	GameChallenge = "challenge"
	// gameTeam is a 2 vs 2 game: not in the engine (Q-03, parked).
	gameTeam = "team"
)

// gameTypes are the types setup.game takes, for errors.
var gameTypes = []string{GameRegular, GameChallenge}

// SituationGame is the game type of a situation (24.34). Type
// "regular", or no game at all, is a regular game. Type "challenge"
// starts the challenge ID at Difficulty on the built table, as a real
// game does (R-CHAL-01): its final boss, rule card, presents and timer.
// The other fields fit a challenge only.
type SituationGame struct {
	Type       string `json:"type"`
	ID         string `json:"id,omitempty"`
	Difficulty string `json:"difficulty,omitempty"`
	// Presents change a challenge's presents from how it starts them
	// (R-KRAMPUS-05): counters left, destroyed, or the player who
	// gained it at Happy Holidays!.
	Presents []SituationPresent `json:"presents,omitempty"`
	// Holidays: Happy Holidays! has happened (R-KRAMPUS-10); every
	// present not destroyed names the player who gained it.
	Holidays bool `json:"holidays,omitempty"`
	// RuleCard names the rule card's side up, by its ref: a two-sided
	// one (Trick/Treat) starts on its front (R-TRICK-04).
	RuleCard CardRef `json:"rule_card,omitempty"`
}

// SituationPresent is one present of a challenge in setup.game, by its
// kind (cents, loot, treasure, soul). Counters are the counters left
// (1 or more; 0 destroys it, R-KRAMPUS-08, so say destroyed); none
// keeps the start count. To is the seat that gained it as an item at
// Happy Holidays! (R-KRAMPUS-10).
type SituationPresent struct {
	Kind      string `json:"kind"`
	Counters  int    `json:"counters,omitempty"`
	Destroyed bool   `json:"destroyed,omitempty"`
	To        *int   `json:"to,omitempty"`
}

// ExpectedPresent checks one present of a challenge (24.34), by kind:
// the counters on it (0 once it is gone), destroyed, and To, the seat
// that controls it as an item. A field left out is not checked.
type ExpectedPresent struct {
	Kind      string `json:"kind"`
	Counters  *int   `json:"counters,omitempty"`
	Destroyed *bool  `json:"destroyed,omitempty"`
	To        *int   `json:"to,omitempty"`
}

// gameType checks setup.game's type and its fields, before the table
// is built: a regular game is nil.
func (s SituationSetup) gameType() (*SituationGame, error) {
	gm := s.Game
	if gm == nil {
		return nil, nil //nolint:nilnil // no game field: a regular game, no error
	}
	switch gm.Type {
	case GameRegular:
		for _, f := range []struct {
			name string
			set  bool
		}{
			{"id", gm.ID != ""},
			{"difficulty", gm.Difficulty != ""},
			{"presents", gm.Presents != nil},
			{"holidays", gm.Holidays},
			{"rule_card", gm.RuleCard != ""},
		} {
			if f.set {
				return nil, fmt.Errorf("situation: game: %s fits only type %s, not %s", f.name, GameChallenge, GameRegular)
			}
		}
		return nil, nil //nolint:nilnil // a regular game, as with no game field
	case GameChallenge:
		return gm, nil
	case gameTeam:
		return nil, fmt.Errorf("situation: game: type %s is not in the game yet (Q-03, parked); options %v", gameTeam, gameTypes)
	case "":
		return nil, fmt.Errorf("situation: game: type is missing; options %v", gameTypes)
	}
	return nil, fmt.Errorf("situation: game: type %q is no game type; options %v", gm.Type, gameTypes)
}

// findChallenge finds the challenge id among the sets' challenges.
func findChallenge(sets []CardSet, id string) (*ChallengeDef, error) {
	var ids []string
	for _, set := range sets {
		for _, c := range set.Challenges {
			if c.ID == id {
				return c, nil
			}
			ids = append(ids, c.ID)
		}
	}
	slices.Sort(ids)
	if id == "" {
		return nil, fmt.Errorf("situation: game: a challenge needs id; options %v", ids)
	}
	return nil, fmt.Errorf("situation: game: challenge %q is in none of the sets; options %v", id, ids)
}

// setUpGame starts setup.game's challenge on the table built so far,
// as a game's setup does (R-COOP-02, R-CHAL-01): the timer, the rule
// card, then the challenge's own setup (final boss, presents). The
// treasure deck top, with setup.decks, gives the presents. The round
// starts once the table is done (startRound).
func (g *Game) setUpGame(s SituationSetup, sets []CardSet) error {
	gm, err := s.gameType()
	if err != nil || gm == nil {
		return err
	}
	c, err := findChallenge(sets, gm.ID)
	if err != nil {
		return err
	}
	levels := make([]string, len(c.Levels))
	for i, l := range c.Levels {
		levels[i] = l.Difficulty
	}
	switch {
	case gm.Difficulty == "":
		return fmt.Errorf("situation: game: challenge %s needs difficulty; options %v", c.ID, levels)
	case !slices.Contains(levels, gm.Difficulty):
		return fmt.Errorf("situation: game: challenge %s has no difficulty %q; options %v", c.ID, gm.Difficulty, levels)
	case len(s.Players) != challengePlayers:
		return fmt.Errorf("situation: game: challenge %s is for %d players, not %d (R-COOP-01)", c.ID, challengePlayers, len(s.Players))
	case len(s.Bosses) > 0:
		return fmt.Errorf("situation: game: challenge %s puts its own final boss out (R-CHAL-02); leave bosses out", c.ID)
	case len(s.BonusSouls) > 0:
		return fmt.Errorf("situation: game: challenge %s is played without bonus souls (R-COOP-02); leave bonus_souls out", c.ID)
	}
	if err := g.startChallenge(c, gm.Difficulty); err != nil {
		return fmt.Errorf("situation: game: %w", err)
	}
	if err := g.setUpRuleCard(gm.RuleCard); err != nil {
		return err
	}
	return g.setUpPresents(gm.Presents, gm.Holidays)
}

// startRound starts a challenge's first round with the active player's
// turn, once the table is built, as StartChallenge does (R-TURN-15,
// Q-34): objects made from here are this turn's (R-COOP-11).
func (g *Game) startRound() {
	if g.Challenge == nil {
		return
	}
	g.Challenge.Starter = g.Turn.Active
	g.turnTaken(g.Turn.Active)
}

// setUpRuleCard turns the rule card to the side ref names: its front,
// or the back of a two-sided card (R-TRICK-04, R-TRICK-11).
func (g *Game) setUpRuleCard(ref CardRef) error {
	if ref == "" {
		return nil
	}
	id := g.Challenge.Card
	sides := []CardRef{g.Object(id).Card}
	if d, ok := g.cards.find(sides[0]); ok && d.Back != "" {
		sides = append(sides, d.Back)
	}
	switch ref {
	case sides[0]:
		return nil
	case sides[len(sides)-1]:
		g.Object(id).Flipped = true
		return nil
	}
	return fmt.Errorf("situation: game: rule_card %q is no side of the rule card; options %v", ref, sides)
}

// setUpPresents changes the presents from their start (R-KRAMPUS-05):
// counters left, destroyed ones in the treasure discard (R-KRAMPUS-08),
// and after Happy Holidays! each other one an item of the player who
// gained it (R-KRAMPUS-10). Nothing triggers: it is the table as it is.
func (g *Game) setUpPresents(ps []SituationPresent, holidays bool) error {
	c := g.Challenge
	kinds := make([]string, len(c.Presents))
	for i, p := range c.Presents {
		kinds[i] = p.Kind
	}
	if len(kinds) == 0 && (ps != nil || holidays) {
		return fmt.Errorf("situation: game: challenge %s has no presents; leave presents and holidays out", c.ID)
	}
	named := map[string]SituationPresent{}
	for _, sp := range ps {
		if !slices.Contains(kinds, sp.Kind) {
			return fmt.Errorf("situation: game: presents: %q is no present; options %v", sp.Kind, kinds)
		}
		if _, dup := named[sp.Kind]; dup {
			return fmt.Errorf("situation: game: presents: %s is named twice", sp.Kind)
		}
		named[sp.Kind] = sp
		if err := sp.check(holidays, len(g.Players)); err != nil {
			return err
		}
	}
	for i := range c.Presents {
		p := &c.Presents[i]
		sp, ok := named[p.Kind]
		switch {
		case ok && sp.Destroyed:
			p.Object, p.Destroyed = g.discard(p.Object, TreasureDeck), true
		case holidays && !ok, holidays && sp.To == nil:
			return fmt.Errorf("situation: game: holidays: the %s present names no player who gained it (to), or destroyed (R-KRAMPUS-10)", p.Kind)
		case ok && sp.Counters > 0:
			o := g.Object(p.Object)
			o.addCounters("", sp.Counters-o.CountersOf(""))
		}
	}
	if !holidays {
		return nil
	}
	c.Holidays = true
	for i := range c.Presents {
		p := &c.Presents[i]
		if p.Destroyed {
			continue
		}
		g.Object(p.Object).Counters = nil
		p.Flipped = true
		to := PlayerID(*named[p.Kind].To)
		p.Object = g.move(p.Object, Zone{Kind: ZoneInPlay}, to)
		g.enterAsItem(to, p.Object)
	}
	return nil
}

// check checks one present's fields against each other.
func (sp SituationPresent) check(holidays bool, seats int) error {
	at := "situation: game: presents: " + sp.Kind
	switch {
	case sp.Counters < 0:
		return fmt.Errorf("%s: counters %d is below 1; a present at 0 is destroyed (R-KRAMPUS-08)", at, sp.Counters)
	case sp.Destroyed && (sp.Counters > 0 || sp.To != nil):
		return fmt.Errorf("%s: a destroyed present has no counters and goes to nobody (R-KRAMPUS-08)", at)
	case sp.To != nil && !holidays:
		return fmt.Errorf("%s: to: a present is gained only at Happy Holidays! (R-KRAMPUS-10); add holidays", at)
	case sp.To != nil && sp.Counters > 0:
		return fmt.Errorf("%s: a gained present has no counters (R-KRAMPUS-10)", at)
	case sp.To != nil && (*sp.To < 0 || *sp.To >= seats):
		return fmt.Errorf("%s: to %d is not a seat, options 0 to %d", at, *sp.To, seats-1)
	case !sp.Destroyed && sp.Counters == 0 && sp.To == nil:
		return errors.New(at + ": names nothing to change: counters, destroyed or to")
	}
	return nil
}

// checkPresents compares the challenge's presents with want (24.34).
func (g *Game) checkPresents(c *checks, want []ExpectedPresent) {
	if want == nil {
		return
	}
	if g.Challenge == nil || len(g.Challenge.Presents) == 0 {
		c.add("presents: the game has no presents")
		return
	}
	for _, w := range want {
		p := g.present(w.Kind)
		if p == nil {
			kinds := make([]string, len(g.Challenge.Presents))
			for i, p := range g.Challenge.Presents {
				kinds[i] = p.Kind
			}
			c.add("presents: %q is no present, options %v", w.Kind, kinds)
			continue
		}
		who := "present " + w.Kind + " "
		c.int(who+"counters", w.Counters, g.PresentCounters(w.Kind))
		c.bool(who+"destroyed", w.Destroyed, p.Destroyed)
		if w.To != nil {
			o := g.Object(p.Object)
			if o.Zone.Kind != ZoneInPlay || int(o.Controller) != *w.To {
				c.add("%sis not an item of player %d", who, *w.To)
			}
		}
		if w.Counters == nil && w.Destroyed == nil && w.To == nil {
			c.add("%schecks nothing: name counters, destroyed or to", who)
		}
	}
}
