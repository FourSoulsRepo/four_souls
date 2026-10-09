# Turns

Sources: S-OFF Turn Structure, Round; S-RU Структура хода.

## Start phase

1. **R-TURN-01** A turn has three phases: start, action, end. Turns go clockwise (to the left).
2. **R-TURN-02** Recharge step: the active player recharges everything they control; the game's objects recharge too, except shop items. Nobody has priority yet.
3. **R-TURN-03** Then start-of-turn triggers trigger and priority passes around.
4. **R-TURN-04** Loot step: the active player loots 1, then priority passes around and the action phase begins.

## Action phase

5. **R-TURN-05** At the start of the action phase the active player gets one loot play, usable until the end of the turn.
6. **R-TURN-06** With an empty stack, the active player has priority and may declare an attack, declare a purchase, or end the turn.
7. **R-TURN-07** By default the active player may attack once and purchase once per turn.
8. **R-TURN-08** While holding priority, the active player may also use activated abilities and play loot cards; doing so, or declaring an attack, purchase or end, passes priority.
9. **R-TURN-09** The action phase ends only when the active player ends the turn or an ability ends it; the empty-stack rule of R-STACK-07 does not move it on.

## End phase

10. **R-TURN-10** End-of-turn triggers trigger, then priority passes around.
11. **R-TURN-11** The active player discards down to their max hand size (10 by default). Nobody has priority from here on.
12. **R-TURN-12** *(later sets)* If a monster died this turn, the active player may put a room into discard; an empty room slot is refilled.
13. **R-TURN-13** The turn passes to the next player. Every object with HP heals to full, dead players included, and "till end of turn" effects end.
14. **R-TURN-14** "End the turn" effects start the end phase but skip none of its steps and remove nothing from the stack.

## Rounds

15. **R-TURN-15** A round runs from the start of a player's turn to the start of their next turn. For the game's objects it is counted from the starting player.
16. **R-TURN-16** Extra turns and skipped turns do not change where a round starts or ends; a skipped turn still counts as taken.
