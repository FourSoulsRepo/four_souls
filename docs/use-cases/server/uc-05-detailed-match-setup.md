# UC-05: Detailed match setup

**Module:** `internal/server`
**Status:** In progress
**Actors:** Host, players
**Goal:** Characters are banned and picked the way the host chose
**Preconditions:** The host creates a game (UC-03)

## Main scenario (happy path)

1. The host opens "Detailed" and sets the options (GS-10).
2. Options: picking (random or draft), draft size, ban rounds, ban timer, host bans, bonus souls.
3. The server checks them against the card sets and the seats.
4. When the table starts, host-banned characters leave the pool (GS-02).
5. Ban rounds go in snake order: 1-2-3-4, then 4-3-2-1 (GS-04).
6. Everyone sees every ban as it happens.
7. Random picking: each player gets one random character from the pool.
8. Draft: each player sees only their own offers and picks one (GS-01, GS-03).
9. The game starts with the chosen characters.

## Alternative scenarios and errors

* **3a. Options that cannot work:** refused with `bad_setup` and the reason.
* **5a. A ban out of turn:** refused with `not_your_turn`.
* **5b. A ban of a card not in the pool:** refused with `bad_card`.
* **5c. The ban timer runs out:** that player makes no ban this round.
* **8a. A pick that is not one of the offers:** refused with `bad_card`.

## Postconditions

* Every player has a character; banned characters were never dealt.
