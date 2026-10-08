// Package legal holds the fan-game notice shown by the app and the server (L-01, L-02).
package legal

import "strings"

// Part is a piece of a notice line; it is a link when URL is set.
type Part struct {
	Text string `json:"text"`
	URL  string `json:"url,omitempty"`
}

// Line is one line of the notice.
// A struct, not a slice type, so Wails can generate TypeScript for it.
type Line struct {
	Parts []Part `json:"parts"`
}

const (
	officialSite = "https://foursouls.com"
	officialShop = "https://maestromedia.com/collections/binding-of-isaac-four-souls"
)

// Notice is the fan-game notice.
// TODO(owner): replace the rights line with the exact line from the official rulebook.
var Notice = []Line{
	{Parts: []Part{{Text: "Four Souls is an unofficial, free, fan-made game."}}},
	{Parts: []Part{
		{Text: "The Binding of Isaac: Four Souls is designed by "},
		{Text: "Edmund McMillen", URL: officialSite},
		{Text: " and published by "},
		{Text: "Maestro Media", URL: officialShop},
		{Text: "."},
	}},
	{Parts: []Part{{Text: "All rights to the game, its rules, card texts and artwork belong to their owners."}}},
	{Parts: []Part{{Text: "This project is not affiliated with or endorsed by them."}}},
	{Parts: []Part{
		{Text: "Please support the official game: "},
		{Text: "buy it here", URL: officialShop},
		{Text: "."},
	}},
}

// PlainText renders the notice for a console, with links in parentheses.
func PlainText() string {
	var b strings.Builder
	for _, line := range Notice {
		for _, p := range line.Parts {
			b.WriteString(p.Text)
			if p.URL != "" {
				b.WriteString(" (" + p.URL + ")")
			}
		}
		b.WriteString("\n")
	}
	return b.String()
}
