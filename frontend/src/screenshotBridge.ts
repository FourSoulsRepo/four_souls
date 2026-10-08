// Loaded only when the page URL has ?screenshot=1 (make screenshots).
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

const w = window as unknown as {
  go: { main: { App: Record<string, () => Promise<unknown>> } };
  runtime: { BrowserOpenURL: (url: string) => void };
};

w.go = {
  main: {
    App: {
      Versions: () => Promise.resolve({ app: 'demo', engine: 'demo' }),
      Notice: () => Promise.resolve(notice),
    },
  },
};
w.runtime = {
  BrowserOpenURL: () => undefined,
};

export {};
