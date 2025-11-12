import React from "react";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";

export function AnalyticsView() {
  return (
    <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-3">
      <Card>
        <CardHeader><CardTitle className="text-sm">% with S1_antiplag ≥85% confirmed</CardTitle></CardHeader>
        <CardContent><div className="text-2xl font-bold">—</div></CardContent>
      </Card>
      <Card>
        <CardHeader><CardTitle className="text-sm">Median days in W2</CardTitle></CardHeader>
        <CardContent><div className="text-2xl font-bold">—</div></CardContent>
      </Card>
      <Card>
        <CardHeader><CardTitle className="text-sm">Bottleneck node this month</CardTitle></CardHeader>
        <CardContent><div className="text-sm text-muted-foreground">Coming soon</div></CardContent>
      </Card>
      <Card>
        <CardHeader><CardTitle className="text-sm">RP required: N students</CardTitle></CardHeader>
        <CardContent><div className="text-2xl font-bold">—</div></CardContent>
      </Card>
    </div>
  );
}

