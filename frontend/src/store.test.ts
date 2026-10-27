import { beforeEach, describe, expect, it, vi } from 'vitest';
import { checkSeat, getBacks, getCards, getDice, getFonts, getIcons, getMats, getSettings, getSounds, saveSettings, serverStatus } from './bridge';
import { language, setCardLanguages, setLanguage } from './i18n';
import { send } from './net';
import { noHostSetup } from './hostSetup';
import type { Envelope, GameView, Update } from './protocol';
import { messageOf, seatOf, useApp } from './store';

vi.mock('./bridge', () => ({
  checkSeat: vi.fn(),
  getBacks: vi.fn(),
  getCards: vi.fn(),
  getDice: vi.fn(),
  getFonts: vi.fn(),
  getIcons: vi.fn(),
  getMats: vi.fn(),
  getSettings: vi.fn(),
  getSounds: vi.fn(),
  saveSettings: vi.fn(),
  serverStatus: vi.fn(),
}));
vi.mock('./net', () => ({ drop: vi.fn(), send: vi.fn() }));

const initial = useApp.getState();

beforeEach(() => {
  useApp.setState(initial, true);
  vi.resetAllMocks();
});

describe('screens', () => {
  it('starts on the splash', () => {
    expect(useApp.getState().screen).toBe('splash');
  });

  it('show switches the screen and clears the error', () => {
    useApp.setState({ error: 'old' });
    useApp.getState().show('settings');
    expect(useApp.getState()).toMatchObject({ screen: 'settings', error: '' });
  });
});

describe('settings', () => {
  it('loads what Go has', async () => {
    vi.mocked(getSettings).mockResolvedValue({ nickname: 'Ann', recent: [], animation: 'normal', fullscreen: false, mat: 'basement', dice: 'tear', no_end_turn_ask: false, rejoin: { address: '', token: '', game: '', via: '' as const }, host: noHostSetup, lan_address: '', language: '' });
    await useApp.getState().loadSettings();
    expect(useApp.getState().settings.nickname).toBe('Ann');
  });

  it('keeps the defaults when loading fails', async () => {
    vi.mocked(getSettings).mockRejectedValue('no bindings');
    await useApp.getState().loadSettings();
    expect(useApp.getState().settings.nickname).toBe('');
  });

  it('keeps what Go stored, not what was typed', async () => {
    vi.mocked(saveSettings).mockResolvedValue({ nickname: 'Ann', recent: [], animation: 'normal', fullscreen: false, mat: 'basement', dice: 'tear', no_end_turn_ask: false, rejoin: { address: '', token: '', game: '', via: '' as const }, host: noHostSetup, lan_address: '', language: '' });
    expect(await useApp.getState().saveSettings({ nickname: '  Ann  ', recent: [], animation: 'normal', fullscreen: false, mat: 'basement', dice: 'tear', no_end_turn_ask: false, rejoin: { address: '', token: '', game: '', via: '' as const }, host: noHostSetup, lan_address: '', language: '' })).toBe(true);
    expect(useApp.getState().settings.nickname).toBe('Ann');
  });

  it('names and pictures the cards in the saved language at once (16.5, 16.6)', async () => {
    const images = { en: '/cardimg/en/b2/isaac.webp', ru: '/cardimg/ru/b2/isaac.webp' };
    const isaac = { ref: 'isaac', set: 'b2', type: 'Character Card', name: 'Isaac', names: { en: 'Isaac', ru: 'Исаак' }, effects: [], footnotes: [], stats: [], rewards: [], image: images.en, images };
    useApp.setState({ cards: new Map([['isaac', isaac]]) });
    setCardLanguages(['en', 'ru']);
    const ru = { nickname: 'Ann', recent: [], animation: 'normal' as const, fullscreen: false, mat: '', dice: 'tear', no_end_turn_ask: false, rejoin: { address: '', token: '', game: '', via: '' as const }, host: noHostSetup, lan_address: '', language: 'ru' };
    vi.mocked(saveSettings).mockImplementation((s) => Promise.resolve(s));
    expect(await useApp.getState().saveSettings(ru)).toBe(true);
    expect(useApp.getState().cards.get('isaac')).toMatchObject({ name: 'Исаак', image: images.ru });
    expect(await useApp.getState().saveSettings({ ...ru, language: 'en' })).toBe(true);
    expect(useApp.getState().cards.get('isaac')).toMatchObject({ name: 'Isaac', image: images.en });
  });

  // A language with card names only (16.10): listed once the cards load,
  // and a setting saved before that names the cards in it then.
  it('lists a card-only language and names the cards in it', async () => {
    setCardLanguages([]);
    setLanguage('xx');
    const isaac = { ref: 'isaac', set: 'b2', type: 'Character Card', name: 'Isaac', names: { en: 'Isaac', xx: 'Xisaac' }, effects: [], footnotes: [], stats: [], rewards: [], image: '' };
    vi.mocked(getCards).mockResolvedValue([isaac]);
    vi.mocked(getDice).mockResolvedValue([]);
    vi.mocked(getMats).mockResolvedValue([]);
    vi.mocked(getBacks).mockResolvedValue([]);
    vi.mocked(getIcons).mockResolvedValue({ coin: '', sword: '', banner: '', sky: [] });
    vi.mocked(getFonts).mockResolvedValue({ title: '', body: '', cyrillic: '' });
    vi.mocked(getSounds).mockResolvedValue({ attackHit: '', effectCast: '', effectImpact: '', diceRoll: '', lootDraw: '', playerDied: '', treasureGain: '' });
    expect(useApp.getState().languages).not.toContain('xx');
    await useApp.getState().loadCards();
    expect(useApp.getState().languages).toContain('xx');
    expect(language()).toBe('xx');
    expect(useApp.getState().cards.get('isaac')?.name).toBe('Xisaac');
    setLanguage('en');
  });

  it('shows the Go error and keeps the old settings', async () => {
    useApp.setState({ settings: { nickname: 'Ann', recent: [], animation: 'normal', fullscreen: false, mat: 'basement', dice: 'tear', no_end_turn_ask: false, rejoin: { address: '', token: '', game: '', via: '' as const }, host: noHostSetup, lan_address: '', language: '' } });
    vi.mocked(saveSettings).mockRejectedValue('nickname must be 1–20 characters');
    expect(await useApp.getState().saveSettings({ nickname: ' ', recent: [], animation: 'normal', fullscreen: false, mat: 'basement', dice: 'tear', no_end_turn_ask: false, rejoin: { address: '', token: '', game: '', via: '' as const }, host: noHostSetup, lan_address: '', language: '' })).toBe(false);
    expect(useApp.getState()).toMatchObject({
      settings: { nickname: 'Ann' },
      error: 'nickname must be 1–20 characters',
    });
  });
});

describe('messageOf', () => {
  it('reads strings, errors and anything else', () => {
    expect(messageOf('bad')).toBe('bad');
    expect(messageOf(new Error('worse'))).toBe('worse');
    expect(messageOf(42)).toBe('something went wrong');
  });
});

describe('seatOf', () => {
  it('takes the seat from the table, or from the welcome after a restart', () => {
    const welcome = { protocol: 1, app: 'a', engine: 'e', token: 't', seat: 1 };
    const table = { game: 'g1', you: 0, seats: null, sets: null };
    expect(seatOf({ lobby: { ...initial.lobby, welcome, table } })).toBe(0);
    expect(seatOf({ lobby: { ...initial.lobby, welcome } })).toBe(1);
    expect(seatOf({ lobby: { ...initial.lobby, welcome: { ...welcome, seat: -1 } } })).toBe(0);
  });
});

describe('recent servers (N-18)', () => {
  const settings = { nickname: 'Ann', recent: [{ address: 'a', name: '' }, { address: 'b', name: '' }], animation: 'normal' as const, fullscreen: false, mat: 'basement', dice: 'tear', no_end_turn_ask: false, rejoin: { address: '', token: '', game: '', via: '' as const }, host: noHostSetup, lan_address: '', language: '' };

  beforeEach(() => {
    vi.mocked(saveSettings).mockImplementation((s) => Promise.resolve(s));
    useApp.setState({ settings, address: 'b' });
  });

  it('a welcome names the server in Recent', async () => {
    useApp.getState().receive({ type: 'welcome', data: { protocol: 5, app: '0.1.0', engine: '0.1.0', token: 't', seat: -1, name: 'Basement' } });
    await vi.waitFor(() => {
      expect(useApp.getState().settings.recent).toEqual([{ address: 'a', name: '' }, { address: 'b', name: 'Basement' }]);
    });
  });

  it('a welcome without a change saves nothing', async () => {
    useApp.getState().receive({ type: 'welcome', data: { protocol: 5, app: '0.1.0', engine: '0.1.0', token: 't', seat: -1 } });
    await Promise.resolve();
    expect(saveSettings).not.toHaveBeenCalled();
  });

  it('forgetServer drops the address', async () => {
    await useApp.getState().forgetServer('a');
    expect(useApp.getState().settings.recent).toEqual([{ address: 'b', name: '' }]);
  });
});

describe('server notices (N-20)', () => {
  it('a lost game goes back to the game list and shows the notice', () => {
    useApp.setState({ screen: 'game', connected: true });
    useApp.getState().receive({ type: 'notice', data: { kind: 'game_lost', game: 'g1' } });
    expect(useApp.getState()).toMatchObject({ screen: 'join', game: null, notice: { kind: 'game_lost' } });
    useApp.getState().endNotice();
    expect(useApp.getState().notice).toBeNull();
  });

  it('a newer engine keeps the game on screen', () => {
    useApp.setState({ screen: 'game', connected: true });
    useApp.getState().receive({ type: 'notice', data: { kind: 'engine_updated' } });
    expect(useApp.getState()).toMatchObject({ screen: 'game', notice: { kind: 'engine_updated' } });
  });

  it('restarting is remembered for the badge', () => {
    useApp.getState().receive({ type: 'restarting' });
    expect(useApp.getState().restarting).toBe(true);
  });

  it('a welcome forgets the intent sent on the old connection (19.86)', () => {
    useApp.setState({ screen: 'game', sent: 7 });
    useApp.getState().receive({ type: 'welcome', data: { protocol: 5, app: '0.1.0', engine: '0.1.0', token: 't', seat: 0 } });
    expect(useApp.getState().sent).toBeNull();
  });
});

describe('the saved seat (15.99)', () => {
  const rejoin = { address: '192.168.1.5', token: 't0k3n', game: 'g1', via: 'lan' as const };
  const settings = { nickname: 'Ann', recent: [], animation: 'normal' as const, fullscreen: false, mat: '', dice: 'tear', no_end_turn_ask: false, rejoin, host: noHostSetup, lan_address: '', language: '' };
  const none = { address: '', token: '', game: '', via: '' };

  beforeEach(() => {
    vi.mocked(saveSettings).mockImplementation((s) => Promise.resolve(s));
    useApp.setState({ settings });
  });

  it('a live seat is offered with its path; the server is asked by its http address', async () => {
    vi.mocked(checkSeat).mockResolvedValue({ live: true, game: 'g1', host: 'Bo', started: true });
    await useApp.getState().checkRejoin();
    expect(checkSeat).toHaveBeenCalledWith(['http://192.168.1.5:4774'], 't0k3n');
    expect(useApp.getState().rejoinable).toEqual({ host: 'Bo', via: 'lan' });
  });

  it('a dead seat is forgotten', async () => {
    vi.mocked(checkSeat).mockResolvedValue({ live: false, game: '', host: '', started: false });
    await useApp.getState().checkRejoin();
    expect(useApp.getState().rejoinable).toBeNull();
    expect(useApp.getState().settings.rejoin).toEqual(none);
  });

  it('no answer keeps the seat for later, without a button', async () => {
    vi.mocked(checkSeat).mockRejectedValue(new Error('server: no answer'));
    useApp.setState({ rejoinable: { host: 'Bo', via: 'lan' } });
    await useApp.getState().checkRejoin();
    expect(useApp.getState().rejoinable).toBeNull();
    expect(saveSettings).not.toHaveBeenCalled();
  });

  it('nothing saved: nothing to ask', async () => {
    useApp.setState({ settings: { ...settings, rejoin: { ...none, via: '' as const } } });
    await useApp.getState().checkRejoin();
    expect(checkSeat).not.toHaveBeenCalled();
  });
});

describe('Recent servers (15.100)', () => {
  const recent = [
    { address: 'games.example.com', name: 'Up' },
    { address: '192.168.1.9', name: '' },
  ];
  const settings = { nickname: 'Ann', recent, animation: 'normal' as const, fullscreen: false, mat: '', dice: 'tear', no_end_turn_ask: false, rejoin: { address: '', token: '', game: '', via: '' as const }, host: noHostSetup, lan_address: '', language: '' };

  beforeEach(() => {
    useApp.setState({ settings });
  });

  it('each row gets its answer: open games, or down', async () => {
    vi.mocked(serverStatus).mockImplementation((bases) =>
      bases[0]?.startsWith('https://games.example.com') ? Promise.resolve({ open: 2 }) : Promise.reject(new Error('server: no answer')),
    );
    await useApp.getState().checkRecent();
    expect(serverStatus).toHaveBeenCalledWith(['http://192.168.1.9:4774']);
    expect(useApp.getState().recentStatus).toEqual({ 'games.example.com': { open: 2 }, '192.168.1.9': 'down' });
  });

  it('rows show at once; a slow answer fills its row alone', async () => {
    let answer: (s: { open: number }) => void = () => undefined;
    vi.mocked(serverStatus).mockImplementation((bases) =>
      bases[0]?.startsWith('https://') ? new Promise((resolve) => (answer = resolve)) : Promise.resolve({ open: 0 }),
    );
    useApp.setState({ recentStatus: { '192.168.1.9': 'down' } });
    const done = useApp.getState().checkRecent();
    expect(useApp.getState().recentStatus).toEqual({});
    await vi.waitFor(() => {
      expect(useApp.getState().recentStatus).toEqual({ '192.168.1.9': { open: 0 } });
    });
    answer({ open: 1 });
    await done;
    expect(useApp.getState().recentStatus['games.example.com']).toEqual({ open: 1 });
  });

  it('a late answer of an older check is dropped', async () => {
    let first: (s: { open: number }) => void = () => undefined;
    vi.mocked(serverStatus).mockImplementationOnce(() => new Promise((resolve) => (first = resolve)));
    vi.mocked(serverStatus).mockResolvedValue({ open: 3 });
    const old = useApp.getState().checkRecent();
    await useApp.getState().checkRecent();
    first({ open: 9 });
    await old;
    expect(useApp.getState().recentStatus['games.example.com']).toEqual({ open: 3 });
  });
});

describe('rollback (ADR 020)', () => {
  const view = (turn: number): GameView => ({
    players: [],
    shop: null,
    monsters: null,
    deck_sizes: [1, 1, 1, 0],
    discards: [null, null, null, null],
    stack: null,
    turn: { active: 0, number: turn, step: 4, loot_plays: 0, attacks: 0, purchases: 0 },
    waiting: { kind: 1, player: 0 },
  });
  const update = (u: Partial<Update> & { step: number }): Envelope => ({
    type: 'update',
    data: { events: null, view: view(1), allowed: null, seats: [{ seat: 1, name: 'Bo', connected: true }], rollback_mode: 'god', rollback_from: 1, ...u },
  });

  it('jumps to the restored table, greys the undone rows and adds a row', async () => {
    const { receive } = useApp.getState();
    receive(update({ step: 0 }));
    receive(update({ step: 2, events: [{ kind: 'attack_declared', player: 1 }] }));
    receive(update({ step: 3, events: [{ kind: 'purchase_declared', player: 1 }], view: view(2) }));
    await vi.waitFor(() => {
      expect(useApp.getState().game?.step).toBe(3);
    });
    // Each row keeps the first step of its update.
    expect(useApp.getState().history.map((e) => e.step)).toEqual([1, 3]);
    expect(useApp.getState().jumped).toBe(false);
    receive(update({ step: 4, rollback: { seat: 1, to: 3, undid: 1 } }));
    await vi.waitFor(() => {
      expect(useApp.getState().game?.step).toBe(4);
    });
    const s = useApp.getState();
    expect(s.jumped).toBe(true);
    expect(s.history.map((e) => [e.text, e.undone ?? false])).toEqual([
      ['Bo declares an attack', false],
      ['Bo declares a purchase', true],
      ['Bo rolls back 1 step', false],
    ]);
    // The undone move is forgotten; the rollback's own has no name (17.7).
    expect(s.moves).toEqual([
      { from: 0, to: 0 },
      { from: 1, to: 2, text: 'Bo declares an attack' },
      { from: 4, to: 4 },
    ]);
    // The next update animates again; its rows start after the rollback.
    receive(update({ step: 5, events: [{ kind: 'attack_declared', player: 1 }] }));
    await vi.waitFor(() => {
      expect(useApp.getState().game?.step).toBe(5);
    });
    expect(useApp.getState().jumped).toBe(false);
    expect(useApp.getState().history.at(-1)?.step).toBe(5);
  });

  it('a Soft vote keeps the step; its end adds a row and the game goes on', async () => {
    const { receive } = useApp.getState();
    const soft = (u: Partial<Update> & { step: number }) => update({ rollback_mode: 'soft', rollback_steps: [{ step: 1, seat: 1 }], ...u });
    useApp.setState({ lobby: { ...useApp.getState().lobby, table: { game: 'g', you: 0, seats: [], sets: [] } } });
    receive(soft({ step: 0 }));
    receive(soft({ step: 1, events: [{ kind: 'attack_declared', player: 1 }] }));
    await vi.waitFor(() => {
      expect(useApp.getState().game?.step).toBe(1);
    });
    // Opening the vote sends the rollback but leaves sent: the step stays.
    useApp.getState().rollBack(1);
    expect(send).toHaveBeenCalledWith('rollback', { step: 1 });
    expect(useApp.getState().sent).toBeNull();
    useApp.getState().rollBack(1, true);
    expect(send).toHaveBeenCalledWith('rollback', { step: 1, ask: true });
    const vote = { step: 1, seat: 0, author: 1, undo: 1, voters: [0], votes: [{ seat: 0, yes: true }], need: 1 };
    receive(soft({ step: 1, rollback_vote: vote }));
    await vi.waitFor(() => {
      expect(useApp.getState().game?.rollback_vote).toEqual(vote);
    });
    receive(soft({ step: 1, rollback_end: { step: 1, seat: 0, why: 'cancelled' } }));
    await vi.waitFor(() => {
      expect(useApp.getState().game?.rollback_vote).toBeUndefined();
    });
    const s = useApp.getState();
    expect(s.history.map((e) => e.text)).toEqual(['Bo declares an attack', 'You cancelled the rollback vote']);
    expect(s.history[0]?.undone).toBeUndefined();
    expect(s.jumped).toBe(false);
  });
});
