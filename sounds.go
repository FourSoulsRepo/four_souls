package main

import "strconv"

// soundsDir holds every sound inside the card images folder.
const soundsDir = "assets/sounds/"

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
	ChoiceDeal    []string `json:"choiceDeal"`    // a shown choice card arrives (20.18)
	ChoicePick    []string `json:"choicePick"`    // a choice card is picked (20.18)
	CardGrab      []string `json:"cardGrab"`      // a hand card is picked up (20.20)
	AimCancel     []string `json:"aimCancel"`     // an aim is cancelled (20.20)
	PlayRefused   []string `json:"playRefused"`   // a card cannot be played now (20.20)
	HistoryHover  []string `json:"historyHover"`  // a history tile is hovered (20.21)
	ChoiceUnpick  []string `json:"choiceUnpick"`  // a picked card is unpicked before Confirm (20.22)
	ChoiceConfirm []string `json:"choiceConfirm"` // a start-of-game pick is confirmed (20.22)
}

// soundFiles are the files of each sound, under the card images folder.
// All live only in the private images. The Fireball's cast and impact
// (15.50) are Hearthstone's, (c) Blizzard. The dice (15.69), loot draw
// (18.7), player death (18.8), treasure gain (18.9) and soul gain
// (18.13, 18.14) sounds come from the owner; their source is unknown.
// The rest are Hearthstone's too, from the reference client of step 20
// (20.17); the attack hits by size among them (20.3).
var soundFiles = Sounds{
	AttackHit:     soundVariants("attack-hit-1", 1),
	EffectCast:    soundVariants("effect-cast", 1),
	EffectImpact:  soundVariants("effect-impact", 1),
	DiceRoll:      soundVariants("dice-roll", 1),
	LootDraw:      soundVariants("loot-draw", 1),
	PlayerDied:    soundVariants("player-died", 1),
	TreasureGain:  soundVariants("treasure-gain", 1),
	SoulGain:      soundVariants("soul-gain", 1),
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
	ChoiceDeal:    soundVariants("choice-deal", 1),
	ChoicePick:    soundVariants("choice-pick", 1),
	CardGrab:      soundVariants("card-grab", 1),
	AimCancel:     soundVariants("aim-cancel", 1),
	PlayRefused:   soundVariants("play-refused", 1),
	HistoryHover:  soundVariants("history-hover", 1),
	ChoiceUnpick:  soundVariants("choice-unpick", 1),
	ChoiceConfirm: soundVariants("choice-confirm", 1),
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
		ChoiceDeal:    f(s.ChoiceDeal),
		ChoicePick:    f(s.ChoicePick),
		CardGrab:      f(s.CardGrab),
		AimCancel:     f(s.AimCancel),
		PlayRefused:   f(s.PlayRefused),
		HistoryHover:  f(s.HistoryHover),
		ChoiceUnpick:  f(s.ChoiceUnpick),
		ChoiceConfirm: f(s.ChoiceConfirm),
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
