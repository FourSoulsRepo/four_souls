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
	// Counters are the counters on the card, e.g. charge counters.
	Counters []Counter `json:"counters,omitempty"`
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
	// LootPlays, Attacks and Purchases are what the player may still do
	// now, bonuses included; attacks and purchases only on their turn.
	LootPlays int `json:"loot_plays"`
	Attacks   int `json:"attacks"`
	Purchases int `json:"purchases"`
	// TreasureTop is the top treasure card, when this player may look at
	// it now (There's Options) and the viewer may see it.
	TreasureTop *CardView `json:"treasure_top,omitempty"`
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
	// Options says what each option of a choose prompt points at, one
	// per Waiting.Options; only for whoever sees those options.
	Options []OptionView `json:"options,omitempty"`
	// Passed lists who passed in a row in the open priority window,
	// first passer first; Responding lists who may still respond, the
	// waiting player first (V-11). Both are empty outside a window.
	Passed     []PlayerID `json:"passed,omitempty"`
	Responding []PlayerID `json:"responding,omitempty"`
}

// OptionView is what one option of a choose prompt points at, so a
// client can let the player pick it on the table. An option with none
// of these set is plain text.
type OptionView struct {
	Object   ObjectID `json:"object,omitempty"` // a card on the table or in the chooser's hand
	Player   PlayerID `json:"player"`           // a player, when IsPlayer
	IsPlayer bool     `json:"is_player,omitempty"`
	Deck     bool     `json:"deck,omitempty"`     // the monster or treasure deck
	StackID  int      `json:"stack_id,omitempty"` // a stack item
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
			LootPlays: g.lootPlaysFor(pl.ID),
		}
		if pl.ID == g.Turn.Active {
			pv.Attacks, pv.Purchases = g.attacksLeft(), g.purchasesLeft()
		}
		for _, id := range pl.InPlay {
			pv.InPlay = append(pv.InPlay, g.cardView(id))
		}
		if v.sees(pl.ID) {
			for _, id := range pl.Hand {
				pv.Hand = append(pv.Hand, g.cardView(id))
			}
			if top := g.DeckTop(TreasureDeck, 1); len(top) > 0 && g.peeksTreasure(pl.ID) {
				cv := g.cardView(top[0])
				pv.TreasureTop = &cv
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
	if g.Waiting.Kind == PromptChoose && v.sees(g.Waiting.Player) {
		out.Options = g.optionViews()
	}
	out.Passed, out.Responding = g.responders()
	return out
}

// responders splits the players of the open priority window into those
// who passed in a row and those still to answer, in turn order.
// Priority goes round once; the window closes when all passed
// (R-PRIO-05).
func (g *Game) responders() (passed, responding []PlayerID) {
	if g.Waiting.Kind != PromptPriority || !g.Priority.Open {
		return nil, nil
	}
	n := len(g.Players)
	p := g.Priority.Holder
	for range g.Priority.Passes {
		p = PlayerID((int(p) + n - 1) % n)
		passed = append([]PlayerID{p}, passed...)
	}
	p = g.Priority.Holder
	for range n - g.Priority.Passes {
		responding = append(responding, p)
		p = g.next(p)
	}
	return passed, responding
}

// optionViews maps the open choice's options to what they point at.
func (g *Game) optionViews() []OptionView {
	c := g.Choice
	if c == nil {
		return nil
	}
	out := make([]OptionView, len(g.Waiting.Options))
	for i := range out {
		out[i].Player = NoPlayer
	}
	set := func(i int, f func(o *OptionView)) {
		if i < len(out) {
			f(&out[i])
		}
	}
	switch c.Purpose {
	case ChooseAttackTarget, ChoosePurchase, ChoosePenaltyItem, ChoosePenaltyLoot, ChooseStartingItem, ChooseLowest:
		for i, id := range c.Objects {
			set(i, func(o *OptionView) { o.Object = id })
		}
		if c.Deck {
			set(len(c.Objects), func(o *OptionView) { o.Deck = true })
		}
	case ChooseTarget:
		for i, t := range c.Targets {
			set(i, func(o *OptionView) {
				switch t.Kind {
				case TargetPlayer, TargetOtherPlayer:
					o.Player, o.IsPlayer = t.Player, true
				case TargetDiceRoll, TargetStackAbility:
					o.StackID = t.StackID
				case TargetMonster, TargetMonsterOrPlayer, TargetItem, TargetCurse, TargetYourItem,
					TargetYourHandCard, TargetItemOrSoul:
					o.Object = t.Object
				}
			})
		}
	case ChooseMonsterSlot:
		for i, s := range c.Slots {
			if top, ok := g.Monsters[s].TopOf(); ok {
				set(i, func(o *OptionView) { o.Object = top })
			}
		}
	case ChooseCursed:
		for i, p := range c.Slots {
			set(i, func(o *OptionView) { o.Player, o.IsPlayer = PlayerID(p), true })
		}
	case ChooseReplacement, ChooseTriggerOrder, ChooseMode, ChooseAnswer:
		// Text options.
	}
	return out
}

func (g *Game) cardView(id ObjectID) CardView {
	o := g.Object(id)
	return CardView{
		ID: id, Card: o.Card, Role: o.Role, Charged: o.Charged, Damage: o.Damage,
		Counters: append([]Counter(nil), o.Counters...),
	}
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

// peeksTreasure: it is p's turn and p controls a card that lets them
// look at the top treasure card.
func (g *Game) peeksTreasure(p PlayerID) bool {
	if g.Turn.Active != p {
		return false
	}
	for _, id := range g.Players[p].InPlay {
		if g.Object(id).Role == RoleItem && g.def(id).PeeksTreasure {
			return true
		}
	}
	return false
}
