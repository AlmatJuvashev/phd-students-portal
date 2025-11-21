/**
 * Appendix 7 generator (template-driven)
 *
 * What this does:
 * - Picks a localized DOCX template (RU/KZ/EN) from assets_list.json.
 * - Maps S1 form values (four collections) into arrays of 6-column rows.
 * - Fills the template using docxtemplater loops and downloads the file.
 *
 * Template requirements (in Word):
 * - For each of the four sections (I–IV) insert a 6-column table:
 *   1) №
 *   2) Наименование работы
 *   3) Печатный или на правах рукописи
 *   4) Наименование издательства журнала (№, стр., год) ISBN или № авторского свидетельства
 *   5) Кол-во печатных листов или стр.
 *   6) Фамилия соавторов
 *
 * - The data row must contain ONE docxtemplater loop spanning the entire row:
 *   Section I  (WoS/Scopus):   {#i_rows}{no} | {title} | {format} | {pub_info} | {pages} | {coauthors}{/i_rows}
 *   Section II (KOKSON):       {#ii_rows}{no} | {title} | {format} | {pub_info} | {pages} | {coauthors}{/ii_rows}
 *   Section III (Conferences): {#iii_rows}{no} | {title} | {format} | {pub_info} | {pages} | {coauthors}{/iii_rows}
 *   Section IV (IP):           {#iv_rows}{no} | {title} | {format} | {pub_info} | {pages} | {coauthors}{/iv_rows}
 *
 * - If a section has no rows, docxtemplater removes the loop row automatically.
 */

import PizZip from "pizzip";
import Docxtemplater from "docxtemplater";
import { saveAs } from "file-saver";
import { getAssetUrl } from "@/lib/assets";
import {
  StudentTemplateData,
  buildSnippet,
  injectSnippet,
} from "./student-template";

// Minimal entry shape coming from S1 form (see playbook.json)
type Entry = {
  title?: string;
  format?: string; // 'print' | 'manuscript' | 'electronic' | 'other'
  format_other?: string;

  journal?: string; // journal / proceedings / venue
  year?: string; // 'YYYY' or 'YYYY/YYYY'
  volume_issue?: string; // '12(3)' or similar
  pages_or_sheets?: string;

  doi?: string;
  issn_print?: string;
  issn_online?: string;

  coauthors?: string[];

  indexing?: string;
  indexing_other?: string;

  // IP-only fields
  ip_type?: string;
  ip_type_other?: string;
  certificate_no?: string;
  isbn?: string;
};

type S1Values = {
  wos_scopus?: Entry[];
  kokson?: Entry[];
  conferences?: Entry[];
  ip?: Entry[];
};

// Localized labels for "format" values across locales
const FORMAT_MAP: Record<"ru" | "kz" | "en", Record<string, string>> = {
  ru: {
    print: "Печатный",
    manuscript: "На правах рукописи",
    electronic: "Электронный",
    other: "Другое",
  },
  kz: {
    print: "Баспа",
    manuscript: "Қолжазба",
    electronic: "Электрондық",
    other: "Басқа",
  },
  en: {
    print: "Printed",
    manuscript: "Manuscript",
    electronic: "Electronic",
    other: "Other",
  },
};

// Localized small tokens for column 4 composition
const PAGES_PREFIX: Record<"ru" | "kz" | "en", string> = {
  ru: "стр.",
  kz: "бет",
  en: "pp.",
};
const CERT_PREFIX: Record<"ru" | "kz" | "en", string> = {
  ru: "Свид. №",
  kz: "Куәлік №",
  en: "Cert. No",
};

/**
 * Map 'format' code to label, or '—' for IP rows where format is not applicable.
 */
function fmtFormat(lang: "ru" | "kz" | "en", e: Entry, isIP = false): string {
  if (isIP) return "—";
  const code = (e.format || "").trim();
  if (code === "other") return e.format_other?.trim() || FORMAT_MAP[lang].other;
  return FORMAT_MAP[lang][code] || "";
}

/**
 * Build Column 4: "publisher/journal (№, pages, year) + ISBN or Certificate"
 * - Journal/proceedings/venue
 * - № <volume_issue>
 * - <pages-prefix> <pages_or_sheets>
 * - <year>
 * - ISBN <isbn> / <CertPrefix> <certificate_no>
 */
function fmtPubInfo(lang: "ru" | "kz" | "en", e: Entry): string {
  const parts: string[] = [];
  if (e.journal) parts.push(e.journal.trim());
  if (e.volume_issue) parts.push(`№ ${e.volume_issue.trim()}`);
  if (e.pages_or_sheets)
    parts.push(`${PAGES_PREFIX[lang]} ${e.pages_or_sheets.trim()}`);
  if (e.year) parts.push(e.year.trim());
  if (e.isbn) parts.push(`ISBN ${e.isbn.trim()}`);
  if (e.certificate_no)
    parts.push(`${CERT_PREFIX[lang]} ${e.certificate_no.trim()}`);
  return parts.filter(Boolean).join(", ");
}

/** Column 5: pages/sheets as-is (kept concise for the dedicated column). */
function fmtPages(e: Entry): string {
  return e.pages_or_sheets?.trim() || "";
}

/** Column 6: coauthors, comma-separated (skip blanks). */
function fmtCoauthors(e: Entry): string {
  const arr = Array.isArray(e.coauthors) ? e.coauthors : [];
  return arr
    .map((s) => (s || "").trim())
    .filter(Boolean)
    .join(", ");
}

/**
 * Convert entries list → 6-column rows consumed by docxtemplater loops.
 * no | title | format | pub_info | pages | coauthors
 */
function toRows(
  lang: "ru" | "kz" | "en",
  list: Entry[] = [],
  opts?: { ip?: boolean }
) {
  const isIP = !!opts?.ip;
  return list.map((e, idx) => ({
    no: String(idx + 1),
    title: e.title || "",
    format: fmtFormat(lang, e, isIP),
    pub_info: fmtPubInfo(lang, e),
    pages: fmtPages(e),
    coauthors: fmtCoauthors(e),
  }));
}

/** Select correct template asset id by locale (matches assets_list.json). */
function pickTemplateId(lang: "ru" | "kz" | "en") {
  if (lang === "kz") return "tpl_app7_kz_docx";
  if (lang === "en") return "tpl_app7_en_docx";
  return "tpl_app7_ru_docx";
}

/**
 * Public API: fill Appendix 7 template with current S1 values and download DOCX.
 *
 * Expects the DOCX to contain the 4 loops: i_rows, ii_rows, iii_rows, iv_rows.
 */
export async function generateApp7FromTemplate(
  values: S1Values,
  lang: "ru" | "kz" | "en",
  studentData?: StudentTemplateData
) {
  // 1) Resolve and fetch the template binary
  const assetId = pickTemplateId(lang);
  const url = getAssetUrl(assetId);
  if (!url || url === "#") throw new Error("Appendix 7 template not found");
  const ab = await fetch(url).then((r) => {
    if (!r.ok) throw new Error(`Failed to load template: ${r.statusText}`);
    return r.arrayBuffer();
  });

  // 2) Prepare docxtemplater engine
  const zip = new PizZip(ab);

  // Inject student data snippet if provided
  if (studentData) {
    const docFile = zip.file("word/document.xml");
    if (docFile) {
      const currentXml = docFile.asText();
      const withSnippet = injectSnippet(currentXml, lang);
      zip.file("word/document.xml", withSnippet);
    }
  }

  const doc = new Docxtemplater(zip, { paragraphLoop: true, linebreaks: true });

  // 3) Map S1 form collections to loop rows (6 columns only)
  const i_rows = toRows(lang, values.wos_scopus, { ip: false });
  const ii_rows = toRows(lang, values.kokson, { ip: false });
  const iii_rows = toRows(lang, values.conferences, { ip: false });
  const iv_rows = toRows(lang, values.ip, { ip: true }); // IP: 'format' becomes '—'

  // 4) Inject data and render
  doc.setData({
    i_rows,
    ii_rows,
    iii_rows,
    iv_rows,
    ...(studentData || {}),
  });
  try {
    doc.render();
  } catch (e) {
    // Helpful diagnostics in dev
    console.error("Docxtemplater render error", e);
    throw e;
  }

  // 5) Generate and download
  const out = doc.getZip().generate({
    type: "blob",
    mimeType:
      "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
  });
  const stamp = new Date().toISOString().slice(0, 10).replace(/-/g, "");
  const fname = `Appendix_7_${lang}_${stamp}.docx`;
  saveAs(out, fname);
}
