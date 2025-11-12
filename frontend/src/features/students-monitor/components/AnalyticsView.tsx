import React from "react";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { fetchMonitorAnalytics } from "../api";

export function AnalyticsView({ filters }: { filters: Record<string, any> }) {
  const [data, setData] = React.useState<any>(null);
  const [loading, setLoading] = React.useState(false);
  const [error, setError] = React.useState<string | null>(null);

  React.useEffect(() => {
    setLoading(true); setError(null);
    fetchMonitorAnalytics({
      q: filters.q,
      program: filters.program,
      department: filters.department,
      cohort: filters.cohort,
      advisor_id: filters.advisor_id,
      rp_required: filters.rp_required ? 1 : undefined,
    }).then(setData).catch((e) => setError(String(e))).finally(() => setLoading(false));
  }, [JSON.stringify(filters)]);

  if (loading) return <div className="text-sm text-muted-foreground">Loading analytics…</div>;
  if (error) return <div className="text-sm text-red-600">{error}</div>;

  return (
    <div className="grid grid-cols-1 sm:grid-cols-2 lg:grid-cols-4 gap-3">
      <Card>
        <CardHeader><CardTitle className="text-sm">% with S1_antiplag ≥85% confirmed</CardTitle></CardHeader>
        <CardContent><div className="text-2xl font-bold">{Math.round(data?.antiplag_done_percent || 0)}%</div></CardContent>
      </Card>
      <Card>
        <CardHeader><CardTitle className="text-sm">Median days in W2</CardTitle></CardHeader>
        <CardContent><div className="text-2xl font-bold">{Math.round(data?.w2_median_days || 0)}</div></CardContent>
      </Card>
      <Card>
        <CardHeader><CardTitle className="text-sm">Bottleneck node this month</CardTitle></CardHeader>
        <CardContent>
          <div className="text-sm font-mono">{data?.bottleneck_node_id || '—'}</div>
          <div className="text-xs text-muted-foreground">Count: {data?.bottleneck_count || 0}</div>
        </CardContent>
      </Card>
      <Card>
        <CardHeader><CardTitle className="text-sm">RP required: N students</CardTitle></CardHeader>
        <CardContent><div className="text-2xl font-bold">{data?.rp_required_count || 0}</div></CardContent>
      </Card>
    </div>
  );
}
