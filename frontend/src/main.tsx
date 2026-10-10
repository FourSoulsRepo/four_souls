import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import './style.css';

// ?fake swaps in fake Wails bindings for make e2e and make screenshots.
// Dev server only: production builds drop the fixture entirely.
if (import.meta.env.DEV && new URLSearchParams(window.location.search).has('fake')) {
  await import('./fakeBridge');
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
