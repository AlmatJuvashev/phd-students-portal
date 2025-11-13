import React from "react";
import { useTranslation } from "react-i18next";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Progress } from "@/components/ui/progress";
import { Checkbox } from "@/components/ui/checkbox";
import { AlertCircle } from "lucide-react";
import type { MonitorStudent } from "../api";
import { stageLabel } from "../utils";

export function StudentsTableView({
  rows,
  onView,
  selected,
  onToggle,
  onToggleAll,
}: {
  rows: MonitorStudent[];
  onView: (s: MonitorStudent) => void;
  selected: Set<string>;
  onToggle: (id: string, checked: boolean) => void;
  onToggleAll: (checked: boolean) => void;
}) {
  const { i18n } = useTranslation("common");

  const getRowClass = (student: MonitorStudent) => {
    if (student.overdue) return "bg-red-50/50 hover:bg-red-50";
    // Add more status colors based on progress or other criteria
    return "hover:bg-muted/50";
  };

  return (
    <Card className="border-border shadow-sm">
      <CardContent className="p-0">
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead className="border-b bg-muted/30">
              <tr>
                <th className="p-4 text-left">
                  <Checkbox
                    checked={selected.size > 0 && selected.size === rows.length}
                    onCheckedChange={(checked) => onToggleAll(!!checked)}
                    aria-label="Select all"
                  />
                </th>
                <th className="p-4 text-left text-sm font-medium text-muted-foreground">
                  Student
                </th>
                <th className="p-4 text-left text-sm font-medium text-muted-foreground">
                  Program · Department
                </th>
                <th className="p-4 text-left text-sm font-medium text-muted-foreground">
                  Advisors · Cohort
                </th>
                <th className="p-4 text-left text-sm font-medium text-muted-foreground">
                  Current Stage
                </th>
                <th className="p-4 text-left text-sm font-medium text-muted-foreground">
                  Progress
                </th>
                <th className="p-4 text-left text-sm font-medium text-muted-foreground">
                  Due Next
                </th>
                <th className="p-4 text-left text-sm font-medium text-muted-foreground">
                  Last Update
                </th>
                <th className="p-4 text-right text-sm font-medium text-muted-foreground">
                  Actions
                </th>
              </tr>
            </thead>
            <tbody>
              {rows.map((r) => (
                  <tr
                    key={r.id}
                    className={`border-b transition-colors cursor-pointer ${getRowClass(
                      r
                    )}`}
                    onClick={() => onView(r)}
                  >
                  <td className="p-4" onClick={(e) => e.stopPropagation()}>
                    <Checkbox
                      checked={selected.has(r.id)}
                      onCheckedChange={(checked) => onToggle(r.id, !!checked)}
                      aria-label={`Select ${r.name}`}
                    />
                  </td>
                  <td className="p-4">
                    <div className="flex items-center gap-3">
                      <Avatar className="h-10 w-10 border-2 border-border">
                        <AvatarFallback className="bg-primary/10 text-primary font-medium">
                          {r.name
                            .split(" ")
                            .map((n) => n[0])
                            .join("")
                            .slice(0, 2)}
                        </AvatarFallback>
                      </Avatar>
                      <div>
                        <div className="font-medium text-foreground">
                          {r.name}
                        </div>
                        <div className="text-sm text-muted-foreground">
                          {r.email || r.phone || r.id}
                        </div>
                      </div>
                    </div>
                  </td>
                  <td className="p-4">
                    <div className="text-sm text-foreground">
                      {r.program || "—"}
                    </div>
                    <div className="text-sm text-muted-foreground">
                      {r.department || "—"}
                    </div>
                  </td>
                  <td className="p-4">
                    <div className="flex flex-wrap gap-1.5 mb-1.5">
                      {(r.advisors || []).map((advisor, idx) => (
                        <Badge
                          key={idx}
                          variant="outline"
                          className="text-xs bg-muted/20"
                        >
                          {advisor.name}
                        </Badge>
                      ))}
                      {(!r.advisors || r.advisors.length === 0) && (
                        <span className="text-xs text-muted-foreground">
                          No advisor
                        </span>
                      )}
                    </div>
                    <div className="text-sm text-muted-foreground">
                      {r.cohort || "—"}
                    </div>
                  </td>
                  <td className="p-4">
                      <Badge className="bg-primary/10 text-primary border border-primary/20">
                        {stageLabel(r.current_stage, i18n.language)}
                      </Badge>
                    {r.rp_required && (
                      <Badge
                        variant="outline"
                        className="ml-1.5 text-xs bg-amber-50 text-amber-700 border-amber-200"
                      >
                        RP
                      </Badge>
                    )}
                    {typeof r.stage_done === "number" &&
                      typeof r.stage_total === "number" && (
                        <div className="text-xs text-muted-foreground mt-1.5">
                          {r.stage_done}/{r.stage_total} nodes done
                        </div>
                      )}
                  </td>
                  <td className="p-4">
                    <div className="flex items-center gap-2 mb-1">
                      <Progress
                        value={r.overall_progress_pct || 0}
                        className="h-2 flex-1 max-w-[120px]"
                      />
                      <span className="text-sm font-medium text-foreground min-w-[38px] text-right">
                        {Math.round(r.overall_progress_pct || 0)}%
                      </span>
                    </div>
                  </td>
                  <td className="p-4">
                    <div
                      className={`flex items-center gap-1.5 text-sm ${
                        r.overdue
                          ? "text-red-600 font-medium"
                          : "text-foreground"
                      }`}
                    >
                      {r.overdue && <AlertCircle className="h-4 w-4" />}
                      {r.due_next || "—"}
                    </div>
                  </td>
                  <td className="p-4">
                    <div className="text-sm text-muted-foreground">
                      {r.last_update
                        ? new Date(r.last_update).toLocaleString(i18n.language)
                        : "—"}
                    </div>
                  </td>
                  <td className="p-4" onClick={(e) => e.stopPropagation()}>
                    <div className="flex items-center justify-end gap-1">
                      <Button
                        variant="ghost"
                        size="sm"
                        className="h-8 px-2 hover:bg-primary/10 hover:text-primary"
                        onClick={() => onView(r)}
                      >
                        View
                      </Button>
                      <Button
                        variant="ghost"
                        size="sm"
                        className="h-8 px-2 hover:bg-primary/10 hover:text-primary"
                      >
                        Notify
                      </Button>
                    </div>
                  </td>
                </tr>
              ))}
              {rows.length === 0 && (
                <tr>
                  <td
                    colSpan={9}
                    className="p-8 text-center text-muted-foreground"
                  >
                    No students match your filters.
                  </td>
                </tr>
              )}
            </tbody>
          </table>
        </div>
      </CardContent>
    </Card>
  );
}
