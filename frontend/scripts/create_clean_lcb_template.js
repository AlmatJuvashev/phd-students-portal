import PizZip from "pizzip";
import fs from "fs";
import path from "path";
import { fileURLToPath } from 'url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

// Create a clean LCB request template
const createCleanLCBTemplate = () => {
  const xml = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">
  <w:body>
    <w:p><w:pPr><w:jc w:val="right"/></w:pPr><w:r><w:t xml:space="preserve">Председателю ЛКБ</w:t></w:r></w:p>
    <w:p><w:pPr><w:jc w:val="right"/></w:pPr><w:r><w:t xml:space="preserve">КазНМУ им. С.Д. Асфендиярова</w:t></w:r></w:p>
    <w:p><w:pPr><w:jc w:val="right"/></w:pPr><w:r><w:t xml:space="preserve">от PhD докторанта:</w:t></w:r></w:p>
    <w:p><w:pPr><w:jc w:val="right"/></w:pPr><w:r><w:rPr><w:b/></w:rPr><w:t>[%student_full_name%]</w:t></w:r></w:p>
    <w:p><w:pPr><w:jc w:val="right"/></w:pPr><w:r><w:t xml:space="preserve">Специальность: </w:t></w:r><w:r><w:rPr><w:b/></w:rPr><w:t>[%student_specialty%]</w:t></w:r></w:p>
    <w:p><w:pPr><w:jc w:val="right"/></w:pPr><w:r><w:t xml:space="preserve">Тема диссертации: </w:t></w:r><w:r><w:rPr><w:b/></w:rPr><w:t>[%dissertation_topic%]</w:t></w:r></w:p>
    <w:p><w:pPr><w:jc w:val="right"/></w:pPr><w:r><w:t xml:space="preserve">Научный руководитель: </w:t></w:r><w:r><w:rPr><w:b/></w:rPr><w:t>[%student_supervisors%]</w:t></w:r></w:p>
    <w:p><w:pPr><w:spacing w:after="200"/></w:pPr></w:p>
    <w:p><w:pPr><w:jc w:val="center"/></w:pPr><w:r><w:rPr><w:b/></w:rPr><w:t>ЗАЯВЛЕНИЕ</w:t></w:r></w:p>
    <w:p><w:pPr><w:jc w:val="center"/></w:pPr><w:r><w:rPr><w:b/></w:rPr><w:t xml:space="preserve">на защиту диссертации</w:t></w:r></w:p>
    <w:p><w:pPr><w:spacing w:after="200"/></w:pPr></w:p>
    <w:p><w:pPr><w:ind w:firstLine="567"/><w:jc w:val="both"/></w:pPr><w:r><w:t xml:space="preserve">Прошу направить мою диссертацию на защиту в Локальную комиссию по биоэтике в соответствии с установленными правилами.</w:t></w:r></w:p>
    <w:p><w:pPr><w:spacing w:after="200"/></w:pPr></w:p>
    <w:p><w:pPr><w:ind w:firstLine="567"/></w:pPr><w:r><w:rPr><w:b/></w:rPr><w:t>Приложения:</w:t></w:r></w:p>
    <w:p><w:pPr><w:ind w:left="567"/></w:pPr><w:r><w:rPr><w:i/></w:rPr><w:t xml:space="preserve">1. Несброшюрованная диссертация (1 экз.)</w:t></w:r></w:p>
    <w:p><w:pPr><w:ind w:left="567"/></w:pPr><w:r><w:rPr><w:i/></w:rPr><w:t xml:space="preserve">2. Отзывы консультантов</w:t></w:r></w:p>
    <w:p><w:pPr><w:ind w:left="567"/></w:pPr><w:r><w:rPr><w:i/></w:rPr><w:t xml:space="preserve">3. Список публикаций (Приложение 7)</w:t></w:r></w:p>
    <w:p><w:pPr><w:ind w:left="567"/></w:pPr><w:r><w:rPr><w:i/></w:rPr><w:t xml:space="preserve">4. Выписка Ученого совета</w:t></w:r></w:p>
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
  const outputPath = path.resolve(__dirname, "../public/templates/Letter_to_LCB_ru.docx");
  
  // Backup the old one
  const backupPath = path.resolve(__dirname, "../public/templates/Letter_to_LCB_ru.docx.backup");
  if (fs.existsSync(outputPath)) {
    fs.copyFileSync(outputPath, backupPath);
    console.log("Backed up old template to:", backupPath);
  }
  
  fs.writeFileSync(outputPath, buffer);
  console.log("Created clean LCB template (RU) at:", outputPath);
};

const createCleanLCBTemplateKZ = () => {
  const xml = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">
  <w:body>
    <w:p><w:pPr><w:jc w:val="center"/><w:rPr><w:b/></w:rPr></w:pPr><w:r><w:rPr><w:b/></w:rPr><w:t xml:space="preserve">ҚазМҰУ-дың жергілікті биоэтика комиссиясына (ЛБК) </w:t></w:r><w:r><w:rPr><w:b/></w:rPr><w:br/></w:r><w:r><w:rPr><w:b/></w:rPr><w:t>өтініш-хат</w:t></w:r></w:p>
    <w:p><w:pPr><w:jc w:val="center"/><w:rPr><w:b/></w:rPr></w:pPr></w:p>
    <w:p><w:pPr><w:spacing w:after="0"/></w:pPr><w:r><w:rPr><w:b/></w:rPr><w:t>E-mail: </w:t></w:r><w:hyperlink r:id="rId6" w:history="1"><w:r><w:rPr><w:rStyle w:val="aff8"/></w:rPr><w:t>lcb@kaznmu.kz</w:t></w:r></w:hyperlink></w:p>
    <w:p><w:r><w:rPr><w:b/></w:rPr><w:t>Тема: ЛБК-ның «қорғауға» қорытындысын беру туралы өтініш</w:t></w:r></w:p>
    <w:p></w:p>
    <w:p><w:r><w:rPr><w:b/></w:rPr><w:t>Мәтін үлгісі:</w:t></w:r></w:p>
    <w:p><w:t>Құрметті ЛБК мүшелері!</w:t></w:p>
    <w:p><w:t>Диссертациялық жұмысым бойынша ЛБК-ның «қорғауға» қорытындысын берулеріңізді сұраймын.</w:t></w:p>
    <w:p></w:p>
    <w:p><w:t xml:space="preserve">Докторант туралы мәліметтер:</w:t></w:p>
    <w:p><w:t>ТАӘ: [%student_full_name%]</w:t></w:p>
    <w:p><w:t>Мамандығы (дайындық бағыты): [%student_specialty%]</w:t></w:p>
    <w:p><w:t>Диссертация тақырыбы: [%dissertation_topic%]</w:t></w:p>
    <w:p><w:t>Ғылыми жетекші: [%student_supervisors%]</w:t></w:p>
    <w:p></w:p>
    <w:p><w:t>Бұрын берілген ЛБК (алғашқы) қорытындысы:</w:t></w:p>
    <w:p><w:t>Хаттама нөмірі: ____________________   Күні: «[%day%]» [%month%] [%year%] ж.</w:t></w:p>
    <w:p></w:p>
    <w:p><w:t>Материалдарды қарап, ЛБК-ның «қорғауға» қорытындысын берулеріңізді өтінемін.</w:t></w:p>
    <w:p></w:p>
    <w:p><w:t>Қосымшалар:</w:t></w:p>
    <w:p><w:pPr><w:ind w:left="720"/></w:pPr><w:r><w:rPr><w:i/></w:rPr><w:t>1. ЛБК-ның алғашқы қорытындысы (көшірме)</w:t></w:r></w:p>
    <w:p><w:pPr><w:ind w:left="720"/></w:pPr><w:r><w:rPr><w:i/></w:rPr><w:t>2. Диссертация (түптелмеген, 1 дана) / немесе файл сілтемесі</w:t></w:r></w:p>
    <w:p><w:pPr><w:ind w:left="720"/></w:pPr><w:r><w:rPr><w:i/></w:rPr><w:t>3. ЛБК талабы бойынша өзге құжаттар</w:t></w:r></w:p>
    <w:p><w:t>Кері байланыс үшін байланыс деректері: [%student_phone%] / [%student_email%]</w:t></w:p>
    <w:p></w:p>
    <w:p><w:t>Күні: «[%day%]» [%month%] [%year%] ж.</w:t></w:p>
  </w:body>
</w:document>`;

  const zip = new PizZip();
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
  <Relationship Id="rId6" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/hyperlink" Target="mailto:lcb@kaznmu.kz" TargetMode="External"/>
</Relationships>`);

  const buffer = zip.generate({ type: "nodebuffer" });
  const outputPath = path.resolve(__dirname, "../public/templates/Letter_to_LCB_kz.docx");
  
  const backupPath = path.resolve(__dirname, "../public/templates/Letter_to_LCB_kz.docx.backup");
  if (fs.existsSync(outputPath)) {
    fs.copyFileSync(outputPath, backupPath);
    console.log("Backed up old template to:", backupPath);
  }
  
  fs.writeFileSync(outputPath, buffer);
  console.log("Created clean LCB template (KZ) at:", outputPath);
};

const createCleanLCBTemplateEN = () => {
  const xml = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">
  <w:body>
    <w:p><w:pPr><w:jc w:val="center"/><w:rPr><w:b/></w:rPr></w:pPr><w:r><w:rPr><w:b/></w:rPr><w:t xml:space="preserve">Request Letter To the Local Bioethics Committee (LBC) of KazNMU</w:t></w:r></w:p>
    <w:p><w:pPr><w:spacing w:after="0"/></w:pPr><w:r><w:rPr><w:b/></w:rPr><w:t>E-mail: </w:t></w:r><w:hyperlink r:id="rId6" w:history="1"><w:r><w:rPr><w:rStyle w:val="aff8"/></w:rPr><w:t>lcb@kaznmu.kz</w:t></w:r></w:hyperlink></w:p>
    <w:p><w:r><w:rPr><w:b/></w:rPr><w:t>Subject: Request for issuance of LBC conclusion “for defense”</w:t></w:r></w:p>
    <w:p></w:p>
    <w:p><w:r><w:rPr><w:b/></w:rPr><w:t>Text Template:</w:t></w:r></w:p>
    <w:p><w:t>Dear members of the LBC,</w:t></w:p>
    <w:p><w:t>I kindly request the issuance of the LBC conclusion “for defense” for my dissertation.</w:t></w:p>
    <w:p></w:p>
    <w:p><w:t>Doctoral student information:</w:t></w:p>
    <w:p><w:t>Full name: [%student_full_name%]</w:t></w:p>
    <w:p><w:t>Specialty (field of study): [%student_specialty%]</w:t></w:p>
    <w:p><w:t>Dissertation title: [%dissertation_topic%]</w:t></w:p>
    <w:p><w:t>Scientific advisor: [%student_supervisors%]</w:t></w:p>
    <w:p></w:p>
    <w:p><w:t>Previously issued LBC (initial) conclusion:</w:t></w:p>
    <w:p><w:t>Protocol number: ____________________   Date: “[%day%]” [%month%] [%year%]</w:t></w:p>
    <w:p></w:p>
    <w:p><w:t>Please review the materials and issue the LBC conclusion “for defense”.</w:t></w:p>
    <w:p></w:p>
    <w:p><w:t>Attachments:</w:t></w:p>
    <w:p><w:pPr><w:ind w:left="720"/></w:pPr><w:r><w:rPr><w:i/></w:rPr><w:t>1. Initial LBC conclusion (copy)</w:t></w:r></w:p>
    <w:p><w:pPr><w:ind w:left="720"/></w:pPr><w:r><w:rPr><w:i/></w:rPr><w:t>2. Dissertation (unbound, 1 copy) / or link to file</w:t></w:r></w:p>
    <w:p><w:pPr><w:ind w:left="720"/></w:pPr><w:r><w:rPr><w:i/></w:rPr><w:t>3. Other documents as required by LBC</w:t></w:r></w:p>
    <w:p><w:t>Contact for feedback: [%student_phone%] / [%student_email%]</w:t></w:p>
    <w:p></w:p>
    <w:p><w:t>Date: “[%day%]” [%month%] [%year%]</w:t></w:p>
  </w:body>
</w:document>`;

  const zip = new PizZip();
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
  <Relationship Id="rId6" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/hyperlink" Target="mailto:lcb@kaznmu.kz" TargetMode="External"/>
</Relationships>`);

  const buffer = zip.generate({ type: "nodebuffer" });
  const outputPath = path.resolve(__dirname, "../public/templates/Letter_to_LCB_en.docx");
  
  const backupPath = path.resolve(__dirname, "../public/templates/Letter_to_LCB_en.docx.backup");
  if (fs.existsSync(outputPath)) {
    fs.copyFileSync(outputPath, backupPath);
    console.log("Backed up old template to:", backupPath);
  }
  
  fs.writeFileSync(outputPath, buffer);
  console.log("Created clean LCB template (EN) at:", outputPath);
};

createCleanLCBTemplate();
createCleanLCBTemplateKZ();
createCleanLCBTemplateEN();
