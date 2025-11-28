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
  publications?: Array<{
    no: string;
    title: string;
    authors: string;
    journal: string;
    volume_issue: string;
    vol: string;
    year: string;
    issn_print: string;
    issn_online: string;
    doi: string;
  }>;
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
  data: StudentTemplateData
) {
  let out = xml;
  (Object.keys(data) as Array<keyof StudentTemplateData>).forEach((key) => {
    // Skip non-string values (e.g., publications array handled by docxtemplater)
    if (typeof data[key] !== "string") return;
    
    const pattern = new RegExp(`{{${key}}}`, "g");
    const encoded = encodeValue(data[key] as string);
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
  // Add cache-busting parameter to ensure fresh template
  const cacheBustUrl = url.includes('?') ? `${url}&_t=${Date.now()}` : `${url}?_t=${Date.now()}`;
  console.log("[template] fetching from:", cacheBustUrl);

  if (url.includes("normocontrol_letter")) {
    console.log("[template] NCGNT letter detected", {
      url,
      dataKeys: Object.keys(data),
      studentName: data.student_full_name,
      dissertationTopic: data.dissertation_topic,
    });
  }

  const arrayBuffer = await fetch(cacheBustUrl, { cache: 'no-store' }).then((res) => {
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

  const rawXml = docFile.asText();
  console.log("[template] RAW XML (before normalization) length:", rawXml.length);
  console.log("[template] RAW XML first 500 chars:", rawXml.substring(0, 500));
  
  // Check for split tags in RAW XML (before normalization)
  const rawSplitOpen = rawXml.match(/<w:t[^>]*>\s*\{\s*<\/w:t>/g);
  const rawSplitClose = rawXml.match(/<w:t[^>]*>\s*\}\s*<\/w:t>/g);
  console.log("[template] RAW split open braces:", rawSplitOpen?.length || 0);
  console.log("[template] RAW split close braces:", rawSplitClose?.length || 0);
  
  if (rawSplitOpen && rawSplitOpen.length > 0) {
    console.log("[template] RAW examples of split braces:", rawSplitOpen.slice(0, 2));
  }
  
  const currentXml = normalizeXml(rawXml);
  console.log("[template] raw XML length:", currentXml.length);
  console.log("[template] raw XML first 500 chars:", currentXml.substring(0, 500));
  
  // Check for split tags in raw XML
  const splitOpenBraces = currentXml.match(/<w:t[^>]*>\s*\{\s*<\/w:t>/g);
  const splitCloseBraces = currentXml.match(/<w:t[^>]*>\s*\}\s*<\/w:t>/g);
  console.log("[template] split open braces found:", splitOpenBraces?.length || 0);
  console.log("[template] split close braces found:", splitCloseBraces?.length || 0);
  
  if (splitOpenBraces && splitOpenBraces.length > 0) {
    console.log("[template] examples of split open braces:", splitOpenBraces.slice(0, 3));
  }
  
  const hasTokens = currentXml.includes("{{") || currentXml.includes("[%");
  const containsStudentName = currentXml.includes("student_full_name");
  console.log("[template] token detection", {
    hasTokens,
    containsStudentName,
    xmlSnippet: currentXml.substring(
      Math.max(0, currentXml.indexOf("student") - 50),
      currentXml.indexOf("student") + 50
    ),
    dataKeys: Object.keys(data),
  });

  let blob: Blob;
  if (hasTokens) {
    // Create a NEW zip with normalized XML to ensure docxtemplater reads clean XML
    zip.file("word/document.xml", currentXml);
    const normalizedZipBuffer = zip.generate({ type: "uint8array" });
    const normalizedZip = new PizZip(normalizedZipBuffer);
    console.log("[template] created new zip with normalized XML");
    
    // Use docxtemplater to fill existing placeholders (handles split runs)
    try {
      const doc = new Docxtemplater(normalizedZip, {
        paragraphLoop: true,
        linebreaks: true,
        delimiters: { start: "[%", end: "%]" },
      });
      
      doc.setData(data);
      doc.render();
      blob = doc.getZip().generate({ type: "blob" });
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
    sanitizeFileName(
      fileLabel || 
      asset.title?.[locale] || 
      asset.title?.ru || 
      asset.title?.en || 
      "Zayavlenie_v_OMiD"
    ) + ".docx";
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

  // Extract and format publications from S1 data
  const publications: StudentTemplateData["publications"] = [];
  const sections = ["wos_scopus", "kokson", "conferences"] as const;
  
  console.log("[buildTemplateData] Full profile structure:", JSON.stringify(profile, null, 2));
  
  console.log("[buildTemplateData] Extracting publications from profile:", JSON.stringify({
    hasData: !!data,
    profileKeys: profile ? Object.keys(profile) : [],
    formKeys: profile?.form ? Object.keys(profile.form) : [],
    dataKeys: data ? Object.keys(data) : [],
    sections: sections.map(s => ({
      key: s,
      hasSection: !!data?.[s],
      count: Array.isArray(data?.[s]) ? data[s].length : 0
    }))
  }));
  
  sections.forEach((sectionKey) => {
    const entries = data?.[sectionKey];
    if (Array.isArray(entries)) {
      console.log(`[buildTemplateData] Processing section ${sectionKey}, entries:`, entries.length);
      entries.forEach((entry: any, index) => {
        console.log(`[buildTemplateData] ${sectionKey}[${index}] raw:`, JSON.stringify(entry));
        if (entry?.title) {
          // Extract volume from volume_issue (e.g., "12(3)" -> "12")
          const volumeMatch = entry.volume_issue?.match(/^(\d+)/);
          const vol = volumeMatch ? volumeMatch[1] : "";
          
          // Format authors: student + coauthors
          const coauthorsList = Array.isArray(entry.coauthors)
            ? entry.coauthors.filter((c: string) => c?.trim()).join(", ")
            : "";
          const authors = coauthorsList
            ? `${fullName}, ${coauthorsList}`
            : fullName;
          
          publications.push({
            no: String(publications.length + 1),
            title: entry.title || "",
            authors,
            journal: entry.journal || "",
            volume_issue: entry.volume_issue || "",
            vol,
            year: entry.year || "",
            issn_print: entry.issn_print || "",
            issn_online: entry.issn_online || "",
            doi: entry.doi || "",
          });
        }
      });
    }
  });

  console.log("[buildTemplateData] Publications extracted:", {
    total: publications.length,
    sample: publications.slice(0, 2).map(p => ({
      no: p.no,
      title: p.title.substring(0, 50),
      hasAuthors: !!p.authors,
      hasJournal: !!p.journal
    }))
  });

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
    publications: publications.length > 0 ? publications : undefined,
  };
}
