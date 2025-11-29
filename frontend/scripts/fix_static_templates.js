import PizZip from "pizzip";
import fs from "fs";
import path from "path";
import { fileURLToPath } from 'url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

const templatesDir = path.resolve(__dirname, "../public/templates");

const filesToFix = [
  "normocontrol_letter.docx",
  "tpl_app5_ru.docx"
];

const fixTemplates = () => {
  filesToFix.forEach(file => {
    const filePath = path.join(templatesDir, file);
    console.log(`Fixing ${file}...`);
    
    try {
      if (!fs.existsSync(filePath)) {
        console.error(`[ERROR] File not found: ${filePath}`);
        return;
      }

      const content = fs.readFileSync(filePath);
      const zip = new PizZip(content);
      
      let docXml = zip.file("word/document.xml")?.asText();
      
      if (!docXml) {
        console.error(`[ERROR] Could not read document.xml from ${file}`);
        return;
      }
      
      // Replace {{ with [% and }} with %]
      // We use a regex to replace all occurrences
      const fixedXml = docXml.replace(/\{\{/g, "[%").replace(/\}\}/g, "%]");
      
      if (docXml === fixedXml) {
        console.log(`[INFO] No changes needed for ${file}`);
      } else {
        zip.file("word/document.xml", fixedXml);
        const buffer = zip.generate({ type: "nodebuffer" });
        fs.writeFileSync(filePath, buffer);
        console.log(`[SUCCESS] Updated ${file}`);
      }
      
    } catch (e) {
      console.error(`[ERROR] Failed to fix ${file}:`, e.message);
    }
  });
};

fixTemplates();
