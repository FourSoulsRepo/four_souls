import { useLayoutEffect, useRef } from 'react';
import { startWaits } from '../choice';
import { feel } from '../feel';
import { t } from '../i18n';
import type { GameView } from '../protocol';
import { useApp, viewScale } from '../store';
import { still } from './lunge';
import './Throw.css';

// StartWait tells everyone else who makes a start-of-game choice now,
// low on the board: "Bo is choosing", one line per chooser (20.22, ADR
// 027 "Start of game" 16.5). No window: the table stays usable. It
// grows in as it shows.
export function StartWait({ view, you, name }: { view: GameView; you: number; name: (p: number) => string }) {
  const lines = startWaits(view, you, name);
  const ref = useRef<HTMLDivElement>(null);
  const shown = lines.length > 0;
  useLayoutEffect(() => {
    const el = ref.current;
    if (!shown || !el || still()) {
      return;
    }
    const s = feel.start;
    // anim: start-wait
    const a = el.animate([{ transform: 'scale(0)' }, { transform: 'scale(1)' }], { duration: s.wait * 1000 * viewScale(useApp.getState()), easing: s.waitCurve });
    return () => {
      a.cancel();
    };
  }, [shown]);
  if (!shown) {
    return null;
  }
  return (
    <div className="start-wait-layer">
      <div ref={ref} className="start-wait" role="status" aria-label={t('start.waiting_label')}>
        {lines.map((l) => (
          <p key={l}>{l}</p>
        ))}
      </div>
    </div>
  );
}
