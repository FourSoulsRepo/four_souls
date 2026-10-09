package rulesengine

// ViewerKind is who looks at the game (A-07, SP-01, SP-02).
type ViewerKind int

// The viewer kinds.
const (
	ViewPlayer    ViewerKind = iota // a seated player: sees their own hand
	ViewSpectator                   // public information only
	ViewJudge                       // sees every hand, never deck order or the RNG
)

// Viewer is one person looking at the game.
type Viewer struct {
	Kind   ViewerKind
	Player PlayerID // for ViewPlayer
}

// CardView is a visible card.
type CardView struct {
	ID      ObjectID `json:"id"`
	Card    CardRef  `json:"card"`
	Role    Role     `json:"role,omitempty"`
	Charged bool     `json:"charged"`
	Damage  int      `json:"damage,omitempty"`
}

// PlayerView is what a viewer sees of one player.
type PlayerView struct {
	ID        PlayerID   `json:"id"`
	Character CardView   `json:"character"`
	HP        int        `json:"hp"`
	ATK       int        `json:"atk"`
	Cents     int        `json:"cents"`
	Souls     int        `json:"souls"`
	Dead      bool       `json:"dead,omitempty"`
	InPlay    []CardView `json:"in_play"`
	HandSize  int        `json:"hand_size"`
	Hand      []CardView `json:"hand,omitempty"` // only when the viewer may see it
}

// SlotView is a table slot: the top card and the covered ones below.
type SlotView struct {
	Top     *CardView  `json:"top,omitempty"`
	Covered []CardView `json:"covered,omitempty"`
}

// GameView is everything one viewer may see.
type GameView struct {
	Players    []PlayerView          `json:"players"`
	Shop       []SlotView            `json:"shop"`
	Monsters   []SlotView            `json:"monsters"`
	DeckSizes  [deckCount]int        `json:"deck_sizes"`
	Discards   [deckCount][]CardView `json:"discards"`
	BonusSouls []CardView            `json:"bonus_souls,omitempty"`
	Stack      []StackItem           `json:"stack"`
	Turn       Turn                  `json:"turn"`
	Waiting    Prompt                `json:"waiting"`
	Over       bool                  `json:"over,omitempty"`
	Winners    []PlayerID            `json:"winners,omitempty"`
}

// sees reports whether the viewer may see player p's hand.
func (v Viewer) sees(p PlayerID) bool {
	return v.Kind == ViewJudge || (v.Kind == ViewPlayer && v.Player == p)
}

// View builds what a viewer sees. It is the only way state reaches
// clients (A-07): hidden hands, deck contents and order, and the RNG
// never leave the engine.
func (g *Game) View(v Viewer) GameView {
	out := GameView{Turn: g.Turn, Over: g.Over, Winners: g.Winners, Stack: append([]StackItem(nil), g.Stack...)}
	for _, pl := range g.Players {
		pv := PlayerView{
			ID: pl.ID, Character: g.cardView(pl.Character), HP: g.PlayerHP(pl.ID), ATK: g.PlayerATK(pl.ID),
			Cents: pl.Cents, Souls: g.SoulValue(pl.ID), Dead: pl.Dead, HandSize: len(pl.Hand),
		}
		for _, id := range pl.InPlay {
			pv.InPlay = append(pv.InPlay, g.cardView(id))
		}
		if v.sees(pl.ID) {
			for _, id := range pl.Hand {
				pv.Hand = append(pv.Hand, g.cardView(id))
			}
		}
		out.Players = append(out.Players, pv)
	}
	out.Shop = g.slotViews(g.Shop)
	out.Monsters = g.slotViews(g.Monsters)
	for d := range deckCount {
		out.DeckSizes[d] = len(g.Decks[d])
		for _, id := range g.Discards[d] {
			out.Discards[d] = append(out.Discards[d], g.cardView(id))
		}
	}
	for _, id := range g.BonusSouls {
		out.BonusSouls = append(out.BonusSouls, g.cardView(id))
	}
	out.Waiting = g.Waiting
	if g.Waiting.Kind == PromptChoose && !v.sees(g.Waiting.Player) {
		// Options may name hidden cards (e.g. which loot card to discard).
		out.Waiting.Options = nil
	}
	return out
}

func (g *Game) cardView(id ObjectID) CardView {
	o := g.Object(id)
	return CardView{ID: id, Card: o.Card, Role: o.Role, Charged: o.Charged, Damage: o.Damage}
}

func (g *Game) slotViews(slots []Slot) []SlotView {
	out := make([]SlotView, len(slots))
	for i, s := range slots {
		for j, id := range s.Cards {
			cv := g.cardView(id)
			if j == len(s.Cards)-1 {
				out[i].Top = &cv
			} else {
				out[i].Covered = append(out[i].Covered, cv) // covered cards are public (R-ZONE-12)
			}
		}
	}
	return out
}

// FilterEvents returns the events a viewer may see; a private event (such
// as a looted card) shows its card only to its player and judges.
func FilterEvents(events []Event, v Viewer) []Event {
	out := make([]Event, 0, len(events))
	for _, e := range events {
		if e.Private && !v.sees(e.Player) {
			e.Card, e.Object = "", 0
		}
		out = append(out, e)
	}
	return out
}
