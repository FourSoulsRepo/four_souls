# Changelog

Each release lists what changed for players.
The release notes take the section of their version; trial releases take Unreleased.
Rules engine versions are in each release's notes.

## Unreleased

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
