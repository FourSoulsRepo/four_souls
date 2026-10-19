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
var ratios = [ratioGroupCount]struct {
	name  string
	count int
}{
	GroupPassiveItems:         {"Passive items", 44},
	GroupActiveItems:          {"Active items", 40},
	GroupPaidItems:            {"Paid items", 10},
	GroupOneUseItems:          {"One-use items", 5},
	GroupSoulItems:            {"Soul items", 1},
	GroupPennies:              {"Pennies", 2},
	GroupTwoCents:             {"2¢ loot", 6},
	GroupThreeCents:           {"3¢ loot", 11},
	GroupFourCents:            {"4¢ loot", 12},
	GroupNickels:              {"Nickels", 6},
	GroupButterBeans:          {"Butter Beans", 5},
	GroupBombs:                {"Bombs", 6},
	GroupBatteries:            {"Batteries", 6},
	GroupDiceShardsSoulHearts: {"Dice Shards and Soul Hearts", 5},
	GroupPillsRunes:           {"Pills and Runes", 6},
	GroupWildcards:            {"Tarot and wildcards", 23},
	GroupTrinkets:             {"Trinkets", 11},
	GroupLostSouls:            {"Lost Souls", 1},
	GroupBasicMonsters:        {"Basic monsters", 30},
	GroupCursedMonsters:       {"Cursed monsters", 9},
	GroupHolyCharmedMonsters:  {"Holy and Charmed monsters", 9},
	GroupBosses:               {"Bosses", 30},
	GroupEpicBosses:           {"Epic bosses", 1},
	GroupGoodEvents:           {"Good events", 8},
	GroupBadEvents:            {"Bad events", 8},
	GroupCurses:               {"Curses", 5},
}

// Name is the group's name for players, e.g. "Butter Beans".
func (r RatioGroup) Name() string {
	if r == NoRatioGroup || r >= ratioGroupCount {
		return "no group"
	}
	return ratios[r].name
}

// Count is how many cards of the group a deck holds with deck ratios.
func (r RatioGroup) Count() int {
	if r >= ratioGroupCount {
		return 0
	}
	return ratios[r].count
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
