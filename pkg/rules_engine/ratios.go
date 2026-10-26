package rulesengine

import (
	"fmt"
	"slices"
	"strings"
)

// RatioGroup is a card's group in the official deck ratios (R-SETUP-17,
// GS-07). The groups follow the site's card types.
type RatioGroup uint8

// The deck ratio groups; NoRatioGroup is every card outside the
// treasure, loot and monster decks.
const (
	NoRatioGroup RatioGroup = iota
	GroupPassiveItems
	GroupActiveItems
	GroupPaidItems
	GroupOneUseItems
	GroupSoulItems
	GroupPennies
	GroupTwoCents
	GroupThreeCents
	GroupFourCents
	GroupNickels
	GroupButterBeans
	GroupBombs
	GroupBatteries
	GroupDiceShardsSoulHearts
	GroupPillsRunes
	GroupWildcards
	GroupTrinkets
	GroupLostSouls
	GroupBasicMonsters
	GroupCursedMonsters
	GroupHolyCharmedMonsters
	GroupBosses
	GroupEpicBosses
	GroupGoodEvents
	GroupBadEvents
	GroupCurses
	ratioGroupCount
)

// ratios are the official counts per group: 100 cards in each deck
// (S-OFF Deck Ratios).
var ratios = [ratioGroupCount]int{
	GroupPassiveItems:         44,
	GroupActiveItems:          40,
	GroupPaidItems:            10,
	GroupOneUseItems:          5,
	GroupSoulItems:            1,
	GroupPennies:              2,
	GroupTwoCents:             6,
	GroupThreeCents:           11,
	GroupFourCents:            12,
	GroupNickels:              6,
	GroupButterBeans:          5,
	GroupBombs:                6,
	GroupBatteries:            6,
	GroupDiceShardsSoulHearts: 5,
	GroupPillsRunes:           6,
	GroupWildcards:            23,
	GroupTrinkets:             11,
	GroupLostSouls:            1,
	GroupBasicMonsters:        30,
	GroupCursedMonsters:       9,
	GroupHolyCharmedMonsters:  9,
	GroupBosses:               30,
	GroupEpicBosses:           1,
	GroupGoodEvents:           8,
	GroupBadEvents:            8,
	GroupCurses:               5,
}

// Name is the group's code, e.g. "butter_beans"; players read it from
// its key (ratioKeys, ADR 021).
func (r RatioGroup) Name() string {
	if r >= ratioGroupCount {
		r = NoRatioGroup
	}
	k := ratioKeys[r].String()
	return k[strings.LastIndexByte(k, '.')+1:]
}

// Count is how many cards of the group a deck holds with deck ratios.
func (r RatioGroup) Count() int {
	if r >= ratioGroupCount {
		return 0
	}
	return ratios[r]
}

// RatioShort is a group with too few cards for the deck ratios.
type RatioShort struct {
	Group RatioGroup
	Have  int
	Need  int
}

func (s RatioShort) String() string {
	return fmt.Sprintf("%s %d of %d", s.Group.Name(), s.Have, s.Need)
}

// Msg is the shortfall for the player: "Butter Beans 3 of 5".
func (s RatioShort) Msg() Msg {
	g := s.Group
	if g >= ratioGroupCount {
		g = NoRatioGroup
	}
	return ratioKeys[g].Msg(N("have", s.Have), N("need", s.Need))
}

// RatioShortfall lists the groups that have too few cards for the deck
// ratios in a game with these sets, players and bans; nil means the
// ratios fit. Cards left out of the game do not count (R-CARD-32, GS-13).
func RatioShortfall(s Setup) ([]RatioShort, error) {
	idx, err := newCardIndex(s.Sets...)
	if err != nil {
		return nil, fmt.Errorf("setup: %w", err)
	}
	pools, err := ratioPools(idx, s)
	if err != nil {
		return nil, err
	}
	return shortfall(pools), nil
}

// ratioPools are the cards each group may take, one ref per copy, in a
// fixed order.
func ratioPools(idx cardIndex, s Setup) ([ratioGroupCount][]CardRef, error) {
	var pools [ratioGroupCount][]CardRef
	for i, d := range idx.defs {
		if !inDecks(d, s) || d.Kind == RoomCard {
			continue
		}
		if _, hasDeck := d.Kind.Deck(); !hasDeck {
			continue
		}
		g := idx.groups[i]
		if g == NoRatioGroup || int(g) >= len(pools) {
			return pools, fmt.Errorf("setup: card %s has no deck ratio group", d.Ref)
		}
		for range idx.copies(i) {
			pools[g] = append(pools[g], d.Ref)
		}
	}
	return pools, nil
}

func shortfall(pools [ratioGroupCount][]CardRef) []RatioShort {
	var out []RatioShort
	for g := GroupPassiveItems; g < ratioGroupCount; g++ {
		if have := len(pools[g]); have < g.Count() {
			out = append(out, RatioShort{Group: g, Have: have, Need: g.Count()})
		}
	}
	return out
}

// inDecks reports whether a card may go into a deck of this game.
func inDecks(d CardDef, s Setup) bool {
	if d.Outside || (d.ThreePlus && s.Players < 3) { // R-CARD-32
		return false
	}
	if slices.Contains(s.Banned, d.Ref) { // GS-13
		return false
	}
	return d.Kind != RoomCard || s.Rooms // R-SETUP-16
}

// addRatioDecks fills the treasure, loot and monster decks with a random
// pick of each group's count (R-SETUP-17).
func (g *Game) addRatioDecks(s Setup) error {
	pools, err := ratioPools(g.cards, s)
	if err != nil {
		return err
	}
	if short := shortfall(pools); len(short) > 0 {
		parts := make([]string, len(short))
		for i, sh := range short {
			parts[i] = sh.String()
		}
		return fmt.Errorf("setup: too few cards for the deck ratios: %s", strings.Join(parts, ", "))
	}
	for grp := GroupPassiveItems; grp < ratioGroupCount; grp++ {
		pool := pools[grp]
		g.RNG.Shuffle(len(pool), func(i, j int) { pool[i], pool[j] = pool[j], pool[i] })
		for _, ref := range pool[:grp.Count()] {
			d, _ := g.cards.find(ref)
			deck, _ := d.Kind.Deck()
			g.AddToDeck(deck, ref)
		}
	}
	return nil
}
