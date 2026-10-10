// Network messages (ADR 006), mirroring internal/protocol/protocol.go.
// Only what the client uses so far; protocol.test.ts checks these types
// against Go's golden file. Go sends a nil list as null.

export const protocolVersion = 1;
export const defaultPort = 4774;

export interface Envelope {
  type: string;
  id?: number;
  data?: unknown;
}

export interface Hello {
  protocol: number;
  app: string;
  name: string;
  role: 'player' | 'spectator' | 'judge';
  sets: string[];
  token?: string;
}

export interface Welcome {
  protocol: number;
  app: string;
  engine: string;
  token: string;
  seat: number;
}

export interface ProtocolError {
  code: string;
  message: string;
  rule?: string;
}

export interface GameInfo {
  id: string;
  host: string;
  seats: number;
  taken: number;
  sets: string[] | null;
  started: boolean;
}

export interface SaveInfo {
  id: string;
  host: string;
  seats: string[] | null;
  saved: number; // Unix ms
  engine: string;
}

export interface Games {
  games: GameInfo[] | null;
  saves?: SaveInfo[];
}

// Options is the detailed match setup (GS-10).
export interface Options {
  picking?: 'random' | 'draft';
  draft_size?: number;
  ban_rounds?: number;
  ban_timer?: number; // seconds; 0: no timer
  host_bans?: string[];
  no_bonus_souls?: boolean;
  response_timer?: number; // seconds; 0: off, else at least 60
}

export interface Create {
  seats: number;
  sets: string[];
  options?: Options;
}

export interface Join {
  game: string;
}

export interface Ready {
  ready: boolean;
}

export interface Load {
  save: string;
}

export interface Seat {
  seat: number;
  name: string;
  connected: boolean;
  ready?: boolean;
  kicked?: boolean;
}

export interface Table {
  game: string;
  you: number;
  seats: Seat[] | null;
  sets: string[] | null;
  options?: Options;
  loaded?: string;
}

export interface Saved {
  save: string;
}

// wsURL turns what a player types ("192.168.1.5", "host:4774",
// "[::1]:4774") into the server's WebSocket URL, or null.
export function wsURL(address: string): string | null {
  const a = address.trim();
  const close = a.startsWith('[') ? a.indexOf(']') + 1 : 0; // IPv6 in brackets
  const colon = a.indexOf(':', close);
  const host = colon < 0 ? a : a.slice(0, colon);
  const portText = colon < 0 ? String(defaultPort) : a.slice(colon + 1);
  if (!/^(\[[0-9a-fA-F:.]+\]|[\w.-]+)$/.test(host) || !/^\d{1,5}$/.test(portText)) {
    return null;
  }
  const port = Number(portText);
  if (port < 1 || port > 65535) {
    return null;
  }
  return `ws://${host}:${String(port)}/ws`;
}

// errorText is what the player reads for a server error.
export function errorText(e: ProtocolError): string {
  switch (e.code) {
    case 'client_outdated':
      return 'This app is older than the server. Update the app to play here.';
    case 'server_outdated':
      return 'The server is older than this app. Ask the host to update.';
    case 'missing_sets':
      return `You are missing a card set: ${e.message}`;
    default:
      return e.message;
  }
}
