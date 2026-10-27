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
  CheckSeat as wailsCheckSeat,
  DeckRatios as wailsDeckRatios,
  Dice as wailsDice,
  FindGames as wailsFindGames,
  Fonts as wailsFonts,
  Icons as wailsIcons,
  Mats as wailsMats,
  Notice as wailsNotice,
  SaveRecord as wailsSaveRecord,
  SaveSettings as wailsSaveSettings,
  ServerStatus as wailsServerStatus,
  Settings as wailsSettings,
  Sounds as wailsSounds,
  Versions as wailsVersions,
} from '../../wailsjs/go/main/App';
import { settings } from '../../wailsjs/go/models';
import { Info as hostInfo, Start as hostStart, Stop as hostStop } from '../../wailsjs/go/main/Host';
import { BrowserOpenURL, Quit, WindowFullscreen, WindowUnfullscreen } from '../../wailsjs/runtime/runtime';

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
  // recent are the last servers joined, newest first.
  recent: RecentServer[];
  // animation is the animation speed (ST-05).
  animation: Speed;
  fullscreen: boolean; // ST-01
  mat: string; // ST-03, V-24: a MatInfo or color mat id, '' for the default; Go checks it
  dice: string; // V-19, a DiceSkin id; Go checks it
  // no_end_turn_ask ends the turn without asking while plays are left (V-46).
  no_end_turn_ask: boolean;
  // rejoin is the last game joined, so a restarted app takes its seat
  // back (N-08).
  rejoin: Rejoin;
  // host is the last hosted setup, for the next host screen; no password.
  host: HostSetup;
  // lan_address is the address last typed on the LAN screen (15.97).
  lan_address: string;
  // language is the interface language (ADR 021); '' follows the OS.
  language: string;
}

// RecentServer is a server joined lately; name is the one it gave in
// its welcome, '' to show the address (N-18).
export interface RecentServer {
  address: string;
  name: string;
}

export interface Rejoin {
  address: string;
  token: string;
  game: string;
  // via is the Play path it was joined by (15.99); '' from an older app.
  via: '' | 'lan' | 'online';
}

// fromGo reads Go's settings; Go checks the choices, so anything else is
// only possible from an older app.
function fromGo(s: {
  nickname: string;
  recent: RecentServer[];
  animation: string;
  fullscreen: boolean;
  mat: string;
  dice: string;
  no_end_turn_ask?: boolean;
  rejoin: Omit<Rejoin, 'via'> & { via?: string };
  host?: HostSetup;
  lan_address?: string;
  language?: string;
}): Settings {
  return {
    nickname: s.nickname,
    recent: s.recent.map((r) => ({ address: r.address, name: r.name })),
    animation: speed(s.animation),
    fullscreen: s.fullscreen,
    mat: s.mat,
    dice: s.dice,
    no_end_turn_ask: s.no_end_turn_ask ?? false,
    rejoin: { address: s.rejoin.address, token: s.rejoin.token, game: s.rejoin.game, via: via(s.rejoin.via) },
    host: s.host ? { ...noHostSetup, ...s.host } : noHostSetup,
    lan_address: s.lan_address ?? '',
    language: s.language ?? '',
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

function via(s: string | undefined): Rejoin['via'] {
  return s === 'lan' || s === 'online' ? s : '';
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

// deckRatiosBlocked says why a game with these sets and players can't
// use the deck ratios, or '' when it can (GS-07).
export function deckRatiosBlocked(sets: string[], players: number): Promise<string> {
  return wailsDeckRatios(sets, players);
}

// getCards returns every card the app knows (7.4).
export function getCards(): Promise<CardInfo[]> {
  return wailsCards();
}

// DiceSkin is one look of the six-sided die (V-19).
export interface DiceSkin {
  id: string;
  name: string;
  set: string; // "Holy" or "Unholy", for the settings menu
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
  banner: string;
  // sky is the menu background (15.42): the sky, then its star layers,
  // far to near; empty without the images.
  sky: string[];
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

// Sounds are the damage sounds (15.50), the dice sound (15.69), the loot
// draw sound (18.7), the player death sound (18.8) and the treasure gain
// sound (18.9); each URL is '' without the private images, and that
// sound stays silent.
export interface Sounds {
  attackHit: string;
  effectCast: string;
  effectImpact: string;
  diceRoll: string;
  lootDraw: string;
  playerDied: string;
  treasureGain: string;
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
// (N-13, ADR 009).
export interface LanServer {
  address: string;
  protocol: number;
  app: string;
  games: LanGame[];
}

export interface LanGame {
  id: string;
  host: string;
  seats: number;
  taken: number;
  locked: boolean; // joining needs the password (N-14)
  // started: in progress; only who held a seat can rejoin (15.97).
  started: boolean;
}

// findGames searches the local network for about a second.
export async function findGames(): Promise<LanServer[]> {
  const found = await wailsFindGames();
  return found.map((f) => ({
    address: f.address,
    protocol: f.protocol,
    app: f.app,
    games: f.games.map((g) => ({
      id: g.id,
      host: g.host,
      seats: g.seats,
      taken: g.taken,
      locked: g.locked ?? false,
      started: g.started ?? false,
    })),
  }));
}

// SeatInfo says whether a token still holds a seat on a server (15.99).
export interface SeatInfo {
  live: boolean;
  game: string;
  host: string;
  started: boolean;
}

// checkSeat asks the server at bases (httpBases) whether token still
// holds a seat; it rejects when no base answers.
export async function checkSeat(bases: string[], token: string): Promise<SeatInfo> {
  const s = await wailsCheckSeat(bases, token);
  return { live: s.live, game: s.game ?? '', host: s.host ?? '', started: s.started ?? false };
}

// serverStatus asks the server at bases (httpBases) how many games wait
// for players (15.100); it rejects when no base answers.
export async function serverStatus(bases: string[]): Promise<{ open: number }> {
  const s = await wailsServerStatus(bases);
  return { open: s.open };
}

// setFullscreen switches the window (ST-01).
export function setFullscreen(on: boolean): void {
  if (on) {
    WindowFullscreen();
  } else {
    WindowUnfullscreen();
  }
}

// quitApp closes the app (M-02).
export function quitApp(): void {
  Quit();
}

// Opens an https link in the system browser; anything else is ignored.
export function openExternal(url: string): void {
  if (url.startsWith('https://')) {
    BrowserOpenURL(url);
  }
}
