import { AnimatePresence, motion, useIsPresent, type HTMLMotionProps } from 'motion/react';
import type { ReactNode } from 'react';
import { scaleFor } from '../playback';
import { useApp } from '../store';

// useRowScale is how long row moves take, as a factor of their normal
// time; 0 when animations are off (ST-05).
function useRowScale(): number {
  return useApp((s) => scaleFor(s.settings.animation, 0, 0));
}

// Rows holds the moving rows of the stack panel or the history strip
// (15.51, V-39). A leaving row stays until it has faded out. With
// animations off, rows come and go at once.
export function Rows({ children }: { children: ReactNode }) {
  const scale = useRowScale();
  return scale > 0 ? <AnimatePresence initial={false}>{children}</AnimatePresence> : children;
}

// Row is one row in Rows, under a stable key: a new row fades in, the
// others slide to their new place, and a leaving row fades out before
// the rest close the gap. A leaving row is no place for effects or a
// drop (data-zone, data-stack), and screen readers skip it.
export function Row(props: HTMLMotionProps<'li'>) {
  const scale = useRowScale();
  const present = useIsPresent();
  const gone = present ? {} : { 'data-zone': undefined, 'data-stack': undefined, 'aria-hidden': true, style: { ...props.style, pointerEvents: 'none' as const } };
  return (
    // anim: row-move
    <motion.li
      {...props}
      {...gone}
      layout={scale > 0 ? 'position' : false}
      initial={scale > 0 ? { opacity: 0 } : false}
      animate={{ opacity: 1 }}
      exit={{ opacity: 0 }}
      transition={{ duration: 0.2 * scale, layout: { duration: 0.25 * scale } }}
    />
  );
}
