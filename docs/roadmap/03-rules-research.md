# Step 3. Rules research

Goal: every rules source downloaded by scripts.
Our own Base Game rules text written with stable IDs.

Ideas: R-01 – R-07.

---

### [x] 3.1 Where research data lives

1. Goal: decide storage before downloading.
2. Tasks:
   1. Raw downloads are copyrighted texts.
   2. Download code lives in `FourSoulsRepo/tools` (no download code here).
   3. Downloads go to its git-ignored `.cache/research/<source>/`.
   4. Digests (Reddit, X, Discord) live there too.
   5. Our own rules text is committed in the engine module.
   6. Notes on the sources: `pkg/rules_engine/docs/sources.md`.
3. Done when:
   1. Folders exist and are documented.

### [x] 3.2 Official rules

1. Goal: foursouls.com rules and FAQ, downloaded.
2. Tasks:
   1. Reuse the browser tool (2.1) for Cloudflare.
   2. Save rules pages and FAQ as clean Markdown.
   3. Record the download date and URL per page.
3. Done when:
   1. Rules and FAQ convert without lost sections.
4. Result:
   1. Tools: `scraper rules -site official` (HTML to Markdown).
   2. Overview, Quickstart (3,000 words), Extended Rulebook (18,500 words).
   3. The FAQ page is only a question form; FAQs live in the rulebook.
   4. Summary in `pkg/rules_engine/docs/sources.md`.

### [x] 3.3 Russian translation

1. Goal: foursouls.ru rules, downloaded.
2. Tasks:
   1. Same pipeline as 3.2.
   2. Note sections that differ from the English rules.
3. Done when:
   1. A diff list of differences exists.
4. Result:
   1. Tools: `scraper rules -site ru`.
   2. Differences: `pkg/rules_engine/docs/sources.md`, section S-RU.
   3. Russian adds card rulings in its FAQ; English points to the site.
   4. One rule differs (co-op timer on death); goes to 3.8.

### [ ] 3.4 Reddit rulings

1. Goal: rules questions from r/FourSouls.
2. Tasks:
   1. Use Reddit's public JSON endpoints, rate-limited.
   2. Collect posts with the rules-question flair.
   3. Keep question, top answers, author, date, link.
   4. Mark answers by Yuggy (Rules Tzar) as high trust.
3. Done when:
   1. A digest file lists questions with links.

### [ ] 3.5 Twitter / X rulings

1. Goal: rulings by Ed, Yuggy, Kizzycocoa.
2. Tasks:
   1. X blocks plain scraping; use the browser tool, logged in.
   2. Start from known tweet links (R-07).
   3. Search each account for ruling threads.
   4. Keep text, date, link.
3. Done when:
   1. All R-07 tweets are saved with text.

### [ ] 3.6 Discord rulings

1. Goal: knowledge from "Турнирный сервер".
2. Tasks:
   1. No bot scraping (Discord terms).
   2. Owner exports or pastes relevant threads.
   3. Agent turns them into the same digest format.
3. Done when:
   1. Pasted threads are in the digest.
4. **Owner:** provide the threads.

### [x] 3.7 Our rules text

1. Goal: our Base Game rules inside the engine (R-04).
2. Tasks:
   1. Folder: `pkg/rules_engine/docs/rules/`.
   2. One file per topic: setup, turn, stack, combat, death, shop, win.
   3. Each rule has a stable ID, e.g. `R-COMBAT-03`.
   4. Each rule lists its sources.
   5. Written in our own words.
3. Done when:
   1. Every Base Game mechanic has a rule.
   2. IDs are unique; a test checks this.
4. Result:
   1. 237 rules in 14 topics (SETUP, CARD, ZONE, ABIL, MECH, TURN, PRIO).
   2. Also STACK, DICE, ATK, SHOP, DEATH, WIN, BARTER.
   3. Written from the official Extended Rulebook, checked against S-RU.
   4. `rules_doc_test.go`: ids well formed, unique, gap-free, refs resolve.
   5. Found: no mulligan rule (R-06); bartering of ¢ needs game support.
   6. Found: official variants Mini-draft (2 players) and Eden Only.

### [ ] 3.8 Open questions

1. Goal: list what the sources do not settle (R-03).
2. Tasks:
   1. One file of open cases: question, sources, options.
   2. Include mulligan (R-06) and first player (GS-06).
   3. Include each old ruling from R-07 to re-check.
   4. Include team play rules for 2 vs 2 (GS-05).
3. Done when:
   1. **Owner** has resolved or parked each case.
   2. Resolved cases wait for the payload format (4.11).
