import { AnimatePresence, motion } from 'motion/react';
import { useEffect, useLayoutEffect, useRef, useState, type RefObject } from 'react';
import { getVersions, quitApp, type Versions } from '../bridge';
import { t } from '../i18n';
import { useEscape } from '../keys';
import { scaleFor } from '../playback';
import { useApp } from '../store';
import './MainMenu.css';

// The main menu: Play opens LAN, Online and Learn to play in place,
// under itself (M-03, M-06). Features still to come show "Soon".
export function MainMenu() {
  const [versions, setVersions] = useState<Versions | null>(null);
  const show = useApp((s) => s.show);
  const openPath = useApp((s) => s.openPath);
  const open = useApp((s) => s.playOpen);
  const setOpen = useApp((s) => s.setPlayOpen);
  const scale = useApp((s) => scaleFor(s.settings.animation, 0, 0));
  const [nav, height] = useOpenHeight();
  const rejoinable = useApp((s) => s.rejoinable);
  const rejoin = useApp((s) => s.rejoin);
  const checkRejoin = useApp((s) => s.checkRejoin);
  const seat = useApp((s) => s.settings.rejoin.token);

  useEffect(() => {
    getVersions().then(setVersions, () => {
      setVersions(null);
    });
  }, []);

  useEscape(() => {
    setOpen(false);
  });

  // Is the game this app left still on? Never waits: the menu shows
  // at once, the button when the answer comes.
  useEffect(() => {
    void checkRejoin();
  }, [checkRejoin, seat]);

  return (
    <div className="menu">
      <Title />
      <nav className="menu-buttons" ref={nav} style={{ height }}>
        <button
          type="button"
          aria-expanded={open}
          onClick={() => {
            setOpen(!open);
          }}
        >
          {t('menu.play')}
        </button>
        <AnimatePresence initial={false}>
          {open && (
            // anim: menu-group
            <motion.div
              key="play"
              className="menu-group"
              role="group"
              aria-label={t('menu.play')}
              initial={scale > 0 ? { height: 0, opacity: 0 } : false}
              animate={{ height: 'auto', opacity: 1 }}
              exit={{ height: 0, opacity: 0 }}
              transition={{ duration: 0.18 * scale }}
            >
              <div className="menu-group-buttons">
              <button
                type="button"
                onClick={() => {
                  openPath('lan');
                }}
              >
                {t('menu.lan')}
              </button>
              <button
                type="button"
                onClick={() => {
                  openPath('online');
                }}
              >
                {t('menu.online')}
              </button>
              <Soon name={t('menu.learn')} />
              </div>
            </motion.div>
          )}
        </AnimatePresence>
        <Soon name={t('menu.records')} />
        <button
          type="button"
          onClick={() => {
            show('settings');
          }}
        >
          {t('menu.settings')}
        </button>
        <button type="button" onClick={quitApp}>
          {t('menu.exit')}
        </button>
      </nav>
      {rejoinable && (
        // In the corner: the menu is centred, and a button that came
        // with the answer moved the others under the pointer (19.49).
        <div className="menu-buttons menu-rejoin">
          <button
            type="button"
            onClick={() => {
              void rejoin();
            }}
          >
            {t('menu.rejoin', { game: rejoinLabel(rejoinable) })}
          </button>
        </div>
      )}
      {versions && (
        <div className="menu-versions">
          {t('menu.versions', { app: versions.app, engine: versions.engine })}
        </div>
      )}
    </div>
  );
}

// rejoinLabel names the game to rejoin and how it was joined (M-05).
function rejoinLabel(r: { host: string; via: '' | 'lan' | 'online' }): string {
  const game = r.host ? t('menu.rejoin_host', { host: r.host }) : t('menu.rejoin_yours');
  return r.via ? `${game} · ${r.via === 'lan' ? t('menu.lan') : t('menu.online')}` : game;
}

// playGroup is how many buttons the Play group holds.
const playGroup = 3;

// Soon is a menu button for a feature still to come (M-07).
function Soon({ name }: { name: string }) {
  return (
    <button type="button" disabled aria-label={t('menu.soon_label', { name })}>
      {name}
      <span className="soon">{t('menu.soon')}</span>
    </button>
  );
}

// useOpenHeight is the menu's height with the Play group open, kept
// while it is closed too: the menu is centred, so a menu that grew
// moved Play and the title up, and the room keeps Exit on screen. The
// group's buttons are as tall as the others.
function useOpenHeight(): [RefObject<HTMLElement | null>, number | undefined] {
  const ref = useRef<HTMLElement>(null);
  const [height, setHeight] = useState<number | undefined>(undefined);
  useLayoutEffect(() => {
    const nav = ref.current;
    if (!nav) {
      return;
    }
    const buttons = [...nav.children].filter((c): c is HTMLElement => c instanceof HTMLButtonElement);
    const measure = () => {
      const gap = parseFloat(getComputedStyle(nav).rowGap) || 0;
      const sizes = buttons.map(boxHeight);
      const each = Math.max(...sizes);
      const closed = sizes.reduce((h, b) => h + b, 0) + gap * (buttons.length - 1);
      setHeight(closed + playGroup * (each + gap));
    };
    measure();
    // The game font loads late and changes the button size.
    const watch = new ResizeObserver(measure);
    for (const b of buttons) {
      watch.observe(b);
    }
    return () => {
      watch.disconnect();
    };
  }, []);
  return [ref, height];
}

// boxHeight is an element's laid-out height, in fractions of a pixel:
// offsetHeight rounds it, and the bounding box grows with the tilt.
function boxHeight(el: HTMLElement): number {
  const cs = getComputedStyle(el);
  const sides = ['padding-top', 'padding-bottom', 'border-top-width', 'border-bottom-width'];
  const extra = cs.boxSizing === 'border-box' ? 0 : sides.reduce((h, k) => h + parseFloat(cs.getPropertyValue(k)), 0);
  return parseFloat(cs.height) + extra;
}

// Title is the foursouls.ru banner (15.41), or the text title without the
// private images or when the banner fails to load. The text shows until
// the banner has loaded: the menu is centred, and a title that changed
// size after it showed moved the buttons under the pointer (19.49).
function Title() {
  const image = useApp((s) => s.icons.banner);
  const [loaded, setLoaded] = useState('');
  useEffect(() => {
    if (image === '') {
      return;
    }
    const img = new Image();
    img.onload = () => {
      setLoaded(image);
    };
    img.src = image;
    return () => {
      img.onload = null;
    };
  }, [image]);
  if (image === '' || image !== loaded) {
    return <h1>{t('app.title')}</h1>;
  }
  return (
    <h1 className="menu-banner">
      <img src={image} alt={t('app.title')} draggable={false} />
    </h1>
  );
}
