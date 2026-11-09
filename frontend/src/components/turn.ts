// A new turn's banner (20.8, ADR 027 "Turn start"). It lives in the
// effect layer of its own (lunge.ts), so it outlives its effect:
// playback goes on after feel.turn.wait, and the turn's draw plays
// while the banner still shows. It is visual only: it takes no pointer
// and locks nothing; the buttons below stay live. The viewer's own turn
// plays the turn-start sound as it shows (20.17).
import { feel } from '../feel';
import { playSound } from '../sounds';
import { layer, still } from './lunge';

// turnFrames are the banner's keyframes over its life: it fades in as
// it grows to full size, punches down a little, holds, then shrinks and
// fades. With less motion it only fades in and out.
export function turnFrames(calm: boolean): Keyframe[] {
  const { grow, from, punch, punchTime, hide, out, curve } = feel.turn;
  const total = hide + out;
  const k = (s: number) => `scale(${String(calm ? 1 : s)})`;
  return [
    { offset: 0, transform: k(from), opacity: 0, easing: curve },
    { offset: grow / total, transform: k(1), opacity: 1, easing: curve },
    { offset: (grow + punchTime) / total, transform: k(punch), opacity: 1, easing: 'linear' },
    { offset: hide / total, transform: k(punch), opacity: 1, easing: curve },
    { offset: 1, transform: k(from), opacity: 0 },
  ];
}

// turnBanner shows the banner and calls done after feel.turn.wait. A
// banner still on screen from a turn before goes at once.
export function turnBanner(text: string, own: boolean, scale: number, done: () => void): void {
  const host = layer();
  for (const old of host.querySelectorAll('.effect-turn')) {
    old.remove();
  }
  const box = document.createElement('div');
  box.className = 'effect-turn';
  box.dataset.testid = 'turn-banner';
  const label = document.createElement('div');
  label.className = 'effect-banner effect-turn-text';
  label.textContent = text;
  box.append(label);
  host.append(box);
  if (own) {
    playSound('turnStart');
  }
  const { hide, out, wait } = feel.turn;
  // anim: turn-banner
  const a = label.animate(turnFrames(still()), { duration: (hide + out) * 1000 * scale, fill: 'forwards', id: 'turn-banner' });
  const end = () => {
    box.remove();
  };
  a.onfinish = end;
  a.oncancel = end;
  window.setTimeout(done, wait * 1000 * scale);
}
