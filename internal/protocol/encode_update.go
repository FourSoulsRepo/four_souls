package protocol

import (
	"strconv"

	engine "github.com/FourSoulsRepo/rules_engine"
)

// appendUpdate writes an update as encoding/json writes it: the fields in
// struct order, each with its omitempty rule (18.17). The view is written
// by hand; the other parts by encoding/json.
func appendUpdate(b []byte, u Update) []byte {
	b = append(b, `{"step":`...)
	b = strconv.AppendInt(b, int64(u.Step), 10)
	b = append(b, `,"events":`...)
	b = appendEvents(b, u.Events)
	b = append(b, `,"view":`...)
	b = appendView(b, u.View)
	b = append(b, `,"allowed":`...)
	b = appendIntents(b, u.Allowed)
	if len(u.Aims) > 0 {
		b = appendMarshaled(b, `,"aims":`, u.Aims)
	}
	if len(u.Blocked) > 0 {
		b = appendMarshaled(b, `,"blocked":`, u.Blocked)
	}
	b = append(b, `,"seats":`...)
	b = appendSeats(b, u.Seats)
	if u.Deadline != 0 {
		b = append(b, `,"deadline":`...)
		b = strconv.AppendInt(b, u.Deadline, 10)
	}
	if u.WaitTurn {
		b = append(b, `,"wait_turn":true`...)
	}
	if u.Pause != nil {
		b = appendMarshaled(b, `,"pause":`, u.Pause)
	}
	if u.RollbackMode != "" {
		b = appendMarshaled(b, `,"rollback_mode":`, u.RollbackMode)
	}
	if u.RollbackFrom != 0 {
		b = append(b, `,"rollback_from":`...)
		b = strconv.AppendInt(b, int64(u.RollbackFrom), 10)
	}
	if u.Rollback != nil {
		b = appendMarshaled(b, `,"rollback":`, u.Rollback)
	}
	if u.RollbackLast != nil {
		b = appendMarshaled(b, `,"rollback_last":`, u.RollbackLast)
	}
	if u.RollbackVote != nil {
		b = appendMarshaled(b, `,"rollback_vote":`, u.RollbackVote)
	}
	if u.RollbackEnd != nil {
		b = appendMarshaled(b, `,"rollback_end":`, u.RollbackEnd)
	}
	if u.RollbackCut != 0 {
		b = append(b, `,"rollback_cut":`...)
		b = strconv.AppendInt(b, int64(u.RollbackCut), 10)
	}
	return append(b, '}')
}

// appendView writes the game view as encoding/json writes it. Its lists of
// players, slots and cards are written by hand; the rest by encoding/json.
func appendView(b []byte, v engine.GameView) []byte {
	b = append(b, `{"players":`...)
	b = appendPlayers(b, v.Players)
	b = append(b, `,"shop":`...)
	b = appendSlots(b, v.Shop)
	b = append(b, `,"monsters":`...)
	b = appendSlots(b, v.Monsters)
	b = append(b, `,"deck_sizes":`...)
	b = appendInts(b, v.DeckSizes[:])
	b = append(b, `,"discards":`...)
	b = appendCardPtrs(b, v.Discards[:])
	b = append(b, `,"discard_sizes":`...)
	b = appendInts(b, v.DiscardSizes[:])
	if len(v.BonusSouls) > 0 {
		b = append(b, `,"bonus_souls":`...)
		b = appendCards(b, v.BonusSouls)
	}
	b = append(b, `,"stack":`...)
	b = appendStack(b, v.Stack)
	b = append(b, `,"turn":`...)
	b = appendTurn(b, v.Turn)
	b = append(b, `,"waiting":`...)
	b = appendPrompt(b, v.Waiting)
	if v.Over {
		b = append(b, `,"over":true`...)
	}
	if len(v.Winners) > 0 {
		b = appendMarshaled(b, `,"winners":`, v.Winners)
	}
	b = append(b, `,"win_souls":`...)
	b = strconv.AppendInt(b, int64(v.WinSouls), 10)
	if v.Revealed != nil {
		b = append(b, `,"revealed":`...)
		b = appendCard(b, *v.Revealed)
	}
	if v.Subject != nil {
		b = append(b, `,"subject":`...)
		b = appendCard(b, *v.Subject)
	}
	if len(v.Ordering) > 0 {
		b = append(b, `,"ordering":`...)
		b = appendCards(b, v.Ordering)
	}
	if len(v.Options) > 0 {
		b = appendMarshaled(b, `,"options":`, v.Options)
	}
	if v.Activating != nil {
		b = appendMarshaled(b, `,"activating":`, v.Activating)
	}
	if len(v.Modes) > 0 {
		b = appendMarshaled(b, `,"modes":`, v.Modes)
	}
	if len(v.Passed) > 0 {
		b = appendMarshaled(b, `,"passed":`, v.Passed)
	}
	if len(v.Responding) > 0 {
		b = appendMarshaled(b, `,"responding":`, v.Responding)
	}
	if v.Attack != nil {
		b = appendMarshaled(b, `,"attack":`, v.Attack)
	}
	if v.Buying {
		b = append(b, `,"buying":true`...)
	}
	if v.FirstRoll != nil {
		b = appendMarshaled(b, `,"first_roll":`, v.FirstRoll)
	}
	return append(b, '}')
}

// appendPlayers writes a list of player views; nil is null.
func appendPlayers(b []byte, ps []engine.PlayerView) []byte {
	if ps == nil {
		return append(b, "null"...)
	}
	b = append(b, '[')
	for i, p := range ps {
		if i > 0 {
			b = append(b, ',')
		}
		b = appendPlayer(b, p)
	}
	return append(b, ']')
}

// appendSlots writes a list of table slots; nil is null.
func appendSlots(b []byte, ss []engine.SlotView) []byte {
	if ss == nil {
		return append(b, "null"...)
	}
	b = append(b, '[')
	for i, s := range ss {
		if i > 0 {
			b = append(b, ',')
		}
		b = appendSlot(b, s)
	}
	return append(b, ']')
}

// appendEnvelope writes the start of an envelope, as encoding/json writes
// Envelope: type, id when it is not 0, then the data key. The caller
// writes the data and the closing brace.
func appendEnvelope(b []byte, typ string, id int) []byte {
	b = append(b, `{"type":`...)
	b = appendString(b, typ)
	if id != 0 {
		b = append(b, `,"id":`...)
		b = strconv.AppendInt(b, int64(id), 10)
	}
	return append(b, `,"data":`...)
}
