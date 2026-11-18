// Loaded only when the page URL has ?fake (make e2e, make screenshots).
// Installs fake Wails bindings so Vite can render screens without the Go runtime.
// Keep the notice in sync with internal/legal/notice.go.

import b2 from '../../pkg/card_db/data/b2.json?raw';
import requiem from '../../pkg/card_db/data/r.json?raw';
import challenge from '../../pkg/card_db/data/challenge.json?raw';
import type { RecordSeat } from './bridge';
import { fakeStep, type FakeRecording } from './fakeReplay';
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
  // ?noendask ends the turn without asking (V-46); ?scale=150 starts at
  // that interface scale (ST-10).
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
    host: { seats: 0, sets: null, detailed: false, response_timer: 0, no_bonus_souls: false, no_deck_ratios: false, room_deck: false, random_seats: false, rollback: '', picking: '', draft_size: 0, mode: '', challenge: '', difficulty: '' },
    lan_address: '',
    // make screenshots LOCALE=qps starts every page in that language.
    language: import.meta.env.VITE_LOCALE ?? '',
    volume: 100,
    scale: Number(params.get('scale') ?? 100),
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
    room_deck?: boolean;
    random_seats?: boolean;
    rollback?: string;
    picking?: string;
    draft_size?: number;
    mode?: string;
    challenge?: string;
    difficulty?: string;
  };
  lan_address?: string;
  language?: string;
  volume?: number;
  scale?: number;
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
  if (s.scale !== undefined && s.scale !== 0 && (!Number.isInteger(s.scale) || s.scale < 75 || s.scale > 200)) {
    return Promise.reject(new Error(`unknown setting value: scale ${String(s.scale)}`));
  }
  const saved = { ...s, scale: s.scale === undefined || s.scale === 0 ? 100 : s.scale, nickname, recent: s.recent.slice(0, 5), lan_address: (s.lan_address ?? '').trim() };
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
  {
    ...none,
    ref: 'lost_soul',
    set: 'b2',
    type: 'Lost Soul Card',
    name: 'Lost Soul',
    effects: ['When this enters play, it becomes a soul.'],
    image: '',
  },
];

interface DataCard {
  id: string;
  version?: number;
  type: string;
  text: Partial<Record<string, { name: string; effects?: string[]; footnotes?: { title: string; text: string }[] }>>;
  translators?: Record<string, string[]>;
  artists?: { role: string; name: string }[];
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
// The challenge cards (ADR 028): a challenge's boss and rule cards, as
// Go's Assets.Cards lists them; never a playable set (14.25).
const challengeData = JSON.parse(challenge) as DataSet;

// translated gives a card's names in every language, its translators
// and artists, as Go's Assets.Cards does (16.5, 20.24); the cards above
// take them from b2 too.
function translated(ref: string, name: string): { names: Record<string, string>; translators: Record<string, string[]>; artists: string[] } {
  const c = b2Data.cards.find((d) => d.id === ref && !d.version);
  const names = Object.fromEntries(Object.entries(c?.text ?? {}).flatMap(([lang, text]) => (text?.name ? [[lang, text.name]] : [])));
  return { names: { en: name, ...names }, translators: c?.translators ?? {}, artists: c ? artistNames(c, false) : [] };
}

function realCards(): unknown[] {
  return [b2Data, requiemData, challengeData].flatMap((data) =>
    data.cards
      .filter((c) => !c.reprint && !cards.some((s) => s.ref === c.id))
      .flatMap((c) => {
        const front = { ...realCard(data.set.code, c), ...(c.back ? { back: c.back.id } : {}) };
        // A back face is a card of its own, as in Go's Assets.Cards.
        return c.back ? [front, { ...realCard(data.set.code, { ...c, ...c.back, version: 0 }, true), front: c.id }] : [front];
      }),
  );
}

// artistNames picks a face's credited artists as Go's Card.ArtistNames
// does (CD-14, 20.24): a back face takes the "Back " roles, a front the
// rest; never the card back's artist or a name that is not public.
function artistNames(c: DataCard, backFace: boolean): string[] {
  const names: string[] = [];
  for (const a of c.artists ?? []) {
    const name = a.name.trim();
    const skip = a.role === 'Card back artist' || name === '' || name.toLowerCase() === 'private';
    if (!skip && a.role.startsWith('Back ') === backFace && !names.includes(name)) {
      names.push(name);
    }
  }
  return names;
}

// realCard is one card of the card data as Go's Assets.Cards gives it;
// backFace marks a two-sided card's back (ADR 023).
function realCard(set: string, c: DataCard, backFace = false): Record<string, unknown> {
  const names = Object.fromEntries(Object.entries(c.text).flatMap(([lang, text]) => (text?.name ? [[lang, text.name]] : [])));
  return {
    names: { en: c.text.en?.name ?? c.id, ...names },
    translators: c.translators ?? {},
    artists: artistNames(c, backFace),
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
  // Wails puts each bound struct under its Go package: go.<package>.<Struct> (ADR 030).
  go: Record<string, Record<string, Record<string, (...args: never[]) => Promise<unknown>>>>;
  runtime: Record<string, (...args: never[]) => void>;
  // fakeFullscreen lets tests see the window mode.
  fakeFullscreen: boolean;
  // fakeQuit lets tests see that the app was asked to close.
  fakeQuit: boolean;
  // fakeLan is what the LAN search finds, when a test sets it.
  fakeLan?: unknown[];
  // fakeUpdate is what the update check gives (18.24): a status as Go's
  // update.Status, or { error } with Go's message text. Unset: an update.
  fakeUpdate?: object;
  // fakeDownload ends the running fake download: with no error it
  // succeeds, else it fails with that text.
  fakeDownload?: { finish: (error?: string) => void };
  // fakeOpened is the page OpenReleases opened.
  fakeOpened?: string;
  // fakeEmit sends a Wails event to the page, e.g. update:progress.
  fakeEmit: (name: string, ...data: unknown[]) => void;
  // fakeRecords is the records folder (M-01): a test may replace it
  // before the Records screen opens; Delete takes from it.
  fakeRecords: FakeRecord[];
  // fakeExported is the last record Export copied.
  fakeExported?: string;
  // fakePick is the file the Open a file dialog picks; unset cancels.
  fakePick?: string;
  // fakeFiles are whole records by path, as Go's Open and Step give
  // them: a test reads a real record from disk and plays it (10.2).
  fakeFiles?: Partial<Record<string, FakeFile>>;
};

// FakeFile is a record as Open and Step give it.
interface FakeFile {
  summary: unknown;
  steps: unknown[];
}

// FakeRecord is a records.Entry as Go lists it.
interface FakeRecord {
  path: string;
  name: string;
  started: number;
  players: RecordSeat[];
  finished: boolean;
  winners: number[]; // seats
  steps: number;
  challenge: string;
  why: string;
}

// fakeRecord is a listed record in the fake records folder: players
// in seat order.
function fakeRecord(name: string, started: string, players: string[], more: Partial<FakeRecord> = {}): FakeRecord {
  const seats = players.map((n, seat) => ({ seat, name: n }));
  return { path: `/fake/records/${name}.fsrec`, name: `${name}.fsrec`, started: Date.parse(started), players: seats, finished: false, winners: [], steps: 0, challenge: '', why: '', ...more };
}

// The fake records folder: one of each kind the screen shows.
const fakeRecords = [
  fakeRecord('2026-11-16-200000-g3_Ann_Bo', '2026-11-16T20:00:00Z', ['Ann', 'Bo'], { steps: 12 }),
  fakeRecord('2026-11-15-180000-g2_Ann_Bo_Cy', '2026-11-15T18:00:00Z', ['Ann', 'Bo', 'Cy'], { finished: true, winners: [1], steps: 840 }),
  fakeRecord('2026-10-01-120000-g1_Ann_Bo', '2026-10-01T12:00:00Z', ['Ann', 'Bo'], { finished: true, winners: [0], steps: 500, why: 'before_replays' }),
  fakeRecord('from-the-future', '', [], { started: 0, why: 'newer' }),
];

// goError is an error as a Go binding sends a player's message.
function goError(key: string): Error {
  return new Error(JSON.stringify({ key }));
}

// opened is the record Open opened, for Step (10.2).
let opened: FakeRecording | FakeFile | null = null;

// fakeSummary is what Open gives for a record: a turn every ten steps,
// as its steps have them (fakeReplay.ts).
function fakeSummary(path: string): Promise<unknown> {
  const e = w.fakeRecords.find((r) => r.path === path) ?? (path === w.fakePick ? fakeRecord('picked', '2026-11-14T10:00:00Z', ['Dee', 'Eve'], { path, finished: true, winners: [1], steps: 40 }) : undefined);
  if (!e) {
    return Promise.reject(goError('error.record_unreadable'));
  }
  if (e.why) {
    return Promise.reject(goError(`error.record_${e.why}`));
  }
  const seats = e.players;
  const turns = Array.from({ length: Math.ceil(e.steps / 10) }, (_, i) => ({ number: i + 1, active: i % seats.length, step: i * 10 }));
  const winners = e.winners;
  opened = { players: seats.length, steps: e.steps, winners };
  return Promise.resolve({
    path,
    game: e.name,
    server: '',
    started: e.started,
    seats,
    sets: ['b2'],
    challenge: e.challenge,
    difficulty: '',
    cards: [],
    steps: e.steps,
    turns,
    end: e.finished ? { ended: e.started + 3_600_000, finished: true, winners } : undefined,
  });
}

// The update the fake check finds unless a test sets fakeUpdate.
const fakeReleases = 'https://github.com/FourSoulsRepo/four_souls/releases';
const newerRelease = {
  current: 'demo',
  latest: 'v9.9.9',
  notes: '## What is new\n\n- A new card set.\n- Faster games.',
  url: `${fakeReleases}/tag/v9.9.9`,
  releases: fakeReleases,
  asset: 'four-souls-9.9.9-macos-universal.zip',
  newer: true,
  canInstall: true,
  reason: '',
};
let checkedUrl = '';

// Wails events: what the page listens to, by name.
const listeners = new Map<string, Set<(...data: unknown[]) => void>>();

w.go = {
  // A normal start: never a --diagnose probe (ADR 035).
  probe: {
    Probe: {
      Variant: () => Promise.resolve(''),
      Report: () => Promise.resolve(),
    },
  },
  about: {
    About: {
      Versions: () => Promise.resolve({ app: 'demo', engine: 'demo' }),
      Notice: () => Promise.resolve(notice),
    },
  },
  prefs: {
    Prefs: {
      Settings: () => Promise.resolve(loadSettings()),
      SaveSettings: saveSettings,
    },
  },
  setup: {
    Setup: {
      // ?requiem adds Requiem, the only set the deck ratios need; other
      // tests keep the older sets.
      CardSets: () => Promise.resolve(['b2', 'g2', 'fsp2', ...(params.has('requiem') ? ['r'] : [])]),
      Challenges: () =>
        Promise.resolve([
          { id: 'krampus', difficulties: ['normal', 'hard', 'ultra'] },
          { id: 'trick_or_treat', difficulties: ['normal', 'hard', 'ultra'] },
        ]),
      // Go's answer for today's card data (Setup.DeckRatios): the ratios
      // need Requiem, and Four Souls+ too in a 2-player game.
      DeckRatios: (sets: string[], players: number) =>
        Promise.resolve(
          sets.includes('r') && (sets.includes('fsp2') || players >= 3)
            ? ''
            : JSON.stringify([{ key: 'refuse.deck_ratios', args: { n: players } }, { key: 'refuse.ratio.butter_beans', args: { have: 3, need: 5 } }]),
        ),
    },
  },
  join: {
    Join: {
      // The real server's answer, like Go's Join.CheckSeat: the first
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
      // Like Go's Join.ServerStatus: the first base that answers.
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
  },
  records: {
    Records: {
      SaveRecord: (link: string) => Promise.resolve(`/fake/${new URL(link).searchParams.get('game') ?? ''}.fsrec`),
      List: () => Promise.resolve(w.fakeRecords.map((r) => ({ ...r }))),
      Open: (path: string) => {
        const file = Object.entries(w.fakeFiles ?? {}).find(([p]) => p === path)?.[1];
        opened = file ?? null;
        return file ? Promise.resolve(file.summary) : fakeSummary(path);
      },
      Step: (index: number) => {
        const s = opened && 'summary' in opened ? opened.steps.at(index) : opened ? fakeStep(opened, index) : undefined;
        return s ? Promise.resolve(s) : Promise.reject(new Error(`records: no such step: ${String(index)}`));
      },
      Delete: (path: string) => {
        if (!path.endsWith('.fsrec')) {
          return Promise.reject(goError('error.record_not_a_record'));
        }
        w.fakeRecords = w.fakeRecords.filter((r) => r.path !== path);
        return Promise.resolve();
      },
      Export: (path: string) => {
        w.fakeExported = path;
        return Promise.resolve(`/fake/export/${path.split('/').pop() ?? ''}`);
      },
      PickFile: () => Promise.resolve(w.fakePick ?? ''),
    },
  },
  assets: {
    Assets: {
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
          ['loot', 'treasure', 'monster', 'room'].map((kind) => ({ kind, image: `/cardimg/assets/backs/${kind}.webp` })),
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
      Sounds: () => {
        const one = (name: string): string[] => [`/cardimg/assets/sounds/${name}.mp3`];
        const three = (name: string): string[] => ['1', '2', '3'].map((i) => `/cardimg/assets/sounds/${name}-${i}.mp3`);
        return Promise.resolve({
          attackHit: one('attack-hit-1'),
          effectCast: one('effect-cast'),
          effectImpact: one('effect-impact'),
          diceRoll: one('dice-roll'),
          lootDraw: one('loot-draw'),
          playerDied: one('player-died'),
          treasureGain: one('treasure-gain'),
          soulGain: one('soul-gain'),
          attackHit2: one('attack-hit-2'),
          attackHit3: one('attack-hit-3'),
          lootLand: three('loot-land'),
          lootDrawOther: three('loot-draw-other'),
          dragBed: one('drag-bed'),
          dragAir: one('drag-air'),
          turnStart: one('turn-start'),
          cardPlayOther: three('card-play-other'),
          endTurnCue: one('end-turn-cue'),
          timerTick: one('timer-tick'),
          timerLast: one('timer-last'),
          choiceDeal: one('choice-deal'),
          choicePick: one('choice-pick'),
          cardGrab: one('card-grab'),
          aimCancel: one('aim-cancel'),
          playRefused: one('play-refused'),
          historyHover: one('history-hover'),
          choiceUnpick: one('choice-unpick'),
          choiceConfirm: one('choice-confirm'),
        });
      },
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
                { id: 'halloween', name: 'Halloween', image: '/cardimg/assets/mats/halloween.webp' },
                { id: 'christmas', name: 'Christmas', image: '/cardimg/assets/mats/christmas.webp' },
              ],
        ),
    },
  },
  host: {
    Host: {
      Start: () => Promise.resolve(hostInfo),
      Info: () => Promise.resolve(hostInfo),
      Stop: () => Promise.resolve(),
    },
  },
  updater: {
    // No network: a test sets the check's answer and sends the download's
    // progress events itself (fakeEmit), then ends it (fakeDownload).
    Updater: {
      Check: () => {
        const got = (w.fakeUpdate ?? newerRelease) as { error?: string; url?: string };
        if (got.error !== undefined) {
          return Promise.reject(new Error(got.error));
        }
        checkedUrl = got.url ?? '';
        return Promise.resolve(got);
      },
      Download: () =>
        new Promise<void>((resolve, reject) => {
          w.fakeDownload = {
            finish: (error?: string) => {
              if (error === undefined) {
                resolve();
              } else {
                reject(new Error(error));
              }
            },
          };
          w.fakeEmit('update:progress', { done: 0, total: 0 });
        }),
      ApplyAndQuit: () => {
        w.fakeQuit = true;
        return Promise.resolve();
      },
      OpenReleases: () => {
        w.fakeOpened = checkedUrl || fakeReleases;
        return Promise.resolve();
      },
    },
  },
};
w.fakeFullscreen = false;
w.fakeQuit = false;
w.fakeRecords = fakeRecords;
w.fakeEmit = (name, ...data) => {
  for (const callback of listeners.get(name) ?? []) {
    callback(...data);
  }
};
w.runtime = {
  BrowserOpenURL: () => undefined,
  // Wails' EventsOn: it returns the call that stops listening.
  EventsOnMultiple: (name: string, callback: (...data: unknown[]) => void) => {
    const set = listeners.get(name) ?? new Set();
    set.add(callback);
    listeners.set(name, set);
    return () => {
      set.delete(callback);
    };
  },
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
