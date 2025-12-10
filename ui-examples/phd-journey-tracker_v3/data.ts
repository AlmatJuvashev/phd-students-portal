import { Playbook } from './types';

// The full JSON provided by the user
export const PHD_PLAYBOOK: Playbook = {
  "playbook_id": "phd-doctorant.kz",
  "version": "1.1.0",
  // ... (abbreviated properties handled by type matching)
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
      "title": {
        "ru": "I — Подготовка",
        "kz": "I — Дайындық",
        "en": "I — Preparation"
      },
      "order": 1,
      "nodes": [
        {
          "id": "S1_profile",
          "title": {
            "ru": "Профиль докторанта",
            "kz": "Докторант профилі",
            "en": "Doctoral profile"
          },
          "type": "form",
          "state": "done" // Initial state override handled in component
        },
        {
          "id": "S1_text_ready",
          "title": {
            "ru": "Текст диссертации подготовлен",
            "kz": "Диссертация мәтіні дайын",
            "en": "Dissertation draft ready"
          },
          "type": "confirmTask",
          "state": "done"
        },
        {
          "id": "S0_antiplagiat",
          "module": "I",
          "title": {
            "ru": "Справка на антиплагиат",
            "kz": "Антиплагиат анықтамасы",
            "en": "Anti-plagiarism certificate"
          },
          "type": "confirmTask",
          "state": "done"
        },
        {
          "id": "S1_publications_list",
          "module": "I",
          "title": {
            "ru": "Список публикаций",
            "kz": "Жарияланымдар тізімі",
            "en": "Publications List"
          },
          "type": "form",
          "description": {
            "ru": "Заполните данные о публикациях и загрузите подписанный список.",
            "kz": "Жарияланымдар туралы деректерді толтырыңыз және қол қойылған тізімді жүктеңіз.",
            "en": "Fill in publication data and upload the signed list."
          },
          "state": "done"
        }
      ],
      "status": "completed",
      "progress": 100
    },
    {
      "id": "W2",
      "title": {
        "ru": "II — Предварительная экспертиза (НК)",
        "kz": "II — Алдын ала сараптама (ҒК)",
        "en": "II — Pre-examination (Scientific Committee)"
      },
      "order": 2,
      "nodes": [
        {
          "id": "E1_apply_omid",
          "title": {
            "ru": "Заявка в ОМиД",
            "kz": "ОМиД-ке өтініш",
            "en": "Application to OMiD"
          },
          "type": "confirmTask",
          "description": {
            "ru": "Подтвердите подготовку и отправку заявления на предварительную экспертизу в ОМиД.",
            "kz": "ОМиД-ке алдын ала сараптамаға өтініш дайындалып, жіберілгенін растаңыз.",
            "en": "Confirm the OMiD preliminary review application is prepared and submitted."
          },
          "state": "done"
        },
        {
          "id": "E3_hearing_nk",
          "title": {
            "ru": "Заслушивание НК",
            "kz": "ҒК тыңдауы",
            "en": "SC Hearing"
          },
          "type": "form",
          "description": {
            "ru": "Подтвердите результаты после заслушивания.",
            "kz": "Тыңдаудан кейінгі нәтижелерді растаңыз.",
            "en": "Confirm post-hearing results."
          },
          "state": "done"
        }
      ],
      "status": "completed",
      "progress": 100
    },
    {
      "id": "W3",
      "title": {
        "ru": "III — RP (условно)",
        "kz": "III — RP (шартты)",
        "en": "III — RP (conditional)"
      },
      "order": 3,
      "nodes": [
        {
          "id": "RP1_overview_actualization",
          "module": "III",
          "title": {
            "ru": "Актуализация Research Proposal: обзор и условия",
            "kz": "Research Proposal өзектендіру: шолу және шарттар",
            "en": "Research Proposal actualization: overview & conditions"
          },
          "type": "info",
          "description": {
            "ru": "Если со дня выпуска прошло 3 года и более...",
            "kz": "Бітіру күнінен 3 жыл және одан көп өтсе...",
            "en": "If 3+ years have passed since graduation..."
          },
          "state": "done"
        },
        {
          "id": "RP2_sc_hearing_prep",
          "module": "III",
          "title": {
            "ru": "Заслушивание Research Proposal",
            "kz": "Ғылыми комитетте Research Proposal тыңдауы",
            "en": "SC hearing of the Research Proposal"
          },
          "description": {
            "ru": "Подготовьтесь к заслушиванию и загрузите выписку из протокола.",
            "kz": "Тыңдауға дайындалыңыз және хаттамадан үзіндіні жүктеңіз.",
            "en": "Prepare for the hearing and upload the protocol extract."
          },
          "type": "form",
          "state": "done"
        },
        {
          "id": "RP3_pre_expertise_application",
          "module": "III",
          "title": {
            "ru": "Заявление в ОМиД о прохождении предварительной экспертизы",
            "kz": "ОМиД-ке алдын ала сараптамадан өту туралы өтініш",
            "en": "Application to OMiD for preliminary expertise"
          },
          "type": "confirmTask",
          "state": "done"
        }
      ],
      "status": "completed",
      "progress": 100
    },
    {
      "id": "W4",
      "title": {
        "ru": "IV — Подача в Диссертационный совет (ДС)",
        "kz": "IV — Диссертациялық кеңеске тапсыру",
        "en": "IV — Submission to Dissertation Council"
      },
      "order": 4,
      "nodes": [
        {
          "id": "D1_normokontrol_ncste",
          "module": "IV",
          "title": {
            "ru": "НЦГНТЭ: нормоконтроль оформления диссертации",
            "kz": "НЦҒНТЭ: диссертацияны рәсімдеудің нормоконтролі",
            "en": "NCSTE: dissertation formatting check (norm control)"
          },
          "type": "confirmTask",
          "state": "active",
          "requirements": {
            "uploads": [
              {
                "key": "dissertation_docx",
                "label": { "ru": "Файл диссертации (DOCX)", "kz": "Диссертация файлы (DOCX)", "en": "Dissertation manuscript (DOCX)" },
              },
              {
                "key": "ncste_receipt",
                "label": { "ru": "Квитанция об оплате", "kz": "Төлем түбіртегі", "en": "Payment receipt" },
              }
            ]
          }
        },
        {
          "id": "IV_rector_application",
          "module": "IV",
          "title": {
            "ru": "Заявление ректору",
            "kz": "Ректорға өтініш",
            "en": "Letter to Rector"
          },
          "type": "confirmTask",
          "state": "locked"
        },
        {
          "id": "IV3_publication_certificate_ncste",
          "module": "IV",
          "title": {
            "ru": "НЦГНТЭ: справка о публикациях",
            "kz": "НЦҒНТЭ: жарияланымдар туралы анықтама",
            "en": "NCSTE: publication certificate"
          },
          "type": "confirmTask",
          "state": "locked"
        },
        {
          "id": "D2_apply_to_ds",
          "module": "IV",
          "title": {
            "ru": "Пакет документов в Диссертационный совет",
            "kz": "Диссертациялық кеңеске құжаттар топтамасы",
            "en": "Document package for Dissertation Council"
          },
          "description": {
            "ru": "Сбор и передача полного комплекта документов.",
            "kz": "Құжаттар топтамасын жинап тапсыру.",
            "en": "Assemble and submit the full document package."
          },
          "type": "form",
          "state": "locked"
        }
      ],
      "status": "active",
      "progress": 10
    },
    {
      "id": "W5",
      "title": {
        "ru": "V — Восстановление / Дооформление",
        "kz": "V — Қалпына келтіру / Құжаттарды толықтыру",
        "en": "V — Restoration / Completion"
      },
      "order": 5,
      "nodes": [
        {
          "id": "V1_reinstatement_package",
          "module": "V",
          "title": {
            "ru": "Восстановление на защиту",
            "kz": "Қорғауға қалпына келтіру",
            "en": "Reinstatement for defense"
          },
          "type": "form",
          "state": "locked"
        }
      ],
      "status": "locked",
      "progress": 0
    },
    {
      "id": "W6",
      "title": {
        "ru": "VI — После принятия ДС",
        "kz": "VI — ДК қабылдағаннан кейін",
        "en": "VI — After DC acceptance"
      },
      "order": 6,
      "nodes": [
        {
          "id": "A1_post_acceptance_overview",
          "module": "VI",
          "title": {
            "ru": "После принятия документов ДС",
            "kz": "ДС құжаттарын қабылдағаннан кейін",
            "en": "After DC acceptance"
          },
          "type": "info",
          "state": "locked"
        }
      ],
      "status": "locked",
      "progress": 0
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
          "state": "locked"
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
          "state": "locked"
        },
        {
          "id": "VI3_ncste_state_registration",
          "module": "VI",
          "title": {
            "ru": "Государственная регистрация",
            "kz": "Мемлекеттік тіркеу",
            "en": "State registration"
          },
          "type": "form",
          "state": "locked"
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
          "state": "locked"
        }
      ],
      "status": "locked",
      "progress": 0
    }
  ]
} as unknown as Playbook; // Casting because we are skipping some deep raw fields for brevity in this demo