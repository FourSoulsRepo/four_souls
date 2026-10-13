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
  // locked: joining needs the game password (N-14).
  locked?: boolean;
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
  // password, if set, is needed to join (N-14).
  password?: string;
}

export interface Join {
  game: string;
  password?: string;
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

// recordURL is where the server serves a finished game's record to one
// of its players (RP-05).
export function recordURL(address: string, game: string, token: string): string | null {
  const ws = wsURL(address);
  if (!ws) {
    return null;
  }
  const q = new URLSearchParams({ game, token });
  return `${ws.replace(/^ws:/, 'http:').replace(/\/ws$/, '/record')}?${q.toString()}`;
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
    case 'refused':
      return e.rule ? `${e.message} (${e.rule})` : e.message;
    default:
      return e.message;
  }
}

// Engine types (pkg/rules_engine), as updates carry them. Enums are
// numbers in Go's order; see the constants below.

// Deck order in deck_sizes and discards.
export const Deck = { treasure: 0, loot: 1, monster: 2, room: 3 } as const;

// Card roles (CardView.role).
export const Role = { none: 0, character: 1, item: 2, monster: 3, event: 4, soul: 5, curse: 6 } as const;

export interface CardView {
  id: number;
  card: string;
  role?: number;
  charged: boolean;
  damage?: number;
  counters?: { name?: string; count: number }[];
  stats?: MonsterStats; // only on monsters
  // abilities: each ability's words by the activate intent's choice;
  // only on a played card with two or more to pick from (C-05).
  abilities?: string[];
}

// counted is how many counters a card carries, of any kind.
export function counted(c: CardView): number {
  return (c.counters ?? []).reduce((n, x) => n + x.count, 0);
}

// MonsterStats are a monster's stats now and as printed (C-03).
export interface MonsterStats {
  hp: number;
  dc: number;
  atk: number;
  printed_hp: number;
  printed_dc: number;
  printed_atk: number;
}

export interface SlotView {
  top?: CardView;
  covered?: CardView[];
}

export interface PlayerView {
  id: number;
  character: CardView;
  hp: number;
  atk: number;
  cents: number;
  souls: number;
  dead?: boolean;
  in_play: CardView[] | null;
  hand_size: number;
  hand?: CardView[]; // only the viewer's own hand
  treasure_top?: CardView;
  // What the player may still do now; attacks and purchases on their turn.
  loot_plays: number;
  attacks: number;
  purchases: number;
}

export interface Turn {
  active: number;
  number: number;
  step: number;
  loot_plays: number;
  attacks: number;
  purchases: number;
  attack_rolls?: number;
  bonus_loot_used?: number;
  bonus_attacks_used?: number;
  bonus_purchases_used?: number;
  must_attack?: number;
  must_attacks?: number;
  must_attack_deck?: number;
  deck_attacks?: number;
  end_declared?: boolean;
  death_end?: boolean;
  entered?: boolean;
}

export interface Target {
  object?: number;
  player: number;
  is_player?: boolean;
}

export interface AbilityRef {
  card: string;
  index: number;
}

export interface Chosen {
  kind: number;
  player: number;
  object?: number;
  stack_id?: number;
  spent?: boolean;
}

export interface StackItem {
  id: number;
  kind: number;
  controller: number;
  source: number;
  card?: string;
  roll?: number;
  amount?: number;
  label?: string;
  target: Target;
  attack?: boolean;
  ability: AbilityRef;
  mode?: number;
  targets?: Chosen[];
  roll_for: AbilityRef;
  roll_result?: number;
  event_player?: number;
  event_amount?: number;
  event_stack?: number;
  event_object?: number;
  reward?: number;
  checked?: number;
  // text is the used ability's words, or its chosen mode's (15.21).
  text?: string;
  // source_card and target_card name the cards of source and target
  // while anyone may see them (15.22).
  source_card?: string;
  target_card?: string;
}

export interface Prompt {
  kind: number;
  player: number;
  count?: number;
  options?: string[];
  purpose?: number;
  text?: string;
}

export interface GameView {
  players: PlayerView[];
  shop: SlotView[] | null;
  monsters: SlotView[] | null;
  deck_sizes: number[];
  discards: (CardView[] | null)[];
  bonus_souls?: CardView[];
  stack: StackItem[] | null;
  turn: Turn;
  waiting: Prompt;
  over?: boolean;
  winners?: number[];
  // options says what each choice option points at; only for the chooser.
  options?: OptionView[];
  // Who passed in a row in the open priority window, and who may still
  // respond, the waiting player first (V-11).
  passed?: number[];
  responding?: number[];
}

// OptionView is what one option of a choose prompt points at.
export interface OptionView {
  object?: number;
  player: number;
  is_player?: boolean;
  deck?: boolean;
  stack_id?: number;
}

export interface Intent {
  player: number;
  kind: number;
  objects?: number[];
  choice?: number;
}

export interface GameEvent {
  kind: string;
  player: number;
  object?: number;
  card?: string;
  amount?: number;
  text?: string;
  private?: boolean;
  // seen_by: who else sees a private card (the giver of a hand card).
  seen_by?: number[];
  prev?: number;
  source?: number;
  stack_id?: number;
  // item: the stack item that resolved, on stack_resolved (15.22).
  item?: StackItem;
}

export interface Vote {
  seat: number;
  kick: boolean;
}

export interface Pause {
  away: number[];
  votes: Vote[] | null;
}

// Update is the state for this client after a step (N-04).
export interface Update {
  step: number;
  events: GameEvent[] | null;
  view: GameView;
  allowed: Intent[] | null;
  seats: Seat[] | null;
  deadline?: number; // Unix ms; the response timer
  skip_all?: boolean;
  wait_turn?: boolean;
  pause?: Pause;
}
