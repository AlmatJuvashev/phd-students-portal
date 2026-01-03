import React from "react";
import { useQuery, useMutation } from "@tanstack/react-query";
import { useNavigate } from "react-router-dom";
import { fetchMonitorStudents, runBatchAnalysis, type MonitorStudent } from "./api";
import { FiltersBar, type Filters } from "./components/FiltersBar";
import { StudentsTableView } from "./components/StudentsTableView";
import { KanbanView } from "./components/KanbanView";
import { AnalyticsView } from "./components/AnalyticsView";
import { Card, CardHeader, CardTitle } from "@/components/ui/card";
import { Tabs, TabsList, TabsTrigger } from "@/components/ui/tabs";
import { Modal } from "@/components/ui/modal";
import { Input } from "@/components/ui/input";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { useTranslation } from "react-i18next";
import { Activity, Play } from "lucide-react";

// Simple toast fallback since sonner is missing
const toast = {
  success: (msg: string) => alert(msg),
  error: (msg: string) => alert(msg)
};

const BatchAnalysisButton = ({ onComplete }: { onComplete: () => void }) => {
  const { mutate, isPending } = useMutation({
    mutationFn: runBatchAnalysis,
    onSuccess: () => {
      toast.success("Risk analysis started successfully");
      onComplete();
    },
    onError: () => {
      toast.error("Failed to start risk analysis");
    }
  });

  return (
    <Button 
      variant="outline" 
      onClick={() => mutate()} 
      disabled={isPending}
      className="gap-2"
    >
      {isPending ? <Activity className="animate-spin h-4 w-4" /> : <Play className="h-4 w-4" />}
      Run Analysis
    </Button>
  );
};

export function StudentsMonitorPage() {
  const { t } = useTranslation("common");
  const [filters, setFilters] = React.useState<Filters>({});
  const {
    data = [],
    isLoading,
    refetch,
  } = useQuery<MonitorStudent[]>({
    queryKey: ["monitor", "students", filters],
    queryFn: () =>
      fetchMonitorStudents({
        q: filters.q,
        program: filters.program,
        department: filters.department,
        cohort: filters.cohort,
        rp_required: filters.rp_required ? 1 : undefined,
        due_from: filters.due_from,
        due_to: filters.due_to,
      }),
  });
  const [tab, setTab] = React.useState<"table" | "kanban" | "analytics">(
    "table"
  );
  const [selected, setSelected] = React.useState<Set<string>>(new Set());
  const [bulkOpen, setBulkOpen] = React.useState(false);
  const [bulkTitle, setBulkTitle] = React.useState("");
  const [bulkMessage, setBulkMessage] = React.useState("");
  const [bulkDue, setBulkDue] = React.useState("");
  const navigate = useNavigate();

  // CSV export: listen to event from FiltersBar and export current rows
  React.useEffect(() => {
    function onExport(e: any) {
      const rows = data || [];
      const head = [
        "id",
        "name",
        "email",
        "phone",
        "program",
        "department",
        "cohort",
        "current_stage",
        "stage_done",
        "stage_total",
        "overall_progress_pct",
        "last_update",
      ];
      const lines = [head.join(",")].concat(
        rows.map((r) =>
          [
            r.id,
            r.name,
            r.email || "",
            r.phone || "",
            r.program || "",
            r.department || "",
            r.cohort || "",
            r.current_stage || "",
            String(r.stage_done || ""),
            String(r.stage_total || ""),
            String(Math.round(r.overall_progress_pct || 0)),
            r.last_update || "",
          ]
            .map((v) => `"${String(v).replace(/"/g, '""')}"`)
            .join(",")
        )
      );
      const blob = new Blob([lines.join("\n")], {
        type: "text/csv;charset=utf-8;",
      });
      const url = URL.createObjectURL(blob);
      const a = document.createElement("a");
      a.href = url;
      a.download = `students-monitor.csv`;
      a.click();
      URL.revokeObjectURL(url);
    }
    window.addEventListener("students-monitor:export", onExport as any);
    return () =>
      window.removeEventListener("students-monitor:export", onExport as any);
  }, [data]);

  // Bulk reminder open handler
  React.useEffect(() => {
    function onOpenBulk() {
      setBulkOpen(true);
    }
    window.addEventListener("students-monitor:bulk-reminder", onOpenBulk);
    return () =>
      window.removeEventListener("students-monitor:bulk-reminder", onOpenBulk);
  }, []);

  React.useEffect(() => {
    setSelected((prev) => {
      const next = new Set<string>();
      data.forEach((student) => {
        if (prev.has(student.id)) {
          next.add(student.id);
        }
      });
      return next;
    });
  }, [data]);

  async function sendBulkReminder() {
    const ids =
      selected.size > 0 ? Array.from(selected) : (data || []).map((r) => r.id);
    if (ids.length === 0) {
      setBulkOpen(false);
      return;
    }
    await (
      await import("./api")
    ).postReminders({
      student_ids: ids,
      title: bulkTitle,
      message: bulkMessage,
      due_at: bulkDue || undefined,
    });
    setBulkOpen(false);
    setBulkTitle("");
    setBulkMessage("");
    setBulkDue("");
  }

  return (
    <div className="min-h-screen bg-background">
      {/* Header */}
      <div className="sticky top-0 z-50 bg-background/95 backdrop-blur supports-[backdrop-filter]:bg-background/80 border-b">
        <div className="px-8 py-6">
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-4">
              <h1 className="text-3xl font-bold tracking-tight">
                {t("admin.monitor.title", { defaultValue: "Student Monitoring" })}
              </h1>
              <Badge
                variant="secondary"
                className="h-6 px-2.5 text-sm font-medium"
              >
                {t("admin.monitor.count_label", {
                  defaultValue: "{{count}} students",
                  count: data.length,
                })}
              </Badge>
            </div>
            <div className="flex gap-2">
               <BatchAnalysisButton onComplete={() => refetch()} />
               <Button onClick={() => setBulkOpen(true)}>
                 {t("admin.monitor.bulk_action", "Bulk Action")}
               </Button>
            </div>
          </div>
        </div>
      </div>

      <FiltersBar
        value={filters}
        onChange={setFilters}
        onRefresh={() => refetch()}
      />

      {/* Tabs */}
      <div className="px-8 py-4 border-b">
        <Tabs value={tab} onValueChange={(v) => setTab(v as any)}>
          <TabsList className="grid w-full max-w-md grid-cols-3">
            <TabsTrigger value="table">
              {t("admin.monitor.tabs.table", { defaultValue: "Table View" })}
            </TabsTrigger>
            <TabsTrigger value="kanban">
              {t("admin.monitor.tabs.kanban", { defaultValue: "Kanban Board" })}
            </TabsTrigger>
            <TabsTrigger value="analytics">
              {t("admin.monitor.tabs.analytics", { defaultValue: "Analytics" })}
            </TabsTrigger>
          </TabsList>
        </Tabs>
      </div>

      {/* Content */}
      <div className="px-8 py-6">
        {isLoading ? (
          <div className="text-center py-12 text-muted-foreground">
            {t("admin.monitor.loading", {
              defaultValue: "Loading students...",
            })}
          </div>
        ) : data.length === 0 ? (
          <div className="text-center py-12 text-muted-foreground">
            {t("admin.monitor.empty", {
              defaultValue: "No students match your filters.",
            })}
          </div>
        ) : tab === "table" ? (
          <StudentsTableView
            rows={data}
            selected={selected}
            onToggle={(id, checked) =>
              setSelected((prev) => {
                const next = new Set(prev);
                if (checked) next.add(id);
                else next.delete(id);
                return next;
              })
            }
            onToggleAll={(checked) =>
              setSelected(checked ? new Set(data.map((d) => d.id)) : new Set())
            }
          />
        ) : tab === "kanban" ? (
          <KanbanView rows={data} />
        ) : (
          <AnalyticsView filters={filters} />
        )}
      </div>

      <Modal open={bulkOpen} onClose={() => setBulkOpen(false)}>
        <div className="space-y-3">
          <div className="text-sm font-semibold">
            {t("admin.monitor.bulk.title", {
              defaultValue: "New reminder for {{count}} students",
              count: data.length || 0,
            })}
          </div>
          <Input
            placeholder={t("admin.monitor.bulk.title_placeholder", {
              defaultValue: "Title",
            })}
            value={bulkTitle}
            onChange={(e) => setBulkTitle(e.target.value)}
          />
          <Input
            placeholder={t("admin.monitor.bulk.message_placeholder", {
              defaultValue: "Message (optional)",
            })}
            value={bulkMessage}
            onChange={(e) => setBulkMessage(e.target.value)}
          />
          <input
            type="datetime-local"
            value={bulkDue}
            onChange={(e) => setBulkDue(e.target.value)}
            className="border rounded px-2 py-1 w-full"
            aria-label={t("admin.monitor.bulk.due_label", {
              defaultValue: "Due date",
            })}
          />
          <div className="flex justify-end gap-2">
            <Button variant="ghost" onClick={() => setBulkOpen(false)}>
              {t("common.cancel", { defaultValue: "Cancel" })}
            </Button>
            <Button onClick={sendBulkReminder} disabled={!bulkTitle}>
              {t("admin.monitor.bulk.send", { defaultValue: "Send" })}
            </Button>
          </div>
        </div>
      </Modal>
    </div>
  );
}

export default StudentsMonitorPage;
