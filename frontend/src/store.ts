// The app's one store (ADR 008). Screens read it; only bridge results
// and server messages write it.
import { create } from 'zustand';
import {
  setFullscreen,
  getCards,
  getCardSets,
  getSettings,
  getVersions,
  saveSettings,
  startHosting,
  stopHosting,
  type HostInfo,
  type Settings,
} from './bridge';
import type { CardInfo } from './cards';
import { entriesFor, type Entry } from './history';
import { emptyLobby, leftTable, lost, onMessage, remember, type Screen, type Session } from './lobby';
import { connect, disconnect, send } from './net';
import { createPlayback, effectsFor, type Effect } from './playback';
import { protocolVersion, wsURL, type Create, type Envelope, type Hello, type Intent, type Update } from './protocol';

export type { Screen } from './lobby';

// Zoom is the card shown large: on hover, or after a long press on touch
// (B-06). side is where it shows, away from the pointer.
export interface Zoom {
  ref: string;
  touch: boolean;
  side: 'left' | 'right';
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
  // host runs the server in this app and opens a game on it.
  host: (game: Omit<Create, 'sets'>) => Promise<void>;
  joinGame: (game: string) => void;
  loadSave: (save: string) => void;
  setReady: (ready: boolean) => void;
  leaveTable: () => Promise<void>;
  // quit closes the connection, stops hosting, and goes to the menu.
  quit: () => Promise<void>;
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
        s.lobby.table?.you ?? -1,
        names(u),
        (s.history.at(-1)?.key ?? 0) + 1,
      );
      set({ game: u, behind: playback.pending() > 0, history: [...s.history, ...entries] });
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
  // stopPlayback drops queued updates and the effect on screen.
  const stopPlayback = () => {
    playback.clear();
    get().playing?.done();
    set({ history: [] });
  };

  // open connects to address and says hello; the server answers with
  // welcome and its game list.
  async function open(address: string): Promise<boolean> {
    const url = wsURL(address);
    if (!url) {
      set({ error: 'Type an address like 192.168.1.5 or 192.168.1.5:4774.' });
      return false;
    }
    try {
      const [versions, sets] = await Promise.all([getVersions(), getCardSets()]);
      await connect(url, get().receive, () => {
        stopPlayback();
        set((s) => ({ ...lost(s), connected: false }));
      });
      set({ connected: true, cardSets: sets, lobby: emptyLobby, error: '' });
      const hello: Hello = {
        protocol: protocolVersion,
        app: versions.app,
        name: get().settings.nickname,
        role: 'player',
        sets,
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
    },
    act: (intent) => {
      send('intent', { intent });
      set({ sent: get().game?.step ?? null, selected: null, error: '' });
    },
    skipAll: (on) => {
      send('skip_all', { on });
    },
    select: (id) => {
      set({ selected: id });
    },
    screen: 'splash',
    lobby: emptyLobby,
    game: null,
    error: '',
    settings: { nickname: '', recent: [], animation: 'normal', fullscreen: false, theme: 'dark', mat: 'felt' },
    hosting: null,
    connected: false,
    cardSets: [],
    cards: new Map(),
    zoom: null,
    loadCards: async () => {
      try {
        const list = await getCards();
        set({ cards: new Map(list.map((c) => [c.ref, c])) });
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
    joinServer: async (address) => {
      if (await open(address)) {
        const s = get().settings;
        void get().saveSettings({ ...s, recent: remember(s.recent, address) });
      }
    },
    host: async (game) => {
      let info: HostInfo;
      try {
        info = await startHosting();
      } catch (e) {
        set({ error: messageOf(e) });
        return;
      }
      set({ hosting: info });
      if (await open(`127.0.0.1:${String(info.port)}`)) {
        send('create', { ...game, sets: get().cardSets });
      }
    },
    joinGame: (game) => {
      send('join', { game });
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
