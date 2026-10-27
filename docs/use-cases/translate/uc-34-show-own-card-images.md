# UC-34: Show own card images

**Module:** `internal/assets`, `pkg/card_db`
**Status:** Implemented
**Actors:** Translator
**Goal:** Card images in their language show in the app, from their own folder
**Preconditions:** Own images as `<lang>/<set>/<file>`; the language has card names

## Main scenario (happy path)

1. The translator keeps images in an own folder or repo.
   1. E.g. `de/b2/the_d6.webp`, named as in the card data.
   2. Access to the private images repo is not needed.
2. They start the app with `FOUR_SOULS_IMAGES` set to that folder.
   1. `wails dev -m` or a release build.
3. They pick their language in Settings.
4. A card with an own image shows that image.
5. Other cards fall back as before.
   1. The built-in image in that language, if any.
   2. Then the English image.

## Alternative scenarios and errors

* **2a. The folder is not there:** a log line; the app starts without it.
* **2b. No private images (public clone):** own images still show.
  * Other cards show as text cards.
* **4a. A link or path leading out of the folder:** it is not read.
* **4b. Other folders there (`assets/`, `data/`):** ignored.
* **4c. Images added while the app runs:** they show after a restart.
* **3a. The language has no card names or interface file:** not in Settings.
  * Images alone do not list a language (UC-32).

## Postconditions

* Nothing is written; the folder is only read.
* Card data and other images are unchanged.
