// The real dedicated server the e2e tests play on (playwright.config.ts
// starts it). The fake Host binding points at it (?server=PORT).
export const serverPort = 4799;

// appURL opens the app with fake bindings, no splash, hosting on the
// test server.
export const appURL = `/?fake&nosplash&server=${String(serverPort)}`;
