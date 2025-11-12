import React from "react";
import { useTranslation } from "react-i18next";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import { Badge } from "@/components/ui/badge";
import type { MonitorStudent } from "../api";

export function StudentsTableView({ rows, onOpen }: { rows: MonitorStudent[]; onOpen: (s: MonitorStudent) => void }) {
  const { i18n } = useTranslation('common');
  return (
    <Card>
      <CardContent className="p-0">
        <div className="overflow-x-auto">
          <table className="w-full text-sm">
            <thead className="border-b bg-muted/30">
              <tr className="text-left">
                <th className="py-2 px-3">Student</th>
                <th className="py-2 px-3">Program · Department · Cohort</th>
                <th className="py-2 px-3">Stage</th>
                <th className="py-2 px-3">Advisors</th>
                <th className="py-2 px-3">Overall</th>
                <th className="py-2 px-3">Due next</th>
                <th className="py-2 px-3">Last update</th>
                <th className="py-2 px-3">Actions</th>
              </tr>
            </thead>
            <tbody>
              {rows.map((r) => (
                <tr key={r.id} className="border-b hover:bg-muted/20">
                  <td className="py-2 px-3">
                    <div className="font-medium">{r.name}</div>
                    <div className="text-xs text-muted-foreground">{r.email || r.phone || "—"}</div>
                  </td>
                  <td className="py-2 px-3">
                    <div>{[r.program, r.department].filter(Boolean).join(" · ") || "—"}</div>
                    <div className="text-xs text-muted-foreground">{r.cohort || "—"}</div>
                  </td>
                  <td className="py-2 px-3">
                    <div className="text-xs font-medium">{stageLabel(r.current_stage, i18n.language)}</div>
                    {typeof r.stage_done === 'number' && typeof r.stage_total === 'number' && (
                      <div className="flex items-center gap-2 min-w-[120px]">
                        <div className="flex-1 bg-muted/40 rounded-full h-1 overflow-hidden">
                          <div className="bg-primary h-1" style={{ width: `${Math.round(((r.stage_done||0)/(r.stage_total||1))*100)}%` }} />
                        </div>
                        <span className="text-xs tabular-nums">{r.stage_done}/{r.stage_total}</span>
                      </div>
                    )}
                  </td>
                  <td className="py-2 px-3">
                    <div className="flex flex-wrap gap-1">
                      {(r.advisors || []).map((a) => (
                        <Badge key={a.id} variant="secondary">{a.name}</Badge>
                      ))}
                      {(!r.advisors || r.advisors.length === 0) && <span className="text-xs text-muted-foreground">—</span>}
                    </div>
                  </td>
                  <td className="py-2 px-3">
                    <div className="flex items-center gap-2 min-w-[140px]">
                      <div className="flex-1 bg-muted/40 rounded-full h-2 overflow-hidden">
                        <div className="bg-primary h-2" style={{ width: `${Math.round(r.overall_progress_pct || 0)}%` }} />
                      </div>
                      <span className="tabular-nums w-10 text-right">{Math.round(r.overall_progress_pct || 0)}%</span>
                    </div>
                  </td>
                  <td className="py-2 px-3">
                    <div className="flex items-center gap-2">
                      <span>{r.due_next ? new Date(r.due_next).toLocaleDateString() : '—'}</span>
                      {r.overdue ? <span title="Overdue" className="text-red-600">●</span> : null}
                    </div>
                  </td>
                  <td className="py-2 px-3">{r.last_update ? new Date(r.last_update).toLocaleString() : "—"}</td>
                  <td className="py-2 px-3">
                    <Button size="sm" variant="outline" onClick={() => onOpen(r)}>View</Button>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </CardContent>
    </Card>
  );
}

function stageLabel(s?: string, lang: string = 'en') {
  const en: Record<string,string> = {
    W1: "I — Preparation",
    W2: "II — Pre-examination",
    W3: "III — RP",
    W4: "IV — Submission to DC",
    W5: "V — Restoration",
    W6: "VI — After DC acceptance",
    W7: "VII — Defense & Post-defense",
  };
  const ru: Record<string,string> = {
    W1: "I — Подготовка",
    W2: "II — Предварительная экспертиза",
    W3: "III — RP (условно)",
    W4: "IV — Подача в ДС",
    W5: "V — Восстановление",
    W6: "VI — После принятия ДС",
    W7: "VII — Защита и После защиты",
  };
  const kz: Record<string,string> = {
    W1: "I — Дайындық",
    W2: "II — Алдын ала сараптама",
    W3: "III — RP",
    W4: "IV — ДК-ға тапсыру",
    W5: "V — Қалпына келтіру",
    W6: "VI — ДК қабылдағаннан кейін",
    W7: "VII — Қорғау және одан кейін",
  };
  const map = lang.startsWith('ru') ? ru : lang.startsWith('kz') ? kz : en;
  if (!s) return '—';
  return map[s] || s;
}
