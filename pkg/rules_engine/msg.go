package rulesengine

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"sync"
)

// Key names a line the player reads (ADR 021, item 9). The engine
// never holds the words: they are in the client's
// frontend/src/locales/en.json, under the same key. Declare every key
// once, at package level, with NewKey. A Key is no string, so a bare
// string can't stand in for one.
type Key struct{ name string }

// String is the key's name, e.g. "ask.yes".
func (k Key) String() string { return k.name }

// MarshalText writes the key as its name.
func (k Key) MarshalText() ([]byte, error) { return []byte(k.name), nil }

// UnmarshalText reads a key's name.
func (k *Key) UnmarshalText(b []byte) error {
	k.name = string(b)
	return nil
}

var (
	keysMu sync.Mutex
	keys   = map[Key]bool{}
	keyRe  = regexp.MustCompile(`^[a-z0-9_]+(\.[a-z0-9_]+)+$`)
)

// NewKey declares a key the engine, a card or the server sends, so the
// registry lists it (Keys). Declaring the same key twice is fine.
// Areas: refuse, prompt, event, ask, error; a card's own lines are
// card.<card id>.<name> (ADR 021, item 5).
func NewKey(name string) Key {
	if !keyRe.MatchString(name) {
		panic("rulesengine: bad message key " + strconv.Quote(name))
	}
	keysMu.Lock()
	defer keysMu.Unlock()
	keys[Key{name}] = true
	return Key{name}
}

// Keys lists every declared key, sorted. A Go test writes them to
// frontend/src/locales/server-keys.json.
func Keys() []Key {
	keysMu.Lock()
	defer keysMu.Unlock()
	out := make([]Key, 0, len(keys))
	for k := range keys {
		out = append(out, k)
	}
	slices.SortFunc(out, func(a, b Key) int { return strings.Compare(a.name, b.name) })
	return out
}

// Declared reports whether k was declared with NewKey.
func Declared(k Key) bool {
	keysMu.Lock()
	defer keysMu.Unlock()
	return keys[k]
}

// Msg is a line for the player: a key and the values of its
// placeholders (ADR 021, item 9). Each client shows it in its own
// language. The zero Msg is no line.
type Msg struct {
	Key  Key  `json:"key,omitzero"`
	Args Args `json:"args,omitempty"`
}

// Msg is the line k with these values.
func (k Key) Msg(args ...Arg) Msg { return Msg{Key: k, Args: args} }

// IsZero reports whether m is no line.
func (m Msg) IsZero() bool { return m.Key.name == "" }

// Arg is one value of a Msg: a number, a plain string, a card or a
// player (ADR 021, item 6). Build it with N, Str, CardArg or PlayerArg.
type Arg struct {
	Name  string
	Value any // int, string, CardValue or PlayerValue
}

// CardValue is a card in a message, shown by its name in the viewer's
// language: {"card": "the_d6"}.
type CardValue struct {
	Card CardRef `json:"card"`
}

// PlayerValue is a player in a message, shown by their name:
// {"player": 1}.
type PlayerValue struct {
	Player PlayerID `json:"player"`
}

// N is a number value; n also picks a plural form.
func N(name string, n int) Arg { return Arg{Name: name, Value: n} }

// Str is a plain string value, shown as it is.
func Str(name, s string) Arg { return Arg{Name: name, Value: s} }

// CardArg is a card value.
func CardArg(name string, ref CardRef) Arg { return Arg{Name: name, Value: CardValue{Card: ref}} }

// PlayerArg is a player value.
func PlayerArg(name string, p PlayerID) Arg { return Arg{Name: name, Value: PlayerValue{Player: p}} }

// Args are a message's values in order. They encode as one JSON object,
// in this order, so the saved state has no Go map (ADR 005).
type Args []Arg

// MarshalJSON writes the values as one object.
func (a Args) MarshalJSON() ([]byte, error) {
	var b bytes.Buffer
	b.WriteByte('{')
	for i, arg := range a {
		if i > 0 {
			b.WriteByte(',')
		}
		name, err := json.Marshal(arg.Name)
		if err != nil {
			return nil, fmt.Errorf("msg: %w", err)
		}
		value, err := json.Marshal(arg.Value)
		if err != nil {
			return nil, fmt.Errorf("msg: %w", err)
		}
		b.Write(name)
		b.WriteByte(':')
		b.Write(value)
	}
	b.WriteByte('}')
	return b.Bytes(), nil
}

// UnmarshalJSON reads the values back in their order.
func (a *Args) UnmarshalJSON(data []byte) error {
	dec := json.NewDecoder(bytes.NewReader(data))
	if _, err := dec.Token(); err != nil { // {
		return fmt.Errorf("msg: %w", err)
	}
	var out Args
	for dec.More() {
		tok, err := dec.Token()
		if err != nil {
			return fmt.Errorf("msg: %w", err)
		}
		name, ok := tok.(string)
		if !ok {
			return errArg
		}
		var raw json.RawMessage
		if err = dec.Decode(&raw); err != nil {
			return fmt.Errorf("msg: %w", err)
		}
		v, err := argValue(raw)
		if err != nil {
			return err
		}
		out = append(out, Arg{Name: name, Value: v})
	}
	*a = out
	return nil
}

var errArg = errors.New("msg: a value is a number, a string, a card or a player")

func argValue(raw json.RawMessage) (any, error) {
	switch {
	case len(raw) == 0:
		return nil, errArg
	case raw[0] == '"':
		var s string
		err := json.Unmarshal(raw, &s)
		return s, err //nolint:wrapcheck // the decoder's own error
	case raw[0] == '{':
		var probe map[string]json.RawMessage
		if err := json.Unmarshal(raw, &probe); err != nil {
			return nil, fmt.Errorf("msg: %w", err)
		}
		if _, ok := probe["card"]; ok {
			var c CardValue
			err := json.Unmarshal(raw, &c)
			return c, err //nolint:wrapcheck // the decoder's own error
		}
		if _, ok := probe["player"]; ok {
			var p PlayerValue
			err := json.Unmarshal(raw, &p)
			return p, err //nolint:wrapcheck // the decoder's own error
		}
		return nil, errArg
	default:
		var n int
		if err := json.Unmarshal(raw, &n); err != nil {
			return nil, errArg
		}
		return n, nil
	}
}

// String is the message in a plain, language-free form, for logs and
// test failures: the key, then each value.
func (m Msg) String() string {
	return plain(m, func(p PlayerID) string { return "player " + strconv.Itoa(int(p)+1) })
}

// Plain is how tests and situation payloads name an option of a choose
// prompt, with no language: a card option is its ref ("bomb"), a
// player option "player 2 (cain)", a number option its digits, a
// card's own words (a "choose one-" option) the words, an ability
// "card: words"; any other line
// is its key, then each value ("ask.slot_card 1 fly").
func (g *Game) Plain(m Msg) string {
	return plain(m, func(p PlayerID) string {
		if p < 0 || int(p) >= len(g.Players) {
			return "player " + strconv.Itoa(int(p)+1)
		}
		return "player " + strconv.Itoa(int(p)+1) + " (" + string(g.Object(g.Players[p].Character).Card) + ")"
	})
}

// PlainOptions are the options of the open prompt in their plain form.
func (g *Game) PlainOptions() []string {
	out := make([]string, len(g.Waiting.Options))
	for i, o := range g.Waiting.Options {
		out[i] = g.Plain(o)
	}
	return out
}

// OptionIndex is the first option of the open prompt whose plain form
// is label, or -1.
func (g *Game) OptionIndex(label string) int {
	for i, o := range g.Waiting.Options {
		if g.Plain(o) == label {
			return i
		}
	}
	return -1
}

func plain(m Msg, player func(PlayerID) string) string {
	value := func(v any) string {
		switch v := v.(type) {
		case int:
			return strconv.Itoa(v)
		case string:
			return v
		case CardValue:
			return string(v.Card)
		case PlayerValue:
			return player(v.Player)
		default:
			return fmt.Sprint(v)
		}
	}
	if len(m.Args) == 1 && (m.Key == AskCard || m.Key == AskPlayer || m.Key == AskNumber || m.Key == AskWords) {
		return value(m.Args[0].Value)
	}
	if len(m.Args) == 2 && m.Key == AskAbility {
		return value(m.Args[0].Value) + ": " + value(m.Args[1].Value)
	}
	parts := []string{m.Key.name}
	for _, a := range m.Args {
		parts = append(parts, value(a.Value))
	}
	return strings.Join(parts, " ")
}

// Shared option and question keys (ADR 021, item 5.5).
var (
	AskCard   = NewKey("ask.card")   // an option that is a card: {card}
	AskPlayer = NewKey("ask.player") // an option that is a player: {player}
	AskNumber = NewKey("ask.number") // an option that is a number: {n}
	AskYes    = NewKey("ask.yes")
	AskNo     = NewKey("ask.no")
	AskNone   = NewKey("ask.none")
	AskDone   = NewKey("ask.done")
	AskCancel = NewKey("ask.cancel")
	AskKeep   = NewKey("ask.keep")
	// AskSlot is a monster slot by number: {n}; AskSlotCard with its
	// top card: {n}, {card}.
	AskSlot     = NewKey("ask.slot")
	AskSlotCard = NewKey("ask.slot_card")
	// AskAbility is an ability by its card and words: {card}, {text}.
	AskAbility = NewKey("ask.ability")
	// AskWords is a card's own words, e.g. a "choose one-" option: {text}.
	AskWords = NewKey("ask.words")
	// The decks by kind, e.g. "treasure deck"; see DeckOption.
	askDecks = [deckCount]Key{
		NewKey("ask.deck.treasure"), NewKey("ask.deck.loot"),
		NewKey("ask.deck.monster"), NewKey("ask.deck.room"),
	}
	// The discard piles by kind, e.g. "loot discard"; see DiscardOption.
	askDiscards = [deckCount]Key{
		NewKey("ask.discard.treasure"), NewKey("ask.discard.loot"),
		NewKey("ask.discard.monster"), NewKey("ask.discard.room"),
	}
)

// CardOption is an option that is a card.
func CardOption(ref CardRef) Msg { return AskCard.Msg(CardArg("card", ref)) }

// PlayerOption is an option that is a player.
func PlayerOption(p PlayerID) Msg { return AskPlayer.Msg(PlayerArg("player", p)) }

// NumberOption is an option that is a number.
func NumberOption(n int) Msg { return AskNumber.Msg(N("n", n)) }

// SlotOption is a monster slot by number (from 1), with its top card if
// it has one.
func SlotOption(n int, top CardRef) Msg {
	if top == "" {
		return AskSlot.Msg(N("n", n))
	}
	return AskSlotCard.Msg(N("n", n), CardArg("card", top))
}

// DeckOption is a deck: "treasure deck".
func DeckOption(d DeckKind) Msg { return askDecks[d].Msg() }

// DiscardOption is a discard pile: "loot discard".
func DiscardOption(d DeckKind) Msg { return askDiscards[d].Msg() }

// Options turns keys into options with no values.
func Options(ks ...Key) []Msg {
	out := make([]Msg, len(ks))
	for i, k := range ks {
		out[i] = k.Msg()
	}
	return out
}
