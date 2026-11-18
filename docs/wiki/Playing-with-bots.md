# Playing with bots

Not enough friends for a game? The host seats bots on the free seats. One player and three bots make a full table. Bots play in any network game: on your own computer's server (**Play → Host**) or at a table on a server (**Play → Online**, then **Create**).

## Add bots in the lobby

After you host, the table screen lists the seats. As the host:

1. Each free seat shows **Add bot**, with one button per level: **Standard** and **Hard**.
2. Press a level. The seat now shows a robot, **Bot 1** (Bot 2, Bot 3 for the next ones) and a level menu. A bot is always ready.
3. To change its level, pick another one in its menu. To take it away and free the seat, press its **×**.
4. Mark yourself ready. When every seat is taken and every player is ready, the game starts, even if you are the only person at the table.

The other players at the table see each bot with its robot, "Bot 1 (Standard)" and Ready, and no bot buttons: only the host seats and removes bots.

A friend who joins later takes a free seat, never a bot's. To make room for a friend, remove a bot before the game starts. The game list counts bots as players.

No **Add bot** on a free seat? Then:

* you are not the host;
* the server allows no bots: a server's operator chooses the levels (see [Dedicated server on a VPS](Dedicated-server-on-a-VPS)), and a server that allows only Standard shows one button; the app's own server allows both;
* the table continues a saved game: its seats wait for its players.

## Standard or hard

* **Standard** plays its turn well enough: it buys, attacks when the odds are good and plays loot that helps it now. It picks a move in a few milliseconds.
* **Hard** thinks ahead. It plays many short games in its head before each move, remembers cards it saw, and keeps a card that can answer for when the leading player needs stopping. It thinks for up to about 1.5 seconds a move; a busy server may make it play a move as a standard bot.

A bot sees only what a player in its seat sees: never another hand or the order of a deck. You can mix levels at one table.

Bots wait a moment before each move so you can follow them: about 0.7 seconds for a pass or an answer, 1 to 2 seconds for an action on their own turn.

## At the table

A bot's name is **Bot N (level)** after a robot icon, wherever the game names a player: its area on the table, the turn badge, the history, the stack, the banners, the vote window and the game end.

Each level has its own robot: a Standard bot's is friendly, with round eyes and an antenna; a Hard bot's has horns, angry eyes and a jagged mouth. A replay shows the bots the same way.

* A bot is never away or disconnected, and the game never pauses for it.
* Bots do not vote on a kick or a Soft rollback.
* A bot stays in its seat through a save, a server restart and a match record.

## Challenges

A bot can be your partner in a challenge: seat it on the second seat. It plays for the team and wins or loses with you. See [Challenges](Challenges).
