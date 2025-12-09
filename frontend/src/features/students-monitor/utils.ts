export function stageLabel(stage?: string, lang = 'en'): string {
  const en: Record<string, string> = {
    W1: 'I — Preparation',
    W2: 'II — Pre-examination',
    W3: 'III — RP',
    W4: 'IV — Submission to DC',
    W5: 'V — Restoration',
    W6: 'VI — After DC acceptance',
    W7: 'VII — Defense & Post-defense',
  };
  const ru: Record<string, string> = {
    W1: 'I — Подготовка',
    W2: 'II — Предварительная экспертиза',
    W3: 'III — RP',
    W4: 'IV — Подача в ДС',
    W5: 'V — Восстановление',
    W6: 'VI — После принятия ДС',
    W7: 'VII — Защита и После защиты',
  };
  const kz: Record<string, string> = {
    W1: 'I — Дайындық',
    W2: 'II — Алдын ала сараптама',
    W3: 'III — RP',
    W4: 'IV — ДК-ға тапсыру',
    W5: 'V — Қалпына келтіру',
    W6: 'VI — ДК қабылдағаннан кейін',
    W7: 'VII — Қорғау және одан кейін',
  };
  const map = lang.startsWith('ru') ? ru : lang.startsWith('kz') ? kz : en;
  if (!stage) return '—';
  return map[stage] || stage;
}
