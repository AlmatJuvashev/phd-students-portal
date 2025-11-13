import React from "react";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { fetchStudentJourney, type JourneyNode, fetchDeadlines, putDeadline, patchStudentNodeState } from "../api";
import { StatusChip } from "./StatusChip";

type Props = {
  studentId: string;
  studentName?: string;
};

export function StudentJourneyPanel({ studentId }: Props) {
  const [nodes, setNodes] = React.useState<JourneyNode[]>([]);
  const [deadlines, setDeadlines] = React.useState<Record<string, string>>({});
  const [loading, setLoading] = React.useState(false);

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
        {!loading && nodes.length === 0 && (
          <p className="text-sm text-muted-foreground">No journey data yet.</p>
        )}
        {!loading && nodes.map((node) => (
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
