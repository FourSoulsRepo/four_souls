// The app's one store (ADR 008). Screens read it; only bridge results
// and server messages write it.
import { create } from 'zustand';
import type { Choice } from './actions';
import {
  findGames,
  setFullscreen,
  getCards,
  getDice,
  getBacks,
  getFonts,
  getSounds,
  getIcons,
  getMats,
  getCardSets,
  saveRecord,
  getSettings,
  getVersions,
  checkSeat,
  saveSettings,
  serverStatus,
  startHosting,
  stopHosting,
  type DiceSkin,
  type HostInfo,
  type Icons,
  type LanServer,
  type MatInfo,
  type Rejoin,
  type Settings,
} from './bridge';
import { noHostSetup } from './hostSetup';
import { inLanguage, type CardInfo, type Kind } from './cards';
import { loadFonts } from './fonts';
import { language, say, setLanguage, t } from './i18n';
import { entriesFor, lookedAt, moveOf, restartCut, rolledBack, voteEnded, type Entry, type Looked, type Move, type Outcome } from './history';
import type { LinePart } from './stack';
import type { Box } from './zoom';
import { emptyLobby, forgotten, leftTable, lost, onMessage, remember, renamed, type Screen, type ServerUp, type Session } from './lobby';
import { connectFirst, disconnect, drop, send } from './net';
import { claimEffect, nextBeat } from './animations';
import { lastRollAfter, type LastRoll } from './dice';
import { createPlayback, effectsFor, scaleFor, slotZone, type Death, type Effect } from './playback';
import { loadSounds } from './sounds';
import { endsGame, lostDialogAfter, retryDelay } from './reconnect';
import {
  httpBases,
  plainNames,
  protocolVersion,
  recordURL,
  wsURLs,
  type Create,
  type Envelope,
  type Hello,
  type Intent,
  type Msg,
  type Notice,
  type OptionView,
  type Presence,
  type Update,
  type Welcome,
} from './protocol';

export type { Screen } from './lobby';

// Zoom is the card shown large: on hover, or after a long press on touch
// (B-06). side is where it shows, away from the pointer.
export interface Zoom {
  ref: string;
  touch: boolean;
  // at is the hovered card or row; the zoom opens beside it (15.27).
  at: Box;
  // rise: a hand card's zoom rises from it instead (H-13).
  rise?: boolean;
  // why says why the player can't use the card now (15.29, F-17).
  why?: string;
  // caption goes below the card: who owns it, and the used ability's
  // words, highlighted (S-06, S-07).
  caption?: ZoomCaption;
  // outcome: what a history headline did, beside the card (15.90).
  outcome?: Outcome[];
  // stats: why the card's stats are not its print, or what an effect of
  // this turn does, below the card (15.92, V-41).
  stats?: string[];
}

export interface ZoomCaption {
  owner?: string;
  color?: string;
  ability?: string;
  // line says what a stack item without ability words did: whom an
  // attack roll was against, who dealt damage to whom (15.22).
  line?: LinePart[];
}

// Aiming is a click-started aim (15.25): the card, its play or use,
// the targets it may pick, and why anything else is wrong (15.33).
// A card already placed aims at its target prompt's options instead
// of a play or use (15.34).
export interface Aiming {
  source: number;
  intent?: Intent;
  choice?: Choice;
  targets: OptionView[];
  reason?: string;
}

// Why is a reason shown where the player clicked (15.29, F-18): a
// wrong aim, or a card or action they can't use now. n restarts its
// fade when the same reason shows again.
export interface Why {
  text: string;
  x: number;
  y: number;
  n: number;
}

// Aimed is a played card's target, kept until the target prompt.
export interface Aimed {
  step: number;
  target: OptionView;
}

// Playing is the effect on screen now; done ends it (playback.ts).
export interface Playing {
  effect: Effect;
  scale: number;
  key: number;
  done: () => void;
}

// Placed is a card a place effect left on its slot, at a screen point,
// until the next view is shown (19.3).
export interface Placed {
  key: number;
  card: string;
  x: number;
  y: number;
}

interface AppState extends Session {
  // playing is the event effect on screen, if any.
  playing: Playing | null;
  placed: Placed[];
  // left is the sources of stack items whose card flew away; they hide
  // until the next view is shown (19.41).
  left: number[];
  // lifted is the table cards a dead card's flight took off their slot;
  // they hide until the next view is shown (19.73).
  lifted: number[];
  // history is the match so far, as this viewer saw it (7.9).
  history: Entry[];
  // moves name what each update did, for the Undo button (17.7); a
  // game without a rollback mode keeps none.
  moves: Move[];
  // lastRoll is the last roll shown, for the small die by the stack (15.63).
  lastRoll: LastRoll | null;
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
  // aiming is the card the player aims with; its targets glow (15.25).
  aiming: Aiming | null;
  setAiming: (aiming: Aiming | null) => void;
  // why is the reason shown by the last refused click, if any.
  why: Why | null;
  showWhy: (text: string, x: number, y: number) => void;
  endWhy: () => void;
  // aimed is the target a card was aimed at, at the step its play was
  // sent; the target prompt that follows picks it (15.25).
  aimed: Aimed | null;
  // playAt sends a play and remembers its target.
  playAt: (intent: Intent, target: OptionView) => void;
  // takeAimed returns aimed and forgets it.
  takeAimed: () => Aimed | null;
  // waitTurn turns "wait for my turn" on or off (N-16).
  waitTurn: (on: boolean) => void;
  // rollBack asks to undo a step and every step after it (ADR 020). In
  // Soft it opens a vote; ask is the step's author asking the others.
  rollBack: (step: number, ask?: boolean) => void;
  // rollbackVote answers the open Soft vote (ADR 020, item 5).
  rollbackVote: (yes: boolean) => void;
  // rollbackCancel ends the vote this player opened (item 5.13).
  rollbackCancel: () => void;
  // jumped is set while the table shown came from a rollback: it shows
  // at once, with no animation (item 8.5).
  jumped: boolean;
  // presence is what the other seats point at now, by seat (15.36).
  presence: Record<number, Presence>;
  // sendPresence shares what this player points at (ADR 016).
  sendPresence: (p: Omit<Presence, 'seat'>) => void;
  // address is the server this app is connected to, for reconnecting.
  address: string;
  // reconnect says hello again with the seat's token (N-08).
  reconnect: () => Promise<void>;
  // reconnecting is set while the app reconnects by itself to a game
  // whose connection dropped (N-20); lostLong once it took long enough
  // for the "Connection lost" dialog.
  reconnecting: boolean;
  lostLong: boolean;
  // restarting is set when the server said it restarts (N-20).
  restarting: boolean;
  // notice is the server notice to show; endNotice closes it (N-20).
  notice: Notice | null;
  endNotice: () => void;
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
  // playOpen: the main menu's Play group shows (M-06); it stays open
  // while the player is on its screens.
  playOpen: boolean;
  setPlayOpen: (open: boolean) => void;
  // path is the Play screen the player came by (M-03): its join screen
  // shows the LAN games or the Online servers.
  path: Path;
  openPath: (path: Path) => void;
  // rejoinable is the saved seat's game while its server says the seat
  // is still on: the menu offers Rejoin (15.99). null: none.
  rejoinable: { host: string; via: Rejoin['via'] } | null;
  // checkRejoin asks the saved seat's server whether it is still on; a
  // dead seat is forgotten, no answer keeps it for later.
  checkRejoin: () => Promise<void>;
  // rejoin takes the saved seat back, by the path it was joined by.
  rejoin: () => Promise<void>;
  loadSettings: () => Promise<void>;
  // saveSettings reports whether Go stored them.
  saveSettings: (s: Settings) => Promise<boolean>;
  // saveNickname saves the nickname if it changed (ST-06).
  saveNickname: (nickname: string) => Promise<boolean>;
  // joinServer connects to address and shows its game list. The
  // address goes to Recent unless recent is false (the LAN screen); a
  // token takes that seat back (15.97).
  joinServer: (address: string, opts?: { recent?: boolean; token?: string }) => Promise<void>;
  // forgetServer drops an address from the Recent list (N-18).
  forgetServer: (address: string) => Promise<void>;
  // recentStatus is the last check of each Recent server, by address
  // (15.100); an address is missing while it is asked.
  recentStatus: Partial<Record<string, ServerUp>>;
  // checkRecent asks every Recent server at once whether it is up; each
  // answer fills its row as it comes.
  checkRecent: () => Promise<void>;
  // lan is the last search for games on the local network (N-13).
  lan: { searching: boolean; found: LanServer[]; error: string };
  findLanGames: () => Promise<void>;
  // host runs the server in this app and opens a game on it; sets
  // defaults to every set this app has.
  host: (game: Omit<Create, 'sets'> & { sets?: string[] }) => Promise<void>;
  // createGame opens a game on the connected server, which runs
  // elsewhere (N-17); sets defaults to every set this app has.
  createGame: (game: Omit<Create, 'sets'> & { sets?: string[] }) => void;
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
  // quit closes the connection, stops hosting, and goes to the menu,
  // or to another screen: the join screen, one level up from its games.
  quit: (to?: Screen) => Promise<void>;
}

// Path is a way to play from the main menu's Play group (M-03).
export type Path = 'lan' | 'online';

// noSeat is a forgotten saved seat.
const noSeat: Rejoin = { address: '', token: '', game: '', via: '' };

// seatOf is this player's seat: from the table, or, for an app that came
// back to a running game after a restart, from the welcome.
export function seatOf(s: Pick<Session, 'lobby'>): number {
  return s.lobby.table?.you ?? Math.max(s.lobby.welcome?.seat ?? 0, 0);
}

// messageOf turns a rejected call into a message for the player. Wails
// rejects with Go's error string: a message as JSON when Go has words
// for it (ADR 021), else the error's own words.
export function messageOf(e: unknown): string {
  const text = typeof e === 'string' ? e : e instanceof Error ? e.message : undefined;
  if (text === undefined) {
    return t('error.unknown');
  }
  const m = goMsgs(text)?.[0];
  return m ? say(m, plainNames) : text;
}

// deckRatiosText says why the deck ratios can't be used (GS-07): Go's
// head line, with the short groups' lines joined into it.
export function deckRatiosText(why: string): string {
  const msgs = goMsgs(why) ?? [];
  const head = msgs.at(0);
  if (!head) {
    return why;
  }
  const list = msgs
    .slice(1)
    .map((g) => say(g, plainNames))
    .join(', ');
  return say({ ...head, args: { ...head.args, groups: list } }, plainNames);
}

// goMsgs reads messages Go sent as JSON text: one, or a list; else null.
export function goMsgs(text: string): Msg[] | null {
  if (!text.startsWith('{"key"') && !text.startsWith('[{"key"')) {
    return null;
  }
  try {
    const v: unknown = JSON.parse(text);
    const list = Array.isArray(v) ? v : [v];
    return list.filter((m): m is Msg => typeof m === 'object' && m !== null && typeof (m as Msg).key === 'string');
  } catch {
    return null;
  }
}

// presenceLife is how long a presence shows without being renewed, in
// ms; presenceRenew is how often a sender renews a still one.
export const presenceLife = 5000;
export const presenceRenew = 2000;

// viewScale is how long the table's own animations take as the view
// changes (a card's pop-in, its effect icons, the moving rows): 0 when
// animations are off, or while a rollback's table shows (ADR 020).
export function viewScale(s: { jumped: boolean; settings: Pick<Settings, 'animation'> }): number {
  return s.jumped ? 0 : scaleFor(s.settings.animation, 0, 0);
}

// dropConnection is for tests: it fails the connection as the network
// would, so reconnecting can be checked.
export const dropConnection = drop;

// cardsInLanguage names the cards in the language t() uses now.
function cardsInLanguage(cards: Map<string, CardInfo>): Map<string, CardInfo> {
  return new Map([...cards].map(([ref, c]) => [ref, inLanguage(c, language())]));
}

export const useApp = create<AppState>()((set, get) => {
  let effectKey = 0;
  // names reads players from an update and cards from the card data.
  const names = (u: Update) => ({
    player: (p: number) => u.seats?.find((s) => s.seat === p)?.name ?? t('table.player_n', { n: p + 1 }),
    card: (ref: string) => get().cards.get(ref)?.name ?? ref,
  });
  // diedAt is where each dead card died, by its held object, until it
  // flies to the discard or a soul (19.73).
  const diedAt = new Map<number, Death>();
  const playback = createPlayback({
    begin: nextBeat,
    show: (u) => {
      const s = get();
      // A rollback greys the undone entries and adds its row (ADR 020).
      // A vote that ended with no rollback gets a row too (item 5).
      // A server restart that cut the history gets a row (item 10.6).
      const kept = u.rollback_cut ? restartCut(s.history, u.rollback_cut, u.view.turn.number) : s.history;
      const history = u.rollback
        ? rolledBack(kept, u.rollback, u.step, u.view.turn.number, seatOf(s), names(u))
        : u.rollback_end
          ? voteEnded(kept, u.rollback_end, u.view.turn.number, seatOf(s), names(u))
          : kept;
      const entries = entriesFor(
        u.events ?? [],
        s.game?.view.turn.number ?? u.view.turn.number,
        seatOf(s),
        names(u),
        (history.at(-1)?.key ?? 0) + 1,
        u.view,
      );
      // Each entry keeps the first step of its update: a rollback's target.
      const first = s.game ? s.game.step + 1 : u.step;
      // The moves, for the Undo button: a rollback forgets those it
      // undid (17.7).
      const r = u.rollback;
      const moves = u.rollback_mode
        ? [...s.moves.filter((m) => !r || m.from < r.to), moveOf(u.events ?? [], entries, first, u.step, seatOf(s), names(u))]
        : [];
      const looked = lookedAt(u.events ?? [], seatOf(s));
      // A rollback jumps to the restored table: no animation, and what
      // the undone steps left on screen goes (item 8.5).
      const jump = u.rollback
        ? { lastRoll: null, looking: [], aiming: null, aimed: null, selected: null, sent: null, zoom: null }
        : { lastRoll: lastRollAfter(s.lastRoll, u.events ?? []), ...(looked.length > 0 ? { looking: looked } : {}) };
      set({
        game: u,
        placed: [],
        left: [],
        lifted: [],
        behind: playback.pending() > 0,
        jumped: u.rollback !== undefined,
        history: [...history, ...entries.map((e) => ({ ...e, step: first }))],
        moves,
        ...jump,
      });
    },
    run: (effect, scale) =>
      new Promise((resolve) => {
        // A card that leaves the stack is gone from it, even when its
        // fly is skipped (19.41).
        const leaves = effect.kind === 'fly' ? effect.leaves : undefined;
        if (leaves !== undefined) {
          set((s) => ({ left: [...s.left, leaves] }));
        }
        const lifts = effect.kind === 'fly' ? effect.lifts : undefined;
        if (lifts !== undefined) {
          set((s) => ({ lifted: [...s.lifted, lifts] }));
        }
        // A clashing effect is logged and skipped (19.6).
        if (!claimEffect(effect)) {
          resolve();
          return;
        }
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
    // A rollback plays no effects (ADR 020, item 8.5); the dead cards
    // it undid are forgotten.
    effects: (u) => {
      if (u.rollback) {
        diedAt.clear();
        return [];
      }
      return effectsFor(u.events ?? [], {
        card: (ref) => get().cards.get(ref),
        name: names(u).player,
        character: (p) => u.view.players.find((x) => x.id === p)?.character.id,
        slot: (id) => slotZone(u.view, id),
        before: (id) => {
          const view = get().game?.view;
          return view ? slotZone(view, id) : undefined;
        },
        diedAt,
        seen: (ref) => (get().game?.events ?? []).some((e) => e.kind === 'card_revealed' && e.card === ref),
        roll: (stack) => get().game?.view.stack?.find((it) => it.id === stack)?.roll,
        you: seatOf(get()),
      });
    },
    speed: () => get().settings.animation,
    shown: () => get().game,
  });
  // keepSeat keeps the table's seat token in the settings, so the app
  // can take the seat back after a restart (N-08).
  const keepSeat = () => {
    const s = get();
    const r = { address: s.address, token: s.lobby.welcome?.token ?? '', game: s.lobby.table?.game ?? '', via: s.path };
    const old = s.settings.rejoin;
    if (r.token && (r.address !== old.address || r.token !== old.token || r.game !== old.game || r.via !== old.via)) {
      void s.saveSettings({ ...s.settings, rejoin: r });
    }
  };
  // endedGame is the forgotten seat's game, for its record download.
  let endedGame = '';
  // forgetSeat drops the saved seat: its game ended, or the player left
  // it on purpose (15.99). A crash or a closed window keeps it.
  const forgetSeat = () => {
    const s = get().settings;
    set({ rejoinable: null });
    if (s.rejoin.token) {
      endedGame = s.rejoin.game;
      void get().saveSettings({ ...s, rejoin: noSeat });
    }
  };
  // seated: the player sits at a table or plays a game.
  const seated = () => get().screen === 'game' || get().lobby.table !== null;

  // serverURL is the WebSocket URL that opened for address: wss:// or
  // ws:// (8.18). The record download follows its scheme.
  let serverURL = '';
  // remembering is the save that put the joined address in Recent; the
  // server's name waits for it, so neither save undoes the other.
  let remembering: Promise<unknown> = Promise.resolve();
  // statusRound counts the Recent checks; a late answer of an older one
  // is dropped.
  let statusRound = 0;
  // nameServer keeps the name from the server's welcome on its Recent
  // entry (N-18).
  const nameServer = (address: string, name: string) => {
    void remembering.then(() => {
      const s = get().settings;
      const recent = renamed(s.recent, address, name);
      if (recent !== s.recent) {
        void get().saveSettings({ ...s, recent });
      }
    });
  };

  // A seat's presence goes when it is not renewed in time: a lost
  // message never leaves an arrow behind. Senders renew it sooner.
  const presenceTimers = new Map<number, number>();
  const forget = (seat: number) => {
    window.clearTimeout(presenceTimers.get(seat));
    presenceTimers.delete(seat);
    set((s) => ({ presence: Object.fromEntries(Object.entries(s.presence).filter(([k]) => Number(k) !== seat)) }));
  };
  const clearPresence = () => {
    for (const t of presenceTimers.values()) {
      window.clearTimeout(t);
    }
    presenceTimers.clear();
    set({ presence: {} });
  };

  // stopPlayback drops queued updates and the effect on screen.
  const stopPlayback = () => {
    playback.clear();
    get().playing?.done();
    clearPresence();
    diedAt.clear();
    set({ history: [], moves: [], lastRoll: null, looking: [], placed: [], left: [], lifted: [] });
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

  // retry reconnects by itself to a game whose connection dropped: at
  // once, then ever less often; the badge turns into the "Connection
  // lost" dialog after a while (N-20, 8.24).
  const retry = { timer: 0, dialog: 0, attempt: 0, busy: false };
  const retryLater = () => {
    if (!get().reconnecting) {
      retry.attempt = 0;
      window.clearTimeout(retry.dialog);
      retry.dialog = window.setTimeout(() => {
        set({ lostLong: true });
      }, lostDialogAfter());
      set({ reconnecting: true, lostLong: false });
    }
    window.clearTimeout(retry.timer);
    retry.timer = window.setTimeout(() => void retryNow(), retryDelay(retry.attempt++));
  };
  const retryNow = async () => {
    const token = get().lobby.welcome?.token;
    if (!get().reconnecting || retry.busy || !token) {
      return;
    }
    window.clearTimeout(retry.timer);
    retry.busy = true;
    const ok = await open(get().address, token);
    retry.busy = false;
    if (!get().reconnecting) {
      return; // the player left meanwhile
    }
    if (ok) {
      endRetry();
    } else {
      retryLater();
    }
  };
  const endRetry = () => {
    if (get().reconnecting) {
      window.clearTimeout(retry.timer);
      window.clearTimeout(retry.dialog);
    }
    set({ reconnecting: false, lostLong: false, restarting: false });
  };

  async function open(address: string, reconnect?: string): Promise<boolean> {
    // A restarted app still holds the seat of its last game there.
    const rejoin = get().settings.rejoin;
    const token = reconnect ?? (rejoin.address === address && rejoin.token ? rejoin.token : undefined);
    const urls = wsURLs(address);
    if (!urls) {
      set({ error: t('error.bad_address') });
      return false;
    }
    try {
      const [versions, sets] = await Promise.all([getVersions(), getCardSets()]);
      serverURL = await connectFirst(urls, get().receive, () => {
        // A game keeps its table and history, to reconnect to.
        if (get().screen !== 'game') {
          stopPlayback();
        }
        clearPresence();
        set((s) => ({ ...lost(s), connected: false }));
        const s = get();
        if (s.screen === 'game' && !s.game?.view.over && s.lobby.welcome?.token) {
          retryLater();
        }
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
        dice: get().settings.dice,
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
    placed: [],
    left: [],
    lifted: [],
    history: [],
    moves: [],
    lastRoll: null,
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
        if ((env.data as Update).view.over) {
          forgetSeat();
        }
      }
      if (env.type === 'error') {
        set({ sent: null, aimed: null });
      }
      if (env.type === 'table') {
        set({ aiming: null, aimed: null });
        keepSeat();
      }
      if (env.type === 'welcome') {
        nameServer(get().address, (env.data as Welcome).name ?? '');
      }
      if (env.type === 'restarting') {
        set({ restarting: true });
      }
      if (env.type === 'notice') {
        const n = env.data as Notice;
        set({ notice: n });
        if (endsGame(n)) {
          forgetSeat();
          endRetry();
          stopPlayback();
          set((s) => ({ screen: 'join', lobby: { ...s.lobby, table: null }, game: null }));
        }
      }
      if (env.type === 'presence') {
        const p = env.data as Presence;
        if (!p.from && !p.at && !p.point) {
          forget(p.seat);
          return;
        }
        window.clearTimeout(presenceTimers.get(p.seat));
        presenceTimers.set(
          p.seat,
          window.setTimeout(() => {
            forget(p.seat);
          }, presenceLife),
        );
        set((s) => ({ presence: { ...s.presence, [p.seat]: p } }));
      }
    },
    act: (intent) => {
      send('intent', { intent });
      set({ sent: get().game?.step ?? null, selected: null, error: '' });
    },
    aiming: null,
    setAiming: (aiming) => {
      set({ aiming, selected: null });
    },
    why: null,
    showWhy: (text, x, y) => {
      set((s) => ({ why: { text, x, y, n: (s.why?.n ?? 0) + 1 } }));
    },
    endWhy: () => {
      set({ why: null });
    },
    aimed: null,
    playAt: (intent, target) => {
      get().act(intent);
      set({ aiming: null, aimed: { step: get().game?.step ?? 0, target } });
    },
    takeAimed: () => {
      const { aimed } = get();
      set({ aimed: null });
      return aimed;
    },
    waitTurn: (on) => {
      send('wait_turn', { on });
    },
    rollBack: (step, ask) => {
      send('rollback', ask ? { step, ask } : { step });
      // Godmode's rollback is a new step; a Soft vote keeps the step,
      // so it would never clear sent.
      const g = get().game;
      set({ sent: g?.rollback_mode === 'god' ? g.step : get().sent, error: '' });
    },
    rollbackVote: (yes) => {
      send('rollback_vote', { yes });
    },
    rollbackCancel: () => {
      send('rollback_cancel', {});
    },
    jumped: false,
    presence: {},
    sendPresence: (p) => {
      if (get().lobby.welcome?.presence) {
        send('presence', p);
      }
    },
    address: '',
    reconnect: async () => {
      const token = get().lobby.welcome?.token;
      if (!token) {
        return;
      }
      if (get().reconnecting) {
        await retryNow();
        return;
      }
      await open(get().address, token);
    },
    reconnecting: false,
    lostLong: false,
    restarting: false,
    notice: null,
    endNotice: () => {
      set({ notice: null });
    },
    vote: (kick) => {
      send('vote', { kick });
    },
    saveRecord: async () => {
      const s = get();
      const link = recordURL(serverURL, s.lobby.table?.game ?? (s.settings.rejoin.game || endedGame), s.lobby.welcome?.token ?? '');
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
      mat: '',
      dice: 'tear',
      no_end_turn_ask: false,
      rejoin: { address: '', token: '', game: '', via: '' },
      host: noHostSetup,
      lan_address: '',
      language: '',
    },
    hosting: null,
    lan: { searching: false, found: [], error: '' },
    connected: false,
    cardSets: [],
    cards: new Map(),
    dice: [],
    mats: [],
    backs: new Map(),
    icons: { coin: '', sword: '', banner: '', sky: [] },
    zoom: null,
    loadCards: async () => {
      try {
        const [list, dice, mats, backs, icons, fonts, sounds] = await Promise.all([
          getCards(),
          getDice(),
          getMats(),
          getBacks(),
          getIcons(),
          getFonts(),
          getSounds(),
        ]);
        loadFonts(fonts);
        loadSounds(sounds);
        set({
          cards: new Map(list.map((c) => [c.ref, inLanguage(c, language())])),
          dice,
          mats,
          backs: new Map(backs.filter((b) => b.image !== '').map((b) => [b.kind, b.image])),
          icons,
        });
      } catch (e) {
        set({ error: t('error.card_data', { error: messageOf(e) }) });
      }
    },
    setZoom: (zoom) => {
      set({ zoom });
    },
    show: (screen) => {
      set({ screen, error: '', zoom: null });
    },
    playOpen: false,
    setPlayOpen: (open) => {
      set({ playOpen: open });
    },
    path: 'lan',
    openPath: (path) => {
      set({ path, screen: 'join', error: '', zoom: null });
    },
    rejoinable: null,
    checkRejoin: async () => {
      const r = get().settings.rejoin;
      const bases = r.token ? httpBases(r.address) : null;
      if (!bases) {
        set({ rejoinable: null });
        return;
      }
      try {
        const seat = await checkSeat(bases, r.token);
        if (get().settings.rejoin.token !== r.token) {
          return; // another seat was saved meanwhile
        }
        if (seat.live) {
          set({ rejoinable: { host: seat.host, via: r.via } });
        } else {
          forgetSeat();
        }
      } catch {
        set({ rejoinable: null });
      }
    },
    rejoin: async () => {
      const r = get().settings.rejoin;
      if (!r.token) {
        return;
      }
      set({ rejoinable: null });
      get().openPath(r.via || 'online');
      await get().joinServer(r.address, { token: r.token, recent: r.via !== 'lan' });
    },
    loadSettings: async () => {
      try {
        const s = await getSettings();
        setLanguage(s.language);
        set({ settings: s, cards: cardsInLanguage(get().cards) });
      } catch {
        // Keep the defaults; Go already falls back to them.
      }
    },
    saveSettings: async (s) => {
      try {
        const before = get().settings;
        const saved = await saveSettings(s);
        // The language too: t() reads it when the screens draw again.
        setLanguage(saved.language);
        set({ settings: saved, error: '', cards: cardsInLanguage(get().cards) });
        // Settings apply at once, without a restart.
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
    joinServer: async (address, opts) => {
      if ((await open(address, opts?.token)) && opts?.recent !== false) {
        const s = get().settings;
        remembering = get().saveSettings({ ...s, recent: remember(s.recent, address) });
      }
    },
    forgetServer: async (address) => {
      const s = get().settings;
      await get().saveSettings({ ...s, recent: forgotten(s.recent, address) });
    },
    recentStatus: {},
    checkRecent: async () => {
      const round = ++statusRound;
      set({ recentStatus: {} });
      await Promise.all(
        get().settings.recent.map(async ({ address }) => {
          const bases = httpBases(address);
          let up: ServerUp = 'down';
          try {
            up = bases ? await serverStatus(bases) : 'down';
          } catch {
            // no address answered: down
          }
          if (round === statusRound) {
            set({ recentStatus: { ...get().recentStatus, [address]: up } });
          }
        }),
      );
    },
    host: async (game) => {
      const info = await startHost();
      if (info && (await open(`127.0.0.1:${String(info.port)}`))) {
        send('create', { ...game, sets: game.sets ?? get().cardSets });
      }
    },
    createGame: (game) => {
      send('create', { ...game, sets: game.sets ?? get().cardSets });
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
      forgetSeat();
      stopPlayback();
      set((s) => leftTable(s));
    },
    quit: async (to = 'menu') => {
      if (seated()) {
        forgetSeat();
      }
      endRetry();
      disconnect();
      stopPlayback();
      const hosting = get().hosting;
      set({ connected: false, hosting: null, lobby: emptyLobby, game: null, screen: to, error: '' });
      if (hosting) {
        await stopHosting().catch(() => undefined);
      }
    },
  };
});
