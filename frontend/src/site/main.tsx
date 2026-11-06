// Entry of the rules website (ADR 026): vite build --mode site.
// It never imports src/bridge: there is no Wails here (lint checks it).
import { StrictMode } from 'react';
import { createRoot } from 'react-dom/client';
import { Site } from './Site';
import './site.css';

const container = document.getElementById('root');
if (!container) {
  throw new Error('root element is missing');
}

createRoot(container).render(
  <StrictMode>
    <Site />
  </StrictMode>,
);
