# 27. Game feel: lessons from a reference client

* **Status:** Accepted
* **Date:** 2026-11-06
* **Authors:** @HardDie

---

## Context

1. Players say the table works but feels flat (V-54).
2. Hearthstone is the feel we aim for (V-13, V-22).
3. Videos and screenshots show the look, not the timing.
4. The owner has Fanstone installed, a Hearthstone-based fan client.
   1. Unity with Mono: its C# code decompiles to readable source.
   2. Its timings and rules are plain constants in that code.
5. We study it only for numbers and ideas.
   1. No code, images or names are copied.
   2. The decompiled source stays in `FourSoulsRepo/tools/.cache/`.
   3. Tool: `reverse/decompile.sh` in the tools repo (ILSpy).
   4. Sounds are the exception (owner, 2026-11-06); see option 4.
6. Times below are in seconds at normal speed.
   1. Our animation speed setting scales them, as today.
   2. "easeOutExpo" etc. are the standard easing curves.

## Considered options

1. **Decompile and measure the reference client**
   1. Won. Exact numbers and the rules behind them.
   2. ILSpy fits Mono builds; Ghidra is for native code.
2. **Tune by eye from videos**
   1. Lost. Guesses; 0.12 s and 0.3 s look alike in a video.
3. **Port its code or assets**
   1. Lost. Not ours to ship; the repo is public and MIT.
4. **Take its sounds into the private images**
   1. Won for sounds only (owner, 2026-11-06).
   2. They are Hearthstone's, (c) Blizzard, like the three of 15.50.
   3. They ship like the card images: private submodule, never public.
   4. Extracted by a tools-repo script from its Unity bundles (20.17).

## Decision

Use option 1, and option 4 for sounds.
The findings below are the target feel for step 20.
Each sub-step of step 20 takes the part it needs.

### Principles

1. Feedback never blocks the game.
   1. The damage number reports done on its first frame.
   2. It stays on screen 3 s while the next events play.
2. The pointer gets an instant answer; the game gets a paced one.
   1. A hovered hand card snaps to its big pose, no tween.
   2. A hovered table card shows large only after 0.4 s.
3. Fast strike, slow settle.
   1. The attack itself takes 0.12 s, linear.
   2. Cards then ease back over 0.5 – 1 s, easeOutExpo.
4. Size of the effect follows size of the event.
   1. The impact effect is picked by the attacker's power.
   2. Shakes come in three strengths: 0.1, 0.5, 1.
5. You never see your own action replayed at you.
   1. Big card only for other players' plays.
   2. Spectators see every play, since they cannot act.
6. Repeats speed up.
   1. Draws 1–3 at full time; then each one 11 % faster.
   2. Floor: 0.556 of full time.
7. Important moments are never skipped.
   1. A lagging client drops minor animations (minion fights).
   2. Fights with a hero always play.
8. One sound per moment.
   1. Two cards die together: one death sound.
   2. A voice line ducks the rest to 20 % (0.7 s in, 0.7 s out).
   3. The same clip has a playback limit.
9. Pressure shows only when it matters.
   1. The turn timer "rope" shows only in the last 20 s.
   2. Each second then ticks; the last tick sounds different.

### Hand

1. Fan only from 4 cards; 1–3 cards lie straight.
2. Fan angle: 40° plus 2° per card, spread evenly.
3. Ends drop on a parabola: `(i + 0.5 − n/2)² / (6n)`.
4. Spacing shrinks to fit a max width; cards never leave it.
5. Hover:
   1. The card snaps upright, 1.5× larger, raised; no tween.
   2. Its neighbours do not move apart.
   3. Keyword tips appear on the side toward the screen middle.
6. Re-layout: scale and turn 0.25 s, move 0.5 s, easeOutExpo.
7. Opponent hand re-layout: move 1.5 s, slower and calmer.
8. Checked in 20.5 (`ZoneHand`, `InputManager`, `BigCard`).
   1. Fan from 4 cards: `IsHandScrunched` is `cardCount > 3`.
   2. Own hand: 40° + 2° per card; an opponent's: 40° flat.
   3. Each card turns about its own centre; x is evenly spaced.
   4. The drop adds `sin(angle) × spacing / 2` on the right side only.
   5. A hovered hand card is set in place at once; others keep their tween.
   6. A board card's big card waits `m_MouseOverDelay`, 0.4 s.
   7. A hand card's shows at once, with no delay.
   8. The big card scales in over `m_ScaleSec`, 0.15 s.
   9. A new hovered card restarts the reference's 0.4 s timer.
9. Built in 20.5.
   1. Own hand only; other players' hidden hands keep their fan.
   2. Fan: `fanPose` in `feel.ts`; each card turns about its centre.
   3. Drop unit: 128 px, tuned by eye; no right-side extra.
   4. Hand hover: the zoom rises from the card at once, as before.
   5. The hovered card keeps its place; no neighbour moves.
   6. Table hover: zoom after 0.4 s, grows from the card in 0.15 s.
   7. Its curve is expo-out, as iTween's default.
   8. Ours: while a zoom shows, the next card's shows at once.
   9. Also within 0.3 s after one closed (`feel.zoom.grace`).
   10. The step file asks for that; the reference waits again.
   11. Speed "off" and reduced motion: no growth; the delay stays.
   12. Touch keeps its long press; history names zoom at once.

### Draw

1. The card rises from the deck, then flies to a "show" spot.
2. Path: deck → just above the deck → show spot, 1.5 s.
   1. Easing: easeInSine then easeOutExpo.
   2. Turn face up over 1.35 s, after 0.15 s.
   3. Scale up over 0.75 s, after 0.15 s.
3. Then it joins the hand with the normal re-layout.
4. Sounds: one when it leaves the deck, one when it lands.
5. It becomes playable 0.3 s after landing.
6. Draws queue: one card in flight at a time per player.
7. Opponent draws: back only, 0.2 s + 0.6 s waits, no show spot.
8. The path's easeInSineOutExpo runs as easeOutExpo (checked in 20.1).
   1. Its switch compares the time with the start value.
   2. A path eases from 0, so the sine-in half never plays.
   3. Our `sineInExpoOut` in `feel.ts` is therefore expo-out.
9. Checked in 20.7 (`Card`, `ZoneHand`, `GameState`, `PowerProcessor`).
   1. The path's middle point is just above the deck.
   2. Turn and scale have no ease set: iTween's default, expo-out.
   3. At the show spot the card waits for its tweens, then joins the hand.
   4. That move is the hand's "slow" re-layout: 0.5 s, expo-in.
   5. `add_card_to_hand` plays as that move starts; then 0.3 s to playable.
   6. Pace: the draw counter, from 0; under 3 is 1, then −0.111 each.
   7. Counter 6 and on: 0.556.
   8. Own and opponent draws count apart.
   9. Both reset when a new action's task list starts.
   10. Opponent: 1.75 s path, quart-in-out; the game goes on after 0.8 s.
10. Built in 20.7.
   1. Own: show spot 1 s, expo-out; hand 0.5 s, expo-in: 1.5 s in all.
   2. The step asks for 1.5 s; the reference takes 2 s, so the first leg is 1 s.
   3. Turn and growth: after 0.1 s, over 0.9 s and 0.5 s, expo-out.
   4. The path rises toward a point 80 px above the deck: ours, by eye.
   5. Show spot: 50 % across, 45 % down the window, 200 px wide.
   6. `loot-land` plays on arrival in the hand, not as the move starts.
   7. Its hand card takes no pointer for 0.3 s once the view shows it.
   8. Pace counts each player's draws in one update (`drawPace`).
   9. Another player's: 0.8 s, quart-in-out, back only, at the deck's size.
   10. 0.8 s is the reference's wait before the game goes on.
   11. A spectator sees every draw as another player's.
   12. Reduced motion: an own card shows at the show spot and fades.
   13. Reduced motion, another player's: only its sound.
   14. Speed "off": no draw, as every effect.
   15. `sineInExpoOut` stays expo-out; no sine-in phase is used.
   16. From the loot discard another player's card shows its face (owner, 20.16).
   17. That card is public; it flies face up, never turns.

### Drag

1. The held card tilts by its own motion.
   1. Pitch from vertical motion, roll from horizontal.
   2. Gain: 25° per table unit moved per frame.
   3. Limits: pitch ±40°, roll ±45°.
   4. It springs back to flat: smooth-damp, 0.2 s.
2. Two looping sounds follow the drag.
   1. A soft "magic" bed fades in to 15 % over 0.5 s.
   2. An "air" whoosh: volume rises with speed, up to 50 %.
   3. Both fade out in 0.2 s on drop.
3. While held, the cost shows on the mana crystals.
4. Checked in 20.6 (`DragRotator`, `InputManager`, `DragCardSoundEffects`).
   1. The gain is degrees per world unit moved in one frame.
   2. Pitch from forward motion, roll from sideways.
   3. The edge the card moves toward leans away.
   4. Rest 2 s times a 0.1 fudge: a smooth-damp of 0.2 s.
   5. Collection and pack drags use gains 3 and 4.5; a held game card 25.
   6. Bed: smooth-damp up to 0.15 over 0.5 s.
   7. Air target: ln(0.5 × units moved per frame + 0.92).
   8. Air eases by smooth-damp 0.04 s, at most 1 per second.
   9. Air plays clamped to 0 – 0.5; when still the log is below 0.
   10. Drop: both smooth-damp to 0 over 0.2 s, then stop.
5. Built in 20.6.
   1. Tilt per frame: `tiltStep` in `drag.ts`; numbers in `feel.drag`.
   2. One table unit is 60 px: ours, tuned by eye.
   3. Speed counts per 60 fps frame; any frame rate plays alike.
   4. The card's inner box turns, perspective 600 px.
   5. Its green aura turns with it.
   6. Speed "off" and reduced motion: the card stays flat.
   7. Loops: Web Audio, a gain node each, then a master volume gain.
   8. The bed fades in linearly over 0.5 s; both fade out linearly in 0.2 s.
   9. The loops still play with speed "off": sound, not motion.
   10. WebKit keeps the MP3 encoder delay: 13 ms before, 14–28 ms after.
   11. `loopBounds` skips it; Chromium trims it itself.

### Attack

1. Windup: the attacker lifts (its own short animation, ≤ 1 s).
2. Strike: the attacker flies into the target in 0.12 s, linear.
   1. The impact point is the target's own place (checked in 20.19).
   2. The target stays; only the attacker moves.
3. On contact:
   1. Damage applies at that exact moment.
   2. The impact effect plays; bigger for more power.
   3. Nearby cards shake away from the impact point.
4. Return: the normal zone re-layout, 1 s, easeOutExpo.
5. A hero turns to face the target in 0.3 s; back in 0.3 s.
6. A declared but not yet resolved attack moves half way.
7. Impact point, checked in 20.2 (`AttackSpellController`).
   1. The target "moves" to its own place; the attacker stops short (20.19).
   2. A proposed attack's point is half way to the target.
   3. The impact effect plays at the target's place at contact.
   4. Contact is the attacker's move ending, not a timer.
8. Built in 20.2.
   1. Windup 0.25 s, ease-out: lift to 1.1×, draw back 6 %.
   2. Strike 0.12 s, linear; only the attacker moves (20.19).
   3. Settle 0.5 s, expo-out: the low end of principle 3.
   4. Our burst carries "−N", so it pops on the target's home.
   5. The target's new HP shows on the hit; it never left (20.19).
   6. Reduced motion: no strike; the hit lands in place.
9. Shakes, checked in 20.3 (`MinionShake`, `CameraShakeMgr`).
   1. Strengths 0.1, 0.5 and 1 scale the push away from the impact.
   2. The push itself is an animator clip; its size is not in code.
   3. A shake may have a radius; cards past it stay still.
   4. A screen shake jitters at random, fading linearly to 0.
   5. Its usual size: 0.1–0.15 units, over 0.4–0.9 s.
   6. Players can turn screen shake off in the reference's options.
   7. Its impact effect is picked by attacker power, not damage.
10. Hits by size, built in 20.3.
    1. Picked by damage: size 1 for 1, 2 for 2–3, 3 for 4 and more.
    2. Damage is what Four Souls shows; attacks are mostly 1.
    3. Burst scale 1, 1.25, 1.5; sizes 2 and 3 add a gold glow.
    4. Nearby cards: pushed 12 px × 0.1, 0.5 or 1, then settle; 0.35 s.
    5. Nearby: centre within 2.5 target sizes of the target's centre.
    6. Lethal hit or player death: the table jitters 6 px, 0.4 s.
    7. Lethal: the target's HP is 0 in the update, or it died in it.
    8. Sound by size: `attack-hit`, `attack-hit-2`, `attack-hit-3`.
    9. Other damage keeps its fireball sound; size sets burst and shake.
    10. Speed "off" and reduced motion: no shake.
11. Impact point, checked in 20.19 (code and prefab data).
    1. Impact point: source + staging × (target − source).
    2. Staging is 1 in both prefabs: `AttackSpellController` and its Battlegrounds minion one.
    3. A proposed attack halves it: the only half-way case.
    4. The target tweens to the impact point: its own place, so it stays.
    5. Source offset −0.25 in both prefabs.
    6. The attacker's centre stops 0.5 − (−0.25) = 0.75 of its length short.
    7. Its length is along the way, from its own bounds.
    8. Cards of one size then overlap by a quarter of that length.
    9. Strike tween 0.12 s, linear, hero and minion alike.
    10. A hero's move back: 0.88 s, expo-out.
12. Only the attacker flies, built in 20.19.
    1. The target card stays home and seen; no copy of it moves.
    2. The attacker's centre stops 0.75 of its length short (`feel.lunge.short`).
    3. A player's panel as target: unchanged, 0.35 of its short side.
    4. Windup, contact, burst, sound, shakes and the number: unchanged.
    5. Cards near it still shake (20.3); the target itself never did.

### Damage and death

1. Damage number: grows 0 → 1 over 1 s, easeOutElastic.
2. It holds 2 s, then fades out over 1 s.
3. Deaths in one step play together.
4. A dying card's death effect waits 0.3 s for keyword deaths.
5. Checked in 20.4 (`DamageSplatSpell`, `Card.OnMetaData`).
   1. The 2 s hold counts from the splat's start, not after its growth.
   2. So it shows 3 s in all: grow 1 s, full until 2 s, fade 1 s.
   3. It reports done one frame after it starts; the game goes on.
   4. A card has one splat: a new hit restarts it with the new number.
6. Built in 20.4.
   1. The burst "−N" grows 0 → 1 in 1 s, elastic-out (`elasticOut`).
   2. Its fade starts at 2 s and ends at 3 s (`feel.burst`).
   3. Playback goes on at once: the hit's effect ends as it starts.
   4. Bursts on one card stack: middle, above, below, two above, …
   5. One burst height apart, so none covers another.
   6. Stacking is ours: we keep each number, not only the last.
   7. It stays on the target's home, as 20.2 put it.
   8. Reduced motion: no growth; full size at once, then the fade.
   9. Speed "off": no burst, as every effect.

### Others' plays

1. Another player's card shows large in the middle.
2. Its effect starts 1 s after the big card shows.
3. The big card stays 3 s (minion) or 4 s (spell).
4. A transformed card halves that time.
5. A spectator sees the chosen option highlighted for 1.75 s.
6. Checked in 20.9 (`HistoryManager`, `HistoryCard`).
   1. A spell's big card moves from the hand card's place to the middle.
   2. That path takes 1 s; it grows from 60 % to full over it.
   3. No ease set: iTween's default, expo-out.
   4. `play_card_from_hand` plays as the move starts.
   5. The game goes on 1 s after the move starts.
   6. Display time counts from the move's start: 3 s minion, 4 s spell.
   7. At its end the big card is removed at once; no exit motion.
   8. A newer big card removes the old one at once.
7. Built in 20.9.
   1. Times in `feel.played`: rise 1 s, expo-out; wait 1 s; stay 3 s.
   2. Stay is 3 s, as the step file asks; the reference's spells take 4 s.
   3. It grows from the hand card's size, not from 60 %.
   4. Ours: at 3 s it shrinks into the stack panel in 0.4 s, as 15.74 did.
   5. The label "Ann plays" fades as the leave starts.
   6. The card outlives its effect: it is the store's `big`, not `playing`.
   7. A newer big card sends it away at once, as the reference.
   8. Ours: so does a later middle effect or a move of the card itself.
   9. Those are reveal, dice, banner, turn banner and an own draw.
   10. `card-play-other` plays as it shows; never for an own play.
   11. A viewer without a seat (you −1) sees every play.
   12. Spectators and replays are steps 9 and 10; they get it from that rule.
   13. Reduced motion: it fades in and out in the middle.
   14. Speed "off": no big card, as every effect.

### History

1. Tiles in a side strip; each new tile pushes the rest, 1 s.
2. Hover a tile: the source card large.
   1. Its targets in a row beside it, scaled to fit.
   2. Each target shows its damage or death.
   3. A soft sound plays on hover.
3. Checked in 20.21 (`HistoryManager`, `HistoryCard`, `HistoryItem`).
   1. Also `PowerProcessor`, `InputManager`, `CardTileListDisplay`.
   2. Prefab data via UnityPy: the tiles, the strip, its `SoundDef`.
   3. Item 2.1 holds for up to 4 targets; 8 below covers more.
4. What makes a tile.
   1. A played card or hero power: one tile; its target goes inside it.
   2. An attack: one tile; the defender goes inside it.
   3. A trigger the server marks for the history; a secret.
   4. Fatigue and burned cards: a tile each.
   5. Damage and deaths make no tile of their own.
   6. They mark the cards inside the tile that caused them.
   7. A tile comes once its whole block has played.
   8. No turn marks, no full log, no button.
5. Size and count.
   1. Play tiles are full height; attack and trigger tiles half.
   2. Tiles stack down from the top, gaps of a quarter tile.
   3. The strip holds about 6 full or 11 half tiles.
   4. When one no longer fits, the oldest goes at once.
   5. No scrolling: older plays are gone.
6. A new tile.
   1. It starts far left, off the board, level with the top slot.
   2. It slides into the top slot in 1 s (iTween's default, expo-out).
   3. The others move down a slot in the same 1 s.
   4. No fade and no sound.
   5. While the pointer is on the strip, no tile moves.
   6. New tiles wait off the board till the pointer leaves.
7. What a tile shows.
   1. The card's art, cropped; a frame around it.
   2. The frame is tinted by owner: blue for own, red for the opponent's.
   3. Attack tiles carry a sword; trigger tiles a trigger icon.
   4. No damage, death or count marks on the tile itself.
8. Hover.
   1. Instant: no delay, no grow.
   2. The strip is one hit area; the nearest tile edge picks the tile.
   3. So the gaps between tiles never close the card.
   4. The source card shows large right of the strip, by the tile.
   5. A chosen mode's card has its owner card left of it.
   6. Then a separator: swords for an attack, else an arrow.
   7. Targets right of it: 4 in one row, 8 in two, then three rows.
   8. They shrink to fit two card widths: at most 1, 0.5 or 0.3.
   9. Each shows its damage splat, a skull if it died, or a flame.
   10. The source card shows its own damage or death the same way.
   11. A mark shows at once, stays 1.5 s, then fades out in 0.5 s.
   12. Keyword tips stand beside the large card.
   13. The rest of the screen dims and greys; the tiles stay bright.
   14. In 0.5 s, ease-in-out cubic; out in 0.1 s.
   15. Each newly hovered tile plays `history_event_mouseover`.
   16. Clip `history_event_mouseover.wav`, volume 0.6; never ducked.
   17. While on the strip, all other sound ducks to 20 %.
   18. In 0.3 s, out in 1 s, linear.
   19. No hover while the own aim arrow is out.
9. Click and touch.
   1. A click on a tile does nothing.
   2. Touch: a finger held on the strip shows the tile under it.
   3. Sliding along switches tiles; lifting the finger hides it.
10. Ours, at 20.21.
    1. A tile per resolved stack item, buy and kill (UC-15 8).
    2. So a Bomb! kill is three tiles: the play, the damage, the kill.
    3. An attack is a tile per roll and per damage, then the kill.
    4. A tile comes after its update's effects, as the reference's.
    5. One size: the stack's icon (card art, die, blood drop, skull).
    6. Badges: sword (attack), bolt (loot), coin (buy), skull (kill).
    7. A damage tile shows its number; no owner tint on any tile.
    8. A line under the first tile of each turn.
    9. About 7 fit in the test window; the rest scroll.
    10. No scroll bar and no sign that more is below.
    11. Nothing drops off; the History button opens the full log.
    12. A new tile fades in in 0.2 s; the others slide in 0.25 s.
    13. That is `row-move`, shared with the stack panel.
    14. Tiles move under a resting pointer too.
    15. Hover: instant; the card large beside the tile.
    16. A caption under it: the ability or line, then the owner in color.
    17. What it did: an arrow, then cards in two columns, hand size.
    18. Damage: a starburst with the number; death: dimmed, a skull.
    19. Marks stay while hovered.
    20. A play's popout shows its target plain, without later marks.
    21. An attack roll's popout: the monster large; no attacker.
    22. Each tile is its own hit area; the 3 px gaps close the card.
    23. No sound, no dimming, no ducking.
    24. Touch: a tap shows the card till the next tap.
11. Differences and recommendations.
    1. What makes a tile.
       1. Ours: one per stack item. Ref: one per play or attack.
       2. Recommend: keep ours; in Four Souls each item can be answered.
    2. A play's result.
       1. Ours: Bomb! points at a plain Peep. Ref: Peep with its splat, skull.
       2. Recommend: a play's popout adds the marks its own items left.
       3. Damage it dealt and a death it caused, on its target.
    3. Tile size.
       1. Ours: one size. Ref: plays full, attacks and triggers half.
       2. Recommend: rolls and damage half height; the rest full.
       3. More fit; plays, kills and buys stand out.
    4. Owner tint.
       1. Ours: none. Ref: blue frame for own, red for others'.
       2. Recommend: a thin frame in the owner's player color.
    5. Full strip.
       1. Ours: scrolls, with no sign. Ref: the oldest drops off.
       2. Recommend: keep the scroll; fade the bottom edge when more is below.
    6. A new tile.
       1. Ours: fades in 0.2 s. Ref: slides in from the left, 1 s, expo-out.
       2. Recommend: slide in from the board's left edge, expo-out.
       3. The rest move down in the same time; 1 s, tuned by eye.
       4. Speed scales it; "off" places it at once.
    7. Still while hovered.
       1. Ours: tiles move under the pointer. Ref: none move.
       2. Recommend: hold the strip still while the pointer is on it.
       3. New tiles come in when it leaves.
    8. Hit area.
       1. Ours: per tile; gaps close the card. Ref: the nearest tile.
       2. Recommend: the strip as one hit area, nearest tile wins.
    9. Hover sound.
       1. Ours: none. Ref: `history_event_mouseover.wav`.
       2. Recommend: that clip on each newly hovered tile (rule 8).
    10. Ducking on hover.
        1. Ours: none. Ref: all else to 20 % while on the strip.
        2. Recommend: skip; a turn or timer cue must stay heard.
    11. Dimming on hover.
        1. Ours: none. Ref: the screen dims and greys, 0.5 s in.
        2. Recommend: try a light dim of the board, 0.1 s out.
        3. Skip it if it flashes when sliding along the strip.
    12. Separator.
        1. Ours: always an arrow. Ref: swords for an attack.
        2. Recommend: swords for an attack's roll, damage and kill.
    13. Attack roll popout.
        1. Ours: the monster large. Ref: the attacker large, then the defender.
        2. Recommend: the attacker's character large, swords, the monster.
        3. The die and its number stay in the caption.
    14. Target layout.
        1. Ours: two columns, at most 220 px wide.
        2. Ref: one row up to 4, then rows, shrunk to fit.
        3. Recommend: one row up to 4, then two rows; hand size at most.
    15. Mark fade.
        1. Ours: marks stay. Ref: gone after 2 s.
        2. Recommend: keep ours; a player reads at their own pace.
    16. Turn line, badges, History button.
        1. Ours: all three. Ref: none.
        2. Recommend: keep ours; Four Souls turns hold many plays.
    17. Touch.
        1. Ours: tap, till the next tap. Ref: hold and slide.
        2. Recommend: keep ours.
    18. Hover while aiming.
        1. Ours: allowed. Ref: off while the arrow is out.
        2. Recommend: keep ours; a look back can help the aim.
12. Built in 20.21 (owner's picks: all; 1, 15–18 kept as ours).
    1. Marks (11.2): `linkMarks` in the store, as entries come in.
    2. A damage marks the play with its item's source, else its card.
    3. Damages add up; a kill turns the mark into a skull.
    4. A target the play did not name joins its list.
    5. An attack's damage marks the attack roll at the same monster.
    6. It looks back within the turn, 200 entries at most.
    7. Size (11.3): rolls and damage half height, by the item's kind.
    8. Frame (11.4): 2 px inside the tile, the owner's player color.
    9. Owner: the item's controller, the buyer, the killer.
    10. Full strip (11.5): the foot fades over 1.2 tile sides.
    11. Only while more is below; a scroll to the end clears it.
    12. New tile (11.6): 1 s, expo-out, from 1.3 tile widths left.
    13. It is clipped at the strip's edge, so it comes out of the board's edge.
    14. The list moves down as one element in the same time.
    15. Speed scales it ("fast" 0.5 s); "off", reduced motion: at once.
    16. Held (11.7): the tiles stay as the pointer found them.
    17. A rollback still greys a held tile in place.
    18. Hit area (11.8): the tile whose edge is nearest the pointer.
    19. Sound (11.9): `historyHover` once per newly hovered tile.
    20. Dim (11.11): the board to 30 % black, 0.5 s cubic in-out; out 0.1 s.
    21. Kept: it did not flash sliding along the strip.
    22. The dim follows the whole hit area, not a tile; gaps keep it.
    23. Speed "off" and reduced motion: it shows and clears at once.
    24. Separator (11.12): swords for an attack's roll, damage and kill.
    25. Attack roll (11.13): the attacker's character, swords, the monster.
    26. The monster carries the attack's damage and death too (11.2).
    27. Targets (11.14): one row up to 4, then two rows.
    28. Cards are hand size, shrunk past four a row.
    29. Speed: each tile draws again only when its entry changes (19.138).

### Turn start

1. A strict order; input stays locked until it ends.
2. "Your turn" banner and its sound.
   1. Fade in 0.25 s while it grows to full size.
   2. Then a small "punch": down to 98 % in 0.15 s.
   3. It stays 1.5 s, then shrinks and fades in 0.25 s.
3. Mana fills and exhausted cards ready, during the banner.
4. Wait 1 s for the banner, then the turn's draw.
5. Then start-of-turn effects.
6. Only then: input unlocks, End Turn becomes active.
7. Checked in 20.8 (`TurnStartIndicator`, `TurnStartManager`).
   1. The banner grows from 10 % to full in 0.25 s while it fades in.
   2. Then the punch: down to 98 % in 0.15 s.
   3. Its 1.5 s counts from its start, not from the punch.
   4. Then it shrinks to almost nothing and fades in 0.25 s: 1.75 s in all.
   5. No ease set: iTween's default, expo-out.
   6. The sound plays as the banner shows; own turn only.
   7. Mana and readied cards change as the banner shows, not later.
   8. The draw waits 1 s from the banner's start.
   9. Then 0.35 s more if effects follow the draw.
   10. An opponent's turn shows no banner in the reference.
8. Built in 20.8.
   1. Banner: `feel.turn`; 10 % → full in 0.25 s, punch to 98 % in 0.15 s.
   2. It starts to fade at 1.5 s and is gone at 1.75 s, as the reference.
   3. The step file's "stays 1.5 s" read as the reference's 1.5 s from start.
   4. Playback goes on after 1 s; the banner stays above the next effects.
   5. It lives in the effect layer (`.effects-trail`); no new layer.
   6. The recharge lights the player's spent cards as the banner shows.
   7. Ours dim a spent card instead of turning it; it lights, not turns.
   8. Lit cards: charged in the update's view, kept until the view shows.
   9. Every player's turn has a banner; only the own plays `turn-start`.
   10. Order: our rules put start-of-turn triggers before the loot step.
   11. The table follows the events: banner, start-of-turn, then the draw.
   12. No input lock: the next view shows at 1 s; its buttons are live.
   13. Reduced motion: it only fades in and out.
   14. Speed "off": no banner, as every effect.
   15. No 0.35 s wait after the draw; playback paces each effect itself.

### End Turn button

1. States: waiting, your turn, nothing left to do.
2. Nothing left: the button glows and changes color.
3. A short cue plays 1.5 s later, once per turn.
4. The cue is skipped if a play becomes possible meanwhile.
5. Checked in 20.10 (`EndTurnButton`).
   1. "No more plays": every option but End turn is unplayable.
   2. The cue is tried once per turn: a flag set as the state first comes.
   3. The flag resets when the own turn starts.
   4. After 1.5 s the cue plays only if the state still holds.
   5. A skipped cue is not tried again that turn.
   6. Green highlight intensity: 1.1, and 1.4 on hover; 0.15 s, linear.
   7. With plays left, a hover shows a white highlight instead.
   8. Ending the turn sets the button to waiting at once.
   9. Press and release clicks exist; not taken (no clip traced in 20.17).
6. Built in 20.10.
   1. Looks: `endTurnLook` in `actions.ts`: `waiting`, `turn`, `done`.
   2. `done`: the update allows only End turn and Pass, and the viewer may act.
   3. The button turns green; its halo is larger and never below 55 %.
   4. The halo still breathes over 1.6 s each way, as in 19.76.
   5. Hover brightens the halo to 130 % in 0.15 s, linear: ours, by eye.
   6. The cue: `endTurnCue.ts`, `feel.endTurn.cue`, 1.5 s, once per turn.
   7. The 1.5 s is a wait, not a motion: the speed setting does not scale it.
   8. The cue plays with speed "off" too, as the drag loops.
   9. Speed "off" and reduced motion: the halo stands still; no hover tween.

### Counters

1. Crystals are added, readied and spent one at a time.
2. 0.2 s between each crystal; the number flips in 0.1 s.
3. Checked in 20.11 (`ManaCrystalMgr`, `ManaCounter`).
   1. Spawn, ready and spend each wait 0.2 s per crystal queued.
   2. A spend's first crystal goes at once; spawn and ready wait 0.2 s first.
   3. No speed-up for many crystals; ours is the step file's.
   4. The 0.1 s is a text fade as the phone's gem flips; not per step.
   5. `ManaCounter` sets its number in one go; it has no count of its own.
   6. A held card marks its cost's crystals "proposed" (`ProposeManaCrystalUsage`).
   7. Cancelled, they turn back to ready.
4. Built in 20.11.
   1. `Count` and `useCount` in `Count.tsx`; times in `feel.count`.
   2. Cents of every player, soul crystals, health badges.
   3. N steps, 0.2 s apart, the first at once, as a spend.
   4. More than 5: the gap is 1 s / N, so the count fits 1 s (`countTimes`).
   5. Each step flips the number in, 90° about its middle, 0.1 s, expo-out: ours.
   6. A crystal flips as it fills or empties.
   7. A change mid-count goes on from the number on screen.
   8. Speed "off" and a rollback's table: the number jumps.
   9. Reduced motion: it still counts, with no flip.
   10. A lunge target's card stays seen at home (20.19); its count runs on the hit.
   11. Buy preview: the purchase prompt sends each option's `price`.
   12. Hovering an option shows "−N" on the own coins; red when short.
   13. The engine's `priceOf` gives both the preview and the payment.

### Catch-up

1. Checked in 20.12 (`SpellController`, `AttackSpellController`, `GameState`).
   1. The client counts the frame time it lost while behind.
   2. Some animations have a worth in seconds; past it they are skipped.
   3. A minion fight is worth 0.8 s; a trigger 0.2 s.
   4. A fight with a hero is worth 0: it always plays.
   5. The rest has no worth: it always plays, at full time.
   6. A skip spends its worth from the lost time.
   7. The threshold comes from the server; a flag can limit it to low-end devices.
2. Built in 20.12.
   1. Ours counts the updates still waiting, not lost time.
   2. Each effect has a weight: minor, normal, key (`weightOf`).
   3. Key: a player's death, a soul, the game's end, the own turn start.
   4. Minor: popups, a stolen coin, another's draw, a flight with no sound.
   5. A reveal that holds its card weighs as much as the move after it.
   6. From 3 waiting: minor effects are dropped (`feel.catchUp`).
   7. From 5 waiting, or 7 effects in one update: the rest at 0.4 of their time.
   8. From 9 waiting: only key effects play, at 0.4.
   9. Before, from 9 waiting every effect was skipped; 3–8 all at 0.4.
   10. Speed "off" still plays nothing, key effects too.

### Other players live

1. Their pointer is sent at most every 0.35 s.
   1. Every 0.25 s while they aim an arrow.
2. Between samples, the client tweens; no jumps.
3. The hand card they hover lifts slightly, 0.15 s.
4. A table card they hover shows a hover glow.
5. A card they hold floats to their play spot in 1 s.
   1. It stands upright and drifts slowly, 10 s per sway.
6. Their aim arrow lingers 0.25 s after they let go.
7. Checked in 20.13 (`RemoteActionHandler`, `ZoneHand`, `TargetReticleManager`).
   1. Sends wait 0.35 s; 0.25 s while the arrow's target changes.
   2. A change in the gap is sent when it ends: the state then, not a queue.
   3. An opponent's hovered hand card moves 0.3 units toward the viewer.
   4. That move takes 0.15 s, expo-out; no scale, no turn.
   5. A hovered board card gets an "opponent mouse over" look.
   6. That look's color is in the card's prefab, not in code.
   7. A held card moves to the opponent's play spot in 1 s.
   8. No ease set: iTween's default, expo-out.
   9. Then it drifts on a path, 10 s each way, linear, ping-pong.
   10. A hidden card stands upright over 5 s, ease-in-out sine.
   11. The arrow head closes 0.1 of the way to the sample each frame.
8. Built in 20.13.
   1. Times in `feel.others`; the server keeps the same gaps.
   2. Server: a sample that comes too soon waits for its gap.
   3. A newer one replaces it, so the last state always arrives.
   4. The reference's sender does the same; ours checks it on the server too.
   5. Client: at most every 0.35 s, 0.25 s while aiming at all.
   6. The head follows by `followStep`: 0.1 per 60 fps frame, any frame rate.
   7. Hovered hand card: its back lifts 14 px toward the table, 0.15 s, expo-out.
   8. 14 px is ours, by eye: the top seat's backs hang off the window.
   9. Hovered table card: a glow in the player's own color.
   10. Red stays for what an aim arrow is over.
   11. Hover is still shared on the own turn only (ADR 016).
   12. A dragged hand card: new presence field `held`, its hand slot.
   13. It is shared on any turn, like an aim: a play in the making.
   14. Its back floats 56 px out and grows to 1.15× in 1 s, expo-out.
   15. Ours: it stays by the hand; we have no play spot on the board.
   16. Then it sways ±3° about the fan's pivot, 10 s per side.
   17. Lifted backs come over their own character (`layer: local`).
   18. Speed "off" and reduced motion: no tween, no sway; the head jumps.

### Emotes

1. At least 4 s between emotes.
2. Many in 5 s: the gap grows to 15 s, then 45 s.

### Sounds

1. Code names a sound prefab; its `SoundDef` names the clips.
   1. Several clips: one picked at random, equal weight.
   2. No random pitch or volume on the sounds we take.
2. Attack impact by attacker power, four prefabs.
   1. 0–5: `FX_Minion_AttackImpact`; 6–7: `…Mid`; 8+: `…Large`.
   2. From power 3 a tavern crowd cheers too; we skip it.
3. Draw: one sound leaves the deck, one lands in hand.
   1. An opponent's draw has one sound of its own.
4. Another player's big card plays a card-play sound.
5. "Your turn" alert; "Job's done!" when nothing is left.
6. Timer: a tick each second, a final tick at zero.
7. Clip per moment: step 20 file, sub-step 20.17.
8. A clip sits in its bundle under its GUID (checked in 20.17).
   1. The catalog in `asset_manifest.unity3d` names the bundle.
9. Checked in 20.14 (`SoundManager`, `SoundConfig` prefab data).
   1. Ducking: a sound's category may duck other categories while it plays.
   2. Each ducked category: a level, a time in, a time out, an ease.
   3. In: from the level now; out once the trigger stops, from where it got to.
   4. Several triggers on one category: the lowest level wins.
   5. All the eases are linear.
   6. 20 %, 0.7 s, 0.7 s: the class defaults and a big card's intro.
   7. That intro ducks music and ambience, not the effects.
   8. "Job's done!" is a voice line: it ducks only ambience, to 50 %.
   9. The "your turn" alert is an effect: it ducks nothing.
   10. Playback limit: a clip group may hold N sounds in their first 10 %.
   11. A new start over the limit is dropped, or stops a lower one.
   12. Some sounds also limit their own copies for a set time.
10. Built in 20.14.
    1. Limit: one sound never starts twice within 80 ms (`feel.mix.gap`).
    2. A start within it is dropped, as the reference's.
    3. Dice at speed "fast" land 64 ms apart: some clatters drop.
    4. Deaths of one update: one death sound, on the first skull.
    5. Each player's skull still shakes the table.
    6. Key cue: only `end-turn-cue`, our one voice line.
    7. It ducks every other sound: to 20 %, 0.7 s in, 0.7 s out, linear.
    8. Ours ducks the effects too; the reference's voice does not.
    9. The step file asks for it; the owner may tune it at 20.16.
    10. One-shots stay HTML audio; their volume is level × duck.
    11. Loops reach the master through a `rest` gain node, ducked alike.
    12. A cue counts as ended 5 s after its start at most.
    13. Sound, not motion: the speed setting does not change it.
    14. A one-shot's element gives its media player back as it ends (19.135).
    15. Browsers cap media players per page; a long game hit the cap.
    16. At most 16 one-shots play at once (`feel.mix.players`).
    17. A new one past that ends the oldest.
11. Size 1 hit since the owner's review (2026-11-08).
    1. `attack-hit-1` is `FX_Minion_AttackImpact`, power 0–5.
    2. It replaces WoW's `unarmedparrymetala` (15.50).

### Response timer

1. Checked in 20.15 (`TurnTimer`).
   1. The rope waits until 20 s are left (`m_RopeCapSeconds`).
   2. Then it burns down, linear, to empty at 0.
   3. The countdown text turns another color at the same time.
   4. Ticks: when the whole seconds left (rounded up) go down.
   5. The count starts at 20, so the first tick is at 19 s.
   6. At 0 the final tick plays instead.
   7. A jump of several seconds ticks once.
   8. Ticks play only where a game mode turns them on.
2. Built in 20.15.
   1. `feel.timer`: shows in the last 20 s; the fuse burns linear.
   2. The seconds after the prompt's title, and a fuse along its foot.
   3. Ticks 19 … 1, then `timer-last` at 0 (`timer.ts`).
   4. Ours: only the player it waits for hears the ticks.
   5. Others see the fuse and the seconds, silent.
   6. A wake-up at each new second, not a 250 ms poll.
   7. On the own prompt the seconds turn dark red to read on gold.
   8. A clock, not an effect: the speed setting does not scale it.
   9. Speed "off" and reduced motion: the fuse steps once a second.
   10. No server change: 5 s allowance and restarts stay.
   11. The Soft vote's timer keeps its full countdown.

### Choosing cards

1. Checked in 20.18 (`ChoiceCardMgr`, `Banner`, `InputManager`, `Card`).
   1. Also its prefab data and three `SoundDef`s, via UnityPy.
   2. A "general" choice: the server lists cards; the chooser picks.
   3. The cards are real card objects, not copies in a dialog.
2. Where they show.
   1. One row, centred on a fixed spot above the board.
   2. 1–3 cards: full size, gaps of about a quarter card.
   3. 4: full size, gap a fifth; 5: 85 %, gap an eighth.
   4. 6 and more: 70 %, gap a tenth; never a second row.
   5. No dimmed backdrop; the board stays in view around them.
3. How they come in, normally.
   1. Each starts at its place, at size 0.
   2. It grows to full in 0.2 s, expo-out (iTween's default).
   3. All cards at once; no stagger, no sound.
4. How they come in from a deck (a card tag: "show from deck").
   1. Each starts on the deck, face down.
   2. Path: deck → 3.6 units above it → its place, 1.5 s.
   3. Its ease is sine-in-expo-out, so expo-out in effect (see Draw).
   4. It turns face up over 1.5 s, after 0.094 s.
   5. Cards leave 0.09 s apart; the first gap is 0.19 s.
   6. Each plays `FX_GameStart09_CardsOntoTable` 0.04 s after it leaves.
5. The banner.
   1. "Choose One" by default; a card may give its own words.
   2. It shows at once above the cards; it has floating motes.
   3. Only one card's banner grows from 0, in 0.5 s.
   4. A headline and an optional caption line.
6. Hover and pick.
   1. A pickable card glows; hover switches it to a brighter glow.
   2. No hover growth in this code; the card is already big.
   3. One pick: a click on the card sends it at once.
   4. The pick plays `HeroDropItem1` (pitch 0.9–1.1).
   5. Nothing else to press: no name buttons, no confirm.
7. Several picks.
   1. A click toggles a card; a chosen card gets an X mark.
   2. A Confirm button sends; End Turn sends too.
   3. Buttons click with `UI_MouseClick_01` (pitch 0.9–1.1).
8. Hide and show.
   1. A "Hide" button beside the cards.
   2. Hide: cards shrink to 0 in 0.2 s; banner and Confirm go.
   3. The button then reads "Show".
   4. Show: cards grow back in 0.2 s, at their places, no deck flight.
9. After the pick.
   1. The unpicked cards vanish at once.
   2. The picked card stays and then moves where the effect sends it.
   3. A card tag may play a "selected" effect and an "unselected" one.
   4. Dredge does: steam on the kept card, smoke on the rest.
   5. Dredge is "look at the bottom 3, put one on top": our flow.
10. Other players.
    1. They see the chooser's cards too, at 1.5 / 2.85 = 53 % size.
    2. At the chooser's side of the board; the server decides the faces.
    3. A pick that comes before they show: they still show for 1 s.
11. Ours, at 20.18.
    1. `OneOnTop` (The Moon, IV The Emperor, IX The Hermit): a choose prompt.
    2. Its options are the 5 cards; the rest go to the bottom, random.
    3. Card options also: Eden's start, Friendly Sack, Red Key.
    4. `PutBackInOrder`, `reorderADeck`, Sleight of Hand, Hagalaz.
    5. Also Spelunker Hat: an order prompt (19.16).
    6. Choose: `ChoiceDialog`; order: `OrderDialog` (UC-13, UC-18).
    7. Both over a dark backdrop, 66 %, with a title line.
    8. Choice cards: 110–180 px wide, 18 px gaps, wrap to rows.
    9. Each card has its name as a red button under it.
    10. Order: dashed places, 1 · top first; cards at table size, 80 px.
    11. No entry motion, no sound, no exit motion.
    12. "View board" hides at once; a bar at the bottom brings it back.
    13. Hover opens the zoom, which covers the cards beside it.
    14. Others see only "Waiting for Ann" and the question.
    15. After The Moon's pick, a "You look at" dialog shows all 5 again.
12. Differences and recommendations.
    1. Entry from the deck.
       1. Ours: all at once. Ref: each flies from the deck, turns face up.
       2. Recommend for deck looks: fly from that deck, 0.09 s apart.
       3. About 1 s each, expo-out; speed "off" shows them at once.
    2. Entry otherwise.
       1. Ours: at once. Ref: grow from 0 in 0.2 s, expo-out.
       2. Recommend for other card choices: the same 0.2 s growth.
    3. Deal sound.
       1. Ours: none. Ref: one card-onto-table sound per card.
       2. Recommend: `FX_GameStart09_CardsOntoTable` per card (rule 8).
    4. Size and rows.
       1. Ours: fixed size, wraps. Ref: one row; shrinks from 5 cards.
       2. Recommend: one row up to 6; 85 % at 5, 70 % at 6.
    5. Backdrop.
       1. Ours: the board dims to 34 %. Ref: no dimming.
       2. Recommend: keep ours; the dim says "nothing else takes clicks".
    6. Title.
       1. Ours: the question, at once. Ref: "Choose One", at once.
       2. Recommend: keep ours; no change.
    7. Name buttons.
       1. Ours: a red name button under each card. Ref: the card only.
       2. Recommend: drop the buttons for card options.
       3. The card keeps its name as its label (keys, e2e, readers).
    8. Hover.
       1. Ours: the zoom covers the neighbours. Ref: a brighter glow only.
       2. Recommend: a brighter glow and a small lift on hover.
       3. The zoom only after the usual delay, or on right click.
    9. Pick sound.
       1. Ours: none. Ref: `HeroDropItem1`, pitch 0.9–1.1.
       2. Recommend: that clip on a pick (rule 8); no random pitch.
    10. After the pick.
        1. Ours: the dialog closes at once. Ref: rest vanish; pick moves on.
        2. Recommend for deck looks: the cards fly back to the deck.
        3. The kept card to the top first; the rest slide under, 0.5 s.
        4. For the order prompt: back in the set order, top last.
    11. The look after The Moon.
        1. Ours: a "You look at" dialog repeats the 5 cards. Ref: nothing.
        2. Recommend: no look dialog after a choice that showed the cards.
        3. As 19.16 did for the order prompt; the history line stays.
    12. Hide and show.
        1. Ours: hides at once. Ref: shrink and grow, 0.2 s each.
        2. Recommend: the same 0.2 s shrink and growth; keep our bar.
    13. Other players.
        1. Ours: text only. Ref: the chooser's cards, small, at their side.
        2. Recommend: small backs over the deck being looked at.
        3. Faces stay hidden; they leave with the cards' flight back.
    14. Several picks.
        1. Ours: discard toggles a glow; the button counts.
        2. Ref: an X mark per chosen card and a Confirm button.
        3. Recommend: keep ours; it already says how many are left.
    15. The order prompt.
        1. The reference has none; Dredge picks one card only.
        2. Recommend: keep our places; take only the motion above.
        3. Cards at the choice size, not table size, if they fit one row.
13. Built in 20.18 (owner's picks: all but 5, 6, 14).
    1. Protocol: `waiting.look` names the deck and the count, for everyone.
    2. It comes from `Question.Looks` (`OneOnTop`) or an order's cards.
    3. Additive: protocol version stays; old saves just lack it.
    4. Deal (12.1): 1 s, sine-in-expo-out, a lift of 80 px above the deck.
    5. It turns face up over 0.9 s after 0.094 s.
    6. Gaps: the second leaves 0.19 s after the first, then 0.09 s each.
    7. `choiceDeal` 0.04 s after each card leaves (12.3).
    8. Growth (12.2): 0.2 s, expo-out, all at once.
    9. Row (12.4): 1–3 full, gap 1/4; 4 full, 1/5; 5 at 85 %, 1/8; 6 at 70 %, 1/10.
    10. More than 6 wrap, at 70 %; the full size stays 110–180 px.
    11. Buttons (12.7): the card is a button, named by the card.
    12. Hover (12.8): an 8 px lift and a brighter glow, 0.15 s.
    13. Zoom: the usual 0.4 s wait; a right click opens it at once.
    14. Pick (12.9): `choicePick`, no random pitch.
    15. Back (12.10): 0.5 s each, ease-in-out, turning face down.
    16. The kept card first; the rest 0.15 s later, 0.06 s apart, under it.
    17. An order: the bottom card first, the top card last.
    18. Look window (12.11): bug 19.130; the history line stays.
    19. Hide and show (12.12): 0.2 s shrink with a fade; 0.2 s growth back.
    20. Others (12.13): one small back per card on the deck, fanned.
    21. They rise in 0.5 s, expo-out, and sink back in 0.5 s.
    22. Others hear no deal sound; the reference was not checked for it.
    23. Order (12.15): choice size up to 6 cards, else table size.
    24. Speed "off" and reduced motion: no deal, growth, flight or rise.
    25. Discards grow in and take the row sizes too.
14. Top or bottom (19.147, owner's picks: card motion and a tag, look-5 too).
    1. Not in the reference; Dredge has no bottom.
    2. Before: only the history told if a looked-at card went under.
    3. Engine: `Question.Peek` shows others a one-card look (Cain's Eye).
    4. Sack Head and XII. The Hanged Man now use `KeepOrBottom`: a `top_choice`.
    5. `ToBottomRandom` emits one public `deck_bottom` for its deck cards.
    6. Kept on top: the card drops onto the pile, 0.5 s, expo-out.
    7. Bottom: it slides out beside the pile, 0.35 s, expo-out.
    8. Then back under it, 0.4 s, ease-in-out; the pile hides it as it goes.
    9. It slides toward the middle of the screen, so it stays on screen.
    10. A tag on the pile, "Top" or "Bottom ×N", fades over 1.5 s.
    11. Others' cards start at their backs over the deck (13.20).
    12. The chooser's own start in a small fan above the pile, face up.
    13. After The Moon's or Eden's pick the window sends the rest under.
    14. The chooser then sees the tag only; the cards already moved.
    15. At most 5 cards move; the tag counts them all.
    16. Reduced motion: the tag only; speed "off": nothing.

### Playing from the hand

1. Checked in 20.20 (`InputManager`, `ZoneHand`, `TargetReticleManager`).
   1. Also `Card`, `PlayErrors`, `GameplayErrorManager`.
   2. Clip names from three `SoundDef`s, via UnityPy.
   3. Tilt and drag loops: see "Drag"; others' view: "Other players live".
2. Picking a card up.
   1. A press and a move of 20 px either way picks the card up.
   2. A click (press and release on the card) picks it up too.
   3. No menu: the card itself is the action.
   4. Picked up, it follows the pointer, button down or not.
   5. It tweens to a held size, 0.7, in 0.2 s (iTween's default, expo-out).
   6. Grab sound: `FX_MinionSummon01_DrawFromHand_01`.
   7. Clip `FX_MinionSummon01_DrawFromHand_01.wav`; no random pitch.
   8. Any own hand card can be picked up on the own turn.
   9. Not during a choice or another card's aim.
3. Where it plays.
   1. The board is one hit box; above the hand is "over the board".
   2. Over it the card's own glow switches to an "over play area" state.
   3. No zone lights up; the card's glow is the only sign.
   4. Its look is prefab data, not traced.
   5. The next release over the board plays it: drag-drop or click-click.
   6. A release over the hand puts it back.
4. Putting it back.
   1. Right click or Esc cancels; so does a release over the hand.
   2. The card goes home with the hand's re-layout: 0.5 s, expo-out.
   3. Loops fade out (see "Drag"); no sound of its own.
5. A card that needs a target.
   1. It is picked up like any card.
   2. On entering the board it stops there; the aim starts at once.
   3. The arrow starts at the player's hero, not at the card.
   4. The pointer hides; the arrow's head carries the effect's text.
   5. Drag loops stop as the aim starts.
   6. Back below the board: the aim ends, the card is held again.
   7. A release on a valid target plays it: drag-to-target.
   8. A release on nothing keeps the aim; a click then picks.
   9. A wrong target: an error line and sounds; the aim goes on.
   10. Right click, Esc, or a click on the source cancels the aim.
   11. Cancel sound: `CancelAttack`, clip `cancel_attack.wav`.
   12. Its pitch varies 0.8 – 1.0.
6. A card with modes ("Choose One").
   1. Dropped on the board, its mode cards show (see "Choosing cards").
   2. A click on the card again, right click, or Esc cancels.
   3. A mode that needs a target aims from the hero, as above.
7. A card that cannot be played now.
   1. It is picked up all the same.
   2. As it enters the board it drops back to the hand.
   3. An error line shows in the middle for 0.1 s per letter.
   4. With `UI_no_can_do` (clip `UI_no_can_do.wav`, volume 0.7).
   5. The hero says a short error line ("Not enough mana").
   6. Not the own turn: nothing is picked up, no message.
8. What a play does.
   1. The game is told at once; the card's move is predicted locally.
   2. Minion: summoned at the slot under the pointer.
   3. Spell: its effects start where it was dropped; prefab data.
   4. The hand closes the gap with its re-layout: 0.5 s, expo-out.
9. Touch.
   1. A tap does not pick a card up.
   2. A drag must rise 80 px above its start.
   3. Back below those 80 px, the card returns to the hand.
   4. A tap on a board card that aims: "drag to play".
10. Keys: only Esc, to cancel; no key plays a card.
11. Ours, at 20.20.
    1. Plain loot and trinkets: drag onto the board plays (UC-13 2).
    2. A click opens a menu: Play, Cancel (UC-13 2a).
    3. The menu holds only Play: no hand card has an ability.
    4. The engine lists abilities only for cards in play and rooms.
    5. A drag starts after 8 px; the card is at hand size × 1.08 at once.
    6. Over the board: no change on the card; the board halo stays.
    7. Dropped off the board: back home in 150 ms, ease-out.
    8. Dropped on it: the card vanishes; a stack row fades in.
    9. Until the server answers it snaps back to the fan meanwhile.
    10. The hand closes the gap at once; no re-layout tween.
    11. A target card: a click lifts it 40 px in the fan; no drag.
    12. The arrow runs from that card; the second click picks (2e).
    13. A click on the source does nothing (owner, 15.33).
    14. A mode card: dragged or Play; a dialog asks the mode (2f, 2h).
    15. Unplayable: not draggable; a click shows why (UC-22).
    16. Sounds: the drag loops only; no grab, cancel or refusal sound.
    17. Touch: tap opens the menu; drag in any direction; tap-tap aims.
    18. Long press zooms (20.5).
    19. Keys: Esc cancels an aim; a hand card takes no focus.
    20. So no key plays a hand card.
12. Differences and recommendations.
    1. A click on a playable card.
       1. Ours: a menu, Play and Cancel. Ref: picks the card up.
       2. Recommend: a click picks up; it follows the pointer.
       3. A click on the board plays; on the hand, puts it back.
    2. The hand card menu.
       1. Ours: it only ever holds Play. Ref: none.
       2. Recommend: no menu for hand cards; table cards keep theirs.
    3. Held size.
       1. Ours: × 1.08 at once. Ref: tweens to its held size in 0.2 s.
       2. Recommend: grow to × 1.08 in 0.2 s, expo-out.
    4. Over the board.
       1. Ours: no sign. Ref: the card's glow changes.
       2. Recommend: a brighter, gold glow on the card over the board.
       3. Back to green when it leaves; no zone glow.
    5. Going back.
       1. Ours: 150 ms, ease-out. Ref: hand re-layout, 0.5 s, expo-out.
       2. Recommend: 0.5 s, expo-out, from where it was let go.
    6. Grab sound.
       1. Ours: none. Ref: `FX_MinionSummon01_DrawFromHand_01`.
       2. Recommend: that clip on pick-up (rule 8).
    7. Target cards: drag.
       1. Ours: click only. Ref: drag; the aim starts at the board's edge.
       2. Recommend: allow the drag; the card stops at the board's edge.
       3. The arrow takes over there; back down resumes the drag.
       4. Release on a target plays; elsewhere the aim stays.
       5. Click-click stays as today.
    8. Arrow start.
       1. Ours: from the card. Ref: from the hero.
       2. Recommend: keep ours; in Four Souls the card is the source.
       3. With 7: from where the card stopped.
    9. Click on the aiming card.
       1. Ours: nothing. Ref: cancels.
       2. Recommend: keep ours (owner, 15.33).
    10. Cancel sound.
        1. Ours: none. Ref: `cancel_attack.wav`, pitch 0.8 – 1.0.
        2. Recommend: that clip on a cancelled aim (rule 8).
        3. No random pitch, as 20.17.
    11. A card that cannot be played.
        1. Ours: no drag; a click says why.
        2. Ref: picked up; bounces back with an error line.
        3. Recommend: keep ours; add `UI_no_can_do` to the "why" note.
        4. No voice: we have no hero lines.
    12. After the drop.
        1. Ours: snaps back to the fan till the server answers; then gone.
        2. Ref: shown as played at once.
        3. Recommend: hold it at the drop point till the answer.
        4. Then it shrinks into its stack row, 0.4 s, expo-out.
        5. Refused: it goes home as in 5.
    13. Hand gap.
        1. Ours: closes at once. Ref: 0.5 s, expo-out.
        2. Recommend: the fan closes in 0.5 s, expo-out (FLIP).
    14. Mode cards.
        1. Ours: a dialog. Ref: mode cards on the table.
        2. Recommend: no change here; see "Choosing cards" 12.
    15. Touch.
        1. Ours: tap opens the menu. Ref: no tap; drag up 80 px.
        2. Recommend: a tap lifts the card, as an aim lifts it.
        3. A tap on the board plays; on the hand, puts it back.
        4. Drag stays; it starts after 8 px in any direction.
    16. Keys.
        1. Ours and ref: no key plays a hand card.
        2. Recommend: hand cards take focus; Enter lifts, Enter plays.
        3. Esc puts it back; a lower priority than the rest.
13. Built in 20.20 (owner's picks: all but 16).
    1. Click (12.1): a click picks a card up; it follows the pointer, button up.
    2. The next click on the board plays it; anywhere else puts it back.
    3. A click on a button on the board (End turn) puts it back too.
    4. Right click and Esc put it back.
    5. Menu (12.2): hand cards have none; table cards keep theirs.
    6. Held size (12.3): 1.08× in 0.2 s, expo-out, on drag or click.
    7. Over the board (12.4): the card's aura turns a brighter gold.
    8. Back (12.5): 0.5 s, expo-out, from where it was, to its fan pose.
    9. Grab (12.6): `cardGrab` on every pick-up, and on a click that aims.
    10. Target drag (12.7): the card stops at its last place off the board.
    11. The aim starts there; the arrow runs from the card (12.8 kept).
    12. Back below the board: the aim ends silently; the drag goes on.
    13. Released on a target it plays; elsewhere the aim and the card stay.
    14. Kept (12.9, 12.11, 12.14): our source click, "why" note, mode dialog.
    15. Cancel (12.10): `aimCancel` on Esc, right click or empty space.
    16. Not on a new update's end of the aim, nor on the drag going back.
    17. Refusal (12.11): `playRefused` with every "why" note.
    18. Also when the server refuses a played card.
    19. After the drop (12.12): the card waits where it was let go.
    20. The answer is the next step; a refusal clears `sent` at the same step.
    21. Left the hand: a copy shrinks into its stack row, 0.4 s, expo-out.
    22. It ends as tall as the row; it fades over the last tenth of its eased way.
    23. No row: its place on the table, else the loot discard.
    24. Still in the hand (a mode asked): it goes home as in 8.
    25. Aimed already (its target asked next): it waits on.
    26. Hand gap (12.13): the rest slide in 0.5 s, expo-out, by left edge.
    27. Only when a card left; a held or waiting card keeps its place.
    28. Touch (12.15): a tap lifts the card as an aim does (40 px, 1.08×).
    29. A tap on the board plays it; elsewhere puts it back; drag stays.
    30. Speed "off" and reduced motion: no growth, flight, shrink or slide.
    31. The wait at the drop point stays: it is no motion.
    32. Times: `feel.play`; none is scaled by the animation speed.
    33. Keys (12.16): dropped by the owner; no key picks up or plays a card.

### Start of game

1. Checked in 20.22 (`MulliganManager`, `MulliganButton`, `MulliganTimer`).
   1. Also `MulliganReplaceLabel`, `ZoneHand`, `Board`, `GameEntity`.
   2. Also `InputManager` (clicks, hover, touch during the screen).
   3. Ten `SoundDef`s via UnityPy; words from `Strings/enUS/GAMEPLAY.txt`.
   4. The screen is "Starting Hand": each player keeps or replaces cards.
   5. One class, two modes: replace (many marks) and choose-one (one pick).
   6. Choose-one is the hero pick of another mode; its lobby was not traced.
2. Before the deal.
   1. The board lights drop at once: a dark blue ambient, the key light off.
   2. A hero intro plays first: both heroes, "VS", announcer lines.
   3. Each hero then drops into its place, `FX_MulliganCoin01_HeroCoinDrop`.
   4. A camera shake of 0.6 s as they land.
   5. Start-hand music begins with the intro.
   6. The deal starts 1 s after the screen loads, once the intro ends.
3. The deal into the middle.
   1. A fixed box across the middle of the board holds the cards.
   2. Its width is split evenly: one slot per card.
   3. Each card flies from the deck: up 3.6 units, then to its slot.
   4. 1.5 s, sine-in-expo-out; it grows to 1.1× a hand card.
   5. It turns face up over 1.5 s, after 0.094 s.
   6. Gaps: 0.19 s after the first card, then 0.09 s.
   7. `FX_GameStart09_CardsOntoTable` 0.04 s after each card leaves.
   8. The same motion, times and sound as "Choosing cards" 4.
4. Who goes first, on the same screen.
   1. 1.1 s after the deal a coin flips at a fixed spot.
   2. Its sound: `FX_MulliganCoin03_CoinFlip`, volume 0.7.
   3. A text by the coin: "You go first" or "You get an extra card".
   4. It shows after 1.8 s, fades in 0.25 s, rises for 2 s.
   5. Then it fades out in 1 s.
   6. 2 s after the flip the second player gets one more card.
   7. That card plays `FX_GameStart20_CardDealSingle`.
   8. The first player's deal left one slot free.
   9. 1.75 s later their cards close up: 0.93 s, ease-in-out cubic.
   10. The picks open 0.6 s after that.
   11. So the first player is known before anyone picks.
5. The title.
   1. A banner at a fixed spot above the cards.
   2. Headline "Starting Hand", caption "Keep or Replace Cards".
   3. It comes when the picks open, not during the deal.
   4. Nothing to pick: "Starting Hand" alone.
   5. Then it grows from 0 in 0.5 s, holds 4 s, shrinks in 0.5 s.
6. Marks, in replace mode.
   1. A click on a card toggles its mark; nothing is sent yet.
   2. A marked card gets a big X over it.
   3. A label "REPLACED" sits under it.
   4. Mark: `HeroDropItem1`; unmark: `GM_ChatWarning`.
   5. Both vary pitch 0.9 – 1.1.
   6. Hover: keyword notes beside the card; no zoom.
   7. Touch: only a short tap marks; its limit is prefab data.
7. Picks, in choose-one mode.
   1. A click picks that card and clears the others.
   2. The same two sounds as marks.
   3. Confirm stays disabled until a card is picked.
   4. Some picks ask twice: Confirm reads "Waiting", with Cancel beside it.
8. The button.
   1. One "Confirm", centred under the cards, at the hand's height.
   2. It comes in with its own birth animation (prefab).
   3. Replace mode: always enabled; no marks keeps every card.
   4. Pressed: disabled, then its death animation.
   5. A game option may drop it; a click then sends at once.
9. The timer.
   1. Text ":SS", rounded seconds, beside Confirm.
   2. Only while the player can still pick.
   3. The server sets its end; at 0 the screen confirms as it stands.
   4. No rope and no ticks in this code.
10. Confirm, in replace mode.
    1. `FX_GameStart28_CardDismissWoosh2_v2`; the X marks and banner go.
    2. Labels turn and shrink to 0 in 0.5 s, expo-in, 0.05 s apart.
    3. Each marked card flies to the deck, dipping down and aside first.
    4. 1.5 s, ease-out cubic, 0.5 s apart.
    5. New cards come back the same path reversed.
    6. 1.5 s, ease-in cubic, 0.5 s apart.
    7. `FX_GameStart30_CardReplaceSingle` 0.5 s after each starts.
    8. 1 s later both decks shuffle: `FX_MulliganCoin09_DeckShuffle`.
11. Confirm, in choose-one mode.
    1. `BG_SelectHero`; the cards do not fly to a deck.
12. Waiting for the other player.
    1. The own cards stay in the middle.
    2. A banner: "Opponent Still Choosing...".
    3. It grows from 0.01 to 1.4× in 0.4 s.
    4. It sits low, by the own hero, its glow moved for that place.
    5. The board stays dark until both are done.
    6. The other's starting cards lie face down in their hand.
    7. Their marks are never shown; only their hand's cards move.
13. Closing.
    1. Labels, buttons, banner and coin text go at once.
    2. The kept cards go from the middle into the hand by its layout.
    3. 0.5 s move, 0.25 s turn and scale, expo-out, all together.
    4. The second player's extra coin card is summoned into the hand.
    5. The lights rise over 3 s, ease-in-out quad.
    6. Game music takes over from the start-hand music.
    7. 1.5 s later the first turn starts.
14. Sounds, from their `SoundDef`s; one clip each.
    1. Deal: `FX_GameStart09_CardsOntoTable` (ours: `choice-deal`).
    2. Mark or pick: `HeroDropItem1`, pitch 0.9 – 1.1 (ours: `choice-pick`).
    3. Unmark: `GM_ChatWarning`, pitch 0.9 – 1.1.
    4. Confirm, replace: `FX_GameStart28_CardDismissWoosh2_v2`.
    5. Confirm, choose one: `BG_SelectHero`.
    6. A replacement lands: `FX_GameStart30_CardReplaceSingle`.
    7. Deck shuffle: `FX_MulliganCoin09_DeckShuffle`.
    8. Coin: `FX_MulliganCoin03_CoinFlip`, volume 0.7.
    9. Extra card: `FX_GameStart20_CardDealSingle`.
    10. Music: a start-hand playlist, a softer one once picks open.
15. Ours, at 20.22.
    1. The table opens complete: characters, items, hands, cents, slots.
    2. No deal motion and no sound at the game's start.
    3. The badge says "Before turn 1 · First player roll" from the start.
    4. Eden (UC-26 6a) gets the 20.18 choice dialog: 3 treasure cards.
    5. Title "Choose your starting item"; the board dims; "View board".
    6. The cards grow in, 0.2 s: no deal, as the choice sends no `look`.
    7. Hover opens the zoom; a click picks at once, with `choicePick`.
    8. No Confirm; a response timer, if set, shows its fuse (20.15).
    9. After the pick the dialog goes; the item is in play, no flight.
    10. The other two go to the treasure deck's bottom, unseen.
    11. Others see "Waiting for Bo · Choose your starting item"; no cards.
    12. The roll list shows from the start.
    13. Others see "Ann throwing…" there while Bo still picks.
    14. UC-26 6a.2 says the roll shows only after her choice.
    15. Several Edens pick one after another, in seat order (engine queue).
    16. Then the roll: a Throw button mid-board, one seat at a time.
    17. Dice by each name; a banner names who goes first; turn 1 starts.
    18. The draft (14.11): server done in 6.6; all pick at the same time.
    19. Each sees only their offers; others see who has picked, not what.
16. Differences and recommendations.
    1. A screen of its own.
       1. Ours: the usual choice dialog. Ref: a "Starting Hand" screen.
       2. Recommend: build a start-of-game mode of the 20.18 choice view.
       3. A banner above the cards: headline and caption ("Choosing cards" 5.4).
       4. Words for the owner, e.g. "Starting item" / "Eden keeps one".
    2. The deal.
       1. Ours: grow in 0.2 s. Ref: dealt from the deck, a sound each.
       2. Recommend: build; Eden's cards are the treasure deck's top 3.
       3. The engine sends `waiting.look` (treasure, 3) for her choice.
       4. The 20.18 deal, sounds and others' backs then play as they are.
    3. Pick, then Confirm.
       1. Ours: a click picks at once. Ref (choose one): select, then Confirm.
       2. Recommend: build, for start-of-game picks only.
       3. An eternal item lasts the whole game; a misclick costs much.
       4. The picked card lifts and keeps its bright glow.
       5. Other choices keep one click (20.18).
    4. Pick sounds.
       1. Ours: `choicePick`. Ref: pick, unpick and confirm sounds.
       2. Recommend with 3: `GM_ChatWarning` on unpick (rule 8).
       3. And `BG_SelectHero` on Confirm.
    5. Others' view.
       1. Ours: a line "Waiting for Bo". Ref: a banner over the board.
       2. Recommend: build a banner, low on the board: "Bo is choosing".
       3. No dialog for them: the table stays usable.
       4. The small backs on the treasure deck come with 2 (20.18 13.20).
    6. The roll list during the pick.
       1. Ours: "Ann throwing…" while Bo picks; UC-26 6a.2 says otherwise.
       2. Recommend: build a fix; no roll list until the choices are made.
       3. Log it in `19-bugs.md` first, as the workflow says.
    7. After the pick.
       1. Ours: the item just appears. Ref: kept cards slide to the hand.
       2. Recommend: build; the picked card flies to her items.
       3. 0.5 s, expo-out, from its place in the row; all see it.
       4. The other two fly back to the deck (20.18 12.10).
    8. The dim.
       1. Ours: the choice dim, gone at once. Ref: dark board, lit in 3 s.
       2. Recommend: keep ours, as 20.18 left it.
    9. Who goes first.
       1. Ours: the roll after the picks. Ref: the coin before the picks.
       2. Recommend: keep the order; R-SETUP-09 puts Eden first.
       3. Later: the roll on the same screen, so the start reads as one.
    10. Timer.
        1. Ours: the 20.15 fuse. Ref: ":SS" beside Confirm.
        2. Recommend: keep ours; one timer look in the game.
    11. Hover.
        1. Ours: the zoom after 0.4 s. Ref: keyword notes only.
        2. Recommend: keep ours (20.18).
    12. Several Edens.
        1. Ours: in seat order. Ref: everyone at once.
        2. Recommend: keep; the banner of 5 names who picks now.
    13. The deal of the whole table.
        1. Ours: the table opens complete. Ref: an intro, then a deal.
        2. Recommend: later; a dealing step of its own.
    14. Music.
        1. Ours: none. Ref: start-hand and game playlists.
        2. Recommend: keep; the game has no music.
17. What the draft pick (14.11) can reuse.
    1. The start-of-game screen of 16.1: banner, row, Confirm.
    2. Pick, then Confirm (16.3), with its sounds (16.4).
    3. The waiting banner (16.5), as "Ann has picked" per seat.
    4. The 20.18 row: up to 5 offers fit one row, 85 % at 5.
    5. Not the deal: offers come from no deck; they grow in (20.18 12.2).
    6. The flight of 16.7: the picked character to its seat.
    7. Text cards when there are no images, as every choice card.
18. Where the 20.18 choice view fits.
    1. Row sizes, growth, deal, hover, glow, pick sound: as they are.
    2. Its deal needs only `waiting.look` from the server.
    3. Others' backs and the flight back to the deck: as they are.
    4. New for this screen: banner, Confirm, picked state, flight to items.
19. Several start-of-game choices at once (owner, 2026-11-09; changes 16.12).
    1. Every player with one chooses at the same time, as the reference.
    2. The engine keeps one prompt; it names all who may answer now.
    3. `Game.StartChoices`: one open choice per such player, seat order.
    4. Each gets the next cards from the treasure deck's top: no card twice.
    5. `Waiting` shows the first; `Waiting.Players` lists all choosers.
    6. `PromptFor(p)`: p's own choice; `Allowed`, `check` and the view use it.
    7. So each chooser sees `waiting.player` as themselves; others see the first.
    8. An answer closes that player's choice; the rest stay open.
    9. The first-player roll starts after the last (19.142).
    10. The record keeps the answers in their order: replays stay the same.
    11. Timer and kicked seats: the server answers for each chooser.
    12. Additive: `players` is new; the protocol version stays.
    13. Engine 0.5.2: older records cannot replay (`startVersion`).
    14. The draft (14.11) is a server step, not an engine prompt; not changed.
20. Built in 20.22 (owner's picks: all of 16, with 16.12 changed by 19).
    1. 16.1: a start-of-game mode of the 20.18 choice view (`confirmFirst`).
    2. Banner: headline "Starting item", caption "Eden keeps one".
    3. It grows from 0 in 0.5 s, expo-out, as the screen opens; it stays.
    4. 16.2: the engine sends `waiting.look` (treasure, 3) for each Eden.
    5. So the 20.18 deal, its sounds and the others' backs play as they are.
    6. 16.3: a click picks a card; it lifts 14 px and keeps its bright glow.
    7. A click on the picked card unpicks it; another card takes the pick.
    8. Confirm under the cards, disabled until a card is picked.
    9. 16.4: `choicePick` on a pick, `choiceUnpick` on an unpick.
    10. `choiceConfirm` on Confirm; no random pitch, as 20.18.
    11. 16.5: others see "Bo is choosing" low on the board, one line per chooser.
    12. It grows from 0 in 0.4 s, expo-out; no window, no dim.
    13. 16.6: bug 19.142, fixed; the badge reads "Before turn 1 · Start of game".
    14. 16.7: the picked card waits at its place after Confirm.
    15. The playback then flies it to its player: 0.5 s, expo-out, silent.
    16. Others see it fly from the treasure deck; the rest fly back as 20.18.
    17. A kept card waits 4 s at most, if no flight comes.
    18. 16.12 changed: several Edens choose at the same time (item 19).
    19. The window keeps its picks while another player's pick comes in.
    20. Kept as they are: 16.8 to 16.11, 16.13, 16.14.
    21. Speed "off" and reduced motion: no growth, deal or flight.
    22. For the draft (14.11): `confirmFirst`, `pickToggle`, `Confirm`.

### Target arrow

1. Checked in 20.23 (`TargetReticleManager`, `TARGET_ARROW_TYPE`).
   1. Also `TARGET_RETICLE_TYPE`, `TargetDamageIndicator`, `InputManager`.
   2. Also `RemoteActionHandler`, `Card`, `CardTextBuilder`, `Entity`.
   3. Prefab data via UnityPy: `Target_Arrow_Bullseye`, `Target_Arrow_Link`.
   4. And `TargetDamageIndicator`, with their meshes and materials.
   5. Units are the reference's world units; ours: 60 px each (see "Drag").
2. Kinds.
   1. Three: own ("friendly"), another player's ("enemy"), and static.
   2. Static: a fixed arrow to a chosen target, one other game mode only.
   3. All three look the same: one red arrow.
   4. The kind changes behaviour, not color.
   5. Some cards swap the head: a "?" for random targets, two special reticles.
3. The look.
   1. A head, a chain of links behind it, and a bullseye.
   2. Links: long dashes, 0.56 wide, 0.99 long, one every 1.2 units.
   3. Head: 1.01 wide, 0.73 long; its tip sits at the pointer.
   4. Bullseye: a round target, 1.13 across, centred on the pointer.
   5. Color is in the textures: red, bright red-orange top, dark red rims.
   6. The materials stay white; no tint by kind or state.
   7. Unlit; drawn on a layer over the cards.
4. The curve.
   1. A parabola from the source to the head, lifted 1.5 units at its middle.
   2. The lift is the same for any length.
   3. Each link turns along the curve's slope; so does the head.
   4. Links: floor(length / 1.2) + 1, at most 15.
   5. Under 1.2 units: no links; the head keeps a fixed tilt.
   6. Past 15 links the part near the source stays empty.
5. The march.
   1. The links move toward the head at 0.5 units/s.
   2. One link's place repeats every 2.4 s.
   3. The link by the head fades out over its last 1.2 units.
   4. Its opacity: (distance / 1.2)^6, so it goes late and fast.
   5. The link at the tail fades in: (offset / 1.2)^2.
6. Following the pointer.
   1. Own arrow: the head is at the pointer each frame; no damping.
   2. The head turns to face away from the source.
   3. The system pointer hides while the own arrow shows.
   4. Over an own hand card the arrow hides and the pointer shows.
   5. Another player's arrow: see "Other players live" 7.11.
7. Over a valid target.
   1. The bullseye shows at once; off at once when the pointer leaves.
   2. Valid targets carry a "valid target" highlight all the aim.
   3. The hovered one turns to a stronger "valid target, mouse over" state.
   4. Both looks are prefab data; not traced.
   5. A big card never opens on hover while an own arrow shows.
8. Over an invalid target.
   1. Nothing changes: no bullseye, no other mark.
   2. A click there: an error line and sounds ("Playing from the hand" 5.9).
9. Appear and go.
   1. Created: drawn at once at the pointer; no tween.
   2. Ended: hidden at once; no tween.
   3. Its prefabs hold an animation component with no clips.
   4. No sound of its own: no audio in its prefabs or code.
   5. Cancel: `CancelAttack`, as "Playing from the hand" 5.11.
   6. Another player's arrow lingers 0.25 s ("Other players live" 6).
10. What the head carries.
    1. A short text: the card's own "target arrow text".
    2. White letters on a yellow-orange burst; text box about 3 units wide.
    3. Just above the head, 0.75 units toward the top.
    4. Touch: shown by the source instead.
    5. Empty text: no burst, no text.
    6. A choose-one mode with no such text shows its card text.
    7. Another player's arrow: never any text.
    8. Attacks: none; only other modes' abilities and objects get one.
11. Where it starts.
    1. A spell: the player's hero ("Playing from the hand" 5.3).
    2. A minion's own aim and an attack: that card on the board.
    3. A hero power: its button.
    4. Another player's: their held card, or their hero for a hidden one.
12. Attack aims.
    1. A click or drag on an attacker starts an arrow from it.
    2. The attacker plays its own attack sounds; Four Souls has none.
    3. Its attack spell starts with the aim ("Attack" 1, windup).
    4. A cancel sends that spell to its cancel state.
    5. Taunts on the other side show.
    6. No text on the head.
13. Another player's arrow (`RemoteActionHandler`).
    1. Made only when they point at a card; not at empty board.
    2. Its head goes to that card's middle, not to their pointer.
    3. It starts at its source and glides out (0.1 per frame).
    4. No bullseye, no text.
    5. It ends 0.25 s after they stop.
14. Sounds: no new clip.
    1. The arrow plays none; cancel is `cancel_attack.wav` (20.20).
15. Ours, at 20.23.
    1. One `AimArrow` for the own aim and others' arrows (15.32).
    2. Square blocks 28 × 26 px, one every 36 px; straight, no curve.
    3. Head 48 long, 60 wide; tip at the pointer; red gradient, dark rim.
    4. Blocks march 36 px per 0.5 s: 72 px/s, at the animation speed.
    5. Tail block fades in, head block fades out, both linear.
    6. The whole arrow at 85 % opacity, with a drop shadow.
    7. Valid targets glow red (`hl-aim`) all the aim.
    8. Over one, a ring with cross ticks: 68 px, in from 1.5× over 0.25 s.
    9. Over an invalid one: nothing; a click gives the why note.
    10. The system pointer stays; a hand pointer over cards.
    11. The hover zoom opens after 0.4 s over any card, aim or not.
    12. It covers the target and the head (shot 02 and 03).
    13. The arrow shows only after the first pointer move; none on touch.
    14. It ends at once; no sound but `aimCancel` on a cancel.
    15. The head carries nothing; a banner says "Choose a target".
    16. Starts at the card: lifted in the fan, or stopped at the board's edge.
    17. Others: from their lifted hand back to the card or point.
    18. Their head snaps to a card's middle, else to their pointer.
    19. A new arrow starts at its first sample; no glide out.
    20. Over any card their ring shows and the card gets a red frame.
    21. Their arrow goes at once; 70 % opacity, under the windows.
    22. Attacks: no arrow.
    23. Attack, then priority passes, then a table pick: blue glow, dim.
16. Differences and recommendations.
    1. The curve.
       1. Ours: straight. Ref: an arc, links and head along its slope.
       2. Recommend: build a 2D arc bowing up the screen.
       3. Peak at the middle; about 1.5 units (90 px), less on short arrows.
       4. Blocks and head turn with the slope.
    2. The links.
       1. Ours: square blocks every 36 px. Ref: long dashes, a short gap.
       2. Recommend: build dashes about 1.8 times longer than wide.
       3. Keep our colors, rim and shadow.
    3. The march.
       1. Ours: 72 px/s, linear fades. Ref: 0.5 units/s (30 px/s).
       2. Recommend: build the slower march and the ref's fade curves.
       3. Head end ^6, tail ^2.
    4. Zoom while aiming.
       1. Ours: the zoom opens and hides the target. Ref: no big card.
       2. Recommend: build; no hover zoom while the own arrow shows.
       3. Owner: how to read a target then (see the open questions).
    5. The pointer.
       1. Ours: the system pointer over the head. Ref: hidden.
       2. Recommend: build; no cursor while aiming with a mouse.
    6. First frame.
       1. Ours: no arrow until the pointer moves. Ref: at once.
       2. Recommend: build; draw it from the click point.
    7. Head text.
       1. Ours: none; the banner names the card. Ref: a short line on a burst.
       2. Recommend: later; it needs a short line per card.
       3. Owner: or the damage only, e.g. "1", for damage cards.
    8. The hovered valid target.
       1. Ours: ring only. Ref: bullseye and a stronger highlight.
       2. Recommend: build; the hovered target's red glow brightens.
       3. The ring: keep ours; its 0.25 s grow is ours, by eye.
    9. Invalid target.
       1. Same in both: nothing on hover; a click explains.
       2. Recommend: keep ours.
    10. Others' arrow start.
        1. Ours: at the first sample. Ref: glides out from its source.
        2. Recommend: build; a new head starts at the source.
    11. Others' arrow end.
        1. Ours: at once. Ref: lingers 0.25 s ("Other players live" 6).
        2. Recommend: build the linger.
    12. Others' ring.
        1. Ours: the ring over any card. Ref: no bullseye on theirs.
        2. Recommend: build; drop the ring there, keep the red frame.
        3. Our ring means "a valid target" only in the own aim.
    13. Others pointing at nothing.
        1. Ours: an arrow to their pointer. Ref: only at a card.
        2. Recommend: keep ours; it shows the aim in the making.
    14. Attack aims.
        1. Ours: a table pick, no arrow. Ref: an arrow from the attacker.
        2. Recommend: owner's pick; build or keep.
        3. Build: the arrow from the own character during the attack pick.
        4. Others see it as any aim; the blue pick glow stays.
    15. Where it starts.
        1. Ours: the card. Ref: the hero for a spell.
        2. Recommend: keep ours ("Playing from the hand" 12.8).
    16. Head over the own hand.
        1. Ref: the arrow hides. Ours: no target is a hand card.
        2. Recommend: keep ours.
    17. Head variants and the static arrow.
        1. Ref: "?" heads for random targets; a kept arrow in one mode.
        2. Recommend: later.
    18. Sounds.
        1. Ref: none of the arrow's own. Recommend: keep; no new clip.
17. Built in 20.23 (owner picks, 2026-11-09).
    1. Every "build" of 16; 4 changed, 14 built; 7 stays later.
    2. Times and sizes in `feel.aim`; `feel.zoom.rest` for the zoom.
    3. Arc: a quadratic curve, bowing up the screen (`arc` in `aim.ts`).
    4. Its middle: 90 px off the line, or a quarter of a short line.
    5. Only the lift across the line bends it: full from 20° off upright.
    6. Below that it fades to none; an upright arrow stays straight.
    7. So the bow never flips side as the pointer crosses upright.
    8. As the reference's lift, seen from its camera.
    9. Dashes: 40 long, 22 across, one every 48 px; by eye.
    10. Our colors, rim and shadow stay; dashes and head turn with the curve.
    11. March: 30 px/s, the reference's 0.5 units/s; speed scales it.
    12. Fades: the head end (left / step)^6, the tail (way / step)^2.
    13. Each dash is one WAAPI loop over one step; all start at one time.
    14. The chain ends 24 px short of the tip, under the head's back.
    15. By eye: at the head's full length a gap showed before the head.
    16. Own arrow: drawn at once from the last mouse place; no move needed.
    17. It follows its card each frame, e.g. while the card lifts.
    18. The system pointer hides while it shows (`html.aim-cursor`).
    19. The hovered valid target gets `data-aim-over`: a brighter glow.
    20. The ring stays on the own aim; its 0.25 s grow is `feel.aim.ring`.
    21. Zoom (4, owner's change): none while the arrow moves.
    22. A card zooms once the arrow rests still on it for 0.8 s.
    23. 0.8 s: twice the 0.4 s table hover; by eye.
    24. Any card, target or not: a rest means "let me read this".
    25. Any pointer move closes the zoom; the next rest opens it again.
    26. It opens beside the card; the arrow's layer is above the zoom.
    27. So the head on the target stays in view.
    28. Others: a new head glides out from its source (`followStep`).
    29. A stopped arrow stays 0.25 s where it was drawn.
    30. No ring on others' arrows; the red frame stays.
    31. Attack (14, owner): the pick of what to attack draws the arrow.
    32. From the own character to the pointer; the blue glow stays.
    33. The hovered option's blue glow brightens; the ring shows.
    34. The click still answers the pick, as in 18.11.
    35. The monster deck counts as the deck option's target.
    36. Others see it as any aim: `presence.from` is the character.
    37. Speed "off" and reduced motion: no march, no glide; it still shows.
    38. No new sound.

## Consequences

### Positive

1. Step 20 has measured targets, not guesses.
2. The principles fit our playback queue (ADR 012).
3. Most changes are timing and easing, cheap to try.

### Negative and risks

1. Four Souls is not Hearthstone.
   1. 2–4 players; anyone may hold priority, not only the active one.
   2. A "locked" turn start only delays our view; never the server.
   3. Our turn banner is visual only; it never blocks input (owner).
2. Shakes and elastic motion can tire some players.
   1. Animation speed "off" and reduced motion must turn them off.
3. New sounds must pass the loudness check (ADR 022).
   1. Ducking is a mixer change, not a per-sound volume.
4. A timeline in DOM may not hit 0.12 s exactly at low frame rates.
5. Public builds get none of the new sounds; they stay silent.

### Neutral

1. Times are scaled by our animation speed setting.
2. Each new animation follows ADR 012: tag, row, claim.
3. The decompiled source never enters this repo.
