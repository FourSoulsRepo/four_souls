package protocol

import (
	"strconv"

	engine "github.com/FourSoulsRepo/rules_engine"
)

// appendEvents writes a list of events; nil is null.
func appendEvents(b []byte, es []engine.Event) []byte {
	if es == nil {
		return append(b, "null"...)
	}
	b = append(b, '[')
	for i, e := range es {
		if i > 0 {
			b = append(b, ',')
		}
		b = appendEvent(b, e)
	}
	return append(b, ']')
}

// appendEvent writes an event as encoding/json writes it.
func appendEvent(b []byte, e engine.Event) []byte {
	b = append(b, `{"kind":`...)
	b = appendString(b, string(e.Kind))
	b = append(b, `,"player":`...)
	b = strconv.AppendInt(b, int64(e.Player), 10)
	b = appendIntIfSet(b, "object", int(e.Object))
	if e.Card != "" {
		b = append(b, `,"card":`...)
		b = appendString(b, string(e.Card))
	}
	b = appendIntIfSet(b, "amount", e.Amount)
	if e.Text != "" {
		b = append(b, `,"text":`...)
		b = appendString(b, e.Text)
	}
	if !e.Label.IsZero() {
		b = appendMarshaled(b, `,"label":`, e.Label)
	}
	if e.Private {
		b = append(b, `,"private":true`...)
	}
	if len(e.SeenBy) > 0 {
		b = appendMarshaled(b, `,"seen_by":`, e.SeenBy)
	}
	b = appendIntIfSet(b, "prev", int(e.Prev))
	b = appendIntIfSet(b, "source", int(e.Source))
	b = appendIntIfSet(b, "stack_id", e.StackID)
	if len(e.Dice) > 0 {
		b = appendMarshaled(b, `,"dice":`, e.Dice)
	}
	b = appendIntIfSet(b, "by", e.By)
	if e.Item != nil {
		b = appendMarshaled(b, `,"item":`, e.Item)
	}
	if len(e.Sources) > 0 {
		b = appendMarshaled(b, `,"sources":`, e.Sources)
	}
	return append(b, '}')
}

// appendIntents writes a list of intents; nil is null.
func appendIntents(b []byte, ins []engine.Intent) []byte {
	if ins == nil {
		return append(b, "null"...)
	}
	b = append(b, '[')
	for i, in := range ins {
		if i > 0 {
			b = append(b, ',')
		}
		b = append(b, `{"player":`...)
		b = strconv.AppendInt(b, int64(in.Player), 10)
		b = append(b, `,"kind":`...)
		b = strconv.AppendInt(b, int64(in.Kind), 10)
		if len(in.Objects) > 0 {
			b = append(b, `,"objects":`...)
			b = append(b, '[')
			for j, o := range in.Objects {
				if j > 0 {
					b = append(b, ',')
				}
				b = strconv.AppendInt(b, int64(o), 10)
			}
			b = append(b, ']')
		}
		b = appendIntIfSet(b, "choice", in.Choice)
		b = append(b, '}')
	}
	return append(b, ']')
}

// appendSeats writes the seats of the table; nil is null.
func appendSeats(b []byte, ss []Seat) []byte {
	if ss == nil {
		return append(b, "null"...)
	}
	b = append(b, '[')
	for i, s := range ss {
		if i > 0 {
			b = append(b, ',')
		}
		b = append(b, `{"seat":`...)
		b = strconv.AppendInt(b, int64(s.Seat), 10)
		b = append(b, `,"name":`...)
		b = appendString(b, s.Name)
		b = append(b, `,"connected":`...)
		b = strconv.AppendBool(b, s.Connected)
		if s.Ready {
			b = append(b, `,"ready":true`...)
		}
		if s.Kicked {
			b = append(b, `,"kicked":true`...)
		}
		if s.Dice != "" {
			b = append(b, `,"dice":`...)
			b = appendString(b, s.Dice)
		}
		b = append(b, '}')
	}
	return append(b, ']')
}
