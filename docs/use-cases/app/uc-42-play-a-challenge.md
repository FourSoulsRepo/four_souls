# UC-42: Play a challenge

**Module:** `frontend`, `rules_engine`
**Status:** Implemented
**Actors:** Two players, together
**Goal:** Both players see and play a challenge to its end: Krampus or Trick or Treat
**Preconditions:** The host started a challenge game (UC-41)

## Main scenario (happy path)

1. Both players see the challenge table (ADR 028, Decision 11).
   1. Pilfering Krampus lies left of the monster slots, landscape.
   2. The challenge panel stands where the bonus souls would be.
   3. In it: the rule card of the difficulty, landscape, zoomable.
   4. Under it: the timer, a d8 with the rounds left.
   5. Then the 4 presents, face down, on the treasure back.
   6. Each present names its kind: ¢, Loot, Treasure, Soul.
   7. Each shows its counters as a counter badge.
2. Each turn starts with the active player's roll (R-KRAMPUS-12).
   1. The stack names Pilfering Krampus as the roll's card.
   2. Counters come off a present; a "−N" rises from it.
3. A present reaches 0 counters (What's in the Box?, R-KRAMPUS-08).
   1. It flies to the treasure discard; its place shows ✗.
   2. Minions come from the monster deck, one by one.
   3. They lie in a row under the monster slots.
4. On Normal, the turn can't end before a random minion attack.
   1. Attack goes to a random minion; the game picks it.
   2. Beside another owed attack, a dialog offers "a random minion".
5. While souls are too few or a minion is left, Krampus is guarded.
   1. Attacking, it is no target: it dims while the targets glow.
6. A monster's reward or a soul could be forgone (All I Want for Christmas).
   1. The gainer gets a yes/no question.
   2. Yes puts that many counters on the present; "+N" rises from it.
7. The players hold enough souls (Happy Holidays!, R-KRAMPUS-10).
   1. The presents turn face up.
   2. The active player picks who gains each.
   3. A gained present is an item in that player's area; its place shows ✓.
8. Both players attack Krampus once no minion is left.
   1. It glows as a target; the stack shows the odds.
9. Krampus dies: both players see "Challenge won!".
10. The timer ticks down at the end of each round.
11. Trick or Treat: Mischievous The Haunt lies in the final boss place (14.33).
    1. Its trick counters show as a counter badge: 12, 14 or 14.
    2. The panel shows the Trick/Treat card, Trick up; no presents.
12. As each turn starts, the active player may flip it (R-TRICK-11).
    1. A yes/no question names the side it would show.
    2. Yes turns the card in the panel to its other side.
13. A player who takes damage gets a fear counter (R-TRICK-07).
    1. Fear shows by the player's name: a dot, the number, "Fear".
    2. Its tooltip: monsters +1 evasion per 2 on the active player.
    3. Monsters show the raised evasion; the zoom names the Haunt.
14. Each attack on a monster starts with the side's roll (R-TRICK-13, R-TRICK-15).
    1. The active player throws it, then the attack roll.
    2. Changed rewards show on the monster as a ribbon and under its zoom.
    3. Gained abilities show as the copy mark on the monster.
    4. Hard rolls 2x, Ultra Hard 3x: one mark per gained monster (19.157).
    5. Hover on a mark zooms the gained monster; the monster's zoom names them all.
15. Trick side: a monster dies; the Haunt loses a trick counter (R-TRICK-12).
    1. The Haunt's roll and the trick counter come at once.
    2. The active player orders them in the trigger order dialog.
16. Souls reach the trick counters: the Haunt can be attacked (R-TRICK-06).
17. The Haunt dies: both players see "Challenge won!".

## Alternative scenarios and errors

* **3a. The fourth present is destroyed:** "Challenge lost", and why under it.
* **3b. A minion is an event:** it resolves and goes to discard; no slot stays.
* **4a. Hard and Ultra Hard:** no forced attack; extra rolls per minion.
* **6a. The player answers No:** they gain the reward as usual.
* **10a. The timer reaches 0:** "Challenge lost", and "The timer ran out" under it.
* **Card images missing:** text cards; the boss and rule card are landscape.
* **Card images missing:** the presents show the plain treasure back.
* **12a. The player answers No:** the card keeps its side.
* **13a. Ultra Hard:** each damage puts fear on each player.
* **15a. Treat side:** a kill removes no trick counter.
  * Each gain of rewards removes a fear counter; the active player picks whose.
* **16a. Trick or Treat is lost only by the timer:** "Challenge lost", "The timer ran out".
* **Card images missing:** the Haunt, Trick and Treat are landscape text cards.

## Postconditions

* The game is over for both players; the record can be saved.
