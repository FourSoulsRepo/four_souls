import { useEffect } from 'react';
import './Splash.css';

const SPLASH_MS = 3000;

interface Props {
  onDone: () => void;
}

export function Splash({ onDone }: Props) {
  useEffect(() => {
    const timer = window.setTimeout(onDone, SPLASH_MS);
    return () => {
      window.clearTimeout(timer);
    };
  }, [onDone]);

  return (
    <div className="splash" onClick={onDone}>
      <h1>Four Souls</h1>
    </div>
  );
}
