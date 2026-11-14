import { useEffect, useState } from 'react';
import {
  checkUpdate,
  downloadUpdate,
  getVersions,
  installUpdate,
  onUpdateProgress,
  openReleases,
  type UpdateProgress,
  type UpdateStatus,
} from '../bridge';
import { t } from '../i18n';
import { messageOf, useApp } from '../store';
import { manualReason, percentOf, progressText, reasonText } from '../update';
import './screens.css';

// phase is where the update is: a check, a download, then the install.
type Phase = 'idle' | 'checking' | 'checked' | 'downloading' | 'downloaded' | 'installing';

// UpdateScreen installs a newer release (18.24, ADR 029, UC-43). It
// checks only when the player asks (owner, 2026-11-14).
export function UpdateScreen() {
  const show = useApp((s) => s.show);
  const [current, setCurrent] = useState('');
  const [phase, setPhase] = useState<Phase>('idle');
  const [status, setStatus] = useState<UpdateStatus | null>(null);
  const [progress, setProgress] = useState<UpdateProgress | null>(null);
  const [error, setError] = useState('');
  const busy = phase === 'checking' || phase === 'downloading' || phase === 'installing';

  useEffect(() => {
    getVersions().then(
      (v) => {
        setCurrent(v.app);
      },
      () => {
        setCurrent('');
      },
    );
  }, []);

  useEffect(() => onUpdateProgress(setProgress), []);

  const check = async () => {
    setPhase('checking');
    setError('');
    setStatus(null);
    setProgress(null);
    try {
      setStatus(await checkUpdate());
      setPhase('checked');
    } catch (e) {
      setError(messageOf(e));
      setPhase('idle');
    }
  };

  const download = async () => {
    setPhase('downloading');
    setError('');
    setProgress(null);
    try {
      await downloadUpdate();
      setPhase('downloaded');
    } catch (e) {
      setError(messageOf(e));
      setProgress(null);
      setPhase('checked');
    }
  };

  const install = async () => {
    setPhase('installing');
    setError('');
    try {
      await installUpdate();
    } catch (e) {
      setError(messageOf(e));
      setPhase('downloaded');
    }
  };

  const version = status ? status.current : current;
  const loading = phase === 'downloading' || phase === 'downloaded' || phase === 'installing';
  const manual = error !== '' || (status !== null && manualReason(status.reason));
  return (
    <div className="screen">
      <h1>{t('update.title')}</h1>
      <div className="panel update-panel">
        {version && <p className="update-version">{t('update.current', { version })}</p>}
        {status && <p className="update-version">{t('update.latest', { version: status.latest })}</p>}
        {status && !status.canInstall && <p className="update-reason">{reasonText(status.reason)}</p>}
        {status?.notes && (
          <>
            <h2 className="update-notes-title" id="update-notes">
              {t('update.notes')}
            </h2>
            <div className="update-notes" role="region" aria-labelledby="update-notes">
              {status.notes}
            </div>
          </>
        )}
        {loading && progress && <ProgressBar progress={progress} />}
        <div className="row update-buttons">
          <button
            type="button"
            className={status ? undefined : 'primary'}
            disabled={busy}
            onClick={() => {
              void check();
            }}
          >
            {t('update.check')}
          </button>
          {status?.canInstall && phase !== 'downloaded' && phase !== 'installing' && (
            <button
              type="button"
              className="primary"
              disabled={busy}
              onClick={() => {
                void download();
              }}
            >
              {t('update.download')}
            </button>
          )}
          {(phase === 'downloaded' || phase === 'installing') && (
            <button
              type="button"
              className="primary"
              disabled={busy}
              onClick={() => {
                void install();
              }}
            >
              {t('update.install')}
            </button>
          )}
          {manual && (
            <button type="button" onClick={openReleases}>
              {t('update.releases')}
            </button>
          )}
        </div>
      </div>
      <p className="status" role="status">
        {error || statusLine(phase, status)}
      </p>
      <button
        type="button"
        onClick={() => {
          show('menu');
        }}
      >
        {t('update.back')}
      </button>
    </div>
  );
}

// statusLine says what the screen is doing, or what comes next.
function statusLine(phase: Phase, status: UpdateStatus | null): string {
  switch (phase) {
    case 'checking':
      return t('update.checking');
    case 'downloading':
      return t('update.downloading');
    case 'downloaded':
      return t('update.downloaded');
    case 'installing':
      return t('update.installing');
    case 'checked':
      return status?.canInstall ? t('update.ready') : '';
    case 'idle':
      return '';
  }
}

// ProgressBar is the download's bar, with percent and MB under it; an
// unknown total shows MB only and an empty bar.
function ProgressBar({ progress }: { progress: UpdateProgress }) {
  const percent = percentOf(progress);
  const text = progressText(progress);
  return (
    <div className="update-progress">
      <div
        className="update-bar"
        role="progressbar"
        aria-label={t('update.progress_label')}
        aria-valuemin={0}
        aria-valuemax={100}
        aria-valuenow={percent ?? undefined}
        aria-valuetext={text}
      >
        <div className="update-fill" style={{ width: `${String(percent ?? 0)}%` }} />
      </div>
      <span className="muted">{text}</span>
    </div>
  );
}
