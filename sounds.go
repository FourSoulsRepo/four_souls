package main

// The damage sounds inside the card images folder (15.50), all from
// Hearthstone: an unarmed hit for attacks, and the Fireball spell's cast
// and impact for damage from effects. They belong to Blizzard, so they
// live only in the private images. So do the dice sound (15.69), the
// loot draw sound (18.7), the player death sound (18.8), the treasure
// gain sound (18.9) and the soul gain sound, for every soul (18.13,
// 18.14): their source is unknown.
const (
	attackHitSound    = "assets/sounds/attack-hit.mp3"
	effectCastSound   = "assets/sounds/effect-cast.mp3"
	effectImpactSound = "assets/sounds/effect-impact.mp3"
	diceRollSound     = "assets/sounds/dice-roll.mp3"
	lootDrawSound     = "assets/sounds/loot-draw.mp3"
	playerDiedSound   = "assets/sounds/player-died.mp3"
	treasureGainSound = "assets/sounds/treasure-gain.mp3"
	soulGainSound     = "assets/sounds/soul-gain.mp3"
)

// Sounds are the sound URLs; each is "" when its file is missing, and
// that sound stays silent.
type Sounds struct {
	AttackHit    string `json:"attackHit"`    // an attack lands, card to card
	EffectCast   string `json:"effectCast"`   // an effect's shot leaves its source
	EffectImpact string `json:"effectImpact"` // the shot lands
	DiceRoll     string `json:"diceRoll"`     // a thrown die lands on the table
	LootDraw     string `json:"lootDraw"`     // a looted card flies to the hand
	PlayerDied   string `json:"playerDied"`   // a player's death resolved
	TreasureGain string `json:"treasureGain"` // a gained treasure flies to the player
	SoulGain     string `json:"soulGain"`     // a soul is gained: flies to the player, or the Lost Soul turns to one
}

// Sounds returns the sound URLs.
func (a *App) Sounds() Sounds {
	return Sounds{
		AttackHit:    a.imageURL(attackHitSound),
		EffectCast:   a.imageURL(effectCastSound),
		EffectImpact: a.imageURL(effectImpactSound),
		DiceRoll:     a.imageURL(diceRollSound),
		LootDraw:     a.imageURL(lootDrawSound),
		PlayerDied:   a.imageURL(playerDiedSound),
		TreasureGain: a.imageURL(treasureGainSound),
		SoulGain:     a.imageURL(soulGainSound),
	}
}
