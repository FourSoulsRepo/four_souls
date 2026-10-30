package protocol

import (
	"encoding/json"
	"strconv"

	engine "github.com/FourSoulsRepo/rules_engine"
)

// The hand-written encoders of the update (18.17). Each writes the bytes
// encoding/json writes for its type, field by field, so the wire format
// stays the same (ADR 006). Parts not written here are marshalled by
// encoding/json and appended; a golden test checks every part against it.

// appendCard writes a card view as encoding/json writes it.
func appendCard(b []byte, c engine.CardView) []byte {
	b = append(b, `{"id":`...)
	b = strconv.AppendInt(b, int64(c.ID), 10)
	b = append(b, `,"card":`...)
	b = appendString(b, string(c.Card))
	if c.Role != 0 {
		b = append(b, `,"role":`...)
		b = strconv.AppendInt(b, int64(c.Role), 10)
	}
	b = append(b, `,"charged":`...)
	b = strconv.AppendBool(b, c.Charged)
	if c.Damage != 0 {
		b = append(b, `,"damage":`...)
		b = strconv.AppendInt(b, int64(c.Damage), 10)
	}
	if len(c.Counters) > 0 {
		b = appendMarshaled(b, `,"counters":`, c.Counters)
	}
	if c.Stats != nil {
		b = appendMarshaled(b, `,"stats":`, c.Stats)
	}
	if len(c.Abilities) > 0 {
		b = appendMarshaled(b, `,"abilities":`, c.Abilities)
	}
	if c.Souls != 0 {
		b = append(b, `,"souls":`...)
		b = strconv.AppendInt(b, int64(c.Souls), 10)
	}
	if len(c.Dice) > 0 {
		b = appendMarshaled(b, `,"dice":`, c.Dice)
	}
	if c.Flipped {
		b = append(b, `,"flipped":true`...)
	}
	if c.Copy != "" {
		b = append(b, `,"copy":`...)
		b = appendString(b, string(c.Copy))
	}
	return append(b, '}')
}

// appendMarshaled writes key and then v as encoding/json writes it. An
// error cannot happen for the types passed here.
func appendMarshaled(b []byte, key string, v any) []byte {
	raw, err := json.Marshal(v)
	if err != nil {
		panic(err) // our own plain types always marshal
	}
	b = append(b, key...)
	return append(b, raw...)
}

// appendString writes s as a JSON string. Plain text is written as it is;
// anything else takes encoding/json, which escapes it the same way.
func appendString(b []byte, s string) []byte {
	for i := range len(s) {
		c := s[i]
		if c < 0x20 || c >= 0x7f || c == '"' || c == '\\' || c == '<' || c == '>' || c == '&' {
			raw, err := json.Marshal(s)
			if err != nil {
				panic(err) // a string always marshals
			}
			return append(b, raw...)
		}
	}
	b = append(b, '"')
	b = append(b, s...)
	return append(b, '"')
}
