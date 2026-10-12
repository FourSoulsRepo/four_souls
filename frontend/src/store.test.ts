import { beforeEach, describe, expect, it, vi } from 'vitest';
import { getSettings, saveSettings } from './bridge';
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
    vi.mocked(getSettings).mockResolvedValue({ nickname: 'Ann', recent: [], animation: 'normal', fullscreen: false, theme: 'dark', mat: 'felt', dice: 'tear', rejoin: { address: '', token: '', game: '' } });
    await useApp.getState().loadSettings();
    expect(useApp.getState().settings.nickname).toBe('Ann');
  });

  it('keeps the defaults when loading fails', async () => {
    vi.mocked(getSettings).mockRejectedValue('no bindings');
    await useApp.getState().loadSettings();
    expect(useApp.getState().settings.nickname).toBe('');
  });

  it('keeps what Go stored, not what was typed', async () => {
    vi.mocked(saveSettings).mockResolvedValue({ nickname: 'Ann', recent: [], animation: 'normal', fullscreen: false, theme: 'dark', mat: 'felt', dice: 'tear', rejoin: { address: '', token: '', game: '' } });
    expect(await useApp.getState().saveSettings({ nickname: '  Ann  ', recent: [], animation: 'normal', fullscreen: false, theme: 'dark', mat: 'felt', dice: 'tear', rejoin: { address: '', token: '', game: '' } })).toBe(true);
    expect(useApp.getState().settings.nickname).toBe('Ann');
  });

  it('shows the Go error and keeps the old settings', async () => {
    useApp.setState({ settings: { nickname: 'Ann', recent: [], animation: 'normal', fullscreen: false, theme: 'dark', mat: 'felt', dice: 'tear', rejoin: { address: '', token: '', game: '' } } });
    vi.mocked(saveSettings).mockRejectedValue('nickname must be 1–20 characters');
    expect(await useApp.getState().saveSettings({ nickname: ' ', recent: [], animation: 'normal', fullscreen: false, theme: 'dark', mat: 'felt', dice: 'tear', rejoin: { address: '', token: '', game: '' } })).toBe(false);
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
