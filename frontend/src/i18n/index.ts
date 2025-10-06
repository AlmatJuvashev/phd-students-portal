import i18n from "i18next";
import { initReactI18next } from "react-i18next";
import LanguageDetector from "i18next-browser-languagedetector";

import ru from "@/locales/ru/common.json";
import kz from "@/locales/kz/common.json";
import en from "@/locales/en/common.json";

i18n
  .use(LanguageDetector)
  .use(initReactI18next)
  .init({
    fallbackLng: "ru",
    supportedLngs: ["ru", "kz", "en"],
    resources: {
      ru: { common: ru },
      kz: { common: kz },
      en: { common: en },
    },
    interpolation: { escapeValue: false },
    detection: {
      order: ["querystring", "localStorage", "navigator", "htmlTag"],
      caches: ["localStorage"],
    },
    defaultNS: "common",
    ns: ["common"],
  });

export default i18n;

