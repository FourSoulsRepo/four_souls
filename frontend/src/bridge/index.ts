// The only module that talks to Wails.
// Everything else imports from here, so a move to Wails v3 or a browser
// build replaces this folder and nothing else (A-14).
import type { CardInfo } from '../cards';
import type { Speed } from '../playback';
import {
  Cards as wailsCards,
  CardSets as wailsCardSets,
  Notice as wailsNotice,
  SaveRecord as wailsSaveRecord,
  SaveSettings as wailsSaveSettings,
  Settings as wailsSettings,
  Versions as wailsVersions,
} from '../../wailsjs/go/main/App';
import { settings } from '../../wailsjs/go/models';
import { Info as hostInfo, Start as hostStart, Stop as hostStop } from '../../wailsjs/go/main/Host';
import { BrowserOpenURL, WindowFullscreen, WindowUnfullscreen } from '../../wailsjs/runtime/runtime';

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
  fullscreen: boolean; // ST-01
  theme: Theme; // ST-02, menus and panels
  mat: Mat; // ST-03
  // rejoin is the last game joined, so a restarted app takes its seat
  // back (N-08).
  rejoin: Rejoin;
}

export interface Rejoin {
  address: string;
  token: string;
  game: string;
}

export type Theme = 'dark' | 'light';
export type Mat = 'felt' | 'basement' | 'cathedral';

// fromGo reads Go's settings; Go checks the choices, so anything else is
// only possible from an older app.
function fromGo(s: {
  nickname: string;
  recent: string[];
  animation: string;
  fullscreen: boolean;
  theme: string;
  mat: string;
  rejoin: Rejoin;
}): Settings {
  return {
    nickname: s.nickname,
    recent: s.recent,
    animation: speed(s.animation),
    fullscreen: s.fullscreen,
    theme: s.theme === 'light' ? 'light' : 'dark',
    mat: s.mat === 'basement' || s.mat === 'cathedral' ? s.mat : 'felt',
    rejoin: { address: s.rejoin.address, token: s.rejoin.token, game: s.rejoin.game },
  };
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

function speed(s: string): Speed {
  return s === 'fast' || s === 'slow' || s === 'off' ? s : 'normal';
}

export async function getSettings(): Promise<Settings> {
  return fromGo(await wailsSettings());
}

// Go checks and cleans the settings; the result is what it stored.
// A rejection carries Go's message for the player.
export async function saveSettings(s: Settings): Promise<Settings> {
  return fromGo(await wailsSaveSettings(settings.Settings.createFrom(s)));
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

// saveRecord asks where to save a finished game's record and downloads
// it; '' means the player cancelled (RP-05).
export function saveRecord(link: string): Promise<string> {
  return wailsSaveRecord(link);
}

// setFullscreen switches the window (ST-01).
export function setFullscreen(on: boolean): void {
  if (on) {
    WindowFullscreen();
  } else {
    WindowUnfullscreen();
  }
}

// Opens an https link in the system browser; anything else is ignored.
export function openExternal(url: string): void {
  if (url.startsWith('https://')) {
    BrowserOpenURL(url);
  }
}
