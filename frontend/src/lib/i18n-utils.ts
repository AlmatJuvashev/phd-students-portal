export const getLocalized = (jsonString: string | undefined | null, lang: string): string => {
  if (!jsonString) return '';
  try {
    // Attempt to parse as JSON
    const parsed = JSON.parse(jsonString);
    if (typeof parsed === 'object' && parsed !== null) {
      // Logic: exact match -> fallback language (en) -> first available key
      return parsed[lang] || parsed['en'] || parsed['ru'] || parsed['kz'] || Object.values(parsed)[0] || '';
    }
    return String(parsed);
  } catch (e) {
    // If not JSON, return as is
    return jsonString;
  }
};
