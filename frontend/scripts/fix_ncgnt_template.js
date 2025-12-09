import PizZip from "pizzip";
import fs from "fs";
import path from "path";
import { fileURLToPath } from 'url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

// Fix NCGNT normocontrol letter template - simplified version without hyperlink
const fixNCGNTTemplate = () => {
  // Simplified XML without hyperlink to avoid rId issues
  const xml = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">
  <w:body>
    <w:p><w:pPr><w:jc w:val="center"/></w:pPr><w:r><w:rPr><w:b/></w:rPr><w:t>Пример письма в АО «НЦГНТЭ» (для нормоконтроля)</w:t></w:r></w:p>
    <w:p></w:p>
    <w:p><w:r><w:rPr><w:b/></w:rPr><w:t xml:space="preserve">Email: astana@ncste.kz</w:t></w:r></w:p>
    <w:p><w:r><w:rPr><w:b/></w:rPr><w:t>Тема письма:</w:t></w:r><w:r><w:br/><w:t>Запрос на проведение нормоконтроля диссертации</w:t></w:r></w:p>
    <w:p></w:p>
    <w:p><w:r><w:rPr><w:b/></w:rPr><w:t>Текст письма:</w:t></w:r></w:p>
    <w:p><w:r><w:t>Уважаемые сотрудники АО «Национальный центр государственной научно-технической экспертизы»!</w:t></w:r></w:p>
    <w:p><w:r><w:t xml:space="preserve">Прошу провести нормоконтроль диссертационной работы на тему</w:t></w:r></w:p>
    <w:p><w:r><w:t>«{{dissertation_topic}}»</w:t></w:r></w:p>
    <w:p><w:r><w:t xml:space="preserve">по специальности {{student_specialty}}.</w:t></w:r></w:p>
    <w:p></w:p>
    <w:p><w:r><w:t xml:space="preserve">Диссертация подготовлена в формате DOCX и прилагается к данному письму.</w:t></w:r></w:p>
    <w:p><w:r><w:t xml:space="preserve">Также прилагаю квитанцию об оплате услуги нормоконтроля.</w:t></w:r></w:p>
    <w:p></w:p>
    <w:p><w:r><w:t>Контактные данные для обратной связи:</w:t></w:r></w:p>
    <w:p><w:r><w:t xml:space="preserve">Ф.И.О.: {{student_full_name}}</w:t></w:r></w:p>
    <w:p><w:r><w:t xml:space="preserve">Докторантура: {{student_program}}</w:t></w:r></w:p>
    <w:p><w:r><w:t xml:space="preserve">Телефон: {{student_phone}}</w:t></w:r></w:p>
    <w:p><w:r><w:t xml:space="preserve">Электронная почта: {{student_email}}</w:t></w:r></w:p>
    <w:p></w:p>
    <w:p><w:r><w:t>С уважением,</w:t></w:r></w:p>
    <w:p><w:r><w:t>{{student_full_name}}</w:t></w:r></w:p>
    <w:p><w:r><w:t>{{day}} {{month}} {{year}} г.</w:t></w:r></w:p>
  </w:body>
</w:document>`;

  const outputPath = path.resolve(__dirname, "../public/templates/normocontrol_letter.docx");
  
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
  console.log("Fixed NCGNT template (simplified XML) at:", outputPath);
};

fixNCGNTTemplate();
