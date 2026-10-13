// The only module that talks to Wails.
// Everything else imports from here, so a move to Wails v3 or a browser
// build replaces this folder and nothing else (A-14).
import { noHostSetup, type HostSetup } from '../hostSetup';
import type { CardInfo, Kind } from '../cards';
import type { Speed } from '../playback';
import {
  Backs as wailsBacks,
  Cards as wailsCards,
  CardSets as wailsCardSets,
  Dice as wailsDice,
  FindGames as wailsFindGames,
  Fonts as wailsFonts,
  Icons as wailsIcons,
  Mats as wailsMats,
  Notice as wailsNotice,
  SaveRecord as wailsSaveRecord,
  SaveSettings as wailsSaveSettings,
  Settings as wailsSettings,
  Sounds as wailsSounds,
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
  mat: string; // ST-03, V-24: a MatInfo or color mat id, '' for the default; Go checks it
  dice: string; // V-19, a DiceSkin id; Go checks it
  // rejoin is the last game joined, so a restarted app takes its seat
  // back (N-08).
  rejoin: Rejoin;
  // host is the last hosted setup, for the next host screen; no password.
  host: HostSetup;
}

export interface Rejoin {
  address: string;
  token: string;
  game: string;
}

export type Theme = 'dark' | 'light';

// fromGo reads Go's settings; Go checks the choices, so anything else is
// only possible from an older app.
function fromGo(s: {
  nickname: string;
  recent: string[];
  animation: string;
  fullscreen: boolean;
  theme: string;
  mat: string;
  dice: string;
  rejoin: Rejoin;
  host?: HostSetup;
}): Settings {
  return {
    nickname: s.nickname,
    recent: s.recent,
    animation: speed(s.animation),
    fullscreen: s.fullscreen,
    theme: s.theme === 'light' ? 'light' : 'dark',
    mat: s.mat,
    dice: s.dice,
    rejoin: { address: s.rejoin.address, token: s.rejoin.token, game: s.rejoin.game },
    host: s.host ? { ...noHostSetup, ...s.host } : noHostSetup,
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

// DiceSkin is one look of the six-sided die (V-19).
export interface DiceSkin {
  id: string;
  name: string;
  // image is the texture URL (see Die); '' when missing: a plain die.
  image: string;
}

// getDice returns every dice look the app knows.
export function getDice(): Promise<DiceSkin[]> {
  return wailsDice();
}

// CardBack is the back image of one deck (V-26).
export interface CardBack {
  kind: Kind;
  image: string; // URL; '' when missing: a plain back
}

// getBacks returns the deck backs; images are '' without the private images.
export async function getBacks(): Promise<CardBack[]> {
  // Go sends only deck kinds, so the cast is safe.
  return (await wailsBacks()) as CardBack[];
}

// Icons are the game's own small images (15.38); each URL is '' without
// the private images, and the app draws its own.
export interface Icons {
  coin: string;
  sword: string;
}

// getIcons returns the icon URLs.
export function getIcons(): Promise<Icons> {
  return wailsIcons();
}

// Fonts are the game's own fonts (15.20); each URL is '' without the
// private images, and the app uses the system font.
export interface Fonts {
  title: string;
  body: string;
  cyrillic: string;
}

// getFonts returns the font URLs.
export function getFonts(): Promise<Fonts> {
  return wailsFonts();
}

// Sounds are the private fireball sounds (15.50); each URL is '' without
// the private images, and the app plays the card hit sound instead.
export interface Sounds {
  fireballCast: string;
  fireballImpact: string;
}

// getSounds returns the sound URLs.
export function getSounds(): Promise<Sounds> {
  return wailsSounds();
}

// MatInfo is one game mat image (ST-03, V-24).
export interface MatInfo {
  id: string;
  name: string;
  image: string; // URL; '' for a color mat
}

// getMats returns the mat images; empty without the private images.
export function getMats(): Promise<MatInfo[]> {
  return wailsMats();
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

// LanServer is a hosting app found on the local network, with its games
// that wait for players (N-13, ADR 009).
export interface LanServer {
  address: string;
  protocol: number;
  app: string;
  games: LanGame[];
}

export interface LanGame {
  host: string;
  seats: number;
  taken: number;
  locked: boolean; // joining needs the password (N-14)
}

// findGames searches the local network for about a second.
export async function findGames(): Promise<LanServer[]> {
  const found = await wailsFindGames();
  return found.map((f) => ({
    address: f.address,
    protocol: f.protocol,
    app: f.app,
    games: f.games.map((g) => ({ host: g.host, seats: g.seats, taken: g.taken, locked: g.locked ?? false })),
  }));
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
