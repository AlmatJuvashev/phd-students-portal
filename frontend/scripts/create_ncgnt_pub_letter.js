import PizZip from "pizzip";
import fs from "fs";
import path from "path";
import { fileURLToPath } from 'url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

// Create NCGNT publication certificate letter template with TABLE PER PUBLICATION
// SUPER SAFE XML: Re-typed tags, simplified structure
const createNCGNTPublicationLetterTemplate = () => {
  const xml = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">
  <w:body>
    <w:p><w:pPr><w:jc w:val="center"/></w:pPr><w:r><w:rPr><w:b/></w:rPr><w:t>Пример письма в АО «НЦГНТЭ»</w:t></w:r></w:p>
    <w:p><w:pPr><w:jc w:val="center"/></w:pPr><w:r><w:rPr><w:b/></w:rPr><w:t>для получения справки о публикациях</w:t></w:r></w:p>
    <w:p><w:pPr><w:spacing w:after="200"/></w:pPr></w:p>
    <w:p><w:r><w:rPr><w:b/></w:rPr><w:t>Email: astana@ncste.kz</w:t></w:r></w:p>
    <w:p><w:r><w:rPr><w:b/></w:rPr><w:t>Тема письма:</w:t></w:r></w:p>
    <w:p><w:r><w:t>Запрос на получение справки о публикациях в индексируемых журналах</w:t></w:r></w:p>
    <w:p><w:pPr><w:spacing w:after="200"/></w:pPr></w:p>
    <w:p><w:r><w:rPr><w:b/></w:rPr><w:t>Текст письма:</w:t></w:r></w:p>
    <w:p><w:r><w:t>Уважаемые сотрудники АО «Национальный центр государственной научно-технической экспертизы»!</w:t></w:r></w:p>
    <w:p><w:pPr><w:spacing w:after="200"/></w:pPr></w:p>
    
    <w:p><w:r><w:t>Прошу выдать справку о публикациях в индексируемых журналах для PhD докторанта:</w:t></w:r></w:p>
    
    <w:p><w:r><w:rPr><w:b/></w:rPr><w:t>{{student_full_name}}</w:t></w:r></w:p>
    
    <w:p><w:pPr><w:spacing w:after="200"/></w:pPr></w:p>
    <w:p><w:r><w:rPr><w:b/></w:rPr><w:t>Сведения о публикациях:</w:t></w:r></w:p>
    <w:p><w:pPr><w:spacing w:after="100"/></w:pPr></w:p>
    
    <w:p><w:r><w:t>{#publications}</w:t></w:r></w:p>
    
    <w:tbl>
      <w:tblPr>
        <w:tblW w:w="5000" w:type="pct"/>
        <w:tblBorders>
          <w:top w:val="single" w:sz="4"/>
          <w:left w:val="single" w:sz="4"/>
          <w:bottom w:val="single" w:sz="4"/>
          <w:right w:val="single" w:sz="4"/>
          <w:insideH w:val="single" w:sz="4"/>
          <w:insideV w:val="single" w:sz="4"/>
        </w:tblBorders>
      </w:tblPr>
      <w:tr>
        <w:tc><w:p><w:r><w:rPr><w:b/></w:rPr><w:t>{no}. Название статьи:</w:t></w:r></w:p></w:tc>
        <w:tc><w:p><w:r><w:t>{title}</w:t></w:r></w:p></w:tc>
      </w:tr>
      <w:tr>
        <w:tc><w:p><w:r><w:rPr><w:b/></w:rPr><w:t>Авторы:</w:t></w:r></w:p></w:tc>
        <w:tc><w:p><w:r><w:t>{authors}</w:t></w:r></w:p></w:tc>
      </w:tr>
      <w:tr>
        <w:tc><w:p><w:r><w:rPr><w:b/></w:rPr><w:t>Название журнала:</w:t></w:r></w:p></w:tc>
        <w:tc><w:p><w:r><w:t>{journal}</w:t></w:r></w:p></w:tc>
      </w:tr>
      <w:tr>
        <w:tc><w:p><w:r><w:rPr><w:b/></w:rPr><w:t>Номер журнала:</w:t></w:r></w:p></w:tc>
        <w:tc><w:p><w:r><w:t>{volume_issue}</w:t></w:r></w:p></w:tc>
      </w:tr>
      <w:tr>
        <w:tc><w:p><w:r><w:rPr><w:b/></w:rPr><w:t>Том:</w:t></w:r></w:p></w:tc>
        <w:tc><w:p><w:r><w:t>{vol}</w:t></w:r></w:p></w:tc>
      </w:tr>
      <w:tr>
        <w:tc><w:p><w:r><w:rPr><w:b/></w:rPr><w:t>Год:</w:t></w:r></w:p></w:tc>
        <w:tc><w:p><w:r><w:t>{year}</w:t></w:r></w:p></w:tc>
      </w:tr>
      <w:tr>
        <w:tc><w:p><w:r><w:rPr><w:b/></w:rPr><w:t>ISSN (печ.):</w:t></w:r></w:p></w:tc>
        <w:tc><w:p><w:r><w:t>{issn_print}</w:t></w:r></w:p></w:tc>
      </w:tr>
      <w:tr>
        <w:tc><w:p><w:r><w:rPr><w:b/></w:rPr><w:t>ISSN (онлайн):</w:t></w:r></w:p></w:tc>
        <w:tc><w:p><w:r><w:t>{issn_online}</w:t></w:r></w:p></w:tc>
      </w:tr>
      <w:tr>
        <w:tc><w:p><w:r><w:rPr><w:b/></w:rPr><w:t>DOI:</w:t></w:r></w:p></w:tc>
        <w:tc><w:p><w:r><w:t>{doi}</w:t></w:r></w:p></w:tc>
      </w:tr>
    </w:tbl>
    
    <w:p><w:r><w:t>{/publications}</w:t></w:r></w:p>

    <w:p><w:pPr><w:spacing w:after="200"/></w:pPr></w:p>
    <w:p><w:r><w:t>К письму прилагаю удостоверение личности в цифровом формате.</w:t></w:r></w:p>
    <w:p><w:pPr><w:spacing w:after="200"/></w:pPr></w:p>
    <w:p><w:r><w:rPr><w:b/></w:rPr><w:t>Контактные данные:</w:t></w:r></w:p>
    <w:p><w:r><w:t>Ф.И.О.: </w:t></w:r><w:r><w:t>{{student_full_name}}</w:t></w:r></w:p>
    <w:p><w:r><w:t>Телефон: </w:t></w:r><w:r><w:t>{{student_phone}}</w:t></w:r></w:p>
    <w:p><w:r><w:t>Электронная почта: </w:t></w:r><w:r><w:t>{{student_email}}</w:t></w:r></w:p>
    <w:p><w:pPr><w:spacing w:after="200"/></w:pPr></w:p>
    <w:p><w:r><w:t>С уважением,</w:t></w:r></w:p>
    <w:p><w:r><w:t>{{student_full_name}}</w:t></w:r></w:p>
    <w:p><w:r><w:t>{{day}} {{month}} {{year}} г.</w:t></w:r></w:p>
  </w:body>
</w:document>`;

  const outputPath = path.resolve(__dirname, "../public/templates/letter_to_ncste_publication_certificate_ru.docx");
  
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
  console.log("Created NCGNT publication certificate letter template with SUPER SAFE XML at:", outputPath);
};

createNCGNTPublicationLetterTemplate();
