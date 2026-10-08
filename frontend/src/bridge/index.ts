// The only module that talks to Wails.
// Everything else imports from here, so a move to Wails v3 or a browser
// build replaces this folder and nothing else (A-14).
import { Versions as wailsVersions } from '../../wailsjs/go/main/App';

export interface Versions {
  app: string;
  engine: string;
}

export async function getVersions(): Promise<Versions> {
  const v = await wailsVersions();
  return { app: v.app, engine: v.engine };
}
