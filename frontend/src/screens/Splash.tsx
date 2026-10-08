import { useEffect, useState } from 'react';
import { getNotice, openExternal, type NoticeLine, type NoticePart } from '../bridge';
import './Splash.css';

// How long the fan-game notice stays on screen (L-01).
const SPLASH_MS = 5000;

// The notice is static text, so its content makes a stable key.
function partKey(part: NoticePart): string {
  return `${part.text}|${part.url ?? ''}`;
}

function lineKey(line: NoticeLine): string {
  return line.map(partKey).join('');
}

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
        {notice.map((line) => (
          <p key={lineKey(line)}>
            {line.map((part) =>
              part.url ? (
                <a
                  key={partKey(part)}
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
                <span key={partKey(part)}>{part.text}</span>
              ),
            )}
          </p>
        ))}
      </div>
    </div>
  );
}
