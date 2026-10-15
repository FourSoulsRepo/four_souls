import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import './style.css';

// ?fake swaps in fake Wails bindings for make e2e and make screenshots.
// Dev server only: production builds drop the fixture entirely.
if (import.meta.env.DEV && new URLSearchParams(window.location.search).has('fake')) {
  await import('./fakeBridge');
}
const { App } = await import('./App');
// ?screen=cards opens the card gallery (dev server only, for tests).
if (import.meta.env.DEV && new URLSearchParams(window.location.search).get('screen') === 'cards') {
  const { useApp } = await import('./store');
  useApp.setState({ screen: 'cards' });
}

// With fake bindings, tests reach the store to play event bursts.
if (import.meta.env.DEV && new URLSearchParams(window.location.search).has('fake')) {
  const { useApp, dropConnection } = await import('./store');
  const { clashes } = await import('./animations');
  (window as unknown as { fourSouls: unknown }).fourSouls = { store: useApp, drop: dropConnection, clashes };
}

const container = document.getElementById('root');
if (!container) {
  throw new Error('root element is missing');
}

createRoot(container).render(
  <StrictMode>
    <App />
  </StrictMode>,
);
