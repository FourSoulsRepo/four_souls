import { useEffect } from 'react';
import { useApp } from '../store';

// How long a reason stays, in milliseconds.
const whyMs = 2500;

// WhyNote shows the reason of the last refused click where the player
// clicked (15.29, F-18): a wrong aim (15.33), or a card or turn action
// they can't use now. It fades after a while; a new click shows it
// again.
export function WhyNote() {
  const why = useApp((s) => s.why);
  const endWhy = useApp((s) => s.endWhy);
  useEffect(() => {
    if (!why) {
      return;
    }
    const t = setTimeout(endWhy, whyMs);
    return () => {
      clearTimeout(t);
    };
  }, [why, endWhy]);
  if (!why) {
    return null;
  }
  return (
    <div key={why.n} className="aim-why" role="status" style={{ transform: `translate(${String(why.x)}px, ${String(why.y)}px)` }}>
      <div>{why.text}</div>
    </div>
  );
}
