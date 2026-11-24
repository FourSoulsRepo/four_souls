# Changelog

Each release lists what changed for players.
The release notes take the section of their version; trial releases take Unreleased.
Rules engine versions are in each release's notes.

## Unreleased

- Lord of the Pit, The Bible's "Cancel your attack.", IBS's "Cancel your attack on a monster." and Emergency Contact need an attack to cancel; with none they are not offered.
- ? Card offers only the ↷ abilities it could use now: it no longer copies Lord of the Pit with no attack to cancel, nor an ability with no valid target; with none it can't be played.
- Boss Rush!'s extra attack goes only at one of the two bosses it revealed; in one slot only the top boss. Friendly Ball's and Lust for Blood!'s extra attack likewise only at their monster or player. If the top boss of one slot dies first, the boss under it is still one of them. Once none of them can be attacked, the extra attack is no longer owed. The attack button and the red frames name them: "You must attack: Monstro or Widow".
- Red Key during an attack and Friendly Ball's first ability with a card under it now say the card forbids it (R-ABIL-08), not an unpaid cost. Friendly Ball's second ability with nothing under it has no valid target, and ? Card no longer offers to copy it.
- Rules engine 0.5.11.

## 0.6.4 - 2026-11-20

- A character that takes more damage (Ghostly Trickery, Daddy Haunt) shows a red "+N" badge above its heart; the zoomed card lists where it comes from.
- The game menu and its settings now cover big cards and dice; they no longer show on top.
- Linux: `--diagnose` checks the GStreamer plugins for sound; without the good plugins the game hangs at its first sound. The Linux packages to install now include `gst-plugins-good` and `gst-libav` (see README.txt).

## 0.6.3 - 2026-11-19

- Linux: the game always draws with the GPU; on a tested laptop the menu went from 6 to 60 frames a second. If the window stays black, start it with `FOURSOULS_GPU_POLICY=ondemand` (or `never`).
- `--diagnose` no longer names every GPU "Apple GPU" on Linux and macOS; WebKit hides the name, and the report now says so.

## 0.6.2 - 2026-11-18

- Diagnostics: start the app with `--diagnose` to see how fast it draws and what slows it. On Linux it tries each known WebKitGTK switch, names the fastest, and finds problems such as no GPU access or no MP3 sound. The report is also saved to a file; see the wiki page Diagnostics.

## 0.6.1 - 2026-11-18

- Boil gives 3¢, as printed, instead of 3 loot cards.
- A choice of a dice value shows dice, not buttons (God Head, Crystal Ball, I - The Magician, Red Candle, picking one of many dice).
- On a dedicated server, a new game after leaving one no longer misses its Throw and other buttons.
- A refused forced attack, and a soul that lost its abilities, now name the rule that applies.
- Rules engine 0.5.6.

## 0.6.0 - 2026-11-18

- Bots: the host seats a bot on any free seat, standard or hard (lobby: add bots per seat); one player and three bots make a full table. A bot sees only what a player in its seat sees.
- Bots play challenges as your partner, stay through saves, server restarts and match records, and never vote on a rollback or a kick.
- Dedicated server: `-bots` (`off`, `standard`, `all`) sets the levels a host may seat; `-bot_think` (100 ms to 10 s, default 1.5 s) a hard bot's time to think a move. At most half the CPU cores think for hard bots at once; a hard bot that waits plays that move as a standard one.
- Rules engine 0.5.5: records and saves from 0.5.4 and older do not replay or load (decks now tell where their cards lie).

## 0.5.1 - 2026-11-17

- A card that asks for a deck shows the decks as card backs, with their names, instead of buttons; a discard pile shows its top card (e.g. The Curse, Sleight of Hand).

## 0.5.0 - 2026-11-17

- Records: a Records screen in the main menu lists hosted games and saved records, newest first; open a file, export or delete one.
- Replays: Watch plays a record on the game table, with play and pause, steps, turns, a slider and speeds from 0.5x to 4x.
- In a replay, Show hands shows every hand or none; click a player's character card to see the table from their seat; the replay bar hides with a button.
- A record's cards that this app lacks show as text cards.
- The card zoom credits the card's artists.
- Host screen: the Detailed setup options sit in three groups: Characters and seats, Cards, Turns and undo.
- Fear shows on the character card.
- A monster you must attack, and the monster deck while an attack on it is owed, get a dashed frame; the banner says how many attacks are left.
- Who helps the next attack, and who shares its rewards, shows for everyone (Sibling Rivalry).
- A roll's result on the stack shows the number; a seat kicked by a vote says "Left for good".
- A finished game leaves the server's game list at once.
- The history says when a deck is shuffled.
- The history shows the ¢ that Greed! and Auction Gavel move; Greed! no longer takes ¢ from a player with Rib of Greed.
- A finished event card, such as a chest or Troll Bombs, flies to the monster discard instead of vanishing.
- Rules engine 0.5.4: games saved with an older version cannot be continued.
- A game saved, undone or kept over a server restart during a card's second question, such as Gimpy's "loot 1, then discard", goes on with that question; before, the card asked its first question again, and some cards stopped the game.
- A dedicated server updated to this version ends the games it kept over the restart; their players are told so.

## 0.4.3 - 2026-11-15

- An Update button in the main menu: it checks for a new version, downloads it with a progress bar, and installs it on restart. You keep your system and edition (with card images or for-poor). This is the last version you download by hand.

## 0.4.2 - 2026-11-14

- Russian card images for the Base Game: with Russian chosen, its cards show in Russian; other sets stay in English.
- Trick or Treat: a monster that gains another monster's abilities shows a copy mark for each one; hover it to see that monster.

## 0.4.1 - 2026-11-11

- Challenges: a co-op game mode for 2 players, picked in the host screen, with Normal, Hard and Ultra Hard.
- How the Krampus Stole Christmas: beat Pilfering Krampus and his minions before all four presents are destroyed or the timer runs out.
- Trick or Treat: beat Mischievous The Haunt before the timer runs out, while its tricks and your fear make it harder.
- Halloween and Christmas game mats, to play the challenges in the mood.
- Requiem rooms: the host can turn on the room deck in Detailed setup; the active player uses the room, and the end phase may change it.
- Attack a player or an item: The Siren, Punching Bag and other Requiem cards offer them as attack targets; extra attack rolls name whose roll it is.
- Rooms and other landscape cards are as sharp as portrait cards.
- Golden Trinket and Glitch with a copied ↷ ability count as ↷ items: Charged Baby and other recharge cards can recharge them, and death or Curse of Impulse deactivates them.
- An attack on Pride, Rag Mega, Dogma and other "when an attack is declared on this" monsters: their ability resolves before the attack die is thrown.
- An attack on a player (The Siren): the attack roll names the player, not their character card.
- Rules engine 0.5.3.

## 0.4.0 - 2026-11-09

- Draft: in Detailed setup the host can offer each player 2 to 5 characters to pick from; one random character each stays the default. Each offer shows the character's starting item; Eden chooses hers when the game starts.
- Start of game: a screen for start-of-game choices, such as Eden's item; several players choose at the same time, and the first-player roll waits until all are done.
- Interface scale: a setting from 75 to 200 %; the whole window grows or shrinks with it.
- Settings open from the game Menu too, without leaving the game.
- The game feels livelier: fast strikes, hits that shake by their size, a damage number that stays, drag tilt, and new sounds, one per moment.
- Others' plays are paced so you can follow them; their hands update live, and the table catches up by importance after a pause.
- Choosing cards, playing a card from the hand, and the target arrow are reworked; only the attacker flies to its target.
- The History strip slides new tiles in and holds them; it keeps the newest 40, so long games stay fast.
- The response timer shows only when a response is possible; End turn shows its states.
- A look at a deck shows on the board whether each card went to the top or the bottom.
- Glitch's copy mark shows and zooms the card it copies.
- Recharging or deactivating an item offers only the items the effect allows; Holy Mom's Eye no longer offers passive items.
- Soul of Envy no longer goes on the stack again for the player who gains it.
- Headless Horseman: its prevented death no longer heals it.
- Online: a game waiting for players on a server no longer closes after about a minute.
- LAN: Join takes you to the host's game instead of the game list.
- Rules website: the rules and card checks run in the browser, at <https://foursoulsrepo.github.io/four_souls/>.
- Rules engine 0.5.2: saved games of older versions can't be loaded.

## 0.3.3 - 2026-10-31

- Dedicated server: one CPU serves more games at once, each step costs less work.
- Dedicated server: a slow connection is cut off less often; the queue of each player holds 256 messages, up from 64.
- Dedicated server: the match records are written once a second instead of at every step; they are about half the size.
- Dedicated server: the log shows the CPUs it runs on, and any request that holds a game for more than 100 ms.

## 0.3.2 - 2026-10-29

- Security: the app is built with Go 1.27.3, which fixes security problems found in Go 1.27.2.
- The app's libraries are updated to their latest versions.

## 0.3.1 - 2026-10-29

- The host can seat the players in a random order: Detailed setup, Random seats. The seats are shuffled when the game starts; without it, seats follow the order players joined.
- A dead monster's slot refills after its reward roll, as the rules say. A monster that enters on that roll, such as Holy Dinga, no longer reacts to it.
- An attack on a monster from the deck: its "when an attack is declared" ability goes on the stack before the dice are thrown.
- Rules engine 0.4.2: match records of 0.4.1 and older can't be replayed.

## 0.3.0 - 2026-10-27

- The host can end a game for everyone from the game Menu: End game, then confirm. It works on a dedicated server too, where leaving only paused the game.
- In a game with rollback, one Undo button under the History title undoes the latest move, press by press; its tooltip names it, "Undo: Bo plays A Penny!". A played card and the target it asks are one move. Greyed, it says why, such as "Nothing to undo since a die was rolled".
- "Choose one-" cards such as Dagaz ask like The Bone: the card large, one line per option; an option with nothing to target is greyed.
- Online: each Recent server shows whether it is up and how many games there have a free seat.
- A look at other players' hands, such as XXI. The World, shows each hand in its own row under its owner's name.
- When a monster that covered another dies, the card below shows up in its slot instead of flying in from the monster deck.
- A dead monster flies from its slot to the monster discard, or to the player who gains its soul, after its skull, instead of just appearing there.
- A declared end of turn shows on top of the stack, "Ann ends the turn", like an attack or a buy, so the others see why they may respond.
- With nothing left to do, End turn breathes a soft green glow instead of showing a green frame.
- What the server says (reasons, card questions and options, roll labels, errors) shows in each player's own language; two players at one table may use two languages. A player option names the player instead of "player 2 (isaac)".
- Network protocol 6: older apps can't join newer games, and the other way round.
- Rules engine 0.4.0: match records and saved games of older versions can't be replayed or loaded.
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
