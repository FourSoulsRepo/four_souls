/// <reference types="vitest/config" />
import {defineConfig} from 'vite'
import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    // e2e runs (playwright.config.ts) watch no files: a Wails build or an
    // edit elsewhere would reload every test page and drop its game (19.12).
    ...(process.env.FOUR_SOULS_E2E === '1' ? { watch: null, hmr: false } : {}),
    fs: {
      // The dev-only fake bridge reads the real card text (src/fakeBridge.ts).
      allow: ['.', '../pkg/card_db/data'],
    },
  },
  test: {
    // The animation inventory reads the CSS as text (animations.test.ts).
    css: true,
  },
})
