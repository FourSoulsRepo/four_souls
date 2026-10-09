# Rules sources

Where our rules text (R-04) comes from, and how the sources differ.
Raw copies are copyrighted; they stay in the tools repo cache.

## Where the raw copies live

1. Downloaded by `FourSoulsRepo/tools` (`scraper rules`).
2. Stored in its git-ignored `.cache/research/<source>/`.
3. Each file starts with its source URL and download date.
4. Nothing from them is copied here word for word.

## Sources

| ID | Source | Pages | Trust |
|----|--------|-------|-------|
| S-OFF | foursouls.com/rules | Overview, Quickstart, Extended Rulebook | Highest |
| S-RU | foursouls.ru/rules | Russian mirror of the same three pages | High; see differences |
| S-RED | r/FourSouls | Rules questions; Yuggy answers rank high | Medium |
| S-X | Tweets by Edmund McMillen, Yuggy, Kizzycocoa | Rulings | High for Ed and Yuggy |
| S-DIS | Discord "Турнирный сервер" | Threads pasted by the owner | Medium |

Order of trust when sources disagree:

1. Official Extended Rulebook and its Errata.
2. Edmund McMillen and Yuggy (Rules Tzar) rulings.
3. Russian translation.
4. Reddit and Discord discussions.
5. Anything unresolved becomes an open case for the owner (R-03).

## Official site (S-OFF)

1. Downloaded 2026-10-09.
2. Overview: about 200 words.
3. Quickstart Guide: about 3,000 words.
4. Extended Rulebook: about 18,500 words in 135 sections.
5. Main sections:
   1. Setup, Card Types, Game Zones, Abilities, Effects.
   2. Specific Mechanics (largest: 4,700 words, A–Z keywords).
   3. Turn Structure, Attacking, Purchasing, Refilling Slots.
   4. Death, Bartering, Dice Rolls, Priority, The Stack, Winning.
   5. Solitaire/Co-op, Challenges, Deck Ratios, Variants.
   6. Errata, Key Terms, FAQ, Credits.
6. The FAQ page holds only a question form.
7. Card pages carry no FAQs for the Base Game V2 cards.

## Reddit (S-RED)

1. Downloaded 2026-10-09 with a logged-in session (public JSON is blocked).
2. Flair "Gameplay Question"; Reddit search returns the newest 249.
3. Covers 2025-09-13 to 2026-10-09.
4. Answers come from regular community members; none by Yuggy.
5. Trust: medium; use as hints and for open cases, not as rulings.

## X (S-X)

1. Downloaded 2026-10-09 with a logged-in session; pages rendered in Chrome.
2. The 9 tweets of R-07, each with its thread; the ruling is marked.
3. They date from 2018–2019, before the second edition.
4. Re-checked against current rules in open-questions.md.

## Russian translation (S-RU)

1. Downloaded 2026-10-09.
2. Same three pages; 136 sections against 135.
3. Most sections are full translations.
   1. Russian runs about 75% of the English word count.
4. Structure differences:
   1. Anatomy of a Card has four extra parts in Russian.
      1. Challenge difficulty symbol, original card name.
      2. Artist, card number.
   2. Only in English: "Additional Names".
   3. Only in English: counter subsections ({HP}, {ATK}, Eternal).
   4. Abilities: Russian lists Loot Abilities before Activated.
5. Content differences:
   1. Attacking: Russian splits rules into bullets with card examples.
      1. Examples: Bob's Brain, Pride, Mom's Shadows, Ambush!.
      2. The rules themselves match the English text.
   2. FAQ: English only points to the website.
      1. Russian holds card rulings (Two of Clubs, Ambush!, Keeper's Sack).
      2. Also: multi-dice choices, Isaac's Tears with coin flips.
      3. These are useful extra rulings; check each against English.
   3. Solitaire/Co-op: English starts with an older duplicate summary.
      1. Russian drops it; the subsections match.
      2. One rule differs: the timer drops on a death.
      3. English: only a death "during their turn"; Russian: any death.
      4. Listed as an open case (3.8).
   4. Errata: identical (left in English).
