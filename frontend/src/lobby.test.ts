import { describe, expect, it } from 'vitest';
import { autoJoinStep, botLead, botLevelName, botLevels, botNames, canCreate, namePieces, seatNameOf, challengeText, isAway, isTableHost, seatName, emptyLobby, lanRows, oneOpenGame, leftTable, forgotten, lost, onMessage, optionsText, remember, renamed, serverLabel, upText, type Session } from './lobby';
import type { Seat, Table, Welcome } from './protocol';

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

  it('keeps the draft offers until the game starts (14.11)', () => {
    const setup = { phase: 'pick' as const, turn: -1, pool: [], banned: [], offers: ['isaac', 'eden'], items: ['the_d6', ''], picked: [false, false], seats: [] };
    const atTable = onMessage(start, { type: 'table', data: table });
    const picking = onMessage(atTable, { type: 'setup', data: setup });
    expect(picking).toMatchObject({ screen: 'game', lobby: { table, setup } });
    const playing = onMessage(picking, { type: 'update', data: {} });
    expect(playing.lobby.setup).toBeNull();
    expect(leftTable(picking).lobby.setup).toBeNull();
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
  it('tells whether a server is up and its open games (15.100)', () => {
    expect(upText(undefined)).toBe('Checking…');
    expect(upText('down')).toBe('Down');
    expect(upText({ open: 0 })).toBe('0 open games');
    expect(upText({ open: 1 })).toBe('1 open game');
    expect(upText({ open: 3 })).toBe('3 open games');
  });
});

describe('optionsText', () => {
  it('sums up the simple and the detailed setup', () => {
    expect(optionsText(undefined)).toBe('Random characters');
    expect(optionsText({ picking: 'draft', draft_size: 4, ban_rounds: 2, response_timer: 90, no_bonus_souls: true })).toBe(
      'Draft of 4 · 2 ban rounds · 90 s to respond · no bonus souls · deck ratios',
    );
    expect(optionsText({ no_deck_ratios: true })).toBe('Random characters');
    // The rollback mode (ADR 020).
    expect(optionsText({ no_deck_ratios: true, rollback: 'god' })).toBe('Random characters · Godmode rollback');
    expect(optionsText({ no_deck_ratios: true, random_seats: true })).toBe('Random characters · random seats');
    // The room deck (14.13).
    expect(optionsText({ no_deck_ratios: true, room_deck: true })).toBe('Random characters · room deck');
    // A challenge comes first; its bonus souls are always off (ADR 028).
    expect(optionsText({ no_deck_ratios: true, no_bonus_souls: true, challenge: 'krampus', difficulty: 'ultra' })).toBe(
      'Challenge: How the Krampus Stole Christmas, Ultra Hard · Random characters',
    );
  });
});

describe('challengeText', () => {
  it('names the challenge and its difficulty; an unknown one by its ID', () => {
    expect(challengeText('krampus', 'normal')).toBe('Challenge: How the Krampus Stole Christmas, Normal');
    expect(challengeText('krampus', 'hard')).toBe('Challenge: How the Krampus Stole Christmas, Hard');
    expect(challengeText('newer', 'hard')).toBe('Challenge: newer, Hard');
  });

  it('names Trick or Treat (14.33)', () => {
    expect(challengeText('trick_or_treat', 'ultra')).toBe('Challenge: Trick or Treat, Ultra Hard');
  });
});

describe('lanRows', () => {
  const none = { address: '', token: '', game: '', via: '' as const };
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
      {
        key: '192.168.1.5:4774#0',
        address: '192.168.1.5:4774',
        text: "🔒 Ann's game · 1/3 players · 192.168.1.5:4774",
        joinable: true,
        rejoin: false,
        game: '',
        host: 'Ann',
        locked: true,
      },
      {
        key: '192.168.1.9:4774#0',
        address: '192.168.1.9:4774',
        text: "Bo's game · 1/2 players · 192.168.1.9:4774 · needs app 2.0.0",
        joinable: false,
        rejoin: false,
        game: '',
        host: 'Bo',
        locked: false,
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
    const mine = lanRows(found, 1, { address: '192.168.1.5', token: 't', game: 'g1', via: 'lan' as const });
    expect(mine[0]).toMatchObject({ joinable: true, rejoin: true });
    expect(mine[1]).toMatchObject({ joinable: false, rejoin: false });
    // A seat without a token holds nothing.
    expect(lanRows(found, 1, { ...none, game: 'g1' })[0]).toMatchObject({ rejoin: false });
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

describe('LAN Join (19.149)', () => {
  const g = { id: 'g1', host: 'Ann', seats: 2, taken: 1, sets: ['b2'], started: false };
  const games = (list: unknown[]) => ({ type: 'games', data: { games: list } });

  it('a typed address takes only one waiting, open, unlocked game', () => {
    expect(oneOpenGame([g])).toBe('g1');
    expect(oneOpenGame([])).toBe('');
    expect(oneOpenGame([g, { ...g, id: 'g2' }])).toBe('');
    expect(oneOpenGame([{ ...g, started: true }])).toBe('');
    expect(oneOpenGame([{ ...g, taken: 2 }])).toBe('');
    expect(oneOpenGame([{ ...g, locked: true }])).toBe('');
  });

  it("a row's Join sends its game and password once the games come", () => {
    const a = { game: 'g7', password: 'pw', sent: false };
    const step = autoJoinStep(a, games([g]));
    expect(step).toEqual({ next: { ...a, sent: true }, join: { game: 'g7', password: 'pw' }, fail: false });
    // A later list changes nothing: the join went out.
    expect(autoJoinStep({ ...a, sent: true }, games([]))).toEqual({ next: { ...a, sent: true }, fail: false });
  });

  it("a typed address joins its one open game; otherwise the Games list shows", () => {
    const typed = { game: '', password: '', sent: false };
    expect(autoJoinStep(typed, games([g]))).toEqual({ next: { game: 'g1', password: '', sent: true }, join: { game: 'g1' }, fail: false });
    expect(autoJoinStep(typed, games([g, { ...g, id: 'g2' }]))).toEqual({ next: null, fail: false });
    expect(autoJoinStep(typed, { type: 'games', data: { games: null } })).toEqual({ next: null, fail: false });
  });

  it('an error fails it; a seat or a game ends it; the welcome waits', () => {
    const a = { game: 'g1', password: '', sent: true };
    expect(autoJoinStep(a, { type: 'error', data: { code: 'full', message: 'Full' } })).toEqual({ next: null, fail: true });
    expect(autoJoinStep(a, { type: 'table', data: table })).toEqual({ next: null, fail: false });
    expect(autoJoinStep(a, { type: 'update', data: {} })).toEqual({ next: null, fail: false });
    expect(autoJoinStep(a, { type: 'welcome', data: {} })).toEqual({ next: a, fail: false });
  });
});

describe('bots at the table (ADR 030)', () => {
  const bot: Seat = { seat: 1, name: 'Bot 1', connected: false, ready: true, bot: 'standard', bot_n: 1 };
  const free: Seat = { seat: 2, name: '', connected: false };
  const ann: Seat = { seat: 0, name: 'Ann', connected: true };
  const bo: Seat = { seat: 3, name: 'Bo', connected: true };
  const withBots: Table = { game: 'g1', you: 0, seats: [ann, bot, free, bo], sets: ['b2'] };
  const welcome: Welcome = { protocol: 7, app: '0.1.0', engine: '0.1.0', token: '', seat: -1, bots: ['standard', 'hard'] };

  it('names a bot by its number and level, a player by the name', () => {
    expect(seatName(bot)).toBe('Bot 1 (Standard)');
    expect(seatName({ ...bot, bot: 'hard', bot_n: 3 })).toBe('Bot 3 (Hard)');
    expect(seatName(ann)).toBe('Ann');
  });

  it("shows a newer server's level by its ID", () => {
    expect(botLevelName('hard')).toBe('Hard');
    expect(botLevelName('ultra')).toBe('ultra');
  });

  it('never shows a bot as away, though it has no connection', () => {
    expect(isAway(bot)).toBe(false);
    expect(isAway({ ...ann, connected: false })).toBe(true);
    expect(isAway(free)).toBe(false);
  });

  it('gives the host the levels the server allows', () => {
    expect(botLevels(withBots, welcome)).toEqual(['standard', 'hard']);
    expect(botLevels(withBots, { ...welcome, bots: ['standard'] })).toEqual(['standard']);
  });

  it('gives no buttons when the server allows no bots', () => {
    expect(botLevels(withBots, { ...welcome, bots: undefined })).toEqual([]);
    expect(botLevels(withBots, { ...welcome, bots: [] })).toEqual([]);
    expect(botLevels(withBots, null)).toEqual([]);
  });

  it('gives no buttons to the other players', () => {
    expect(isTableHost({ ...withBots, you: 3 })).toBe(false);
    expect(botLevels({ ...withBots, you: 3 }, welcome)).toEqual([]);
  });

  it("gives no buttons at a saved game's table: its seats wait for its players", () => {
    expect(botLevels({ ...withBots, loaded: 'g1' }, welcome)).toEqual([]);
  });

  it('makes the first player in seat order the host once seat 0 is free', () => {
    // The host left a table without bots; the server moves Bo to seat 0
    // on the first bot (Decision 10.12). A bot never hosts.
    const left: Table = { ...withBots, you: 3, seats: [{ ...free, seat: 0 }, bot, free, bo] };
    expect(isTableHost(left)).toBe(true);
    expect(botLevels(left, welcome)).toEqual(['standard', 'hard']);
  });
});

describe('bots at the game table (21.12)', () => {
  const ann: Seat = { seat: 0, name: 'Ann', connected: true };
  const bot1: Seat = { seat: 1, name: 'Bot 1', connected: false, bot: 'standard', bot_n: 1 };
  const bot2: Seat = { seat: 2, name: 'Bot 2', connected: false, bot: 'hard', bot_n: 2 };
  const seats = [ann, bot1, bot2];

  it('names a seat as the table does, a bot translated with its level', () => {
    expect(seatNameOf(seats, 0)).toBe('Ann');
    expect(seatNameOf(seats, 2)).toBe('Bot 2 (Hard)');
    expect(seatNameOf(seats, 3)).toBe('Player 4');
    expect(seatNameOf(null, 0)).toBe('Player 1');
  });

  it('lists the names the bots show, never a player', () => {
    expect(botNames(seats)).toEqual(['Bot 1 (Standard)', 'Bot 2 (Hard)']);
    expect(botNames([ann])).toEqual([]);
    expect(botNames(undefined)).toEqual([]);
  });

  it('marks every bot name in a line, the words around it plain', () => {
    const bots = botNames(seats);
    expect(namePieces('Bot 1 (Standard) attacks Ann', bots)).toEqual([{ text: 'Bot 1 (Standard)', bot: true }, { text: ' attacks Ann' }]);
    expect(namePieces('Ann, Bot 2 (Hard) and Bot 1 (Standard) win!', bots)).toEqual([
      { text: 'Ann, ' },
      { text: 'Bot 2 (Hard)', bot: true },
      { text: ' and ' },
      { text: 'Bot 1 (Standard)', bot: true },
      { text: ' win!' },
    ]);
    expect(namePieces("Bot 1 (Standard)'s turn", bots)).toEqual([{ text: 'Bot 1 (Standard)', bot: true }, { text: "'s turn" }]);
  });

  it('leaves a line with no bot name whole', () => {
    expect(namePieces('Ann attacks Bot 3', botNames(seats))).toEqual([{ text: 'Ann attacks Bot 3' }]);
    expect(namePieces('Ann attacks', [])).toEqual([{ text: 'Ann attacks' }]);
    expect(namePieces('', botNames(seats))).toEqual([]);
  });

  it("keeps the robot with a bot name's first word; the rest may wrap", () => {
    expect(botLead('Bot 1 (Standard)')).toEqual(['Bot', ' 1 (Standard)']);
    expect(botLead('Robot')).toEqual(['Robot', '']);
  });

  it('tries the longest name first, so one inside another never splits it', () => {
    expect(namePieces('Bot 1 (Standard) Max', ['Bot 1 (Standard)', 'Bot 1 (Standard) Max'])).toEqual([{ text: 'Bot 1 (Standard) Max', bot: true }]);
  });
});
