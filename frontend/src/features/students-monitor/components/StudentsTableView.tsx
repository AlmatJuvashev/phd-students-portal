import React from "react";
import { useNavigate } from "react-router-dom";
import { useTranslation } from "react-i18next";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Checkbox } from "@/components/ui/checkbox";
import { AlertCircle } from "lucide-react";
import type { MonitorStudent } from "../api";
import { stageLabel } from "../utils";

const STAGES = ["W1", "W2", "W3", "W4", "W5", "W6", "W7"];

type StageStatus = "completed" | "in-progress" | "at-risk" | "pending";

function getStageHistory(
  student: MonitorStudent
): Array<{ stage: string; status: StageStatus }> {
  const currentStageIndex = STAGES.indexOf(student.current_stage || "W1");

  return STAGES.map((stage, idx) => {
    // Skip W3 if RP not required
    if (stage === "W3" && !student.rp_required) {
      return null;
    }

    let status: StageStatus = "pending";

    if (idx < currentStageIndex) {
      status = "completed";
    } else if (idx === currentStageIndex) {
      // Current stage - check if at risk
      status = student.overdue ? "at-risk" : "in-progress";
    }

    return { stage, status };
  }).filter(Boolean) as Array<{ stage: string; status: StageStatus }>;
}

export function StudentsTableView({
  rows,
  selected,
  onToggle,
  onToggleAll,
}: {
  rows: MonitorStudent[];
  selected: Set<string>;
  onToggle: (id: string, checked: boolean) => void;
  onToggleAll: (checked: boolean) => void;
}) {
  const navigate = useNavigate();
  const { i18n } = useTranslation("common");

  const getRowClass = (student: MonitorStudent, index: number) => {
    const baseClass = index % 2 === 0 ? "bg-white" : "bg-muted/5";
    if (student.overdue) return `${baseClass} hover:bg-red-50/50`;
    return `${baseClass} hover:bg-muted/30`;
  };

  const handleRowClick = (studentId: string) => {
    navigate(`/admin/students-monitor/${studentId}`);
  };

  return (
    <Card className="border-border shadow-sm">
      <CardContent className="p-0">
        <div className="overflow-x-auto">
          <table className="w-full">
            <thead className="border-b bg-muted/30">
              <tr>
                <th className="py-3 px-4 text-left text-xs font-semibold text-muted-foreground uppercase tracking-wide">
                  <Checkbox
                    checked={selected.size > 0 && selected.size === rows.length}
                    onCheckedChange={(checked) => onToggleAll(!!checked)}
                    aria-label="Select all"
                  />
                </th>
                <th className="py-3 px-4 text-left text-xs font-semibold text-muted-foreground uppercase tracking-wide">
                  Student
                </th>
                <th className="py-3 px-4 text-left text-xs font-semibold text-muted-foreground uppercase tracking-wide">
                  Program · Department
                </th>
                <th className="py-3 px-4 text-left text-xs font-semibold text-muted-foreground uppercase tracking-wide">
                  Stage
                </th>
                <th className="py-3 px-4 text-left text-xs font-semibold text-muted-foreground uppercase tracking-wide">
                  Progress
                </th>
                <th className="py-3 px-4 text-center text-xs font-semibold text-muted-foreground uppercase tracking-wide">
                  Cohort
                </th>
                <th className="py-3 px-4 text-center text-xs font-semibold text-muted-foreground uppercase tracking-wide">
                  Due
                </th>
              </tr>
            </thead>
            <tbody>
              {rows.map((student, index) => (
                <tr
                  key={student.id}
                  className={`border-b border-border/50 cursor-pointer transition-colors ${getRowClass(
                    student,
                    index
                  )}`}
                  onClick={() => handleRowClick(student.id)}
                >
                  <td
                    className="py-2.5 px-4"
                    onClick={(e) => e.stopPropagation()}
                  >
                    <Checkbox
                      checked={selected.has(student.id)}
                      onCheckedChange={(checked) =>
                        onToggle(student.id, !!checked)
                      }
                      aria-label={`Select ${student.name}`}
                    />
                  </td>
                  <td className="py-2.5 px-4">
                    <div className="flex items-center gap-3">
                      <Avatar className="h-8 w-8 border border-border">
                        <AvatarFallback className="bg-primary/10 text-primary text-xs font-medium">
                          {student.name
                            .split(" ")
                            .map((n) => n[0])
                            .join("")
                            .slice(0, 2)}
                        </AvatarFallback>
                      </Avatar>
                      <div>
                        <div className="font-medium text-sm text-foreground">
                          {student.name}
                        </div>
                        <div className="text-xs text-muted-foreground">
                          {student.id}
                        </div>
                      </div>
                    </div>
                  </td>
                  <td className="py-2.5 px-4">
                    <div className="text-sm text-foreground">
                      {student.program?.split(" ")[0] || "—"}
                    </div>
                    <div className="text-xs text-muted-foreground">
                      {student.department || "—"}
                    </div>
                  </td>
                  <td className="py-2.5 px-4">
                    <div className="text-sm font-medium text-primary">
                      {student.current_stage || "—"}
                    </div>
                    <div className="text-xs text-muted-foreground">
                      {typeof student.stage_done === "number" &&
                      typeof student.stage_total === "number"
                        ? `${student.stage_done}/${student.stage_total} nodes`
                        : "—"}
                    </div>
                  </td>
                  <td className="py-2.5 px-4">
                    <div className="flex items-center gap-1.5">
                      {getStageHistory(student).map((stage, idx) => (
                        <div
                          key={idx}
                          className={`h-5 w-7 rounded-sm flex items-center justify-center text-xs font-medium transition-colors ${
                            stage.status === "completed"
                              ? "bg-green-600 text-white"
                              : stage.status === "in-progress"
                              ? "bg-primary text-primary-foreground"
                              : stage.status === "at-risk"
                              ? "bg-red-600 text-white"
                              : "bg-muted/30 text-muted-foreground"
                          }`}
                          title={`${stage.stage}: ${stage.status}`}
                        >
                          {stage.stage.replace("W", "")}
                        </div>
                      ))}
                      <span className="ml-2 text-sm font-semibold text-foreground">
                        {Math.round(student.overall_progress_pct || 0)}%
                      </span>
                    </div>
                  </td>
                  <td className="py-2.5 px-4 text-center">
                    <div className="text-sm text-foreground">
                      {student.cohort || "—"}
                    </div>
                  </td>
                  <td className="py-2.5 px-4 text-center">
                    <div
                      className={`text-sm font-medium ${
                        student.overdue ? "text-red-600" : "text-foreground"
                      }`}
                    >
                      {student.due_next || "—"}
                    </div>
                    {student.overdue && (
                      <div className="text-xs text-red-600 flex items-center justify-center gap-1">
                        <AlertCircle className="h-3 w-3" />
                        Overdue
                      </div>
                    )}
                  </td>
                </tr>
              ))}
              {rows.length === 0 && (
                <tr>
                  <td
                    colSpan={7}
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
