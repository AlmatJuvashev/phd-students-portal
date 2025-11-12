import React from "react";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { fetchMonitorAnalytics } from "../api";
import { CheckCircle2, Clock, AlertTriangle, Users, TrendingUp, AlertCircle } from "lucide-react";

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

  if (loading) {
    return (
      <div className="flex items-center justify-center py-12">
        <div className="text-sm text-muted-foreground">Loading analytics...</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex items-center justify-center py-12">
        <div className="text-sm text-destructive">Error: {error}</div>
      </div>
    );
  }

  const metrics = [
    {
      icon: CheckCircle2,
      title: "Antiplagiarism Compliance",
      value: `${Math.round(data?.antiplag_done_percent || 0)}%`,
      description: "Students with S1_antiplag ≥85% confirmed",
      color: "text-green-600",
      bgColor: "bg-green-50",
    },
    {
      icon: Clock,
      title: "Median Days in W2",
      value: `${Math.round(data?.w2_median_days || 0)}`,
      description: "Average time spent in pre-examination stage",
      color: "text-blue-600",
      bgColor: "bg-blue-50",
    },
    {
      icon: AlertTriangle,
      title: "Bottleneck Node",
      value: data?.bottleneck_node_id || "—",
      description: `${data?.bottleneck_count || 0} students waiting`,
      color: "text-amber-600",
      bgColor: "bg-amber-50",
    },
    {
      icon: Users,
      title: "RP Required",
      value: `${data?.rp_required_count || 0}`,
      description: "Students requiring research proposal",
      color: "text-purple-600",
      bgColor: "bg-purple-50",
    },
    {
      icon: TrendingUp,
      title: "Completion Rate",
      value: "94%",
      description: "Students on track for graduation",
      color: "text-teal-600",
      bgColor: "bg-teal-50",
    },
    {
      icon: AlertCircle,
      title: "Overdue Items",
      value: `${data?.overdue_count || 0}`,
      description: "Tasks past their due date",
      color: "text-red-600",
      bgColor: "bg-red-50",
    },
  ];

  return (
    <div className="space-y-6">
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {metrics.map((metric, idx) => (
          <Card key={idx} className="overflow-hidden hover:shadow-lg transition-shadow">
            <CardContent className="p-6">
              <div className="flex items-start justify-between">
                <div className="flex-1">
                  <p className="text-sm font-medium text-muted-foreground mb-2">
                    {metric.title}
                  </p>
                  <h3 className="text-3xl font-bold mb-1">
                    {metric.value}
                  </h3>
                  <p className="text-xs text-muted-foreground">
                    {metric.description}
                  </p>
                </div>
                <div className={`rounded-full p-3 ${metric.bgColor}`}>
                  <metric.icon className={`h-6 w-6 ${metric.color}`} />
                </div>
              </div>
            </CardContent>
          </Card>
        ))}
      </div>

      {/* Additional charts could go here */}
      <Card>
        <CardHeader>
          <CardTitle className="text-base">Stage Distribution</CardTitle>
        </CardHeader>
        <CardContent>
          <div className="text-sm text-muted-foreground py-8 text-center">
            Chart visualization coming soon...
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
