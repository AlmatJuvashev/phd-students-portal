import React from "react";
import { useQuery } from "@tanstack/react-query";
import { api } from "@/api/client";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { TooltipProvider, Tooltip, TooltipTrigger, TooltipContent } from "@/components/ui/tooltip";
import { useTranslation } from "react-i18next";

type Row = {
  id: string;
  name: string;
  email: string;
  role: string;
  progress: {
    completed_nodes: number;
    total_nodes: number;
    percent: number;
    current_node_id?: string;
    last_submission_at?: string;
  };
};

export function StudentProgress() {
  const { t } = useTranslation("common");
  const { data = [], isLoading, refetch } = useQuery<Row[]>({
    queryKey: ["admin", "student-progress"],
    queryFn: () => api("/admin/student-progress"),
  });
  const [q, setQ] = React.useState("");

  const rows = React.useMemo(() => {
    if (!q) return data;
    const s = q.toLowerCase();
    return data.filter((r) => r.name.toLowerCase().includes(s) || r.email.toLowerCase().includes(s));
  }, [q, data]);

  return (
    <div className="max-w-6xl mx-auto space-y-6">
      <div className="flex items-center justify-between gap-2 flex-wrap">
        <div>
          <h2 className="text-2xl font-bold">{t('admin.student_progress.title', 'Student Progress')}</h2>
          <p className="text-muted-foreground">{t('admin.student_progress.subtitle', 'Overview of students’ dissertation map completion.')}</p>
        </div>
        <div className="flex items-center gap-2">
          <Input placeholder={t('admin.student_progress.search', 'Search by name or email')} value={q} onChange={(e) => setQ(e.target.value)} />
          <Button variant="outline" onClick={() => refetch()}>{t('admin.student_progress.refresh', 'Refresh')}</Button>
        </div>
      </div>

      <Card>
        <CardHeader>
          <CardTitle>{t('admin.student_progress.students', 'Students')} ({rows.length})</CardTitle>
        </CardHeader>
        <CardContent>
          {isLoading ? (
            <div className="text-sm text-muted-foreground">{t('common.loading', 'Loading…')}</div>
          ) : rows.length === 0 ? (
            <div className="text-sm text-muted-foreground">{t('admin.student_progress.empty', 'No students found.')}</div>
          ) : (
            <div className="overflow-x-auto">
              <table className="w-full text-sm">
                <thead>
                  <tr className="border-b text-left">
                    <th className="py-2 px-3">{t('table.name', 'Name')}</th>
                    <th className="py-2 px-3">{t('table.email', 'Email')}</th>
                    <th className="py-2 px-3">{t('map.progress', 'Progress')}</th>
                    <th className="py-2 px-3">{t('table.current_node', 'Current Node')}</th>
                    <th className="py-2 px-3">{t('table.last_activity', 'Last Activity')}</th>
                    <th className="py-2 px-3">{t('table.actions', 'Actions')}</th>
                  </tr>
                </thead>
                <tbody>
                  {rows.map((r) => {
                    const pct = Math.round(r.progress.percent || 0);
                    const last = r.progress.last_submission_at
                      ? new Date(r.progress.last_submission_at)
                      : null;
                    return (
                      <tr key={r.id} className="border-b">
                        <td className="py-2 px-3 font-medium">{r.name}</td>
                        <td className="py-2 px-3">{r.email}</td>
                        <td className="py-2 px-3">
                          <TooltipProvider>
                            <Tooltip>
                              <TooltipTrigger asChild>
                                <div className="flex items-center gap-2 min-w-[180px] cursor-help">
                                  <div className="flex-1 bg-muted/40 rounded-full h-2 overflow-hidden">
                                    <div
                                      className="bg-primary h-2"
                                      style={{ width: `${Math.min(100, Math.max(0, pct))}%` }}
                                    />
                                  </div>
                                  <span className="tabular-nums w-10 text-right">{pct}%</span>
                                </div>
                              </TooltipTrigger>
                              <TooltipContent>
                                {t('admin.student_progress.tooltip_progress', '{{done}} of {{total}} nodes done ({{pct}}%)', { done: r.progress.completed_nodes, total: r.progress.total_nodes, pct })}
                              </TooltipContent>
                            </Tooltip>
                          </TooltipProvider>
                        </td>
                        <td className="py-2 px-3 font-mono text-xs">
                          {r.progress.current_node_id || "—"}
                        </td>
                        <td className="py-2 px-3">
                          {last ? last.toLocaleString() : "—"}
                        </td>
                        <td className="py-2 px-3">
                          <Button size="sm" variant="outline" onClick={() => (window.location.href = "/journey")}>{t('admin.student_progress.view_journey', 'View Journey')}</Button>
                        </td>
                      </tr>
                    );
                  })}
                </tbody>
              </table>
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}

export default StudentProgress;
