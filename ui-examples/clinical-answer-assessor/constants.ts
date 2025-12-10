
import { AnswerRecord, ModelConfig, Protocol, Topic, PendingAnswer } from './types';

export const MOCK_ANSWERS: AnswerRecord[] = [
  {
    user_id: 101,
    name_rus: "Акушерство и гинекология",
    topic: "Токсикоз у беременных",
    question: "Полный клинический разбор случая: Токсикоз беременных. Заполните все разделы в соответствии с клиническим протоколом.",
    answer: "Пациентку необходимо госпитализировать. Провести инфузионную терапию...",
    comment: "Хорошо, но не указаны дозировки.",
    score: 72,
    ai_score: 75,
    examiner_id: 456,
    examiner_name: "Иванова И.И.",
    attempt_id: "2024-residency-001",
    created_at: "2024-02-10",
    strengths: [
      "Правильно определена необходимость госпитализации.",
      "Верно указаны основные группы препаратов для инфузионной терапии."
    ],
    weaknesses: [
      "Отсутствует классификация степени тяжести.",
      "Не указаны конкретные дозировки и скорость инфузии.",
      "Не описана тактика при неэффективности консервативной терапии."
    ],
    section_scores: [
      { id: 'r1', score: 8, max: 10, feedback: "Определение сформулировано верно." },
      { id: 'r2', score: 5, max: 10, feedback: "Классификация приведена не полностью." },
      { id: 'r3', score: 9, max: 10, feedback: "Симптомы описаны подробно." },
      { id: 'r7', score: 12, max: 20, feedback: "Лечение описано поверхностно, нет дозировок." },
    ]
  },
  {
    user_id: 102,
    name_rus: "Акушерство и гинекология",
    topic: "Токсикоз у беременных",
    question: "Опишите тактику ведения пациентки с тяжелым токсикозом...",
    answer: "Срочное родоразрешение вне зависимости от срока.",
    comment: "Неверно. Сначала стабилизация.",
    score: 30,
    ai_score: 35,
    examiner_id: 456,
    examiner_name: "Иванова И.И.",
    attempt_id: "2024-residency-002",
    created_at: "2024-02-11"
  },
  {
    user_id: 103,
    name_rus: "Кардиология",
    topic: "Инфаркт миокарда",
    question: "Алгоритм действий при ОКС с подъемом ST",
    answer: "Аспирин, клопидогрел, гепарин, вызов скорой...",
    comment: "Полный ответ.",
    score: 95,
    ai_score: 92,
    examiner_id: 789,
    examiner_name: "Петров П.П.",
    attempt_id: "2024-residency-003",
    created_at: "2024-02-12"
  },
  {
    user_id: 104,
    name_rus: "Кардиология",
    topic: "Инфаркт миокарда",
    question: "Алгоритм действий при ОКС с подъемом ST",
    answer: "Дать нитроглицерин и ждать.",
    comment: "Критически недостаточно.",
    score: 20,
    ai_score: 25,
    examiner_id: 789,
    examiner_name: "Петров П.П.",
    attempt_id: "2024-residency-004",
    created_at: "2024-02-12"
  },
  {
    user_id: 105,
    name_rus: "Акушерство и гинекология",
    topic: "HELLP-синдром",
    question: "Диагностические критерии HELLP",
    answer: "Гемолиз, повышение печеночных ферментов, тромбоцитопения.",
    comment: "Верно.",
    score: 88,
    ai_score: 90,
    examiner_id: 111,
    examiner_name: "Сидорова С.С.",
    attempt_id: "2024-residency-005",
    created_at: "2024-02-13"
  }
];

// Mock for Grading Queue
export const MOCK_PENDING_ANSWERS: PendingAnswer[] = [
    {
        id: 'p1',
        student_id: 9942,
        category: 'Оториноларингология',
        question: 'Острый средний отит. Этиология и патогенез.',
        answer: 'Воспаление слизистой оболочки барабанной полости. Возбудители: пневмококк, гемофильная палочка...',
        submitted_at: '2024-02-25T10:30:00'
    },
    {
        id: 'p2',
        student_id: 5521,
        category: 'Кардиология',
        question: 'Дифференциальная диагностика болей в грудной клетке.',
        answer: 'Необходимо исключить ОКС, ТЭЛА, расслоение аорты. При ОКС боль давящая, иррадиирует в левую руку...',
        submitted_at: '2024-02-25T11:15:00'
    },
    {
        id: 'p3',
        student_id: 1234,
        category: 'Терапия',
        question: 'Лечение железодефицитной анемии.',
        answer: 'Препараты железа перорально (сорбифер, мальтофер). Контроль ферритина через 3 месяца.',
        submitted_at: '2024-02-26T09:00:00'
    }
];

// Generate more mock data for stats
for (let i = 0; i < 20; i++) {
  MOCK_ANSWERS.push({
    user_id: 200 + i,
    name_rus: "Терапия",
    topic: "Пневмония",
    question: "Антибиотикотерапия при внебольничной пневмонии",
    answer: "Амоксициллин/клавуланат...",
    comment: "Ок",
    score: 60 + Math.floor(Math.random() * 40), // Random score 60-100
    ai_score: 65 + Math.floor(Math.random() * 30), // Random AI score
    examiner_id: 111, // Sidorova is lenient
    examiner_name: "Сидорова С.С.",
    attempt_id: `2024-gen-${i}`,
    created_at: "2024-02-14"
  });
}

export const PROTOCOLS: Protocol[] = [
  { id: '1', name: 'NRCHD: Токсикоз беременных (2023)', type: 'PDF', language: 'RU', status: 'indexed', last_updated: '2023-10-15', active: true },
  { id: '2', name: 'NRCHD: Артериальная гипертензия', type: 'PDF', language: 'RU', status: 'indexed', last_updated: '2022-05-20', active: true },
  { id: '3', name: 'AHA/ACC: STEMI Guidelines', type: 'Text', language: 'KZ', status: 'not_indexed', last_updated: '2023-01-10', active: false },
  { id: '4', name: 'MOH: Сахарный диабет 2 типа', type: 'PDF', language: 'RU', status: 'indexed', last_updated: '2024-01-05', active: true },
];

export const LLM_MODELS: ModelConfig[] = [
  { id: 'qwen-2-72b', name: 'Qwen 2 72B', tags: ['High Accuracy', 'Multilingual'] },
  { id: 'qwen-2-32b', name: 'Qwen 2 32B', tags: ['Fast', 'Low Cost'] },
  { id: 'gpt-4o', name: 'OpenAI GPT-4o', tags: ['Benchmark Leader', 'Expensive'] },
  { id: 'deepseek-r1', name: 'DeepSeek R1', tags: ['Reasoning Specialist'] },
];

export const TOPICS: Topic[] = [
  {
    id: 't1',
    title: 'Токсикоз у беременных',
    description: 'Диагностика и лечение рвоты беременных различной степени тяжести.',
    question_count: 12,
    objectives: ['Оценка степени тяжести', 'Инфузионная терапия', 'Показания к прерыванию'],
    questions: [
        { 
          id: 'q1', 
          text: 'Полный клинический разбор случая: Токсикоз беременных. Заполните все разделы в соответствии с клиническим протоколом.', 
          difficulty: 'Advanced', 
          estimated_time_mins: 45,
          rubric: [
            { id: 'r1', criteria: 'Определение заболевания', description: 'Сформулировано определение заболевания, приведены основные этиологические факторы и описаны механизмы патогенеза.', max_score: 10 },
            { id: 'r2', criteria: 'Классификация заболевания', description: 'Приведены актуальные классификации заболевания.', max_score: 10 },
            { id: 'r3', criteria: 'Основные симптомы и синдромы', description: 'Описаны основные симптомы и синдромы, характерные для данного заболевания.', max_score: 10 },
            { id: 'r4', criteria: 'Диагностические критерии', description: 'Установлены диагностические критерии, основанные на клинических проявлениях и анамнезе пациента.', max_score: 10 },
            { id: 'r5', criteria: 'Лабораторные и инструментальные исследования', description: 'Представлен список необходимых лабораторных и инструментальных исследований с описанием ожидаемых результатов и их значимости для диагностики.', max_score: 10 },
            { id: 'r6', criteria: 'Дифференциальная диагностика', description: 'Выполнена дифференциальная диагностика с учетом клинических и лабораторно-инструментальных данных.', max_score: 10 },
            { id: 'r7', criteria: 'Лечение', description: 'Описано лечение, соответствующее современным клиническим протоколам, с указанием кратности и продолжительности применения лекарственных средств.', max_score: 20 },
            { id: 'r8', criteria: 'Тактика ведения пациента', description: 'Определена тактика ведения пациента, включая показания и противопоказания к оперативному вмешательству.', max_score: 10 },
            { id: 'r9', criteria: 'Осложнения, профилактика и прогноз', description: 'Описаны возможные осложнения заболевания, предложены методы их профилактики и дан прогноз.', max_score: 10 },
          ]
        },
        { id: 'q2', text: 'Составьте план инфузионной терапии при токсикозе средней тяжести.', difficulty: 'Intermediate', estimated_time_mins: 15 },
    ]
  },
  {
    id: 't2',
    title: 'Инфаркт миокарда (ОКС)',
    description: 'Алгоритмы ведения пациентов с острым коронарным синдромом.',
    question_count: 25,
    objectives: ['ЭКГ диагностика', 'Тромболизис', 'Антикоагулянты'],
    questions: [
        { id: 'q3', text: 'Опишите ЭКГ признаки переднего распространенного инфаркта миокарда.', difficulty: 'Intermediate', estimated_time_mins: 12 },
        { id: 'q4', text: 'Тактика ведения пациента с кардиогенным шоком.', difficulty: 'Advanced', estimated_time_mins: 20 },
    ]
  },
  {
    id: 't3',
    title: 'Сахарный диабет 2 типа',
    description: 'Современные подходы к гликемическому контролю.',
    question_count: 18,
    objectives: ['Метформин и ингибиторы SGLT2', 'Инсулинотерапия', 'Осложнения'],
    questions: [
        { id: 'q5', text: 'Критерии диагностики сахарного диабета 2 типа.', difficulty: 'Beginner', estimated_time_mins: 8 },
    ]
  }
];
