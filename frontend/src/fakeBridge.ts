// Loaded only when the page URL has ?fake (make e2e, make screenshots).
// Installs fake Wails bindings so Vite can render screens without the Go runtime.
// Keep the notice in sync with internal/legal/notice.go.

import b2 from '../../pkg/card_db/data/b2.json?raw';

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
  // ?anim=off starts with animations off; ?firstrun starts without a
  // nickname, as on the first start.
  const params = new URLSearchParams(window.location.search);
  const animation = params.get('anim') ?? 'normal';
  const nickname = params.has('firstrun') ? '' : 'Player';
  const fresh = { nickname, recent: [], animation, fullscreen: false, theme: 'dark', mat: 'felt' };
  return JSON.parse(localStorage.getItem(settingsKey) ?? JSON.stringify(fresh));
}

// Same checks as internal/settings (Go counts runes, so code points), so
// tests see the Go error.
interface FakeSettings {
  nickname: string;
  recent: string[];
  animation: string;
  fullscreen: boolean;
  theme: string;
  mat: string;
}

function saveSettings(s: FakeSettings): Promise<unknown> {
  const nickname = s.nickname.trim();
  if (s.nickname !== '' && (nickname === '' || Array.from(nickname).length > 20)) {
    return Promise.reject(new Error('nickname must be 1–20 characters'));
  }
  const saved = { ...s, nickname, recent: s.recent.slice(0, 5) };
  localStorage.setItem(settingsKey, JSON.stringify(saved));
  return Promise.resolve(saved);
}

// A few cards: one with an image, one whose image fails to load, and
// some without images, so tests see both faces of the Card component.
// The real card text follows them, without images, as in a public clone.
const isaacImage =
  'data:image/svg+xml,' +
  encodeURIComponent(
    '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 63 88"><rect width="63" height="88" fill="#7a5c2e"/>' +
      '<text x="31" y="46" font-size="9" text-anchor="middle" fill="#fff">Isaac art</text></svg>',
  );
const none = { footnotes: [], stats: [], rewards: [], effects: [] };
const cards = [
  {
    ...none,
    ref: 'isaac',
    set: 'b2',
    type: 'Character Card',
    name: 'Isaac',
    effects: ['{Tap Effect}Play an additional loot card this turn.'],
    stats: [
      { name: 'HP', value: '2' },
      { name: 'ATK', value: '1' },
    ],
    image: isaacImage,
  },
  {
    ...none,
    ref: 'the_d6',
    set: 'b2',
    type: 'Eternal Treasure Card',
    name: 'The D6',
    effects: ['{Tap Effect}Choose a dice roll. Its controller rerolls it.', 'At the end of your turn, recharge this.'],
    footnotes: [{ title: '-Eternal-', text: "This can't be destroyed or put into discard." }],
    image: '/cardimg/b2/the_d6.webp',
  },
  {
    ...none,
    ref: 'big_spider',
    set: 'b2',
    type: 'Basic Monster Card',
    name: 'Big Spider',
    effects: ['When this dies, the active player may attack the monster deck an additional time.'],
    stats: [
      { name: 'HP', value: '3' },
      { name: 'DC', value: '4' },
      { name: 'ATK', value: '1' },
    ],
    rewards: [{ count: '1x', kind: 'Loot' }],
    image: '',
  },
  {
    ...none,
    ref: 'a_dime',
    set: 'b2',
    type: '10 Cent Card',
    name: 'A Dime!!',
    effects: ['Gain 10¢.'],
    image: '',
  },
  {
    ...none,
    ref: 'bonus_soul',
    set: 'b2',
    type: 'Bonus Soul Card',
    name: 'Soul of Greed',
    effects: ['If you have 25¢ or more, gain this soul.'],
    image: '',
  },
];

interface DataCard {
  id: string;
  version?: number;
  type: string;
  text: { en?: { name: string; effects?: string[]; footnotes?: { title: string; text: string }[] } };
  stats?: { name: string; value: string }[];
  rewards?: { count: string; kind: string }[];
}

function realCards(): unknown[] {
  const data = JSON.parse(b2) as { set: { code: string }; cards: DataCard[] };
  return data.cards
    .filter((c) => !cards.some((s) => s.ref === c.id))
    .map((c) => ({
      ref: c.version && c.version > 1 ? `${c.id}@${String(c.version)}` : c.id,
      set: data.set.code,
      type: c.type,
      name: c.text.en?.name ?? c.id,
      effects: c.text.en?.effects ?? [],
      footnotes: c.text.en?.footnotes ?? [],
      stats: c.stats ?? [],
      rewards: c.rewards ?? [],
      image: '',
    }));
}

// Hosting points at a real dedicated server the e2e tests start
// (?server=PORT), so two pages reach the same lobby.
const hostPort = Number(new URLSearchParams(window.location.search).get('server') ?? '4774');
const hostInfo = { running: true, port: hostPort, addresses: ['127.0.0.1'] };

const w = window as unknown as {
  go: { main: Record<string, Record<string, (...args: never[]) => Promise<unknown>>> };
  runtime: Record<string, (...args: never[]) => void>;
  // fakeFullscreen lets tests see the window mode.
  fakeFullscreen: boolean;
};

w.go = {
  main: {
    App: {
      Versions: () => Promise.resolve({ app: 'demo', engine: 'demo' }),
      Notice: () => Promise.resolve(notice),
      Settings: () => Promise.resolve(loadSettings()),
      SaveSettings: saveSettings,
      CardSets: () => Promise.resolve(['b2']),
      Cards: () => Promise.resolve([...cards, ...realCards()]),
      SaveRecord: (link: string) => Promise.resolve(`/fake/${new URL(link).searchParams.get('game') ?? ''}.fsrec`),
    },
    Host: {
      Start: () => Promise.resolve(hostInfo),
      Info: () => Promise.resolve(hostInfo),
      Stop: () => Promise.resolve(),
    },
  },
};
w.fakeFullscreen = false;
w.runtime = {
  BrowserOpenURL: () => undefined,
  WindowFullscreen: () => {
    w.fakeFullscreen = true;
  },
  WindowUnfullscreen: () => {
    w.fakeFullscreen = false;
  },
};

export {};
