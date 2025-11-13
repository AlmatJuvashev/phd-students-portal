export function computeTimeAgo(targetISO?: string, now: Date = new Date()) {
  if (!targetISO) return null;
  const target = new Date(targetISO);
  if (isNaN(target.getTime())) return null;
  const diffMs = now.getTime() - target.getTime();
  const sec = Math.max(1, Math.round(Math.abs(diffMs) / 1000));
  const units: Array<[string, number]> = [
    ["years", 31536000],
    ["months", 2592000],
    ["weeks", 604800],
    ["days", 86400],
    ["hours", 3600],
    ["minutes", 60],
    ["seconds", 1],
  ];
  for (const [unit, u] of units) {
    if (sec >= u) {
      const value = Math.floor(sec / u);
      return { value, unit } as { value: number; unit: string };
    }
  }
  return { value: 1, unit: "seconds" };
}

