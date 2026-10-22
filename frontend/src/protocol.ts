// Network messages (ADR 006), mirroring internal/protocol/protocol.go.
// Only what the client uses so far; protocol.test.ts checks these types
// against Go's golden file. Go sends a nil list as null.

export const protocolVersion = 4;
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
  dice?: string; // the player's dice look (V-19)
}

export interface Welcome {
  protocol: number;
  app: string;
  engine: string;
  token: string;
  seat: number;
  // create is set when this app may open a game on the server (N-17).
  create?: boolean;
  // presence is set when the server relays live aiming and hover; the
  // app sends none otherwise (15.35, ADR 016).
  presence?: boolean;
  // name is the server's name, shown instead of its address (N-18).
  name?: string;
}

// Presence is what a player points at now (15.35, ADR 016): the card
// they aim from, and the spot the arrow or pointer is over. Each one
// replaces the sender's last; an empty one means nothing. The server
// fills in seat.
export interface Presence {
  seat: number;
  from?: Spot;
  at?: Spot;
  // point is the pointer on the board frame, 0 to 1, while over no spot.
  point?: { x: number; y: number };
}

// Spot is one thing on the table; exactly one field is set. hand is
// the sender's own hand slot, from 1, never the card (A-16).
export interface Spot {
  object?: number;
  hand?: number;
  player?: number;
  deck?: number;
  stack_id?: number;
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
  no_deck_ratios?: boolean;
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
  // dice is the seat's dice look (V-19); none draws the default one.
  dice?: string;
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

// Notice tells the player something the game screen does not show
// (N-20): engine_updated, game_lost, game_idle.
export interface Notice {
  kind: string;
  game?: string;
  // reason of game_idle: 'empty' (nobody connected) or 'long' (24 hours).
  reason?: string;
}

// lanName and ipv4 are hosts never behind an HTTPS proxy: a bare one
// goes straight to plain ws:// (N-19).
const lanName = /\.(local|lan|home|internal|home\.arpa)$/i;
const ipv4 = /^[\d.]+$/;

// wsURLs turns what a player types into the server's WebSocket URLs to
// try in turn, or null (8.18, N-19):
//   - "192.168.1.5", "host:4774", "[::1]:4774", "localhost", "pc.lan":
//     ws:// on that port, 4774 by default;
//   - "wss://games.example.com", "ws://host:5000", "https://…",
//     "http://…": that scheme only; wss:// on 443 by default;
//   - "games.example.com": wss:// on 443 first, for a server behind an
//     HTTPS proxy (ADR 017), then ws:// on 4774.
export function wsURLs(address: string): string[] | null {
  let a = address.trim();
  let secure: boolean | null = null; // null: try both
  const scheme = /^(wss?|https?):\/\//i.exec(a);
  if (scheme) {
    secure = /^(wss|https):/i.test(a);
    a = a.slice(scheme[0].length);
    const tail = ['/ws/', '/ws', '/'].find((t) => a.endsWith(t));
    if (tail) {
      a = a.slice(0, -tail.length);
    }
  }
  const close = a.startsWith('[') ? a.indexOf(']') + 1 : 0; // IPv6 in brackets
  const colon = a.indexOf(':', close);
  const host = colon < 0 ? a : a.slice(0, colon);
  const portText = colon < 0 ? '' : a.slice(colon + 1);
  if (!/^(\[[0-9a-fA-F:.]+\]|[\w.-]+)$/.test(host) || (colon >= 0 && !/^\d{1,5}$/.test(portText))) {
    return null;
  }
  const port = portText ? Number(portText) : null;
  if (port !== null && (port < 1 || port > 65535)) {
    return null;
  }
  const plain = `ws://${host}:${String(port ?? defaultPort)}/ws`;
  const tls = port === null ? `wss://${host}/ws` : `wss://${host}:${String(port)}/ws`;
  if (secure !== null) {
    return [secure ? tls : plain];
  }
  const domain = port === null && host.includes('.') && !host.startsWith('[') && !ipv4.test(host) && !lanName.test(host);
  return domain ? [tls, plain] : [plain];
}

// recordURL is where the server at a WebSocket URL serves a finished
// game's record to one of its players (RP-05): http:// for ws://,
// https:// for wss://.
export function recordURL(ws: string, game: string, token: string): string | null {
  if (!/^wss?:\/\/.+\/ws$/.test(ws)) {
    return null;
  }
  const q = new URLSearchParams({ game, token });
  return `${ws.replace(/^ws/, 'http').replace(/\/ws$/, '/record')}?${q.toString()}`;
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
  // souls: the soul value the card adds to its owner, on a soul in play
  // or a character; older servers send none (V-43).
  souls?: number;
  // dice: the die numbers the card keeps, e.g. a stored roll (15.67).
  dice?: number[];
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
  sources?: StatSource[]; // what changes HP, DC and ATK now (V-41)
}

// StatSource is one source of a stat change (V-41). card is "" or absent
// when the source is hidden; older servers send no sources.
export interface StatSource {
  stat: "hp" | "atk" | "dc" | "roll" | "attack_roll" | "attacked_dc";
  source?: number;
  card?: string;
  amount: number;
  lasts: "in_play" | "turn" | "next_roll";
  when?: string; // a condition, e.g. "first attack roll"
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
  hand?: CardView[]; // the viewer's own hand, or one revealed to all (Guppy's Eye)
  treasure_top?: CardView;
  // What the player may still do now; attacks and purchases on their turn.
  loot_plays: number;
  attacks: number;
  purchases: number;
  // The character's printed stats and what changes them (V-41); older
  // servers send none.
  printed_hp?: number;
  printed_atk?: number;
  sources?: StatSource[];
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
  // revealed is the card a deck attack drew, until its slot is picked (19.9).
  revealed?: CardView;
  // subject is the card the open question is about, e.g. the deck top
  // a player looks at; only for the chooser (19.11).
  subject?: CardView;
  // ordering are the cards of an open order prompt, in their order now;
  // only for the chooser (19.16).
  ordering?: CardView[];
  stack: StackItem[] | null;
  turn: Turn;
  waiting: Prompt;
  over?: boolean;
  winners?: number[];
  // win_souls is the soul value that wins; older servers send none (V-30).
  win_souls?: number;
  // options says what each choice option points at; only for the chooser.
  options?: OptionView[];
  // activating is the card whose mode or targets the chooser picks now
  // (15.34); only for the chooser.
  activating?: number;
  // Who passed in a row in the open priority window, and who may still
  // respond, the waiting player first (V-11).
  passed?: number[];
  responding?: number[];
  // attack is the active player's attack while it lasts; buying a
  // declared purchase (15.30). Older servers send neither.
  attack?: AttackView;
  buying?: boolean;
  // first_roll is the first-player roll before turn 1: who still rolls,
  // in seat order, and this round's dice so far (GS-06, 15.43).
  first_roll?: FirstRoll;
}

// FirstRoll is the first-player roll in progress: the lowest goes
// first; a tie among the lowest rolls again, with only those players.
export interface FirstRoll {
  // round counts from 1; each tie starts the next.
  round: number;
  players: number[];
  rolls?: number[];
}

// AttackView is an attack in progress: its target once chosen, and
// whether it was declared on the monster deck.
export interface AttackView {
  target?: number;
  deck?: boolean;
}

// OptionView is what one option of a choose prompt points at.
export interface OptionView {
  object?: number;
  player: number;
  is_player?: boolean;
  deck?: boolean;
  stack_id?: number;
  // outcome is what the option leads to, when the engine can tell
  // (replacement order, 19.17).
  outcome?: OutcomeView;
}

// OutcomeView is what an order of replacements gives (19.17): what then
// happens, the cards whose effect no longer applies, or that the player
// orders the rest next.
export interface OutcomeView {
  actions?: { kind: number; player: number; amount: number; card?: string }[];
  unused?: string[];
  again?: boolean;
}

// Aim is what the allowed intent at index intent asks to target first.
export interface Aim {
  intent: number;
  targets: OptionView[];
  reason?: string; // what it may target, for a wrong click (15.33)
}

// Block says why the intent of that kind, on object if a card, can't be
// sent now (15.29).
export interface Block {
  kind: number;
  object: number;
  reason: string;
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
  // dice: every die of a roll of several dice on dice_rolled; the first
  // is amount (19.39).
  dice?: number[];
  // by: the player who dealt damage, plus 1; 0 or none for a monster
  // or nobody (15.77). On a death: the killer (15.88).
  by?: number;
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
  // grace is when a kick vote opens, Unix ms, after a server restart
  // (N-20); absent: open now.
  grace?: number;
}

// Update is the state for this client after a step (N-04).
export interface Update {
  step: number;
  events: GameEvent[] | null;
  view: GameView;
  allowed: Intent[] | null;
  aims?: Aim[]; // the first targets of allowed plays and uses (15.25)
  blocked?: Block[]; // why other cards and turn actions can't be used (15.29)
  seats: Seat[] | null;
  deadline?: number; // Unix ms; the response timer
  wait_turn?: boolean;
  pause?: Pause;
}
