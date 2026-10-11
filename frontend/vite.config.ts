import {defineConfig} from 'vite'
import react from '@vitejs/plugin-react'

// https://vitejs.dev/config/
export default defineConfig({
  plugins: [react()],
  server: {
    fs: {
      // The dev-only fake bridge reads the real card text (src/fakeBridge.ts).
      allow: ['.', '../pkg/card_db/data'],
    },
  },
})
