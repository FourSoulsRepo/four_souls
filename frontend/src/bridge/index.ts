// The only module that talks to Wails.
// Everything else imports from here, so a move to Wails v3 or a browser
// build replaces this folder and nothing else (A-14).
import { Notice as wailsNotice, Versions as wailsVersions } from '../../wailsjs/go/main/App';
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

export async function getVersions(): Promise<Versions> {
  const v = await wailsVersions();
  return { app: v.app, engine: v.engine };
}

export async function getNotice(): Promise<NoticeLine[]> {
  const lines = await wailsNotice();
  return lines.map((line) => line.parts.map((p) => ({ text: p.text, url: p.url })));
}

// Opens an https link in the system browser; anything else is ignored.
export function openExternal(url: string): void {
  if (url.startsWith('https://')) {
    BrowserOpenURL(url);
  }
}
