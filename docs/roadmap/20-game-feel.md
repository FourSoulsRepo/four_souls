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
12. Order: 20.18 runs before the step review, 20.16 (owner, 2026-11-08).
13. Order: 20.19 runs before the step review too (owner, 2026-11-08).
14. Order: 20.20 runs before the step review too (owner, 2026-11-08).
15. Order: 20.21 runs before the step review too (owner, 2026-11-08).
16. Order: 20.22 runs after 20.18 and before the step review (owner, 2026-11-08).
17. Order: 20.23 runs after 20.20 and before the step review (owner, 2026-11-08).

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
   6. Sound by size: `attack-hit-1`, `attack-hit-2`, `attack-hit-3` (20.17).
      1. Size 1 is the reference's too since the owner's review (20.16).
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
      1. From the loot discard: face up, as it is public (owner, 20.16).
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

### [x] 20.9 Others' plays, paced

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

### [x] 20.10 End turn states

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

### [x] 20.11 Counters that count

1. Goal: cents, souls and HP change one step at a time.
2. Tasks:
   1. A change of N plays as N steps, 0.2 s apart.
   2. Above 5 steps: steps speed up to fit 1 s.
   3. Each step flips the digit in 0.1 s.
   4. A buy previews its cost on the cents counter while hovered.
3. Done when:
   1. Vitest: step times for N = 1, 5, 20.
   2. e2e: a +3¢ shows 1, 2, 3.

### [x] 20.12 Catch-up by importance

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

### [x] 20.13 Others' hands live

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

### [x] 20.14 One sound per moment

1. Goal: sound stays clear when many things happen at once.
2. Tasks:
   1. Deaths in one update: one death sound.
   2. A playback limit per sound: no two starts within 80 ms.
   3. Key cues duck the rest to 20 %: 0.7 s in, 0.7 s out.
   4. Ducking lives in the mixer; no per-sound volume (ADR 022).
3. Done when:
   1. Vitest: three deaths queue one sound.
   2. Vitest: duck curve values.

### [x] 20.15 Response timer only when it matters

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

### [x] 20.16 Step review

1. Goal: the owner plays and compares with the reference.
2. Tasks:
   1. A full 2-player game recorded on video for the owner.
   2. ADR 027 updated with every tuned number.
   3. ADR 027 becomes Accepted for what was built.
3. Done when:
   1. Owner approval.
4. Owner notes (2026-11-09)
   1. Step 20 marked solved by the owner.
   2. No video; the owner judges the feel in play tests.
   3. Tuned numbers are in each "Built in 20.N" item of ADR 027.
   4. Open bugs 19.132, 19.140, 19.143 are fixed during play tests.

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
| `attack-hit-1` | `FX_Minion_AttackImpact.wav` | `784c4afefc057624fa90cab6478eed93` | `AttackSpellController.m_ImpactDefHandles`, power 0–5 | 20.3 |
| `choice-deal` | `FX_GameStart09_CardsOntoTable.wav` | `d6db11e6fe15fb6448aa73618050217c` | ADR 027, 20.18 research | 20.18 |
| `choice-pick` | `HeroDropItem1.wav` | `667e129aef9fb854cb4e74466fe7e942` | ADR 027, 20.18 research | 20.18 |
| `card-grab` | `FX_MinionSummon01_DrawFromHand_01.wav` | `630f053471521454c96fe96837040c13` | ADR 027, 20.20 research | 20.20 |
| `aim-cancel` | `cancel_attack.wav` | `066172a886131dc4b96649adffddf9cb` | ADR 027, 20.20 research | 20.20 |
| `play-refused` | `UI_no_can_do.wav` | `9c034c03dde8fd8458679423cb554253` | ADR 027, 20.20 research | 20.20 |
| `history-hover` | `history_event_mouseover.wav` | `5a250f9b221350647b46287650dd530c` | ADR 027, 20.21 research | 20.21 |
| `choice-unpick` | `GM_ChatWarning.wav` | `2a6b3e2ef8f202b47bb1d4b9fc734d6b` | ADR 027 "Start of game" 14.3 | 20.22 |
| `choice-confirm` | `BG_SelectHero.wav` | `fcecbadc52023e14ebf008249961886a` | ADR 027 "Start of game" 14.5 | 20.22 |

Kept as it is: `loot-draw` (draw out).
The reference's own: `draw_card_1..3.wav`.
Swapping it is the owner's call; not part of this step.

`attack-hit-1` replaced `attack-hit` (owner, 2026-11-08).

1. The old file, WoW's `unarmedparrymetala`, did not match in style.
2. Its clip found by name among the audio bundles' `AudioClip`s.
3. The Mid and Large GUIDs found the same way match the table.

### [x] 20.18 Choosing cards, compared with the reference

1. Goal: picking among shown cards feels like the reference's choice.
   1. Owner, 2026-11-08: the reference's "choose one of 3".
   2. Ours: look at the top cards of a deck, put one on top, the rest on the bottom.
   3. E.g. XVIII The Moon: top 5 loot, 1 on top, the rest on the bottom.
   4. The same flow for monster cards, e.g. top 4 monsters.
2. Tasks:
   1. Read the reference's choice flow (rule 10).
      1. Code: `ChoiceCardMgr.cs`, its banner, `CHOICE_TYPE`.
      2. How the cards come in: from where, spacing, size, times, curves.
      3. How one is picked: hover, click, confirm, sounds.
      4. What happens to the rest: where they go, how they leave.
      5. Hiding the choice to look at the table, and coming back.
      6. Asset data (prefabs, `SoundDef`) via UnityPy, as in 20.17.
   2. Write the findings into ADR 027, a new "Choosing cards" section.
   3. Compare with our flow, every card that uses it.
      1. Engine prompts and the client's pick view.
      2. List each difference with a recommendation.
   4. The owner picks what to build; then it is built here.
   5. A sound it needs comes as rule 8 says.
   6. Owner picks (2026-11-08): every recommendation.
      1. ADR 027 "Choosing cards" item 12, all 15.
      2. Kept as they are: 5, 6, 14 (as recommended).
3. Done when:
   1. ADR 027 has the findings and the differences.
   2. The owner's picks are built, with before and after screenshots.
   3. e2e: the pick flow, set up so it never depends on the deal.

### [x] 20.19 Only the attacker flies

1. Goal: the attacker flies into the target; the target stays and takes the hit.
   1. Owner, 2026-11-08, from a recorded Hearthstone match.
   2. There only the card that deals the damage moves.
   3. Ours since 20.2: both cards meet half way.
2. Findings so far (`AttackSpellController`, checked 2026-11-08):
   1. Impact point: source + `m_ImpactStagingPoint` × (target − source).
   2. `m_ImpactStagingPoint` is 1 by default: the target's own place.
   3. The target "moves" to that point, so it stays where it is.
   4. Only a proposed attack uses 0.5: half way, both move.
   5. The attacker stops short by `m_SourceImpactOffset` (−0.25 by default).
   6. ADR 027 "Attack" items 2.1, 2.2 and 7.1 read the half-way case as the normal one.
3. Tasks:
   1. Check both fields in the prefab data with UnityPy (rule 10).
   2. Check how far the attacker overlaps the target at contact.
   3. Strike: only the attacker's copy moves; the target card stays home.
   4. Contact, burst, sound, shakes and the hit number stay at contact.
   5. The target may still get the 20.3 shake; it is no flight.
   6. A player panel as target: unchanged, the attacker goes to it.
   7. Fix ADR 027 "Attack" in place; note what was checked.
   8. Update the `lunge` row in `docs/wiki/Animations.md`.
4. Done when:
   1. e2e: the target card does not move during the strike.
   2. e2e: the hit still lands at contact, after the strike starts.
   3. Before and after screenshots of windup, contact and settle.

### [x] 20.20 Playing a card from the hand, compared with the reference

1. Goal: playing a hand card feels like the reference, without a menu.
   1. Owner, 2026-11-08: today a click opens a menu with a "Play" button.
   2. Ours: a drag onto the board plays a loot card (UC-13 2).
   3. Ours: a click opens Play, Use, Cancel (UC-13 2a).
   4. Ours: a card that needs a target aims by two clicks (UC-13 2e, 2f).
2. Tasks:
   1. Read the reference's play flow (rule 10).
      1. Click on a hand card: does it pick the card up to follow the pointer?
      2. Where a drop or a second click plays it; the play zone and its glow.
      3. How a play is cancelled: drop on the hand, right click, Esc.
      4. A card that needs a target: where the arrow starts, when.
      5. What the card does on a play: its move, times, curves, sounds.
      6. A card that cannot be played now: what a click or drag shows.
   2. Write the findings into ADR 027, a new "Playing from the hand" section.
   3. Compare with our flow, every kind of hand card.
      1. Plain loot, loot with a target, loot with a mode, trinkets.
      2. Hand cards with abilities: today the menu's "Use".
      3. Touch and keyboard: what replaces a drag there.
      4. List each difference with a recommendation.
   4. The owner picks what to build; then it is built here.
   5. A sound it needs comes as rule 8 says.
   6. Owner picks (2026-11-08): every recommendation but 16.
      1. ADR 027 "Playing from the hand" item 12.
      2. 16, keyboard play: dropped; no key plays a card.
      3. 12 holds the card till the answer, as the reference.
3. Done when:
   1. ADR 027 has the findings and the differences.
   2. The owner's picks are built, with before and after screenshots.
   3. UC-13 matches the new flow.
   4. e2e: each kind of play, set up so it never depends on the deal.

### [x] 20.21 History strip, compared with the reference

1. Goal: the left history strip works like the reference's.
   1. Owner, 2026-11-08.
   2. Ours: icons along the board's left edge, newest on top (UC-15 8, 9).
   3. Ours: hover shows the card large, the used ability and the owner.
   4. ADR 027 "History" has only a short first look (items 1, 2).
2. Tasks:
   1. Read the reference's history strip (rule 10).
      1. Which plays make a tile; how many tiles fit; what drops off.
      2. A new tile: where it comes from, how the rest move, times, curves.
      3. What a tile shows: art, frame, badges, damage or death marks.
      4. Hover: the large card, its targets, their damage or death, sounds.
      5. Hover delays, and what happens on click or touch.
      6. Asset data (prefabs, `SoundDef`) via UnityPy, as in 20.17.
   2. Write the findings into ADR 027 "History", appended items.
   3. Compare with our strip and the full history panel.
      1. Which entries get an icon; the turn line; the History button.
      2. List each difference with a recommendation.
   4. The owner picks what to build; then it is built here.
   5. A sound it needs comes as rule 8 says.
   6. Owner picks (2026-11-08): every recommendation.
      1. ADR 027 "History" item 11, all 18.
3. Done when:
   1. ADR 027 has the findings and the differences.
   2. The owner's picks are built, with before and after screenshots.
   3. UC-15 matches the new strip.
   4. e2e: the strip's new behaviour, set up so it never depends on the deal.

### [x] 20.22 A start-of-game screen, like the reference's starting hand

1. Goal: the game opens with a clear screen for choices made before turn 1.
   1. Owner, 2026-11-08: the reference shows "Starting hand" before the first turn.
   2. Ours first: Eden picks her eternal item at setup (UC-26 6a).
   3. Today that is a plain choice dialog, then the first-player roll.
   4. Later the same screen serves the draft pick (14.11); not built here.
2. Tasks:
   1. Read the reference's starting hand screen (rule 10).
      1. Code: `MulliganManager.cs`, `MulliganButton.cs`, `MulliganTimer.cs`.
      2. How it opens: the deal into the middle, sizes, times, curves.
      3. The title, the button, the timer, picks and their marks.
      4. What the other player sees meanwhile.
      5. How it closes: where the kept cards go, the move into the game.
      6. Sounds and music changes; asset data via UnityPy, as in 20.17.
   2. Write the findings into ADR 027, a new "Start of game" section.
   3. Compare with our setup: Eden's pick, the waiting players, the roll after it.
      1. List each difference with a recommendation.
      2. Note what the draft pick (14.11) can reuse.
   4. The owner picks what to build; then it is built here.
   5. Reuse the 20.18 choice view where it fits.
   6. A sound it needs comes as rule 8 says.
   7. Owner picks (2026-11-09): every recommendation of "Start of game" 16.
      1. Banner words as proposed: a title, and "Bo is choosing" for others.
      2. Confirm is for Eden now and for the draft (14.11) later.
      3. 12 changed: several Edens choose at the same time.
      4. 6 is bug 19.142; it is fixed here.
3. Done when:
   1. ADR 027 has the findings and the differences.
   2. The owner's picks are built, with before and after screenshots.
   3. UC-26 matches the new screen.
   4. e2e: Eden's pick on the new screen, set up so it never depends on the deal.

### [x] 20.23 Target arrow, compared with the reference

1. Goal: aiming at a target looks and feels like the reference's arrow.
   1. Owner, 2026-11-08.
   2. Ours: a chain of red blocks with a big head (UC-13 2e, 15.32).
   3. Ours: a red ring on a valid target; a short note says what to pick.
   4. Others see it live since 20.13; the head glides to each sample.
2. Tasks:
   1. Read the reference's arrow (rule 10).
      1. Code: `TargetReticleManager.cs`, `TARGET_ARROW_TYPE.cs`.
      2. Its look: the shaft, its segments, the head, colors by kind.
      3. How it follows the pointer: curve, lift, speed, damping.
      4. Over a valid target and over an invalid one: what changes.
      5. How it appears and goes away: times, curves, sounds.
      6. What the head carries: text, damage, an icon.
      7. Asset data (prefabs, `SoundDef`) via UnityPy, as in 20.17.
   2. Write the findings into ADR 027, a new "Target arrow" section.
   3. Compare with ours: own aim, others' arrows, attack aims.
      1. List each difference with a recommendation.
   4. The owner picks what to build; then it is built here.
   5. A sound it needs comes as rule 8 says.
   6. Owner picks (2026-11-09): every recommendation of "Target arrow" 16.
      1. Plus 14: attack aims get the arrow, from the attacker.
      2. 4 changed: the zoom shows when the arrow rests still on a target.
      3. The rest time is tuned by eye; any move closes the zoom.
      4. 7 (head text) stays later, as recommended.
      5. The arc: about 90 px at the middle, tuned by eye.
3. Done when:
   1. ADR 027 has the findings and the differences.
   2. The owner's picks are built, with before and after screenshots.
   3. UC-13 2e matches the new arrow.
   4. e2e: aiming with the new arrow, set up so it never depends on the deal.
