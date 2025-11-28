import PizZip from "pizzip";
import fs from "fs";
import path from "path";
import { fileURLToPath } from 'url';

const __filename = fileURLToPath(import.meta.url);
const __dirname = path.dirname(__filename);

// Create a clean NCGNT normocontrol letter template
const createCleanNCGNTTemplate = () => {
  const xml = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<w:document xmlns:w="http://schemas.openxmlformats.org/wordprocessingml/2006/main">
  <w:body>
    <w:p><w:pPr><w:pStyle w:val="ac"/><w:rPr><w:b/><w:bCs/><w:color w:val="ADADAD" w:themeColor="background2" w:themeShade="BF"/></w:rPr></w:pPr><w:r><w:rPr><w:b/><w:bCs/><w:color w:val="ADADAD" w:themeColor="background2" w:themeShade="BF"/></w:rPr><w:t>Пример письма в АО «НЦГНТЭ» (для нормоконтроля)</w:t></w:r></w:p>
    <w:p><w:pPr><w:pStyle w:val="ac"/><w:rPr><w:rStyle w:val="ad"/><w:rFonts w:eastAsiaTheme="majorEastAsia"/></w:rPr></w:pPr></w:p>
    <w:p><w:pPr><w:pStyle w:val="ac"/><w:rPr><w:rStyle w:val="ad"/><w:b w:val="0"/><w:bCs w:val="0"/></w:rPr></w:pPr><w:r><w:rPr><w:rStyle w:val="ad"/><w:rFonts w:eastAsiaTheme="majorEastAsia"/></w:rPr><w:t xml:space="preserve">Email: </w:t></w:r><w:hyperlink r:id="rId6" w:history="1"><w:r><w:rPr><w:rStyle w:val="ad"/></w:rPr><w:t>astana@ncste.kz</w:t></w:r></w:hyperlink></w:p>
    <w:p><w:pPr><w:pStyle w:val="ac"/></w:pPr><w:r><w:rPr><w:rStyle w:val="ad"/><w:rFonts w:eastAsiaTheme="majorEastAsia"/><w:color w:val="ADADAD" w:themeColor="background2" w:themeShade="BF"/></w:rPr><w:t>Тема письма:</w:t></w:r><w:r><w:br/><w:t>Запрос на проведение нормоконтроля диссертации</w:t></w:r></w:p>
    <w:p><w:pPr><w:pStyle w:val="ac"/><w:rPr><w:color w:val="ADADAD" w:themeColor="background2" w:themeShade="BF"/></w:rPr></w:pPr><w:r><w:rPr><w:rStyle w:val="ad"/><w:rFonts w:eastAsiaTheme="majorEastAsia"/><w:color w:val="ADADAD" w:themeColor="background2" w:themeShade="BF"/></w:rPr><w:t>Текст письма:</w:t></w:r></w:p>
    <w:p><w:pPr><w:pStyle w:val="ac"/></w:pPr><w:r><w:t>Уважаемые сотрудники АО «Национальный центр государственной научно-технической экспертизы»!</w:t></w:r></w:p>
    <w:p><w:pPr><w:pStyle w:val="ac"/></w:pPr><w:r><w:t xml:space="preserve">Прошу провести </w:t></w:r><w:r><w:rPr><w:rStyle w:val="ad"/><w:rFonts w:eastAsiaTheme="majorEastAsia"/></w:rPr><w:t>нормоконтроль диссертационной работы</w:t></w:r><w:r><w:t xml:space="preserve"> на тему</w:t></w:r><w:r><w:br/></w:r><w:r><w:rPr><w:rStyle w:val="ad"/><w:rFonts w:eastAsiaTheme="majorEastAsia"/></w:rPr><w:t>«{{dissertation_topic}}»</w:t></w:r><w:r><w:br/><w:t xml:space="preserve">по специальности </w:t></w:r><w:r><w:rPr><w:rStyle w:val="ad"/><w:rFonts w:eastAsiaTheme="majorEastAsia"/></w:rPr><w:t>{{student_specialty}}</w:t></w:r><w:r><w:t>.</w:t></w:r></w:p>
    <w:p><w:pPr><w:pStyle w:val="ac"/></w:pPr><w:r><w:t xml:space="preserve">Диссертация подготовлена в формате </w:t></w:r><w:r><w:rPr><w:rStyle w:val="ad"/><w:rFonts w:eastAsiaTheme="majorEastAsia"/></w:rPr><w:t>DOCX</w:t></w:r><w:r><w:t xml:space="preserve"> и прилагается к данному письму.</w:t></w:r><w:r><w:br/><w:t xml:space="preserve">Также прилагаю </w:t></w:r><w:r><w:rPr><w:rStyle w:val="ad"/><w:rFonts w:eastAsiaTheme="majorEastAsia"/></w:rPr><w:t>квитанцию об оплате услуги нормоконтроля</w:t></w:r><w:r><w:t>.</w:t></w:r></w:p>
    <w:p><w:pPr><w:pStyle w:val="ac"/></w:pPr><w:r><w:t>Контактные данные для обратной связи:</w:t></w:r><w:r><w:br/></w:r><w:r><w:rPr><w:rStyle w:val="ad"/><w:rFonts w:eastAsiaTheme="majorEastAsia"/></w:rPr><w:t>Ф.И.О.:</w:t></w:r><w:r><w:t xml:space="preserve"> {{student_full_name}}</w:t></w:r><w:r><w:br/></w:r><w:r><w:rPr><w:rStyle w:val="ad"/><w:rFonts w:eastAsiaTheme="majorEastAsia"/></w:rPr><w:t>Докторантура:</w:t></w:r><w:r><w:t xml:space="preserve"> {{student_program}}</w:t></w:r><w:r><w:br/></w:r><w:r><w:rPr><w:rStyle w:val="ad"/><w:rFonts w:eastAsiaTheme="majorEastAsia"/></w:rPr><w:t>Телефон:</w:t></w:r><w:r><w:t xml:space="preserve"> {{student_phone}}</w:t></w:r><w:r><w:br/></w:r><w:r><w:rPr><w:rStyle w:val="ad"/><w:rFonts w:eastAsiaTheme="majorEastAsia"/></w:rPr><w:t>Электронная почта:</w:t></w:r><w:r><w:t xml:space="preserve"> {{student_email}}</w:t></w:r></w:p>
    <w:p><w:pPr><w:pStyle w:val="ac"/></w:pPr><w:r><w:t>С уважением,</w:t></w:r><w:r><w:br/><w:t>{{student_full_name}}</w:t></w:r><w:r><w:br/><w:t>{{day}} {{month}} {{year}} г.</w:t></w:r></w:p>
  </w:body>
</w:document>`;

  const outputPath = path.resolve(__dirname, "../public/templates/normocontrol_letter.docx");
  
  // Load existing file to preserve styles/themes
  if (!fs.existsSync(outputPath)) {
    console.error("Error: Template file not found at", outputPath);
    return;
  }

  const content = fs.readFileSync(outputPath);
  const zip = new PizZip(content);
  
  // Update ONLY document.xml
  zip.file("word/document.xml", xml);
  
  const buffer = zip.generate({ type: "nodebuffer" });
  fs.writeFileSync(outputPath, buffer);
  console.log("Updated NCGNT template (preserving styles) at:", outputPath);
};

createCleanNCGNTTemplate();
