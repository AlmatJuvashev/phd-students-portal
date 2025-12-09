import JSZip from "jszip";
import { saveAs } from "file-saver";
import { generateStudentTemplateDoc } from "./student-template";
import type { PublicAsset, StudentTemplateData, Locale } from "@/types";

/**
 * Downloads all templates for a node as a single ZIP file
 * 
 * @param assets - Array of template assets to generate
 * @param templateData - Student profile data for template filling
 * @param locale - Current locale (ru/kz/en)
 * @param nodeId - Node identifier for filename
 */
export async function downloadAllTemplatesAsZIP(
  assets: PublicAsset[],
  templateData: StudentTemplateData,
  locale: Locale,
  nodeId: string
): Promise<void> {
  const zip = new JSZip();
  const errors: string[] = [];
  let successCount = 0;

  // Generate each template and add to ZIP
  for (const asset of assets) {
    try {
      console.log(`[zip-templates] Generating ${asset.id}...`);
      
      // Call existing template generation function
      // This returns nothing (downloads directly), so we need to modify approach
      // We'll need to extract the blob generation logic
      const blob = await generateTemplateBlob(asset, templateData, locale);
      
      if (blob) {
        const filename = sanitizeFilename(asset, locale);
        zip.file(filename, blob);
        successCount++;
        console.log(`[zip-templates] Added ${filename} to ZIP`);
      }
    } catch (err) {
      const errorMsg = err instanceof Error ? err.message : String(err);
      console.error(`[zip-templates] Failed to generate ${asset.id}:`, errorMsg);
      errors.push(`${asset.id}: ${errorMsg}`);
    }
  }

  // Show errors if any templates failed
  if (errors.length > 0) {
    console.warn(`[zip-templates] ${errors.length} template(s) failed:`, errors);
    // Could show toast notification here
  }

  if (successCount === 0) {
    throw new Error("Failed to generate any templates");
  }

  // Generate ZIP file
  console.log(`[zip-templates] Generating ZIP with ${successCount} files...`);
  const zipBlob = await zip.generateAsync({ 
    type: "blob",
    compression: "DEFLATE",
    compressionOptions: { level: 6 }
  });

  // Download ZIP
  const timestamp = new Date().toISOString().slice(0, 10).replace(/-/g, "");
  const filename = `templates_${nodeId}_${timestamp}.zip`;
  saveAs(zipBlob, filename);
  
  console.log(`[zip-templates] Downloaded ${filename} (${successCount} files)`);
}

/**
 * Generate template blob without triggering download
 * This is a modified version of generateStudentTemplateDoc that returns the blob
 */
async function generateTemplateBlob(
  asset: PublicAsset,
  data: StudentTemplateData,
  locale: Locale
): Promise<Blob | null> {
  const { generateStudentTemplateDocBlob } = await import("./student-template");
  return generateStudentTemplateDocBlob({ asset, data, locale });
}

/**
 * Sanitize asset title to create valid filename
 * 
 * @param asset - Template asset
 * @param locale - Current locale
 * @returns Clean filename with .docx extension
 */
function sanitizeFilename(asset: PublicAsset, locale: Locale): string {
  // Get title from asset
  const title = asset.title?.[locale] || asset.id;
  
  // Remove special characters, keep alphanumeric and spaces
  const clean = title
    .replace(/[^a-zA-Zа-яА-ЯёЁ0-9\s]/g, "_")
    .replace(/\s+/g, "_")
    .replace(/_+/g, "_") // Collapse multiple underscores
    .replace(/^_|_$/g, "") // Trim underscores from ends
    .substring(0, 100); // Limit length
  
  return `${clean}.docx`;
}
