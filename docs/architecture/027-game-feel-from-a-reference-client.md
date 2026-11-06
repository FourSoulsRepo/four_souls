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

### Damage and death

1. Damage number: grows 0 → 1 over 1 s, easeOutElastic.
2. It holds 2 s, then fades out over 1 s.
3. Deaths in one step play together.
4. A dying card's death effect waits 0.3 s for keyword deaths.

### Others' plays

1. Another player's card shows large in the middle.
2. Its effect starts 1 s after the big card shows.
3. The big card stays 3 s (minion) or 4 s (spell).
4. A transformed card halves that time.
5. A spectator sees the chosen option highlighted for 1.75 s.

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

### End Turn button

1. States: waiting, your turn, nothing left to do.
2. Nothing left: the button glows and changes color.
3. A short cue plays 1.5 s later, once per turn.
4. The cue is skipped if a play becomes possible meanwhile.

### Counters

1. Crystals are added, readied and spent one at a time.
2. 0.2 s between each crystal; the number flips in 0.1 s.

### Other players live

1. Their pointer is sent at most every 0.35 s.
   1. Every 0.25 s while they aim an arrow.
2. Between samples, the client tweens; no jumps.
3. The hand card they hover lifts slightly, 0.15 s.
4. A table card they hover shows a hover glow.
5. A card they hold floats to their play spot in 1 s.
   1. It stands upright and drifts slowly, 10 s per sway.
6. Their aim arrow lingers 0.25 s after they let go.

### Emotes

1. At least 4 s between emotes.
2. Many in 5 s: the gap grows to 15 s, then 45 s.

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
