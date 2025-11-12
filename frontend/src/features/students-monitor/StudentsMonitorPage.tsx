import React from "react";
import { useQuery } from "@tanstack/react-query";
import { fetchMonitorStudents, type MonitorStudent } from "./api";
import { FiltersBar, type Filters } from "./components/FiltersBar";
import { StudentsTableView } from "./components/StudentsTableView";
import { KanbanView } from "./components/KanbanView";
import { AnalyticsView } from "./components/AnalyticsView";
import { Card, CardHeader, CardTitle } from "@/components/ui/card";
import { Tabs } from "@/components/ui/tabs";
import { StudentDetailDrawer } from "./components/StudentDetailDrawer";

export function StudentsMonitorPage() {
  const [filters, setFilters] = React.useState<Filters>({});
  const { data = [], isLoading, refetch } = useQuery<MonitorStudent[]>({
    queryKey: ["monitor", "students", filters],
    queryFn: () => fetchMonitorStudents({
      q: filters.q,
      program: filters.program,
      department: filters.department,
      cohort: filters.cohort,
      rp_required: filters.rp_required ? 1 : undefined,
      overdue: filters.overdue ? 1 : undefined,
      due_from: filters.due_from,
      due_to: filters.due_to,
    }),
  });
  const [tab, setTab] = React.useState<"table"|"kanban"|"analytics">("table");
  const [detail, setDetail] = React.useState<MonitorStudent | null>(null);

  return (
    <div className="space-y-4">
      <div className="flex items-center justify-between">
        <h2 className="text-2xl font-bold">Students Progress</h2>
      </div>
      <FiltersBar value={filters} onChange={setFilters} onRefresh={() => refetch()} />
      <Card>
        <CardHeader>
          <CardTitle className="text-sm">View</CardTitle>
        </CardHeader>
      </Card>
      <div className="flex items-center gap-2 text-sm">
        <button className={`px-3 py-1.5 rounded ${tab==='table'?'bg-muted':''}`} onClick={() => setTab('table')}>Table</button>
        <button className={`px-3 py-1.5 rounded ${tab==='kanban'?'bg-muted':''}`} onClick={() => setTab('kanban')}>Kanban</button>
        <button className={`px-3 py-1.5 rounded ${tab==='analytics'?'bg-muted':''}`} onClick={() => setTab('analytics')}>Analytics</button>
      </div>

      {isLoading ? (
        <div className="text-sm text-muted-foreground">Loading…</div>
      ) : data.length === 0 ? (
        <div className="text-sm text-muted-foreground">No students match your filters.</div>
      ) : tab === 'table' ? (
        <StudentsTableView rows={data} onOpen={(s) => setDetail(s)} />
      ) : tab === 'kanban' ? (
        <KanbanView rows={data} />
      ) : (
        <AnalyticsView />
      )}

      <StudentDetailDrawer open={!!detail} onOpenChange={(b) => !b && setDetail(null)} student={detail ? { id: detail.id, name: detail.name, program: detail.program, department: detail.department, advisors: detail.advisors as any } : null} />
    </div>
  );
}

export default StudentsMonitorPage;
