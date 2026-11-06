# UC-38: Check a situation on the website

**Module:** `frontend/src/site`, `cmd/website`, `internal/website`
**Status:** Implemented
**Actors:** Player, rules judge, developer
**Goal:** The engine says whether a move is allowed, and by which rule
**Preconditions:** A browser with `DecompressionStream`; the built site

## Main scenario (happy path)

1. The player opens a check link: `#/en/check/<data>`.
2. The page loads the rules engine.
   1. It shows a progress bar while `site.wasm.gz` downloads.
   2. It unpacks and starts it in the browser; no server answers.
3. The page unpacks the link into the situation payload.
   1. The payload shows in the text box.
4. The engine checks the situation.
5. The page shows the answer.
   1. Allowed or Refused, with the payload's title.
   2. A refusal: its reason in the page's language.
   3. The rule ID that decided it, as a link.
   4. With `expect` in the payload: whether the answer matches.
   5. An allowed move: a list of its events.
6. The player clicks the rule ID.
7. The rules page opens at that rule, scrolled to it and highlighted.

## Alternative scenarios and errors

* **1a. No link:** the player pastes a payload and presses Check.
  * Or presses one of the example buttons.
  * The address becomes the situation's link.
* **1b. No language in the link:** the browser's language, else English.
* **3a. The link is cut or broken:** the page says so.
* **4a. The payload is not valid:** the engine's reason shows, in English.
* **5a. Copy link:** the page copies a link to the payload in the box.
  * The link also shows in a field, to copy by hand.
  * Above 2000 characters, a warning: some chats cut long links.
* **5b. Copy payload:** the page copies the payload text.
* **5c. The browser does not allow copying:** the page says so.
* **2a. The browser lacks `DecompressionStream`:** the page says to update it.
* **2b. The engine does not load:** the page shows why.
  * The rules page still works; it needs no engine.
* **7a. The link names a rule the site lacks:** the page says so.

## Postconditions

* The answer shows; nothing was sent to a server.
* The address is a link that opens the same situation.
