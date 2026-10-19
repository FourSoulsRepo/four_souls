// Checks the TS protocol types against Go's golden file
// (internal/protocol/testdata/golden.jsonl). When Go changes a message,
// the golden file changes and this test fails until the types follow.
import { describe, expect, it } from 'vitest';
import golden from '../../internal/protocol/testdata/golden.jsonl?raw';
import {
  errorText,
  recordURL,
  wsURL,
  type Create,
  type Envelope,
  type Games,
  type Hello,
  type Intent,
  type Join,
  type Load,
  type ProtocolError,
  type Ready,
  type Saved,
  type Table,
  type Update,
  type Welcome,
} from './protocol';

// One expected envelope per golden line the client models. `satisfies`
// makes the compiler check each against its type.
const modeled: Record<string, Envelope> = {
  hello: {
    type: 'hello',
    id: 1,
    data: { protocol: 3, app: '0.1.0', name: 'Ann', role: 'player', sets: ['b2'], dice: 'tear' } satisfies Hello,
  },
  welcome: {
    type: 'welcome',
    id: 1,
    data: { protocol: 3, app: '0.1.0', engine: '0.1.0', token: 't0k3n', seat: 0, presence: true } satisfies Welcome,
  },
  error: {
    type: 'error',
    id: 2,
    data: { code: 'refused', message: 'no loot play available', rule: 'R-CARD-08' } satisfies ProtocolError,
  },
  intent: { type: 'intent', id: 2, data: { intent: { player: 0, kind: 0 } satisfies Intent } },
  update: {
    type: 'update',
    data: {
      step: 1,
      events: [{ kind: 'looted', player: 0, card: 'a_dime' }],
      view: {
        players: [
          {
            id: 0,
            character: { id: 310, card: 'isaac', role: 1, charged: false },
            hp: 2,
            atk: 1,
            cents: 3,
            souls: 0,
            in_play: [{ id: 311, card: 'the_d6', role: 2, charged: true }],
            hand_size: 3,
            loot_plays: 0,
            attacks: 0,
            purchases: 0,
          },
          {
            id: 1,
            character: { id: 312, card: 'cain', role: 1, charged: true },
            hp: 2,
            atk: 1,
            cents: 3,
            souls: 0,
            in_play: [{ id: 313, card: 'sleight_of_hand', role: 2, charged: true }],
            hand_size: 3,
            hand: [
              { id: 317, card: 'three_cents', charged: false },
              { id: 318, card: 'four_cents', charged: false },
              { id: 319, card: 'three_cents', charged: false },
            ],
            loot_plays: 0,
            attacks: 0,
            purchases: 0,
          },
        ],
        shop: [
          { top: { id: 305, card: 'synthoil', role: 2, charged: true } },
          { top: { id: 306, card: 'breakfast', role: 2, charged: true } },
        ],
        monsters: [
          {
            top: {
              id: 308,
              card: 'pin',
              role: 3,
              charged: true,
              stats: { hp: 2, dc: 4, atk: 1, printed_hp: 2, printed_dc: 4, printed_atk: 1 },
            },
          },
          {
            top: {
              id: 309,
              card: 'pestilence',
              role: 3,
              charged: true,
              stats: { hp: 4, dc: 4, atk: 1, printed_hp: 4, printed_dc: 4, printed_atk: 1 },
            },
          },
        ],
        deck_sizes: [96, 98, 101, 0],
        discards: [null, null, null, null],
        stack: null,
        turn: { active: 1, number: 1, step: 1, loot_plays: 0, attacks: 0, purchases: 0, entered: true },
        waiting: { kind: 1, player: 1 },
        win_souls: 4,
        responding: [1, 0],
      },
      allowed: [
        { player: 1, kind: 0 },
        { player: 1, kind: 7, objects: [312] },
        { player: 1, kind: 7, objects: [313] },
      ],
      aims: [
        {
          intent: 0,
          targets: [
            { object: 320, player: -1 },
            { player: 1, is_player: true },
          ],
          reason: 'Pick a monster or a player.',
        },
      ],
      blocked: [
        { kind: 5, object: 0, reason: 'Not your turn' },
        { kind: 3, object: 314, reason: 'Loot is played on your turn' },
      ],
      seats: [
        { seat: 0, name: 'Ann', connected: true, dice: 'tear' },
        { seat: 1, name: 'Bo', connected: false },
      ],
    } satisfies Update,
  },
  list: { type: 'list', id: 4, data: {} },
  games: {
    type: 'games',
    id: 4,
    data: {
      games: [{ id: 'g1', host: 'Ann', seats: 4, taken: 2, sets: ['b2'], started: false }],
    } satisfies Games,
  },
  create: {
    type: 'create',
    id: 5,
    data: {
      seats: 3,
      sets: ['b2'],
      options: {
        picking: 'draft',
        draft_size: 3,
        ban_rounds: 1,
        ban_timer: 30,
        host_bans: ['eden'],
        response_timer: 90,
      },
    } satisfies Create,
  },
  join: { type: 'join', id: 6, data: { game: 'g1' } satisfies Join },
  ready: { type: 'ready', id: 7, data: { ready: true } satisfies Ready },
  leave: { type: 'leave', id: 8, data: {} },
  table: {
    type: 'table',
    data: {
      game: 'g1',
      you: 1,
      seats: [
        { seat: 0, name: 'Ann', connected: true, ready: true, dice: 'tear' },
        { seat: 1, name: 'Bo', connected: true },
      ],
      sets: ['b2'],
      options: { picking: 'random' },
    } satisfies Table,
  },
  saved: { type: 'saved', data: { save: 'g1' } satisfies Saved },
  load: { type: 'load', id: 14, data: { save: 'g1' } satisfies Load },
  wait_turn: { type: 'wait_turn', id: 15, data: { on: true } },
  // A second line of one type is keyed by type and ID.
  'intent#16': { type: 'intent', id: 16, data: { intent: { player: 0, kind: 9 } satisfies Intent } },
};

// model is the expected envelope of a golden line, if the client models it.
function model(env: Envelope): Envelope | undefined {
  return modeled[`${env.type}#${String(env.id)}`] ?? modeled[env.type];
}

// Messages the client does not read or send yet; later steps model them.
const notYet = ['resync', 'setup', 'ban', 'pick', 'skip_all', 'vote', 'save', 'presence'];

const lines = golden
  .split('\n')
  .filter((l) => l.trim() !== '')
  .map((l) => JSON.parse(l) as Envelope);

describe('golden protocol file', () => {
  it('has every message type either modeled or listed as not yet', () => {
    for (const env of lines) {
      expect(model(env) !== undefined || notYet.includes(env.type), env.type).toBe(true);
    }
  });

  it.each(lines.filter((env) => model(env) !== undefined))('$type $id matches the TS types', (env) => {
    expect(env).toEqual(model(env));
  });
});

describe('wsURL', () => {
  it('adds the default port and the path', () => {
    expect(wsURL(' 192.168.1.5 ')).toBe('ws://192.168.1.5:4774/ws');
    expect(wsURL('host.lan:5000')).toBe('ws://host.lan:5000/ws');
    expect(wsURL('[::1]:4774')).toBe('ws://[::1]:4774/ws');
  });

  it('rejects what is not an address', () => {
    for (const bad of ['', 'a b', 'http://x', 'x:0', 'x:70000', 'x:y']) {
      expect(wsURL(bad), bad).toBeNull();
    }
  });
});

describe('errorText', () => {
  it('explains an outdated server', () => {
    expect(errorText({ code: 'server_outdated', message: '' })).toBe(
      'The server is older than this app. Ask the host to update.',
    );
  });
});

describe('recordURL', () => {
  it('points at the record download of the same server', () => {
    expect(recordURL('192.168.1.5', 'g1', 't0k3n')).toBe('http://192.168.1.5:4774/record?game=g1&token=t0k3n');
    expect(recordURL('not an address', 'g1', 't')).toBeNull();
  });
});
