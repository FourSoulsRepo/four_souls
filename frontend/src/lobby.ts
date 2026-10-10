// The lobby state and how server messages change it. Pure, so it is
// tested without a socket.
import type { Envelope, GameInfo, Games, Options, ProtocolError, Saved, SaveInfo, Table, Update, Welcome } from './protocol';
import { errorText } from './protocol';

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
// just before the close says more than "lost".
export function lost(s: Session): Session {
  return { screen: 'join', lobby: emptyLobby, game: null, error: s.error || 'Lost the connection to the server.' };
}

// remember puts an address first in the recent list, once.
export function remember(recent: string[], address: string, max = 5): string[] {
  const a = address.trim();
  return [a, ...recent.filter((r) => r !== a)].slice(0, max);
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
  return parts.join(' · ');
}
