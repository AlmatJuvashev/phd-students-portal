// lib/playbook.ts
export type RoleId = "student" | "advisor" | "secretary" | "chair" | "admin";

export type ActionKind = "form" | "outcome" | "wait" | "gateway";

export type FieldDef = {
  key: string;
  required?: boolean;
  type?: string;
  label?: Record<string, string>;
  placeholder?: Record<string, string>;
  options?: Array<{ value: string; label?: Record<string, string> }>;
  other_key?: string; // for select with 'other'
};

export type UploadDef = {
  key: string;
  mime?: string[];
  required?: boolean;
  label?: Record<string, string>;
  accept?: string;
};

export type Playbook = {
  playbook_id: string;
  version: string;
  worlds: Array<{
    id: string;
    title: Record<string, string>;
    order: number;
    nodes: NodeDef[];
  }>;
  roles?: Array<{ id: RoleId; label: Record<string, string> }>;
  conditions?: Array<{ id: string; expr: string }>;
};

export type NodeDef = {
  id: string;
  title: Record<string, string>;
  type:
    | "form"
    | "upload"
    | "decision"
    | "meeting"
    | "waiting"
    | "external"
    | "boss"
    | "gateway"
    | "info"
    | "confirmTask"; // custom simple confirm step
  who_can_complete: RoleId[];
  prerequisites?: string[];
  next?: string[];
  outcomes?: Array<{
    value: string;
    label?: Record<string, string>;
    next: string[];
  }>;
  condition?: string; // like "rp_required"
  timer?: { duration_days: number; start_on: string };
  requirements?: {
    fields?: FieldDef[];
    uploads?: UploadDef[];
    validations?: Array<{ rule: string; source?: string }>;
    notes?: string;
    checklist?: string[]; // external sub-steps
  };
  outputs?: Array<{ key: string; type: "upload" | "auto_generated" }>;

  // NEW (optional) — explicit classification override
  actionHints?: ActionKind[]; // e.g., ["outcome","upload"] for hearing with minutes
};

// Simple helper to pick RU title by default
import i18n from "i18next";

function humanizeKey(key: string) {
  if (!key) return "";
  const s = key
    .replace(/[_-]+/g, " ")
    .replace(/([a-z])([A-Z])/g, "$1 $2")
    .toLowerCase();
  return s.replace(/(^|\s)\S/g, (c) => c.toUpperCase());
}

export const t = (obj?: Record<string, string>, fallback = "") => {
  const lang = (i18n?.language as "ru" | "kz" | "en") || "ru";
  const val = obj?.[lang] ?? obj?.en ?? obj?.ru ?? obj?.kz;
  if (val) return val;
  // Try i18n dictionary fallback for field keys: fields.<key>
  if (fallback) {
    const key = `fields.${fallback}`;
    const dictVal = (i18n as any)?.t?.(key);
    if (dictVal && dictVal !== key) return dictVal as string;
  }
  // Humanize technical fallback keys like "full_name"
  return humanizeKey(fallback);
};

// safeText: accepts string | string[] | locale-map of strings or string[]
// Returns a safe string for rendering, picking current locale when applicable.
export function safeText(
  value: unknown,
  fallback = "",
  joiner: string = "\n"
): string {
  const lang = (i18n?.language as "ru" | "kz" | "en") || "ru";
  if (value == null) return fallback;
  if (typeof value === "string") return value;
  if (Array.isArray(value)) {
    return value.filter((x) => typeof x === "string").join(joiner) || fallback;
  }
  if (typeof value === "object") {
    const obj = value as Record<string, any>;
    const localized = obj[lang] ?? obj.en ?? obj.ru ?? obj.kz;
    if (typeof localized === "string") return localized;
    if (Array.isArray(localized))
      return localized.filter((x) => typeof x === "string").join(joiner);
  }
  return fallback;
}

// Build a fast lookup
export function indexPlaybook(pb: Playbook) {
  const nodeById = new Map<string, NodeDef>();
  pb.worlds.forEach((w) => w.nodes.forEach((n) => nodeById.set(n.id, n)));
  return { nodeById };
}

// Basic progress per world (done/submitted states would come from server; here mocked)
export type NodeState =
  | "locked"
  | "active"
  | "submitted"
  | "waiting"
  | "needs_fixes"
  | "done";
export type NodeVM = NodeDef & { worldId: string; state: NodeState };

/**
 * computeNodeStates — applies prerequisite-based unlocking logic.
 * Only unlocks nodes whose prerequisites are all "done" AND that have no runtime conditions,
 * OR nodes with conditions that are satisfied by the provided activeConditions set.
 */
export function computeNodeStates(
  pb: Playbook,
  rawStateByNodeId: Record<string, NodeState> = {},
  activeConditions: Set<string> = new Set()
): Record<string, NodeState> {
  let hasChanges = false;
  const computed: Record<string, NodeState> = {};

  // Copy all existing states
  for (const key in rawStateByNodeId) {
    computed[key] = rawStateByNodeId[key];
  }

  // Helper to check if all prerequisites are done
  const allPrereqsDone = (prereqs: string[] | undefined): boolean => {
    // Empty array means "no prerequisites required" - can be activated
    if (prereqs && prereqs.length === 0) return true;
    // Undefined or missing prerequisites means node is not ready to be auto-activated
    if (!prereqs) return false;
    // Check if all prerequisites are completed
    return prereqs.every((prereqId) => computed[prereqId] === "done");
  };

  // Iterate through all nodes and compute states
  let changed = true;
  let iterations = 0;
  const maxIterations = 100; // safety limit

  while (changed && iterations < maxIterations) {
    changed = false;
    iterations++;

    pb.worlds.forEach((world) => {
      world.nodes.forEach((node) => {
        const currentState = computed[node.id] ?? "locked";

        // Only process locked nodes
        if (currentState === "locked") {
          // Check if node has a condition
          if (node.condition) {
            // Only activate if condition is satisfied AND prerequisites are met
            if (
              activeConditions.has(node.condition) &&
              allPrereqsDone(node.prerequisites)
            ) {
              computed[node.id] = "active";
              changed = true;
              hasChanges = true;
            }
            return;
          }

          // Check if all prerequisites are done
          if (allPrereqsDone(node.prerequisites)) {
            computed[node.id] = "active";
            changed = true;
            hasChanges = true;
          }
        }
      });
    });
  }

  // Return original object if nothing changed to maintain reference equality
  return hasChanges ? computed : rawStateByNodeId;
}

export function toViewModel(
  pb: Playbook,
  stateByNodeId: Record<string, NodeState> = {},
  activeConditions: Set<string> = new Set()
): { worlds: Array<{ id: string; title: string; nodes: NodeVM[] }> } {
  // Compute states with prerequisite logic before building view model
  const computedStates = computeNodeStates(pb, stateByNodeId, activeConditions);

  const worlds = [...pb.worlds]
    .sort((a, b) => a.order - b.order)
    .map((w) => ({
      id: w.id,
      title: t(w.title, w.id),
      nodes: w.nodes.map((n) => ({
        ...n,
        worldId: w.id,
        state: computedStates[n.id] ?? "locked",
      })),
    }));
  return { worlds };
}

// Compute simple edges: prerequisites → node, plus outcomes branches
export function edgesForWorld(pb: Playbook, worldId: string) {
  const w = pb.worlds.find((x) => x.id === worldId);
  if (!w) return [];
  const edges: Array<{
    from: string;
    to: string;
    kind: "default" | "conditional" | "outcome";
  }> = [];
  const ids = new Set(w.nodes.map((n) => n.id));
  w.nodes.forEach((n) => {
    n.prerequisites?.forEach((pr) => {
      if (ids.has(pr)) edges.push({ from: pr, to: n.id, kind: "default" });
    });
    n.outcomes?.forEach((o) =>
      o.next.forEach((nx) => {
        if (ids.has(nx)) edges.push({ from: n.id, to: nx, kind: "outcome" });
      })
    );
    if (n.condition && n.next)
      n.next.forEach((nx) => {
        if (ids.has(nx))
          edges.push({ from: n.id, to: nx, kind: "conditional" });
      });
  });
  return edges;
}

// -------- Action detector (switcher input) ----------
export function detectActionKinds(n: NodeDef): ActionKind[] {
  // Explicit read-only informational nodes should render like a gateway/info screen
  if (n.type === "info") return ["gateway"]; // no actions, just read-only content

  if (n.actionHints?.length) return n.actionHints;

  const kinds: ActionKind[] = [];

  const hasFields = !!n.requirements?.fields?.length;
  const hasOutcomes =
    !!n.outcomes?.length ||
    n.type === "decision" ||
    n.type === "meeting" ||
    n.type === "boss";
  const hasTimer = !!n.timer;
  const isWaiting = n.type === "waiting";
  const isGateway = n.type === "gateway";

  if (hasFields) kinds.push("form");
  if (hasOutcomes) kinds.push("outcome");
  if (hasTimer || isWaiting) kinds.push("wait");
  if (isGateway) kinds.push("gateway");

  return kinds.length ? kinds : ["gateway"]; // default read-only
}
