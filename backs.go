package main

// backsDir holds the card back images inside the card images folder.
const backsDir = "assets/backs"

// backKinds are the decks with a back image, as the frontend's card kinds.
var backKinds = []string{"loot", "treasure", "monster"}

// CardBack is the back image of one deck (V-26).
type CardBack struct {
	Kind string `json:"kind"` // "loot", "treasure" or "monster"
	// Image is the image URL; "" when the file is missing, and the back
	// is drawn plain.
	Image string `json:"image"`
}

// Backs returns the back of every deck, for face-down cards.
func (a *App) Backs() []CardBack {
	out := make([]CardBack, 0, len(backKinds))
	for _, kind := range backKinds {
		out = append(out, CardBack{Kind: kind, Image: a.imageURL(backsDir + "/" + kind + ".webp")})
	}
	return out
}
