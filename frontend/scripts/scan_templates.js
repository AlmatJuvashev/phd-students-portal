import PizZip from "pizzip";
import fs from "fs";
import path from "path";
import { fileURLToPath } from 'url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

const templatesDir = path.resolve(__dirname, "../public/templates");

const scanTemplates = () => {
  const files = fs.readdirSync(templatesDir).filter(f => f.endsWith(".docx"));
  
  console.log(`Scanning ${files.length} templates for '{{' delimiters...`);
  
  files.forEach(file => {
    const filePath = path.join(templatesDir, file);
    try {
      const content = fs.readFileSync(filePath);
      const zip = new PizZip(content);
      const docXml = zip.file("word/document.xml")?.asText();
      
      if (docXml && docXml.includes("{{")) {
        console.log(`[FOUND] ${file} contains '{{'`);
      }
    } catch (e) {
      console.error(`[ERROR] Failed to read ${file}:`, e.message);
    }
  });
};

scanTemplates();
