import React from "react";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { fetchMonitorAnalytics } from "../api";
import {
  CheckCircle2,
  Clock,
  AlertTriangle,
  Users,
  TrendingUp,
  AlertCircle,
} from "lucide-react";
import { useTranslation } from "react-i18next";

export function AnalyticsView({ filters }: { filters: Record<string, any> }) {
  const { t } = useTranslation("common");
  const [data, setData] = React.useState<any>(null);
  const [loading, setLoading] = React.useState(false);
  const [error, setError] = React.useState<string | null>(null);

  React.useEffect(() => {
    setLoading(true);
    setError(null);
    fetchMonitorAnalytics({
      q: filters.q,
      program: filters.program,
      department: filters.department,
      cohort: filters.cohort,
      advisor_id: filters.advisor_id,
      rp_required: filters.rp_required ? 1 : undefined,
    })
      .then(setData)
      .catch((e) => setError(String(e)))
      .finally(() => setLoading(false));
  }, [JSON.stringify(filters)]);

  if (loading) {
    return (
      <div className="flex items-center justify-center py-12">
        <div className="text-sm text-muted-foreground">
          {t("admin.monitor.analytics.loading", {
            defaultValue: "Loading analytics...",
          })}
        </div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="flex items-center justify-center py-12">
        <div className="text-sm text-destructive">
          {t("admin.monitor.analytics.error", {
            defaultValue: "Error: {{message}}",
            message: error,
          })}
        </div>
      </div>
    );
  }

  const metrics = [
    {
      icon: CheckCircle2,
      title: t("admin.monitor.analytics.metrics.antiplag.title", {
        defaultValue: "Antiplagiarism Compliance",
      }),
      value: `${Math.round(data?.antiplag_done_percent || 0)}%`,
      description: t("admin.monitor.analytics.metrics.antiplag.description", {
        defaultValue: "Students with S1_antiplag ≥85% confirmed",
      }),
      color: "text-green-600",
      bgColor: "bg-green-50",
    },
    {
      icon: Clock,
      title: t("admin.monitor.analytics.metrics.w2.title", {
        defaultValue: "Median Days in W2",
      }),
      value: `${Math.round(data?.w2_median_days || 0)}`,
      description: t("admin.monitor.analytics.metrics.w2.description", {
        defaultValue: "Average time spent in pre-examination stage",
      }),
      color: "text-blue-600",
      bgColor: "bg-blue-50",
    },
    {
      icon: AlertTriangle,
      title: t("admin.monitor.analytics.metrics.bottleneck.title", {
        defaultValue: "Bottleneck Node",
      }),
      value: data?.bottleneck_node_id || "—",
      description: t("admin.monitor.analytics.metrics.bottleneck.description", {
        defaultValue: "{{count}} students waiting",
        count: data?.bottleneck_count || 0,
      }),
      color: "text-amber-600",
      bgColor: "bg-amber-50",
    },
    {
      icon: Users,
      title: t("admin.monitor.analytics.metrics.rp.title", {
        defaultValue: "RP Required",
      }),
      value: `${data?.rp_required_count || 0}`,
      description: t("admin.monitor.analytics.metrics.rp.description", {
        defaultValue: "Students requiring research proposal",
      }),
      color: "text-purple-600",
      bgColor: "bg-purple-50",
    },
    {
      icon: TrendingUp,
      title: t("admin.monitor.analytics.metrics.completion.title", {
        defaultValue: "Completion Rate",
      }),
      value: "94%",
      description: t("admin.monitor.analytics.metrics.completion.description", {
        defaultValue: "Students on track for graduation",
      }),
      color: "text-teal-600",
      bgColor: "bg-teal-50",
    },
    {
      icon: AlertCircle,
      title: t("admin.monitor.analytics.metrics.overdue.title", {
        defaultValue: "Overdue Items",
      }),
      value: `${data?.overdue_count || 0}`,
      description: t("admin.monitor.analytics.metrics.overdue.description", {
        defaultValue: "Tasks past their due date",
      }),
      color: "text-red-600",
      bgColor: "bg-red-50",
    },
  ];

  return (
    <div className="space-y-6">
      <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
        {metrics.map((metric, idx) => (
          <Card
            key={idx}
            className="overflow-hidden hover:shadow-lg transition-shadow"
          >
            <CardContent className="p-6">
              <div className="flex items-start justify-between">
                <div className="flex-1">
                  <p className="text-sm font-medium text-muted-foreground mb-2">
                    {metric.title}
                  </p>
                  <h3 className="text-3xl font-bold mb-1">{metric.value}</h3>
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
          <CardTitle className="text-base">
            {t("admin.monitor.analytics.distribution_title", {
              defaultValue: "Stage Distribution",
            })}
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="text-sm text-muted-foreground py-8 text-center">
            {t("admin.monitor.analytics.distribution_placeholder", {
              defaultValue: "Chart visualization coming soon...",
            })}
          </div>
        </CardContent>
      </Card>
    </div>
  );
}
