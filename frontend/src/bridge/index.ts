// The only module that talks to Wails.
// Everything else imports from here, so a move to Wails v3 or a browser
// build replaces this folder and nothing else (A-14).
import {
  Notice as wailsNotice,
  SaveSettings as wailsSaveSettings,
  Settings as wailsSettings,
  Versions as wailsVersions,
} from '../../wailsjs/go/main/App';
import { BrowserOpenURL } from '../../wailsjs/runtime/runtime';

export interface Versions {
  app: string;
  engine: string;
}

export interface NoticePart {
  text: string;
  url?: string;
}

export type NoticeLine = NoticePart[];

// Settings are kept by Go in the user's config folder.
export interface Settings {
  nickname: string;
}

export async function getVersions(): Promise<Versions> {
  const v = await wailsVersions();
  return { app: v.app, engine: v.engine };
}

export async function getNotice(): Promise<NoticeLine[]> {
  const lines = await wailsNotice();
  return lines.map((line) => line.parts.map((p) => ({ text: p.text, url: p.url })));
}

export async function getSettings(): Promise<Settings> {
  const s = await wailsSettings();
  return { nickname: s.nickname };
}

// Go checks and cleans the settings; the result is what it stored.
// A rejection carries Go's message for the player.
export async function saveSettings(s: Settings): Promise<Settings> {
  const saved = await wailsSaveSettings({ nickname: s.nickname });
  return { nickname: saved.nickname };
}

// Opens an https link in the system browser; anything else is ignored.
export function openExternal(url: string): void {
  if (url.startsWith('https://')) {
    BrowserOpenURL(url);
  }
}
