import { Playbook, ChatRoom, ChatUser } from './types';

// --- CHAT MOCK DATA ---

const USERS: Record<string, ChatUser> = {
  me: { id: 'me', name: 'Alikhan', role: 'student', isOnline: true },
  advisor1: { id: 'adv1', name: 'Prof. Ivanov', role: 'advisor', isOnline: true },
  advisor2: { id: 'adv2', name: 'Dr. Smith', role: 'advisor', isOnline: false },
  secretary: { id: 'sec', name: 'Elena Petrovna', role: 'secretary', isOnline: true },
  admin: { id: 'adm', name: 'OmiD Support', role: 'admin', isOnline: false },
};

export const MOCK_CHAT_ROOMS: ChatRoom[] = [
  {
    id: 'room-1',
    name: 'Scientific Supervisors',
    type: 'group',
    participants: [USERS.me, USERS.advisor1, USERS.advisor2],
    unreadCount: 2,
    lastMessage: {
      id: 'm1',
      senderId: 'adv1',
      content: 'Please review the comments on Chapter 2 by Friday.',
      timestamp: new Date(Date.now() - 1000 * 60 * 15).toISOString(), // 15 mins ago
      status: 'read'
    },
    messages: [
      {
        id: 'm-prev',
        senderId: 'me',
        content: 'I have uploaded the latest draft of the Methodology section.',
        timestamp: new Date(Date.now() - 1000 * 60 * 60 * 2).toISOString(),
        status: 'read',
        attachments: [{ id: 'a1', type: 'file', name: 'Methodology_v3.docx', url: '#', size: '2.4 MB' }]
      },
      {
        id: 'm1',
        senderId: 'adv1',
        content: 'Please review the comments on Chapter 2 by Friday.',
        timestamp: new Date(Date.now() - 1000 * 60 * 15).toISOString(),
        status: 'read'
      }
    ]
  },
  {
    id: 'room-2',
    name: 'Elena Petrovna (Secretary)',
    type: 'private',
    participants: [USERS.me, USERS.secretary],
    unreadCount: 0,
    lastMessage: {
      id: 'm2',
      senderId: 'me',
      content: 'Thank you, I will bring the printed copies tomorrow.',
      timestamp: new Date(Date.now() - 1000 * 60 * 60 * 24).toISOString(), // 1 day ago
      status: 'read'
    },
    messages: [
      {
        id: 'm2-1',
        senderId: 'sec',
        content: 'Alikhan, do not forget the NCSTE certificate.',
        timestamp: new Date(Date.now() - 1000 * 60 * 60 * 25).toISOString(),
        status: 'read'
      },
      {
        id: 'm2',
        senderId: 'me',
        content: 'Thank you, I will bring the printed copies tomorrow.',
        timestamp: new Date(Date.now() - 1000 * 60 * 60 * 24).toISOString(),
        status: 'read'
      }
    ]
  },
  {
    id: 'room-3',
    name: 'PhD Announcements',
    type: 'channel',
    participants: [USERS.admin],
    unreadCount: 5,
    lastMessage: {
      id: 'm3',
      senderId: 'adm',
      content: 'REMINDER: Deadline for preliminary defense applications is Oct 15.',
      timestamp: new Date(Date.now() - 1000 * 60 * 60 * 48).toISOString(), // 2 days ago
      status: 'read'
    },
    messages: [
      {
        id: 'm3',
        senderId: 'adm',
        content: 'REMINDER: Deadline for preliminary defense applications is Oct 15.',
        timestamp: new Date(Date.now() - 1000 * 60 * 60 * 48).toISOString(),
        status: 'read'
      }
    ]
  }
];


// The full JSON provided by the user
export const PHD_PLAYBOOK: Playbook = {
  "playbook_id": "phd-doctorant.kz",
  "version": "1.1.0",
  "ui": {
    "worlds_palette": [
      "#0ea5e9",
      "#22c55e",
      "#f59e0b",
      "#ef4444",
      "#a855f7",
      "#14b8a6",
      "#64748b"
    ],
    "icons": {
      "form": "lucide:form-input",
      "upload": "lucide:upload",
      "decision": "lucide:git-merge",
      "meeting": "lucide:users",
      "waiting": "lucide:hourglass",
      "external": "lucide:external-link",
      "boss": "lucide:trophy"
    }
  },
  "worlds": [
    {
      "id": "W1",
      "title": { "ru": "I — Подготовка", "kz": "I — Дайындық", "en": "I — Preparation" },
      "order": 1,
      "nodes": [
        { "id": "S1_profile", "title": { "ru": "Профиль докторанта", "kz": "Докторант профилі", "en": "Doctoral profile" }, "type": "form", "state": "done" },
        { "id": "S1_text_ready", "title": { "ru": "Текст диссертации подготовлен", "kz": "Диссертация мәтіні дайын", "en": "Dissertation draft ready" }, "type": "confirmTask", "state": "done" },
        { "id": "S0_antiplagiat", "module": "I", "title": { "ru": "Справка на антиплагиат", "kz": "Антиплагиат анықтамасы", "en": "Anti-plagiarism certificate" }, "type": "confirmTask", "state": "done" },
        { "id": "S1_publications_list", "module": "I", "title": { "ru": "Список публикаций", "kz": "Жарияланымдар тізімі", "en": "Publications List" }, "type": "form", "description": { "ru": "Заполните данные о публикациях...", "kz": "Жарияланымдар туралы деректер...", "en": "Fill in publication data..." }, "state": "done" }
      ],
      "status": "completed",
      "progress": 100
    },
    {
      "id": "W2",
      "title": { "ru": "II — Предварительная экспертиза (НК)", "kz": "II — Алдын ала сараптама (ҒК)", "en": "II — Pre-examination (SC)" },
      "order": 2,
      "nodes": [
        { "id": "E1_apply_omid", "title": { "ru": "Заявка в ОМиД", "kz": "ОМиД-ке өтініш", "en": "Application to OMiD" }, "type": "confirmTask", "state": "done" },
        { "id": "E3_hearing_nk", "title": { "ru": "Заслушивание НК", "kz": "ҒК тыңдауы", "en": "SC Hearing" }, "type": "form", "state": "done" }
      ],
      "status": "completed",
      "progress": 100
    },
    {
      "id": "W3",
      "title": { "ru": "III — RP (условно)", "kz": "III — RP (шартты)", "en": "III — RP (conditional)" },
      "order": 3,
      "nodes": [
        { "id": "RP1_overview_actualization", "module": "III", "title": { "ru": "Актуализация Research Proposal", "kz": "Research Proposal өзектендіру", "en": "Research Proposal actualization" }, "type": "info", "state": "done" },
        { "id": "RP2_sc_hearing_prep", "module": "III", "title": { "ru": "Заслушивание Research Proposal", "kz": "Research Proposal тыңдауы", "en": "SC hearing of RP" }, "type": "form", "state": "done" },
        { "id": "RP3_pre_expertise_application", "module": "III", "title": { "ru": "Заявление в ОМиД", "kz": "ОМиД-ке өтініш", "en": "Application to OMiD" }, "type": "confirmTask", "state": "done" }
      ],
      "status": "completed",
      "progress": 100
    },
    {
      "id": "W4",
      "title": { "ru": "IV — Подача в Диссертационный совет", "kz": "IV — Диссертациялық кеңеске тапсыру", "en": "IV — Submission to DC" },
      "order": 4,
      "nodes": [
        { "id": "D1_normokontrol_ncste", "module": "IV", "title": { "ru": "НЦГНТЭ: нормоконтроль", "kz": "НЦҒНТЭ: нормоконтроль", "en": "NCSTE: norm control" }, "type": "confirmTask", "state": "done" },
        { "id": "IV_rector_application", "module": "IV", "title": { "ru": "Заявление ректору", "kz": "Ректорға өтініш", "en": "Letter to Rector" }, "type": "confirmTask", "state": "done" },
        { "id": "IV3_publication_certificate_ncste", "module": "IV", "title": { "ru": "НЦГНТЭ: справка о публикациях", "kz": "НЦҒНТЭ: жарияланымдар анықтамасы", "en": "NCSTE: publication certificate" }, "type": "confirmTask", "state": "done" },
        { "id": "D2_apply_to_ds", "module": "IV", "title": { "ru": "Пакет документов в ДС", "kz": "ДК құжаттар топтамасы", "en": "DC Document package" }, "type": "form", "state": "done" }
      ],
      "status": "completed",
      "progress": 100
    },
    {
      "id": "W5",
      "title": { "ru": "V — Восстановление", "kz": "V — Қалпына келтіру", "en": "V — Restoration" },
      "order": 5,
      "nodes": [
        { "id": "V1_reinstatement_package", "module": "V", "title": { "ru": "Восстановление на защиту", "kz": "Қорғауға қалпына келтіру", "en": "Reinstatement for defense" }, "type": "form", "state": "done" }
      ],
      "status": "completed",
      "progress": 100
    },
    {
      "id": "W6",
      "title": { "ru": "VI — После принятия ДС", "kz": "VI — ДК қабылдағаннан кейін", "en": "VI — After DC acceptance" },
      "order": 6,
      "nodes": [
        { "id": "A1_post_acceptance_overview", "module": "VI", "title": { "ru": "После принятия документов", "kz": "Құжаттар қабылданғаннан кейін", "en": "After acceptance" }, "type": "info", "state": "done" }
      ],
      "status": "completed",
      "progress": 100
    },
    {
      "id": "W7",
      "title": {
        "ru": "VII — Защита и После защиты",
        "kz": "VII — Қорғау және Қорғаудан кейін",
        "en": "VII — Defense & Post-defense"
      },
      "order": 7,
      "nodes": [
        {
          "id": "VI1_post_defense_overview",
          "module": "VI",
          "title": {
            "ru": "После защиты",
            "kz": "Қорғаудан кейін",
            "en": "After the defense"
          },
          "type": "info",
          "description": {
            "ru": "После успешной защиты начинается административный этап, который выполняет Ученый секретарь.",
            "kz": "Сәтті қорғаудан кейін ғылыми хатшы орындайтын әкімшілік кезең басталады.",
            "en": "After a successful defense, the administrative phase begins."
          },
          "state": "done"
        },
        {
          "id": "VI2_library_deposits",
          "module": "VI",
          "title": {
            "ru": "Сдача печатных экземпляров",
            "kz": "Кітапханаларға тапсыру",
            "en": "Library deposits"
          },
          "type": "form",
          "description": {
            "ru": "Отметьте подготовку и сдачу комплектов в три библиотеки.",
            "kz": "Үш кітапханаға жиынтықтарды дайындау және тапсыруды белгілеңіз.",
            "en": "Mark preparation and delivery of sets to the three libraries."
          },
          "state": "done",
          "requirements": {
            "fields": [
              { "key": "chk_print_5_hardbound", "type": "boolean", "required": true, "label": { "ru": "5 экз. в твёрдом переплёте (прошитые)", "kz": "5 дана қатты түптелген (тігілген)", "en": "5 hardbound copies (sewn)" } },
              { "key": "chk_cd_abstracts_pdf", "type": "boolean", "required": true, "label": { "ru": "CD с аннотациями (3 PDF)", "kz": "Аннотациялармен CD (3 PDF)", "en": "CD with abstracts (3 PDFs)" } },
              { "key": "chk_cd_hard_case", "type": "boolean", "required": true, "label": { "ru": "CD в твёрдом чехле + вкладыш", "kz": "CD қатты қаптамада + вкладыш", "en": "CD in hard case + insert" } },
              { "key": "chk_delivered_nal", "type": "boolean", "required": true, "label": { "ru": "Сдано в Нац. академ. библиотеку", "kz": "Ұлттық академиялық кітапханаға тапсырылды", "en": "Delivered to Nat. Academic Library" } },
              { "key": "chk_delivered_nl", "type": "boolean", "required": true, "label": { "ru": "Сдано в Национальную библиотеку", "kz": "Ұлттық кітапханаға тапсырылды", "en": "Delivered to National Library" } },
              { "key": "chk_delivered_kaznmu", "type": "boolean", "required": true, "label": { "ru": "Сдано в библиотеку КазНМУ", "kz": "ҚазҰМУ кітапханасына тапсырылды", "en": "Delivered to KazNMU Library" } },
              { "key": "chk_receipts", "type": "boolean", "required": true, "label": { "ru": "Справки получены", "kz": "Анықтамалар алынды", "en": "Receipts obtained" } }
            ]
          }
        },
        {
          "id": "VI3_ncste_state_registration",
          "module": "VI",
          "title": {
            "ru": "Государственная регистрация (НЦГНТЭ)",
            "kz": "Мемлекеттік тіркеу (ҰҒТАО)",
            "en": "State registration (NCSTE)"
          },
          "type": "form",
          "description": {
            "ru": "Сформируйте пакет и передайте через учёного секретаря.",
            "kz": "Пакетті қалыптастырып, ғылыми хатшы арқылы тапсырыңыз.",
            "en": "Assemble the package and submit via the DC Secretary."
          },
          "state": "done",
          "requirements": {
            "fields": [
              { "key": "grp_print", "type": "note", "label": { "ru": "Печатные материалы", "kz": "Баспадағы материалдар", "en": "Printed materials" } },
              { "key": "chk_unbound", "type": "boolean", "required": true, "label": { "ru": "Диссертация непереплётная + подпись", "kz": "Тігілмеген диссертация + қол", "en": "Unbound dissertation + signature" } },
              { "key": "chk_print_abstracts", "type": "boolean", "required": true, "label": { "ru": "Аннотации (по 1 экз.)", "kz": "Аннотациялар (1 данадан)", "en": "Abstracts (1 copy each)" } },
              { "key": "chk_ukd_dek", "type": "boolean", "required": true, "label": { "ru": "УКД и ДЕК (2 экз., плотная бумага)", "kz": "УКД және ДЕК (2 дана, тығыз қағаз)", "en": "UKD & DEK (2 copies, heavy paper)" } },
              { "key": "chk_id_copy", "type": "boolean", "required": true, "label": { "ru": "Копия удостоверения", "kz": "Жеке куәлік көшірмесі", "en": "ID copy" } },
              { "key": "chk_pub_list", "type": "boolean", "required": true, "label": { "ru": "Заверенный список публикаций", "kz": "Расталған жарияланымдар тізімі", "en": "Certified publication list" } },
              { "key": "grp_cd", "type": "note", "label": { "ru": "Электронные (CD)", "kz": "Электрондық (CD)", "en": "Electronic (CD)" } },
              { "key": "chk_cd_thesis", "type": "boolean", "required": true, "label": { "ru": "Диссертация (Word)", "kz": "Диссертация (Word)", "en": "Dissertation (Word)" } },
              { "key": "chk_cd_abstracts", "type": "boolean", "required": true, "label": { "ru": "Аннотации (Word)", "kz": "Аннотациялар (Word)", "en": "Abstracts (Word)" } },
              { "key": "chk_cd_cards", "type": "boolean", "required": true, "label": { "ru": "УКД и ДЕК (Word)", "kz": "УКД және ДЕК (Word)", "en": "UKD & DEK (Word)" } },
              { "key": "chk_pub_list_pdf", "type": "boolean", "required": true, "label": { "ru": "Заверенный список публикаций (PDF) на CD", "kz": "Расталған жарияланымдар тізімі (PDF) — CD-да", "en": "Certified publications list (PDF) on CD" } },
              { "key": "chk_submitted", "type": "boolean", "required": true, "label": { "ru": "Передано секретарю", "kz": "Хатшыға тапсырылды", "en": "Submitted to Secretary" } }
            ]
          }
        },
        {
          "id": "VI_attestation_file",
          "module": "VI",
          "title": {
            "ru": "Аттестационное дело",
            "kz": "Аттестаттау ісі",
            "en": "Attestation file"
          },
          "type": "form",
          "description": {
            "ru": "Проверьте наличие всех документов для аттестационного дела (31 пункт).",
            "kz": "Аттестаттау ісі үшін барлық құжаттардың болуын тексеріңіз.",
            "en": "Check all documents for the attestation file (31 items)."
          },
          "state": "active",
          "requirements": {
            "fields": [
              { "key": "chk_inventory", "type": "boolean", "required": true, "label": { "ru": "1) Опись документов", "kz": "1) Құжаттар тізімі", "en": "1) Document inventory" } },
              { "key": "chk_app_rector", "type": "boolean", "required": true, "label": { "ru": "2) Заявление ректору", "kz": "2) Ректорға өтініш", "en": "2) Letter to Rector" } },
              { "key": "chk_app_ds", "type": "boolean", "required": true, "label": { "ru": "3) Заявление в ДС", "kz": "3) ДК-ге өтініш", "en": "3) Letter to DC" } },
              { "key": "chk_stud_info", "type": "boolean", "required": true, "label": { "ru": "5) Сведения о докторанте", "kz": "5) Докторант мәліметтері", "en": "5) Student info" } },
              { "key": "chk_personal_sheet", "type": "boolean", "required": true, "label": { "ru": "6) Личный листок (кадры)", "kz": "6) Жеке парақ", "en": "6) Personal sheet" } },
              { "key": "chk_diplomas", "type": "boolean", "required": true, "label": { "ru": "7) Копии дипломов (нотариус)", "kz": "7) Дипломдар көшірмесі", "en": "7) Diploma copies" } },
              { "key": "chk_transcript", "type": "boolean", "required": true, "label": { "ru": "8) Транскрипт", "kz": "8) Транскрипт", "en": "8) Transcript" } },
              { "key": "chk_topic_extract", "type": "boolean", "required": true, "label": { "ru": "9) Выписка о теме/консультантах", "kz": "9) Тақырып туралы үзінді", "en": "9) Topic extract" } },
              { "key": "chk_pub_copies", "type": "boolean", "required": true, "label": { "ru": "10) Копии публикаций", "kz": "10) Жарияланымдар көшірмесі", "en": "10) Publication copies" } },
              { "key": "chk_abstracts", "type": "boolean", "required": true, "label": { "ru": "11) Аннотации (3 языка)", "kz": "11) Аннотациялар", "en": "11) Abstracts" } },
              { "key": "chk_sc_extract", "type": "boolean", "required": true, "label": { "ru": "12) Выписка НК (положит.)", "kz": "12) ҒК үзіндісі", "en": "12) SC extract" } },
              { "key": "chk_ethics", "type": "boolean", "required": true, "label": { "ru": "13) Заключение ЛЭК", "kz": "13) ЛЭК қорытындысы", "en": "13) Ethics conclusion" } },
              { "key": "chk_antiplag_ncste", "type": "boolean", "required": true, "label": { "ru": "14) Антиплагиат НЦГНТЭ", "kz": "14) ҰҒТАО антиплагиат", "en": "14) NCSTE anti-plagiarism" } },
              { "key": "chk_reviews_adv", "type": "boolean", "required": true, "label": { "ru": "16) Отзывы консультантов", "kz": "16) Кеңесшілер пікірлері", "en": "16) Advisor reviews" } },
              { "key": "chk_reviews_off", "type": "boolean", "required": true, "label": { "ru": "17) Отзывы рецензентов", "kz": "17) Рецензент пікірлері", "en": "17) Reviewer reports" } },
              { "key": "chk_vote_proto", "type": "boolean", "required": true, "label": { "ru": "18) Протокол счётной комиссии", "kz": "18) Санау комиссиясы хаттамасы", "en": "18) Vote protocol" } },
              { "key": "chk_defense_proto", "type": "boolean", "required": true, "label": { "ru": "21) Протокол заседания ДС", "kz": "21) ДК отырысы хаттамасы", "en": "21) Defense protocol" } },
              { "key": "chk_video", "type": "boolean", "required": true, "label": { "ru": "22) Видеозапись", "kz": "22) Бейнежазба", "en": "22) Video recording" } },
              { "key": "chk_reg_cards", "type": "boolean", "required": true, "label": { "ru": "23/24) Учётные карточки", "kz": "23/24) Есептік карталар", "en": "23/24) Reg cards" } },
              { "key": "chk_libs", "type": "boolean", "required": true, "label": { "ru": "25-27) Справки библиотек", "kz": "25-27) Кітапхана анықтамалары", "en": "25-27) Library receipts" } }
            ]
          }
        }
      ],
      "status": "active",
      "progress": 75
    }
  ]
} as unknown as Playbook;