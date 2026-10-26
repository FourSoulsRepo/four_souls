// The lobby state and how server messages change it. Pure, so it is
// tested without a socket.
import type { Envelope, GameInfo, Games, Options, ProtocolError, RollbackMode, Saved, SaveInfo, Table, Update, Welcome } from './protocol';
import { errorText } from './protocol';
import type { LanServer, RecentServer, Rejoin } from './bridge';

// 'cards' is a dev-only card gallery for tests (?screen=cards).
export type Screen = 'splash' | 'welcome' | 'menu' | 'settings' | 'host' | 'join' | 'table' | 'game' | 'cards';

export interface Lobby {
  welcome: Welcome | null;
  games: GameInfo[];
  saves: SaveInfo[];
  // table is the table the player sits at; null in the game list.
  table: Table | null;
  // saved is the ID the host saved the game under; '' if none.
  saved: string;
}

export const emptyLobby: Lobby = { welcome: null, games: [], saves: [], table: null, saved: '' };

// Session is the part of the store that server messages change.
export interface Session {
  screen: Screen;
  lobby: Lobby;
  // game is the update on screen; playback shows each update after its
  // effects (playback.ts). null before the game.
  game: Update | null;
  // error is the last failure to show the player; '' if none.
  error: string;
}

// onMessage returns the session after one server message.
export function onMessage(s: Session, env: Envelope): Session {
  const l = s.lobby;
  switch (env.type) {
    case 'welcome':
      return { ...s, lobby: { ...l, welcome: env.data as Welcome } };
    case 'games': {
      const g = env.data as Games;
      return { ...s, lobby: { ...l, games: g.games ?? [], saves: g.saves ?? [] } };
    }
    case 'table':
      return { screen: 'table', lobby: { ...l, table: env.data as Table, saved: '' }, game: null, error: '' };
    case 'setup':
      return s.screen === 'game' ? s : { ...s, screen: 'game', error: '' };
    case 'update':
      // The update itself goes to playback, which shows it in turn.
      return s.screen === 'game' ? s : { ...s, screen: 'game', error: '' };
    case 'saved':
      return { ...s, screen: 'join', lobby: { ...l, table: null, saved: (env.data as Saved).save }, game: null };
    case 'error':
      return { ...s, error: errorText(env.data as ProtocolError) };
    default:
      return s;
  }
}

// leftTable is the session after the player leaves the table: back to
// the server's game list.
export function leftTable(s: Session): Session {
  return { screen: 'join', lobby: { ...s.lobby, table: null }, game: null, error: '' };
}

// lost is the session after the connection dropped. A server error sent
// just before the close says more than "lost". During a game the table
// stays, so the player can reconnect to the same seat (N-08).
export function lost(s: Session): Session {
  const error = s.error || 'Lost the connection to the server.';
  if (s.screen === 'game') {
    return { ...s, error };
  }
  return { screen: 'join', lobby: emptyLobby, game: null, error };
}

// canCreate reports whether the Games list offers "Create a game": the
// server allows it (a dedicated one), and this app does not host. A
// friend of an app's host never opens games there (N-17).
export function canCreate(lobby: Lobby, hosting: boolean): boolean {
  return !hosting && lobby.welcome?.create === true;
}

// remember puts an address first in the recent list, once; it keeps
// the name the server gave before.
export function remember(recent: RecentServer[], address: string, max = 5): RecentServer[] {
  const a = address.trim();
  const name = recent.find((r) => r.address === a)?.name ?? '';
  return [{ address: a, name }, ...recent.filter((r) => r.address !== a)].slice(0, max);
}

// renamed gives the address's entry the name from the server's welcome
// (N-18). It returns the same list when nothing changes, so nothing is
// saved; an address not in the list stays out.
export function renamed(recent: RecentServer[], address: string, name: string): RecentServer[] {
  const a = address.trim();
  if (!recent.some((r) => r.address === a && r.name !== name)) {
    return recent;
  }
  return recent.map((r) => (r.address === a ? { ...r, name } : r));
}

// forgotten drops the address from the recent list (N-18).
export function forgotten(recent: RecentServer[], address: string): RecentServer[] {
  return recent.filter((r) => r.address !== address);
}

// serverLabel is how a recent server shows: its name, else its address.
export function serverLabel(r: RecentServer): string {
  return r.name || r.address;
}

// ServerUp is one check of a server (15.100): its open games while it
// answers, 'down' when no address of it answered.
export type ServerUp = { open: number } | 'down';

// upText is a Recent row's status: still asking, down, or how many
// games wait for players.
export function upText(up: ServerUp | undefined): string {
  if (up === undefined) {
    return 'Checking…';
  }
  if (up === 'down') {
    return 'Down';
  }
  return up.open === 1 ? '1 open game' : `${String(up.open)} open games`;
}

// optionsText sums up the match setup for the table.
export function optionsText(o: Options | undefined): string {
  const parts = [o?.picking === 'draft' ? `Draft of ${String(o.draft_size ?? 3)}` : 'Random characters'];
  if (o?.ban_rounds) {
    parts.push(`${String(o.ban_rounds)} ban round${o.ban_rounds > 1 ? 's' : ''}`);
  }
  if (o?.response_timer) {
    parts.push(`${String(o.response_timer)} s to respond`);
  }
  if (o?.no_bonus_souls) {
    parts.push('no bonus souls');
  }
  // The server turns the deck ratios off when the sets are too small.
  if (o && !o.no_deck_ratios) {
    parts.push('deck ratios');
  }
  if (o?.rollback) {
    parts.push(rollbackText(o.rollback));
  }
  return parts.join(' · ');
}

// rollbackText names a rollback mode for the game list and the table
// (ADR 020, item 1.4).
export function rollbackText(mode: RollbackMode): string {
  return mode === 'god' ? 'Godmode rollback' : 'Soft rollback';
}

// LanRow is one game found on the local network, as the LAN screen
// lists it (N-13).
export interface LanRow {
  key: string;
  address: string;
  text: string;
  joinable: boolean;
  // rejoin: a started game this app held a seat in (15.97); joining
  // takes the seat back.
  rejoin: boolean;
}

// lanRows lists the found games. A server that speaks another protocol
// version shows, but cannot be joined; neither can a full or started
// game, unless this app held a seat in it (rejoin is its saved seat).
export function lanRows(found: LanServer[], protocol: number, rejoin: Rejoin): LanRow[] {
  return found.flatMap((f) =>
    f.games.map((g, i) => {
      const same = f.protocol === protocol;
      const version = same ? '' : ` · needs app ${f.app}`;
      const mine = rejoin.token !== '' && g.id !== '' && g.id === rejoin.game;
      const open = !g.started && g.taken < g.seats;
      const state = g.started ? ' · In progress' : '';
      return {
        key: `${f.address}#${String(i)}`,
        address: f.address,
        text: `${g.locked ? '🔒 ' : ''}${g.host}'s game · ${String(g.taken)}/${String(g.seats)} players${state} · ${f.address}${version}`,
        joinable: same && (open || mine),
        rejoin: same && mine,
      };
    }),
  );
}
