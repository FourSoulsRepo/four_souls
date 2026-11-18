import type { Page } from '@playwright/test';

// giveLooks gives the page two dice looks of plain colors, so a test
// tells them apart without the private textures, and picks one as the
// viewer's own. The game's seats keep the looks sent in hello.
export async function giveLooks(page: Page, own: string): Promise<Record<string, string>> {
  return page.evaluate((own) => {
    interface S {
      settings: { dice: string };
    }
    const store = (window as unknown as { fourSouls: { store: { getState: () => S; setState: (p: object) => void } } }).fourSouls.store;
    const paint = (color: string) => {
      const c = document.createElement('canvas');
      c.width = c.height = 96;
      const g = c.getContext('2d');
      if (g) {
        g.fillStyle = color;
        g.fillRect(0, 0, 96, 96);
      }
      return c.toDataURL();
    };
    const images = { tear: paint('#b33'), gold: paint('#db3') };
    store.setState({
      dice: Object.entries(images).map(([id, image]) => ({ id, name: id, image })),
      settings: { ...store.getState().settings, dice: own },
    });
    return images;
  }, own);
}
