import i18n from "i18next";
import { initReactI18next } from "react-i18next";
import LanguageDetector from "i18next-browser-languagedetector";

import ru from "@/locales/ru/common.json";
import kz from "@/locales/kz/common.json";
import en from "@/locales/en/common.json";
import guides_ru from "@/locales/ru/guides.json";
import guides_kz from "@/locales/kz/guides.json";
import guides_en from "@/locales/en/guides.json";

i18n
  .use(LanguageDetector)
  .use(initReactI18next)
  .init({
    fallbackLng: "ru",
    supportedLngs: ["ru", "kz", "en"],
    resources: {
      ru: { common: ru, guides: guides_ru },
      kz: { common: kz, guides: guides_kz },
      en: { common: en, guides: guides_en },
    },
    interpolation: { escapeValue: false },
    detection: {
      order: ["querystring", "localStorage", "navigator", "htmlTag"],
      caches: ["localStorage"],
    },
    defaultNS: "common",
    ns: ["common", "guides"],
  });

export default i18n;
