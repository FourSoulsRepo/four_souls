// The rules view (12.2): every rules file on one page, with a topic
// index. A rule route scrolls to its rule and highlights it.
import { useEffect } from 'react';
import { t } from '../i18n';
import { RulesText } from './RulesText';
import { rulesHash } from './route';
import { hasRule, topics } from './rules';

export function RulesView({ lang, rule }: { lang: string; rule: string }) {
  useEffect(() => {
    const el = rule === '' ? null : document.getElementById(`rule-${rule}`);
    if (el) {
      el.scrollIntoView({ block: 'start' });
    } else {
      window.scrollTo(0, 0);
    }
  }, [rule]);

  const ruleHref = (id: string) => rulesHash(lang, id);

  return (
    <div className="rules">
      {rule !== '' && !hasRule(rule) && <p className="note warn">{t('site.rules.unknown', { rule })}</p>}
      <h1>{t('site.rules.title')}</h1>
      <p className="intro">{t('site.rules.intro')}</p>
      <nav className="topics" aria-label={t('site.rules.topics')}>
        <ul>
          {topics.map((tp) => (
            <li key={tp.file}>{tp.rules[0] ? <a href={ruleHref(tp.rules[0])}>{tp.title}</a> : tp.title}</li>
          ))}
        </ul>
      </nav>
      {topics.map((tp) => (
        <section key={tp.file} className="topic">
          <h2>{tp.title}</h2>
          <RulesText blocks={tp.blocks} ruleHref={ruleHref} target={rule} />
        </section>
      ))}
    </div>
  );
}
