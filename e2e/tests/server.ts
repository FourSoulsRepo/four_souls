// The real dedicated servers the e2e tests play on (playwright.config.ts
// starts them): one per browser project, so their games never meet. The
// fake Host binding points at the project's server (?server=PORT).
import { test } from '@playwright/test';
import { serverPorts } from './ports';

// serverPort is the test server of the running project.
export function serverPort(): number {
  return serverPorts[test.info().project.name] ?? 4799;
}

// appURL opens the app with fake bindings, no splash, hosting on the
// project's test server, or on port.
export function appURL(port = serverPort()): string {
  return `/?fake&nosplash&server=${String(port)}`;
}
