// Package setup is the binding of the host screen's game setup: the card
// sets, the challenges and the deck ratios (ADR 030).
package setup

import (
	"github.com/FourSoulsRepo/four_souls/internal/msgerr"
	engine "github.com/FourSoulsRepo/rules_engine"
	"github.com/FourSoulsRepo/rules_engine/cards"
)

// errDeckRatios is the line the player reads when the ratios fail to
// count (ADR 021); its words are in frontend/src/locales/en.json.
var errDeckRatios = engine.NewKey("error.deck_ratios")

// Setup tells the host screen what this app can set up.
type Setup struct{}

// New makes the binding.
func New() *Setup {
	return &Setup{}
}

// CardSets returns the codes of the card sets this app has; the client
// sends them when it connects (CD-03). The challenge set is not among
// them: it is never a playable set (ADR 028).
func (*Setup) CardSets() []string {
	var codes []string
	for _, s := range cards.Sets() {
		codes = append(codes, s.Code)
	}
	return codes
}

// ChallengeInfo is a challenge the host screen offers (ADR 028,
// Decision 10.7): its ID and difficulties, easiest first. The client
// names them; a server that lacks one refuses it at create.
type ChallengeInfo struct {
	ID           string   `json:"id"`
	Difficulties []string `json:"difficulties"`
}

// Challenges lists the challenges this app knows, like CardSets lists
// its sets.
func (*Setup) Challenges() []ChallengeInfo {
	var out []ChallengeInfo
	for _, c := range cards.Challenges() {
		info := ChallengeInfo{ID: c.ID}
		for _, l := range c.Levels {
			info.Difficulties = append(info.Difficulties, l.Difficulty)
		}
		out = append(out, info)
	}
	return out
}

// DeckRatios says why a game with these sets and players can't use the
// deck ratios, or "" when it can (GS-07). The host screen shows it as
// the reason the box is off. The reason is a JSON list of messages
// (ADR 021): the head, {n} players, then one line per short group.
func (*Setup) DeckRatios(sets []string, players int) string {
	s := engine.Setup{Players: players}
	for _, code := range sets {
		if set, ok := cards.Find(code); ok {
			s.Sets = append(s.Sets, set)
		}
	}
	short, err := engine.RatioShortfall(s)
	if err != nil {
		return msgerr.JSON([]engine.Msg{errDeckRatios.Msg(engine.Str("error", err.Error()))})
	}
	if len(short) == 0 {
		return ""
	}
	out := []engine.Msg{engine.RefuseDeckRatios.Msg(engine.N("n", players))}
	for _, sh := range short {
		out = append(out, sh.Msg())
	}
	return msgerr.JSON(out)
}
