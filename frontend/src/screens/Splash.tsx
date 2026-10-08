import { useEffect, useState } from 'react';
import { getNotice, openExternal, type NoticeLine } from '../bridge';
import './Splash.css';

// How long the fan-game notice stays on screen (L-01).
const SPLASH_MS = 5000;

interface Props {
  onDone: () => void;
}

export function Splash({ onDone }: Props) {
  const [notice, setNotice] = useState<NoticeLine[]>([]);

  useEffect(() => {
    getNotice().then(setNotice, () => {
      setNotice([]);
    });
    const timer = window.setTimeout(onDone, SPLASH_MS);
    return () => {
      window.clearTimeout(timer);
    };
  }, [onDone]);

  return (
    <div className="splash">
      <h1>Four Souls</h1>
      <div className="splash-notice">
        {notice.map((line, i) => (
          <p key={i}>
            {line.map((part, j) =>
              part.url ? (
                <a
                  key={j}
                  href={part.url}
                  onClick={(e) => {
                    e.preventDefault();
                    if (part.url) {
                      openExternal(part.url);
                    }
                  }}
                >
                  {part.text}
                </a>
              ) : (
                <span key={j}>{part.text}</span>
              ),
            )}
          </p>
        ))}
      </div>
    </div>
  );
}
