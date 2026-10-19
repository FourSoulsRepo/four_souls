# UC-20: Play with deck ratios

**Module:** `internal/server`, `rules_engine`, `frontend`
**Status:** Implemented
**Actors:** Host
**Goal:** Each deck holds 100 cards in the official ratios (GS-07)
**Preconditions:** The host creates a game (UC-03)

## Main scenario (happy path)

1. The host ticks the card sets and the number of players.
2. The app checks each ratio group against the chosen sets.
3. All four sets fit: "Deck ratios" in detailed setup is ticked.
4. The host keeps it on; it is on even without detailed setup.
5. The server checks the groups again for the sets, players and host bans.
6. The engine picks each group's count at random from its copies.
7. Treasure, loot and monster decks get 100 cards each.
8. The rest of the cards stay out of the game.
9. The table shows "deck ratios" in its options.
10. The record header says `deck_ratios`.

## Alternative scenarios and errors

* **2a. Too few cards in a group:** the box is unticked and locked.
  * Its tooltip names the short groups, e.g. "Butter Beans 4 of 5".
  * The Base Game alone is short; so are Gold Box and Four Souls+ without Requiem.
  * 2 players drop the 3+ cards; some sets fit 3 players but not 2.
* **4a. The host unticks the box:** the decks hold every copy.
  * The app remembers the choice for the next host screen.
* **5a. The server finds a group short:** it plays with whole decks.
  * Host bans can make a group short.
  * The table then shows no "deck ratios".

## Postconditions

* Each deck group holds exactly its official count.
* Players learn nothing about which cards are in the decks (A-16).
