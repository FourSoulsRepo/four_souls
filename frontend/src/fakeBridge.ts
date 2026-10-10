// Loaded only when the page URL has ?fake (make e2e, make screenshots).
// Installs fake Wails bindings so Vite can render screens without the Go runtime.
// Keep the notice in sync with internal/legal/notice.go.

interface Part {
  text: string;
  url?: string;
}

const site = 'https://foursouls.com';
const shop = 'https://maestromedia.com/collections/binding-of-isaac-four-souls';

const notice: { parts: Part[] }[] = [
  { parts: [{ text: 'Four Souls is an unofficial, free, fan-made game.' }] },
  {
    parts: [
      { text: 'The Binding of Isaac: Four Souls is designed by ' },
      { text: 'Edmund McMillen', url: site },
      { text: ' and published by ' },
      { text: 'Maestro Media', url: shop },
      { text: '.' },
    ],
  },
  { parts: [{ text: 'All rights to the game, its rules, card texts and artwork belong to their owners.' }] },
  { parts: [{ text: 'This project is not affiliated with or endorsed by them.' }] },
  {
    parts: [
      { text: 'Please support the official game: ' },
      { text: 'buy it here', url: shop },
      { text: '.' },
    ],
  },
];

// Settings live in localStorage, so a page reload acts like an app restart.
const settingsKey = 'fake-settings';

function loadSettings(): unknown {
  return JSON.parse(localStorage.getItem(settingsKey) ?? '{"nickname":""}');
}

// Same checks as internal/settings (Go counts runes, so code points), so
// tests see the Go error.
function saveSettings(s: { nickname: string }): Promise<unknown> {
  const nickname = s.nickname.trim();
  if (s.nickname !== '' && (nickname === '' || Array.from(nickname).length > 20)) {
    return Promise.reject(new Error('nickname must be 1–20 characters'));
  }
  const saved = { nickname };
  localStorage.setItem(settingsKey, JSON.stringify(saved));
  return Promise.resolve(saved);
}

const w = window as unknown as {
  go: { main: { App: Record<string, (...args: never[]) => Promise<unknown>> } };
  runtime: { BrowserOpenURL: (url: string) => void };
};

w.go = {
  main: {
    App: {
      Versions: () => Promise.resolve({ app: 'demo', engine: 'demo' }),
      Notice: () => Promise.resolve(notice),
      Settings: () => Promise.resolve(loadSettings()),
      SaveSettings: saveSettings,
    },
  },
};
w.runtime = {
  BrowserOpenURL: () => undefined,
};

export {};
