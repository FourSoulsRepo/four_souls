// Table hotkeys (15.19, 15.20): B buys, A attacks, E ends the turn,
// Space passes (F-15).
import { useEffect, useRef } from 'react';

// useKey calls run when key is pressed on the table. Keys typed into a
// text field, or pressed inside a dialog or menu, are left alone. Space
// also stops a focused button from clicking, so it only ever passes.
export function useKey(key: string, run: () => void): void {
  const runRef = useRef(run);
  useEffect(() => {
    runRef.current = run;
  });
  useEffect(() => {
    const ours = (e: KeyboardEvent) =>
      e.key.toLowerCase() === key && !e.ctrlKey && !e.metaKey && !e.altKey && !aside(e.target);
    const onDown = (e: KeyboardEvent) => {
      if (!ours(e)) {
        return;
      }
      if (key === ' ') {
        e.preventDefault();
      }
      if (!e.repeat) {
        runRef.current();
      }
    };
    const onUp = (e: KeyboardEvent) => {
      if (key === ' ' && ours(e)) {
        e.preventDefault();
      }
    };
    window.addEventListener('keydown', onDown, true);
    window.addEventListener('keyup', onUp, true);
    return () => {
      window.removeEventListener('keydown', onDown, true);
      window.removeEventListener('keyup', onUp, true);
    };
  }, [key]);
}

// aside: the key goes to a text field, a dialog, or a menu, not to the table.
function aside(target: EventTarget | null): boolean {
  if (!(target instanceof HTMLElement)) {
    return false;
  }
  return (
    target.isContentEditable ||
    /^(INPUT|TEXTAREA|SELECT)$/.test(target.tagName) ||
    target.closest('[role="dialog"], [role="menu"]') !== null
  );
}
