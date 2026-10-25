# Changelog

Each release lists what changed for players.
The release notes take the section of their version; trial releases take Unreleased.
Rules engine versions are in each release's notes.

## Unreleased

- Online: each Recent server shows whether it is up and how many games there have a free seat.
## 0.2.1 - 2026-10-22

The first public release.
0.2.0 was never published: its builds without card images failed.

- If the app closed during a game, the main menu offers "Rejoin" while that game is still on. Leaving a game on purpose, or its end, forgets it.
- Online: type a server's address or pick it from Recent; it no longer lists LAN games. Back and Esc go one screen up.
- LAN: Host a game sits on top; games on your network show up by themselves, without "Search again". A game in progress stays listed, and if your app closed during it, Rejoin takes your seat back. The address field remembers the last address you typed.
- Main menu: Play opens LAN, Online and Learn to play right under it; Records and Learn to play are marked "Soon".
- The board glows green whenever the game waits for you, on your turn or out of it; an hourglass by a name shows everyone whom the game waits for.
- A resolved stack row no longer shows again for a moment while the rows slide.
- A changed health, dice or attack is green when it went up and red when it went down; hover the card to see the sum, such as "ATK 3 = 1 printed + 1 Brimstone + 1 XI. Strength (this turn)".
- Effects of this turn show as small icons on the character or monster; hover one to see its card and what it does.
- An aimed attack says what the roll needs and what a hit deals: "Hits on 4+ (DC 5 − 1 Meat!) · ATK 1"; a thrown roll names what changes it, "(+1 Meat! = 2)", and so does the History.
- A bonus to the next attack roll shows a small die on its icon; the icon pops when the roll uses it, and this turn's icons shrink away when the turn ends.
- The app is now named FourSouls: `FourSouls.app`, `FourSouls.exe`, and `FourSouls` on Linux. Settings carry over.
- Stack and history strip rows no longer jump: new rows fade in, resolved ones fade out, and the rest slide into place.
- The stack panel shows an open attack or buy on top, "Bo attacks Pooter" or "Bo buys…", until it ends.
- Hovering a history strip tile shows what the play did: its target, the damage, the dead card with a skull, and the rewards.
- The history strip shows each buy with a coin and each kill with a skull, on the card that was bought or killed.
- The History says what damage, deaths and rolls did, "Bo deals 2 damage to Gaper", instead of "Roll — damage".
- The History names every play: who attacked what, who bought which item, who killed which monster, souls, thefts, cents, deaths, and effects that failed.
- The "Skip all" button is gone.
- "Wait for my turn" stops as soon as another player targets you or your cards, so you can answer before they steal from you.
- It also stops for cards that pick their victim only as they resolve, such as XV. The Devil, Joker and Dead Bird.
- Wording: the turn counters say "Shop" instead of "Buys"; roll texts say "dice", so "die" only means death.
- Dice: a thrown die is a 3D cube that tumbles onto the table and clatters as it lands; every die has light and shade.
- End turn asks first while you can still play something, and lists what is left; Settings can turn the question off.
- A dedicated server can have a name (`-name`); players see it instead of its address.
- Join a game: remove a server from Recent with its × button.
- A revealed card grows out of its deck; a new monster or shop card then flies from the middle into its slot, and any other card goes back onto its deck.
- A played loot card rises from its player's hand to the stack, and flies from the stack to the discard when it resolves.
- Stolen cents fly as a coin from the player robbed to the thief.
- A changed roll shows the old value, a shot from the card that changed it, then the new value.
- Damage from a card that already left the table flies from the player who dealt it.
- The main menu shows the Four Souls logo of foursouls.ru at the top.
- The menu screens drift on the starry sky of foursouls.ru.
- A match record names the dedicated server the game ran on.
- The main menu has an Exit button.
- The light theme is gone; the app has one dark theme, and Settings no longer has a Theme field.
- The Green felt game mat is gone; a table that used it gets the default mat.
- A thrown die rests with its result on top, seen from above, with two sides still showing.
- Before turn 1 every player throws a die for the first player; everyone sees each die, ties roll again, and a banner names who goes first. Cain still goes first without a roll.
- Records from older versions no longer replay: the first-player roll is now part of every record.
- Dedicated server: every option can also be set by an ENV variable, e.g. `FOUR_SOULS_PORT`.
- The dedicated server has its own downloads, `four-souls-server-<system>-<cpu>`; the game archives no longer include it. Its Linux build runs on any distribution, Alpine included.
- Join a game: type just a domain, e.g. `games.example.com`; the app tries a secure connection (`wss://`) first, then the plain port 4774. A typed `wss://` or `ws://` is kept.
- Dedicated server in Docker: `deploy/docker/Dockerfile` builds an Alpine image from the latest release; `/healthz` tells Docker the server is healthy.
- Dedicated server with HTTPS: `deploy/docker/docker-compose.yml` runs it behind Caddy, which gets the certificate by itself; players type only the domain.
- Dedicated server guides: pick a systemd service, Docker, or Docker Compose with Caddy; each says how to update.
- Each release also has "for-poor" archives: the game without card images, with cards shown as text.
- Network protocol 4: older apps can't join newer games, and the other way round.
- Rules engine 0.3.0: match records and saved games of older versions can't be replayed or loaded.

## 0.1.1 - 2026-10-19

- Rules: Sacred Heart and similar "would roll" items no longer change a roll that was already changed.
- Rules: an item swapped away (Decoy) can still be destroyed by an ability aimed at it.
- Rules: abilities that choose an item can choose shop items (The D20, Placebo, Glass Cannon).
- Rules: a player who already died this turn can't be chosen while a living player can (Yuggy's ruling).
- Rules: a dead player still takes damage for "when you take damage" items, without losing HP again.
- Rules: cards with the 3+ player symbol are left out of 2-player games.
- Several copies of the app can run on one computer; each one is its own player.
- "Put them back in any order" (Sleight of Hand, The Compass, The Map): one dialog shows every card's place, top first; confirm when done.
- Dice: every roll waits for its player to press Throw (or Space); the others see who is throwing.
- Dice: a thrown die falls onto the table, bounces, and tumbles to its result.
- Dice: the last roll stays as a small die by the stack; hover it to see who rolled and what for.
- Dice: everyone sees a roll in the roller's own dice look.
- Dice: Settings has ten dice looks, Holy and Unholy, with a preview of every face.
- Souls: gained soul cards hide behind the soul crystals; hover a crystal (or long-tap it) to see its card.
- Curses sit at the far left of their player, apart from the items.
- Network protocol 3: older apps can't join newer games, and the other way round.
- Rules engine 0.2.0: match records and saved games of older versions can't be replayed or loaded.

## 0.1.0 - 2026-10-12

The first release (MVP): the Base Game for 2–4 players over a LAN.

- Network play for 2–4 players: host a game in the app, or run the dedicated server.
- The Base Game: every card, with the official rules and FAQ.
- A Hearthstone-like table: cards fly, dice roll, a history and a stack panel.
- Drag and drop or tap to play; works with mouse and touch.
- Pause when a player drops; vote to wait or go on; reconnect to the same seat.
- Save a game and continue it later.
- Match records to keep.
- Settings: nickname, theme, game mat, animation speed, fullscreen.
