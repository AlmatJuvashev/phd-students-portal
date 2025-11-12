import React from "react";
import { Sheet, SheetContent, SheetTitle } from "@/components/ui/sheet";
import { Badge } from "@/components/ui/badge";
import { fetchStudentJourney, type JourneyNode, patchStudentNodeState, fetchDeadlines, putDeadline } from "../api";
import { Button } from "@/components/ui/button";
import { StatusChip } from "./StatusChip";

export function StudentDetailDrawer({ open, onOpenChange, student }: { open: boolean; onOpenChange: (b: boolean) => void; student: { id: string; name: string; program?: string; department?: string; advisors?: { id: string; name: string }[] } | null }) {
  const [nodes, setNodes] = React.useState<JourneyNode[]>([]);
  const [loading, setLoading] = React.useState(false);
  const [deadlines, setDeadlines] = React.useState<Record<string, string>>({});
  React.useEffect(() => {
    if (!open || !student) return;
    setLoading(true);
    Promise.all([
      fetchStudentJourney(student.id).then((d) => setNodes(d.nodes)),
      fetchDeadlines(student.id).then(list => {
        const m: Record<string,string> = {}; list.forEach(it => { m[it.node_id] = it.due_at }); return m;
      }).then(m => setDeadlines(m))
    ]).finally(() => setLoading(false));
  }, [open, student?.id]);

  async function confirm(nodeId: string) {
    if (!student) return;
    await patchStudentNodeState(student.id, nodeId, "done");
    const d = await fetchStudentJourney(student.id);
    setNodes(d.nodes);
  }

  async function setDue(nodeId: string, due: string) {
    if (!student) return;
    await putDeadline(student.id, nodeId, due);
    setDeadlines(prev => ({ ...prev, [nodeId]: due }));
  }

  return (
    <Sheet open={open} onOpenChange={onOpenChange}>
      <SheetContent side="right" className="w-[96vw] sm:w-[800px] lg:w-[900px] p-0 overflow-y-auto">
        <div className="sticky top-0 z-10 h-16 flex items-center px-6 border-b bg-background/95 backdrop-blur">
          <div className="flex-1">
            <SheetTitle className="text-xl">{student?.name}</SheetTitle>
            <div className="text-sm text-muted-foreground mt-0.5">
              {student?.program} · {student?.department}
            </div>
          </div>
        </div>
        <div className="p-6 space-y-6">
          {/* Advisors */}
          <div className="flex flex-wrap gap-2">
            {(student?.advisors || []).map((a) => (
              <Badge key={a.id} variant="secondary">{a.name}</Badge>
            ))}
          </div>
          {/* Nodes */}
          <div className="space-y-3">
            <div className="text-base font-semibold">Stage Checklist</div>
            {loading ? (
              <div className="text-sm text-muted-foreground py-8 text-center">Loading journey data...</div>
            ) : (
              <div className="space-y-3">
                {nodes.map((n) => (
                  <div key={n.node_id} className="border rounded-lg p-4 space-y-3 hover:shadow-md transition-shadow">
                    <div className="flex items-start justify-between gap-4">
                      <div className="flex-1 min-w-0">
                        <div className="font-mono text-sm font-medium mb-1">{n.node_id}</div>
                        <div className="text-xs text-muted-foreground space-y-0.5">
                          <div>Updated: {n.updated_at ? new Date(n.updated_at).toLocaleString() : "—"}</div>
                          <div>Due: {deadlines[n.node_id] ? new Date(deadlines[n.node_id]).toLocaleString() : '—'}</div>
                        </div>
                      </div>
                      <StatusChip state={n.state} />
                    </div>
                    
                    <div className="flex flex-wrap items-center gap-2">
                      <input
                        type="datetime-local"
                        aria-label="Set due date"
                        className="flex-1 min-w-[200px] border rounded-md px-3 py-2 text-sm"
                        value={deadlines[n.node_id] ? deadlines[n.node_id].slice(0,16) : ''}
                        onChange={(e) => setDue(n.node_id, e.target.value)}
                      />
                      <Button size="sm" variant="outline" onClick={() => confirm(n.node_id)}>
                        Mark as Done
                      </Button>
                    </div>
                  </div>
                ))}
                {nodes.length === 0 && (
                  <div className="text-sm text-muted-foreground py-8 text-center">No journey data available.</div>
                )}
              </div>
            )}
          </div>
        </div>
      </SheetContent>
    </Sheet>
  );
}
