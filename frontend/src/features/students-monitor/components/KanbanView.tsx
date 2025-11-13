import React from "react";
import type { MonitorStudent } from "../api";
import { Card, CardHeader, CardTitle, CardContent } from "@/components/ui/card";
import { Avatar, AvatarFallback } from "@/components/ui/avatar";
import { Progress } from "@/components/ui/progress";
import { Badge } from "@/components/ui/badge";
import { Clock, AlertCircle } from "lucide-react";

const STAGES = [
  { id: "W1", label: { EN: "Preparation", RU: "Подготовка", KZ: "Дайындық" } },
  {
    id: "W2",
    label: { EN: "Pre-examination", RU: "Предзащита", KZ: "Алдын ала қорғау" },
  },
  { id: "W3", label: { EN: "RP", RU: "РП", KZ: "ЖЖ" } },
  {
    id: "W4",
    label: { EN: "Submission to DC", RU: "Подача в ДС", KZ: "ДК-ға тапсыру" },
  },
  {
    id: "W5",
    label: { EN: "Restoration", RU: "Доработка", KZ: "Қалпына келтіру" },
  },
  {
    id: "W6",
    label: {
      EN: "After DC acceptance",
      RU: "После ДС",
      KZ: "ДК қабылдауынан кейін",
    },
  },
  {
    id: "W7",
    label: { EN: "Defense & Post-defense", RU: "Защита", KZ: "Қорғау" },
  },
];

type Language = "EN" | "RU" | "KZ";

export function KanbanView({
  rows,
  language = "EN",
}: {
  rows: MonitorStudent[];
  language?: Language;
}) {
  // Bucket students by current stage
  const byStage: Record<string, MonitorStudent[]> = Object.fromEntries(
    STAGES.map((s) => [s.id, []])
  );

  rows.forEach((r) => {
    const stage = r.current_stage || "W1";
    if (byStage[stage]) {
      byStage[stage].push(r);
    }
  });

  // Check if any student has rp_required = true to show W3
  const showW3 = rows.some((r) => r.rp_required);
  const visibleStages = showW3 ? STAGES : STAGES.filter((s) => s.id !== "W3");

  const getStatusVariant = (
    student: MonitorStudent
  ): "default" | "secondary" | "destructive" => {
    if (student.overdue) return "destructive";
    const pct = student.overall_progress_pct || 0;
    if (pct < 30) return "secondary";
    return "default";
  };

  return (
    <div className="overflow-x-auto pb-4">
      <div className="flex gap-4 min-w-max">
        {visibleStages.map((stage) => (
          <Card key={stage.id} className="w-80 flex-shrink-0">
            <CardHeader className="pb-3">
              <div className="flex items-center justify-between">
                <CardTitle className="text-sm font-semibold">
                  {stage.label[language]}
                </CardTitle>
                <Badge variant="outline" className="text-xs">
                  {byStage[stage.id].length}
                </Badge>
              </div>
            </CardHeader>
            <CardContent className="space-y-3 max-h-[600px] overflow-y-auto">
              {byStage[stage.id].length === 0 ? (
                <div className="text-center py-8 text-sm text-muted-foreground">
                  No students in this stage
                </div>
              ) : (
                byStage[stage.id].map((student) => (
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
                                Progress
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

                          {student.due_next && (
                            <div
                              className={`flex items-center gap-1.5 mt-3 text-xs ${
                                student.overdue
                                  ? "text-destructive"
                                  : "text-muted-foreground"
                              }`}
                            >
                              {student.overdue ? (
                                <AlertCircle className="h-3.5 w-3.5" />
                              ) : (
                                <Clock className="h-3.5 w-3.5" />
                              )}
                              <span>
                                {student.overdue ? "Overdue:" : "Due:"}{" "}
                                {student.due_next}
                              </span>
                            </div>
                          )}

                          {student.rp_required && (
                            <Badge variant="outline" className="mt-2 text-xs">
                              RP Required
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
