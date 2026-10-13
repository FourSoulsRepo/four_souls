package main

// The damage sounds inside the card images folder (15.50), all from
// Hearthstone: an unarmed hit for attacks, and the Fireball spell's cast
// and impact for damage from effects. They belong to Blizzard, so they
// live only in the private images.
const (
	attackHitSound    = "assets/sounds/attack-hit.mp3"
	effectCastSound   = "assets/sounds/effect-cast.mp3"
	effectImpactSound = "assets/sounds/effect-impact.mp3"
)

// Sounds are the sound URLs; each is "" when its file is missing, and
// that sound stays silent.
type Sounds struct {
	AttackHit    string `json:"attackHit"`    // an attack lands, card to card
	EffectCast   string `json:"effectCast"`   // an effect's shot leaves its source
	EffectImpact string `json:"effectImpact"` // the shot lands
}

// Sounds returns the sound URLs.
func (a *App) Sounds() Sounds {
	return Sounds{
		AttackHit:    a.imageURL(attackHitSound),
		EffectCast:   a.imageURL(effectCastSound),
		EffectImpact: a.imageURL(effectImpactSound),
	}
}
