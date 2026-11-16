# UC-45: Watch a record

**Module:** `bindings/records`, `internal/replay`, `frontend`
**Status:** Implemented
**Actors:** Player
**Goal:** The player watches a saved match on the table, at their own pace
**Preconditions:** The Records screen lists a format 2 record (UC-44)

## Main scenario (happy path)

1. The player presses **Watch** on a row.
2. The app opens the record without the rules engine (ADR 031).
3. The table shows the record's first step at once.
   1. Seat 0 sits at the bottom; nobody is "you".
   2. The words name every player: "Ann's turn", not "Your turn".
   3. Every hand shows: the record keeps the judge view.
4. The replay bar shows at the bottom left.
   1. "Step i / n · Turn t" and a list of the turns.
   2. A slider over all steps.
   3. Step back, Play or Pause, Step forward.
   4. Speeds 0.5×, 1×, 2×, 4×; 1× at first.
   5. **Back to records**.
5. The player presses **Step forward**.
   1. The step's events play as effects, as in a live game.
   2. Then the table shows the step's view.
6. The player presses **Step back**.
   1. The table shows the step before at once, with no effects.
   2. The history drops the rows of the steps after it.
7. The player picks a turn in the list, or drags the slider.
   1. The table shows that step at once, with no effects.
8. The player picks a speed and presses **Play**.
   1. Each step plays as in 5, one after another.
   2. After a step's effects: 0.8 s if it had events, else 0.2 s.
   3. The speed divides the effects' times and these pauses.
9. Autoplay reaches the last step and stops.
   1. Play and Step forward turn off.
10. The player presses **Back to records** or Esc.
    1. The Records screen lists the records again.

## Alternative scenarios and errors

* **3a. The table can't act:** no hand card plays or drags.
  * No End turn, Pass, turn actions, Undo or game menu.
  * No response timer and no end-of-turn cue.
* **5a. Step forward during autoplay:** autoplay pauses first.
* **6a. Any jump during effects:** the effects stop; the view shows at once.
* **8a. Pause:** autoplay stops on the step shown.
* **8b. Animations off in Settings:** steps show without effects.
* **9a. A finished record:** its last step plays the winner's banner.
* **\*a. A step can't be read:** the bar says why; the table stays.

## Postconditions

* The record is still open in the app until another one opens.
* Nothing was sent to any server.
