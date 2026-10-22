import { describe, expect, it } from 'vitest';
import { canCreate, emptyLobby, lanRows, leftTable, forgotten, lost, onMessage, optionsText, remember, renamed, serverLabel, type Session } from './lobby';
import type { Table } from './protocol';

const start: Session = { screen: 'join', lobby: emptyLobby, game: null, error: '' };

const table: Table = {
  game: 'g1',
  you: 0,
  seats: [{ seat: 0, name: 'Ann', connected: true }],
  sets: ['b2'],
};

describe('onMessage', () => {
  it('keeps the game list, reading null lists from Go as empty', () => {
    const s = onMessage(start, { type: 'games', data: { games: null } });
    expect(s.lobby).toMatchObject({ games: [], saves: [] });
    const game = { id: 'g1', host: 'Ann', seats: 2, taken: 1, sets: ['b2'], started: false };
    expect(onMessage(s, { type: 'games', data: { games: [game] } }).lobby.games).toEqual([game]);
  });

  it('goes to the table when seated, and clears an old error', () => {
    const s = onMessage({ ...start, error: 'old' }, { type: 'table', data: table });
    expect(s).toMatchObject({ screen: 'table', error: '', lobby: { table } });
  });

  it('goes to the game when the setup or the first update comes', () => {
    const atTable = onMessage(start, { type: 'table', data: table });
    expect(onMessage(atTable, { type: 'setup', data: {} }).screen).toBe('game');
    expect(onMessage(atTable, { type: 'update', data: {} }).screen).toBe('game');
  });

  it('leaves showing updates to playback, and drops the game at a new table', () => {
    expect(onMessage(start, { type: 'update', data: { step: 1 } }).game).toBeNull();
    const playing = { ...start, game: { step: 1 } as Session['game'] };
    expect(onMessage(playing, { type: 'table', data: table }).game).toBeNull();
  });

  it('goes back to the game list after a save', () => {
    const atTable = onMessage(start, { type: 'table', data: table });
    const s = onMessage(atTable, { type: 'saved', data: { save: 'g1' } });
    expect(s).toMatchObject({ screen: 'join', lobby: { table: null, saved: 'g1' } });
  });

  it('shows server errors in plain words', () => {
    const s = onMessage(start, { type: 'error', data: { code: 'client_outdated', message: 'protocol 0 < 1' } });
    expect(s.error).toBe('This app is older than the server. Update the app to play here.');
    expect(onMessage(start, { type: 'error', data: { code: 'room_full', message: 'the game is full' } }).error).toBe(
      'the game is full',
    );
  });

  it('ignores messages it does not know', () => {
    expect(onMessage(start, { type: 'something_new' })).toBe(start);
  });
});

describe('leaving', () => {
  it('leftTable goes back to the game list', () => {
    const s = leftTable(onMessage(start, { type: 'table', data: table }));
    expect(s).toMatchObject({ screen: 'join', lobby: { table: null } });
  });

  it('lost keeps a server error sent just before the close', () => {
    expect(lost(start).error).toBe('Lost the connection to the server.');
    expect(lost({ ...start, error: 'This app is older than the server.' }).error).toBe(
      'This app is older than the server.',
    );
    expect(lost({ ...start, screen: 'table' }).screen).toBe('join');
  });

  it('lost keeps a running game, to reconnect to', () => {
    const playing = onMessage(start, { type: 'table', data: table });
    const s = lost({ ...playing, screen: 'game' });
    expect(s).toMatchObject({ screen: 'game', lobby: { table } });
  });
});

// at makes recent entries without names.
const at = (...addresses: string[]) => addresses.map((address) => ({ address, name: '' }));

describe('remember', () => {
  it('puts the address first, once, and keeps five', () => {
    expect(remember(at('a', 'b'), ' b ')).toEqual(at('b', 'a'));
    expect(remember(at('a', 'b', 'c', 'd', 'e'), 'f')).toEqual(at('f', 'a', 'b', 'c', 'd'));
  });
  it('keeps the name the server gave', () => {
    expect(remember([{ address: 'a', name: '' }, { address: 'b', name: 'Basement' }], 'b')).toEqual([
      { address: 'b', name: 'Basement' },
      { address: 'a', name: '' },
    ]);
  });
});

describe('server names (N-18)', () => {
  it('names the entry of the address', () => {
    expect(renamed(at('a', 'b'), ' b ', 'Basement')).toEqual([
      { address: 'a', name: '' },
      { address: 'b', name: 'Basement' },
    ]);
  });
  it('keeps the same list when nothing changes', () => {
    const recent = [{ address: 'a', name: 'Basement' }];
    expect(renamed(recent, 'a', 'Basement')).toBe(recent);
    expect(renamed(recent, 'b', 'Cathedral')).toBe(recent);
  });
  it('drops a name the server no longer gives', () => {
    expect(renamed([{ address: 'a', name: 'Basement' }], 'a', '')).toEqual(at('a'));
  });
  it('forgets an address', () => {
    expect(forgotten(at('a', 'b', 'c'), 'b')).toEqual(at('a', 'c'));
  });
  it('shows the name, else the address', () => {
    expect(serverLabel({ address: 'a', name: 'Basement' })).toBe('Basement');
    expect(serverLabel({ address: 'a', name: '' })).toBe('a');
  });
});

describe('optionsText', () => {
  it('sums up the simple and the detailed setup', () => {
    expect(optionsText(undefined)).toBe('Random characters');
    expect(optionsText({ picking: 'draft', draft_size: 4, ban_rounds: 2, response_timer: 90, no_bonus_souls: true })).toBe(
      'Draft of 4 · 2 ban rounds · 90 s to respond · no bonus souls · deck ratios',
    );
    expect(optionsText({ no_deck_ratios: true })).toBe('Random characters');
  });
});

describe('lanRows', () => {
  const none = { address: '', token: '', game: '' };
  const game = { id: '', host: 'Ann', seats: 3, taken: 1, locked: false, started: false };

  it('lists each waiting game with its address and lock; another version cannot be joined', () => {
    const rows = lanRows(
      [
        { address: '192.168.1.5:4774', protocol: 1, app: '1.0.0', games: [{ ...game, locked: true }] },
        { address: '192.168.1.9:4774', protocol: 2, app: '2.0.0', games: [{ ...game, host: 'Bo', seats: 2 }] },
      ],
      1,
      none,
    );
    expect(rows).toEqual([
      { key: '192.168.1.5:4774#0', address: '192.168.1.5:4774', text: "🔒 Ann's game · 1/3 players · 192.168.1.5:4774", joinable: true, rejoin: false },
      {
        key: '192.168.1.9:4774#0',
        address: '192.168.1.9:4774',
        text: "Bo's game · 1/2 players · 192.168.1.9:4774 · needs app 2.0.0",
        joinable: false,
        rejoin: false,
      },
    ]);
  });

  it('marks a started game in progress; only who held a seat rejoins it (15.97)', () => {
    const found = [
      {
        address: '192.168.1.5:4774',
        protocol: 1,
        app: '1.0.0',
        games: [
          { ...game, id: 'g1', taken: 3, started: true },
          { ...game, id: 'g2', taken: 3 },
        ],
      },
    ];
    const [started, full] = lanRows(found, 1, none);
    expect(started).toMatchObject({ text: "Ann's game · 3/3 players · In progress · 192.168.1.5:4774", joinable: false, rejoin: false });
    expect(full).toMatchObject({ joinable: false, rejoin: false });
    const mine = lanRows(found, 1, { address: '192.168.1.5', token: 't', game: 'g1' });
    expect(mine[0]).toMatchObject({ joinable: true, rejoin: true });
    expect(mine[1]).toMatchObject({ joinable: false, rejoin: false });
    // A seat without a token holds nothing.
    expect(lanRows(found, 1, { address: '', token: '', game: 'g1' })[0]).toMatchObject({ rejoin: false });
  });
});

describe('canCreate', () => {
  const welcome = { protocol: 1, app: '', engine: '', token: 't', seat: -1 };

  it('offers "Create a game" when the server allows it (N-17)', () => {
    expect(canCreate({ ...emptyLobby, welcome: { ...welcome, create: true } }, false)).toBe(true);
  });

  it('hides it on a friend\'s hosted game, which never allows it', () => {
    expect(canCreate({ ...emptyLobby, welcome }, false)).toBe(false);
    expect(canCreate(emptyLobby, false)).toBe(false);
  });

  it('hides it from the host, who uses the host screen', () => {
    expect(canCreate({ ...emptyLobby, welcome: { ...welcome, create: true } }, true)).toBe(false);
  });
});
