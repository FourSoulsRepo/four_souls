/// <reference types="vitest/config" />
import {defineConfig} from 'vite'
import react from '@vitejs/plugin-react'

// The rules website (ADR 026) is a second entry: vite --mode site and
// vite build --mode site. It is served under /four_souls/ as on GitHub
// Pages; make site puts site.wasm.gz and wasm_exec.js into .cache/site,
// its public folder, and renames the page to index.html. The app's
// build (wails build, npm run build) stays as it was.
const site = {
  base: '/four_souls/',
  publicDir: '../.cache/site',
  build: {
    outDir: 'dist-site',
    rolldownOptions: { input: 'site.html' },
  },
}

// https://vitejs.dev/config/
export default defineConfig(({ mode }) => ({
  plugins: [react()],
  ...(mode === 'site' ? site : {}),
  server: {
    // e2e runs (playwright.config.ts) watch no files: a Wails build or an
    // edit elsewhere would reload every test page and drop its game (19.12).
    ...(process.env.FOUR_SOULS_E2E === '1' ? { watch: null, hmr: false } : {}),
    fs: {
      // The dev-only fake bridge reads the real card text (src/fakeBridge.ts);
      // the site reads the rules text and the sample payloads.
      allow: ['.', '../pkg/card_db/data', '../pkg/rules_engine/docs/rules', '../internal/website/testdata', '../.cache/site'],
    },
  },
  test: {
    // The animation inventory reads the CSS as text (animations.test.ts).
    css: true,
    // Vitest 5 keeps transformed files on disk across runs. A file with
    // import.meta.glob<T>( is cached too, so its file list goes stale:
    // the inventory tests then miss a new file (19.124). It saves 0.2 s.
    fsModuleCache: false,
  },
}))
