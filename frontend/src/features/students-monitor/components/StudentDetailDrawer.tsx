import React from "react";
import { Sheet, SheetContent, SheetTitle } from "@/components/ui/sheet";
import { Badge } from "@/components/ui/badge";
import { fetchStudentJourney, type JourneyNode, patchStudentNodeState } from "../api";
import { Button } from "@/components/ui/button";
import { StatusChip } from "./StatusChip";

export function StudentDetailDrawer({ open, onOpenChange, student }: { open: boolean; onOpenChange: (b: boolean) => void; student: { id: string; name: string; program?: string; department?: string; advisors?: { id: string; name: string }[] } | null }) {
  const [nodes, setNodes] = React.useState<JourneyNode[]>([]);
  const [loading, setLoading] = React.useState(false);
  React.useEffect(() => {
    if (!open || !student) return;
    setLoading(true);
    fetchStudentJourney(student.id).then((d) => setNodes(d.nodes)).finally(() => setLoading(false));
  }, [open, student?.id]);

  async function confirm(nodeId: string) {
    if (!student) return;
    await patchStudentNodeState(student.id, nodeId, "done");
    const d = await fetchStudentJourney(student.id);
    setNodes(d.nodes);
  }

  return (
    <Sheet open={open} onOpenChange={onOpenChange}>
      <SheetContent side="right" className="w-[96vw] sm:w-[640px] p-0">
        <div className="h-14 flex items-center px-4 border-b">
          <SheetTitle className="flex-1">{student?.name}</SheetTitle>
          <div className="text-xs text-muted-foreground">
            {student?.program} · {student?.department}
          </div>
        </div>
        <div className="p-4 space-y-4">
          {/* Advisors */}
          <div className="flex flex-wrap gap-2">
            {(student?.advisors || []).map((a) => (
              <Badge key={a.id} variant="secondary">{a.name}</Badge>
            ))}
          </div>
          {/* Nodes */}
          <div className="space-y-2">
            <div className="text-sm font-medium">Stage checklist</div>
            {loading ? (
              <div className="text-sm text-muted-foreground">Loading…</div>
            ) : (
              <div className="space-y-2">
                {nodes.map((n) => (
                  <div key={n.node_id} className="border rounded p-2 flex items-center justify-between">
                    <div>
                      <div className="font-mono text-xs">{n.node_id}</div>
                      <div className="text-xs text-muted-foreground">Updated: {n.updated_at ? new Date(n.updated_at).toLocaleString() : "—"}</div>
                    </div>
                    <div className="flex items-center gap-2">
                      <StatusChip state={n.state} />
                      <Button size="sm" variant="outline" onClick={() => confirm(n.node_id)}>Confirm</Button>
                    </div>
                  </div>
                ))}
                {nodes.length === 0 && <div className="text-sm text-muted-foreground">No journey data.</div>}
              </div>
            )}
          </div>
        </div>
      </SheetContent>
    </Sheet>
  );
}

