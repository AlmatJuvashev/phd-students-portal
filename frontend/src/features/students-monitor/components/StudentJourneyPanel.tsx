import React from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { fetchStudentJourney, type JourneyNode, fetchDeadlines, putDeadline, patchStudentNodeState } from "../api";
import { StatusChip } from "./StatusChip";

type Props = {
  studentId: string;
  currentStage?: string; // W1..W7; if provided, show only nodes in this stage
  studentName?: string;
};

export function StudentJourneyPanel({ studentId, currentStage }: Props) {
  const [nodes, setNodes] = React.useState<JourneyNode[]>([]);
  const [deadlines, setDeadlines] = React.useState<Record<string, string>>({});
  const [loading, setLoading] = React.useState(false);
  const [stageNodeIds, setStageNodeIds] = React.useState<string[] | null>(null);

  React.useEffect(() => {
    if (!studentId) {
      setNodes([]);
      setDeadlines({});
      return;
    }
    setLoading(true);
    Promise.all([
      fetchStudentJourney(studentId).then((data) => setNodes(data.nodes)),
      fetchDeadlines(studentId).then((list) => {
        const map: Record<string, string> = {};
        list.forEach((item) => { map[item.node_id] = item.due_at; });
        setDeadlines(map);
      }),
    ]).finally(() => setLoading(false));
  }, [studentId]);

  // Load playbook to resolve which node ids belong to the current stage (world)
  React.useEffect(() => {
    let mounted = true;
    if (!currentStage) {
      setStageNodeIds(null);
      return;
    }
    import("@/playbooks/playbook.json").then((mod: any) => {
      if (!mounted) return;
      const pb = (mod && (mod.default || mod)) as any;
      const world = (pb.worlds || pb.Worlds || []).find((w: any) => w.id === currentStage || w.ID === currentStage);
      if (world) {
        const ids = (world.nodes || world.Nodes || []).map((n: any) => n.id || n.ID);
        setStageNodeIds(ids);
      } else {
        setStageNodeIds(null);
      }
    });
    return () => { mounted = false };
  }, [currentStage]);

  const confirmNode = async (nodeId: string) => {
    await patchStudentNodeState(studentId, nodeId, "done");
    const { nodes } = await fetchStudentJourney(studentId);
    setNodes(nodes);
  };

  const updateDeadline = async (nodeId: string, value: string) => {
    await putDeadline(studentId, nodeId, value);
    setDeadlines((prev) => ({ ...prev, [nodeId]: value }));
  };

  return (
    <Card>
      <CardHeader>
        <CardTitle className="text-sm">Stage Checklist</CardTitle>
      </CardHeader>
      <CardContent className="space-y-2">
        {loading && <p className="text-sm text-muted-foreground">Loading journey…</p>}
        {!loading && filteredNodes(nodes, stageNodeIds).length === 0 && (
          <p className="text-sm text-muted-foreground">No journey data yet.</p>
        )}
        {!loading && filteredNodes(nodes, stageNodeIds).map((node) => (
          <div key={node.node_id} className="border rounded p-3 flex items-start justify-between gap-3">
            <div>
              <div className="font-mono text-xs">{node.node_id}</div>
              <div className="text-xs text-muted-foreground">Updated: {node.updated_at ? new Date(node.updated_at).toLocaleString() : "—"}</div>
              {deadlines[node.node_id] && (
                <div className="text-xs text-muted-foreground">Due: {new Date(deadlines[node.node_id]).toLocaleString()}</div>
              )}
            </div>
            <div className="flex flex-col items-end gap-2">
              <StatusChip state={node.state} />
              <div className="flex items-center gap-1">
                <input
                  className="border rounded px-2 py-1 text-xs"
                  type="datetime-local"
                  value={deadlines[node.node_id] ? deadlines[node.node_id].slice(0, 16) : ""}
                  onChange={(e) => updateDeadline(node.node_id, e.target.value)}
                />
                <Button size="xs" onClick={() => confirmNode(node.node_id)}>Confirm</Button>
              </div>
            </div>
          </div>
        ))}
      </CardContent>
    </Card>
  );
}

export default StudentJourneyPanel;

function filteredNodes(all: JourneyNode[], stageNodeIds: string[] | null) {
  if (!stageNodeIds || stageNodeIds.length === 0) return all;
  const set = new Set(stageNodeIds);
  return all.filter((n) => set.has(n.node_id));
}
