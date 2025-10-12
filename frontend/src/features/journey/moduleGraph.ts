import type { Playbook, NodeDef } from "@/lib/playbook";

export type WorldInfo = {
  id: string;
  nodes: NodeDef[];
};

export function getWorlds(pb: Playbook): WorldInfo[] {
  return pb.worlds.map((w) => ({ id: w.id, nodes: w.nodes }));
}

// Heuristic: a terminal node is one whose next[] is empty or points to a node in another world
export function detectTerminalNodeIds(pb: Playbook): Record<string, string[]> {
  const worldByNode: Record<string, string> = {};
  for (const w of pb.worlds) {
    for (const n of w.nodes) worldByNode[n.id] = w.id;
  }
  const terminals: Record<string, string[]> = {};
  for (const w of pb.worlds) {
    const list: string[] = [];
    for (const n of w.nodes) {
      const next = n.next || [];
      if (!next || next.length === 0) {
        list.push(n.id);
        continue;
      }
      const leavesWorld = next.some((id) => worldByNode[id] && worldByNode[id] !== w.id);
      if (leavesWorld) list.push(n.id);
    }
    terminals[w.id] = list;
  }
  return terminals;
}

export function firstNodeIdOfNextWorld(pb: Playbook, worldId: string): string | null {
  const idx = pb.worlds.findIndex((w) => w.id === worldId);
  if (idx < 0) return null;
  const next = pb.worlds[idx + 1];
  if (!next || !next.nodes.length) return null;
  return next.nodes[0]?.id ?? null;
}

