// Loaded only when the page URL has ?fake (make e2e, make screenshots).
// Installs fake Wails bindings so Vite can render screens without the Go runtime.
// Keep the notice in sync with internal/legal/notice.go.

import b2 from '../../pkg/card_db/data/b2.json?raw';
import requiem from '../../pkg/card_db/data/r.json?raw';
import { protocolVersion } from './protocol';

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
  // ?anim=off starts with animations off; ?mat=basement starts on that
  // mat; ?firstrun starts without a nickname, as on the first start;
  // ?noendask ends the turn without asking (V-46).
  const params = new URLSearchParams(window.location.search);
  const animation = params.get('anim') ?? 'normal';
  const mat = params.get('mat') ?? '';
  const nickname = params.has('firstrun') ? '' : 'Player';
  const fresh = {
    nickname,
    recent: [],
    animation,
    fullscreen: false,
    mat,
    dice: 'tear',
    no_end_turn_ask: params.has('noendask'),
    rejoin: { address: '', token: '', game: '', via: '' },
    host: { seats: 0, sets: null, detailed: false, response_timer: 0, no_bonus_souls: false, no_deck_ratios: false, rollback: '' },
    lan_address: '',
    // make screenshots LOCALE=qps starts every page in that language.
    language: import.meta.env.VITE_LOCALE ?? '',
    volume: 100,
  };
  return JSON.parse(localStorage.getItem(settingsKey) ?? JSON.stringify(fresh));
}

// Same checks as internal/settings (Go counts runes, so code points), so
// tests see the Go error.
interface FakeSettings {
  nickname: string;
  recent: { address: string; name: string }[];
  animation: string;
  fullscreen: boolean;
  mat: string;
  dice: string;
  no_end_turn_ask?: boolean;
  rejoin: { address: string; token: string; game: string; via?: string };
  host?: {
    seats: number;
    sets: string[] | null;
    detailed: boolean;
    response_timer: number;
    no_bonus_souls: boolean;
    no_deck_ratios: boolean;
    rollback?: string;
  };
  lan_address?: string;
  language?: string;
  volume?: number;
}

function saveSettings(s: FakeSettings): Promise<unknown> {
  const nickname = s.nickname.trim();
  if (s.nickname !== '' && (nickname === '' || Array.from(nickname).length > 20)) {
    return Promise.reject(new Error('nickname must be 1–20 characters'));
  }
  if (s.language && !/^[a-z]{2,3}$/.test(s.language)) {
    return Promise.reject(new Error(`unknown setting value: "${s.language}"`));
  }
  if (s.volume !== undefined && (!Number.isInteger(s.volume) || s.volume < 0 || s.volume > 100)) {
    return Promise.reject(new Error(`unknown setting value: ${String(s.volume)}`));
  }
  const saved = { ...s, nickname, recent: s.recent.slice(0, 5), lan_address: (s.lan_address ?? '').trim() };
  localStorage.setItem(settingsKey, JSON.stringify(saved));
  return Promise.resolve(saved);
}

// A few cards: one with an image, one whose image fails to load, and
// some without images, so tests see both faces of the Card component.
// The real card text follows them, with images at their real paths; the
// e2e tests serve those only for screenshots, else they fail to load and
// the cards show their text, as in a public clone.
// Isaac also has a Russian test image, so tests see the image follow the
// language (16.6); no real card has one yet.
const artImage = (text: string): string =>
  'data:image/svg+xml,' +
  encodeURIComponent(
    '<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 63 88"><rect width="63" height="88" fill="#7a5c2e"/>' +
      `<text x="31" y="46" font-size="9" text-anchor="middle" fill="#fff">${text}</text></svg>`,
  );
const isaacImage = artImage('Isaac art');
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
    images: { en: isaacImage, ru: artImage('Isaac ru') },
  },
  {
    ...none,
    ref: 'the_d6',
    set: 'b2',
    type: 'Eternal Treasure Card',
    name: 'The D6',
    effects: ['{Tap Effect}Choose a dice roll. Its controller rerolls it.', 'At the end of your turn, recharge this.'],
    footnotes: [{ title: '-Eternal-', text: "This can't be destroyed or put into discard." }],
    // A path with no file, so its image always fails, screenshots too.
    image: '/cardimg/en/b2/missing_the_d6.webp',
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
  text: Partial<Record<string, { name: string; effects?: string[]; footnotes?: { title: string; text: string }[] }>>;
  translators?: Record<string, string[]>;
  stats?: { name: string; value: string }[];
  rewards?: { count: string; kind: string }[];
  images?: string[]; // one per copy; the first stands for the card
  reprint?: boolean;
  // back is a two-sided card's other face (ADR 023).
  back?: Pick<DataCard, 'id' | 'text' | 'stats' | 'rewards' | 'images'>;
}

interface DataSet {
  set: { code: string };
  cards: DataCard[];
}

const b2Data = JSON.parse(b2) as DataSet;
// Requiem has the two-sided cards (14.15); its cards are only dealt
// with ?requiem.
const requiemData = JSON.parse(requiem) as DataSet;

// translated gives a card's names in every language and its translators,
// as Go's App.Cards does (16.5); the cards above take them from b2 too.
function translated(ref: string, name: string): { names: Record<string, string>; translators: Record<string, string[]> } {
  const c = b2Data.cards.find((d) => d.id === ref && !d.version);
  const names = Object.fromEntries(Object.entries(c?.text ?? {}).flatMap(([lang, text]) => (text?.name ? [[lang, text.name]] : [])));
  return { names: { en: name, ...names }, translators: c?.translators ?? {} };
}

function realCards(): unknown[] {
  return [b2Data, requiemData].flatMap((data) =>
    data.cards
      .filter((c) => !c.reprint && !cards.some((s) => s.ref === c.id))
      .flatMap((c) => {
        const front = { ...realCard(data.set.code, c), ...(c.back ? { back: c.back.id } : {}) };
        // A back face is a card of its own, as in Go's App.Cards.
        return c.back ? [front, { ...realCard(data.set.code, { ...c, ...c.back, version: 0 }), front: c.id }] : [front];
      }),
  );
}

// realCard is one card of the card data as Go's App.Cards gives it.
function realCard(set: string, c: DataCard): Record<string, unknown> {
  const names = Object.fromEntries(Object.entries(c.text).flatMap(([lang, text]) => (text?.name ? [[lang, text.name]] : [])));
  return {
    names: { en: c.text.en?.name ?? c.id, ...names },
    translators: c.translators ?? {},
    ref: c.version && c.version > 1 ? `${c.id}@${String(c.version)}` : c.id,
    set,
    type: c.type,
    name: c.text.en?.name ?? c.id,
    effects: c.text.en?.effects ?? [],
    footnotes: c.text.en?.footnotes ?? [],
    stats: c.stats ?? [],
    rewards: c.rewards ?? [],
    // English only: the fake bridge cannot see which other images exist.
    image: c.images?.[0] ? `/cardimg/en/${c.images[0]}` : '',
    images: c.images?.[0] ? { en: `/cardimg/en/${c.images[0]}` } : {},
  };
}

// ?cardlang=xx names Isaac in the card-only language xx, which has no
// interface file, as card data with text.xx would (16.10). The real card
// files never get test names.
function withTestLanguage(list: unknown[]): unknown[] {
  const lang = params.get('cardlang');
  if (!lang) {
    return list;
  }
  return list.map((c) => {
    const card = c as { ref: string; names: Record<string, string> };
    return card.ref === 'isaac' ? { ...card, names: { ...card.names, [lang]: `Isaac (${lang})` } } : card;
  });
}

// Hosting points at a real dedicated server the e2e tests start
// (?server=PORT), so two pages reach the same lobby.
const params = new URLSearchParams(window.location.search);
const hostPort = Number(params.get('server') ?? '4774');
const hostInfo = { running: true, port: hostPort, addresses: ['127.0.0.1'] };

const w = window as unknown as {
  go: { main: Record<string, Record<string, (...args: never[]) => Promise<unknown>>> };
  runtime: Record<string, (...args: never[]) => void>;
  // fakeFullscreen lets tests see the window mode.
  fakeFullscreen: boolean;
  // fakeQuit lets tests see that the app was asked to close.
  fakeQuit: boolean;
  // fakeLan is what the LAN search finds, when a test sets it.
  fakeLan?: unknown[];
};

w.go = {
  main: {
    App: {
      Versions: () => Promise.resolve({ app: 'demo', engine: 'demo' }),
      Notice: () => Promise.resolve(notice),
      Settings: () => Promise.resolve(loadSettings()),
      SaveSettings: saveSettings,
      // ?requiem adds Requiem, the only set the deck ratios need; other
      // tests keep the older sets.
      CardSets: () => Promise.resolve(['b2', 'g2', 'fsp2', ...(params.has('requiem') ? ['r'] : [])]),
      // Go's answer for today's card data (App.DeckRatios): the ratios
      // need Requiem, and Four Souls+ too in a 2-player game.
      DeckRatios: (sets: string[], players: number) =>
        Promise.resolve(
          sets.includes('r') && (sets.includes('fsp2') || players >= 3)
            ? ''
            : JSON.stringify([{ key: 'refuse.deck_ratios', args: { n: players } }, { key: 'refuse.ratio.butter_beans', args: { have: 3, need: 5 } }]),
        ),
      Cards: () => Promise.resolve(withTestLanguage([...cards.map((c) => ({ ...c, ...translated(c.ref, c.name) })), ...realCards()])),
      // The dice textures at their real paths; like the mat, the e2e
      // tests serve them only for screenshots, else they fail: plain dice.
      // ?nodice has no textures, as Go without the private images.
      Dice: () =>
        Promise.resolve(
          [
            ['tear', 'Tear', 'holy'],
            ['blue_fly', 'Blue Fly', 'holy'],
            ['cent', 'Cent', 'holy'],
            ['moms_knife', "Mom's Knife", 'holy'],
            ['the_bible', 'The Bible', 'holy'],
            ['brimstone', 'Brimstone', 'unholy'],
            ['guppy_paw', "Guppy's Paw", 'unholy'],
            ['pentagram', 'Pentagram', 'unholy'],
            ['poop', 'Poop', 'unholy'],
            ['tinted_rock', 'Tinted Rock', 'unholy'],
          ].map(([id = '', name, set = '']) => ({
            id,
            name,
            set: set === 'holy' ? 'Holy' : 'Unholy',
            image: params.has('nodice') ? '' : `/cardimg/assets/dices/${set}_set/d6_${id}.webp`,
          })),
        ),
      // Card backs at their real paths; like the mat, the e2e tests serve
      // them only for screenshots, else they fail to load: plain backs.
      Backs: () =>
        Promise.resolve(
          ['loot', 'treasure', 'monster'].map((kind) => ({ kind, image: `/cardimg/assets/backs/${kind}.webp` })),
        ),
      // The coin and sword icons, the banner and the sky at their real
      // paths; like the backs, they load only in screenshots, else the
      // drawn ones, the text title and a plain background show.
      Icons: () =>
        Promise.resolve({
          coin: '/cardimg/assets/icons/coin.webp',
          sword: '/cardimg/assets/icons/sword.webp',
          banner: '/cardimg/assets/banner.webp',
          sky: ['sky', 'stars-1', 'stars-2', 'stars-3', 'stars-4', 'stars-5'].map((n) => `/cardimg/assets/background/${n}.webp`),
        }),
      // The fonts at their real paths; they load only in screenshots, else
      // the system font shows.
      Fonts: () =>
        Promise.resolve({
          title: '/cardimg/assets/fonts/font-souls-title.ttf',
          body: '/cardimg/assets/fonts/font-souls-body.ttf',
          cyrillic: '/cardimg/assets/fonts/tboi4s-original.ttf',
        }),
      // The sounds at their real paths; like the fonts, they load
      // only in screenshots, else they fail to load and stay silent.
      Sounds: () =>
        Promise.resolve({
          attackHit: '/cardimg/assets/sounds/attack-hit.mp3',
          effectCast: '/cardimg/assets/sounds/effect-cast.mp3',
          effectImpact: '/cardimg/assets/sounds/effect-impact.mp3',
          diceRoll: '/cardimg/assets/sounds/dice-roll.mp3',
          lootDraw: '/cardimg/assets/sounds/loot-draw.mp3',
          playerDied: '/cardimg/assets/sounds/player-died.mp3',
          treasureGain: '/cardimg/assets/sounds/treasure-gain.mp3',
        }),
      // The mat images, at their real paths; the e2e tests serve them only
      // for screenshots, else they fail to load as in a public clone.
      // ?nomats lists none, as Go does without the private images.
      Mats: () =>
        Promise.resolve(
          params.has('nomats')
            ? []
            : [
                { id: 'basement', name: 'Basement', image: '/cardimg/assets/mats/basement.webp' },
                { id: 'cathedral', name: 'Cathedral', image: '/cardimg/assets/mats/cathedral.webp' },
              ],
        ),
      // The real server's answer, like Go's App.CheckSeat: the first
      // base that answers.
      CheckSeat: async (bases: string[], token: string) => {
        for (const b of bases) {
          try {
            const r = await fetch(`${b}/seat?token=${encodeURIComponent(token)}`);
            if (r.ok) {
              return (await r.json()) as unknown;
            }
          } catch {
            // the next base
          }
        }
        throw new Error('server: no answer');
      },
      // Like Go's App.ServerStatus: the first base that answers.
      ServerStatus: async (bases: string[]) => {
        for (const b of bases) {
          try {
            const r = await fetch(`${b}/status`);
            if (r.ok) {
              return (await r.json()) as unknown;
            }
          } catch {
            // the next base
          }
        }
        throw new Error('server: no answer');
      },
      SaveRecord: (link: string) => Promise.resolve(`/fake/${new URL(link).searchParams.get('game') ?? ''}.fsrec`),
      // The LAN holds one hosting app: the test server (?server=PORT).
      // ?nolan finds nothing; a test sets window.fakeLan to list others,
      // in this app's protocol unless it names one.
      FindGames: () =>
        Promise.resolve(
          w.fakeLan?.map((f) => ({ protocol: protocolVersion, ...(f as object) })) ??
            (params.has('nolan')
              ? []
              : [{ address: `127.0.0.1:${String(hostPort)}`, protocol: protocolVersion, app: 'demo', games: [{ id: '', host: 'Lan', seats: 2, taken: 1 }] }]),
        ),
    },
    Host: {
      Start: () => Promise.resolve(hostInfo),
      Info: () => Promise.resolve(hostInfo),
      Stop: () => Promise.resolve(),
    },
  },
};
w.fakeFullscreen = false;
w.fakeQuit = false;
w.runtime = {
  BrowserOpenURL: () => undefined,
  Quit: () => {
    w.fakeQuit = true;
  },
  WindowFullscreen: () => {
    w.fakeFullscreen = true;
  },
  WindowUnfullscreen: () => {
    w.fakeFullscreen = false;
  },
};

export {};
