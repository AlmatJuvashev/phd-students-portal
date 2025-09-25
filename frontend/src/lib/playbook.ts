// lib/playbook.ts
export type RoleId = "student" | "advisor" | "secretary" | "chair" | "admin";

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
    | "gateway";
  who_can_complete: RoleId[];
  prerequisites?: string[];
  next?: string[];
  outcomes?: Array<{ value: string; next: string[] }>;
  condition?: string; // like "rp_required"
  timer?: { duration_days: number; start_on: string };
  requirements?: {
    fields?: Array<{ key: string; required?: boolean; type?: string }>;
    uploads?: Array<{ key: string; mime?: string[]; required?: boolean }>;
    validations?: Array<{ rule: string; source?: string }>;
    notes?: string;
  };
  outputs?: Array<{ key: string; type: "upload" | "auto_generated" }>;
};

// Simple helper to pick RU title by default
export const t = (obj?: Record<string, string>, fallback = "") =>
  obj?.ru ?? obj?.en ?? obj?.kz ?? fallback;

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

export function toViewModel(
  pb: Playbook,
  stateByNodeId: Record<string, NodeState> = {}
): { worlds: Array<{ id: string; title: string; nodes: NodeVM[] }> } {
  const worlds = [...pb.worlds]
    .sort((a, b) => a.order - b.order)
    .map((w) => ({
      id: w.id,
      title: t(w.title, w.id),
      nodes: w.nodes.map((n) => ({
        ...n,
        worldId: w.id,
        state: stateByNodeId[n.id] ?? "locked",
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
