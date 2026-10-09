import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import './style.css';

// ?screenshot=1 swaps in fake Wails bindings for make screenshots.
// Dev server only: production builds drop the fixture entirely.
if (import.meta.env.DEV && new URLSearchParams(window.location.search).has('screenshot')) {
  await import('./screenshotBridge');
}
const { App } = await import('./App');

const container = document.getElementById('root');
if (!container) {
  throw new Error('root element is missing');
}

createRoot(container).render(
  <StrictMode>
    <App />
  </StrictMode>,
);
