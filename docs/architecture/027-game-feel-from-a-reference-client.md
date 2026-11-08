# 27. Game feel: lessons from a reference client

* **Status:** Proposed
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
2. Strike: attacker and target meet in 0.12 s, linear.
   1. The meeting point is between them, not at the target.
   2. The target also moves toward it.
3. On contact:
   1. Damage applies at that exact moment.
   2. The impact effect plays; bigger for more power.
   3. Nearby cards shake away from the impact point.
4. Return: the normal zone re-layout, 1 s, easeOutExpo.
5. A hero turns to face the target in 0.3 s; back in 0.3 s.
6. A declared but not yet resolved attack moves half way.
7. Impact point, checked in 20.2 (`AttackSpellController`).
   1. The target moves to it; the attacker stops at its own edge.
   2. A proposed attack's point is half way to the target.
   3. The impact effect plays at the target's place at contact.
   4. Contact is the attacker's move ending, not a timer.
8. Built in 20.2.
   1. Windup 0.25 s, ease-out: lift to 1.1×, draw back 6 %.
   2. Strike 0.12 s, linear; a card target meets it half way.
   3. Settle 0.5 s, expo-out: the low end of principle 3.
   4. Our burst carries "−N", so it pops on the target's home.
   5. The target's new HP shows once its copy is home.
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
   10. A lunge target's card hides while its copy settles; its count runs meanwhile.
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
