# UC-41: Host a challenge

**Module:** `frontend`, `internal/settings`, `internal/server`
**Status:** Implemented
**Actors:** Host (a player), one friend
**Goal:** A co-op challenge waits in the lobby for the second player
**Preconditions:** The host screen is open (UC-03)

## Main scenario (happy path)

1. The host picks "Challenge" in "Game mode" (ADR 028).
   1. "Regular game" is the default.
   2. It is in the simple setup, above the players.
2. Two selects show beside it: the challenge and its difficulty.
   1. The challenges come from the app (`getChallenges`).
   2. Now only "How the Krampus Stole Christmas".
   3. Difficulties: Normal, Hard, Ultra Hard; Normal first.
3. The players are fixed at 2; the select is locked.
   1. A line under it says a challenge is for exactly 2 players.
4. The host presses Host (or Create on a server).
5. The create message has 2 seats and the challenge and difficulty.
   1. They go without "Detailed setup" too.
6. The game list shows "Challenge: How the Krampus Stole Christmas, Hard".
7. The friend joins (UC-04); the table line names the challenge too.
8. Both get ready; the game starts with the final boss (UC-05, item 17).

## Alternative scenarios and errors

* **1a. Detailed setup:** every option stays but bonus souls.
  * Sets, picking or draft, timer, random seats, deck ratios, rooms, rollback.
  * Bonus souls are hidden; the server plays a challenge without them.
* **1b. The last setup:** the host screen starts from the last host.
  * It keeps the mode, the challenge and the difficulty.
  * A settings file from before them starts with "Regular game".
  * A remembered challenge the app lacks: the first one, at Normal.
* **1c. The app knows no challenge:** "Challenge" can't be picked.
* **1d. Back to "Regular game":** the players select is free again.
  * It shows the count chosen before the challenge.
* **5a. The server lacks the challenge:** refused (`error.no_challenge`, UC-03 4e).

## Postconditions

* A 2-seat challenge game waits in the lobby with its difficulty.
