import { beforeEach, describe, expect, it, vi } from 'vitest';
import { getSettings, saveSettings } from './bridge';
import { messageOf, useApp } from './store';

vi.mock('./bridge', () => ({
  getSettings: vi.fn(),
  saveSettings: vi.fn(),
}));
vi.mock('./net', () => ({}));

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
    vi.mocked(getSettings).mockResolvedValue({ nickname: 'Ann', recent: [] });
    await useApp.getState().loadSettings();
    expect(useApp.getState().settings.nickname).toBe('Ann');
  });

  it('keeps the defaults when loading fails', async () => {
    vi.mocked(getSettings).mockRejectedValue('no bindings');
    await useApp.getState().loadSettings();
    expect(useApp.getState().settings.nickname).toBe('');
  });

  it('keeps what Go stored, not what was typed', async () => {
    vi.mocked(saveSettings).mockResolvedValue({ nickname: 'Ann', recent: [] });
    expect(await useApp.getState().saveSettings({ nickname: '  Ann  ', recent: [] })).toBe(true);
    expect(useApp.getState().settings.nickname).toBe('Ann');
  });

  it('shows the Go error and keeps the old settings', async () => {
    useApp.setState({ settings: { nickname: 'Ann', recent: [] } });
    vi.mocked(saveSettings).mockRejectedValue('nickname must be 1–20 characters');
    expect(await useApp.getState().saveSettings({ nickname: ' ', recent: [] })).toBe(false);
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
