package rulesengine

import "fmt"

// CardKind is the printed card type, as far as the rules care.
type CardKind int

// The card kinds (R-CARD).
const (
	CharacterCard CardKind = iota
	TreasureCard
	LootCard
	MonsterCard // has a stat block: a monster in play (R-CARD-09)
	EventCard   // monster-deck card without stats (R-CARD-09)
	BonusSoulCard
	RoomCard
)

// Deck returns the deck a card kind belongs to; ok is false for kinds
// without a deck (R-ZONE-04).
func (k CardKind) Deck() (DeckKind, bool) {
	switch k {
	case TreasureCard:
		return TreasureDeck, true
	case LootCard:
		return LootDeck, true
	case MonsterCard, EventCard:
		return MonsterDeck, true
	case RoomCard:
		return RoomDeck, true
	case CharacterCard, BonusSoulCard:
		return 0, false
	default:
		return 0, false
	}
}

// CardDef is what the engine knows about one card version.
// Effects are added as blocks (step 4.7); display data lives in card_db.
type CardDef struct {
	Ref    CardRef
	Kind   CardKind
	Copies int // copies in the set

	HP, ATK, DC int // stat block; DC is the evasion (R-CARD-28)
	Soul        int // soul value, 0 if none (R-CARD-18)

	StartingItem CardRef // characters only (R-CARD-26)
	Eternal      bool    // keyword (R-ABIL-18)
	// Outside marks cards that start outside the game, such as starting
	// items: they are never shuffled into a deck (R-ZONE-13).
	Outside bool

	// Replacements are the card's replacement effects (R-ABIL-29).
	Replacements []Replacement
	// Abilities are activated, loot and triggered abilities (step 4.7).
	Abilities []Ability
	// Statics are static abilities that change stats (R-ABIL-12).
	Statics []Static

	// Rewards are gained by the active player when it dies (R-CARD-31).
	Rewards []Reward
	// Tap marks cards with a ↷ ability; the death penalty deactivates
	// them (R-DEATH-14). Characters always have one.
	Tap bool
}

// CardSet is a set of card definitions, e.g. the Base Game.
type CardSet struct {
	Name  string
	Cards []CardDef
}

// cardIndex finds definitions by ref. It is rebuilt from the sets and is
// not part of the saved state.
type cardIndex struct {
	defs []CardDef
}

func newCardIndex(sets ...CardSet) (cardIndex, error) {
	var idx cardIndex
	for _, s := range sets {
		for _, d := range s.Cards {
			if _, dup := idx.find(d.Ref); dup {
				return cardIndex{}, fmt.Errorf("card %s defined twice", d.Ref)
			}
			idx.defs = append(idx.defs, d)
		}
	}
	return idx, nil
}

func (idx cardIndex) find(ref CardRef) (CardDef, bool) {
	for _, d := range idx.defs {
		if d.Ref == ref {
			return d, true
		}
	}
	return CardDef{}, false
}

// RewardKind is what a reward box gives (R-CARD-31).
type RewardKind int

// The reward kinds.
const (
	RewardCents RewardKind = iota
	RewardLoot
	RewardTreasure
)

// Reward is one line of a reward box, e.g. 3 cents.
type Reward struct {
	Kind   RewardKind
	Amount int
}

func (r Reward) action() ActionKind {
	switch r.Kind {
	case RewardLoot:
		return ActLoot
	case RewardTreasure:
		return ActGainTreasure
	case RewardCents:
		return ActGainCents
	default:
		return ActGainCents
	}
}
