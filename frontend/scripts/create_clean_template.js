import PizZip from "pizzip";
import Docxtemplater from "docxtemplater";
import fs from "fs";
import path from "path";
import { fileURLToPath } from 'url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

// Create a minimal valid .docx with clean placeholders
const createCleanTemplate = () => {
  const xml = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">
  <w:body>
    <w:p><w:pPr><w:jc w:val="right"/></w:pPr><w:r><w:t xml:space="preserve">Ректору КазНМУ им. С.Д. Асфендиярова</w:t></w:r></w:p>
    <w:p><w:pPr><w:jc w:val="right"/></w:pPr><w:r><w:t xml:space="preserve">Шоранову Марату Едигеевичу</w:t></w:r></w:p>
    <w:p><w:pPr><w:jc w:val="right"/></w:pPr><w:r><w:t xml:space="preserve">от PhD докторанта:</w:t></w:r></w:p>
    <w:p><w:pPr><w:jc w:val="right"/></w:pPr><w:r><w:rPr><w:b/></w:rPr><w:t>{{student_full_name}}</w:t></w:r></w:p>
    <w:p><w:pPr><w:jc w:val="right"/></w:pPr><w:r><w:t xml:space="preserve">Специальность: </w:t></w:r><w:r><w:rPr><w:b/></w:rPr><w:t>{{student_specialty}}</w:t></w:r></w:p>
    <w:p><w:pPr><w:jc w:val="right"/></w:pPr><w:r><w:t xml:space="preserve">Тема диссертации: </w:t></w:r><w:r><w:rPr><w:b/></w:rPr><w:t>{{dissertation_topic}}</w:t></w:r></w:p>
    <w:p><w:pPr><w:jc w:val="right"/></w:pPr><w:r><w:t xml:space="preserve">Научный руководитель: </w:t></w:r><w:r><w:rPr><w:b/></w:rPr><w:t>{{student_supervisors}}</w:t></w:r></w:p>
    <w:p><w:pPr><w:spacing w:after="200"/></w:pPr></w:p>
    <w:p><w:pPr><w:jc w:val="center"/></w:pPr><w:r><w:rPr><w:b/></w:rPr><w:t>ЗАЯВЛЕНИЕ</w:t></w:r></w:p>
    <w:p><w:pPr><w:jc w:val="center"/></w:pPr><w:r><w:rPr><w:b/></w:rPr><w:t xml:space="preserve">на проведение предварительной экспертизы диссертации в ОМиД</w:t></w:r></w:p>
    <w:p><w:pPr><w:jc w:val="center"/></w:pPr><w:r><w:rPr><w:b/></w:rPr><w:t>(Отдел мониторинга и диссертаций)</w:t></w:r></w:p>
    <w:p><w:pPr><w:spacing w:after="200"/></w:pPr></w:p>
    <w:p><w:pPr><w:ind w:firstLine="567"/><w:jc w:val="both"/></w:pPr><w:r><w:t xml:space="preserve">Прошу направить мою диссертацию на предварительную экспертизу в ОМиД в соответствии с установленными правилами.</w:t></w:r></w:p>
    <w:p><w:pPr><w:spacing w:after="200"/></w:pPr></w:p>
    <w:p><w:pPr><w:ind w:firstLine="567"/></w:pPr><w:r><w:rPr><w:b/></w:rPr><w:t>Приложения:</w:t></w:r></w:p>
    <w:p><w:pPr><w:ind w:left="567"/></w:pPr><w:r><w:rPr><w:i/></w:rPr><w:t xml:space="preserve">1. Текст диссертации (PDF)</w:t></w:r></w:p>
    <w:p><w:pPr><w:ind w:left="567"/></w:pPr><w:r><w:rPr><w:i/></w:rPr><w:t xml:space="preserve">2. Отчёт антиплагиата</w:t></w:r></w:p>
    <w:p><w:pPr><w:ind w:left="567"/></w:pPr><w:r><w:rPr><w:i/></w:rPr><w:t xml:space="preserve">3. Список публикаций (Приложение 7)</w:t></w:r></w:p>
    <w:p><w:pPr><w:ind w:left="567"/></w:pPr><w:r><w:rPr><w:i/></w:rPr><w:t xml:space="preserve">4. Отзыв научного руководителя</w:t></w:r></w:p>
    <w:p><w:pPr><w:spacing w:after="200"/></w:pPr></w:p>
    <w:tbl>
      <w:tblPr><w:tblW w:w="5000" w:type="pct"/></w:tblPr>
      <w:tr>
        <w:tc><w:p><w:r><w:rPr><w:b/></w:rPr><w:t>Дата</w:t></w:r></w:p><w:p><w:r><w:rPr><w:b/></w:rPr><w:t>«{{day}}» {{month}} {{year}} г.</w:t></w:r></w:p></w:tc>
        <w:tc><w:p><w:r><w:rPr><w:b/></w:rPr><w:t xml:space="preserve">Подпись докторанта: </w:t></w:r></w:p><w:p><w:r><w:rPr><w:b/></w:rPr><w:t>_____________________</w:t></w:r></w:p></w:tc>
      </w:tr>
    </w:tbl>
  </w:body>
</w:document>`;

  const zip = new PizZip();
  
  // Add required files
  zip.file("word/document.xml", xml);
  zip.file("[Content_Types].xml", `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Types xmlns="http://schemas.openxmlformats.org/package/2006/content-types">
  <Default Extension="rels" ContentType="application/vnd.openxmlformats-package.relationships+xml"/>
  <Default Extension="xml" ContentType="application/xml"/>
  <Override PartName="/word/document.xml" ContentType="application/vnd.openxmlformats-officedocument.wordprocessingml.document.main+xml"/>
</Types>`);
  
  zip.file("_rels/.rels", `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
  <Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="word/document.xml"/>
</Relationships>`);

  const buffer = zip.generate({ type: "nodebuffer" });
  const outputPath = path.resolve(__dirname, "../public/templates/Zayavlenie_v_OMiD_ru.docx");
  
  // Backup the old one
  const backupPath = path.resolve(__dirname, "../public/templates/Zayavlenie_v_OMiD_ru.docx.backup");
  if (fs.existsSync(outputPath)) {
    fs.copyFileSync(outputPath, backupPath);
    console.log("Backed up old template to:", backupPath);
  }
  
  fs.writeFileSync(outputPath, buffer);
  console.log("Created clean template at:", outputPath);
};

createCleanTemplate();
