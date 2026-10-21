# UC-01: Start the app with the fan-game notice

**Module:** `internal/legal`, `frontend/src/screens/Splash.tsx`, `frontend/src/screens/MainMenu.tsx`, `frontend/src/components/Sky.tsx`
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
7. The main menu shows the foursouls.ru banner at the top (15.41).
8. The main menu and the game setup screens drift on a starry sky (15.42).
9. The player clicks Exit, the last menu button; the app closes (15.79).

## Alternative scenarios and errors

* **4a. The player clicks a link:** The system browser opens it; the splash keeps running.
* **2a. The notice fails to load:** The splash shows only the title for 5 seconds.
* **7a. No private images, or the banner fails to load:** The menu shows the text title.
  * The text title also shows until the banner has loaded; the buttons do not move.
* **8a. No private images:** The menu screens keep their plain color.
* **8b. Animations off:** The sky stands still.

## Postconditions

* The main menu is open, or the app closed on Exit.
