import type { NodeVM } from "@/lib/playbook";

export type UINodeKind = "formEntry" | "checklist" | "cards" | "info";

export function deriveNodeKind(node: NodeVM): UINodeKind {
  // Explicit info node
  if (node.type === "info") return "info";

  const fields = (node.requirements as any)?.fields as Array<any> | undefined;
  const ui = (node.requirements as any)?.ui_hints || {};

  // Cards layout explicitly requested
  if (ui?.cards_layout) return "cards";

  // Heuristic checklist: mostly boolean fields and at least one required
  if (Array.isArray(fields) && fields.length > 0) {
    const bools = fields.filter((f) => f?.type === "boolean");
    const requiredBools = bools.filter((f) => !!f?.required);
    if (bools.length >= Math.max(1, Math.floor(fields.length * 0.6)) && requiredBools.length > 0) {
      return "checklist";
    }
  }

  return "formEntry";
}

