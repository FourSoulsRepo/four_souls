// The test servers' ports, one per browser project. A file of its own,
// so playwright.config.ts can read it without loading the test runner.
export const serverPorts: Record<string, number> = { chromium: 4799, webkit: 4797 };
