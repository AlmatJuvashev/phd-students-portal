import Docxtemplater from "docxtemplater";
import PizZip from "pizzip";
import { saveAs } from "file-saver";
import type { PublicAsset } from "@/lib/assets";
import { getAssetUrl } from "@/lib/assets";

type Locale = "ru" | "kz" | "en";

export type StudentTemplateData = {
  student_full_name: string;
  student_program: string;
  student_specialty: string;
  student_supervisors: string;
  submission_date: string;
  student_email: string;
  student_phone: string;
  dissertation_topic: string;
  student_department: string;
};

const LABELS: Record<
  Locale,
  {
    name: string;
    program: string;
    specialty: string;
    advisors: string;
    date: string;
    email: string;
    phone: string;
    topic: string;
    department: string;
  }
> = {
  ru: {
    name: "Докторант:",
    program: "Программа:",
    specialty: "Специальность:",
    advisors: "Научные руководители:",
    date: "Дата подачи:",
    email: "Email:",
    phone: "Телефон:",
    topic: "Тема диссертации:",
    department: "Кафедра:",
  },
  kz: {
    name: "Докторант:",
    program: "Бағдарлама:",
    specialty: "Мамандық:",
    advisors: "Ғылыми жетекшілер:",
    date: "Тапсыру күні:",
    email: "Email:",
    phone: "Телефон:",
    topic: "Диссертация тақырыбы:",
    department: "Кафедра:",
  },
  en: {
    name: "Doctoral candidate:",
    program: "Program:",
    specialty: "Specialty:",
    advisors: "Supervisors:",
    date: "Submission date:",
    email: "Email:",
    phone: "Phone:",
    topic: "Dissertation topic:",
    department: "Department:",
  },
};

const XML_NS = `<w:p><w:r><w:t xml:space="preserve">{{label}} </w:t></w:r><w:r><w:t>{{token}}</w:t></w:r></w:p>`;

function escapeXml(value: string) {
  return value
    .replace(/&/g, "&amp;")
    .replace(/</g, "&lt;")
    .replace(/>/g, "&gt;")
    .replace(/"/g, "&quot;")
    .replace(/'/g, "&apos;");
}

export function buildSnippet(locale: Locale) {
  const labels = LABELS[locale] || LABELS.ru;
  const paragraphs = [
    { label: labels.name, token: "student_full_name" },
    { label: labels.program, token: "student_program" },
    { label: labels.specialty, token: "student_specialty" },
    { label: labels.department, token: "student_department" },
    { label: labels.topic, token: "dissertation_topic" },
    { label: labels.advisors, token: "student_supervisors" },
    { label: labels.email, token: "student_email" },
    { label: labels.phone, token: "student_phone" },
    { label: labels.date, token: "submission_date" },
  ];
  return paragraphs
    .map((item) =>
      XML_NS.replace("{{label}}", escapeXml(item.label)).replace(
        "{{token}}",
        `{{${item.token}}}`
      )
    )
    .join("");
}

export function injectSnippet(xml: string, locale: Locale) {
  if (xml.includes("{{student_full_name}}")) return xml;
  const snippet = buildSnippet(locale);
  const firstParagraphEnd = xml.indexOf("</w:p>");
  if (firstParagraphEnd !== -1) {
    const idx = firstParagraphEnd + "</w:p>".length;
    return `${xml.slice(0, idx)}${snippet}${xml.slice(idx)}`;
  }
  const bodyTagStart = xml.indexOf("<w:body");
  if (bodyTagStart === -1) return xml;
  const bodyTagClose = xml.indexOf(">", bodyTagStart);
  if (bodyTagClose === -1) return xml;
  const insertPos = bodyTagClose + 1;
  return `${xml.slice(0, insertPos)}${snippet}${xml.slice(insertPos)}`;
}

function isTemplatableDoc(asset: PublicAsset) {
  const key = asset.storage?.key?.toLowerCase() || "";
  const id = asset.id?.toLowerCase() || "";
  return (
    key.endsWith(".docx") &&
    !id.includes("app7") &&
    asset.mime?.includes("wordprocessingml.document")
  );
}

function sanitizeFileName(name: string) {
  return name.replace(/[<>:"/\\|?*]+/g, "_");
}

export async function generateStudentTemplateDoc({
  asset,
  data,
  locale,
  fileLabel,
}: {
  asset: PublicAsset;
  data: StudentTemplateData;
  locale: Locale;
  fileLabel?: string;
}) {
  if (!isTemplatableDoc(asset)) {
    throw new Error("Asset is not a templatable DOCX");
  }
  const url = getAssetUrl(asset.id);
  if (!url || url === "#") {
    throw new Error("Template asset not found");
  }
  const arrayBuffer = await fetch(url).then((res) => {
    if (!res.ok) throw new Error(`Failed to load template (${res.status})`);
    return res.arrayBuffer();
  });
  const zip = new PizZip(arrayBuffer);
  const docFile = zip.file("word/document.xml");
  if (!docFile) {
    throw new Error("Invalid template structure");
  }
  const currentXml = docFile.asText();
  const withSnippet = injectSnippet(currentXml, locale);
  zip.file("word/document.xml", withSnippet);

  const doc = new Docxtemplater(zip, { paragraphLoop: true, linebreaks: true });
  doc.setData(data);
  doc.render();

  const blob = doc.getZip().generate({ type: "blob" });
  const safeTitle =
    sanitizeFileName(fileLabel || asset.title?.[locale] || asset.id) + ".docx";
  saveAs(blob, safeTitle);
}

export function supportsStudentDocTemplate(asset: PublicAsset) {
  return isTemplatableDoc(asset);
}

export function buildTemplateData(
  user: {
    full_name?: string;
    first_name?: string;
    last_name?: string;
    email?: string;
    phone?: string;
  } | null,
  profile: Record<string, any> | null | undefined,
  locale: "ru" | "kz" | "en"
): StudentTemplateData {
  const data = (profile?.form?.data ?? profile) as
    | Record<string, any>
    | undefined;
  const advisorsValue = data?.advisors_full_names;
  const advisors = Array.isArray(advisorsValue)
    ? advisorsValue
    : typeof advisorsValue === "string"
    ? advisorsValue.split(/\r?\n/)
    : [];
  const fullName =
    data?.full_name ||
    user?.full_name ||
    [user?.first_name, user?.last_name].filter(Boolean).join(" ");
  const program = data?.program || "";
  const specialty = data?.specialty || program;
  const localeMap: Record<string, string> = {
    ru: "ru-RU",
    kz: "kk-KZ",
    en: "en-US",
  };
  const formatter = new Intl.DateTimeFormat(localeMap[locale] || "ru-RU", {
    day: "numeric",
    month: "long",
    year: "numeric",
  });
  const submissionDate = formatter.format(new Date());

  const email = data?.email || user?.email || "";
  const phone = data?.phone || user?.phone || "";
  const topic = data?.dissertation_topic || data?.topic || "";
  const department = data?.department || "";

  return {
    student_full_name: fullName || "",
    student_program: program || "",
    student_specialty: specialty || "",
    student_supervisors: advisors
      .map((a: string) => (a || "").trim())
      .filter(Boolean)
      .join("\n"),
    submission_date: submissionDate,
    student_email: email,
    student_phone: phone,
    dissertation_topic: topic,
    student_department: department,
  };
}
