# Step 16. Localization

Goal: menus, messages, and cards in the player's language.

Ideas: LO-01 – LO-04.

---

### [ ] 16.1 Translation ADR

1. Goal: pick how texts are translated.
2. Tasks:
   1. Compare i18n libraries for React; keep it small (PR-01).
   2. File format for translations.
   3. Fallback to English for missing keys.
   4. Update `NOTICE` with the chosen library.
3. Done when:
   1. ADR accepted by the **Owner**.

### [ ] 16.2 App texts

1. Goal: every UI string comes from translation files.
2. Tasks:
   1. Move menu, screen, and settings texts to keys.
   2. Lint or test: no hard-coded UI text left.
3. Done when:
   1. The app looks the same in English.

### [ ] 16.3 Server messages

1. Goal: each client shows messages in its own language (LO-04).
2. Tasks:
   1. Refusal reasons, prompts, history: key plus values.
   2. Client renders the text; English as fallback.
   3. Dedicated server console stays in English.
3. Done when:
   1. Two clients in different languages play one game.

### [ ] 16.4 Card texts

1. Goal: cards in the chosen language (LO-02).
2. Tasks:
   1. Card names and text from card_db translations.
   2. Missing translation: English text.
   3. Translators shown with the card (CD-06).
3. Done when:
   1. Base Game cards show in Russian.

### [ ] 16.5 Language setting

1. Goal: the player picks the language.
2. Tasks:
   1. Setting in the settings screen; default from the OS.
   2. First extra language: Russian.
   3. Use case: change the language.
3. Done when:
   1. Switching language changes menus, messages, and cards.
