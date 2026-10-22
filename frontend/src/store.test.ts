import { beforeEach, describe, expect, it, vi } from 'vitest';
import { getSettings, saveSettings } from './bridge';
import { noHostSetup } from './hostSetup';
import { messageOf, seatOf, useApp } from './store';

vi.mock('./bridge', () => ({
  getSettings: vi.fn(),
  saveSettings: vi.fn(),
}));
vi.mock('./net', () => ({ drop: vi.fn() }));

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
    vi.mocked(getSettings).mockResolvedValue({ nickname: 'Ann', recent: [], animation: 'normal', fullscreen: false, mat: 'basement', dice: 'tear', no_end_turn_ask: false, rejoin: { address: '', token: '', game: '' }, host: noHostSetup, lan_address: '' });
    await useApp.getState().loadSettings();
    expect(useApp.getState().settings.nickname).toBe('Ann');
  });

  it('keeps the defaults when loading fails', async () => {
    vi.mocked(getSettings).mockRejectedValue('no bindings');
    await useApp.getState().loadSettings();
    expect(useApp.getState().settings.nickname).toBe('');
  });

  it('keeps what Go stored, not what was typed', async () => {
    vi.mocked(saveSettings).mockResolvedValue({ nickname: 'Ann', recent: [], animation: 'normal', fullscreen: false, mat: 'basement', dice: 'tear', no_end_turn_ask: false, rejoin: { address: '', token: '', game: '' }, host: noHostSetup, lan_address: '' });
    expect(await useApp.getState().saveSettings({ nickname: '  Ann  ', recent: [], animation: 'normal', fullscreen: false, mat: 'basement', dice: 'tear', no_end_turn_ask: false, rejoin: { address: '', token: '', game: '' }, host: noHostSetup, lan_address: '' })).toBe(true);
    expect(useApp.getState().settings.nickname).toBe('Ann');
  });

  it('shows the Go error and keeps the old settings', async () => {
    useApp.setState({ settings: { nickname: 'Ann', recent: [], animation: 'normal', fullscreen: false, mat: 'basement', dice: 'tear', no_end_turn_ask: false, rejoin: { address: '', token: '', game: '' }, host: noHostSetup, lan_address: '' } });
    vi.mocked(saveSettings).mockRejectedValue('nickname must be 1–20 characters');
    expect(await useApp.getState().saveSettings({ nickname: ' ', recent: [], animation: 'normal', fullscreen: false, mat: 'basement', dice: 'tear', no_end_turn_ask: false, rejoin: { address: '', token: '', game: '' }, host: noHostSetup, lan_address: '' })).toBe(false);
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
  const settings = { nickname: 'Ann', recent: [{ address: 'a', name: '' }, { address: 'b', name: '' }], animation: 'normal' as const, fullscreen: false, mat: 'basement', dice: 'tear', no_end_turn_ask: false, rejoin: { address: '', token: '', game: '' }, host: noHostSetup, lan_address: '' };

  beforeEach(() => {
    vi.mocked(saveSettings).mockImplementation((s) => Promise.resolve(s));
    useApp.setState({ settings, address: 'b' });
  });

  it('a welcome names the server in Recent', async () => {
    useApp.getState().receive({ type: 'welcome', data: { protocol: 4, app: '0.1.0', engine: '0.1.0', token: 't', seat: -1, name: 'Basement' } });
    await vi.waitFor(() => {
      expect(useApp.getState().settings.recent).toEqual([{ address: 'a', name: '' }, { address: 'b', name: 'Basement' }]);
    });
  });

  it('a welcome without a change saves nothing', async () => {
    useApp.getState().receive({ type: 'welcome', data: { protocol: 4, app: '0.1.0', engine: '0.1.0', token: 't', seat: -1 } });
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
});
