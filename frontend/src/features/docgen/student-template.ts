import PizZip from "pizzip";
import Docxtemplater from "docxtemplater";
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
  day: string;
  month: string;
  year: string;
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

function encodeValue(value: string) {
  const safe = escapeXml(value || "");
  if (!safe.includes("\n")) return safe;
  const parts = safe.split(/\r?\n/);
  return parts.join("</w:t><w:br/><w:t>");
}

function replaceTokens(
  xml: string,
  data: Record<keyof StudentTemplateData, string>
) {
  let out = xml;
  (Object.keys(data) as Array<keyof StudentTemplateData>).forEach((key) => {
    const pattern = new RegExp(`{{${key}}}`, "g");
    const encoded = encodeValue(data[key] || "");
    out = out.replace(pattern, encoded);
  });
  return out;
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
  console.log("[template] start", {
    assetId: asset.id,
    locale,
    url,
    hasData: !!data,
    name: data?.student_full_name,
    specialty: data?.student_specialty,
  });
  const zip = new PizZip(arrayBuffer);
  const docFile = zip.file("word/document.xml");
  if (!docFile) {
    throw new Error("Invalid template structure");
  }
  // Pre-process XML to fix common template errors (duplicate tags)
  const normalizeXml = (xml: string) => {
    let out = xml;
    // Fix split open braces: { ... { -> {{
    // Matches <w:t>{<w:t> ... <w:t>{<w:t> and merges them, allowing for whitespace between tags
    out = out.replace(/(<w:t[^>]*>)\s*\{\s*(<\/w:t>(?:<[^>]+>|\s+)*<w:t[^>]*>)\s*\{\s*(<\/w:t>)/g, "$1{{$3");
    // Fix split close braces: } ... } -> }}
    out = out.replace(/(<w:t[^>]*>)\s*\}\s*(<\/w:t>(?:<[^>]+>|\s+)*<w:t[^>]*>)\s*\}\s*(<\/w:t>)/g, "$1}}$3");
    
    // Specific fix for student_full_name if the above generic one misses (e.g. due to complex nesting)
    // Look for { ... { ... student_full_name
    out = out.replace(/\{\s*(<\/w:t>(?:<[^>]+>|\s+)*<w:t[^>]*>)\s*\{\s*(<\/w:t>(?:<[^>]+>|\s+)*<w:t[^>]*>)\s*student_full_name/g, "{{$1$2student_full_name");
    
    // Fix duplicate open tags: {{...{{ -> {{
    out = out.replace(/({{(?:<[^>]+>)*){{/g, "$1");
    // Fix duplicate close tags: }}...}} -> }}
    out = out.replace(/}}((?:<[^>]+>)*)}}/g, "$1}}");
    
    return out;
  };

  const currentXml = normalizeXml(docFile.asText());
  const hasTokens = currentXml.includes("{{");
  
  // Debug: Find the context around student_full_name in the XML
  const nameIndex = currentXml.indexOf("student_full_name");
  const xmlSnippet = nameIndex !== -1 
    ? currentXml.substring(Math.max(0, nameIndex - 100), Math.min(currentXml.length, nameIndex + 100))
    : "Not found";

  console.log("[template] token detection", {
    hasTokens,
    containsStudentName: currentXml.includes("student_full_name"),
    xmlSnippet,
    dataKeys: Object.keys(data),
  });

  let blob: Blob;
  if (hasTokens) {
    // Use docxtemplater to fill existing placeholders (handles split runs)
    try {
      const doc = new Docxtemplater(zip, {
        paragraphLoop: true,
        linebreaks: true,
      });
      
      // Update the zip with normalized XML before rendering
      doc.loadZip(new PizZip(
        zip.generate({ type: "nodebuffer" })
      ));
      zip.file("word/document.xml", currentXml);
      
      const docNormalized = new Docxtemplater(zip, {
        paragraphLoop: true,
        linebreaks: true,
      });
      docNormalized.setData(data);
      docNormalized.render();
      blob = docNormalized.getZip().generate({ type: "blob" });
      console.log("[template] docxtemplater render success");
    } catch (err) {
      console.error("[template] docxtemplater error, falling back", err);
      const fallbackZip = new PizZip(arrayBuffer);
      // Improved fallback: try to handle simple split runs for known keys
      let fallbackXml = currentXml;
      // ... (existing fallback logic)
      fallbackXml = replaceTokens(fallbackXml, data);
      fallbackZip.file("word/document.xml", fallbackXml);
      blob = fallbackZip.generate({ type: "blob" });
    }
  } else {
    // No tokens: inject snippet once near top, then simple replacement
    const withSnippet = injectSnippet(currentXml, locale);
    const filledXml = replaceTokens(withSnippet, data);
    console.log("[template] replacement", {
      hasName: filledXml.includes(data.student_full_name),
      hasProgram: filledXml.includes(data.student_program),
      hasAdvisors: !!data.student_supervisors,
    });
    zip.file("word/document.xml", filledXml);
    blob = zip.generate({ type: "blob" });
  }

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

  // Get current date parts for separate placeholders
  const now = new Date();
  const day = String(now.getDate()).padStart(2, "0"); // Two digits: "27"
  
  // Get month name based on locale
  const monthFormatter = new Intl.DateTimeFormat(localeMap[locale] || "ru-RU", {
    month: "long",
  });
  const month = monthFormatter.format(now); // e.g., "ноября"
  
  const year = String(now.getFullYear()); // Four digits: "2025"

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
    day,
    month,
    year,
  };
}
