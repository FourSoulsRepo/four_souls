// The Update screen's words and numbers (18.24, ADR 029).
import type { UpdateProgress, UpdateReason } from './bridge';
import { t } from './i18n';

// megabytes is a byte count in MB (2^20 bytes) with one decimal.
export function megabytes(bytes: number): string {
  return (Math.max(bytes, 0) / (1 << 20)).toFixed(1);
}

// percentOf is how much of the download is done, 0 to 100; null when
// the total is unknown.
export function percentOf(p: UpdateProgress): number | null {
  if (p.total <= 0) {
    return null;
  }
  return Math.min(100, Math.max(0, Math.floor((p.done * 100) / p.total)));
}

// progressText is the line under the download bar: percent and MB, or
// MB only when the total is unknown.
export function progressText(p: UpdateProgress): string {
  const percent = percentOf(p);
  if (percent === null) {
    return t('update.progress_mb', { done: megabytes(p.done) });
  }
  return t('update.progress', { percent, done: megabytes(p.done), total: megabytes(p.total) });
}

// reasonText says in plain words why the app can't update itself.
export function reasonText(reason: UpdateReason): string {
  switch (reason) {
    case 'not_version':
      return t('update.reason.not_version');
    case 'not_newer':
      return t('update.reason.not_newer');
    case 'no_archive':
      return t('update.reason.no_archive');
    case 'not_writable':
      return t('update.reason.not_writable');
    case '':
      return '';
    case 'other':
      return t('update.reason.unknown');
  }
}

// manualReason: the player may update by hand from the releases page.
// Already on the latest version, there is nothing to get.
export function manualReason(reason: UpdateReason): boolean {
  return reason !== '' && reason !== 'not_newer';
}
