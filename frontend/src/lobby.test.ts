import { describe, expect, it } from 'vitest';
import { emptyLobby, leftTable, lost, onMessage, optionsText, remember, type Session } from './lobby';
import type { Table } from './protocol';

const start: Session = { screen: 'join', lobby: emptyLobby, error: '' };

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
});

describe('remember', () => {
  it('puts the address first, once, and keeps five', () => {
    expect(remember(['a', 'b'], ' b ')).toEqual(['b', 'a']);
    expect(remember(['a', 'b', 'c', 'd', 'e'], 'f')).toEqual(['f', 'a', 'b', 'c', 'd']);
  });
});

describe('optionsText', () => {
  it('sums up the simple and the detailed setup', () => {
    expect(optionsText(undefined)).toBe('Random characters');
    expect(optionsText({ picking: 'draft', draft_size: 4, ban_rounds: 2, response_timer: 90, no_bonus_souls: true })).toBe(
      'Draft of 4 · 2 ban rounds · 90 s to respond · no bonus souls',
    );
  });
});
