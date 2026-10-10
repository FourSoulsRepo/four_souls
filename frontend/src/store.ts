// The app's one store (ADR 008). Screens read it; only the bridge
// results and, later, the event queue write it.
import { create } from 'zustand';
import { getSettings, saveSettings, type Settings } from './bridge';

export type Screen = 'splash' | 'menu' | 'settings';

interface AppState {
  screen: Screen;
  settings: Settings;
  // error is the last failure to show the player; '' if none.
  error: string;
  show: (screen: Screen) => void;
  loadSettings: () => Promise<void>;
  // saveSettings reports whether Go stored them.
  saveSettings: (s: Settings) => Promise<boolean>;
}

// errorText turns a rejected bridge call into a message for the player.
// Wails rejects with Go's error string.
export function errorText(e: unknown): string {
  if (typeof e === 'string') {
    return e;
  }
  if (e instanceof Error) {
    return e.message;
  }
  return 'something went wrong';
}

export const useApp = create<AppState>()((set) => ({
  screen: 'splash',
  settings: { nickname: '' },
  error: '',
  show: (screen) => {
    set({ screen, error: '' });
  },
  loadSettings: async () => {
    try {
      set({ settings: await getSettings() });
    } catch {
      // Keep the defaults; Go already falls back to them.
    }
  },
  saveSettings: async (s) => {
    try {
      set({ settings: await saveSettings(s), error: '' });
      return true;
    } catch (e) {
      set({ error: errorText(e) });
      return false;
    }
  },
}));
