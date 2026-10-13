package main

// The game's fonts inside the card images folder, downloaded by the tools
// repo (15.20): font-souls for Latin text, its Body face for card text
// (15.24), and the Russian site's card font, which also has Cyrillic.
const (
	titleFont    = "assets/fonts/font-souls-title.ttf"
	bodyFont     = "assets/fonts/font-souls-body.ttf"
	cyrillicFont = "assets/fonts/tboi4s-original.ttf"
)

// Fonts are the game's own font URLs; each is "" when its file is missing,
// and the frontend uses the system font.
type Fonts struct {
	Title    string `json:"title"`    // font-souls Title, Latin only
	Body     string `json:"body"`     // font-souls Body, card text, Latin only
	Cyrillic string `json:"cyrillic"` // the Russian card font, for Cyrillic
}

// Fonts returns the font URLs.
func (a *App) Fonts() Fonts {
	return Fonts{Title: a.imageURL(titleFont), Body: a.imageURL(bodyFont), Cyrillic: a.imageURL(cyrillicFont)}
}
