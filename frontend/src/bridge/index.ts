// The only module that talks to Wails.
// Everything else imports from here, so a move to Wails v3 or a browser
// build replaces this folder and nothing else (A-14).
import type { CardInfo } from '../cards';
import type { Speed } from '../playback';
import {
  Cards as wailsCards,
  CardSets as wailsCardSets,
  Notice as wailsNotice,
  SaveSettings as wailsSaveSettings,
  Settings as wailsSettings,
  Versions as wailsVersions,
} from '../../wailsjs/go/main/App';
import { Info as hostInfo, Start as hostStart, Stop as hostStop } from '../../wailsjs/go/main/Host';
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
  // recent are the last addresses joined, newest first.
  recent: string[];
  // animation is the animation speed (ST-05).
  animation: Speed;
}

// HostInfo tells the host how friends reach the hosted game.
export interface HostInfo {
  port: number;
  addresses: string[];
}

export async function getVersions(): Promise<Versions> {
  const v = await wailsVersions();
  return { app: v.app, engine: v.engine };
}

export async function getNotice(): Promise<NoticeLine[]> {
  const lines = await wailsNotice();
  return lines.map((line) => line.parts.map((p) => ({ text: p.text, url: p.url })));
}

// speed reads Go's animation setting; Go checks it, so anything else is
// only possible from an older app.
function speed(s: string): Speed {
  return s === 'fast' || s === 'slow' || s === 'off' ? s : 'normal';
}

export async function getSettings(): Promise<Settings> {
  const s = await wailsSettings();
  return { nickname: s.nickname, recent: s.recent, animation: speed(s.animation) };
}

// Go checks and cleans the settings; the result is what it stored.
// A rejection carries Go's message for the player.
export async function saveSettings(s: Settings): Promise<Settings> {
  const saved = await wailsSaveSettings({ nickname: s.nickname, recent: s.recent, animation: s.animation });
  return { nickname: saved.nickname, recent: saved.recent, animation: speed(saved.animation) };
}

// getCardSets lists the card set codes this app has (CD-03).
export function getCardSets(): Promise<string[]> {
  return wailsCardSets();
}

// getCards returns every card the app knows (7.4).
export function getCards(): Promise<CardInfo[]> {
  return wailsCards();
}

// startHosting runs the game server inside the app. If it runs already,
// its info is returned.
export async function startHosting(): Promise<HostInfo> {
  try {
    const info = await hostStart();
    return { port: info.port, addresses: info.addresses };
  } catch (e) {
    const info = await hostInfo();
    if (!info.running) {
      throw e;
    }
    return { port: info.port, addresses: info.addresses };
  }
}

export function stopHosting(): Promise<void> {
  return hostStop();
}

// Opens an https link in the system browser; anything else is ignored.
export function openExternal(url: string): void {
  if (url.startsWith('https://')) {
    BrowserOpenURL(url);
  }
}
