# Step 20. Game feel

Goal: the table feels alive like Hearthstone: fast hits, calm settles, clear pacing.

Ideas: V-54, V-13, V-22.
Source: ADR 027 (numbers measured in a Hearthstone-based client).

Rules for this step (owner, 2026-11-06):

1. Never copy code, images or names from the reference.
   1. Only timings, curves and rules from ADR 027.
   2. The decompiled source stays in `FourSoulsRepo/tools/.cache/`.
   3. Sounds are the exception: see rule 8.
2. Sub-steps run one after another, each in a fresh subagent.
3. UI sub-steps leave before and after screenshots.
   1. In `.cache/screenshots/20.N/before/` and `after/`.
   2. Never committed.
4. Every changed or new animation updates `docs/wiki/Animations.md`.
5. Animation speed "off" turns every new motion off.
6. Times in ADR 027 are a start; tune by eye, note the change there.
7. Work goes on the `step-12-rules-website` branch (owner).
8. Sounds come from the reference client (owner, 2026-11-06).
   1. They are Hearthstone's, (c) Blizzard, like the three of 15.50.
   2. They live only in the private images: `assets/sounds/`.
   3. Never in this repo or in a public build.
   4. Our own file names; `NOTICE` lists each source clip name.
   5. Each passes `loudness check` before its commit (ADR 022).
   6. Without the private images the game stays silent, as today.
   7. The extraction tool lives in `FourSoulsRepo/tools` only.
9. Order: 20.17 runs right after 20.1, before 20.2.
10. Unclear how the reference does something: check its code (owner).
    1. Decompiled C#: `FourSoulsRepo/tools/.cache/reverse/fanstone/src/`.
    2. Asset data (prefabs, `SoundDef`): UnityPy, as traced in 20.17.
    3. Read it for numbers and rules; never copy code from it.
    4. Note what was checked in ADR 027.
11. Images repo changes (owner, 2026-11-06).
    1. In `pkg/card_db/images`, on a local branch `step-20-sounds`.
    2. A normal commit: no amend, no push.
    3. Game repo commits never include the submodule pointer.
    4. The owner merges, pushes and bumps the pointer by hand.

Owner notes (2026-11-06):

1. The ADR names the reference client openly; it is a fan game.
2. Step 20 is next in the work order.
3. Sounds the step needs come from Fanstone's files (20.17).

---

### [x] 20.1 Feel constants in one place

1. Goal: every timing and easing from ADR 027 has one home.
2. Tasks:
   1. `frontend/src/feel.ts`: named durations and easing curves.
   2. Curves as CSS strings: expo-out, sine-in-expo-out, elastic-out.
   3. Elastic-out as a keyframe list (CSS has no elastic).
   4. Move today's constants there: lunge, burst, fly, banner.
   5. No visible change in this sub-step.
3. Done when:
   1. A Vitest test checks every curve string parses.
   2. `make frontend e2e` passes; screenshots unchanged.

### [x] 20.2 Fast strike

1. Goal: an attack hits like ADR 027 "Attack".
2. Tasks:
   1. Windup: the attacker lifts and grows a little, about 0.25 s.
   2. Strike: attacker and target meet half way in 0.12 s, linear.
   3. Damage, burst and sound on contact, not at a time offset.
   4. Return: both settle back in 0.5 s, expo-out.
   5. Replace the `lunge` keyframes; keep its id and claim.
3. Done when:
   1. e2e checks the hit comes after the strike starts, not before.
   2. Screenshots of windup, contact and settle.

### [x] 20.3 Hits by size

1. Goal: a bigger hit looks bigger (ADR 027 principle 4).
2. Tasks:
   1. Three hit sizes by damage: 1, 2–3, 4 and more.
   2. Burst size and glow grow with the size.
   3. Cards near the target shake away from it: 0.1, 0.5, 1.
   4. A player's death or a lethal hit: a short table shake.
   5. No shake when speed is "off" or reduced motion is on.
   6. Sound by size: `attack-hit`, `attack-hit-2`, `attack-hit-3` (20.17).
      1. Size 1 keeps today's `attack-hit` (15.50).
      2. No crowd cheer; the reference's is a tavern crowd.
3. Done when:
   1. Vitest covers the size picker.
   2. e2e: a 1-damage and a 4-damage hit, with screenshots.

### [x] 20.4 Damage number that stays

1. Goal: damage reads clearly and never slows the game.
2. Tasks:
   1. The "−N" grows from 0 with an elastic overshoot, 1 s.
   2. It holds 2 s, then fades over 1 s.
   3. Playback goes on after its first frame.
   4. Several hits on one card stack, not overlap.
3. Done when:
   1. e2e: the next effect starts while the number still shows.
   2. Animation clash fixture stays green.

### [x] 20.5 Instant hand hover, patient table hover

1. Goal: the pointer gets the answer it expects (principle 2).
2. Tasks:
   1. Hand card hover: snaps upright and large, no transition.
   2. Neighbours stay in place.
   3. Table card hover: zoom after 0.4 s; grows in 0.15 s.
   4. Moving between table cards inside the zoom keeps it open.
   5. Fan only from 4 cards; angle 40° + 2° per card.
3. Done when:
   1. Vitest: fan angle and drop per card count.
   2. e2e: zoom absent at 0.2 s, present at 0.6 s.

### [x] 20.6 Drag tilt

1. Goal: a dragged card feels held in a hand.
2. Tasks:
   1. Tilt by pointer speed: pitch ±40°, roll ±45°.
   2. Spring back to flat in about 0.2 s when still.
   3. Transform only; no layout change (ADR 002).
   4. Off when speed is "off" or reduced motion is on.
   5. Drag sounds from 20.17 (ADR 027 "Drag").
      1. `drag-bed` loops, fades in to 15 % over 0.5 s.
      2. `drag-air` loops; volume follows speed, up to 50 %.
      3. Both fade out in 0.2 s on drop.
      4. Those levels are mixer gains, like ducking (20.14).
3. Done when:
   1. Vitest covers the tilt function with fixed inputs.
   2. Screenshot of a card mid-drag.

### [x] 20.7 Loot draw with a look

1. Goal: the drawer sees each new loot card before it joins the hand.
2. Tasks:
   1. Own draw: deck → lift → show spot → hand, 1.5 s total.
   2. Turns face up on the way; sound out, sound in.
      1. Out: `loot-draw`, as today; in: `loot-land` (20.17).
   3. Playable 0.3 s after it lands.
   4. Draws 4 and on: each 11 % faster, floor 0.556.
   5. Others' draws: back only, short, no show spot.
      1. One sound: `loot-draw-other` (20.17).
   6. Replace `fly` for own draws; claim as `arrive`.
3. Done when:
   1. Vitest: speed factor per draw index.
   2. e2e: a 3-card draw plays one card at a time.

### [x] 20.8 Turn start sequence

1. Goal: a new turn is one clear, ordered moment.
2. Tasks:
   1. Banner: fade in 0.25 s, small punch, stays 1.5 s, out 0.25 s.
   2. During it: items recharge and turn upright.
   3. Then the turn's loot draw (20.7).
   4. Then start-of-turn effects.
   5. Visual only: no input lock, no change to priority (owner).
   6. Buttons stay live; a click during the banner works as always.
   7. The banner plays `turn-start` (20.17); own turn only.
3. Done when:
   1. e2e: a pass during the banner is accepted.
   2. Screenshots of each phase.

### [ ] 20.9 Others' plays, paced

1. Goal: everyone reads another player's card before it acts.
2. Tasks:
   1. Big card for others' plays only; never for your own.
   2. Spectators and replays: big card for every play.
   3. The effect starts 1 s after the big card shows.
   4. The big card stays 3 s, or until the next one comes.
   5. Check with `played` (15.74); one animation, one claim.
   6. The big card plays `card-play-other` (20.17) as it shows.
      1. Never for your own plays, like the big card.
3. Done when:
   1. e2e: own play shows no big card; another player's does.
   2. Clash fixture green.

### [ ] 20.10 End turn states

1. Goal: End turn tells the player what is left.
2. Tasks:
   1. Three looks: waiting, your turn, nothing left.
   2. Nothing left: stronger glow and color change.
   3. A short cue 1.5 s later, once per turn.
   4. Skip the cue if a play became possible.
   5. The cue sound passes `loudness check` (ADR 022).
   6. The cue is `end-turn-cue` from 20.17 (owner, 2026-11-06).
      1. It is a spoken English line, "Job's done!".
3. Done when:
   1. Vitest: state picker from allowed actions.
   2. e2e: the cue fires once per turn at most.

### [ ] 20.11 Counters that count

1. Goal: cents, souls and HP change one step at a time.
2. Tasks:
   1. A change of N plays as N steps, 0.2 s apart.
   2. Above 5 steps: steps speed up to fit 1 s.
   3. Each step flips the digit in 0.1 s.
   4. A buy previews its cost on the cents counter while hovered.
3. Done when:
   1. Vitest: step times for N = 1, 5, 20.
   2. e2e: a +3¢ shows 1, 2, 3.

### [ ] 20.12 Catch-up by importance

1. Goal: a lagging client skips small things, never big ones.
2. Tasks:
   1. Each effect gets a weight: minor, normal, key.
   2. Key: deaths of players, souls, game end, own turn start.
   3. A long queue drops minor effects first, then shortens normal.
   4. Key effects always play, at least shortened.
   5. Replaces the plain "pending > 8 skips all" in `scaleFor`.
3. Done when:
   1. Vitest: a queue of 20 keeps every key effect.
   2. e2e with a slowed client still shows a death.

### [ ] 20.13 Others' hands live

1. Goal: you see what another player is about to do.
2. Tasks:
   1. Their hovered hand card lifts slightly (0.15 s), back hidden.
   2. Their hovered table card gets a hover glow, in their color.
   3. A card they drag floats up and sways slowly until played.
   4. Pointer samples at most every 0.35 s; 0.25 s while aiming.
   5. Tween between samples; never jump.
   6. Never recorded (N-15), like live aiming.
3. Done when:
   1. Server test: sample rate limit holds.
   2. e2e with two pages: the lift shows on the other page.

### [ ] 20.14 One sound per moment

1. Goal: sound stays clear when many things happen at once.
2. Tasks:
   1. Deaths in one update: one death sound.
   2. A playback limit per sound: no two starts within 80 ms.
   3. Key cues duck the rest to 20 %: 0.7 s in, 0.7 s out.
   4. Ducking lives in the mixer; no per-sound volume (ADR 022).
3. Done when:
   1. Vitest: three deaths queue one sound.
   2. Vitest: duck curve values.

### [ ] 20.15 Response timer only when it matters

1. Goal: time pressure shows only near the end.
2. Tasks:
   1. The timer bar shows only in its last 20 s.
   2. Each second ticks; the last tick sounds different.
      1. `timer-tick` and `timer-last` (20.17).
   3. No server change (owner).
      1. The timer restarts after every intent (6.7).
      2. It already adds 5 s for animations (`animationAllowance`).
3. Done when:
   1. e2e: no timer at 25 s left; shown at 15 s.

### [ ] 20.16 Step review

1. Goal: the owner plays and compares with the reference.
2. Tasks:
   1. A full 2-player game recorded on video for the owner.
   2. ADR 027 updated with every tuned number.
   3. ADR 027 becomes Accepted for what was built.
3. Done when:
   1. Owner approval.

### [x] 20.17 Sounds from the reference client

1. Goal: each moment plays the clip the reference plays for it.
2. Tasks:
   1. Tools repo: `reverse/sounds.sh` extracts Fanstone's audio.
      1. Input: `Data/OSX/*.unity3d` bundles (UnityFS, `AudioClip`).
      2. UnityPy in a venv under `.cache/`; output WAV.
      3. Output stays in `FourSoulsRepo/tools/.cache/`.
      4. Bundles lack the Unity version: set `6000.3.11f1`.
      5. It reads it from `Fanstone.app/Contents/Info.plist`.
   2. Each clip is found by its GUID, not by a guess.
      1. `asset_manifest.unity3d` maps a GUID to its bundle.
      2. Its `base_assets_catalog`: `m_assets` and `m_bundleNames`.
   3. The clips, traced in the decompiled code (table below).
      1. Code names a sound prefab, or a `SoundDef` field.
      2. The prefab's `SoundDef.m_AudioClip` and `m_RandomClips` name the clips.
      3. A traced clip needs no owner pick.
   4. Variants: a sound may have several files.
      1. Files `<name>-1.mp3`, `<name>-2.mp3`, …
      2. `playSound` picks one at random, equal weight, like the reference.
      3. `App.Sounds()` gives a list of URLs per sound.
   5. `loudness level`, then `loudness check` on each file (ADR 022).
      1. Loops are levelled like any sound.
   6. Commit into the images submodule as rule 11 says.
   7. `App.Sounds()`, the bridge and `fakeBridge.ts` get each new name.
      1. Nothing plays them yet; each later sub-step wires its own.
   8. `NOTICE`: each file with its source clip name.
   9. Wiki: `Client.md` sound list and `Build-modes.md`.
3. Done when:
   1. Go test: `App.Sounds()` gives every file of every new name with images.
   2. Without the images every new list is empty.
   3. Vitest: `playSound` picks among variants.
   4. `loudness check` passes on every sound.
   5. `wails build -m -tags embed` and `make e2e` pass without images.

Clips (traced 2026-11-06 in Fanstone 37.0.254477):

| Our name | Clip(s) | Clip GUID(s) | Traced from | Used in |
|----------|---------|--------------|-------------|---------|
| `attack-hit-2` | `FX_Minion_AttackImpactMid.wav` | `d1b25cb27c4f1e74496a7fc861b4d642` | `AttackSpellController.m_ImpactDefHandles`, power 6–7 | 20.3 |
| `attack-hit-3` | `FX_Minion_AttackImpactLarge.wav` | `7de34d0f8af3ba741a5d0c05c250c419` | same, power 8+ | 20.3 |
| `loot-land` ×3 | `add_card_to_hand_1..3.wav` | `399c3b1e667bad040bb99af727fc4a9b`, `50eb80b80f5f83040a12cba144491f4e`, `39f15b2dc9439e546a09a4f325295a69` | `Card.DrawFriendlyCardWithTiming` | 20.7 |
| `loot-draw-other` ×3 | `draw_card_and_add_to_hand_opp_1..3.wav` | `6286e3185a1b9ae43b9c05d179289141`, `84876a2336256f545a52908d553697a1`, `8584838dc81c4fd4583095345b0aeff6` | `Card.DrawUnknownOpponentCard` | 20.7 |
| `drag-bed` | `card_motion_loop_magical.wav` | `3fcf3ef84ba2e55449878be534ac6095` | `DragCardSoundEffects` | 20.6 |
| `drag-air` | `card_motion_loop_air.wav` | `454f4be0fa7176143ad1299b3143ee99` | `DragCardSoundEffects` | 20.6 |
| `turn-start` | `ALERT_YourTurn_0v2.wav` | `ed4d6c88b0ba7d4488c71fe607d87e0e` | `TurnStartManager` | 20.8 |
| `card-play-other` ×3 | `play_card_from_hand_1..3.wav` | `60b55d955f0bbe34eaeee0004eff70ad`, `2beb57fd1a017a44099e48c6f005053c`, `c74ba91347b52f94180dd0061a259f2c` | `HistoryCard` big card, `CastSpellCardFromHandSepll` | 20.9 |
| `end-turn-cue` | `VO_JobsDone.wav` (enUS) | `9cf1adfef5afae941bbf780a83efdf6d` | `EndTurnButton.PlayEndTurnSound` | 20.10 |
| `timer-tick` | `Countdown_Tick.wav` | `d4cbcbff6b2c84843839a2787776a2b8` | `TurnTimer.m_TickSound` | 20.15 |
| `timer-last` | `Countdown_Tick_Final.wav` | `e99b87e5dc68bcc4dbe04e8e05190b6a` | `TurnTimer.m_FinalTickSound` | 20.15 |

Kept as they are: `attack-hit` (size 1) and `loot-draw` (draw out).
The reference's own: `FX_Minion_AttackImpact.wav` and `draw_card_1..3.wav`.
Swapping them is the owner's call; not part of this step.
