import React from "react";
import type { MonitorStudent } from "../api";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";

const STAGES = ["W1","W2","W3","W4","W5","W6","W7"];

export function KanbanView({ rows }: { rows: MonitorStudent[] }) {
  // naive stage assignment: use rp_required to include W3 and bucket by rough progress
  const byStage: Record<string, MonitorStudent[]> = Object.fromEntries(STAGES.map(s => [s, []]));
  rows.forEach(r => {
    const pct = Math.round(r.overall_progress_pct || 0);
    let stage = "W1";
    if (pct >= 85) stage = "W7"; else if (pct >= 70) stage = "W6"; else if (pct >= 55) stage = "W5"; else if (pct >= 40) stage = "W4"; else if (pct >= 25) stage = "W3"; else if (pct >= 10) stage = "W2";
    if (stage === "W3" && r.rp_required === false) stage = "W4"; // skip W3 when not required
    byStage[stage].push(r);
  });
  return (
    <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-7 gap-3">
      {STAGES.map(s => (
        <Card key={s} className="min-h-[280px]">
          <CardHeader><CardTitle className="text-sm">{s}</CardTitle></CardHeader>
          <CardContent className="space-y-2">
            {byStage[s].map(r => (
              <div key={r.id} className="border rounded p-2">
                <div className="text-sm font-medium truncate">{r.name}</div>
                <div className="text-xs text-muted-foreground truncate">{[r.program, r.department].filter(Boolean).join(" · ")}</div>
                <div className="mt-1 bg-muted/40 rounded-full h-1 overflow-hidden">
                  <div className="bg-primary h-1" style={{ width: `${Math.round(r.overall_progress_pct || 0)}%` }} />
                </div>
              </div>
            ))}
            {byStage[s].length === 0 && (
              <div className="text-xs text-muted-foreground">No students</div>
            )}
          </CardContent>
        </Card>
      ))}
    </div>
  );
}

