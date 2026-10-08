# Step 13. Web client

Goal: play and watch from a browser, e.g. on a phone.

Ideas: A-15.

---

### [ ] 13.1 Browser bridge

1. Goal: the React app runs without Wails.
2. Tasks:
   1. Second bridge implementation for the browser.
   2. Hosting and records are hidden in the browser.
3. Done when:
   1. The same frontend builds for app and browser.

### [ ] 13.2 Served by the game server

1. Goal: open the server address in a browser.
2. Tasks:
   1. Server embeds the built web client.
   2. Card images served only if an image folder is next to it.
   3. Otherwise blank cards with text.
3. Done when:
   1. A phone on the same Wi-Fi plays a game.
   2. Use case: join from a browser.

### [ ] 13.3 Phone check

1. Goal: usable on small touch screens.
2. Tasks:
   1. Test layout on phone sizes.
   2. Fix touch issues.
3. Done when:
   1. A full game is playable on a phone.
