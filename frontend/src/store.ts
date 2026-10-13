// The app's one store (ADR 008). Screens read it; only bridge results
// and server messages write it.
import { create } from 'zustand';
import {
  findGames,
  setFullscreen,
  getCards,
  getDice,
  getBacks,
  getFonts,
  getIcons,
  getMats,
  getCardSets,
  saveRecord,
  getSettings,
  getVersions,
  saveSettings,
  startHosting,
  stopHosting,
  type DiceSkin,
  type HostInfo,
  type Icons,
  type LanServer,
  type MatInfo,
  type Settings,
} from './bridge';
import { noHostSetup } from './hostSetup';
import type { CardInfo, Kind } from './cards';
import { loadFonts } from './fonts';
import { entriesFor, lookedAt, type Entry, type Looked } from './history';
import type { LinePart } from './stack';
import { emptyLobby, leftTable, lost, onMessage, remember, type Screen, type Session } from './lobby';
import { connect, disconnect, drop, send } from './net';
import { createPlayback, effectsFor, type Effect } from './playback';
import {
  protocolVersion,
  recordURL,
  wsURL,
  type Create,
  type Envelope,
  type Hello,
  type Intent,
  type MonsterStats,
  type Update,
} from './protocol';

export type { Screen } from './lobby';

// Zoom is the card shown large: on hover, or after a long press on touch
// (B-06). side is where it shows, away from the pointer.
export interface Zoom {
  ref: string;
  touch: boolean;
  side: 'left' | 'right';
  // marks are the zoomed card's badges, as on the table (C-02).
  marks?: Marks;
  // caption goes below the card: who owns it, and the used ability's
  // words, highlighted (S-06, S-07).
  caption?: ZoomCaption;
}

// Marks are the badges a card on the table carries.
export interface Marks {
  damage?: number;
  counters?: number;
  stats?: { hp: number; atk: number };
  monster?: MonsterStats;
}

export interface ZoomCaption {
  owner?: string;
  color?: string;
  ability?: string;
  // line says what a stack item without ability words did: whom an
  // attack roll was against, who dealt damage to whom (15.22).
  line?: LinePart[];
}

// Playing is the effect on screen now; done ends it (playback.ts).
export interface Playing {
  effect: Effect;
  scale: number;
  key: number;
  done: () => void;
}

interface AppState extends Session {
  // playing is the event effect on screen, if any.
  playing: Playing | null;
  // history is the match so far, as this viewer saw it (7.9).
  history: Entry[];
  // looking is the cards an effect lets this player look at, shown
  // until endLook (V-07).
  looking: Looked[];
  endLook: () => void;
  // receive handles one server message.
  receive: (env: Envelope) => void;
  // behind is set while updates wait to be shown: the player acts only
  // on the latest state.
  behind: boolean;
  // sent is the step an intent was sent at; actions wait for the next
  // update, or for the refusal.
  sent: number | null;
  // selected is the card whose action menu is open.
  selected: number | null;
  act: (intent: Intent) => void;
  // skipAll turns "skip all" on or off (N-06).
  skipAll: (on: boolean) => void;
  // address is the server this app is connected to, for reconnecting.
  address: string;
  // reconnect says hello again with the seat's token (N-08).
  reconnect: () => Promise<void>;
  // vote: kick the missing players, or wait for them (N-08).
  vote: (kick: boolean) => void;
  // saveRecord saves the finished game's record; it returns the file,
  // or '' when cancelled or failed (the error says why).
  saveRecord: () => Promise<string>;
  select: (id: number | null) => void;
  settings: Settings;
  // hosting is set while this app runs the game server.
  hosting: HostInfo | null;
  // connected is set while the socket to a server is open.
  connected: boolean;
  // cardSets are the sets this app has, sent in hello.
  cardSets: string[];
  // cards is the display data by ref; empty until loaded.
  cards: Map<string, CardInfo>;
  // dice are the dice looks (V-19); empty until loaded.
  dice: DiceSkin[];
  // mats are the mat images (V-24); empty until loaded or without images.
  mats: MatInfo[];
  // backs are the deck back images by kind (V-26); a missing kind is plain.
  backs: Map<Kind, string>;
  // icons are the game's own icon images (15.38); '' draws our own.
  icons: Icons;
  zoom: Zoom | null;
  loadCards: () => Promise<void>;
  setZoom: (zoom: Zoom | null) => void;
  show: (screen: Screen) => void;
  loadSettings: () => Promise<void>;
  // saveSettings reports whether Go stored them.
  saveSettings: (s: Settings) => Promise<boolean>;
  // saveNickname saves the nickname if it changed (ST-06).
  saveNickname: (nickname: string) => Promise<boolean>;
  // joinServer connects to address and shows its game list.
  joinServer: (address: string) => Promise<void>;
  // lan is the last search for games on the local network (N-13).
  lan: { searching: boolean; found: LanServer[]; error: string };
  findLanGames: () => Promise<void>;
  // host runs the server in this app and opens a game on it; sets
  // defaults to every set this app has.
  host: (game: Omit<Create, 'sets'> & { sets?: string[] }) => Promise<void>;
  // hostSaved runs the server in this app and shows its games and saved
  // games, to continue one (N-11).
  hostSaved: () => Promise<void>;
  // saveGame saves the running game; only the host may (N-11).
  saveGame: () => void;
  // joinGame sits down at a game; a locked one needs its password.
  joinGame: (game: string, password?: string) => void;
  loadSave: (save: string) => void;
  setReady: (ready: boolean) => void;
  leaveTable: () => Promise<void>;
  // quit closes the connection, stops hosting, and goes to the menu.
  quit: () => Promise<void>;
}

// seatOf is this player's seat: from the table, or, for an app that came
// back to a running game after a restart, from the welcome.
export function seatOf(s: Pick<Session, 'lobby'>): number {
  return s.lobby.table?.you ?? Math.max(s.lobby.welcome?.seat ?? 0, 0);
}

// applyTheme colours menus and panels (ST-02).
// Unit tests run without a document.
function applyTheme(s: Settings): void {
  if (typeof document !== 'undefined') {
    document.documentElement.dataset.theme = s.theme;
  }
}

// messageOf turns a rejected call into a message for the player. Wails
// rejects with Go's error string.
export function messageOf(e: unknown): string {
  if (typeof e === 'string') {
    return e;
  }
  if (e instanceof Error) {
    return e.message;
  }
  return 'something went wrong';
}

// dropConnection is for tests: it fails the connection as the network
// would, so reconnecting can be checked.
export const dropConnection = drop;

export const useApp = create<AppState>()((set, get) => {
  let effectKey = 0;
  // names reads players from an update and cards from the card data.
  const names = (u: Update) => ({
    player: (p: number) => u.seats?.find((s) => s.seat === p)?.name ?? `Player ${String(p + 1)}`,
    card: (ref: string) => get().cards.get(ref)?.name ?? ref,
  });
  const playback = createPlayback({
    show: (u) => {
      const s = get();
      const entries = entriesFor(
        u.events ?? [],
        s.game?.view.turn.number ?? u.view.turn.number,
        seatOf(s),
        names(u),
        (s.history.at(-1)?.key ?? 0) + 1,
      );
      const looked = lookedAt(u.events ?? [], seatOf(s));
      set({
        game: u,
        behind: playback.pending() > 0,
        history: [...s.history, ...entries],
        ...(looked.length > 0 ? { looking: looked } : {}),
      });
    },
    run: (effect, scale) =>
      new Promise((resolve) => {
        const key = ++effectKey;
        const done = () => {
          window.clearTimeout(timer);
          if (get().playing?.key === key) {
            set({ playing: null });
          }
          resolve();
        };
        // An effect that never ends (its screen closed) must not block.
        const timer = window.setTimeout(done, 3000 * scale + 1000);
        set({ playing: { effect, scale, key, done } });
      }),
    effects: (u) =>
      effectsFor(u.events ?? [], {
        card: (ref) => get().cards.get(ref),
        name: names(u).player,
      }),
    speed: () => get().settings.animation,
    shown: () => get().game,
  });
  // keepSeat keeps the table's seat token in the settings, so the app
  // can take the seat back after a restart (N-08).
  const keepSeat = () => {
    const s = get();
    const r = { address: s.address, token: s.lobby.welcome?.token ?? '', game: s.lobby.table?.game ?? '' };
    const old = s.settings.rejoin;
    if (r.token && (r.address !== old.address || r.token !== old.token || r.game !== old.game)) {
      void s.saveSettings({ ...s.settings, rejoin: r });
    }
  };

  // stopPlayback drops queued updates and the effect on screen.
  const stopPlayback = () => {
    playback.clear();
    get().playing?.done();
    set({ history: [], looking: [] });
  };

  // open connects to address and says hello; the server answers with
  // welcome and its game list.
  // startHost runs the server in this app; null when it cannot.
  async function startHost(): Promise<HostInfo | null> {
    try {
      const info = await startHosting();
      set({ hosting: info });
      return info;
    } catch (e) {
      set({ error: messageOf(e) });
      return null;
    }
  }

  async function open(address: string, reconnect?: string): Promise<boolean> {
    // A restarted app still holds the seat of its last game there.
    const rejoin = get().settings.rejoin;
    const token = reconnect ?? (rejoin.address === address && rejoin.token ? rejoin.token : undefined);
    const url = wsURL(address);
    if (!url) {
      set({ error: 'Type an address like 192.168.1.5 or 192.168.1.5:4774.' });
      return false;
    }
    try {
      const [versions, sets] = await Promise.all([getVersions(), getCardSets()]);
      await connect(url, get().receive, () => {
        // A game keeps its table and history, to reconnect to.
        if (get().screen !== 'game') {
          stopPlayback();
        }
        set((s) => ({ ...lost(s), connected: false }));
      });
      // A reconnect keeps the table it comes back to.
      set(token ? { connected: true, address, error: '' } : { connected: true, address, cardSets: sets, lobby: emptyLobby, error: '' });
      const hello: Hello = {
        protocol: protocolVersion,
        app: versions.app,
        name: get().settings.nickname,
        role: 'player',
        sets,
        ...(token ? { token } : {}),
      };
      send('hello', hello);
      return true;
    } catch (e) {
      set({ error: messageOf(e) });
      return false;
    }
  }

  return {
    playing: null,
    history: [],
    looking: [],
    endLook: () => {
      set({ looking: [] });
    },
    behind: false,
    sent: null,
    selected: null,
    receive: (env) => {
      if (env.type === 'table' || env.type === 'saved') {
        stopPlayback();
      }
      set((s) => onMessage(s, env));
      if (env.type === 'update') {
        set({ behind: true });
        playback.push(env.data as Update);
      }
      if (env.type === 'error') {
        set({ sent: null });
      }
      if (env.type === 'table') {
        keepSeat();
      }
    },
    act: (intent) => {
      send('intent', { intent });
      set({ sent: get().game?.step ?? null, selected: null, error: '' });
    },
    skipAll: (on) => {
      send('skip_all', { on });
    },
    address: '',
    reconnect: async () => {
      const token = get().lobby.welcome?.token;
      if (token) {
        await open(get().address, token);
      }
    },
    vote: (kick) => {
      send('vote', { kick });
    },
    saveRecord: async () => {
      const s = get();
      const link = recordURL(s.address, s.lobby.table?.game ?? s.settings.rejoin.game, s.lobby.welcome?.token ?? '');
      if (!link) {
        return '';
      }
      try {
        return await saveRecord(link);
      } catch (e) {
        set({ error: messageOf(e) });
        return '';
      }
    },
    select: (id) => {
      set({ selected: id });
    },
    screen: 'splash',
    lobby: emptyLobby,
    game: null,
    error: '',
    settings: {
      nickname: '',
      recent: [],
      animation: 'normal',
      fullscreen: false,
      theme: 'dark',
      mat: '',
      dice: 'tear',
      rejoin: { address: '', token: '', game: '' },
      host: noHostSetup,
    },
    hosting: null,
    lan: { searching: false, found: [], error: '' },
    connected: false,
    cardSets: [],
    cards: new Map(),
    dice: [],
    mats: [],
    backs: new Map(),
    icons: { coin: '' },
    zoom: null,
    loadCards: async () => {
      try {
        const [list, dice, mats, backs, icons, fonts] = await Promise.all([
          getCards(),
          getDice(),
          getMats(),
          getBacks(),
          getIcons(),
          getFonts(),
        ]);
        loadFonts(fonts);
        set({
          cards: new Map(list.map((c) => [c.ref, c])),
          dice,
          mats,
          backs: new Map(backs.filter((b) => b.image !== '').map((b) => [b.kind, b.image])),
          icons,
        });
      } catch (e) {
        set({ error: `Card data: ${messageOf(e)}` });
      }
    },
    setZoom: (zoom) => {
      set({ zoom });
    },
    show: (screen) => {
      set({ screen, error: '', zoom: null });
    },
    loadSettings: async () => {
      try {
        const s = await getSettings();
        set({ settings: s });
        applyTheme(s);
      } catch {
        // Keep the defaults; Go already falls back to them.
      }
    },
    saveSettings: async (s) => {
      try {
        const before = get().settings;
        const saved = await saveSettings(s);
        set({ settings: saved, error: '' });
        // Settings apply at once, without a restart.
        applyTheme(saved);
        if (saved.fullscreen !== before.fullscreen) {
          setFullscreen(saved.fullscreen);
        }
        return true;
      } catch (e) {
        set({ error: messageOf(e) });
        return false;
      }
    },
    saveNickname: async (nickname) => {
      const s = get().settings;
      return nickname === s.nickname || get().saveSettings({ ...s, nickname });
    },
    findLanGames: async () => {
      if (get().lan.searching) {
        return;
      }
      set({ lan: { ...get().lan, searching: true } });
      try {
        set({ lan: { searching: false, found: await findGames(), error: '' } });
      } catch (e) {
        set({ lan: { searching: false, found: [], error: messageOf(e) } });
      }
    },
    joinServer: async (address) => {
      if (await open(address)) {
        const s = get().settings;
        void get().saveSettings({ ...s, recent: remember(s.recent, address) });
      }
    },
    host: async (game) => {
      const info = await startHost();
      if (info && (await open(`127.0.0.1:${String(info.port)}`))) {
        send('create', { ...game, sets: game.sets ?? get().cardSets });
      }
    },
    hostSaved: async () => {
      const info = await startHost();
      if (info && (await open(`127.0.0.1:${String(info.port)}`))) {
        set({ screen: 'join' });
      }
    },
    saveGame: () => {
      send('save');
    },
    joinGame: (game, password) => {
      send('join', password ? { game, password } : { game });
    },
    loadSave: (save) => {
      send('load', { save });
    },
    setReady: (ready) => {
      send('ready', { ready });
    },
    leaveTable: async () => {
      if (get().hosting) {
        await get().quit();
        return;
      }
      send('leave');
      stopPlayback();
      set((s) => leftTable(s));
    },
    quit: async () => {
      disconnect();
      stopPlayback();
      const hosting = get().hosting;
      set({ connected: false, hosting: null, lobby: emptyLobby, game: null, screen: 'menu', error: '' });
      if (hosting) {
        await stopHosting().catch(() => undefined);
      }
    },
  };
});
