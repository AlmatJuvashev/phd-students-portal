import React from "react";
import type { MonitorStudent } from "../api";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Progress } from "@/components/ui/progress";
import { Badge } from "@/components/ui/badge";
import { Clock, AlertCircle } from "lucide-react";
import { useTranslation } from "react-i18next";
import { stageLabel } from "../utils";

const STAGES = ["W1", "W2", "W3", "W4", "W5", "W6", "W7"];

export function KanbanView({ rows }: { rows: MonitorStudent[] }) {
  const { t, i18n } = useTranslation("common");
  const language = i18n.language || "en";
  // Bucket students by current stage
  const byStage: Record<string, MonitorStudent[]> = Object.fromEntries(
    STAGES.map((s) => [s, []])
  );

  rows.forEach((r) => {
    const stage = r.current_stage || "W1";
    if (byStage[stage]) {
      byStage[stage].push(r);
    }
  });

  // Check if any student has rp_required = true to show W3
  const showW3 = rows.some((r) => r.rp_required);
  const visibleStages = showW3 ? STAGES : STAGES.filter((s) => s !== "W3");

  const getBadgeVariant = (
  student: MonitorStudent
): "default" | "secondary" | "outline" | "destructive" => {
  if (student.overall_progress_pct >= 80) return "default";
  if (student.overall_progress_pct >= 50) return "secondary";
  return "outline";
};

  return (
    <div className="overflow-x-auto pb-4">
      <div className="flex gap-4 min-w-max">
        {visibleStages.map((stageId) => (
          <Card key={stageId} className="w-80 flex-shrink-0">
            <CardHeader className="pb-3">
              <div className="flex items-center justify-between">
                <CardTitle className="text-sm font-semibold">
                  {stageLabel(stageId, language)}
                </CardTitle>
                <Badge variant="outline" className="text-xs">
                  {byStage[stageId].length}
                </Badge>
              </div>
            </CardHeader>
            <CardContent className="space-y-3 max-h-[600px] overflow-y-auto">
              {byStage[stageId].length === 0 ? (
                <div className="text-center py-8 text-sm text-muted-foreground">
                  {t("admin.monitor.kanban.empty_stage", {
                    defaultValue: "No students in this stage",
                  })}
                </div>
              ) : (
                byStage[stageId].map((student) => (
                  <Card
                    key={student.id}
                    className="border shadow-sm hover:shadow-md transition-shadow cursor-pointer"
                  >
                    <CardContent className="p-4">
                      <div className="flex items-start gap-3">
                        <Avatar className="h-10 w-10">
                          <AvatarFallback className="bg-primary/10 text-primary text-sm font-medium">
                            {student.name
                              .split(" ")
                              .map((n) => n[0])
                              .join("")
                              .slice(0, 2)}
                          </AvatarFallback>
                        </Avatar>
                        <div className="flex-1 min-w-0">
                          <div className="font-medium text-sm truncate">
                            {student.name}
                          </div>
                          <div className="text-xs text-muted-foreground truncate mt-0.5">
                            {[student.program, student.department]
                              .filter(Boolean)
                              .join(" · ")}
                          </div>

                          {student.advisors && student.advisors.length > 0 && (
                            <div className="flex flex-wrap gap-1 mt-2">
                              {student.advisors.map((adv, idx) => (
                                <Badge
                                  key={idx}
                                  variant="secondary"
                                  className="text-xs px-2 py-0.5"
                                >
                                  {adv.name}
                                </Badge>
                              ))}
                            </div>
                          )}

                          <div className="mt-3 space-y-1.5">
                            <div className="flex items-center justify-between text-xs">
                              <span className="text-muted-foreground">
                                {t("admin.monitor.kanban.progress", {
                                  defaultValue: "Progress",
                                })}
                              </span>
                              <span className="font-medium">
                                {Math.round(student.overall_progress_pct || 0)}%
                              </span>
                            </div>
                            <Progress
                              value={student.overall_progress_pct || 0}
                              className="h-1.5"
                            />
                          </div>

                          {student.rp_required && (
                            <Badge variant="outline" className="mt-2 text-xs">
                              {t("admin.monitor.kanban.rp_required", {
                                defaultValue: "RP Required",
                              })}
                            </Badge>
                          )}
                        </div>
                      </div>
                    </CardContent>
                  </Card>
                ))
              )}
            </CardContent>
          </Card>
        ))}
      </div>
    </div>
  );
}
