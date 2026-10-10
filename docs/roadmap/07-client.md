# Step 7. Client

Goal: a playable desktop table.
Ends with an internal play-test build on the cards done so far.

Ideas: M-01, V-01 – V-08, ST-01 – ST-03, ST-05, ST-06, B-06, N-11.

---

### [ ] 7.1 UI libraries ADR

1. Goal: pick the few libraries we need.
2. Tasks:
   1. State store for server views and events.
   2. Animation library (transform and opacity only).
   3. Drag and drop with pointer events (touch works).
   4. Prefer small, popular, maintained libraries.
3. Done when:
   1. **Owner** accepts the ADR.

### [ ] 7.2 App shell

1. Goal: splash, main menu, settings storage.
2. Tasks:
   1. Splash with the notice (1.6).
   2. Menu: host, join, records (disabled), settings.
   3. Versions in the corner (A-02).
   4. Settings saved by Go in the user config folder.
3. Done when:
   1. Settings survive a restart.

### [ ] 7.3 Host and join screens

1. Goal: get into a game.
2. Tasks:
   1. Host: simple setup; "Detailed" expands the rest (GS-10).
   2. Join: address field, recent addresses.
   3. Lobby: seats, ready, errors like "client outdated".
3. Done when:
   1. Two apps on a LAN reach the same lobby.

### [ ] 7.4 Card component

1. Goal: one card component used everywhere.
2. Tasks:
   1. Image card, or blank card with name, text, stats.
   2. Sizes: hand, table, zoom.
   3. Zoom on hover, and on long press for touch (B-06).
3. Done when:
   1. Missing images show readable blank cards.

### [ ] 7.5 Table layout

1. Goal: Hearthstone-like table (V-01, V-02).
2. Tasks:
   1. Own area at the bottom.
   2. 2 players: opponent on top.
   3. 3–4 players: left, top, right.
   4. Shop, monsters, decks, discards in the middle.
   5. Game mat as background (ST-03).
3. Done when:
   1. Layout fits at 1280×720 and up.

### [ ] 7.6 Event playback and animations

1. Goal: events play in order as animations (V-08).
2. Tasks:
   1. Queue of server events; one plays at a time.
   2. Card moves, flips, damage, dice roll.
   3. Animation speed setting (ST-05).
   4. UI never changes game state on its own.
3. Done when:
   1. Fast event bursts play without glitches.

### [ ] 7.7 Player actions

1. Goal: act through allowed actions only.
2. Tasks:
   1. Highlight what is allowed now.
   2. Drag and drop to play cards; click to pick targets.
   3. Choice dialogs; "look at cards" Hearthstone-style (V-07).
   4. Rejected intents show the reason.
3. Done when:
   1. A full turn is playable with mouse and with touch.

### [ ] 7.8 Stack and responses

1. Goal: clear response moments.
2. Tasks:
   1. Stack panel with pending items.
   2. Pass, "Skip all", cancel skip (N-06).
   3. Response timer display.
3. Done when:
   1. Players always see why the game waits.

### [ ] 7.9 Info panels

1. Goal: never get lost (V-03, V-05, V-06).
2. Tasks:
   1. History panel: icons, hover details, hide button.
   2. One entry per resolved stack; whole match, scrolling.
   3. Turn indicator with turn number.
   4. Player counters computed by the engine.
3. Done when:
   1. History shows only what the viewer may see.

### [ ] 7.10 Settings screen

1. Goal: ST-01, ST-02, ST-03, ST-05, ST-06.
2. Tasks:
   1. Fullscreen or window.
   2. Dark or light theme for menus and panels.
   3. Game mat choice.
   4. Animation speed.
   5. Nickname; asked on the first start.
3. Done when:
   1. Every setting applies without restart where possible.

### [ ] 7.11 Pause, vote, end of match

1. Goal: UI for N-08 and RP-05.
2. Tasks:
   1. Pause overlay with vote: wait or kick.
   2. End screen: winner, save record.
3. Done when:
   1. Use cases: disconnect vote, save record.

### [ ] 7.12 Linux, touch, internal build

1. Goal: check weak spots; first play-test.
2. Tasks:
   1. Test on Linux (WebKitGTK); fix slow animations.
   2. Test touch input.
   3. Internal build with the cards done so far.
3. Done when:
   1. **Owner** plays a LAN game on the internal build.

### [ ] 7.13 Save and continue UI

1. Goal: UI for N-11.
2. Tasks:
   1. Host menu: save the game and leave.
   2. Host screen: list saved games and load one.
   3. Lobby shows saved seats and who is back.
3. Done when:
   1. Use case covers the UI flow.
