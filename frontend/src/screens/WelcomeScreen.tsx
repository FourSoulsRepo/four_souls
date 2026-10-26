import { useState } from 'react';
import { t } from '../i18n';
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
      <h1>{t('welcome.title')}</h1>
      <p className="muted">{t('welcome.hint')}</p>
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
          {t('lobby.nickname')}
          <input
            value={nickname}
            onChange={(e) => {
              setNickname(e.target.value);
            }}
          />
        </label>
        <button type="submit" className="primary">
          {t('welcome.continue')}
        </button>
      </form>
      <p className="status" role="status">
        {error}
      </p>
    </div>
  );
}
