# UC-01: Start the app with the fan-game notice

**Module:** `internal/legal`, `frontend/src/screens/Splash.tsx`
**Status:** Implemented
**Actors:** Player
**Goal:** The player sees that this is an unofficial fan game, then reaches the main menu
**Preconditions:** The app is installed

## Main scenario (happy path)

1. The player starts the app.
2. The splash screen shows the title and the notice.
3. The notice names the designer and the publisher.
4. The names and "buy it here" are links to official resources.
5. After 5 seconds the main menu opens.
6. The main menu shows app and rules engine versions in the corner.

## Alternative scenarios and errors

* **4a. The player clicks a link:** The system browser opens it; the splash keeps running.
* **2a. The notice fails to load:** The splash shows only the title for 5 seconds.

## Postconditions

* The main menu is open.
