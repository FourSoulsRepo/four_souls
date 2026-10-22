// The test servers' ports, one per browser project. A file of its own,
// so playwright.config.ts can read it without loading the test runner.
export const serverPorts: Record<string, number> = { chromium: 4799, webkit: 4797 };

// restartPorts are for the servers restart.spec.ts starts and stops
// itself (N-20), one per browser project.
export const restartPorts: Record<string, number> = { chromium: 4795, webkit: 4793 };
