/// <reference types="vite/client" />

interface ImportMetaEnv {
  // VITE_LOCALE is the language fake settings start in (make screenshots LOCALE=qps).
  readonly VITE_LOCALE?: string;
}

interface ImportMetaEnv {
  // make i18n sets them (translations.test.ts): the language to check
  // in full, and the one area its todo.json holds.
  readonly I18N_LOCALE?: string;
  readonly I18N_AREA?: string;
}
