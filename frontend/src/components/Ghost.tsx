import { counterColor, ghostPath } from '../counters';

// Ghost is fear's small icon (18.26): a ghost-white ghost with an ink
// edge and two eyes, in a character's zoom. Its shape tells fear apart
// without its colour. Decorative: the text beside it names the fear.
export function Ghost({ className }: { className: string }) {
  return (
    <svg className={className} viewBox="0 0 100 100" aria-hidden="true">
      <path d={ghostPath} fill={counterColor('fear')} stroke="#1d120a" strokeWidth="9" strokeLinejoin="round" />
      <ellipse cx="37" cy="44" rx="7" ry="10" fill="#1d120a" />
      <ellipse cx="63" cy="44" rx="7" ry="10" fill="#1d120a" />
    </svg>
  );
}
