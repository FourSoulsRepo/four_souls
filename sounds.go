package main

// The fireball sounds inside the card images folder (15.50): Hearthstone's
// Fireball cast and impact, played on damage that is not an attack. They
// belong to Blizzard, so they live only in the private images; the card
// hit sound is CC0 and ships with the frontend.
const (
	fireballCastSound   = "assets/sounds/fireball-cast.mp3"
	fireballImpactSound = "assets/sounds/fireball-impact.mp3"
)

// Sounds are the private sound URLs; each is "" when its file is missing,
// and the frontend plays the card hit sound instead.
type Sounds struct {
	FireballCast   string `json:"fireballCast"`   // the shot leaves its source
	FireballImpact string `json:"fireballImpact"` // the shot lands
}

// Sounds returns the sound URLs.
func (a *App) Sounds() Sounds {
	return Sounds{FireballCast: a.imageURL(fireballCastSound), FireballImpact: a.imageURL(fireballImpactSound)}
}
