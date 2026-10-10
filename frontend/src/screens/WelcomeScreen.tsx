import { useState } from 'react';
import { useApp } from '../store';
import './screens.css';

// WelcomeScreen asks for a nickname on the first start (ST-06).
export function WelcomeScreen() {
  const saveNickname = useApp((s) => s.saveNickname);
  const show = useApp((s) => s.show);
  const error = useApp((s) => s.error);
  const [nickname, setNickname] = useState('');
  return (
    <div className="screen">
      <h1>Welcome!</h1>
      <p className="muted">Choose the name other players will see.</p>
      <form
        className="row"
        onSubmit={(e) => {
          e.preventDefault();
          void saveNickname(nickname).then((ok) => {
            if (ok) {
              show('menu');
            }
          });
        }}
      >
        <label className="field">
          Nickname
          <input
            value={nickname}
            onChange={(e) => {
              setNickname(e.target.value);
            }}
          />
        </label>
        <button type="submit" className="primary">
          Continue
        </button>
      </form>
      <p className="status" role="status">
        {error}
      </p>
    </div>
  );
}
