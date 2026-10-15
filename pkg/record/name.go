package record

import (
	"strconv"
	"strings"
	"time"
	"unicode"
)

// maxPart is the longest cleaned nickname, in runes: the longest
// nickname the server takes (protocol.MaxName).
const maxPart = 20

// maxFile is the longest cleaned file name without Ext, in runes; four
// seats of maxPart and the time and game ID fit.
const maxFile = 120

// FileName names a record: start time, game ID and the seats'
// nicknames, e.g. 2026-10-10-120000-g3_Ann_Bob.fsrec. Every part is
// cleaned, so nicknames cannot reach outside the records folder.
func FileName(started time.Time, game string, seats []string) string {
	parts := []string{started.UTC().Format("2006-01-02-150405") + "-" + SafeName(game)}
	for i, s := range seats {
		name := SafeName(s)
		if name == "" {
			name = "p" + strconv.Itoa(i+1)
		}
		parts = append(parts, name)
	}
	return strings.Join(parts, "_") + Ext
}

// SafeName keeps the letters and digits of s; every other run becomes
// one "-". It is trimmed of "-" and at most maxPart runes, so it is
// safe as part of a file name on every system; it may be "".
func SafeName(s string) string {
	return clean(s, maxPart, isAlnum)
}

// SafeFileName cleans a suggested record file name, such as one a game
// server sent: only its base, letters, digits, "-" and "_", and Ext.
// It returns "" when nothing is left.
func SafeFileName(name string) string {
	base := name[strings.LastIndexAny(name, `/\`)+1:]
	base = strings.TrimSuffix(base, Ext)
	base = clean(base, maxFile, func(r rune) bool { return isAlnum(r) || r == '_' })
	if strings.Trim(base, "_") == "" {
		return ""
	}
	return base + Ext
}

// clean keeps the runes keep accepts, turns every other run into one
// "-", trims "-", and cuts to most runes.
func clean(s string, most int, keep func(rune) bool) string {
	var b strings.Builder
	n, dash := 0, false
	for _, r := range s {
		if !keep(r) {
			dash = true
			continue
		}
		dash = dash && b.Len() > 0
		if dash && n+2 > most || n+1 > most {
			break
		}
		if dash {
			b.WriteByte('-')
			n++
		}
		dash = false
		b.WriteRune(r)
		n++
	}
	return b.String()
}

// isAlnum reports a letter or digit; marks, symbols, spaces, dots,
// slashes and control or format runes (bidi overrides) are not.
func isAlnum(r rune) bool { return unicode.IsLetter(r) || unicode.IsDigit(r) }
