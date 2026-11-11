# 28. Challenges

* **Status:** Accepted
* **Date:** 2026-11-10
* **Authors:** @HardDie

---

## Context

1. Owner, 2026-11-10: a "Challenge" game mode beside the regular game.
   1. First challenge: How the Krampus Stole Christmas.
   2. Difficulties: Normal, Hard, Ultra Hard.
   3. Co-op only: exactly 2 people, one character each.
   4. Solo (14.27) and the competitive variant (14.28) come later.
2. Rules: `pkg/rules_engine/docs/rules/challenges.md` (14.17).
   1. COOP: co-op setup, timer, shared win and loss.
   2. CHAL: final boss slot, boss armor, minion slots.
   3. KRAMPUS: presents, three global modifiers, per difficulty.
   4. Open cases Q-30 – Q-42 hold the proposals the engine uses.
3. Official challenges come with their own cards.
   1. foursouls.com lists them as origin `challenge` (83 cards).
   2. Krampus: 4 Pilfering Krampus boss cards, one per difficulty.
   3. Krampus: 4 rule cards with the global modifiers.
   4. The boss card holds the per-difficulty rules.
4. The engine knows card sets, decks and three slot rows.
   1. Shop, monster and room slots (`SlotKind`).
   2. Card abilities run only for objects `inPlay` lists.
   3. Saved games keep set names; definitions come back on load.
5. One app holds one seat; the view hides other hands.

## Considered options

1. **A challenge as a card set plus a small definition**
   1. Won. The boss and rule cards are `CardDef`s with abilities.
   2. Load needs no new lookup: the set is in `Game.Sets`.
2. **A challenge as Go code with hooks into the turn**
   1. Lost. A second ability system beside the cards.
3. **Boss and minions inside `Monsters`, with a flag per slot**
   1. Lost. "Monster slot" effects would need to skip them everywhere.
4. **Own slot rows for the final boss and the minions**
   1. Won. Slot choices stay on monster slots by default.
   2. "Monsters in play" goes through one helper over all three rows.
5. **Presents as treasure cards outside the game, with counters**
   1. Won. They are not in play (R-KRAMPUS-06); a list names them.
6. **Presents as items of nobody**
   1. Lost. Item abilities would work; R-KRAMPUS-06 says they don't.

## Decision

Use options 1, 4 and 5.

1. Card data
   1. card_db set `challenge`; only cards of supported challenges.
   2. Krampus: `pilfering_krampus_{norm,hard,ultra,comp}`.
   3. Its rule cards: `how_the_krampus_stole_christmas_{norm,hard,ultra,comp}`.
   4. The set is never offered as a playable set.
   5. The app's set list, the host screen and `checkOptions` skip it.
   6. Images: private repo, folder `challenge/`; without them, text cards.
   7. IDs: the site's slugs, a hyphen as `_` (ADR 004); fixed in 14.19.
   8. `comp` is the competitive variant (14.28): data only for now.
   9. The site prints no type; the tools name them.
      1. `Challenge Boss Card`: the boss, a card with stats.
      2. `Challenge Card`: a rule card.
   10. Tools: `gameCards` lists the cards the game takes; one copy each.
   11. Set order 5, released 2022-12 (ADR 010).
   12. Images are landscape, 839×600.
   13. Trick or Treat: `mischievous_the_haunt_{norm,hard,ultra,comp}`.
   14. Its rule card: `trick_treat_{normhardultra,comp}`, two-sided (ADR 023).
       1. Trick is the front; the backs are `treat_normhardultra`, `treat_comp`.
       2. `normhardultra` serves all three co-op difficulties.
2. Engine kinds
   1. A boss is a `MonsterCard` with `Outside`: never in a deck.
   2. A rule card is a new kind `ChallengeCard`: no deck, never dealt.
3. Engine definition
   1. `ChallengeDef`: ID, difficulties, and per difficulty the boss and rule card.
   2. Also a `Setup func(g, difficulty)` for the challenge's own setup.
   3. Krampus's setup: the boss into a final boss slot, the 4 presents.
   4. Definitions live in `cards/challenge`, beside the sets.
   5. `cards.Challenges()` and `cards.FindChallenge(id)` list them.
   6. `ChallengeDef.Set` is the code of its card set.
   7. `ChallengeDef.Levels`: difficulty, boss and rule card each.
   8. Krampus's ID is `krampus`.
4. Engine setup
   1. `Setup.Challenge *ChallengeDef` and `Setup.Difficulty string`.
   2. Difficulties: `normal`, `hard`, `ultra`.
   3. The challenge set must be among `Setup.Sets`; the caller adds it.
   4. A challenge needs 2 players, no bonus souls (R-COOP-01, R-COOP-02).
   5. 3+ player cards are out, as for any 2-player game (R-COOP-03).
5. Engine state: `Game.Challenge *ChallengeState`, nil in a regular game.
   1. `ID`, `Difficulty`, `Timer` (starts at 8).
   2. `Card`: the rule card's object; its abilities work.
   3. `Presents`: kind (`cents`, `loot`, `treasure`, `soul`) and object.
   4. A present object sits outside the game, face down, with its counters.
   5. Destroyed: to the treasure discard; flipped: an item of a player.
   6. `Holidays`: Happy Holidays! has happened (R-KRAMPUS-10).
   7. `Lost` and `Reason`: why the game ended, as a message key.
   8. New fields are `omitempty`; regular saves stay byte-identical (ADR 019).
   9. `Starter`: rounds count from the starting player (Q-34).
   10. `TurnFrom`, `RoundFrom`: first object ID of each player's last two turns (Q-39).
   11. `RoomTurn`: the turn whose room rule has run (R-COOP-11).
   12. `Present.Destroyed`, `Present.Flipped`: where a present went (14.22).
   13. Presents API: `AddPresent`, `AddPresentCounters`, `DestroyPresent`.
   14. Also `FlipPresents`, `GivePresent`, `Forgo`, `SoulsInPlay`.
6. Abilities
   1. `inPlay` adds the rule card after every slot row.
   2. Boss rules are the boss card's abilities (R-KRAMPUS-11 – 14).
   3. Global modifiers are the rule card's abilities (R-KRAMPUS-07, 08, 10).
   4. New hooks go on `CardDef`, like `InsteadOfDestroy`; no hooks by challenge ID.
   5. Engine calls: `WinChallenge()`, `LoseChallenge(reason)`.
   6. Hooks from 14.22: `Guarded`, `MustAttackRandom`, `NoSpoils`.
   7. Also `Replacement.May`: an optional replacement, asked yes or no.
   8. R-KRAMPUS-14 is the boss card's: the boss card prints it.
   9. The goal (R-KRAMPUS-01) is a rule card ability.
   10. A rule card's triggered abilities are the active player's.
7. Slot rows
   1. New `SlotKind`s: `FinalBossSlot`, `MinionSlot`.
   2. `Game.Bosses` and `Game.Minions`, each `[]Slot`.
   3. One helper lists monsters in play over monsters, bosses, minions.
   4. Attack targets, odds and "choose a monster" use that helper.
   5. "Monster slot" choices (cover, refill, expand) stay on `Monsters`.
   6. A minion slot is removed when its card leaves (R-CHAL-11).
   7. Later minions' `Zone.Index` shifts down with it.
   8. Boss armor is one rule in the damage and kill code (R-CHAL-05).
   9. Helpers `MonstersInPlay()`, `MonsterRows()`: monsters, bosses, minions.
   10. `Held(id)`: a final boss or minion; discard and deck moves skip it.
   11. Engine calls: `PutInFinalBossSlot(card)`, `PutInMinionSlot(id)`, `MinionFromDeck()`.
   12. A dead final boss leaves its slot empty; nothing refills it.
   13. Situations take `bosses` and `minions` rows too.
8. Co-op turn and end
   1. Timer ticks at the end of each round (R-COOP-08, Q-34).
   2. Timer at 0: `LoseChallenge` (R-COOP-09).
   3. The soul target and "you win" effects do nothing (R-CHAL-14).
   4. A win: `Over`, `Winners` both players, event `EvChallengeWon`.
   5. A loss: `Over`, no winners, `Lost`, event `EvChallengeLost`.
   6. Rooms: the co-op room rule after the room step (R-COOP-11).
9. View
   1. `GameView.Challenge`: ID, difficulty, timer, rule card, presents, end reason.
   2. A face-down present shows its kind and counters, never its card.
   3. `GameView.Bosses` and `GameView.Minions`: `[]SlotView`.
   4. In a challenge, each player sees the partner's hand (R-COOP-07).
   5. A present also shows destroyed, flipped, its face-up card, its gainer.
10. Protocol, server, records
   1. `Options.Challenge` (ID) and `Options.Difficulty`.
   2. A challenge table has exactly 2 seats; the server checks it.
   3. The server adds the `challenge` set to the game's sets.
   4. It turns bonus souls off in a challenge.
   5. Record setup: `challenge` and `difficulty` (ADR 007).
   6. Saves and the seamless restart keep them (ADR 019).
   7. The host screen lists the challenges from the app (`App.Challenges`), like card sets.
   8. Each: ID and difficulties; the client names them (`challenge.name.<id>`).
   9. A server that lacks one refuses it at create: `error.no_challenge` (14.23).
   10. Other refusals: `error.no_difficulty`, `error.challenge_seats`; code `bad_setup`.
   11. Table, game list, saves and records keep only the playable set codes.
   12. The challenge in the options brings its set back (`gameSets`, last).
   13. `GameInfo` carries `challenge` and `difficulty`.
   14. A host ban of a challenge card is refused.
   15. A lost challenge is a finished game: record `end` finished, no winners.
11. Client
   1. Host screen: "Game mode" select; then challenge and difficulty selects.
   2. Table: the final boss left of the monster slots.
   3. Minions in a row under the monsters; the row hides when empty.
   4. Presents, timer and rule card in a challenge panel by the decks.
   5. A face-down present shows the treasure back; plain without images.
   6. Game over: "Challenge won" or "Challenge lost" with the reason.
   7. The hidden minion row keeps its place: the board keeps one shape (V-23).
   8. The challenge panel takes the bonus souls' column; a challenge has none.
   9. Boss and rule cards are landscape, by card type (`landscape` in `cards.ts`).
   10. A destroyed present's place shows ✗; a given one's ✓ (14.25).
   11. The timer is a d8 with "Rounds left"; its tooltip says how many.
12. Solo, later (14.27)
   1. Nothing here may assume one seat per app forever.
   2. The engine already treats each character as a player.

## Consequences

### Positive

1. A new challenge is a few cards plus a setup function.
2. Saves, records and replays reuse the set mechanism.
3. Regular games do not change: every new field is empty there.

### Negative and risks

1. Every "monster" lookup must go through the new helper.
   1. A missed one ignores bosses or minions; 14.21 audits all sets.
2. Removing minion slots shifts indices; tests must cover it.
3. Showing the partner's hand changes the view filter for challenges.

### Neutral

1. The 69 other challenge cards stay out of card_db until their challenge.
2. Plain co-op without a challenge is not offered (R-COOP-12).
