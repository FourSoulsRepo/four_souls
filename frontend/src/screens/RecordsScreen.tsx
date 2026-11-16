import { useCallback, useEffect, useState } from 'react';
import { deleteRecord, exportRecord, listRecords, openRecord, pickRecordFile, type RecordEntry } from '../bridge';
import { t } from '../i18n';
import { useEscape } from '../keys';
import { entryLine, recordTitle } from '../records';
import { messageOf, useApp } from '../store';
import './screens.css';

// RecordsScreen manages saved matches (M-01, UC-44): the records
// folder's list, Watch, Export, Delete, and Open a file. Watch puts the
// record on the table (10.2, UC-45).
export function RecordsScreen() {
  const show = useApp((s) => s.show);
  const watchRecord = useApp((s) => s.watch);
  const [records, setRecords] = useState<RecordEntry[] | null>(null);
  const [asking, setAsking] = useState('');
  const [busy, setBusy] = useState(false);
  const [status, setStatus] = useState('');

  const load = useCallback(
    () =>
      listRecords().then(setRecords, (e: unknown) => {
        setRecords([]);
        setStatus(messageOf(e));
      }),
    [],
  );

  useEffect(() => {
    void load();
  }, [load]);

  useEscape(() => {
    show('menu');
  });

  // run does one action at a time and shows its error in the status.
  const run = async (action: () => Promise<void>) => {
    setBusy(true);
    setStatus('');
    try {
      await action();
    } catch (e) {
      setStatus(messageOf(e));
    } finally {
      setBusy(false);
    }
  };

  const watch = (path: string) =>
    run(async () => {
      await watchRecord(await openRecord(path));
    });

  const openFile = () =>
    run(async () => {
      const path = await pickRecordFile();
      if (path) {
        await watchRecord(await openRecord(path));
      }
    });

  const exportOne = (path: string) =>
    run(async () => {
      const to = await exportRecord(path);
      if (to) {
        setStatus(t('records.exported', { path: to }));
      }
    });

  const remove = (path: string) =>
    run(async () => {
      setAsking('');
      await deleteRecord(path);
      await load();
      setStatus(t('records.deleted'));
    });

  return (
    <div className="screen">
      <h1>{t('records.heading')}</h1>
      <section className="panel records-panel">
        {records === null && <span className="muted">{t('records.loading')}</span>}
        {records?.length === 0 && <span className="muted">{t('records.none')}</span>}
        {records && records.length > 0 && (
          <ul className="list records-list" aria-label={t('records.list_label')}>
            {records.map((r) => {
              const name = recordTitle(r);
              return (
                <li key={r.path}>
                  <div className="records-entry">
                    <span className="selectable">{name}</span>
                    <span className="muted">{entryLine(r)}</span>
                  </div>
                  {asking === r.path ? (
                    <div className="row records-buttons" role="group" aria-label={t('records.delete_ask')}>
                      <span>{t('records.delete_ask')}</span>
                      <button
                        type="button"
                        disabled={busy}
                        onClick={() => {
                          void remove(r.path);
                        }}
                      >
                        {t('records.delete_yes')}
                      </button>
                      <button
                        type="button"
                        className="primary"
                        onClick={() => {
                          setAsking('');
                        }}
                      >
                        {t('records.delete_no')}
                      </button>
                    </div>
                  ) : (
                    <div className="row records-buttons">
                      {r.why === '' && (
                        <button
                          type="button"
                          className="primary"
                          disabled={busy}
                          aria-label={t('records.watch_label', { name })}
                          onClick={() => {
                            void watch(r.path);
                          }}
                        >
                          {t('records.watch')}
                        </button>
                      )}
                      <button
                        type="button"
                        disabled={busy}
                        aria-label={t('records.export_label', { name })}
                        onClick={() => {
                          void exportOne(r.path);
                        }}
                      >
                        {t('records.export')}
                      </button>
                      <button
                        type="button"
                        disabled={busy}
                        aria-label={t('records.delete_label', { name })}
                        onClick={() => {
                          setAsking(r.path);
                        }}
                      >
                        {t('records.delete')}
                      </button>
                    </div>
                  )}
                </li>
              );
            })}
          </ul>
        )}
      </section>
      <p className="status" role="status">
        {status}
      </p>
      <div className="row">
        <button
          type="button"
          disabled={busy}
          onClick={() => {
            void openFile();
          }}
        >
          {t('records.open_file')}
        </button>
        <button
          type="button"
          onClick={() => {
            show('menu');
          }}
        >
          {t('records.back')}
        </button>
      </div>
    </div>
  );
}
