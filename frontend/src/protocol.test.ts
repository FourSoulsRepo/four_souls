// Checks the TS protocol types against Go's golden file
// (internal/protocol/testdata/golden.jsonl). When Go changes a message,
// the golden file changes and this test fails until the types follow.
import { describe, expect, it } from 'vitest';
import golden from '../../internal/protocol/testdata/golden.jsonl?raw';
import {
  errorText,
  wsURL,
  type Create,
  type Envelope,
  type Games,
  type Hello,
  type Join,
  type Load,
  type ProtocolError,
  type Ready,
  type Saved,
  type Table,
  type Welcome,
} from './protocol';

// One expected envelope per golden line the client models. `satisfies`
// makes the compiler check each against its type.
const modeled: Record<string, Envelope> = {
  hello: {
    type: 'hello',
    id: 1,
    data: { protocol: 1, app: '0.1.0', name: 'Ann', role: 'player', sets: ['b2'] } satisfies Hello,
  },
  welcome: {
    type: 'welcome',
    id: 1,
    data: { protocol: 1, app: '0.1.0', engine: '0.1.0', token: 't0k3n', seat: 0 } satisfies Welcome,
  },
  error: {
    type: 'error',
    id: 2,
    data: { code: 'refused', message: 'no loot play available', rule: 'R-CARD-08' } satisfies ProtocolError,
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
        { seat: 0, name: 'Ann', connected: true, ready: true },
        { seat: 1, name: 'Bo', connected: true },
      ],
      sets: ['b2'],
      options: { picking: 'random' },
    } satisfies Table,
  },
  saved: { type: 'saved', data: { save: 'g1' } satisfies Saved },
  load: { type: 'load', id: 14, data: { save: 'g1' } satisfies Load },
};

// Messages the client does not read or send yet; later steps model them.
const notYet = ['intent', 'resync', 'update', 'setup', 'ban', 'pick', 'skip_all', 'vote', 'save'];

const lines = golden
  .split('\n')
  .filter((l) => l.trim() !== '')
  .map((l) => JSON.parse(l) as Envelope);

describe('golden protocol file', () => {
  it('has every message type either modeled or listed as not yet', () => {
    for (const env of lines) {
      expect(Object.keys(modeled).includes(env.type) || notYet.includes(env.type), env.type).toBe(true);
    }
  });

  it.each(lines.filter((env) => env.type in modeled))('$type matches the TS types', (env) => {
    expect(env).toEqual(modeled[env.type]);
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
