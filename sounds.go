package main

import "strconv"

// soundsDir holds every sound inside the card images folder.
const soundsDir = "assets/sounds/"

// The damage sounds (15.50), all from Hearthstone: the Fireball
// spell's cast and impact for damage from effects; an attack's hit
// is the reference client's (20.3). They belong to Blizzard, so they live only in the private
// images. So do the dice sound (15.69), the loot draw sound (18.7), the
// player death sound (18.8), the treasure gain sound (18.9) and the soul
// gain sound, for every soul (18.13, 18.14): their source is unknown.
// The rest are Hearthstone's too, taken from the reference client of
// step 20 (20.17); each later sub-step of step 20 plays its own.
const (
	attackHitSound    = soundsDir + "attack-hit-1.mp3"
	effectCastSound   = soundsDir + "effect-cast.mp3"
	effectImpactSound = soundsDir + "effect-impact.mp3"
	diceRollSound     = soundsDir + "dice-roll.mp3"
	lootDrawSound     = soundsDir + "loot-draw.mp3"
	playerDiedSound   = soundsDir + "player-died.mp3"
	treasureGainSound = soundsDir + "treasure-gain.mp3"
	soulGainSound     = soundsDir + "soul-gain.mp3"
)

// Sounds are the sound URLs, a list per sound. A sound with several
// files plays one of them at random (20.17). A file that is missing is
// left out; a sound with none stays silent.
type Sounds struct {
	AttackHit     []string `json:"attackHit"`     // an attack lands, card to card
	EffectCast    []string `json:"effectCast"`    // an effect's shot leaves its source
	EffectImpact  []string `json:"effectImpact"`  // the shot lands
	DiceRoll      []string `json:"diceRoll"`      // a thrown die lands on the table
	LootDraw      []string `json:"lootDraw"`      // a looted card flies to the hand
	PlayerDied    []string `json:"playerDied"`    // a player's death resolved
	TreasureGain  []string `json:"treasureGain"`  // a gained treasure flies to the player
	SoulGain      []string `json:"soulGain"`      // a soul is gained: flies to the player, or the Lost Soul turns to one
	AttackHit2    []string `json:"attackHit2"`    // a bigger hit (20.3)
	AttackHit3    []string `json:"attackHit3"`    // the biggest hit (20.3)
	LootLand      []string `json:"lootLand"`      // an own loot card lands in the hand (20.7)
	LootDrawOther []string `json:"lootDrawOther"` // another player draws a loot card (20.7)
	DragBed       []string `json:"dragBed"`       // a held card's soft loop (20.6)
	DragAir       []string `json:"dragAir"`       // a held card's whoosh loop, by speed (20.6)
	TurnStart     []string `json:"turnStart"`     // the own turn's banner (20.8)
	CardPlayOther []string `json:"cardPlayOther"` // another player's big card shows (20.9)
	EndTurnCue    []string `json:"endTurnCue"`    // nothing left to do this turn (20.10)
	TimerTick     []string `json:"timerTick"`     // a response timer second (20.15)
	TimerLast     []string `json:"timerLast"`     // the response timer's last second (20.15)
}

// soundFiles are the files of each sound, under the card images folder.
var soundFiles = Sounds{
	AttackHit:     []string{attackHitSound},
	EffectCast:    []string{effectCastSound},
	EffectImpact:  []string{effectImpactSound},
	DiceRoll:      []string{diceRollSound},
	LootDraw:      []string{lootDrawSound},
	PlayerDied:    []string{playerDiedSound},
	TreasureGain:  []string{treasureGainSound},
	SoulGain:      []string{soulGainSound},
	AttackHit2:    soundVariants("attack-hit-2", 1),
	AttackHit3:    soundVariants("attack-hit-3", 1),
	LootLand:      soundVariants("loot-land", 3),
	LootDrawOther: soundVariants("loot-draw-other", 3),
	DragBed:       soundVariants("drag-bed", 1),
	DragAir:       soundVariants("drag-air", 1),
	TurnStart:     soundVariants("turn-start", 1),
	CardPlayOther: soundVariants("card-play-other", 3),
	EndTurnCue:    soundVariants("end-turn-cue", 1),
	TimerTick:     soundVariants("timer-tick", 1),
	TimerLast:     soundVariants("timer-last", 1),
}

// soundVariants names a sound's n files: <name>.mp3 alone, or
// <name>-1.mp3 … <name>-n.mp3.
func soundVariants(name string, n int) []string {
	if n == 1 {
		return []string{soundsDir + name + ".mp3"}
	}
	files := make([]string, n)
	for i := range files {
		files[i] = soundsDir + name + "-" + strconv.Itoa(i+1) + ".mp3"
	}
	return files
}

// Sounds returns the sound URLs.
func (a *App) Sounds() Sounds {
	f, s := a.soundURLs, soundFiles
	return Sounds{
		AttackHit:     f(s.AttackHit),
		EffectCast:    f(s.EffectCast),
		EffectImpact:  f(s.EffectImpact),
		DiceRoll:      f(s.DiceRoll),
		LootDraw:      f(s.LootDraw),
		PlayerDied:    f(s.PlayerDied),
		TreasureGain:  f(s.TreasureGain),
		SoulGain:      f(s.SoulGain),
		AttackHit2:    f(s.AttackHit2),
		AttackHit3:    f(s.AttackHit3),
		LootLand:      f(s.LootLand),
		LootDrawOther: f(s.LootDrawOther),
		DragBed:       f(s.DragBed),
		DragAir:       f(s.DragAir),
		TurnStart:     f(s.TurnStart),
		CardPlayOther: f(s.CardPlayOther),
		EndTurnCue:    f(s.EndTurnCue),
		TimerTick:     f(s.TimerTick),
		TimerLast:     f(s.TimerLast),
	}
}

// soundURLs gives the URL of each file that exists; never nil, so the
// frontend gets an empty list, not null.
func (a *App) soundURLs(files []string) []string {
	urls := make([]string, 0, len(files))
	for _, file := range files {
		if u := a.imageURL(file); u != "" {
			urls = append(urls, u)
		}
	}
	return urls
}
