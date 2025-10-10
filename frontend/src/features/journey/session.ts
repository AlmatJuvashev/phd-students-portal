const KEY = "journey_state";

export type JourneyStateMap = Record<string, string>;

export function loadJourneyState(): JourneyStateMap | null {
  try {
    const raw = sessionStorage.getItem(KEY);
    if (!raw) return null;
    return JSON.parse(raw) as JourneyStateMap;
  } catch {
    return null;
  }
}

export function saveJourneyState(next: JourneyStateMap) {
  try {
    sessionStorage.setItem(KEY, JSON.stringify(next));
  } catch {}
}

export function patchJourneyState(patch: JourneyStateMap) {
  const cur = loadJourneyState() || {};
  const next = { ...cur, ...patch };
  saveJourneyState(next);
  return next;
}

