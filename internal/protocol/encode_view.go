package protocol

import (
	"strconv"

	engine "github.com/FourSoulsRepo/rules_engine"
)

// appendCards writes a list of cards. A nil list is null, as encoding/json
// writes a nil slice; an empty one is [].
func appendCards(b []byte, cs []engine.CardView) []byte {
	if cs == nil {
		return append(b, "null"...)
	}
	b = append(b, '[')
	for i, c := range cs {
		if i > 0 {
			b = append(b, ',')
		}
		b = appendCard(b, c)
	}
	return append(b, ']')
}

// appendPlayer writes a player view as encoding/json writes it.
func appendPlayer(b []byte, p engine.PlayerView) []byte {
	b = append(b, `{"id":`...)
	b = strconv.AppendInt(b, int64(p.ID), 10)
	b = append(b, `,"character":`...)
	b = appendCard(b, p.Character)
	b = append(b, `,"hp":`...)
	b = strconv.AppendInt(b, int64(p.HP), 10)
	b = append(b, `,"atk":`...)
	b = strconv.AppendInt(b, int64(p.ATK), 10)
	b = append(b, `,"cents":`...)
	b = strconv.AppendInt(b, int64(p.Cents), 10)
	b = append(b, `,"souls":`...)
	b = strconv.AppendInt(b, int64(p.Souls), 10)
	if p.Dead {
		b = append(b, `,"dead":true`...)
	}
	b = append(b, `,"in_play":`...)
	b = appendCards(b, p.InPlay)
	b = append(b, `,"hand_size":`...)
	b = strconv.AppendInt(b, int64(p.HandSize), 10)
	if len(p.Hand) > 0 {
		b = append(b, `,"hand":`...)
		b = appendCards(b, p.Hand)
	}
	b = append(b, `,"loot_plays":`...)
	b = strconv.AppendInt(b, int64(p.LootPlays), 10)
	b = append(b, `,"attacks":`...)
	b = strconv.AppendInt(b, int64(p.Attacks), 10)
	b = append(b, `,"purchases":`...)
	b = strconv.AppendInt(b, int64(p.Purchases), 10)
	if p.TreasureTop != nil {
		b = append(b, `,"treasure_top":`...)
		b = appendCard(b, *p.TreasureTop)
	}
	b = append(b, `,"printed_hp":`...)
	b = strconv.AppendInt(b, int64(p.PrintedHP), 10)
	b = append(b, `,"printed_atk":`...)
	b = strconv.AppendInt(b, int64(p.PrintedATK), 10)
	if len(p.Sources) > 0 {
		b = appendMarshaled(b, `,"sources":`, p.Sources)
	}
	return append(b, '}')
}

// appendSlot writes a table slot as encoding/json writes it.
func appendSlot(b []byte, s engine.SlotView) []byte {
	b = append(b, '{')
	sep := byte(0)
	if s.Top != nil {
		b = append(b, `"top":`...)
		b = appendCard(b, *s.Top)
		sep = ','
	}
	if len(s.Covered) > 0 {
		if sep != 0 {
			b = append(b, sep)
		}
		b = append(b, `"covered":`...)
		b = appendCards(b, s.Covered)
	}
	return append(b, '}')
}

// appendTurn writes the turn as encoding/json writes it.
func appendTurn(b []byte, t engine.Turn) []byte {
	b = append(b, `{"active":`...)
	b = strconv.AppendInt(b, int64(t.Active), 10)
	b = append(b, `,"number":`...)
	b = strconv.AppendInt(b, int64(t.Number), 10)
	b = append(b, `,"step":`...)
	b = strconv.AppendInt(b, int64(t.Step), 10)
	b = append(b, `,"loot_plays":`...)
	b = strconv.AppendInt(b, int64(t.LootPlays), 10)
	b = append(b, `,"attacks":`...)
	b = strconv.AppendInt(b, int64(t.Attacks), 10)
	b = append(b, `,"purchases":`...)
	b = strconv.AppendInt(b, int64(t.Purchases), 10)
	b = appendIntIfSet(b, "attack_rolls", t.AttackRolls)
	b = appendIntIfSet(b, "bonus_loot_used", t.BonusLootUsed)
	b = appendIntIfSet(b, "bonus_attacks_used", t.BonusAttacksUsed)
	b = appendIntIfSet(b, "bonus_purchases_used", t.BonusPurchasesUsed)
	b = appendIntIfSet(b, "must_attack", int(t.MustAttack))
	b = appendIntIfSet(b, "must_attacks", t.MustAttacks)
	b = appendIntIfSet(b, "must_attack_deck", t.MustAttackDeck)
	b = appendIntIfSet(b, "deck_attacks", t.DeckAttacks)
	b = appendIntIfSet(b, "shop_bought", t.ShopBought)
	if t.FreeShopItem {
		b = append(b, `,"free_shop_item":true`...)
	}
	if t.EndDeclared {
		b = append(b, `,"end_declared":true`...)
	}
	if t.DeathEnd {
		b = append(b, `,"death_end":true`...)
	}
	if len(t.NextHelpers) > 0 {
		b = appendMarshaled(b, `,"next_helpers":`, t.NextHelpers)
	}
	if len(t.NextShare) > 0 {
		b = appendMarshaled(b, `,"next_share":`, t.NextShare)
	}
	if t.MonsterDied {
		b = append(b, `,"monster_died":true`...)
	}
	if t.RoomReplaced {
		b = append(b, `,"room_replaced":true`...)
	}
	if len(t.LateRooms) > 0 {
		b = appendMarshaled(b, `,"late_rooms":`, t.LateRooms)
	}
	if t.Entered {
		b = append(b, `,"entered":true`...)
	}
	return append(b, '}')
}

// appendIntIfSet writes ,"key":n when n is not zero, as an omitempty int.
func appendIntIfSet(b []byte, key string, n int) []byte {
	if n == 0 {
		return b
	}
	b = append(b, `,"`...)
	b = append(b, key...)
	b = append(b, `":`...)
	return strconv.AppendInt(b, int64(n), 10)
}

// appendInts writes a list of ints as a JSON array.
func appendInts(b []byte, ns []int) []byte {
	b = append(b, '[')
	for i, n := range ns {
		if i > 0 {
			b = append(b, ',')
		}
		b = strconv.AppendInt(b, int64(n), 10)
	}
	return append(b, ']')
}

// appendCardPtrs writes a list of optional cards: null for a nil card.
func appendCardPtrs(b []byte, cs []*engine.CardView) []byte {
	b = append(b, '[')
	for i, c := range cs {
		if i > 0 {
			b = append(b, ',')
		}
		if c == nil {
			b = append(b, "null"...)
		} else {
			b = appendCard(b, *c)
		}
	}
	return append(b, ']')
}

// appendPrompt writes the prompt the game waits on as encoding/json writes
// it. The question is left out when it has no line (omitzero).
func appendPrompt(b []byte, p engine.Prompt) []byte {
	b = append(b, `{"kind":`...)
	b = appendMarshaled(b, "", p.Kind)
	b = append(b, `,"player":`...)
	b = strconv.AppendInt(b, int64(p.Player), 10)
	b = appendIntIfSet(b, "count", p.Count)
	if len(p.Options) > 0 {
		b = appendMarshaled(b, `,"options":`, p.Options)
	}
	if p.Purpose != 0 {
		b = appendMarshaled(b, `,"purpose":`, p.Purpose)
	}
	if !p.Text.IsZero() {
		b = appendMarshaled(b, `,"text":`, p.Text)
	}
	return append(b, '}')
}

// appendStack writes the stack. Its items are written by encoding/json;
// an empty stack needs no call: nil is null, an empty list is [].
func appendStack(b []byte, st []engine.StackItem) []byte {
	switch {
	case st == nil:
		return append(b, "null"...)
	case len(st) == 0:
		return append(b, "[]"...)
	}
	return appendMarshaled(b, "", st)
}
