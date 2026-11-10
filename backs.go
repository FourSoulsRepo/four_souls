package main

// backsDir holds the card back images inside the card images folder.
const backsDir = "assets/backs"

// backKinds are the decks with a back image, as the frontend's card kinds.
// The room back is landscape, as the rooms are (14.13).
var backKinds = []string{"loot", "treasure", "monster", "room"}

// CardBack is the back image of one deck (V-26).
type CardBack struct {
	Kind string `json:"kind"` // "loot", "treasure", "monster" or "room"
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
