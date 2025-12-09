import PizZip from "pizzip";
import fs from "fs";
import path from "path";
import { fileURLToPath } from 'url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

// Fix letter to rector defense request template
const fixRectorLetterTemplate = () => {
  // Simplified clean XML structure
  const xml = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">
  <w:body>
    <w:p><w:pPr><w:jc w:val="right"/></w:pPr><w:r><w:t xml:space="preserve">Ректору КазНМУ им. С.Д. Асфендиярова</w:t></w:r></w:p>
    <w:p><w:pPr><w:jc w:val="right"/></w:pPr><w:r><w:t xml:space="preserve">Шоранову Марату Едигеевичу</w:t></w:r></w:p>
    <w:p><w:pPr><w:jc w:val="right"/></w:pPr><w:r><w:t xml:space="preserve">от PhD докторанта:</w:t></w:r></w:p>
    <w:p><w:pPr><w:jc w:val="right"/></w:pPr><w:r><w:rPr><w:b/></w:rPr><w:t>[%student_full_name%]</w:t></w:r></w:p>
    <w:p><w:pPr><w:jc w:val="right"/></w:pPr><w:r><w:t xml:space="preserve">Специальность: </w:t></w:r><w:r><w:rPr><w:b/></w:rPr><w:t>[%student_specialty%]</w:t></w:r></w:p>
    <w:p><w:pPr><w:jc w:val="right"/></w:pPr><w:r><w:t xml:space="preserve">Тема диссертации: </w:t></w:r><w:r><w:rPr><w:b/></w:rPr><w:t>[%dissertation_topic%]</w:t></w:r></w:p>
    <w:p><w:pPr><w:jc w:val="right"/></w:pPr><w:r><w:t xml:space="preserve">Программа: </w:t></w:r><w:r><w:rPr><w:b/></w:rPr><w:t>[%student_program%]</w:t></w:r></w:p>
    <w:p><w:pPr><w:spacing w:after="200"/></w:pPr></w:p>
    <w:p><w:pPr><w:jc w:val="center"/></w:pPr><w:r><w:rPr><w:b/></w:rPr><w:t>ЗАЯВЛЕНИЕ</w:t></w:r></w:p>
    <w:p><w:pPr><w:jc w:val="center"/></w:pPr><w:r><w:rPr><w:b/></w:rPr><w:t xml:space="preserve">о приёме к защите диссертации</w:t></w:r></w:p>
    <w:p><w:pPr><w:spacing w:after="200"/></w:pPr></w:p>
    <w:p><w:pPr><w:ind w:firstLine="567"/><w:jc w:val="both"/></w:pPr><w:r><w:t xml:space="preserve">Прошу принять меня к защите диссертации на соискание степени доктора философии (PhD) по специальности [%student_specialty%] на тему: «[%dissertation_topic%]».</w:t></w:r></w:p>
    <w:p><w:pPr><w:spacing w:after="200"/></w:pPr></w:p>
    <w:p><w:pPr><w:ind w:firstLine="567"/></w:pPr><w:r><w:rPr><w:b/></w:rPr><w:t>Приложения:</w:t></w:r></w:p>
    <w:p><w:pPr><w:ind w:left="567"/></w:pPr><w:r><w:rPr><w:i/></w:rPr><w:t xml:space="preserve">1. Диссертация в электронном виде (PDF)</w:t></w:r></w:p>
    <w:p><w:pPr><w:ind w:left="567"/></w:pPr><w:r><w:rPr><w:i/></w:rPr><w:t xml:space="preserve">2. Заключение ЛКБ</w:t></w:r></w:p>
    <w:p><w:pPr><w:ind w:left="567"/></w:pPr><w:r><w:rPr><w:i/></w:rPr><w:t xml:space="preserve">3. Справка НЦГНТЭ о публикациях</w:t></w:r></w:p>
    <w:p><w:pPr><w:spacing w:after="200"/></w:pPr></w:p>
    <w:tbl>
      <w:tblPr><w:tblW w:w="5000" w:type="pct"/></w:tblPr>
      <w:tr>
        <w:tc><w:p><w:r><w:rPr><w:b/></w:rPr><w:t>Дата</w:t></w:r></w:p><w:p><w:r><w:rPr><w:b/></w:rPr><w:t>«[%day%]» [%month%] [%year%] г.</w:t></w:r></w:p></w:tc>
        <w:tc><w:p><w:r><w:rPr><w:b/></w:rPr><w:t xml:space="preserve">Подпись докторанта: </w:t></w:r></w:p><w:p><w:r><w:rPr><w:b/></w:rPr><w:t>_____________________</w:t></w:r></w:p></w:tc>
      </w:tr>
    </w:tbl>
  </w:body>
</w:document>`;

  const outputPath = path.resolve(__dirname, "../public/templates/letter_to_rector_request_defense_ru.docx");
  
  if (!fs.existsSync(outputPath)) {
    console.error("Error: Template file not found at", outputPath);
    return;
  }

  const content = fs.readFileSync(outputPath);
  const zip = new PizZip(content);
  
  // Update document.xml
  zip.file("word/document.xml", xml);
  
  const buffer = zip.generate({ type: "nodebuffer" });
  fs.writeFileSync(outputPath, buffer);
  console.log("Fixed rector defense letter template at:", outputPath);
};

fixRectorLetterTemplate();
