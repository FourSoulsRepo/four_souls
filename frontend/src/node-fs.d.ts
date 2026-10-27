// The few node:fs calls translations.test.ts makes under Vitest (16.7).
// The frontend has no @types/node: app code runs in the browser, and
// no new dependency is needed for one test file.
declare module 'node:fs' {
  export function existsSync(path: URL): boolean;
  export function mkdirSync(path: URL, options: { recursive: true }): void;
  export function readdirSync(path: URL): string[];
  export function readFileSync(path: URL, encoding: 'utf8'): string;
  export function rmSync(path: URL): void;
  export function writeFileSync(path: URL, data: string): void;
}
